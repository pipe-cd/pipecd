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
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/toolregistry"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

func executeRollbackStage(ctx context.Context, input *sdk.ExecuteStageInput[config.TerraformApplicationSpec], dts []*sdk.DeployTarget[config.TerraformDeployTargetConfig]) sdk.StageStatus {
	lp := input.Client.LogPersister()

	cmd, ok := setupTerraformCommand(
		ctx,
		input.Request.TargetDeploymentSource,
		dts[0].Config,
		lp,
		toolregistry.NewRegistry(input.Client.ToolRegistry()),
	)
	if !ok {
		return sdk.StageStatusFailure
	}

	return ensureRollback(ctx, cmd, lp)
}

func ensureRollback(ctx context.Context, cmd *provider.Terraform, lp sdk.StageLogPersister) sdk.StageStatus {
	if err := cmd.Apply(ctx, lp); err != nil {
		lp.Errorf("Failed to apply changes (%v)", err)
		return sdk.StageStatusFailure
	}

	lp.Success("Successfully rolled back the changes")
	return sdk.StageStatusSuccess
}
