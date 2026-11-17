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

package planpreview

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/pipe-cd/piped-plugin-sdk-go/toolregistry/toolregistrytest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/provider"
)

func TestToResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		planResult provider.PlanResult
		planBuf    *bytes.Buffer
		want       *sdk.GetPlanPreviewResponse
	}{
		{
			name: "no changes",
			planResult: provider.PlanResult{
				Imports:  0,
				Adds:     0,
				Changes:  0,
				Destroys: 0,
			},
			planBuf: bytes.NewBuffer([]byte("No changes.")),
			want: &sdk.GetPlanPreviewResponse{
				Results: []sdk.PlanPreviewResult{
					{
						DeployTarget: "dt-1",
						NoChange:     true,
						Summary:      "No changes were detected",
						DiffLanguage: "hcl",
					},
				},
			},
		},
		{
			name: "with changes",
			planResult: provider.PlanResult{
				Imports:  1,
				Adds:     2,
				Changes:  3,
				Destroys: 4,
			},
			planBuf: bytes.NewBuffer([]byte(`
Terraform will perform the following actions:
<plan-details>
───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.`)),
			want: &sdk.GetPlanPreviewResponse{
				Results: []sdk.PlanPreviewResult{
					{
						DeployTarget: "dt-1",
						NoChange:     false,
						Summary:      "1 to import, 2 to add, 3 to change, 4 to destroy",
						DiffLanguage: "hcl",
						Details: []byte(`
Terraform will perform the following actions:
<plan-details>
───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.`),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := toResponse(tt.planResult, tt.planBuf, "dt-1")

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetPlanPreview_WithoutColor(t *testing.T) {
	t.Parallel()

	appDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(appDir, "main.tf"), []byte(`resource "terraform_data" "example" {
  input = "hello"
}
`), 0o644))

	p := &Plugin{}
	got, err := p.GetPlanPreview(context.Background(), nil, []*sdk.DeployTarget[config.DeployTargetConfig]{{Name: "dt-1"}}, &sdk.GetPlanPreviewInput[config.ApplicationConfigSpec]{
		Request: sdk.GetPlanPreviewRequest[config.ApplicationConfigSpec]{
			TargetDeploymentSource: sdk.DeploymentSource[config.ApplicationConfigSpec]{
				ApplicationDirectory: appDir,
				ApplicationConfig: &sdk.ApplicationConfig[config.ApplicationConfigSpec]{
					Spec: &config.ApplicationConfigSpec{TerraformVersion: "1.12.0"},
				},
			},
		},
		// Plan preview runs outside a deployment, so there is no stage log persister.
		Client: sdk.NewClient(nil, "terraform", "app-id", "", nil, toolregistrytest.NewTestToolRegistry(t)),
		Logger: zap.NewNop(),
	})
	require.NoError(t, err)
	require.Len(t, got.Results, 1)

	result := got.Results[0]
	assert.False(t, result.NoChange)
	assert.Equal(t, "0 to import, 1 to add, 0 to change, 0 to destroy", result.Summary)
	assert.Contains(t, string(result.Details), "Plan: 1 to add, 0 to change, 0 to destroy.")
	assert.NotContains(t, string(result.Details), "\x1b[", "plan preview details must not contain ANSI escape codes")
}
