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

	// Which resources should be considered as the Workload of application.
	// Empty means all Deployments.
	// e.g.
	// - kind: Deployment
	//   name: deployment-name
	// - kind: ReplicationController
	//   name: replication-controller-name
	Workloads []K8sResourceReference `json:"workloads"`

	// TODO: Define fields for KubernetesApplicationSpec.
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

	// The namespace where manifests will be applied.
	Namespace string `json:"namespace,omitempty"`

	// TODO: Define fields for KubernetesDeploymentInput.
}
