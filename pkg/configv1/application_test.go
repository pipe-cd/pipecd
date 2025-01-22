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
	"testing"
	"time"

	"github.com/creasty/defaults"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func newBoolPointer(b bool) *bool {
	return &b
}

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

func TestFindSlackUsersAndGroups(t *testing.T) {
	testcases := []struct {
		name       string
		mentions   []NotificationMention
		event      model.NotificationEventType
		wantUsers  []string
		wantGroups []string
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
			event:     model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGERED,
			wantUsers: []string{"user-1", "user-2"},
		},
		{
			name: "match with both event name and all-events mark",
			mentions: []NotificationMention{
				{
					Event: "DEPLOYMENT_TRIGGERED",
					Slack: []string{"user-1", "user-2"},
				},
				{
					Event:       "*",
					Slack:       []string{"user-1", "user-3"},
					SlackGroups: []string{"group-1", "group-2"},
				},
			},
			event:      model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGERED,
			wantUsers:  []string{"user-1", "user-2", "user-3"},
			wantGroups: []string{"group-1", "group-2"},
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
			event:     model.NotificationEventType_EVENT_DEPLOYMENT_PLANNED,
			wantUsers: []string{"user-1", "user-3"},
		},
		{
			name: "match by all-events mark with slack groups",
			mentions: []NotificationMention{
				{
					Event: "DEPLOYMENT_TRIGGERED",
					Slack: []string{"user-1", "user-2"},
				},
				{
					Event:       "*",
					SlackGroups: []string{"group-1", "group-2"},
				},
			},
			event:      model.NotificationEventType_EVENT_DEPLOYMENT_PLANNED,
			wantGroups: []string{"group-1", "group-2"},
		},
		{
			name: "does not match anything",
			mentions: []NotificationMention{
				{
					Event: "DEPLOYMENT_TRIGGERED",
					Slack: []string{"user-1", "user-2"},
				},
			},
			event:     model.NotificationEventType_EVENT_DEPLOYMENT_PLANNED,
			wantUsers: []string{},
		},
		{
			name: "match an event name with Slack Groups",
			mentions: []NotificationMention{
				{
					Event:       "DEPLOYMENT_PLANNED",
					SlackGroups: []string{"group-1", "group-2"},
				},
			},
			event:      model.NotificationEventType_EVENT_DEPLOYMENT_PLANNED,
			wantGroups: []string{"group-1", "group-2"},
		},
		{
			name: "match an event name with Slack Users and Groups",
			mentions: []NotificationMention{
				{
					Event:       "DEPLOYMENT_PLANNED",
					Slack:       []string{"user-1", "user-2"},
					SlackGroups: []string{"group-1", "group-2"},
				},
			},
			event:      model.NotificationEventType_EVENT_DEPLOYMENT_PLANNED,
			wantUsers:  []string{"user-1", "user-2"},
			wantGroups: []string{"group-1", "group-2"},
		},
		{
			name: "match an event name with Slack Users with new field SlackUsers",
			mentions: []NotificationMention{
				{
					Event:       "DEPLOYMENT_PLANNED",
					SlackUsers:  []string{"user-1", "user-2"},
					Slack:       []string{"user-3", "user-4"},
					SlackGroups: []string{"group-1", "group-2"},
				},
			},
			event:      model.NotificationEventType_EVENT_DEPLOYMENT_PLANNED,
			wantUsers:  []string{"user-1", "user-2", "user-3", "user-4"},
			wantGroups: []string{"group-1", "group-2"},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			n := &DeploymentNotification{
				tc.mentions,
			}
			as := n.FindSlackUsers(tc.event)
			ag := n.FindSlackGroups(tc.event)
			assert.ElementsMatch(t, tc.wantUsers, as)
			assert.ElementsMatch(t, tc.wantGroups, ag)
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
		targets          []string
		wantErr          bool
	}{
		{
			name:             "valid",
			encryptedSecrets: map[string]string{"password": "pw"},
			targets:          []string{"secret.yaml"},
			wantErr:          false,
		},
		{
			name:             "invalid because key is empty",
			encryptedSecrets: map[string]string{"": "pw"},
			targets:          []string{"secret.yaml"},
			wantErr:          true,
		},
		{
			name:             "invalid because value is empty",
			encryptedSecrets: map[string]string{"password": ""},
			targets:          []string{"secret.yaml"},
			wantErr:          true,
		},
		{
			name:             "no target files sepcified",
			encryptedSecrets: map[string]string{"password": "pw"},
			wantErr:          true,
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			s := &SecretEncryption{
				EncryptedSecrets:  tc.encryptedSecrets,
				DecryptionTargets: tc.targets,
			}
			err := s.Validate()
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestValidateAttachment(t *testing.T) {
	testcases := []struct {
		name    string
		sources map[string]string
		targets []string
		wantErr bool
	}{
		{
			name:    "valid",
			sources: map[string]string{"config": "config.yaml"},
			targets: []string{"target.yaml"},
			wantErr: false,
		},
		{
			name:    "invalid because key is empty",
			sources: map[string]string{"": "config-data"},
			targets: []string{"target.yaml"},
			wantErr: true,
		},
		{
			name:    "invalid because value is empty",
			sources: map[string]string{"config": ""},
			targets: []string{"target.yaml"},
			wantErr: true,
		},
		{
			name:    "no target files sepcified",
			sources: map[string]string{"config": "config.yaml"},
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			a := &Attachment{
				Sources: tc.sources,
				Targets: tc.targets,
			}
			err := a.Validate()
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
			expectedSpec: &GenericApplicationSpec{
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
				Planner: DeploymentPlanner{
					AutoRollback: newBoolPointer(true),
				},
				Pipeline: &DeploymentPipeline{},
			},
			expectedError: nil,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.fileName, func(t *testing.T) {
			cfg, err := LoadFromYAML[*GenericApplicationSpec](tc.fileName)
			require.Equal(t, tc.expectedError, err)
			if err == nil {
				assert.Equal(t, tc.expectedKind, cfg.Kind)
				assert.Equal(t, tc.expectedAPIVersion, cfg.APIVersion)
				assert.Equal(t, tc.expectedSpec, cfg.Spec)
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
			expectedSpec: &GenericApplicationSpec{
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
				Planner: DeploymentPlanner{
					AutoRollback: newBoolPointer(true),
				},
				Pipeline: &DeploymentPipeline{},
			},
			expectedError: nil,
		},
		{
			fileName:           "testdata/application/truebydefaultbool-false-explicitly.yaml",
			expectedKind:       KindKubernetesApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &GenericApplicationSpec{
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
				Planner: DeploymentPlanner{
					AutoRollback: newBoolPointer(true),
				},
				Pipeline: &DeploymentPipeline{},
			},
			expectedError: nil,
		},
		{
			fileName:           "testdata/application/truebydefaultbool-true-explicitly.yaml",
			expectedKind:       KindKubernetesApp,
			expectedAPIVersion: "pipecd.dev/v1beta1",
			expectedSpec: &GenericApplicationSpec{
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
				Planner: DeploymentPlanner{
					AutoRollback: newBoolPointer(true),
				},
				Pipeline: &DeploymentPipeline{},
			},
			expectedError: nil,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.fileName, func(t *testing.T) {
			cfg, err := LoadFromYAML[*GenericApplicationSpec](tc.fileName)
			require.Equal(t, tc.expectedError, err)
			if err == nil {
				assert.Equal(t, tc.expectedKind, cfg.Kind)
				assert.Equal(t, tc.expectedAPIVersion, cfg.APIVersion)
				assert.Equal(t, tc.expectedSpec, cfg.Spec)
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
			expectedSpec: &GenericApplicationSpec{
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
				Planner: DeploymentPlanner{
					AutoRollback: newBoolPointer(true),
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
				Pipeline: &DeploymentPipeline{},
			},
			expectedError: nil,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.fileName, func(t *testing.T) {
			cfg, err := LoadFromYAML[*GenericApplicationSpec](tc.fileName)
			require.Equal(t, tc.expectedError, err)
			if err == nil {
				assert.Equal(t, tc.expectedKind, cfg.Kind)
				assert.Equal(t, tc.expectedAPIVersion, cfg.APIVersion)
				assert.Equal(t, tc.expectedSpec, cfg.Spec)
			}
		})
	}
}
func TestGetStageByte(t *testing.T) {
	testcases := []struct {
		name   string
		s      GenericApplicationSpec
		index  int32
		want   []byte
		wantOk bool
	}{
		{
			name:   "pipeline not defined",
			s:      GenericApplicationSpec{},
			index:  0,
			want:   nil,
			wantOk: true,
		},
		{
			name: "valid stage index",
			s: GenericApplicationSpec{
				Pipeline: &DeploymentPipeline{
					Stages: []PipelineStage{
						{
							Name: model.StageK8sSync,
						},
					},
				},
			},
			index:  0,
			want:   []byte(`{"id":"","name":"K8S_SYNC","timeout":"0s","with":null}`),
			wantOk: true,
		},
		{
			name: "invalid stage index",
			s: GenericApplicationSpec{
				Pipeline: &DeploymentPipeline{
					Stages: []PipelineStage{
						{
							Name: model.StageK8sSync,
						},
					},
				},
			},
			index:  1,
			want:   nil,
			wantOk: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			defaults.Set(&tc.s)
			got, ok := tc.s.GetStageByte(tc.index)
			assert.Equal(t, tc.wantOk, ok)
			assert.Equal(t, string(tc.want), string(got))
		})
	}
}
