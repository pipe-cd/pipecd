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

package ecs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiff(t *testing.T) {
	t.Parallel()

	old, err := LoadECSManifest("testdata/", "old_taskdef.yaml", "old_servicedef.yaml")
	require.NoError(t, err)
	require.NotEmpty(t, old)

	new, err := LoadECSManifest("testdata/", "new_taskdef.yaml", "new_servicedef.yaml")
	require.NoError(t, err)
	require.NotEmpty(t, new)

	// Expect to have diff.
	got, err := Diff(old, new)
	require.NoError(t, err)
	require.NotEmpty(t, got)

	// Expect no diff.
	got, err = Diff(old, old)
	require.NoError(t, err)
	require.NotEmpty(t, got)
}

// func TestDiffResult_NoChange(t *testing.T) {
// 	t.Parallel()

// 	old, err := LoadECSManifest("./testdata/", "old_taskdef.yaml", "old_servicedef.yaml")
// 	require.NoError(t, err)
// 	require.NotEmpty(t, old)

// 	new, err := LoadECSManifest("./testdata/", "new_taskdef.yaml", "new_servicedef.yaml")
// 	require.NoError(t, err)
// 	require.NotEmpty(t, new)

// 	result, err := Diff(old, new)
// 	require.NoError(t, err)

// 	got := result.NoChange()
// 	require.False(t, got)
// }

// func TestDiffResult_Render(t *testing.T) {
// 	old, err := LoadECSManifest("./testdata/", "old_taskdef.yaml", "old_servicedef.yaml")
// 	require.NoError(t, err)

// 	new, err := LoadECSManifest("./testdata/", "new_taskdef.yaml", "new_servicedef.yaml")
// 	require.NoError(t, err)

// 	result, err := Diff(old, new)
// 	require.NoError(t, err)

// 	// Not use diff command
// 	opt := DiffRenderOptions{}
// 	got := result.Render(opt)
// 	want := `  spec:
//     template:
//       spec:
//         containers:
//           -
//             #spec.template.spec.containers.0.image
// -           image: gcr.io/pipecd/helloworld:v0.6.0
// +           image: gcr.io/pipecd/helloworld:v0.5.0

// `
// 	require.Equal(t, want, got)

// 	// Use diff command
// 	opt = DiffRenderOptions{UseDiffCommand: true}
// 	got = result.Render(opt)
// 	want = `@@ -18,7 +18,7 @@
//        containers:
//        - args:
//          - server
// -        image: gcr.io/pipecd/helloworld:v0.6.0
// +        image: gcr.io/pipecd/helloworld:v0.5.0
//          ports:
//          - containerPort: 9085
//            name: http1
// `
// 	require.Equal(t, want, got)
// }

// func TestDiffByCommand(t *testing.T) {
// 	t.Parallel()

// 	testcases := []struct {
// 		name        string
// 		command     string
// 		oldManifest string
// 		newManifest string
// 		expected    string
// 		expectedErr bool
// 	}{
// 		{
// 			name:        "no command",
// 			command:     "non-existent-diff",
// 			oldManifest: "testdata/old_manifest.yaml",
// 			newManifest: "testdata/old_manifest.yaml",
// 			expected:    "",
// 			expectedErr: true,
// 		},
// 		{
// 			name:        "no diff",
// 			command:     diffCommand,
// 			oldManifest: "testdata/old_manifest.yaml",
// 			newManifest: "testdata/old_manifest.yaml",
// 			expected:    "",
// 		},
// 		{
// 			name:        "has diff",
// 			command:     diffCommand,
// 			oldManifest: "testdata/old_manifest.yaml",
// 			newManifest: "testdata/new_manifest.yaml",
// 			expected: `@@ -18,7 +18,7 @@
//        containers:
//        - args:
//          - server
// -        image: gcr.io/pipecd/helloworld:v0.6.0
// +        image: gcr.io/pipecd/helloworld:v0.5.0
//          ports:
//          - containerPort: 9085
//            name: http1`,
// 		},
// 	}
// 	for _, tc := range testcases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			old, err := LoadECSManifest("./testdata/", "old_taskdef.yaml", "old_servicedef.yaml")
// 			require.NoError(t, err)

// 			new, err := LoadECSManifest("./testdata/", "new_taskdef.yaml", "new_servicedef.yaml")
// 			require.NoError(t, err)

// 			got, err := diffByCommand(tc.command, old, new)
// 			if tc.expectedErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 			assert.Equal(t, tc.expected, string(got))
// 		})
// 	}
// }
