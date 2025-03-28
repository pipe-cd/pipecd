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

func Test_buildQuickSyncPipeline(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		rollback bool
		expected []sdk.QuickSyncStage
	}{
		{
			name:     "without rollback",
			rollback: false,
			expected: []sdk.QuickSyncStage{
				{
					Name:               StageK8sMultiSync,
					Description:        StageDescriptionK8sMultiSync,
					Rollback:           false,
					Metadata:           make(map[string]string, 0),
					AvailableOperation: sdk.ManualOperationNone,
				},
			},
		},
		{
			name:     "with rollback",
			rollback: true,
			expected: []sdk.QuickSyncStage{
				{
					Name:               StageK8sMultiSync,
					Description:        StageDescriptionK8sMultiSync,
					Rollback:           false,
					Metadata:           make(map[string]string, 0),
					AvailableOperation: sdk.ManualOperationNone,
				},
				{
					Name:               StageK8sMultiRollback,
					Description:        StageDescriptionK8sMultiRollback,
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

			actual := buildQuickSyncPipeline(tt.rollback)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func Test_buildPipelineStages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		stages       []sdk.StageConfig
		autoRollback bool
		expected     []sdk.PipelineStage
	}{
		{
			name: "without auto rollback",
			stages: []sdk.StageConfig{
				{
					Name:  "Stage 1",
					Index: 0,
				},
				{
					Name:  "Stage 2",
					Index: 1,
				},
			},
			autoRollback: false,
			expected: []sdk.PipelineStage{
				{
					Name:               "Stage 1",
					Index:              0,
					Rollback:           false,
					Metadata:           make(map[string]string, 0),
					AvailableOperation: sdk.ManualOperationNone,
				},
				{
					Name:               "Stage 2",
					Index:              1,
					Rollback:           false,
					Metadata:           make(map[string]string, 0),
					AvailableOperation: sdk.ManualOperationNone,
				},
			},
		},
		{
			name: "with auto rollback",
			stages: []sdk.StageConfig{
				{
					Name:  "Stage 1",
					Index: 0,
				},
				{
					Name:  "Stage 2",
					Index: 1,
				},
			},
			autoRollback: true,
			expected: []sdk.PipelineStage{
				{
					Name:               "Stage 1",
					Index:              0,
					Rollback:           false,
					Metadata:           make(map[string]string, 0),
					AvailableOperation: sdk.ManualOperationNone,
				},
				{
					Name:               "Stage 2",
					Index:              1,
					Rollback:           false,
					Metadata:           make(map[string]string, 0),
					AvailableOperation: sdk.ManualOperationNone,
				},
				{
					Name:               StageK8sMultiRollback,
					Index:              0,
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

			actual := buildPipelineStages(tt.stages, tt.autoRollback)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
