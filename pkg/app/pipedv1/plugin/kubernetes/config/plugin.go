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

type KubernetesPluginConfig struct {
	// List of helm chart repositories that should be added while starting up.
	ChartRepositories []HelmChartRepository `json:"chartRepositories,omitempty"`
	// List of helm chart registries that should be logged in while starting up.
	ChartRegistries []HelmChartRegistry `json:"chartRegistries,omitempty"`
}

type HelmChartRepositoryType string

const (
	HTTPHelmChartRepository HelmChartRepositoryType = "HTTP"
)

type HelmChartRepository struct {
	// The repository type. Only HTTP is supported.
	Type HelmChartRepositoryType `json:"type" default:"HTTP"`

	// Configuration for HTTP type.
	// The name of the Helm chart repository.
	Name string `json:"name,omitempty"`
	// The address to the Helm chart repository.
	Address string `json:"address,omitempty"`
	// Username used for the repository backed by HTTP basic authentication.
	Username string `json:"username,omitempty"`
	// Password used for the repository backed by HTTP basic authentication.
	Password string `json:"password,omitempty"`
	// Whether to skip TLS certificate checks for the repository or not.
	Insecure bool `json:"insecure"`
}

// HelmChartRegistryType represents the type of Helm chart registry.
type HelmChartRegistryType string

// The registry types that hosts Helm charts.
const (
	OCIHelmChartRegistry HelmChartRegistryType = "OCI"
)

// HelmChartRegistry represents the configuration for a Helm chart registry.
type HelmChartRegistry struct {
	// The registry type. Currently, only OCI is supported.
	Type HelmChartRegistryType `json:"type" default:"OCI"`

	// The address to the Helm chart registry.
	Address string `json:"address"`
	// Username used for the registry authentication.
	Username string `json:"username,omitempty"`
	// Password used for the registry authentication.
	Password string `json:"password,omitempty"`
}

// IsOCI checks if the registry is an OCI registry.
func (r *HelmChartRegistry) IsOCI() bool {
	return r.Type == OCIHelmChartRegistry
}

// KubernetesDeployTargetConfig represents the configuration for a Kubernetes deployment target.
type KubernetesDeployTargetConfig struct {
	// The master URL of the kubernetes cluster.
	// Empty means in-cluster.
	MasterURL string `json:"masterURL,omitempty"`
	// The path to the kubeconfig file.
	// Empty means in-cluster.
	KubeConfigPath string `json:"kubeConfigPath,omitempty"`
	// Version of kubectl will be used.
	KubectlVersion string `json:"kubectlVersion"`
	// Configuration for application resource informer.
	AppStateInformer KubernetesAppStateInformer `json:"appStateInformer"`
}

// KubernetesAppStateInformer represents the configuration for application resource informer.
type KubernetesAppStateInformer struct {
	// Only watches the specified namespace.
	// Empty means watching all namespaces.
	Namespace string `json:"namespace,omitempty"`
	// List of resources that should be added to the watching targets.
	IncludeResources []KubernetesResourceMatcher `json:"includeResources,omitempty"`
	// List of resources that should be ignored from the watching targets.
	ExcludeResources []KubernetesResourceMatcher `json:"excludeResources,omitempty"`
}

// KubernetesResourceMatcher represents the matcher for a Kubernetes resource.
type KubernetesResourceMatcher struct {
	// The APIVersion of the kubernetes resource.
	APIVersion string `json:"apiVersion,omitempty"`
	// The kind name of the kubernetes resource.
	// Empty means all kinds are matching.
	Kind string `json:"kind,omitempty"`
}
