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

package config

// KubernetesDeploymentSpec represents a deployment configuration for Kubernetes application.
type KubernetesDeploymentSpec struct {
	Input         KubernetesDeploymentInput `json:"input"`
	CommitMatcher DeploymentCommitMatcher   `json:"commitMatcher"`
	QuickSync     K8sSyncStageOptions       `json:"quickSync"`
	Pipeline      *DeploymentPipeline       `json:"pipeline"`

	// Which resource should be considered as the Service of application.
	// Empty means the first Service resource will be used.
	Service K8sResourceReference `json:"service"`
	// Which resources should be considered as the Workload of application.
	// Empty means all Deployments.
	// e.g.
	// - kind: Deployment
	//   name: deployment-name
	// - kind: ReplicationController
	//   name: replication-controller-name
	Workloads []K8sResourceReference `json:"workloads"`
	// Which method should be used for traffic routing.
	TrafficRouting *TrafficRouting `json:"trafficRouting"`
}

func (s *KubernetesDeploymentSpec) GetStage(index int32) (PipelineStage, bool) {
	if s.Pipeline == nil {
		return PipelineStage{}, false
	}
	if int(index) >= len(s.Pipeline.Stages) {
		return PipelineStage{}, false
	}
	return s.Pipeline.Stages[index], true
}

// Validate returns an error if any wrong configuration value was found.
func (s *KubernetesDeploymentSpec) Validate() error {
	return nil
}

type KubernetesDeploymentInput struct {
	Manifests      []string `json:"manifests"`
	KubectlVersion string   `json:"kubectlVersion"`

	KustomizeVersion string `json:"kustomizeVersion"`

	HelmChart   *InputHelmChart   `json:"helmChart"`
	HelmOptions *InputHelmOptions `json:"helmOptions"`
	HelmVersion string            `json:"helmVersion"`

	// The namespace where manifests will be applied.
	Namespace string `json:"namespace"`
	// Automatically reverts all changes from all stages when one of them failed.
	// Default is true.
	AutoRollback bool     `json:"autoRollback"`
	Dependencies []string `json:"dependencies,omitempty"`
}

type InputHelmChart struct {
	// Empty means current repository.
	GitRemote string `json:"gitRemote"`
	// The commit SHA or tag for remote git.
	Ref string `json:"ref"`
	// Relative path from the repository root directory to the chart directory.
	Path string `json:"path"`

	// The name of an added Helm chart repository.
	Repository string `json:"repository"`
	Name       string `json:"name"`
	Version    string `json:"version"`
}

type InputHelmOptions struct {
	// By default the release name is equal to the application name.
	ReleaseName string `json:"releaseName"`
	// List of value files should be loaded.
	ValueFiles []string `json:"valueFiles"`
	SetFiles   map[string]string
}

type TrafficRoutingMethod string

const (
	TrafficRoutingMethodPodSelector TrafficRoutingMethod = "podselector"
	TrafficRoutingMethodIstio       TrafficRoutingMethod = "istio"
)

type TrafficRouting struct {
	Method TrafficRoutingMethod `json:"method"`
	Istio  *IstioTrafficRouting `json:"istio"`
}

// DetermineTrafficRoutingMethod determines the routing method should be used based on the TrafficRouting config.
// The default is PodSelector: the way by updating the selector in Service to switching all of traffic.
func DetermineTrafficRoutingMethod(cfg *TrafficRouting) TrafficRoutingMethod {
	if cfg == nil {
		return TrafficRoutingMethodPodSelector
	}
	if cfg.Method == "" {
		return TrafficRoutingMethodPodSelector
	}
	return cfg.Method
}

type IstioTrafficRouting struct {
	EditableRoutes []string `json:"editableRoutes"`
	// TODO: Add a validate to ensure this was configured or using the default value by service name.
	Host           string               `json:"host"`
	VirtualService K8sResourceReference `json:"virtualService"`
}

type K8sResourceReference struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

// K8sSyncStageOptions contains all configurable values for a K8S_SYNC stage.
type K8sSyncStageOptions struct {
	// Whether the PRIMARY variant label should be added to manifests if they were missing.
	AddVariantLabelToSelector bool `json:"addVariantLabelToSelector"`
	// Whether the resources that are no longer defined in Git will be removed.
	Prune bool `json:"prune"`
}

// K8sPrimaryRolloutStageOptions contains all configurable values for a K8S_PRIMARY_ROLLOUT stage.
type K8sPrimaryRolloutStageOptions struct {
	// Suffix that should be used when naming the PRIMARY variant's resources.
	// Default is "primary".
	Suffix string `json:"suffix"`
	// Whether the PRIMARY service should be created.
	CreateService bool `json:"createService"`
	// Whether the PRIMARY variant label should be added to manifests if they were missing.
	AddVariantLabelToSelector bool `json:"addVariantLabelToSelector"`
	// Whether the resources that are no longer defined in Git will be removed.
	Prune bool `json:"prune"`
}

// K8sCanaryRolloutStageOptions contains all configurable values for a K8S_CANARY_ROLLOUT stage.
type K8sCanaryRolloutStageOptions struct {
	// How many pods for CANARY workloads.
	// An integer value can be specified to indicate an absolute value of pod number.
	// Or a string suffixed by "%" to indicate an percentage value compared to the pod number of PRIMARY.
	// Default is 1 pod.
	Replicas Replicas `json:"replicas"`
	// Suffix that should be used when naming the CANARY variant's resources.
	// Default is "canary".
	Suffix string `json:"suffix"`
	// Whether the CANARY service should be created.
	CreateService bool `json:"createService"`
}

// K8sCanaryCleanStageOptions contains all configurable values for a K8S_CANARY_CLEAN stage.
type K8sCanaryCleanStageOptions struct {
}

// K8sBaselineRolloutStageOptions contains all configurable values for a K8S_BASELINE_ROLLOUT stage.
type K8sBaselineRolloutStageOptions struct {
	// How many pods for BASELINE workloads.
	// An integer value can be specified to indicate an absolute value of pod number.
	// Or a string suffixed by "%" to indicate an percentage value compared to the pod number of PRIMARY.
	// Default is 1 pod.
	Replicas Replicas `json:"replicas"`
	// Suffix that should be used when naming the BASELINE variant's resources.
	// Default is "baseline".
	Suffix string `json:"suffix"`
	// Whether the BASELINE service should be created.
	CreateService bool `json:"createService"`
}

// K8sBaselineCleanStageOptions contains all configurable values for a K8S_BASELINE_CLEAN stage.
type K8sBaselineCleanStageOptions struct {
}

// K8sTrafficRoutingStageOptions contains all configurable values for a K8S_TRAFFIC_ROUTING stage.
type K8sTrafficRoutingStageOptions struct {
	// Which variant should receive all traffic.
	// "primary" or "canary" or "baseline" can be populated.
	All string `json:"all"`
	// The percentage of traffic should be routed to PRIMARY variant.
	Primary int `json:"primary"`
	// The percentage of traffic should be routed to CANARY variant.
	Canary int `json:"canary"`
	// The percentage of traffic should be routed to BASELINE variant.
	Baseline int `json:"baseline"`
}

func (opts K8sTrafficRoutingStageOptions) Percentages() (primary, canary, baseline int) {
	switch opts.All {
	case "primary":
		primary = 100
		return
	case "canary":
		canary = 100
		return
	case "baseline":
		baseline = 100
		return
	}
	return opts.Primary, opts.Canary, opts.Baseline
}
