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
	"cmp"
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/toolregistry"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

const (
	defaultKubectlVersion = "1.18.2"
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

func (p *Plugin) Name() string {
	return "kubernetes"
}

func (p *Plugin) Version() string {
	return "0.0.1" // TODO
}

func (p *Plugin) FetchDefinedStages() []string {
	stages := make([]string, 0, len(AllStages))
	for _, s := range AllStages {
		stages = append(stages, string(s))
	}

	return stages
}

// FIXME
func (p *Plugin) BuildPipelineSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	return nil, nil
}

// FIXME
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

// TODO: implement pruing feature
func (p *Plugin) executeK8sSyncStage(ctx context.Context, input *sdk.ExecuteStageInput, dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]) sdk.StageStatus {
	lp := input.Client.LogPersister()
	lp.Info("Start syncing the deployment")

	cfg, err := config.DecodeYAML[*kubeconfig.KubernetesApplicationSpec](input.Request.TargetDeploymentSource.ApplicationConfig)
	if err != nil {
		lp.Errorf("Failed while decoding application config (%v)", err)
		return sdk.StageStatusFailure
	}

	// TODO: find the way to hold the tool registry and loader in the plugin.
	// Currently, we create them every time the stage is executed beucause we can't pass input.Client.toolRegistry to the plugin when starting the plugin.
	toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())
	loader := provider.NewLoader(toolRegistry)

	lp.Infof("Loading manifests at commit %s for handling", input.Request.TargetDeploymentSource.CommitHash)
	manifests, err := p.loadManifests(ctx, &input.Request.Deployment, cfg.Spec, &input.Request.TargetDeploymentSource, loader)
	if err != nil {
		lp.Errorf("Failed while loading manifests (%v)", err)
		return sdk.StageStatusFailure
	}
	lp.Successf("Successfully loaded %d manifests", len(manifests))

	// Because the loaded manifests are read-only
	// we duplicate them to avoid updating the shared manifests data in cache.
	// TODO: implement duplicateManifests function

	// When addVariantLabelToSelector is true, ensure that all workloads
	// have the variant label in their selector.
	var (
		variantLabel   = cfg.Spec.VariantLabel.Key
		primaryVariant = cfg.Spec.VariantLabel.PrimaryValue
	)
	// TODO: treat the stage options specified under "with"
	if cfg.Spec.QuickSync.AddVariantLabelToSelector {
		workloads := findWorkloadManifests(manifests, cfg.Spec.Workloads)
		for _, m := range workloads {
			if err := ensureVariantSelectorInWorkload(m, variantLabel, primaryVariant); err != nil {
				lp.Errorf("Unable to check/set %q in selector of workload %s (%v)", variantLabel+": "+primaryVariant, m.Key().ReadableString(), err)
				return sdk.StageStatusFailure
			}
		}
	}

	// Add variant annotations to all manifests.
	for i := range manifests {
		manifests[i].AddLabels(map[string]string{
			variantLabel: primaryVariant,
		})
		manifests[i].AddAnnotations(map[string]string{
			variantLabel: primaryVariant,
		})
	}

	if err := annotateConfigHash(manifests); err != nil {
		lp.Errorf("Unable to set %q annotation into the workload manifest (%v)", provider.AnnotationConfigHash, err)
		return sdk.StageStatusFailure
	}

	// Get the deploy target config.
	if len(dts) == 0 {
		lp.Error("No deploy target was found")
		return sdk.StageStatusFailure
	}
	deployTargetConfig := dts[0].Config

	// Get the kubectl tool path.
	kubectlPath, err := toolRegistry.Kubectl(ctx, cmp.Or(cfg.Spec.Input.KubectlVersion, deployTargetConfig.KubectlVersion, defaultKubectlVersion))
	if err != nil {
		lp.Errorf("Failed while getting kubectl tool (%v)", err)
		return sdk.StageStatusFailure
	}

	// Create the kubectl wrapper for the target cluster.
	kubectl := provider.NewKubectl(kubectlPath)

	// Create the applier for the target cluster.
	applier := provider.NewApplier(kubectl, cfg.Spec.Input, deployTargetConfig, input.Logger)

	// Start applying all manifests to add or update running resources.
	// TODO: use applyManifests instead of applyManifestsSDK
	if err := applyManifestsSDK(ctx, applier, manifests, cfg.Spec.Input.Namespace, lp); err != nil {
		lp.Errorf("Failed while applying manifests (%v)", err)
		return sdk.StageStatusSuccess
	}

	// TODO: treat the stage options specified under "with"
	if !cfg.Spec.QuickSync.Prune {
		lp.Info("Resource GC was skipped because sync.prune was not configured")
		return sdk.StageStatusSuccess
	}

	// Wait for all applied manifests to be stable.
	// In theory, we don't need to wait for them to be stable before going to the next step
	// but waiting for a while reduces the number of Kubernetes changes in a short time.
	lp.Info("Waiting for the applied manifests to be stable")
	select {
	case <-time.After(15 * time.Second):
		break
	case <-ctx.Done():
		break
	}

	lp.Info("Start finding all running resources but no longer defined in Git")

	namespacedLiveResources, clusterScopedLiveResources, err := provider.GetLiveResources(ctx, kubectl, deployTargetConfig.KubeConfigPath, input.Request.Deployment.ApplicationID)
	if err != nil {
		lp.Errorf("Failed while getting live resources (%v)", err)
		return sdk.StageStatusFailure
	}

	if len(namespacedLiveResources)+len(clusterScopedLiveResources) == 0 {
		lp.Info("There is no data about live resource so no resource will be removed")
		return sdk.StageStatusSuccess
	}

	lp.Successf("Successfully loaded %d live resources", len(namespacedLiveResources)+len(clusterScopedLiveResources))

	removeKeys := provider.FindRemoveResources(manifests, namespacedLiveResources, clusterScopedLiveResources)
	if len(removeKeys) == 0 {
		lp.Info("There are no live resources should be removed")
		return sdk.StageStatusSuccess
	}

	lp.Infof("Start pruning %d resources", len(removeKeys))
	var deletedCount int
	for _, key := range removeKeys {
		if err := kubectl.Delete(ctx, deployTargetConfig.KubeConfigPath, key.Namespace(), key); err != nil {
			if errors.Is(err, provider.ErrNotFound) {
				lp.Infof("Specified resource does not exist, so skip deleting the resource: %s (%v)", key.ReadableString(), err)
				continue
			}
			lp.Errorf("Failed while deleting resource %s (%v)", key.ReadableString(), err)
			continue // continue to delete other resources
		}
		deletedCount++
		lp.Successf("- deleted resource: %s", key.ReadableString())
	}

	lp.Successf("Successfully deleted %d resources", deletedCount)

	return sdk.StageStatusSuccess
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

func (p *Plugin) executeK8sRollbackStage(ctx context.Context, input *sdk.ExecuteStageInput, dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]) sdk.StageStatus {
	lp := input.Client.LogPersister()

	if input.Request.RunningDeploymentSource.CommitHash == "" {
		lp.Errorf("Unable to determine the last deployed commit to rollback. It seems this is the first deployment.")
		return sdk.StageStatusFailure
	}

	lp.Info("Start rolling back the deployment")

	cfg, err := config.DecodeYAML[*kubeconfig.KubernetesApplicationSpec](input.Request.RunningDeploymentSource.ApplicationConfig)
	if err != nil {
		lp.Errorf("Failed while decoding application config (%v)", err)
		return sdk.StageStatusFailure
	}

	lp.Infof("Loading manifests at commit %s for handling", input.Request.RunningDeploymentSource.CommitHash)
	toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())
	manifests, err := p.loadManifests(ctx, &input.Request.Deployment, cfg.Spec, &input.Request.RunningDeploymentSource, provider.NewLoader(toolRegistry))
	if err != nil {
		lp.Errorf("Failed while loading manifests (%v)", err)
		return sdk.StageStatusFailure
	}
	lp.Successf("Successfully loaded %d manifests", len(manifests))

	// Because the loaded manifests are read-only
	// we duplicate them to avoid updating the shared manifests data in cache.
	// TODO: implement duplicateManifests function

	// When addVariantLabelToSelector is true, ensure that all workloads
	// have the variant label in their selector.
	var (
		variantLabel   = cfg.Spec.VariantLabel.Key
		primaryVariant = cfg.Spec.VariantLabel.PrimaryValue
	)
	// TODO: Consider other fields to configure whether to add a variant label to the selector
	// because the rollback stage is executed in both quick sync and pipeline sync strategies.
	if cfg.Spec.QuickSync.AddVariantLabelToSelector {
		workloads := findWorkloadManifests(manifests, cfg.Spec.Workloads)
		for _, m := range workloads {
			if err := ensureVariantSelectorInWorkload(m, variantLabel, primaryVariant); err != nil {
				lp.Errorf("Unable to check/set %q in selector of workload %s (%v)", variantLabel+": "+primaryVariant, m.Key().ReadableString(), err)
				return sdk.StageStatusFailure
			}
		}
	}

	// Add variant annotations to all manifests.
	for i := range manifests {
		manifests[i].AddAnnotations(map[string]string{
			variantLabel: primaryVariant,
		})
	}

	if err := annotateConfigHash(manifests); err != nil {
		lp.Errorf("Unable to set %q annotation into the workload manifest (%v)", provider.AnnotationConfigHash, err)
		return sdk.StageStatusFailure
	}

	// Get the deploy target config.
	if len(dts) == 0 {
		lp.Error("No deploy target was found")
		return sdk.StageStatusFailure
	}
	deployTargetConfig := dts[0].Config

	// Get the kubectl tool path.
	kubectlPath, err := toolRegistry.Kubectl(ctx, cmp.Or(cfg.Spec.Input.KubectlVersion, deployTargetConfig.KubectlVersion, defaultKubectlVersion))
	if err != nil {
		lp.Errorf("Failed while getting kubectl tool (%v)", err)
		return sdk.StageStatusFailure
	}

	// Create the applier for the target cluster.
	applier := provider.NewApplier(provider.NewKubectl(kubectlPath), cfg.Spec.Input, deployTargetConfig, input.Logger)

	// Start applying all manifests to add or update running resources.
	if err := applyManifestsSDK(ctx, applier, manifests, cfg.Spec.Input.Namespace, lp); err != nil {
		lp.Errorf("Failed while applying manifests (%v)", err)
		return sdk.StageStatusFailure
	}

	// TODO: implement prune resources
	// TODO: delete all resources of CANARY variant
	// TODO: delete all resources of BASELINE variant

	return sdk.StageStatusSuccess
}

func (p *Plugin) DetermineVersions(ctx context.Context, _ *sdk.ConfigNone, _ *sdk.Client, input *sdk.DetermineVersionsInput) (*sdk.DetermineVersionsResponse, error) {
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
		Versions: determineVersionsSDK(manifests),
	}, nil
}

// FIXME
func (p *Plugin) DetermineStrategy(context.Context, *sdk.ConfigNone, *sdk.Client, sdk.TODO) (sdk.TODO, error) {
	return sdk.TODO{}, nil
}

func (p *Plugin) BuildQuickSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildQuickSyncStagesInput) (*sdk.BuildQuickSyncStagesResponse, error) {
	return &sdk.BuildQuickSyncStagesResponse{
		Stages: BuildQuickSyncPipeline(input.Request.Rollback),
	}, nil
}
