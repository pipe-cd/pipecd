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

package planpreview

import (
	"context"
	"testing"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/config"
)

const (
	appConfigFile = "testdata/app.pipecd.yaml"
	pluginName    = "ecs"
	deployTarget  = "ecs-dev"
	runningDir    = "testdata/running"
	targetDir     = "testdata/target"
	runningCommit = "abc000"
	targetCommit  = "abc123"
)

func deployTargets(t *testing.T) []*sdk.DeployTarget[config.ECSDeployTargetConfig] {
	t.Helper()
	return []*sdk.DeployTarget[config.ECSDeployTargetConfig]{
		{Name: deployTarget},
	}
}

func makeInput(
	t *testing.T,
	targetAppDir string,
	runningAppDir string, // empty string means first deployment
) *sdk.GetPlanPreviewInput[config.ECSApplicationSpec] {
	t.Helper()

	appCfg := sdk.LoadApplicationConfigForTest[config.ECSApplicationSpec](t, appConfigFile, pluginName)

	targetDS := sdk.DeploymentSource[config.ECSApplicationSpec]{
		ApplicationDirectory: targetAppDir,
		CommitHash:           targetCommit,
		ApplicationConfig:    appCfg,
	}

	input := &sdk.GetPlanPreviewInput[config.ECSApplicationSpec]{
		Request: sdk.GetPlanPreviewRequest[config.ECSApplicationSpec]{
			ApplicationID:          "app-id",
			ApplicationName:        "simple-app",
			DeployTargets:          []string{deployTarget},
			TargetDeploymentSource: targetDS,
		},
		Logger: zaptest.NewLogger(t),
	}

	if runningAppDir != "" {
		input.Request.RunningDeploymentSource = sdk.DeploymentSource[config.ECSApplicationSpec]{
			ApplicationDirectory: runningAppDir,
			CommitHash:           runningCommit,
			ApplicationConfig:    appCfg,
		}
	}

	return input
}

func TestPlugin_GetPlanPreview(t *testing.T) {
	p := &Plugin{}

	tests := []struct {
		name          string
		targetDir     string
		runningDir    string
		wantNoChange  bool
		wantSummary   string
		wantInDetails []string
	}{
		{
			name:         "first deployment (no running source)",
			targetDir:    targetDir,
			runningDir:   "",
			wantNoChange: false,
			wantSummary:  "task definition changed, service definition changed",
			wantInDetails: []string{
				"nginx:2.0",
			},
		},
		{
			name:         "no change (same files)",
			targetDir:    runningDir,
			runningDir:   runningDir,
			wantNoChange: true,
			wantSummary:  "No changes were detected",
		},
		{
			name:         "task definition changed",
			targetDir:    targetDir,
			runningDir:   runningDir,
			wantNoChange: false,
			wantSummary:  "task definition changed",
			wantInDetails: []string{
				"-",
				"+",
				"nginx:1.0",
				"nginx:2.0",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := makeInput(t, tc.targetDir, tc.runningDir)

			resp, err := p.GetPlanPreview(context.Background(), nil, deployTargets(t), input)
			require.NoError(t, err)
			require.Len(t, resp.Results, 1)

			result := resp.Results[0]
			assert.Equal(t, deployTarget, result.DeployTarget)
			assert.Equal(t, tc.wantNoChange, result.NoChange)
			assert.Equal(t, tc.wantSummary, result.Summary)
			assert.Equal(t, "diff", result.DiffLanguage)

			for _, s := range tc.wantInDetails {
				assert.Contains(t, string(result.Details), s)
			}

			if tc.wantNoChange {
				assert.Nil(t, result.Details)
			}
		})
	}
}
