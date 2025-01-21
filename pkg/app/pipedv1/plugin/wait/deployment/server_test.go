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
	"context"
	"testing"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/stretchr/testify/assert"
)

func TestFetchDefinedStages(t *testing.T) {
	t.Parallel()

	server := NewDeploymentService(nil, zap.NewNop(), nil)
	resp, err := server.FetchDefinedStages(context.Background(), &deployment.FetchDefinedStagesRequest{})
	assert.NoError(t, err)
	assert.Equal(t, resp.Stages, []string{"WAIT"})
}

func TestDetermineVersions(t *testing.T) {
	t.Parallel()

	server := NewDeploymentService(nil, zap.NewNop(), nil)
	resp, err := server.DetermineVersions(context.Background(), &deployment.DetermineVersionsRequest{})
	assert.NoError(t, err)
	assert.Equal(t, resp.Versions, []*model.ArtifactVersion{}) // Empty
}

func TestDetermineStrategy(t *testing.T) {
	t.Parallel()

	server := NewDeploymentService(nil, zap.NewNop(), nil)
	resp, err := server.DetermineStrategy(context.Background(), &deployment.DetermineStrategyRequest{})
	assert.NoError(t, err)
	assert.Equal(t, resp.Unsupported, true)
}

func TestBuildQuickSyncStages(t *testing.T) {
	t.Parallel()

	expected := &deployment.BuildQuickSyncStagesResponse{
		Stages: []*model.PipelineStage{
			{
				Id:       "WAIT",
				Name:     "WAIT",
				Desc:     "Wait for the specified duration",
				Rollback: false,
				Status:   model.StageStatus_STAGE_NOT_STARTED_YET,
			},
		},
	}

	server := NewDeploymentService(nil, zap.NewNop(), nil)

	resp, err := server.BuildQuickSyncStages(context.Background(), &deployment.BuildQuickSyncStagesRequest{
		Rollback: false,
	})
	assert.NoError(t, err)
	resp.Stages[0].CreatedAt = 0 // Ignore timestamps
	resp.Stages[0].UpdatedAt = 0
	assert.Equal(t, resp, expected)

	// The response will be the same even if Rollback is true.
	resp, err = server.BuildQuickSyncStages(context.Background(), &deployment.BuildQuickSyncStagesRequest{
		Rollback: true,
	})
	assert.NoError(t, err)
	resp.Stages[0].CreatedAt = 0 // Ignore timestamps
	resp.Stages[0].UpdatedAt = 0
	assert.Equal(t, resp, expected)
}

func TestBuildPipelineSyncStages(t *testing.T) {
	t.Parallel()

	req := &deployment.BuildPipelineSyncStagesRequest{
		Stages: []*deployment.BuildPipelineSyncStagesRequest_StageConfig{
			{
				// ID is empty.
				Name:  "WAIT",
				Index: 0,
			},
			{
				Id:    "stage-2",
				Name:  "WAIT",
				Index: 2,
			},
		},
		Rollback: true,
	}

	expected := &deployment.BuildPipelineSyncStagesResponse{
		Stages: []*model.PipelineStage{
			{
				Id:       "stage-0",
				Name:     "WAIT",
				Desc:     "Wait for the specified duration",
				Index:    0,
				Rollback: false,
				Status:   model.StageStatus_STAGE_NOT_STARTED_YET,
			},
			{
				Id:       "stage-2",
				Name:     "WAIT",
				Desc:     "Wait for the specified duration",
				Index:    2,
				Rollback: false,
				Status:   model.StageStatus_STAGE_NOT_STARTED_YET,
			},
		},
	}

	server := NewDeploymentService(nil, zap.NewNop(), nil)
	resp, err := server.BuildPipelineSyncStages(context.Background(), req)
	assert.NoError(t, err)
	resp.Stages[0].CreatedAt = 0 // Ignore timestamps
	resp.Stages[0].UpdatedAt = 0
	resp.Stages[1].CreatedAt = 0
	resp.Stages[1].UpdatedAt = 0
	assert.Equal(t, resp, expected)
}
