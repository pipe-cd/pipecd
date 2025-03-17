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

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/toolregistry"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/common"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type toolClient interface {
	InstallTool(ctx context.Context, name, version, script string) (string, error)
}

type logPersister interface {
	StageLogPersister(deploymentID, stageID string) logpersister.StageLogPersister
}

type DeploymentService struct {
	deployment.UnimplementedDeploymentServiceServer

	// this field is set with the plugin configuration
	// the plugin configuration is sent from piped while initializing the plugin
	pluginConfig *config.PipedPlugin

	logger       *zap.Logger
	toolRegistry toolRegistry
	loader       loader
	logPersister logPersister
}

// NewDeploymentService creates a new planService.
func NewDeploymentService(
	config *config.PipedPlugin,
	logger *zap.Logger,
	toolClient toolClient,
	logPersister logPersister,
) *DeploymentService {
	toolRegistry := toolregistry.NewRegistry(toolClient)

	return &DeploymentService{
		pluginConfig: config,
		logger:       logger.Named("planner"),
		toolRegistry: toolRegistry,
		loader:       provider.NewLoader(toolRegistry),
		logPersister: logPersister,
	}
}

// Register registers all handling of this service into the specified gRPC server.
func (a *DeploymentService) Register(server *grpc.Server) {
	deployment.RegisterDeploymentServiceServer(server, a)
}

// DetermineStrategy implements deployment.DeploymentServiceServer.
func (a *DeploymentService) DetermineStrategy(ctx context.Context, request *deployment.DetermineStrategyRequest) (*deployment.DetermineStrategyResponse, error) {
	return &deployment.DetermineStrategyResponse{}, nil

}

// DetermineVersions implements deployment.DeploymentServiceServer.
func (a *DeploymentService) DetermineVersions(ctx context.Context, request *deployment.DetermineVersionsRequest) (*deployment.DetermineVersionsResponse, error) {
	return &deployment.DetermineVersionsResponse{}, nil
}

// BuildPipelineSyncStages implements deployment.DeploymentServiceServer.
func (a *DeploymentService) BuildPipelineSyncStages(ctx context.Context, request *deployment.BuildPipelineSyncStagesRequest) (*deployment.BuildPipelineSyncStagesResponse, error) {
	return &deployment.BuildPipelineSyncStagesResponse{}, nil
}

// BuildQuickSyncStages implements deployment.DeploymentServiceServer.
func (a *DeploymentService) BuildQuickSyncStages(ctx context.Context, request *deployment.BuildQuickSyncStagesRequest) (*deployment.BuildQuickSyncStagesResponse, error) {
	return &deployment.BuildQuickSyncStagesResponse{}, nil
}

// FetchDefinedStages implements deployment.DeploymentServiceServer.
func (a *DeploymentService) FetchDefinedStages(context.Context, *deployment.FetchDefinedStagesRequest) (*deployment.FetchDefinedStagesResponse, error) {
	stages := make([]string, 0, len(AllStages))
	for _, s := range AllStages {
		stages = append(stages, string(s))
	}

	return &deployment.FetchDefinedStagesResponse{
		Stages: stages,
	}, nil
}

func (a *DeploymentService) loadManifests(ctx context.Context, deploy *model.Deployment, spec *kubeconfig.KubernetesApplicationSpec, deploymentSource *common.DeploymentSource) ([]provider.Manifest, error) {
	manifests, err := a.loader.LoadManifests(ctx, provider.LoaderInput{
		PipedID:          deploy.GetPipedId(),
		AppID:            deploy.GetApplicationId(),
		CommitHash:       deploymentSource.GetCommitHash(),
		AppName:          deploy.GetApplicationName(),
		AppDir:           deploymentSource.GetApplicationDirectory(),
		ConfigFilename:   deploymentSource.GetApplicationConfigFilename(),
		Manifests:        spec.Input.Manifests,
		Namespace:        spec.Input.Namespace,
		TemplatingMethod: provider.TemplatingMethodNone, // TODO: Implement detection of templating method or add it to the config spec.

		// TODO: Define other fields for LoaderInput
	})

	if err != nil {
		return nil, err
	}

	return manifests, nil
}

// ExecuteStage performs stage-defined tasks.
// It returns stage status after execution without error.
// Error only be raised if the given stage is not supported.
func (a *DeploymentService) ExecuteStage(ctx context.Context, request *deployment.ExecuteStageRequest) (response *deployment.ExecuteStageResponse, _ error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
