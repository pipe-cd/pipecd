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

	"github.com/pipe-cd/piped-plugin-sdk-go/logpersister/logpersistertest"
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
func Test_ContextInfo_BuildEnv(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		ci      *ContextInfo
		want    map[string]string
		wantErr bool
	}{
		{
			name: "success",
			ci: &ContextInfo{
				DeploymentID:        "deployment-id",
				ApplicationID:       "application-id",
				ApplicationName:     "application-name",
				TriggeredAt:         1234567890,
				TriggeredCommitHash: "commit-hash",
				TriggeredCommander:  "commander",
				RepositoryURL:       "repo-url",
				Labels: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
				IsRollback: false,
				Summary:    "summary",
			},
			want: map[string]string{
				"SR_DEPLOYMENT_ID":         "deployment-id",
				"SR_APPLICATION_ID":        "application-id",
				"SR_APPLICATION_NAME":      "application-name",
				"SR_TRIGGERED_AT":          "1234567890",
				"SR_TRIGGERED_COMMIT_HASH": "commit-hash",
				"SR_TRIGGERED_COMMANDER":   "commander",
				"SR_REPOSITORY_URL":        "repo-url",
				"SR_SUMMARY":               "summary",
				"SR_IS_ROLLBACK":           "false",
				"SR_LABELS_KEY1":           "value1",
				"SR_LABELS_KEY2":           "value2",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ci.buildEnv()
			assert.Equal(t, tt.wantErr, err != nil)

			for k, v := range got {
				if k == "SR_CONTEXT_RAW" {
					continue
				}
				assert.Equal(t, tt.want[k], v)
			}

			var gotRaw ContextInfo
			err = json.Unmarshal([]byte(got["SR_CONTEXT_RAW"]), &gotRaw)
			assert.Nil(t, err)
			assert.Equal(t, tt.ci, &gotRaw)
		})
	}
}
func TestPlugin_ExecuteScriptRun(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name string
		req  sdk.ExecuteStageRequest[struct{}]
		lp   sdk.StageLogPersister
		want sdk.StageStatus
	}{
		{
			name: "success",
			req: sdk.ExecuteStageRequest[struct{}]{
				StageName:   stageScriptRun,
				StageConfig: []byte(`{"run": "echo 'success'"}`),
				Deployment: sdk.Deployment{
					ID:            "deployment-1",
					ApplicationID: "app-1",
				},
			},
			lp:   logpersistertest.NewTestLogPersister(t),
			want: sdk.StageStatusSuccess,
		},
		{
			name: "program failed",
			req: sdk.ExecuteStageRequest[struct{}]{
				StageName:   stageScriptRun,
				StageConfig: []byte(`{"run": "exit 1"}`),
				Deployment: sdk.Deployment{
					ID:            "deployment-2",
					ApplicationID: "app-2",
				},
			},
			lp:   logpersistertest.NewTestLogPersister(t),
			want: sdk.StageStatusFailure,
		},
		{
			name: "command failed",
			req: sdk.ExecuteStageRequest[struct{}]{
				StageName:   stageScriptRun,
				StageConfig: []byte(`{"run": "not_runnable"}`),
				Deployment: sdk.Deployment{
					ID:            "deployment-3",
					ApplicationID: "app-3",
				},
			},
			lp:   logpersistertest.NewTestLogPersister(t),
			want: sdk.StageStatusFailure,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp := executeScriptRun(t.Context(), tc.req, tc.lp)
			assert.Equal(t, tc.want, resp)
		})
	}
}
func mustMarshalJSON(t *testing.T, v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal json: %v", err)
	}
	return data
}
