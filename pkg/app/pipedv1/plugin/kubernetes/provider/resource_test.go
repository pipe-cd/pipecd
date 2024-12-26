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
					groupKind: schema.GroupKind{Group: "", Kind: "secret"},
					namespace: "default",
					name:      "old-secret",
				},
				{
					groupKind: schema.GroupKind{Group: "", Kind: "namespace"},
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
