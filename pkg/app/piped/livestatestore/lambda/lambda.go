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

package lambda

import (
	"context"
	"time"

	"go.uber.org/zap"

	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/lambda"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type Store struct {
	store         *store
	logger        *zap.Logger
	interval      time.Duration
	firstSyncedCh chan error
}

type Getter interface {
	GetFunctionManifest(appID string) (provider.FunctionManifest, bool)
	GetState(appID string) (State, bool)

	WaitForReady(ctx context.Context, timeout time.Duration) error
}

type State struct {
	Resources []*model.LambdaResourceState
	Version   model.ApplicationLiveStateVersion
}

func NewStore(cfg *config.PlatformProviderLambdaConfig, platformProvider string, logger *zap.Logger) (*Store, error) {
	logger = logger.Named("lambda").
		With(zap.String("platform-provider", platformProvider))

	client, err := provider.DefaultRegistry().Client(platformProvider, cfg, logger)
	if err != nil {
		return nil, err
	}

	store := &Store{
		store: &store{
			client: client,
			logger: logger.Named("store"),
		},
		interval:      time.Duration(cfg.LiveStateInterval),
		logger:        logger,
		firstSyncedCh: make(chan error, 1),
	}

	return store, nil
}

func (s *Store) Run(ctx context.Context) error {
	s.logger.Info("start running lambda app state store")

	tick := time.NewTicker(s.interval)
	defer tick.Stop()

	// Run the first sync of Lambda resources.
	if err := s.store.run(ctx); err != nil {
		s.firstSyncedCh <- err
		return err
	}

	s.logger.Info("successfully ran the first sync of all lambda resources")
	close(s.firstSyncedCh)

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("lambda app state store has been stopped")
			return nil

		case <-tick.C:
			if err := s.store.run(ctx); err != nil {
				s.logger.Error("failed to sync lambda resources", zap.Error(err))
				continue
			}
			s.logger.Info("successfully synced all lambda resources")
		}
	}
}

func (s *Store) GetFunctionManifest(appID string) (provider.FunctionManifest, bool) {
	return s.store.getFunctionManifest(appID)
}

func (s *Store) GetState(appID string) (State, bool) {
	return s.store.getState(appID)
}

func (s *Store) WaitForReady(ctx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil
	case err := <-s.firstSyncedCh:
		return err
	}
}
