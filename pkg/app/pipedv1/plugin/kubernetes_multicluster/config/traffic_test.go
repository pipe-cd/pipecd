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
	"testing"

	"github.com/pipe-cd/piped-plugin-sdk-go/unit"
	"github.com/stretchr/testify/assert"
)

func TestDetermineKubernetesTrafficRoutingMethod(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		cfg      *KubernetesTrafficRouting
		expected KubernetesTrafficRoutingMethod
	}{
		{
			name:     "nil config returns PodSelector",
			cfg:      nil,
			expected: KubernetesTrafficRoutingMethodPodSelector,
		},
		{
			name:     "empty Method returns PodSelector",
			cfg:      &KubernetesTrafficRouting{Method: ""},
			expected: KubernetesTrafficRoutingMethodPodSelector,
		},
		{
			name:     "explicit podselector returns PodSelector",
			cfg:      &KubernetesTrafficRouting{Method: KubernetesTrafficRoutingMethodPodSelector},
			expected: KubernetesTrafficRoutingMethodPodSelector,
		},
		{
			name:     "explicit istio returns Istio",
			cfg:      &KubernetesTrafficRouting{Method: KubernetesTrafficRoutingMethodIstio},
			expected: KubernetesTrafficRoutingMethodIstio,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := DetermineKubernetesTrafficRoutingMethod(tt.cfg)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestK8sTrafficRoutingStageOptions_Percentages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		opts             K8sTrafficRoutingStageOptions
		expectedPrimary  int
		expectedCanary   int
		expectedBaseline int
	}{
		{
			name:             "all=primary returns 100, 0, 0",
			opts:             K8sTrafficRoutingStageOptions{All: "primary"},
			expectedPrimary:  100,
			expectedCanary:   0,
			expectedBaseline: 0,
		},
		{
			name:             "all=canary returns 0, 100, 0",
			opts:             K8sTrafficRoutingStageOptions{All: "canary"},
			expectedPrimary:  0,
			expectedCanary:   100,
			expectedBaseline: 0,
		},
		{
			name:             "all=baseline returns 0, 0, 100",
			opts:             K8sTrafficRoutingStageOptions{All: "baseline"},
			expectedPrimary:  0,
			expectedCanary:   0,
			expectedBaseline: 100,
		},
		{
			name: "explicit percentages are returned as-is",
			opts: K8sTrafficRoutingStageOptions{
				Primary:  unit.Percentage{Number: 60},
				Canary:   unit.Percentage{Number: 40},
				Baseline: unit.Percentage{Number: 0},
			},
			expectedPrimary:  60,
			expectedCanary:   40,
			expectedBaseline: 0,
		},
		{
			name:             "zero value returns 0, 0, 0",
			opts:             K8sTrafficRoutingStageOptions{},
			expectedPrimary:  0,
			expectedCanary:   0,
			expectedBaseline: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			primary, canary, baseline := tt.opts.Percentages()
			assert.Equal(t, tt.expectedPrimary, primary)
			assert.Equal(t, tt.expectedCanary, canary)
			assert.Equal(t, tt.expectedBaseline, baseline)
		})
	}
}

func TestK8sTrafficRoutingStageOptions_DisplayString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		opts     K8sTrafficRoutingStageOptions
		expected string
	}{
		{
			name:     "all=primary",
			opts:     K8sTrafficRoutingStageOptions{All: "primary"},
			expected: "Primary: 100%, Canary: 0%, Baseline: 0%",
		},
		{
			name: "explicit split",
			opts: K8sTrafficRoutingStageOptions{
				Primary:  unit.Percentage{Number: 60},
				Canary:   unit.Percentage{Number: 40},
				Baseline: unit.Percentage{Number: 0},
			},
			expected: "Primary: 60%, Canary: 40%, Baseline: 0%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.opts.DisplayString()
			assert.Equal(t, tt.expected, got)
		})
	}
}
