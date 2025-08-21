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
	istiov1 "istio.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
		&input.Request.TargetDeploymentSource, loader, input.Logger)
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

func (p *Plugin) executeK8sTrafficRoutingStageIstio(ctx context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], cfg *sdk.ApplicationConfig[kubeconfig.KubernetesApplicationSpec]) sdk.StageStatus {
	lp := input.Client.LogPersister()

	var stageCfg kubeconfig.K8sTrafficRoutingStageOptions
	if len(input.Request.StageConfig) == 0 {
		lp.Error("Stage config is empty, this should not happen")
		return sdk.StageStatusFailure
	}
	if err := json.Unmarshal(input.Request.StageConfig, &stageCfg); err != nil {
		lp.Errorf("Failed while unmarshalling stage config (%v)", err)
		return sdk.StageStatusFailure
	}

	toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())
	loader := provider.NewLoader(toolRegistry)

	lp.Infof("Loading manifests at commit %s", input.Request.TargetDeploymentSource.CommitHash)
	manifests, err := p.loadManifests(ctx, &input.Request.Deployment, cfg.Spec,
		&input.Request.TargetDeploymentSource, loader, input.Logger)
	if err != nil {
		lp.Errorf("Failed while loading manifests (%v)", err)
		return sdk.StageStatusFailure
	}
	lp.Successf("Successfully loaded %d manifests", len(manifests))

	if len(manifests) == 0 {
		lp.Error("There are no kubernetes manifests to handle")
		return sdk.StageStatusFailure
	}

	// Use VirtualService reference from Istio config if specified, otherwise use Service reference
	var vsRef kubeconfig.K8sResourceReference
	if cfg.Spec.TrafficRouting != nil && cfg.Spec.TrafficRouting.Istio != nil && cfg.Spec.TrafficRouting.Istio.VirtualService.Name != "" {
		vsRef = cfg.Spec.TrafficRouting.Istio.VirtualService
	} else {
		vsRef = cfg.Spec.Service
	}

	virtualServices, err := findIstioVirtualServiceManifests(manifests, vsRef)
	if err != nil {
		lp.Errorf("Failed while finding traffic routing manifest: (%v)", err)
		return sdk.StageStatusFailure
	}
	if len(virtualServices) == 0 {
		lp.Error("Unable to find any VirtualService manifest")
		return sdk.StageStatusFailure
	}

	if len(virtualServices) > 1 {
		lp.Infof("Found %d VirtualService manifests, using the first one", len(virtualServices))
	}
	virtualService := virtualServices[0]

	primaryPercent, canaryPercent, baselinePercent := stageCfg.Percentages()

	vs, err := generateVirtualServiceManifest(virtualService, cfg.Spec.Service.Name, cfg.Spec.TrafficRouting.Istio.EditableRoutes, cfg.Spec.VariantLabel, int32(canaryPercent), int32(baselinePercent))
	if err != nil {
		lp.Errorf("Failed while generating VirtualService manifest: (%v)", err)
		return sdk.StageStatusFailure
	}

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

	lp.Infof("Start updating traffic routing to be percentages: primary=%d, canary=%d, baseline=%d",
		primaryPercent,
		canaryPercent,
		baselinePercent,
	)

	if err := applyManifests(ctx, applier, []provider.Manifest{vs},
		cfg.Spec.Input.Namespace, lp); err != nil {
		lp.Errorf("Failed while applying VirtualService manifest (%v)", err)
		return sdk.StageStatusFailure
	}

	lp.Success("Successfully updated traffic routing")
	return sdk.StageStatusSuccess
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

// virtualService is a wrapper around istiov1.VirtualService
// to enable the conversion from unstructured.Unstructured to VirtualService
// We can use this struct across APIVersion v1alpha3 and v1beta1 because v1.VirtualService
// and v1beta1.VirtualService are type alias of v1alpha3.VirtualService.
type virtualService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              istiov1.VirtualService
}

func convertVirtualService(m provider.Manifest) (*virtualService, error) {
	var vs virtualService
	if err := m.ConvertToStructuredObject(&vs); err != nil {
		return nil, err
	}
	return &vs, nil
}

func (vs *virtualService) toManifest() (provider.Manifest, error) {
	return provider.FromStructuredObject(vs)
}

// generateVirtualServiceManifest generates a new VirtualService manifest
// that routes traffic to the specified host with the given percentages.
// It also supports the editableRoutes parameter to specify the routes that
// can be edited by the user.
func generateVirtualServiceManifest(m provider.Manifest, host string, editableRoutes []string, variantLabel kubeconfig.KubernetesVariantLabel, canaryPercent, baselinePercent int32) (provider.Manifest, error) {
	vs, err := convertVirtualService(m)
	if err != nil {
		return provider.Manifest{}, err
	}

	editableMap := make(map[string]struct{}, len(editableRoutes))
	for _, r := range editableRoutes {
		editableMap[r] = struct{}{}
	}

	for _, http := range vs.Spec.Http {
		if len(editableMap) > 0 {
			if _, ok := editableMap[http.Name]; !ok {
				continue
			}
		}

		// Calculate the weight of the other host
		var (
			otherHostWeight int32
			otherHostRoutes = make([]*istiov1.HTTPRouteDestination, 0)
		)
		for _, r := range http.Route {
			if r.Destination != nil && r.Destination.Host != host {
				otherHostWeight += r.Weight
				otherHostRoutes = append(otherHostRoutes, r)
			}
		}

		// Calculate the weight of the variants
		var (
			variantsWeight = 100 - otherHostWeight
			canaryWeight   = canaryPercent * variantsWeight / 100
			baselineWeight = baselinePercent * variantsWeight / 100
			primaryWeight  = variantsWeight - canaryWeight - baselineWeight
			routes         = make([]*istiov1.HTTPRouteDestination, 0, len(otherHostRoutes)+3)
		)

		// Add the primary route
		routes = append(routes, &istiov1.HTTPRouteDestination{
			Destination: &istiov1.Destination{
				Host:   host,
				Subset: variantLabel.PrimaryValue,
			},
			Weight: primaryWeight,
		})

		// Add the canary route
		if canaryWeight > 0 {
			routes = append(routes, &istiov1.HTTPRouteDestination{
				Destination: &istiov1.Destination{
					Host:   host,
					Subset: variantLabel.CanaryValue,
				},
				Weight: canaryWeight,
			})
		}

		// Add the baseline route
		if baselineWeight > 0 {
			routes = append(routes, &istiov1.HTTPRouteDestination{
				Destination: &istiov1.Destination{
					Host:   host,
					Subset: variantLabel.BaselineValue,
				},
				Weight: baselineWeight,
			})
		}

		routes = append(routes, otherHostRoutes...)
		http.Route = routes
	}

	return vs.toManifest()
}
