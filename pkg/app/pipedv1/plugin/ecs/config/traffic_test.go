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

import "testing"

func TestPercentages(t *testing.T) {
	tests := []struct {
		name        string
		opts        ECSTrafficRoutingStageOptions
		wantPrimary int
		wantCanary  int
	}{
		{
			name:        "default: neither set returns 100/0",
			opts:        ECSTrafficRoutingStageOptions{},
			wantPrimary: 100,
			wantCanary:  0,
		},
		{
			name:        "primary set to 80",
			opts:        ECSTrafficRoutingStageOptions{Primary: 80},
			wantPrimary: 80,
			wantCanary:  20,
		},
		{
			name:        "canary set to 30",
			opts:        ECSTrafficRoutingStageOptions{Canary: 30},
			wantPrimary: 70,
			wantCanary:  30,
		},
		{
			name:        "primary set to 100",
			opts:        ECSTrafficRoutingStageOptions{Primary: 100},
			wantPrimary: 100,
			wantCanary:  0,
		},
		{
			name:        "canary set to 100",
			opts:        ECSTrafficRoutingStageOptions{Canary: 100},
			wantPrimary: 0,
			wantCanary:  100,
		},
		{
			name:        "primary takes precedence when both set",
			opts:        ECSTrafficRoutingStageOptions{Primary: 60, Canary: 30},
			wantPrimary: 60,
			wantCanary:  40,
		},
		{
			name:        "primary out of range (0) falls through to canary",
			opts:        ECSTrafficRoutingStageOptions{Primary: 0, Canary: 40},
			wantPrimary: 60,
			wantCanary:  40,
		},
		{
			name:        "primary out of range (>100) falls through to canary",
			opts:        ECSTrafficRoutingStageOptions{Primary: 101, Canary: 40},
			wantPrimary: 60,
			wantCanary:  40,
		},
		{
			name:        "canary out of range (0) returns default",
			opts:        ECSTrafficRoutingStageOptions{Canary: 0},
			wantPrimary: 100,
			wantCanary:  0,
		},
		{
			name:        "canary out of range (>100) returns default",
			opts:        ECSTrafficRoutingStageOptions{Canary: 101},
			wantPrimary: 100,
			wantCanary:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrimary, gotCanary := tt.opts.Percentages()
			if gotPrimary != tt.wantPrimary || gotCanary != tt.wantCanary {
				t.Errorf("Percentages() = (%d, %d), want (%d, %d)", gotPrimary, gotCanary, tt.wantPrimary, tt.wantCanary)
			}
		})
	}
}
