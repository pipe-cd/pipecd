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

package pipedclientfake

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kapetaniosci/pipe/pkg/app/api/service/pipedservice"
	"github.com/kapetaniosci/pipe/pkg/model"
)

type fakeClient struct {
	applications map[string]*model.Application
	deployments  map[string]*model.Deployment
	mu           sync.RWMutex
	logger       *zap.Logger
}

// NewClient returns a new fakeClient.
func NewClient(logger *zap.Logger) *fakeClient {
	var (
		projectID   = "local-project"
		envID       = "dev"
		pipedID     = "local-piped"
		apps        = make(map[string]*model.Application, 0)
		k8sAppNames = map[string]bool{
			"analysis-by-http":       false,
			"analysis-by-log":        false,
			"analysis-by-metrics":    false,
			"analysis-with-baseline": false,
			"bluegreen":              false,
			"canary":                 true,
			"helm-local-chart":       false,
			"helm-remote-chart":      false,
			"helm-remote-git-chart":  false,
			"kustomize-local-base":   false,
			"kustomize-remote-base":  false,
			"mesh-envoy-bluegreen":   false,
			"mesh-envoy-canary":      false,
			"mesh-istio-bluegreen":   false,
			"mesh-istio-canary":      false,
			"multi-steps-canary":     false,
			"simple":                 false,
			"wait-approval":          false,
		}
	)

	// Register applications for pipe-debug repository.
	for name, enable := range k8sAppNames {
		apps[name] = &model.Application{
			Id:        projectID + "/" + envID + "/" + name,
			Name:      name,
			EnvId:     envID,
			PipedId:   pipedID,
			ProjectId: projectID,
			Kind:      model.ApplicationKind_KUBERNETES,
			GitPath: &model.ApplicationGitPath{
				RepoId: "pipe-debug",
				Path:   "k8s/" + name,
			},
			Disabled: !enable,
		}
	}

	return &fakeClient{
		applications: apps,
		deployments:  map[string]*model.Deployment{},
		logger:       logger.Named("fake-piped-client"),
	}
}

// Close closes the connection to server.
func (c *fakeClient) Close() error {
	c.logger.Info("fakeClient client is closing")
	return nil
}

// Ping is periodically sent by piped to report its status/stats to API.
// The received stats will be written to the cache immediately.
// The cache data may be lost anytime so we need a singleton Persister
// to persist those data into datastore every n minutes.
func (c *fakeClient) Ping(ctx context.Context, req *pipedservice.PingRequest, opts ...grpc.CallOption) (*pipedservice.PingResponse, error) {
	c.logger.Info("fake client received Ping rpc", zap.Any("request", req))
	return &pipedservice.PingResponse{}, nil
}

// ReportPipedMeta is sent by piped while starting up to report its metadata
// such as configured cloud providers.
func (c *fakeClient) ReportPipedMeta(ctx context.Context, req *pipedservice.ReportPipedMetaRequest, opts ...grpc.CallOption) (*pipedservice.ReportPipedMetaResponse, error) {
	c.logger.Info("fake client received ReportPipedMeta rpc", zap.Any("request", req))
	return &pipedservice.ReportPipedMetaResponse{}, nil
}

// ListApplications returns a list of registered applications
// that should be managed by the requested piped.
// Disabled applications should not be included in the response.
// Piped uses this RPC to fetch and sync the application configuration into its local database.
func (c *fakeClient) ListApplications(ctx context.Context, req *pipedservice.ListApplicationsRequest, opts ...grpc.CallOption) (*pipedservice.ListApplicationsResponse, error) {
	c.logger.Info("fake client received ListApplications rpc", zap.Any("request", req))
	apps := make([]*model.Application, 0, len(c.applications))
	for _, app := range c.applications {
		if app.Disabled {
			continue
		}
		apps = append(apps, app)
	}
	return &pipedservice.ListApplicationsResponse{
		Applications: apps,
	}, nil
}

// ListNotCompletedDeployments returns a list of not completed deployments
// which are managed by this piped.
// DeploymentController component uses this RPC to spawns/syncs its local deployment executors.
func (c *fakeClient) ListNotCompletedDeployments(ctx context.Context, req *pipedservice.ListNotCompletedDeploymentsRequest, opts ...grpc.CallOption) (*pipedservice.ListNotCompletedDeploymentsResponse, error) {
	c.logger.Info("fake client received ListNotCompletedDeployments rpc", zap.Any("request", req))
	c.mu.RLock()
	defer c.mu.RUnlock()

	deployments := make([]*model.Deployment, 0, len(c.deployments))
	for _, d := range c.deployments {
		if model.IsCompletedDeployment(d.Status) {
			continue
		}
		deployments = append(deployments, d.Clone())
	}
	return &pipedservice.ListNotCompletedDeploymentsResponse{
		Deployments: deployments,
	}, nil
}

// GetMostRecentDeployment returns the most recent deployment of the given application.
func (c *fakeClient) GetMostRecentDeployment(ctx context.Context, req *pipedservice.GetMostRecentDeploymentRequest, opts ...grpc.CallOption) (*pipedservice.GetMostRecentDeploymentResponse, error) {
	c.logger.Info("fake client received GetMostRecentDeployment rpc", zap.Any("request", req))

	var (
		requiredStatus model.DeploymentStatus
		hasStatus      bool
		mostRecent     *model.Deployment
	)
	if req.Status != nil {
		hasStatus = true
		requiredStatus = model.DeploymentStatus(req.Status.Value)
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, d := range c.deployments {
		if hasStatus && d.Status != requiredStatus {
			continue
		}
		if mostRecent != nil && !d.TriggerBefore(mostRecent) {
			continue
		}
		mostRecent = d
	}

	if mostRecent == nil {
		return nil, status.Error(codes.NotFound, "")
	}

	return &pipedservice.GetMostRecentDeploymentResponse{
		Deployment: mostRecent.Clone(),
	}, nil
}

// CreateDeployment creates/triggers a new deployment for an application
// that is managed by this piped.
// This will be used by DeploymentTrigger component.
func (c *fakeClient) CreateDeployment(ctx context.Context, req *pipedservice.CreateDeploymentRequest, opts ...grpc.CallOption) (*pipedservice.CreateDeploymentResponse, error) {
	c.logger.Info("fake client received CreateDeployment rpc", zap.Any("request", req))
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.deployments[req.Deployment.Id]; ok {
		return nil, status.Error(codes.AlreadyExists, "")
	}
	c.deployments[req.Deployment.Id] = req.Deployment
	return &pipedservice.CreateDeploymentResponse{}, nil
}

// ReportDeploymentPlanned used by piped to update the status
// of a specific deployment to PLANNED.
func (c *fakeClient) ReportDeploymentPlanned(ctx context.Context, req *pipedservice.ReportDeploymentPlannedRequest, opts ...grpc.CallOption) (*pipedservice.ReportDeploymentPlannedResponse, error) {
	c.logger.Info("fake client received ReportDeploymentPlanned rpc", zap.Any("request", req))
	c.mu.Lock()
	defer c.mu.Unlock()

	d, ok := c.deployments[req.DeploymentId]
	if !ok {
		return nil, status.Error(codes.NotFound, "deployment was not found")
	}

	s := model.DeploymentStatus_DEPLOYMENT_PLANNED
	if !model.CanUpdateDeploymentStatus(d.Status, s) {
		msg := fmt.Sprintf("invalid status, cur = %s, req = %s", d.Status.String(), s.String())
		return nil, status.Error(codes.FailedPrecondition, msg)
	}

	if req.Description != "" {
		d.Description = req.Description
	}
	d.Status = s
	d.StatusDescription = req.StatusDescription
	if len(req.Stages) > 0 {
		d.Stages = req.Stages
	}

	return &pipedservice.ReportDeploymentPlannedResponse{}, nil
}

// ReportDeploymentRunning used by piped to update the status
// of a specific deployment to RUNNING.
func (c *fakeClient) ReportDeploymentRunning(ctx context.Context, req *pipedservice.ReportDeploymentRunningRequest, opts ...grpc.CallOption) (*pipedservice.ReportDeploymentRunningResponse, error) {
	c.logger.Info("fake client received ReportDeploymentRunning rpc", zap.Any("request", req))
	c.mu.Lock()
	defer c.mu.Unlock()

	d, ok := c.deployments[req.DeploymentId]
	if !ok {
		return nil, status.Error(codes.NotFound, "deployment was not found")
	}

	s := model.DeploymentStatus_DEPLOYMENT_RUNNING
	if !model.CanUpdateDeploymentStatus(d.Status, s) {
		msg := fmt.Sprintf("invalid status, cur = %s, req = %s", d.Status.String(), s.String())
		return nil, status.Error(codes.FailedPrecondition, msg)
	}

	d.Status = s
	d.StatusDescription = req.StatusDescription
	return &pipedservice.ReportDeploymentRunningResponse{}, nil
}

// ReportDeploymentCompleted used by piped to update the status
// of a specific deployment to SUCCESS | FAILURE | CANCELLED.
func (c *fakeClient) ReportDeploymentCompleted(ctx context.Context, req *pipedservice.ReportDeploymentCompletedRequest, opts ...grpc.CallOption) (*pipedservice.ReportDeploymentCompletedResponse, error) {
	c.logger.Info("fake client received ReportDeploymentCompleted rpc", zap.Any("request", req))
	c.mu.Lock()
	defer c.mu.Unlock()

	d, ok := c.deployments[req.DeploymentId]
	if !ok {
		return nil, status.Error(codes.NotFound, "deployment was not found")
	}

	if !model.IsCompletedDeployment(req.Status) {
		msg := fmt.Sprintf("invalid status, expected a completed one but got  %s", req.Status.String())
		return nil, status.Error(codes.FailedPrecondition, msg)
	}

	if !model.CanUpdateDeploymentStatus(d.Status, req.Status) {
		msg := fmt.Sprintf("invalid status, cur = %s, req = %s", d.Status.String(), req.Status.String())
		return nil, status.Error(codes.FailedPrecondition, msg)
	}

	d.Status = req.Status
	d.StatusDescription = req.StatusDescription
	d.CompletedAt = req.CompletedAt
	for _, stage := range d.Stages {
		if status, ok := req.StageStatuses[stage.Id]; ok {
			stage.Status = status
		}
	}

	return &pipedservice.ReportDeploymentCompletedResponse{}, nil
}

// SaveDeploymentMetadata used by piped to persist the metadata of a specific deployment.
func (c *fakeClient) SaveDeploymentMetadata(ctx context.Context, req *pipedservice.SaveDeploymentMetadataRequest, opts ...grpc.CallOption) (*pipedservice.SaveDeploymentMetadataResponse, error) {
	c.logger.Info("fake client received SaveDeploymentMetadata rpc", zap.Any("request", req))
	c.mu.Lock()
	defer c.mu.Unlock()

	d, ok := c.deployments[req.DeploymentId]
	if !ok {
		return nil, status.Error(codes.NotFound, "deployment was not found")
	}

	d.Metadata = req.Metadata
	return &pipedservice.SaveDeploymentMetadataResponse{}, nil
}

// SaveStageMetadata used by piped to persist the metadata
// of a specific stage of a deployment.
func (c *fakeClient) SaveStageMetadata(ctx context.Context, req *pipedservice.SaveStageMetadataRequest, opts ...grpc.CallOption) (*pipedservice.SaveStageMetadataResponse, error) {
	c.logger.Info("fake client received SaveStageMetadata rpc", zap.Any("request", req))
	c.mu.Lock()
	defer c.mu.Unlock()

	d, ok := c.deployments[req.DeploymentId]
	if !ok {
		return nil, status.Error(codes.NotFound, "deployment was not found")
	}

	for _, s := range d.Stages {
		if s.Id != req.StageId {
			continue
		}
		s.JsonMetadata = req.JsonMetadata
		return &pipedservice.SaveStageMetadataResponse{}, nil
	}
	return nil, status.Error(codes.NotFound, "stage was not found")
}

// ReportStageLog is sent by piped to save the log of a pipeline stage.
func (c *fakeClient) ReportStageLog(ctx context.Context, req *pipedservice.ReportStageLogRequest, opts ...grpc.CallOption) (*pipedservice.ReportStageLogResponse, error) {
	c.logger.Info("fake client received ReportStageLog rpc", zap.Any("request", req))
	return &pipedservice.ReportStageLogResponse{}, nil
}

// ReportStageStatusChanged used by piped to update the status
// of a specific stage of a deployment.
func (c *fakeClient) ReportStageStatusChanged(ctx context.Context, req *pipedservice.ReportStageStatusChangedRequest, opts ...grpc.CallOption) (*pipedservice.ReportStageStatusChangedResponse, error) {
	c.logger.Info("fake client received ReportStageStatusChanged rpc", zap.Any("request", req))
	c.mu.Lock()
	defer c.mu.Unlock()

	d, ok := c.deployments[req.DeploymentId]
	if !ok {
		return nil, status.Error(codes.NotFound, "deployment was not found")
	}

	for _, s := range d.Stages {
		if s.Id != req.StageId {
			continue
		}
		s.Status = req.Status
		s.RetriedCount = req.RetriedCount
		s.CompletedAt = req.CompletedAt
		return &pipedservice.ReportStageStatusChangedResponse{}, nil
	}
	return nil, status.Error(codes.NotFound, "stage was not found")
}

// ListUnhandledCommands is periodically called by piped to obtain the commands
// that should be handled.
// Whenever an user makes an interaction from WebUI (cancel/approve/retry/sync)
// a new command with a unique identifier will be generated an saved into the datastore.
// Piped uses this RPC to list all still-not-handled commands to handle them,
// then report back the result to server.
// On other side, the web will periodically check the command status and feedback the result to user.
// In the future, we may need a solution to remove all old-handled commands from datastore for space.
func (c *fakeClient) ListUnhandledCommands(ctx context.Context, req *pipedservice.ListUnhandledCommandsRequest, opts ...grpc.CallOption) (*pipedservice.ListUnhandledCommandsResponse, error) {
	c.logger.Info("fake client received ListUnhandledCommands rpc", zap.Any("request", req))
	return &pipedservice.ListUnhandledCommandsResponse{}, nil
}

// ReportCommandHandled is called by piped to mark a specific command as handled.
// The request payload will contain the handle status as well as any additional result data.
// The handle result should be updated to both datastore and cache (for reading from web).
func (c *fakeClient) ReportCommandHandled(ctx context.Context, req *pipedservice.ReportCommandHandledRequest, opts ...grpc.CallOption) (*pipedservice.ReportCommandHandledResponse, error) {
	c.logger.Info("fake client received ReportCommandHandled rpc", zap.Any("request", req))
	return &pipedservice.ReportCommandHandledResponse{}, nil
}

// ReportApplicationLiveState is periodically sent by piped to refresh the current state of application.
// This may contain a full tree of application resources for Kubernetes application.
// The tree data will be written into filestore and the cache inmmediately.
func (c *fakeClient) ReportApplicationLiveState(ctx context.Context, req *pipedservice.ReportApplicationLiveStateRequest, opts ...grpc.CallOption) (*pipedservice.ReportApplicationLiveStateResponse, error) {
	c.logger.Info("fake client received ReportApplicationLiveState rpc", zap.Any("request", req))
	return &pipedservice.ReportApplicationLiveStateResponse{}, nil
}

// ReportAppStateEvents is sent by piped to submit one or multiple events
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
func (c *fakeClient) ReportAppStateEvents(ctx context.Context, req *pipedservice.ReportAppStateEventsRequest, opts ...grpc.CallOption) (*pipedservice.ReportAppStateEventsResponse, error) {
	c.logger.Info("fake client received ReportAppStateEvents rpc", zap.Any("request", req))
	return &pipedservice.ReportAppStateEventsResponse{}, nil
}

var _ pipedservice.PipedServiceClient = (*fakeClient)(nil)
