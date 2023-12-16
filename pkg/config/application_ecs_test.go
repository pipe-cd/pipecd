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
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestECSApplicationConfig(t *testing.T) {
	testcases := []struct {
		fileName                                   string
		expectedKind                               Kind
		expectedAPIVersion                         string
		expectedLaunchType                         string
		expectedELBListenerRuleSelectorIsSpecified bool
		expectedSpec                               interface{}
		expectedError                              error
	}{
		{
			fileName:           "testdata/application/ecs-app.yaml",
			expectedKind:       KindECSApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedELBListenerRuleSelectorIsSpecified: false,
			expectedSpec: &ECSApplicationSpec{
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
				Input: ECSDeploymentInput{
					ServiceDefinitionFile: "/path/to/servicedef.yaml",
					TaskDefinitionFile:    "/path/to/taskdef.yaml",
					TargetGroups: ECSTargetGroups{
						Primary: json.RawMessage(`{"containerName":"web","containerPort":80,"targetGroupArn":"arn:aws:elasticloadbalancing:xyz"}`),
					},
					LaunchType:        "FARGATE",
					AutoRollback:      newBoolPointer(true),
					RunStandaloneTask: newBoolPointer(true),
					AccessType:        "ELB",
				},
			},
			expectedError: nil,
		},
		{
			fileName:           "testdata/application/ecs-app-service-discovery.yaml",
			expectedKind:       KindECSApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedELBListenerRuleSelectorIsSpecified: false,
			expectedSpec: &ECSApplicationSpec{
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
				Input: ECSDeploymentInput{
					ServiceDefinitionFile: "/path/to/servicedef.yaml",
					TaskDefinitionFile:    "/path/to/taskdef.yaml",
					LaunchType:            "FARGATE",
					AutoRollback:          newBoolPointer(true),
					RunStandaloneTask:     newBoolPointer(true),
					AccessType:            "SERVICE_DISCOVERY",
				},
			},
			expectedError: nil,
		},
		{
			fileName:           "testdata/application/ecs-app-invalid-access-type.yaml",
			expectedKind:       KindECSApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedELBListenerRuleSelectorIsSpecified: false,
			expectedSpec: &ECSApplicationSpec{
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
				Input: ECSDeploymentInput{
					ServiceDefinitionFile: "/path/to/servicedef.yaml",
					TaskDefinitionFile:    "/path/to/taskdef.yaml",
					LaunchType:            "FARGATE",
					AutoRollback:          newBoolPointer(true),
					RunStandaloneTask:     newBoolPointer(true),
					AccessType:            "XXX",
				},
			},
			expectedError: fmt.Errorf("invalid accessType: XXX"),
		},
		{
			fileName:           "testdata/application/ecs-app-elb-selector.yaml",
			expectedKind:       KindECSApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedELBListenerRuleSelectorIsSpecified: true,
			expectedSpec: &ECSApplicationSpec{
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
				Input: ECSDeploymentInput{
					ServiceDefinitionFile: "/path/to/servicedef.yaml",
					TaskDefinitionFile:    "/path/to/taskdef.yaml",
					TargetGroups: ECSTargetGroups{
						Primary: json.RawMessage(`{"containerName":"web","containerPort":80,"targetGroupArn":"arn:aws:elasticloadbalancing:xyz"}`),
					},
					LaunchType:        "FARGATE",
					AutoRollback:      newBoolPointer(true),
					RunStandaloneTask: newBoolPointer(true),
					AccessType:        "ELB",
					ListenerRuleSelector: ELBListenerRuleSelector{
						ListenerRuleArn: "arn:aws:elasticloadbalancing:ap-northeast-1:<account-id>:listener-rule/app/<elb-name>/xxx/yyy/zzz",
					},
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
				assert.Equal(t, tc.expectedELBListenerRuleSelectorIsSpecified, cfg.ECSApplicationSpec.Input.ListenerRuleSelector.IsSpecified())
			}
		})
	}
}
