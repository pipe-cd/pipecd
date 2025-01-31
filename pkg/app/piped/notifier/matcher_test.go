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

package notifier

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestMatch(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name      string
		config    config.NotificationRoute
		matchings map[model.NotificationEvent]bool
	}{
		{
			name:   "empty config",
			config: config.NotificationRoute{},
			matchings: map[model.NotificationEvent]bool{
				{}: true,
				{Type: model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGERED}: true,
			},
		},
		{
			name: "filter by event",
			config: config.NotificationRoute{
				Events: []string{
					"DEPLOYMENT_TRIGGERED",
				},
				IgnoreEvents: []string{
					"DEPLOYMENT_ROLLING_BACK",
				},
			},
			matchings: map[model.NotificationEvent]bool{
				{
					Type: model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGERED,
				}: true,
				{
					Type: model.NotificationEventType_EVENT_DEPLOYMENT_ROLLING_BACK,
				}: false,
			},
		},
		{
			name: "filter by group",
			config: config.NotificationRoute{
				Groups: []string{
					"DEPLOYMENT",
				},
				IgnoreGroups: []string{
					"APPLICATION",
				},
			},
			matchings: map[model.NotificationEvent]bool{
				{
					Type: model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGERED,
				}: true,
				{
					Type: model.NotificationEventType_EVENT_APPLICATION_SYNCED,
				}: false,
			},
		},
		{
			name: "filter by app",
			config: config.NotificationRoute{
				Apps: []string{
					"canary",
				},
				IgnoreApps: []string{
					"bluegreen",
				},
			},
			matchings: map[model.NotificationEvent]bool{
				{
					Type: model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGERED,
					Metadata: &model.NotificationEventDeploymentTriggered{
						Deployment: &model.Deployment{
							ApplicationName: "canary",
						},
					},
				}: true,
				{
					Type: model.NotificationEventType_EVENT_DEPLOYMENT_PLANNED,
					Metadata: &model.NotificationEventDeploymentTriggered{
						Deployment: &model.Deployment{
							ApplicationName: "bluegreen",
						},
					},
				}: false,
				{
					Type: model.NotificationEventType_EVENT_DEPLOYMENT_SUCCEEDED,
					Metadata: &model.NotificationEventDeploymentTriggered{
						Deployment: &model.Deployment{
							ApplicationName: "not-specified",
						},
					},
				}: false,
				{
					Type:     model.NotificationEventType_EVENT_PIPED_STARTED,
					Metadata: &model.NotificationEventPipedStarted{},
				}: true,
				{
					Type: model.NotificationEventType_EVENT_DEPLOYMENT_CANCELLED,
					Metadata: &model.NotificationEventDeploymentCancelled{
						Deployment: &model.Deployment{
							ApplicationName: "bluegreen",
						},
					},
				}: false,
				{
					Type: model.NotificationEventType_EVENT_DEPLOYMENT_WAIT_APPROVAL,
					Metadata: &model.NotificationEventDeploymentWaitApproval{
						Deployment: &model.Deployment{
							ApplicationName: "bluegreen",
						},
					},
				}: false,
				{
					Type: model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGER_FAILED,
					Metadata: &model.NotificationEventDeploymentTriggerFailed{
						Application: &model.Application{
							Name: "bluegreen",
						},
					},
				}: false,
				{
					Type: model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGER_FAILED,
					Metadata: &model.NotificationEventDeploymentTriggerFailed{
						Application: &model.Application{
							Name: "canary",
						},
					},
				}: true,
				{
					Type: model.NotificationEventType_EVENT_DEPLOYMENT_STARTED,
					Metadata: &model.NotificationEventDeploymentStarted{
						Deployment: &model.Deployment{
							ApplicationName: "canary",
						},
					},
				}: true,
			},
		},
		{
			name: "filter by labels",
			config: config.NotificationRoute{
				Labels: map[string]string{
					"env":  "dev",
					"team": "pipecd",
				},
				IgnoreLabels: map[string]string{
					"env":  "local",
					"team": "not-pipecd",
				},
			},
			matchings: map[model.NotificationEvent]bool{
				{
					Type: model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGERED,
					Metadata: &model.NotificationEventDeploymentTriggered{
						Deployment: &model.Deployment{
							Labels: map[string]string{
								"team":    "pipecd",
								"env":     "dev",
								"project": "pipecd",
							},
						},
					},
				}: true,
				{
					Type: model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGERED,
					Metadata: &model.NotificationEventDeploymentTriggered{
						Deployment: &model.Deployment{},
					},
				}: false,
				{
					Type: model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGERED,
					Metadata: &model.NotificationEventDeploymentTriggered{
						Deployment: &model.Deployment{
							Labels: map[string]string{
								"team": "pipecd",
								"env":  "prod",
							},
						},
					},
				}: false,
				{
					Type: model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGERED,
					Metadata: &model.NotificationEventDeploymentTriggered{
						Deployment: &model.Deployment{
							Labels: map[string]string{
								"env": "stg",
							},
						},
					},
				}: false,
				{
					Type: model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGERED,
					Metadata: &model.NotificationEventDeploymentTriggered{
						Deployment: &model.Deployment{
							Labels: map[string]string{
								"env":  "local",
								"team": "pipecd",
							},
						},
					},
				}: false,
				{
					Type:     model.NotificationEventType_EVENT_PIPED_STARTED,
					Metadata: &model.NotificationEventPipedStarted{},
				}: true,
				{
					Type: model.NotificationEventType_EVENT_DEPLOYMENT_CANCELLED,
					Metadata: &model.NotificationEventDeploymentCancelled{
						Deployment: &model.Deployment{
							Labels: map[string]string{
								"env":  "stg",
								"team": "pipecd",
							},
						},
					},
				}: false,
				{
					Type: model.NotificationEventType_EVENT_DEPLOYMENT_WAIT_APPROVAL,
					Metadata: &model.NotificationEventDeploymentWaitApproval{
						Deployment: &model.Deployment{
							Labels: map[string]string{
								"env":  "stg",
								"team": "pipecd",
							},
						},
					},
				}: false,
				{
					Type: model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGER_FAILED,
					Metadata: &model.NotificationEventDeploymentTriggerFailed{
						Application: &model.Application{
							Labels: map[string]string{
								"env":  "stg",
								"team": "pipecd",
							},
						},
					},
				}: false,
				{
					Type: model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGER_FAILED,
					Metadata: &model.NotificationEventDeploymentTriggerFailed{
						Application: &model.Application{
							Labels: map[string]string{
								"env":  "dev",
								"team": "pipecd",
							},
						},
					},
				}: true,
				{
					Type: model.NotificationEventType_EVENT_DEPLOYMENT_STARTED,
					Metadata: &model.NotificationEventDeploymentStarted{
						Deployment: &model.Deployment{
							Labels: map[string]string{
								"env":  "stg",
								"team": "pipecd",
							},
						},
					},
				}: false,
				{
					Type: model.NotificationEventType_EVENT_DEPLOYMENT_STARTED,
					Metadata: &model.NotificationEventDeploymentStarted{
						Deployment: &model.Deployment{
							Labels: map[string]string{
								"env":  "dev",
								"team": "pipecd",
							},
						},
					},
				}: true,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			matcher := newMatcher(tc.config)
			for event, expected := range tc.matchings {
				got := matcher.Match(event)
				assert.Equal(t, expected, got, event.Type.String())
			}
		})
	}
}
