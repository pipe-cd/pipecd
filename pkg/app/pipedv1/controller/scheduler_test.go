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

package controller

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/deploysource"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/model"
	pluginapi "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
)

func TestDetermineStageStatus(t *testing.T) {
	testcases := []struct {
		name     string
		sig      StopSignalType
		ori      model.StageStatus
		got      model.StageStatus
		expected model.StageStatus
	}{
		{
			name:     "No stop signal, should get got status",
			sig:      StopSignalNone,
			ori:      model.StageStatus_STAGE_RUNNING,
			got:      model.StageStatus_STAGE_SUCCESS,
			expected: model.StageStatus_STAGE_SUCCESS,
		}, {
			name:     "Terminated signal given, should get original status",
			sig:      StopSignalTerminate,
			ori:      model.StageStatus_STAGE_RUNNING,
			got:      model.StageStatus_STAGE_SKIPPED,
			expected: model.StageStatus_STAGE_RUNNING,
		}, {
			name:     "Timeout signal given, should get failed status",
			sig:      StopSignalTimeout,
			ori:      model.StageStatus_STAGE_RUNNING,
			got:      model.StageStatus_STAGE_RUNNING,
			expected: model.StageStatus_STAGE_FAILURE,
		}, {
			name:     "Cancel signal given, should get cancelled status",
			sig:      StopSignalCancel,
			ori:      model.StageStatus_STAGE_RUNNING,
			got:      model.StageStatus_STAGE_RUNNING,
			expected: model.StageStatus_STAGE_CANCELLED,
		}, {
			name:     "Unknown signal type given, should get failed status",
			sig:      StopSignalType("unknown"),
			ori:      model.StageStatus_STAGE_RUNNING,
			got:      model.StageStatus_STAGE_RUNNING,
			expected: model.StageStatus_STAGE_FAILURE,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := determineStageStatus(tc.sig, tc.ori, tc.got)
			assert.Equal(t, tc.expected, got)
		})
	}
}

type fakeExecutorPluginClient struct {
	pluginapi.PluginClient
}

func (m *fakeExecutorPluginClient) ExecuteStage(ctx context.Context, req *deployment.ExecuteStageRequest, opts ...grpc.CallOption) (*deployment.ExecuteStageResponse, error) {
	return &deployment.ExecuteStageResponse{
		Status: model.StageStatus_STAGE_SUCCESS,
	}, nil
}

type fakeApiClient struct {
	apiClient
}

func (f *fakeApiClient) ReportStageStatusChanged(ctx context.Context, req *pipedservice.ReportStageStatusChangedRequest, opts ...grpc.CallOption) (*pipedservice.ReportStageStatusChangedResponse, error) {
	return nil, nil
}

func TestExecuteStage(t *testing.T) {
	logger := zaptest.NewLogger(t)

	testcases := []struct {
		name              string
		deployment        *model.Deployment
		stageStatuses     map[string]model.StageStatus
		applicationConfig *config.GenericApplicationSpec
		expected          model.StageStatus
	}{
		{
			name: "stage not started yet, everything go right",
			deployment: &model.Deployment{
				Stages: []*model.PipelineStage{
					{
						Id:     "stage-id",
						Name:   "stage-name",
						Index:  0,
						Status: model.StageStatus_STAGE_NOT_STARTED_YET,
					},
				},
			},
			stageStatuses: map[string]model.StageStatus{
				"stage-id": model.StageStatus_STAGE_NOT_STARTED_YET,
			},
			expected: model.StageStatus_STAGE_SUCCESS,
		},
		{
			name: "stage is rollback but base stage not started yet, should not trigger anything",
			deployment: &model.Deployment{
				Stages: []*model.PipelineStage{
					{
						Id:       "stage-rollback-id",
						Name:     "stage-rollback-name",
						Status:   model.StageStatus_STAGE_NOT_STARTED_YET,
						Rollback: true,
						Metadata: map[string]string{
							"baseStageId": "stage-id",
						},
					},
				},
			},
			stageStatuses: map[string]model.StageStatus{
				"stage-id": model.StageStatus_STAGE_NOT_STARTED_YET,
			},
			expected: model.StageStatus_STAGE_NOT_STARTED_YET,
		},
		{
			name: "stage is rollback but base stage is skipped, should not trigger anything",
			deployment: &model.Deployment{
				Stages: []*model.PipelineStage{
					{
						Id:       "stage-rollback-id",
						Name:     "stage-rollback-name",
						Status:   model.StageStatus_STAGE_NOT_STARTED_YET,
						Rollback: true,
						Metadata: map[string]string{
							"baseStageId": "stage-id",
						},
					},
				},
			},
			stageStatuses: map[string]model.StageStatus{
				"stage-id": model.StageStatus_STAGE_SKIPPED,
			},
			expected: model.StageStatus_STAGE_NOT_STARTED_YET,
		},
		{
			name: "stage which can not be handled by the current scheduler, should be set as failed",
			deployment: &model.Deployment{
				Stages: []*model.PipelineStage{
					{
						Id:     "stage-id",
						Name:   "stage-name-not-found",
						Status: model.StageStatus_STAGE_NOT_STARTED_YET,
					},
				},
			},
			stageStatuses: map[string]model.StageStatus{
				"stage-id": model.StageStatus_STAGE_NOT_STARTED_YET,
			},
			expected: model.StageStatus_STAGE_FAILURE,
		},
		{
			name: "stage without config, should be success",
			deployment: &model.Deployment{
				Stages: []*model.PipelineStage{
					{
						Id:     "stage-id",
						Name:   "stage-name",
						Index:  0,
						Status: model.StageStatus_STAGE_NOT_STARTED_YET,
					},
				},
			},
			stageStatuses: map[string]model.StageStatus{
				"stage-id": model.StageStatus_STAGE_NOT_STARTED_YET,
			},
			applicationConfig: &config.GenericApplicationSpec{
				Pipeline: &config.DeploymentPipeline{
					Stages: []config.PipelineStage{},
				},
			},
			expected: model.StageStatus_STAGE_SUCCESS,
		},
	}

	sig, handler := NewStopSignal()
	defer handler.Terminate()

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := &scheduler{
				apiClient:  &fakeApiClient{},
				targetDSP:  &fakeDeploySourceProvider{},
				runningDSP: &fakeDeploySourceProvider{},
				stageBasedPluginsMap: map[string]pluginapi.PluginClient{
					"stage-name": &fakeExecutorPluginClient{},
				},
				genericApplicationConfig: &config.GenericApplicationSpec{
					Pipeline: &config.DeploymentPipeline{
						Stages: []config.PipelineStage{
							{ID: "stage-id", Name: "stage-name"},
						},
					},
				},
				deployment:    tc.deployment,
				stageStatuses: tc.stageStatuses,
				logger:        logger,
				nowFunc:       time.Now,
			}

			if tc.applicationConfig != nil {
				s.genericApplicationConfig = tc.applicationConfig
			}

			finalStatus := s.executeStage(sig, s.deployment.Stages[0])
			assert.Equal(t, tc.expected, finalStatus)
		})
	}
}

type fakeDeploySourceProvider struct {
	deploysource.Provider
}

func (f *fakeDeploySourceProvider) Get(ctx context.Context, logWriter io.Writer) (*deploysource.DeploySource, error) {
	return &deploysource.DeploySource{}, nil
}

func TestExecuteStage_SignalTerminated(t *testing.T) {
	logger := zaptest.NewLogger(t)
	sig, handler := NewStopSignal()

	s := &scheduler{
		apiClient:  &fakeApiClient{},
		targetDSP:  &fakeDeploySourceProvider{},
		runningDSP: &fakeDeploySourceProvider{},
		stageBasedPluginsMap: map[string]pluginapi.PluginClient{
			"stage-name": &fakeExecutorPluginClient{},
		},
		genericApplicationConfig: &config.GenericApplicationSpec{
			Pipeline: &config.DeploymentPipeline{
				Stages: []config.PipelineStage{
					{ID: "stage-id", Name: "stage-name"},
				},
			},
		},
		deployment: &model.Deployment{
			Stages: []*model.PipelineStage{
				{
					Id:     "stage-id",
					Name:   "stage-name",
					Index:  0,
					Status: model.StageStatus_STAGE_NOT_STARTED_YET,
				},
			},
		},
		stageStatuses: map[string]model.StageStatus{},
		logger:        logger,
		nowFunc:       time.Now,
	}

	handler.Terminate()
	finalStatus := s.executeStage(sig, s.deployment.Stages[0])
	assert.Equal(t, model.StageStatus_STAGE_FAILURE, finalStatus)
}

func TestExecuteStage_SignalCancelled(t *testing.T) {
	logger := zaptest.NewLogger(t)
	sig, handler := NewStopSignal()

	s := &scheduler{
		apiClient:  &fakeApiClient{},
		targetDSP:  &fakeDeploySourceProvider{},
		runningDSP: &fakeDeploySourceProvider{},
		stageBasedPluginsMap: map[string]pluginapi.PluginClient{
			"stage-name": &fakeExecutorPluginClient{},
		},
		genericApplicationConfig: &config.GenericApplicationSpec{
			Pipeline: &config.DeploymentPipeline{
				Stages: []config.PipelineStage{
					{ID: "stage-id", Name: "stage-name"},
				},
			},
		},
		deployment: &model.Deployment{
			Stages: []*model.PipelineStage{
				{
					Id:     "stage-id",
					Name:   "stage-name",
					Index:  0,
					Status: model.StageStatus_STAGE_NOT_STARTED_YET,
				},
			},
		},
		stageStatuses: map[string]model.StageStatus{},
		logger:        logger,
		nowFunc:       time.Now,
	}

	handler.Cancel()
	finalStatus := s.executeStage(sig, s.deployment.Stages[0])
	assert.Equal(t, model.StageStatus_STAGE_FAILURE, finalStatus)
}
