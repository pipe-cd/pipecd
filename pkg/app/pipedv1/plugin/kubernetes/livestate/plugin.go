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
	"context"
	"fmt"
	"strings"
	"sync"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime/schema"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/pipe-cd/piped-plugin-sdk-go/diff"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/livestate/store"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/toolregistry"
)

var (
	_ sdk.LivestatePlugin[sdk.ConfigNone, kubeconfig.KubernetesDeployTargetConfig, kubeconfig.KubernetesApplicationSpec] = (*Plugin)(nil)
	_ sdk.Initializer[sdk.ConfigNone, kubeconfig.KubernetesDeployTargetConfig]                                           = (*Plugin)(nil)
)

type Plugin struct {
	store       *store.Store
	initialized sync.Once
}

// Initialize implements sdk.Initializer.
func (p *Plugin) Initialize(ctx context.Context, input *sdk.InitializeInput[sdk.ConfigNone, kubeconfig.KubernetesDeployTargetConfig]) error {
	var err error

	p.initialized.Do(func() {
		p.store, err = store.Run(ctx, input.DeployTargets, input.Logger)
		if err != nil {
			err = fmt.Errorf("failed to run livestate store: %w", err)
		}
	})

	return err
}

// GetLivestate implements sdk.LivestatePlugin.
func (p *Plugin) GetLivestate(ctx context.Context, _ *sdk.ConfigNone, deployTargets []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], input *sdk.GetLivestateInput[kubeconfig.KubernetesApplicationSpec]) (*sdk.GetLivestateResponse, error) {
	if len(deployTargets) != 1 {
		return nil, fmt.Errorf("only 1 deploy target is allowed but got %d", len(deployTargets))
	}

	deployTarget := deployTargets[0]

	managedResources, err := p.store.ManagedResources(ctx, deployTarget.Name, input.Request.ApplicationID)
	if err != nil {
		input.Logger.Error("Failed to get livestate", zap.Error(err))
		return nil, err
	}

	watchingResourceKinds, err := p.store.WatchingResourceKinds(deployTarget.Name)
	if err != nil {
		input.Logger.Error("Failed to get watching resource kinds", zap.Error(err))
		return nil, err
	}

	cfg, err := input.Request.DeploymentSource.AppConfig()
	if err != nil {
		input.Logger.Error("Failed while loading application config", zap.Error(err))
		return nil, err
	}

	// TODO: find the way to hold the tool registry and loader in the plugin.
	// Currently, we create them every time the stage is executed beucause we can't pass input.Client.toolRegistry to the plugin when starting the plugin.
	toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())

	manifests, err := p.loadManifests(ctx, input, cfg.Spec, provider.NewLoader(toolRegistry))
	if err != nil {
		input.Logger.Error("Failed to load manifests", zap.Error(err))
		return nil, err
	}

	// Calculate SyncState by comparing live manifests with desired manifests
	// TODO: Implement drift detection ignore configs
	diffResult, err := provider.DiffList(
		filterIgnoringManifests(managedResources), // live manifests are already filtered by watchingResourceKinds in the store
		filterIgnoringManifests(onlyWatchingResourceKinds(manifests, watchingResourceKinds)),
		input.Logger,
		diff.WithEquateEmpty(),
		diff.WithIgnoreAddingMapKeys(),
		diff.WithCompareNumberAndNumericString(),
	)
	if err != nil {
		input.Logger.Error("Failed to calculate diff", zap.Error(err))
		return nil, err
	}

	liveManifests, err := p.store.Livestate(ctx, deployTarget.Name, input.Request.ApplicationID)
	if err != nil {
		input.Logger.Error("Failed to get livestate", zap.Error(err))
		return nil, err
	}

	resourceStates := make([]sdk.ResourceState, 0, len(liveManifests))
	for _, manifest := range liveManifests {
		resourceStates = append(resourceStates, manifest.ToResourceState(deployTarget.Name))
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

func filterIgnoringManifests(manifests []provider.Manifest) []provider.Manifest {
	out := make([]provider.Manifest, 0, len(manifests))
	for _, m := range manifests {
		annotations := m.GetAnnotations()
		if annotations[provider.LabelIgnoreDriftDirection] == provider.IgnoreDriftDetectionTrue {
			continue
		}
		out = append(out, m)
	}
	return out
}

func onlyWatchingResourceKinds(manifests []provider.Manifest, watchingResourceKinds []schema.GroupVersionKind) []provider.Manifest {
	watchingMap := make(map[schema.GroupVersionKind]struct{}, len(watchingResourceKinds))
	for _, k := range watchingResourceKinds {
		watchingMap[k] = struct{}{}
	}

	filtered := make([]provider.Manifest, 0, len(manifests))
	for _, m := range manifests {
		_, ok := watchingMap[m.GroupVersionKind()]
		if ok {
			filtered = append(filtered, m)
		}
	}
	return filtered
}

type loader interface {
	// LoadManifests renders and loads all manifests for application.
	LoadManifests(ctx context.Context, input provider.LoaderInput) ([]provider.Manifest, error)
}

// TODO: share this implementation with the deployment plugin
func (p *Plugin) loadManifests(ctx context.Context, input *sdk.GetLivestateInput[kubeconfig.KubernetesApplicationSpec], spec *kubeconfig.KubernetesApplicationSpec, loader loader) ([]provider.Manifest, error) {
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
