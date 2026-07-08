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

package config

import (
	"encoding/json"

	"github.com/creasty/defaults"
	"github.com/pipe-cd/piped-plugin-sdk-go/unit"
)

type CloudRunApplicationSpec struct {
	// Input for CloudRun deployment such as docker image...
	Input CloudRunDeploymentInput `json:"input"`
	// Configuration for quick sync.
	QuickSync CloudRunSyncStageOptions `json:"quickSync"`
}

func (s *CloudRunApplicationSpec) UnmarshalJSON(data []byte) error {
	type alias CloudRunApplicationSpec

	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	*s = CloudRunApplicationSpec(a)
	if err := defaults.Set(s); err != nil {
		return err
	}

	return nil
}

type CloudRunDeploymentInput struct {
	// The name of service manifest file placing in application directory.
	// Default is service.yaml
	ServiceManifestFile string `json:"serviceManifestFile" default:"service.yaml"`
}

// CloudRunSyncStageOptions contains all configurable values for a CLOUDRUN_SYNC stage.
type CloudRunSyncStageOptions struct {
}

// CloudRunPromoteStageOptions contains all configurable values for a CLOUDRUN_PROMOTE stage.
type CloudRunPromoteStageOptions struct {
	// Percentage of traffic should be routed to the new version.
	Percent unit.Percentage `json:"percent"`
}

type CloudRunDeployTargetConfig struct {
	// The GCP project hosting the CloudRun service.
	Project string `json:"project"`
	// The region of running CloudRun service.
	Region string `json:"region"`
	// The path to the service account file for accessing CloudRun service.
	CredentialsFile string `json:"credentialsFile,omitempty"`
}
