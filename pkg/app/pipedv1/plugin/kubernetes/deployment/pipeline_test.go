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
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

func Test_buildQuickSyncPipeline(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name         string
		autoRollback bool
		expected     []*model.PipelineStage
	}{
		{
			name:         "without auto rollback",
			autoRollback: false,
			expected: []*model.PipelineStage{
				{
					Id:        PredefinedStageK8sSync,
					Name:      StageK8sSync.String(),
					Desc:      "Sync by applying all manifests",
					Index:     0,
					Rollback:  false,
					Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
					Metadata:  nil,
					CreatedAt: now.Unix(),
					UpdatedAt: now.Unix(),
				},
			},
		},
		{
			name:         "with auto rollback",
			autoRollback: true,
			expected: []*model.PipelineStage{
				{
					Id:        PredefinedStageK8sSync,
					Name:      StageK8sSync.String(),
					Desc:      "Sync by applying all manifests",
					Index:     0,
					Rollback:  false,
					Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
					Metadata:  nil,
					CreatedAt: now.Unix(),
					UpdatedAt: now.Unix(),
				},
				{
					Id:        PredefinedStageRollback,
					Name:      StageK8sRollback.String(),
					Desc:      "Rollback the deployment",
					Rollback:  true,
					Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
					CreatedAt: now.Unix(),
					UpdatedAt: now.Unix(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := buildQuickSyncPipeline(tt.autoRollback, now)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestBuildPipelineStages(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name         string
		stages       []*deployment.BuildPipelineSyncStagesRequest_StageConfig
		autoRollback bool
		expected     []*model.PipelineStage
	}{
		{
			name: "without auto rollback",
			stages: []*deployment.BuildPipelineSyncStagesRequest_StageConfig{
				{
					Id:    "stage-1",
					Name:  "Stage 1",
					Desc:  "Description 1",
					Index: 0,
				},
				{
					Id:    "stage-2",
					Name:  "Stage 2",
					Desc:  "Description 2",
					Index: 1,
				},
			},
			autoRollback: false,
			expected: []*model.PipelineStage{
				{
					Id:        "stage-1",
					Name:      "Stage 1",
					Desc:      "Description 1",
					Index:     0,
					Rollback:  false,
					Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
					CreatedAt: now.Unix(),
					UpdatedAt: now.Unix(),
				},
				{
					Id:        "stage-2",
					Name:      "Stage 2",
					Desc:      "Description 2",
					Index:     1,
					Rollback:  false,
					Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
					CreatedAt: now.Unix(),
					UpdatedAt: now.Unix(),
				},
			},
		},
		{
			name: "with auto rollback",
			stages: []*deployment.BuildPipelineSyncStagesRequest_StageConfig{
				{
					Id:    "stage-1",
					Name:  "Stage 1",
					Desc:  "Description 1",
					Index: 0,
				},
				{
					Id:    "stage-2",
					Name:  "Stage 2",
					Desc:  "Description 2",
					Index: 1,
				},
			},
			autoRollback: true,
			expected: []*model.PipelineStage{
				{
					Id:        "stage-1",
					Name:      "Stage 1",
					Desc:      "Description 1",
					Index:     0,
					Rollback:  false,
					Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
					CreatedAt: now.Unix(),
					UpdatedAt: now.Unix(),
				},
				{
					Id:        "stage-2",
					Name:      "Stage 2",
					Desc:      "Description 2",
					Index:     1,
					Rollback:  false,
					Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
					CreatedAt: now.Unix(),
					UpdatedAt: now.Unix(),
				},
				{
					Id:        PredefinedStageRollback,
					Name:      StageK8sRollback.String(),
					Desc:      "Rollback the deployment",
					Index:     0,
					Rollback:  true,
					Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
					CreatedAt: now.Unix(),
					UpdatedAt: now.Unix(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := buildPipelineStages(tt.stages, tt.autoRollback, now)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestBuildQuickSyncPipeline(t *testing.T) {
	t.Parallel()

	now := time.Now()

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
					Name:               predefinedStages[PredefinedStageK8sSync].Name,
					Description:        predefinedStages[PredefinedStageK8sSync].Desc,
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
					Name:               predefinedStages[PredefinedStageK8sSync].Name,
					Description:        predefinedStages[PredefinedStageK8sSync].Desc,
					Rollback:           false,
					Metadata:           make(map[string]string, 0),
					AvailableOperation: sdk.ManualOperationNone,
				},
				{
					Name:               predefinedStages[PredefinedStageRollback].Name,
					Description:        predefinedStages[PredefinedStageRollback].Desc,
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

			actual := BuildQuickSyncPipeline(tt.rollback, now)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
