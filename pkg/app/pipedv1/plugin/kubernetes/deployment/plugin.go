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

package deployment

import (
	"context"
	"errors"

	"go.uber.org/zap"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/toolregistry"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

// Plugin implements the sdk.DeploymentPlugin interface.
type Plugin struct {
}

type loader interface {
	// LoadManifests renders and loads all manifests for application.
	LoadManifests(ctx context.Context, input provider.LoaderInput) ([]provider.Manifest, error)
}

type toolRegistry interface {
	Kubectl(ctx context.Context, version string) (string, error)
	Kustomize(ctx context.Context, version string) (string, error)
	Helm(ctx context.Context, version string) (string, error)
}

var _ sdk.DeploymentPlugin[sdk.ConfigNone, kubeconfig.KubernetesDeployTargetConfig] = (*Plugin)(nil)

// Name returns the name of this plugin.
func (p *Plugin) Name() string {
	return "kubernetes"
}

// Version returns the version of this plugin.
func (p *Plugin) Version() string {
	return "0.0.1" // TODO
}

// FetchDefinedStages returns the defined stages for this plugin.
func (p *Plugin) FetchDefinedStages() []string {
	stages := make([]string, 0, len(AllStages))
	for _, s := range AllStages {
		stages = append(stages, string(s))
	}

	return stages
}

// BuildPipelineSyncStages returns the stages for the pipeline sync strategy.
func (p *Plugin) BuildPipelineSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	return &sdk.BuildPipelineSyncStagesResponse{
		Stages: buildPipelineStages(input.Request.Stages, input.Request.Rollback),
	}, nil
}

// ExecuteStage executes the stage.
func (p *Plugin) ExecuteStage(ctx context.Context, _ *sdk.ConfigNone, dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], input *sdk.ExecuteStageInput) (*sdk.ExecuteStageResponse, error) {
	switch input.Request.StageName {
	case StageK8sSync.String():
		return &sdk.ExecuteStageResponse{
			Status: p.executeK8sSyncStage(ctx, input, dts),
		}, nil
	case StageK8sRollback.String():
		return &sdk.ExecuteStageResponse{
			Status: p.executeK8sRollbackStage(ctx, input, dts),
		}, nil
	default:
		return nil, errors.New("unimplemented or unsupported stage")
	}
}

func (p *Plugin) loadManifests(ctx context.Context, deploy *sdk.Deployment, spec *kubeconfig.KubernetesApplicationSpec, deploymentSource *sdk.DeploymentSource, loader loader) ([]provider.Manifest, error) {
	manifests, err := loader.LoadManifests(ctx, provider.LoaderInput{
		PipedID:          deploy.PipedID,
		AppID:            deploy.ApplicationID,
		CommitHash:       deploymentSource.CommitHash,
		AppName:          deploy.ApplicationName,
		AppDir:           deploymentSource.ApplicationDirectory,
		ConfigFilename:   deploymentSource.ApplicationConfigFilename,
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

// DetermineVersions determines the versions of the application.
func (p *Plugin) DetermineVersions(ctx context.Context, _ *sdk.ConfigNone, input *sdk.DetermineVersionsInput) (*sdk.DetermineVersionsResponse, error) {
	logger := input.Logger

	cfg, err := config.DecodeYAML[*kubeconfig.KubernetesApplicationSpec](input.Request.DeploymentSource.ApplicationConfig)
	if err != nil {
		logger.Error("Failed while decoding application config", zap.Error(err))
		return nil, err
	}

	manifests, err := p.loadManifests(ctx, &input.Request.Deployment, cfg.Spec, &input.Request.DeploymentSource, provider.NewLoader(toolregistry.NewRegistry(input.Client.ToolRegistry())))

	if err != nil {
		logger.Error("Failed while loading manifests", zap.Error(err))
		return nil, err
	}

	return &sdk.DetermineVersionsResponse{
		Versions: determineVersions(manifests),
	}, nil
}

// DetermineStrategy determines the strategy for the deployment.
func (p *Plugin) DetermineStrategy(ctx context.Context, _ *sdk.ConfigNone, input *sdk.DetermineStrategyInput) (*sdk.DetermineStrategyResponse, error) {
	logger := input.Logger
	loader := provider.NewLoader(toolregistry.NewRegistry(input.Client.ToolRegistry()))

	cfg, err := config.DecodeYAML[*kubeconfig.KubernetesApplicationSpec](input.Request.TargetDeploymentSource.ApplicationConfig)
	if err != nil {
		logger.Error("Failed while decoding application config", zap.Error(err))
		return nil, err
	}

	runnings, err := p.loadManifests(ctx, &input.Request.Deployment, cfg.Spec, &input.Request.RunningDeploymentSource, loader)

	if err != nil {
		logger.Error("Failed while loading running manifests", zap.Error(err))
		return nil, err
	}

	targets, err := p.loadManifests(ctx, &input.Request.Deployment, cfg.Spec, &input.Request.TargetDeploymentSource, loader)

	if err != nil {
		logger.Error("Failed while loading target manifests", zap.Error(err))
		return nil, err
	}

	strategy, summary := determineStrategy(runnings, targets, cfg.Spec.Workloads, logger)

	return &sdk.DetermineStrategyResponse{
		Strategy: strategy,
		Summary:  summary,
	}, nil
}

// BuildQuickSyncStages returns the stages for the quick sync strategy.
func (p *Plugin) BuildQuickSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildQuickSyncStagesInput) (*sdk.BuildQuickSyncStagesResponse, error) {
	return &sdk.BuildQuickSyncStagesResponse{
		Stages: buildQuickSyncPipeline(input.Request.Rollback),
	}, nil
}
