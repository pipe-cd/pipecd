// Copyright 2023 The PipeCD Authors.
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
	"errors"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type Store interface {
	ListUnhandledCommands(ctx context.Context, pipedID string) ([]*model.Command, error)
	AddCommand(ctx context.Context, command *model.Command) error
	GetCommand(ctx context.Context, id string) (*model.Command, error)
	UpdateCommandHandled(ctx context.Context, id string, status model.CommandStatus, metadata map[string]string, unhandledAt int64) error
}

type store struct {
	backend datastore.CommandStore
	cache   *commandCache
	logger  *zap.Logger
}

func NewStore(w datastore.Commander, ds datastore.DataStore, c cache.Cache, logger *zap.Logger) Store {
	return &store{
		backend: datastore.NewCommandStore(ds, w),
		cache: &commandCache{
			backend: c,
		},
		logger: logger,
	}
}

func (s *store) ListUnhandledCommands(ctx context.Context, pipedID string) ([]*model.Command, error) {
	opts := datastore.ListOptions{
		Filters: []datastore.ListFilter{
			{
				Field:    "PipedId",
				Operator: datastore.OperatorEqual,
				Value:    pipedID,
			},
			{
				Field:    "Status",
				Operator: datastore.OperatorEqual,
				Value:    model.CommandStatus_COMMAND_NOT_HANDLED_YET,
			},
		},
	}
	commands, err := s.backend.List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return commands, nil
}

func (s *store) AddCommand(ctx context.Context, command *model.Command) error {
	if err := s.backend.Add(ctx, command); err != nil {
		s.logger.Error("failed to put command to datastore", zap.Error(err))
		return err
	}

	if err := s.cache.Put(command.Id, command); err != nil {
		s.logger.Error("failed to put command to cache", zap.Error(err))
	}
	return nil
}

func (s *store) GetCommand(ctx context.Context, id string) (*model.Command, error) {
	cacheResp, err := s.cache.Get(id)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		s.logger.Error("failed to get command from cache", zap.Error(err))
	}
	if cacheResp != nil {
		return cacheResp, nil
	}

	dsResp, err := s.backend.Get(ctx, id)
	if err != nil {
		s.logger.Error("failed to get command from datastore", zap.Error(err))
		return nil, err
	}

	if err := s.cache.Put(id, dsResp); err != nil {
		s.logger.Error("failed to put command to cache", zap.Error(err))
	}
	return dsResp, nil
}

func (s *store) UpdateCommandHandled(ctx context.Context, id string, status model.CommandStatus, metadata map[string]string, handledAt int64) error {
	if err := s.backend.UpdateStatus(ctx, id, status, metadata, handledAt); err != nil {
		return err
	}

	cmd, err := s.backend.Get(ctx, id)
	if err != nil {
		s.logger.Error("failed to get command from datastore", zap.Error(err))
	}

	if err := s.cache.Put(id, cmd); err != nil {
		s.logger.Error("failed to put command to cache", zap.Error(err))
	}
	return nil
}
