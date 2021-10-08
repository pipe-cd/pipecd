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
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipe/pkg/app/api/applicationlivestatestore"
	"github.com/pipe-cd/pipe/pkg/app/api/commandstore"
	"github.com/pipe-cd/pipe/pkg/app/api/service/webservice"
	"github.com/pipe-cd/pipe/pkg/app/api/stagelogstore"
	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/cache/memorycache"
	"github.com/pipe-cd/pipe/pkg/cache/rediscache"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/filestore"
	"github.com/pipe-cd/pipe/pkg/git"
	"github.com/pipe-cd/pipe/pkg/insight/insightstore"
	"github.com/pipe-cd/pipe/pkg/model"
	"github.com/pipe-cd/pipe/pkg/redis"
	"github.com/pipe-cd/pipe/pkg/rpc/rpcauth"
)

type encrypter interface {
	Encrypt(text string) (string, error)
}

// WebAPI implements the behaviors for the gRPC definitions of WebAPI.
type WebAPI struct {
	applicationStore          datastore.ApplicationStore
	environmentStore          datastore.EnvironmentStore
	deploymentStore           datastore.DeploymentStore
	pipedStore                datastore.PipedStore
	projectStore              datastore.ProjectStore
	apiKeyStore               datastore.APIKeyStore
	stageLogStore             stagelogstore.Store
	applicationLiveStateStore applicationlivestatestore.Store
	commandStore              commandstore.Store
	insightStore              insightstore.Store
	encrypter                 encrypter

	appProjectCache        cache.Cache
	deploymentProjectCache cache.Cache
	pipedProjectCache      cache.Cache
	envProjectCache        cache.Cache
	insightCache           cache.Cache

	projectsInConfig map[string]config.ControlPlaneProject
	logger           *zap.Logger
}

// NewWebAPI creates a new WebAPI instance.
func NewWebAPI(
	ctx context.Context,
	ds datastore.DataStore,
	fs filestore.Store,
	sls stagelogstore.Store,
	alss applicationlivestatestore.Store,
	cmds commandstore.Store,
	is insightstore.Store,
	rd redis.Redis,
	projs map[string]config.ControlPlaneProject,
	encrypter encrypter,
	logger *zap.Logger) *WebAPI {
	a := &WebAPI{
		applicationStore:          datastore.NewApplicationStore(ds),
		environmentStore:          datastore.NewEnvironmentStore(ds),
		deploymentStore:           datastore.NewDeploymentStore(ds),
		pipedStore:                datastore.NewPipedStore(ds),
		projectStore:              datastore.NewProjectStore(ds),
		apiKeyStore:               datastore.NewAPIKeyStore(ds),
		stageLogStore:             sls,
		applicationLiveStateStore: alss,
		commandStore:              cmds,
		insightStore:              is,
		projectsInConfig:          projs,
		encrypter:                 encrypter,
		appProjectCache:           memorycache.NewTTLCache(ctx, 24*time.Hour, 3*time.Hour),
		deploymentProjectCache:    memorycache.NewTTLCache(ctx, 24*time.Hour, 3*time.Hour),
		pipedProjectCache:         memorycache.NewTTLCache(ctx, 24*time.Hour, 3*time.Hour),
		envProjectCache:           memorycache.NewTTLCache(ctx, 24*time.Hour, 3*time.Hour),
		insightCache:              rediscache.NewTTLCache(rd, 3*time.Hour),
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
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
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
		return nil, status.Error(codes.AlreadyExists, "The environment already exists")
	}
	if err != nil {
		a.logger.Error("failed to create environment", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to create environment")
	}

	return &webservice.AddEnvironmentResponse{}, nil
}

func (a *WebAPI) UpdateEnvironmentDesc(ctx context.Context, req *webservice.UpdateEnvironmentDescRequest) (*webservice.UpdateEnvironmentDescResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *WebAPI) ListEnvironments(ctx context.Context, req *webservice.ListEnvironmentsRequest) (*webservice.ListEnvironmentsResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	opts := datastore.ListOptions{
		Filters: []datastore.ListFilter{
			{
				Field:    "ProjectId",
				Operator: datastore.OperatorEqual,
				Value:    claims.Role.ProjectId,
			},
		},
	}
	envs, err := a.environmentStore.ListEnvironments(ctx, opts)
	if err != nil {
		a.logger.Error("failed to get environments", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to get environments")
	}

	return &webservice.ListEnvironmentsResponse{
		Environments: envs,
	}, nil
}

func (a *WebAPI) EnableEnvironment(ctx context.Context, req *webservice.EnableEnvironmentRequest) (*webservice.EnableEnvironmentResponse, error) {
	if err := a.updateEnvironmentEnable(ctx, req.EnvironmentId, true); err != nil {
		return nil, err
	}
	return &webservice.EnableEnvironmentResponse{}, nil
}

func (a *WebAPI) DisableEnvironment(ctx context.Context, req *webservice.DisableEnvironmentRequest) (*webservice.DisableEnvironmentResponse, error) {
	if err := a.updateEnvironmentEnable(ctx, req.EnvironmentId, false); err != nil {
		return nil, err
	}
	return &webservice.DisableEnvironmentResponse{}, nil
}

// DeleteEnvironment deletes the given environment and all applications that belong to it.
// It returns a FailedPrecondition error if any Piped is still using that environment.
func (a *WebAPI) DeleteEnvironment(ctx context.Context, req *webservice.DeleteEnvironmentRequest) (*webservice.DeleteEnvironmentResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	if err := a.validateEnvBelongsToProject(ctx, req.EnvironmentId, claims.Role.ProjectId); err != nil {
		return nil, err
	}
	// Check if no Piped has permission to the given environment.
	pipeds, err := a.pipedStore.ListPipeds(ctx, datastore.ListOptions{
		Filters: []datastore.ListFilter{
			{
				Field:    "ProjectId",
				Operator: datastore.OperatorEqual,
				Value:    claims.Role.ProjectId,
			},
			{
				Field:    "EnvIds",
				Operator: datastore.OperatorContains,
				Value:    req.EnvironmentId,
			},
			{
				Field:    "Disabled",
				Operator: datastore.OperatorEqual,
				Value:    false,
			},
		},
	})
	if err != nil {
		a.logger.Error("failed to fetch Pipeds linked to the given environment",
			zap.String("env-id", req.EnvironmentId),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "Failed to validate the deletion operation")
	}
	if len(pipeds) > 0 {
		pipedNames := make([]string, 0, len(pipeds))
		for _, p := range pipeds {
			pipedNames = append(pipedNames, p.Name)
		}
		return nil, status.Errorf(
			codes.FailedPrecondition,
			"Found Pipeds linked the environment to be deleted. Please remove this environment from all Pipeds (%s) on the Piped settings page",
			strings.Join(pipedNames, ","),
		)
	}

	// Delete all applications that belongs to the given env.
	apps, _, err := a.applicationStore.ListApplications(ctx, datastore.ListOptions{
		Filters: []datastore.ListFilter{
			{
				Field:    "ProjectId",
				Operator: datastore.OperatorEqual,
				Value:    claims.Role.ProjectId,
			},
			{
				Field:    "EnvId",
				Operator: datastore.OperatorEqual,
				Value:    req.EnvironmentId,
			},
		},
	})
	if err != nil {
		a.logger.Error("failed to fetch applications that belongs to the given environment",
			zap.String("env-id", req.EnvironmentId),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "Failed to fetch applications that belongs to the given environment")
	}
	for _, app := range apps {
		if app.ProjectId != claims.Role.ProjectId {
			continue
		}
		err := a.applicationStore.DeleteApplication(ctx, app.Id)
		if err == nil {
			continue
		}
		switch err {
		case datastore.ErrNotFound:
			return nil, status.Error(codes.Internal, "The application is not found")
		case datastore.ErrInvalidArgument:
			return nil, status.Error(codes.InvalidArgument, "Invalid value to delete")
		default:
			a.logger.Error("failed to delete the application",
				zap.String("application-id", app.Id),
				zap.Error(err),
			)
			return nil, status.Error(codes.Internal, "Failed to delete the application")
		}
	}

	if err := a.environmentStore.DeleteEnvironment(ctx, req.EnvironmentId); err != nil {
		switch err {
		case datastore.ErrNotFound:
			return nil, status.Error(codes.NotFound, "The environment is not found")
		case datastore.ErrInvalidArgument:
			return nil, status.Error(codes.InvalidArgument, "Invalid value to delete")
		default:
			a.logger.Error("failed to delete the environment",
				zap.String("env-id", req.EnvironmentId),
				zap.Error(err),
			)
			return nil, status.Error(codes.Internal, "Failed to delete the environment")
		}
	}

	return &webservice.DeleteEnvironmentResponse{}, nil
}

func (a *WebAPI) updateEnvironmentEnable(ctx context.Context, envID string, enable bool) error {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return err
	}

	if err := a.validateEnvBelongsToProject(ctx, envID, claims.Role.ProjectId); err != nil {
		return err
	}

	var updater func(context.Context, string) error
	if enable {
		updater = a.environmentStore.EnableEnvironment
	} else {
		updater = a.environmentStore.DisableEnvironment
	}

	if err := updater(ctx, envID); err != nil {
		switch err {
		case datastore.ErrNotFound:
			return status.Error(codes.NotFound, "The environment is not found")
		case datastore.ErrInvalidArgument:
			return status.Error(codes.InvalidArgument, "Invalid value for update")
		default:
			a.logger.Error("failed to update the environment",
				zap.String("env-id", envID),
				zap.Error(err),
			)
			return status.Error(codes.Internal, "Failed to update the environment")
		}
	}
	return nil
}

// validateEnvBelongsToProject checks if the given piped belongs to the given project.
// It gives back error unless the env belongs to the project.
func (a *WebAPI) validateEnvBelongsToProject(ctx context.Context, envID, projectID string) error {
	eid, err := a.envProjectCache.Get(envID)
	if err == nil {
		if projectID != eid {
			return status.Error(codes.PermissionDenied, "Requested environment doesn't belong to the project you logged in")
		}
		return nil
	}

	env, err := getEnvironment(ctx, a.environmentStore, envID, a.logger)
	if err != nil {
		return err
	}
	a.envProjectCache.Put(envID, env.ProjectId)

	if projectID != env.ProjectId {
		return status.Error(codes.PermissionDenied, "Requested environment doesn't belong to the project you logged in")
	}
	return nil
}

func (a *WebAPI) RegisterPiped(ctx context.Context, req *webservice.RegisterPipedRequest) (*webservice.RegisterPipedResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	key, keyHash, err := model.GeneratePipedKey()
	if err != nil {
		a.logger.Error("failed to generate piped key", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to generate the piped key")
	}

	piped := model.Piped{
		Id:        uuid.New().String(),
		Name:      req.Name,
		Desc:      req.Desc,
		ProjectId: claims.Role.ProjectId,
		EnvIds:    req.EnvIds,
		Status:    model.Piped_OFFLINE,
	}
	if err := piped.AddKey(keyHash, claims.Subject, time.Now()); err != nil {
		return nil, status.Error(codes.FailedPrecondition, fmt.Sprintf("Failed to create key: %v", err))
	}

	err = a.pipedStore.AddPiped(ctx, &piped)
	if errors.Is(err, datastore.ErrAlreadyExists) {
		return nil, status.Error(codes.AlreadyExists, "The piped already exists")
	}
	if err != nil {
		a.logger.Error("failed to register piped", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to register piped")
	}
	return &webservice.RegisterPipedResponse{
		Id:  piped.Id,
		Key: key,
	}, nil
}

func (a *WebAPI) UpdatePiped(ctx context.Context, req *webservice.UpdatePipedRequest) (*webservice.UpdatePipedResponse, error) {
	updater := func(ctx context.Context, pipedID string) error {
		return a.pipedStore.UpdatePiped(ctx, req.PipedId, func(p *model.Piped) error {
			p.Name = req.Name
			p.Desc = req.Desc
			p.EnvIds = req.EnvIds
			return nil
		})
	}
	if err := a.updatePiped(ctx, req.PipedId, updater); err != nil {
		return nil, err
	}

	return &webservice.UpdatePipedResponse{}, nil
}

func (a *WebAPI) RecreatePipedKey(ctx context.Context, req *webservice.RecreatePipedKeyRequest) (*webservice.RecreatePipedKeyResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	key, keyHash, err := model.GeneratePipedKey()
	if err != nil {
		a.logger.Error("failed to generate piped key", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to generate the piped key")
	}

	updater := func(ctx context.Context, pipedID string) error {
		return a.pipedStore.AddKey(ctx, pipedID, keyHash, claims.Subject, time.Now())
	}
	if err := a.updatePiped(ctx, req.Id, updater); err != nil {
		return nil, err
	}

	return &webservice.RecreatePipedKeyResponse{
		Key: key,
	}, nil
}

func (a *WebAPI) DeleteOldPipedKeys(ctx context.Context, req *webservice.DeleteOldPipedKeysRequest) (*webservice.DeleteOldPipedKeysResponse, error) {
	if _, err := rpcauth.ExtractClaims(ctx); err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	updater := func(ctx context.Context, pipedID string) error {
		return a.pipedStore.DeleteOldKeys(ctx, pipedID)
	}
	if err := a.updatePiped(ctx, req.PipedId, updater); err != nil {
		return nil, err
	}

	return &webservice.DeleteOldPipedKeysResponse{}, nil
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
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return err
	}

	if err := a.validatePipedBelongsToProject(ctx, pipedID, claims.Role.ProjectId); err != nil {
		return err
	}

	if err := updater(ctx, pipedID); err != nil {
		switch err {
		case datastore.ErrNotFound:
			return status.Error(codes.InvalidArgument, "The piped is not found")
		case datastore.ErrInvalidArgument:
			return status.Error(codes.InvalidArgument, "Invalid value for update")
		default:
			a.logger.Error("failed to update the piped",
				zap.String("piped-id", pipedID),
				zap.Error(err),
			)
			// TODO: Improve error handling, instead of considering all as Internal error like this
			// we should check the error type to decide to pass its message to the web client or just a generic message.
			return status.Error(codes.Internal, "Failed to update the piped")
		}
	}
	return nil
}

// TODO: Consider using piped-stats to decide piped connection status.
func (a *WebAPI) ListPipeds(ctx context.Context, req *webservice.ListPipedsRequest) (*webservice.ListPipedsResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	opts := datastore.ListOptions{
		Filters: []datastore.ListFilter{
			{
				Field:    "ProjectId",
				Operator: datastore.OperatorEqual,
				Value:    claims.Role.ProjectId,
			},
		},
	}

	if req.Options != nil {
		if req.Options.Enabled != nil {
			opts.Filters = append(opts.Filters, datastore.ListFilter{
				Field:    "Disabled",
				Operator: datastore.OperatorEqual,
				Value:    !req.Options.Enabled.GetValue(),
			})
		}
	}

	pipeds, err := a.pipedStore.ListPipeds(ctx, opts)
	if err != nil {
		a.logger.Error("failed to get pipeds", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to get pipeds")
	}

	// Redact all sensitive data inside piped message before sending to the client.
	for i := range pipeds {
		pipeds[i].RedactSensitiveData()
	}

	return &webservice.ListPipedsResponse{
		Pipeds: pipeds,
	}, nil
}

func (a *WebAPI) GetPiped(ctx context.Context, req *webservice.GetPipedRequest) (*webservice.GetPipedResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	piped, err := getPiped(ctx, a.pipedStore, req.PipedId, a.logger)
	if err != nil {
		return nil, err
	}
	if err := a.validatePipedBelongsToProject(ctx, req.PipedId, claims.Role.ProjectId); err != nil {
		return nil, err
	}

	// Redact all sensitive data inside piped message before sending to the client.
	piped.RedactSensitiveData()

	return &webservice.GetPipedResponse{
		Piped: piped,
	}, nil
}

func (a *WebAPI) UpdatePipedDesiredVersion(ctx context.Context, req *webservice.UpdatePipedDesiredVersionRequest) (*webservice.UpdatePipedDesiredVersionResponse, error) {
	updater := func(ctx context.Context, pipedID string) error {
		return a.pipedStore.UpdatePiped(ctx, pipedID, func(p *model.Piped) error {
			p.DesiredVersion = req.Version
			return nil
		})
	}
	for _, pipedID := range req.PipedIds {
		if err := a.updatePiped(ctx, pipedID, updater); err != nil {
			return nil, err
		}
	}

	return &webservice.UpdatePipedDesiredVersionResponse{}, nil
}

// validatePipedBelongsToProject checks if the given piped belongs to the given project.
// It gives back error unless the piped belongs to the project.
func (a *WebAPI) validatePipedBelongsToProject(ctx context.Context, pipedID, projectID string) error {
	pid, err := a.pipedProjectCache.Get(pipedID)
	if err == nil {
		if pid != projectID {
			return status.Error(codes.PermissionDenied, "Requested piped doesn't belong to the project you logged in")
		}
		return nil
	}

	piped, err := getPiped(ctx, a.pipedStore, pipedID, a.logger)
	if err != nil {
		return err
	}
	a.pipedProjectCache.Put(pipedID, piped.ProjectId)

	if piped.ProjectId != projectID {
		return status.Error(codes.PermissionDenied, "Requested piped doesn't belong to the project you logged in")
	}
	return nil
}

// TODO: Validate the specified piped to ensure that it belongs to the specified environment.
func (a *WebAPI) AddApplication(ctx context.Context, req *webservice.AddApplicationRequest) (*webservice.AddApplicationResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	piped, err := getPiped(ctx, a.pipedStore, req.PipedId, a.logger)
	if err != nil {
		return nil, err
	}

	if piped.ProjectId != claims.Role.ProjectId {
		return nil, status.Error(codes.InvalidArgument, "Requested piped does not belong to your project")
	}

	gitpath, err := makeGitPath(
		req.GitPath.Repo.Id,
		req.GitPath.Path,
		req.GitPath.ConfigFilename,
		piped,
		a.logger,
	)
	if err != nil {
		return nil, err
	}

	app := model.Application{
		Id:            uuid.New().String(),
		Name:          req.Name,
		EnvId:         req.EnvId,
		PipedId:       req.PipedId,
		ProjectId:     claims.Role.ProjectId,
		GitPath:       gitpath,
		Kind:          req.Kind,
		CloudProvider: req.CloudProvider,
		Description:   req.Description,
	}
	err = a.applicationStore.AddApplication(ctx, &app)
	if errors.Is(err, datastore.ErrAlreadyExists) {
		return nil, status.Error(codes.AlreadyExists, "The application already exists")
	}
	if err != nil {
		a.logger.Error("failed to create application", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to create application")
	}

	return &webservice.AddApplicationResponse{
		ApplicationId: app.Id,
	}, nil
}

func (a *WebAPI) UpdateApplication(ctx context.Context, req *webservice.UpdateApplicationRequest) (*webservice.UpdateApplicationResponse, error) {
	updater := func(app *model.Application) error {
		app.Name = req.Name
		app.EnvId = req.EnvId
		app.PipedId = req.PipedId
		app.Kind = req.Kind
		app.CloudProvider = req.CloudProvider
		return nil
	}

	if err := a.updateApplication(ctx, req.ApplicationId, req.PipedId, updater); err != nil {
		return nil, err
	}
	return &webservice.UpdateApplicationResponse{}, nil
}

func (a *WebAPI) UpdateApplicationDescription(ctx context.Context, req *webservice.UpdateApplicationDescriptionRequest) (*webservice.UpdateApplicationDescriptionResponse, error) {
	updater := func(app *model.Application) error {
		app.Description = req.Description
		return nil
	}

	if err := a.updateApplication(ctx, req.ApplicationId, "", updater); err != nil {
		return nil, err
	}
	return &webservice.UpdateApplicationDescriptionResponse{}, nil
}

func (a *WebAPI) updateApplication(ctx context.Context, id, pipedID string, updater func(app *model.Application) error) error {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return err
	}

	// Ensure that the specified piped is assignable for this application.
	if pipedID != "" {
		piped, err := getPiped(ctx, a.pipedStore, pipedID, a.logger)
		if err != nil {
			return err
		}

		if piped.ProjectId != claims.Role.ProjectId {
			return status.Error(codes.InvalidArgument, "Requested piped does not belong to your project")
		}
	}

	err = a.applicationStore.UpdateApplication(ctx, id, updater)
	if err != nil {
		a.logger.Error("failed to update application", zap.Error(err))
		return status.Error(codes.Internal, "Failed to update application")
	}

	return nil
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

func (a *WebAPI) DeleteApplication(ctx context.Context, req *webservice.DeleteApplicationRequest) (*webservice.DeleteApplicationResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	if err := a.validateAppBelongsToProject(ctx, req.ApplicationId, claims.Role.ProjectId); err != nil {
		return nil, err
	}

	if err := a.applicationStore.DeleteApplication(ctx, req.ApplicationId); err != nil {
		switch err {
		case datastore.ErrNotFound:
			return nil, status.Error(codes.NotFound, "The application is not found")
		case datastore.ErrInvalidArgument:
			return nil, status.Error(codes.InvalidArgument, "Invalid value to delete")
		default:
			a.logger.Error("failed to delete the application",
				zap.String("application-id", req.ApplicationId),
				zap.Error(err),
			)
			return nil, status.Error(codes.Internal, "Failed to delete the application")
		}
	}

	return &webservice.DeleteApplicationResponse{}, nil
}

func (a *WebAPI) updateApplicationEnable(ctx context.Context, appID string, enable bool) error {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return err
	}

	if err := a.validateAppBelongsToProject(ctx, appID, claims.Role.ProjectId); err != nil {
		return err
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
			return status.Error(codes.NotFound, "The application is not found")
		case datastore.ErrInvalidArgument:
			return status.Error(codes.InvalidArgument, "Invalid value for update")
		default:
			a.logger.Error("failed to update the application",
				zap.String("application-id", appID),
				zap.Error(err),
			)
			return status.Error(codes.Internal, "Failed to update the application")
		}
	}
	return nil
}

func (a *WebAPI) ListApplications(ctx context.Context, req *webservice.ListApplicationsRequest) (*webservice.ListApplicationsResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	orders := []datastore.Order{
		{
			Field:     "UpdatedAt",
			Direction: datastore.Desc,
		},
		{
			Field:     "Id",
			Direction: datastore.Asc,
		},
	}
	filters := []datastore.ListFilter{
		{
			Field:    "ProjectId",
			Operator: datastore.OperatorEqual,
			Value:    claims.Role.ProjectId,
		},
	}
	if o := req.Options; o != nil {
		if o.Enabled != nil {
			filters = append(filters, datastore.ListFilter{
				Field:    "Disabled",
				Operator: datastore.OperatorEqual,
				Value:    !o.Enabled.GetValue(),
			})
		}
		// Allowing multiple so that it can do In Query later.
		// Currently only the first value is used.
		if len(o.Kinds) > 0 {
			filters = append(filters, datastore.ListFilter{
				Field:    "Kind",
				Operator: datastore.OperatorEqual,
				Value:    o.Kinds[0],
			})
		}
		if len(o.SyncStatuses) > 0 {
			filters = append(filters, datastore.ListFilter{
				Field:    "SyncState.Status",
				Operator: datastore.OperatorEqual,
				Value:    o.SyncStatuses[0],
			})
		}
		if len(o.EnvIds) > 0 {
			filters = append(filters, datastore.ListFilter{
				Field:    "EnvId",
				Operator: datastore.OperatorEqual,
				Value:    o.EnvIds[0],
			})
		}
		if o.Name != "" {
			filters = append(filters, datastore.ListFilter{
				Field:    "Name",
				Operator: datastore.OperatorEqual,
				Value:    o.Name,
			})
		}
	}

	apps, _, err := a.applicationStore.ListApplications(ctx, datastore.ListOptions{
		Filters: filters,
		Orders:  orders,
	})
	if err != nil {
		a.logger.Error("failed to get applications", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to get applications")
	}

	return &webservice.ListApplicationsResponse{
		Applications: apps,
	}, nil
}

func (a *WebAPI) SyncApplication(ctx context.Context, req *webservice.SyncApplicationRequest) (*webservice.SyncApplicationResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	app, err := getApplication(ctx, a.applicationStore, req.ApplicationId, a.logger)
	if err != nil {
		return nil, err
	}

	if claims.Role.ProjectId != app.ProjectId {
		return nil, status.Error(codes.InvalidArgument, "Requested application does not belong to your project")
	}

	cmd := model.Command{
		Id:            uuid.New().String(),
		PipedId:       app.PipedId,
		ApplicationId: app.Id,
		ProjectId:     app.ProjectId,
		Type:          model.Command_SYNC_APPLICATION,
		Commander:     claims.Subject,
		SyncApplication: &model.Command_SyncApplication{
			ApplicationId: app.Id,
			SyncStrategy:  req.SyncStrategy,
		},
	}
	if err := addCommand(ctx, a.commandStore, &cmd, a.logger); err != nil {
		return nil, err
	}

	return &webservice.SyncApplicationResponse{
		CommandId: cmd.Id,
	}, nil
}

func (a *WebAPI) GetApplication(ctx context.Context, req *webservice.GetApplicationRequest) (*webservice.GetApplicationResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	app, err := getApplication(ctx, a.applicationStore, req.ApplicationId, a.logger)
	if err != nil {
		return nil, err
	}

	if app.ProjectId != claims.Role.ProjectId {
		return nil, status.Error(codes.InvalidArgument, "Requested application does not belong to your project")
	}

	return &webservice.GetApplicationResponse{
		Application: app,
	}, nil
}

func (a *WebAPI) GenerateApplicationSealedSecret(ctx context.Context, req *webservice.GenerateApplicationSealedSecretRequest) (*webservice.GenerateApplicationSealedSecretResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	piped, err := getPiped(ctx, a.pipedStore, req.PipedId, a.logger)
	if err != nil {
		return nil, err
	}
	if err := a.validatePipedBelongsToProject(ctx, req.PipedId, claims.Role.ProjectId); err != nil {
		return nil, err
	}

	se := model.GetSecretEncryptionInPiped(piped)
	pubkey, err := getEncriptionKey(se)
	if err != nil {
		return nil, err
	}

	ciphertext, err := encrypt(req.Data, pubkey, req.Base64Encoding, a.logger)
	if err != nil {
		return nil, err
	}

	return &webservice.GenerateApplicationSealedSecretResponse{
		Data: ciphertext,
	}, nil
}

// validateAppBelongsToProject checks if the given application belongs to the given project.
// It gives back error unless the application belongs to the project.
func (a *WebAPI) validateAppBelongsToProject(ctx context.Context, appID, projectID string) error {
	pid, err := a.appProjectCache.Get(appID)
	if err == nil {
		if pid != projectID {
			return status.Error(codes.PermissionDenied, "Requested application doesn't belong to the project you logged in")
		}
		return nil
	}

	app, err := getApplication(ctx, a.applicationStore, appID, a.logger)
	if err != nil {
		return err
	}
	a.appProjectCache.Put(appID, app.ProjectId)

	if app.ProjectId != projectID {
		return status.Error(codes.PermissionDenied, "Requested application doesn't belong to the project you logged in")
	}
	return nil
}

func (a *WebAPI) ListDeployments(ctx context.Context, req *webservice.ListDeploymentsRequest) (*webservice.ListDeploymentsResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	orders := []datastore.Order{
		{
			Field:     "UpdatedAt",
			Direction: datastore.Desc,
		},
		{
			Field:     "Id",
			Direction: datastore.Asc,
		},
	}
	filters := []datastore.ListFilter{
		{
			Field:    "ProjectId",
			Operator: datastore.OperatorEqual,
			Value:    claims.Role.ProjectId,
		},
	}
	if o := req.Options; o != nil {
		// Allowing multiple so that it can do In Query later.
		// Currently only the first value is used.
		if len(o.Statuses) > 0 {
			filters = append(filters, datastore.ListFilter{
				Field:    "Status",
				Operator: datastore.OperatorEqual,
				Value:    o.Statuses[0],
			})
		}
		if len(o.Kinds) > 0 {
			filters = append(filters, datastore.ListFilter{
				Field:    "Kind",
				Operator: datastore.OperatorEqual,
				Value:    o.Kinds[0],
			})
		}
		if len(o.ApplicationIds) > 0 {
			filters = append(filters, datastore.ListFilter{
				Field:    "ApplicationId",
				Operator: datastore.OperatorEqual,
				Value:    o.ApplicationIds[0],
			})
		}
		if len(o.EnvIds) > 0 {
			filters = append(filters, datastore.ListFilter{
				Field:    "EnvId",
				Operator: datastore.OperatorEqual,
				Value:    o.EnvIds[0],
			})
		}
		if o.ApplicationName != "" {
			filters = append(filters, datastore.ListFilter{
				Field:    "ApplicationName",
				Operator: datastore.OperatorEqual,
				Value:    o.ApplicationName,
			})
		}
	}

	deployments, cursor, err := a.deploymentStore.ListDeployments(ctx, datastore.ListOptions{
		Filters: filters,
		Orders:  orders,
		Limit:   int(req.PageSize),
		Cursor:  req.Cursor,
	})
	if err != nil {
		a.logger.Error("failed to get deployments", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to get deployments")
	}
	return &webservice.ListDeploymentsResponse{
		Deployments: deployments,
		Cursor:      cursor,
	}, nil
}

func (a *WebAPI) GetDeployment(ctx context.Context, req *webservice.GetDeploymentRequest) (*webservice.GetDeploymentResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	deployment, err := getDeployment(ctx, a.deploymentStore, req.DeploymentId, a.logger)
	if err != nil {
		return nil, err
	}

	if claims.Role.ProjectId != deployment.ProjectId {
		return nil, status.Error(codes.InvalidArgument, "Requested deployment does not belong to your project")
	}

	return &webservice.GetDeploymentResponse{
		Deployment: deployment,
	}, nil
}

// validateDeploymentBelongsToProject checks if the given deployment belongs to the given project.
// It gives back error unless the deployment belongs to the project.
func (a *WebAPI) validateDeploymentBelongsToProject(ctx context.Context, deploymentID, projectID string) error {
	pid, err := a.deploymentProjectCache.Get(deploymentID)
	if err == nil {
		if pid != projectID {
			return status.Error(codes.PermissionDenied, "Requested deployment doesn't belong to the project you logged in")
		}
		return nil
	}

	deployment, err := getDeployment(ctx, a.deploymentStore, deploymentID, a.logger)
	if err != nil {
		return err
	}
	a.deploymentProjectCache.Put(deploymentID, deployment.ProjectId)

	if deployment.ProjectId != projectID {
		return status.Error(codes.PermissionDenied, "Requested deployment doesn't belong to the project you logged in")
	}
	return nil
}

func (a *WebAPI) GetStageLog(ctx context.Context, req *webservice.GetStageLogRequest) (*webservice.GetStageLogResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	if err := a.validateDeploymentBelongsToProject(ctx, req.DeploymentId, claims.Role.ProjectId); err != nil {
		return nil, err
	}

	blocks, completed, err := a.stageLogStore.FetchLogs(ctx, req.DeploymentId, req.StageId, req.RetriedCount, req.OffsetIndex)
	if errors.Is(err, stagelogstore.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "The stage log not found")
	}
	if err != nil {
		a.logger.Error("failed to get stage logs", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to get stage logs")
	}

	return &webservice.GetStageLogResponse{
		Blocks:    blocks,
		Completed: completed,
	}, nil
}

func (a *WebAPI) CancelDeployment(ctx context.Context, req *webservice.CancelDeploymentRequest) (*webservice.CancelDeploymentResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	deployment, err := getDeployment(ctx, a.deploymentStore, req.DeploymentId, a.logger)
	if err != nil {
		return nil, err
	}

	if claims.Role.ProjectId != deployment.ProjectId {
		return nil, status.Error(codes.InvalidArgument, "Requested deployment does not belong to your project")
	}

	if model.IsCompletedDeployment(deployment.Status) {
		return nil, status.Errorf(codes.FailedPrecondition, "could not cancel the deployment because it was already completed")
	}

	cmd := model.Command{
		Id:            uuid.New().String(),
		PipedId:       deployment.PipedId,
		ApplicationId: deployment.ApplicationId,
		ProjectId:     deployment.ProjectId,
		DeploymentId:  req.DeploymentId,
		Type:          model.Command_CANCEL_DEPLOYMENT,
		Commander:     claims.Subject,
		CancelDeployment: &model.Command_CancelDeployment{
			DeploymentId:    req.DeploymentId,
			ForceRollback:   req.ForceRollback,
			ForceNoRollback: req.ForceNoRollback,
		},
	}
	if err := addCommand(ctx, a.commandStore, &cmd, a.logger); err != nil {
		return nil, err
	}

	return &webservice.CancelDeploymentResponse{
		CommandId: cmd.Id,
	}, nil
}

func (a *WebAPI) ApproveStage(ctx context.Context, req *webservice.ApproveStageRequest) (*webservice.ApproveStageResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	deployment, err := getDeployment(ctx, a.deploymentStore, req.DeploymentId, a.logger)
	if err != nil {
		return nil, err
	}
	if err := a.validateApprover(deployment.Stages, claims.Subject, req.StageId); err != nil {
		return nil, err
	}
	if err := a.validateDeploymentBelongsToProject(ctx, req.DeploymentId, claims.Role.ProjectId); err != nil {
		return nil, err
	}
	stage, ok := deployment.StageStatusMap()[req.StageId]
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, "The stage was not found in the deployment")
	}
	if model.IsCompletedStage(stage) {
		return nil, status.Errorf(codes.FailedPrecondition, "Could not approve the stage because it was already completed")
	}

	commandID := uuid.New().String()
	cmd := model.Command{
		Id:            commandID,
		PipedId:       deployment.PipedId,
		ApplicationId: deployment.ApplicationId,
		ProjectId:     deployment.ProjectId,
		DeploymentId:  req.DeploymentId,
		StageId:       req.StageId,
		Type:          model.Command_APPROVE_STAGE,
		Commander:     claims.Subject,
		ApproveStage: &model.Command_ApproveStage{
			DeploymentId: req.DeploymentId,
			StageId:      req.StageId,
		},
	}
	if err := addCommand(ctx, a.commandStore, &cmd, a.logger); err != nil {
		return nil, err
	}

	return &webservice.ApproveStageResponse{
		CommandId: commandID,
	}, nil
}

// No error means that the given commander is valid.
func (a *WebAPI) validateApprover(stages []*model.PipelineStage, commander, stageID string) error {
	var approvers []string
	for _, s := range stages {
		if s.Id == stageID {
			approvers = strings.Split(s.Metadata["Approvers"], ",")
		}
	}
	if len(approvers) == 0 {
		// Anyone can approve the deployment pipeline
		return nil
	}
	for _, ap := range approvers {
		if ap == commander {
			return nil
		}
	}
	return status.Error(codes.FailedPrecondition, "Could not approve the stage because you aren't an approver.")
}

func (a *WebAPI) GetApplicationLiveState(ctx context.Context, req *webservice.GetApplicationLiveStateRequest) (*webservice.GetApplicationLiveStateResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	if err := a.validateAppBelongsToProject(ctx, req.ApplicationId, claims.Role.ProjectId); err != nil {
		return nil, err
	}

	snapshot, err := a.applicationLiveStateStore.GetStateSnapshot(ctx, req.ApplicationId)
	if errors.Is(err, filestore.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "Application live state not found")
	}
	if err != nil {
		a.logger.Error("failed to get application live state", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to get application live state")
	}
	return &webservice.GetApplicationLiveStateResponse{
		Snapshot: snapshot,
	}, nil
}

// GetProject gets the specified porject without sensitive data.
func (a *WebAPI) GetProject(ctx context.Context, req *webservice.GetProjectRequest) (*webservice.GetProjectResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	project, err := a.getProject(ctx, claims.Role.ProjectId)
	if err != nil {
		return nil, err
	}

	// Redact all sensitive data inside project message before sending to the client.
	project.RedactSensitiveData()

	return &webservice.GetProjectResponse{
		Project: project,
	}, nil
}

func (a *WebAPI) getProject(ctx context.Context, projectID string) (*model.Project, error) {
	if p, ok := a.projectsInConfig[projectID]; ok {
		return &model.Project{
			Id:   p.Id,
			Desc: p.Desc,
			StaticAdmin: &model.ProjectStaticUser{
				Username:     p.StaticAdmin.Username,
				PasswordHash: p.StaticAdmin.PasswordHash,
			},
		}, nil
	}

	project, err := a.projectStore.GetProject(ctx, projectID)
	if errors.Is(err, datastore.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "The project is not found")
	}
	if err != nil {
		a.logger.Error("failed to get project", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to get project")
	}
	return project, nil
}

// UpdateProjectStaticAdmin updates the static admin user settings.
func (a *WebAPI) UpdateProjectStaticAdmin(ctx context.Context, req *webservice.UpdateProjectStaticAdminRequest) (*webservice.UpdateProjectStaticAdminResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	if _, ok := a.projectsInConfig[claims.Role.ProjectId]; ok {
		return nil, status.Error(codes.FailedPrecondition, "Failed to update a debug project specified in the control-plane configuration")
	}

	if err := a.projectStore.UpdateProjectStaticAdmin(ctx, claims.Role.ProjectId, req.Username, req.Password); err != nil {
		a.logger.Error("failed to update static admin", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to update static admin")
	}
	return &webservice.UpdateProjectStaticAdminResponse{}, nil
}

// EnableStaticAdmin enables static admin login.
func (a *WebAPI) EnableStaticAdmin(ctx context.Context, req *webservice.EnableStaticAdminRequest) (*webservice.EnableStaticAdminResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	if _, ok := a.projectsInConfig[claims.Role.ProjectId]; ok {
		return nil, status.Error(codes.FailedPrecondition, "Failed to update a debug project specified in the control-plane configuration")
	}

	if err := a.projectStore.EnableStaticAdmin(ctx, claims.Role.ProjectId); err != nil {
		a.logger.Error("failed to enable static admin login", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to enable static admin login")
	}
	return &webservice.EnableStaticAdminResponse{}, nil
}

// DisableStaticAdmin disables static admin login.
func (a *WebAPI) DisableStaticAdmin(ctx context.Context, req *webservice.DisableStaticAdminRequest) (*webservice.DisableStaticAdminResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	if _, ok := a.projectsInConfig[claims.Role.ProjectId]; ok {
		return nil, status.Error(codes.FailedPrecondition, "Failed to update a debug project specified in the control-plane configuration")
	}

	if err := a.projectStore.DisableStaticAdmin(ctx, claims.Role.ProjectId); err != nil {
		a.logger.Error("failed to disenable static admin login", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to disenable static admin login")
	}
	return &webservice.DisableStaticAdminResponse{}, nil
}

// UpdateProjectSSOConfig updates the sso settings.
func (a *WebAPI) UpdateProjectSSOConfig(ctx context.Context, req *webservice.UpdateProjectSSOConfigRequest) (*webservice.UpdateProjectSSOConfigResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	if _, ok := a.projectsInConfig[claims.Role.ProjectId]; ok {
		return nil, status.Error(codes.FailedPrecondition, "Failed to update a debug project specified in the control-plane configuration")
	}

	if err := req.Sso.Encrypt(a.encrypter); err != nil {
		a.logger.Error("failed to encrypt sensitive data in sso configurations", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to encrypt sensitive data in sso configurations")
	}

	if err := a.projectStore.UpdateProjectSSOConfig(ctx, claims.Role.ProjectId, req.Sso); err != nil {
		a.logger.Error("failed to update project single sign on settings", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to update project single sign on settings")
	}
	return &webservice.UpdateProjectSSOConfigResponse{}, nil
}

// UpdateProjectRBACConfig updates the sso settings.
func (a *WebAPI) UpdateProjectRBACConfig(ctx context.Context, req *webservice.UpdateProjectRBACConfigRequest) (*webservice.UpdateProjectRBACConfigResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	if _, ok := a.projectsInConfig[claims.Role.ProjectId]; ok {
		return nil, status.Error(codes.FailedPrecondition, "Failed to update a debug project specified in the control-plane configuration")
	}

	if err := a.projectStore.UpdateProjectRBACConfig(ctx, claims.Role.ProjectId, req.Rbac); err != nil {
		a.logger.Error("failed to update project single sign on settings", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to update project single sign on settings")
	}
	return &webservice.UpdateProjectRBACConfigResponse{}, nil
}

// GetMe gets information about the current user.
func (a *WebAPI) GetMe(ctx context.Context, req *webservice.GetMeRequest) (*webservice.GetMeResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	return &webservice.GetMeResponse{
		Subject:     claims.Subject,
		AvatarUrl:   claims.AvatarURL,
		ProjectId:   claims.Role.ProjectId,
		ProjectRole: claims.Role.ProjectRole,
	}, nil
}

func (a *WebAPI) GetCommand(ctx context.Context, req *webservice.GetCommandRequest) (*webservice.GetCommandResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	cmd, err := getCommand(ctx, a.commandStore, req.CommandId, a.logger)
	if err != nil {
		return nil, err
	}

	if claims.Role.ProjectId != cmd.ProjectId {
		return nil, status.Error(codes.InvalidArgument, "Requested command does not belong to your project")
	}

	return &webservice.GetCommandResponse{
		Command: cmd,
	}, nil
}

func (a *WebAPI) ListDeploymentConfigTemplates(ctx context.Context, req *webservice.ListDeploymentConfigTemplatesRequest) (*webservice.ListDeploymentConfigTemplatesResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	app, err := getApplication(ctx, a.applicationStore, req.ApplicationId, a.logger)
	if err != nil {
		return nil, err
	}
	if err := a.validateAppBelongsToProject(ctx, req.ApplicationId, claims.Role.ProjectId); err != nil {
		return nil, err
	}

	var templates []*webservice.DeploymentConfigTemplate
	switch app.Kind {
	case model.ApplicationKind_KUBERNETES:
		templates = k8sDeploymentConfigTemplates
	case model.ApplicationKind_TERRAFORM:
		templates = terraformDeploymentConfigTemplates
	case model.ApplicationKind_LAMBDA:
		templates = lambdaDeploymentConfigTemplates
	case model.ApplicationKind_CLOUDRUN:
		templates = cloudrunDeploymentConfigTemplates
	case model.ApplicationKind_ECS:
		templates = ecsDeploymentConfigTemplates
	default:
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Unknown application kind %v", app.Kind))
	}
	for _, t := range templates {
		g := app.GetGitPath()
		filename := g.ConfigFilename
		if filename == "" {
			filename = ".pipe.yaml"
		}
		t.FileCreationUrl, err = git.MakeFileCreationURL(g.Repo.Remote, g.Path, g.Repo.Branch, filename, t.Content)
		if err != nil {
			a.logger.Error("failed to make a link to create a file", zap.Error(err))
			return nil, status.Error(codes.Internal, "Failed to make a link to create a file")
		}
	}

	if len(req.Labels) == 0 {
		return &webservice.ListDeploymentConfigTemplatesResponse{Templates: templates}, nil
	}

	filtered := filterDeploymentConfigTemplates(templates, req.Labels)
	return &webservice.ListDeploymentConfigTemplatesResponse{Templates: filtered}, nil
}

// Returns the one from the given templates with all the specified labels.
func filterDeploymentConfigTemplates(templates []*webservice.DeploymentConfigTemplate, labels []webservice.DeploymentConfigTemplateLabel) []*webservice.DeploymentConfigTemplate {
	filtered := make([]*webservice.DeploymentConfigTemplate, 0, len(templates))
L:
	for _, template := range templates {
		for _, l := range labels {
			if !template.HasLabel(l) {
				continue L
			}
		}
		filtered = append(filtered, template)
	}
	return filtered
}

func (a *WebAPI) GenerateAPIKey(ctx context.Context, req *webservice.GenerateAPIKeyRequest) (*webservice.GenerateAPIKeyResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	id := uuid.New().String()
	key, hash, err := model.GenerateAPIKey(id)
	if err != nil {
		a.logger.Error("failed to generate API key", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to generate API key")
	}

	apiKey := model.APIKey{
		Id:        id,
		Name:      req.Name,
		KeyHash:   hash,
		ProjectId: claims.Role.ProjectId,
		Role:      req.Role,
		Creator:   claims.Subject,
	}

	err = a.apiKeyStore.AddAPIKey(ctx, &apiKey)
	if errors.Is(err, datastore.ErrAlreadyExists) {
		return nil, status.Error(codes.AlreadyExists, "The API key already exists")
	}
	if err != nil {
		a.logger.Error("failed to create API key", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to create API key")
	}

	return &webservice.GenerateAPIKeyResponse{
		Key: key,
	}, nil
}

func (a *WebAPI) DisableAPIKey(ctx context.Context, req *webservice.DisableAPIKeyRequest) (*webservice.DisableAPIKeyResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	if err := a.apiKeyStore.DisableAPIKey(ctx, req.Id, claims.Role.ProjectId); err != nil {
		switch err {
		case datastore.ErrNotFound:
			return nil, status.Error(codes.InvalidArgument, "The API key is not found")
		case datastore.ErrInvalidArgument:
			return nil, status.Error(codes.InvalidArgument, "Invalid value for update")
		default:
			a.logger.Error("failed to disable the API key",
				zap.String("apikey-id", req.Id),
				zap.Error(err),
			)
			return nil, status.Error(codes.Internal, "Failed to disable the API key")
		}
	}

	return &webservice.DisableAPIKeyResponse{}, nil
}

func (a *WebAPI) ListAPIKeys(ctx context.Context, req *webservice.ListAPIKeysRequest) (*webservice.ListAPIKeysResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	opts := datastore.ListOptions{
		Filters: []datastore.ListFilter{
			{
				Field:    "ProjectId",
				Operator: datastore.OperatorEqual,
				Value:    claims.Role.ProjectId,
			},
		},
	}

	if req.Options != nil {
		if req.Options.Enabled != nil {
			opts.Filters = append(opts.Filters, datastore.ListFilter{
				Field:    "Disabled",
				Operator: datastore.OperatorEqual,
				Value:    !req.Options.Enabled.GetValue(),
			})
		}
	}

	apiKeys, err := a.apiKeyStore.ListAPIKeys(ctx, opts)
	if err != nil {
		a.logger.Error("failed to list API keys", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to list API keys")
	}

	// Redact all sensitive data inside API key before sending to the client.
	for i := range apiKeys {
		apiKeys[i].RedactSensitiveData()
	}

	return &webservice.ListAPIKeysResponse{
		Keys: apiKeys,
	}, nil
}

// GetInsightData returns the accumulated insight data.
func (a *WebAPI) GetInsightData(ctx context.Context, req *webservice.GetInsightDataRequest) (*webservice.GetInsightDataResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	count := int(req.DataPointCount)
	from := time.Unix(req.RangeFrom, 0)

	chunks, err := insightstore.LoadChunksFromCache(a.insightCache, claims.Role.ProjectId, req.ApplicationId, req.MetricsKind, req.Step, from, count)
	if err != nil {
		a.logger.Error("failed to load chunks from cache", zap.Error(err))

		chunks, err = a.insightStore.LoadChunks(ctx, claims.Role.ProjectId, req.ApplicationId, req.MetricsKind, req.Step, from, count)
		if err != nil {
			a.logger.Error("failed to load chunks from insightstore", zap.Error(err))
			return nil, err
		}
		if err := insightstore.PutChunksToCache(a.insightCache, chunks); err != nil {
			a.logger.Error("failed to put chunks to cache", zap.Error(err))
		}
	}

	idp, err := chunks.ExtractDataPoints(req.Step, from, count)
	if err != nil {
		a.logger.Error("failed to extract data points from chunks", zap.Error(err))
	}

	var updateAt int64
	for _, c := range chunks {
		accumulatedTo := c.GetAccumulatedTo()
		if accumulatedTo > updateAt {
			updateAt = accumulatedTo
		}
	}

	return &webservice.GetInsightDataResponse{
		UpdatedAt:  updateAt,
		DataPoints: idp,
		Type:       model.InsightResultType_MATRIX,
		Matrix: []*model.InsightSampleStream{
			{
				DataPoints: idp,
			},
		},
	}, nil
}

func (a *WebAPI) GetInsightApplicationCount(ctx context.Context, req *webservice.GetInsightApplicationCountRequest) (*webservice.GetInsightApplicationCountResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	// TODO: Cache application counts in the cache service.
	c, err := a.insightStore.LoadApplicationCounts(ctx, claims.Role.ProjectId)
	if err != nil {
		if err == filestore.ErrNotFound {
			return nil, status.Error(codes.NotFound, "Not found")
		}
		a.logger.Error("failed to load application counts", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to load application counts")
	}

	counts := make([]*model.InsightApplicationCount, 0, len(c.Counts))
	for i := range c.Counts {
		counts = append(counts, &c.Counts[i])
	}

	return &webservice.GetInsightApplicationCountResponse{
		Counts:    counts,
		UpdatedAt: c.UpdatedAt,
	}, nil
}
