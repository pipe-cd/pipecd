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
	"strings"
	"testing"

	"github.com/pipe-cd/piped-plugin-sdk-go/logpersister/logpersistertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/common"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
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

func (m *mockStagePlugin) ExecuteStage(ctx context.Context, config *struct{}, targets []*DeployTarget[struct{}], input *ExecuteStageInput[struct{}]) (*ExecuteStageResponse, error) {
	return &ExecuteStageResponse{
		Status: m.result,
	}, m.err
}

func newTestStagePluginServiceServer(t *testing.T, plugin *mockStagePlugin) *StagePluginServiceServer[struct{}, struct{}, struct{}] {
	return &StagePluginServiceServer[struct{}, struct{}, struct{}]{
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

			config := strings.TrimSpace(`
apiVersion: pipecd.dev/v1beta1
kind: Appilcation
spec: {}
`)

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
					RunningDeploymentSource: &common.DeploymentSource{
						ApplicationDirectory:      "app-dir",
						CommitHash:                "commit-hash",
						ApplicationConfig:         []byte(config),
						ApplicationConfigFilename: "app-config-filename",
					},
					TargetDeploymentSource: &common.DeploymentSource{
						ApplicationDirectory:      "app-dir",
						CommitHash:                "commit-hash",
						ApplicationConfig:         []byte(config),
						ApplicationConfigFilename: "app-config-filename",
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
	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.Empty(t, response.GetStages())
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

func TestNewDetermineVersionsRequest(t *testing.T) {
	validConfig := strings.TrimSpace(`
apiVersion: pipecd.dev/v1beta1
kind: Appilcation
spec: {}
`)

	tests := []struct {
		name     string
		request  *deployment.DetermineVersionsRequest
		expected DetermineVersionsRequest[struct{}]
	}{
		{
			name: "valid request",
			request: &deployment.DetermineVersionsRequest{
				Input: &deployment.PlanPluginInput{
					Deployment: &model.Deployment{
						Id:              "deployment-id",
						ApplicationId:   "app-id",
						ApplicationName: "app-name",
						PipedId:         "piped-id",
						ProjectId:       "project-id",
						CreatedAt:       1234567890,
						Trigger: &model.DeploymentTrigger{
							Commander: "triggered-by",
						},
					},
					TargetDeploymentSource: &common.DeploymentSource{
						ApplicationDirectory:      "app-dir",
						CommitHash:                "commit-hash",
						ApplicationConfig:         []byte(validConfig),
						ApplicationConfigFilename: "app-config-filename",
					},
				},
			},
			expected: DetermineVersionsRequest[struct{}]{
				Deployment: Deployment{
					ID:              "deployment-id",
					ApplicationID:   "app-id",
					ApplicationName: "app-name",
					PipedID:         "piped-id",
					ProjectID:       "project-id",
					TriggeredBy:     "triggered-by",
					CreatedAt:       1234567890,
				},
				DeploymentSource: DeploymentSource[struct{}]{
					ApplicationDirectory:      "app-dir",
					CommitHash:                "commit-hash",
					ApplicationConfig:         nil, // ApplicationConfig is a pointer, so nil is fine for this test
					ApplicationConfigFilename: "app-config-filename",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := newDetermineVersionsRequest[struct{}]("test-plugin", tt.request)
			require.NoError(t, err)
			assert.Equal(t, tt.expected.Deployment, result.Deployment)
			assert.Equal(t, tt.expected.DeploymentSource.ApplicationDirectory, result.DeploymentSource.ApplicationDirectory)
			assert.Equal(t, tt.expected.DeploymentSource.CommitHash, result.DeploymentSource.CommitHash)
			assert.Equal(t, tt.expected.DeploymentSource.ApplicationConfigFilename, result.DeploymentSource.ApplicationConfigFilename)
		})
	}
}

func TestArtifactVersion_toModel(t *testing.T) {
	// Use the default value for expected "Kind" for now.
	// They will be removed after deleting Kind from model.ArtifactVersion.

	tests := []struct {
		name     string
		version  ArtifactVersion
		expected *model.ArtifactVersion
	}{
		{
			name: "container image",
			version: ArtifactVersion{
				Version: "v1.0.0",
				Name:    "nginx",
				URL:     "https://example.com/nginx:v1.0.0",
			},
			expected: &model.ArtifactVersion{
				Kind:    model.ArtifactVersion_UNKNOWN,
				Version: "v1.0.0",
				Name:    "nginx",
				Url:     "https://example.com/nginx:v1.0.0",
			},
		},
		{
			name: "s3 object",
			version: ArtifactVersion{
				Version: "v1.0.0",
				Name:    "backup",
				URL:     "s3://bucket/backup/v1.0.0",
			},
			expected: &model.ArtifactVersion{
				Kind:    model.ArtifactVersion_UNKNOWN,
				Version: "v1.0.0",
				Name:    "backup",
				Url:     "s3://bucket/backup/v1.0.0",
			},
		},
		{
			name: "git source",
			version: ArtifactVersion{
				Version: "commit-hash",
				Name:    "repo",
				URL:     "https://github.com/repo/commit/commit-hash",
			},
			expected: &model.ArtifactVersion{
				Kind:    model.ArtifactVersion_UNKNOWN,
				Version: "commit-hash",
				Name:    "repo",
				Url:     "https://github.com/repo/commit/commit-hash",
			},
		},
		{
			name: "terraform module",
			version: ArtifactVersion{
				Version: "v1.0.0",
				Name:    "module",
				URL:     "https://registry.terraform.io/modules/module/v1.0.0",
			},
			expected: &model.ArtifactVersion{
				Kind:    model.ArtifactVersion_UNKNOWN,
				Version: "v1.0.0",
				Name:    "module",
				Url:     "https://registry.terraform.io/modules/module/v1.0.0",
			},
		},
		{
			name: "unknown kind",
			version: ArtifactVersion{
				Version: "v1.0.0",
				Name:    "unknown",
				URL:     "https://example.com/unknown:v1.0.0",
			},
			expected: &model.ArtifactVersion{
				Kind:    model.ArtifactVersion_UNKNOWN,
				Version: "v1.0.0",
				Name:    "unknown",
				Url:     "https://example.com/unknown:v1.0.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.version.toModel()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDetermineVersionsResponse_toModel(t *testing.T) {
	// Use the default value for expected "Kind" for now.
	// They will be removed after deleting Kind from model.ArtifactVersion.

	tests := []struct {
		name     string
		response DetermineVersionsResponse
		expected []*model.ArtifactVersion
	}{
		{
			name: "single version",
			response: DetermineVersionsResponse{
				Versions: []ArtifactVersion{
					{
						Version: "v1.0.0",
						Name:    "nginx",
						URL:     "https://example.com/nginx:v1.0.0",
					},
				},
			},
			expected: []*model.ArtifactVersion{
				{
					Kind:    model.ArtifactVersion_UNKNOWN,
					Version: "v1.0.0",
					Name:    "nginx",
					Url:     "https://example.com/nginx:v1.0.0",
				},
			},
		},
		{
			name: "multiple versions",
			response: DetermineVersionsResponse{
				Versions: []ArtifactVersion{
					{
						Version: "v1.0.0",
						Name:    "nginx",
						URL:     "https://example.com/nginx:v1.0.0",
					},
					{
						Version: "v1.0.0",
						Name:    "backup",
						URL:     "s3://bucket/backup/v1.0.0",
					},
				},
			},
			expected: []*model.ArtifactVersion{
				{
					Kind:    model.ArtifactVersion_UNKNOWN,
					Version: "v1.0.0",
					Name:    "nginx",
					Url:     "https://example.com/nginx:v1.0.0",
				},
				{
					Kind:    model.ArtifactVersion_UNKNOWN,
					Version: "v1.0.0",
					Name:    "backup",
					Url:     "s3://bucket/backup/v1.0.0",
				},
			},
		},
		{
			name: "empty versions",
			response: DetermineVersionsResponse{
				Versions: []ArtifactVersion{},
			},
			expected: []*model.ArtifactVersion{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.toModel()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewDetermineStrategyRequest(t *testing.T) {
	t.Parallel()

	validConfig := strings.TrimSpace(`
apiVersion: pipecd.dev/v1beta1
kind: Appilcation
spec: {}
	`)

	tests := []struct {
		name     string
		request  *deployment.DetermineStrategyRequest
		expected DetermineStrategyRequest[struct{}]
	}{
		{
			name: "valid request",
			request: &deployment.DetermineStrategyRequest{
				Input: &deployment.PlanPluginInput{
					Deployment: &model.Deployment{
						Id:              "deployment-id",
						ApplicationId:   "app-id",
						ApplicationName: "app-name",
						PipedId:         "piped-id",
						ProjectId:       "project-id",
						CreatedAt:       1234567890,
						Trigger: &model.DeploymentTrigger{
							Commander: "triggered-by",
						},
					},
					RunningDeploymentSource: &common.DeploymentSource{
						ApplicationDirectory:      "app-dir",
						CommitHash:                "commit-hash-1",
						ApplicationConfig:         []byte(validConfig),
						ApplicationConfigFilename: "app-config-filename",
					},
					TargetDeploymentSource: &common.DeploymentSource{
						ApplicationDirectory:      "app-dir",
						CommitHash:                "commit-hash-2",
						ApplicationConfig:         []byte(validConfig),
						ApplicationConfigFilename: "app-config-filename",
					},
				},
			},
			expected: DetermineStrategyRequest[struct{}]{
				Deployment: Deployment{
					ID:              "deployment-id",
					ApplicationID:   "app-id",
					ApplicationName: "app-name",
					PipedID:         "piped-id",
					ProjectID:       "project-id",
					TriggeredBy:     "triggered-by",
					CreatedAt:       1234567890,
				},
				RunningDeploymentSource: DeploymentSource[struct{}]{
					ApplicationDirectory:      "app-dir",
					CommitHash:                "commit-hash-1",
					ApplicationConfig:         nil,
					ApplicationConfigFilename: "app-config-filename",
				},
				TargetDeploymentSource: DeploymentSource[struct{}]{
					ApplicationDirectory:      "app-dir",
					CommitHash:                "commit-hash-2",
					ApplicationConfig:         nil,
					ApplicationConfigFilename: "app-config-filename",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, _ := newDetermineStrategyRequest[struct{}]("test-plugin", tt.request)
			assert.Equal(t, tt.expected.Deployment, result.Deployment)
		})
	}
}

func TestNewDetermineStrategyResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		response  *DetermineStrategyResponse
		expected  *deployment.DetermineStrategyResponse
		expectErr bool
	}{
		{
			name: "valid quick sync strategy",
			response: &DetermineStrategyResponse{
				Strategy: SyncStrategyQuickSync,
				Summary:  "quick sync strategy",
			},
			expected: &deployment.DetermineStrategyResponse{
				SyncStrategy: model.SyncStrategy_QUICK_SYNC,
				Summary:      "quick sync strategy",
			},
			expectErr: false,
		},
		{
			name: "valid pipeline sync strategy",
			response: &DetermineStrategyResponse{
				Strategy: SyncStrategyPipelineSync,
				Summary:  "pipeline sync strategy",
			},
			expected: &deployment.DetermineStrategyResponse{
				SyncStrategy: model.SyncStrategy_PIPELINE,
				Summary:      "pipeline sync strategy",
			},
			expectErr: false,
		},
		{
			name: "invalid strategy",
			response: &DetermineStrategyResponse{
				Strategy: SyncStrategy(999),
				Summary:  "invalid strategy",
			},
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := newDetermineStrategyResponse(tt.response)
			if tt.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSyncStrategy_toModelEnum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		strategy  SyncStrategy
		expected  model.SyncStrategy
		expectErr bool
	}{
		{
			name:      "quick sync",
			strategy:  SyncStrategyQuickSync,
			expected:  model.SyncStrategy_QUICK_SYNC,
			expectErr: false,
		},
		{
			name:      "pipeline sync",
			strategy:  SyncStrategyPipelineSync,
			expected:  model.SyncStrategy_PIPELINE,
			expectErr: false,
		},
		{
			name:      "invalid strategy",
			strategy:  SyncStrategy(999),
			expected:  0,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := tt.strategy.toModelEnum()
			if tt.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
