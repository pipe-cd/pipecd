// Copyright 2023 The PipeCD Authors.
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

package terraform

import (
	"context"

	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/terraform"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type deployExecutor struct {
	executor.Input

	repoDir       string
	appDir        string
	vars          []string
	terraformPath string
	appCfg        *config.TerraformApplicationSpec
}

func (e *deployExecutor) Execute(sig executor.StopSignal) model.StageStatus {
	providerCfg, found := findPlatformProvider(&e.Input)
	if !found {
		return model.StageStatus_STAGE_FAILURE
	}

	ctx := sig.Context()
	ds, err := e.TargetDSP.Get(ctx, e.LogPersister)
	if err != nil {
		e.LogPersister.Errorf("Failed to prepare target deploy source data (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	e.appCfg = ds.ApplicationConfig.TerraformApplicationSpec
	if e.appCfg == nil {
		e.LogPersister.Error("Malformed application configuration: missing TerraformApplicationSpec")
		return model.StageStatus_STAGE_FAILURE
	}

	e.repoDir = ds.RepoDir
	e.appDir = ds.AppDir

	e.vars = make([]string, 0, len(providerCfg.Vars)+len(e.appCfg.Input.Vars))
	e.vars = append(e.vars, providerCfg.Vars...)
	e.vars = append(e.vars, e.appCfg.Input.Vars...)

	var (
		originalStatus = e.Stage.Status
		status         model.StageStatus
	)

	var ok bool
	e.terraformPath, ok = findTerraform(ctx, e.appCfg.Input.TerraformVersion, e.LogPersister)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	switch model.Stage(e.Stage.Name) {
	case model.StageTerraformSync:
		status = e.ensureSync(ctx)

	case model.StageTerraformPlan:
		status = e.ensurePlan(ctx)

	case model.StageTerraformApply:
		status = e.ensureApply(ctx)

	default:
		e.LogPersister.Errorf("Unsupported stage %s for terraform application", e.Stage.Name)
		return model.StageStatus_STAGE_FAILURE
	}

	return executor.DetermineStageStatus(sig.Signal(), originalStatus, status)
}

func (e *deployExecutor) ensureSync(ctx context.Context) model.StageStatus {
	var (
		flags = e.appCfg.Input.CommandFlags
		envs  = e.appCfg.Input.CommandEnvs
		cmd   = provider.NewTerraform(
			e.terraformPath,
			e.appDir,
			provider.WithVars(e.vars),
			provider.WithVarFiles(e.appCfg.Input.VarFiles),
			provider.WithAdditionalFlags(flags.Shared, flags.Init, flags.Plan, flags.Apply),
			provider.WithAdditionalEnvs(envs.Shared, envs.Init, envs.Plan, envs.Apply),
		)
	)

	if ok := showUsingVersion(ctx, cmd, e.LogPersister); !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	if err := cmd.Init(ctx, e.LogPersister); err != nil {
		e.LogPersister.Errorf("Failed to init (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	if ok := selectWorkspace(ctx, cmd, e.appCfg.Input.Workspace, e.LogPersister); !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	planResult, err := cmd.Plan(ctx, e.LogPersister)
	if err != nil {
		e.LogPersister.Errorf("Failed to plan (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	if planResult.NoChanges() {
		e.LogPersister.Info("No changes to apply")
		return model.StageStatus_STAGE_SUCCESS
	}

	e.LogPersister.Infof("Detected %d add, %d change, %d destroy. Those changes will be applied automatically.", planResult.Adds, planResult.Changes, planResult.Destroys)

	if err := cmd.Apply(ctx, e.LogPersister); err != nil {
		e.LogPersister.Errorf("Failed to apply changes (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	e.LogPersister.Success("Successfully applied changes")
	return model.StageStatus_STAGE_SUCCESS
}

func (e *deployExecutor) ensurePlan(ctx context.Context) model.StageStatus {
	var (
		flags = e.appCfg.Input.CommandFlags
		envs  = e.appCfg.Input.CommandEnvs
		cmd   = provider.NewTerraform(
			e.terraformPath,
			e.appDir,
			provider.WithVars(e.vars),
			provider.WithVarFiles(e.appCfg.Input.VarFiles),
			provider.WithAdditionalFlags(flags.Shared, flags.Init, flags.Plan, flags.Apply),
			provider.WithAdditionalEnvs(envs.Shared, envs.Init, envs.Plan, envs.Apply),
		)
	)

	if ok := showUsingVersion(ctx, cmd, e.LogPersister); !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	if err := cmd.Init(ctx, e.LogPersister); err != nil {
		e.LogPersister.Errorf("Failed to init (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	if ok := selectWorkspace(ctx, cmd, e.appCfg.Input.Workspace, e.LogPersister); !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	planResult, err := cmd.Plan(ctx, e.LogPersister)
	if err != nil {
		e.LogPersister.Errorf("Failed to plan (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	if planResult.NoChanges() {
		e.LogPersister.Success("No changes to apply")
		if e.StageConfig.TerraformPlanStageOptions.ExitOnNoChanges {
			return model.StageStatus_STAGE_EXITED
		}
		return model.StageStatus_STAGE_SUCCESS
	}

	e.LogPersister.Successf("Detected %d add, %d change, %d destroy.", planResult.Adds, planResult.Changes, planResult.Destroys)
	return model.StageStatus_STAGE_SUCCESS
}

func (e *deployExecutor) ensureApply(ctx context.Context) model.StageStatus {
	var (
		flags = e.appCfg.Input.CommandFlags
		envs  = e.appCfg.Input.CommandEnvs
		cmd   = provider.NewTerraform(
			e.terraformPath,
			e.appDir,
			provider.WithVars(e.vars),
			provider.WithVarFiles(e.appCfg.Input.VarFiles),
			provider.WithAdditionalFlags(flags.Shared, flags.Init, flags.Plan, flags.Apply),
			provider.WithAdditionalEnvs(envs.Shared, envs.Init, envs.Plan, envs.Apply),
		)
	)

	if ok := showUsingVersion(ctx, cmd, e.LogPersister); !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	if err := cmd.Init(ctx, e.LogPersister); err != nil {
		e.LogPersister.Errorf("Failed to init (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	if ok := selectWorkspace(ctx, cmd, e.appCfg.Input.Workspace, e.LogPersister); !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	if err := cmd.Apply(ctx, e.LogPersister); err != nil {
		e.LogPersister.Errorf("Failed to apply changes (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	e.LogPersister.Success("Successfully applied changes")
	return model.StageStatus_STAGE_SUCCESS
}
