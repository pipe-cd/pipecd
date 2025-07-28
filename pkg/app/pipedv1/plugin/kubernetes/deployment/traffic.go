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
	"encoding/json"
	"fmt"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/toolregistry"
)

func (p *Plugin) executeK8sTrafficRoutingStage(ctx context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]) sdk.StageStatus {
	lp := input.Client.LogPersister()
	lp.Info("Start routing the traffic")

	cfg, err := input.Request.TargetDeploymentSource.AppConfig()
	if err != nil {
		lp.Errorf("Failed while loading application config (%v)", err)
		return sdk.StageStatusFailure
	}

	switch kubeconfig.DetermineKubernetesTrafficRoutingMethod(cfg.Spec.TrafficRouting) {
	case kubeconfig.KubernetesTrafficRoutingMethodPodSelector:
		return p.executeK8sTrafficRoutingStagePodSelector(ctx, input, dts, cfg)
	case kubeconfig.KubernetesTrafficRoutingMethodIstio:
		return p.executeK8sTrafficRoutingStageIstio(ctx, input, dts, cfg)
	default:
		lp.Errorf("Unknown traffic routing method: %s", cfg.Spec.TrafficRouting.Method)
		return sdk.StageStatusFailure
	}
}

func (p *Plugin) executeK8sTrafficRoutingStagePodSelector(ctx context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], cfg *sdk.ApplicationConfig[kubeconfig.KubernetesApplicationSpec]) sdk.StageStatus {
	lp := input.Client.LogPersister()

	// 1. Parse stage configuration
	var stageCfg kubeconfig.K8sTrafficRoutingStageOptions
	if len(input.Request.StageConfig) == 0 {
		lp.Error("Stage config is empty, this should not happen")
		return sdk.StageStatusFailure
	}
	if err := json.Unmarshal(input.Request.StageConfig, &stageCfg); err != nil {
		lp.Errorf("Failed while unmarshalling stage config (%v)", err)
		return sdk.StageStatusFailure
	}

	// 2. Get traffic percentages
	primaryPercent, canaryPercent, baselinePercent := stageCfg.Percentages()

	// 3. Validate percentages for PodSelector method
	if baselinePercent > 0 {
		lp.Error("PodSelector method does not support baseline variant")
		return sdk.StageStatusFailure
	}

	var targetVariant string
	switch {
	case primaryPercent == 100:
		targetVariant = cfg.Spec.VariantLabel.PrimaryValue
	case canaryPercent == 100:
		targetVariant = cfg.Spec.VariantLabel.CanaryValue
	default:
		lp.Errorf("PodSelector requires either primary or canary to be 100%% (primary=%d, canary=%d)",
			primaryPercent, canaryPercent)
		return sdk.StageStatusFailure
	}

	// 4. Create tool registry and loader
	toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())
	loader := provider.NewLoader(toolRegistry)

	// 5. Load manifests
	lp.Infof("Loading manifests at commit %s", input.Request.TargetDeploymentSource.CommitHash)
	manifests, err := p.loadManifests(ctx, &input.Request.Deployment, cfg.Spec,
		&input.Request.TargetDeploymentSource, loader)
	if err != nil {
		lp.Errorf("Failed while loading manifests (%v)", err)
		return sdk.StageStatusFailure
	}
	lp.Successf("Successfully loaded %d manifests", len(manifests))

	if len(manifests) == 0 {
		lp.Error("There are no kubernetes manifests to handle")
		return sdk.StageStatusFailure
	}

	// 6. Find service manifests
	services := findManifests(provider.KindService, cfg.Spec.Service.Name, manifests)
	if len(services) == 0 {
		lp.Error("Unable to find any Service manifest")
		return sdk.StageStatusFailure
	}
	if len(services) > 1 {
		lp.Infof("Found %d Service manifests, using the first one", len(services))
	}

	// 7. Deep copy the Service manifest
	service := services[0].DeepCopy()

	// 8. Validate that Service has variant label in its selector
	variantLabel := cfg.Spec.VariantLabel.Key
	primaryVariant := cfg.Spec.VariantLabel.PrimaryValue
	if err := checkVariantSelectorInService(service, variantLabel, primaryVariant); err != nil {
		lp.Errorf("Service validation failed: %v", err)
		return sdk.StageStatusFailure
	}

	// 9. Update Service selector
	if err := updateServiceSelector(service, variantLabel, targetVariant); err != nil {
		lp.Errorf("Failed to update Service selector: %v", err)
		return sdk.StageStatusFailure
	}

	// 10. Metadata saving is not needed (already saved during Plan)

	// 11. Get deploy target and tools
	if len(dts) == 0 {
		lp.Error("No deploy target was found")
		return sdk.StageStatusFailure
	}
	deployTargetConfig := dts[0].Config

	kubectlPath, err := toolRegistry.Kubectl(ctx,
		cmp.Or(cfg.Spec.Input.KubectlVersion, deployTargetConfig.KubectlVersion))
	if err != nil {
		lp.Errorf("Failed while getting kubectl tool (%v)", err)
		return sdk.StageStatusFailure
	}

	kubectl := provider.NewKubectl(kubectlPath)
	applier := provider.NewApplier(kubectl, cfg.Spec.Input, deployTargetConfig, input.Logger)

	// 12. Apply the updated Service
	lp.Infof("Updating Service to route traffic to %s variant", targetVariant)
	if err := applyManifests(ctx, applier, []provider.Manifest{service},
		cfg.Spec.Input.Namespace, lp); err != nil {
		lp.Errorf("Failed while applying Service manifest (%v)", err)
		return sdk.StageStatusFailure
	}

	lp.Success("Successfully updated traffic routing")
	return sdk.StageStatusSuccess
}

func (p *Plugin) executeK8sTrafficRoutingStageIstio(_ context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], _ []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], _ *sdk.ApplicationConfig[kubeconfig.KubernetesApplicationSpec]) sdk.StageStatus {
	lp := input.Client.LogPersister()
	lp.Error("Traffic routing by Istio is not yet implemented")
	return sdk.StageStatusFailure
}

func checkVariantSelectorInService(m provider.Manifest, variantLabel, variant string) error {
	value, ok, err := m.NestedString("spec", "selector", variantLabel)
	if err != nil {
		return fmt.Errorf("failed to get spec.selector.%s: %w", variantLabel, err)
	}
	if !ok {
		return fmt.Errorf("missing %s key in spec.selector", variantLabel)
	}

	if value != variant {
		return fmt.Errorf("require %s but got %s for %s key in spec.selector",
			variant, value, variantLabel)
	}

	return nil
}

func updateServiceSelector(m provider.Manifest, variantLabel, targetVariant string) error {
	return m.AddStringMapValues(
		map[string]string{variantLabel: targetVariant},
		"spec", "selector",
	)
}

func findIstioVirtualServiceManifests(manifests []provider.Manifest, ref kubeconfig.K8sResourceReference) ([]provider.Manifest, error) {
	const (
		istioNetworkingGroup    = "networking.istio.io"
		istioVirtualServiceKind = "VirtualService"
	)

	if ref.Kind != "" && ref.Kind != istioVirtualServiceKind {
		return nil, fmt.Errorf("support only %q kind for VirtualService reference", istioVirtualServiceKind)
	}

	out := make([]provider.Manifest, 0, len(manifests))
	for _, m := range manifests {
		if m.GroupVersionKind().Group != istioNetworkingGroup {
			continue
		}
		if m.Kind() != istioVirtualServiceKind {
			continue
		}
		if ref.Name != "" && m.Name() != ref.Name {
			continue
		}
		out = append(out, m)
	}

	return out, nil
}
