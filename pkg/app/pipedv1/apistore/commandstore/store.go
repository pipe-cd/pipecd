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

package commandstore

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type apiClient interface {
	ListUnhandledCommands(ctx context.Context, in *pipedservice.ListUnhandledCommandsRequest, opts ...grpc.CallOption) (*pipedservice.ListUnhandledCommandsResponse, error)
	ReportCommandHandled(ctx context.Context, in *pipedservice.ReportCommandHandledRequest, opts ...grpc.CallOption) (*pipedservice.ReportCommandHandledResponse, error)
}

type Store interface {
	Run(ctx context.Context) error
	Lister() Lister
	StageCommandHandledReporter() StageCommandHandledReporter
}

// Lister helps list commands.
// All objects returned here must be treated as read-only.
type Lister interface {
	ListApplicationCommands() []model.ReportableCommand
	ListDeploymentCommands() []model.ReportableCommand
	ListBuildPlanPreviewCommands() []model.ReportableCommand
	ListPipedCommands() []model.ReportableCommand

	// ListStageCommands returns all stage commands of the given deployment and stage.
	// If the command type is not supported, it returns an error.
	ListStageCommands(deploymentID, stageID string, commandType model.Command_Type) ([]*model.Command, error)
}

type StageCommandHandledReporter interface {
	// ReportCommandsHandled reports all stage commands of the given deployment and stage as handled successfully.
	// This should be called on piped side, not plugin side, to ensure reporting is called correctly.
	ReportCommandsHandled(ctx context.Context, deploymentID, stageID string) error
}

// stageCommandMap is a map of stage commands. Keys are deploymentID and stageID.
type stageCommandMap map[string]map[string][]*model.Command

type store struct {
	apiClient    apiClient
	syncInterval time.Duration
	// TODO: Using atomic for storing a map of all commands
	// instead of some separate lists + mutex as the current.
	applicationCommands  []model.ReportableCommand
	deploymentCommands   []model.ReportableCommand
	planPreviewCommands  []model.ReportableCommand
	pipedCommands        []model.ReportableCommand
	stageApproveCommands stageCommandMap
	stageSkipCommands    stageCommandMap
	handledCommands      map[string]time.Time
	mu                   sync.RWMutex
	gracePeriod          time.Duration
	logger               *zap.Logger
}

var (
	defaultSyncInterval = 5 * time.Second
	staleCommandPeriod  = 10 * time.Minute
)

// NewStore creates a new command store instance.
// This watches/fetches new commands from the control plane
// and then notifies them to the registered subscribers.
func NewStore(apiClient apiClient, gracePeriod time.Duration, logger *zap.Logger) Store {
	return &store{
		apiClient:       apiClient,
		syncInterval:    defaultSyncInterval,
		handledCommands: make(map[string]time.Time),
		gracePeriod:     gracePeriod,
		logger:          logger.Named("command-store"),
	}
}

// Run starts watching and notifying the new commands.
func (s *store) Run(ctx context.Context) error {
	s.logger.Info("start running command store")

	syncTicker := time.NewTicker(s.syncInterval)
	defer syncTicker.Stop()

	cleanHandledCommandTicker := time.NewTicker(10 * time.Minute)
	defer cleanHandledCommandTicker.Stop()

	for {
		select {
		case <-syncTicker.C:
			s.sync(ctx)

		case now := <-cleanHandledCommandTicker.C:
			s.cleanHandledCommands(now)

		case <-ctx.Done():
			s.logger.Info("command store has been stopped")
			return nil
		}
	}
}

func (s *store) Lister() Lister {
	return s
}

func (s *store) StageCommandHandledReporter() StageCommandHandledReporter {
	return s
}

func (s *store) sync(ctx context.Context) error {
	resp, err := s.apiClient.ListUnhandledCommands(ctx, &pipedservice.ListUnhandledCommandsRequest{})
	if err != nil {
		s.logger.Error("failed to list unhandled commands", zap.Error(err))
		return err
	}

	var (
		applicationCommands  = make([]model.ReportableCommand, 0)
		deploymentCommands   = make([]model.ReportableCommand, 0)
		planPreviewCommands  = make([]model.ReportableCommand, 0)
		pipedCommands        = make([]model.ReportableCommand, 0)
		stageApproveCommands stageCommandMap
		stageSkipCommands    stageCommandMap
	)
	for _, cmd := range resp.Commands {
		switch cmd.Type {
		case model.Command_SYNC_APPLICATION, model.Command_UPDATE_APPLICATION_CONFIG, model.Command_CHAIN_SYNC_APPLICATION:
			applicationCommands = append(applicationCommands, s.makeReportableCommand(cmd))
		case model.Command_CANCEL_DEPLOYMENT:
			deploymentCommands = append(deploymentCommands, s.makeReportableCommand(cmd))
		case model.Command_BUILD_PLAN_PREVIEW:
			planPreviewCommands = append(planPreviewCommands, s.makeReportableCommand(cmd))
		case model.Command_RESTART_PIPED:
			pipedCommands = append(pipedCommands, s.makeReportableCommand(cmd))
		case model.Command_APPROVE_STAGE:
			stageApproveCommands.append(cmd)
		case model.Command_SKIP_STAGE:
			stageSkipCommands.append(cmd)
		}
	}

	s.mu.Lock()
	s.applicationCommands = applicationCommands
	s.deploymentCommands = deploymentCommands
	s.planPreviewCommands = planPreviewCommands
	s.pipedCommands = pipedCommands
	s.stageApproveCommands = stageApproveCommands
	s.stageSkipCommands = stageSkipCommands
	s.mu.Unlock()

	return nil
}

func (s *store) cleanHandledCommands(now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	handledCommands := make(map[string]time.Time, len(s.handledCommands))
	for k, v := range s.handledCommands {
		if now.Sub(v) > staleCommandPeriod {
			continue
		}
		handledCommands[k] = v
	}
	s.handledCommands = handledCommands
}

func (s *store) ListApplicationCommands() []model.ReportableCommand {
	s.mu.RLock()
	defer s.mu.RUnlock()

	commands := make([]model.ReportableCommand, 0, len(s.applicationCommands))
	for _, cmd := range s.applicationCommands {
		if _, ok := s.handledCommands[cmd.Id]; ok {
			continue
		}
		commands = append(commands, cmd)
	}
	return commands
}

func (s *store) ListDeploymentCommands() []model.ReportableCommand {
	s.mu.RLock()
	defer s.mu.RUnlock()

	commands := make([]model.ReportableCommand, 0, len(s.deploymentCommands))
	for _, cmd := range s.deploymentCommands {
		if _, ok := s.handledCommands[cmd.Id]; ok {
			continue
		}
		commands = append(commands, cmd)
	}
	return commands
}

func (s *store) ListBuildPlanPreviewCommands() []model.ReportableCommand {
	s.mu.RLock()
	defer s.mu.RUnlock()

	commands := make([]model.ReportableCommand, 0, len(s.planPreviewCommands))
	for _, cmd := range s.planPreviewCommands {
		if _, ok := s.handledCommands[cmd.Id]; ok {
			continue
		}
		commands = append(commands, cmd)
	}
	return commands
}

func (s *store) ListPipedCommands() []model.ReportableCommand {
	s.mu.RLock()
	defer s.mu.RUnlock()

	commands := make([]model.ReportableCommand, 0, len(s.pipedCommands))
	for _, cmd := range s.pipedCommands {
		if _, ok := s.handledCommands[cmd.Id]; ok {
			continue
		}
		commands = append(commands, cmd)
	}
	return commands
}

func (s *store) makeReportableCommand(c *model.Command) model.ReportableCommand {
	return model.ReportableCommand{
		Command: c,
		Report: func(ctx context.Context, status model.CommandStatus, metadata map[string]string, output []byte) error {
			return s.reportCommandHandled(ctx, c, status, metadata, output)
		},
	}
}

func (s *store) reportCommandHandled(ctx context.Context, c *model.Command, status model.CommandStatus, metadata map[string]string, output []byte) error {
	now := time.Now()

	s.mu.Lock()
	s.handledCommands[c.Id] = now
	s.mu.Unlock()

	_, err := s.apiClient.ReportCommandHandled(ctx, &pipedservice.ReportCommandHandledRequest{
		CommandId: c.Id,
		Status:    status,
		Metadata:  metadata,
		HandledAt: now.Unix(),
		Output:    output,
	})
	return err
}

func (s *store) ListStageCommands(deploymentID, stageID string, commandType model.Command_Type) ([]*model.Command, error) {
	var list stageCommandMap
	switch commandType {
	case model.Command_APPROVE_STAGE:
		list = s.stageApproveCommands
	case model.Command_SKIP_STAGE:
		list = s.stageSkipCommands
	default:
		s.logger.Error("invalid command type", zap.String("commandType", commandType.String()))
		return nil, fmt.Errorf("invalid command type: %v", commandType.String())
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return list[deploymentID][stageID], nil
}

func (s *store) ReportCommandsHandled(ctx context.Context, deploymentID, stageID string) error {
	maps := []stageCommandMap{s.stageApproveCommands, s.stageSkipCommands}
	for _, m := range maps {
		for _, c := range m[deploymentID][stageID] {
			if err := s.reportCommandHandled(ctx, c, model.CommandStatus_COMMAND_SUCCEEDED, nil, nil); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m stageCommandMap) append(c *model.Command) {
	deploymentID := c.DeploymentId
	stageID := c.StageId
	if _, ok := m[deploymentID]; !ok {
		m[deploymentID] = make(map[string][]*model.Command)
	}
	m[deploymentID][stageID] = append(m[deploymentID][stageID], c)
}
