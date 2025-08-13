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
	"context"
	"encoding/json"
	"testing"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestBuildPipelineSyncStages(t *testing.T) {
	t.Parallel()
	p := &plugin{}
	ctx := context.Background()

	testcases := []struct {
		name     string
		input    *sdk.BuildPipelineSyncStagesInput
		expected *sdk.BuildPipelineSyncStagesResponse
		wantErr  bool
	}{
		{
			name: "valid stage with approvers",
			input: &sdk.BuildPipelineSyncStagesInput{
				Request: sdk.BuildPipelineSyncStagesRequest{
					Stages: []sdk.StageConfig{
						{
							Index: 0,
							Name:  stageWaitApproval,
							Config: mustMarshalJSON(t, waitApprovalStageOptions{
								Approvers:      []string{"alice", "bob"},
								MinApproverNum: 1,
							}),
						},
					},
				},
			},
			expected: &sdk.BuildPipelineSyncStagesResponse{
				Stages: []sdk.PipelineStage{
					{
						Index:               0,
						Name:                stageWaitApproval,
						Rollback:            false,
						Metadata:            map[string]string{},
						AvailableOperation:  sdk.ManualOperationApprove,
						AuthorizedOperators: []string{"alice", "bob"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid config should return error",
			input: &sdk.BuildPipelineSyncStagesInput{
				Request: sdk.BuildPipelineSyncStagesRequest{
					Stages: []sdk.StageConfig{
						{
							Index:  1,
							Name:   stageWaitApproval,
							Config: []byte("not-json"),
						},
					},
				},
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp, err := p.BuildPipelineSyncStages(ctx, &struct{}{}, tc.input)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, resp)
			}
		})
	}
}

func Test_FetchDefinedStages(t *testing.T) {
	p := &plugin{}
	want := []string{stageWaitApproval}
	got := p.FetchDefinedStages()

	assert.Equal(t, want, got)
}

func mustMarshalJSON(t *testing.T, v interface{}) []byte {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	return data
}
