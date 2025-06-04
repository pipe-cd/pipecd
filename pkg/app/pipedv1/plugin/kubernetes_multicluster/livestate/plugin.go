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

package livestate

import (
	"cmp"
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/pipe-cd/piped-plugin-sdk-go/diff"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/provider"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/toolregistry"
)

type Plugin struct{}

func (p Plugin) GetLivestate(ctx context.Context, _ *sdk.ConfigNone, deployTargets []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], input *sdk.GetLivestateInput[kubeconfig.KubernetesApplicationSpec]) (*sdk.GetLivestateResponse, error) {
	cfg, err := input.Request.DeploymentSource.AppConfig()
	if err != nil {
		input.Logger.Error("Failed to load application config", zap.Error(err))
		return nil, err
	}

	logger := input.Logger

	type targetConfig struct {
		deployTarget *sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]
		multiTarget  *kubeconfig.KubernetesMultiTarget
	}

	deployTargetMap := make(map[string]*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], 0)
	targetConfigs := make([]targetConfig, 0, len(deployTargets))

	// prevent the deployment when its deployTarget is not found in the piped config
	for _, target := range deployTargets {
		deployTargetMap[target.Name] = target
	}

	// If no multi-targets are specified, sync to all deploy targets.
	if len(cfg.Spec.Input.MultiTargets) == 0 {
		for _, dt := range deployTargets {
			targetConfigs = append(targetConfigs, targetConfig{
				deployTarget: dt,
				multiTarget:  nil,
			})
		}
	} else {
		// Sync to the specified multi-targets.
		for _, multiTarget := range cfg.Spec.Input.MultiTargets {
			dt, ok := deployTargetMap[multiTarget.Target.Name]
			if !ok {
				logger.Info("Ignore multi target '%s': not matched any deployTarget", zap.String("multiTargetName", multiTarget.Target.Name))
				continue
			}

			targetConfigs = append(targetConfigs, targetConfig{
				deployTarget: dt,
				multiTarget:  &multiTarget,
			})
		}
	}

	// TODO: find the way to hold the tool registry and loader in the plugin.
	// Currently, we create them every time the stage is executed beucause we can't pass input.Client.toolRegistry to the plugin when starting the plugin.
	toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())

	liveStates := make([]sdk.ApplicationLiveState, 0, len(targetConfigs))
	syncStates := make([]sdk.ApplicationSyncState, 0, len(targetConfigs))
	for _, tc := range targetConfigs {
		// Get the kubectl tool path.
		kubectlPath, err := toolRegistry.Kubectl(ctx, cmp.Or(cfg.Spec.Input.KubectlVersion, tc.deployTarget.Config.KubectlVersion))
		if err != nil {
			input.Logger.Error("Failed to get kubectl tool", zap.Error(err))
			return nil, err
		}

		// Create the kubectl wrapper for the target cluster.
		kubectl := provider.NewKubectl(kubectlPath)

		// TODO: We need to implement including/excluding resources.
		// ref; https://pipecd.dev/docs-v0.50.x/user-guide/managing-piped/configuration-reference/#kubernetesappstateinformer
		namespacedLiveResources, clusterScopedLiveResources, err := provider.GetLiveResources(ctx, kubectl, tc.deployTarget.Config.KubeConfigPath, input.Request.ApplicationID)
		if err != nil {
			input.Logger.Error("Failed to get live resources", zap.Error(err))
			return nil, err
		}

		liveState := p.makeAppLivestate(namespacedLiveResources, clusterScopedLiveResources, tc.deployTarget)
		liveStates = append(liveStates, liveState)

		liveManifests := make([]provider.Manifest, 0, len(namespacedLiveResources)+len(clusterScopedLiveResources))
		liveManifests = append(liveManifests, namespacedLiveResources...)
		liveManifests = append(liveManifests, clusterScopedLiveResources...)

		manifests, err := p.loadManifests(ctx, input, cfg.Spec, provider.NewLoader(toolRegistry), tc.multiTarget)
		if err != nil {
			input.Logger.Error("Failed to load manifests", zap.Error(err))
			return nil, err
		}

		syncState, err := p.makeAppSyncState(liveManifests, manifests, tc.deployTarget, input.Request.DeploymentSource.CommitHash, logger)
		if err != nil {
			input.Logger.Error("Failed to make app sync state", zap.Error(err), zap.String("deployTarget", tc.deployTarget.Name))
			return nil, err
		}
		syncStates = append(syncStates, syncState)
	}

	appLiveState := sdk.ApplicationLiveState{}
	for _, ls := range liveStates {
		appLiveState.Resources = append(appLiveState.Resources, ls.Resources...)
	}

	appSyncState := sdk.ApplicationSyncState{}
	for _, ss := range syncStates {
		appSyncState.Reason = fmt.Sprintf("%s\n%s", appSyncState.Reason, ss.Reason)
		appSyncState.ShortReason = fmt.Sprintf("%s\n%s", appSyncState.ShortReason, ss.ShortReason)
	}
	appSyncState.Status = calculateSyncStatus(syncStates)

	return &sdk.GetLivestateResponse{
		LiveState: appLiveState,
		SyncState: appSyncState,
	}, nil
}

func (p Plugin) makeAppLivestate(namespacedLiveResources, clusterScopedLiveResources []provider.Manifest, dt *sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]) sdk.ApplicationLiveState {
	resourceStates := make([]sdk.ResourceState, 0, len(namespacedLiveResources)+len(clusterScopedLiveResources))
	for _, m := range namespacedLiveResources {
		resourceStates = append(resourceStates, m.ToResourceState(dt.Name))
	}
	for _, m := range clusterScopedLiveResources {
		resourceStates = append(resourceStates, m.ToResourceState(dt.Name))
	}

	return sdk.ApplicationLiveState{
		Resources: resourceStates,
	}
}

func (p Plugin) makeAppSyncState(liveManifests, gitManifests []provider.Manifest, dt *sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], commit string, logger *zap.Logger) (sdk.ApplicationSyncState, error) {
	// Calculate SyncState by comparing live manifests with desired manifests
	// TODO: Implement drift detection ignore configs
	diffResult, err := provider.DiffList(liveManifests, gitManifests, logger,
		diff.WithEquateEmpty(),
		diff.WithIgnoreAddingMapKeys(),
		diff.WithCompareNumberAndNumericString(),
	)
	if err != nil {
		return sdk.ApplicationSyncState{}, err
	}

	return calculateSyncState(diffResult, commit, dt), nil
}

func calculateSyncState(diffResult *provider.DiffListResult, commit string, dt *sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]) sdk.ApplicationSyncState {
	if diffResult.NoChanges() {
		return sdk.ApplicationSyncState{
			Status:      sdk.ApplicationSyncStateSynced,
			ShortReason: "",
			Reason:      "",
		}
	}

	total := len(diffResult.Adds) + len(diffResult.Deletes) + len(diffResult.Changes)
	shortReason := fmt.Sprintf("There are %d manifests not synced (%d adds, %d deletes, %d changes)",
		total,
		len(diffResult.Adds),
		len(diffResult.Deletes),
		len(diffResult.Changes),
	)

	if len(commit) > 7 {
		commit = commit[:7]
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("Diff between the defined state in Git at commit %s and actual state in cluster: %s\n\n", commit, dt.Name))
	b.WriteString("--- Actual   (LiveState)\n+++ Expected (Git)\n\n")

	details := diffResult.Render(provider.DiffRenderOptions{
		MaskSecret:          true,
		MaskConfigMap:       true,
		MaxChangedManifests: 3,
	})
	b.WriteString(details)

	return sdk.ApplicationSyncState{
		Status:      sdk.ApplicationSyncStateOutOfSync,
		ShortReason: shortReason,
		Reason:      b.String(),
	}
}

// calculateSyncStatus returns the highest-priority sync status among the given states.
// Priority: InvalidConfig > Unknown > OutOfSync > Synced.
func calculateSyncStatus(states []sdk.ApplicationSyncState) sdk.ApplicationSyncStatus {
	var (
		hasInvalidConfig bool
		hasUnknown       bool
		hasOutOfSync     bool
	)
	for _, state := range states {
		switch state.Status {
		case sdk.ApplicationSyncStateInvalidConfig:
			hasInvalidConfig = true
		case sdk.ApplicationSyncStateUnknown:
			hasUnknown = true
		case sdk.ApplicationSyncStateOutOfSync:
			hasOutOfSync = true
		}
	}

	if hasInvalidConfig {
		return sdk.ApplicationSyncStateInvalidConfig
	}

	if hasUnknown {
		return sdk.ApplicationSyncStateUnknown
	}

	if hasOutOfSync {
		return sdk.ApplicationSyncStateOutOfSync
	}

	return sdk.ApplicationSyncStateSynced
}

type loader interface {
	// LoadManifests renders and loads all manifests for application.
	LoadManifests(ctx context.Context, input provider.LoaderInput) ([]provider.Manifest, error)
}

// TODO: share this implementation with the deployment plugin
func (p Plugin) loadManifests(ctx context.Context, input *sdk.GetLivestateInput[kubeconfig.KubernetesApplicationSpec], spec *kubeconfig.KubernetesApplicationSpec, loader loader, multiTarget *kubeconfig.KubernetesMultiTarget) ([]provider.Manifest, error) {
	// override values if multiTarget has value.
	manifestPathes := spec.Input.Manifests
	if multiTarget != nil {
		if len(multiTarget.Manifests) > 0 {
			manifestPathes = multiTarget.Manifests
		}
	}

	manifests, err := loader.LoadManifests(ctx, provider.LoaderInput{
		PipedID:          input.Request.PipedID,
		AppID:            input.Request.ApplicationID,
		CommitHash:       input.Request.DeploymentSource.CommitHash,
		AppName:          input.Request.ApplicationName,
		AppDir:           input.Request.DeploymentSource.ApplicationDirectory,
		ConfigFilename:   input.Request.DeploymentSource.ApplicationConfigFilename,
		Manifests:        manifestPathes,
		Namespace:        spec.Input.Namespace,
		KustomizeVersion: spec.Input.KustomizeVersion,
		KustomizeOptions: spec.Input.KustomizeOptions,
		HelmVersion:      spec.Input.HelmVersion,
		HelmChart:        spec.Input.HelmChart,
		HelmOptions:      spec.Input.HelmOptions,
	})

	if err != nil {
		return nil, err
	}

	return manifests, nil
}
