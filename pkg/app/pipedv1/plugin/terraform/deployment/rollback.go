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

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/provider"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

// TODO: add test
func (p *Plugin) executeRollbackStage(ctx context.Context, input *sdk.ExecuteStageInput[config.ApplicationConfigSpec], dts []*sdk.DeployTarget[config.DeployTargetConfig]) sdk.StageStatus {
	lp := input.Client.LogPersister()
	rds := input.Request.RunningDeploymentSource

	if rds.CommitHash == "" {
		lp.Errorf("Unable to determine the last deployed commit to rollback. It seems this is the first deployment.")
		return sdk.StageStatusFailure
	}

	cmd, err := provider.NewTerraformCommand(ctx, input.Client, rds, dts[0])
	if err != nil {
		return sdk.StageStatusFailure
	}

	lp.Infof("Start rolling back to the state defined at commit %s", rds.CommitHash)
	if err = cmd.Apply(ctx, lp); err != nil {
		lp.Errorf("Failed to apply changes (%v)", err)
		return sdk.StageStatusFailure
	}

	lp.Success("Successfully rolled back the changes")
	return sdk.StageStatusSuccess
}
