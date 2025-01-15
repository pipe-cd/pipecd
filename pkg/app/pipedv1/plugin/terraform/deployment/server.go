// Copyright 2024 The PipeCD Authors.
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

package deployment

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	tfconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/toolregistry"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister"
	"github.com/pipe-cd/pipecd/pkg/plugin/signalhandler"
)

type toolClient interface {
	InstallTool(ctx context.Context, name, version, script string) (path string, err error)
}

type toolRegistry interface {
	Terraform(ctx context.Context, version string) (path string, err error)
}

type logPersister interface {
	StageLogPersister(deploymentID, stageID string) logpersister.StageLogPersister
}

type DeploymentServiceServer struct {
	deployment.UnimplementedDeploymentServiceServer

	// this field is set with the plugin configuration
	// the plugin configuration is sent from piped while initializing the plugin
	pluginConfig *config.PipedPlugin
	// deployTargetConfig might be empty. e.g. when it's not specified in the piped config.
	// For now, this plugin supports up to one deploy target.
	deployTargetConfig tfconfig.TerraformDeployTargetConfig

	logger       *zap.Logger
	toolRegistry toolRegistry
	logPersister logPersister
}

// NewDeploymentServiceServer creates a new DeploymentServiceServer of Terraform plugin.
func NewDeploymentServiceServer(
	config *config.PipedPlugin,
	logger *zap.Logger,
	toolClient toolClient,
	logPersister logPersister,
) (*DeploymentServiceServer, error) {
	toolRegistry := toolregistry.NewRegistry(toolClient)

	deployTargetConfig := tfconfig.TerraformDeployTargetConfig{}
	if len(config.DeployTargets) > 0 {
		var err error
		if deployTargetConfig, err = tfconfig.ParseDeployTargetConfig(config.DeployTargets[0]); err != nil {
			return nil, err
		}
	}

	return &DeploymentServiceServer{
		pluginConfig:       config,
		deployTargetConfig: deployTargetConfig,
		logger:             logger.Named("terraform-plugin"),
		toolRegistry:       toolRegistry,
		logPersister:       logPersister,
	}, nil
}

// Register registers all handling of this service into the specified gRPC server.
func (s *DeploymentServiceServer) Register(server *grpc.Server) {
	deployment.RegisterDeploymentServiceServer(server, s)
}

// DetermineStrategy implements deployment.DeploymentServiceServer.
func (s *DeploymentServiceServer) DetermineStrategy(ctx context.Context, request *deployment.DetermineStrategyRequest) (*deployment.DetermineStrategyResponse, error) {
	return &deployment.DetermineStrategyResponse{
		Unsupported: true,
	}, nil
}

// DetermineVersions implements deployment.DeploymentServiceServer.
func (s *DeploymentServiceServer) DetermineVersions(ctx context.Context, request *deployment.DetermineVersionsRequest) (*deployment.DetermineVersionsResponse, error) {
	return &deployment.DetermineVersionsResponse{
		Versions: nil,
	}, nil
}

// BuildPipelineSyncStages implements deployment.DeploymentServiceServer.
func (s *DeploymentServiceServer) BuildPipelineSyncStages(ctx context.Context, request *deployment.BuildPipelineSyncStagesRequest) (*deployment.BuildPipelineSyncStagesResponse, error) {
	now := time.Now()
	stages := buildPipelineStages(request.GetStages(), request.GetRollback(), now)
	return &deployment.BuildPipelineSyncStagesResponse{
		Stages: stages,
	}, nil
}

// BuildQuickSyncStages implements deployment.DeploymentServiceServer.
func (s *DeploymentServiceServer) BuildQuickSyncStages(ctx context.Context, request *deployment.BuildQuickSyncStagesRequest) (*deployment.BuildQuickSyncStagesResponse, error) {
	now := time.Now()
	stages := buildQuickSyncStages(request.GetRollback(), now)
	return &deployment.BuildQuickSyncStagesResponse{
		Stages: stages,
	}, nil
}

// FetchDefinedStages implements deployment.DeploymentServiceServer.
func (s *DeploymentServiceServer) FetchDefinedStages(context.Context, *deployment.FetchDefinedStagesRequest) (*deployment.FetchDefinedStagesResponse, error) {
	return &deployment.FetchDefinedStagesResponse{
		Stages: allStages,
	}, nil
}

// ExecuteStage performs stage-defined tasks.
// It returns stage status after execution without error.
// An error will be returned only if the given stage is not supported.
func (s *DeploymentServiceServer) ExecuteStage(ctx context.Context, request *deployment.ExecuteStageRequest) (response *deployment.ExecuteStageResponse, _ error) {
	lp := s.logPersister.StageLogPersister(request.GetInput().GetDeployment().GetId(), request.GetInput().GetStage().GetId())
	defer func() {
		// When termination signal received and the stage is not completed yet, we should not mark the log persister as completed.
		// This can occur when the piped is shutting down while the stage is still running.
		if !response.GetStatus().IsCompleted() && signalhandler.Terminated() {
			return
		}
		lp.Complete(time.Minute)
	}()

	status, err := s.executeStage(ctx, lp, request.GetInput())
	if err != nil {
		return nil, err
	}
	return &deployment.ExecuteStageResponse{
		Status: status,
	}, nil
}
