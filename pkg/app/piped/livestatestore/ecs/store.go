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
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
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
	serviceDefinition types.Service
	taskDefinision    types.TaskDefinition
	// states            []*model.EcsApplicationLiveState
	version model.ApplicationLiveStateVersion
}

func (s *store) run(ctx context.Context) error {
	clusterArns, err := s.fetchClusterArns(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch clusterArns: %w", err)
	}

	const maxResults = 5000
	allServices := make([]types.Service, 0, maxResults)
	allTasks := make([]types.Task, 0, maxResults)
	for _, clusterArn := range clusterArns {
		services, err := s.fetchServices(ctx, clusterArn)
		if err != nil {
			s.logger.Error("failed to fetch service. clusterArn", zap.Error(err))
			continue
		}
		allServices = append(allServices, services...)

		tasks, err := s.fetchTasks(ctx, clusterArn)
		if err != nil {
			s.logger.Error("failed to fetch task. clusterArn", zap.Error(err))
			continue
		}
		allTasks = append(allTasks, tasks...)
	}

	apps := s.buildAppMap(ctx, allServices, allTasks)
	s.apps.Store(apps)

	return nil
}

func (s *store) buildAppMap(ctx context.Context, services []types.Service, tasks []types.Task) map[string]app {
	maxLen := int32(math.Max(float64(len(services)), float64(len(tasks))))
	apps, now := make(map[string]app, maxLen), time.Now()
	version := model.ApplicationLiveStateVersion{
		Timestamp: now.Unix(),
	}

	for _, service := range services {
		for _, tag := range service.Tags {
			if appID := *tag.Key; *tag.Key == provider.LabelApplication {
				if a, ok := apps[appID]; ok {
					a.serviceDefinition = service
				} else {
					apps[appID] = app{
						serviceDefinition: service,
						version:           version,
					}
				}
			}
		}
	}

	for _, task := range tasks {
		taskDefinition, err := s.client.DescribeTaskDefinition(ctx, *task.TaskDefinitionArn)
		if err != nil {
			s.logger.Error("failed to load taskDefinition", zap.Error(err))
			continue
		}
		for _, tag := range task.Tags {
			if appID := *tag.Key; *tag.Key == provider.LabelApplication {
				if a, ok := apps[appID]; ok {
					a.taskDefinision = *taskDefinition
				} else {
					apps[appID] = app{
						taskDefinision: *taskDefinition,
						version:        version,
					}
				}
			}
		}
	}
	return apps
}

func (s *store) loadApps() map[string]app {
	apps := s.apps.Load()
	if apps == nil {
		return nil
	}
	return apps.(map[string]app)
}

func (s *store) fetchClusterArns(ctx context.Context) ([]string, error) {
	const maxResults = 100
	var cursor string
	clusterArns := make([]string, 0, maxResults)

	for {
		v, next, err := s.client.ListClusters(ctx, maxResults, cursor)
		if err != nil {
			return nil, err
		}
		clusterArns = append(clusterArns, v...)
		if next == nil {
			break
		}
		cursor = *next
	}

	return clusterArns, nil
}

func (s *store) fetchServices(ctx context.Context, clusterArn string) ([]types.Service, error) {
	const maxCapacity = 5000
	const maxResults = 100
	var cursor string
	serviceArns := make([]string, 0, maxCapacity)

	for {
		v, next, err := s.client.ListServices(ctx, clusterArn, maxResults, cursor)
		if err != nil {
			return nil, err
		}
		serviceArns = append(serviceArns, v...)
		if next == nil {
			break
		}
		cursor = *next
	}

	services, err := s.client.DescribeServices(ctx, serviceArns, clusterArn)
	if err != nil {
		return nil, fmt.Errorf("failed to get list of service. cluster: %s: %w", clusterArn, err)
	}

	return services, nil
}

func (s *store) fetchTasks(ctx context.Context, clusterArn string) ([]types.Task, error) {
	const maxCapacity = 5000
	const maxResults = 100
	var cursor string
	taskArns := make([]string, 0, maxCapacity)

	for {
		v, next, err := s.client.ListServices(ctx, clusterArn, maxResults, cursor)
		if err != nil {
			return nil, err
		}
		taskArns = append(taskArns, v...)
		if next == nil {
			break
		}
		cursor = *next
	}

	tasks, err := s.client.DescribeTasks(ctx, taskArns, clusterArn)
	if err != nil {
		return nil, fmt.Errorf("failed to get list of service. cluster: %s: %w", clusterArn, err)
	}

	return tasks, nil
}
