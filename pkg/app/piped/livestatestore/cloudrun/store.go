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
	"sync"

	provider "github.com/pipe-cd/pipecd/pkg/app/piped/cloudprovider/cloudrun"
	"github.com/pipe-cd/pipecd/pkg/config"
	"go.uber.org/zap"
)

type store struct {
	config        *config.CloudProviderCloudRunConfig
	cloudProvider string
	apps          map[string]provider.ServiceManifest
	mu            sync.RWMutex
	logger        *zap.Logger
	client        provider.Client
}

func (s *store) addApp(sm provider.ServiceManifest) {
	appID := sm.Labels()[provider.LabelApplication]
	if appID == "" {
		return
	}
	s.mu.Lock()
	s.apps[appID] = sm
	s.mu.Unlock()
}

func (s *store) run(ctx context.Context) error {
	client, err := provider.DefaultRegistry().Client(ctx, s.cloudProvider, s.config, s.logger)
	if err != nil {
		s.logger.Error("failed to new cloudrun client: %v", zap.Error(err))
		return err
	}

	const maxLimit = 500
	var cursor string
	svc := make([]*provider.Service, 0, maxLimit)
	for {
		ops := &provider.ListOptions{
			Limit:         maxLimit,
			LabelSelector: provider.MakeManagedByPipedLabel(),
			Cursor:        cursor,
		}
		// Cloud Run Admin API rate Limits
		// https://cloud.google.com/run/quotas#api
		v, next, err := client.List(ctx, ops)
		if err != nil {
			s.logger.Error("failed to list cloudrun services: %v", zap.Error(err))
			return err
		}
		svc = append(svc, v...)
		if next == "" {
			break
		}
		cursor = next
	}
	for i := range svc {
		sm, err := svc[i].ServiceManifest()
		if err != nil {
			s.logger.Error("failed to load cloudrun service into service manifest: %v", zap.Error(err))
			continue
		}
		s.addApp(sm)
	}
	return nil
}

func (s *store) GetAppLiveServiceManifest(appID string) provider.ServiceManifest {
	s.mu.RLock()
	sv, ok := s.apps[appID]
	s.mu.RUnlock()

	if !ok {
		return provider.ServiceManifest{}
	}

	return sv
}
