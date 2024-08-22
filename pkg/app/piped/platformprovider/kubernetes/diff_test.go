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

package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/diff"
)

func TestGroupManifests(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name               string
		olds               []Manifest
		news               []Manifest
		expectedAdds       []Manifest
		expectedDeletes    []Manifest
		expectedNewChanges []Manifest
		expectedOldChanges []Manifest
	}{
		{
			name: "empty list",
		},
		{
			name: "only adds",
			news: []Manifest{
				{Key: ResourceKey{Name: "b"}},
				{Key: ResourceKey{Name: "a"}},
			},
			expectedAdds: []Manifest{
				{Key: ResourceKey{Name: "a"}},
				{Key: ResourceKey{Name: "b"}},
			},
		},
		{
			name: "only deletes",
			olds: []Manifest{
				{Key: ResourceKey{Name: "b"}},
				{Key: ResourceKey{Name: "a"}},
			},
			expectedDeletes: []Manifest{
				{Key: ResourceKey{Name: "a"}},
				{Key: ResourceKey{Name: "b"}},
			},
		},
		{
			name: "only inters",
			olds: []Manifest{
				{Key: ResourceKey{Name: "b"}},
				{Key: ResourceKey{Name: "a"}},
			},
			news: []Manifest{
				{Key: ResourceKey{Name: "a"}},
				{Key: ResourceKey{Name: "b"}},
			},
			expectedNewChanges: []Manifest{
				{Key: ResourceKey{Name: "a"}},
				{Key: ResourceKey{Name: "b"}},
			},
			expectedOldChanges: []Manifest{
				{Key: ResourceKey{Name: "a"}},
				{Key: ResourceKey{Name: "b"}},
			},
		},
		{
			name: "all kinds",
			olds: []Manifest{
				{Key: ResourceKey{Name: "b"}},
				{Key: ResourceKey{Name: "a"}},
				{Key: ResourceKey{Name: "c"}},
			},
			news: []Manifest{
				{Key: ResourceKey{Name: "a"}},
				{Key: ResourceKey{Name: "d"}},
				{Key: ResourceKey{Name: "b"}},
			},
			expectedAdds: []Manifest{
				{Key: ResourceKey{Name: "d"}},
			},
			expectedDeletes: []Manifest{
				{Key: ResourceKey{Name: "c"}},
			},
			expectedNewChanges: []Manifest{
				{Key: ResourceKey{Name: "a"}},
				{Key: ResourceKey{Name: "b"}},
			},
			expectedOldChanges: []Manifest{
				{Key: ResourceKey{Name: "a"}},
				{Key: ResourceKey{Name: "b"}},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			adds, deletes, newChanges, oldChanges := groupManifests(tc.olds, tc.news)
			assert.Equal(t, tc.expectedAdds, adds)
			assert.Equal(t, tc.expectedDeletes, deletes)
			assert.Equal(t, tc.expectedNewChanges, newChanges)
			assert.Equal(t, tc.expectedOldChanges, oldChanges)
		})
	}
}

func TestDiffByCommand(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		command     string
		manifests   string
		expected    string
		expectedErr bool
	}{
		{
			name:        "no command",
			command:     "non-existent-diff",
			manifests:   "testdata/diff_by_command_no_change.yaml",
			expected:    "",
			expectedErr: true,
		},
		{
			name:      "no diff",
			command:   diffCommand,
			manifests: "testdata/diff_by_command_no_change.yaml",
			expected:  "",
		},
		{
			name:      "has diff",
			command:   diffCommand,
			manifests: "testdata/diff_by_command.yaml",
			expected: `@@ -6,7 +6,7 @@
     pipecd.dev/managed-by: piped
   name: simple
 spec:
-  replicas: 2
+  replicas: 3
   selector:
     matchLabels:
       app: simple
@@ -18,6 +18,7 @@
       containers:
       - args:
         - a
+        - d
         - b
         - c
         image: gcr.io/pipecd/first:v1.0.0
@@ -26,7 +27,6 @@
         - containerPort: 9085
       - args:
         - xx
-        - yy
         - zz
         image: gcr.io/pipecd/second:v1.0.0
         name: second`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			manifests, err := LoadManifestsFromYAMLFile(tc.manifests)
			require.NoError(t, err)
			require.Equal(t, 2, len(manifests))

			got, err := diffByCommand(tc.command, manifests[0], manifests[1])
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expected, string(got))
		})
	}
}

func TestDiff(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		manifests     string
		expected      string
		diffNum       int
		falsePositive bool
	}{
		{
			name: "Secret no diff 1",
			manifests: `apiVersion: apps/v1
kind: Secret
metadata:
  name: secret-management
---
apiVersion: apps/v1
kind: Secret
metadata:
  name: secret-management
`,
			expected: "",
			diffNum:  0,
		},
		{
			name: "Secret no diff 2",
			manifests: `apiVersion: apps/v1
kind: Secret
metadata:
  name: secret-management
data:
  password: hoge
stringData:
  foo: bar
---
apiVersion: apps/v1
kind: Secret
metadata:
  name: secret-management
data:
  password: hoge
stringData:
  foo: bar
`,
			expected: "",
			diffNum:  0,
		},
		{
			name: "Secret no diff with merge",
			manifests: `apiVersion: apps/v1
kind: Secret
metadata:
  name: secret-management
data:
  password: hoge
  foo: YmFy
---
apiVersion: apps/v1
kind: Secret
metadata:
  name: secret-management
data:
  password: hoge
stringData:
  foo: bar
`,
			expected: "",
			diffNum:  0,
		},
		{
			name: "Secret no diff override false-positive",
			manifests: `apiVersion: apps/v1
kind: Secret
metadata:
  name: secret-management
data:
  password: hoge
  foo: YmFy
---
apiVersion: apps/v1
kind: Secret
metadata:
  name: secret-management
data:
  password: hoge
  foo: Zm9v
stringData:
  foo: bar
`,
			expected:      "",
			diffNum:       0,
			falsePositive: true,
		},
		{
			name: "Secret has diff",
			manifests: `apiVersion: apps/v1
kind: Secret
metadata:
  name: secret-management
data:
  foo: YmFy
---
apiVersion: apps/v1
kind: Secret
metadata:
  name: secret-management
data:
  password: hoge
stringData:
  foo: bar
`,
			expected: `  #data
+ data:
+   password: hoge

`,
			diffNum: 1,
		},
		{
			name: "Pod no diff 1",
			manifests: `apiVersion: v1
kind: Pod
metadata:
  name: static-web
  labels:
    role: myrole
spec:
  containers:
    - name: web
      image: nginx
      resources:
        limits:
          memory: "2Gi"
---
apiVersion: v1
kind: Pod
metadata:
  name: static-web
  labels:
    role: myrole
spec:
  containers:
    - name: web
      image: nginx
      ports:
      resources:
        limits:
          memory: "2Gi"
`,
			expected:      "",
			diffNum:       0,
			falsePositive: false,
		},
		{
			name: "Pod no diff 2",
			manifests: `apiVersion: v1
kind: Pod
metadata:
  name: static-web
  labels:
    role: myrole
spec:
  containers:
    - name: web
      image: nginx
      resources:
        limits:
          memory: "1536Mi"
---
apiVersion: v1
kind: Pod
metadata:
  name: static-web
  labels:
    role: myrole
spec:
  containers:
    - name: web
      image: nginx
      ports:
      resources:
        limits:
          memory: "1.5Gi"
`,
			expected:      "",
			diffNum:       0,
			falsePositive: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			manifests, err := ParseManifests(tc.manifests)
			require.NoError(t, err)
			require.Equal(t, 2, len(manifests))
			old, new := manifests[0], manifests[1]

			result, err := Diff(old, new, zap.NewNop(), diff.WithEquateEmpty(), diff.WithIgnoreAddingMapKeys(), diff.WithCompareNumberAndNumericString())
			require.NoError(t, err)

			renderer := diff.NewRenderer(diff.WithLeftPadding(1))
			ds := renderer.Render(result.Nodes())
			if tc.falsePositive {
				assert.NotEqual(t, tc.diffNum, result.NumNodes())
				assert.NotEqual(t, tc.expected, ds)
			} else {
				assert.Equal(t, tc.diffNum, result.NumNodes())
				assert.Equal(t, tc.expected, ds)
			}
		})
	}
}

func TestLoadAndDiff(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name           string
		manifests      [2]string
		noChangeAssert assert.BoolAssertionFunc
	}{
		{
			name: "no diff",
			manifests: [2]string{
				`apiVersion: v1
kind: Service
metadata:
  name: simple
spec:
  ports:
  - port: 9085
    protocol: TCP
    targetPort: 9085
  selector:
    app: simple
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: simple
  name: simple
spec:
  replicas: 2
  selector:
    matchLabels:
      app: simple
      pipecd.dev/variant: primary
  template:
    metadata:
      annotations:
        sidecar.istio.io/inject: "false"
      labels:
        app: simple
        pipecd.dev/variant: primary
    spec:
      containers:
      - args:
        - server
        env:
        - name: JVM_MEM_MIN
          value: 750m
        - name: JVM_MEM_MAX
          value: 2000m
        image: ghcr.io/pipe-cd/helloworld:v0.32.0
        lifecycle:
          preStop:
            exec:
              command:
              - sh
              - -c
              - sleep 20
        name: helloworld
        ports:
        - containerPort: 9085
        resources:
          limits:
            cpu: "3"
            memory: 1.5Gi
          requests:
            cpu: 150m
            memory: 1Gi
`,
				`apiVersion: v1
kind: Service
metadata:
  name: simple
spec:
  ports:
  - port: 9085
    protocol: TCP
    targetPort: 9085
  selector:
    app: simple
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: simple
  name: simple
spec:
  replicas: 2
  selector:
    matchLabels:
      app: simple
      pipecd.dev/variant: primary
  template:
    metadata:
      annotations:
        sidecar.istio.io/inject: "false"
      labels:
        app: simple
        pipecd.dev/variant: primary
    spec:
      containers:
      - args:
        - server
        env:
        - name: JVM_MEM_MIN
          value: 750m
        - name: JVM_MEM_MAX
          value: 2000m
        image: ghcr.io/pipe-cd/helloworld:v0.32.0
        lifecycle:
          preStop:
            exec:
              command:
              - sh
              - -c
              - sleep 20
        name: helloworld
        ports:
        - containerPort: 9085
        resources:
          limits:
            cpu: "3"
            memory: 1.5Gi
          requests:
            cpu: 150m
            memory: 1Gi
`,
			},
			noChangeAssert: assert.True,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var manifests [2][]Manifest
			manifests[0], _ = ParseManifests(tc.manifests[0])
			manifests[1], _ = ParseManifests(tc.manifests[1])

			result, err := DiffList(manifests[0], manifests[1], zap.NewNop(), diff.WithEquateEmpty(), diff.WithIgnoreAddingMapKeys(), diff.WithCompareNumberAndNumericString())
			require.NoError(t, err)

			tc.noChangeAssert(t, result.NoChange())
		})
	}
}
