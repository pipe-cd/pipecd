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

// CloudRunApplicationSpec represents an application configuration for CloudRun application.
type CloudRunApplicationSpec struct {
	GenericApplicationSpec
	// Input for CloudRun deployment such as docker image...
	Input CloudRunDeploymentInput `json:"input"`
	// Configuration for quick sync.
	QuickSync CloudRunSyncStageOptions `json:"quickSync"`
}

// Validate returns an error if any wrong configuration value was found.
func (s *CloudRunApplicationSpec) Validate() error {
	if err := s.GenericApplicationSpec.Validate(); err != nil {
		return err
	}
	return nil
}

type CloudRunDeploymentInput struct {
	// The name of service manifest file placing in application directory.
	// Default is service.yaml
	ServiceManifestFile string `json:"serviceManifestFile"`
	// Automatically reverts to the previous state when the deployment is failed.
	// Default is true.
	AutoRollback *bool `json:"autoRollback,omitempty" default:"true"`
}

// CloudRunSyncStageOptions contains all configurable values for a CLOUDRUN_SYNC stage.
type CloudRunSyncStageOptions struct {
}

// CloudRunPromoteStageOptions contains all configurable values for a CLOUDRUN_PROMOTE stage.
type CloudRunPromoteStageOptions struct {
	// Percentage of traffic should be routed to the new version.
	Percent Percentage `json:"percent"`
}
