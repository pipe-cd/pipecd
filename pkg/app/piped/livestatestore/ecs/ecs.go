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

package ecs

import (
	"context"
	"time"

	"go.uber.org/zap"

	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/ecs"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type Store struct {
	logger *zap.Logger
}

type Getter interface {
	GetManifests(appID string) (provider.ECSManifests, bool)
	GetState(appID string) (State, bool)

	WaitForReady(ctx context.Context, timeout time.Duration) error
}

type State struct {
	Resources []*model.CloudRunResourceState
	Version   model.ApplicationLiveStateVersion
}

func NewStore(cfg *config.PlatformProviderECSConfig, platformProvider string, logger *zap.Logger) *Store {
	logger = logger.Named("ecs").
		With(zap.String("platform-provider", platformProvider))

	return &Store{
		logger: logger,
	}
}

func (s *Store) Run(ctx context.Context) error {
	s.logger.Info("start running ecs app state store")

	s.logger.Info("ecs app state store has been stopped")
	return nil
}

func (s *Store) GetManifests(appID string) (provider.ECSManifests, bool) {
	panic("unimplemented")
}

func (s *Store) GetState(appID string) (State, bool) {
	panic("unimplemented")
}

func (s *Store) WaitForReady(ctx context.Context, timeout time.Duration) error {
	panic("unimplemented")
}
