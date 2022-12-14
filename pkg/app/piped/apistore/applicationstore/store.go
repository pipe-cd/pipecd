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

package applicationstore

import (
	"context"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/model"
)

// Lister helps list and get application.
// All objects returned here must be treated as read-only.
type Lister interface {
	// List lists all applications that should be handled by this piped.
	// All disabled applications will be ignored.
	List() []*model.Application
	// ListByPlatformProvider lists all applications for a given cloud provider name.
	ListByPlatformProvider(name string) []*model.Application
	// Get retrieves a specifiec deployment for the given id.
	Get(id string) (*model.Application, bool)
}

type apiClient interface {
	ListApplications(ctx context.Context, in *pipedservice.ListApplicationsRequest, opts ...grpc.CallOption) (*pipedservice.ListApplicationsResponse, error)
}

type Store interface {
	// Run starts syncing the application list with the control-plane.
	Run(ctx context.Context) error
	// Lister returns a lister for retrieving applications.
	Lister() Lister
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

// Run starts syncing the application list with the control-plane.
func (s *store) Run(ctx context.Context) error {
	s.logger.Info("start running application store")

	syncTicker := time.NewTicker(s.syncInterval)
	defer syncTicker.Stop()

	// Do first sync without waiting the first ticker.
	s.sync(ctx)

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

// Lister returns a lister for retrieving applications.
func (s *store) Lister() Lister {
	return s
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

// List lists all applications that should be handled by this piped.
// All disabled applications will be ignored.
func (s *store) List() []*model.Application {
	apps := s.applicationList.Load()
	if apps == nil {
		return nil
	}
	return apps.([]*model.Application)
}

// ListByPlatformProvider lists all applications for a given platform provider name.
func (s *store) ListByPlatformProvider(name string) []*model.Application {
	list := s.applicationList.Load()
	if list == nil {
		return nil
	}

	var (
		apps = list.([]*model.Application)
		out  = make([]*model.Application, 0, len(apps))
	)
	for _, app := range apps {
		if app.PlatformProvider == name {
			out = append(out, app)
		}
	}
	return out
}

// Get retrieves a specific deployment for the given id.
func (s *store) Get(id string) (*model.Application, bool) {
	apps := s.applicationMap.Load()
	if apps == nil {
		return nil, false
	}

	app, ok := apps.(map[string]*model.Application)[id]
	return app, ok
}
