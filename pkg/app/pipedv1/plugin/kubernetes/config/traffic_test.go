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

	"github.com/pipe-cd/piped-plugin-sdk-go/unit"
)

func TestDetermineKubernetesTrafficRoutingMethod(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			got := DetermineKubernetesTrafficRoutingMethod(tt.cfg)
			if got != tt.want {
				t.Errorf("DetermineKubernetesTrafficRoutingMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestK8sTrafficRoutingStageOptions_Percentages(t *testing.T) {
	t.Parallel()

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
				Primary:  unit.Percentage{Number: 50},
				Canary:   unit.Percentage{Number: 30},
				Baseline: unit.Percentage{Number: 20},
			},
			wantPrimary:  50,
			wantCanary:   30,
			wantBaseline: 20,
		},
		{
			name: "custom split with only primary and canary",
			opts: K8sTrafficRoutingStageOptions{
				Primary: unit.Percentage{Number: 80},
				Canary:  unit.Percentage{Number: 20},
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
				Primary: unit.Percentage{Number: 60},
				Canary:  unit.Percentage{Number: 40},
			},
			wantPrimary:  60,
			wantCanary:   40,
			wantBaseline: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotPrimary, gotCanary, gotBaseline := tt.opts.Percentages()
			assert.Equal(t, tt.wantPrimary, gotPrimary, "primary percentage")
			assert.Equal(t, tt.wantCanary, gotCanary, "canary percentage")
			assert.Equal(t, tt.wantBaseline, gotBaseline, "baseline percentage")
		})
	}
}

func TestK8sTrafficRoutingStageOptions_DisplayString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts K8sTrafficRoutingStageOptions
		want string
	}{
		{
			name: "all traffic to primary",
			opts: K8sTrafficRoutingStageOptions{
				All: "primary",
			},
			want: "Primary: 100%, Canary: 0%, Baseline: 0%",
		},
		{
			name: "all traffic to canary",
			opts: K8sTrafficRoutingStageOptions{
				All: "canary",
			},
			want: "Primary: 0%, Canary: 100%, Baseline: 0%",
		},
		{
			name: "all traffic to baseline",
			opts: K8sTrafficRoutingStageOptions{
				All: "baseline",
			},
			want: "Primary: 0%, Canary: 0%, Baseline: 100%",
		},
		{
			name: "custom split with all percentages",
			opts: K8sTrafficRoutingStageOptions{
				Primary:  unit.Percentage{Number: 50},
				Canary:   unit.Percentage{Number: 30},
				Baseline: unit.Percentage{Number: 20},
			},
			want: "Primary: 50%, Canary: 30%, Baseline: 20%",
		},
		{
			name: "custom split with only primary and canary",
			opts: K8sTrafficRoutingStageOptions{
				Primary: unit.Percentage{Number: 80},
				Canary:  unit.Percentage{Number: 20},
			},
			want: "Primary: 80%, Canary: 20%, Baseline: 0%",
		},
		{
			name: "custom split with only primary",
			opts: K8sTrafficRoutingStageOptions{
				Primary: unit.Percentage{Number: 100},
			},
			want: "Primary: 100%, Canary: 0%, Baseline: 0%",
		},
		{
			name: "custom split with only baseline",
			opts: K8sTrafficRoutingStageOptions{
				Baseline: unit.Percentage{Number: 100},
			},
			want: "Primary: 0%, Canary: 0%, Baseline: 100%",
		},
		{
			name: "empty options should return all zeros",
			opts: K8sTrafficRoutingStageOptions{},
			want: "Primary: 0%, Canary: 0%, Baseline: 0%",
		},
		{
			name: "invalid 'all' value should use percentages",
			opts: K8sTrafficRoutingStageOptions{
				All:     "invalid",
				Primary: unit.Percentage{Number: 60},
				Canary:  unit.Percentage{Number: 40},
			},
			want: "Primary: 60%, Canary: 40%, Baseline: 0%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.opts.DisplayString()
			assert.Equal(t, tt.want, got)
		})
	}
}
