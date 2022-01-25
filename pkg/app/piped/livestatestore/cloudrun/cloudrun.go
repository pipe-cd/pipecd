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

package cloudrun

import (
	"context"
	"time"

	"go.uber.org/zap"

	provider "github.com/pipe-cd/pipecd/pkg/app/piped/cloudprovider/cloudrun"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type applicationLister interface {
	List() []*model.Application
}

type Store struct {
	store    *store
	logger   *zap.Logger
	interval time.Duration
}

type Getter interface {
	GetAppLiveServiceManifest(appID string) provider.ServiceManifest
}

func NewStore(cfg *config.CloudProviderCloudRunConfig, cloudProvider string, appLister applicationLister, logger *zap.Logger) *Store {
	logger = logger.Named("cloudrun").
		With(zap.String("cloud-provider", cloudProvider))

	return &Store{
		store: &store{
			apps:   make(map[string]provider.ServiceManifest),
			logger: logger.Named("store"),
		},
		interval: 15 * time.Second,
		logger:   logger,
	}
}

func (s *Store) Run(ctx context.Context) error {
	s.logger.Info("start running cloudrun app state store")

	tick := time.NewTicker(s.interval)
	defer tick.Stop()

L:
	for {
		select {
		case <-ctx.Done():
			break L

		case <-tick.C:
			if err := s.store.run(ctx); err != nil {
				continue
			}
			s.logger.Info("successfully synced all cloudrun services")
		}
	}

	s.logger.Info("cloudrun app state store has been stopped")
	return nil
}

func (s *Store) GetAppLiveServiceManifest(appID string) provider.ServiceManifest {
	return s.store.GetAppLiveServiceManifest(appID)
}
