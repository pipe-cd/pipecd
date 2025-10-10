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
	"encoding/json"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/provider"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"go.uber.org/zap"
)

// TODO: add test
func (p *Plugin) executePlanStage(ctx context.Context, input *sdk.ExecuteStageInput[config.ApplicationConfigSpec], dts []*sdk.DeployTarget[config.DeployTargetConfig]) sdk.StageStatus {
	slp, err := input.Client.StageLogPersister()
	if err != nil {
		input.Logger.Error("No stage log persister available", zap.Error(err))
		return sdk.StageStatusFailure
	}
	cmd, err := provider.NewTerraformCommand(ctx, input.Client, input.Request.TargetDeploymentSource, dts[0])
	if err != nil {
		slp.Errorf("Failed to initialize Terraform command (%v)", err)
		return sdk.StageStatusFailure
	}

	stageConfig := config.TerraformPlanStageOptions{}
	if err := json.Unmarshal(input.Request.StageConfig, &stageConfig); err != nil {
		slp.Errorf("Failed to unmarshal stage config (%v)", err)
		return sdk.StageStatusFailure
	}

	planResult, err := cmd.Plan(ctx, slp)
	if err != nil {
		slp.Errorf("Failed to plan (%v)", err)
		return sdk.StageStatusFailure
	}

	if planResult.NoChanges() {
		slp.Success("No changes to apply")
		if stageConfig.ExitOnNoChanges {
			return sdk.StageStatusExited
		}
		return sdk.StageStatusSuccess
	}

	slp.Successf("Detected %d import, %d add, %d change, %d destroy.", planResult.Imports, planResult.Adds, planResult.Changes, planResult.Destroys)
	return sdk.StageStatusSuccess
}
