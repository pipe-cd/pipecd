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

func loadManifests(appDir, taskDefFile, serviceDefFile string) (ECSManifests, error) {
	taskDef, err := LoadTaskDefinition(appDir, taskDefFile)
	if err != nil {
		return ECSManifests{}, err
	}
	serviceDef, err := LoadServiceDefinition(appDir, serviceDefFile)
	if err != nil {
		return ECSManifests{}, err
	}

	return ECSManifests{
		TaskDefinition:    &taskDef,
		ServiceDefinition: &serviceDef,
	}, nil
}

func TestDiff(t *testing.T) {
	t.Parallel()

	old, err := loadManifests("testdata/", "old_taskdef.yaml", "old_servicedef.yaml")
	require.NoError(t, err)

	new, err := loadManifests("testdata/", "new_taskdef.yaml", "new_servicedef.yaml")
	require.NoError(t, err)

	// Expect to have change.
	diff, err := Diff(old, new)
	require.NoError(t, err)
	noChange := diff.NoChange()
	require.False(t, noChange)

	// Expect no change.
	diff, err = Diff(old, old)
	require.NoError(t, err)
	noChange = diff.NoChange()
	require.True(t, noChange)
}

func TestDiffResult_Render(t *testing.T) {
	old, err := loadManifests("testdata/", "old_taskdef.yaml", "old_servicedef.yaml")
	require.NoError(t, err)

	new, err := loadManifests("testdata/", "new_taskdef.yaml", "new_servicedef.yaml")
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
	expected = `# 1. ServiceDefinition
@@ -10,7 +10,7 @@
 DeploymentController:
   Type: EXTERNAL
 Deployments: null
-DesiredCount: 2
+DesiredCount: 3
 EnableECSManagedTags: true
 EnableExecuteCommand: false
 Events: null

# 2. TaskDefinition
@@ -17,7 +17,7 @@
   FirelensConfiguration: null
   HealthCheck: null
   Hostname: null
-  Image: XXXX.dkr.ecr.ap-northeast-1.amazonaws.com/nginx:1
+  Image: XXXX.dkr.ecr.ap-northeast-1.amazonaws.com/nginx:2
   Interactive: null
   Links: null
   LinuxParameters: null
`
	require.Equal(t, expected, actual)
}
