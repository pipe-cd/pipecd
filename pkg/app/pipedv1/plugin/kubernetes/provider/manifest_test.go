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
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestManifest_AddStringMapValues(t *testing.T) {
	tests := []struct {
		name     string
		initial  map[string]interface{}
		values   map[string]string
		fields   []string
		expected map[string]interface{}
	}{
		{
			name: "add new values to empty map",
			initial: map[string]interface{}{
				"metadata": map[string]interface{}{
					"annotations": map[string]interface{}{},
				},
			},
			values: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			fields: []string{"metadata", "annotations"},
			expected: map[string]interface{}{
				"metadata": map[string]interface{}{
					"annotations": map[string]interface{}{
						"key1": "value1",
						"key2": "value2",
					},
				},
			},
		},
		{
			name: "override existing values",
			initial: map[string]interface{}{
				"metadata": map[string]interface{}{
					"annotations": map[string]interface{}{
						"key1": "oldvalue1",
					},
				},
			},
			values: map[string]string{
				"key1": "newvalue1",
				"key2": "value2",
			},
			fields: []string{"metadata", "annotations"},
			expected: map[string]interface{}{
				"metadata": map[string]interface{}{
					"annotations": map[string]interface{}{
						"key1": "newvalue1",
						"key2": "value2",
					},
				},
			},
		},
		{
			name: "add values to non-existing map",
			initial: map[string]interface{}{
				"metadata": map[string]interface{}{},
			},
			values: map[string]string{
				"key1": "value1",
			},
			fields: []string{"metadata", "annotations"},
			expected: map[string]interface{}{
				"metadata": map[string]interface{}{
					"annotations": map[string]interface{}{
						"key1": "value1",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifest := Manifest{
				Body: &unstructured.Unstructured{
					Object: tt.initial,
				},
			}
			err := manifest.AddStringMapValues(tt.values, tt.fields...)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tt.expected, manifest.Body.Object); diff != "" {
				t.Errorf("unexpected result (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFindConfigsAndSecrets(t *testing.T) {
	tests := []struct {
		name      string
		manifests []string
		want      map[ResourceKey]Manifest
	}{
		{
			name: "find ConfigMap and Secret",
			manifests: []string{
				`
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-config
  namespace: default
data:
  key: value
`,
				`
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  namespace: default
data:
  key: dmFsdWU=
`,
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
  spec:
    containers:
    - name: nginx
    image: nginx:1.19.3
`,
			},
			want: map[ResourceKey]Manifest{
				{
					APIVersion: "v1",
					Kind:       "ConfigMap",
					Name:       "my-config",
					Namespace:  "default",
				}: mustParseManifests(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-config
  namespace: default
data:
  key: value
`)[0],
				{
					APIVersion: "v1",
					Kind:       "Secret",
					Name:       "my-secret",
					Namespace:  "default",
				}: mustParseManifests(t, `
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  namespace: default
data:
  key: dmFsdWU=
`)[0],
			},
		},
		{
			name: "no ConfigMap or Secret",
			manifests: []string{
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
  spec:
    containers:
    - name: nginx
    image: nginx:1.19.3
`,
			},
			want: map[ResourceKey]Manifest{},
		},
		{
			name: "only ConfigMap",
			manifests: []string{
				`
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-config
  namespace: default
data:
  key: value
`,
			},
			want: map[ResourceKey]Manifest{
				{
					APIVersion: "v1",
					Kind:       "ConfigMap",
					Name:       "my-config",
					Namespace:  "default",
				}: mustParseManifests(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-config
  namespace: default
data:
  key: value
`)[0],
			},
		},
		{
			name: "only Secret",
			manifests: []string{
				`
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  namespace: default
data:
  key: dmFsdWU=
`,
			},
			want: map[ResourceKey]Manifest{
				{
					APIVersion: "v1",
					Kind:       "Secret",
					Name:       "my-secret",
					Namespace:  "default",
				}: mustParseManifests(t, `
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  namespace: default
data:
  key: dmFsdWU=
`)[0],
			},
		},
		{
			name: "non-default namespace",
			manifests: []string{
				`
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-config
  namespace: custom-namespace
data:
  key: value
`,
				`
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  namespace: custom-namespace
data:
  key: dmFsdWU=
`,
			},
			want: map[ResourceKey]Manifest{
				{
					APIVersion: "v1",
					Kind:       "ConfigMap",
					Name:       "my-config",
					Namespace:  "custom-namespace",
				}: mustParseManifests(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-config
  namespace: custom-namespace
data:
  key: value
`)[0],
				{
					APIVersion: "v1",
					Kind:       "Secret",
					Name:       "my-secret",
					Namespace:  "custom-namespace",
				}: mustParseManifests(t, `
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  namespace: custom-namespace
data:
  key: dmFsdWU=
`)[0],
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var manifests []Manifest
			for _, data := range tt.manifests {
				manifests = append(manifests, mustParseManifests(t, data)...)
			}
			got := FindConfigsAndSecrets(manifests)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFindSameManifests(t *testing.T) {
	tests := []struct {
		name string
		olds []string
		news []string
		want []WorkloadPair
	}{
		{
			name: "single updated workload",
			olds: []string{
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
`,
			},
			news: []string{
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.4
`,
			},
			want: []WorkloadPair{
				{
					Old: mustParseManifests(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
`)[0],
					New: mustParseManifests(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.4
`)[0],
				},
			},
		},
		{
			name: "multiple updated workloads",
			olds: []string{
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
`,
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-deployment
spec:
  template:
    spec:
      containers:
      - name: redis
        image: redis:6.0.9
`,
			},
			news: []string{
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.4
`,
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-deployment
spec:
  template:
    spec:
      containers:
      - name: redis
        image: redis:6.0.10
`,
			},
			want: []WorkloadPair{
				{
					Old: mustParseManifests(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
`)[0],
					New: mustParseManifests(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.4
`)[0],
				},
				{
					Old: mustParseManifests(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-deployment
spec:
  template:
    spec:
      containers:
      - name: redis
        image: redis:6.0.9
`)[0],
					New: mustParseManifests(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-deployment
spec:
  template:
    spec:
      containers:
      - name: redis
        image: redis:6.0.10
`)[0],
				},
			},
		},
		{
			name: "no updated workloads",
			olds: []string{
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
`,
			},
			news: []string{
				`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-deployment
spec:
  template:
    spec:
      containers:
      - name: redis
        image: redis:7.0.0
`,
			},
			want: []WorkloadPair{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldManifests := mustParseManifests(t, strings.Join(tt.olds, "\n---\n"))
			newManifests := mustParseManifests(t, strings.Join(tt.news, "\n---\n"))
			got := FindSameManifests(oldManifests, newManifests)
			assert.Equal(t, tt.want, got)
		})
	}
}
