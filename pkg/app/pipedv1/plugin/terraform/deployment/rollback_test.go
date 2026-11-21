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

package deployment

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/pipe-cd/piped-plugin-sdk-go/logpersister/logpersistertest"
	"github.com/pipe-cd/piped-plugin-sdk-go/toolregistry/toolregistrytest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/config"
)

func TestPlugin_executeRollbackStage(t *testing.T) {
	t.Parallel()

	dts := []*sdk.DeployTarget[config.DeployTargetConfig]{{Name: "dt-1"}}

	t.Run("failure when stage log persister is unavailable", func(t *testing.T) {
		t.Parallel()

		p := &Plugin{}
		got := p.executeRollbackStage(context.Background(), &sdk.ExecuteStageInput[config.ApplicationConfigSpec]{
			Client: sdk.NewClient(nil, "terraform", "app-id", "", nil, toolregistrytest.NewTestToolRegistry(t)),
			Logger: zap.NewNop(),
		}, dts)
		assert.Equal(t, sdk.StageStatusFailure, got)
	})

	t.Run("failure when running deployment commit hash is empty", func(t *testing.T) {
		t.Parallel()

		p := &Plugin{}
		got := p.executeRollbackStage(context.Background(), &sdk.ExecuteStageInput[config.ApplicationConfigSpec]{
			Request: sdk.ExecuteStageRequest[config.ApplicationConfigSpec]{
				RunningDeploymentSource: sdk.DeploymentSource[config.ApplicationConfigSpec]{
					CommitHash: "",
				},
			},
			Client: sdk.NewClient(nil, "terraform", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), toolregistrytest.NewTestToolRegistry(t)),
			Logger: zap.NewNop(),
		}, dts)
		assert.Equal(t, sdk.StageStatusFailure, got)
	})

	t.Run("success rolling back to running deployment source", func(t *testing.T) {
		t.Parallel()

		appDir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(appDir, "main.tf"), []byte(`resource "terraform_data" "example" {
  input = "previous"
}
`), 0o644))

		p := &Plugin{}
		got := p.executeRollbackStage(context.Background(), &sdk.ExecuteStageInput[config.ApplicationConfigSpec]{
			Request: sdk.ExecuteStageRequest[config.ApplicationConfigSpec]{
				RunningDeploymentSource: sdk.DeploymentSource[config.ApplicationConfigSpec]{
					CommitHash:           "abc1234",
					ApplicationDirectory: appDir,
					ApplicationConfig: &sdk.ApplicationConfig[config.ApplicationConfigSpec]{
						Spec: &config.ApplicationConfigSpec{TerraformVersion: "1.12.0"},
					},
				},
			},
			Client: sdk.NewClient(nil, "terraform", "app-id", "stage-id", logpersistertest.NewTestLogPersister(t), toolregistrytest.NewTestToolRegistry(t)),
			Logger: zap.NewNop(),
		}, dts)
		require.Equal(t, sdk.StageStatusSuccess, got)
		state, err := os.ReadFile(filepath.Join(appDir, "terraform.tfstate"))
		require.NoError(t, err)
		assert.Contains(t, string(state), `"previous"`)
	})
}
