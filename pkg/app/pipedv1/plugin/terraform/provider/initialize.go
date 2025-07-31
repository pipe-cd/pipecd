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

package provider

import (
	"context"
	"errors"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/toolregistry"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

// NewTerraformCommand initializes a Terraform command, including `terraform init`.
func NewTerraformCommand(ctx context.Context, client *sdk.Client, ds sdk.DeploymentSource[config.ApplicationConfigSpec], dt *sdk.DeployTarget[config.DeployTargetConfig]) (*Terraform, error) {
	var (
		appSpec = ds.ApplicationConfig.Spec
		flags   = appSpec.CommandFlags
		envs    = appSpec.CommandEnvs
		lp      = client.LogPersister()
	)
	tr := toolregistry.NewRegistry(client.ToolRegistry())
	terraformPath, err := tr.Terraform(ctx, appSpec.TerraformVersion)
	if err != nil {
		lp.Errorf("Failed to find terraform (%v)", err)
		return nil, err
	}

	cmd := newTerraform(
		terraformPath,
		ds.ApplicationDirectory,
		WithVars(mergeVars(dt.Config.Vars, appSpec.Vars)),
		WithVarFiles(appSpec.VarFiles),
		WithAdditionalFlags(flags.Shared, flags.Init, flags.Plan, flags.Apply),
		WithAdditionalEnvs(envs.Shared, envs.Init, envs.Plan, envs.Apply),
	)

	if ok := showUsingVersion(ctx, cmd, lp); !ok {
		return nil, errors.New("failed to show using version")
	}

	if err := cmd.init(ctx, lp); err != nil {
		lp.Errorf("Failed to execute 'terraform init' (%v)", err)
		return nil, err
	}

	if ok := selectWorkspace(ctx, cmd, appSpec.Workspace, lp); !ok {
		return nil, errors.New("failed to select workspace")
	}

	return cmd, nil
}

func mergeVars(deployTargetVars []string, appVars []string) []string {
	// TODO: Validate duplication
	mergedVars := make([]string, 0, len(deployTargetVars)+len(appVars))
	mergedVars = append(mergedVars, deployTargetVars...)
	mergedVars = append(mergedVars, appVars...)
	return mergedVars
}

func showUsingVersion(ctx context.Context, cmd *Terraform, lp sdk.StageLogPersister) bool {
	version, err := cmd.version(ctx)
	if err != nil {
		lp.Errorf("Failed to check terraform version (%v)", err)
		return false
	}
	lp.Infof("Using terraform version %q to execute the terraform commands", version)
	return true
}

func selectWorkspace(ctx context.Context, cmd *Terraform, workspace string, lp sdk.StageLogPersister) bool {
	if workspace == "" {
		return true
	}
	if err := cmd.selectWorkspace(ctx, workspace); err != nil {
		lp.Errorf("Failed to select workspace %q (%v). You might need to create the workspace before using by command %q", workspace, err, "terraform workspace new "+workspace)
		return false
	}
	lp.Infof("Selected workspace %q", workspace)
	return true
}
