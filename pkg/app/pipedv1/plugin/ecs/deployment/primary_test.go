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

package deployment

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrimaryRollout(t *testing.T) {
	t.Parallel()

	var (
		clusterArn = "arn:aws:ecs:us-east-1:123456789012:cluster/my-cluster"
		serviceArn = "arn:aws:ecs:us-east-1:123456789012:service/my-cluster/my-service"
		taskDefArn = "arn:aws:ecs:us-east-1:123456789012:task-definition/my-family:1"
		tsArn1     = "arn:aws:ecs:us-east-1:123456789012:task-set/my-cluster/my-service/ecs-svc:1"
		newTSArn   = "arn:aws:ecs:us-east-1:123456789012:task-set/my-cluster/my-service/ecs-svc:2"

		baseTaskDef = types.TaskDefinition{
			Family:            aws.String("my-family"),
			TaskDefinitionArn: aws.String(taskDefArn),
		}
		baseServiceDef = types.Service{
			ClusterArn:   aws.String(clusterArn),
			ServiceName:  aws.String("my-service"),
			ServiceArn:   aws.String(serviceArn),
			DesiredCount: 2,
		}
		registeredTD = &types.TaskDefinition{
			Family:            aws.String("my-family"),
			TaskDefinitionArn: aws.String(taskDefArn),
		}
		updatedService = &types.Service{
			ClusterArn:   aws.String(clusterArn),
			ServiceName:  aws.String("my-service"),
			ServiceArn:   aws.String(serviceArn),
			DesiredCount: 2,
		}
		newTaskSet  = &types.TaskSet{TaskSetArn: aws.String(newTSArn)}
		prevTaskSet = types.TaskSet{
			TaskSetArn: aws.String(tsArn1),
			ClusterArn: aws.String(clusterArn),
			ServiceArn: aws.String(serviceArn),
		}
		primaryLB = &types.LoadBalancer{
			TargetGroupArn: aws.String("arn:aws:elasticloadbalancing:us-east-1:123:targetgroup/primary/abc"),
		}
	)

	testcases := []struct {
		name       string
		taskDef    types.TaskDefinition
		serviceDef types.Service
		primary    *types.LoadBalancer
		client     *mockECSClient
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "success: existing service, no previous task sets",
			taskDef:    baseTaskDef,
			serviceDef: baseServiceDef,
			client:     happyPathClient(registeredTD, updatedService, newTaskSet, []types.TaskSet{}),
		},
		{
			name:       "success: existing service, previous task set is deleted",
			taskDef:    baseTaskDef,
			serviceDef: baseServiceDef,
			client:     happyPathClient(registeredTD, updatedService, newTaskSet, []types.TaskSet{prevTaskSet}),
		},
		{
			name:       "success: service does not exist, new service is created",
			taskDef:    baseTaskDef,
			serviceDef: baseServiceDef,
			client: func() *mockECSClient {
				m := happyPathClient(registeredTD, updatedService, newTaskSet, []types.TaskSet{})
				m.ServiceExistsFunc = func(_ context.Context, _, _ string) (bool, error) { return false, nil }
				m.CreateServiceFunc = func(_ context.Context, _ types.Service) (*types.Service, error) {
					svc := *updatedService
					return &svc, nil
				}
				return m
			}(),
		},
		{
			name:       "success: with primary ELB target group at scale 100",
			taskDef:    baseTaskDef,
			serviceDef: baseServiceDef,
			primary:    primaryLB,
			client: func() *mockECSClient {
				m := happyPathClient(registeredTD, updatedService, newTaskSet, []types.TaskSet{})
				m.CreateTaskSetFunc = func(_ context.Context, _ types.Service, _ types.TaskDefinition, lb *types.LoadBalancer, scale float64) (*types.TaskSet, error) {
					assert.Equal(t, primaryLB, lb)
					assert.Equal(t, float64(100), scale)
					return newTaskSet, nil
				}
				return m
			}(),
		},
		{
			name:       "fail: RegisterTaskDefinition error",
			taskDef:    baseTaskDef,
			serviceDef: baseServiceDef,
			client: func() *mockECSClient {
				m := happyPathClient(registeredTD, updatedService, newTaskSet, []types.TaskSet{})
				m.RegisterTaskDefinitionFunc = func(_ context.Context, _ types.TaskDefinition) (*types.TaskDefinition, error) {
					return nil, errors.New("register error")
				}
				return m
			}(),
			wantErr:    true,
			wantErrMsg: "failed to apply task definition",
		},
		{
			name:       "fail: ServiceExists error",
			taskDef:    baseTaskDef,
			serviceDef: baseServiceDef,
			client: func() *mockECSClient {
				m := happyPathClient(registeredTD, updatedService, newTaskSet, []types.TaskSet{})
				m.ServiceExistsFunc = func(_ context.Context, _, _ string) (bool, error) {
					return false, errors.New("describe error")
				}
				return m
			}(),
			wantErr:    true,
			wantErrMsg: "failed to apply service definition",
		},
		{
			name:       "fail: service is DRAINING (not updatable)",
			taskDef:    baseTaskDef,
			serviceDef: baseServiceDef,
			client: func() *mockECSClient {
				m := happyPathClient(registeredTD, updatedService, newTaskSet, []types.TaskSet{})
				m.GetServiceStatusFunc = func(_ context.Context, _, _ string) (string, error) {
					return "DRAINING", nil
				}
				return m
			}(),
			wantErr:    true,
			wantErrMsg: "failed to apply service definition",
		},
		{
			name:       "fail: UpdateService error during apply",
			taskDef:    baseTaskDef,
			serviceDef: baseServiceDef,
			client: func() *mockECSClient {
				m := happyPathClient(registeredTD, updatedService, newTaskSet, []types.TaskSet{})
				m.UpdateServiceFunc = func(_ context.Context, _ types.Service) (*types.Service, error) {
					return nil, errors.New("update error")
				}
				return m
			}(),
			wantErr:    true,
			wantErrMsg: "failed to apply service definition",
		},
		{
			name:       "fail: GetPrimaryTaskSet error",
			taskDef:    baseTaskDef,
			serviceDef: baseServiceDef,
			client: func() *mockECSClient {
				m := happyPathClient(registeredTD, updatedService, newTaskSet, []types.TaskSet{})
				m.GetPrimaryTaskSetFunc = func(_ context.Context, _ types.Service) (*types.TaskSet, error) {
					return nil, errors.New("get primary task set error")
				}
				return m
			}(),
			wantErr:    true,
			wantErrMsg: "failed to get current primary taskset",
		},
		{
			name:       "fail: CreateTaskSet error",
			taskDef:    baseTaskDef,
			serviceDef: baseServiceDef,
			client: func() *mockECSClient {
				m := happyPathClient(registeredTD, updatedService, newTaskSet, []types.TaskSet{})
				m.CreateTaskSetFunc = func(_ context.Context, _ types.Service, _ types.TaskDefinition, _ *types.LoadBalancer, _ float64) (*types.TaskSet, error) {
					return nil, errors.New("create task set error")
				}
				return m
			}(),
			wantErr:    true,
			wantErrMsg: "failed to create primary taskset for service",
		},
		{
			name:       "fail: UpdateServicePrimaryTaskSet error",
			taskDef:    baseTaskDef,
			serviceDef: baseServiceDef,
			client: func() *mockECSClient {
				m := happyPathClient(registeredTD, updatedService, newTaskSet, []types.TaskSet{})
				m.UpdateServicePrimaryTaskSetFunc = func(_ context.Context, _ types.Service, _ types.TaskSet) (*types.TaskSet, error) {
					return nil, errors.New("update primary error")
				}
				return m
			}(),
			wantErr:    true,
			wantErrMsg: "failed to create primary taskset for service",
		},
		{
			name:       "fail: DeleteTaskSet error",
			taskDef:    baseTaskDef,
			serviceDef: baseServiceDef,
			client: func() *mockECSClient {
				m := happyPathClient(registeredTD, updatedService, newTaskSet, []types.TaskSet{prevTaskSet})
				m.DeleteTaskSetFunc = func(_ context.Context, _ types.TaskSet) error {
					return errors.New("delete error")
				}
				return m
			}(),
			wantErr:    true,
			wantErrMsg: "failed to delete old primary taskset",
		},
		{
			name:       "fail: WaitServiceStable error",
			taskDef:    baseTaskDef,
			serviceDef: baseServiceDef,
			client: func() *mockECSClient {
				m := happyPathClient(registeredTD, updatedService, newTaskSet, []types.TaskSet{})
				m.WaitServiceStableFunc = func(_ context.Context, _, _ string) error {
					return errors.New("wait stable error")
				}
				return m
			}(),
			wantErr:    true,
			wantErrMsg: "did not reach stable state",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := primaryRollout(context.Background(), &fakeLogPersister{}, tc.client, tc.taskDef, tc.serviceDef, tc.primary)

			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErrMsg)
				return
			}
			require.NoError(t, err)
		})
	}
}
