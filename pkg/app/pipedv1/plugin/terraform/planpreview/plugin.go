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
	"fmt"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/provider"
	"go.uber.org/zap"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

var (
	_ sdk.PlanPreviewPlugin[sdk.ConfigNone, config.DeployTargetConfig, config.ApplicationConfigSpec] = (*Plugin)(nil)
)

type Plugin struct{}

// GetPlanPreview implements sdk.PlanPreviewPlugin.
func (p *Plugin) GetPlanPreview(ctx context.Context, _ *sdk.ConfigNone, dts []*sdk.DeployTarget[config.DeployTargetConfig], input *sdk.GetPlanPreviewInput[config.ApplicationConfigSpec]) (*sdk.GetPlanPreviewResponse, error) {
	if len(dts) != 1 {
		return nil, fmt.Errorf("only 1 deploy target is allowed but got %d", len(dts))
	}
	dt := dts[0]

	cmd, err := provider.NewTerraformCommand(ctx, input.Client, input.Request.TargetDeploymentSource, dt)
	if err != nil {
		input.Logger.Error("Failed to initialize Terraform command", zap.Error(err))
		return nil, err
	}

	buf := &bytes.Buffer{}
	planResult, err := cmd.Plan(ctx, buf)
	if err != nil {
		input.Logger.Error("Failed to execute plan", zap.Error(err))
		return nil, err
	}

	return toResponse(planResult, buf), nil
}

func toResponse(planResult provider.PlanResult, planBuf *bytes.Buffer) *sdk.GetPlanPreviewResponse {
	if planResult.NoChanges() {
		return &sdk.GetPlanPreviewResponse{
			Summary:  "No changes were detected",
			NoChange: true,
			Details:  []byte("No changes were detected"),
		}
	}

	return &sdk.GetPlanPreviewResponse{
		Summary:  fmt.Sprintf("%d to import, %d to add, %d to change, %d to destroy", planResult.Imports, planResult.Adds, planResult.Changes, planResult.Destroys),
		NoChange: false,
		Details:  planBuf.Bytes(),
	}
}
