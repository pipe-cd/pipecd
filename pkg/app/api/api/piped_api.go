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

package api

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kapetaniosci/pipe/pkg/app/api/service/pipedservice"
	"github.com/kapetaniosci/pipe/pkg/datastore"
)

// PipedAPI implements the behaviors for the gRPC definitions of PipedAPI.
type PipedAPI struct {
	applicationStore datastore.ApplicationStore
	commandStore     datastore.CommandStore
	deploymentStore  datastore.DeploymentStore
	environmentStore datastore.EnvironmentStore
	pipedStatsStore  datastore.PipedStatsStore
	pipedStore       datastore.PipedStore
	projectStore     datastore.ProjectStore

	logger *zap.Logger
}

// NewPipedAPI creates a new PipedAPI instance.
func NewPipedAPI(ds datastore.DataStore, logger *zap.Logger) *PipedAPI {
	a := &PipedAPI{
		applicationStore: datastore.NewApplicationStore(ds),
		commandStore:     datastore.NewCommandStore(ds),
		deploymentStore:  datastore.NewDeploymentStore(ds),
		environmentStore: datastore.NewEnvironmentStore(ds),
		pipedStatsStore:  datastore.NewPipedStatsStore(ds),
		pipedStore:       datastore.NewPipedStore(ds),
		projectStore:     datastore.NewProjectStore(ds),
		logger:           logger.Named("piped-api"),
	}
	return a
}

// Register registers all handling of this service into the specified gRPC server.
func (a *PipedAPI) Register(server *grpc.Server) {
	pipedservice.RegisterPipedServiceServer(server, a)
}

// Ping is periodically sent by piped to report its status/stats to API.
// The received stats will be written to the cache immediately.
// The cache data may be lost anytime so we need a singleton Persister
// to persist those data into datastore every n minutes.
func (a *PipedAPI) Ping(ctx context.Context, req *pipedservice.PingRequest) (*pipedservice.PingResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ListApplications returns a list of registered applications
// that should be managed by the requested piped.
// Disabled applications should not be included in the response.
// Piped uses this RPC to fetch and sync the application configuration into its local database.
func (a *PipedAPI) ListApplications(ctx context.Context, req *pipedservice.ListApplicationsRequest) (*pipedservice.ListApplicationsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ListNotCompletedDeployments returns a list of not completed deployments
// which are managed by this piped.
// DeploymentController component uses this RPC to spawns/syncs its local deployment executors.
func (a *PipedAPI) ListNotCompletedDeployments(ctx context.Context, req *pipedservice.ListNotCompletedDeploymentsRequest) (*pipedservice.ListNotCompletedDeploymentsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// CreateDeployment creates/triggers a new deployment for an application
// that is managed by this piped.
// This will be used by DeploymentTrigger component.
func (a *PipedAPI) CreateDeployment(ctx context.Context, req *pipedservice.CreateDeploymentRequest) (*pipedservice.CreateDeploymentResponse, error) {
	err := a.deploymentStore.AddDeployment(ctx, req.Deployment)
	if errors.Cause(err) == datastore.ErrAlreadyExists {
		return nil, status.Error(codes.AlreadyExists, "deployment already exists")
	}
	if err != nil {
		a.logger.Error("failed to create deployment", zap.Error(err))
		return nil, err
	}
	return &pipedservice.CreateDeploymentResponse{}, nil
}

// ReportDeploymentPlanned used by piped to update the status
// of a specific deployment to PLANNED.
func (a *PipedAPI) ReportDeploymentPlanned(ctx context.Context, req *pipedservice.ReportDeploymentPlannedRequest) (*pipedservice.ReportDeploymentPlannedResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ReportDeploymentRunning used by piped to update the status
// of a specific deployment to RUNNING.
func (a *PipedAPI) ReportDeploymentRunning(ctx context.Context, req *pipedservice.ReportDeploymentRunningRequest) (*pipedservice.ReportDeploymentRunningResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ReportDeploymentCompleted used by piped to update the status
// of a specific deployment to SUCCESS | FAILURE | CANCELLED.
func (a *PipedAPI) ReportDeploymentCompleted(ctx context.Context, req *pipedservice.ReportDeploymentCompletedRequest) (*pipedservice.ReportDeploymentCompletedResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// SaveDeploymentMetadata used by piped to persist the metadata of a specific deployment.
func (a *PipedAPI) SaveDeploymentMetadata(ctx context.Context, req *pipedservice.SaveDeploymentMetadataRequest) (*pipedservice.SaveDeploymentMetadataResponse, error) {
	err := a.deploymentStore.PutDeploymentMetadata(ctx, req.DeploymentId, req.Metadata)
	if errors.Cause(err) == datastore.ErrNotFound {
		return nil, status.Error(codes.InvalidArgument, "deployment is not found")
	}
	if err != nil {
		a.logger.Error("failed to save deployment metadata",
			zap.String("deployment-id", req.DeploymentId),
			zap.Error(err),
		)
		return nil, err
	}
	return &pipedservice.SaveDeploymentMetadataResponse{}, nil
}

// SaveStageMetadata used by piped to persist the metadata
// of a specific stage of a deployment.
func (a *PipedAPI) SaveStageMetadata(ctx context.Context, req *pipedservice.SaveStageMetadataRequest) (*pipedservice.SaveStageMetadataResponse, error) {
	err := a.deploymentStore.PutDeploymentStageMetadata(ctx, req.DeploymentId, req.StageId, req.JsonMetadata)
	if err != nil {
		switch errors.Cause(err) {
		case datastore.ErrNotFound:
			return nil, status.Error(codes.InvalidArgument, "deployment is not found")
		case datastore.ErrInvalidArgument:
			return nil, status.Error(codes.InvalidArgument, "invalid value for update")
		default:
			a.logger.Error("failed to save deployment stage metadata",
				zap.String("deployment-id", req.DeploymentId),
				zap.String("stage-id", req.StageId),
				zap.Error(err),
			)
			return nil, err
		}
	}
	return &pipedservice.SaveStageMetadataResponse{}, nil
}

// ReportStageLog is sent by piped to save the log of a pipeline stage.
func (a *PipedAPI) ReportStageLog(ctx context.Context, req *pipedservice.ReportStageLogRequest) (*pipedservice.ReportStageLogResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ReportStageStatusChanged used by piped to update the status
// of a specific stage of a deployment.
func (a *PipedAPI) ReportStageStatusChanged(ctx context.Context, req *pipedservice.ReportStageStatusChangedRequest) (*pipedservice.ReportStageStatusChangedResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ListUnhandledCommands is periodically called by piped to obtain the commands
// that should be handled.
// Whenever an user makes an interaction from WebUI (cancel/approve/retry/sync)
// a new command with a unique identifier will be generated an saved into the datastore.
// Piped uses this RPC to list all still-not-handled commands to handle them,
// then report back the result to server.
// On other side, the web will periodically check the command status and feedback the result to user.
// In the future, we may need a solution to remove all old-handled commands from datastore for space.
func (a *PipedAPI) ListUnhandledCommands(ctx context.Context, req *pipedservice.ListUnhandledCommandsRequest) (*pipedservice.ListUnhandledCommandsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ReportCommandHandled is called by piped to mark a specific command as handled.
// The request payload will contain the handle status as well as any additional result data.
// The handle result should be updated to both datastore and cache (for reading from web).
func (a *PipedAPI) ReportCommandHandled(ctx context.Context, req *pipedservice.ReportCommandHandledRequest) (*pipedservice.ReportCommandHandledResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ReportApplicationState is periodically sent by piped to refresh the current state of application.
// This may contain a full tree of application resources for Kubernetes application.
// The tree data will be written into filestore and the cache inmmediately.
func (a *PipedAPI) ReportApplicationState(ctx context.Context, req *pipedservice.ReportApplicationStateRequest) (*pipedservice.ReportApplicationStateResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
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
func (a *PipedAPI) ReportAppStateEvents(ctx context.Context, req *pipedservice.ReportAppStateEventsRequest) (*pipedservice.ReportAppStateEventsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}
