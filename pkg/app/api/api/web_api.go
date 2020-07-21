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
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipe/pkg/app/api/applicationlivestatestore"
	"github.com/pipe-cd/pipe/pkg/app/api/commandstore"
	"github.com/pipe-cd/pipe/pkg/app/api/service/webservice"
	"github.com/pipe-cd/pipe/pkg/app/api/stagelogstore"
	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/model"
	"github.com/pipe-cd/pipe/pkg/rpc/rpcauth"
)

// PipedAPI implements the behaviors for the gRPC definitions of WebAPI.
type WebAPI struct {
	applicationStore          datastore.ApplicationStore
	environmentStore          datastore.EnvironmentStore
	deploymentStore           datastore.DeploymentStore
	pipedStore                datastore.PipedStore
	stageLogStore             stagelogstore.Store
	applicationLiveStateStore applicationlivestatestore.Store
	commandStore              commandstore.Store

	logger *zap.Logger
}

// NewWebAPI creates a new WebAPI instance.
func NewWebAPI(ds datastore.DataStore, sls stagelogstore.Store, alss applicationlivestatestore.Store, cmds commandstore.Store, logger *zap.Logger) *WebAPI {
	a := &WebAPI{
		applicationStore:          datastore.NewApplicationStore(ds),
		environmentStore:          datastore.NewEnvironmentStore(ds),
		deploymentStore:           datastore.NewDeploymentStore(ds),
		pipedStore:                datastore.NewPipedStore(ds),
		stageLogStore:             sls,
		applicationLiveStateStore: alss,
		commandStore:              cmds,
		logger:                    logger.Named("web-api"),
	}
	return a
}

// Register registers all handling of this service into the specified gRPC server.
func (a *WebAPI) Register(server *grpc.Server) {
	webservice.RegisterWebServiceServer(server, a)
}

func (a *WebAPI) AddEnvironment(ctx context.Context, req *webservice.AddEnvironmentRequest) (*webservice.AddEnvironmentResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		return nil, err
	}
	env := model.Environment{
		Id:        uuid.New().String(),
		Name:      req.Name,
		Desc:      req.Desc,
		ProjectId: claims.Role.ProjectId,
	}
	err = a.environmentStore.AddEnvironment(ctx, &env)
	if errors.Is(err, datastore.ErrAlreadyExists) {
		return nil, status.Error(codes.AlreadyExists, "environment already exists")
	}
	if err != nil {
		a.logger.Error("failed to create environment", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create environment")
	}
	return &webservice.AddEnvironmentResponse{}, nil
}

func (a *WebAPI) UpdateEnvironmentDesc(ctx context.Context, req *webservice.UpdateEnvironmentDescRequest) (*webservice.UpdateEnvironmentDescResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *WebAPI) ListEnvironments(ctx context.Context, req *webservice.ListEnvironmentsRequest) (*webservice.ListEnvironmentsResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		return nil, err
	}
	opts := datastore.ListOptions{
		Filters: []datastore.ListFilter{
			{
				Field:    "ProjectId",
				Operator: "==",
				Value:    claims.Role.ProjectId,
			},
		},
	}
	envs, err := a.environmentStore.ListEnvironments(ctx, opts)
	if err != nil {
		a.logger.Error("failed to get environments", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get environments")
	}

	return &webservice.ListEnvironmentsResponse{
		Environments: envs,
	}, nil
}

func (a *WebAPI) RegisterPiped(ctx context.Context, req *webservice.RegisterPipedRequest) (*webservice.RegisterPipedResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		return nil, err
	}

	key, keyHash, err := model.GeneratePipedKey()
	if err != nil {
		a.logger.Error("failed to generate piped key", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to generate the piped key")
	}
	id := uuid.New().String()
	piped := model.Piped{
		Id:        id,
		Name:      req.Name,
		Desc:      req.Desc,
		KeyHash:   keyHash,
		ProjectId: claims.Role.ProjectId,
	}
	err = a.pipedStore.AddPiped(ctx, &piped)
	if errors.Is(err, datastore.ErrAlreadyExists) {
		return nil, status.Error(codes.AlreadyExists, "piped already exists")
	}
	if err != nil {
		a.logger.Error("failed to register piped", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to register piped")
	}
	return &webservice.RegisterPipedResponse{
		Id:  id,
		Key: key,
	}, nil
}

func (a *WebAPI) RecreatePipedKey(ctx context.Context, req *webservice.RecreatePipedKeyRequest) (*webservice.RecreatePipedKeyResponse, error) {
	key, keyHash, err := model.GeneratePipedKey()
	if err != nil {
		a.logger.Error("failed to generate piped key", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to generate the piped key")
	}

	updater := func(ctx context.Context, pipedID string) error {
		return a.pipedStore.UpdateKeyHash(ctx, pipedID, keyHash)
	}
	if err := a.updatePiped(ctx, req.Id, updater); err != nil {
		return nil, err
	}

	return &webservice.RecreatePipedKeyResponse{
		Key: key,
	}, nil
}

func (a *WebAPI) EnablePiped(ctx context.Context, req *webservice.EnablePipedRequest) (*webservice.EnablePipedResponse, error) {
	if err := a.updatePiped(ctx, req.PipedId, a.pipedStore.EnablePiped); err != nil {
		return nil, err
	}
	return &webservice.EnablePipedResponse{}, nil
}

func (a *WebAPI) DisablePiped(ctx context.Context, req *webservice.DisablePipedRequest) (*webservice.DisablePipedResponse, error) {
	if err := a.updatePiped(ctx, req.PipedId, a.pipedStore.DisablePiped); err != nil {
		return nil, err
	}
	return &webservice.DisablePipedResponse{}, nil
}

func (a *WebAPI) updatePiped(ctx context.Context, pipedID string, updater func(context.Context, string) error) error {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		return err
	}

	piped, err := a.getPiped(ctx, pipedID)
	if err != nil {
		return err
	}

	if claims.Role.ProjectId != piped.ProjectId {
		return status.Error(codes.PermissionDenied, "The current project does not have requested piped")
	}

	if err := updater(ctx, pipedID); err != nil {
		switch err {
		case datastore.ErrNotFound:
			return status.Error(codes.InvalidArgument, "piped is not found")
		case datastore.ErrInvalidArgument:
			return status.Error(codes.InvalidArgument, "invalid value for update")
		default:
			a.logger.Error("failed to update the piped",
				zap.String("piped-id", pipedID),
				zap.Error(err),
			)
			return status.Error(codes.Internal, "failed to update the piped ")
		}
	}
	return nil
}

func (a *WebAPI) ListPipeds(ctx context.Context, req *webservice.ListPipedsRequest) (*webservice.ListPipedsResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		return nil, err
	}
	opts := datastore.ListOptions{
		Filters: []datastore.ListFilter{
			{
				Field:    "ProjectId",
				Operator: "==",
				Value:    claims.Role.ProjectId,
			},
		},
	}

	if req.Options != nil {
		if req.Options.Enabled != nil {
			opts.Filters = append(opts.Filters, datastore.ListFilter{
				Field:    "Disabled",
				Operator: "==",
				Value:    !req.Options.Enabled.GetValue(),
			})
		}
	}
	pipeds, err := a.pipedStore.ListPipeds(ctx, opts)
	if err != nil {
		a.logger.Error("failed to get pipeds", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get pipeds")
	}
	// Returning a response that does not contain sensitive data like KeyHash is more safer.
	webPipeds := make([]*webservice.Piped, len(pipeds))
	for i := range pipeds {
		webPipeds[i] = webservice.MakePiped(pipeds[i])
		if req.WithStatus {
			// TODO: While reducing the number of query and get ping connection status from piped_stats
			webPipeds[i].Status = webservice.PipedConnectionStatus_PIPED_CONNECTION_ONLINE
		}
	}
	return &webservice.ListPipedsResponse{
		Pipeds: webPipeds,
	}, nil
}

func (a *WebAPI) GetPiped(ctx context.Context, req *webservice.GetPipedRequest) (*webservice.GetPipedResponse, error) {
	piped, err := a.getPiped(ctx, req.PipedId)
	if err != nil {
		return nil, err
	}
	return &webservice.GetPipedResponse{
		Piped: webservice.MakePiped(piped),
	}, nil
}

func (a *WebAPI) getPiped(ctx context.Context, pipedID string) (*model.Piped, error) {
	piped, err := a.pipedStore.GetPiped(ctx, pipedID)
	if errors.Is(err, datastore.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "piped is not found")
	}
	if err != nil {
		a.logger.Error("failed to get piped", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get piped")
	}
	return piped, nil
}

func (a *WebAPI) AddApplication(ctx context.Context, req *webservice.AddApplicationRequest) (*webservice.AddApplicationResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		return nil, err
	}

	app := model.Application{
		Id:            uuid.New().String(),
		Name:          req.Name,
		EnvId:         req.EnvId,
		PipedId:       req.PipedId,
		ProjectId:     claims.Role.ProjectId,
		GitPath:       req.GitPath,
		Kind:          req.Kind,
		CloudProvider: req.CloudProvider,
	}
	err = a.applicationStore.AddApplication(ctx, &app)
	if errors.Is(err, datastore.ErrAlreadyExists) {
		return nil, status.Error(codes.AlreadyExists, "application already exists")
	}
	if err != nil {
		a.logger.Error("failed to create application", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create application")
	}

	return &webservice.AddApplicationResponse{}, nil
}

func (a *WebAPI) EnableApplication(ctx context.Context, req *webservice.EnableApplicationRequest) (*webservice.EnableApplicationResponse, error) {
	if err := a.updateApplicationEnable(ctx, req.ApplicationId, true); err != nil {
		return nil, err
	}
	return &webservice.EnableApplicationResponse{}, nil
}

func (a *WebAPI) DisableApplication(ctx context.Context, req *webservice.DisableApplicationRequest) (*webservice.DisableApplicationResponse, error) {
	if err := a.updateApplicationEnable(ctx, req.ApplicationId, false); err != nil {
		return nil, err
	}
	return &webservice.DisableApplicationResponse{}, nil
}

func (a *WebAPI) updateApplicationEnable(ctx context.Context, appID string, enable bool) error {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		return err
	}
	app, err := a.getApplication(ctx, appID)
	if err != nil {
		return err
	}
	if app.ProjectId != claims.Role.ProjectId {
		return status.Error(codes.PermissionDenied, "The current project does not have requested application")
	}

	var updater func(context.Context, string) error
	if enable {
		updater = a.applicationStore.EnableApplication
	} else {
		updater = a.applicationStore.DisableApplication
	}

	if err := updater(ctx, appID); err != nil {
		switch err {
		case datastore.ErrNotFound:
			return status.Error(codes.InvalidArgument, "application is not found")
		case datastore.ErrInvalidArgument:
			return status.Error(codes.InvalidArgument, "invalid value for update")
		default:
			a.logger.Error("failed to update the application",
				zap.String("application-id", appID),
				zap.Error(err),
			)
			return status.Error(codes.Internal, "failed to update the piped ")
		}
	}
	return nil
}

func (a *WebAPI) ListApplications(ctx context.Context, req *webservice.ListApplicationsRequest) (*webservice.ListApplicationsResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		return nil, err
	}
	filters := []datastore.ListFilter{
		{
			Field:    "ProjectId",
			Operator: "==",
			Value:    claims.Role.ProjectId,
		},
	}
	if o := req.Options; o != nil {
		if o.Enabled != nil {
			filters = append(filters, datastore.ListFilter{
				Field:    "Disabled",
				Operator: "==",
				Value:    !o.Enabled.GetValue(),
			})
		}
		// TODO: Support multiple filtering of Kinds, SyncStatuses and EnvIds
		//   Because some datastores do not support IN query filter,
		//   currently only the first value is used.
		if len(o.Kinds) > 0 {
			filters = append(filters, datastore.ListFilter{
				Field:    "Kind",
				Operator: "==",
				Value:    o.Kinds[0],
			})
		}
		if len(o.SyncStatuses) > 0 {
			filters = append(filters, datastore.ListFilter{
				Field:    "SyncState.Status",
				Operator: "==",
				Value:    o.SyncStatuses[0],
			})
		}
		if len(o.EnvIds) > 0 {
			filters = append(filters, datastore.ListFilter{
				Field:    "EnvId",
				Operator: "==",
				Value:    o.EnvIds[0],
			})
		}
	}

	apps, err := a.applicationStore.ListApplications(ctx, datastore.ListOptions{
		Filters: filters,
	})
	if err != nil {
		a.logger.Error("failed to get applications", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get applications")
	}

	return &webservice.ListApplicationsResponse{
		Applications: apps,
	}, nil
}

func (a *WebAPI) SyncApplication(ctx context.Context, req *webservice.SyncApplicationRequest) (*webservice.SyncApplicationResponse, error) {
	app, err := a.getApplication(ctx, req.ApplicationId)
	if err != nil {
		return nil, err
	}

	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		return nil, err
	}
	if app.ProjectId != claims.Role.ProjectId {
		return nil, status.Error(codes.PermissionDenied, "The current project does not have requested application")
	}

	commandID := uuid.New().String()
	cmd := model.Command{
		Id:            commandID,
		PipedId:       app.PipedId,
		ApplicationId: app.Id,
		Type:          model.Command_SYNC_APPLICATION,
		Commander:     "anonymous", // TODO: Getting value from login user.
		SyncApplication: &model.Command_SyncApplication{
			ApplicationId: req.ApplicationId,
		},
	}
	if err := a.addCommand(ctx, &cmd); err != nil {
		return nil, err
	}
	return &webservice.SyncApplicationResponse{
		CommandId: commandID,
	}, nil
}

func (a *WebAPI) addCommand(ctx context.Context, cmd *model.Command) error {
	if err := a.commandStore.AddCommand(ctx, cmd); err != nil {
		a.logger.Error("failed to create command", zap.Error(err))
		return status.Error(codes.Internal, "failed to create command")
	}
	return nil
}

func (a *WebAPI) GetApplication(ctx context.Context, req *webservice.GetApplicationRequest) (*webservice.GetApplicationResponse, error) {
	app, err := a.getApplication(ctx, req.ApplicationId)
	if err != nil {
		return nil, err
	}
	return &webservice.GetApplicationResponse{
		Application: app,
	}, nil
}

func (a *WebAPI) getApplication(ctx context.Context, id string) (*model.Application, error) {
	app, err := a.applicationStore.GetApplication(ctx, id)
	if errors.Is(err, datastore.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "application is not found")
	}
	if err != nil {
		a.logger.Error("failed to get application", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get application")
	}
	return app, nil
}

func (a *WebAPI) ListDeployments(ctx context.Context, req *webservice.ListDeploymentsRequest) (*webservice.ListDeploymentsResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		return nil, err
	}
	// TODO: Support pagination for Deployment list
	filters := []datastore.ListFilter{
		{
			Field:    "ProjectId",
			Operator: "==",
			Value:    claims.Role.ProjectId,
		},
	}
	if o := req.Options; o != nil {
		// TODO: Support IN Query filter for Deployments fields
		//   Because some datastores do not support IN query filter,
		//   currently only the first value is used.
		if len(o.Statuses) > 0 {
			filters = append(filters, datastore.ListFilter{
				Field:    "Status",
				Operator: "==",
				Value:    o.Statuses[0],
			})
		}
		if len(o.Kinds) > 0 {
			filters = append(filters, datastore.ListFilter{
				Field:    "Kind",
				Operator: "==",
				Value:    o.Kinds[0],
			})
		}
		if len(o.ApplicationIds) > 0 {
			filters = append(filters, datastore.ListFilter{
				Field:    "ApplicationId",
				Operator: "==",
				Value:    o.ApplicationIds[0],
			})
		}
		if len(o.EnvIds) > 0 {
			filters = append(filters, datastore.ListFilter{
				Field:    "EnvId",
				Operator: "==",
				Value:    o.EnvIds[0],
			})
		}
	}

	deployments, err := a.deploymentStore.ListDeployments(ctx, datastore.ListOptions{
		Filters: filters,
	})
	if err != nil {
		a.logger.Error("failed to get deployments", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get deployments")
	}
	return &webservice.ListDeploymentsResponse{
		Deployments: deployments,
	}, nil
}

func (a *WebAPI) GetDeployment(ctx context.Context, req *webservice.GetDeploymentRequest) (*webservice.GetDeploymentResponse, error) {
	deployment, err := a.getDeployment(ctx, req.DeploymentId)
	if err != nil {
		return nil, err
	}
	return &webservice.GetDeploymentResponse{
		Deployment: deployment,
	}, nil
}

func (a *WebAPI) getDeployment(ctx context.Context, id string) (*model.Deployment, error) {
	deployment, err := a.deploymentStore.GetDeployment(ctx, id)
	if errors.Is(err, datastore.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "deployment is not found")
	}
	if err != nil {
		a.logger.Error("failed to get deployment", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get deployment")
	}
	return deployment, nil
}

func (a *WebAPI) GetStageLog(ctx context.Context, req *webservice.GetStageLogRequest) (*webservice.GetStageLogResponse, error) {
	blocks, completed, err := a.stageLogStore.FetchLogs(ctx, req.DeploymentId, req.StageId, req.RetriedCount, req.OffsetIndex)
	if errors.Is(err, stagelogstore.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "stage log not found")
	}
	if err != nil {
		a.logger.Error("failed to get stage logs", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get stage logs")
	}

	return &webservice.GetStageLogResponse{
		Blocks:    blocks,
		Completed: completed,
	}, nil
}

func (a *WebAPI) CancelDeployment(ctx context.Context, req *webservice.CancelDeploymentRequest) (*webservice.CancelDeploymentResponse, error) {
	deployment, err := a.getDeployment(ctx, req.DeploymentId)
	if err != nil {
		return nil, err
	}
	if model.IsCompletedDeployment(deployment.Status) {
		return nil, status.Errorf(codes.FailedPrecondition, "could not cancel the deployment because it was already completed")
	}

	commandID := uuid.New().String()
	cmd := model.Command{
		Id:            commandID,
		PipedId:       deployment.PipedId,
		ApplicationId: deployment.ApplicationId,
		DeploymentId:  req.DeploymentId,
		Type:          model.Command_CANCEL_DEPLOYMENT,
		Commander:     "anonymous",
		CancelDeployment: &model.Command_CancelDeployment{
			DeploymentId:    req.DeploymentId,
			WithoutRollback: req.WithoutRollback,
		},
	}
	if err := a.addCommand(ctx, &cmd); err != nil {
		return nil, err
	}
	return &webservice.CancelDeploymentResponse{
		CommandId: commandID,
	}, nil
}

func (a *WebAPI) ApproveStage(ctx context.Context, req *webservice.ApproveStageRequest) (*webservice.ApproveStageResponse, error) {
	deployment, err := a.getDeployment(ctx, req.DeploymentId)
	if err != nil {
		return nil, err
	}
	stage, ok := deployment.StageStatusMap()[req.StageId]
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, "stage was not found in the deployment")
	}
	if model.IsCompletedStage(stage) {
		return nil, status.Errorf(codes.FailedPrecondition, "could not approve the stage because it was already completed")
	}

	commandID := uuid.New().String()
	cmd := model.Command{
		Id:            commandID,
		PipedId:       deployment.PipedId,
		ApplicationId: deployment.ApplicationId,
		DeploymentId:  req.DeploymentId,
		StageId:       req.StageId,
		Type:          model.Command_APPROVE_STAGE,
		Commander:     "anonymous",
		ApproveStage: &model.Command_ApproveStage{
			DeploymentId: req.DeploymentId,
			StageId:      req.StageId,
		},
	}
	if err := a.addCommand(ctx, &cmd); err != nil {
		return nil, err
	}
	return &webservice.ApproveStageResponse{
		CommandId: commandID,
	}, nil
}

func (a *WebAPI) GetApplicationLiveState(ctx context.Context, req *webservice.GetApplicationLiveStateRequest) (*webservice.GetApplicationLiveStateResponse, error) {
	snapshot, err := a.applicationLiveStateStore.GetStateSnapshot(ctx, req.ApplicationId)
	if err != nil {
		a.logger.Error("failed to get application live state", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get application live state")
	}
	return &webservice.GetApplicationLiveStateResponse{
		Snapshot: snapshot,
	}, nil
}

func (a *WebAPI) GetProject(ctx context.Context, req *webservice.GetProjectRequest) (*webservice.GetProjectResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *WebAPI) GetMe(ctx context.Context, req *webservice.GetMeRequest) (*webservice.GetMeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *WebAPI) GetCommand(ctx context.Context, req *webservice.GetCommandRequest) (*webservice.GetCommandResponse, error) {
	cmd, err := a.commandStore.GetCommand(ctx, req.CommandId)
	if errors.Is(err, datastore.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "command is not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get command")
	}
	return &webservice.GetCommandResponse{
		Command: cmd,
	}, nil
}

func (a *WebAPI) ListDeploymentConfigTemplates(ctx context.Context, req *webservice.ListDeploymentConfigTemplatesRequest) (*webservice.ListDeploymentConfigTemplatesResponse, error) {
	res := &webservice.ListDeploymentConfigTemplatesResponse{Templates: []*webservice.DeploymentConfigTemplate{
		{
			ApplicationKind: model.ApplicationKind_KUBERNETES,
			Name:            "Canary Deployment",
			Content:         k8sCanaryDeploymentConfigTemplate,
		},
		{
			ApplicationKind: model.ApplicationKind_KUBERNETES,
			Name:            "Blue/Green Deployment",
			Content:         k8sBluebreenDeploymentConfigTemplate,
		},
	}}
	return res, nil
}
