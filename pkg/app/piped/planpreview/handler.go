// Copyright 2024 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package planpreview

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	metrics "github.com/pipe-cd/pipecd/pkg/app/piped/planpreview/planpreviewmetrics"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/regexpool"
)

const (
	defaultWorkerNum              = 3
	defaultCommandQueueBufferSize = 10
	defaultCommandCheckInterval   = 5 * time.Second
	defaultCommandHandleTimeout   = 5 * time.Minute
)

type options struct {
	workerNum              int
	commandQueueBufferSize int
	commandCheckInterval   time.Duration
	commandHandleTimeout   time.Duration
	logger                 *zap.Logger
}

type Option func(*options)

func WithWorkerNum(n int) Option {
	return func(opts *options) {
		opts.workerNum = n
	}
}

func WithCommandQueueBufferSize(s int) Option {
	return func(opts *options) {
		opts.commandQueueBufferSize = s
	}
}

func WithCommandCheckInterval(i time.Duration) Option {
	return func(opts *options) {
		opts.commandCheckInterval = i
	}
}

func WithCommandHandleTimeout(t time.Duration) Option {
	return func(opts *options) {
		opts.commandHandleTimeout = t
	}
}

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

type apiClient interface {
	GetApplicationMostRecentDeployment(ctx context.Context, req *pipedservice.GetApplicationMostRecentDeploymentRequest, opts ...grpc.CallOption) (*pipedservice.GetApplicationMostRecentDeploymentResponse, error)
}

type applicationLister interface {
	List() []*model.Application
}

type commandLister interface {
	ListBuildPlanPreviewCommands() []model.ReportableCommand
}

type secretDecrypter interface {
	Decrypt(string) (string, error)
}

type Handler struct {
	gitClient     gitClient
	commandLister commandLister

	commandCh    chan model.ReportableCommand
	prevCommands map[string]struct{}

	options        *options
	builderFactory func() Builder
	logger         *zap.Logger
}

func NewHandler(
	gc gitClient,
	ac apiClient,
	cl commandLister,
	al applicationLister,
	cg lastTriggeredCommitGetter,
	sd secretDecrypter,
	appManifestsCache cache.Cache,
	cfg *config.PipedSpec,
	opts ...Option,
) *Handler {

	opt := &options{
		workerNum:              defaultWorkerNum,
		commandQueueBufferSize: defaultCommandQueueBufferSize,
		commandCheckInterval:   defaultCommandCheckInterval,
		commandHandleTimeout:   defaultCommandHandleTimeout,
		logger:                 zap.NewNop(),
	}
	for _, o := range opts {
		o(opt)
	}

	h := &Handler{
		gitClient:     gc,
		commandLister: cl,
		commandCh:     make(chan model.ReportableCommand, opt.commandQueueBufferSize),
		prevCommands:  map[string]struct{}{},
		options:       opt,
		logger:        opt.logger.Named("plan-preview-handler"),
	}

	regexPool := regexpool.DefaultPool()
	h.builderFactory = func() Builder {
		return newBuilder(gc, ac, al, cg, sd, appManifestsCache, regexPool, cfg, h.logger)
	}

	return h
}

// Run starts running Handler until the given context has done.
func (h *Handler) Run(ctx context.Context) error {
	h.logger.Info("start running planpreview handler")

	startWorker := func(ctx context.Context, cmdCh <-chan model.ReportableCommand) {
		h.logger.Info("started a worker for handling plan-preview command")
		for {
			select {
			case cmd := <-cmdCh:
				h.handleCommand(ctx, cmd)

			case <-ctx.Done():
				h.logger.Info("a worker has been stopped")
				return
			}
		}
	}

	h.logger.Info(fmt.Sprintf("spawn %d worker to handle commands", h.options.workerNum))
	for i := 0; i < h.options.workerNum; i++ {
		go startWorker(ctx, h.commandCh)
	}

	commandTicker := time.NewTicker(h.options.commandCheckInterval)
	defer commandTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			h.logger.Info("planpreview handler has been stopped")
			return nil

		case <-commandTicker.C:
			h.enqueueNewCommands(ctx)
		}
	}
}

func (h *Handler) enqueueNewCommands(ctx context.Context) {
	h.logger.Debug("fetching unhandled commands to enqueue")

	commands := h.commandLister.ListBuildPlanPreviewCommands()
	if len(commands) == 0 {
		h.logger.Debug("there is no command to enqueue")
		return
	}

	news := make([]model.ReportableCommand, 0, len(commands))
	cmds := make(map[string]struct{}, len(commands))
	for _, cmd := range commands {
		cmds[cmd.Id] = struct{}{}
		if _, ok := h.prevCommands[cmd.Id]; !ok {
			news = append(news, cmd)
		}
	}

	h.logger.Info("fetched unhandled commands to enqueue",
		zap.Any("pre-commands", h.prevCommands),
		zap.Any("commands", cmds),
		zap.Int("news", len(news)),
	)

	if len(news) == 0 {
		h.logger.Info("there is no new command to enqueue")
		return
	}

	h.prevCommands = cmds
	metrics.ReceivedCommands(len(news))
	h.logger.Info(fmt.Sprintf("will enqueue %d new commands", len(news)))

	for _, cmd := range news {
		select {
		case h.commandCh <- cmd:
			h.logger.Info("queued a new new command", zap.String("command", cmd.Id))

		case <-ctx.Done():
			return
		}
	}
}

func (h *Handler) handleCommand(ctx context.Context, cmd model.ReportableCommand) {
	start := time.Now()
	logger := h.logger.With(
		zap.String("command", cmd.Id),
	)
	logger.Info("received a plan-preview command to handle")

	result := &model.PlanPreviewCommandResult{
		CommandId: cmd.Id,
		PipedId:   cmd.PipedId,
	}

	reportError := func(err error) {
		metrics.HandledCommand(metrics.StatusFailure, time.Since(start))

		result.Error = err.Error()
		output, err := json.Marshal(result)
		if err != nil {
			logger.Error("failed to marshal command result", zap.Error(err))
		}

		if err := cmd.Report(ctx, model.CommandStatus_COMMAND_FAILED, nil, output); err != nil {
			logger.Error("failed to report command status", zap.Error(err))
			return
		}
		logger.Info("successfully reported a failure command")
	}

	if cmd.BuildPlanPreview == nil {
		reportError(fmt.Errorf("malformed command"))
		return
	}

	timeout := time.Duration(cmd.BuildPlanPreview.Timeout) * time.Second
	if timeout == 0 {
		timeout = h.options.commandHandleTimeout
	}
	buildCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	b := h.builderFactory()
	appResults, err := b.Build(buildCtx, cmd.Id, *cmd.BuildPlanPreview)
	if err != nil {
		reportError(err)
		return
	}

	result.Results = appResults
	output, err := json.Marshal(result)
	if err != nil {
		reportError(fmt.Errorf("failed to marshal command result (%w)", err))
		return
	}

	if err := cmd.Report(ctx, model.CommandStatus_COMMAND_SUCCEEDED, nil, output); err != nil {
		metrics.HandledCommand(metrics.StatusFailure, time.Since(start))
		logger.Error("failed to report command status", zap.Error(err))
		return
	}

	metrics.HandledCommand(metrics.StatusSuccess, time.Since(start))
	logger.Info("successfully reported a success command")
}
