// Copyright 2026 The PipeCD Authors.
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

package livestate

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/provider"
)

// ECSFetcher fetches the current state of an ECS application's resources
type ECSFetcher struct {
	client provider.Client
}

// QueryResources fetches the live Service, its PipeCD-managed TaskSets, and all Tasks for the given service descriptor
func (w *ECSFetcher) FetchResources(ctx context.Context, service types.Service) (*queryResourcesResult, error) {
	liveService, err := w.client.DescribeService(ctx, service)
	if err != nil {
		return nil, err
	}

	// GetServiceTaskSets filters out DRAINING task sets and task sets not created by Pipecd,
	// so the result only contains task sets that Pipecd is responsible for
	tasksets, err := w.client.GetServiceTaskSets(ctx, *liveService)
	if err != nil {
		return nil, err
	}

	// Tasks are fetched for the whole service rather than per task set.
	// Grouping tasks under their parent task set is deferred to the caller.
	tasks, err := w.client.GetTasks(ctx, *liveService)
	if err != nil {
		return nil, err
	}

	return &queryResourcesResult{
		Service:  liveService,
		TaskSets: tasksets,
		Tasks:    tasks,
	}, nil
}

// queryResourcesResult holds the raw AWS objects for a single ECS application.
//
// Tasks are a flat list, not group by taskset ID yet
type queryResourcesResult struct {
	Service  *types.Service
	TaskSets []types.TaskSet
	Tasks    []types.Task
}
