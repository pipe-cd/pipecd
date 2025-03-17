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

func TestDiffList(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		live     string
		desired  string
		wantAdds int
		wantDels int
		wantMods int
	}{
		{
			name: "no changes",
			live: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment
spec:
  replicas: 3`,
			desired: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment
spec:
  replicas: 3`,
			wantAdds: 0,
			wantDels: 0,
			wantMods: 0,
		},
		{
			name: "one addition",
			live: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment-1
spec:
  replicas: 3`,
			desired: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment-1
spec:
  replicas: 3
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment-2
spec:
  replicas: 3`,
			wantAdds: 1,
			wantDels: 0,
			wantMods: 0,
		},
		{
			name: "one deletion",
			live: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment-1
spec:
  replicas: 3
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment-2
spec:
  replicas: 3`,
			desired: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment-1
spec:
  replicas: 3`,
			wantAdds: 0,
			wantDels: 1,
			wantMods: 0,
		},
		{
			name: "one modification",
			live: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment
spec:
  replicas: 3`,
			desired: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment
spec:
  replicas: 5`,
			wantAdds: 0,
			wantDels: 0,
			wantMods: 1,
		},
		{
			name: "mixed changes",
			live: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment-1
spec:
  replicas: 3
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment-2
spec:
  replicas: 3`,
			desired: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment-1
spec:
  replicas: 5
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment-3
spec:
  replicas: 3`,
			wantAdds: 1,
			wantDels: 1,
			wantMods: 1,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			liveManifests, err := ParseManifests(tc.live)
			require.NoError(t, err)
			desiredManifests, err := ParseManifests(tc.desired)
			require.NoError(t, err)

			result, err := DiffList(liveManifests, desiredManifests, zap.NewNop(), diff.WithEquateEmpty(), diff.WithIgnoreAddingMapKeys(), diff.WithCompareNumberAndNumericString())
			require.NoError(t, err)

			assert.Equal(t, tc.wantAdds, len(result.Adds))
			assert.Equal(t, tc.wantDels, len(result.Deletes))
			assert.Equal(t, tc.wantMods, len(result.Changes))
			assert.Equal(t, tc.wantAdds+tc.wantDels+tc.wantMods == 0, result.NoChanges())
			assert.Equal(t, tc.wantAdds+tc.wantDels+tc.wantMods, result.TotalOutOfSync())
		})
	}
}

func TestGroupManifests(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		olds        string
		news        string
		wantAdds    int
		wantDeletes int
		wantChanges int
	}{
		{
			name:        "empty lists",
			olds:        "",
			news:        "",
			wantAdds:    0,
			wantDeletes: 0,
			wantChanges: 0,
		},
		{
			name: "only additions",
			olds: "",
			news: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment-1
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment-2`,
			wantAdds:    2,
			wantDeletes: 0,
			wantChanges: 0,
		},
		{
			name: "only deletions",
			olds: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment-1
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment-2`,
			news:        "",
			wantAdds:    0,
			wantDeletes: 2,
			wantChanges: 0,
		},
		{
			name: "only changes",
			olds: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment-1
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment-2`,
			news: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment-1
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment-2`,
			wantAdds:    0,
			wantDeletes: 0,
			wantChanges: 2,
		},
		{
			name: "mixed changes",
			olds: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment-1
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment-2`,
			news: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment-2
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment-3`,
			wantAdds:    1,
			wantDeletes: 1,
			wantChanges: 1,
		},
		{
			name: "different resource types with same name",
			olds: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-1
---
apiVersion: v1
kind: Service
metadata:
  name: test-1`,
			news: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-1
---
apiVersion: v1
kind: Service
metadata:
  name: test-1`,
			wantAdds:    0,
			wantDeletes: 0,
			wantChanges: 2,
		},
		{
			name: "different namespaces with same name",
			olds: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-1
  namespace: ns1
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-1
  namespace: ns2`,
			news: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-1
  namespace: ns1
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-1
  namespace: ns2`,
			wantAdds:    0,
			wantDeletes: 0,
			wantChanges: 2,
		},
		{
			name: "old list larger than new list",
			olds: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-1
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-2
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-3`,
			news: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-2`,
			wantAdds:    0,
			wantDeletes: 2,
			wantChanges: 1,
		},
		{
			name: "new list larger than old list",
			olds: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-2`,
			news: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-1
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-2
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-3`,
			wantAdds:    2,
			wantDeletes: 0,
			wantChanges: 1,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var olds, news []Manifest
			var err error

			if tc.olds != "" {
				olds, err = ParseManifests(tc.olds)
				require.NoError(t, err)
			}

			if tc.news != "" {
				news, err = ParseManifests(tc.news)
				require.NoError(t, err)
			}

			adds, deletes, newChanges, oldChanges := groupManifests(olds, news)
			assert.Equal(t, tc.wantAdds, len(adds))
			assert.Equal(t, tc.wantDeletes, len(deletes))
			assert.Equal(t, tc.wantChanges, len(newChanges))
			assert.Equal(t, len(newChanges), len(oldChanges))
		})
	}
}
