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
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestParseManifests(t *testing.T) {
	maker := func(name, kind string, metadata map[string]interface{}) Manifest {
		return Manifest{
			Key: ResourceKey{
				APIVersion: "v1",
				Kind:       kind,
				Name:       name,
				Namespace:  "default",
			},
			u: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       kind,
					"metadata":   metadata,
				},
			},
		}
	}

	testcases := []struct {
		name      string
		manifests string
		want      []Manifest
	}{
		{
			name: "empty1",
		},
		{
			name:      "empty2",
			manifests: "---",
		},
		{
			name:      "empty3",
			manifests: "\n---",
		},
		{
			name:      "empty4",
			manifests: "\n---\n",
		},
		{
			name:      "multiple empty manifests",
			manifests: "---\n---\n---\n---\n---\n",
		},
		{
			name: "one manifest",
			manifests: `---
apiVersion: v1
kind: ConfigMap
metadata:
  name: envoy-config
  creationTimestamp: "2022-12-09T01:23:45Z"
`,
			want: []Manifest{
				maker("envoy-config", "ConfigMap", map[string]interface{}{
					"name":              "envoy-config",
					"creationTimestamp": "2022-12-09T01:23:45Z",
				}),
			},
		},
		{
			name: "contains new line at the end of file",
			manifests: `
apiVersion: v1
kind: Kind1
metadata:
  name: config
  extra: |
    single-new-line
`,
			want: []Manifest{
				maker("config", "Kind1", map[string]interface{}{
					"name":  "config",
					"extra": "single-new-line\n",
				}),
			},
		},
		{
			name: "not contains new line at the end of file",
			manifests: `
apiVersion: v1
kind: Kind1
metadata:
  name: config
  extra: |
    no-new-line`,
			want: []Manifest{
				maker("config", "Kind1", map[string]interface{}{
					"name":  "config",
					"extra": "no-new-line",
				}),
			},
		},
		{
			name: "multiple manifests",
			manifests: `
apiVersion: v1
kind: Kind1
metadata:
  name: config1
  extra: |-
    no-new-line
---
apiVersion: v1
kind: Kind2
metadata:
  name: config2
  extra: |
    single-new-line-1
---
apiVersion: v1
kind: Kind3
metadata:
  name: config3
  extra: |
    single-new-line-2


---
apiVersion: v1
kind: Kind4
metadata:
  name: config4
  extra: |+
    multiple-new-line-1


---
apiVersion: v1
kind: Kind5
metadata:
  name: config5
  extra: |+
    multiple-new-line-2


`,
			want: []Manifest{
				maker("config1", "Kind1", map[string]interface{}{
					"name":  "config1",
					"extra": "no-new-line",
				}),
				maker("config2", "Kind2", map[string]interface{}{
					"name":  "config2",
					"extra": "single-new-line-1\n",
				}),
				maker("config3", "Kind3", map[string]interface{}{
					"name":  "config3",
					"extra": "single-new-line-2\n",
				}),
				maker("config4", "Kind4", map[string]interface{}{
					"name":  "config4",
					"extra": "multiple-new-line-1\n\n\n",
				}),
				maker("config5", "Kind5", map[string]interface{}{
					"name":  "config5",
					"extra": "multiple-new-line-2\n\n\n",
				}),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := ParseManifests(tc.manifests)
			require.NoError(t, err)
			assert.ElementsMatch(t, m, tc.want)
		})
	}
}
