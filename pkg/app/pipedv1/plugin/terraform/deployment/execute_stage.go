// Copyright 2024 The PipeCD Authors.
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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tfconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/provider"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister"
)

type deployExecutor struct {
	appDir        string
	vars          []string
	terraformPath string
	appCfg        *tfconfig.TerraformApplicationSpec
}

// Memo: Copied from pkg/app/piped/executor/terraform/deploy.go > Execute()
func (s *DeploymentServiceServer) executeStage(ctx context.Context, slp logpersister.StageLogPersister, input *deployment.ExecutePluginInput) (model.StageStatus, error) {
	cfg, err := config.DecodeYAML[*tfconfig.TerraformApplicationSpec](input.GetTargetDeploymentSource().GetApplicationConfig())
	if err != nil {
		slp.Errorf("Failed while decoding application config (%v)", err)
		return model.StageStatus_STAGE_FAILURE, err
	}

	e := &deployExecutor{
		appCfg: cfg.Spec,
	}
	e.appDir = string(input.GetTargetDeploymentSource().GetApplicationDirectory())
	e.vars = mergeVars(&s.deployTargetConfig, cfg.Spec)
	e.terraformPath, err = s.toolRegistry.Terraform(ctx, cfg.Spec.Input.TerraformVersion)
	if err != nil {
		return model.StageStatus_STAGE_FAILURE, err
	}

	switch input.GetStage().GetName() {
	case stageTerraformSync.String():
		return e.ensureSync(ctx, cfg.Spec, slp), nil
	case stageTerraformRollback.String():
		// TODO: implement executeTerraformRollbackStage
		// return s.executeTerraformRollbackStage(ctx, lp, input), nil
		return model.StageStatus_STAGE_FAILURE, nil
	// TODO: Add PLAN, APPLY Stages
	// case stageTerraformPlan.String():
	// case stageTerraformApply.String():
	default:
		return model.StageStatus_STAGE_FAILURE, status.Error(codes.InvalidArgument, "unsupported stage")
	}
}

// Memo: Copied from pkg/app/piped/executor/terraform/deploy.go > ensureSync()
func (e *deployExecutor) ensureSync(ctx context.Context, cfg *tfconfig.TerraformApplicationSpec, slp logpersister.StageLogPersister) model.StageStatus {
	var (
		flags = cfg.Input.CommandFlags
		envs  = cfg.Input.CommandEnvs
		tfcmd = provider.NewTerraform(
			e.terraformPath,
			e.appDir,
			provider.WithVars(e.vars),
			provider.WithVarFiles(cfg.Input.VarFiles),
			provider.WithAdditionalFlags(flags.Shared, flags.Init, flags.Plan, flags.Apply),
			provider.WithAdditionalEnvs(envs.Shared, envs.Init, envs.Plan, envs.Apply),
		)
	)

	if ok := showUsingVersion(ctx, tfcmd, slp); !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	if err := tfcmd.Init(ctx, slp); err != nil {
		slp.Errorf("Failed to init (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	if ok := selectWorkspace(ctx, tfcmd, cfg.Input.Workspace, slp); !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	planResult, err := tfcmd.Plan(ctx, slp)
	if err != nil {
		slp.Errorf("Failed to plan (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	if planResult.NoChanges() {
		slp.Info("No changes to apply")
		return model.StageStatus_STAGE_SUCCESS
	}

	slp.Infof("Detected %d import, %d add, %d change, %d destroy. Those changes will be applied automatically.", planResult.Imports, planResult.Adds, planResult.Changes, planResult.Destroys)

	if err := tfcmd.Apply(ctx, slp); err != nil {
		slp.Errorf("Failed to apply changes (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	slp.Success("Successfully applied changes")
	return model.StageStatus_STAGE_SUCCESS
}

func mergeVars(dtCfg *tfconfig.TerraformDeployTargetConfig, appSpec *tfconfig.TerraformApplicationSpec) []string {
	vars := make([]string, 0, len(dtCfg.Vars)+len(appSpec.Input.Vars))
	vars = append(vars, dtCfg.Vars...)
	vars = append(vars, appSpec.Input.Vars...)
	return vars
}

func showUsingVersion(ctx context.Context, tfcmd *provider.Terraform, slp logpersister.StageLogPersister) (ok bool) {
	version, err := tfcmd.Version(ctx)
	if err != nil {
		slp.Errorf("Failed to check terraform version (%v)", err)
		return false
	}
	slp.Infof("Using terraform version %q to execute the terraform commands", version)
	return true
}

func selectWorkspace(ctx context.Context, tfcmd *provider.Terraform, workspace string, slp logpersister.StageLogPersister) bool {
	if workspace == "" {
		return true
	}
	if err := tfcmd.SelectWorkspace(ctx, workspace); err != nil {
		slp.Errorf("Failed to select workspace %q (%v). You might need to create the workspace before using by command %q", workspace, err, "terraform workspace new "+workspace)
		return false
	}
	slp.Infof("Selected workspace %q", workspace)
	return true
}
