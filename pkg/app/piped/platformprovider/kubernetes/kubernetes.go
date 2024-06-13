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

package kubernetes

import (
	"errors"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/pipe-cd/pipecd/pkg/config"
)

var (
	ErrNotFound = errors.New("not found")
)

const (
	LabelManagedBy            = "pipecd.dev/managed-by"             // Always be piped.
	LabelPiped                = "pipecd.dev/piped"                  // The id of piped handling this application.
	LabelApplication          = "pipecd.dev/application"            // The application this resource belongs to.
	LabelCommitHash           = "pipecd.dev/commit-hash"            // Hash value of the deployed commit.
	LabelOriginalAPIVersion   = "pipecd.dev/original-api-version"   // The api version defined in git configuration. e.g. apps/v1
	LabelIgnoreDriftDirection = "pipecd.dev/ignore-drift-detection" // Whether the drift detection should ignore this resource.
	LabelSyncReplace          = "pipecd.dev/sync-by-replace"        // Use replace instead of apply.
	LabelServerSideApply      = "pipecd.dev/server-side-apply"      // Use server side apply instead of client side apply.
	AnnotationConfigHash      = "pipecd.dev/config-hash"            // The hash value of all mouting config resources.
	AnnotationOrder           = "pipecd.dev/order"                  // The order number of resource used to sort them before using.

	ManagedByPiped           = "piped"
	IgnoreDriftDetectionTrue = "true"
	UseReplaceEnabled        = "enabled"
	UseServerSideApply       = "true"

	kustomizationFileName = "kustomization.yaml"
)

// GetIsNamespacedResources return the map to determine whether the given GroupVersionKind is namespaced or not.
// The key is GroupVersionKind and the value is a boolean value.
// This function will get the information from the Kubernetes cluster using the given PlatformProviderKubernetesConfig.
func GetIsNamespacedResources(cp *config.PlatformProviderKubernetesConfig) (map[schema.GroupVersionKind]bool, error) {
	kubeConfig, err := clientcmd.BuildConfigFromFlags(cp.MasterURL, cp.KubeConfigPath)
	if err != nil {
		return nil, err
	}
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	groupResources, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		return nil, err
	}

	isNamespacedResources := make(map[schema.GroupVersionKind]bool)
	for _, gr := range groupResources {
		for _, resource := range gr.APIResources {
			gvk := schema.FromAPIVersionAndKind(gr.GroupVersion, resource.Kind)
			isNamespacedResources[gvk] = resource.Namespaced
		}
	}

	return isNamespacedResources, nil
}
