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
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/provider"
)

// fakeClient embeds provider.Client so that only the methods used by the fetcher need to be implemented.
type fakeClient struct {
	provider.Client

	service         *types.Service
	taskSets        []types.TaskSet
	taskSetsFetched bool
}

func (c *fakeClient) DescribeService(context.Context, types.Service) (*types.Service, error) {
	return c.service, nil
}

func (c *fakeClient) GetServiceTaskSets(context.Context, types.Service) ([]types.TaskSet, error) {
	c.taskSetsFetched = true
	return c.taskSets, nil
}

func (c *fakeClient) GetTasks(context.Context, types.Service) ([]types.Task, error) {
	return []types.Task{{TaskArn: aws.String("arn:task/1")}}, nil
}

func TestFetchResources_TaskSets(t *testing.T) {
	t.Parallel()

	taskSets := []types.TaskSet{{TaskSetArn: aws.String("arn:taskset/1")}}

	tests := []struct {
		name             string
		controller       *types.DeploymentController
		wantTaskSets     int
		wantTaskSetsCall bool
	}{
		{
			name:         "EXTERNAL fetches task sets",
			controller:   &types.DeploymentController{Type: types.DeploymentControllerTypeExternal},
			wantTaskSets: 1, wantTaskSetsCall: true,
		},
		{
			name:       "ECS does not fetch task sets",
			controller: &types.DeploymentController{Type: types.DeploymentControllerTypeEcs},
		},
		{
			name:       "no controller is treated as ECS",
			controller: nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := &fakeClient{
				service:  &types.Service{ServiceName: aws.String("svc"), DeploymentController: tc.controller},
				taskSets: taskSets,
			}
			f := &ECSFetcher{client: client}

			got, err := f.FetchResources(context.Background(), types.Service{})
			assert.NoError(t, err)
			assert.Len(t, got.TaskSets, tc.wantTaskSets)
			assert.Equal(t, tc.wantTaskSetsCall, client.taskSetsFetched)
			// Tasks are always fetched.
			assert.Len(t, got.Tasks, 1)
		})
	}
}
