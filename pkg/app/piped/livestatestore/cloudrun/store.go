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

package cloudrun

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	provider "github.com/pipe-cd/pipecd/pkg/app/piped/cloudprovider/cloudrun"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type store struct {
	apps   atomic.Value
	logger *zap.Logger
	client provider.Client
}

type app struct {
	service provider.ServiceManifest
	// The states of service and all its active revsions which may handle the traffic.
	states  []*model.CloudRunResourceState
	version model.ApplicationLiveStateVersion
}

func (s *store) run(ctx context.Context) error {
	svcs, err := s.getAllServices(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all cloudrun services: %w", err)
	}
	revs := s.listRevisionsFromServices(ctx, svcs)

	// Update apps to the latest.
	s.storeApps(svcs, revs)

	return nil
}

func (s *store) storeApps(svcs []*provider.Service, revs map[string][]*provider.Revision) {
	apps := make(map[string]app, len(svcs))
	for i := range svcs {
		sm, err := svcs[i].ServiceManifest()
		if err != nil {
			s.logger.Error("failed to load cloudrun service into service manifest", zap.Error(err))
			continue
		}

		appID, ok := sm.AppID()
		if !ok {
			continue
		}

		id, _ := svcs[i].UID()
		now := time.Now()
		apps[appID] = app{
			service: sm,
			states:  provider.MakeResourceStates(svcs[i], revs[id], now),
			version: model.ApplicationLiveStateVersion{
				Timestamp: now.Unix(),
			},
		}
	}
	s.apps.Store(apps)
}

func (s *store) getAllServices(ctx context.Context) ([]*provider.Service, error) {
	const maxLimit = 500
	var cursor string
	svcs := make([]*provider.Service, 0, maxLimit)
	for {
		ops := &provider.ListOptions{
			Limit:         maxLimit,
			LabelSelector: provider.MakeManagedByPipedSelector(),
			Cursor:        cursor,
		}
		// Cloud Run Admin API rate Limits.
		// https://cloud.google.com/run/quotas#api
		v, next, err := s.client.List(ctx, ops)
		if err != nil {
			return nil, err
		}
		svcs = append(svcs, v...)
		if next == "" {
			break
		}
		cursor = next
	}
	return svcs, nil
}

func (s *store) listRevisionsFromServices(ctx context.Context, svcs []*provider.Service) map[string][]*provider.Revision {
	revs := make(map[string][]*provider.Revision, len(svcs))
	for i := range svcs {
		id, ok := svcs[i].UID()
		if !ok {
			continue
		}
		names := svcs[i].ActiveRevisionNames()
		if rs := s.getMultiRevisions(ctx, names); len(rs) > 0 {
			revs[id] = rs
		}
	}
	return revs
}

func (s *store) getMultiRevisions(ctx context.Context, names []string) []*provider.Revision {
	ops := &provider.ListRevisionsOptions{
		LabelSelector: provider.MakeRevisionNamesSelector(names),
	}
	v, _, err := s.client.ListRevisions(ctx, ops)
	if err != nil {
		s.logger.Error("failed to list cloudrun revisions", zap.Error(err))
	}
	return v
}

func (s *store) loadApps() map[string]app {
	apps := s.apps.Load()
	if apps == nil {
		return nil
	}
	return apps.(map[string]app)
}

func (s *store) getServiceManifest(appID string) (provider.ServiceManifest, bool) {
	apps := s.loadApps()
	if apps == nil {
		return provider.ServiceManifest{}, false
	}

	app, ok := apps[appID]
	if !ok {
		return provider.ServiceManifest{}, false
	}

	return app.service, true
}

func (s *store) getState(appID string) (State, bool) {
	apps := s.loadApps()
	if apps == nil {
		return State{}, false
	}

	app, ok := apps[appID]
	if !ok {
		return State{}, false
	}

	state := State{
		Resources: app.states,
		Version:   app.version,
	}
	return state, true
}
