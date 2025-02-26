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
	return &BuildPipelineSyncStagesResponse{}, nil
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
			name:           "cancelled",
			stage:          "stage1",
			status:         StageStatusCancelled,
			expectedStatus: model.StageStatus_STAGE_CANCELLED,
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
