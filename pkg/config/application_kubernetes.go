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

package config

// KubernetesApplicationSpec represents an application configuration for Kubernetes application.
type KubernetesApplicationSpec struct {
	GenericApplicationSpec
	// Input for Kubernetes deployment such as kubectl version, helm version, manifests filter...
	Input KubernetesDeploymentInput `json:"input"`
	// Configuration for quick sync.
	QuickSync K8sSyncStageOptions `json:"quickSync"`
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
	TrafficRouting *KubernetesTrafficRouting `json:"trafficRouting"`
	// The label will be configured to variant manifests used to distinguish them.
	VariantLabel KubernetesVariantLabel `json:"variantLabel"`
	// List of route configurations to resolve the platform provider for application resources.
	// Each resource will be checked over the match conditions of each route.
	// If matches, it will be applied to the route's provider,
	// otherwise, it will be fallen through the next route to check.
	// Any resource which does not match any specified route will be applied
	// to the default platform provider which had been specified while registering the application.
	ResourceRoutes []KubernetesResourceRoute `json:"resourceRoutes"`
}

// Validate returns an error if any wrong configuration value was found.
func (s *KubernetesApplicationSpec) Validate() error {
	if err := s.GenericApplicationSpec.Validate(); err != nil {
		return err
	}
	return nil
}

type KubernetesVariantLabel struct {
	// The key of the label.
	// Default is pipecd.dev/variant.
	Key string `json:"key" default:"pipecd.dev/variant"`
	// The label value for PRIMARY variant.
	// Default is primary.
	PrimaryValue string `json:"primaryValue" default:"primary"`
	// The label value for CANARY variant.
	// Default is canary.
	CanaryValue string `json:"canaryValue" default:"canary"`
	// The label value for BASELINE variant.
	// Default is baseline.
	BaselineValue string `json:"baselineValue" default:"baseline"`
}

// KubernetesDeploymentInput represents needed input for triggering a Kubernetes deployment.
type KubernetesDeploymentInput struct {
	// List of manifest files in the application directory used to deploy.
	// Empty means all manifest files in the directory will be used.
	Manifests []string `json:"manifests"`
	// Version of kubectl will be used.
	KubectlVersion string `json:"kubectlVersion"`

	// Version of kustomize will be used.
	KustomizeVersion string `json:"kustomizeVersion"`
	// List of options that should be used by Kustomize commands.
	KustomizeOptions map[string]string `json:"kustomizeOptions"`

	// Version of helm will be used.
	HelmVersion string `json:"helmVersion"`
	// Where to fetch helm chart.
	HelmChart *InputHelmChart `json:"helmChart"`
	// Configurable parameters for helm commands.
	HelmOptions *InputHelmOptions `json:"helmOptions"`

	// The namespace where manifests will be applied.
	Namespace string `json:"namespace"`

	// Automatically reverts all deployment changes on failure.
	// Default is true.
	AutoRollback *bool `json:"autoRollback,omitempty" default:"true"`

	// Automatically create a new namespace if it does not exist.
	// Default is false.
	AutoCreateNamespace *bool `json:"autoCreateNamespace,omitempty" default:"false"`
}

type InputHelmChart struct {
	// Git remote address where the chart is placing.
	// Empty means the same repository.
	GitRemote string `json:"gitRemote"`
	// The commit SHA or tag for remote git.
	Ref string `json:"ref"`
	// Relative path from the repository root directory to the chart directory.
	Path string `json:"path"`

	// The name of an added Helm Chart Repository.
	Repository string `json:"repository"`
	Name       string `json:"name"`
	Version    string `json:"version"`
	// Whether to skip TLS certificate checks for the repository or not.
	// This option will automatically set the value of HelmChartRepository.Insecure.
	Insecure bool `json:"-"`
}

type InputHelmOptions struct {
	// The release name of helm deployment.
	// By default the release name is equal to the application name.
	ReleaseName string `json:"releaseName"`
	// List of value files should be loaded.
	ValueFiles []string `json:"valueFiles"`
	// List of file path for values.
	SetFiles map[string]string
	// Set of supported Kubernetes API versions.
	APIVersions []string `json:"apiVersions"`
	// Kubernetes version used for Capabilities.KubeVersion
	KubeVersion string `json:"kubeVersion"`
}

type KubernetesTrafficRoutingMethod string

const (
	KubernetesTrafficRoutingMethodPodSelector KubernetesTrafficRoutingMethod = "podselector"
	KubernetesTrafficRoutingMethodIstio       KubernetesTrafficRoutingMethod = "istio"
	KubernetesTrafficRoutingMethodSMI         KubernetesTrafficRoutingMethod = "smi"
)

type KubernetesTrafficRouting struct {
	Method KubernetesTrafficRoutingMethod `json:"method"`
	Istio  *IstioTrafficRouting           `json:"istio"`
}

// DetermineKubernetesTrafficRoutingMethod determines the routing method should be used based on the TrafficRouting config.
// The default is PodSelector: the way by updating the selector in Service to switching all of traffic.
func DetermineKubernetesTrafficRoutingMethod(cfg *KubernetesTrafficRouting) KubernetesTrafficRoutingMethod {
	if cfg == nil {
		return KubernetesTrafficRoutingMethodPodSelector
	}
	if cfg.Method == "" {
		return KubernetesTrafficRoutingMethodPodSelector
	}
	return cfg.Method
}

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

type K8sResourceReference struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

// K8sSyncStageOptions contains all configurable values for a K8S_SYNC stage.
type K8sSyncStageOptions struct {
	// Whether the PRIMARY variant label should be added to manifests if they were missing.
	AddVariantLabelToSelector bool `json:"addVariantLabelToSelector"`
	// Whether the resources that are no longer defined in Git should be removed or not.
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
	// Whether the resources that are no longer defined in Git should be removed or not.
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
	// List of patches used to customize manifests for CANARY variant.
	Patches []K8sResourcePatch
}

type K8sResourcePatch struct {
	Target K8sResourcePatchTarget `json:"target"`
	Ops    []K8sResourcePatchOp   `json:"ops"`
}

type K8sResourcePatchTarget struct {
	K8sResourceReference
	// In case you want to manipulate the YAML or JSON data specified in a field
	// of the manifest, specify that field's path. The string value of that field
	// will be used as input for the patch operations.
	// Otherwise, the whole manifest will be the target of patch operations.
	DocumentRoot string `json:"documentRoot"`
}

type K8sResourcePatchOpName string

const (
	K8sResourcePatchOpYAMLReplace = "yaml-replace"
)

type K8sResourcePatchOp struct {
	// The operation type.
	// This must be one of "yaml-replace", "yaml-add", "yaml-remove", "json-replace" or "text-regex".
	// Default is "yaml-replace".
	Op K8sResourcePatchOpName `json:"op" default:"yaml-replace"`
	// The path string pointing to the manipulated field.
	// E.g. "$.spec.foos[0].bar"
	Path string `json:"path"`
	// The value string whose content will be used as new value for the field.
	Value string `json:"value"`
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
	Primary Percentage `json:"primary"`
	// The percentage of traffic should be routed to CANARY variant.
	Canary Percentage `json:"canary"`
	// The percentage of traffic should be routed to BASELINE variant.
	Baseline Percentage `json:"baseline"`
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
	return opts.Primary.Int(), opts.Canary.Int(), opts.Baseline.Int()
}

type KubernetesResourceRoute struct {
	Provider KubernetesProviderMatcher       `json:"provider"`
	Match    *KubernetesResourceRouteMatcher `json:"match"`
}

type KubernetesResourceRouteMatcher struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type KubernetesProviderMatcher struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
}
