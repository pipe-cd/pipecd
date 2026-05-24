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

		namespacedLiveResources, clusterScopedLiveResources, err := provider.GetLiveResources(ctx, kubectl, tc.deployTarget.Config.KubeConfigPath, input.Request.ApplicationID)
		if err != nil {
			input.Logger.Error("Failed to get live resources", zap.Error(err))
			return nil, err
		}

		// Filter live resources by include/exclude rules from the deploy target's AppStateInformer config.
		informer := tc.deployTarget.Config.AppStateInformer
		namespacedLiveResources = filterByAppStateInformer(namespacedLiveResources, informer)
		clusterScopedLiveResources = filterByAppStateInformer(clusterScopedLiveResources, informer)

		liveState := p.makeAppLivestate(namespacedLiveResources, clusterScopedLiveResources, tc.deployTarget)
		liveStates = append(liveStates, liveState)

		liveManifests := make([]provider.Manifest, 0, len(namespacedLiveResources)+len(clusterScopedLiveResources))
		liveManifests = append(liveManifests, namespacedLiveResources...)
		liveManifests = append(liveManifests, clusterScopedLiveResources...)

		manifests, err := p.loadManifests(ctx, input, cfg.Spec, provider.NewLoader(toolRegistry), input.Logger, tc.multiTarget)
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
	// Filter out resources annotated to be ignored from drift detection.
	diffResult, err := provider.DiffList(
		filterIgnoringManifests(liveManifests),
		filterIgnoringManifests(gitManifests),
		logger,
		diff.WithEquateEmpty(),
		diff.WithIgnoreAddingMapKeys(),
		diff.WithCompareNumberAndNumericString(),
	)
	if err != nil {
		return sdk.ApplicationSyncState{}, err
	}

	return calculateSyncState(diffResult, commit, dt), nil
}

// filterIgnoringManifests removes manifests that are annotated to be excluded from drift detection.
// Resources with the annotation pipecd.dev/ignore-drift-detection=true are skipped.
func filterIgnoringManifests(manifests []provider.Manifest) []provider.Manifest {
	out := make([]provider.Manifest, 0, len(manifests))
	for _, m := range manifests {
		if m.GetAnnotations()[provider.LabelIgnoreDriftDirection] == provider.IgnoreDriftDetectionTrue {
			continue
		}
		out = append(out, m)
	}
	return out
}

// filterByAppStateInformer filters live resources based on include/exclude rules
// from the deploy target's AppStateInformer config.
// If IncludeResources is set, only resources matching at least one entry are kept.
// If ExcludeResources is set, resources matching any entry are removed.
// If neither is set, all resources are returned unchanged.
func filterByAppStateInformer(manifests []provider.Manifest, informer kubeconfig.KubernetesAppStateInformer) []provider.Manifest {
	if len(informer.IncludeResources) == 0 && len(informer.ExcludeResources) == 0 {
		return manifests
	}

	out := make([]provider.Manifest, 0, len(manifests))
	for _, m := range manifests {
		if len(informer.IncludeResources) > 0 {
			included := false
			for _, matcher := range informer.IncludeResources {
				if matchesResourceMatcher(m, matcher) {
					included = true
					break
				}
			}
			if !included {
				continue
			}
		}

		excluded := false
		for _, matcher := range informer.ExcludeResources {
			if matchesResourceMatcher(m, matcher) {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}

		out = append(out, m)
	}
	return out
}

// matchesResourceMatcher returns true if the manifest matches the given resource matcher.
// An empty APIVersion or Kind in the matcher acts as a wildcard for that field.
func matchesResourceMatcher(m provider.Manifest, matcher kubeconfig.KubernetesResourceMatcher) bool {
	if matcher.APIVersion != "" && m.APIVersion() != matcher.APIVersion {
		return false
	}
	if matcher.Kind != "" && m.Kind() != matcher.Kind {
		return false
	}
	return true
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

func (p Plugin) loadManifests(ctx context.Context, input *sdk.GetLivestateInput[kubeconfig.KubernetesApplicationSpec], spec *kubeconfig.KubernetesApplicationSpec, loader loader, logger *zap.Logger, multiTarget *kubeconfig.KubernetesMultiTarget) ([]provider.Manifest, error) {
	// override values if multiTarget has value.
	manifestPathes := spec.Input.Manifests
	kustomizeDir := ""
	if multiTarget != nil {
		if len(multiTarget.Manifests) > 0 {
			manifestPathes = multiTarget.Manifests
		}
		kustomizeDir = multiTarget.KustomizeDir
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
			provider.LabelPiped:       input.Request.PipedID,
			provider.LabelApplication: input.Request.ApplicationID,
			provider.LabelCommitHash:  input.Request.DeploymentSource.CommitHash,
		})
		manifests[i].AddAnnotations(map[string]string{
			provider.LabelManagedBy:          provider.ManagedByPiped,
			provider.LabelPiped:              input.Request.PipedID,
			provider.LabelApplication:        input.Request.ApplicationID,
			provider.LabelOriginalAPIVersion: manifests[i].APIVersion(),
			provider.LabelResourceKey:        manifests[i].Key().String(),
			provider.LabelCommitHash:         input.Request.DeploymentSource.CommitHash,
		})
	}

	return manifests, nil
}
