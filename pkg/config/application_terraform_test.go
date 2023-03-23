// Copyright 2023 The PipeCD Authors.
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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestTerraformApplicationtConfig(t *testing.T) {
	testcases := []struct {
		fileName           string
		expectedKind       Kind
		expectedAPIVersion string
		expectedSpec       interface{}
		expectedError      error
	}{
		{
			fileName:           "testdata/application/terraform-app-empty.yaml",
			expectedKind:       KindTerraformApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &TerraformApplicationSpec{
				GenericApplicationSpec: GenericApplicationSpec{
					Timeout: Duration(6 * time.Hour),
					Trigger: Trigger{
						OnCommit: OnCommit{
							Disabled: false,
						},
						OnCommand: OnCommand{
							Disabled: false,
						},
						OnOutOfSync: OnOutOfSync{
							Disabled:  newBoolPointer(true),
							MinWindow: Duration(5 * time.Minute),
						},
						OnChain: OnChain{
							Disabled: newBoolPointer(true),
						},
					},
				},
				Input: TerraformDeploymentInput{},
			},
			expectedError: nil,
		},
		{
			fileName:           "testdata/application/terraform-app.yaml",
			expectedKind:       KindTerraformApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &TerraformApplicationSpec{
				GenericApplicationSpec: GenericApplicationSpec{
					Timeout: Duration(6 * time.Hour),
					Trigger: Trigger{
						OnCommit: OnCommit{
							Disabled: false,
						},
						OnCommand: OnCommand{
							Disabled: false,
						},
						OnOutOfSync: OnOutOfSync{
							Disabled:  newBoolPointer(true),
							MinWindow: Duration(5 * time.Minute),
						},
						OnChain: OnChain{
							Disabled: newBoolPointer(true),
						},
					},
				},
				Input: TerraformDeploymentInput{
					Workspace:        "dev",
					TerraformVersion: "0.12.23",
				},
			},
			expectedError: nil,
		},
		{
			fileName:           "testdata/application/terraform-app-secret-management.yaml",
			expectedKind:       KindTerraformApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &TerraformApplicationSpec{
				GenericApplicationSpec: GenericApplicationSpec{
					Timeout: Duration(6 * time.Hour),
					Trigger: Trigger{
						OnCommit: OnCommit{
							Disabled: false,
						},
						OnCommand: OnCommand{
							Disabled: false,
						},
						OnOutOfSync: OnOutOfSync{
							Disabled:  newBoolPointer(false),
							MinWindow: Duration(5 * time.Minute),
						},
						OnChain: OnChain{
							Disabled: newBoolPointer(true),
						},
					},
					Encryption: &SecretEncryption{
						EncryptedSecrets: map[string]string{
							"serviceAccount": "ENCRYPTED_DATA_GENERATED_FROM_WEB",
						},
						DecryptionTargets: []string{
							"service-account.yaml",
						},
					},
				},
				Input: TerraformDeploymentInput{
					Workspace:        "dev",
					TerraformVersion: "0.12.23",
				},
			},
			expectedError: nil,
		},
		{
			fileName:           "testdata/application/terraform-app-with-approval.yaml",
			expectedKind:       KindTerraformApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &TerraformApplicationSpec{
				GenericApplicationSpec: GenericApplicationSpec{
					Pipeline: &DeploymentPipeline{
						Stages: []PipelineStage{
							{
								Name:                      model.StageTerraformPlan,
								TerraformPlanStageOptions: &TerraformPlanStageOptions{},
							},
							{
								Name: model.StageWaitApproval,
								WaitApprovalStageOptions: &WaitApprovalStageOptions{
									Approvers:      []string{"foo", "bar"},
									Timeout:        Duration(6 * time.Hour),
									MinApproverNum: 1,
								},
							},
							{
								Name:                       model.StageTerraformApply,
								TerraformApplyStageOptions: &TerraformApplyStageOptions{},
							},
						},
					},
					Timeout: Duration(6 * time.Hour),
					Trigger: Trigger{
						OnCommit: OnCommit{
							Disabled: false,
						},
						OnCommand: OnCommand{
							Disabled: false,
						},
						OnOutOfSync: OnOutOfSync{
							Disabled:  newBoolPointer(true),
							MinWindow: Duration(5 * time.Minute),
						},
						OnChain: OnChain{
							Disabled: newBoolPointer(true),
						},
					},
				},
				Input: TerraformDeploymentInput{
					Workspace:        "dev",
					TerraformVersion: "0.12.23",
				},
			},
			expectedError: nil,
		},
		{
			fileName:           "testdata/application/terraform-app-with-exit.yaml",
			expectedKind:       KindTerraformApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &TerraformApplicationSpec{
				GenericApplicationSpec: GenericApplicationSpec{
					Pipeline: &DeploymentPipeline{
						Stages: []PipelineStage{
							{
								Name: model.StageTerraformPlan,
								TerraformPlanStageOptions: &TerraformPlanStageOptions{
									ExitOnNoChanges: true,
								},
							},
							{
								Name: model.StageWaitApproval,
								WaitApprovalStageOptions: &WaitApprovalStageOptions{
									Approvers:      []string{"foo", "bar"},
									Timeout:        Duration(6 * time.Hour),
									MinApproverNum: 1,
								},
							},
							{
								Name:                       model.StageTerraformApply,
								TerraformApplyStageOptions: &TerraformApplyStageOptions{},
							},
						},
					},
					Timeout: Duration(6 * time.Hour),
					Trigger: Trigger{
						OnCommit: OnCommit{
							Disabled: false,
						},
						OnCommand: OnCommand{
							Disabled: false,
						},
						OnOutOfSync: OnOutOfSync{
							Disabled:  newBoolPointer(true),
							MinWindow: Duration(5 * time.Minute),
						},
						OnChain: OnChain{
							Disabled: newBoolPointer(true),
						},
					},
				},
				Input: TerraformDeploymentInput{
					Workspace:        "dev",
					TerraformVersion: "0.12.23",
				},
			},
			expectedError: nil,
		},
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
