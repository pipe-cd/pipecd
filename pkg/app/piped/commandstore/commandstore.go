// Copyright 2020 The PipeCD Authors.
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
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/kapetaniosci/pipe/pkg/app/api/service/pipedservice"
	"github.com/kapetaniosci/pipe/pkg/model"
)

type apiClient interface {
	ListUnhandledCommands(ctx context.Context, in *pipedservice.ListUnhandledCommandsRequest, opts ...grpc.CallOption) (*pipedservice.ListUnhandledCommandsResponse, error)
	ReportCommandHandled(ctx context.Context, in *pipedservice.ReportCommandHandledRequest, opts ...grpc.CallOption) (*pipedservice.ReportCommandHandledResponse, error)
}

type Store interface {
	Run(ctx context.Context) error
	ListApplicationCommands() []*model.Command
	ListDeploymentCommands(deploymentID string) []*model.Command
	ReportCommandHandled(ctx context.Context, c *model.Command, status model.CommandStatus, metadata map[string]string) error
}

type store struct {
	apiClient           apiClient
	syncInterval        time.Duration
	applicationCommands []*model.Command
	deploymentCommands  []*model.Command
	handledCommands     map[string]time.Time
	mu                  sync.RWMutex
	gracePeriod         time.Duration
	logger              *zap.Logger
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

func (s *store) sync(ctx context.Context) error {
	resp, err := s.apiClient.ListUnhandledCommands(ctx, &pipedservice.ListUnhandledCommandsRequest{})
	if err != nil {
		s.logger.Error("failed to list unhandled commands", zap.Error(err))
		return err
	}

	applicationCommands := make([]*model.Command, 0, len(resp.Commands))
	deploymentCommands := make([]*model.Command, 0, len(resp.Commands))
	for _, cmd := range resp.Commands {
		switch cmd.Type {
		case model.CommandType_COMMAND_APPLICATION:
			applicationCommands = append(applicationCommands, cmd)
		case model.CommandType_COMMAND_DEPLOYMENT:
			deploymentCommands = append(deploymentCommands, cmd)
		}
	}

	s.mu.Lock()
	s.applicationCommands = applicationCommands
	s.deploymentCommands = deploymentCommands
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

func (s *store) ListApplicationCommands() []*model.Command {
	s.mu.RLock()
	defer s.mu.RUnlock()

	commands := make([]*model.Command, 0, len(s.applicationCommands))
	for _, cmd := range s.applicationCommands {
		if _, ok := s.handledCommands[cmd.Id]; ok {
			continue
		}
		commands = append(commands, cmd)
	}
	return commands
}

func (s *store) ListDeploymentCommands(deploymentID string) []*model.Command {
	s.mu.RLock()
	defer s.mu.RUnlock()

	commands := make([]*model.Command, 0, len(s.deploymentCommands))
	for _, cmd := range s.deploymentCommands {
		if _, ok := s.handledCommands[cmd.Id]; ok {
			continue
		}
		if cmd.DeploymentId != deploymentID {
			continue
		}
		commands = append(commands, cmd)
	}
	return commands
}

func (s *store) ReportCommandHandled(ctx context.Context, c *model.Command, status model.CommandStatus, metadata map[string]string) error {
	now := time.Now()

	s.mu.Lock()
	s.handledCommands[c.Id] = now
	s.mu.Unlock()

	_, err := s.apiClient.ReportCommandHandled(ctx, &pipedservice.ReportCommandHandledRequest{
		CommandId: c.Id,
		Status:    status,
		Metadata:  metadata,
		HandledAt: now.Unix(),
	})
	return err
}
