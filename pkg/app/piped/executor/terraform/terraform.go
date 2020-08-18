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
	"fmt"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/terraform"
	"github.com/pipe-cd/pipe/pkg/app/piped/executor"
	"github.com/pipe-cd/pipe/pkg/app/piped/toolregistry"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

type Executor struct {
	executor.Input

	config *config.TerraformDeploymentSpec
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
	e.config = e.DeploymentConfig.TerraformDeploymentSpec
	if e.config == nil {
		e.LogPersister.Error("Malformed deployment configuration: missing TerraformDeploymentSpec")
		return model.StageStatus_STAGE_FAILURE
	}

	var (
		ctx            = sig.Context()
		originalStatus = e.Stage.Status
		status         model.StageStatus
	)

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

func (e *Executor) ensureSync(ctx context.Context) model.StageStatus {
	// terraform init -no-color '-var-file=simplegcs/terraform.tfvars' '-lock=false' simplegcs
	// terraform workspace select -no-color default simplegcs
	// terraform validate -no-color simplegcs
	// terraform plan -no-color '-var-file=simplegcs/terraform.tfvars' '-lock=false' simplegcs
	// terraform apply -no-color -auto-approve '-input=false' '-var-file=simplegcs/terraform.tfvars' simplegcs
	return model.StageStatus_STAGE_SUCCESS
}

func (e *Executor) ensurePlan(ctx context.Context) model.StageStatus {
	return model.StageStatus_STAGE_SUCCESS
}

func (e *Executor) ensureApply(ctx context.Context) model.StageStatus {
	return model.StageStatus_STAGE_SUCCESS
}

func (e *Executor) ensureRollback(ctx context.Context) model.StageStatus {
	return model.StageStatus_STAGE_SUCCESS
}

func (e *Executor) findTerraform(ctx context.Context, version string) (*provider.Terraform, error) {
	path, installed, err := toolregistry.DefaultRegistry().Terraform(ctx, version)
	if err != nil {
		return nil, fmt.Errorf("no terraform %s (%v)", version, err)
	}
	if installed {
		e.LogPersister.Infof("Terraform %s has just been installed because of no pre-installed binary for that version", version)
	}
	return provider.NewTerraform(version, path), nil
}
