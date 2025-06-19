// Copyright 2025 The PipeCD Authors.
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

package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pipe-cd/piped-plugin-sdk-go/signalhandler"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
)

// DeploymentPlugin is the interface that be implemented by a full-spec deployment plugin.
// This kind of plugin should implement all methods to manage resources and execute stages.
// The Config and DeployTargetConfig are the plugin's config defined in piped's config.
type DeploymentPlugin[Config, DeployTargetConfig, ApplicationConfigSpec any] interface {
	StagePlugin[Config, DeployTargetConfig, ApplicationConfigSpec]

	// DetermineVersions determines the versions of the resources that will be deployed.
	DetermineVersions(context.Context, *Config, *DetermineVersionsInput[ApplicationConfigSpec]) (*DetermineVersionsResponse, error)
	// DetermineStrategy determines the strategy to deploy the resources.
	// This is called when the strategy was not determined by common logic, including judging by the pipeline length, whether it is the first deployment, and so on.
	// It should return (nil, nil) if the plugin does not have specific logic for DetermineStrategy.
	DetermineStrategy(context.Context, *Config, *DetermineStrategyInput[ApplicationConfigSpec]) (*DetermineStrategyResponse, error)
	// BuildQuickSyncStages builds the stages that will be executed during the quick sync process.
	BuildQuickSyncStages(context.Context, *Config, *BuildQuickSyncStagesInput) (*BuildQuickSyncStagesResponse, error)
}

// StagePlugin is the interface implemented by a plugin that focuses on executing generic stages.
// This kind of plugin may not implement quick sync stages.
// The Config and DeployTargetConfig are the plugin's config defined in piped's config.
type StagePlugin[Config, DeployTargetConfig, ApplicationConfigSpec any] interface {
	// FetchDefinedStages returns the list of stages that the plugin can execute.
	FetchDefinedStages() []string
	// BuildPipelineSyncStages builds the stages that will be executed by the plugin.
	// The built pipeline includes non-rollback (defined in the application config) and rollback stages.
	// The request contains only non-rollback stages whose names are listed in FetchDefinedStages() of this plugin.
	//
	// Note about the response indexes:
	//  - For a non-rollback stage, use the index given by the request remaining the execution order.
	//  - For a rollback stage, use one of the indexes given by the request.
	//  - The indexes of the response stages must not be duplicated across non-rollback stages and rollback stages.
	//    A non-rollback stage and a rollback stage can have the same index.
	// For example, given request indexes are {2,4,5}, then
	//  - Non-rollback stages indexes must be {2,4,5}
	//  - Rollback stages indexes must be selected from {2,4,5}.  For a deploymentPlugin, using only {2} is recommended.
	BuildPipelineSyncStages(context.Context, *Config, *BuildPipelineSyncStagesInput) (*BuildPipelineSyncStagesResponse, error)
	// ExecuteStage executes the given stage.
	ExecuteStage(context.Context, *Config, []*DeployTarget[DeployTargetConfig], *ExecuteStageInput[ApplicationConfigSpec]) (*ExecuteStageResponse, error)
}

// DeploymentPluginServiceServer is the gRPC server that handles requests from the piped.
type DeploymentPluginServiceServer[Config, DeployTargetConfig, ApplicationConfigSpec any] struct {
	deployment.UnimplementedDeploymentServiceServer
	commonFields

	base          DeploymentPlugin[Config, DeployTargetConfig, ApplicationConfigSpec]
	config        Config
	deployTargets map[string]*DeployTarget[DeployTargetConfig]
}

// Register registers the server to the given gRPC server.
func (s *DeploymentPluginServiceServer[Config, DeployTargetConfig, ApplicationConfigSpec]) Register(server *grpc.Server) {
	deployment.RegisterDeploymentServiceServer(server, s)
}

// setFields sets the common fields and configs to the server.
func (s *DeploymentPluginServiceServer[Config, DeployTargetConfig, ApplicationConfigSpec]) setFields(fields commonFields) error {
	s.commonFields = fields

	cfg := fields.config
	if cfg.Config != nil {
		if err := json.Unmarshal(cfg.Config, &s.config); err != nil {
			s.logger.Fatal("failed to unmarshal the plugin config", zap.Error(err))
			return err
		}
	}

	s.deployTargets = make(map[string]*DeployTarget[DeployTargetConfig], len(cfg.DeployTargets))
	for _, dt := range cfg.DeployTargets {
		var sdkDt DeployTargetConfig
		if err := json.Unmarshal(dt.Config, &sdkDt); err != nil {
			s.logger.Fatal("failed to unmarshal deploy target config", zap.Error(err))
			return err
		}
		s.deployTargets[dt.Name] = &DeployTarget[DeployTargetConfig]{
			Name:   dt.Name,
			Labels: dt.Labels,
			Config: sdkDt,
		}
	}

	return nil
}

func (s *DeploymentPluginServiceServer[Config, DeployTargetConfig, ApplicationConfigSpec]) FetchDefinedStages(context.Context, *deployment.FetchDefinedStagesRequest) (*deployment.FetchDefinedStagesResponse, error) {
	return &deployment.FetchDefinedStagesResponse{Stages: s.base.FetchDefinedStages()}, nil
}
func (s *DeploymentPluginServiceServer[Config, DeployTargetConfig, ApplicationConfigSpec]) DetermineVersions(ctx context.Context, request *deployment.DetermineVersionsRequest) (*deployment.DetermineVersionsResponse, error) {
	client := &Client{
		base:          s.client,
		pluginName:    s.name,
		applicationID: request.GetInput().GetDeployment().GetApplicationId(),
		deploymentID:  request.GetInput().GetDeployment().GetId(),
		toolRegistry:  s.toolRegistry,
	}

	req, err := newDetermineVersionsRequest[ApplicationConfigSpec](s.name, request)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to parse deployment source: %v", err)
	}
	input := &DetermineVersionsInput[ApplicationConfigSpec]{
		Request: req,
		Client:  client,
		Logger:  s.logger,
	}

	versions, err := s.base.DetermineVersions(ctx, &s.config, input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to determine versions: %v", err)
	}
	return &deployment.DetermineVersionsResponse{
		Versions: versions.toModel(),
	}, nil
}
func (s *DeploymentPluginServiceServer[Config, DeployTargetConfig, ApplicationConfigSpec]) DetermineStrategy(ctx context.Context, request *deployment.DetermineStrategyRequest) (*deployment.DetermineStrategyResponse, error) {
	client := &Client{
		base:          s.client,
		pluginName:    s.name,
		applicationID: request.GetInput().GetDeployment().GetApplicationId(),
		deploymentID:  request.GetInput().GetDeployment().GetId(),
		toolRegistry:  s.toolRegistry,
	}

	req, err := newDetermineStrategyRequest[ApplicationConfigSpec](s.name, request)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to parse deployment source: %v", err)
	}
	input := &DetermineStrategyInput[ApplicationConfigSpec]{
		Request: req,
		Client:  client,
		Logger:  s.logger,
	}

	response, err := s.base.DetermineStrategy(ctx, &s.config, input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to determine strategy: %v", err)
	}
	if response == nil {
		// If the plugin does not have specific logic to determine strategy,
		// use PipelineSync by default.
		response = &DetermineStrategyResponse{
			Strategy: SyncStrategyPipelineSync,
			Summary:  "Use PipelineSync because no other logic was matched",
		}
	}
	return newDetermineStrategyResponse(response)
}
func (s *DeploymentPluginServiceServer[Config, DeployTargetConfig, ApplicationConfigSpec]) BuildPipelineSyncStages(ctx context.Context, request *deployment.BuildPipelineSyncStagesRequest) (*deployment.BuildPipelineSyncStagesResponse, error) {
	client := &Client{
		base:       s.client,
		pluginName: s.name,
	}
	return buildPipelineSyncStages(ctx, s.name, s.base, &s.config, client, request, s.logger)
}
func (s *DeploymentPluginServiceServer[Config, DeployTargetConfig, ApplicationConfigSpec]) BuildQuickSyncStages(ctx context.Context, request *deployment.BuildQuickSyncStagesRequest) (*deployment.BuildQuickSyncStagesResponse, error) {
	input := &BuildQuickSyncStagesInput{
		Request: BuildQuickSyncStagesRequest{
			Rollback: request.GetRollback(),
		},
		Client: &Client{
			base:       s.client,
			pluginName: s.name,
		},
		Logger: s.logger,
	}

	response, err := s.base.BuildQuickSyncStages(ctx, &s.config, input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build quick sync stages: %v", err)
	}
	return newQuickSyncStagesResponse(s.name, time.Now(), response), nil
}
func (s *DeploymentPluginServiceServer[Config, DeployTargetConfig, ApplicationConfigSpec]) ExecuteStage(ctx context.Context, request *deployment.ExecuteStageRequest) (response *deployment.ExecuteStageResponse, _ error) {
	lp := s.logPersister.StageLogPersister(request.GetInput().GetDeployment().GetId(), request.GetInput().GetStage().GetId())
	defer func() {
		// When termination signal received and the stage is not completed yet, we should not mark the log persister as completed.
		// This can occur when the piped is shutting down while the stage is still running.
		if !response.GetStatus().IsCompleted() && signalhandler.Terminated() {
			return
		}
		lp.Complete(time.Minute)
	}()

	client := &Client{
		base:          s.client,
		pluginName:    s.name,
		applicationID: request.GetInput().GetDeployment().GetApplicationId(),
		deploymentID:  request.GetInput().GetDeployment().GetId(),
		stageID:       request.GetInput().GetStage().GetId(),
		logPersister:  lp,
		toolRegistry:  s.toolRegistry,
	}

	// Get the deploy targets set on the deployment from the piped plugin config.
	dtNames := request.GetInput().GetDeployment().GetDeployTargets(s.commonFields.config.Name)
	deployTargets := make([]*DeployTarget[DeployTargetConfig], 0, len(dtNames))
	for _, name := range dtNames {
		dt, ok := s.deployTargets[name]
		if !ok {
			return nil, status.Errorf(codes.Internal, "the deploy target %s is not found in the piped plugin config", name)
		}

		deployTargets = append(deployTargets, dt)
	}

	return executeStage(ctx, s.name, s.base, &s.config, deployTargets, client, request, s.logger)
}

// StagePluginServiceServer is the gRPC server that handles requests from the piped.
type StagePluginServiceServer[Config, DeployTargetConfig, ApplicationConfigSpec any] struct {
	deployment.UnimplementedDeploymentServiceServer
	commonFields

	base   StagePlugin[Config, DeployTargetConfig, ApplicationConfigSpec]
	config Config
}

// Register registers the server to the given gRPC server.
func (s *StagePluginServiceServer[Config, DeployTargetConfig, ApplicationConfigSpec]) Register(server *grpc.Server) {
	deployment.RegisterDeploymentServiceServer(server, s)
}

// setFields sets the common fields and configs to the server.
func (s *StagePluginServiceServer[Config, DeployTargetConfig, ApplicationConfigSpec]) setFields(fields commonFields) error {
	s.commonFields = fields

	cfg := fields.config
	if cfg.Config != nil {
		if err := json.Unmarshal(cfg.Config, &s.config); err != nil {
			s.logger.Fatal("failed to unmarshal the plugin config", zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *StagePluginServiceServer[Config, DeployTargetConfig, ApplicationConfigSpec]) FetchDefinedStages(context.Context, *deployment.FetchDefinedStagesRequest) (*deployment.FetchDefinedStagesResponse, error) {
	return &deployment.FetchDefinedStagesResponse{Stages: s.base.FetchDefinedStages()}, nil
}
func (s *StagePluginServiceServer[Config, DeployTargetConfig, ApplicationConfigSpec]) DetermineVersions(context.Context, *deployment.DetermineVersionsRequest) (*deployment.DetermineVersionsResponse, error) {
	return &deployment.DetermineVersionsResponse{}, nil
}
func (s *StagePluginServiceServer[Config, DeployTargetConfig, ApplicationConfigSpec]) DetermineStrategy(context.Context, *deployment.DetermineStrategyRequest) (*deployment.DetermineStrategyResponse, error) {
	return &deployment.DetermineStrategyResponse{Unsupported: true}, nil
}
func (s *StagePluginServiceServer[Config, DeployTargetConfig, ApplicationConfigSpec]) BuildPipelineSyncStages(ctx context.Context, request *deployment.BuildPipelineSyncStagesRequest) (*deployment.BuildPipelineSyncStagesResponse, error) {
	client := &Client{
		base:       s.client,
		pluginName: s.name,
	}

	return buildPipelineSyncStages(ctx, s.name, s.base, &s.config, client, request, s.logger)
}
func (s *StagePluginServiceServer[Config, DeployTargetConfig, ApplicationConfigSpec]) BuildQuickSyncStages(context.Context, *deployment.BuildQuickSyncStagesRequest) (*deployment.BuildQuickSyncStagesResponse, error) {
	// Return an empty response in case the plugin does not support the QuickSync strategy.
	return &deployment.BuildQuickSyncStagesResponse{}, nil
}
func (s *StagePluginServiceServer[Config, DeployTargetConfig, ApplicationConfigSpec]) ExecuteStage(ctx context.Context, request *deployment.ExecuteStageRequest) (response *deployment.ExecuteStageResponse, _ error) {
	lp := s.logPersister.StageLogPersister(request.GetInput().GetDeployment().GetId(), request.GetInput().GetStage().GetId())
	defer func() {
		// When termination signal received and the stage is not completed yet, we should not mark the log persister as completed.
		// This can occur when the piped is shutting down while the stage is still running.
		if !response.GetStatus().IsCompleted() && signalhandler.Terminated() {
			return
		}
		lp.Complete(time.Minute)
	}()

	client := &Client{
		base:          s.client,
		pluginName:    s.name,
		applicationID: request.GetInput().GetDeployment().GetApplicationId(),
		deploymentID:  request.GetInput().GetDeployment().GetId(),
		stageID:       request.GetInput().GetStage().GetId(),
		logPersister:  lp,
		toolRegistry:  s.toolRegistry,
	}

	return executeStage(ctx, s.name, s.base, &s.config, nil, client, request, s.logger) // TODO: pass the deployTargets
}

// buildPipelineSyncStages builds the stages that will be executed by the plugin.
func buildPipelineSyncStages[Config, DeployTargetConfig, ApplicationConfigSpec any](ctx context.Context, pluginName string, plugin StagePlugin[Config, DeployTargetConfig, ApplicationConfigSpec], config *Config, client *Client, request *deployment.BuildPipelineSyncStagesRequest, logger *zap.Logger) (*deployment.BuildPipelineSyncStagesResponse, error) {
	resp, err := plugin.BuildPipelineSyncStages(ctx, config, newPipelineSyncStagesInput(request, client, logger))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build pipeline sync stages: %v", err)
	}
	return newPipelineSyncStagesResponse(pluginName, time.Now(), request, resp)
}

func executeStage[Config, DeployTargetConfig, ApplicationConfigSpec any](
	ctx context.Context,
	pluginName string,
	plugin StagePlugin[Config, DeployTargetConfig, ApplicationConfigSpec],
	config *Config,
	deployTargets []*DeployTarget[DeployTargetConfig],
	client *Client,
	request *deployment.ExecuteStageRequest,
	logger *zap.Logger,
) (*deployment.ExecuteStageResponse, error) {
	targetDeploymentSource, err := newDeploymentSource[ApplicationConfigSpec](pluginName, request.GetInput().GetTargetDeploymentSource())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create target deployment source: %v", err)
	}

	// running deploy source is empty on the first deployment
	runningDeploymentSource := DeploymentSource[ApplicationConfigSpec]{}
	if request.GetInput().GetRunningDeploymentSource() != nil {
		runningDeploymentSource, err = newDeploymentSource[ApplicationConfigSpec](pluginName, request.GetInput().GetRunningDeploymentSource())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create running deployment source: %v", err)
		}
	}

	in := &ExecuteStageInput[ApplicationConfigSpec]{
		Request: ExecuteStageRequest[ApplicationConfigSpec]{
			StageName:               request.GetInput().GetStage().GetName(),
			StageConfig:             request.GetInput().GetStageConfig(),
			RunningDeploymentSource: runningDeploymentSource,
			TargetDeploymentSource:  targetDeploymentSource,
			Deployment:              newDeployment(request.GetInput().GetDeployment()),
		},
		Client: client,
		Logger: logger,
	}

	resp, err := plugin.ExecuteStage(ctx, config, deployTargets, in)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to execute stage: %v", err)
	}

	return &deployment.ExecuteStageResponse{
		Status: resp.Status.toModelEnum(),
	}, nil
}

// ManualOperation represents the manual operation that the user can perform.
type ManualOperation int

const (
	// ManualOperationNone indicates that there is no manual operation.
	ManualOperationNone ManualOperation = iota
	// ManualOperationSkip indicates that the manual operation is to skip the stage.
	ManualOperationSkip
	// ManualOperationApprove indicates that the manual operation is to approve the stage.
	ManualOperationApprove
)

// toModelEnum converts the ManualOperation to the model.ManualOperation.
func (o ManualOperation) toModelEnum() model.ManualOperation {
	switch o {
	case ManualOperationNone:
		return model.ManualOperation_MANUAL_OPERATION_NONE
	case ManualOperationSkip:
		return model.ManualOperation_MANUAL_OPERATION_SKIP
	case ManualOperationApprove:
		return model.ManualOperation_MANUAL_OPERATION_APPROVE
	default:
		return model.ManualOperation_MANUAL_OPERATION_UNKNOWN
	}
}

// newPipelineSyncStagesInput converts the request to the internal representation.
func newPipelineSyncStagesInput(request *deployment.BuildPipelineSyncStagesRequest, client *Client, logger *zap.Logger) *BuildPipelineSyncStagesInput {
	stages := make([]StageConfig, 0, len(request.Stages))
	for _, s := range request.GetStages() {
		stages = append(stages, StageConfig{
			Index:  int(s.GetIndex()),
			Name:   s.GetName(),
			Config: s.GetConfig(),
		})
	}
	req := BuildPipelineSyncStagesRequest{
		Rollback: request.GetRollback(),
		Stages:   stages,
	}
	return &BuildPipelineSyncStagesInput{
		Request: req,
		Client:  client,
		Logger:  logger,
	}
}

// newPipelineSyncStagesResponse converts the response to the external representation.
func newPipelineSyncStagesResponse(pluginName string, now time.Time, request *deployment.BuildPipelineSyncStagesRequest, response *BuildPipelineSyncStagesResponse) (*deployment.BuildPipelineSyncStagesResponse, error) {
	// Convert the request stages to a map for easier access.
	requestStages := make(map[int]*deployment.BuildPipelineSyncStagesRequest_StageConfig, len(request.GetStages()))
	for _, s := range request.GetStages() {
		requestStages[int(s.GetIndex())] = s
	}

	stages := make([]*model.PipelineStage, 0, len(response.Stages))
	for _, s := range response.Stages {
		// Find the corresponding stage in the request.
		requestStage, ok := requestStages[s.Index]
		if !ok {
			return nil, status.Errorf(codes.Internal, "missing stage with index %d in the request, it's unexpected behavior of the plugin", s.Index)
		}
		id := requestStage.GetId()
		if id == "" {
			id = fmt.Sprintf("%s-stage-%d", pluginName, s.Index)
		}

		stages = append(stages, s.toModel(id, requestStage.GetDesc(), now))
	}
	return &deployment.BuildPipelineSyncStagesResponse{
		Stages: stages,
	}, nil
}

// newQuickSyncStagesResponse converts the response to the external representation.
func newQuickSyncStagesResponse(pluginName string, now time.Time, response *BuildQuickSyncStagesResponse) *deployment.BuildQuickSyncStagesResponse {
	stages := make([]*model.PipelineStage, 0, len(response.Stages))
	for i, s := range response.Stages {
		id := fmt.Sprintf("%s-stage-%d", pluginName, i)
		stages = append(stages, s.toModel(id, now))
	}
	return &deployment.BuildQuickSyncStagesResponse{
		Stages: stages,
	}
}

// BuildPipelineSyncStagesInput is the input for the BuildPipelineSyncStages method.
type BuildPipelineSyncStagesInput struct {
	// Request is the request to build pipeline sync stages.
	Request BuildPipelineSyncStagesRequest
	// Client is the client to interact with the piped.
	Client *Client
	// Logger is the logger to log the events.
	Logger *zap.Logger
}

// BuildPipelineSyncStagesRequest is the request to build pipeline sync stages.
// Rollback indicates whether the stages for rollback are requested.
type BuildPipelineSyncStagesRequest struct {
	// Rollback indicates whether the stages for rollback are requested.
	Rollback bool
	// Stages contains the stage names and their configurations.
	Stages []StageConfig
}

// BuildQuickSyncStagesInput is the input for the BuildQuickSyncStages method.
type BuildQuickSyncStagesInput struct {
	// Request is the request to build pipeline sync stages.
	Request BuildQuickSyncStagesRequest
	// Client is the client to interact with the piped.
	Client *Client
	// Logger is the logger to log the events.
	Logger *zap.Logger
}

// BuildQuickSyncStagesRequest is the request to build quick sync stages.
// Rollback indicates whether the stages for rollback are requested.
type BuildQuickSyncStagesRequest struct {
	// Rollback indicates whether the stages for rollback are requested.
	Rollback bool
}

// StageConfig represents the configuration of a stage.
type StageConfig struct {
	// Index is the order of the stage in the pipeline.
	Index int
	// Name is the name of the stage.
	// It must be one of the stages returned by FetchDefinedStages.
	Name string
	// Config is the configuration of the stage.
	// It should be marshaled to JSON bytes.
	// The plugin should unmarshal it to the appropriate struct.
	Config []byte
}

// BuildPipelineSyncStagesResponse is the response of the request to build pipeline sync stages.
type BuildPipelineSyncStagesResponse struct {
	Stages []PipelineStage
}

// BuildQuickSyncStagesResponse is the response of the request to build quick sync stages.
type BuildQuickSyncStagesResponse struct {
	Stages []QuickSyncStage
}

// PipelineStage represents a stage in the pipeline.
type PipelineStage struct {
	// Index is the order of the stage in the pipeline.
	// The value must be one of the index of the stage in the request.
	// The rollback stage should have the same index as the original stage.
	Index int
	// Name is the name of the stage.
	// It must be one of the stages returned by FetchDefinedStages.
	Name string
	// Rollback indicates whether the stage is for rollback.
	Rollback bool
	// Metadata contains the metadata of the stage.
	Metadata map[string]string
	// AvailableOperation indicates the manual operation that the user can perform.
	AvailableOperation ManualOperation
}

func (p *PipelineStage) toModel(id, description string, now time.Time) *model.PipelineStage {
	return &model.PipelineStage{
		Id:                 id,
		Name:               p.Name,
		Desc:               description,
		Index:              int32(p.Index),
		Status:             model.StageStatus_STAGE_NOT_STARTED_YET,
		StatusReason:       "", // TODO: set the reason
		Metadata:           p.Metadata,
		Rollback:           p.Rollback,
		CreatedAt:          now.Unix(),
		UpdatedAt:          now.Unix(),
		AvailableOperation: p.AvailableOperation.toModelEnum(),
	}
}

// QuickSyncStage represents a stage in the pipeline.
type QuickSyncStage struct {
	// Name is the name of the stage.
	// It must be one of the stages returned by FetchDefinedStages.
	Name string
	// Description is the description of the stage.
	Description string
	// Rollback indicates whether the stage is for rollback.
	Rollback bool
	// Metadata contains the metadata of the stage.
	Metadata map[string]string
	// AvailableOperation indicates the manual operation that the user can perform.
	AvailableOperation ManualOperation
}

func (p *QuickSyncStage) toModel(id string, now time.Time) *model.PipelineStage {
	return &model.PipelineStage{
		Id:                 id,
		Name:               p.Name,
		Desc:               p.Description,
		Index:              0,
		Status:             model.StageStatus_STAGE_NOT_STARTED_YET,
		StatusReason:       "", // TODO: set the reason
		Metadata:           p.Metadata,
		Rollback:           p.Rollback,
		CreatedAt:          now.Unix(),
		UpdatedAt:          now.Unix(),
		AvailableOperation: p.AvailableOperation.toModelEnum(),
	}
}

// ExecuteStageInput is the input for the ExecuteStage method.
type ExecuteStageInput[ApplicationConfigSpec any] struct {
	// Request is the request to execute a stage.
	Request ExecuteStageRequest[ApplicationConfigSpec]
	// Client is the client to interact with the piped.
	Client *Client
	// Logger is the logger to log the events.
	Logger *zap.Logger
}

// ExecuteStageRequest is the request to execute a stage.
type ExecuteStageRequest[ApplicationConfigSpec any] struct {
	// The name of the stage to execute.
	StageName string
	// Json encoded configuration of the stage.
	StageConfig []byte

	// RunningDeploymentSource is the source of the running deployment.
	RunningDeploymentSource DeploymentSource[ApplicationConfigSpec]

	// TargetDeploymentSource is the source of the target deployment.
	TargetDeploymentSource DeploymentSource[ApplicationConfigSpec]

	// The deployment that the stage is running.
	Deployment Deployment
}

// Deployment represents the deployment that the stage is running. This is read-only.
type Deployment struct {
	// ID is the unique identifier of the deployment.
	ID string
	// ApplicationID is the unique identifier of the application.
	ApplicationID string
	// ApplicationName is the name of the application.
	ApplicationName string
	// PipedID is the unique identifier of the piped that is running the deployment.
	PipedID string
	// ProjectID is the unique identifier of the project that the application belongs to.
	ProjectID string
	// TriggeredBy is the name of the entity that triggered the deployment.
	TriggeredBy string
	// CreatedAt is the time when the deployment was created.
	CreatedAt int64
	// RepositoryURL is the repo remote path
	RepositoryURL string
	// Summary is the simple description about what this deployment does
	Summary string
	// Labels are custom attributes to identify applications
	Labels map[string]string
}

// newDeployment converts the model.Deployment to the internal representation.
func newDeployment(deployment *model.Deployment) Deployment {
	return Deployment{
		ID:              deployment.GetId(),
		ApplicationID:   deployment.GetApplicationId(),
		ApplicationName: deployment.GetApplicationName(),
		PipedID:         deployment.GetPipedId(),
		ProjectID:       deployment.GetProjectId(),
		TriggeredBy:     deployment.TriggeredBy(),
		CreatedAt:       deployment.GetCreatedAt(),
		RepositoryURL:   deployment.GetGitPath().GetRepo().GetRemote(),
		Summary:         deployment.GetSummary(),
		Labels:          deployment.GetLabels(),
	}
}

// ExecuteStageResponse is the response of the request to execute a stage.
type ExecuteStageResponse struct {
	Status StageStatus
}

// StageStatus represents the current status of a stage of a deployment.
type StageStatus int

const (
	_ StageStatus = iota
	// StageStatusSuccess indicates that the stage succeeded.
	StageStatusSuccess
	// StageStatusFailure indicates that the stage failed.
	StageStatusFailure
	// StageStatusExited can be used when the stage succeeded and exit the pipeline without executing the following stages.
	StageStatusExited

	// StageStatusSkipped // TODO: If SDK can handle whole skipping, this is unnecessary.
)

// toModelEnum converts the StageStatus to the model.StageStatus.
// It returns model.StageStatus_STAGE_FAILURE if the given value is invalid.
func (o StageStatus) toModelEnum() model.StageStatus {
	switch o {
	case StageStatusSuccess:
		return model.StageStatus_STAGE_SUCCESS
	case StageStatusFailure:
		return model.StageStatus_STAGE_FAILURE
	case StageStatusExited:
		return model.StageStatus_STAGE_EXITED
	default:
		return model.StageStatus_STAGE_FAILURE
	}
}

func (o StageStatus) String() string {
	switch o {
	case StageStatusSuccess:
		return model.StageStatus_STAGE_SUCCESS.String()
	case StageStatusFailure:
		return model.StageStatus_STAGE_FAILURE.String()
	case StageStatusExited:
		return model.StageStatus_STAGE_EXITED.String()
	default:
		return model.StageStatus_STAGE_FAILURE.String()
	}
}

// StageCommand represents a command for a stage.
type StageCommand struct {
	Commander string
	Type      CommandType
}

// CommandType represents the type of the command.
type CommandType int32

const (
	CommandTypeApproveStage CommandType = iota
	CommandTypeSkipStage
)

// newStageCommand converts the model.Command to the internal representation.
func newStageCommand(c *model.Command) (StageCommand, error) {
	switch c.Type {
	case model.Command_APPROVE_STAGE:
		return StageCommand{
			Commander: c.GetCommander(),
			Type:      CommandTypeApproveStage,
		}, nil
	case model.Command_SKIP_STAGE:
		return StageCommand{
			Commander: c.GetCommander(),
			Type:      CommandTypeSkipStage,
		}, nil
	default:
		return StageCommand{}, fmt.Errorf("invalid command type: %d", c.Type)
	}
}

// DetermineVersionsInput is the input for the DetermineVersions method.
type DetermineVersionsInput[ApplicationConfigSpec any] struct {
	// Request is the request to determine versions.
	Request DetermineVersionsRequest[ApplicationConfigSpec]
	// Client is the client to interact with the piped.
	Client *Client
	// Logger is the logger to log the events.
	Logger *zap.Logger
}

// DetermineVersionsRequest is the request to determine versions.
type DetermineVersionsRequest[ApplicationConfigSpec any] struct {
	// Deloyment is the deployment that the versions will be determined.
	Deployment Deployment
	// DeploymentSource is the source of the deployment.
	DeploymentSource DeploymentSource[ApplicationConfigSpec]
}

// newDetermineVersionsRequest converts the common.DetermineVersionsRequest to the internal representation.
func newDetermineVersionsRequest[ApplicationConfigSpec any](pluginName string, request *deployment.DetermineVersionsRequest) (DetermineVersionsRequest[ApplicationConfigSpec], error) {
	ds, err := newDeploymentSource[ApplicationConfigSpec](pluginName, request.GetInput().GetTargetDeploymentSource())
	if err != nil {
		return DetermineVersionsRequest[ApplicationConfigSpec]{}, fmt.Errorf("failed to parse target deployment source: %w", err)
	}
	return DetermineVersionsRequest[ApplicationConfigSpec]{
		Deployment:       newDeployment(request.GetInput().GetDeployment()),
		DeploymentSource: ds,
	}, nil
}

// DetermineVersionsResponse is the response of the request to determine versions.
type DetermineVersionsResponse struct {
	// Versions contains the versions of the resources.
	Versions []ArtifactVersion
}

// toModel converts the DetermineVersionsResponse to the model.ArtifactVersion.
func (r *DetermineVersionsResponse) toModel() []*model.ArtifactVersion {
	versions := make([]*model.ArtifactVersion, 0, len(r.Versions))
	for _, v := range r.Versions {
		versions = append(versions, v.toModel())
	}
	return versions
}

// ArtifactVersion represents the version of an artifact.
type ArtifactVersion struct {
	// Version is the version of the artifact.
	Version string
	// Name is the name of the artifact.
	Name string
	// URL is the URL of the artifact.
	URL string
}

// toModel converts the ArtifactVersion to the model.ArtifactVersion.
func (v *ArtifactVersion) toModel() *model.ArtifactVersion {
	return &model.ArtifactVersion{
		Version: v.Version,
		Name:    v.Name,
		Url:     v.URL,
	}
}

// DetermineStrategyInput is the input for the DetermineStrategy method.
type DetermineStrategyInput[ApplicationConfigSpec any] struct {
	// Request is the request to determine the strategy.
	Request DetermineStrategyRequest[ApplicationConfigSpec]
	// Client is the client to interact with the piped.
	Client *Client
	// Logger is the logger to log the events.
	Logger *zap.Logger
}

// DetermineStrategyRequest is the request to determine the strategy.
type DetermineStrategyRequest[ApplicationConfigSpec any] struct {
	// Deployment is the deployment that the strategy will be determined.
	Deployment Deployment
	// RunningDeploymentSource is the source of the running deployment.
	RunningDeploymentSource DeploymentSource[ApplicationConfigSpec]
	// TargetDeploymentSource is the source of the target deployment.
	TargetDeploymentSource DeploymentSource[ApplicationConfigSpec]
}

// newDetermineStrategyRequest converts the common.DetermineStrategyRequest to the internal representation.
func newDetermineStrategyRequest[ApplicationConfigSpec any](pluginName string, request *deployment.DetermineStrategyRequest) (DetermineStrategyRequest[ApplicationConfigSpec], error) {
	rds, err := newDeploymentSource[ApplicationConfigSpec](pluginName, request.GetInput().GetRunningDeploymentSource())
	if err != nil {
		return DetermineStrategyRequest[ApplicationConfigSpec]{}, fmt.Errorf("failed to parse running deployment source: %w", err)
	}
	tds, err := newDeploymentSource[ApplicationConfigSpec](pluginName, request.GetInput().GetTargetDeploymentSource())
	if err != nil {
		return DetermineStrategyRequest[ApplicationConfigSpec]{}, fmt.Errorf("failed to parse target deployment source: %w", err)
	}
	return DetermineStrategyRequest[ApplicationConfigSpec]{
		Deployment:              newDeployment(request.GetInput().GetDeployment()),
		RunningDeploymentSource: rds,
		TargetDeploymentSource:  tds,
	}, nil
}

// DetermineStrategyResponse is the response of the request to determine the strategy.
type DetermineStrategyResponse struct {
	// Strategy is the strategy to deploy the resources.
	Strategy SyncStrategy
	// Summary is the summary of the strategy.
	Summary string
}

// newDetermineStrategyResponse converts the response to the external representation.
func newDetermineStrategyResponse(response *DetermineStrategyResponse) (*deployment.DetermineStrategyResponse, error) {
	strategy, err := response.Strategy.toModelEnum()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert the strategy: %v", err)
	}
	return &deployment.DetermineStrategyResponse{
		Summary:      response.Summary,
		SyncStrategy: strategy,
	}, nil
}

// SyncStrategy represents the strategy to deploy the resources.
type SyncStrategy int

const (
	_ SyncStrategy = iota
	// SyncStrategyQuickSync indicates that the resources will be deployed using the quick sync strategy.
	SyncStrategyQuickSync
	// SyncStrategyPipelineSync indicates that the resources will be deployed using the pipeline sync strategy.
	SyncStrategyPipelineSync
)

// toModelEnum converts the SyncStrategy to the model.SyncStrategy.
// It returns an error if the given value is invalid.
func (s SyncStrategy) toModelEnum() (model.SyncStrategy, error) {
	switch s {
	case SyncStrategyQuickSync:
		return model.SyncStrategy_QUICK_SYNC, nil
	case SyncStrategyPipelineSync:
		return model.SyncStrategy_PIPELINE, nil
	default:
		return 0, fmt.Errorf("invalid sync strategy: %d", s)
	}
}
