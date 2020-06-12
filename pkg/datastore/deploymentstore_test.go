// Copyright 2020 The PipeCD Authors.
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

package datastore

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipe/pkg/model"
)

func TestDeploymentToPlannedUpdater(t *testing.T) {
	expectedDesc := "updated-description"
	expectedStatusDesc := "update-status-desc"
	expectedStages := []*model.PipelineStage{
		{
			Id:    "stage-id1",
			Name:  "stage1",
			Desc:  "desc1",
			Index: 1,
			Requires: []string{
				"requires1",
			},
			Status:   model.StageStatus_STAGE_SUCCESS,
			Metadata: map[string]string{"meta": "value"},
		},
	}

	d := model.Deployment{
		Id:                "deployment-id",
		Description:       "description",
		StatusDescription: "status-description",
		Status:            model.DeploymentStatus_DEPLOYMENT_PENDING,
		Stages:            []*model.PipelineStage{},
	}

	updater := DeploymentToPlannedUpdater(
		expectedDesc,
		expectedStatusDesc,
		expectedStages,
	)
	err := updater(&d)
	require.NoError(t, err)
	assert.Equal(t, model.DeploymentStatus_DEPLOYMENT_PLANNED, d.Status)
	assert.Equal(t, expectedDesc, d.Description)
	assert.Equal(t, expectedStatusDesc, d.StatusDescription)
	assert.Equal(t, expectedStages, d.Stages)
}

func TestDeploymentStatusUpdater(t *testing.T) {
	var (
		expectedStatus     = model.DeploymentStatus_DEPLOYMENT_RUNNING
		expectedStatusDesc = "update-status-desc"
		d                  = model.Deployment{
			Id:                "deployment-id",
			StatusDescription: "status-description",
			Status:            model.DeploymentStatus_DEPLOYMENT_PENDING,
		}
	)

	updater := DeploymentStatusUpdater(expectedStatus, expectedStatusDesc)
	err := updater(&d)
	require.NoError(t, err)
	assert.Equal(t, expectedStatus, d.Status)
	assert.Equal(t, expectedStatusDesc, d.StatusDescription)
}

func TestDeploymentToCompletedUpdater(t *testing.T) {
	now := time.Now()
	testcases := []struct {
		name          string
		deployment    model.Deployment
		status        model.DeploymentStatus
		stageStatuses map[string]model.StageStatus
		statusDesc    string
		completedAt   int64

		expectedDeployment model.Deployment
		expectedErr        error
	}{
		{
			name: "invalid complete status",
			deployment: model.Deployment{
				Id:                "deployment-id",
				StatusDescription: "status-description",
				Status:            model.DeploymentStatus_DEPLOYMENT_PENDING,
			},
			status:      model.DeploymentStatus_DEPLOYMENT_RUNNING,
			statusDesc:  "updated-status-desc",
			completedAt: now.Unix(),

			expectedErr: ErrInvalidArgument,
		},
		{
			name: "valid complete status and updated fields",
			deployment: model.Deployment{
				Id:                "deployment-id",
				StatusDescription: "status-desc",
				Status:            model.DeploymentStatus_DEPLOYMENT_PENDING,
				Stages: []*model.PipelineStage{
					{
						Id:       "stage-id1",
						Name:     "stage1",
						Desc:     "desc1",
						Index:    1,
						Status:   model.StageStatus_STAGE_SUCCESS,
						Metadata: map[string]string{"meta": "value"},
					},
					{
						Id:       "stage-id2",
						Name:     "stage2",
						Desc:     "desc2",
						Index:    2,
						Status:   model.StageStatus_STAGE_RUNNING,
						Metadata: map[string]string{"meta": "value"},
					},
				},
			},
			status:     model.DeploymentStatus_DEPLOYMENT_SUCCESS,
			statusDesc: "updated-status-desc",
			stageStatuses: map[string]model.StageStatus{
				"stage-id2": model.StageStatus_STAGE_SUCCESS,
			},
			completedAt: now.Unix(),

			expectedDeployment: model.Deployment{
				Id:                "deployment-id",
				StatusDescription: "updated-status-desc",
				Status:            model.DeploymentStatus_DEPLOYMENT_SUCCESS,
				Stages: []*model.PipelineStage{
					{
						Id:       "stage-id1",
						Name:     "stage1",
						Desc:     "desc1",
						Index:    1,
						Status:   model.StageStatus_STAGE_SUCCESS,
						Metadata: map[string]string{"meta": "value"},
					},
					{
						Id:       "stage-id2",
						Name:     "stage2",
						Desc:     "desc2",
						Index:    2,
						Status:   model.StageStatus_STAGE_SUCCESS,
						Metadata: map[string]string{"meta": "value"},
					},
				},
				CompletedAt: now.Unix(),
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			updater := DeploymentToCompletedUpdater(tc.status, tc.stageStatuses, tc.statusDesc, tc.completedAt)
			err := updater(&tc.deployment)
			if err != nil {
				if tc.expectedErr == nil {
					assert.NoError(t, err)
					return
				}
				assert.Error(t, errors.Unwrap(err), tc.expectedErr)
				return
			}
			assert.Equal(t, tc.expectedDeployment, tc.deployment)
		})
	}
}

func TestStageStatusChangedUpdater(t *testing.T) {
	now := time.Now()
	testcases := []struct {
		name         string
		deployment   model.Deployment
		stageID      string
		status       model.StageStatus
		statusDesc   string
		retriedCount int32
		completedAt  int64

		expectedDeployment model.Deployment
		expectedErr        error
	}{
		{
			name: "stageID not found",
			deployment: model.Deployment{
				Id:                "deployment-id",
				StatusDescription: "status-description",
				Status:            model.DeploymentStatus_DEPLOYMENT_PENDING,
				Stages: []*model.PipelineStage{
					{
						Id: "stage-id1",
					},
				},
			},
			stageID: "not-found-stage-id",

			expectedErr: ErrInvalidArgument,
		},
		{
			name: "update target stage status",
			deployment: model.Deployment{
				Id:                "deployment-id",
				StatusDescription: "status-desc",
				Status:            model.DeploymentStatus_DEPLOYMENT_RUNNING,
				Stages: []*model.PipelineStage{
					{
						Id:           "stage-id1",
						Name:         "stage1",
						Desc:         "desc1",
						Index:        1,
						Status:       model.StageStatus_STAGE_RUNNING,
						Metadata:     map[string]string{"meta": "value"},
						RetriedCount: 1,
					},
				},
			},
			stageID:      "stage-id1",
			status:       model.StageStatus_STAGE_SUCCESS,
			statusDesc:   "updated-status-desc",
			retriedCount: 2,
			completedAt:  now.Unix(),

			expectedDeployment: model.Deployment{
				Id:                "deployment-id",
				StatusDescription: "updated-status-desc",
				Status:            model.DeploymentStatus_DEPLOYMENT_RUNNING,
				Stages: []*model.PipelineStage{
					{
						Id:           "stage-id1",
						Name:         "stage1",
						Desc:         "desc1",
						Index:        1,
						Status:       model.StageStatus_STAGE_SUCCESS,
						Metadata:     map[string]string{"meta": "value"},
						RetriedCount: 2,
						CompletedAt:  now.Unix(),
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			updater := StageStatusChangedUpdater(tc.stageID, tc.status, tc.statusDesc, tc.retriedCount, tc.completedAt)
			err := updater(&tc.deployment)
			if err != nil {
				if tc.expectedErr == nil {
					assert.NoError(t, err)
					return
				}
				assert.Error(t, errors.Unwrap(err), tc.expectedErr)
				return
			}
			assert.Equal(t, tc.expectedDeployment, tc.deployment)
		})
	}
}
