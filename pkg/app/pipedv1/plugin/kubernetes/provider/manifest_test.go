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
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

func TestManifest_AddStringMapValues(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

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
	t.Parallel()

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
			t.Parallel()

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
	t.Parallel()

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
			t.Parallel()

			oldManifests := mustParseManifests(t, strings.Join(tt.olds, "\n---\n"))
			newManifests := mustParseManifests(t, strings.Join(tt.news, "\n---\n"))
			got := FindSameManifests(oldManifests, newManifests)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestManifest_IsDeployment(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			manifest := mustParseManifests(t, strings.TrimSpace(tt.manifest))[0]
			got := manifest.IsDeployment()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestManifest_IsStatefulSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		manifest string
		want     bool
	}{
		{
			name: "is statefulset",
			manifest: `
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: my-statefulset
  namespace: default
spec:
  serviceName: "nginx"
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.19.3
`,
			want: true,
		},
		{
			name: "is not statefulset",
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
			want: false,
		},
		{
			name: "is not statefulset with custom apigroup",
			manifest: `
apiVersion: custom.io/v1
kind: StatefulSet
metadata:
  name: custom-statefulset
spec:
  serviceName: "custom"
  replicas: 1
  selector:
    matchLabels:
      app: custom
  template:
    metadata:
      labels:
        app: custom
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
			t.Parallel()
			manifest := mustParseManifests(t, strings.TrimSpace(tt.manifest))[0]
			got := manifest.IsStatefulSet()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestManifest_IsSecret(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			manifest := mustParseManifests(t, strings.TrimSpace(tt.manifest))[0]
			got := manifest.IsSecret()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestManifest_IsConfigMap(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			manifest := mustParseManifests(t, strings.TrimSpace(tt.manifest))[0]
			got := manifest.IsConfigMap()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestManifest_IsService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		manifest string
		want     bool
	}{
		{
			name: "is service",
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
			want: true,
		},
		{
			name: "is not service",
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
			want: false,
		},
		{
			name: "is not service with custom apigroup",
			manifest: `
apiVersion: custom.io/v1
kind: Service
metadata:
  name: custom-service
spec:
  selector:
    app: custom
  ports:
  - protocol: TCP
    port: 80
    targetPort: 80
`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			manifest := mustParseManifests(t, strings.TrimSpace(tt.manifest))[0]
			got := manifest.IsService()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestManifest_IsWorkload(t *testing.T) {
	t.Parallel()

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
			name: "is replicaset",
			manifest: `
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: nginx-replicaset
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
			name: "is daemonset",
			manifest: `
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: nginx-daemonset
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
			name: "is pod",
			manifest: `
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
spec:
  containers:
  - name: nginx
    image: nginx:1.19.3
`,
			want: true,
		},
		{
			name: "is not workload (service)",
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
			name: "is not workload (custom apigroup)",
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
			t.Parallel()

			manifest := mustParseManifests(t, strings.TrimSpace(tt.manifest))[0]
			got := manifest.IsWorkload()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsManagedByPiped(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			gotResult := tc.manifest.IsManagedByPiped()
			assert.Equal(t, tc.wantResult, gotResult)
		})
	}
}

func TestManifest_ToResourceState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		manifest     Manifest
		deployTarget string
		want         sdk.ResourceState
	}{
		{
			name: "no owner references",
			manifest: Manifest{
				body: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"metadata": map[string]interface{}{
							"name":              "nginx-deployment",
							"namespace":         "default",
							"uid":               "12345",
							"creationTimestamp": "2023-10-01T00:00:00Z",
						},
						"kind":       "Deployment",
						"apiVersion": "apps/v1",
					},
				},
			},
			deployTarget: "target1",
			want: sdk.ResourceState{
				ID:                "12345",
				Name:              "nginx-deployment",
				ParentIDs:         nil,
				HealthStatus:      sdk.ResourceHealthStateUnknown,
				HealthDescription: "The number of desired replicas is unspecified",
				ResourceType:      "Deployment",
				ResourceMetadata: map[string]string{
					"Namespace":   "default",
					"API Version": "apps/v1",
					"Kind":        "Deployment",
				},
				DeployTarget: "target1",
				CreatedAt:    time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC).Local(),
			},
		},
		{
			name: "with owner references",
			manifest: Manifest{
				body: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"metadata": map[string]interface{}{
							"name":              "nginx-deployment",
							"namespace":         "default",
							"uid":               "12345",
							"creationTimestamp": "2023-10-01T00:00:00Z",
							"ownerReferences": []interface{}{
								map[string]interface{}{
									"apiVersion": "apps/v1",
									"kind":       "ReplicaSet",
									"name":       "nginx-replicaset",
									"uid":        "67890",
								},
							},
						},
						"kind":       "Deployment",
						"apiVersion": "apps/v1",
					},
				},
			},
			deployTarget: "target2",
			want: sdk.ResourceState{
				ID:                "12345",
				Name:              "nginx-deployment",
				ParentIDs:         []string{"67890"},
				HealthStatus:      sdk.ResourceHealthStateUnknown,
				HealthDescription: "The number of desired replicas is unspecified",
				ResourceType:      "Deployment",
				ResourceMetadata: map[string]string{
					"Namespace":   "default",
					"API Version": "apps/v1",
					"Kind":        "Deployment",
				},
				DeployTarget: "target2",
				CreatedAt:    time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC).Local(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.manifest.ToResourceState(tt.deployTarget)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestManifest_ConvertToStructuredObject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		yaml    string
		want    interface{}
		wantErr bool
	}{
		{
			name: "ConfigMap conversion",
			yaml: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  namespace: default
data:
  key: value
`,
			want: &corev1.ConfigMap{
				TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "ConfigMap"},
				ObjectMeta: metav1.ObjectMeta{Name: "test-config", Namespace: "default"},
				Data:       map[string]string{"key": "value"},
			},
			wantErr: false,
		},
		{
			name: "Secret conversion",
			yaml: `
apiVersion: v1
kind: Secret
metadata:
  name: test-secret
  namespace: default
data:
  password: cGFzc3dvcmQ=
  username: dXNlcg==
`,
			want: &corev1.Secret{
				TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Secret"},
				ObjectMeta: metav1.ObjectMeta{Name: "test-secret", Namespace: "default"},
				Data: map[string][]byte{
					"password": []byte("password"),
					"username": []byte("user"),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			manifests := mustParseManifests(t, tt.yaml)
			require.NotEmpty(t, manifests)
			require.Len(t, manifests, 1)

			switch want := tt.want.(type) {
			case *corev1.ConfigMap:
				var got corev1.ConfigMap
				err := manifests[0].ConvertToStructuredObject(&got)
				if tt.wantErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				assert.Equal(t, want.Name, got.Name)
				assert.Equal(t, want.Namespace, got.Namespace)
				assert.Equal(t, want.Data, got.Data)
			case *corev1.Secret:
				var got corev1.Secret
				err := manifests[0].ConvertToStructuredObject(&got)
				if tt.wantErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				assert.Equal(t, want.Name, got.Name)
				assert.Equal(t, want.Namespace, got.Namespace)
				assert.Equal(t, want.Data, got.Data)
			default:
				t.Fatalf("unsupported want type: %T", tt.want)
			}
		})
	}
}

func TestDeepCopyManifests(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		yaml   string
		mutate func(orig, copy []Manifest)
	}{
		{
			name: "deep copy: changing label in copy does not affect original",
			yaml: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  labels:
    foo: bar
`,
			mutate: func(orig, copy []Manifest) {
				copy[0].AddLabels(map[string]string{"foo": "baz"})
			},
		},
		{
			name: "deep copy: changing annotation in original does not affect copy",
			yaml: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  annotations:
    a: b
`,
			mutate: func(orig, copy []Manifest) {
				orig[0].AddAnnotations(map[string]string{"a": "c"})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			orig := mustParseManifests(t, tt.yaml)
			copy := DeepCopyManifests(orig)

			require.Equal(t, orig, copy, "copy should be equal to original before mutation")

			// Mutate as per test case
			tt.mutate(orig, copy)

			// After mutation, the original and copy should differ
			assert.NotEqual(t, orig, copy, "mutation should not affect the other slice")
		})
	}
}

func TestManifest_DeepCopyWithName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		yaml      string
		newName   string
		mutate    func(orig, copy *Manifest)
		checkOrig func(orig Manifest)
		checkCopy func(copy Manifest)
	}{
		{
			name: "deep copy with new name does not affect original",
			yaml: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: original-name
  labels:
    foo: bar
`,
			newName: "copied-name",
			mutate:  nil,
			checkOrig: func(orig Manifest) {
				assert.Equal(t, "original-name", orig.Name())
			},
			checkCopy: func(copy Manifest) {
				assert.Equal(t, "copied-name", copy.Name())
			},
		},
		{
			name: "mutate copy does not affect original",
			yaml: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: original-name
  labels:
    foo: bar
`,
			newName: "copied-name",
			mutate: func(orig, copy *Manifest) {
				copy.AddLabels(map[string]string{"foo": "baz"})
			},
			checkOrig: func(orig Manifest) {
				assert.Equal(t, "bar", orig.body.GetLabels()["foo"], "original label should remain unchanged")
			},
			checkCopy: func(copy Manifest) {
				assert.Equal(t, "baz", copy.body.GetLabels()["foo"], "copy label should be updated")
			},
		},
		{
			name: "mutate original does not affect copy",
			yaml: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: original-name
  labels:
    foo: bar
`,
			newName: "copied-name",
			mutate: func(orig, copy *Manifest) {
				orig.AddLabels(map[string]string{"foo": "baz"})
			},
			checkOrig: func(orig Manifest) {
				assert.Equal(t, "baz", orig.body.GetLabels()["foo"], "original label should be updated")
			},
			checkCopy: func(copy Manifest) {
				assert.Equal(t, "bar", copy.body.GetLabels()["foo"], "copy label should remain unchanged")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			manifests := mustParseManifests(t, tt.yaml)
			require.Len(t, manifests, 1)
			orig := manifests[0]
			copy := orig.DeepCopyWithName(tt.newName)

			if tt.mutate != nil {
				tt.mutate(&orig, &copy)
			}

			tt.checkOrig(orig)
			tt.checkCopy(copy)
		})
	}
}

func TestFromStructuredObject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		input          any
		wantKind       string
		wantAPIVersion string
		wantName       string
		wantErr        bool
		wantData       map[string]any
	}{
		{
			name: "ConfigMap",
			input: &corev1.ConfigMap{
				TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "ConfigMap"},
				ObjectMeta: metav1.ObjectMeta{Name: "test-config", Namespace: "default"},
				Data:       map[string]string{"key": "value"},
			},
			wantKind:       "ConfigMap",
			wantAPIVersion: "v1",
			wantName:       "test-config",
			wantErr:        false,
			wantData:       map[string]any{"key": "value"},
		},
		{
			name: "Secret",
			input: &corev1.Secret{
				TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Secret"},
				ObjectMeta: metav1.ObjectMeta{Name: "test-secret", Namespace: "default"},
				Data:       map[string][]byte{"password": []byte("password")},
			},
			wantKind:       "Secret",
			wantAPIVersion: "v1",
			wantName:       "test-secret",
			wantErr:        false,
			wantData:       map[string]any{"password": "cGFzc3dvcmQ="},
		},
		{
			name:           "invalid object",
			input:          struct{ Foo string }{Foo: "bar"},
			wantKind:       "",
			wantAPIVersion: "",
			wantName:       "",
			wantErr:        true,
			wantData:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m, err := FromStructuredObject(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantKind, m.Kind())
			assert.Equal(t, tt.wantAPIVersion, m.APIVersion())
			assert.Equal(t, tt.wantName, m.Name())
			assert.Equal(t, tt.wantData, m.body.Object["data"])
		})
	}
}
