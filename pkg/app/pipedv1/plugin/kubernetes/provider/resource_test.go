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

package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestResourceKey_normalizeNamespace(t *testing.T) {
	tests := []struct {
		name        string
		resourceKey ResourceKey
		expected    ResourceKey
	}{
		{
			name: "default namespace",
			resourceKey: ResourceKey{
				groupKind: schema.GroupKind{Group: "apps", Kind: "Deployment"},
				namespace: DefaultNamespace,
				name:      "test-deployment",
			},
			expected: ResourceKey{
				groupKind: schema.GroupKind{Group: "apps", Kind: "Deployment"},
				namespace: "",
				name:      "test-deployment",
			},
		},
		{
			name: "non-default namespace",
			resourceKey: ResourceKey{
				groupKind: schema.GroupKind{Group: "apps", Kind: "Deployment"},
				namespace: "custom-namespace",
				name:      "test-deployment",
			},
			expected: ResourceKey{
				groupKind: schema.GroupKind{Group: "apps", Kind: "Deployment"},
				namespace: "custom-namespace",
				name:      "test-deployment",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.resourceKey.normalizeNamespace()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestFindRemoveResources(t *testing.T) {
	tests := []struct {
		name                           string
		manifestsYAML                  string
		namespacedLiveResourcesYAML    string
		clusterScopedLiveResourcesYAML string
		expectedRemoveKeys             []ResourceKey
	}{
		{
			name: "find remove resources",
			manifestsYAML: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-configmap
  namespace: default
  annotations:
    "pipecd.dev/managed-by": "piped"
---
apiVersion: v1
kind: Secret
metadata:
  name: test-secret
  namespace: default
  annotations:
    "pipecd.dev/managed-by": "piped"
`,
			namespacedLiveResourcesYAML: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-configmap
  namespace: default
  annotations:
    "pipecd.dev/managed-by": "piped"
---
apiVersion: v1
kind: Secret
metadata:
  name: test-secret
  namespace: default
  annotations:
    "pipecd.dev/managed-by": "piped"
---
apiVersion: v1
kind: Secret
metadata:
  name: old-secret
  namespace: default
  annotations:
    "pipecd.dev/managed-by": "piped"
`,
			clusterScopedLiveResourcesYAML: `
apiVersion: v1
kind: Namespace
metadata:
  name: test-namespace
  annotations:
    "pipecd.dev/managed-by": "piped"
`,
			expectedRemoveKeys: []ResourceKey{
				{
					groupKind: schema.GroupKind{Group: "", Kind: "Secret"},
					namespace: "default",
					name:      "old-secret",
				},
				{
					groupKind: schema.GroupKind{Group: "", Kind: "Namespace"},
					namespace: "",
					name:      "test-namespace",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifests := mustParseManifests(t, tt.manifestsYAML)

			namespacedLiveResources := mustParseManifests(t, tt.namespacedLiveResourcesYAML)

			clusterScopedLiveResources := mustParseManifests(t, tt.clusterScopedLiveResourcesYAML)

			removeKeys := FindRemoveResources(manifests, namespacedLiveResources, clusterScopedLiveResources)
			assert.ElementsMatch(t, tt.expectedRemoveKeys, removeKeys)
		})
	}
}
