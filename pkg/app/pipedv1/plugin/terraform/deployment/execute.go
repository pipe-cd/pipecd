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
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/toolregistry"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

func executeSyncStage(ctx context.Context, input *sdk.ExecuteStageInput[config.TerraformApplicationSpec], dts []*sdk.DeployTarget[config.TerraformDeployTargetConfig]) sdk.StageStatus {
	lp := input.Client.LogPersister()

	var stageCfg config.TerraformApplyStageOptions
	if len(input.Request.StageConfig) > 0 {
		// TODO: this is a temporary solution to support the stage options specified under "with"
		// When the stage options under "with" are empty, we cannot detect whether the stage is a quick sync stage or not.
		// So we have to add a new field to the sdk.ExecuteStageRequest or sdk.Deployment to indicate that the deployment is a quick sync strategy or in a pipeline sync strategy.
		if err := json.Unmarshal(input.Request.StageConfig, &stageCfg); err != nil {
			lp.Errorf("Failed while unmarshalling stage config (%v)", err)
			return sdk.StageStatusFailure
		}
	} else {
		stageCfg = input.Request.TargetDeploymentSource.ApplicationConfig.Spec.QuickSync
	}

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

	return ensureSync(ctx, cmd, lp)
}

func executePlanStage(ctx context.Context, input *sdk.ExecuteStageInput[config.TerraformApplicationSpec], dts []*sdk.DeployTarget[config.TerraformDeployTargetConfig]) sdk.StageStatus {
	lp := input.Client.LogPersister()

	var stageCfg config.TerraformPlanStageOptions
	if err := json.Unmarshal(input.Request.StageConfig, &stageCfg); err != nil {
		lp.Errorf("Failed while unmarshalling stage config (%v)", err)
		return sdk.StageStatusFailure
	}

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

	return ensurePlan(ctx, cmd, lp, stageCfg.ExitOnNoChanges)
}

func executeApplyStage(ctx context.Context, input *sdk.ExecuteStageInput[config.TerraformApplicationSpec], dts []*sdk.DeployTarget[config.TerraformDeployTargetConfig]) sdk.StageStatus {
	lp := input.Client.LogPersister()

	var stageCfg config.TerraformApplyStageOptions
	if err := json.Unmarshal(input.Request.StageConfig, &stageCfg); err != nil {
		lp.Errorf("Failed while unmarshalling stage config (%v)", err)
		return sdk.StageStatusFailure
	}

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

	return ensureApply(ctx, cmd, lp)
}

func setupTerraformCommand(ctx context.Context, ds sdk.DeploymentSource[config.TerraformApplicationSpec], dtCfg config.TerraformDeployTargetConfig, lp sdk.StageLogPersister, tr *toolregistry.Registry) (tfcmd *provider.Terraform, ok bool) {
	if ds.ApplicationConfig == nil {
		lp.Error("Malformed application configuration: missing TerraformApplicationSpec")
		return nil, false
	}
	in := ds.ApplicationConfig.Spec.Input

	vars := make([]string, 0, len(dtCfg.Vars)+len(in.Vars))
	vars = append(vars, dtCfg.Vars...)
	vars = append(vars, in.Vars...)

	tfPath, err := tr.Terraform(ctx, in.TerraformVersion)
	if err != nil {
		lp.Errorf("Failed to get terraform command: %v", err)
		return nil, false
	}

	var (
		flags = in.CommandFlags
		envs  = in.CommandEnvs
		cmd   = provider.NewTerraform(
			tfPath,
			ds.ApplicationDirectory,
			provider.WithVars(vars),
			provider.WithVarFiles(in.VarFiles),
			provider.WithAdditionalFlags(flags.Shared, flags.Init, flags.Plan, flags.Apply),
			provider.WithAdditionalEnvs(envs.Shared, envs.Init, envs.Plan, envs.Apply),
		)
	)

	if ok := showUsingVersion(ctx, cmd, lp); !ok {
		return nil, false
	}

	if err := cmd.Init(ctx, lp); err != nil {
		lp.Errorf("Failed to init terraform (%v)", err)
		return nil, false
	}

	if ok := selectWorkspace(ctx, cmd, in.Workspace, lp); !ok {
		return nil, false
	}

	return cmd, true
}

func ensureSync(ctx context.Context, cmd *provider.Terraform, lp sdk.StageLogPersister) sdk.StageStatus {
	planResult, err := cmd.Plan(ctx, lp)
	if err != nil {
		lp.Errorf("Failed to plan (%v)", err)
		return sdk.StageStatusFailure
	}

	if planResult.NoChanges() {
		lp.Info("No changes to apply")
		return sdk.StageStatusSuccess
	}

	lp.Infof("Detected %d import, %d add, %d change, %d destroy. Those changes will be applied automatically.", planResult.Imports, planResult.Adds, planResult.Changes, planResult.Destroys)

	if err := cmd.Apply(ctx, lp); err != nil {
		lp.Errorf("Failed to apply changes (%v)", err)
		return sdk.StageStatusFailure
	}

	lp.Success("Successfully applied changes")
	return sdk.StageStatusSuccess
}

func ensurePlan(ctx context.Context, cmd *provider.Terraform, lp sdk.StageLogPersister, exitOnNoChanges bool) sdk.StageStatus {
	planResult, err := cmd.Plan(ctx, lp)
	if err != nil {
		lp.Errorf("Failed to plan (%v)", err)
		return sdk.StageStatusFailure
	}

	if planResult.NoChanges() {
		lp.Success("No changes to apply")
		if exitOnNoChanges {
			return sdk.StageStatusExited
		}
		return sdk.StageStatusSuccess
	}

	lp.Successf("Detected %d import, %d add, %d change, %d destroy.", planResult.Imports, planResult.Adds, planResult.Changes, planResult.Destroys)
	return sdk.StageStatusSuccess
}

func ensureApply(ctx context.Context, cmd *provider.Terraform, lp sdk.StageLogPersister) sdk.StageStatus {
	if err := cmd.Apply(ctx, lp); err != nil {
		lp.Errorf("Failed to apply changes (%v)", err)
		return sdk.StageStatusFailure
	}

	lp.Success("Successfully applied changes")
	return sdk.StageStatusSuccess
}

func showUsingVersion(ctx context.Context, cmd *provider.Terraform, lp sdk.StageLogPersister) bool {
	version, err := cmd.Version(ctx)
	if err != nil {
		lp.Errorf("Failed to check terraform version (%v)", err)
		return false
	}
	lp.Infof("Using terraform version %q to execute the terraform commands", version)
	return true
}

func selectWorkspace(ctx context.Context, cmd *provider.Terraform, workspace string, lp sdk.StageLogPersister) bool {
	if workspace == "" {
		return true
	}
	if err := cmd.SelectWorkspace(ctx, workspace); err != nil {
		lp.Errorf("Failed to select workspace %q (%v). You might need to create the workspace before using by command %q", workspace, err, "terraform workspace new "+workspace)
		return false
	}
	lp.Infof("Selected workspace %q", workspace)
	return true
}
