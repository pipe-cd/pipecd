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
	"encoding/json"

	"github.com/creasty/defaults"
)

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

func (k *KubernetesDeployTargetConfig) UnmarshalJSON(data []byte) error {
	type alias KubernetesDeployTargetConfig

	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	*k = KubernetesDeployTargetConfig(a)
	if err := defaults.Set(k); err != nil {
		return err
	}

	return nil
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
