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

package grpcapi

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipe/pkg/app/api/analysisresultstore"
	"github.com/pipe-cd/pipe/pkg/app/api/applicationlivestatestore"
	"github.com/pipe-cd/pipe/pkg/app/api/commandstore"
	"github.com/pipe-cd/pipe/pkg/app/api/service/pipedservice"
	"github.com/pipe-cd/pipe/pkg/app/api/stagelogstore"
	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/cache/memorycache"
	"github.com/pipe-cd/pipe/pkg/cache/rediscache"
	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/filestore"
	"github.com/pipe-cd/pipe/pkg/model"
	"github.com/pipe-cd/pipe/pkg/redis"
	"github.com/pipe-cd/pipe/pkg/rpc/rpcauth"
)

// PipedAPI implements the behaviors for the gRPC definitions of PipedAPI.
type PipedAPI struct {
	applicationStore          datastore.ApplicationStore
	deploymentStore           datastore.DeploymentStore
	deploymentChainStore      datastore.DeploymentChainStore
	environmentStore          datastore.EnvironmentStore
	pipedStore                datastore.PipedStore
	projectStore              datastore.ProjectStore
	eventStore                datastore.EventStore
	stageLogStore             stagelogstore.Store
	applicationLiveStateStore applicationlivestatestore.Store
	analysisResultStore       analysisresultstore.Store
	commandStore              commandstore.Store
	commandOutputPutter       commandOutputPutter

	appPipedCache        cache.Cache
	deploymentPipedCache cache.Cache
	envProjectCache      cache.Cache
	pipedStatCache       cache.Cache
	redis                redis.Redis

	webBaseURL string
	logger     *zap.Logger
}

// NewPipedAPI creates a new PipedAPI instance.
func NewPipedAPI(ctx context.Context, ds datastore.DataStore, sls stagelogstore.Store, alss applicationlivestatestore.Store, las analysisresultstore.Store, cs commandstore.Store, hc cache.Cache, rd redis.Redis, cop commandOutputPutter, webBaseURL string, logger *zap.Logger) *PipedAPI {
	a := &PipedAPI{
		applicationStore:          datastore.NewApplicationStore(ds),
		deploymentStore:           datastore.NewDeploymentStore(ds),
		deploymentChainStore:      datastore.NewDeploymentChainStore(ds),
		environmentStore:          datastore.NewEnvironmentStore(ds),
		pipedStore:                datastore.NewPipedStore(ds),
		projectStore:              datastore.NewProjectStore(ds),
		eventStore:                datastore.NewEventStore(ds),
		stageLogStore:             sls,
		applicationLiveStateStore: alss,
		analysisResultStore:       las,
		commandStore:              cs,
		commandOutputPutter:       cop,
		appPipedCache:             memorycache.NewTTLCache(ctx, 24*time.Hour, 3*time.Hour),
		deploymentPipedCache:      memorycache.NewTTLCache(ctx, 24*time.Hour, 3*time.Hour),
		envProjectCache:           memorycache.NewTTLCache(ctx, 24*time.Hour, 3*time.Hour),
		pipedStatCache:            hc,
		redis:                     rd,
		webBaseURL:                webBaseURL,
		logger:                    logger.Named("piped-api"),
	}
	return a
}

// Register registers all handling of this service into the specified gRPC server.
func (a *PipedAPI) Register(server *grpc.Server) {
	pipedservice.RegisterPipedServiceServer(server, a)
}

// Ping is periodically sent to report its realtime status/stats to control-plane.
// The received stats will be pushed to the metrics collector.
// Note: This rpc is deprecated, use ReportStat instead.
func (a *PipedAPI) Ping(ctx context.Context, req *pipedservice.PingRequest) (*pipedservice.PingResponse, error) {
	return &pipedservice.PingResponse{}, nil
	// return nil, status.Error(codes.Unimplemented, "")
}

// ReportStat is periodically sent to report its realtime status/stats to control-plane.
// The received stats will be pushed to the metrics collector.
func (a *PipedAPI) ReportStat(ctx context.Context, req *pipedservice.ReportStatRequest) (*pipedservice.ReportStatResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	val, err := json.Marshal(model.PipedStat{PipedId: pipedID, Metrics: req.PipedStats, Timestamp: now})
	if err != nil {
		a.logger.Error("failed to store the reported piped stat",
			zap.String("piped-id", pipedID),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to encode the reported piped stat")
	}

	if err := a.pipedStatCache.Put(pipedID, val); err != nil {
		a.logger.Error("failed to store the reported piped stat",
			zap.String("piped-id", pipedID),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to store the reported piped stat")
	}
	return &pipedservice.ReportStatResponse{}, nil
}

// ReportPipedMeta is sent by piped while starting up to report its metadata
// such as configured cloud providers.
func (a *PipedAPI) ReportPipedMeta(ctx context.Context, req *pipedservice.ReportPipedMetaRequest) (*pipedservice.ReportPipedMetaResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	connStatus := model.Piped_ONLINE

	if err = a.pipedStore.UpdatePiped(ctx, pipedID, datastore.PipedMetadataUpdater(req.CloudProviders, req.Repositories, connStatus, req.SecretEncryption, req.Version, now)); err != nil {
		switch err {
		case datastore.ErrNotFound:
			return nil, status.Error(codes.InvalidArgument, "piped is not found")
		case datastore.ErrInvalidArgument:
			return nil, status.Error(codes.InvalidArgument, "invalid value for update")
		default:
			a.logger.Error("failed to update the piped metadata",
				zap.String("piped-id", pipedID),
				zap.Error(err),
			)
			return nil, status.Error(codes.Internal, "failed to update the piped metadata")
		}
	}

	piped, err := getPiped(ctx, a.pipedStore, pipedID, a.logger)
	if err != nil {
		return nil, err
	}
	return &pipedservice.ReportPipedMetaResponse{
		Name:       piped.Name,
		WebBaseUrl: a.webBaseURL,
	}, nil
}

// GetEnvironment finds and returns the environment for the specified ID.
func (a *PipedAPI) GetEnvironment(ctx context.Context, req *pipedservice.GetEnvironmentRequest) (*pipedservice.GetEnvironmentResponse, error) {
	projectID, _, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}
	if err := a.validateEnvBelongsToProject(ctx, req.Id, projectID); err != nil {
		return nil, err
	}

	env, err := a.environmentStore.GetEnvironment(ctx, req.Id)
	if errors.Is(err, datastore.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "environment is not found")
	}
	if err != nil {
		a.logger.Error("failed to get environment", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get environment")
	}
	return &pipedservice.GetEnvironmentResponse{
		Environment: env,
	}, nil
}

// ListApplications returns a list of registered applications
// that should be managed by the requested piped.
// Disabled applications should not be included in the response.
// Piped uses this RPC to fetch and sync the application configuration into its local database.
func (a *PipedAPI) ListApplications(ctx context.Context, req *pipedservice.ListApplicationsRequest) (*pipedservice.ListApplicationsResponse, error) {
	projectID, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}
	opts := datastore.ListOptions{
		Filters: []datastore.ListFilter{
			{
				Field:    "ProjectId",
				Operator: datastore.OperatorEqual,
				Value:    projectID,
			},
			{
				Field:    "PipedId",
				Operator: datastore.OperatorEqual,
				Value:    pipedID,
			},
			{
				Field:    "Disabled",
				Operator: datastore.OperatorEqual,
				Value:    false,
			},
		},
	}
	// TODO: Support pagination in ListApplications
	apps, _, err := a.applicationStore.ListApplications(ctx, opts)
	if err != nil {
		a.logger.Error("failed to fetch applications", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to fetch applications")
	}
	return &pipedservice.ListApplicationsResponse{
		Applications: apps,
	}, nil
}

// ReportApplicationSyncState is used to update the sync status of an application.
func (a *PipedAPI) ReportApplicationSyncState(ctx context.Context, req *pipedservice.ReportApplicationSyncStateRequest) (*pipedservice.ReportApplicationSyncStateResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}
	if err := a.validateAppBelongsToPiped(ctx, req.ApplicationId, pipedID); err != nil {
		return nil, err
	}

	err = a.applicationStore.PutApplicationSyncState(ctx, req.ApplicationId, req.State)
	if err != nil {
		switch err {
		case datastore.ErrNotFound:
			return nil, status.Error(codes.InvalidArgument, "application is not found")
		case datastore.ErrInvalidArgument:
			return nil, status.Error(codes.InvalidArgument, "invalid value for update")
		default:
			a.logger.Error("failed to update application sync state",
				zap.String("application-id", req.ApplicationId),
				zap.Error(err),
			)
			return nil, status.Error(codes.Internal, "failed to update the application sync state")
		}
	}
	return &pipedservice.ReportApplicationSyncStateResponse{}, nil
}

// ReportApplicationDeployingStatus is used to report whether the specified application is deploying or not.
func (a *PipedAPI) ReportApplicationDeployingStatus(ctx context.Context, req *pipedservice.ReportApplicationDeployingStatusRequest) (*pipedservice.ReportApplicationDeployingStatusResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}
	if err := a.validateAppBelongsToPiped(ctx, req.ApplicationId, pipedID); err != nil {
		return nil, err
	}

	err = a.applicationStore.UpdateApplication(ctx, req.ApplicationId, func(app *model.Application) error {
		app.Deploying = req.Deploying
		return nil
	})
	if err == nil {
		return &pipedservice.ReportApplicationDeployingStatusResponse{}, nil
	}

	switch err {
	case datastore.ErrNotFound:
		return nil, status.Error(codes.InvalidArgument, "application is not found")
	case datastore.ErrInvalidArgument:
		return nil, status.Error(codes.InvalidArgument, "invalid value for update")
	default:
		a.logger.Error("failed to update deploying status of application",
			zap.String("application-id", req.ApplicationId),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to update deploying status of application")
	}
}

// ReportApplicationMostRecentDeployment is used to update the basic information about
// the most recent deployment of a specific application.
func (a *PipedAPI) ReportApplicationMostRecentDeployment(ctx context.Context, req *pipedservice.ReportApplicationMostRecentDeploymentRequest) (*pipedservice.ReportApplicationMostRecentDeploymentResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}
	if err := a.validateAppBelongsToPiped(ctx, req.ApplicationId, pipedID); err != nil {
		return nil, err
	}

	err = a.applicationStore.PutApplicationMostRecentDeployment(ctx, req.ApplicationId, req.Status, req.Deployment)
	if err != nil {
		switch err {
		case datastore.ErrNotFound:
			return nil, status.Error(codes.InvalidArgument, "application is not found")
		case datastore.ErrInvalidArgument:
			return nil, status.Error(codes.InvalidArgument, "invalid value for update")
		default:
			a.logger.Error("failed to update application completed deployment",
				zap.String("application-id", req.ApplicationId),
				zap.Error(err),
			)
			return nil, status.Error(codes.Internal, "failed to update the application completed deployment")
		}
	}
	return &pipedservice.ReportApplicationMostRecentDeploymentResponse{}, nil
}

// GetApplicationMostRecentDeployment returns the most recent deployment of the given application.
func (a *PipedAPI) GetApplicationMostRecentDeployment(ctx context.Context, req *pipedservice.GetApplicationMostRecentDeploymentRequest) (*pipedservice.GetApplicationMostRecentDeploymentResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}
	if err := a.validateAppBelongsToPiped(ctx, req.ApplicationId, pipedID); err != nil {
		return nil, err
	}

	app, err := a.applicationStore.GetApplication(ctx, req.ApplicationId)
	if errors.Is(err, datastore.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "application is not found")
	}
	if err != nil {
		a.logger.Error("failed to get application", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get application")
	}

	if req.Status == model.DeploymentStatus_DEPLOYMENT_SUCCESS && app.MostRecentlySuccessfulDeployment != nil {
		return &pipedservice.GetApplicationMostRecentDeploymentResponse{Deployment: app.MostRecentlySuccessfulDeployment}, nil
	}

	if req.Status == model.DeploymentStatus_DEPLOYMENT_PENDING && app.MostRecentlyTriggeredDeployment != nil {
		return &pipedservice.GetApplicationMostRecentDeploymentResponse{Deployment: app.MostRecentlyTriggeredDeployment}, nil
	}

	return nil, status.Error(codes.NotFound, "deployment is not found")
}

func (a *PipedAPI) GetDeployment(ctx context.Context, req *pipedservice.GetDeploymentRequest) (*pipedservice.GetDeploymentResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}

	deployment, err := getDeployment(ctx, a.deploymentStore, req.Id, a.logger)
	if err != nil {
		return nil, err
	}

	if deployment.PipedId != pipedID {
		return nil, status.Error(codes.PermissionDenied, "requested deployment doesn't belong to the piped")
	}

	return &pipedservice.GetDeploymentResponse{
		Deployment: deployment,
	}, nil
}

// ListNotCompletedDeployments returns a list of not completed deployments
// which are managed by this piped.
// DeploymentController component uses this RPC to spawns/syncs its local deployment executors.
func (a *PipedAPI) ListNotCompletedDeployments(ctx context.Context, req *pipedservice.ListNotCompletedDeploymentsRequest) (*pipedservice.ListNotCompletedDeploymentsResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}

	opts := datastore.ListOptions{
		Filters: []datastore.ListFilter{
			{
				Field:    "PipedId",
				Operator: datastore.OperatorEqual,
				Value:    pipedID,
			},
			// TODO: Change to simple conditional clause without using OR clause for portability
			// Note: firestore does not support OR operator.
			// See more: https://firebase.google.com/docs/firestore/query-data/queries?hl=en
			{
				Field:    "Status",
				Operator: datastore.OperatorIn,
				Value:    model.GetNotCompletedDeploymentStatuses(),
			},
		},
	}

	deployments, cursor, err := a.deploymentStore.ListDeployments(ctx, opts)
	if err != nil {
		a.logger.Error("failed to fetch deployments", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to fetch deployments")
	}
	return &pipedservice.ListNotCompletedDeploymentsResponse{
		Deployments: deployments,
		Cursor:      cursor,
	}, nil
}

// CreateDeployment creates/triggers a new deployment for an application
// that is managed by this piped.
// This will be used by DeploymentTrigger component.
func (a *PipedAPI) CreateDeployment(ctx context.Context, req *pipedservice.CreateDeploymentRequest) (*pipedservice.CreateDeploymentResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}
	if err := a.validateAppBelongsToPiped(ctx, req.Deployment.ApplicationId, pipedID); err != nil {
		return nil, err
	}

	err = a.deploymentStore.AddDeployment(ctx, req.Deployment)
	if errors.Is(err, datastore.ErrAlreadyExists) {
		return nil, status.Error(codes.AlreadyExists, "deployment already exists")
	}
	if err != nil {
		a.logger.Error("failed to create deployment", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create deployment")
	}

	// If the deployment doesn't belong to another chain, return immediately.
	if req.Deployment.DeploymentChainId == "" {
		return &pipedservice.CreateDeploymentResponse{}, nil
	}

	// Otherwise, add the created deployment ref to its chain block model.
	if err = a.deploymentChainStore.UpdateDeploymentChain(ctx, req.Deployment.DeploymentChainId, datastore.DeploymentChainAddDeploymentToBlock(req.Deployment)); err != nil {
		return nil, status.Error(codes.Internal, "failed to add deployment ref to its chain block")
	}
	return &pipedservice.CreateDeploymentResponse{}, nil
}

// ReportDeploymentPlanned used by piped to update the status
// of a specific deployment to PLANNED.
func (a *PipedAPI) ReportDeploymentPlanned(ctx context.Context, req *pipedservice.ReportDeploymentPlannedRequest) (*pipedservice.ReportDeploymentPlannedResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}
	if err := a.validateDeploymentBelongsToPiped(ctx, req.DeploymentId, pipedID); err != nil {
		return nil, err
	}

	updater := datastore.DeploymentToPlannedUpdater(req.Summary, req.StatusReason, req.RunningCommitHash, req.Version, req.Stages)
	err = a.deploymentStore.UpdateDeployment(ctx, req.DeploymentId, updater)
	if err != nil {
		switch err {
		case datastore.ErrNotFound:
			return nil, status.Error(codes.InvalidArgument, "deployment is not found")
		case datastore.ErrInvalidArgument:
			return nil, status.Error(codes.InvalidArgument, "invalid value for update")
		default:
			a.logger.Error("failed to update deployment to be planned",
				zap.String("deployment-id", req.DeploymentId),
				zap.Error(err),
			)
			return nil, status.Error(codes.Internal, "failed to update deployment to be planned")
		}
	}
	return &pipedservice.ReportDeploymentPlannedResponse{}, nil
}

// ReportDeploymentStatusChanged is used to update the status
// of a specific deployment to RUNNING or ROLLING_BACK.
func (a *PipedAPI) ReportDeploymentStatusChanged(ctx context.Context, req *pipedservice.ReportDeploymentStatusChangedRequest) (*pipedservice.ReportDeploymentStatusChangedResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}
	if err := a.validateDeploymentBelongsToPiped(ctx, req.DeploymentId, pipedID); err != nil {
		return nil, err
	}

	updater := datastore.DeploymentStatusUpdater(req.Status, req.StatusReason)
	err = a.deploymentStore.UpdateDeployment(ctx, req.DeploymentId, updater)
	if err != nil {
		switch err {
		case datastore.ErrNotFound:
			return nil, status.Error(codes.InvalidArgument, "deployment is not found")
		case datastore.ErrInvalidArgument:
			return nil, status.Error(codes.InvalidArgument, "invalid value for update")
		default:
			a.logger.Error("failed to update deployment status",
				zap.String("deployment-id", req.DeploymentId),
				zap.Error(err),
			)
			return nil, status.Error(codes.Internal, "failed to update deployment status")
		}
	}
	return &pipedservice.ReportDeploymentStatusChangedResponse{}, nil
}

// ReportDeploymentCompleted used by piped to update the status
// of a specific deployment to SUCCESS | FAILURE | CANCELLED.
func (a *PipedAPI) ReportDeploymentCompleted(ctx context.Context, req *pipedservice.ReportDeploymentCompletedRequest) (*pipedservice.ReportDeploymentCompletedResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}
	if err := a.validateDeploymentBelongsToPiped(ctx, req.DeploymentId, pipedID); err != nil {
		return nil, err
	}

	updater := datastore.DeploymentToCompletedUpdater(req.Status, req.StageStatuses, req.StatusReason, req.CompletedAt)
	err = a.deploymentStore.UpdateDeployment(ctx, req.DeploymentId, updater)
	if err != nil {
		switch err {
		case datastore.ErrNotFound:
			return nil, status.Error(codes.InvalidArgument, "deployment is not found")
		case datastore.ErrInvalidArgument:
			return nil, status.Error(codes.InvalidArgument, "invalid value for update")
		default:
			a.logger.Error("failed to update deployment to be completed",
				zap.String("deployment-id", req.DeploymentId),
				zap.Error(err),
			)
			return nil, status.Error(codes.Internal, "failed to update deployment to be completed")
		}
	}
	return &pipedservice.ReportDeploymentCompletedResponse{}, nil
}

// SaveDeploymentMetadata used by piped to persist the metadata of a specific deployment.
func (a *PipedAPI) SaveDeploymentMetadata(ctx context.Context, req *pipedservice.SaveDeploymentMetadataRequest) (*pipedservice.SaveDeploymentMetadataResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}
	if err := a.validateDeploymentBelongsToPiped(ctx, req.DeploymentId, pipedID); err != nil {
		return nil, err
	}

	err = a.deploymentStore.PutDeploymentMetadata(ctx, req.DeploymentId, req.Metadata)
	if errors.Is(err, datastore.ErrNotFound) {
		return nil, status.Error(codes.InvalidArgument, "deployment is not found")
	}
	if err != nil {
		a.logger.Error("failed to save deployment metadata",
			zap.String("deployment-id", req.DeploymentId),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to save deployment metadata")
	}
	return &pipedservice.SaveDeploymentMetadataResponse{}, nil
}

// SaveStageMetadata used by piped to persist the metadata
// of a specific stage of a deployment.
func (a *PipedAPI) SaveStageMetadata(ctx context.Context, req *pipedservice.SaveStageMetadataRequest) (*pipedservice.SaveStageMetadataResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}
	if err := a.validateDeploymentBelongsToPiped(ctx, req.DeploymentId, pipedID); err != nil {
		return nil, err
	}

	err = a.deploymentStore.PutDeploymentStageMetadata(ctx, req.DeploymentId, req.StageId, req.Metadata)
	if err != nil {
		switch errors.Unwrap(err) {
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
			return nil, status.Error(codes.Internal, "failed to save deployment stage metadata")
		}
	}
	return &pipedservice.SaveStageMetadataResponse{}, nil
}

// ReportStageLogs is sent by piped to save the log of a pipeline stage.
func (a *PipedAPI) ReportStageLogs(ctx context.Context, req *pipedservice.ReportStageLogsRequest) (*pipedservice.ReportStageLogsResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}
	if err := a.validateDeploymentBelongsToPiped(ctx, req.DeploymentId, pipedID); err != nil {
		return nil, err
	}

	err = a.stageLogStore.AppendLogs(ctx, req.DeploymentId, req.StageId, req.RetriedCount, req.Blocks)
	if errors.Is(err, stagelogstore.ErrAlreadyCompleted) {
		return nil, status.Error(codes.FailedPrecondition, "could not append the logs because the stage was already completed")
	}
	if err != nil {
		a.logger.Error("failed to append logs", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to append logs")
	}
	return &pipedservice.ReportStageLogsResponse{}, nil
}

// ReportStageLogsFromLastCheckpoint is used to save the full logs from the most recently saved point.
func (a *PipedAPI) ReportStageLogsFromLastCheckpoint(ctx context.Context, req *pipedservice.ReportStageLogsFromLastCheckpointRequest) (*pipedservice.ReportStageLogsFromLastCheckpointResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}
	if err := a.validateDeploymentBelongsToPiped(ctx, req.DeploymentId, pipedID); err != nil {
		return nil, err
	}

	err = a.stageLogStore.AppendLogsFromLastCheckpoint(ctx, req.DeploymentId, req.StageId, req.RetriedCount, req.Blocks, req.Completed)
	if errors.Is(err, stagelogstore.ErrAlreadyCompleted) {
		return nil, status.Error(codes.FailedPrecondition, "could not append the logs because the stage was already completed")
	}
	if err != nil {
		a.logger.Error("failed to append logs", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to append logs")
	}
	return &pipedservice.ReportStageLogsFromLastCheckpointResponse{}, nil
}

// ReportStageStatusChanged used by piped to update the status
// of a specific stage of a deployment.
func (a *PipedAPI) ReportStageStatusChanged(ctx context.Context, req *pipedservice.ReportStageStatusChangedRequest) (*pipedservice.ReportStageStatusChangedResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}
	if err := a.validateDeploymentBelongsToPiped(ctx, req.DeploymentId, pipedID); err != nil {
		return nil, err
	}

	updater := datastore.StageStatusChangedUpdater(req.StageId, req.Status, req.StatusReason, req.Requires, req.Visible, req.RetriedCount, req.CompletedAt)
	err = a.deploymentStore.UpdateDeployment(ctx, req.DeploymentId, updater)
	if err != nil {
		switch err {
		case datastore.ErrNotFound:
			return nil, status.Error(codes.InvalidArgument, "deployment is not found")
		case datastore.ErrInvalidArgument:
			return nil, status.Error(codes.InvalidArgument, "invalid value for update")
		default:
			a.logger.Error("failed to update stage status",
				zap.String("deployment-id", req.DeploymentId),
				zap.String("stage-id", req.StageId),
				zap.Error(err),
			)
			return nil, status.Error(codes.Internal, "failed to update stage status")
		}
	}
	return &pipedservice.ReportStageStatusChangedResponse{}, nil
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
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}

	cmds, err := a.commandStore.ListUnhandledCommands(ctx, pipedID)
	if err != nil {
		a.logger.Error("failed to fetch unhandled commands", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to unhandled commands")
	}
	return &pipedservice.ListUnhandledCommandsResponse{
		Commands: cmds,
	}, nil
}

// ReportCommandHandled is called by piped to mark a specific command as handled.
// The request payload will contain the handle status as well as any additional result data.
// The handle result should be updated to both datastore and cache (for reading from web).
func (a *PipedAPI) ReportCommandHandled(ctx context.Context, req *pipedservice.ReportCommandHandledRequest) (*pipedservice.ReportCommandHandledResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}

	cmd, err := a.getCommand(ctx, req.CommandId)
	if err != nil {
		return nil, err
	}
	if pipedID != cmd.PipedId {
		return nil, status.Error(codes.PermissionDenied, "The current piped does not have requested command")
	}

	if len(req.Output) > 0 {
		if err := a.commandOutputPutter.Put(ctx, req.CommandId, req.Output); err != nil {
			a.logger.Error("failed to store output of command",
				zap.String("command_id", req.CommandId),
				zap.Error(err),
			)
			return nil, status.Error(codes.Internal, "Failed to store output of command")
		}
	}

	err = a.commandStore.UpdateCommandHandled(ctx, req.CommandId, req.Status, req.Metadata, req.HandledAt)
	if err != nil {
		switch err {
		case datastore.ErrNotFound:
			return nil, status.Error(codes.NotFound, "command is not found")
		case datastore.ErrInvalidArgument:
			return nil, status.Error(codes.InvalidArgument, "invalid value for update")
		default:
			a.logger.Error("failed to update command",
				zap.String("command-id", req.CommandId),
				zap.Error(err),
			)
			return nil, status.Error(codes.Internal, "failed to update command")
		}
	}
	return &pipedservice.ReportCommandHandledResponse{}, nil
}

func (a *PipedAPI) getCommand(ctx context.Context, pipedID string) (*model.Command, error) {
	cmd, err := a.commandStore.GetCommand(ctx, pipedID)
	if errors.Is(err, datastore.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "command is not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get command")
	}
	return cmd, nil
}

// ReportApplicationLiveState is periodically sent to correct full state of an application.
// For kubernetes application, this contains a full tree of its kubernetes resources.
// The tree data should be written into filestore immediately and then the state in cache should be refreshsed too.
func (a *PipedAPI) ReportApplicationLiveState(ctx context.Context, req *pipedservice.ReportApplicationLiveStateRequest) (*pipedservice.ReportApplicationLiveStateResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}
	if err := a.validateAppBelongsToPiped(ctx, req.Snapshot.ApplicationId, pipedID); err != nil {
		return nil, err
	}

	if err := a.applicationLiveStateStore.PutStateSnapshot(ctx, req.Snapshot); err != nil {
		return nil, status.Error(codes.Internal, "failed to report application live state")
	}
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
func (a *PipedAPI) ReportApplicationLiveStateEvents(ctx context.Context, req *pipedservice.ReportApplicationLiveStateEventsRequest) (*pipedservice.ReportApplicationLiveStateEventsResponse, error) {
	a.applicationLiveStateStore.PatchKubernetesApplicationLiveState(ctx, req.KubernetesEvents)
	// TODO: Patch Terraform application live state
	// TODO: Patch Cloud Run application live state
	// TODO: Patch Lambda application live state
	return &pipedservice.ReportApplicationLiveStateEventsResponse{}, nil
}

// GetLatestEvent returns the latest event that meets the given conditions.
func (a *PipedAPI) GetLatestEvent(ctx context.Context, req *pipedservice.GetLatestEventRequest) (*pipedservice.GetLatestEventResponse, error) {
	projectID, _, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}

	// Try to fetch the most recently registered event that has the given parameters.
	opts := datastore.ListOptions{
		Limit: 1,
		Filters: []datastore.ListFilter{
			{
				Field:    "ProjectId",
				Operator: datastore.OperatorEqual,
				Value:    projectID,
			},
			{
				Field:    "Name",
				Operator: datastore.OperatorEqual,
				Value:    req.Name,
			},
			{
				Field:    "EventKey",
				Operator: datastore.OperatorEqual,
				Value:    model.MakeEventKey(req.Name, req.Labels),
			},
		},
		Orders: []datastore.Order{
			{
				Field:     "CreatedAt",
				Direction: datastore.Desc,
			},
			{
				Field:     "Id",
				Direction: datastore.Asc,
			},
		},
	}
	events, err := a.eventStore.ListEvents(ctx, opts)
	if err != nil {
		a.logger.Error("failed to list events", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list event")
	}
	if len(events) == 0 {
		return nil, status.Error(codes.NotFound, "no events found")
	}

	return &pipedservice.GetLatestEventResponse{
		Event: events[0],
	}, nil
}

// ListEvents returns a list of Events inside the given range.
func (a *PipedAPI) ListEvents(ctx context.Context, req *pipedservice.ListEventsRequest) (*pipedservice.ListEventsResponse, error) {
	projectID, _, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}

	// Build options based on the request.
	opts := datastore.ListOptions{
		Filters: []datastore.ListFilter{
			{
				Field:    "ProjectId",
				Operator: datastore.OperatorEqual,
				Value:    projectID,
			},
		},
	}
	if req.From > 0 {
		opts.Filters = append(opts.Filters, datastore.ListFilter{
			Field:    "CreatedAt",
			Operator: datastore.OperatorGreaterThanOrEqual,
			Value:    req.From,
		})
	}
	if req.To > 0 {
		opts.Filters = append(opts.Filters, datastore.ListFilter{
			Field:    "CreatedAt",
			Operator: datastore.OperatorLessThan,
			Value:    req.To,
		})
	}
	switch req.Order {
	case pipedservice.ListOrder_ASC:
		opts.Orders = []datastore.Order{
			{
				Field:     "CreatedAt",
				Direction: datastore.Asc,
			},
			{
				Field:     "Id",
				Direction: datastore.Asc,
			},
		}
	case pipedservice.ListOrder_DESC:
		opts.Orders = []datastore.Order{
			{
				Field:     "CreatedAt",
				Direction: datastore.Desc,
			},
			{
				Field:     "Id",
				Direction: datastore.Asc,
			},
		}
	}

	events, err := a.eventStore.ListEvents(ctx, opts)
	if err != nil {
		a.logger.Error("failed to list events", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list events")
	}
	return &pipedservice.ListEventsResponse{
		Events: events,
	}, nil
}

func (a *PipedAPI) ReportEventsHandled(ctx context.Context, req *pipedservice.ReportEventsHandledRequest) (*pipedservice.ReportEventsHandledResponse, error) {
	_, _, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}

	for _, id := range req.EventIds {
		if err := a.eventStore.MarkEventHandled(ctx, id); err != nil {
			switch err {
			case datastore.ErrNotFound:
				return nil, status.Errorf(codes.NotFound, "event %q is not found", id)
			default:
				a.logger.Error("failed to mark event as handled",
					zap.String("event-id", id),
					zap.Error(err),
				)
				return nil, status.Errorf(codes.Internal, "failed to mark event %q as handled", id)
			}
		}
	}
	return &pipedservice.ReportEventsHandledResponse{}, nil
}

func (a *PipedAPI) GetLatestAnalysisResult(ctx context.Context, req *pipedservice.GetLatestAnalysisResultRequest) (*pipedservice.GetLatestAnalysisResultResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}
	if err := a.validateAppBelongsToPiped(ctx, req.ApplicationId, pipedID); err != nil {
		return nil, err
	}

	result, err := a.analysisResultStore.GetLatestAnalysisResult(ctx, req.ApplicationId)
	if errors.Is(err, filestore.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "the most recent analysis result is not found")
	}
	if err != nil {
		a.logger.Error("failed to get the most recent analysis result", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get the most recent analysis result")
	}
	return &pipedservice.GetLatestAnalysisResultResponse{
		AnalysisResult: result,
	}, nil
}

func (a *PipedAPI) PutLatestAnalysisResult(ctx context.Context, req *pipedservice.PutLatestAnalysisResultRequest) (*pipedservice.PutLatestAnalysisResultResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}
	if err := a.validateAppBelongsToPiped(ctx, req.ApplicationId, pipedID); err != nil {
		return nil, err
	}

	err = a.analysisResultStore.PutLatestAnalysisResult(ctx, req.ApplicationId, req.AnalysisResult)
	if err != nil {
		a.logger.Error("failed to put the most recent analysis result", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to put the most recent analysis result")
	}
	return &pipedservice.PutLatestAnalysisResultResponse{}, nil
}

func (a *PipedAPI) GetDesiredVersion(ctx context.Context, _ *pipedservice.GetDesiredVersionRequest) (*pipedservice.GetDesiredVersionResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}
	piped, err := getPiped(ctx, a.pipedStore, pipedID, a.logger)
	if err != nil {
		return nil, err
	}
	return &pipedservice.GetDesiredVersionResponse{
		Version: piped.DesiredVersion,
	}, nil
}

func (a *PipedAPI) UpdateApplicationConfigurations(ctx context.Context, req *pipedservice.UpdateApplicationConfigurationsRequest) (*pipedservice.UpdateApplicationConfigurationsResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}
	// Scan all of them to guarantee in advance that there is no invalid request.
	for _, appInfo := range req.Applications {
		if err := a.validateAppBelongsToPiped(ctx, appInfo.Id, pipedID); err != nil {
			return nil, err
		}
	}
	// TODO: Consider bulk-updating multiple apps
	for _, appInfo := range req.Applications {
		updater := func(app *model.Application) error {
			app.Name = appInfo.Name
			app.Kind = appInfo.Kind
			app.Labels = appInfo.Labels
			return nil
		}
		if err := a.applicationStore.UpdateApplication(ctx, appInfo.Id, updater); err != nil {
			a.logger.Error("failed to update application", zap.Error(err))
			return nil, status.Error(codes.Internal, "failed to update application")
		}
	}
	return &pipedservice.UpdateApplicationConfigurationsResponse{}, nil
}

func (a *PipedAPI) ReportUnregisteredApplicationConfigurations(ctx context.Context, req *pipedservice.ReportUnregisteredApplicationConfigurationsRequest) (*pipedservice.ReportUnregisteredApplicationConfigurationsResponse, error) {
	projectID, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}

	key := makeUnregisteredAppsCacheKey(projectID)
	c := rediscache.NewHashCache(a.redis, key)
	// Cache a slice of *model.ApplicationInfo.
	if err := c.Put(pipedID, req.Applications); err != nil {
		return nil, status.Error(codes.Internal, "failed to put the unregistered apps to the cache")
	}

	return &pipedservice.ReportUnregisteredApplicationConfigurationsResponse{}, nil
}

// CreateDeploymentChain creates a new deployment chain object and all required commands to
// trigger deployment for applications in the chain.
func (a *PipedAPI) CreateDeploymentChain(ctx context.Context, req *pipedservice.CreateDeploymentChainRequest) (*pipedservice.CreateDeploymentChainResponse, error) {
	firstDeployment := req.FirstDeployment
	projectID, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}
	if err := a.validateAppBelongsToPiped(ctx, firstDeployment.ApplicationId, pipedID); err != nil {
		return nil, err
	}

	buildChainNodes := func(matcher *pipedservice.CreateDeploymentChainRequest_ApplicationMatcher) ([]*model.ChainNode, []*model.Application, error) {
		filters := []datastore.ListFilter{
			{
				Field:    "ProjectId",
				Operator: datastore.OperatorEqual,
				Value:    projectID,
			},
		}

		if matcher.Name != "" {
			filters = append(filters, datastore.ListFilter{
				Field:    "Name",
				Operator: datastore.OperatorEqual,
				Value:    matcher.Name,
			})
		}

		// TODO: Support find node apps by appKind and appLabels.

		apps, _, err := a.applicationStore.ListApplications(ctx, datastore.ListOptions{
			Filters: filters,
		})
		if err != nil {
			return nil, nil, err
		}

		nodes := make([]*model.ChainNode, 0, len(apps))
		for _, app := range apps {
			nodes = append(nodes, &model.ChainNode{
				ApplicationRef: &model.ChainApplicationRef{
					ApplicationId:   app.Id,
					ApplicationName: app.Name,
				},
			})
		}

		return nodes, apps, nil
	}

	chainBlocks := make([]*model.ChainBlock, 0, len(req.Matchers)+1)
	// Add the first deployment which created by piped as the first block of the chain.
	chainBlocks = append(chainBlocks, &model.ChainBlock{
		Nodes: []*model.ChainNode{
			{
				ApplicationRef: &model.ChainApplicationRef{
					ApplicationId:   firstDeployment.ApplicationId,
					ApplicationName: firstDeployment.ApplicationName,
				},
				DeploymentRef: &model.ChainDeploymentRef{
					DeploymentId: firstDeployment.Id,
				},
			},
		},
		StartedAt: time.Now().Unix(),
	})

	blockAppsMap := make(map[int][]*model.Application, len(req.Matchers))
	for i, filter := range req.Matchers {
		nodes, blockApps, err := buildChainNodes(filter)
		if err != nil {
			return nil, err
		}

		blockAppsMap[i+1] = blockApps
		chainBlocks = append(chainBlocks, &model.ChainBlock{
			Nodes:     nodes,
			StartedAt: time.Now().Unix(),
		})
	}

	dc := model.DeploymentChain{
		Id:        uuid.New().String(),
		ProjectId: projectID,
		Blocks:    chainBlocks,
	}

	// Create a new deployment chain instance to control newly triggered deployment chain.
	if err := a.deploymentChainStore.AddDeploymentChain(ctx, &dc); err != nil {
		a.logger.Error("failed to create deployment chain", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to trigger new deployment chain")
	}

	firstDeployment.DeploymentChainId = dc.Id
	// Trigger new deployment for the first application by store first deployment to datastore.
	if err := a.deploymentStore.AddDeployment(ctx, firstDeployment); err != nil {
		a.logger.Error("failed to create deployment", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to trigger new deployment for the first application in chain")
	}

	// Make sync application command for applications of the chain.
	for blockIndex, apps := range blockAppsMap {
		for _, app := range apps {
			cmd := model.Command{
				Id:            uuid.New().String(),
				PipedId:       app.PipedId,
				ApplicationId: app.Id,
				ProjectId:     app.ProjectId,
				Commander:     dc.Id,
				Type:          model.Command_CHAIN_SYNC_APPLICATION,
				ChainSyncApplication: &model.Command_ChainSyncApplication{
					DeploymentChainId: dc.Id,
					BlockIndex:        uint32(blockIndex),
					ApplicationId:     app.Id,
					SyncStrategy:      model.SyncStrategy_AUTO,
				},
			}

			if err := addCommand(ctx, a.commandStore, &cmd, a.logger); err != nil {
				a.logger.Error("failed to create command to trigger application in chain", zap.Error(err))
				return nil, status.Error(codes.Internal, "failed to command to trigger for applications in chain")
			}
		}
	}

	return &pipedservice.CreateDeploymentChainResponse{}, nil
}

// InChainDeploymentPlannable hecks the completion and status of the previous block in the deployment chain.
// An in chain deployment is treated as plannable in case:
// - It's the first deployment of its deployment chain.
// - All deployments of its previous block in chain are at DEPLOYMENT_SUCCESS state.
func (a *PipedAPI) InChainDeploymentPlannable(ctx context.Context, req *pipedservice.InChainDeploymentPlannableRequest) (*pipedservice.InChainDeploymentPlannableResponse, error) {
	_, pipedID, _, err := rpcauth.ExtractPipedToken(ctx)
	if err != nil {
		return nil, err
	}
	if err := a.validateDeploymentBelongsToPiped(ctx, req.Deployment.Id, pipedID); err != nil {
		return nil, err
	}

	dc, err := a.deploymentChainStore.GetDeploymentChain(ctx, req.Deployment.DeploymentChainId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "unable to find the deployment chain which this deployment belongs to")
	}

	// Deployment of blocks[0] in the chain means it's the first deployment of the chain;
	// hence it should be processed without any lock.
	if req.Deployment.DeploymentChainBlockIndex == 0 {
		return &pipedservice.InChainDeploymentPlannableResponse{
			Plannable: true,
		}, nil
	}

	if req.Deployment.DeploymentChainBlockIndex >= uint32(len(dc.Blocks)) {
		return nil, status.Error(codes.InvalidArgument, "invalid deployment with chain block index provided")
	}

	prevBlock := dc.Blocks[req.Deployment.DeploymentChainBlockIndex-1]
	plannable := true
	for _, node := range prevBlock.Nodes {
		// TODO: Consider add deployment status to the deployment ref in the deployment chain model
		// instead of fetching deployment model here.
		dp, err := a.deploymentStore.GetDeployment(ctx, node.DeploymentRef.DeploymentId)
		if err != nil {
			return nil, status.Error(codes.Internal, "unable to process previous block nodes in deployment chain")
		}
		if model.IsSuccessfullyCompletedDeployment(dp.Status) {
			plannable = false
			break
		}
	}

	return &pipedservice.InChainDeploymentPlannableResponse{
		Plannable: plannable,
	}, nil
}

// validateAppBelongsToPiped checks if the given application belongs to the given piped.
// It gives back an error unless the application belongs to the piped.
func (a *PipedAPI) validateAppBelongsToPiped(ctx context.Context, appID, pipedID string) error {
	pid, err := a.appPipedCache.Get(appID)
	if err == nil {
		if pid != pipedID {
			return status.Error(codes.PermissionDenied, "requested application doesn't belong to the piped")
		}
		return nil
	}

	app, err := a.applicationStore.GetApplication(ctx, appID)
	if errors.Is(err, datastore.ErrNotFound) {
		return status.Error(codes.NotFound, "the application is not found")
	}
	if err != nil {
		a.logger.Error("failed to get application", zap.Error(err))
		return status.Error(codes.Internal, "failed to get application")
	}
	a.appPipedCache.Put(appID, app.PipedId)

	if app.PipedId != pipedID {
		return status.Error(codes.PermissionDenied, "requested application doesn't belong to the piped")
	}
	return nil
}

// validateDeploymentBelongsToPiped checks if the given deployment belongs to the given piped.
// It gives back an error unless the deployment belongs to the piped.
func (a *PipedAPI) validateDeploymentBelongsToPiped(ctx context.Context, deploymentID, pipedID string) error {
	pid, err := a.deploymentPipedCache.Get(deploymentID)
	if err == nil {
		if pid != pipedID {
			return status.Error(codes.PermissionDenied, "requested deployment doesn't belong to the piped")
		}
		return nil
	}

	deployment, err := a.deploymentStore.GetDeployment(ctx, deploymentID)
	if errors.Is(err, datastore.ErrNotFound) {
		return status.Error(codes.NotFound, "the deployment is not found")
	}
	if err != nil {
		a.logger.Error("failed to get deployment", zap.Error(err))
		return status.Error(codes.Internal, "failed to get deployment")
	}
	a.deploymentPipedCache.Put(deploymentID, deployment.PipedId)

	if deployment.PipedId != pipedID {
		return status.Error(codes.PermissionDenied, "requested deployment doesn't belong to the piped")
	}
	return nil
}

// validateEnvBelongsToProject checks if the given environment belongs to the given project.
// It gives back an error unless the environment belongs to the project.
func (a *PipedAPI) validateEnvBelongsToProject(ctx context.Context, envID, projectID string) error {
	pid, err := a.envProjectCache.Get(envID)
	if err == nil {
		if pid != projectID {
			return status.Error(codes.PermissionDenied, "requested environment doesn't belong to the project")
		}
		return nil
	}

	env, err := a.environmentStore.GetEnvironment(ctx, envID)
	if errors.Is(err, datastore.ErrNotFound) {
		return status.Error(codes.NotFound, "the environment is not found")
	}
	if err != nil {
		a.logger.Error("failed to get environment", zap.Error(err))
		return status.Error(codes.Internal, "failed to get environment")
	}
	a.envProjectCache.Put(envID, env.ProjectId)

	if env.ProjectId != projectID {
		return status.Error(codes.PermissionDenied, "requested environment doesn't belong to the project")
	}
	return nil
}
