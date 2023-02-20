// Copyright 2023 The PipeCD Authors.
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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v29/github"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/app/server/applicationlivestatestore"
	"github.com/pipe-cd/pipecd/pkg/app/server/commandstore"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/webservice"
	"github.com/pipe-cd/pipecd/pkg/app/server/stagelogstore"
	"github.com/pipe-cd/pipecd/pkg/app/server/unregisteredappstore"
	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/cache/memorycache"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/insight"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcauth"
)

type encrypter interface {
	Encrypt(text string) (string, error)
}

type webApiApplicationStore interface {
	Add(ctx context.Context, app *model.Application) error
	Get(ctx context.Context, id string) (*model.Application, error)
	List(ctx context.Context, opts datastore.ListOptions) ([]*model.Application, string, error)
	Delete(ctx context.Context, id string) error
	UpdateConfiguration(ctx context.Context, id, pipedID, platformProvider, configFilename string) error
	Enable(ctx context.Context, id string) error
	Disable(ctx context.Context, id string) error
}

type webApiDeploymentChainStore interface {
	Get(ctx context.Context, id string) (*model.DeploymentChain, error)
	List(ctx context.Context, opts datastore.ListOptions) ([]*model.DeploymentChain, string, error)
}

type webApiDeploymentStore interface {
	Get(ctx context.Context, id string) (*model.Deployment, error)
	List(ctx context.Context, opts datastore.ListOptions) ([]*model.Deployment, string, error)
}

type webApiPipedStore interface {
	Add(ctx context.Context, piped *model.Piped) error
	Get(ctx context.Context, id string) (*model.Piped, error)
	List(ctx context.Context, opts datastore.ListOptions) ([]*model.Piped, error)
	AddKey(ctx context.Context, id, keyHash, creator string, createdAt time.Time) error
	DeleteOldKeys(ctx context.Context, id string) error
	UpdateInfo(ctx context.Context, id, name, desc string) error
	EnablePiped(ctx context.Context, id string) error
	DisablePiped(ctx context.Context, id string) error
	UpdateDesiredVersion(ctx context.Context, id, version string) error
}

type webApiProjectStore interface {
	Get(ctx context.Context, id string) (*model.Project, error)
	UpdateProjectStaticAdmin(ctx context.Context, id, username, password string) error
	EnableStaticAdmin(ctx context.Context, id string) error
	DisableStaticAdmin(ctx context.Context, id string) error
	UpdateProjectSSOConfig(ctx context.Context, id string, sso *model.ProjectSSOConfig) error
	UpdateProjectRBACConfig(ctx context.Context, id string, sso *model.ProjectRBACConfig) error
	AddProjectRBACRole(ctx context.Context, id, name string, policies []*model.ProjectRBACPolicy) error
	UpdateProjectRBACRole(ctx context.Context, id, name string, policies []*model.ProjectRBACPolicy) error
	DeleteProjectRBACRole(ctx context.Context, id, name string) error
	AddProjectUserGroup(ctx context.Context, id, sso, role string) error
	DeleteProjectUserGroup(ctx context.Context, id, sso string) error
}

type webApiEventStore interface {
	List(ctx context.Context, opts datastore.ListOptions) ([]*model.Event, string, error)
}

type webApiAPIKeyStore interface {
	Add(ctx context.Context, k *model.APIKey) error
	List(ctx context.Context, opts datastore.ListOptions) ([]*model.APIKey, error)
	Disable(ctx context.Context, id, projectID string) error
}

// TODO: all structs webApi should be webAPI
type webApiAPIKeyLastUsedStore interface { //nolint:stylecheck
	Get(k string) (interface{}, error)
}

// WebAPI implements the behaviors for the gRPC definitions of WebAPI.
type WebAPI struct {
	webservice.UnimplementedWebServiceServer

	applicationStore          webApiApplicationStore
	deploymentChainStore      webApiDeploymentChainStore
	deploymentStore           webApiDeploymentStore
	pipedStore                webApiPipedStore
	projectStore              webApiProjectStore
	apiKeyStore               webApiAPIKeyStore
	apiKeyLastUsedStore       webApiAPIKeyLastUsedStore
	eventStore                webApiEventStore
	stageLogStore             stagelogstore.Store
	applicationLiveStateStore applicationlivestatestore.Store
	commandStore              commandstore.Store
	insightProvider           insight.Provider
	unregisteredAppStore      unregisteredappstore.Store
	encrypter                 encrypter
	githubCli                 *github.Client

	appProjectCache        cache.Cache
	deploymentProjectCache cache.Cache
	pipedProjectCache      cache.Cache
	pipedStatCache         cache.Cache

	projectsInConfig map[string]config.ControlPlaneProject
	logger           *zap.Logger
}

// NewWebAPI creates a new WebAPI instance.
func NewWebAPI(
	ctx context.Context,
	ds datastore.DataStore,
	sc cache.Cache,
	sls stagelogstore.Store,
	alss applicationlivestatestore.Store,
	uas unregisteredappstore.Store,
	akluc cache.Cache,
	ip insight.Provider,
	psc cache.Cache,
	projs map[string]config.ControlPlaneProject,
	encrypter encrypter,
	logger *zap.Logger,
) *WebAPI {
	w := datastore.WebCommander
	a := &WebAPI{
		applicationStore:          datastore.NewApplicationStore(ds, w),
		deploymentChainStore:      datastore.NewDeploymentChainStore(ds, w),
		deploymentStore:           datastore.NewDeploymentStore(ds, w),
		pipedStore:                datastore.NewPipedStore(ds, w),
		projectStore:              datastore.NewProjectStore(ds, w),
		apiKeyStore:               datastore.NewAPIKeyStore(ds, w),
		apiKeyLastUsedStore:       akluc,
		eventStore:                datastore.NewEventStore(ds, w),
		stageLogStore:             sls,
		applicationLiveStateStore: alss,
		commandStore:              commandstore.NewStore(w, ds, sc, logger),
		insightProvider:           ip,
		unregisteredAppStore:      uas,
		projectsInConfig:          projs,
		encrypter:                 encrypter,
		githubCli:                 github.NewClient(nil),
		appProjectCache:           memorycache.NewTTLCache(ctx, 24*time.Hour, 3*time.Hour),
		deploymentProjectCache:    memorycache.NewTTLCache(ctx, 24*time.Hour, 3*time.Hour),
		pipedProjectCache:         memorycache.NewTTLCache(ctx, 24*time.Hour, 3*time.Hour),
		pipedStatCache:            psc,
		logger:                    logger.Named("web-api"),
	}
	return a
}

// Register registers all handling of this service into the specified gRPC server.
func (a *WebAPI) Register(server *grpc.Server) {
	webservice.RegisterWebServiceServer(server, a)
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
	}
	if err := piped.AddKey(keyHash, claims.Subject, time.Now()); err != nil {
		return nil, status.Error(codes.FailedPrecondition, fmt.Sprintf("Failed to create key: %v", err))
	}

	if err = a.pipedStore.Add(ctx, &piped); err != nil {
		return nil, gRPCStoreError(err, fmt.Sprintf("add piped %s", piped.Id))
	}

	return &webservice.RegisterPipedResponse{
		Id:  piped.Id,
		Key: key,
	}, nil
}

func (a *WebAPI) UpdatePiped(ctx context.Context, req *webservice.UpdatePipedRequest) (*webservice.UpdatePipedResponse, error) {
	updater := func(ctx context.Context, pipedID string) error {
		return a.pipedStore.UpdateInfo(ctx, req.PipedId, req.Name, req.Desc)
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
		return gRPCStoreError(err, fmt.Sprintf("update piped %s", pipedID))
	}
	return nil
}

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

	pipeds, err := a.pipedStore.List(ctx, opts)
	if err != nil {
		return nil, gRPCStoreError(err, "list pipeds")
	}

	// Check piped connection status if necessary.
	// The connection status of piped determined by its submitted stat in pipedStatCache.
	if req.WithStatus {
		for i := range pipeds {
			pipedStatus, err := getPipedStatus(a.pipedStatCache, pipeds[i].Id)
			if err != nil {
				a.logger.Error("failed to get or unmarshal piped stat", zap.Error(err))
				pipedStatus = model.Piped_UNKNOWN
			}
			pipeds[i].Status = pipedStatus
		}
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
		return a.pipedStore.UpdateDesiredVersion(ctx, pipedID, req.Version)
	}
	for _, pipedID := range req.PipedIds {
		if err := a.updatePiped(ctx, pipedID, updater); err != nil {
			return nil, err
		}
	}

	return &webservice.UpdatePipedDesiredVersionResponse{}, nil
}

func (a *WebAPI) RestartPiped(ctx context.Context, req *webservice.RestartPipedRequest) (*webservice.RestartPipedResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	piped, err := getPiped(ctx, a.pipedStore, req.PipedId, a.logger)
	if err != nil {
		return nil, err
	}

	if claims.Role.ProjectId != piped.ProjectId {
		return nil, status.Error(codes.PermissionDenied, "Requested Piped does not belong to your project")
	}

	cmd := model.Command{
		Id:        uuid.New().String(),
		PipedId:   piped.Id,
		ProjectId: piped.ProjectId,
		Type:      model.Command_RESTART_PIPED,
		Commander: claims.Subject,
		RestartPiped: &model.Command_RestartPiped{
			PipedId: piped.Id,
		},
	}
	if err := addCommand(ctx, a.commandStore, &cmd, a.logger); err != nil {
		return nil, err
	}

	return &webservice.RestartPipedResponse{
		CommandId: cmd.Id,
	}, nil
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

func (a *WebAPI) ListUnregisteredApplications(ctx context.Context, _ *webservice.ListUnregisteredApplicationsRequest) (*webservice.ListUnregisteredApplicationsResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	allApps, err := a.unregisteredAppStore.ListApplications(ctx, claims.Role.ProjectId)
	if errors.Is(err, cache.ErrNotFound) {
		return &webservice.ListUnregisteredApplicationsResponse{}, nil
	}
	if err != nil {
		a.logger.Error("failed to get unregistered apps", zap.Error(err))
		return nil, gRPCStoreError(err, "get unregistered apps")
	}

	if len(allApps) == 0 {
		return &webservice.ListUnregisteredApplicationsResponse{}, nil
	}

	sort.Slice(allApps, func(i, j int) bool {
		return allApps[i].Path < allApps[j].Path
	})

	return &webservice.ListUnregisteredApplicationsResponse{
		Applications: allApps,
	}, nil
}

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
		return nil, status.Error(codes.PermissionDenied, "Requested piped does not belong to your project")
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
		Id:               uuid.New().String(),
		Name:             req.Name,
		PipedId:          req.PipedId,
		ProjectId:        claims.Role.ProjectId,
		GitPath:          gitpath,
		Kind:             req.Kind,
		PlatformProvider: req.PlatformProvider,
		CloudProvider:    req.PlatformProvider,
		Description:      req.Description,
		Labels:           req.Labels,
	}
	if err = a.applicationStore.Add(ctx, &app); err != nil {
		return nil, gRPCStoreError(err, fmt.Sprintf("add application %s", app.Id))
	}

	return &webservice.AddApplicationResponse{
		ApplicationId: app.Id,
	}, nil
}

func (a *WebAPI) UpdateApplication(ctx context.Context, req *webservice.UpdateApplicationRequest) (*webservice.UpdateApplicationResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	piped, err := getPiped(ctx, a.pipedStore, req.PipedId, a.logger)
	if err != nil {
		return nil, gRPCStoreError(err, fmt.Sprintf("failed to get piped %s", req.PipedId))
	}

	if piped.ProjectId != claims.Role.ProjectId {
		return nil, status.Error(codes.PermissionDenied, "Requested piped does not belong to your project")
	}

	if err := a.applicationStore.UpdateConfiguration(ctx, req.ApplicationId, req.PipedId, req.PlatformProvider, req.ConfigFilename); err != nil {
		return nil, gRPCStoreError(err, fmt.Sprintf("failed to update application %s", req.ApplicationId))
	}

	return &webservice.UpdateApplicationResponse{}, nil
}

func (a *WebAPI) EnableApplication(ctx context.Context, req *webservice.EnableApplicationRequest) (*webservice.EnableApplicationResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	if err := a.validateAppBelongsToProject(ctx, req.ApplicationId, claims.Role.ProjectId); err != nil {
		return nil, err
	}

	if err := a.applicationStore.Enable(ctx, req.ApplicationId); err != nil {
		return nil, gRPCStoreError(err, fmt.Sprintf("enable application %s", req.ApplicationId))
	}
	return &webservice.EnableApplicationResponse{}, nil
}

func (a *WebAPI) DisableApplication(ctx context.Context, req *webservice.DisableApplicationRequest) (*webservice.DisableApplicationResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	if err := a.validateAppBelongsToProject(ctx, req.ApplicationId, claims.Role.ProjectId); err != nil {
		return nil, err
	}

	if err := a.applicationStore.Disable(ctx, req.ApplicationId); err != nil {
		return nil, gRPCStoreError(err, fmt.Sprintf("disable application %s", req.ApplicationId))
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

	if err := a.applicationStore.Delete(ctx, req.ApplicationId); err != nil {
		return nil, gRPCStoreError(err, fmt.Sprintf("delete application %s", req.ApplicationId))
	}

	return &webservice.DeleteApplicationResponse{}, nil
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
		if o.Name != "" {
			filters = append(filters, datastore.ListFilter{
				Field:    "Name",
				Operator: datastore.OperatorEqual,
				Value:    o.Name,
			})
		}
	}

	apps, _, err := a.applicationStore.List(ctx, datastore.ListOptions{
		Filters: filters,
		Orders:  orders,
	})
	if err != nil {
		return nil, gRPCStoreError(err, "list applications")
	}

	if len(req.Options.Labels) == 0 {
		return &webservice.ListApplicationsResponse{
			Applications: apps,
		}, nil
	}

	// NOTE: Filtering by labels is done by the application-side because we need to create composite indexes for every combination in the filter.
	filtered := make([]*model.Application, 0, len(apps))
	for _, a := range apps {
		if a.ContainLabels(req.Options.Labels) {
			filtered = append(filtered, a)
		}
	}
	return &webservice.ListApplicationsResponse{
		Applications: filtered,
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
		return nil, status.Error(codes.PermissionDenied, "Requested application does not belong to your project")
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
		return nil, status.Error(codes.PermissionDenied, "Requested application does not belong to your project")
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

	pubkey, err := getEncriptionKey(piped.SecretEncryption)
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
		{
			Field:    "UpdatedAt",
			Operator: datastore.OperatorGreaterThanOrEqual,
			Value:    req.PageMinUpdatedAt,
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
		if o.ApplicationName != "" {
			filters = append(filters, datastore.ListFilter{
				Field:    "ApplicationName",
				Operator: datastore.OperatorEqual,
				Value:    o.ApplicationName,
			})
		}
	}

	pageSize := int(req.PageSize)
	options := datastore.ListOptions{
		Filters: filters,
		Orders:  orders,
		Limit:   pageSize,
		Cursor:  req.Cursor,
	}
	deployments, cursor, err := a.deploymentStore.List(ctx, options)
	if err != nil {
		a.logger.Error("failed to get deployments", zap.Error(err))
		return nil, gRPCStoreError(err, "get deployments")
	}
	labels := req.Options.Labels
	if len(labels) == 0 || len(deployments) == 0 {
		return &webservice.ListDeploymentsResponse{
			Deployments: deployments,
			Cursor:      cursor,
		}, nil
	}

	// Start filtering them by labels.
	//
	// NOTE: Filtering by labels is done by the application-side because we need to create composite indexes for every combination in the filter.
	// We don't want to depend on any other search engine, that's why it filters here.
	filtered := make([]*model.Deployment, 0, len(deployments))
	for _, d := range deployments {
		if d.ContainLabels(labels) {
			filtered = append(filtered, d)
		}
	}
	// Stop running additional queries for more data, and return filtered deployments immediately with
	// current cursor if the size before filtering is already less than the page size.
	if len(deployments) < pageSize {
		return &webservice.ListDeploymentsResponse{
			Deployments: filtered,
			Cursor:      cursor,
		}, nil
	}
	// Repeat the query until the number of filtered deployments reaches the page size,
	// or until it finishes scanning to page_min_updated_at.
	for len(filtered) < pageSize {
		options.Cursor = cursor
		deployments, cursor, err = a.deploymentStore.List(ctx, options)
		if err != nil {
			a.logger.Error("failed to get deployments", zap.Error(err))
			return nil, gRPCStoreError(err, "get deployments")
		}
		if len(deployments) == 0 {
			break
		}
		for _, d := range deployments {
			if d.ContainLabels(labels) {
				filtered = append(filtered, d)
			}
		}
		// We've already specified UpdatedAt >= req.PageMinUpdatedAt, so we need to check just equality.
		if deployments[len(deployments)-1].UpdatedAt == req.PageMinUpdatedAt {
			break
		}
	}
	// TODO: Think about possibility that the response of ListDeployments exceeds the page size
	return &webservice.ListDeploymentsResponse{
		Deployments: filtered,
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
		return nil, status.Error(codes.PermissionDenied, "Requested deployment does not belong to your project")
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
	if err != nil {
		a.logger.Error("failed to get stage logs", zap.Error(err))
		return nil, gRPCStoreError(err, "get stage logs")
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
		return nil, status.Error(codes.PermissionDenied, "Requested deployment does not belong to your project")
	}

	if deployment.Status.IsCompleted() {
		return nil, status.Error(codes.FailedPrecondition, "could not cancel the deployment because it was already completed")
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

func (a *WebAPI) SkipStage(ctx context.Context, req *webservice.SkipStageRequest) (*webservice.SkipStageResponse, error) {
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
		return nil, status.Error(codes.PermissionDenied, "Requested deployment does not belong to your project")
	}
	stage, ok := deployment.Stage(req.StageId)
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, "The stage was not found in the deployment")
	}
	if !stage.IsSkippable() {
		return nil, status.Error(codes.FailedPrecondition, fmt.Sprintf("SKIP STAGE is not supported for stage %q", stage.Name))
	}
	if stage.Status.IsCompleted() {
		return nil, status.Error(codes.FailedPrecondition, "Could not skip the stage because it was already completed")
	}

	commandID := uuid.New().String()
	cmd := model.Command{
		Id:            commandID,
		PipedId:       deployment.PipedId,
		ApplicationId: deployment.ApplicationId,
		ProjectId:     deployment.ProjectId,
		DeploymentId:  req.DeploymentId,
		StageId:       req.StageId,
		Type:          model.Command_SKIP_STAGE,
		Commander:     claims.Subject,
		SkipStage: &model.Command_SkipStage{
			DeploymentId: req.DeploymentId,
			StageId:      req.StageId,
		},
	}
	if err := addCommand(ctx, a.commandStore, &cmd, a.logger); err != nil {
		return nil, err
	}

	return &webservice.SkipStageResponse{
		CommandId: commandID,
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
	if err := validateApprover(deployment.Stages, claims.Subject, req.StageId); err != nil {
		return nil, err
	}
	if err := a.validateDeploymentBelongsToProject(ctx, req.DeploymentId, claims.Role.ProjectId); err != nil {
		return nil, err
	}
	stage, ok := deployment.StageMap()[req.StageId]
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, "The stage was not found in the deployment")
	}
	if stage.Status.IsCompleted() {
		return nil, status.Error(codes.FailedPrecondition, "Could not approve the stage because it was already completed")
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
func validateApprover(stages []*model.PipelineStage, commander, stageID string) error {
	var approvers []string
	for _, s := range stages {
		if s.Id != stageID {
			continue
		}
		if as := s.Metadata["Approvers"]; as != "" {
			approvers = strings.Split(as, ",")
		}
		break
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
	return status.Error(codes.PermissionDenied, fmt.Sprintf("You can't approve this deployment because you (%s) are not in the approver list: %v", commander, approvers))
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
	if err != nil {
		a.logger.Error("failed to get application live state", zap.Error(err))
		return nil, gRPCStoreError(err, "get application live state")
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
		p := &model.Project{
			Id:   p.Id,
			Desc: p.Desc,
			StaticAdmin: &model.ProjectStaticUser{
				Username:     p.StaticAdmin.Username,
				PasswordHash: p.StaticAdmin.PasswordHash,
			},
		}
		p.SetBuiltinRBACRoles()
		return p, nil
	}

	project, err := a.projectStore.Get(ctx, projectID)
	if err != nil {
		return nil, gRPCStoreError(err, fmt.Sprintf("get project %s", projectID))
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
		return nil, gRPCStoreError(err, "update static admin")
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
		return nil, gRPCStoreError(err, "enable static admin login")
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
		a.logger.Error("failed to disable static admin login", zap.Error(err))
		return nil, gRPCStoreError(err, "disable static admin login")
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
		return nil, gRPCStoreError(err, "encrypt sensitive data in sso configurations")
	}

	if err := a.projectStore.UpdateProjectSSOConfig(ctx, claims.Role.ProjectId, req.Sso); err != nil {
		a.logger.Error("failed to update project single sign on settings", zap.Error(err))
		return nil, gRPCStoreError(err, "update project single sign on settings")
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
		return nil, gRPCStoreError(err, "update project single sign on settings")
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
		Subject:   claims.Subject,
		AvatarUrl: claims.AvatarURL,
		ProjectId: claims.Role.ProjectId,
	}, nil
}

func (a *WebAPI) AddProjectRBACRole(ctx context.Context, req *webservice.AddProjectRBACRoleRequest) (*webservice.AddProjectRBACRoleResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	if _, ok := a.projectsInConfig[claims.Role.ProjectId]; ok {
		return nil, status.Error(codes.FailedPrecondition, "Failed to update a debug project specified in the control-plane configuration")
	}

	if err := a.projectStore.AddProjectRBACRole(ctx, claims.Role.ProjectId, req.Name, req.Policies); err != nil {
		a.logger.Error("failed to add rbac role", zap.Error(err))
		return nil, gRPCStoreError(err, "add rbac role")
	}
	return &webservice.AddProjectRBACRoleResponse{}, nil
}

func (a *WebAPI) UpdateProjectRBACRole(ctx context.Context, req *webservice.UpdateProjectRBACRoleRequest) (*webservice.UpdateProjectRBACRoleResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	if _, ok := a.projectsInConfig[claims.Role.ProjectId]; ok {
		return nil, status.Error(codes.FailedPrecondition, "Failed to update a debug project specified in the control-plane configuration")
	}

	if err := a.projectStore.UpdateProjectRBACRole(ctx, claims.Role.ProjectId, req.Name, req.Policies); err != nil {
		a.logger.Error("failed to update rbac role", zap.Error(err))
		return nil, gRPCStoreError(err, "update rbac role")
	}
	return &webservice.UpdateProjectRBACRoleResponse{}, nil
}

func (a *WebAPI) DeleteProjectRBACRole(ctx context.Context, req *webservice.DeleteProjectRBACRoleRequest) (*webservice.DeleteProjectRBACRoleResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	if _, ok := a.projectsInConfig[claims.Role.ProjectId]; ok {
		return nil, status.Error(codes.FailedPrecondition, "Failed to update a debug project specified in the control-plane configuration")
	}

	if err := a.projectStore.DeleteProjectRBACRole(ctx, claims.Role.ProjectId, req.Name); err != nil {
		a.logger.Error("failed to delete rbac role", zap.Error(err))
		return nil, gRPCStoreError(err, "delete rbac role")
	}
	return &webservice.DeleteProjectRBACRoleResponse{}, nil
}

func (a *WebAPI) AddProjectUserGroup(ctx context.Context, req *webservice.AddProjectUserGroupRequest) (*webservice.AddProjectUserGroupResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	if _, ok := a.projectsInConfig[claims.Role.ProjectId]; ok {
		return nil, status.Error(codes.FailedPrecondition, "Failed to update a debug project specified in the control-plane configuration")
	}

	if err := a.projectStore.AddProjectUserGroup(ctx, claims.Role.ProjectId, req.SsoGroup, req.Role); err != nil {
		a.logger.Error("failed to add user group", zap.Error(err))
		return nil, gRPCStoreError(err, "add user group")
	}
	return &webservice.AddProjectUserGroupResponse{}, nil
}

func (a *WebAPI) DeleteProjectUserGroup(ctx context.Context, req *webservice.DeleteProjectUserGroupRequest) (*webservice.DeleteProjectUserGroupResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	if _, ok := a.projectsInConfig[claims.Role.ProjectId]; ok {
		return nil, status.Error(codes.FailedPrecondition, "Failed to update a debug project specified in the control-plane configuration")
	}

	if err := a.projectStore.DeleteProjectUserGroup(ctx, claims.Role.ProjectId, req.SsoGroup); err != nil {
		a.logger.Error("failed to delete user group", zap.Error(err))
		return nil, gRPCStoreError(err, "delete user group")
	}
	return &webservice.DeleteProjectUserGroupResponse{}, nil
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
		return nil, status.Error(codes.PermissionDenied, "Requested command does not belong to your project")
	}

	return &webservice.GetCommandResponse{
		Command: cmd,
	}, nil
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
		Id:         id,
		Name:       req.Name,
		KeyHash:    hash,
		ProjectId:  claims.Role.ProjectId,
		Role:       req.Role,
		Creator:    claims.Subject,
		LastUsedAt: 0,
	}

	if err = a.apiKeyStore.Add(ctx, &apiKey); err != nil {
		return nil, gRPCStoreError(err, fmt.Sprintf("add API key %s", apiKey.Id))
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

	if err := a.apiKeyStore.Disable(ctx, req.Id, claims.Role.ProjectId); err != nil {
		return nil, gRPCStoreError(err, fmt.Sprintf("disable API key %s", req.Id))
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

	apiKeys, err := a.apiKeyStore.List(ctx, opts)
	if err != nil {
		return nil, gRPCStoreError(err, "list API keys")
	}

	for i := range apiKeys {
		// Get LastUsedTIme from Redis
		if lastUsedAt, error := a.apiKeyLastUsedStore.Get(apiKeys[i].Id); error == nil {
			apiKeys[i].LastUsedAt = bytes2int64(lastUsedAt.([]byte))
		} else {
			apiKeys[i].LastUsedAt = 0
		}
		// Redact all sensitive data inside API key before sending to the client.
		apiKeys[i].RedactSensitiveData()
	}

	return &webservice.ListAPIKeysResponse{
		Keys: apiKeys,
	}, nil
}

func bytes2int64(bytes []byte) int64 {
	var numString string
	for i := range bytes {
		numString += string(bytes[i])
	}
	num, _ := strconv.ParseInt(numString, 10, 64)
	return num
}

// GetInsightData returns the accumulated insight data.
func (a *WebAPI) GetInsightData(ctx context.Context, req *webservice.GetInsightDataRequest) (*webservice.GetInsightDataResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	var points []*model.InsightDataPoint

	switch req.MetricsKind {
	case model.InsightMetricsKind_DEPLOYMENT_FREQUENCY:
		points, err = a.insightProvider.GetDeploymentFrequencyDataPoints(
			ctx,
			claims.Role.ProjectId,
			req.ApplicationId,
			req.Labels,
			req.RangeFrom,
			req.RangeTo,
			req.Resolution,
		)

	case model.InsightMetricsKind_CHANGE_FAILURE_RATE:
		points, err = a.insightProvider.GetDeploymentChangeFailureRateDataPoints(
			ctx,
			claims.Role.ProjectId,
			req.ApplicationId,
			req.Labels,
			req.RangeFrom,
			req.RangeTo,
			req.Resolution,
		)

	default:
		return nil, status.Error(codes.Unimplemented, fmt.Sprintf("The insight metrics %s is not implemented yet", req.MetricsKind.String()))
	}

	if err != nil {
		a.logger.Error("failed to get insight data", zap.Error(err))
		return nil, gRPCStoreError(err, "get insight data")
	}

	return &webservice.GetInsightDataResponse{
		Type: model.InsightResultType_MATRIX,
		Matrix: []*model.InsightSampleStream{
			{DataPoints: points},
		},
	}, nil
}

func (a *WebAPI) GetInsightApplicationCount(ctx context.Context, req *webservice.GetInsightApplicationCountRequest) (*webservice.GetInsightApplicationCountResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	c, err := a.insightProvider.GetApplicationCounts(ctx, claims.Role.ProjectId)
	if err != nil {
		a.logger.Error("failed to load application counts", zap.Error(err))
		return nil, gRPCStoreError(err, "load application counts")
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

func (a *WebAPI) ListDeploymentChains(ctx context.Context, req *webservice.ListDeploymentChainsRequest) (*webservice.ListDeploymentChainsResponse, error) {
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
		{
			Field:    "UpdatedAt",
			Operator: datastore.OperatorGreaterThan,
			Value:    req.PageMinUpdatedAt,
		},
	}
	// TODO: Support filter list deployment chain with options.

	pageSize := int(req.PageSize)
	options := datastore.ListOptions{
		Filters: filters,
		Orders:  orders,
		Limit:   pageSize,
		Cursor:  req.Cursor,
	}

	deploymentChains, cursor, err := a.deploymentChainStore.List(ctx, options)
	if err != nil {
		return nil, gRPCStoreError(err, "list deployment chains")
	}

	return &webservice.ListDeploymentChainsResponse{
		DeploymentChains: deploymentChains,
		Cursor:           cursor,
	}, nil
}

func (a *WebAPI) GetDeploymentChain(ctx context.Context, req *webservice.GetDeploymentChainRequest) (*webservice.GetDeploymentChainResponse, error) {
	claims, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	dc, err := a.deploymentChainStore.Get(ctx, req.DeploymentChainId)
	if err != nil {
		return nil, gRPCStoreError(err, fmt.Sprintf("get deployment chain %s", req.DeploymentChainId))
	}

	if claims.Role.ProjectId != dc.ProjectId {
		return nil, status.Error(codes.PermissionDenied, "Requested deployment chain does not belong to your project")
	}

	return &webservice.GetDeploymentChainResponse{
		DeploymentChain: dc,
	}, nil
}

func (a *WebAPI) ListEvents(ctx context.Context, req *webservice.ListEventsRequest) (*webservice.ListEventsResponse, error) {
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
		{
			Field:    "UpdatedAt",
			Operator: datastore.OperatorGreaterThanOrEqual,
			Value:    req.PageMinUpdatedAt,
		},
	}
	if o := req.Options; o != nil {
		if o.Name != "" {
			filters = append(filters, datastore.ListFilter{
				Field:    "Name",
				Operator: datastore.OperatorEqual,
				Value:    o.Name,
			})
		}
		// Allowing multiple so that it can do In Query later.
		// Currently only the first value is used.
		if len(o.Statuses) > 0 {
			filters = append(filters, datastore.ListFilter{
				Field:    "Status",
				Operator: datastore.OperatorEqual,
				Value:    o.Statuses[0],
			})
		}
	}

	pageSize := int(req.PageSize)
	options := datastore.ListOptions{
		Filters: filters,
		Orders:  orders,
		Limit:   pageSize,
		Cursor:  req.Cursor,
	}
	events, cursor, err := a.eventStore.List(ctx, options)
	if err != nil {
		return nil, gRPCStoreError(err, "list events")
	}

	labels := req.Options.Labels
	if len(labels) == 0 || len(events) == 0 {
		return &webservice.ListEventsResponse{
			Events: events,
			Cursor: cursor,
		}, nil
	}

	// Start filtering them by labels.
	//
	// NOTE: Filtering by labels is done by the application-side because we need to create composite indexes for every combination in the filter.
	// We don't want to depend on any other search engine, that's why it filters here.
	filtered := make([]*model.Event, 0, len(events))
	for _, e := range events {
		if e.ContainLabels(labels) {
			filtered = append(filtered, e)
		}
	}
	// Stop running additional queries for more data, and return filtered events immediately with
	// current cursor if the size before filtering is already less than the page size.
	if len(events) < pageSize {
		return &webservice.ListEventsResponse{
			Events: filtered,
			Cursor: cursor,
		}, nil
	}
	// Repeat the query until the number of filtered events reaches the page size,
	// or until it finishes scanning to page_min_updated_at.
	for len(filtered) < pageSize {
		options.Cursor = cursor
		events, cursor, err = a.eventStore.List(ctx, options)
		if err != nil {
			a.logger.Error("failed to get events", zap.Error(err))
			return nil, gRPCStoreError(err, "Failed to get events")
		}
		if len(events) == 0 {
			break
		}
		for _, e := range events {
			if e.ContainLabels(labels) {
				filtered = append(filtered, e)
			}
		}
		// We've already specified UpdatedAt >= req.PageMinUpdatedAt, so we need to check just equality.
		if events[len(events)-1].UpdatedAt == req.PageMinUpdatedAt {
			break
		}
	}
	// TODO: Think about possibility that the response of ListEvents exceeds the page size
	return &webservice.ListEventsResponse{
		Events: filtered,
		Cursor: cursor,
	}, nil
}

func (a *WebAPI) ListReleasedVersions(ctx context.Context, req *webservice.ListReleasedVersionsRequest) (*webservice.ListReleasedVersionsResponse, error) {
	_, err := rpcauth.ExtractClaims(ctx)
	if err != nil {
		a.logger.Error("failed to authenticate the current user", zap.Error(err))
		return nil, err
	}

	releases, _, err := a.githubCli.Repositories.ListReleases(ctx, "pipe-cd", "pipecd", nil)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to list released versions")
	}

	if len(releases) == 0 {
		return &webservice.ListReleasedVersionsResponse{}, nil
	}

	versions := make([]string, 0, len(releases))
	for _, release := range releases {
		// Ignore pre-release tagged or draft release.
		if *release.Prerelease || *release.Draft {
			continue
		}
		versions = append(versions, *release.TagName)
	}

	return &webservice.ListReleasedVersionsResponse{
		Versions: versions,
	}, nil
}
