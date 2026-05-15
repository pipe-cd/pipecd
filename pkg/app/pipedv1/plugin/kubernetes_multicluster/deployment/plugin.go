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

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/provider"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/toolregistry"
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

var _ sdk.DeploymentPlugin[sdk.ConfigNone, kubeconfig.KubernetesDeployTargetConfig, kubeconfig.KubernetesApplicationSpec] = (*Plugin)(nil)

// FetchDefinedStages returns the defined stages for this plugin.
func (p *Plugin) FetchDefinedStages() []string {
	return allStages
}

// BuildPipelineSyncStages returns the stages for the pipeline sync strategy.
func (p *Plugin) BuildPipelineSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	return &sdk.BuildPipelineSyncStagesResponse{
		Stages: buildPipelineStages(input.Request.Stages, input.Request.Rollback),
	}, nil
}

// ExecuteStage executes the stage.
func (p *Plugin) ExecuteStage(ctx context.Context, _ *sdk.ConfigNone, dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec]) (*sdk.ExecuteStageResponse, error) {
	switch input.Request.StageName {
	case StageK8sMultiSync:
		return &sdk.ExecuteStageResponse{
			Status: p.executeK8sMultiSyncStage(ctx, input, dts),
		}, nil
	case StageK8sMultiRollback:
		return &sdk.ExecuteStageResponse{
			Status: p.executeK8sMultiRollbackStage(ctx, input, dts),
		}, nil
	case StageK8sMultiCanaryRollout:
		return &sdk.ExecuteStageResponse{Status: p.executeK8sMultiCanaryRolloutStage(ctx, input, dts)}, nil
	case StageK8sMultiCanaryClean:
		return &sdk.ExecuteStageResponse{
			Status: p.executeK8sMultiCanaryCleanStage(ctx, input, dts),
		}, nil
	case StageK8sMultiPrimaryRollout:
		return &sdk.ExecuteStageResponse{Status: p.executeK8sMultiPrimaryRolloutStage(ctx, input, dts)}, nil
	case StageK8sMultiBaselineRollout:
		return &sdk.ExecuteStageResponse{Status: p.executeK8sMultiBaselineRolloutStage(ctx, input, dts)}, nil
	case StageK8sMultiBaselineClean:
		return &sdk.ExecuteStageResponse{
			Status: p.executeK8sMultiBaselineCleanStage(ctx, input, dts),
		}, nil
	case StageK8sMultiTrafficRouting:
		return &sdk.ExecuteStageResponse{Status: p.executeK8sMultiTrafficRoutingStage(ctx, input, dts)}, nil
	default:
		return nil, errors.New("unimplemented or unsupported stage")
	}
}

func (p *Plugin) loadManifests(ctx context.Context, deploy *sdk.Deployment, spec *kubeconfig.KubernetesApplicationSpec, deploymentSource *sdk.DeploymentSource[kubeconfig.KubernetesApplicationSpec], loader loader, logger *zap.Logger, multiTarget *kubeconfig.KubernetesMultiTarget) ([]provider.Manifest, error) {
	// Override manifest paths and kustomize directory from the per-target config if set.
	manifestPathes := spec.Input.Manifests
	kustomizeDir := ""
	if multiTarget != nil {
		if len(multiTarget.Manifests) > 0 {
			manifestPathes = multiTarget.Manifests
		}
		kustomizeDir = multiTarget.KustomizeDir
	}

	manifests, err := loader.LoadManifests(ctx, provider.LoaderInput{
		PipedID:          deploy.PipedID,
		AppID:            deploy.ApplicationID,
		CommitHash:       deploymentSource.CommitHash,
		AppName:          deploy.ApplicationName,
		AppDir:           deploymentSource.ApplicationDirectory,
		ConfigFilename:   deploymentSource.ApplicationConfigFilename,
		Manifests:        manifestPathes,
		Namespace:        spec.Input.Namespace,
		KustomizeVersion: spec.Input.KustomizeVersion,
		KustomizeDir:     kustomizeDir,
		KustomizeOptions: spec.Input.KustomizeOptions,
		HelmVersion:      spec.Input.HelmVersion,
		HelmChart:        spec.Input.HelmChart,
		HelmOptions:      spec.Input.HelmOptions,
		Logger:           logger,
	})

	if err != nil {
		return nil, err
	}

	// Add builtin labels and annotations for tracking application live state.
	for i := range manifests {
		manifests[i].AddLabels(map[string]string{
			provider.LabelManagedBy:   provider.ManagedByPiped,
			provider.LabelPiped:       deploy.PipedID,
			provider.LabelApplication: deploy.ApplicationID,
			provider.LabelCommitHash:  deploymentSource.CommitHash,
		})

		manifests[i].AddAnnotations(map[string]string{
			provider.LabelManagedBy:          provider.ManagedByPiped,
			provider.LabelPiped:              deploy.PipedID,
			provider.LabelApplication:        deploy.ApplicationID,
			provider.LabelOriginalAPIVersion: manifests[i].APIVersion(),
			provider.LabelResourceKey:        manifests[i].Key().String(),
			provider.LabelCommitHash:         deploymentSource.CommitHash,
		})
	}

	return manifests, nil
}

// DetermineVersions determines the versions of the application.
func (p *Plugin) DetermineVersions(ctx context.Context, _ *sdk.ConfigNone, input *sdk.DetermineVersionsInput[kubeconfig.KubernetesApplicationSpec]) (*sdk.DetermineVersionsResponse, error) {
	logger := input.Logger

	cfg, err := input.Request.DeploymentSource.AppConfig()
	if err != nil {
		logger.Error("Failed while loading application config", zap.Error(err))
		return nil, err
	}

	loader := provider.NewLoader(toolregistry.NewRegistry(input.Client.ToolRegistry()))
	multiTargets := cfg.Spec.Input.MultiTargets

	if len(multiTargets) == 0 {
		manifests, err := p.loadManifests(ctx, &input.Request.Deployment, cfg.Spec, &input.Request.DeploymentSource, loader, logger, &kubeconfig.KubernetesMultiTarget{})
		if err != nil {
			logger.Error("Failed while loading manifests", zap.Error(err))
			return nil, err
		}
		return &sdk.DetermineVersionsResponse{Versions: determineVersions(manifests)}, nil
	}

	// Collect manifests from all targets, deduplicating by resource key so the
	// same image is not counted twice when multiple targets share the same manifests.
	seen := make(map[provider.ResourceKey]struct{})
	var allManifests []provider.Manifest
	for i := range multiTargets {
		mt := &multiTargets[i]
		manifests, err := p.loadManifests(ctx, &input.Request.Deployment, cfg.Spec, &input.Request.DeploymentSource, loader, logger, mt)
		if err != nil {
			logger.Error("Failed while loading manifests", zap.String("target", mt.Target.Name), zap.Error(err))
			return nil, err
		}
		for _, m := range manifests {
			key := m.Key()
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			allManifests = append(allManifests, m)
		}
	}

	return &sdk.DetermineVersionsResponse{
		Versions: determineVersions(allManifests),
	}, nil
}

// DetermineStrategy determines the strategy for the deployment.
func (p *Plugin) DetermineStrategy(ctx context.Context, _ *sdk.ConfigNone, input *sdk.DetermineStrategyInput[kubeconfig.KubernetesApplicationSpec]) (*sdk.DetermineStrategyResponse, error) {
	logger := input.Logger
	loader := provider.NewLoader(toolregistry.NewRegistry(input.Client.ToolRegistry()))

	cfg, err := input.Request.TargetDeploymentSource.AppConfig()
	if err != nil {
		logger.Error("Failed while loading application config", zap.Error(err))
		return nil, err
	}

	multiTargets := cfg.Spec.Input.MultiTargets

	if len(multiTargets) == 0 {
		runnings, err := p.loadManifests(ctx, &input.Request.Deployment, cfg.Spec, &input.Request.RunningDeploymentSource, loader, logger, &kubeconfig.KubernetesMultiTarget{})
		if err != nil {
			logger.Error("Failed while loading running manifests", zap.Error(err))
			return nil, err
		}
		targets, err := p.loadManifests(ctx, &input.Request.Deployment, cfg.Spec, &input.Request.TargetDeploymentSource, loader, logger, &kubeconfig.KubernetesMultiTarget{})
		if err != nil {
			logger.Error("Failed while loading target manifests", zap.Error(err))
			return nil, err
		}
		strategy, summary := determineStrategy(runnings, targets, cfg.Spec.Workloads, logger)
		return &sdk.DetermineStrategyResponse{Strategy: strategy, Summary: summary}, nil
	}

	// Evaluate strategy for each configured target independently.
	// If any target requires PipelineSync, the overall deployment must use PipelineSync.
	for i := range multiTargets {
		mt := &multiTargets[i]
		runnings, err := p.loadManifests(ctx, &input.Request.Deployment, cfg.Spec, &input.Request.RunningDeploymentSource, loader, logger, mt)
		if err != nil {
			logger.Error("Failed while loading running manifests", zap.String("target", mt.Target.Name), zap.Error(err))
			return nil, err
		}
		targets, err := p.loadManifests(ctx, &input.Request.Deployment, cfg.Spec, &input.Request.TargetDeploymentSource, loader, logger, mt)
		if err != nil {
			logger.Error("Failed while loading target manifests", zap.String("target", mt.Target.Name), zap.Error(err))
			return nil, err
		}
		strategy, summary := determineStrategy(runnings, targets, cfg.Spec.Workloads, logger)
		if strategy == sdk.SyncStrategyPipelineSync {
			return &sdk.DetermineStrategyResponse{Strategy: strategy, Summary: summary}, nil
		}
	}

	return &sdk.DetermineStrategyResponse{
		Strategy: sdk.SyncStrategyQuickSync,
		Summary:  "Quick sync by applying all manifests",
	}, nil
}

// BuildQuickSyncStages returns the stages for the quick sync strategy.
func (p *Plugin) BuildQuickSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildQuickSyncStagesInput) (*sdk.BuildQuickSyncStagesResponse, error) {
	return &sdk.BuildQuickSyncStagesResponse{
		Stages: buildQuickSyncPipeline(input.Request.Rollback),
	}, nil
}
