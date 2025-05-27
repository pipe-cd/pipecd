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

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/toolregistry"
	"github.com/pipe-cd/pipecd/pkg/plugin/diff"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

type Plugin struct{}

// GetLivestate implements sdk.LivestatePlugin.
func (p Plugin) GetLivestate(ctx context.Context, _ *sdk.ConfigNone, deployTargets []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], input *sdk.GetLivestateInput[kubeconfig.KubernetesApplicationSpec]) (*sdk.GetLivestateResponse, error) {
	if len(deployTargets) != 1 {
		return nil, fmt.Errorf("only 1 deploy target is allowed but got %d", len(deployTargets))
	}

	deployTarget := deployTargets[0]
	deployTargetConfig := deployTarget.Config

	cfg, err := input.Request.DeploymentSource.AppConfig()
	if err != nil {
		input.Logger.Error("Failed while loading application config", zap.Error(err))
		return nil, err
	}

	// TODO: find the way to hold the tool registry and loader in the plugin.
	// Currently, we create them every time the stage is executed beucause we can't pass input.Client.toolRegistry to the plugin when starting the plugin.
	toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())

	// Get the kubectl tool path.
	kubectlPath, err := toolRegistry.Kubectl(ctx, cmp.Or(cfg.Spec.Input.KubectlVersion, deployTargetConfig.KubectlVersion))
	if err != nil {
		input.Logger.Error("Failed to get kubectl tool", zap.Error(err))
		return nil, err
	}

	// Create the kubectl wrapper for the target cluster.
	kubectl := provider.NewKubectl(kubectlPath)

	// TODO: We need to implement including/excluding resources.
	// ref; https://pipecd.dev/docs-v0.50.x/user-guide/managing-piped/configuration-reference/#kubernetesappstateinformer
	namespacedLiveResources, clusterScopedLiveResources, err := provider.GetLiveResources(ctx, kubectl, deployTargetConfig.KubeConfigPath, input.Request.ApplicationID)
	if err != nil {
		input.Logger.Error("Failed to get live resources", zap.Error(err))
		return nil, err
	}

	resourceStates := make([]sdk.ResourceState, 0, len(namespacedLiveResources)+len(clusterScopedLiveResources))
	for _, m := range namespacedLiveResources {
		resourceStates = append(resourceStates, m.ToResourceState(deployTarget.Name))
	}
	for _, m := range clusterScopedLiveResources {
		resourceStates = append(resourceStates, m.ToResourceState(deployTarget.Name))
	}

	manifests, err := p.loadManifests(ctx, input, cfg.Spec, provider.NewLoader(toolRegistry))
	if err != nil {
		input.Logger.Error("Failed to load manifests", zap.Error(err))
		return nil, err
	}

	liveManifests := make([]provider.Manifest, 0, len(namespacedLiveResources)+len(clusterScopedLiveResources))
	liveManifests = append(liveManifests, namespacedLiveResources...)
	liveManifests = append(liveManifests, clusterScopedLiveResources...)

	// Calculate SyncState by comparing live manifests with desired manifests
	// TODO: Implement drift detection ignore configs
	diffResult, err := provider.DiffList(liveManifests, manifests, input.Logger,
		diff.WithEquateEmpty(),
		diff.WithIgnoreAddingMapKeys(),
		diff.WithCompareNumberAndNumericString(),
	)
	if err != nil {
		input.Logger.Error("Failed to calculate diff", zap.Error(err))
		return nil, err
	}

	syncState := calculateSyncState(diffResult, input.Request.DeploymentSource.CommitHash)

	return &sdk.GetLivestateResponse{
		LiveState: sdk.ApplicationLiveState{
			Resources: resourceStates,
		},
		SyncState: syncState,
	}, nil
}

func calculateSyncState(diffResult *provider.DiffListResult, commit string) sdk.ApplicationSyncState {
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
	b.WriteString(fmt.Sprintf("Diff between the defined state in Git at commit %s and actual state in cluster:\n\n", commit))
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

type loader interface {
	// LoadManifests renders and loads all manifests for application.
	LoadManifests(ctx context.Context, input provider.LoaderInput) ([]provider.Manifest, error)
}

// TODO: share this implementation with the deployment plugin
func (p Plugin) loadManifests(ctx context.Context, input *sdk.GetLivestateInput[kubeconfig.KubernetesApplicationSpec], spec *kubeconfig.KubernetesApplicationSpec, loader loader) ([]provider.Manifest, error) {
	manifests, err := loader.LoadManifests(ctx, provider.LoaderInput{
		PipedID:          input.Request.PipedID,
		AppID:            input.Request.ApplicationID,
		CommitHash:       input.Request.DeploymentSource.CommitHash,
		AppName:          input.Request.ApplicationName,
		AppDir:           input.Request.DeploymentSource.ApplicationDirectory,
		ConfigFilename:   input.Request.DeploymentSource.ApplicationConfigFilename,
		Manifests:        spec.Input.Manifests,
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
