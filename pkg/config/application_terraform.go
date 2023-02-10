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

package config

// TerraformApplicationSpec represents an application configuration for Terraform application.
type TerraformApplicationSpec struct {
	GenericApplicationSpec
	// Input for Terraform deployment such as terraform version, workspace...
	Input TerraformDeploymentInput `json:"input"`
	// Configuration for quick sync.
	QuickSync TerraformApplyStageOptions `json:"quickSync"`
}

// Validate returns an error if any wrong configuration value was found.
func (s *TerraformApplicationSpec) Validate() error {
	if err := s.GenericApplicationSpec.Validate(); err != nil {
		return err
	}
	return nil
}

type TerraformDeploymentInput struct {
	// The terraform workspace name.
	// Empty means "default" workpsace.
	Workspace string `json:"workspace,omitempty"`
	// The version of terraform should be used.
	// Empty means the pre-installed version will be used.
	TerraformVersion string `json:"terraformVersion,omitempty"`
	// List of variables that will be set directly on terraform commands with "-var" flag.
	// The variable must be formatted by "key=value" as below:
	// "image_id=ami-abc123"
	// 'image_id_list=["ami-abc123","ami-def456"]'
	// 'image_id_map={"us-east-1":"ami-abc123","us-east-2":"ami-def456"}'
	Vars []string `json:"vars,omitempty"`
	// List of variable files that will be set on terraform commands with "-var-file" flag.
	VarFiles []string `json:"varFiles,omitempty"`
	// Automatically reverts all changes from all stages when one of them failed.
	// Default is false.
	AutoRollback bool `json:"autoRollback"`
	// List of additional flags will be used while executing terraform commands.
	CommandFlags TerraformCommandFlags `json:"commandFlags"`
	// List of additional environment variables will be used while executing terraform commands.
	CommandEnvs TerraformCommandEnvs `json:"commandEnvs"`
}

// TerraformSyncStageOptions contains all configurable values for a TERRAFORM_SYNC stage.
type TerraformSyncStageOptions struct {
	// How many times to retry applying terraform changes.
	Retries int `json:"retries"`
}

// TerraformPlanStageOptions contains all configurable values for a TERRAFORM_PLAN stage.
type TerraformPlanStageOptions struct {
	// Exit the pipeline if the result is "No Changes" with success status.
	ExitOnNoChanges bool `json:"exitOnNoChanges"`
}

// TerraformApplyStageOptions contains all configurable values for a TERRAFORM_APPLY stage.
type TerraformApplyStageOptions struct {
	// How many times to retry applying terraform changes.
	Retries int `json:"retries"`
}

// TerraformCommandFlags contains all additional flags will be used while executing terraform commands.
type TerraformCommandFlags struct {
	Shared []string `json:"shared"`
	Init   []string `json:"init"`
	Plan   []string `json:"plan"`
	Apply  []string `json:"apply"`
}

// TerraformCommandEnvs contains all additional environment variables will be used while executing terraform commands.
type TerraformCommandEnvs struct {
	Shared []string `json:"shared"`
	Init   []string `json:"init"`
	Plan   []string `json:"plan"`
	Apply  []string `json:"apply"`
}
