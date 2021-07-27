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
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			result: PlanPreviewResult{},
			expected: `<!-- pipecd-plan-preview-->
[![PLAN_PREVIEW](https://img.shields.io/static/v1?label=PipeCD&message=Plan_Preview&color=success&style=flat)](https://pipecd.dev/docs/user-guide/plan-preview/)

Ran plan-preview against head commit abc of this pull request. PipeCD detected ` + "`0`" + ` updated application. It means no deployment will be triggered once this pull request got merged.
`,
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
							EnvID:                "env-id-1",
							EnvURL:               "env-url-1",
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
			expected: `<!-- pipecd-plan-preview-->
[![PLAN_PREVIEW](https://img.shields.io/static/v1?label=PipeCD&message=Plan_Preview&color=success&style=flat)](https://pipecd.dev/docs/user-guide/plan-preview/)

Ran plan-preview against head commit abc of this pull request. PipeCD detected ` + "`1`" + ` updated applications and here are their plan results. Once this pull request got merged their deployments will be triggered to run as these estimations.

## app: [app-name-1](app-url-1), env: [](env-url-1), kind: app-kind-1
Sync strategy: PIPELINE
Summary: plan-summary-1

<details>
<summary>Details (Click me)</summary>
<p>

` + "```" + ` diff
plan-details-1
` + "```" + `
</p>
</details>
`,
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
							EnvID:                "env-id-2",
							EnvURL:               "env-url-2",
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
							EnvID:                "env-id-1",
							EnvURL:               "env-url-1",
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
							EnvID:                "env-id-3",
							EnvURL:               "env-url-3",
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
			expected: `<!-- pipecd-plan-preview-->
[![PLAN_PREVIEW](https://img.shields.io/static/v1?label=PipeCD&message=Plan_Preview&color=success&style=flat)](https://pipecd.dev/docs/user-guide/plan-preview/)

Ran plan-preview against head commit abc of this pull request. PipeCD detected ` + "`3`" + ` updated applications and here are their plan results. Once this pull request got merged their deployments will be triggered to run as these estimations.

## app: [app-name-1](app-url-1), env: [](env-url-1), kind: app-kind-1
Sync strategy: PIPELINE
Summary: plan-summary-1

<details>
<summary>Details (Click me)</summary>
<p>

` + "```" + ` diff
plan-details-1
` + "```" + `
</p>
</details>

## No resource changes were detected but the following apps will also be triggered

### ` + "`PIPELINE`" + `

- app: [app-name-3](app-url-3), env: [](env-url-3), kind: app-kind-3

### ` + "`QUICK_SYNC`" + `

- app: [app-name-2](app-url-2), env: [](env-url-2), kind: app-kind-2
`,
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
							EnvID:                "env-id-1",
							EnvURL:               "env-url-1",
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
							EnvID:                "env-id-2",
							EnvURL:               "env-url-2",
							ApplicationKind:      "app-kind-2",
							ApplicationDirectory: "app-dir-2",
						},
						Reason:      "reason-2",
						PlanDetails: "",
					},
				},
			},
			expected: `<!-- pipecd-plan-preview-->
[![PLAN_PREVIEW](https://img.shields.io/static/v1?label=PipeCD&message=Plan_Preview&color=orange&style=flat)](https://pipecd.dev/docs/user-guide/plan-preview/)

Ran plan-preview against head commit abc of this pull request. PipeCD detected ` + "`1`" + ` updated applications and here are their plan results. Once this pull request got merged their deployments will be triggered to run as these estimations.

## app: [app-name-1](app-url-1), env: [](env-url-1), kind: app-kind-1
Sync strategy: PIPELINE
Summary: plan-summary-1

<details>
<summary>Details (Click me)</summary>
<p>

` + "```" + ` diff
plan-details-1
` + "```" + `
</p>
</details>

---

## NOTE

**An error occurred while building plan-preview for the following applications**

## app: [app-name-2](app-url-2), env: [](env-url-2), kind: app-kind-2
Reason: reason-2

<details>
<summary>Details (Click me)</summary>
<p>

` + "```" + ` diff

` + "```" + `
</p>
</details>
`,
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
							EnvID:                "env-id-1",
							EnvURL:               "env-url-1",
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
			expected: `<!-- pipecd-plan-preview-->
[![PLAN_PREVIEW](https://img.shields.io/static/v1?label=PipeCD&message=Plan_Preview&color=orange&style=flat)](https://pipecd.dev/docs/user-guide/plan-preview/)

Ran plan-preview against head commit abc of this pull request. PipeCD detected ` + "`1`" + ` updated applications and here are their plan results. Once this pull request got merged their deployments will be triggered to run as these estimations.

## app: [app-name-1](app-url-1), env: [](env-url-1), kind: app-kind-1
Sync strategy: PIPELINE
Summary: plan-summary-1

<details>
<summary>Details (Click me)</summary>
<p>

` + "```" + ` diff
plan-details-1
` + "```" + `
</p>
</details>

---

## NOTE

**An error occurred while building plan-preview for applications of the following Pipeds**

## piped: [piped-id-1](piped-url-1)
Reason: piped-reason-1

`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := makeCommentBody(&tc.event, &tc.result)
			assert.Equal(t, tc.expected, got)
		})
	}
}
