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

package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

func TestBuildPipelineStages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    *sdk.BuildPipelineSyncStagesInput
		expected []sdk.PipelineStage
	}{
		{
			name: "without auto rollback",
			input: &sdk.BuildPipelineSyncStagesInput{
				Request: sdk.BuildPipelineSyncStagesRequest{
					Rollback: false,
					Stages: []sdk.StageConfig{
						{
							Index:  0,
							Name:   "Stage 1",
							Config: []byte(""),
						},
						{
							Index:  1,
							Name:   "Stage 2",
							Config: []byte(""),
						},
					},
				},
			},
			expected: []sdk.PipelineStage{
				{
					Index:              0,
					Name:               "Stage 1",
					Rollback:           false,
					Metadata:           make(map[string]string, 0),
					AvailableOperation: sdk.ManualOperationNone,
				},
				{
					Index:              1,
					Name:               "Stage 2",
					Rollback:           false,
					Metadata:           make(map[string]string, 0),
					AvailableOperation: sdk.ManualOperationNone,
				},
			},
		},
		{
			name: "with auto rollback",
			input: &sdk.BuildPipelineSyncStagesInput{
				Request: sdk.BuildPipelineSyncStagesRequest{
					Rollback: true,
					Stages: []sdk.StageConfig{
						{
							Index:  0,
							Name:   "Stage 1",
							Config: []byte(""),
						},
						{
							Index:  1,
							Name:   "Stage 2",
							Config: []byte(""),
						},
					},
				},
			},
			expected: []sdk.PipelineStage{
				{
					Index:              0,
					Name:               "Stage 1",
					Rollback:           false,
					Metadata:           make(map[string]string, 0),
					AvailableOperation: sdk.ManualOperationNone,
				},
				{
					Index:              1,
					Name:               "Stage 2",
					Rollback:           false,
					Metadata:           make(map[string]string, 0),
					AvailableOperation: sdk.ManualOperationNone,
				},
				{
					Index:              0,
					Name:               StageK8sMultiRollback,
					Rollback:           true,
					Metadata:           make(map[string]string, 0),
					AvailableOperation: sdk.ManualOperationNone,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := BuildPipelineStages(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
