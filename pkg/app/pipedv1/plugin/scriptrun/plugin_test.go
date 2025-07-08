// Copyright 2025 The PipeCD Authors.
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
	"encoding/json"
	"testing"
	"time"

	"github.com/pipe-cd/piped-plugin-sdk-go/unit"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestBuildPipelineSyncStages(t *testing.T) {
	t.Parallel()
	p := &plugin{}
	ctx := t.Context()

	testcases := []struct {
		name     string
		input    *sdk.BuildPipelineSyncStagesInput
		expected *sdk.BuildPipelineSyncStagesResponse
	}{
		{
			name: "should generate one stage if onRollback is not defined",
			input: &sdk.BuildPipelineSyncStagesInput{
				Request: sdk.BuildPipelineSyncStagesRequest{
					Stages: []sdk.StageConfig{
						{
							Index: 0,
							Name:  stageScriptRun,
							Config: mustMarshalJSON(t, &scriptRunStageOptions{
								Timeout: unit.Duration(1 * time.Minute),
								Run:     "echo 'hello'",
							}),
						},
					},
				},
			},
			expected: &sdk.BuildPipelineSyncStagesResponse{
				Stages: []sdk.PipelineStage{
					{Index: 0, Name: stageScriptRun, Rollback: false, Metadata: map[string]string{}},
				},
			},
		},
		{
			name: "should generate one stage if onRollback is not defined even if request.Rollback is true",
			input: &sdk.BuildPipelineSyncStagesInput{
				Request: sdk.BuildPipelineSyncStagesRequest{
					Stages: []sdk.StageConfig{
						{
							Index: 0,
							Name:  stageScriptRun,
							Config: mustMarshalJSON(t, &scriptRunStageOptions{
								Timeout: unit.Duration(1 * time.Minute),
								Run:     "echo 'hello'",
							}),
						},
					},
					Rollback: true,
				},
			},
			expected: &sdk.BuildPipelineSyncStagesResponse{
				Stages: []sdk.PipelineStage{
					{Index: 0, Name: stageScriptRun, Rollback: false, Metadata: map[string]string{}},
				},
			},
		},
		{
			name: "should generate two stages if onRollback is defined",
			input: &sdk.BuildPipelineSyncStagesInput{
				Request: sdk.BuildPipelineSyncStagesRequest{
					Stages: []sdk.StageConfig{
						{
							Index: 0,
							Name:  stageScriptRun,
							Config: mustMarshalJSON(t, &scriptRunStageOptions{
								Timeout:    unit.Duration(1 * time.Minute),
								Run:        "echo 'hello'",
								OnRollback: "echo 'rollback'",
							}),
						},
					},
				},
			},
			expected: &sdk.BuildPipelineSyncStagesResponse{
				Stages: []sdk.PipelineStage{
					{Index: 0, Name: stageScriptRun, Rollback: false, Metadata: map[string]string{}},
					{Index: 0, Name: stageScriptRunRollback, Rollback: true, Metadata: map[string]string{}},
				},
			},
		},
		{
			name: "index should map 1 to 1 for rollback stage",
			input: &sdk.BuildPipelineSyncStagesInput{
				Request: sdk.BuildPipelineSyncStagesRequest{
					Stages: []sdk.StageConfig{
						{
							Index: 0,
							Name:  stageScriptRun,
							Config: mustMarshalJSON(t, &scriptRunStageOptions{
								Timeout:    unit.Duration(1 * time.Minute),
								Run:        "echo 'hello 0'",
								OnRollback: "echo 'rollback 0'",
							}),
						},
						{
							Index: 2,
							Name:  stageScriptRun,
							Config: mustMarshalJSON(t, &scriptRunStageOptions{
								Timeout:    unit.Duration(1 * time.Minute),
								Run:        "echo 'hello 2'",
								OnRollback: "echo 'rollback 2'",
							}),
						},
					},
				},
			},
			expected: &sdk.BuildPipelineSyncStagesResponse{
				Stages: []sdk.PipelineStage{
					{Index: 0, Name: stageScriptRun, Rollback: false, Metadata: map[string]string{}},
					{Index: 0, Name: stageScriptRunRollback, Rollback: true, Metadata: map[string]string{}},
					{Index: 2, Name: stageScriptRun, Rollback: false, Metadata: map[string]string{}},
					{Index: 2, Name: stageScriptRunRollback, Rollback: true, Metadata: map[string]string{}},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp, err := p.BuildPipelineSyncStages(ctx, &struct{}{}, tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, resp)

		})
	}
}
func Test_FetchDefinedStages(t *testing.T) {
	p := &plugin{}
	want := []string{"SCRIPT_RUN", "SCRIPT_RUN_ROLLBACK"}
	got := p.FetchDefinedStages()

	assert.Equal(t, want, got)
}
func mustMarshalJSON(t *testing.T, v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal json: %v", err)
	}
	return data
}
