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

// Close closes the connection to server.
func (c *fakeClient) Close() error {
	c.logger.Info("fakeClient client is closing")
	return nil
}

// Ping is periodically sent by runner to report its status/stats to API.
// The received stats will be written to the cache immediately.
// The cache data may be lost anytime so we need a singleton Persister
// to persist those data into datastore every n minutes.
func (c *fakeClient) Ping(ctx context.Context, in *runnerservice.PingRequest, opts ...grpc.CallOption) (*runnerservice.PingResponse, error) {
	c.logger.Info("received Ping rpc", zap.Any("request", in))
	return &runnerservice.PingResponse{}, nil
}

// ListApplications returns a list of registered applications
// that should be managed by the requested runner.
// Disabled applications should not be included in the response.
// Runner uses this RPC to fetch and sync the application configuration into its local database.
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

// CreateDeployment creates/triggers a new deployment for an application
// that is managed by this runner.
// This will be used by DeploymentTrigger component.
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

// ListNotCompletedDeployments returns a list of not completed deployments
// which are managed by this runner.
// DeploymentController component uses this RPC to spawns/syncs its local deployment executors.
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

// SaveStageMetadata used by runner to persist the metadata
// of a specific stage of a deployment.
func (c *fakeClient) SaveStageMetadata(ctx context.Context, in *runnerservice.SaveStageMetadataRequest, opts ...grpc.CallOption) (*runnerservice.SaveStageMetadataResponse, error) {
	c.logger.Info("received SaveStageMetadata rpc", zap.Any("request", in))
	return nil, nil
}

// ReportStageStatusChanged used by runner to update the status
// of a specific stage of a deployment.
func (c *fakeClient) ReportStageStatusChanged(ctx context.Context, in *runnerservice.ReportStageStatusChangedRequest, opts ...grpc.CallOption) (*runnerservice.ReportStageStatusChangedResponse, error) {
	c.logger.Info("received ReportStageStatusChanged rpc", zap.Any("request", in))
	return nil, nil
}

// ReportStageLog is sent by runner to save the log of a pipeline stage.
func (c *fakeClient) ReportStageLog(ctx context.Context, in *runnerservice.ReportStageLogRequest, opts ...grpc.CallOption) (*runnerservice.ReportStageLogResponse, error) {
	c.logger.Info("received ReportStageLog rpc", zap.Any("request", in))
	return nil, nil
}

// ReportDeploymentCompleted used by runner to send the final state
// of the pipeline that has just been completed.
func (c *fakeClient) ReportDeploymentCompleted(ctx context.Context, in *runnerservice.ReportDeploymentCompletedRequest, opts ...grpc.CallOption) (*runnerservice.ReportDeploymentCompletedResponse, error) {
	c.logger.Info("received ReportDeploymentCompleted rpc", zap.Any("request", in))
	return nil, nil
}

// GetCommands is periodically called by runner to obtain the commands
// that should be handled.
// Whenever an user makes an interaction from WebUI (cancel/approve/retry/sync)
// a new command with a unique identifier will be generated an saved into the datastore.
// Runner uses this RPC to list all still-not-handled commands to handle them,
// then report back the result to server.
// On other side, the web will periodically check the command status and feedback the result to user.
// In the future, we may need a solution to remove all old-handled commands from datastore for space.
func (c *fakeClient) GetCommands(ctx context.Context, in *runnerservice.GetCommandsRequest, opts ...grpc.CallOption) (*runnerservice.GetCommandsResponse, error) {
	c.logger.Info("received GetCommands rpc", zap.Any("request", in))
	return nil, nil
}

// ReportCommandHandled is called by runner to mark a specific command as handled.
// The request payload will contain the handle status as well as any additional result data.
// The handle result should be updated to both datastore and cache (for reading from web).
func (c *fakeClient) ReportCommandHandled(ctx context.Context, in *runnerservice.ReportCommandHandledRequest, opts ...grpc.CallOption) (*runnerservice.ReportCommandHandledResponse, error) {
	c.logger.Info("received ReportCommandHandled rpc", zap.Any("request", in))
	return nil, nil
}

// ReportApplicationState is periodically sent by runner to refresh the current state of application.
// This may contain a full tree of application resources for Kubernetes application.
// The tree data will be written into filestore and the cache inmmediately.
func (c *fakeClient) ReportApplicationState(ctx context.Context, in *runnerservice.ReportApplicationStateRequest, opts ...grpc.CallOption) (*runnerservice.ReportApplicationStateResponse, error) {
	c.logger.Info("received ReportApplicationState rpc", zap.Any("request", in))
	return nil, nil
}

// ReportAppStateEvents is sent by runner to submit one or multiple events
// about the changes of application state.
// Control plane uses the received events to update the state of application-resource-tree.
// We want to start by a simple solution at this initial stage of development,
// so the API server just handles as below:
// - loads the releated application-resource-tree from filestore
// - checks and builds new state for the application-resource-tree
// - updates new state into fielstore and cache (cache data is for reading while handling web requests)
// In the future, we may want to redesign the behavior of this RPC by using pubsub/queue pattern.
// After receiving the events, all of them will be publish into a queue immediately,
// and then another Handler service will pick them inorder to apply to build new state.
// By that way we can control the traffic to the datastore in a better way.
func (c *fakeClient) ReportAppStateEvents(ctx context.Context, in *runnerservice.ReportAppStateEventsRequest, opts ...grpc.CallOption) (*runnerservice.ReportAppStateEventsResponse, error) {
	c.logger.Info("received ReportAppStateEvents rpc", zap.Any("request", in))
	return nil, nil
}

var _ runnerservice.RunnerServiceClient = (*fakeClient)(nil)
