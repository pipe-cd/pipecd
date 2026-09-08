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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKubernetesApplicationSpecValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		spec    KubernetesApplicationSpec
		wantErr string
	}{
		{
			name: "valid spec",
			spec: validKubernetesApplicationSpec(),
		},
		{
			name: "helmChart and kustomizeOptions are mutually exclusive",
			spec: KubernetesApplicationSpec{
				Input: KubernetesDeploymentInput{
					HelmChart:        &InputHelmChart{Path: "./chart"},
					KustomizeOptions: map[string]string{"load-restrictor": "LoadRestrictionsNone"},
				},
				VariantLabel: validKubernetesVariantLabel(),
			},
			wantErr: "helmChart and kustomizeOptions are mutually exclusive",
		},
		{
			name: "unsupported traffic routing method",
			spec: KubernetesApplicationSpec{
				VariantLabel: validKubernetesVariantLabel(),
				TrafficRouting: &KubernetesTrafficRouting{
					Method: "unknown",
				},
			},
			wantErr: `unsupported trafficRouting.method "unknown"`,
		},
		{
			name: "empty variant label key",
			spec: KubernetesApplicationSpec{
				VariantLabel: KubernetesVariantLabel{
					PrimaryValue:  "primary",
					CanaryValue:   "canary",
					BaselineValue: "baseline",
				},
			},
			wantErr: "variantLabel.key must not be empty",
		},
		{
			name: "empty primary variant label value",
			spec: KubernetesApplicationSpec{
				VariantLabel: KubernetesVariantLabel{
					Key:           "pipecd.dev/variant",
					CanaryValue:   "canary",
					BaselineValue: "baseline",
				},
			},
			wantErr: "variantLabel.primaryValue must not be empty",
		},
		{
			name: "empty canary variant label value",
			spec: KubernetesApplicationSpec{
				VariantLabel: KubernetesVariantLabel{
					Key:           "pipecd.dev/variant",
					PrimaryValue:  "primary",
					BaselineValue: "baseline",
				},
			},
			wantErr: "variantLabel.canaryValue must not be empty",
		},
		{
			name: "empty baseline variant label value",
			spec: KubernetesApplicationSpec{
				VariantLabel: KubernetesVariantLabel{
					Key:          "pipecd.dev/variant",
					PrimaryValue: "primary",
					CanaryValue:  "canary",
				},
			},
			wantErr: "variantLabel.baselineValue must not be empty",
		},
		{
			name: "duplicate variant label values",
			spec: KubernetesApplicationSpec{
				VariantLabel: KubernetesVariantLabel{
					Key:           "pipecd.dev/variant",
					PrimaryValue:  "stable",
					CanaryValue:   "stable",
					BaselineValue: "baseline",
				},
			},
			wantErr: "variantLabel primaryValue, canaryValue, and baselineValue must be unique",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.spec.Validate()
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestKubernetesApplicationSpecUnmarshalJSONSetsDefaults(t *testing.T) {
	t.Parallel()

	var spec KubernetesApplicationSpec
	err := json.Unmarshal([]byte(`{}`), &spec)

	assert.NoError(t, err)
	assert.NoError(t, spec.Validate())
	assert.Equal(t, validKubernetesVariantLabel(), spec.VariantLabel)
}

func validKubernetesApplicationSpec() KubernetesApplicationSpec {
	return KubernetesApplicationSpec{
		VariantLabel: validKubernetesVariantLabel(),
		TrafficRouting: &KubernetesTrafficRouting{
			Method: KubernetesTrafficRoutingMethodPodSelector,
		},
	}
}

func validKubernetesVariantLabel() KubernetesVariantLabel {
	return KubernetesVariantLabel{
		Key:           "pipecd.dev/variant",
		PrimaryValue:  "primary",
		CanaryValue:   "canary",
		BaselineValue: "baseline",
	}
}
