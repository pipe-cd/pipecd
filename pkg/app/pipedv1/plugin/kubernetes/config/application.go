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

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/creasty/defaults"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
)

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
	// Version of kubectl will be used.
	KubectlVersion string `json:"kubectlVersion,omitempty"`

	// The namespace where manifests will be applied.
	Namespace string `json:"namespace,omitempty"`

	// Automatically create a new namespace if it does not exist.
	// Default is false.
	AutoCreateNamespace bool `json:"autoCreateNamespace,omitempty"`

	// TODO: Define fields for KubernetesDeploymentInput.
	MultiTargets []KubernetesMultiTarget `json:"multiTargets,omitempty"`
}

type KubernetesMultiTarget struct {
	Target         KubernetesMultiTargetDeployTarget `json:"target"`
	Manifests      []string                          `json:"manifests,omitempty"`
	KubectlVersion string                            `json:"kubectlVersion,omitempty"`
	KustomizeDir   string                            `json:"kustomizeDir,omitempty"`
}

type KubernetesMultiTargetDeployTarget struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
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

type KubernetesDeployTargetConfig struct {
	Name string `json:"-"`
	// The master URL of the kubernetes cluster.
	// Empty means in-cluster.
	MasterURL string `json:"masterURL,omitempty"`
	// The path to the kubeconfig file.
	// Empty means in-cluster.
	KubeConfigPath string `json:"kubeConfigPath,omitempty"`
	// Version of kubectl will be used.
	KubectlVersion string `json:"kubectlVersion"`
}

// K8sSyncStageOptions contains all configurable values for a K8S_SYNC stage.
type K8sSyncStageOptions struct {
	// Whether the PRIMARY variant label should be added to manifests if they were missing.
	AddVariantLabelToSelector bool `json:"addVariantLabelToSelector"`
	// Whether the resources that are no longer defined in Git should be removed or not.
	Prune bool `json:"prune"`
}

// FindDeployTarget finds the deploy target configuration by the given name.
func FindDeployTarget(cfg *config.PipedPlugin, name string) (KubernetesDeployTargetConfig, error) {
	if cfg == nil {
		return KubernetesDeployTargetConfig{}, errors.New("missing plugin configuration")
	}

	deployTarget := cfg.FindDeployTarget(name)
	if deployTarget == nil {
		return KubernetesDeployTargetConfig{}, errors.New("missing deploy target configuration")
	}

	var targetCfg KubernetesDeployTargetConfig

	if err := json.Unmarshal(deployTarget.Config, &targetCfg); err != nil { // TODO: not decode here but in the initialization of the plugin.
		return KubernetesDeployTargetConfig{}, fmt.Errorf("failed to unmarshal deploy target configuration: %w", err)
	}

	if err := defaults.Set(&targetCfg); err != nil {
		return KubernetesDeployTargetConfig{}, fmt.Errorf("failed to set default values for deploy target configuration: %w", err)
	}

	targetCfg.Name = deployTarget.Name
	return targetCfg, nil
}
