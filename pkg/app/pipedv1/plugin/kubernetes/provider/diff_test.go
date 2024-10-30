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
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/plugin/diff"
)

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
