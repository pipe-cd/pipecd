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
	"strings"

	istiov1alpha3 "istio.io/api/networking/v1alpha3"
	istiov1beta1 "istio.io/api/networking/v1beta1"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

func (e *Executor) ensureTrafficRouting(ctx context.Context) model.StageStatus {
	options := e.StageConfig.K8sTrafficRoutingStageOptions
	if options == nil {
		e.LogPersister.AppendError(fmt.Sprintf("Malformed configuration for stage %s", e.Stage.Name))
		return model.StageStatus_STAGE_FAILURE
	}

	manifests, err := e.loadManifests(ctx)
	if err != nil {
		e.LogPersister.AppendError(fmt.Sprintf("Failed while loading manifests (%v)", err))
		return model.StageStatus_STAGE_FAILURE
	}

	if len(manifests) == 0 {
		e.LogPersister.AppendError("There are no kubernetes manifests to handle")
		return model.StageStatus_STAGE_FAILURE
	}

	var (
		cfg                                            = e.config.TrafficRouting
		primaryPercent, canaryPercent, baselinePercent = options.Percentages()
	)
	if cfg == nil {
		return e.ensurePodTrafficRouting(ctx, primaryPercent, canaryPercent, manifests, nil)
	}

	switch cfg.Method {
	case config.TrafficRoutingMethodIstio:
		return e.ensureIstioTrafficRouting(ctx, canaryPercent, baselinePercent, manifests, cfg.Istio)

	default:
		return e.ensurePodTrafficRouting(ctx, primaryPercent, canaryPercent, manifests, cfg.Pod)
	}
}

func (e *Executor) rollbackTraffic(ctx context.Context) error {
	return nil
}

func (e *Executor) ensurePodTrafficRouting(ctx context.Context, primaryPercent, canaryPercent int, manifests []provider.Manifest, cfg *config.PodTrafficRouting) model.StageStatus {
	variant := primaryVariant
	if canaryPercent > primaryPercent {
		variant = canaryVariant
	}
	e.LogPersister.AppendInfo(fmt.Sprintf("All traffic should be routed to %s", variant))
	e.LogPersister.AppendError("Not implemented")
	return model.StageStatus_STAGE_FAILURE
}

func (e *Executor) ensureIstioTrafficRouting(ctx context.Context, canaryPercent, baselinePercent int, manifests []provider.Manifest, cfg *config.IstioTrafficRouting) model.StageStatus {
	if cfg == nil {
		cfg = &config.IstioTrafficRouting{}
	}

	manifest, ok := findIstioVirtualServiceManifest(manifests, cfg.VirtualService)
	if !ok {
		e.LogPersister.AppendError("There is no VirtualService manifests to handle")
		return model.StageStatus_STAGE_FAILURE
	}

	var err error
	if strings.HasPrefix(manifest.Key.APIVersion, "v1alpha3") {
		err = generateVirtualServiceManifestV1Alpha3(manifest, cfg.Host, cfg.EditableRoutes, int32(canaryPercent), int32(baselinePercent))
	} else {
		err = generateVirtualServiceManifest(manifest, cfg.Host, cfg.EditableRoutes, int32(canaryPercent), int32(baselinePercent))
	}
	if err != nil {
		e.LogPersister.AppendError(fmt.Sprintf("Unable to generate VirtualService manifest: %v", err))
		return model.StageStatus_STAGE_FAILURE
	}

	e.LogPersister.AppendInfo("Start updating traffic routing...")
	if err = e.provider.ApplyManifest(ctx, manifest); err != nil {
		e.LogPersister.AppendError(fmt.Sprintf("Failed to apply manifest: %s (%v)", manifest.Key.ReadableString(), err))
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.AppendSuccess(fmt.Sprintf("Applied manifest: %s", manifest.Key.ReadableString()))

	e.LogPersister.AppendSuccess("Successfully updated traffic routing")
	return model.StageStatus_STAGE_SUCCESS
}

func findIstioVirtualServiceManifest(manifests []provider.Manifest, cfg config.K8sVariantVirtualService) (provider.Manifest, bool) {
	const (
		istioNetworkingAPIVersionPrefix = "networking.istio.io/"
		istioVirtualServiceKind         = "VirtualService"
	)
	_, name, ok := config.ParseVariantResourceReference(cfg.Reference)
	if !ok {
		return provider.Manifest{}, false
	}

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
		return m, true
	}

	return provider.Manifest{}, false
}

func generateVirtualServiceManifest(m provider.Manifest, host string, editableRoutes []string, canaryPercent, baselinePercent int32) error {
	spec, err := m.GetSpec()
	if err != nil {
		return err
	}

	vs := istiov1beta1.VirtualService{}
	data, err := json.Marshal(spec)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &vs); err != nil {
		return err
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

	return m.SetStructuredSpec(vs)
}

func generateVirtualServiceManifestV1Alpha3(m provider.Manifest, host string, editableRoutes []string, canaryPercent, baselinePercent int32) error {
	spec, err := m.GetSpec()
	if err != nil {
		return err
	}

	vs := istiov1alpha3.VirtualService{}
	data, err := json.Marshal(spec)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &vs); err != nil {
		return err
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

	return m.SetStructuredSpec(vs)
}
