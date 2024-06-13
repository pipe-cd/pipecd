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

package ecs

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/pipe-cd/pipecd/pkg/model"
	"go.uber.org/zap"

	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/ecs"
)

type store struct {
	apps   atomic.Value
	logger *zap.Logger
	client provider.Client
}

type app struct {
	// ServiceDefinition and its primary taskset's TaskDefinition.
	manifests provider.ECSManifests

	// States of services, tasksets, and tasks.
	// NOTE: Standalone tasks are NOT included yet.
	states  []*model.ECSResourceState
	version model.ApplicationLiveStateVersion
}

func (s *store) run(ctx context.Context) error {
	apps := map[string]app{}
	now := time.Now()
	version := model.ApplicationLiveStateVersion{
		Timestamp: now.Unix(),
	}

	clusters, err := s.client.ListClusters(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch ECS clusters: %w", err)
	}

	for _, cluster := range clusters {
		services, err := s.client.GetServices(ctx, cluster)
		if err != nil {
			return fmt.Errorf("failed to fetch ECS services of cluster %s: %w", cluster, err)
		}

		for _, service := range services {
			taskSetTasks := make(map[string][]*types.Task, len(service.TaskSets))
			var primaryTaskDef *types.TaskDefinition
			for _, taskSet := range service.TaskSets {
				if *taskSet.Status == "PRIMARY" {
					primaryTaskDef, err = s.client.GetTaskDefinition(ctx, *taskSet.TaskDefinition)
					if err != nil {
						return fmt.Errorf("failed to fetch ECS task definition %s: %w", *taskSet.TaskDefinition, err)
					}
				}

				tasks, err := s.client.GetTaskSetTasks(ctx, taskSet)
				if err != nil {
					return fmt.Errorf("failed to fetch ECS tasks of task set %s: %w", *taskSet.TaskSetArn, err)
				}
				taskSetTasks[*taskSet.TaskSetArn] = tasks
			}

			apps[*service.ServiceArn] = app{
				manifests: provider.ECSManifests{
					ServiceDefinition: service,
					TaskDefinition:    primaryTaskDef,
				},
				states:  provider.MakeServiceResourceStates(service, taskSetTasks),
				version: version,
			}
		}
	}

	// 3. Store the apps
	s.apps.Store(apps)

	return nil
}

func (s *store) loadApps() map[string]app {
	apps := s.apps.Load()
	if apps == nil {
		return nil
	}
	return apps.(map[string]app)
}

func (s *store) getManifests(appID string) (provider.ECSManifests, bool) {
	apps := s.loadApps()
	if apps == nil {
		return provider.ECSManifests{}, false
	}

	app, ok := apps[appID]
	if !ok {
		return provider.ECSManifests{}, false
	}

	return app.manifests, true
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

	return State{
		Resources: app.states,
		Version:   app.version,
	}, true
}
