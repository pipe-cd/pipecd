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

// CloudRunDeploymentSpec represents a deployment configuration for CloudRun application.
type CloudRunDeploymentSpec struct {
	// Input for CloudRun deployment such as docker image...
	Input CloudRunDeploymentInput `json:"input"`
	// Configuration for quick sync.
	QuickSync CloudRunSyncStageOptions `json:"quickSync"`
	// Pipline for deploying progressively.
	Pipeline *DeploymentPipeline `json:"pipeline"`
}

func (s *CloudRunDeploymentSpec) GetStage(index int32) (PipelineStage, bool) {
	if s.Pipeline == nil {
		return PipelineStage{}, false
	}
	if int(index) >= len(s.Pipeline.Stages) {
		return PipelineStage{}, false
	}
	return s.Pipeline.Stages[index], true
}

// Validate returns an error if any wrong configuration value was found.
func (s *CloudRunDeploymentSpec) Validate() error {
	return nil
}

type CloudRunDeploymentInput struct {
	// The name of service manifest file placing in application configuration directory.
	// Default is service.yaml
	ServiceManifestFile string `json:"serviceManifestFile"`
	// Automatically reverts all changes from all stages when one of them failed.
	// Default is true.
	AutoRollback bool `json:"autoRollback"`
}

// CloudRunSyncStageOptions contains all configurable values for a CLOUDRUN_SYNC stage.
type CloudRunSyncStageOptions struct {
}

// CloudRunPromoteStageOptions contains all configurable values for a CLOUDRUN_PROMOTE stage.
type CloudRunPromoteStageOptions struct {
	Percent int `json:"percent"`
}
