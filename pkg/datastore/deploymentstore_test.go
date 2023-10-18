// Copyright 2023 The PipeCD Authors.
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
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestDeploymentToPlannedUpdater(t *testing.T) {
	var (
		expectedDesc                  = "updated-summary"
		expectedStatusDesc            = "update-status-desc"
		expectedRunningCommitHash     = "update-running-commit-hash"
		expectedRunningConfigFilename = "update-running-config-filename"
		expectedVersion               = "update-version"
		expectedVersions              = []*model.ArtifactVersion{
			{
				Kind:    model.ArtifactVersion_CONTAINER_IMAGE,
				Version: "update-version",
				Name:    "update-image-name",
				Url:     "dummy-registry/update-image-name:update-version",
			},
		}
		expectedStages = []*model.PipelineStage{
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

		d = model.Deployment{
			Id:           "deployment-id",
			Summary:      "summary",
			StatusReason: "status-reason",
			Status:       model.DeploymentStatus_DEPLOYMENT_PENDING,
			Stages:       []*model.PipelineStage{},
		}

		updater = toPlannedUpdateFunc(
			expectedDesc,
			expectedStatusDesc,
			expectedRunningCommitHash,
			expectedRunningConfigFilename,
			expectedVersion,
			expectedVersions,
			expectedStages,
		)
	)

	err := updater(&d)
	require.NoError(t, err)
	assert.Equal(t, model.DeploymentStatus_DEPLOYMENT_PLANNED, d.Status)
	assert.Equal(t, expectedDesc, d.Summary)
	assert.Equal(t, expectedStatusDesc, d.StatusReason)
	assert.Equal(t, expectedRunningCommitHash, d.RunningCommitHash)
	assert.Equal(t, expectedRunningConfigFilename, d.RunningConfigFilename)
	assert.Equal(t, expectedVersion, d.Version)
	assert.Equal(t, expectedVersions, d.Versions)
	assert.Equal(t, expectedStages, d.Stages)
}

func TestDeploymentStatusUpdater(t *testing.T) {
	var (
		expectedStatus     = model.DeploymentStatus_DEPLOYMENT_RUNNING
		expectedStatusDesc = "update-status-desc"
		d                  = model.Deployment{
			Id:           "deployment-id",
			StatusReason: "status-reason",
			Status:       model.DeploymentStatus_DEPLOYMENT_PENDING,
		}
	)

	updater := statusUpdateFunc(expectedStatus, expectedStatusDesc)
	err := updater(&d)
	require.NoError(t, err)
	assert.Equal(t, expectedStatus, d.Status)
	assert.Equal(t, expectedStatusDesc, d.StatusReason)
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
				Id:           "deployment-id",
				StatusReason: "status-reason",
				Status:       model.DeploymentStatus_DEPLOYMENT_PENDING,
			},
			status:      model.DeploymentStatus_DEPLOYMENT_RUNNING,
			statusDesc:  "updated-status-desc",
			completedAt: now.Unix(),

			expectedErr: ErrInvalidArgument,
		},
		{
			name: "valid complete status and updated fields",
			deployment: model.Deployment{
				Id:           "deployment-id",
				StatusReason: "status-desc",
				Status:       model.DeploymentStatus_DEPLOYMENT_PENDING,
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
				Id:           "deployment-id",
				StatusReason: "updated-status-desc",
				Status:       model.DeploymentStatus_DEPLOYMENT_SUCCESS,
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
			updater := toCompletedUpdateFunc(tc.status, tc.stageStatuses, tc.statusDesc, tc.completedAt)
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
		requires     []string
		visible      bool
		retriedCount int32
		completedAt  int64

		expectedDeployment model.Deployment
		expectedErr        error
	}{
		{
			name: "stageID not found",
			deployment: model.Deployment{
				Id:           "deployment-id",
				StatusReason: "status-reason",
				Status:       model.DeploymentStatus_DEPLOYMENT_PENDING,
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
				Id:           "deployment-id",
				StatusReason: "status-desc",
				Status:       model.DeploymentStatus_DEPLOYMENT_RUNNING,
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
			requires:     []string{"stage-1"},
			visible:      true,
			retriedCount: 2,
			completedAt:  now.Unix(),

			expectedDeployment: model.Deployment{
				Id:           "deployment-id",
				StatusReason: "status-desc",
				Status:       model.DeploymentStatus_DEPLOYMENT_RUNNING,
				Stages: []*model.PipelineStage{
					{
						Id:           "stage-id1",
						Name:         "stage1",
						Desc:         "desc1",
						Index:        1,
						Status:       model.StageStatus_STAGE_SUCCESS,
						StatusReason: "updated-status-desc",
						Requires:     []string{"stage-1"},
						Visible:      true,
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
			updater := stageStatusUpdateFunc(tc.stageID, tc.status, tc.statusDesc, tc.requires, tc.visible, tc.retriedCount, tc.completedAt)
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

func TestAddDeployment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		name       string
		deployment *model.Deployment
		dsFactory  func(*model.Deployment) DataStore
		wantErr    bool
	}{
		{
			name:       "Invalid deployment",
			deployment: &model.Deployment{},
			dsFactory:  func(d *model.Deployment) DataStore { return nil },
			wantErr:    true,
		},
		{
			name: "Valid deployment",
			deployment: &model.Deployment{
				Id:              "id",
				ApplicationId:   "app-id",
				ApplicationName: "app-name",
				PipedId:         "piped-id",
				ProjectId:       "project-id",
				Kind:            model.ApplicationKind_KUBERNETES,
				GitPath: &model.ApplicationGitPath{
					Repo: &model.ApplicationGitRepository{Id: "id"},
					Path: "path",
				},
				PlatformProvider: "platform-provider",
				Trigger: &model.DeploymentTrigger{
					Commit: &model.Commit{
						Hash:      "hash",
						Message:   "message",
						Author:    "author",
						Branch:    "branch",
						CreatedAt: 1,
					},
					Timestamp: 1,
				},
				Status: model.DeploymentStatus_DEPLOYMENT_PENDING,

				CompletedAt: 1,
				CreatedAt:   1,
				UpdatedAt:   1,
			},
			dsFactory: func(d *model.Deployment) DataStore {
				ds := NewMockDataStore(ctrl)
				ds.EXPECT().Create(gomock.Any(), gomock.Any(), d.Id, d)
				return ds
			},
			wantErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewDeploymentStore(tc.dsFactory(tc.deployment), TestCommander)
			err := s.Add(context.Background(), tc.deployment)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestGetDeployment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		name    string
		id      string
		ds      DataStore
		wantErr bool
	}{
		{
			name: "successful fetch from datastore",
			id:   "id",
			ds: func() DataStore {
				ds := NewMockDataStore(ctrl)
				ds.EXPECT().
					Get(gomock.Any(), gomock.Any(), "id", &model.Deployment{}).
					Return(nil)
				return ds
			}(),
			wantErr: false,
		},
		{
			name: "failed fetch from datastore",
			id:   "id",
			ds: func() DataStore {
				ds := NewMockDataStore(ctrl)
				ds.EXPECT().
					Get(gomock.Any(), gomock.Any(), "id", &model.Deployment{}).
					Return(fmt.Errorf("err"))
				return ds
			}(),
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewDeploymentStore(tc.ds, TestCommander)
			_, err := s.Get(context.Background(), tc.id)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestListDeployments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		name    string
		opts    ListOptions
		ds      DataStore
		wantErr bool
	}{
		{
			name: "iterator done",
			opts: ListOptions{},
			ds: func() DataStore {
				it := NewMockIterator(ctrl)
				it.EXPECT().
					Next(&model.Deployment{}).
					Return(ErrIteratorDone)

				ds := NewMockDataStore(ctrl)
				ds.EXPECT().
					Find(gomock.Any(), gomock.Any(), ListOptions{}).
					Return(it, nil)
				return ds
			}(),
			wantErr: false,
		},
		{
			name: "unexpected error occurred",
			opts: ListOptions{},
			ds: func() DataStore {
				it := NewMockIterator(ctrl)
				it.EXPECT().
					Next(&model.Deployment{}).
					Return(fmt.Errorf("err"))

				ds := NewMockDataStore(ctrl)
				ds.EXPECT().
					Find(gomock.Any(), gomock.Any(), ListOptions{}).
					Return(it, nil)
				return ds
			}(),
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewDeploymentStore(tc.ds, TestCommander)
			_, _, err := s.List(context.Background(), tc.opts)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestMergeMetadata(t *testing.T) {
	testcases := []struct {
		name     string
		ori      map[string]string
		new      map[string]string
		expected map[string]string
	}{
		{
			name:     "both are empty",
			expected: map[string]string{},
		},
		{
			name: "ori map is empty",
			new: map[string]string{
				"key-1": "value-1",
			},
			expected: map[string]string{
				"key-1": "value-1",
			},
		},
		{
			name: "new map is empty",
			ori: map[string]string{
				"key-1": "value-1",
			},
			expected: map[string]string{
				"key-1": "value-1",
			},
		},
		{
			name: "there is a same key",
			ori: map[string]string{
				"key-1": "value-1",
				"key-2": "value-2",
			},
			new: map[string]string{
				"key-2": "value-22",
				"key-3": "value-3",
			},
			expected: map[string]string{
				"key-1": "value-1",
				"key-2": "value-22",
				"key-3": "value-3",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := mergeMetadata(tc.ori, tc.new)
			assert.Equal(t, got, tc.expected)
		})
	}
}
