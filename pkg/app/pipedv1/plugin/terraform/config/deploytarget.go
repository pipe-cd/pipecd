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

import (
	"encoding/json"
	"fmt"

	"github.com/creasty/defaults"

	pipedconfig "github.com/pipe-cd/pipecd/pkg/configv1"
)

// TerraformDeployTargetConfig represents PipedDeployTarget.Config for Terraform plugin.
type TerraformDeployTargetConfig struct {
	// List of variables that will be set directly on terraform commands with "-var" flag.
	// The variable must be formatted by "key=value" as below:
	// "image_id=ami-abc123"
	// 'image_id_list=["ami-abc123","ami-def456"]'
	// 'image_id_map={"us-east-1":"ami-abc123","us-east-2":"ami-def456"}'
	Vars []string `json:"vars,omitempty"`
	// Enable drift detection.
	// TODO: This is a temporary option because Terraform drift detection is buggy and has performance issues. This will be possibly removed in the future release.
	DriftDetectionEnabled *bool `json:"driftDetectionEnabled" default:"true"`
}

func ParseDeployTargetConfig(deployTarget pipedconfig.PipedDeployTarget) (TerraformDeployTargetConfig, error) {
	var cfg TerraformDeployTargetConfig

	if err := json.Unmarshal(deployTarget.Config, &cfg); err != nil {
		return TerraformDeployTargetConfig{}, fmt.Errorf("failed to unmarshal deploy target configuration: %w", err)
	}

	if err := defaults.Set(&cfg); err != nil {
		return TerraformDeployTargetConfig{}, fmt.Errorf("failed to set default values for deploy target configuration: %w", err)
	}

	return cfg, nil
}
