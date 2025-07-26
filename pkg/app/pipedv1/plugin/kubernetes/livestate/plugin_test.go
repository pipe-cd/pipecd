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

package livestate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/pipe-cd/piped-plugin-sdk-go/diff"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
)

func makeTestManifest(t *testing.T, yaml string) provider.Manifest {
	t.Helper()
	manifests, err := provider.ParseManifests(yaml)
	require.NoError(t, err)
	require.Len(t, manifests, 1)
	return manifests[0]
}

func makeTestDiffChange(t *testing.T, oldYAML, newYAML string) provider.DiffListChange {
	t.Helper()
	old := makeTestManifest(t, oldYAML)
	new := makeTestManifest(t, newYAML)

	oldData, err := old.MarshalJSON()
	require.NoError(t, err)
	oldUnstructured := unstructured.Unstructured{}
	err = oldUnstructured.UnmarshalJSON(oldData)
	require.NoError(t, err)

	newData, err := new.MarshalJSON()
	require.NoError(t, err)
	newUnstructured := unstructured.Unstructured{}
	err = newUnstructured.UnmarshalJSON(newData)
	require.NoError(t, err)

	result, err := diff.DiffUnstructureds(oldUnstructured, newUnstructured, old.Key().String(),
		diff.WithEquateEmpty(),
		diff.WithIgnoreAddingMapKeys(),
		diff.WithCompareNumberAndNumericString(),
	)
	require.NoError(t, err)

	return provider.DiffListChange{
		Old:  old,
		New:  new,
		Diff: result,
	}
}

func TestCalculateSyncState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		diffResult *provider.DiffListResult
		commitHash string
		want       sdk.ApplicationSyncState
	}{
		{
			name: "all resources are in sync",
			diffResult: &provider.DiffListResult{
				Adds:    []provider.Manifest{},
				Deletes: []provider.Manifest{},
				Changes: []provider.DiffListChange{},
			},
			commitHash: "1234567",
			want: sdk.ApplicationSyncState{
				Status:      sdk.ApplicationSyncStateSynced,
				ShortReason: "",
				Reason:      "",
			},
		},
		{
			name: "changed one resource",
			diffResult: &provider.DiffListResult{
				Changes: []provider.DiffListChange{
					makeTestDiffChange(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
      - name: app
        image: nginx:1.19
`, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment
  namespace: default
spec:
  replicas: 3
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
      - name: app
        image: nginx:1.20
`),
				},
			},
			commitHash: "1234567",
			want: sdk.ApplicationSyncState{
				Status:      sdk.ApplicationSyncStateOutOfSync,
				ShortReason: "There are 1 manifests not synced (0 adds, 0 deletes, 1 changes)",
				Reason: `Diff between the defined state in Git at commit 1234567 and actual state in cluster:

--- Actual   (LiveState)
+++ Expected (Git)

# 1. name="test-deployment", kind="Deployment", namespace="default", apiGroup="apps"

  spec:
    #spec.replicas
-   replicas: 1
+   replicas: 3

    template:
      spec:
        containers:
          -
            #spec.template.spec.containers.0.image
-           image: nginx:1.19
+           image: nginx:1.20


`,
			},
		},
		{
			name: "changed two resources",
			diffResult: &provider.DiffListResult{
				Changes: []provider.DiffListChange{
					makeTestDiffChange(t, `
apiVersion: v1
kind: Service
metadata:
  name: test-service
  namespace: default
spec:
  ports:
  - port: 80
    targetPort: 8080
  selector:
    app: test
`, `
apiVersion: v1
kind: Service
metadata:
  name: test-service
  namespace: default
spec:
  ports:
  - port: 443
    targetPort: 8443
  selector:
    app: test
`),
					makeTestDiffChange(t, `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: test-ingress
  namespace: default
spec:
  rules:
  - host: old.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: test-service
            port:
              number: 80
`, `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: test-ingress
  namespace: default
spec:
  rules:
  - host: new.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: test-service
            port:
              number: 443
`),
				},
			},
			commitHash: "1234567",
			want: sdk.ApplicationSyncState{
				Status:      sdk.ApplicationSyncStateOutOfSync,
				ShortReason: "There are 2 manifests not synced (0 adds, 0 deletes, 2 changes)",
				Reason: `Diff between the defined state in Git at commit 1234567 and actual state in cluster:

--- Actual   (LiveState)
+++ Expected (Git)

# 1. name="test-service", kind="Service", namespace="default", apiGroup=""

  spec:
    ports:
      -
        #spec.ports.0.port
-       port: 80
+       port: 443

        #spec.ports.0.targetPort
-       targetPort: 8080
+       targetPort: 8443


# 2. name="test-ingress", kind="Ingress", namespace="default", apiGroup="networking.k8s.io"

  spec:
    rules:
      -
        #spec.rules.0.host
-       host: old.example.com
+       host: new.example.com

        http:
          paths:
            - backend:
                service:
                  port:
                    #spec.rules.0.http.paths.0.backend.service.port.number
-                   number: 80
+                   number: 443


`,
			},
		},
		{
			name: "resource deletion and addition",
			diffResult: &provider.DiffListResult{
				Adds: []provider.Manifest{
					makeTestManifest(t, `
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: test-pvc
  namespace: default
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
`),
				},
				Deletes: []provider.Manifest{
					makeTestManifest(t, `
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
  namespace: default
spec:
  containers:
  - name: app
    image: nginx:1.19
`),
				},
			},
			commitHash: "1234567",
			want: sdk.ApplicationSyncState{
				Status:      sdk.ApplicationSyncStateOutOfSync,
				ShortReason: "There are 2 manifests not synced (1 adds, 1 deletes, 0 changes)",
				Reason: `Diff between the defined state in Git at commit 1234567 and actual state in cluster:

--- Actual   (LiveState)
+++ Expected (Git)

- 1. name="test-pod", kind="Pod", namespace="default", apiGroup=""

+ 2. name="test-pvc", kind="PersistentVolumeClaim", namespace="default", apiGroup=""

`,
			},
		},
		{
			name: "config map data is masked",
			diffResult: &provider.DiffListResult{
				Changes: []provider.DiffListChange{
					makeTestDiffChange(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  namespace: default
data:
  key: old-value
`, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  namespace: default
data:
  key: new-value
`),
				},
			},
			commitHash: "1234567",
			want: sdk.ApplicationSyncState{
				Status:      sdk.ApplicationSyncStateOutOfSync,
				ShortReason: "There are 1 manifests not synced (0 adds, 0 deletes, 1 changes)",
				Reason: `Diff between the defined state in Git at commit 1234567 and actual state in cluster:

--- Actual   (LiveState)
+++ Expected (Git)

# 1. name="test-config", kind="ConfigMap", namespace="default", apiGroup=""

  data:
    #data.key
-   key: *****
+   key: *****


`,
			},
		},
		{
			name: "secret data is masked",
			diffResult: &provider.DiffListResult{
				Changes: []provider.DiffListChange{
					makeTestDiffChange(t, `
apiVersion: v1
kind: Secret
metadata:
  name: test-secret
  namespace: default
data:
  username: YWRtaW4=
  password: c2VjcmV0
`, `
apiVersion: v1
kind: Secret
metadata:
  name: test-secret
  namespace: default
data:
  username: dXNlcg==
  password: bmV3c2VjcmV0
`),
				},
			},
			commitHash: "1234567",
			want: sdk.ApplicationSyncState{
				Status:      sdk.ApplicationSyncStateOutOfSync,
				ShortReason: "There are 1 manifests not synced (0 adds, 0 deletes, 1 changes)",
				Reason: `Diff between the defined state in Git at commit 1234567 and actual state in cluster:

--- Actual   (LiveState)
+++ Expected (Git)

# 1. name="test-secret", kind="Secret", namespace="default", apiGroup=""

  data:
    #data.password
-   password: *****
+   password: *****

    #data.username
-   username: *****
+   username: *****


`,
			},
		},
		{
			name: "maximum three changes are shown",
			diffResult: &provider.DiffListResult{
				Changes: []provider.DiffListChange{
					makeTestDiffChange(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: config1
  namespace: default
data:
  key: value1
`, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: config1
  namespace: default
data:
  key: new-value1
`),
					makeTestDiffChange(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: config2
  namespace: default
data:
  key: value2
`, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: config2
  namespace: default
data:
  key: new-value2
`),
					makeTestDiffChange(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: config3
  namespace: default
data:
  key: value3
`, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: config3
  namespace: default
data:
  key: new-value3
`),
					makeTestDiffChange(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: config4
  namespace: default
data:
  key: value4
`, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: config4
  namespace: default
data:
  key: new-value4
`),
				},
			},
			commitHash: "1234567",
			want: sdk.ApplicationSyncState{
				Status:      sdk.ApplicationSyncStateOutOfSync,
				ShortReason: "There are 4 manifests not synced (0 adds, 0 deletes, 4 changes)",
				Reason: `Diff between the defined state in Git at commit 1234567 and actual state in cluster:

--- Actual   (LiveState)
+++ Expected (Git)

# 1. name="config1", kind="ConfigMap", namespace="default", apiGroup=""

  data:
    #data.key
-   key: *****
+   key: *****


# 2. name="config2", kind="ConfigMap", namespace="default", apiGroup=""

  data:
    #data.key
-   key: *****
+   key: *****


# 3. name="config3", kind="ConfigMap", namespace="default", apiGroup=""

  data:
    #data.key
-   key: *****
+   key: *****


... (omitted 1 other changed manifests)
`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := calculateSyncState(tt.diffResult, tt.commitHash)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFilterIgnoringManifests(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		manifests []provider.Manifest
		want      []provider.Manifest
	}{
		{
			name:      "empty slice",
			manifests: []provider.Manifest{},
			want:      []provider.Manifest{},
		},
		{
			name: "manifest without ignore annotation",
			manifests: []provider.Manifest{
				makeTestManifest(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  namespace: default
data:
  key: value
`),
			},
			want: []provider.Manifest{
				makeTestManifest(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  namespace: default
data:
  key: value
`),
			},
		},
		{
			name: "manifest with ignore annotation set to true - should be filtered out",
			manifests: []provider.Manifest{
				makeTestManifest(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  namespace: default
  annotations:
    pipecd.dev/ignore-drift-detection: "true"
data:
  key: value
`),
			},
			want: []provider.Manifest{},
		},
		{
			name: "manifest with ignore annotation set to false - should be included",
			manifests: []provider.Manifest{
				makeTestManifest(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  namespace: default
  annotations:
    pipecd.dev/ignore-drift-detection: "false"
data:
  key: value
`),
			},
			want: []provider.Manifest{
				makeTestManifest(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  namespace: default
  annotations:
    pipecd.dev/ignore-drift-detection: "false"
data:
  key: value
`),
			},
		},
		{
			name: "manifest with ignore annotation set to other value - should be included",
			manifests: []provider.Manifest{
				makeTestManifest(t, `
apiVersion: v1
kind: Service
metadata:
  name: test-service
  namespace: default
  annotations:
    pipecd.dev/ignore-drift-detection: "some-other-value"
spec:
  selector:
    app: test
  ports:
  - port: 80
`),
			},
			want: []provider.Manifest{
				makeTestManifest(t, `
apiVersion: v1
kind: Service
metadata:
  name: test-service
  namespace: default
  annotations:
    pipecd.dev/ignore-drift-detection: "some-other-value"
spec:
  selector:
    app: test
  ports:
  - port: 80
`),
			},
		},
		{
			name: "mixed manifests - some filtered, some included",
			manifests: []provider.Manifest{
				makeTestManifest(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: include-config
  namespace: default
data:
  key: value
`),
				makeTestManifest(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: ignore-config
  namespace: default
  annotations:
    pipecd.dev/ignore-drift-detection: "true"
data:
  key: value
`),
				makeTestManifest(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment
  namespace: default
  annotations:
    pipecd.dev/ignore-drift-detection: "false"
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
      - name: app
        image: nginx:1.19
`),
				makeTestManifest(t, `
apiVersion: v1
kind: Secret
metadata:
  name: ignore-secret
  namespace: default
  annotations:
    pipecd.dev/ignore-drift-detection: "true"
data:
  username: dXNlcg==
`),
			},
			want: []provider.Manifest{
				makeTestManifest(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: include-config
  namespace: default
data:
  key: value
`),
				makeTestManifest(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment
  namespace: default
  annotations:
    pipecd.dev/ignore-drift-detection: "false"
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
      - name: app
        image: nginx:1.19
`),
			},
		},
		{
			name: "all manifests should be filtered out",
			manifests: []provider.Manifest{
				makeTestManifest(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: ignore-config1
  namespace: default
  annotations:
    pipecd.dev/ignore-drift-detection: "true"
data:
  key: value1
`),
				makeTestManifest(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: ignore-config2
  namespace: default
  annotations:
    pipecd.dev/ignore-drift-detection: "true"
data:
  key: value2
`),
			},
			want: []provider.Manifest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := filterIgnoringManifests(tt.manifests)
			assert.Equal(t, len(tt.want), len(got), "expected number of manifests should match")

			// Compare each manifest by their key and content
			for i, wantManifest := range tt.want {
				require.Less(t, i, len(got), "got slice should have enough elements")

				assert.Equal(t, wantManifest.Key().String(), got[i].Key().String(),
					"manifest keys should match at index %d", i)

				wantJSON, err := wantManifest.MarshalJSON()
				require.NoError(t, err)
				gotJSON, err := got[i].MarshalJSON()
				require.NoError(t, err)
				assert.JSONEq(t, string(wantJSON), string(gotJSON),
					"manifest content should match at index %d", i)
			}
		})
	}
}

func TestOnlyWatchingResourceKinds(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                  string
		manifests             []provider.Manifest
		watchingResourceKinds []schema.GroupVersionKind
		want                  []provider.Manifest
	}{
		{
			name:                  "empty manifests",
			manifests:             []provider.Manifest{},
			watchingResourceKinds: []schema.GroupVersionKind{{Group: "apps", Version: "v1", Kind: "Deployment"}},
			want:                  []provider.Manifest{},
		},
		{
			name: "empty watching resource kinds",
			manifests: []provider.Manifest{
				makeTestManifest(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  namespace: default
data:
  key: value
`),
			},
			watchingResourceKinds: []schema.GroupVersionKind{},
			want:                  []provider.Manifest{},
		},
		{
			name: "no matching resource kinds",
			manifests: []provider.Manifest{
				makeTestManifest(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  namespace: default
data:
  key: value
`),
				makeTestManifest(t, `
apiVersion: v1
kind: Secret
metadata:
  name: test-secret
  namespace: default
data:
  username: dXNlcg==
`),
			},
			watchingResourceKinds: []schema.GroupVersionKind{
				{Group: "apps", Version: "v1", Kind: "Deployment"},
				{Group: "networking.k8s.io", Version: "v1", Kind: "Ingress"},
			},
			want: []provider.Manifest{},
		},
		{
			name: "some matching resource kinds",
			manifests: []provider.Manifest{
				makeTestManifest(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  namespace: default
data:
  key: value
`),
				makeTestManifest(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
      - name: app
        image: nginx:1.19
`),
				makeTestManifest(t, `
apiVersion: v1
kind: Secret
metadata:
  name: test-secret
  namespace: default
data:
  username: dXNlcg==
`),
			},
			watchingResourceKinds: []schema.GroupVersionKind{
				{Group: "", Version: "v1", Kind: "ConfigMap"},
				{Group: "apps", Version: "v1", Kind: "Deployment"},
			},
			want: []provider.Manifest{
				makeTestManifest(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  namespace: default
data:
  key: value
`),
				makeTestManifest(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
      - name: app
        image: nginx:1.19
`),
			},
		},
		{
			name: "all manifests match watching resource kinds",
			manifests: []provider.Manifest{
				makeTestManifest(t, `
apiVersion: v1
kind: Service
metadata:
  name: test-service
  namespace: default
spec:
  selector:
    app: test
  ports:
  - port: 80
`),
				makeTestManifest(t, `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: test-ingress
  namespace: default
spec:
  rules:
  - host: example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: test-service
            port:
              number: 80
`),
			},
			watchingResourceKinds: []schema.GroupVersionKind{
				{Group: "", Version: "v1", Kind: "Service"},
				{Group: "networking.k8s.io", Version: "v1", Kind: "Ingress"},
			},
			want: []provider.Manifest{
				makeTestManifest(t, `
apiVersion: v1
kind: Service
metadata:
  name: test-service
  namespace: default
spec:
  selector:
    app: test
  ports:
  - port: 80
`),
				makeTestManifest(t, `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: test-ingress
  namespace: default
spec:
  rules:
  - host: example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: test-service
            port:
              number: 80
`),
			},
		},
		{
			name: "duplicate resource kinds in watching list",
			manifests: []provider.Manifest{
				makeTestManifest(t, `
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
  namespace: default
spec:
  containers:
  - name: app
    image: nginx:1.19
`),
			},
			watchingResourceKinds: []schema.GroupVersionKind{
				{Group: "", Version: "v1", Kind: "Pod"},
				{Group: "", Version: "v1", Kind: "Pod"}, // duplicate
				{Group: "apps", Version: "v1", Kind: "Deployment"},
			},
			want: []provider.Manifest{
				makeTestManifest(t, `
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
  namespace: default
spec:
  containers:
  - name: app
    image: nginx:1.19
`),
			},
		},
		{
			name: "mixed resource kinds with some matches",
			manifests: []provider.Manifest{
				makeTestManifest(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: config1
  namespace: default
data:
  key: value1
`),
				makeTestManifest(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deploy1
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
      - name: app
        image: nginx:1.19
`),
				makeTestManifest(t, `
apiVersion: v1
kind: Secret
metadata:
  name: secret1
  namespace: default
data:
  username: dXNlcg==
`),
				makeTestManifest(t, `
apiVersion: v1
kind: Service
metadata:
  name: service1
  namespace: default
spec:
  selector:
    app: test
  ports:
  - port: 80
`),
			},
			watchingResourceKinds: []schema.GroupVersionKind{
				{Group: "", Version: "v1", Kind: "ConfigMap"},
				{Group: "", Version: "v1", Kind: "Service"},
			},
			want: []provider.Manifest{
				makeTestManifest(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: config1
  namespace: default
data:
  key: value1
`),
				makeTestManifest(t, `
apiVersion: v1
kind: Service
metadata:
  name: service1
  namespace: default
spec:
  selector:
    app: test
  ports:
  - port: 80
`),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := onlyWatchingResourceKinds(tt.manifests, tt.watchingResourceKinds)
			assert.Equal(t, len(tt.want), len(got), "expected number of manifests should match")

			// Compare each manifest by their key and content
			for i, wantManifest := range tt.want {
				require.Less(t, i, len(got), "got slice should have enough elements")

				assert.Equal(t, wantManifest.Key().String(), got[i].Key().String(),
					"manifest keys should match at index %d", i)

				wantJSON, err := wantManifest.MarshalJSON()
				require.NoError(t, err)
				gotJSON, err := got[i].MarshalJSON()
				require.NoError(t, err)
				assert.JSONEq(t, string(wantJSON), string(gotJSON),
					"manifest content should match at index %d", i)
			}
		})
	}
}
