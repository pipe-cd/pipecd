// Copyright 2021 The PipeCD Authors.
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

package main

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/*
var testdata embed.FS

func TestMakeCommentBody(t *testing.T) {
	testcases := []struct {
		name     string
		event    githubEvent
		result   PlanPreviewResult
		expected string
	}{
		{
			name: "no changes",
			event: githubEvent{
				HeadCommit: "abc",
			},
			result:   PlanPreviewResult{},
			expected: "testdata/comment-no-changes.txt",
		},
		{
			name: "only changed app",
			event: githubEvent{
				HeadCommit: "abc",
			},
			result: PlanPreviewResult{
				Applications: []ApplicationResult{
					{
						ApplicationInfo: ApplicationInfo{
							ApplicationID:        "app-id-1",
							ApplicationName:      "app-name-1",
							ApplicationURL:       "app-url-1",
							EnvValue:             "env-value-1",
							ApplicationKind:      "app-kind-1",
							ApplicationDirectory: "app-dir-1",
						},
						SyncStrategy: "PIPELINE",
						PlanSummary:  "plan-summary-1",
						PlanDetails:  "plan-details-1",
						NoChange:     false,
					},
				},
			},
			expected: "testdata/comment-only-changed-app.txt",
		},
		{
			name: "has no diff apps",
			event: githubEvent{
				HeadCommit: "abc",
			},
			result: PlanPreviewResult{
				Applications: []ApplicationResult{
					{
						ApplicationInfo: ApplicationInfo{
							ApplicationID:        "app-id-2",
							ApplicationName:      "app-name-2",
							ApplicationURL:       "app-url-2",
							EnvValue:             "env-value-2",
							ApplicationKind:      "app-kind-2",
							ApplicationDirectory: "app-dir-2",
						},
						SyncStrategy: "QUICK_SYNC",
						PlanSummary:  "",
						PlanDetails:  "",
						NoChange:     true,
					},
					{
						ApplicationInfo: ApplicationInfo{
							ApplicationID:        "app-id-1",
							ApplicationName:      "app-name-1",
							ApplicationURL:       "app-url-1",
							EnvValue:             "env-value-1",
							ApplicationKind:      "app-kind-1",
							ApplicationDirectory: "app-dir-1",
						},
						SyncStrategy: "PIPELINE",
						PlanSummary:  "plan-summary-1",
						PlanDetails:  "plan-details-1",
						NoChange:     false,
					},
					{
						ApplicationInfo: ApplicationInfo{
							ApplicationID:        "app-id-3",
							ApplicationName:      "app-name-3",
							ApplicationURL:       "app-url-3",
							EnvValue:             "env-value-3",
							ApplicationKind:      "app-kind-3",
							ApplicationDirectory: "app-dir-3",
						},
						SyncStrategy: "PIPELINE",
						PlanSummary:  "",
						PlanDetails:  "",
						NoChange:     true,
					},
				},
			},
			expected: "testdata/comment-has-no-diff-apps.txt",
		},
		{
			name: "no env",
			event: githubEvent{
				HeadCommit: "abc",
			},
			result: PlanPreviewResult{
				Applications: []ApplicationResult{
					{
						ApplicationInfo: ApplicationInfo{
							ApplicationID:        "app-id-1",
							ApplicationName:      "app-name-1",
							ApplicationURL:       "app-url-1",
							ApplicationKind:      "app-kind-1",
							ApplicationDirectory: "app-dir-1",
						},
						SyncStrategy: "PIPELINE",
						PlanSummary:  "plan-summary-1",
						PlanDetails:  "plan-details-1",
						NoChange:     false,
					},
				},
			},
			expected: "testdata/comment-no-env.txt",
		},
		{
			name: "has failed app",
			event: githubEvent{
				HeadCommit: "abc",
			},
			result: PlanPreviewResult{
				Applications: []ApplicationResult{
					{
						ApplicationInfo: ApplicationInfo{
							ApplicationID:        "app-id-1",
							ApplicationName:      "app-name-1",
							ApplicationURL:       "app-url-1",
							EnvValue:             "env-value-1",
							ApplicationKind:      "app-kind-1",
							ApplicationDirectory: "app-dir-1",
						},
						SyncStrategy: "PIPELINE",
						PlanSummary:  "plan-summary-1",
						PlanDetails:  "plan-details-1",
						NoChange:     false,
					},
				},
				FailureApplications: []FailureApplication{
					{
						ApplicationInfo: ApplicationInfo{
							ApplicationID:        "app-id-2",
							ApplicationName:      "app-name-2",
							ApplicationURL:       "app-url-2",
							EnvValue:             "env-value-2",
							ApplicationKind:      "app-kind-2",
							ApplicationDirectory: "app-dir-2",
						},
						Reason:      "reason-2",
						PlanDetails: "",
					},
				},
			},
			expected: "testdata/comment-has-failed-app.txt",
		},
		{
			name: "has failed piped",
			event: githubEvent{
				HeadCommit: "abc",
			},
			result: PlanPreviewResult{
				Applications: []ApplicationResult{
					{
						ApplicationInfo: ApplicationInfo{
							ApplicationID:        "app-id-1",
							ApplicationName:      "app-name-1",
							ApplicationURL:       "app-url-1",
							EnvValue:             "env-value-1",
							ApplicationKind:      "app-kind-1",
							ApplicationDirectory: "app-dir-1",
						},
						SyncStrategy: "PIPELINE",
						PlanSummary:  "plan-summary-1",
						PlanDetails:  "plan-details-1",
						NoChange:     false,
					},
				},
				FailurePipeds: []FailurePiped{
					{
						PipedInfo: PipedInfo{
							PipedID:  "piped-id-1",
							PipedURL: "piped-url-1",
						},
						Reason: "piped-reason-1",
					},
				},
			},
			expected: "testdata/comment-has-failed-piped.txt",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			expected, err := testdata.ReadFile(tc.expected)
			require.NoError(t, err)

			got := makeCommentBody(&tc.event, &tc.result)
			assert.Equal(t, string(expected), got)
		})
	}
}
