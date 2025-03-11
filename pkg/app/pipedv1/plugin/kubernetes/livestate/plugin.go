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

	"go.uber.org/zap"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/toolregistry"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

type Plugin struct{}

// GetLivestate implements sdk.LivestatePlugin.
func (p Plugin) GetLivestate(ctx context.Context, _ sdk.ConfigNone, deployTargets []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], input *sdk.GetLivestateInput) (*sdk.GetLivestateResponse, error) {
	if len(deployTargets) != 1 {
		return nil, fmt.Errorf("only 1 deploy target is allowed but got %d", len(deployTargets))
	}

	deployTarget := deployTargets[0]
	deployTargetConfig := deployTarget.Config

	cfg, err := config.DecodeYAML[*kubeconfig.KubernetesApplicationSpec](input.Request.DeploymentSource.ApplicationConfig)
	if err != nil {
		input.Logger.Error("Failed to decode the application spec", zap.Error(err))
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

	return &sdk.GetLivestateResponse{
		LiveState: sdk.ApplicationLiveState{
			Resources:    resourceStates,
			HealthStatus: sdk.ApplicationHealthStateUnknown, // TODO: Implement health status calculation
		},
		SyncState: sdk.ApplicationSyncState{}, // TODO: Implement sync state calculation
	}, nil
}

// Name implements sdk.LivestatePlugin.
func (p Plugin) Name() string {
	return "kubernetes" // TODO: make this constant to share with deployment plugin
}

// Version implements sdk.LivestatePlugin.
func (p Plugin) Version() string {
	return "0.0.1" // TODO: make this constant to share with deployment plugin
}
