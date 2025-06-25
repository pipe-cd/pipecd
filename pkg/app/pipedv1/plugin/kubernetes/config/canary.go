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

import "github.com/pipe-cd/piped-plugin-sdk-go/unit"

// K8sCanaryRolloutStageOptions contains all configurable values for a K8S_CANARY_ROLLOUT stage.
type K8sCanaryRolloutStageOptions struct {
	// How many pods for CANARY workloads.
	// An integer value can be specified to indicate an absolute value of pod number.
	// Or a string suffixed by "%" to indicate an percentage value compared to the pod number of PRIMARY.
	// Default is 1 pod.
	Replicas unit.Replicas `json:"replicas"`
	// Suffix that should be used when naming the CANARY variant's resources.
	// Default is "canary".
	Suffix string `json:"suffix" default:"canary"`
	// Whether the CANARY service should be created.
	CreateService bool `json:"createService"`
	// List of patches used to customize manifests for CANARY variant.
	Patches []K8sResourcePatch
}

// K8sCanaryCleanStageOptions contains all configurable values for a K8S_CANARY_CLEAN stage.
type K8sCanaryCleanStageOptions struct {
}
