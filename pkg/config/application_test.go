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

func TestHasStage(t *testing.T) {
	testcases := []struct {
		name  string
		s     GenericApplicationSpec
		stage model.Stage
		want  bool
	}{
		{
			name:  "no pipeline configured",
			s:     GenericApplicationSpec{},
			stage: model.StageK8sSync,
			want:  false,
		},
		{
			name: "given one doesn't exist",
			s: GenericApplicationSpec{
				Pipeline: &DeploymentPipeline{
					Stages: []PipelineStage{
						{
							Name: model.StageK8sSync,
						},
					},
				},
			},
			stage: model.StageK8sPrimaryRollout,
			want:  false,
		},
		{
			name: "given one exists",
			s: GenericApplicationSpec{
				Pipeline: &DeploymentPipeline{
					Stages: []PipelineStage{
						{
							Name: model.StageK8sSync,
						},
					},
				},
			},
			stage: model.StageK8sSync,
			want:  true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.s.HasStage(tc.stage)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestValidateWaitApprovalStageOptions(t *testing.T) {
	testcases := []struct {
		name           string
		minApproverNum int
		wantErr        bool
	}{
		{
			name:           "valid",
			minApproverNum: 1,
			wantErr:        false,
		},
		{
			name:           "invalid",
			minApproverNum: -1,
			wantErr:        true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			w := &WaitApprovalStageOptions{
				MinApproverNum: tc.minApproverNum,
			}
			err := w.Validate()
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestFindSlackAccounts(t *testing.T) {
	testcases := []struct {
		name     string
		mentions []NotificationMention
		event    model.NotificationEventType
		want     []string
	}{
		{
			name: "match an event name",
			mentions: []NotificationMention{
				{
					Event: "DEPLOYMENT_TRIGGERED",
					Slack: []string{"user-1", "user-2"},
				},
				{
					Event: "DEPLOYMENT_PLANNED",
					Slack: []string{"user-3", "user-4"},
				},
			},
			event: model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGERED,
			want:  []string{"user-1", "user-2"},
		},
		{
			name: "match with both event name and all-events mark",
			mentions: []NotificationMention{
				{
					Event: "DEPLOYMENT_TRIGGERED",
					Slack: []string{"user-1", "user-2"},
				},
				{
					Event: "*",
					Slack: []string{"user-1", "user-3"},
				},
			},
			event: model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGERED,
			want:  []string{"user-1", "user-2", "user-3"},
		},
		{
			name: "match by all-events mark",
			mentions: []NotificationMention{
				{
					Event: "DEPLOYMENT_TRIGGERED",
					Slack: []string{"user-1", "user-2"},
				},
				{
					Event: "*",
					Slack: []string{"user-1", "user-3"},
				},
			},
			event: model.NotificationEventType_EVENT_DEPLOYMENT_PLANNED,
			want:  []string{"user-1", "user-3"},
		},
		{
			name: "does not match anything",
			mentions: []NotificationMention{
				{
					Event: "DEPLOYMENT_TRIGGERED",
					Slack: []string{"user-1", "user-2"},
				},
			},
			event: model.NotificationEventType_EVENT_DEPLOYMENT_PLANNED,
			want:  []string{},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			n := &DeploymentNotification{
				tc.mentions,
			}
			as := n.FindSlackAccounts(tc.event)
			assert.ElementsMatch(t, tc.want, as)
		})
	}
}

func TestValidateAnalysisTemplateRef(t *testing.T) {
	testcases := []struct {
		name    string
		tplName string
		wantErr bool
	}{
		{
			name:    "valid",
			tplName: "name",
			wantErr: false,
		},
		{
			name:    "invalid due to empty template name",
			tplName: "",
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			a := &AnalysisTemplateRef{
				Name: tc.tplName,
			}
			err := a.Validate()
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestValidateEncryption(t *testing.T) {
	testcases := []struct {
		name             string
		encryptedSecrets map[string]string
		wantErr          bool
	}{
		{
			name:             "valid",
			encryptedSecrets: map[string]string{"password": "pw"},
			wantErr:          false,
		},
		{
			name:             "invalid because key is empty",
			encryptedSecrets: map[string]string{"": "pw"},
			wantErr:          true,
		},
		{
			name:             "invalid because value is empty",
			encryptedSecrets: map[string]string{"password": ""},
			wantErr:          true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := &SecretEncryption{
				EncryptedSecrets: tc.encryptedSecrets,
			}
			err := s.Validate()
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestValidateMentions(t *testing.T) {
	testcases := []struct {
		name    string
		event   string
		slack   []string
		wantErr bool
	}{
		{
			name:    "valid",
			event:   "DEPLOYMENT_TRIGGERED",
			slack:   []string{"user-1", "user-2"},
			wantErr: false,
		},
		{
			name:    "valid",
			event:   "*",
			slack:   []string{"user-1", "user-2"},
			wantErr: false,
		},
		{
			name:    "invalid because of non-existent event",
			event:   "event-1",
			slack:   []string{"user-1", "user-2"},
			wantErr: true,
		},
		{
			name:    "invalid because of missing event",
			event:   "",
			slack:   []string{"user-1", "user-2"},
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			m := &NotificationMention{
				Event: tc.event,
				Slack: tc.slack,
			}
			err := m.Validate()
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestGenericTriggerConfiguration(t *testing.T) {
	testcases := []struct {
		fileName           string
		expectedKind       Kind
		expectedAPIVersion string
		expectedSpec       interface{}
		expectedError      error
	}{
		{
			fileName:           "testdata/application/generic-trigger.yaml",
			expectedKind:       KindKubernetesApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &KubernetesApplicationSpec{
				GenericApplicationSpec: GenericApplicationSpec{
					Timeout: Duration(6 * time.Hour),
					Trigger: Trigger{
						OnCommit: OnCommit{
							Disabled: false,
							Paths: []string{
								"deployment.yaml",
							},
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
				Input: KubernetesDeploymentInput{
					AutoRollback: newBoolPointer(true),
				},
				VariantLabel: KubernetesVariantLabel{
					Key:           "pipecd.dev/variant",
					PrimaryValue:  "primary",
					BaselineValue: "baseline",
					CanaryValue:   "canary",
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

func TestTrueByDefaultBoolConfiguration(t *testing.T) {
	testcases := []struct {
		fileName           string
		expectedKind       Kind
		expectedAPIVersion string
		expectedSpec       interface{}
		expectedError      error
	}{
		{
			fileName:           "testdata/application/truebydefaultbool-not-specified.yaml",
			expectedKind:       KindKubernetesApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &KubernetesApplicationSpec{
				GenericApplicationSpec: GenericApplicationSpec{
					Timeout: Duration(6 * time.Hour),
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
				Input: KubernetesDeploymentInput{
					AutoRollback: newBoolPointer(true),
				},
				VariantLabel: KubernetesVariantLabel{
					Key:           "pipecd.dev/variant",
					PrimaryValue:  "primary",
					BaselineValue: "baseline",
					CanaryValue:   "canary",
				},
			},
			expectedError: nil,
		},
		{
			fileName:           "testdata/application/truebydefaultbool-false-explicitly.yaml",
			expectedKind:       KindKubernetesApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &KubernetesApplicationSpec{
				GenericApplicationSpec: GenericApplicationSpec{
					Timeout: Duration(6 * time.Hour),
					Trigger: Trigger{
						OnOutOfSync: OnOutOfSync{
							Disabled:  newBoolPointer(false),
							MinWindow: Duration(5 * time.Minute),
						},
						OnChain: OnChain{
							Disabled: newBoolPointer(true),
						},
					},
				},
				Input: KubernetesDeploymentInput{
					AutoRollback: newBoolPointer(false),
				},
				VariantLabel: KubernetesVariantLabel{
					Key:           "pipecd.dev/variant",
					PrimaryValue:  "primary",
					BaselineValue: "baseline",
					CanaryValue:   "canary",
				},
			},
			expectedError: nil,
		},
		{
			fileName:           "testdata/application/truebydefaultbool-true-explicitly.yaml",
			expectedKind:       KindKubernetesApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &KubernetesApplicationSpec{
				GenericApplicationSpec: GenericApplicationSpec{
					Timeout: Duration(6 * time.Hour),
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
				Input: KubernetesDeploymentInput{
					AutoRollback: newBoolPointer(true),
				},
				VariantLabel: KubernetesVariantLabel{
					Key:           "pipecd.dev/variant",
					PrimaryValue:  "primary",
					BaselineValue: "baseline",
					CanaryValue:   "canary",
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

func TestGenericPostSyncConfiguration(t *testing.T) {
	testcases := []struct {
		fileName           string
		expectedKind       Kind
		expectedAPIVersion string
		expectedSpec       interface{}
		expectedError      error
	}{
		{
			fileName:           "testdata/application/generic-postsync.yaml",
			expectedKind:       KindKubernetesApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &KubernetesApplicationSpec{
				GenericApplicationSpec: GenericApplicationSpec{
					Timeout: Duration(6 * time.Hour),
					Trigger: Trigger{
						OnOutOfSync: OnOutOfSync{
							Disabled:  newBoolPointer(true),
							MinWindow: Duration(5 * time.Minute),
						},
						OnChain: OnChain{
							Disabled: newBoolPointer(true),
						},
					},
					PostSync: &PostSync{
						DeploymentChain: &DeploymentChain{
							ApplicationMatchers: []ChainApplicationMatcher{
								{
									Name: "app-1",
								},
								{
									Labels: map[string]string{
										"env": "staging",
										"foo": "bar",
									},
								},
								{
									Kind: "ECSApp",
								},
							},
						},
					},
				},
				Input: KubernetesDeploymentInput{
					AutoRollback: newBoolPointer(true),
				},
				VariantLabel: KubernetesVariantLabel{
					Key:           "pipecd.dev/variant",
					PrimaryValue:  "primary",
					BaselineValue: "baseline",
					CanaryValue:   "canary",
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

func TestCustomSyncConfig(t *testing.T) {
	testcases := []struct {
		fileName           string
		expectedKind       Kind
		expectedAPIVersion string
		expectedSpec       interface{}
		expectedError      error
	}{
		{
			fileName:           "testdata/application/custom-stage.yaml",
			expectedKind:       KindLambdaApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &LambdaApplicationSpec{
				GenericApplicationSpec: GenericApplicationSpec{
					Timeout: Duration(6 * time.Hour),
					Pipeline: &DeploymentPipeline{
						Stages: []PipelineStage{
							{
								Name: model.StageCustomSync,
								Desc: "deploy by sam",
								CustomSyncOptions: &CustomSyncOptions{
									Timeout: Duration(6 * time.Hour),
									Envs: map[string]string{
										"AWS_PROFILE": "default",
									},
									Run: "sam build\nsam deploy -g --profile $AWS_PROFILE\n",
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
