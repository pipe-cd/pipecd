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

package config

// LambdaApplicationSpec represents an application configuration for Lambda application.
type LambdaApplicationSpec struct {
	GenericApplicationSpec
	// Input for Lambda deployment such as where to fetch source code...
	Input LambdaDeploymentInput `json:"input"`
	// Configuration for quick sync.
	QuickSync LambdaSyncStageOptions `json:"quickSync"`
}

// Validate returns an error if any wrong configuration value was found.
func (s *LambdaApplicationSpec) Validate() error {
	if err := s.GenericApplicationSpec.Validate(); err != nil {
		return err
	}
	return nil
}

type LambdaDeploymentInput struct {
	// The name of service manifest file placing in application directory.
	// Default is function.yaml
	FunctionManifestFile string `json:"functionManifestFile" default:"function.yaml"`
	// Automatically reverts all changes from all stages when one of them failed.
	// Default is true.
	AutoRollback *bool `json:"autoRollback,omitempty" default:"true"`
}

// LambdaSyncStageOptions contains all configurable values for a LAMBDA_SYNC stage.
type LambdaSyncStageOptions struct {
}

// LambdaCanaryRolloutStageOptions contains all configurable values for a LAMBDA_CANARY_ROLLOUT stage.
type LambdaCanaryRolloutStageOptions struct {
}

// LambdaPromoteStageOptions contains all configurable values for a LAMBDA_PROMOTE stage.
type LambdaPromoteStageOptions struct {
	// Percentage of traffic should be routed to the new version.
	Percent Percentage `json:"percent"`
}
