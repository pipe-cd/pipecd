// Copyright 2026 The PipeCD Authors.
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

package provider

import (
	"k8s.io/client-go/tools/clientcmd"
)

// ResolveDefaultNamespace returns the namespace where kubectl places a namespaced
// resource whose manifest does not specify one, for the given kubeconfig.
//
// It follows the same rules as kubectl: the namespace of the current context
// in the kubeconfig, or the namespace of the service account when running
// inside a cluster without a kubeconfig, and "default" otherwise.
func ResolveDefaultNamespace(kubeconfigPath string) string {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.ExplicitPath = kubeconfigPath

	ns, _, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{}).Namespace()
	if err != nil || ns == "" {
		return DefaultNamespace
	}
	return ns
}
