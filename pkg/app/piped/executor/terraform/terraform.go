// Copyright 2020 The PipeCD Authors.
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

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/terraform"
	"github.com/pipe-cd/pipe/pkg/app/piped/deploysource"
	"github.com/pipe-cd/pipe/pkg/app/piped/executor"
	"github.com/pipe-cd/pipe/pkg/app/piped/toolregistry"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

type Executor struct {
	executor.Input

	cloudProviderConfig *config.CloudProviderTerraformConfig

	repoDir       string
	appDir        string
	config        *config.TerraformDeploymentSpec
	terraformPath string
	vars          []string
}

type registerer interface {
	Register(stage model.Stage, f executor.Factory) error
	RegisterRollback(kind model.ApplicationKind, f executor.Factory) error
}

func Register(r registerer) {
	f := func(in executor.Input) executor.Executor {
		return &Executor{
			Input: in,
		}
	}
	r.Register(model.StageTerraformSync, f)
	r.Register(model.StageTerraformPlan, f)
	r.Register(model.StageTerraformApply, f)

	r.RegisterRollback(model.ApplicationKind_TERRAFORM, f)
}

func (e *Executor) Execute(sig executor.StopSignal) model.StageStatus {
	cloudProviderName := e.Application.CloudProvider
	if cloudProviderName == "" {
		e.LogPersister.Error("This application configuration was missing CloudProvider name")
		return model.StageStatus_STAGE_FAILURE
	}

	cpConfig, ok := e.PipedConfig.FindCloudProvider(cloudProviderName, model.CloudProviderTerraform)
	if !ok {
		e.LogPersister.Errorf("The specified cloud provider %q was not found in piped configuration", cloudProviderName)
		return model.StageStatus_STAGE_FAILURE
	}
	e.cloudProviderConfig = cpConfig.TerraformConfig

	var ds *deploysource.DeploySource
	var err error
	ctx := sig.Context()

	if model.Stage(e.Stage.Name) == model.StageRollback {
		ds, err = e.RunningDSP.Get(ctx, e.LogPersister)
		if err != nil {
			e.LogPersister.Errorf("Failed to prepare running deploy source data (%v)", err)
			return model.StageStatus_STAGE_FAILURE
		}
	} else {
		ds, err = e.TargetDSP.Get(ctx, e.LogPersister)
		if err != nil {
			e.LogPersister.Errorf("Failed to prepare target deploy source data (%v)", err)
			return model.StageStatus_STAGE_FAILURE
		}
	}

	e.repoDir = ds.RepoDir
	e.appDir = ds.AppDir
	e.config = ds.DeploymentConfig.TerraformDeploymentSpec
	if e.config == nil {
		e.LogPersister.Error("Malformed deployment configuration: missing TerraformDeploymentSpec")
		return model.StageStatus_STAGE_FAILURE
	}

	e.vars = make([]string, 0, len(e.cloudProviderConfig.Vars)+len(e.config.Input.Vars))
	e.vars = append(e.vars, e.cloudProviderConfig.Vars...)
	e.vars = append(e.vars, e.config.Input.Vars...)

	var (
		originalStatus = e.Stage.Status
		status         model.StageStatus
	)

	execPath, ok := e.findTerraform(ctx, e.config.Input.TerraformVersion)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}
	e.terraformPath = execPath

	switch model.Stage(e.Stage.Name) {
	case model.StageTerraformSync:
		status = e.ensureSync(ctx)

	case model.StageTerraformPlan:
		status = e.ensurePlan(ctx)

	case model.StageTerraformApply:
		status = e.ensureApply(ctx)

	case model.StageRollback:
		status = e.ensureRollback(ctx)

	default:
		e.LogPersister.Errorf("Unsupported stage %s for cloudrun application", e.Stage.Name)
		return model.StageStatus_STAGE_FAILURE
	}

	return executor.DetermineStageStatus(sig.Signal(), originalStatus, status)
}

func (e *Executor) showUsingVersion(ctx context.Context, cmd *provider.Terraform) bool {
	version, err := cmd.Version(ctx)
	if err != nil {
		e.LogPersister.Errorf("Failed to check terraform version (%v)", err)
		return false
	}
	e.LogPersister.Infof("Using terraform version %q to execute the terraform commands", version)
	return true
}

func (e *Executor) selectWorkspace(ctx context.Context, cmd *provider.Terraform) bool {
	workspace := e.config.Input.Workspace
	if workspace == "" {
		return true
	}
	if err := cmd.SelectWorkspace(ctx, workspace); err != nil {
		e.LogPersister.Errorf("Failed to select workspace %q (%v). You might need to create the workspace before using by command %q", workspace, err, "terraform workspace new "+workspace)
		return false
	}
	e.LogPersister.Infof("Selected workspace %q", workspace)
	return true
}

func (e *Executor) findTerraform(ctx context.Context, version string) (string, bool) {
	path, installed, err := toolregistry.DefaultRegistry().Terraform(ctx, version)
	if err != nil {
		e.LogPersister.Errorf("Unable to find required terraform %q (%v)", version, err)
		return "", false
	}
	if installed {
		e.LogPersister.Infof("Terraform %q has just been installed to %q because of no pre-installed binary for that version", version, path)
	}
	return path, true
}
