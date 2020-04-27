// Copyright 2020 The PipeCD Authors.
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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppConfig(t *testing.T) {
	testcases := []struct {
		fileName           string
		expectedKind       Kind
		expectedAPIVersion string
		expectedSpec       interface{}
		expectedError      error
	}{
		{
			fileName:      "testdata/application/terraform-app-missing-destination.yaml",
			expectedError: fmt.Errorf("spec.destination for terraform application is required"),
		},
		{
			fileName:           "testdata/application/terraform-app.yaml",
			expectedKind:       KindTerraformApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &TerraformAppSpec{
				Input: &TerraformAppInput{
					Workspace:        "dev",
					TerraformVersion: "0.12.23",
				},
				Pipeline:    nil,
				Destination: "dev-terraform",
			},
			expectedError: nil,
		},
		{
			fileName:           "testdata/application/terraform-app-with-approval.yaml",
			expectedKind:       KindTerraformApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &TerraformAppSpec{
				Input: &TerraformAppInput{
					Workspace:        "dev",
					TerraformVersion: "0.12.23",
				},
				Pipeline: &AppPipeline{
					Stages: []PipelineStage{
						PipelineStage{
							Name:                      StageTerraformPlan,
							TerraformPlanStageOptions: &TerraformPlanStageOptions{},
						},
						PipelineStage{
							Name: StageWaitApproval,
							WaitApprovalStageOptions: &WaitApprovalStageOptions{
								Approvers: []string{"foo", "bar"},
							},
						},
						PipelineStage{
							Name:                       StageTerraformApply,
							TerraformApplyStageOptions: &TerraformApplyStageOptions{},
						},
					},
				},
				Destination: "dev-terraform",
			},
			expectedError: nil,
		},
		// {
		// 	fileName:           "testdata/application/k8s-app-helm.yaml",
		// 	expectedKind:       KindK8sHelmApp,
		// 	expectedAPIVersion: "pipecd.dev/v1beta1",
		// 	expectedSpec: &K8sHelmAppSpec{
		// 		Input: &K8sHelmAppInput{
		// 			Chart:       "git@github.com:org/config-repo.git:charts/demoapp?ref=v1.0.0",
		// 			ValueFiles:  []string{"values.yaml"},
		// 			HelmVersion: "3.1.1",
		// 		},
		// 		Pipeline: nil,
		// 	},
		// 	expectedError: nil,
		// },
		// {
		// 	fileName:           "testdata/application/k8s-app-canary.yaml",
		// 	expectedKind:       KindKubernetesApp,
		// 	expectedAPIVersion: "pipecd.dev/v1beta1",
		// 	expectedSpec: &K8sAppSpec{
		// 		Pipeline: &AppPipeline{
		// 			Stages: []PipelineStage{
		// 				PipelineStage{
		// 					Name: StageK8sCanaryOut,
		// 					K8sCanaryOutStageOptions: &K8sCanaryOutStageOptions{
		// 						Weight: 10,
		// 					},
		// 					Timeout:   Duration(10 * time.Minute),
		// 					PostDelay: Duration(time.Minute),
		// 				},
		// 				PipelineStage{
		// 					Name: StageWaitApproval,
		// 					ApprovalStageOptions: &ApprovalStageOptions{
		// 						Approvers: []string{"foo", "bar"},
		// 					},
		// 				},
		// 				PipelineStage{
		// 					Name:                   StageK8sRollout,
		// 					K8sRolloutStageOptions: &K8sRolloutStageOptions{},
		// 				},
		// 				PipelineStage{
		// 					Name:                    StageK8sCanaryIn,
		// 					K8sCanaryInStageOptions: &K8sCanaryInStageOptions{},
		// 				},
		// 			},
		// 		},
		// 	},
		// 	expectedError: nil,
		// },
	}
	for _, tc := range testcases {
		t.Run(tc.fileName, func(t *testing.T) {
			cfg, err := LoadFromYAML(tc.fileName)
			require.Equal(t, tc.expectedError, err)
			if err == nil {
				assert.Equal(t, tc.expectedKind, cfg.Kind)
				assert.Equal(t, tc.expectedAPIVersion, cfg.APIVersion)
				assert.Equal(t, tc.expectedSpec, cfg.spec)
			}
		})
	}
}
