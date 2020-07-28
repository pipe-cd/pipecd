// Copyright 2020 The PipeCD Authors.
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

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

const (
	primaryMetadataKey  = "primary-percentage"
	canaryMetadataKey   = "canary-percentage"
	baselineMetadataKey = "baseline-percentage"
)

func (e *Executor) ensureTrafficRouting(ctx context.Context) model.StageStatus {
	var (
		commitHash = e.Deployment.Trigger.Commit.Hash
		options    = e.StageConfig.K8sTrafficRoutingStageOptions
	)
	if options == nil {
		e.LogPersister.AppendErrorf("Malformed configuration for stage %s", e.Stage.Name)
		return model.StageStatus_STAGE_FAILURE
	}

	// Load the manifests at the triggered commit.
	e.LogPersister.AppendInfof("Loading manifests at commit %s for handling", commitHash)
	manifests, err := e.loadManifests(ctx)
	if err != nil {
		e.LogPersister.AppendErrorf("Failed while loading manifests (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.AppendSuccessf("Successfully loaded %d manifests", len(manifests))

	if len(manifests) == 0 {
		e.LogPersister.AppendError("There are no kubernetes manifests to handle")
		return model.StageStatus_STAGE_FAILURE
	}

	// Decide traffic routing percentage for all variants.
	primaryPercent, canaryPercent, baselinePercent := options.Percentages()
	e.saveTrafficRoutingMetadata(ctx, primaryPercent, canaryPercent, baselinePercent)

	// Find traffic routing manifests.
	trafficRoutingManifests, err := e.findTrafficRoutingManifests(manifests, e.config.TrafficRouting)
	if err != nil {
		e.LogPersister.AppendErrorf("Failed while finding traffic routing manifest: (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	switch len(trafficRoutingManifests) {
	case 1:
		break
	case 0:
		e.LogPersister.AppendErrorf("Unable to find any traffic routing manifests")
		return model.StageStatus_STAGE_FAILURE
	default:
		e.LogPersister.AppendInfof(
			"Detected %d traffic routing manifests but only the first one (%s) will be used",
			len(trafficRoutingManifests),
			trafficRoutingManifests[0].Key.ReadableString(),
		)
	}
	trafficRoutingManifest := trafficRoutingManifests[0]

	// In case we are routing by Pod, the service manifest must contain variantLabel inside its selector.
	if e.config.TrafficRouting == nil || e.config.TrafficRouting.Method == config.TrafficRoutingMethodPod || e.config.TrafficRouting.Method == "" {
		if err := checkVariantSelectorInService(trafficRoutingManifest, primaryVariant); err != nil {
			e.LogPersister.AppendErrorf("Traffic routing by Pod requires %q inside the selector of Service manifest but it was unable to check that field in manifest %s (%v)",
				variantLabel+": "+primaryVariant,
				trafficRoutingManifest.Key.ReadableString(),
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
		e.config.TrafficRouting,
	)
	if err != nil {
		e.LogPersister.AppendErrorf("Unable generate traffic routing manifest: (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	// Add builtin annotations for tracking application live state.
	e.addBuiltinAnnontations([]provider.Manifest{trafficRoutingManifest}, primaryVariant, commitHash)

	e.LogPersister.AppendInfof("Start updating traffic routing to be percentages: primary=%d, canary=%d, baseline=%d",
		primaryPercent,
		canaryPercent,
		baselinePercent,
	)
	if err := e.applyManifests(ctx, []provider.Manifest{trafficRoutingManifest}); err != nil {
		return model.StageStatus_STAGE_FAILURE
	}

	e.LogPersister.AppendSuccess("Successfully updated traffic routing")
	return model.StageStatus_STAGE_SUCCESS
}

func (e *Executor) findTrafficRoutingManifests(manifests []provider.Manifest, cfg *config.TrafficRouting) ([]provider.Manifest, error) {
	if cfg != nil && cfg.Method == config.TrafficRoutingMethodIstio {
		istioConfig := cfg.Istio
		if istioConfig == nil {
			istioConfig = &config.IstioTrafficRouting{}
		}
		return findIstioVirtualServiceManifests(manifests, istioConfig.VirtualService)
	}

	var podConfig config.PodTrafficRouting
	if cfg != nil && cfg.Pod != nil {
		podConfig = *cfg.Pod
	}

	// Find out the service which be updated the selector.
	_, serviceName, ok := config.ParseVariantResourceReference(podConfig.Service.Reference)
	if !ok {
		return nil, fmt.Errorf("malformed Service reference %q", podConfig.Service.Reference)
	}

	return findManifests(provider.KindService, serviceName, manifests), nil
}

func (e *Executor) generateTrafficRoutingManifest(manifest provider.Manifest, primaryPercent, canaryPercent, baselinePercent int, cfg *config.TrafficRouting) (provider.Manifest, error) {
	if cfg != nil && cfg.Method == config.TrafficRoutingMethodIstio {
		istioConfig := cfg.Istio
		if istioConfig == nil {
			istioConfig = &config.IstioTrafficRouting{}
		}

		if strings.HasPrefix(manifest.Key.APIVersion, "v1alpha3") {
			return generateVirtualServiceManifestV1Alpha3(manifest, istioConfig.Host, istioConfig.EditableRoutes, int32(canaryPercent), int32(baselinePercent))
		}
		return generateVirtualServiceManifest(manifest, istioConfig.Host, istioConfig.EditableRoutes, int32(canaryPercent), int32(baselinePercent))
	}

	// Because the loaded maninests are read-only
	// so we duplicate them to avoid updating the shared manifests data in cache.
	manifest = duplicateManifest(manifest, "")

	// Determine which variant will receive 100% percent of traffic.
	var variant string
	switch {
	case primaryPercent == 100:
		variant = primaryVariant
	case canaryPercent == 100:
		variant = canaryVariant
	default:
		return manifest, fmt.Errorf("traffic routing by pod requires either PRIMARY or CANARY must be 100 (primary=%d, canary=%d)", primaryPercent, canaryPercent)
	}

	if err := manifest.AddStringMapValues(map[string]string{variantLabel: variant}, "spec", "selector"); err != nil {
		return manifest, fmt.Errorf("unable to update selector for service %q because of: %v", manifest.Key.Name, err)
	}

	return manifest, nil
}

func (e *Executor) saveTrafficRoutingMetadata(ctx context.Context, primary, canary, baseline int) {
	metadata := map[string]string{
		primaryMetadataKey:  strconv.FormatInt(int64(primary), 10),
		canaryMetadataKey:   strconv.FormatInt(int64(canary), 10),
		baselineMetadataKey: strconv.FormatInt(int64(baseline), 10),
	}
	if err := e.MetadataStore.SetStageMetadata(ctx, e.Stage.Id, metadata); err != nil {
		e.Logger.Error("failed to save traffic routing percentages to metadata", zap.Error(err))
	}
}

func findIstioVirtualServiceManifests(manifests []provider.Manifest, cfg config.K8sResourceReference) ([]provider.Manifest, error) {
	const (
		istioNetworkingAPIVersionPrefix = "networking.istio.io/"
		istioVirtualServiceKind         = "VirtualService"
	)
	_, name, ok := config.ParseVariantResourceReference(cfg.Reference)
	if !ok {
		return nil, fmt.Errorf("malformed VirtualService reference: %s", cfg.Reference)
	}

	var out []provider.Manifest
	for _, m := range manifests {
		if !strings.HasPrefix(m.Key.APIVersion, istioNetworkingAPIVersionPrefix) {
			continue
		}
		if m.Key.Kind != istioVirtualServiceKind {
			continue
		}
		if name != "" && m.Key.Name != name {
			continue
		}
		out = append(out, m)
	}

	return out, nil
}

func generateVirtualServiceManifest(m provider.Manifest, host string, editableRoutes []string, canaryPercent, baselinePercent int32) (provider.Manifest, error) {
	// Because the loaded maninests are read-only
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
		routes = append(routes,
			&istiov1beta1.HTTPRouteDestination{
				Destination: &istiov1beta1.Destination{
					Host:   host,
					Subset: primaryVariant,
				},
				Weight: primaryWeight,
			},
			&istiov1beta1.HTTPRouteDestination{
				Destination: &istiov1beta1.Destination{
					Host:   host,
					Subset: canaryVariant,
				},
				Weight: canaryWeight,
			},
			&istiov1beta1.HTTPRouteDestination{
				Destination: &istiov1beta1.Destination{
					Host:   host,
					Subset: baselineVariant,
				},
				Weight: baselineWeight,
			},
		)
		routes = append(routes, otherHostRoutes...)
		http.Route = routes
	}

	if err := m.SetStructuredSpec(vs); err != nil {
		return m, err
	}

	return m, nil
}

func generateVirtualServiceManifestV1Alpha3(m provider.Manifest, host string, editableRoutes []string, canaryPercent, baselinePercent int32) (provider.Manifest, error) {
	// Because the loaded maninests are read-only
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
		routes = append(routes,
			&istiov1alpha3.HTTPRouteDestination{
				Destination: &istiov1alpha3.Destination{
					Host:   host,
					Subset: primaryVariant,
				},
				Weight: primaryWeight,
			},
			&istiov1alpha3.HTTPRouteDestination{
				Destination: &istiov1alpha3.Destination{
					Host:   host,
					Subset: canaryVariant,
				},
				Weight: canaryWeight,
			},
			&istiov1alpha3.HTTPRouteDestination{
				Destination: &istiov1alpha3.Destination{
					Host:   host,
					Subset: baselineVariant,
				},
				Weight: baselineWeight,
			},
		)
		routes = append(routes, otherHostRoutes...)
		http.Route = routes
	}

	if err := m.SetStructuredSpec(vs); err != nil {
		return m, err
	}

	return m, nil
}

func checkVariantSelectorInService(m provider.Manifest, variant string) error {
	selector, err := m.GetNestedStringMap("spec", "selector")
	if err != nil {
		return err
	}

	value, ok := selector[variantLabel]
	if !ok {
		return fmt.Errorf("missing %s key in the selector", variantLabel)
	}

	if value != variant {
		return fmt.Errorf("want %s but got %s for %s key in the selector", variant, value, variantLabel)
	}
	return nil
}
