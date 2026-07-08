// Copyright 2026 The PipeCD Authors.
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

// ECSTrafficRoutingStageOptions contains all configurable values for an ECS_TRAFFIC_ROUTING stage.
type ECSTrafficRoutingStageOptions struct {
	// Canary represents the percentage of traffic to route to the canary variant.
	// If set, primary will be 100 - canary.
	Canary int `json:"canary,omitempty"`
	// Primary represents the percentage of traffic to route to the primary variant.
	// If set, canary will be 100 - primary.
	Primary int `json:"primary,omitempty"`
}

// Percentages returns the traffic split between primary and canary.
// If neither is set, primary gets 100% by default.
func (opts ECSTrafficRoutingStageOptions) Percentages() (primary, canary int) {
	if opts.Primary > 0 && opts.Primary <= 100 {
		return opts.Primary, 100 - opts.Primary
	}
	if opts.Canary > 0 && opts.Canary <= 100 {
		return 100 - opts.Canary, opts.Canary
	}
	return 100, 0
}
