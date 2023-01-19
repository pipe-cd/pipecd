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
	"sync"
	"sync/atomic"

	"go.uber.org/zap"

	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/ecs"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type applicationLister interface {
	List() []*model.Application
}

type store struct {
	apps   atomic.Value
	logger *zap.Logger
	client provider.Client
	mu     sync.RWMutex
}

type app struct {
	taskDefinisionManifest provider.Manifest
	serviceManifest        provider.Manifest
	states                 []*model.EcsApplicationLiveState
	version                model.ApplicationLiveStateVersion
}

func (s *store) Run(ctx context.Context) error {
	s.logger.Info("start running ecs app state store")

	s.logger.Info("ecs app state store has been stopped")
	return nil
}

func (s *store) loadApps() map[string]app {
	apps := s.apps.Load()
	if apps == nil {
		return nil
	}
	return apps.(map[string]app)
}
