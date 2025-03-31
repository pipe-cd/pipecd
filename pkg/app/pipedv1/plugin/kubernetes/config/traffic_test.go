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

	"github.com/stretchr/testify/assert"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
)

func TestDetermineKubernetesTrafficRoutingMethod(t *testing.T) {
	tests := []struct {
		name string
		cfg  *KubernetesTrafficRouting
		want KubernetesTrafficRoutingMethod
	}{
		{
			name: "nil config should return pod selector method",
			cfg:  nil,
			want: KubernetesTrafficRoutingMethodPodSelector,
		},
		{
			name: "empty method should return pod selector method",
			cfg: &KubernetesTrafficRouting{
				Method: "",
			},
			want: KubernetesTrafficRoutingMethodPodSelector,
		},
		{
			name: "pod selector method should be returned when specified",
			cfg: &KubernetesTrafficRouting{
				Method: KubernetesTrafficRoutingMethodPodSelector,
			},
			want: KubernetesTrafficRoutingMethodPodSelector,
		},
		{
			name: "istio method should be returned when specified",
			cfg: &KubernetesTrafficRouting{
				Method: KubernetesTrafficRoutingMethodIstio,
			},
			want: KubernetesTrafficRoutingMethodIstio,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetermineKubernetesTrafficRoutingMethod(tt.cfg)
			if got != tt.want {
				t.Errorf("DetermineKubernetesTrafficRoutingMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestK8sTrafficRoutingStageOptions_Percentages(t *testing.T) {
	tests := []struct {
		name         string
		opts         K8sTrafficRoutingStageOptions
		wantPrimary  int
		wantCanary   int
		wantBaseline int
	}{
		{
			name: "all traffic to primary",
			opts: K8sTrafficRoutingStageOptions{
				All: "primary",
			},
			wantPrimary:  100,
			wantCanary:   0,
			wantBaseline: 0,
		},
		{
			name: "all traffic to canary",
			opts: K8sTrafficRoutingStageOptions{
				All: "canary",
			},
			wantPrimary:  0,
			wantCanary:   100,
			wantBaseline: 0,
		},
		{
			name: "all traffic to baseline",
			opts: K8sTrafficRoutingStageOptions{
				All: "baseline",
			},
			wantPrimary:  0,
			wantCanary:   0,
			wantBaseline: 100,
		},
		{
			name: "custom split with all percentages",
			opts: K8sTrafficRoutingStageOptions{
				Primary:  config.Percentage{Number: 50},
				Canary:   config.Percentage{Number: 30},
				Baseline: config.Percentage{Number: 20},
			},
			wantPrimary:  50,
			wantCanary:   30,
			wantBaseline: 20,
		},
		{
			name: "custom split with only primary and canary",
			opts: K8sTrafficRoutingStageOptions{
				Primary: config.Percentage{Number: 80},
				Canary:  config.Percentage{Number: 20},
			},
			wantPrimary:  80,
			wantCanary:   20,
			wantBaseline: 0,
		},
		{
			name:         "empty options should return all zeros",
			opts:         K8sTrafficRoutingStageOptions{},
			wantPrimary:  0,
			wantCanary:   0,
			wantBaseline: 0,
		},
		{
			name: "invalid 'all' value should use percentages",
			opts: K8sTrafficRoutingStageOptions{
				All:     "invalid",
				Primary: config.Percentage{Number: 60},
				Canary:  config.Percentage{Number: 40},
			},
			wantPrimary:  60,
			wantCanary:   40,
			wantBaseline: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrimary, gotCanary, gotBaseline := tt.opts.Percentages()
			assert.Equal(t, tt.wantPrimary, gotPrimary, "primary percentage")
			assert.Equal(t, tt.wantCanary, gotCanary, "canary percentage")
			assert.Equal(t, tt.wantBaseline, gotBaseline, "baseline percentage")
		})
	}
}
