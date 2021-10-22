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

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipe/pkg/model"
)

func TestHasStage(t *testing.T) {
	testcases := []struct {
		name  string
		s     GenericDeploymentSpec
		stage model.Stage
		want  bool
	}{
		{
			name:  "no pipeline configured",
			s:     GenericDeploymentSpec{},
			stage: model.StageK8sSync,
			want:  false,
		},
		{
			name: "given one doesn't exist",
			s: GenericDeploymentSpec{
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
			s: GenericDeploymentSpec{
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

func TestFindSlackAccounts(t *testing.T) {
	testcases := []struct {
		name string
		mentions []NotificationMention
		event model.NotificationEventType
		want []string
	}{
		{
			name: "the applicable event",
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
			want: []string{"user-1", "user-2"},
		},
		{
			name: "the applicable event and all events",
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
			want: []string{"user-1", "user-2", "user-3"},
		},
		{
			name: "not the appricable event and all events",
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
			want: []string{"user-1", "user-3"},
		},
		{
			name: "not the appricable event",
			mentions: []NotificationMention{
				{
					Event: "DEPLOYMENT_TRIGGERED",
					Slack: []string{"user-1", "user-2"},
				},
			},
			event: model.NotificationEventType_EVENT_DEPLOYMENT_PLANNED,
			want: []string{},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			n := &DeploymentNotification {
				tc.mentions,
			}
			as := n.FindSlackAccounts(tc.event)
			assert.ElementsMatch(t, tc.want, as)
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
