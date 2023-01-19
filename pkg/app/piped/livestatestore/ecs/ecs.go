// Copyright 2022 The PipeCD Authors.
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

	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/ecs"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
	"go.uber.org/zap"
)

type Store struct {
	store         *store
	logger        *zap.Logger
	interval      time.Duration
	firstSyncedCh chan error
}

type Getter interface {
	GetAppLiveManifests(appID string) []provider.Manifest
}

type AppState struct {
	Resources []*model.EcsApplicationLiveState
	Version   model.ApplicationLiveStateVersion
}

func NewStore(
	ctx context.Context,
	cfg *config.PlatformProviderECSConfig,
	platformProvider string,
	pipcdConfig *config.PipedSpec,
	logger *zap.Logger,
) *Store {
	logger = logger.Named("ecs").
		With(zap.String("cloud-provider", platformProvider))

	return &Store{
		store: &store{
			logger: logger.Named("store"),
		},
		interval:      15 * time.Second,
		logger:        logger,
		firstSyncedCh: make(chan error, 1),
	}
}

func (s *Store) Run(ctx context.Context) error {
	return nil
}
