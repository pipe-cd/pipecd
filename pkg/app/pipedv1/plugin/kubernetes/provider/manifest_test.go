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
	"k8s.io/apimachinery/pkg/runtime/schema"
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
				body: &unstructured.Unstructured{
					Object: tt.initial,
				},
			}
			err := manifest.AddStringMapValues(tt.values, tt.fields...)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tt.expected, manifest.body.Object); diff != "" {
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
					groupKind: schema.ParseGroupKind("ConfigMap"),
					name:      "my-config",
					namespace: "default",
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
					groupKind: schema.ParseGroupKind("Secret"),
					name:      "my-secret",
					namespace: "default",
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
					groupKind: schema.ParseGroupKind("ConfigMap"),
					name:      "my-config",
					namespace: "default",
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
					groupKind: schema.ParseGroupKind("Secret"),
					name:      "my-secret",
					namespace: "default",
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
					groupKind: schema.ParseGroupKind("ConfigMap"),
					name:      "my-config",
					namespace: "custom-namespace",
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
					groupKind: schema.ParseGroupKind("Secret"),
					name:      "my-secret",
					namespace: "custom-namespace",
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

func TestManifest_IsDeployment(t *testing.T) {
	tests := []struct {
		name     string
		manifest string
		want     bool
	}{
		{
			name: "is deployment",
			manifest: `
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
			want: true,
		},
		{
			name: "is not deployment",
			manifest: `
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  selector:
    app: nginx
  ports:
  - protocol: TCP
    port: 80
    targetPort: 80
`,
			want: false,
		},
		{
			name: "is not deployment with custom apigroup",
			manifest: `
apiVersion: custom.io/v1
kind: Deployment
metadata:
  name: custom-deployment
spec:
  template:
    spec:
      containers:
      - name: custom
        image: custom:1.0.0
`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifest := mustParseManifests(t, strings.TrimSpace(tt.manifest))[0]
			got := manifest.IsDeployment()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestManifest_IsSecret(t *testing.T) {
	tests := []struct {
		name     string
		manifest string
		want     bool
	}{
		{
			name: "is secret",
			manifest: `
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  namespace: default
data:
  key: dmFsdWU=
`,
			want: true,
		},
		{
			name: "is not secret",
			manifest: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-config
  namespace: default
data:
  key: value
`,
			want: false,
		},
		{
			name: "is not secret with custom apigroup",
			manifest: `
apiVersion: custom.io/v1
kind: Secret
metadata:
  name: custom-secret
data:
  key: dmFsdWU=
`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifest := mustParseManifests(t, strings.TrimSpace(tt.manifest))[0]
			got := manifest.IsSecret()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestManifest_IsConfigMap(t *testing.T) {
	tests := []struct {
		name     string
		manifest string
		want     bool
	}{
		{
			name: "is configmap",
			manifest: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-config
  namespace: default
data:
  key: value
`,
			want: true,
		},
		{
			name: "is not configmap",
			manifest: `
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  namespace: default
data:
  key: dmFsdWU=
`,
			want: false,
		},
		{
			name: "is not configmap with custom apigroup",
			manifest: `
apiVersion: custom.io/v1
kind: ConfigMap
metadata:
  name: custom-config
data:
  key: value
`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifest := mustParseManifests(t, strings.TrimSpace(tt.manifest))[0]
			got := manifest.IsConfigMap()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsManagedByPiped(t *testing.T) {
	testcases := []struct {
		name       string
		manifest   Manifest
		wantResult bool
	}{
		{
			name: "managed by Piped",
			manifest: Manifest{
				body: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"metadata": map[string]interface{}{
							"annotations": map[string]interface{}{
								LabelManagedBy: ManagedByPiped,
							},
						},
					},
				},
			},
			wantResult: true,
		},
		{
			name: "not managed by Piped",
			manifest: Manifest{
				body: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"metadata": map[string]interface{}{
							"annotations": map[string]interface{}{
								"some-other-label": "some-value",
							},
						},
					},
				},
			},
			wantResult: false,
		},
		{
			name: "has owner references",
			manifest: Manifest{
				body: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"metadata": map[string]interface{}{
							"annotations": map[string]interface{}{
								LabelManagedBy: ManagedByPiped,
							},
							"ownerReferences": []interface{}{
								map[string]interface{}{
									"apiVersion": "v1",
									"kind":       "ReplicaSet",
									"name":       "example-replicaset",
								},
							},
						},
					},
				},
			},
			wantResult: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotResult := tc.manifest.IsManagedByPiped()
			assert.Equal(t, tc.wantResult, gotResult)
		})
	}
}
