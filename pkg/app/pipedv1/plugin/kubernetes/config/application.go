// Copyright 2024 The PipeCD Authors.
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

// K8sResourceReference represents a reference to a Kubernetes resource.
// It is used to specify the resources which are treated as the workload of an application.
type K8sResourceReference struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

// KubernetesApplicationSpec represents an application configuration for Kubernetes application.
type KubernetesApplicationSpec struct {
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

	// The label will be configured to variant manifests used to distinguish them.
	VariantLabel KubernetesVariantLabel `json:"variantLabel"`

	// Which method should be used for traffic routing.
	TrafficRouting *KubernetesTrafficRouting `json:"trafficRouting"`
}

func (s *KubernetesApplicationSpec) Validate() error {
	// TODO: Validate KubernetesApplicationSpec fields.
	return nil
}

// KubernetesDeploymentInput represents needed input for triggering a Kubernetes deployment.
type KubernetesDeploymentInput struct {
	// List of manifest files in the application directory used to deploy.
	// Empty means all manifest files in the directory will be used.
	Manifests []string `json:"manifests,omitempty"`
	// Version of kubectl will be used.
	KubectlVersion string `json:"kubectlVersion,omitempty"`

	// Version of kustomize will be used.
	KustomizeVersion string `json:"kustomizeVersion,omitempty"`
	// List of options that should be used by Kustomize commands.
	KustomizeOptions map[string]string `json:"kustomizeOptions,omitempty"`

	// Version of helm will be used.
	HelmVersion string `json:"helmVersion,omitempty"`
	// Where to fetch helm chart.
	HelmChart *InputHelmChart `json:"helmChart,omitempty"`
	// Configurable parameters for helm commands.
	HelmOptions *InputHelmOptions `json:"helmOptions,omitempty"`

	// The namespace where manifests will be applied.
	Namespace string `json:"namespace,omitempty"`

	// Automatically create a new namespace if it does not exist.
	// Default is false.
	AutoCreateNamespace bool `json:"autoCreateNamespace,omitempty"`
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

// K8sResourcePatch represents a patch operation for a Kubernetes resource.
type K8sResourcePatch struct {
	// The target of the patch operation.
	Target K8sResourcePatchTarget `json:"target"`
	// The operations to be performed on the target.
	Ops []K8sResourcePatchOp `json:"ops"`
}

// K8sResourcePatchTarget represents the target of a patch operation for a Kubernetes resource.
type K8sResourcePatchTarget struct {
	// The reference to the Kubernetes resource.
	K8sResourceReference
	// In case you want to manipulate the YAML or JSON data specified in a field
	// of the manifest, specify that field's path. The string value of that field
	// will be used as input for the patch operations.
	// Otherwise, the whole manifest will be the target of patch operations.
	DocumentRoot string `json:"documentRoot"`
}

// K8sResourcePatchOpName represents the name of a patch operation for a Kubernetes resource.
type K8sResourcePatchOpName string

const (
	// K8sResourcePatchOpYAMLReplace is the name of the patch operation that replaces the target with a new YAML document.
	K8sResourcePatchOpYAMLReplace = "yaml-replace"
)

// K8sResourcePatchOp represents a patch operation for a Kubernetes resource.
type K8sResourcePatchOp struct {
	// The operation type.
	// Currently, only "yaml-replace" is supported.
	// Default is "yaml-replace".
	// TODO: support "yaml-add", "yaml-remove", "json-replace" and "text-regex".
	Op K8sResourcePatchOpName `json:"op" default:"yaml-replace"`
	// The path string pointing to the manipulated field.
	// E.g. "$.spec.foos[0].bar"
	Path string `json:"path"`
	// The value string whose content will be used as new value for the field.
	Value string `json:"value"`
}
