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

package planpreview

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestReadableResultString(t *testing.T) {
	testcases := []struct {
		name     string
		results  []*model.PlanPreviewCommandResult
		expected string
	}{
		{
			name:    "empty",
			results: []*model.PlanPreviewCommandResult{},
			expected: `
There are no updated applications. It means no deployment will be triggered once this pull request got merged.
`,
		},
		{
			name: "there is only a plannable application",
			results: []*model.PlanPreviewCommandResult{
				{
					CommandId: "command-2",
					PipedId:   "piped-2",
					PipedUrl:  "https://pipecd.dev/piped-2",
					Results: []*model.ApplicationPlanPreviewResult{
						{
							ApplicationId:   "app-1",
							ApplicationName: "app-1",
							ApplicationUrl:  "https://pipecd.dev/app-1",
							ApplicationKind: model.ApplicationKind_KUBERNETES,
							Labels:          map[string]string{"env": "env-1"},
							SyncStrategy:    model.SyncStrategy_QUICK_SYNC,
							PlanSummary:     []byte("2 manifests will be added, 1 manifest will be deleted and 5 manifests will be changed"),
							PlanDetails:     []byte("changes-1"),
						},
					},
				},
			},
			expected: `
Here are plan-preview for 1 application:

1. app: app-1, env: env-1, kind: KUBERNETES
  sync strategy: QUICK_SYNC
  summary: 2 manifests will be added, 1 manifest will be deleted and 5 manifests will be changed
  details:

  ---DETAILS_BEGIN---
changes-1
  ---DETAILS_END---
`,
		},
		{
			name: "there is only a failure application",
			results: []*model.PlanPreviewCommandResult{
				{
					CommandId: "command-2",
					PipedId:   "piped-2",
					PipedUrl:  "https://pipecd.dev/piped-2",
					Results: []*model.ApplicationPlanPreviewResult{
						{
							ApplicationId:   "app-2",
							ApplicationName: "app-2",
							ApplicationUrl:  "https://pipecd.dev/app-2",
							ApplicationKind: model.ApplicationKind_TERRAFORM,
							Labels:          map[string]string{"env": "env-2"},
							Error:           "wrong application configuration",
						},
					},
				},
			},
			expected: `
NOTE: An error occurred while building plan-preview for the following application:

1. app: app-2, env: env-2, kind: TERRAFORM
  reason: wrong application configuration
`,
		},
		{
			name: "there is only a failure piped",
			results: []*model.PlanPreviewCommandResult{
				{
					CommandId: "command-1",
					PipedId:   "piped-1",
					PipedUrl:  "https://pipecd.dev/piped-1",
					Error:     "failed to clone",
				},
			},
			expected: `
NOTE: An error occurred while building plan-preview for applications of the following Piped:

1. piped: piped-1
  reason: failed to clone
`,
		},
		{
			name: "all kinds",
			results: []*model.PlanPreviewCommandResult{
				{
					CommandId: "command-1",
					PipedId:   "piped-1",
					PipedUrl:  "https://pipecd.dev/piped-1",
					Error:     "failed to clone",
				},
				{
					CommandId: "command-2",
					PipedId:   "piped-2",
					PipedUrl:  "https://pipecd.dev/piped-2",
					Results: []*model.ApplicationPlanPreviewResult{
						{
							ApplicationId:   "app-1",
							ApplicationName: "app-1",
							ApplicationUrl:  "https://pipecd.dev/app-1",
							ApplicationKind: model.ApplicationKind_KUBERNETES,
							Labels:          map[string]string{"env": "env-1"},
							SyncStrategy:    model.SyncStrategy_QUICK_SYNC,
							PlanSummary:     []byte("2 manifests will be added, 1 manifest will be deleted and 5 manifests will be changed"),
							PlanDetails:     []byte("changes-1"),
						},
						{
							ApplicationId:   "app-2",
							ApplicationName: "app-2",
							ApplicationUrl:  "https://pipecd.dev/app-2",
							ApplicationKind: model.ApplicationKind_TERRAFORM,
							Labels:          map[string]string{"env": "env-2"},
							SyncStrategy:    model.SyncStrategy_PIPELINE,
							PlanSummary:     []byte("1 to add, 2 to change, 0 to destroy"),
							PlanDetails:     []byte("changes-2"),
						},
						{
							ApplicationId:   "app-3",
							ApplicationName: "app-3",
							ApplicationUrl:  "https://pipecd.dev/app-3",
							ApplicationKind: model.ApplicationKind_TERRAFORM,
							Labels:          map[string]string{"env": "env-3"},
							Error:           "wrong application configuration",
						},
						{
							ApplicationId:   "app-4",
							ApplicationName: "app-4",
							ApplicationUrl:  "https://pipecd.dev/app-4",
							ApplicationKind: model.ApplicationKind_CLOUDRUN,
							Error:           "missing key",
						},
					},
				},
				{
					CommandId: "command-3",
					PipedId:   "piped-3",
					PipedUrl:  "https://pipecd.dev/piped-3",
					Error:     "failed to checkout branch",
				},
			},
			expected: `
Here are plan-preview for 2 applications:

1. app: app-1, env: env-1, kind: KUBERNETES
  sync strategy: QUICK_SYNC
  summary: 2 manifests will be added, 1 manifest will be deleted and 5 manifests will be changed
  details:

  ---DETAILS_BEGIN---
changes-1
  ---DETAILS_END---

2. app: app-2, env: env-2, kind: TERRAFORM
  sync strategy: PIPELINE
  summary: 1 to add, 2 to change, 0 to destroy
  details:

  ---DETAILS_BEGIN---
changes-2
  ---DETAILS_END---

NOTE: An error occurred while building plan-preview for the following 2 applications:

1. app: app-3, env: env-3, kind: TERRAFORM
  reason: wrong application configuration

2. app: app-4, kind: CLOUDRUN
  reason: missing key

NOTE: An error occurred while building plan-preview for applications of the following 2 Pipeds:

1. piped: piped-1
  reason: failed to clone

2. piped: piped-3
  reason: failed to checkout branch
`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			printResults(tc.results, &buf, "")

			assert.Equal(t, tc.expected, buf.String())
		})
	}
}
