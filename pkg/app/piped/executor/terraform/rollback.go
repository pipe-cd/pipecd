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

package terraform

import (
	"context"

	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/terraform"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type rollbackExecutor struct {
	executor.Input
}

func (e *rollbackExecutor) Execute(sig executor.StopSignal) model.StageStatus {
	var (
		ctx            = sig.Context()
		originalStatus = e.Stage.Status
		status         model.StageStatus
	)

	switch model.Stage(e.Stage.Name) {
	case model.StageRollback:
		status = e.ensureRollback(ctx)

	default:
		e.LogPersister.Errorf("Unsupported stage %s for terraform application", e.Stage.Name)
		return model.StageStatus_STAGE_FAILURE
	}

	return executor.DetermineStageStatus(sig.Signal(), originalStatus, status)
}

func (e *rollbackExecutor) ensureRollback(ctx context.Context) model.StageStatus {
	// There is nothing to do if this is the first deployment.
	if e.Deployment.RunningCommitHash == "" {
		e.LogPersister.Errorf("Unable to determine the last deployed commit to rollback. It seems this is the first deployment.")
		return model.StageStatus_STAGE_FAILURE
	}

	providerCfg, found := findPlatformProvider(&e.Input)
	if !found {
		return model.StageStatus_STAGE_FAILURE
	}

	ds, err := e.RunningDSP.Get(ctx, e.LogPersister)
	if err != nil {
		e.LogPersister.Errorf("Failed to prepare running deploy source data (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	appCfg := ds.ApplicationConfig.TerraformApplicationSpec
	if appCfg == nil {
		e.LogPersister.Error("Malformed application configuration: missing TerraformApplicationSpec")
		return model.StageStatus_STAGE_FAILURE
	}

	terraformPath, ok := findTerraform(ctx, appCfg.Input.TerraformVersion, e.LogPersister)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	vars := make([]string, 0, len(providerCfg.Vars)+len(appCfg.Input.Vars))
	vars = append(vars, providerCfg.Vars...)
	vars = append(vars, appCfg.Input.Vars...)

	e.LogPersister.Infof("Start rolling back to the state defined at commit %s", e.Deployment.RunningCommitHash)
	var (
		flags = appCfg.Input.CommandFlags
		envs  = appCfg.Input.CommandEnvs
		cmd   = provider.NewTerraform(
			terraformPath,
			ds.AppDir,
			provider.WithVars(vars),
			provider.WithVarFiles(appCfg.Input.VarFiles),
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

	if ok := selectWorkspace(ctx, cmd, appCfg.Input.Workspace, e.LogPersister); !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	if err := cmd.Apply(ctx, e.LogPersister); err != nil {
		e.LogPersister.Errorf("Failed to apply changes (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	e.LogPersister.Success("Successfully rolled back the changes")
	return model.StageStatus_STAGE_SUCCESS
}
