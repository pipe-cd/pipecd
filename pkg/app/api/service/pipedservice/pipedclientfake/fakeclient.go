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

	"github.com/pipe-cd/pipe/pkg/app/api/service/pipedservice"
	"github.com/pipe-cd/pipe/pkg/model"
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

	// Register applications for debug repository.
	for name, enable := range k8sAppNames {
		app := &model.Application{
			Id:            projectID + "/" + envID + "/" + name,
			Name:          name,
			EnvId:         envID,
			PipedId:       pipedID,
			ProjectId:     projectID,
			Kind:          model.ApplicationKind_KUBERNETES,
			CloudProvider: "kubernetes-default",
			GitPath: &model.ApplicationGitPath{
				Repo: &model.ApplicationGitRepository{
					Id:     "debug",
					Remote: "git@github.com:pipe-cd/debug.git",
					Branch: "master",
				},
				Path: "kubernetes/" + name,
			},
			Disabled: !enable,
		}
		apps[app.Id] = app
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

// Ping is periodically sent to report its realtime status/stats to control-plane.
// The received stats will be pushed to the metrics collector.
func (c *fakeClient) Ping(ctx context.Context, req *pipedservice.PingRequest, opts ...grpc.CallOption) (*pipedservice.PingResponse, error) {
	c.logger.Info("fake client received Ping rpc", zap.Any("request", req))
	return &pipedservice.PingResponse{}, nil
}

// ReportStat is periodically sent to report its realtime status/stats to control-plane.
// The received stats will be pushed to the metrics collector.
func (c *fakeClient) ReportStat(ctx context.Context, req *pipedservice.ReportStatRequest, opts ...grpc.CallOption) (*pipedservice.ReportStatResponse, error) {
	c.logger.Info("fake client received ReportStat rpc", zap.Any("request", req))
	return &pipedservice.ReportStatResponse{}, nil
}

// ReportPipedMeta is sent by piped while starting up to report its metadata
// such as configured cloud providers.
func (c *fakeClient) ReportPipedMeta(ctx context.Context, req *pipedservice.ReportPipedMetaRequest, opts ...grpc.CallOption) (*pipedservice.ReportPipedMetaResponse, error) {
	c.logger.Info("fake client received ReportPipedMeta rpc", zap.Any("request", req))
	return &pipedservice.ReportPipedMetaResponse{}, nil
}

// GetEnvironment finds and returns the environment for the specified ID.
func (c *fakeClient) GetEnvironment(ctx context.Context, req *pipedservice.GetEnvironmentRequest, opts ...grpc.CallOption) (*pipedservice.GetEnvironmentResponse, error) {
	c.logger.Info("fake client received GetEnvironment rpc", zap.Any("request", req))
	return &pipedservice.GetEnvironmentResponse{
		Environment: &model.Environment{
			Id:   "dev",
			Name: "dev",
		},
	}, nil
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

// ReportApplicationSyncState is used to update the sync status of an application.
func (c *fakeClient) ReportApplicationSyncState(ctx context.Context, req *pipedservice.ReportApplicationSyncStateRequest, opts ...grpc.CallOption) (*pipedservice.ReportApplicationSyncStateResponse, error) {
	c.logger.Info("fake client received ReportApplicationSyncState rpc", zap.Any("request", req))
	c.mu.RLock()
	defer c.mu.RUnlock()

	app, ok := c.applications[req.ApplicationId]
	if !ok {
		return nil, status.Error(codes.NotFound, "application was not found")
	}
	app.SyncState = req.State

	return &pipedservice.ReportApplicationSyncStateResponse{}, nil
}

// ReportApplicationDeployingStatus is used to report whether the specified application is deploying or not.
func (c *fakeClient) ReportApplicationDeployingStatus(_ context.Context, req *pipedservice.ReportApplicationDeployingStatusRequest, _ ...grpc.CallOption) (*pipedservice.ReportApplicationDeployingStatusResponse, error) {
	c.logger.Info("fake client received ReportApplicationDeployingStatus rpc", zap.Any("request", req))
	c.mu.RLock()
	defer c.mu.RUnlock()

	app, ok := c.applications[req.ApplicationId]
	if !ok {
		return nil, status.Error(codes.NotFound, "application was not found")
	}
	app.Deploying = req.Deploying

	return &pipedservice.ReportApplicationDeployingStatusResponse{}, nil
}

// ReportApplicationMostRecentDeployment is used to update the basic information about
// the most recent deployment of a specific application.
func (c *fakeClient) ReportApplicationMostRecentDeployment(ctx context.Context, req *pipedservice.ReportApplicationMostRecentDeploymentRequest, opts ...grpc.CallOption) (*pipedservice.ReportApplicationMostRecentDeploymentResponse, error) {
	c.logger.Info("fake client received ReportApplicationMostRecentDeployment rpc", zap.Any("request", req))

	c.mu.RLock()
	defer c.mu.RUnlock()

	app, ok := c.applications[req.ApplicationId]
	if !ok {
		return nil, status.Error(codes.NotFound, "application was not found")
	}

	switch req.Status {
	case model.DeploymentStatus_DEPLOYMENT_SUCCESS:
		app.MostRecentlySuccessfulDeployment = req.Deployment

	case model.DeploymentStatus_DEPLOYMENT_PENDING:
		app.MostRecentlyTriggeredDeployment = req.Deployment
	}

	return &pipedservice.ReportApplicationMostRecentDeploymentResponse{}, nil
}

// GetApplicationMostRecentDeployment returns the most recent deployment of the given application.
func (c *fakeClient) GetApplicationMostRecentDeployment(ctx context.Context, req *pipedservice.GetApplicationMostRecentDeploymentRequest, opts ...grpc.CallOption) (*pipedservice.GetApplicationMostRecentDeploymentResponse, error) {
	c.logger.Info("fake client received GetApplicationMostRecentDeployment rpc", zap.Any("request", req))

	c.mu.RLock()
	defer c.mu.RUnlock()

	app, ok := c.applications[req.ApplicationId]
	if !ok {
		return nil, status.Error(codes.NotFound, "application was not found")
	}

	if req.Status == model.DeploymentStatus_DEPLOYMENT_SUCCESS && app.MostRecentlySuccessfulDeployment != nil {
		return &pipedservice.GetApplicationMostRecentDeploymentResponse{Deployment: app.MostRecentlySuccessfulDeployment}, nil
	}

	if req.Status == model.DeploymentStatus_DEPLOYMENT_PENDING && app.MostRecentlyTriggeredDeployment != nil {
		return &pipedservice.GetApplicationMostRecentDeploymentResponse{Deployment: app.MostRecentlyTriggeredDeployment}, nil
	}

	return nil, status.Error(codes.NotFound, "")
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

	if req.Summary != "" {
		d.Summary = req.Summary
	}
	d.Status = s
	d.StatusReason = req.StatusReason
	d.RunningCommitHash = req.RunningCommitHash
	d.Version = req.Version
	if len(req.Stages) > 0 {
		d.Stages = req.Stages
	}

	return &pipedservice.ReportDeploymentPlannedResponse{}, nil
}

// ReportDeploymentStatusChanged is used to update the status
// of a specific deployment to RUNNING or ROLLING_BACK.
func (c *fakeClient) ReportDeploymentStatusChanged(ctx context.Context, req *pipedservice.ReportDeploymentStatusChangedRequest, opts ...grpc.CallOption) (*pipedservice.ReportDeploymentStatusChangedResponse, error) {
	c.logger.Info("fake client received ReportDeploymentStatusChanged rpc", zap.Any("request", req))
	c.mu.Lock()
	defer c.mu.Unlock()

	d, ok := c.deployments[req.DeploymentId]
	if !ok {
		return nil, status.Error(codes.NotFound, "deployment was not found")
	}

	if !model.CanUpdateDeploymentStatus(d.Status, req.Status) {
		msg := fmt.Sprintf("invalid status, cur = %s, req = %s", d.Status.String(), req.Status.String())
		return nil, status.Error(codes.FailedPrecondition, msg)
	}

	d.Status = req.Status
	d.StatusReason = req.StatusReason
	return &pipedservice.ReportDeploymentStatusChangedResponse{}, nil
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
	d.StatusReason = req.StatusReason
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
		s.Metadata = req.Metadata
		return &pipedservice.SaveStageMetadataResponse{}, nil
	}
	return nil, status.Error(codes.NotFound, "stage was not found")
}

// ReportStageLogs is sent by piped to save the log of a pipeline stage.
func (c *fakeClient) ReportStageLogs(ctx context.Context, req *pipedservice.ReportStageLogsRequest, opts ...grpc.CallOption) (*pipedservice.ReportStageLogsResponse, error) {
	c.logger.Info("fake client received ReportStageLogs rpc", zap.Any("request", req))
	return &pipedservice.ReportStageLogsResponse{}, nil
}

// ReportStageLogsFromLastCheckpoint is used to save the full logs from the most recently saved point.
func (c *fakeClient) ReportStageLogsFromLastCheckpoint(ctx context.Context, req *pipedservice.ReportStageLogsFromLastCheckpointRequest, opts ...grpc.CallOption) (*pipedservice.ReportStageLogsFromLastCheckpointResponse, error) {
	c.logger.Info("fake client received ReportStageLogsFromLastCheckpoint rpc", zap.Any("request", req))
	return &pipedservice.ReportStageLogsFromLastCheckpointResponse{}, nil
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
		s.Visible = req.Visible
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

// ReportApplicationLiveState is periodically sent to correct full state of an application.
// For kubernetes application, this contains a full tree of its kubernetes resources.
// The tree data should be written into filestore immediately and then the state in cache should be refreshsed too.
func (c *fakeClient) ReportApplicationLiveState(ctx context.Context, req *pipedservice.ReportApplicationLiveStateRequest, opts ...grpc.CallOption) (*pipedservice.ReportApplicationLiveStateResponse, error) {
	c.logger.Info("fake client received ReportApplicationLiveState rpc", zap.Any("request", req))
	return &pipedservice.ReportApplicationLiveStateResponse{}, nil
}

// ReportApplicationLiveStateEvents is sent by piped to submit one or multiple events
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
func (c *fakeClient) ReportApplicationLiveStateEvents(ctx context.Context, req *pipedservice.ReportApplicationLiveStateEventsRequest, opts ...grpc.CallOption) (*pipedservice.ReportApplicationLiveStateEventsResponse, error) {
	c.logger.Info("fake client received ReportApplicationLiveStateEvents rpc", zap.Any("request", req))
	return &pipedservice.ReportApplicationLiveStateEventsResponse{}, nil
}

func (c *fakeClient) GetLatestEvent(ctx context.Context, req *pipedservice.GetLatestEventRequest, opts ...grpc.CallOption) (*pipedservice.GetLatestEventResponse, error) {
	c.logger.Info("fake client received GetLatestEvent rpc", zap.Any("request", req))
	return &pipedservice.GetLatestEventResponse{
		Event: &model.Event{
			Id:        "dev",
			Name:      "dev",
			ProjectId: "dev",
		},
	}, nil
}

func (c *fakeClient) ListEvents(ctx context.Context, req *pipedservice.ListEventsRequest, opts ...grpc.CallOption) (*pipedservice.ListEventsResponse, error) {
	c.logger.Info("fake client received ListEvents rpc", zap.Any("request", req))
	return &pipedservice.ListEventsResponse{}, nil
}

func (a *fakeClient) GetLatestAnalysisResult(ctx context.Context, req *pipedservice.GetLatestAnalysisResultRequest, opts ...grpc.CallOption) (*pipedservice.GetLatestAnalysisResultResponse, error) {
	a.logger.Info("fake client received GetLatestAnalysisResult rpc", zap.Any("request", req))
	return &pipedservice.GetLatestAnalysisResultResponse{}, nil
}

func (a *fakeClient) PutLatestAnalysisResult(ctx context.Context, req *pipedservice.PutLatestAnalysisResultRequest, opts ...grpc.CallOption) (*pipedservice.PutLatestAnalysisResultResponse, error) {
	a.logger.Info("fake client received PutLatestAnalysisResult rpc", zap.Any("request", req))
	return &pipedservice.PutLatestAnalysisResultResponse{}, nil
}

var _ pipedservice.PipedServiceClient = (*fakeClient)(nil)
