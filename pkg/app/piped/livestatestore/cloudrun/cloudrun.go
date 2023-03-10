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

package cloudrun

import (
	"context"
	"time"

	"go.uber.org/zap"

	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/cloudrun"
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
	GetState(appID string) (State, bool)
	GetServiceManifest(appID string) (provider.ServiceManifest, bool)

	WaitForReady(ctx context.Context, timeout time.Duration) error
}

type State struct {
	Resources []*model.CloudRunResourceState
	Version   model.ApplicationLiveStateVersion
}

func NewStore(ctx context.Context, cfg *config.PlatformProviderCloudRunConfig, platformProvider string, logger *zap.Logger) (*Store, error) {
	logger = logger.Named("cloudrun").
		With(zap.String("platform-provider", platformProvider))

	client, err := provider.DefaultRegistry().Client(ctx, platformProvider, cfg, logger)
	if err != nil {
		return nil, err
	}

	store := &Store{
		store: &store{
			client: client,
			logger: logger.Named("store"),
		},
		interval:      15 * time.Second,
		logger:        logger,
		firstSyncedCh: make(chan error, 1),
	}

	return store, nil
}

func (s *Store) Run(ctx context.Context) error {
	s.logger.Info("start running cloudrun app state store")

	tick := time.NewTicker(s.interval)
	defer tick.Stop()

	// Run the first sync cloudrun services.
	if err := s.store.run(ctx); err != nil {
		s.firstSyncedCh <- err
		return err
	}

	s.logger.Info("successfully the first synced all cloudrun services")
	close(s.firstSyncedCh)

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("cloudrun app state store has been stopped")
			return nil

		case <-tick.C:
			if err := s.store.run(ctx); err != nil {
				s.logger.Error("failed to sync cloudrun services", zap.Error(err))
				continue
			}
			s.logger.Info("successfully synced all cloudrun services")
		}
	}
}

func (s *Store) GetServiceManifest(appID string) (provider.ServiceManifest, bool) {
	return s.store.getServiceManifest(appID)
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
