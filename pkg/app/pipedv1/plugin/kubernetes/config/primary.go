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

// K8sPrimaryRolloutStageOptions contains all configurable values for a K8S_PRIMARY_ROLLOUT stage.
type K8sPrimaryRolloutStageOptions struct {
	// Suffix that should be used when naming the PRIMARY variant's resources.
	// Default is "primary".
	Suffix string `json:"suffix" default:"primary"`
	// Whether the PRIMARY service should be created.
	CreateService bool `json:"createService"`
	// Whether the PRIMARY variant label should be added to manifests if they were missing.
	AddVariantLabelToSelector bool `json:"addVariantLabelToSelector"`
	// Whether the resources that are no longer defined in Git should be removed or not.
	Prune bool `json:"prune"`
}
