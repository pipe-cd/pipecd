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

package sdk

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/zap/zaptest"

	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister/logpersistertest"
)

type mockStagePlugin struct {
	result StageStatus
	err    error
}

func (m *mockStagePlugin) Name() string {
	return "mockStagePlugin"
}

func (m *mockStagePlugin) Version() string {
	return "v1.0.0"
}

func (m *mockStagePlugin) FetchDefinedStages() []string {
	return []string{"stage1", "stage2"}
}

func (m *mockStagePlugin) BuildPipelineSyncStages(ctx context.Context, config *struct{}, input *BuildPipelineSyncStagesInput) (*BuildPipelineSyncStagesResponse, error) {
	return &BuildPipelineSyncStagesResponse{
		Stages: []PipelineStage{
			{
				Index: 0,
				Name:  "stage1",
			},
			{
				Index: 1,
				Name:  "stage2",
			},
			{
				Index:    0,
				Name:     "rollback",
				Rollback: true,
			},
		},
	}, m.err
}

func (m *mockStagePlugin) ExecuteStage(ctx context.Context, config *struct{}, targets []*DeployTarget[struct{}], input *ExecuteStageInput) (*ExecuteStageResponse, error) {
	return &ExecuteStageResponse{
		Status: m.result,
	}, m.err
}

func newTestStagePluginServiceServer(t *testing.T, plugin *mockStagePlugin) *StagePluginServiceServer[struct{}, struct{}] {
	return &StagePluginServiceServer[struct{}, struct{}]{
		base: plugin,
		commonFields: commonFields{
			logger:       zaptest.NewLogger(t),
			logPersister: logpersistertest.NewTestLogPersister(t),
		},
	}
}

func TestStagePluginServiceServer_ExecuteStage(t *testing.T) {
	tests := []struct {
		name           string
		stage          string
		status         StageStatus
		err            error
		expectedStatus model.StageStatus
		expectErr      bool
	}{
		{
			name:           "success",
			stage:          "stage1",
			status:         StageStatusSuccess,
			expectedStatus: model.StageStatus_STAGE_SUCCESS,
		},
		{
			name:           "failure",
			stage:          "stage2",
			status:         StageStatusFailure,
			expectedStatus: model.StageStatus_STAGE_FAILURE,
		},
		{
			name:           "exited",
			stage:          "stage2",
			status:         StageStatusExited,
			expectedStatus: model.StageStatus_STAGE_EXITED,
		},
		{
			name:      "error",
			stage:     "unknown",
			err:       errors.New("some error"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plugin := &mockStagePlugin{
				result: tt.status,
				err:    tt.err,
			}
			server := newTestStagePluginServiceServer(t, plugin)

			request := &deployment.ExecuteStageRequest{
				Input: &deployment.ExecutePluginInput{
					Stage: &model.PipelineStage{
						Name: tt.stage,
					},
					Deployment: &model.Deployment{
						Trigger: &model.DeploymentTrigger{
							Commit: &model.Commit{},
						},
					},
				},
			}

			response, err := server.ExecuteStage(context.Background(), request)
			if (err != nil) != tt.expectErr {
				t.Fatalf("unexpected error: %v", err)
			}

			if response.GetStatus() != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, response.Status)
			}
		})
	}
}

func TestStagePluginServiceServer_BuildPipelineSyncStages(t *testing.T) {
	tests := []struct {
		name         string
		request      *deployment.BuildPipelineSyncStagesRequest
		err          error
		expectStages int
		expectErr    bool
	}{
		{
			name: "success",
			request: &deployment.BuildPipelineSyncStagesRequest{
				Stages: []*deployment.BuildPipelineSyncStagesRequest_StageConfig{
					{
						Index: 0,
						Name:  "stage1",
					},
					{
						Index: 1,
						Name:  "stage2",
					},
				},
				Rollback: true,
			},
			err:          nil,
			expectStages: 3,
			expectErr:    false,
		},
		{
			name:      "error on plugin",
			request:   &deployment.BuildPipelineSyncStagesRequest{},
			err:       errors.New("some error"),
			expectErr: true,
		},
		{
			name: "returned non-valid stage index",
			request: &deployment.BuildPipelineSyncStagesRequest{
				Stages: []*deployment.BuildPipelineSyncStagesRequest_StageConfig{
					{
						Index: 100,
						Name:  "unknown",
					},
				},
			},
			err:       nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plugin := &mockStagePlugin{
				err: tt.err,
			}

			server := newTestStagePluginServiceServer(t, plugin)

			response, err := server.BuildPipelineSyncStages(context.Background(), tt.request)
			if (err != nil) != tt.expectErr {
				t.Fatalf("unexpected error: %v", err)
			}

			if !tt.expectErr && response == nil {
				t.Errorf("expected non-nil response")
			}

			if response != nil && len(response.GetStages()) != tt.expectStages {
				t.Errorf("expected %d stages, got %d", tt.expectStages, len(response.GetStages()))
			}
		})
	}
}

func TestStagePluginServiceServer_FetchDefinedStages(t *testing.T) {
	plugin := &mockStagePlugin{}
	server := newTestStagePluginServiceServer(t, plugin)

	request := &deployment.FetchDefinedStagesRequest{}
	response, err := server.FetchDefinedStages(context.Background(), request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stages := response.GetStages()
	if len(stages) != 2 {
		t.Errorf("expected 2 stages, got %d", len(stages))
	}

	expectedStages := []string{"stage1", "stage2"}
	for i, stage := range stages {
		if stage != expectedStages[i] {
			t.Errorf("expected stage %s, got %s", expectedStages[i], stage)
		}
	}
}

func TestStagePluginServiceServer_DetermineVersions(t *testing.T) {
	plugin := &mockStagePlugin{}
	server := newTestStagePluginServiceServer(t, plugin)

	request := &deployment.DetermineVersionsRequest{}
	response, err := server.DetermineVersions(context.Background(), request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response == nil {
		t.Errorf("expected non-nil response")
	}

	if len(response.GetVersions()) != 0 {
		t.Errorf("expected 0 versions, got %d", len(response.GetVersions()))
	}
}

func TestStagePluginServiceServer_DetermineStrategy(t *testing.T) {
	plugin := &mockStagePlugin{}
	server := newTestStagePluginServiceServer(t, plugin)

	request := &deployment.DetermineStrategyRequest{}
	response, err := server.DetermineStrategy(context.Background(), request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response == nil {
		t.Errorf("expected non-nil response")
	}

	if !response.GetUnsupported() {
		t.Errorf("expected unsupported strategy")
	}
}

func TestStagePluginServiceServer_BuildQuickSyncStages(t *testing.T) {
	plugin := &mockStagePlugin{}
	server := newTestStagePluginServiceServer(t, plugin)

	request := &deployment.BuildQuickSyncStagesRequest{}
	response, err := server.BuildQuickSyncStages(context.Background(), request)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if response != nil {
		t.Errorf("expected nil response, got %v", response)
	}
}

func TestStageStatus_toModelEnum(t *testing.T) {
	tests := []struct {
		name     string
		status   StageStatus
		expected model.StageStatus
	}{
		{
			name:     "success",
			status:   StageStatusSuccess,
			expected: model.StageStatus_STAGE_SUCCESS,
		},
		{
			name:     "failure",
			status:   StageStatusFailure,
			expected: model.StageStatus_STAGE_FAILURE,
		},
		{
			name:     "exited",
			status:   StageStatusExited,
			expected: model.StageStatus_STAGE_EXITED,
		},
		{
			name:     "unknown",
			status:   StageStatus(999),
			expected: model.StageStatus_STAGE_FAILURE,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.toModelEnum()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestManualOperation_toModelEnum(t *testing.T) {
	tests := []struct {
		name     string
		op       ManualOperation
		expected model.ManualOperation
	}{
		{
			name:     "none",
			op:       ManualOperationNone,
			expected: model.ManualOperation_MANUAL_OPERATION_NONE,
		},
		{
			name:     "skip",
			op:       ManualOperationSkip,
			expected: model.ManualOperation_MANUAL_OPERATION_SKIP,
		},
		{
			name:     "approve",
			op:       ManualOperationApprove,
			expected: model.ManualOperation_MANUAL_OPERATION_APPROVE,
		},
		{
			name:     "unknown",
			op:       ManualOperation(999),
			expected: model.ManualOperation_MANUAL_OPERATION_UNKNOWN,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.op.toModelEnum()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
