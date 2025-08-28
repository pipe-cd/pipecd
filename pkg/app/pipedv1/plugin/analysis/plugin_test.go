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

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/analysis/config"
)

func mustMarshalJSON(t *testing.T, v interface{}) []byte {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestBuildPipelineSyncStages(t *testing.T) {
	t.Parallel()
	p := &plugin{}
	ctx := context.Background()

	testcases := []struct {
		name     string
		input    *sdk.BuildPipelineSyncStagesInput
		expected *sdk.BuildPipelineSyncStagesResponse
	}{
		{
			name: "should generate single analysis stage",
			input: &sdk.BuildPipelineSyncStagesInput{
				Request: sdk.BuildPipelineSyncStagesRequest{
					Stages: []sdk.StageConfig{
						{
							Index:  0,
							Name:   stageAnalysis,
							Config: mustMarshalJSON(t, &config.AnalysisStageOptions{}),
						},
					},
				},
			},
			expected: &sdk.BuildPipelineSyncStagesResponse{
				Stages: []sdk.PipelineStage{
					{
						Index:              0,
						Name:               stageAnalysis,
						Rollback:           false,
						Metadata:           map[string]string{},
						AvailableOperation: sdk.ManualOperationSkip,
					},
				},
			},
		},
		{
			name: "should handle multiple analysis stages",
			input: &sdk.BuildPipelineSyncStagesInput{
				Request: sdk.BuildPipelineSyncStagesRequest{
					Stages: []sdk.StageConfig{
						{
							Index:  0,
							Name:   stageAnalysis,
							Config: mustMarshalJSON(t, &config.AnalysisStageOptions{}),
						},
						{
							Index:  2,
							Name:   stageAnalysis,
							Config: mustMarshalJSON(t, &config.AnalysisStageOptions{}),
						},
					},
				},
			},
			expected: &sdk.BuildPipelineSyncStagesResponse{
				Stages: []sdk.PipelineStage{
					{
						Index:              0,
						Name:               stageAnalysis,
						Rollback:           false,
						Metadata:           map[string]string{},
						AvailableOperation: sdk.ManualOperationSkip,
					},
					{
						Index:              2,
						Name:               stageAnalysis,
						Rollback:           false,
						Metadata:           map[string]string{},
						AvailableOperation: sdk.ManualOperationSkip,
					},
				},
			},
		},
		{
			name: "should generate single analysis stage with rollback",
			input: &sdk.BuildPipelineSyncStagesInput{
				Request: sdk.BuildPipelineSyncStagesRequest{
					Stages: []sdk.StageConfig{
						{
							Index:  0,
							Name:   stageAnalysis,
							Config: mustMarshalJSON(t, &config.AnalysisStageOptions{}),
						},
					},
				},
			},
			expected: &sdk.BuildPipelineSyncStagesResponse{
				Stages: []sdk.PipelineStage{
					{
						Index:              0,
						Name:               stageAnalysis,
						Rollback:           false,
						Metadata:           map[string]string{},
						AvailableOperation: sdk.ManualOperationSkip,
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp, err := p.BuildPipelineSyncStages(ctx, &config.PluginConfig{}, tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, resp)
		})
	}
}

func TestFetchDefinedStages(t *testing.T) {
	p := &plugin{}
	want := []string{"ANALYSIS"}
	got := p.FetchDefinedStages()

	assert.Equal(t, want, got)
}
