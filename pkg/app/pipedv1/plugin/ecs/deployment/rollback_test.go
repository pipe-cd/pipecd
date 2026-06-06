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

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/provider"
)

func TestRollBack(t *testing.T) {
	t.Parallel()

	var (
		clusterArn = "arn:aws:ecs:us-east-1:123456789012:cluster/my-cluster"
		serviceArn = "arn:aws:ecs:us-east-1:123456789012:service/my-cluster/my-service"
		taskDefArn = "arn:aws:ecs:us-east-1:123456789012:task-definition/my-family:2"
		tsArn1     = "arn:aws:ecs:us-east-1:123456789012:task-set/my-cluster/my-service/ecs-svc:1"
		tsArn2     = "arn:aws:ecs:us-east-1:123456789012:task-set/my-cluster/my-service/ecs-svc:2"
		newTSArn   = "arn:aws:ecs:us-east-1:123456789012:task-set/my-cluster/my-service/ecs-svc:3"

		baseTaskDef = types.TaskDefinition{
			Family:            aws.String("my-family"),
			TaskDefinitionArn: aws.String(taskDefArn),
		}
		baseServiceDef = types.Service{
			ClusterArn:  aws.String(clusterArn),
			ServiceName: aws.String("my-service"),
			ServiceArn:  aws.String(serviceArn),
		}
		registeredTD = &types.TaskDefinition{
			Family:            aws.String("my-family"),
			TaskDefinitionArn: aws.String(taskDefArn),
		}
		updatedService = &types.Service{
			ClusterArn:  aws.String(clusterArn),
			ServiceName: aws.String("my-service"),
			ServiceArn:  aws.String(serviceArn),
		}
		newTaskSet = &types.TaskSet{
			TaskSetArn: aws.String(newTSArn),
		}
		prevTaskSet1 = types.TaskSet{
			TaskSetArn: aws.String(tsArn1),
			ClusterArn: aws.String(clusterArn),
			ServiceArn: aws.String(serviceArn),
		}
		prevTaskSet2 = types.TaskSet{
			TaskSetArn: aws.String(tsArn2),
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
			name:       "success: existing service, two previous task sets are deleted",
			taskDef:    baseTaskDef,
			serviceDef: baseServiceDef,
			client:     happyPathClient(registeredTD, updatedService, newTaskSet, []types.TaskSet{prevTaskSet1, prevTaskSet2}),
		},
		{
			name:       "success: with primary ELB target group passed to CreateTaskSet at scale 100",
			taskDef:    baseTaskDef,
			serviceDef: baseServiceDef,
			primary:    primaryLB,
			client: func() *mockECSClient {
				m := happyPathClient(registeredTD, updatedService, newTaskSet, []types.TaskSet{})
				m.CreateTaskSetFunc = func(_ context.Context, _ types.Service, _ types.TaskDefinition, lb *types.LoadBalancer, scale float64) (*types.TaskSet, error) {
					assert.Equal(t, primaryLB, lb)
					assert.Equal(t, float64(100), scale)
					ts := *newTaskSet
					return &ts, nil
				}
				return m
			}(),
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
			name:       "success: stale tags on existing service are removed",
			taskDef:    baseTaskDef,
			serviceDef: baseServiceDef,
			client: func() *mockECSClient {
				m := happyPathClient(registeredTD, updatedService, newTaskSet, []types.TaskSet{})
				m.ListTagsFunc = func(_ context.Context, _ string) ([]types.Tag, error) {
					return []types.Tag{{Key: aws.String("old-tag"), Value: aws.String("old-val")}}, nil
				}
				m.UntagResourceFunc = func(_ context.Context, _ string, keys []string) error {
					assert.Equal(t, []string{"old-tag"}, keys)
					return nil
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
			wantErrMsg: "failed to register task definition my-family",
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
			wantErrMsg: "failed to apply service definition for service my-service",
		},
		{
			name:       "fail: GetServiceStatus error",
			taskDef:    baseTaskDef,
			serviceDef: baseServiceDef,
			client: func() *mockECSClient {
				m := happyPathClient(registeredTD, updatedService, newTaskSet, []types.TaskSet{})
				m.GetServiceStatusFunc = func(_ context.Context, _, _ string) (string, error) {
					return "", errors.New("status error")
				}
				return m
			}(),
			wantErr:    true,
			wantErrMsg: "failed to apply service definition for service my-service",
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
			wantErrMsg: "failed to apply service definition for service my-service",
		},
		{
			name:       "fail: UpdateService error",
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
			wantErrMsg: "failed to apply service definition for service my-service",
		},
		{
			name:       "fail: ListTags error",
			taskDef:    baseTaskDef,
			serviceDef: baseServiceDef,
			client: func() *mockECSClient {
				m := happyPathClient(registeredTD, updatedService, newTaskSet, []types.TaskSet{})
				m.ListTagsFunc = func(_ context.Context, _ string) ([]types.Tag, error) {
					return nil, errors.New("list tags error")
				}
				return m
			}(),
			wantErr:    true,
			wantErrMsg: "failed to apply service definition for service my-service",
		},
		{
			name:       "fail: TagResource error",
			taskDef:    baseTaskDef,
			serviceDef: baseServiceDef,
			client: func() *mockECSClient {
				m := happyPathClient(registeredTD, updatedService, newTaskSet, []types.TaskSet{})
				m.TagResourceFunc = func(_ context.Context, _ string, _ []types.Tag) error {
					return errors.New("tag error")
				}
				return m
			}(),
			wantErr:    true,
			wantErrMsg: "failed to apply service definition for service my-service",
		},
		{
			name:       "fail: CreateService error when service does not exist",
			taskDef:    baseTaskDef,
			serviceDef: baseServiceDef,
			client: func() *mockECSClient {
				m := happyPathClient(registeredTD, updatedService, newTaskSet, []types.TaskSet{})
				m.ServiceExistsFunc = func(_ context.Context, _, _ string) (bool, error) { return false, nil }
				m.CreateServiceFunc = func(_ context.Context, _ types.Service) (*types.Service, error) {
					return nil, errors.New("create error")
				}
				return m
			}(),
			wantErr:    true,
			wantErrMsg: "failed to apply service definition for service my-service",
		},
		{
			name:       "fail: GetServiceTaskSets error",
			taskDef:    baseTaskDef,
			serviceDef: baseServiceDef,
			client: func() *mockECSClient {
				m := happyPathClient(registeredTD, updatedService, newTaskSet, []types.TaskSet{})
				m.GetServiceTaskSetsFunc = func(_ context.Context, _ types.Service) ([]types.TaskSet, error) {
					return nil, errors.New("get task sets error")
				}
				return m
			}(),
			wantErr:    true,
			wantErrMsg: "failed to get task sets for service my-service",
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
			wantErrMsg: "failed to create task set for service my-service",
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
			wantErrMsg: "failed to update primary task set for service my-service",
		},
		{
			name:       "fail: DeleteTaskSet error stops at first failure",
			taskDef:    baseTaskDef,
			serviceDef: baseServiceDef,
			client: func() *mockECSClient {
				m := happyPathClient(registeredTD, updatedService, newTaskSet, []types.TaskSet{prevTaskSet1, prevTaskSet2})
				m.DeleteTaskSetFunc = func(_ context.Context, ts types.TaskSet) error {
					if aws.ToString(ts.TaskSetArn) == tsArn1 {
						return errors.New("delete error")
					}
					return nil
				}
				return m
			}(),
			wantErr:    true,
			wantErrMsg: "failed to delete task set",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := &externalController{}
			err := ctrl.Rollback(context.Background(), &fakeLogPersister{}, tc.client, tc.taskDef, tc.serviceDef, tc.primary)
			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErrMsg)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestRestoreELBWeights(t *testing.T) {
	t.Parallel()

	var (
		primaryARN   = "arn:aws:elasticloadbalancing:us-east-1:123:targetgroup/primary/aaa"
		canaryARN    = "arn:aws:elasticloadbalancing:us-east-1:123:targetgroup/canary/bbb"
		listenerARN1 = "arn:aws:elasticloadbalancing:us-east-1:123:listener/app/my-alb/aaa/bbb"
		listenerARN2 = "arn:aws:elasticloadbalancing:us-east-1:123:listener/app/my-alb/aaa/ccc"
		primaryLB    = &types.LoadBalancer{TargetGroupArn: aws.String(primaryARN)}
	)

	testcases := []struct {
		name       string
		primary    *types.LoadBalancer
		metadata   *mockMetadataStore
		client     *mockECSClient
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:    "success: listener ARNs and canary ARN in metadata, weights restored to 100/0",
			primary: primaryLB,
			metadata: func() *mockMetadataStore {
				m := happyMetadataStore()
				m.GetFunc = func(_ context.Context, key string) (string, bool, error) {
					switch key {
					case currentListenersKey:
						return listenerARN1 + "," + listenerARN2, true, nil
					case canaryTargetGroupArnKey:
						return canaryARN, true, nil
					}
					return "", false, nil
				}
				return m
			}(),
			client: func() *mockECSClient {
				m := &mockECSClient{}
				m.ModifyListenersFunc = func(_ context.Context, listenerArns []string, cfg provider.RoutingTrafficConfig) ([]string, error) {
					assert.Equal(t, []string{listenerARN1, listenerARN2}, listenerArns)
					assert.Equal(t, provider.RoutingTrafficConfig{
						{TargetGroupArn: primaryARN, Weight: 100},
						{TargetGroupArn: canaryARN, Weight: 0},
					}, cfg)
					return []string{"rule-1"}, nil
				}
				return m
			}(),
		},
		{
			name:     "success (no-op): no listener ARNs in metadata (traffic routing never ran)",
			primary:  primaryLB,
			metadata: happyMetadataStore(), // GetFunc returns (_, false, nil) for all keys
			client:   &mockECSClient{},
		},
		{
			name:    "success (no-op): listener ARNs found but no canary ARN in metadata",
			primary: primaryLB,
			metadata: func() *mockMetadataStore {
				m := happyMetadataStore()
				m.GetFunc = func(_ context.Context, key string) (string, bool, error) {
					if key == currentListenersKey {
						return listenerARN1, true, nil
					}
					return "", false, nil
				}
				return m
			}(),
			client: &mockECSClient{},
		},
		{
			name:    "fail: GetDeploymentPluginMetadata error for listener ARNs",
			primary: primaryLB,
			metadata: func() *mockMetadataStore {
				m := happyMetadataStore()
				m.GetFunc = func(_ context.Context, key string) (string, bool, error) {
					if key == currentListenersKey {
						return "", false, errors.New("metadata error")
					}
					return "", false, nil
				}
				return m
			}(),
			client:     &mockECSClient{},
			wantErr:    true,
			wantErrMsg: "failed to get listener ARNs from metadata",
		},
		{
			name:    "fail: GetDeploymentPluginMetadata error for canary ARN",
			primary: primaryLB,
			metadata: func() *mockMetadataStore {
				m := happyMetadataStore()
				m.GetFunc = func(_ context.Context, key string) (string, bool, error) {
					switch key {
					case currentListenersKey:
						return listenerARN1, true, nil
					case canaryTargetGroupArnKey:
						return "", false, errors.New("metadata error")
					}
					return "", false, nil
				}
				return m
			}(),
			client:     &mockECSClient{},
			wantErr:    true,
			wantErrMsg: "failed to get canary target group ARN from metadata",
		},
		{
			name:    "fail: ModifyListeners error",
			primary: primaryLB,
			metadata: func() *mockMetadataStore {
				m := happyMetadataStore()
				m.GetFunc = func(_ context.Context, key string) (string, bool, error) {
					switch key {
					case currentListenersKey:
						return listenerARN1, true, nil
					case canaryTargetGroupArnKey:
						return canaryARN, true, nil
					}
					return "", false, nil
				}
				return m
			}(),
			client: func() *mockECSClient {
				m := &mockECSClient{}
				m.ModifyListenersFunc = func(_ context.Context, _ []string, _ provider.RoutingTrafficConfig) ([]string, error) {
					return nil, errors.New("modify listeners error")
				}
				return m
			}(),
			wantErr:    true,
			wantErrMsg: "failed to restore ELB listener weights",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := restoreELBWeights(context.Background(), &fakeLogPersister{}, tc.metadata, tc.client, tc.primary)
			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErrMsg)
				return
			}
			require.NoError(t, err)
		})
	}
}
