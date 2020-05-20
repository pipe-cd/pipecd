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

package applicationstore

import (
	"context"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/kapetaniosci/pipe/pkg/app/api/service/pipedservice"
	"github.com/kapetaniosci/pipe/pkg/model"
)

type apiClient interface {
	ListApplications(ctx context.Context, in *pipedservice.ListApplicationsRequest, opts ...grpc.CallOption) (*pipedservice.ListApplicationsResponse, error)
}

type Store interface {
	Run(ctx context.Context) error
	ListApplications() []*model.Application
	GetApplication(id string) (*model.Application, bool)
}

type store struct {
	apiClient       apiClient
	applicationMap  atomic.Value
	applicationList atomic.Value
	syncInterval    time.Duration
	gracePeriod     time.Duration
	logger          *zap.Logger
}

var (
	defaultSyncInterval = time.Minute
)

// NewStore creates a new application store instance.
// This syncs with the control plane to keep the list of applications for this runner up-to-date.
func NewStore(apiClient apiClient, gracePeriod time.Duration, logger *zap.Logger) Store {
	return &store{
		apiClient:    apiClient,
		syncInterval: defaultSyncInterval,
		gracePeriod:  gracePeriod,
		logger:       logger.Named("application-store"),
	}
}

// Run starts syncing.
func (s *store) Run(ctx context.Context) error {
	s.logger.Info("start running application store")

	syncTicker := time.NewTicker(s.syncInterval)
	defer syncTicker.Stop()

	for {
		select {
		case <-syncTicker.C:
			s.sync(ctx)

		case <-ctx.Done():
			s.logger.Info("application store has been stopped")
			return nil
		}
	}
}

func (s *store) sync(ctx context.Context) error {
	resp, err := s.apiClient.ListApplications(ctx, &pipedservice.ListApplicationsRequest{})
	if err != nil {
		s.logger.Error("failed to list unhandled application", zap.Error(err))
		return err
	}

	applicationMap := make(map[string]*model.Application, len(resp.Applications))
	for _, app := range resp.Applications {
		applicationMap[app.Id] = app
	}

	s.applicationMap.Store(applicationMap)
	s.applicationList.Store(resp.Applications)
	return nil
}

func (s *store) ListApplications() []*model.Application {
	apps := s.applicationList.Load()
	if apps == nil {
		return nil
	}
	return apps.([]*model.Application)
}

func (s *store) GetApplication(id string) (*model.Application, bool) {
	apps := s.applicationList.Load()
	if apps == nil {
		return nil, false
	}

	app, ok := apps.(map[string]*model.Application)[id]
	return app, ok
}
