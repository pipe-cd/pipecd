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

package planpreview

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/pipe-cd/pipecd/pkg/app/piped/deploysource"
	terraformprovider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/terraform"
	"github.com/pipe-cd/pipecd/pkg/app/piped/toolregistry"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func (b *builder) terraformDiff(
	ctx context.Context,
	app *model.Application,
	targetDSP deploysource.Provider,
	buf *bytes.Buffer,
) (*diffResult, error) {

	cp, ok := b.pipedCfg.FindPlatformProvider(app.PlatformProvider, model.ApplicationKind_TERRAFORM)
	if !ok {
		err := fmt.Errorf("platform provider %s was not found in Piped config", app.PlatformProvider)
		fmt.Fprintln(buf, err.Error())
		return nil, err
	}
	cpCfg := cp.TerraformConfig

	ds, err := targetDSP.Get(ctx, io.Discard)
	if err != nil {
		fmt.Fprintf(buf, "failed to prepare deploy source data at the head commit (%v)\n", err)
		return nil, err
	}

	appCfg := ds.ApplicationConfig.TerraformApplicationSpec
	if appCfg == nil {
		err := fmt.Errorf("missing Terraform spec field in application configuration")
		fmt.Fprintln(buf, err.Error())
		return nil, err
	}

	version := appCfg.Input.TerraformVersion
	terraformPath, installed, err := toolregistry.DefaultRegistry().Terraform(ctx, version)
	if err != nil {
		fmt.Fprintf(buf, "unable to find the specified terraform version %q (%v)\n", version, err)
		return nil, err
	}
	if installed {
		b.logger.Info(fmt.Sprintf("terraform %q has just been installed to %q because of no pre-installed binary for that version", version, terraformPath))
	}

	vars := make([]string, 0, len(cpCfg.Vars)+len(appCfg.Input.Vars))
	vars = append(vars, cpCfg.Vars...)
	vars = append(vars, appCfg.Input.Vars...)
	flags := appCfg.Input.CommandFlags
	envs := appCfg.Input.CommandEnvs

	executor := terraformprovider.NewTerraform(
		terraformPath,
		ds.AppDir,
		terraformprovider.WithoutColor(),
		terraformprovider.WithVars(vars),
		terraformprovider.WithVarFiles(appCfg.Input.VarFiles),
		terraformprovider.WithAdditionalFlags(flags.Shared, flags.Init, flags.Plan, flags.Apply),
		terraformprovider.WithAdditionalEnvs(envs.Shared, envs.Init, envs.Plan, envs.Apply),
	)

	if err := executor.Init(ctx, buf); err != nil {
		fmt.Fprintf(buf, "failed while executing terraform init (%v)\n", err)
		return nil, err
	}

	if ws := appCfg.Input.Workspace; ws != "" {
		if err := executor.SelectWorkspace(ctx, ws); err != nil {
			fmt.Fprintf(buf, "failed to select workspace %q (%v). You might need to create the workspace before using by command %q\n",
				ws,
				err,
				"terraform workspace new "+ws,
			)
			return nil, err
		}
		fmt.Fprintf(buf, "selected workspace %q\n", ws)
	}

	result, err := executor.Plan(ctx, buf)
	if err != nil {
		fmt.Fprintf(buf, "failed while executing terraform plan (%v)\n", err)
		return nil, err
	}

	if result.NoChanges() {
		fmt.Fprintln(buf, "No changes were detected")
		return &diffResult{
			summary:  "No changes were detected",
			noChange: true,
		}, nil
	}

	summary := fmt.Sprintf("%d to import, %d to add, %d to change, %d to destroy", result.Imports, result.Adds, result.Changes, result.Destroys)
	fmt.Fprintln(buf, summary)
	return &diffResult{
		summary: summary,
	}, nil
}
