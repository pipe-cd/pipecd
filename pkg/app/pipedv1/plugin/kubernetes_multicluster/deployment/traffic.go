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
	"strings"

	"golang.org/x/sync/errgroup"
	istiov1 "istio.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/provider"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/toolregistry"
)

func (p *Plugin) executeK8sMultiTrafficRoutingStage(ctx context.Context, input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec], dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]) sdk.StageStatus {
	lp := input.Client.LogPersister()
	lp.Info("Start routing the traffic")

	cfg, err := input.Request.TargetDeploymentSource.AppConfig()
	if err != nil {
		lp.Errorf("Failed while loading application config (%v)", err)
		return sdk.StageStatusFailure
	}

	switch kubeconfig.DetermineKubernetesTrafficRoutingMethod(cfg.Spec.TrafficRouting) {
	case kubeconfig.KubernetesTrafficRoutingMethodPodSelector:
		return p.executeK8sMultiTrafficRoutingStagePodSelector(ctx, input, dts, cfg)
	case kubeconfig.KubernetesTrafficRoutingMethodIstio:
		return p.executeK8sMultiTrafficRoutingStageIstio(ctx, input, dts, cfg)
	default:
		lp.Errorf("Unknown traffic routing method: %s", cfg.Spec.TrafficRouting.Method)
		return sdk.StageStatusFailure
	}
}

func (p *Plugin) executeK8sMultiTrafficRoutingStagePodSelector(
	ctx context.Context,
	input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec],
	dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig],
	cfg *sdk.ApplicationConfig[kubeconfig.KubernetesApplicationSpec],
) sdk.StageStatus {
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

	primaryPercent, canaryPercent, baselinePercent := stageCfg.Percentages()

	// PodSelector does not support baseline and requires one variant to be 100%.
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

	lp.Infof("Routing traffic to %s variant (%s)", targetVariant, stageCfg.DisplayString())

	deployTargetMap := make(map[string]*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], len(dts))
	for _, dt := range dts {
		deployTargetMap[dt.Name] = dt
	}

	type targetConfig struct {
		deployTarget *sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]
		multiTarget  *kubeconfig.KubernetesMultiTarget
	}

	targetConfigs := make([]targetConfig, 0, len(dts))
	if len(cfg.Spec.Input.MultiTargets) == 0 {
		for _, dt := range dts {
			targetConfigs = append(targetConfigs, targetConfig{deployTarget: dt})
		}
	} else {
		for _, mt := range cfg.Spec.Input.MultiTargets {
			dt, ok := deployTargetMap[mt.Target.Name]
			if !ok {
				lp.Infof("Ignore multi target '%s': not matched any deployTarget", mt.Target.Name)
				continue
			}
			targetConfigs = append(targetConfigs, targetConfig{deployTarget: dt, multiTarget: &mt})
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, tc := range targetConfigs {
		eg.Go(func() error {
			lp.Infof("Start updating traffic routing on target %s", tc.deployTarget.Name)
			if err := p.podSelectorTrafficRouting(ctx, input, tc.deployTarget, tc.multiTarget, cfg, targetVariant); err != nil {
				return fmt.Errorf("failed to update traffic routing on target %s: %w", tc.deployTarget.Name, err)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		lp.Errorf("Failed while updating traffic routing (%v)", err)
		return sdk.StageStatusFailure
	}

	lp.Success("Successfully updated traffic routing")
	return sdk.StageStatusSuccess
}

func (p *Plugin) podSelectorTrafficRouting(
	ctx context.Context,
	input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec],
	dt *sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig],
	multiTarget *kubeconfig.KubernetesMultiTarget,
	cfg *sdk.ApplicationConfig[kubeconfig.KubernetesApplicationSpec],
	targetVariant string,
) error {
	lp := input.Client.LogPersister()
	appCfg := cfg.Spec

	toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())
	loader := provider.NewLoader(toolRegistry)

	manifests, err := p.loadManifests(ctx, &input.Request.Deployment, cfg.Spec, &input.Request.TargetDeploymentSource, loader, input.Logger, multiTarget)
	if err != nil {
		return fmt.Errorf("failed while loading manifests: %w", err)
	}

	if len(manifests) == 0 {
		return fmt.Errorf("no kubernetes manifests found")
	}

	services := findManifests(provider.KindService, appCfg.Service.Name, manifests)
	if len(services) == 0 {
		return fmt.Errorf("unable to find any Service manifest")
	}
	if len(services) > 1 {
		lp.Infof("Found %d Service manifests on target %s, using the first one", len(services), dt.Name)
	}

	service := services[0].DeepCopy()

	variantLabel := appCfg.VariantLabel.Key
	primaryVariant := appCfg.VariantLabel.PrimaryValue
	if err := checkVariantSelectorInService(service, variantLabel, primaryVariant); err != nil {
		return fmt.Errorf("service validation failed: %w", err)
	}

	if err := updateServiceSelector(service, variantLabel, targetVariant); err != nil {
		return fmt.Errorf("failed to update Service selector: %w", err)
	}

	// Resolve kubectl version: multiTarget > spec > deployTarget
	kubectlVersion := cmp.Or(appCfg.Input.KubectlVersion, dt.Config.KubectlVersion)
	if multiTarget != nil {
		kubectlVersion = cmp.Or(multiTarget.KubectlVersion, kubectlVersion)
	}

	kubectlPath, err := toolRegistry.Kubectl(ctx, kubectlVersion)
	if err != nil {
		return fmt.Errorf("failed while getting kubectl tool: %w", err)
	}

	kubectl := provider.NewKubectl(kubectlPath)
	applier := provider.NewApplier(kubectl, appCfg.Input, dt.Config, input.Logger)

	if err := applyManifests(ctx, applier, []provider.Manifest{service}, appCfg.Input.Namespace, lp); err != nil {
		return fmt.Errorf("failed while applying Service manifest: %w", err)
	}

	return nil
}

func (p *Plugin) executeK8sMultiTrafficRoutingStageIstio(
	ctx context.Context,
	input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec],
	dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig],
	cfg *sdk.ApplicationConfig[kubeconfig.KubernetesApplicationSpec],
) sdk.StageStatus {
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

	primaryPercent, canaryPercent, baselinePercent := stageCfg.Percentages()
	lp.Infof("Updating traffic routing: %s", stageCfg.DisplayString())

	deployTargetMap := make(map[string]*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], len(dts))
	for _, dt := range dts {
		deployTargetMap[dt.Name] = dt
	}

	type targetConfig struct {
		deployTarget *sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig]
		multiTarget  *kubeconfig.KubernetesMultiTarget
	}

	targetConfigs := make([]targetConfig, 0, len(dts))
	if len(cfg.Spec.Input.MultiTargets) == 0 {
		for _, dt := range dts {
			targetConfigs = append(targetConfigs, targetConfig{deployTarget: dt})
		}
	} else {
		for _, mt := range cfg.Spec.Input.MultiTargets {
			dt, ok := deployTargetMap[mt.Target.Name]
			if !ok {
				lp.Infof("Ignore multi target '%s': not matched any deployTarget", mt.Target.Name)
				continue
			}
			targetConfigs = append(targetConfigs, targetConfig{deployTarget: dt, multiTarget: &mt})
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, tc := range targetConfigs {
		eg.Go(func() error {
			lp.Infof("Start updating Istio traffic routing on target %s", tc.deployTarget.Name)
			if err := p.istioTrafficRouting(ctx, input, tc.deployTarget, tc.multiTarget, cfg, int32(canaryPercent), int32(baselinePercent)); err != nil {
				return fmt.Errorf("failed to update Istio traffic routing on target %s: %w", tc.deployTarget.Name, err)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		lp.Errorf("Failed while updating Istio traffic routing (%v)", err)
		return sdk.StageStatusFailure
	}

	lp.Successf("Successfully updated Istio traffic routing: primary=%d%%, canary=%d%%, baseline=%d%%",
		primaryPercent, canaryPercent, baselinePercent)
	return sdk.StageStatusSuccess
}

func (p *Plugin) istioTrafficRouting(
	ctx context.Context,
	input *sdk.ExecuteStageInput[kubeconfig.KubernetesApplicationSpec],
	dt *sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig],
	multiTarget *kubeconfig.KubernetesMultiTarget,
	cfg *sdk.ApplicationConfig[kubeconfig.KubernetesApplicationSpec],
	canaryPercent, baselinePercent int32,
) error {
	lp := input.Client.LogPersister()
	appCfg := cfg.Spec

	toolRegistry := toolregistry.NewRegistry(input.Client.ToolRegistry())
	loader := provider.NewLoader(toolRegistry)

	manifests, err := p.loadManifests(ctx, &input.Request.Deployment, cfg.Spec, &input.Request.TargetDeploymentSource, loader, input.Logger, multiTarget)
	if err != nil {
		return fmt.Errorf("failed while loading manifests: %w", err)
	}

	if len(manifests) == 0 {
		return fmt.Errorf("no kubernetes manifests found")
	}

	// Determine which VirtualService to use: from Istio config or fall back to service reference.
	var vsRef kubeconfig.K8sResourceReference
	if appCfg.TrafficRouting != nil && appCfg.TrafficRouting.Istio != nil && appCfg.TrafficRouting.Istio.VirtualService.Name != "" {
		vsRef = appCfg.TrafficRouting.Istio.VirtualService
	} else {
		vsRef = appCfg.Service
	}

	virtualServices, err := findIstioVirtualServiceManifests(manifests, vsRef)
	if err != nil {
		return fmt.Errorf("failed while finding VirtualService manifest: %w", err)
	}
	if len(virtualServices) == 0 {
		return fmt.Errorf("unable to find any VirtualService manifest")
	}
	if len(virtualServices) > 1 {
		lp.Infof("Found %d VirtualService manifests on target %s, using the first one", len(virtualServices), dt.Name)
	}

	vs, err := generateVirtualServiceManifest(
		virtualServices[0],
		appCfg.Service.Name,
		appCfg.TrafficRouting.Istio.EditableRoutes,
		appCfg.VariantLabel,
		canaryPercent,
		baselinePercent,
	)
	if err != nil {
		return fmt.Errorf("failed while generating VirtualService manifest: %w", err)
	}

	// Resolve kubectl version: multiTarget > spec > deployTarget
	kubectlVersion := cmp.Or(appCfg.Input.KubectlVersion, dt.Config.KubectlVersion)
	if multiTarget != nil {
		kubectlVersion = cmp.Or(multiTarget.KubectlVersion, kubectlVersion)
	}

	kubectlPath, err := toolRegistry.Kubectl(ctx, kubectlVersion)
	if err != nil {
		return fmt.Errorf("failed while getting kubectl tool: %w", err)
	}

	kubectl := provider.NewKubectl(kubectlPath)
	applier := provider.NewApplier(kubectl, appCfg.Input, dt.Config, input.Logger)

	if err := applyManifests(ctx, applier, []provider.Manifest{vs}, appCfg.Input.Namespace, lp); err != nil {
		return fmt.Errorf("failed while applying VirtualService manifest: %w", err)
	}

	return nil
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
		return fmt.Errorf("require %s but got %s for %s key in spec.selector", variant, value, variantLabel)
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
		if !strings.Contains(m.APIVersion(), istioNetworkingGroup) {
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

// virtualService is a wrapper around istiov1.VirtualService to enable conversion
// from unstructured.Unstructured. Works across APIVersion v1alpha3 and v1beta1
// because v1.VirtualService is a type alias of v1alpha3.VirtualService.
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

// generateVirtualServiceManifest updates a VirtualService manifest with the given
// traffic percentages for canary and baseline variants. The primary variant receives
// the remainder. Routes are filtered by editableRoutes if non-empty.
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

		// Calculate the weight of routes pointing to other hosts (not this service).
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

		// Calculate variant weights, preserving any weight reserved for other hosts.
		var (
			variantsWeight = 100 - otherHostWeight
			canaryWeight   = canaryPercent * variantsWeight / 100
			baselineWeight = baselinePercent * variantsWeight / 100
			primaryWeight  = variantsWeight - canaryWeight - baselineWeight
			routes         = make([]*istiov1.HTTPRouteDestination, 0, len(otherHostRoutes)+3)
		)

		routes = append(routes, &istiov1.HTTPRouteDestination{
			Destination: &istiov1.Destination{
				Host:   host,
				Subset: variantLabel.PrimaryValue,
			},
			Weight: primaryWeight,
		})

		if canaryWeight > 0 {
			routes = append(routes, &istiov1.HTTPRouteDestination{
				Destination: &istiov1.Destination{
					Host:   host,
					Subset: variantLabel.CanaryValue,
				},
				Weight: canaryWeight,
			})
		}

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
