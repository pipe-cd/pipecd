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
	"testing"
	"time"

	"github.com/pipe-cd/pipe/pkg/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLambdaApplicationConfig(t *testing.T) {
	testcases := []struct {
		fileName           string
		expectedKind       Kind
		expectedAPIVersion string
		expectedSpec       interface{}
		expectedError      error
	}{
		{
			fileName:           "testdata/application/lambda-app.yaml",
			expectedKind:       KindLambdaApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &LambdaApplicationSpec{
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
				Input: LambdaDeploymentInput{
					FunctionManifestFile: "function.yaml",
					AutoRollback:         newBoolPointer(true),
				},
			},
			expectedError: nil,
		},
		{
			fileName:           "testdata/application/lambda-app-canary.yaml",
			expectedKind:       KindLambdaApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &LambdaApplicationSpec{
				GenericApplicationSpec: GenericApplicationSpec{
					Timeout: Duration(6 * time.Hour),
					Pipeline: &DeploymentPipeline{
						Stages: []PipelineStage{
							{
								Name:                            model.StageLambdaCanaryRollout,
								LambdaCanaryRolloutStageOptions: &LambdaCanaryRolloutStageOptions{},
							},
							{
								Name: model.StageLambdaPromote,
								LambdaPromoteStageOptions: &LambdaPromoteStageOptions{
									Percent: Percentage{
										Number:    10,
										HasSuffix: false,
									},
								},
							},
							{
								Name: model.StageLambdaPromote,
								LambdaPromoteStageOptions: &LambdaPromoteStageOptions{
									Percent: Percentage{
										Number:    100,
										HasSuffix: false,
									},
								},
							},
						},
					},
					Trigger: Trigger{
						OnOutOfSync: OnOutOfSync{
							Disabled:  newBoolPointer(true),
							MinWindow: Duration(5 * time.Minute),
						},
						OnChain: OnChain{
							Disabled: newBoolPointer(true),
						},
					},
				},
				Input: LambdaDeploymentInput{
					FunctionManifestFile: "function.yaml",
					AutoRollback:         newBoolPointer(true),
				},
			},
			expectedError: nil,
		},
		{
			fileName:           "testdata/application/lambda-app-bluegreen.yaml",
			expectedKind:       KindLambdaApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &LambdaApplicationSpec{
				GenericApplicationSpec: GenericApplicationSpec{
					Timeout: Duration(6 * time.Hour),
					Pipeline: &DeploymentPipeline{
						Stages: []PipelineStage{
							{
								Name:                            model.StageLambdaCanaryRollout,
								LambdaCanaryRolloutStageOptions: &LambdaCanaryRolloutStageOptions{},
							},
							{
								Name: model.StageLambdaPromote,
								LambdaPromoteStageOptions: &LambdaPromoteStageOptions{
									Percent: Percentage{
										Number:    100,
										HasSuffix: false,
									},
								},
							},
						},
					},
					Trigger: Trigger{
						OnOutOfSync: OnOutOfSync{
							Disabled:  newBoolPointer(true),
							MinWindow: Duration(5 * time.Minute),
						},
						OnChain: OnChain{
							Disabled: newBoolPointer(true),
						},
					},
				},
				Input: LambdaDeploymentInput{
					FunctionManifestFile: "function.yaml",
					AutoRollback:         newBoolPointer(true),
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
