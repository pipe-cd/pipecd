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

package config

// TerraformDeploymentSpec represents a deployment configuration for Terraform application.
type TerraformDeploymentSpec struct {
	Input     TerraformDeploymentInput   `json:"input"`
	QuickSync TerraformApplyStageOptions `json:"quickSync"`
	Pipeline  *DeploymentPipeline        `json:"pipeline"`
}

func (s *TerraformDeploymentSpec) GetStage(index int32) (PipelineStage, bool) {
	if s.Pipeline == nil {
		return PipelineStage{}, false
	}
	if int(index) >= len(s.Pipeline.Stages) {
		return PipelineStage{}, false
	}
	return s.Pipeline.Stages[index], true
}

// Validate returns an error if any wrong configuration value was found.
func (s *TerraformDeploymentSpec) Validate() error {
	return nil
}

type TerraformDeploymentInput struct {
	Workspace        string `json:"workspace,omitempty"`
	TerraformVersion string `json:"terraformVersion,omitempty"`
	// Automatically reverts all changes from all stages when one of them failed.
	// Default is false.
	AutoRollback bool     `json:"autoRollback"`
	Dependencies []string `json:"dependencies,omitempty"`
}

// TerraformSyncStageOptions contains all configurable values for a TERRAFORM_SYNC stage.
type TerraformSyncStageOptions struct {
	// How many times to retry applying terraform changes.
	Retries int `json:"retries"`
}

// TerraformPlanStageOptions contains all configurable values for a TERRAFORM_PLAN stage.
type TerraformPlanStageOptions struct {
}

// TerraformApplyStageOptions contains all configurable values for a TERRAFORM_APPLY stage.
type TerraformApplyStageOptions struct {
	// How many times to retry applying terraform changes.
	Retries int `json:"retries"`
}
