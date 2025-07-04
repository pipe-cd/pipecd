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

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/config"
)

func (p *Plugin) executePlanStage(ctx context.Context, input *sdk.ExecuteStageInput[config.ApplicationConfigSpec], dts []*sdk.DeployTarget[config.DeployTargetConfig]) sdk.StageStatus {
	cmd, err := initTerraformCommand(ctx, input.Client, input.Request.TargetDeploymentSource, dts[0])
	if err != nil {
		return sdk.StageStatusFailure
	}

	lp := input.Client.LogPersister()

	stageConfig := config.TerraformPlanStageOptions{}
	if err := json.Unmarshal(input.Request.StageConfig, &stageConfig); err != nil {
		lp.Errorf("Failed to unmarshal stage config (%v)", err)
		return sdk.StageStatusFailure
	}

	planResult, err := cmd.Plan(ctx, lp)
	if err != nil {
		lp.Errorf("Failed to plan (%v)", err)
		return sdk.StageStatusFailure
	}

	if planResult.NoChanges() {
		lp.Success("No changes to apply")
		if stageConfig.ExitOnNoChanges {
			return sdk.StageStatusExited
		}
		return sdk.StageStatusSuccess
	}

	lp.Successf("Detected %d import, %d add, %d change, %d destroy.", planResult.Imports, planResult.Adds, planResult.Changes, planResult.Destroys)
	return sdk.StageStatusSuccess
}
