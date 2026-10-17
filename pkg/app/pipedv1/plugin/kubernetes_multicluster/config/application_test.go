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
	"testing"

	"github.com/stretchr/testify/assert"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
)

func TestFindDeployTarget(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *config.PipedPlugin
		targetName  string
		expected    KubernetesDeployTargetConfig
		expectedErr bool
	}{
		{
			name:        "nil config",
			cfg:         nil,
			targetName:  "target",
			expected:    KubernetesDeployTargetConfig{},
			expectedErr: true,
		},
		{
			name: "missing deploy target",
			cfg: &config.PipedPlugin{
				DeployTargets: []config.PipedDeployTarget{},
			},
			targetName:  "target",
			expected:    KubernetesDeployTargetConfig{},
			expectedErr: true,
		},
		{
			name: "valid deploy target",
			cfg: &config.PipedPlugin{
				DeployTargets: []config.PipedDeployTarget{
					{
						Name: "target",
						Config: json.RawMessage(`{
							"masterURL": "https://example.com",
							"kubeConfigPath": "/path/to/kubeconfig",
							"kubectlVersion": "v1.20.0"
						}`),
					},
				},
			},
			targetName: "target",
			expected: KubernetesDeployTargetConfig{
				MasterURL:      "https://example.com",
				KubeConfigPath: "/path/to/kubeconfig",
				KubectlVersion: "v1.20.0",
			},
			expectedErr: false,
		},
		{
			name: "invalid deploy target config",
			cfg: &config.PipedPlugin{
				DeployTargets: []config.PipedDeployTarget{
					{
						Name:   "target",
						Config: json.RawMessage(`invalid`),
					},
				},
			},
			targetName:  "target",
			expected:    KubernetesDeployTargetConfig{},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FindDeployTarget(tt.cfg, tt.targetName)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestKubernetesApplicationSpecValidate(t *testing.T) {
	tests := []struct {
		name        string
		spec        KubernetesApplicationSpec
		expectedErr string
	}{
		{
			name: "valid empty spec",
			spec: KubernetesApplicationSpec{},
		},
		{
			name: "valid single multiTarget",
			spec: KubernetesApplicationSpec{
				Input: KubernetesDeploymentInput{
					MultiTargets: []KubernetesMultiTarget{
						{Target: KubernetesMultiTargetDeployTarget{Name: "cluster-us"}},
					},
				},
			},
		},
		{
			name: "multiTarget with empty name",
			spec: KubernetesApplicationSpec{
				Input: KubernetesDeploymentInput{
					MultiTargets: []KubernetesMultiTarget{
						{Target: KubernetesMultiTargetDeployTarget{Name: "cluster-us"}},
						{Target: KubernetesMultiTargetDeployTarget{Name: ""}},
					},
				},
			},
			expectedErr: "multiTargets[1].target.name must not be empty",
		},
		{
			name: "helmChart and kustomizeOptions are mutually exclusive",
			spec: KubernetesApplicationSpec{
				Input: KubernetesDeploymentInput{
					HelmChart:        &InputHelmChart{Path: "./charts/myapp"},
					KustomizeOptions: map[string]string{"enable-alpha-plugins": "true"},
				},
			},
			expectedErr: "helmChart and kustomizeOptions are mutually exclusive",
		},
		{
			name: "helmChart without kustomizeOptions is valid",
			spec: KubernetesApplicationSpec{
				Input: KubernetesDeploymentInput{
					HelmChart: &InputHelmChart{Path: "./charts/myapp"},
				},
			},
		},
		{
			name: "kustomizeOptions without helmChart is valid",
			spec: KubernetesApplicationSpec{
				Input: KubernetesDeploymentInput{
					KustomizeOptions: map[string]string{"enable-alpha-plugins": "true"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.spec.Validate()
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
