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
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipe/pkg/app/api/commandstore"
	"github.com/pipe-cd/pipe/pkg/app/api/service/apiservice"
	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/cache/memorycache"
	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/model"
	"github.com/pipe-cd/pipe/pkg/rpc/rpcauth"
)

// API implements the behaviors for the gRPC definitions of API.
type API struct {
	applicationStore    datastore.ApplicationStore
	environmentStore    datastore.EnvironmentStore
	deploymentStore     datastore.DeploymentStore
	pipedStore          datastore.PipedStore
	eventStore          datastore.EventStore
	tagStore            datastore.TagStore
	commandStore        commandstore.Store
	commandOutputGetter commandOutputGetter

	encryptionKeyCache cache.Cache

	webBaseURL string
	logger     *zap.Logger
}

// NewAPI creates a new API instance.
func NewAPI(
	ctx context.Context,
	ds datastore.DataStore,
	cmds commandstore.Store,
	cog commandOutputGetter,
	webBaseURL string,
	logger *zap.Logger,
) *API {
	a := &API{
		applicationStore:    datastore.NewApplicationStore(ds),
		environmentStore:    datastore.NewEnvironmentStore(ds),
		deploymentStore:     datastore.NewDeploymentStore(ds),
		pipedStore:          datastore.NewPipedStore(ds),
		eventStore:          datastore.NewEventStore(ds),
		tagStore:            datastore.NewTagStore(ds),
		commandStore:        cmds,
		commandOutputGetter: cog,
		// Public key is variable but likely to be accessed multiple times in a short period.
		encryptionKeyCache: memorycache.NewTTLCache(ctx, 5*time.Minute, 5*time.Minute),
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

	// TODO: Cache the existing tags.
	//   This is not necessary if you want to pass the Tag model itself to the web client.
	tags, err := getOrCreateTags(ctx, req.TagNames, key.ProjectId, a.tagStore, a.logger)
	if err != nil {
		return nil, err
	}

	app := model.Application{
		Id:            uuid.New().String(),
		Name:          req.Name,
		EnvId:         req.EnvId,
		PipedId:       req.PipedId,
		ProjectId:     key.ProjectId,
		GitPath:       gitpath,
		Kind:          req.Kind,
		CloudProvider: req.CloudProvider,
		Description:   req.Description,
		Tags:          tags,
	}
	err = a.applicationStore.AddApplication(ctx, &app)
	if errors.Is(err, datastore.ErrAlreadyExists) {
		return nil, status.Error(codes.AlreadyExists, "The application already exists")
	}
	if err != nil {
		a.logger.Error("failed to create application", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to create application")
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

	if req.EnvId != "" {
		filters = append(filters, datastore.ListFilter{
			Field:    "EnvId",
			Operator: datastore.OperatorEqual,
			Value:    req.EnvId,
		})
	}
	// Use env-name as listApplications filter only in case env-id is not set.
	if req.EnvId == "" && req.EnvName != "" {
		envListOpts := datastore.ListOptions{
			Filters: []datastore.ListFilter{
				{
					Field:    "ProjectId",
					Operator: datastore.OperatorEqual,
					Value:    key.ProjectId,
				},
				{
					Field:    "Name",
					Operator: datastore.OperatorEqual,
					Value:    req.EnvName,
				},
			},
			Limit: limit,
		}
		envs, err := listEnvironments(ctx, a.environmentStore, envListOpts, a.logger)
		if err != nil {
			return nil, err
		}

		switch len(envs) {
		case 0:
			return nil, status.Error(codes.NotFound, fmt.Sprintf("No environment named as %s", req.EnvName))
		case 1:
			filters = append(filters, datastore.ListFilter{
				Field:    "EnvId",
				Operator: datastore.OperatorEqual,
				Value:    envs[0].Id,
			})
		default:
			envsID := make([]string, 0, len(envs))
			for _, env := range envs {
				envsID = append(envsID, env.Id)
			}
			filters = append(filters, datastore.ListFilter{
				Field:    "EnvId",
				Operator: datastore.OperatorIn,
				Value:    envsID,
			})
		}
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

	apps, cursor, err := listApplications(ctx, a.applicationStore, opts, a.logger)
	if err != nil {
		return nil, err
	}

	return &apiservice.ListApplicationsResponse{
		Applications: apps,
		Cursor:       cursor,
	}, nil
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
			return status.Error(codes.Internal, "Failed to update the piped")
		}
	}
	return nil
}

func (a *API) RegisterEvent(ctx context.Context, req *apiservice.RegisterEventRequest) (*apiservice.RegisterEventResponse, error) {
	key, err := requireAPIKey(ctx, model.APIKey_READ_WRITE, a.logger)
	if err != nil {
		return nil, err
	}
	id := uuid.New().String()

	err = a.eventStore.AddEvent(ctx, model.Event{
		Id:        id,
		Name:      req.Name,
		Data:      req.Data,
		Labels:    req.Labels,
		EventKey:  model.MakeEventKey(req.Name, req.Labels),
		ProjectId: key.ProjectId,
	})
	if errors.Is(err, datastore.ErrAlreadyExists) {
		return nil, status.Error(codes.AlreadyExists, "The event already exists")
	}
	if err != nil {
		a.logger.Error("failed to register event", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to register event")
	}

	return &apiservice.RegisterEventResponse{
		EventId: id,
	}, nil
}

func (a *API) RequestPlanPreview(ctx context.Context, req *apiservice.RequestPlanPreviewRequest) (*apiservice.RequestPlanPreviewResponse, error) {
	key, err := requireAPIKey(ctx, model.APIKey_READ_WRITE, a.logger)
	if err != nil {
		return nil, err
	}

	// TODO: We may need to cache the list of pipeds to reduce load on database.
	// Adding the cache after understanding the real situation from our metrics data.
	pipeds, err := a.pipedStore.ListPipeds(ctx, datastore.ListOptions{
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
		a.logger.Error("failed to list pipeds to request planpreview",
			zap.String("project", key.ProjectId),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "Failed to list pipeds")
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
			return nil, status.Error(codes.FailedPrecondition, fmt.Sprint("Command %s is not a plan preview command", commandID))
		}

		if !cmd.IsHandled() {
			if time.Since(time.Unix(cmd.CreatedAt, 0)) <= commandHandleTimeout {
				return nil, status.Error(codes.NotFound, fmt.Sprintf("Waiting for result of command %s from piped %s", commandID, cmd.PipedId))
			}
			results = append(results, &model.PlanPreviewCommandResult{
				CommandId: cmd.Id,
				PipedId:   cmd.PipedId,
				Error:     fmt.Sprintf("Timed out, maybe the Piped is offline currently."),
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

	// Fetch ouput data to build results.
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
		pubkey, err = getEncriptionKey(model.GetSecretEncryptionInPiped(piped))
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
