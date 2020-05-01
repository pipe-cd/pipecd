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

package runnerclientfake

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kapetaniosci/pipe/pkg/app/api/service/runnerservice"
	"github.com/kapetaniosci/pipe/pkg/model"
)

type fakeClient struct {
	applications map[string]*model.Application
	deployments  map[string]*model.Deployment
	mu           sync.Mutex
	logger       *zap.Logger
}

// NewClient returns a new fakeClient.
func NewClient(logger *zap.Logger) *fakeClient {
	return &fakeClient{
		applications: map[string]*model.Application{
			"fake-app-1": &model.Application{
				Id:        "fake-app-1",
				Name:      "fake-app-1",
				Env:       "fake-env",
				RunnerId:  "fake-runner-id",
				ProjectId: "fake-project-id",
				Kind:      model.ApplicationKind_KUBERNETES,
				GitPath: &model.ApplicationGitPath{
					Host:   "https://github.com",
					Org:    "fake-org",
					Repo:   "fake-repo",
					Branch: "master",
					Path:   "demoapp",
				},
				Disabled: false,
			},
		},
		deployments: map[string]*model.Deployment{},
		logger:      logger.Named("fake-runner-client"),
	}
}

func (c *fakeClient) Close() error {
	c.logger.Info("fakeClient client is closing")
	return nil
}

func (c *fakeClient) Ping(ctx context.Context, in *runnerservice.PingRequest, opts ...grpc.CallOption) (*runnerservice.PingResponse, error) {
	c.logger.Info("received Ping rpc", zap.Any("request", in))
	return &runnerservice.PingResponse{}, nil
}

func (c *fakeClient) ListApplications(ctx context.Context, in *runnerservice.ListApplicationsRequest, opts ...grpc.CallOption) (*runnerservice.ListApplicationsResponse, error) {
	c.logger.Info("received ListApplications rpc", zap.Any("request", in))
	apps := make([]*model.Application, 0, len(c.applications))
	for _, app := range c.applications {
		apps = append(apps, app)
	}
	return &runnerservice.ListApplicationsResponse{
		Applications: apps,
	}, nil
}

func (c *fakeClient) CreateDeployment(ctx context.Context, in *runnerservice.CreateDeploymentRequest, opts ...grpc.CallOption) (*runnerservice.CreateDeploymentResponse, error) {
	c.logger.Info("received CreateDeployment rpc", zap.Any("request", in))
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.deployments[in.Deployment.Id]; ok {
		return nil, status.Error(codes.AlreadyExists, "")
	}
	c.deployments[in.Deployment.Id] = in.Deployment
	return &runnerservice.CreateDeploymentResponse{}, nil
}

func (c *fakeClient) ListNotCompletedDeployments(ctx context.Context, in *runnerservice.ListNotCompletedDeploymentsRequest, opts ...grpc.CallOption) (*runnerservice.ListNotCompletedDeploymentsResponse, error) {
	c.logger.Info("received ListNotCompletedDeployments rpc", zap.Any("request", in))
	deployments := make([]*model.Deployment, 0, len(c.deployments))
	for _, deployment := range c.deployments {
		if !deployment.IsCompleted() {
			continue
		}
		deployments = append(deployments, deployment)
	}
	return &runnerservice.ListNotCompletedDeploymentsResponse{
		Deployments: deployments,
	}, nil
}

func (c *fakeClient) RegisterEvents(ctx context.Context, in *runnerservice.RegisterEventsRequest, opts ...grpc.CallOption) (*runnerservice.RegisterEventsResponse, error) {
	c.logger.Info("received RegisterEvents rpc", zap.Any("request", in))
	return nil, nil
}

func (c *fakeClient) SendStageLog(ctx context.Context, in *runnerservice.SendStageLogRequest, opts ...grpc.CallOption) (*runnerservice.SendStageLogResponse, error) {
	c.logger.Info("received SendStageLog rpc", zap.Any("request", in))
	return nil, nil
}

func (c *fakeClient) GetCommands(ctx context.Context, in *runnerservice.GetCommandsRequest, opts ...grpc.CallOption) (*runnerservice.GetCommandsResponse, error) {
	c.logger.Info("received GetCommands rpc", zap.Any("request", in))
	return nil, nil
}

func (c *fakeClient) ReportCommandHandled(ctx context.Context, in *runnerservice.ReportCommandHandledRequest, opts ...grpc.CallOption) (*runnerservice.ReportCommandHandledResponse, error) {
	c.logger.Info("received ReportCommandHandled rpc", zap.Any("request", in))
	return nil, nil
}

func (c *fakeClient) ReportDeploymentCompleted(ctx context.Context, in *runnerservice.ReportDeploymentCompletedRequest, opts ...grpc.CallOption) (*runnerservice.ReportDeploymentCompletedResponse, error) {
	c.logger.Info("received ReportDeploymentCompleted rpc", zap.Any("request", in))
	return nil, nil
}

func (c *fakeClient) ReportApplicationState(ctx context.Context, in *runnerservice.ReportApplicationStateRequest, opts ...grpc.CallOption) (*runnerservice.ReportApplicationStateResponse, error) {
	c.logger.Info("received ReportApplicationState rpc", zap.Any("request", in))
	return nil, nil
}

var _ runnerservice.RunnerServiceClient = (*fakeClient)(nil)
