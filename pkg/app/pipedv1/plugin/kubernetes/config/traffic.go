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

package config

import (
	"fmt"

	"github.com/pipe-cd/piped-plugin-sdk-go/unit"
)

type KubernetesTrafficRoutingMethod string

const (
	// KubernetesTrafficRoutingMethodPodSelector is the way by updating the selector in Service to switching all of traffic.
	KubernetesTrafficRoutingMethodPodSelector KubernetesTrafficRoutingMethod = "podselector"
	// KubernetesTrafficRoutingMethodIstio is the way by updating the VirtualService to update traffic routing.
	KubernetesTrafficRoutingMethodIstio KubernetesTrafficRoutingMethod = "istio"
)

// KubernetesTrafficRouting represents the traffic routing configuration for a Kubernetes application.
type KubernetesTrafficRouting struct {
	// The method to be used for traffic routing.
	// The default is PodSelector: the way by updating the selector in Service to switching all of traffic.
	Method KubernetesTrafficRoutingMethod `json:"method"`
	// The Istio-specific configuration for traffic routing.
	Istio *IstioTrafficRouting `json:"istio"`
}

// DetermineKubernetesTrafficRoutingMethod determines the routing method should be used based on the TrafficRouting config.
// The default is PodSelector: the way by updating the selector in Service to switching all of traffic.
func DetermineKubernetesTrafficRoutingMethod(cfg *KubernetesTrafficRouting) KubernetesTrafficRoutingMethod {
	if cfg == nil || cfg.Method == "" {
		return KubernetesTrafficRoutingMethodPodSelector
	}
	return cfg.Method
}

// IstioTrafficRouting represents the Istio-specific configuration for traffic routing.
type IstioTrafficRouting struct {
	// List of routes in the VirtualService that can be changed to update traffic routing.
	// Empty means all routes should be updated.
	EditableRoutes []string `json:"editableRoutes"`
	// TODO: Add a validate to ensure this was configured or using the default value by service name.
	// The service host.
	Host string `json:"host"`
	// The reference to VirtualService manifest.
	// Empty means the first VirtualService resource will be used.
	VirtualService K8sResourceReference `json:"virtualService"`
}

// K8sTrafficRoutingStageOptions contains all configurable values for a K8S_TRAFFIC_ROUTING stage.
type K8sTrafficRoutingStageOptions struct {
	// Which variant should receive all traffic.
	// "primary" or "canary" or "baseline" can be populated.
	All string `json:"all"`
	// The percentage of traffic should be routed to PRIMARY variant.
	Primary unit.Percentage `json:"primary"`
	// The percentage of traffic should be routed to CANARY variant.
	Canary unit.Percentage `json:"canary"`
	// The percentage of traffic should be routed to BASELINE variant.
	Baseline unit.Percentage `json:"baseline"`
}

// Percentages returns the primary, canary, and baseline percentages from the K8sTrafficRoutingStageOptions.
func (opts K8sTrafficRoutingStageOptions) Percentages() (primary, canary, baseline int) {
	switch opts.All {
	case "primary":
		return 100, 0, 0
	case "canary":
		return 0, 100, 0
	case "baseline":
		return 0, 0, 100
	}
	return opts.Primary.Int(), opts.Canary.Int(), opts.Baseline.Int()
}

// DisplayString returns the display string for the K8sTrafficRoutingStageOptions.
// This is used to display the traffic routing configuration in the UI.
func (opts K8sTrafficRoutingStageOptions) DisplayString() string {
	primary, canary, baseline := opts.Percentages()
	return fmt.Sprintf("Primary: %d%%, Canary: %d%%, Baseline: %d%%", primary, canary, baseline)
}
