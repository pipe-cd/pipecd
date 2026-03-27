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
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/provider"

	appconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/config"
)

// fakeLogPersister is a no-op log persister for tests
type fakeLogPersister struct{}

func (f *fakeLogPersister) Write(p []byte) (n int, err error) { return len(p), nil }
func (f *fakeLogPersister) Info(string)                       {}
func (f *fakeLogPersister) Infof(string, ...any)              {}
func (f *fakeLogPersister) Success(string)                    {}
func (f *fakeLogPersister) Successf(string, ...any)           {}
func (f *fakeLogPersister) Error(string)                      {}
func (f *fakeLogPersister) Errorf(string, ...any)             {}
func (f *fakeLogPersister) Complete(time.Duration) error      { return nil }

// mockECSClient is used to mock the provider.Client interface
type mockECSClient struct {
	CreateServiceFunc               func(ctx context.Context, service types.Service) (*types.Service, error)
	UpdateServiceFunc               func(ctx context.Context, service types.Service) (*types.Service, error)
	DescribeServiceFunc             func(ctx context.Context, service types.Service) (*types.Service, error)
	GetServiceTaskSetsFunc          func(ctx context.Context, service types.Service) ([]types.TaskSet, error)
	GetPrimaryTaskSetFunc           func(ctx context.Context, service types.Service) (*types.TaskSet, error)
	CreateTaskSetFunc               func(ctx context.Context, service types.Service, taskDefinition types.TaskDefinition, targetGroup *types.LoadBalancer, scale float64) (*types.TaskSet, error)
	UpdateServicePrimaryTaskSetFunc func(ctx context.Context, service types.Service, taskSet types.TaskSet) (*types.TaskSet, error)
	DeleteTaskSetFunc               func(ctx context.Context, taskSet types.TaskSet) error
	GetTasksFunc                    func(ctx context.Context, service types.Service) ([]types.Task, error)
	ServiceExistsFunc               func(ctx context.Context, cluster, serviceName string) (bool, error)
	GetServiceStatusFunc            func(ctx context.Context, cluster, serviceName string) (string, error)
	WaitServiceStableFunc           func(ctx context.Context, cluster, serviceName string) error
	RegisterTaskDefinitionFunc      func(ctx context.Context, taskDef types.TaskDefinition) (*types.TaskDefinition, error)
	RunTaskFunc                     func(ctx context.Context, taskDefinition types.TaskDefinition, clusterArn string, launchType string, awsVpcConfiguration *appconfig.ECSVpcConfiguration, tags []types.Tag) error
	PruneServiceTasksFunc           func(ctx context.Context, service types.Service) error
	ListTagsFunc                    func(ctx context.Context, resourceArn string) ([]types.Tag, error)
	TagResourceFunc                 func(ctx context.Context, resourceArn string, tags []types.Tag) error
	UntagResourceFunc               func(ctx context.Context, resourceArn string, tagKeys []string) error
}

var _ provider.Client = (*mockECSClient)(nil)

func (m *mockECSClient) CreateService(ctx context.Context, service types.Service) (*types.Service, error) {
	return m.CreateServiceFunc(ctx, service)
}
func (m *mockECSClient) UpdateService(ctx context.Context, service types.Service) (*types.Service, error) {
	return m.UpdateServiceFunc(ctx, service)
}
func (m *mockECSClient) DescribeService(ctx context.Context, service types.Service) (*types.Service, error) {
	return m.DescribeServiceFunc(ctx, service)
}
func (m *mockECSClient) GetServiceTaskSets(ctx context.Context, service types.Service) ([]types.TaskSet, error) {
	return m.GetServiceTaskSetsFunc(ctx, service)
}
func (m *mockECSClient) GetPrimaryTaskSet(ctx context.Context, service types.Service) (*types.TaskSet, error) {
	return m.GetPrimaryTaskSetFunc(ctx, service)
}
func (m *mockECSClient) CreateTaskSet(ctx context.Context, service types.Service, taskDefinition types.TaskDefinition, targetGroup *types.LoadBalancer, scale float64) (*types.TaskSet, error) {
	return m.CreateTaskSetFunc(ctx, service, taskDefinition, targetGroup, scale)
}
func (m *mockECSClient) UpdateServicePrimaryTaskSet(ctx context.Context, service types.Service, taskSet types.TaskSet) (*types.TaskSet, error) {
	return m.UpdateServicePrimaryTaskSetFunc(ctx, service, taskSet)
}
func (m *mockECSClient) DeleteTaskSet(ctx context.Context, taskSet types.TaskSet) error {
	return m.DeleteTaskSetFunc(ctx, taskSet)
}
func (m *mockECSClient) GetTasks(ctx context.Context, service types.Service) ([]types.Task, error) {
	return m.GetTasksFunc(ctx, service)
}
func (m *mockECSClient) ServiceExists(ctx context.Context, cluster, serviceName string) (bool, error) {
	return m.ServiceExistsFunc(ctx, cluster, serviceName)
}
func (m *mockECSClient) GetServiceStatus(ctx context.Context, cluster, serviceName string) (string, error) {
	return m.GetServiceStatusFunc(ctx, cluster, serviceName)
}
func (m *mockECSClient) WaitServiceStable(ctx context.Context, cluster, serviceName string) error {
	return m.WaitServiceStableFunc(ctx, cluster, serviceName)
}
func (m *mockECSClient) RegisterTaskDefinition(ctx context.Context, taskDef types.TaskDefinition) (*types.TaskDefinition, error) {
	return m.RegisterTaskDefinitionFunc(ctx, taskDef)
}
func (m *mockECSClient) RunTask(ctx context.Context, taskDefinition types.TaskDefinition, clusterArn string, launchType string, awsVpcConfiguration *appconfig.ECSVpcConfiguration, tags []types.Tag) error {
	return m.RunTaskFunc(ctx, taskDefinition, clusterArn, launchType, awsVpcConfiguration, tags)
}
func (m *mockECSClient) PruneServiceTasks(ctx context.Context, service types.Service) error {
	return m.PruneServiceTasksFunc(ctx, service)
}
func (m *mockECSClient) ListTags(ctx context.Context, resourceArn string) ([]types.Tag, error) {
	return m.ListTagsFunc(ctx, resourceArn)
}
func (m *mockECSClient) TagResource(ctx context.Context, resourceArn string, tags []types.Tag) error {
	return m.TagResourceFunc(ctx, resourceArn, tags)
}
func (m *mockECSClient) UntagResource(ctx context.Context, resourceArn string, tagKeys []string) error {
	return m.UntagResourceFunc(ctx, resourceArn, tagKeys)
}

func happyPathClient(registeredTD *types.TaskDefinition, updatedSvc *types.Service, newTS *types.TaskSet, prevTaskSets []types.TaskSet) *mockECSClient {
	return &mockECSClient{
		RegisterTaskDefinitionFunc: func(_ context.Context, _ types.TaskDefinition) (*types.TaskDefinition, error) {
			td := *registeredTD
			return &td, nil
		},
		ServiceExistsFunc: func(_ context.Context, _, _ string) (bool, error) {
			return true, nil
		},
		GetServiceStatusFunc: func(_ context.Context, _, _ string) (string, error) {
			return "ACTIVE", nil
		},
		UpdateServiceFunc: func(_ context.Context, _ types.Service) (*types.Service, error) {
			svc := *updatedSvc
			return &svc, nil
		},
		ListTagsFunc: func(_ context.Context, _ string) ([]types.Tag, error) {
			return []types.Tag{}, nil
		},
		TagResourceFunc: func(_ context.Context, _ string, _ []types.Tag) error {
			return nil
		},
		GetServiceTaskSetsFunc: func(_ context.Context, _ types.Service) ([]types.TaskSet, error) {
			return prevTaskSets, nil
		},
		GetPrimaryTaskSetFunc: func(_ context.Context, _ types.Service) (*types.TaskSet, error) {
			if len(prevTaskSets) > 0 {
				ts := prevTaskSets[0]
				return &ts, nil
			}
			return nil, nil
		},
		CreateTaskSetFunc: func(_ context.Context, _ types.Service, _ types.TaskDefinition, _ *types.LoadBalancer, _ float64) (*types.TaskSet, error) {
			ts := *newTS
			return &ts, nil
		},
		UpdateServicePrimaryTaskSetFunc: func(_ context.Context, _ types.Service, _ types.TaskSet) (*types.TaskSet, error) {
			ts := *newTS
			return &ts, nil
		},
		DeleteTaskSetFunc: func(_ context.Context, _ types.TaskSet) error {
			return nil
		},
		WaitServiceStableFunc: func(_ context.Context, _, _ string) error {
			return nil
		},
		PruneServiceTasksFunc: func(_ context.Context, _ types.Service) error {
			return nil
		},
	}
}
