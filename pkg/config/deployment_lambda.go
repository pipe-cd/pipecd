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

// LambdaDeploymentSpec represents a deployment configuration for Lambda application.
type LambdaDeploymentSpec struct {
	Input     LambdaDeploymentInput  `json:"input"`
	QuickSync LambdaSyncStageOptions `json:"quickSync"`
	Pipeline  *DeploymentPipeline    `json:"pipeline"`
}

func (s *LambdaDeploymentSpec) GetStage(index int32) (PipelineStage, bool) {
	if s.Pipeline == nil {
		return PipelineStage{}, false
	}
	if int(index) >= len(s.Pipeline.Stages) {
		return PipelineStage{}, false
	}
	return s.Pipeline.Stages[index], true
}

// Validate returns an error if any wrong configuration value was found.
func (s *LambdaDeploymentSpec) Validate() error {
	return nil
}

type LambdaDeploymentInput struct {
	Git  string `json:"git"`
	Path string `json:"path"`
	Ref  string `json:"ref"`
	// Automatically reverts all changes from all stages when one of them failed.
	// Default is true.
	AutoRollback bool `json:"autoRollback"`
}

// LambdaSyncStageOptions contains all configurable values for a CLOUDRUN_SYNC stage.
type LambdaSyncStageOptions struct {
}

// LambdaCanaryRolloutStageOptions contains all configurable values for a CLOUDRUN_CANARY_ROLLOUT stage.
type LambdaCanaryRolloutStageOptions struct {
}

// LambdaTrafficRoutingStageOptions contains all configurable values for a CLOUDRUN_TRAFFIC_ROUTING stage.
type LambdaTrafficRoutingStageOptions struct {
	// Which variant should receive all traffic.
	// This can be either "primary" or "canary".
	All string `json:"all"`
	// The percentage of traffic should be routed to PRIMARY variant.
	Primary int `json:"primary"`
	// The percentage of traffic should be routed to CANARY variant.
	Canary int `json:"canary"`
}

func (opts LambdaTrafficRoutingStageOptions) Percentages() (primary, canary int) {
	switch opts.All {
	case "primary":
		primary = 100
		return
	case "canary":
		canary = 100
		return
	}
	return opts.Primary, opts.Canary
}
