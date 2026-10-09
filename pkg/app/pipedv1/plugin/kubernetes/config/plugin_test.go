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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKubernetesPluginConfig_SetHelmChartInsecure(t *testing.T) {
	t.Parallel()

	cfg := &KubernetesPluginConfig{
		ChartRepositories: []HelmChartRepository{
			{Name: "insecure-repo", Address: "https://insecure.example.com", Insecure: true},
			{Name: "secure-repo", Address: "https://secure.example.com"},
		},
	}

	testcases := []struct {
		name     string
		cfg      *KubernetesPluginConfig
		chart    *InputHelmChart
		expected bool
	}{
		{
			name:     "repository configured as insecure",
			cfg:      cfg,
			chart:    &InputHelmChart{Repository: "insecure-repo", Name: "chart", Version: "1.0.0"},
			expected: true,
		},
		{
			name:     "repository configured as secure",
			cfg:      cfg,
			chart:    &InputHelmChart{Repository: "secure-repo", Name: "chart", Version: "1.0.0"},
			expected: false,
		},
		{
			name:     "repository not in the plugin config",
			cfg:      cfg,
			chart:    &InputHelmChart{Repository: "unknown-repo", Name: "chart", Version: "1.0.0"},
			expected: false,
		},
		{
			name:     "local chart",
			cfg:      cfg,
			chart:    &InputHelmChart{Path: "chart"},
			expected: false,
		},
		{
			name:     "nil plugin config",
			cfg:      nil,
			chart:    &InputHelmChart{Repository: "insecure-repo", Name: "chart", Version: "1.0.0"},
			expected: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.cfg.SetHelmChartInsecure(tc.chart)
			assert.Equal(t, tc.expected, tc.chart.Insecure)
		})
	}

	t.Run("nil chart", func(t *testing.T) {
		t.Parallel()

		assert.NotPanics(t, func() { cfg.SetHelmChartInsecure(nil) })
	})
}
