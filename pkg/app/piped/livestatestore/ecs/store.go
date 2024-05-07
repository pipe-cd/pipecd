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
	manifests provider.ECSManifests

	// service, taskset, task??, taskDef

	// The states of service and all its active revsions which may handle the traffic.
	states  []*model.ECSResourceState
	version model.ApplicationLiveStateVersion
}

func (s *store) run(ctx context.Context) error {

	{
		clusters, err := s.fetchClusters(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch ECS clusters: %w", err)
		}

		for _, cluster := range clusters {
			services, err := s.fetchServices(ctx, cluster)
			if err != nil {
				return fmt.Errorf("failed to fetch ECS services of cluster %s: %w", cluster, err)
			}

			m := make(map[*types.Service]map[string][]*types.Task, len(services))
			for _, service := range services {
				m[service] = make(map[string][]*types.Task, len(service.TaskSets))
				for _, taskSet := range service.TaskSets {
					tasks, err := s.fetchTasks(ctx, taskSet)
					if err != nil {
						return fmt.Errorf("failed to fetch ECS tasks of task set %s: %w", *taskSet.TaskSetArn, err)
					}
					m[service][*taskSet.TaskSetArn] = tasks
				}
			}
			// TODO Convert to apps and store to the Store

		}

	}

	return nil
}

func (s *store) fetchClusters(ctx context.Context) ([]string, error) {
	return s.client.ListClusters(ctx)
}

func (s *store) fetchServices(ctx context.Context, cluster string) ([]*types.Service, error) {
	return s.client.GetServices(ctx, cluster)
}

// func (s *store) fetchTaskSets(ctx context.Context, service types.Service) ([]*types.TaskSet, error) {
// 	return s.client.GetServiceTaskSets(ctx, service)
// }

func (s *store) fetchTasks(ctx context.Context, taskSet types.TaskSet) ([]*types.Task, error) {
	return s.client.GetTaskSetTasks(ctx, taskSet)
}

// func (s *store) fetchTaskDefinition(ctx context.Context, taskDefArn string) (*types.TaskDefinition, error) {
// 	return s.client.GetTaskDefinition(ctx, taskDefArn)
// }
