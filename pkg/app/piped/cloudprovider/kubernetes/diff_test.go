// Copyright 2021 The PipeCD Authors.
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
)

func TestGroupManifests(t *testing.T) {
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
	testcases := []struct {
		name      string
		manifests string
		expected  string
	}{
		{
			name:      "no diff",
			manifests: "testdata/diff_by_command_no_change.yaml",
			expected:  "",
		},
		{
			name:      "has diff",
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
         name: second
`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			manifests, err := LoadManifestsFromYAMLFile(tc.manifests)
			require.NoError(t, err)
			require.Equal(t, 2, len(manifests))

			got, err := diffByCommand(manifests[0], manifests[1])
			require.NoError(t, err)
			assert.Equal(t, tc.expected, string(got))
		})
	}
}
