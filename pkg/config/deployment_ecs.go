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

import "encoding/json"

// ECSDeploymentSpec represents a deployment configuration for ECS application.
type ECSDeploymentSpec struct {
	GenericDeploymentSpec
	// Input for ECS deployment such as where to fetch source code...
	Input ECSDeploymentInput `json:"input"`
	// Configuration for quick sync.
	QuickSync ECSSyncStageOptions `json:"quickSync"`
}

// Validate returns an error if any wrong configuration value was found.
func (s *ECSDeploymentSpec) Validate() error {
	if err := s.GenericDeploymentSpec.Validate(); err != nil {
		return err
	}
	return nil
}

type ECSDeploymentInput struct {
	// The name of service definition file placing in application directory.
	// Default is servicedef.yaml
	ServiceDefinitionFile string `json:"serviceDefinitionFile" default:"servicedef.yaml"`
	// The name of task definition file placing in application directory.
	// Default is taskdef.yaml
	TaskDefinitionFile string `json:"taskDefinitionFile" default:"taskdef.yaml"`
	// ECSTargetGroups
	TargetGroups ECSTargetGroups `json:"targetGroups"`
	// Automatically reverts all changes from all stages when one of them failed.
	// Default is true.
	AutoRollback bool `json:"autoRollback" default:"true"`
}

type ECSTargetGroups struct {
	Primary json.RawMessage `json:"primary"`
	Canary  json.RawMessage `json:"canary"`
}

// ECSSyncStageOptions contains all configurable values for a ECS_SYNC stage.
type ECSSyncStageOptions struct {
}

// ECSCanaryRolloutStageOptions contains all configurable values for a ECS_CANARY_ROLLOUT stage.
type ECSCanaryRolloutStageOptions struct {
}

// ECSPrimaryRolloutStageOptions contains all configurable values for a ECS_PRIMARY_ROLLOUT stage.
type ECSPrimaryRolloutStageOptions struct {
}

// ECSCanaryCleanStageOptions contains all configurable values for a ECS_CANARY_CLEAN stage.
type ECSCanaryCleanStageOptions struct {
}

// ECSTrafficRoutingStageOptions contains all configurable values for a ECS_TRAFFIC_ROUTING stage.
type ECSTrafficRoutingStageOptions struct {
	// The percentage of traffic should be routed to CANARY variant.
	Canary int `json:"canary"`
}
