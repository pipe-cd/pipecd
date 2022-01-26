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

	"go.uber.org/zap"

	provider "github.com/pipe-cd/pipecd/pkg/app/piped/cloudprovider/cloudrun"
)

type store struct {
	apps   atomic.Value
	logger *zap.Logger
	client provider.Client
}

func (s *store) run(ctx context.Context) error {
	const maxLimit = 500
	var cursor string
	svc := make([]*provider.Service, 0, maxLimit)
	for {
		ops := &provider.ListOptions{
			Limit:         maxLimit,
			LabelSelector: provider.MakeManagedByPipedLabel(),
			Cursor:        cursor,
		}
		// Cloud Run Admin API rate Limits.
		// https://cloud.google.com/run/quotas#api
		v, next, err := s.client.List(ctx, ops)
		if err != nil {
			return fmt.Errorf("failed to list cloudrun services: %w", err)
		}
		svc = append(svc, v...)
		if next == "" {
			break
		}
		cursor = next
	}

	// Update apps to the latest.
	s.setApps(svc)

	return nil
}

func (s *store) setApps(svc []*provider.Service) {
	apps := make(map[string]provider.ServiceManifest, len(svc))
	for i := range svc {
		sm, err := svc[i].ServiceManifest()
		if err != nil {
			s.logger.Error("failed to load cloudrun service into service manifest", zap.Error(err))
			continue
		}
		appID := sm.Labels()[provider.LabelApplication]
		apps[appID] = sm
	}
	s.apps.Store(apps)
}

func (s *store) loadApps() map[string]provider.ServiceManifest {
	apps := s.apps.Load()
	if apps == nil {
		return nil
	}
	return apps.(map[string]provider.ServiceManifest)
}

func (s *store) GetServiceManifest(appID string) (provider.ServiceManifest, bool) {
	apps := s.loadApps()
	if apps == nil {
		s.logger.Error("failed to load cloudrun apps")
		return provider.ServiceManifest{}, false
	}

	sm, ok := apps[appID]
	if !ok {
		return provider.ServiceManifest{}, false
	}

	return sm, true
}
