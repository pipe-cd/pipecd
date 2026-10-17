// Copyright 2026 The PipeCD Authors.
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
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/config"
)

func TestLoadServiceDefinitionCommitHashTag(t *testing.T) {
	t.Parallel()

	serviceDef := []byte(`
cluster: arn:aws:ecs:ap-northeast-1:XXXX:cluster/YYYY
serviceName: nginx-service
desiredCount: 2
deploymentController:
  type: EXTERNAL
`)
	runningDir := t.TempDir()
	targetDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(runningDir, "servicedef.yaml"), serviceDef, 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(targetDir, "servicedef.yaml"), serviceDef, 0o600))

	input := &sdk.ExecuteStageInput[config.ECSApplicationSpec]{
		Request: sdk.ExecuteStageRequest[config.ECSApplicationSpec]{
			Deployment: sdk.Deployment{
				PipedID:       "piped-id",
				ApplicationID: "app-id",
			},
			RunningDeploymentSource: sdk.DeploymentSource[config.ECSApplicationSpec]{
				ApplicationDirectory: runningDir,
				CommitHash:           "running-commit",
			},
			TargetDeploymentSource: sdk.DeploymentSource[config.ECSApplicationSpec]{
				ApplicationDirectory: targetDir,
				CommitHash:           "target-commit",
			},
		},
	}

	testcases := []struct {
		name       string
		ds         sdk.DeploymentSource[config.ECSApplicationSpec]
		wantCommit string
	}{
		{
			name:       "target source is tagged with the target commit",
			ds:         input.Request.TargetDeploymentSource,
			wantCommit: "target-commit",
		},
		{
			// The rollback stage re-applies the running source. The live state sync check
			// compares this tag with the latest commit, so it must name the running commit.
			name:       "running source is tagged with the running commit",
			ds:         input.Request.RunningDeploymentSource,
			wantCommit: "running-commit",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			service, err := LoadServiceDefinition(tc.ds, "servicedef.yaml", input)
			require.NoError(t, err)

			tags := make(map[string]string, len(service.Tags))
			for _, tag := range service.Tags {
				tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
			}
			assert.Equal(t, tc.wantCommit, tags[LabelCommitHash])
			assert.Equal(t, ManagedByECSPlugin, tags[LabelManagedBy])
			assert.Equal(t, "piped-id", tags[LabelPiped])
			assert.Equal(t, "app-id", tags[LabelApplication])
		})
	}
}
