// Copyright 2022 The PipeCD Authors.
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
			manifests: "",
		},
		{
			name:      "empty3",
			manifests: "---",
		},
		{
			name:      "empty4",
			manifests: "\n---",
		},
		{
			name:      "empty5",
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
			name: "multiple manifests",
			manifests: `
apiVersion: v1
kind: Kind1
metadata:
  name: config1
---
apiVersion: v1
kind: Kind2
metadata:
  name: config2
---
apiVersion: v1
kind: Kind3
metadata:
  name: config3
			`,
			want: []Manifest{
				maker("config1", "Kind1", map[string]interface{}{"name": "config1"}),
				maker("config2", "Kind2", map[string]interface{}{"name": "config2"}),
				maker("config3", "Kind3", map[string]interface{}{"name": "config3"}),
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
