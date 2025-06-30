package main

import (
	"context"
	"encoding/json"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ContextInfo_BuildEnv(t *testing.T) {
	tests := []struct {
		name   string
		ci     *ContextInfo
		optEnv map[string]string
		want   []string
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
			optEnv: map[string]string{
				"env1": "x1",
				"env2": "x2",
			},
			want: []string{
				"SR_DEPLOYMENT_ID=deployment-id",
				"SR_APPLICATION_ID=application-id",
				"SR_APPLICATION_NAME=application-name",
				"SR_TRIGGERED_AT=1234567890",
				"SR_TRIGGERED_COMMIT_HASH=commit-hash",
				"SR_TRIGGERED_COMMANDER=commander",
				"SR_REPOSITORY_URL=repo-url",
				"SR_SUMMARY=summary",
				"SR_IS_ROLLBACK=false",
				"SR_LABELS_KEY1=value1",
				"SR_LABELS_KEY2=value2",
				"env1=x1",
				"env2=x2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildEnvStr(tt.ci, tt.optEnv)
			assert.Nil(t, err)

			assert.Subset(t, got, tt.want)

		})
	}
}

func TestBuildPipelineSyncStages(t *testing.T) {
	p := &plugin{}
	ctx := context.Background()

	testcases := []struct {
		name      string
		input     *sdk.BuildPipelineSyncStagesInput
		expected  *sdk.BuildPipelineSyncStagesResponse
		expectErr bool
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
								Run: "echo 'hello'",
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
			expectErr: false,
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
			expectErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := p.BuildPipelineSyncStages(ctx, &struct{}{}, tc.input)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, resp)
			}
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
