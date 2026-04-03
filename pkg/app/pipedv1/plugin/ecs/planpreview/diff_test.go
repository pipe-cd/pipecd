// Copyright 2026 The PipeCD Authors.
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
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiffDefinitions(t *testing.T) {
	taskWithImage := func(image string) *types.TaskDefinition {
		return &types.TaskDefinition{
			Family: aws.String("my-task"),
			ContainerDefinitions: []types.ContainerDefinition{
				{Name: aws.String("app"), Image: aws.String(image)},
			},
		}
	}

	tests := []struct {
		name           string
		old            *types.TaskDefinition
		new            *types.TaskDefinition
		wantEmpty      bool
		wantContains   []string
		wantNoRemovals bool
	}{
		{
			name:      "no change",
			old:       taskWithImage("nginx:1.0"),
			new:       taskWithImage("nginx:1.0"),
			wantEmpty: true,
		},
		{
			name: "image tag changed",
			old:  taskWithImage("nginx:1.0"),
			new:  taskWithImage("nginx:2.0"),
			wantContains: []string{
				"nginx:1.0",
				"nginx:2.0",
				"taskdef (running)",
				"taskdef (target)",
			},
		},
		{
			name:           "nil old (first deployment)",
			old:            nil,
			new:            &types.TaskDefinition{Family: aws.String("my-task")},
			wantNoRemovals: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diff, err := diffDefinitions(tc.old, tc.new, "taskdef")
			require.NoError(t, err)

			if tc.wantEmpty {
				assert.Empty(t, diff)
				return
			}

			assert.NotEmpty(t, diff)

			for _, s := range tc.wantContains {
				assert.Contains(t, diff, s)
			}

			if tc.wantNoRemovals {
				for _, line := range strings.Split(diff, "\n") {
					if strings.HasPrefix(line, "---") {
						continue
					}
					assert.False(t, strings.HasPrefix(line, "-"), "unexpected removed line: %q", line)
				}
			}

			t.Logf("diff: %+v", diff)
		})
	}
}

func TestBuildSummary(t *testing.T) {
	tests := []struct {
		name        string
		taskDiff    string
		serviceDiff string
		want        string
	}{
		{
			name:        "only task def changed",
			taskDiff:    "some diff",
			serviceDiff: "",
			want:        "task definition changed",
		},
		{
			name:        "only service def changed",
			taskDiff:    "",
			serviceDiff: "some diff",
			want:        "service definition changed",
		},
		{
			name:        "both changed",
			taskDiff:    "some diff",
			serviceDiff: "some diff",
			want:        "task definition changed, service definition changed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, buildSummary(tc.taskDiff, tc.serviceDiff))
		})
	}
}

func TestBuildDetails(t *testing.T) {
	tests := []struct {
		name        string
		taskDiff    string
		serviceDiff string
		want        []byte
	}{
		{
			name:        "no changes",
			taskDiff:    "",
			serviceDiff: "",
			want:        nil,
		},
		{
			name:        "task diff only",
			taskDiff:    "task-diff\n",
			serviceDiff: "",
			want:        []byte("task-diff\n"),
		},
		{
			name:        "service diff only",
			taskDiff:    "",
			serviceDiff: "service-diff\n",
			want:        []byte("service-diff\n"),
		},
		{
			name:        "both diffs combined with separator",
			taskDiff:    "task-diff\n",
			serviceDiff: "service-diff\n",
			want:        []byte("task-diff\n\nservice-diff\n"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, buildDetails(tc.taskDiff, tc.serviceDiff))
		})
	}
}

func TestToResponse(t *testing.T) {
	tests := []struct {
		name         string
		deployTarget string
		taskDiff     string
		serviceDiff  string
		want         sdk.PlanPreviewResult
	}{
		{
			name:         "no changes",
			deployTarget: "prod",
			taskDiff:     "",
			serviceDiff:  "",
			want: sdk.PlanPreviewResult{
				DeployTarget: "prod",
				NoChange:     true,
				Summary:      "No changes were detected",
				Details:      nil,
				DiffLanguage: "diff",
			},
		},
		{
			name:         "task def changed",
			deployTarget: "prod",
			taskDiff:     "task-diff\n",
			serviceDiff:  "",
			want: sdk.PlanPreviewResult{
				DeployTarget: "prod",
				NoChange:     false,
				Summary:      "task definition changed",
				Details:      []byte("task-diff\n"),
				DiffLanguage: "diff",
			},
		},
		{
			name:         "both changed",
			deployTarget: "prod",
			taskDiff:     "task-diff\n",
			serviceDiff:  "service-diff\n",
			want: sdk.PlanPreviewResult{
				DeployTarget: "prod",
				NoChange:     false,
				Summary:      "task definition changed, service definition changed",
				Details:      []byte("task-diff\n\nservice-diff\n"),
				DiffLanguage: "diff",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := toResponse(tc.deployTarget, tc.taskDiff, tc.serviceDiff)
			require.Len(t, resp.Results, 1)
			assert.Equal(t, tc.want, resp.Results[0])
		})
	}
}
