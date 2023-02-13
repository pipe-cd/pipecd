// Copyright 2023 The PipeCD Authors.
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

package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"
	istiov1alpha3 "istio.io/api/networking/v1alpha3"
	istiov1beta1 "istio.io/api/networking/v1beta1"

	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	primaryMetadataKey  = "primary-percentage"
	canaryMetadataKey   = "canary-percentage"
	baselineMetadataKey = "baseline-percentage"
)

func (e *deployExecutor) ensureTrafficRouting(ctx context.Context) model.StageStatus {
	var (
		commitHash     = e.Deployment.Trigger.Commit.Hash
		options        = e.StageConfig.K8sTrafficRoutingStageOptions
		variantLabel   = e.appCfg.VariantLabel.Key
		primaryVariant = e.appCfg.VariantLabel.PrimaryValue
	)
	if options == nil {
		e.LogPersister.Errorf("Malformed configuration for stage %s", e.Stage.Name)
		return model.StageStatus_STAGE_FAILURE
	}
	method := config.DetermineKubernetesTrafficRoutingMethod(e.appCfg.TrafficRouting)

	// Load the manifests at the triggered commit.
	e.LogPersister.Infof("Loading manifests at commit %s for handling", commitHash)
	manifests, err := loadManifests(
		ctx,
		e.Deployment.ApplicationId,
		e.commit,
		e.AppManifestsCache,
		e.loader,
		e.Logger,
	)
	if err != nil {
		e.LogPersister.Errorf("Failed while loading manifests (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.Successf("Successfully loaded %d manifests", len(manifests))

	if len(manifests) == 0 {
		e.LogPersister.Error("There are no kubernetes manifests to handle")
		return model.StageStatus_STAGE_FAILURE
	}

	// Decide traffic routing percentage for all variants.
	primaryPercent, canaryPercent, baselinePercent := options.Percentages()
	e.saveTrafficRoutingMetadata(ctx, primaryPercent, canaryPercent, baselinePercent)

	// Find traffic routing manifests.
	trafficRoutingManifests, err := findTrafficRoutingManifests(manifests, e.appCfg.Service.Name, e.appCfg.TrafficRouting)
	if err != nil {
		e.LogPersister.Errorf("Failed while finding traffic routing manifest: (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	switch len(trafficRoutingManifests) {
	case 1:
		break
	case 0:
		e.LogPersister.Errorf("Unable to find any traffic routing manifests")
		return model.StageStatus_STAGE_FAILURE
	default:
		e.LogPersister.Infof(
			"Detected %d traffic routing manifests but only the first one (%s) will be used",
			len(trafficRoutingManifests),
			trafficRoutingManifests[0].Key.ReadableLogString(),
		)
	}
	trafficRoutingManifest := trafficRoutingManifests[0]

	// In case we are routing by PodSelector, the service manifest must contain variantLabel inside its selector.
	if method == config.KubernetesTrafficRoutingMethodPodSelector {
		if err := checkVariantSelectorInService(trafficRoutingManifest, variantLabel, primaryVariant); err != nil {
			e.LogPersister.Errorf("Traffic routing by PodSelector requires %q inside the selector of Service manifest but it was unable to check that field in manifest %s (%v)",
				variantLabel+": "+primaryVariant,
				trafficRoutingManifest.Key.ReadableLogString(),
				err,
			)
			return model.StageStatus_STAGE_FAILURE
		}
	}

	trafficRoutingManifest, err = e.generateTrafficRoutingManifest(
		trafficRoutingManifest,
		primaryPercent,
		canaryPercent,
		baselinePercent,
	)
	if err != nil {
		e.LogPersister.Errorf("Unable generate traffic routing manifest: (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	// Add builtin annotations for tracking application live state.
	addBuiltinAnnotations(
		[]provider.Manifest{trafficRoutingManifest},
		variantLabel,
		primaryVariant,
		commitHash,
		e.PipedConfig.PipedID,
		e.Deployment.ApplicationId,
	)

	e.LogPersister.Infof("Start updating traffic routing to be percentages: primary=%d, canary=%d, baseline=%d",
		primaryPercent,
		canaryPercent,
		baselinePercent,
	)
	if err := applyManifests(ctx, e.applierGetter, []provider.Manifest{trafficRoutingManifest}, e.appCfg.Input.Namespace, e.LogPersister); err != nil {
		return model.StageStatus_STAGE_FAILURE
	}

	e.LogPersister.Success("Successfully updated traffic routing")
	return model.StageStatus_STAGE_SUCCESS
}

func findTrafficRoutingManifests(manifests []provider.Manifest, serviceName string, cfg *config.KubernetesTrafficRouting) ([]provider.Manifest, error) {
	method := config.DetermineKubernetesTrafficRoutingMethod(cfg)

	switch method {
	case config.KubernetesTrafficRoutingMethodPodSelector:
		return findManifests(provider.KindService, serviceName, manifests), nil

	case config.KubernetesTrafficRoutingMethodIstio:
		istioConfig := cfg.Istio
		if istioConfig == nil {
			istioConfig = &config.IstioTrafficRouting{}
		}
		return findIstioVirtualServiceManifests(manifests, istioConfig.VirtualService)

	default:
		return nil, fmt.Errorf("unsupport traffic routing method %v", method)
	}
}

func (e *deployExecutor) generateTrafficRoutingManifest(manifest provider.Manifest, primaryPercent, canaryPercent, baselinePercent int) (provider.Manifest, error) {
	// Because the loaded manifests are read-only
	// so we duplicate them to avoid updating the shared manifests data in cache.
	manifest = duplicateManifest(manifest, "")

	// When all traffic should be routed to primary variant
	// we do not need to change the traffic manifest
	// just copy and return the one specified in the target commit.
	if primaryPercent == 100 {
		return manifest, nil
	}

	cfg := e.appCfg.TrafficRouting
	if cfg != nil && cfg.Method == config.KubernetesTrafficRoutingMethodIstio {
		istioConfig := cfg.Istio
		if istioConfig == nil {
			istioConfig = &config.IstioTrafficRouting{}
		}

		if strings.HasPrefix(manifest.Key.APIVersion, "v1alpha3") {
			return e.generateVirtualServiceManifestV1Alpha3(manifest, istioConfig.Host, istioConfig.EditableRoutes, int32(canaryPercent), int32(baselinePercent))
		}
		return e.generateVirtualServiceManifest(manifest, istioConfig.Host, istioConfig.EditableRoutes, int32(canaryPercent), int32(baselinePercent))
	}

	// Determine which variant will receive 100% percent of traffic.
	var variant string
	switch {
	case primaryPercent == 100:
		variant = e.appCfg.VariantLabel.PrimaryValue
	case canaryPercent == 100:
		variant = e.appCfg.VariantLabel.CanaryValue
	default:
		return manifest, fmt.Errorf("traffic routing by pod requires either PRIMARY or CANARY must be 100 (primary=%d, canary=%d)", primaryPercent, canaryPercent)
	}

	variantLabel := e.appCfg.VariantLabel.Key
	if err := manifest.AddStringMapValues(map[string]string{variantLabel: variant}, "spec", "selector"); err != nil {
		return manifest, fmt.Errorf("unable to update selector for service %q because of: %v", manifest.Key.Name, err)
	}

	return manifest, nil
}

func (e *deployExecutor) saveTrafficRoutingMetadata(ctx context.Context, primary, canary, baseline int) {
	metadata := map[string]string{
		primaryMetadataKey:  strconv.FormatInt(int64(primary), 10),
		canaryMetadataKey:   strconv.FormatInt(int64(canary), 10),
		baselineMetadataKey: strconv.FormatInt(int64(baseline), 10),
	}
	if err := e.MetadataStore.Stage(e.Stage.Id).PutMulti(ctx, metadata); err != nil {
		e.Logger.Error("failed to save traffic routing percentages to metadata", zap.Error(err))
	}
}

func findIstioVirtualServiceManifests(manifests []provider.Manifest, ref config.K8sResourceReference) ([]provider.Manifest, error) {
	const (
		istioNetworkingAPIVersionPrefix = "networking.istio.io/"
		istioVirtualServiceKind         = "VirtualService"
	)

	if ref.Kind != "" && ref.Kind != istioVirtualServiceKind {
		return nil, fmt.Errorf("support only %q kind for VirtualService reference", istioVirtualServiceKind)
	}

	var out []provider.Manifest
	for _, m := range manifests {
		if !strings.HasPrefix(m.Key.APIVersion, istioNetworkingAPIVersionPrefix) {
			continue
		}
		if m.Key.Kind != istioVirtualServiceKind {
			continue
		}
		if ref.Name != "" && m.Key.Name != ref.Name {
			continue
		}
		out = append(out, m)
	}

	return out, nil
}

func (e *deployExecutor) generateVirtualServiceManifest(m provider.Manifest, host string, editableRoutes []string, canaryPercent, baselinePercent int32) (provider.Manifest, error) {
	// Because the loaded manifests are read-only
	// so we duplicate them to avoid updating the shared manifests data in cache.
	m = duplicateManifest(m, "")

	spec, err := m.GetSpec()
	if err != nil {
		return m, err
	}

	vs := istiov1beta1.VirtualService{}
	data, err := json.Marshal(spec)
	if err != nil {
		return m, err
	}
	if err := json.Unmarshal(data, &vs); err != nil {
		return m, err
	}

	editableMap := make(map[string]struct{}, len(editableRoutes))
	for _, r := range editableRoutes {
		editableMap[r] = struct{}{}
	}

	for _, http := range vs.Http {
		if len(editableMap) > 0 {
			if _, ok := editableMap[http.Name]; !ok {
				continue
			}
		}

		var (
			otherHostWeight int32
			otherHostRoutes = make([]*istiov1beta1.HTTPRouteDestination, 0)
		)
		for _, r := range http.Route {
			if r.Destination != nil && r.Destination.Host != host {
				otherHostWeight += r.Weight
				otherHostRoutes = append(otherHostRoutes, r)
			}
		}

		var (
			variantsWeight = 100 - otherHostWeight
			canaryWeight   = canaryPercent * variantsWeight / 100
			baselineWeight = baselinePercent * variantsWeight / 100
			primaryWeight  = variantsWeight - canaryWeight - baselineWeight
			routes         = make([]*istiov1beta1.HTTPRouteDestination, 0, len(otherHostRoutes)+3)
		)

		routes = append(routes, &istiov1beta1.HTTPRouteDestination{
			Destination: &istiov1beta1.Destination{
				Host:   host,
				Subset: e.appCfg.VariantLabel.PrimaryValue,
			},
			Weight: primaryWeight,
		})
		if canaryWeight > 0 {
			routes = append(routes, &istiov1beta1.HTTPRouteDestination{
				Destination: &istiov1beta1.Destination{
					Host:   host,
					Subset: e.appCfg.VariantLabel.CanaryValue,
				},
				Weight: canaryWeight,
			})
		}
		if baselineWeight > 0 {
			routes = append(routes, &istiov1beta1.HTTPRouteDestination{
				Destination: &istiov1beta1.Destination{
					Host:   host,
					Subset: e.appCfg.VariantLabel.BaselineValue,
				},
				Weight: baselineWeight,
			})
		}
		routes = append(routes, otherHostRoutes...)
		http.Route = routes
	}

	if err := m.SetStructuredSpec(vs); err != nil {
		return m, err
	}

	return m, nil
}

func (e *deployExecutor) generateVirtualServiceManifestV1Alpha3(m provider.Manifest, host string, editableRoutes []string, canaryPercent, baselinePercent int32) (provider.Manifest, error) {
	// Because the loaded manifests are read-only
	// so we duplicate them to avoid updating the shared manifests data in cache.
	m = duplicateManifest(m, "")

	spec, err := m.GetSpec()
	if err != nil {
		return m, err
	}

	vs := istiov1alpha3.VirtualService{}
	data, err := json.Marshal(spec)
	if err != nil {
		return m, err
	}
	if err := json.Unmarshal(data, &vs); err != nil {
		return m, err
	}

	editableMap := make(map[string]struct{}, len(editableRoutes))
	for _, r := range editableRoutes {
		editableMap[r] = struct{}{}
	}

	for _, http := range vs.Http {
		if len(editableMap) > 0 {
			if _, ok := editableMap[http.Name]; !ok {
				continue
			}
		}

		var (
			otherHostWeight int32
			otherHostRoutes = make([]*istiov1alpha3.HTTPRouteDestination, 0)
		)
		for _, r := range http.Route {
			if r.Destination != nil && r.Destination.Host != host {
				otherHostWeight += r.Weight
				otherHostRoutes = append(otherHostRoutes, r)
			}
		}

		var (
			variantsWeight = 100 - otherHostWeight
			canaryWeight   = canaryPercent * variantsWeight / 100
			baselineWeight = baselinePercent * variantsWeight / 100
			primaryWeight  = variantsWeight - canaryWeight - baselineWeight
			routes         = make([]*istiov1alpha3.HTTPRouteDestination, 0, len(otherHostRoutes)+3)
		)

		routes = append(routes, &istiov1alpha3.HTTPRouteDestination{
			Destination: &istiov1alpha3.Destination{
				Host:   host,
				Subset: e.appCfg.VariantLabel.PrimaryValue,
			},
			Weight: primaryWeight,
		})
		if canaryWeight > 0 {
			routes = append(routes, &istiov1alpha3.HTTPRouteDestination{
				Destination: &istiov1alpha3.Destination{
					Host:   host,
					Subset: e.appCfg.VariantLabel.CanaryValue,
				},
				Weight: canaryWeight,
			})
		}
		if baselineWeight > 0 {
			routes = append(routes, &istiov1alpha3.HTTPRouteDestination{
				Destination: &istiov1alpha3.Destination{
					Host:   host,
					Subset: e.appCfg.VariantLabel.BaselineValue,
				},
				Weight: baselineWeight,
			})
		}
		routes = append(routes, otherHostRoutes...)
		http.Route = routes
	}

	if err := m.SetStructuredSpec(vs); err != nil {
		return m, err
	}

	return m, nil
}

func checkVariantSelectorInService(m provider.Manifest, variantLabel, variant string) error {
	selector, err := m.GetNestedStringMap("spec", "selector")
	if err != nil {
		return err
	}

	value, ok := selector[variantLabel]
	if !ok {
		return fmt.Errorf("missing %s key in spec.selector", variantLabel)
	}

	if value != variant {
		return fmt.Errorf("require %s but got %s for %s key in spec.selector", variant, value, variantLabel)
	}
	return nil
}
