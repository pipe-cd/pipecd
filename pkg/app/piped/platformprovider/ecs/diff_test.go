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

	"github.com/stretchr/testify/assert"
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
	diff, err := Diff(old, new)
	require.NoError(t, err)
	require.NotEmpty(t, diff)

	// Expect no diff for the same manifest.
	diff, err = Diff(old, old)
	require.NoError(t, err)
	require.NotEmpty(t, diff)
}

func TestDiffResult_NoChange(t *testing.T) {
	t.Parallel()

	old, err := LoadECSManifest("testdata/", "old_taskdef.yaml", "old_servicedef.yaml")
	require.NoError(t, err)
	require.NotEmpty(t, old)

	new, err := LoadECSManifest("testdata/", "new_taskdef.yaml", "new_servicedef.yaml")
	require.NoError(t, err)
	require.NotEmpty(t, new)

	diff, err := Diff(old, new)
	require.NoError(t, err)

	// Expect to have change.
	noChange := diff.NoChange()
	require.False(t, noChange)
}

// TODO: Convert each attribute in render result to lowerCamelCase
func TestDiffResult_Render(t *testing.T) {
	old, err := LoadECSManifest("testdata/", "old_taskdef.yaml", "old_servicedef.yaml")
	require.NoError(t, err)

	new, err := LoadECSManifest("testdata/", "new_taskdef.yaml", "new_servicedef.yaml")
	require.NoError(t, err)

	result, err := Diff(old, new)
	require.NoError(t, err)

	// Not using diff command
	opt := DiffRenderOptions{}
	actual := result.Render(opt)
	expected := `  ServiceDefinition:
    #ServiceDefinition.DesiredCount
-   DesiredCount: 2
+   DesiredCount: 3

  TaskDefinition:
    ContainerDefinitions:
      -
        #TaskDefinition.ContainerDefinitions.0.Image
-       Image: XXXX.dkr.ecr.ap-northeast-1.amazonaws.com/nginx:1
+       Image: XXXX.dkr.ecr.ap-northeast-1.amazonaws.com/nginx:2


`
	require.Equal(t, expected, actual)

	// Use diff command
	opt = DiffRenderOptions{UseDiffCommand: true}
	actual = result.Render(opt)
	expected = `@@ -17,7 +17,7 @@
   FirelensConfiguration: null
   HealthCheck: null
   Hostname: null
-  Image: XXXX.dkr.ecr.ap-northeast-1.amazonaws.com/nginx:1
+  Image: XXXX.dkr.ecr.ap-northeast-1.amazonaws.com/nginx:2
   Interactive: null
   Links: null
   LinuxParameters: null
@@ -10,7 +10,7 @@
 DeploymentController:
   Type: EXTERNAL
 Deployments: null
-DesiredCount: 2
+DesiredCount: 3
 EnableECSManagedTags: true
 EnableExecuteCommand: false
 Events: null
`
	require.Equal(t, expected, actual)
}

func TestDiffByCommand(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		command       string
		oldTaskDef    string
		oldServiceDef string
		newTaskDef    string
		newServiceDef string
		expected      string
		expectedErr   bool
	}{
		{
			name:          "no command",
			command:       "non-existent-diff",
			oldTaskDef:    "old_taskdef.yaml",
			oldServiceDef: "old_servicedef.yaml",
			newTaskDef:    "old_taskdef.yaml",
			newServiceDef: "old_servicedef.yaml",
			expected:      "",
			expectedErr:   true,
		},
		{
			name:          "no diff",
			command:       diffCommand,
			oldTaskDef:    "old_taskdef.yaml",
			oldServiceDef: "old_servicedef.yaml",
			newTaskDef:    "old_taskdef.yaml",
			newServiceDef: "old_servicedef.yaml",
			expected:      "\n",
		},
		{
			name:          "has diff",
			command:       diffCommand,
			oldTaskDef:    "old_taskdef.yaml",
			oldServiceDef: "old_servicedef.yaml",
			newTaskDef:    "new_taskdef.yaml",
			newServiceDef: "new_servicedef.yaml",
			expected: `@@ -17,7 +17,7 @@
   FirelensConfiguration: null
   HealthCheck: null
   Hostname: null
-  Image: XXXX.dkr.ecr.ap-northeast-1.amazonaws.com/nginx:1
+  Image: XXXX.dkr.ecr.ap-northeast-1.amazonaws.com/nginx:2
   Interactive: null
   Links: null
   LinuxParameters: null
@@ -10,7 +10,7 @@
 DeploymentController:
   Type: EXTERNAL
 Deployments: null
-DesiredCount: 2
+DesiredCount: 3
 EnableECSManagedTags: true
 EnableExecuteCommand: false
 Events: null`,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			old, err := LoadECSManifest("testdata/", tc.oldTaskDef, tc.oldServiceDef)
			require.NoError(t, err)

			new, err := LoadECSManifest("testdata/", tc.newTaskDef, tc.newServiceDef)
			require.NoError(t, err)

			got, err := diffByCommand(tc.command, old, new)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expected, string(got))
		})
	}
}
