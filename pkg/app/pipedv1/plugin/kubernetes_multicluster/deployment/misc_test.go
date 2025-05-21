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

package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/provider"
)

func TestCheckVariantSelectorInWorkload(t *testing.T) {
	t.Parallel()

	const (
		variantLabel   = "pipecd.dev/variant"
		primaryVariant = "primary"
	)
	testcases := []struct {
		name     string
		manifest string
	}{
		{
			name: "missing variant in selector",
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
spec:
  selector:
    matchLabels:
      app: simple
  template:
    metadata:
      labels:
        app: simple
`,
		},
		{
			name: "missing variant in template labels",
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
spec:
  selector:
    matchLabels:
      app: simple
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: simple
`,
		},
		{
			name: "wrong variant in selector",
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
spec:
  selector:
    matchLabels:
      app: simple
      pipecd.dev/variant: canary
  template:
    metadata:
      labels:
        app: simple
`,
		},
		{
			name: "wrong variant in temlate labels",
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
spec:
  selector:
    matchLabels:
      app: simple
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: simple
        pipecd.dev/variant: canary
`,
		},
		{
			name: "ok",
			manifest: `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
spec:
  selector:
    matchLabels:
      app: simple
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: simple
        pipecd.dev/variant: primary
`,
		},
	}

	expected := `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple
spec:
  selector:
    matchLabels:
      app: simple
      pipecd.dev/variant: primary
  template:
    metadata:
      labels:
        app: simple
        pipecd.dev/variant: primary
`
	generatedManifests, err := provider.ParseManifests(expected)
	require.NoError(t, err)
	require.Equal(t, 1, len(generatedManifests))

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			manifests, err := provider.ParseManifests(tc.manifest)
			require.NoError(t, err)
			require.Equal(t, 1, len(manifests))

			err = ensureVariantSelectorInWorkload(manifests[0], variantLabel, primaryVariant)
			assert.NoError(t, err)
			assert.Equal(t, generatedManifests[0], manifests[0])
		})
	}

}

func TestAddVariantLabelsAndAnnotations(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name         string
		inputYAML    string
		variantLabel string
		variant      string
		wantLabels   map[string]string
		wantAnnots   map[string]string
	}{
		{
			name: "single manifest",
			inputYAML: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
`,
			variantLabel: "pipecd.dev/variant",
			variant:      "primary",
			wantLabels:   map[string]string{"pipecd.dev/variant": "primary"},
			wantAnnots:   map[string]string{"pipecd.dev/variant": "primary"},
		},
		{
			name: "multiple manifests",
			inputYAML: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: config1
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: config2
`,
			variantLabel: "custom/label",
			variant:      "canary",
			wantLabels:   map[string]string{"custom/label": "canary"},
			wantAnnots:   map[string]string{"custom/label": "canary"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			manifests, err := provider.ParseManifests(tc.inputYAML)
			require.NoError(t, err)
			require.NotEmpty(t, manifests)

			addVariantLabelsAndAnnotations(manifests, tc.variantLabel, tc.variant)

			for _, m := range manifests {
				labelsMap, _, err := m.NestedMap("metadata", "labels")
				require.NoError(t, err)
				labels := map[string]string{}
				for k, v := range labelsMap {
					if strVal, ok := v.(string); ok {
						labels[k] = strVal
					}
				}
				for k, v := range tc.wantLabels {
					assert.Equal(t, v, labels[k], "label %q should be %q", k, v)
				}
				annots := m.GetAnnotations()
				for k, v := range tc.wantAnnots {
					assert.Equal(t, v, annots[k], "annotation %q should be %q", k, v)
				}
			}
		})
	}
}
