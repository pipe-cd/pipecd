// Copyright 2022 The PipeCD Authors.
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
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/app/server/commandstore"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/app/server/stagelogstore"
	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/cache/memorycache"
	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/filestore"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcauth"
)

type apiApplicationStore interface {
	Add(ctx context.Context, app *model.Application) error
	Get(ctx context.Context, id string) (*model.Application, error)
	List(ctx context.Context, opts datastore.ListOptions) ([]*model.Application, string, error)
	Delete(ctx context.Context, id string) error
	Disable(ctx context.Context, id string) error
	UpdateConfigFilename(ctx context.Context, id, filename string) error
}

type apiDeploymentStore interface {
	Get(ctx context.Context, id string) (*model.Deployment, error)
}

type apiPipedStore interface {
	Get(ctx context.Context, id string) (*model.Piped, error)
	List(ctx context.Context, opts datastore.ListOptions) ([]*model.Piped, error)
	EnablePiped(ctx context.Context, id string) error
	DisablePiped(ctx context.Context, id string) error
}

type apiEventStore interface {
	Add(ctx context.Context, event model.Event) error
}

type commandOutputGetter interface {
	Get(ctx context.Context, commandID string) ([]byte, error)
}

// API implements the behaviors for the gRPC definitions of API.
type API struct {
	apiservice.UnimplementedAPIServiceServer

	applicationStore    apiApplicationStore
	deploymentStore     apiDeploymentStore
	pipedStore          apiPipedStore
	eventStore          apiEventStore
	commandStore        commandstore.Store
	stageLogStore       stagelogstore.Store
	commandOutputGetter commandOutputGetter

	encryptionKeyCache cache.Cache
	pipedStatCache     cache.Cache

	webBaseURL string
	logger     *zap.Logger
}

// NewAPI creates a new API instance.
func NewAPI(
	ctx context.Context,
	ds datastore.DataStore,
	fs filestore.Store,
	sc cache.Cache,
	cog commandOutputGetter,
	psc cache.Cache,
	webBaseURL string,
	logger *zap.Logger,
) *API {
	w := datastore.PipectlCommander
	a := &API{
		applicationStore:    datastore.NewApplicationStore(ds, w),
		deploymentStore:     datastore.NewDeploymentStore(ds, w),
		pipedStore:          datastore.NewPipedStore(ds, w),
		eventStore:          datastore.NewEventStore(ds, w),
		commandStore:        commandstore.NewStore(w, ds, sc, logger),
		stageLogStore:       stagelogstore.NewStore(fs, sc, logger),
		commandOutputGetter: cog,
		// Public key is variable but likely to be accessed multiple times in a short period.
		encryptionKeyCache: memorycache.NewTTLCache(ctx, 5*time.Minute, 5*time.Minute),
		pipedStatCache:     psc,
		webBaseURL:         webBaseURL,
		logger:             logger.Named("api"),
	}
	return a
}

// Register registers all handling of this service into the specified gRPC server.
func (a *API) Register(server *grpc.Server) {
	apiservice.RegisterAPIServiceServer(server, a)
}

func (a *API) AddApplication(ctx context.Context, req *apiservice.AddApplicationRequest) (*apiservice.AddApplicationResponse, error) {
	key, err := requireAPIKey(ctx, model.APIKey_READ_WRITE, a.logger)
	if err != nil {
		return nil, err
	}

	piped, err := getPiped(ctx, a.pipedStore, req.PipedId, a.logger)
	if err != nil {
		return nil, err
	}

	if key.ProjectId != piped.ProjectId {
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
		Id:               uuid.New().String(),
		Name:             req.Name,
		PipedId:          req.PipedId,
		ProjectId:        key.ProjectId,
		GitPath:          gitpath,
		Kind:             req.Kind,
		PlatformProvider: req.PlatformProvider,
		CloudProvider:    req.PlatformProvider,
		Description:      req.Description,
	}
	if err := a.applicationStore.Add(ctx, &app); err != nil {
		return nil, gRPCStoreError(err, fmt.Sprintf("add application %s", app.Id))
	}

	return &apiservice.AddApplicationResponse{
		ApplicationId: app.Id,
	}, nil
}

func (a *API) SyncApplication(ctx context.Context, req *apiservice.SyncApplicationRequest) (*apiservice.SyncApplicationResponse, error) {
	key, err := requireAPIKey(ctx, model.APIKey_READ_WRITE, a.logger)
	if err != nil {
		return nil, err
	}

	app, err := getApplication(ctx, a.applicationStore, req.ApplicationId, a.logger)
	if err != nil {
		return nil, err
	}

	if key.ProjectId != app.ProjectId {
		return nil, status.Error(codes.InvalidArgument, "Requested application does not belong to your project")
	}

	cmd := model.Command{
		Id:            uuid.New().String(),
		PipedId:       app.PipedId,
		ApplicationId: app.Id,
		ProjectId:     app.ProjectId,
		Type:          model.Command_SYNC_APPLICATION,
		Commander:     key.Id,
		SyncApplication: &model.Command_SyncApplication{
			ApplicationId: app.Id,
			SyncStrategy:  model.SyncStrategy_AUTO,
		},
	}
	if err := addCommand(ctx, a.commandStore, &cmd, a.logger); err != nil {
		return nil, err
	}

	return &apiservice.SyncApplicationResponse{
		CommandId: cmd.Id,
	}, nil
}

func (a *API) GetApplication(ctx context.Context, req *apiservice.GetApplicationRequest) (*apiservice.GetApplicationResponse, error) {
	key, err := requireAPIKey(ctx, model.APIKey_READ_ONLY, a.logger)
	if err != nil {
		return nil, err
	}

	app, err := getApplication(ctx, a.applicationStore, req.ApplicationId, a.logger)
	if err != nil {
		return nil, err
	}

	if app.ProjectId != key.ProjectId {
		return nil, status.Error(codes.InvalidArgument, "Requested application does not belong to your project")
	}

	return &apiservice.GetApplicationResponse{
		Application: app,
	}, nil
}

// ListApplications returns the application list of the project where the caller belongs to.
// Currently, the maximum number of returned applications per request is set to 10.
// The response contains a "cursor" value, which should be passed in the next request in order to get
// the next 10 applications. If the cursor is not provided in the request, only 10 latest applications will be returned.
func (a *API) ListApplications(ctx context.Context, req *apiservice.ListApplicationsRequest) (*apiservice.ListApplicationsResponse, error) {
	key, err := requireAPIKey(ctx, model.APIKey_READ_ONLY, a.logger)
	if err != nil {
		return nil, err
	}

	const limit = 10
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
			Value:    key.ProjectId,
		},
		{
			Field:    "Disabled",
			Operator: datastore.OperatorEqual,
			Value:    req.Disabled,
		},
	}

	if req.Name != "" {
		filters = append(filters, datastore.ListFilter{
			Field:    "Name",
			Operator: datastore.OperatorEqual,
			Value:    req.Name,
		})
	}
	if req.Kind != "" {
		kind, ok := model.ApplicationKind_value[req.Kind]
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "Invalid application kind")
		}
		filters = append(filters, datastore.ListFilter{
			Field:    "Kind",
			Operator: datastore.OperatorEqual,
			Value:    model.ApplicationKind(kind),
		})
	}
	opts := datastore.ListOptions{
		Orders:  orders,
		Filters: filters,
		Limit:   limit,
		Cursor:  req.Cursor,
	}

	apps, cursor, err := a.applicationStore.List(ctx, opts)
	if err != nil {
		return nil, gRPCStoreError(err, "failed to list applications")
	}

	return &apiservice.ListApplicationsResponse{
		Applications: apps,
		Cursor:       cursor,
	}, nil
}

func (a *API) DeleteApplication(ctx context.Context, req *apiservice.DeleteApplicationRequest) (*apiservice.DeleteApplicationResponse, error) {
	key, err := requireAPIKey(ctx, model.APIKey_READ_WRITE, a.logger)
	if err != nil {
		return nil, err
	}

	app, err := getApplication(ctx, a.applicationStore, req.ApplicationId, a.logger)
	if err != nil {
		return nil, err
	}

	if app.ProjectId != key.ProjectId {
		return nil, status.Error(codes.InvalidArgument, "Requested application does not belong to your project")
	}

	if err := a.applicationStore.Delete(ctx, req.ApplicationId); err != nil {
		return nil, gRPCStoreError(err, fmt.Sprintf("delete application %s", app.Id))
	}

	return &apiservice.DeleteApplicationResponse{
		ApplicationId: app.Id,
	}, nil
}

func (a *API) DisableApplication(ctx context.Context, req *apiservice.DisableApplicationRequest) (*apiservice.DisableApplicationResponse, error) {
	key, err := requireAPIKey(ctx, model.APIKey_READ_WRITE, a.logger)
	if err != nil {
		return nil, err
	}

	app, err := getApplication(ctx, a.applicationStore, req.ApplicationId, a.logger)
	if err != nil {
		return nil, err
	}

	if app.ProjectId != key.ProjectId {
		return nil, status.Error(codes.InvalidArgument, "Requested application does not belong to your project")
	}

	if err := a.applicationStore.Disable(ctx, req.ApplicationId); err != nil {
		return nil, gRPCStoreError(err, fmt.Sprintf("disable application %s", req.ApplicationId))
	}

	return &apiservice.DisableApplicationResponse{
		ApplicationId: app.Id,
	}, nil
}

func (a *API) RenameApplicationConfigFile(ctx context.Context, req *apiservice.RenameApplicationConfigFileRequest) (*apiservice.RenameApplicationConfigFileResponse, error) {
	key, err := requireAPIKey(ctx, model.APIKey_READ_WRITE, a.logger)
	if err != nil {
		return nil, err
	}

	for _, appID := range req.ApplicationIds {
		app, err := a.applicationStore.Get(ctx, appID)
		if err != nil {
			return nil, gRPCStoreError(err, fmt.Sprintf("failed to get application %s", appID))
		}
		if app.ProjectId != key.ProjectId {
			return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("requested application %s does not belong to your project", appID))
		}
		if err = a.applicationStore.UpdateConfigFilename(ctx, appID, req.NewFilename); err != nil {
			return nil, gRPCStoreError(err, fmt.Sprintf("failed to update application %s config file name", appID))
		}
	}

	return &apiservice.RenameApplicationConfigFileResponse{}, nil
}

func (a *API) GetDeployment(ctx context.Context, req *apiservice.GetDeploymentRequest) (*apiservice.GetDeploymentResponse, error) {
	key, err := requireAPIKey(ctx, model.APIKey_READ_ONLY, a.logger)
	if err != nil {
		return nil, err
	}

	deployment, err := getDeployment(ctx, a.deploymentStore, req.DeploymentId, a.logger)
	if err != nil {
		return nil, err
	}

	if key.ProjectId != deployment.ProjectId {
		return nil, status.Error(codes.InvalidArgument, "Requested deployment does not belong to your project")
	}

	return &apiservice.GetDeploymentResponse{
		Deployment: deployment,
	}, nil
}

func (a *API) ListStageLogs(ctx context.Context, req *apiservice.ListStageLogsRequest) (*apiservice.ListStageLogsResponse, error) {
	key, err := requireAPIKey(ctx, model.APIKey_READ_ONLY, a.logger)
	if err != nil {
		return nil, err
	}

	deployment, err := getDeployment(ctx, a.deploymentStore, req.DeploymentId, a.logger)
	if err != nil {
		return nil, err
	}

	if key.ProjectId != deployment.ProjectId {
		return nil, status.Error(codes.InvalidArgument, "Requested deployment does not belong to your project")
	}

	stageLogs := map[string]*apiservice.StageLog{}

	for _, stage := range deployment.Stages {
		blocks, completed, err := a.stageLogStore.FetchLogs(ctx, deployment.Id, stage.Id, stage.RetriedCount, 0)
		if err != nil && !errors.Is(err, stagelogstore.ErrNotFound) {
			return nil, err
		}

		if err != nil && stage.Name == model.StageRollback.String() {
			// Rollback is generated automatically and returns nothing if not found
			if stage.Name == model.StageRollback.String() {
				continue
			} else {
				stageLogs[stage.Id] = &apiservice.StageLog{}
			}
		}

		stageLogs[stage.Id] = &apiservice.StageLog{
			Blocks:    blocks,
			Completed: completed,
		}
	}

	return &apiservice.ListStageLogsResponse{
		StageLogs: stageLogs,
	}, nil
}

func (a *API) GetCommand(ctx context.Context, req *apiservice.GetCommandRequest) (*apiservice.GetCommandResponse, error) {
	_, err := requireAPIKey(ctx, model.APIKey_READ_ONLY, a.logger)
	if err != nil {
		return nil, err
	}

	cmd, err := getCommand(ctx, a.commandStore, req.CommandId, a.logger)
	if err != nil {
		return nil, err
	}

	return &apiservice.GetCommandResponse{
		Command: cmd,
	}, nil
}

func (a *API) EnablePiped(ctx context.Context, req *apiservice.EnablePipedRequest) (*apiservice.EnablePipedResponse, error) {
	if err := a.updatePiped(ctx, req.PipedId, a.pipedStore.EnablePiped); err != nil {
		return nil, err
	}
	return &apiservice.EnablePipedResponse{}, nil
}

func (a *API) DisablePiped(ctx context.Context, req *apiservice.DisablePipedRequest) (*apiservice.DisablePipedResponse, error) {
	if err := a.updatePiped(ctx, req.PipedId, a.pipedStore.DisablePiped); err != nil {
		return nil, err
	}
	return &apiservice.DisablePipedResponse{}, nil
}

func (a *API) updatePiped(ctx context.Context, pipedID string, updater func(context.Context, string) error) error {
	key, err := requireAPIKey(ctx, model.APIKey_READ_WRITE, a.logger)
	if err != nil {
		return err
	}

	piped, err := getPiped(ctx, a.pipedStore, pipedID, a.logger)
	if err != nil {
		return err
	}

	if key.ProjectId != piped.ProjectId {
		return status.Error(codes.PermissionDenied, "Requested piped doesn't belong to your project")
	}

	if err := updater(ctx, pipedID); err != nil {
		return gRPCStoreError(err, fmt.Sprintf("update piped %s", pipedID))
	}
	return nil
}

func (a *API) RegisterEvent(ctx context.Context, req *apiservice.RegisterEventRequest) (*apiservice.RegisterEventResponse, error) {
	key, err := requireAPIKey(ctx, model.APIKey_READ_WRITE, a.logger)
	if err != nil {
		return nil, err
	}
	id := uuid.New().String()

	event := model.Event{
		Id:                id,
		Name:              req.Name,
		Data:              req.Data,
		Labels:            req.Labels,
		EventKey:          model.MakeEventKey(req.Name, req.Labels),
		ProjectId:         key.ProjectId,
		Status:            model.EventStatus_EVENT_NOT_HANDLED,
		StatusDescription: fmt.Sprintf("It is going to be replaced by %s", req.Data),
	}
	if err = a.eventStore.Add(ctx, event); err != nil {
		return nil, gRPCStoreError(err, fmt.Sprintf("add event %s", id))
	}

	return &apiservice.RegisterEventResponse{EventId: id}, nil
}

func (a *API) RequestPlanPreview(ctx context.Context, req *apiservice.RequestPlanPreviewRequest) (*apiservice.RequestPlanPreviewResponse, error) {
	key, err := requireAPIKey(ctx, model.APIKey_READ_WRITE, a.logger)
	if err != nil {
		return nil, err
	}

	// TODO: We may need to cache the list of pipeds to reduce load on database.
	// Adding the cache after understanding the real situation from our metrics data.
	pipeds, err := a.pipedStore.List(ctx, datastore.ListOptions{
		Filters: []datastore.ListFilter{
			{
				Field:    "ProjectId",
				Operator: datastore.OperatorEqual,
				Value:    key.ProjectId,
			},
			{
				Field:    "Disabled",
				Operator: datastore.OperatorEqual,
				Value:    false,
			},
		},
	})
	if err != nil {
		return nil, gRPCStoreError(err, "list pipeds")
	}

	repositories := make(map[string]string, len(pipeds))
	for _, p := range pipeds {
		for _, r := range p.Repositories {
			if r.Remote == req.RepoRemoteUrl && r.Branch == req.BaseBranch {
				repositories[p.Id] = r.Id
				break
			}
		}
	}
	if len(repositories) == 0 {
		return &apiservice.RequestPlanPreviewResponse{}, nil
	}

	const commander = "pipectl"
	commands := make([]string, 0, len(repositories))

	for pipedID, repositoryID := range repositories {
		cmd := model.Command{
			Id:        uuid.New().String(),
			PipedId:   pipedID,
			ProjectId: key.ProjectId,
			Type:      model.Command_BUILD_PLAN_PREVIEW,
			Commander: commander,
			BuildPlanPreview: &model.Command_BuildPlanPreview{
				RepositoryId: repositoryID,
				HeadBranch:   req.HeadBranch,
				HeadCommit:   req.HeadCommit,
				BaseBranch:   req.BaseBranch,
				Timeout:      req.Timeout,
			},
		}
		if err := addCommand(ctx, a.commandStore, &cmd, a.logger); err != nil {
			return nil, err
		}
		commands = append(commands, cmd.Id)
	}

	return &apiservice.RequestPlanPreviewResponse{
		Commands: commands,
	}, nil
}

func (a *API) GetPlanPreviewResults(ctx context.Context, req *apiservice.GetPlanPreviewResultsRequest) (*apiservice.GetPlanPreviewResultsResponse, error) {
	key, err := requireAPIKey(ctx, model.APIKey_READ_WRITE, a.logger)
	if err != nil {
		return nil, err
	}

	const (
		freshDuration               = 24 * time.Hour
		defaultCommandHandleTimeout = 5 * time.Minute
	)

	var (
		handledCommands = make([]string, 0, len(req.Commands))
		results         = make([]*model.PlanPreviewCommandResult, 0, len(req.Commands))
	)

	commandHandleTimeout := time.Duration(req.CommandHandleTimeout) * time.Second
	if commandHandleTimeout == 0 {
		commandHandleTimeout = defaultCommandHandleTimeout
	}

	// Validate based on command model stored in datastore.
	for _, commandID := range req.Commands {
		cmd, err := getCommand(ctx, a.commandStore, commandID, a.logger)
		if err != nil {
			return nil, err
		}
		if cmd.ProjectId != key.ProjectId {
			a.logger.Warn("detected a request to get planpreview result of an unowned command",
				zap.String("command", commandID),
				zap.String("command-project-id", cmd.ProjectId),
				zap.String("request-project-id", key.ProjectId),
			)
			return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("The requested command %s does not belong to your project", commandID))
		}
		if cmd.Type != model.Command_BUILD_PLAN_PREVIEW {
			return nil, status.Error(codes.FailedPrecondition, fmt.Sprintf("Command %s is not a plan preview command", commandID))
		}

		if !cmd.IsHandled() {
			pipedStatus, err := getPipedStatus(a.pipedStatCache, cmd.PipedId)
			if err != nil {
				a.logger.Error("failed to get or unmarshal piped stat", zap.Error(err))
				pipedStatus = model.Piped_UNKNOWN
			}

			if pipedStatus != model.Piped_ONLINE {
				results = append(results, &model.PlanPreviewCommandResult{
					CommandId: cmd.Id,
					PipedId:   cmd.PipedId,
					Error:     "Maybe Piped is offline currently.",
				})
				continue
			}

			if time.Since(time.Unix(cmd.CreatedAt, 0)) <= commandHandleTimeout {
				return nil, status.Error(codes.NotFound, fmt.Sprintf("Waiting for result of command %s from piped %s", commandID, cmd.PipedId))
			}

			results = append(results, &model.PlanPreviewCommandResult{
				CommandId: cmd.Id,
				PipedId:   cmd.PipedId,
				Error:     "Timed out, maybe the Piped is offline currently.",
			})
			continue
		}

		// There is no reason to fetch output data of command that has been completed a long time ago.
		// So in order to prevent unintended actions, we disallow that ability.
		if time.Since(time.Unix(cmd.HandledAt, 0)) > freshDuration {
			return nil, status.Error(codes.FailedPrecondition, fmt.Sprintf("The output data for command %s is too old for access", commandID))
		}

		handledCommands = append(handledCommands, commandID)
	}

	// Fetch output data to build results.
	for _, commandID := range handledCommands {
		data, err := a.commandOutputGetter.Get(ctx, commandID)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to retrieve output data of command %s", commandID))
		}

		var result model.PlanPreviewCommandResult
		if err := json.Unmarshal(data, &result); err != nil {
			a.logger.Error("failed to unmarshal planpreview command result",
				zap.String("command", commandID),
				zap.Error(err),
			)
			return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to decode output data of command %s", commandID))
		}

		results = append(results, &result)
	}

	// All URL fields inside the result model are empty.
	// So we fill them before sending to the client.
	for _, r := range results {
		r.FillURLs(a.webBaseURL)
	}

	return &apiservice.GetPlanPreviewResultsResponse{
		Results: results,
	}, nil
}

func (a *API) Encrypt(ctx context.Context, req *apiservice.EncryptRequest) (*apiservice.EncryptResponse, error) {
	_, err := requireAPIKey(ctx, model.APIKey_READ_ONLY, a.logger)
	if err != nil {
		return nil, err
	}

	var pubkey []byte
	if v, err := a.encryptionKeyCache.Get(req.PipedId); err == nil {
		pubkey = v.([]byte)
	}
	if pubkey == nil {
		piped, err := getPiped(ctx, a.pipedStore, req.PipedId, a.logger)
		if err != nil {
			return nil, err
		}
		pubkey, err = getEncriptionKey(piped.SecretEncryption)
		if err != nil {
			return nil, err
		}
		a.encryptionKeyCache.Put(req.PipedId, pubkey)
	}
	ciphertext, err := encrypt(req.Plaintext, pubkey, req.Base64Encoding, a.logger)
	if err != nil {
		return nil, err
	}

	return &apiservice.EncryptResponse{
		Ciphertext: ciphertext,
	}, nil
}

// requireAPIKey checks the existence of an API key inside the given context
// and ensures that it has enough permissions for the give role.
func requireAPIKey(ctx context.Context, role model.APIKey_Role, logger *zap.Logger) (*model.APIKey, error) {
	key, err := rpcauth.ExtractAPIKey(ctx)
	if err != nil {
		return nil, err
	}

	switch key.Role {
	case model.APIKey_READ_WRITE:
		return key, nil

	case model.APIKey_READ_ONLY:
		if role == model.APIKey_READ_ONLY {
			return key, nil
		}
		logger.Warn("detected an API key that has insufficient permissions", zap.String("key", key.Id))
		return nil, status.Error(codes.PermissionDenied, "Permission denied")

	default:
		logger.Warn("detected an API key that has an invalid role", zap.String("key", key.Id))
		return nil, status.Error(codes.PermissionDenied, "Invalid role")
	}
}
