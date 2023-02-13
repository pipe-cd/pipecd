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
	"fmt"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/cloudrun"
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
	svcs, err := s.fetchManagedServices(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch managed services: %w", err)
	}

	revs := make(map[string][]*provider.Revision, len(svcs))
	for _, svc := range svcs {
		id, ok := svc.UID()
		if !ok {
			continue
		}
		names := svc.ActiveRevisionNames()
		rs, err := s.fetchActiveRevisions(ctx, names)
		if err != nil {
			return fmt.Errorf("failed to fetch active revisions: %w", err)
		}
		if len(rs) == 0 {
			continue
		}
		revs[id] = rs
	}

	// Update apps to the latest.
	apps := s.buildAppMap(svcs, revs)
	s.apps.Store(apps)

	return nil
}

func (s *store) buildAppMap(svcs []*provider.Service, revs map[string][]*provider.Revision) map[string]app {
	apps, now := make(map[string]app, len(svcs)), time.Now()
	version := model.ApplicationLiveStateVersion{
		Timestamp: now.Unix(),
	}
	for _, svc := range svcs {
		sm, err := svc.ServiceManifest()
		if err != nil {
			s.logger.Error("failed to load cloudrun service into service manifest", zap.Error(err))
			continue
		}

		appID, ok := sm.AppID()
		if !ok {
			continue
		}

		id, _ := svc.UID()
		apps[appID] = app{
			service: sm,
			states:  provider.MakeResourceStates(svc, revs[id], now),
			version: version,
		}
	}
	return apps
}

func (s *store) fetchManagedServices(ctx context.Context) ([]*provider.Service, error) {
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

func (s *store) fetchActiveRevisions(ctx context.Context, names []string) ([]*provider.Revision, error) {
	ops := &provider.ListRevisionsOptions{
		LabelSelector: provider.MakeRevisionNamesSelector(names),
	}
	v, _, err := s.client.ListRevisions(ctx, ops)
	return v, err
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
