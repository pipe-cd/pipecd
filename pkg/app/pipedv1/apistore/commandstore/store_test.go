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

package commandstore

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/model"
)

// mockAPIClient is a mock implementation of the apiClient interface
type mockAPIClient struct {
	reportCommandHandledCalls []*pipedservice.ReportCommandHandledRequest
	reportCommandHandledError error
}

func (m *mockAPIClient) ListUnhandledCommands(ctx context.Context, in *pipedservice.ListUnhandledCommandsRequest, opts ...grpc.CallOption) (*pipedservice.ListUnhandledCommandsResponse, error) {
	return &pipedservice.ListUnhandledCommandsResponse{}, nil
}

func (m *mockAPIClient) ReportCommandHandled(ctx context.Context, in *pipedservice.ReportCommandHandledRequest, opts ...grpc.CallOption) (*pipedservice.ReportCommandHandledResponse, error) {
	m.reportCommandHandledCalls = append(m.reportCommandHandledCalls, in)
	return &pipedservice.ReportCommandHandledResponse{}, m.reportCommandHandledError
}

func (m *mockAPIClient) getReportCommandHandledCalls() []*pipedservice.ReportCommandHandledRequest {
	return m.reportCommandHandledCalls
}

func (m *mockAPIClient) setReportCommandHandledError(err error) {
	m.reportCommandHandledError = err
}

func TestListStageCommands(t *testing.T) {
	t.Parallel()

	store := store{
		stageCommands: stageCommandsMap{
			"deployment-1": {
				"stage-1": []*model.Command{
					{
						Id:           "command-1",
						DeploymentId: "deployment-1",
						StageId:      "stage-1",
						Type:         model.Command_APPROVE_STAGE,
						Commander:    "commander-1",
					},
					{
						Id:           "command-2",
						DeploymentId: "deployment-1",
						StageId:      "stage-1",
						Type:         model.Command_APPROVE_STAGE,
						Commander:    "commander-2",
					},
					{
						Id:           "command-3",
						DeploymentId: "deployment-1",
						StageId:      "stage-1",
						Type:         model.Command_SKIP_STAGE,
					},
				},
			},
		},
		logger: zap.NewNop(),
	}

	testcases := []struct {
		name         string
		deploymentID string
		stageID      string
		want         []*model.Command
	}{
		{
			name:         "valid arguments",
			deploymentID: "deployment-1",
			stageID:      "stage-1",
			want: []*model.Command{
				{
					Id:           "command-1",
					DeploymentId: "deployment-1",
					StageId:      "stage-1",
					Type:         model.Command_APPROVE_STAGE,
					Commander:    "commander-1",
				},
				{
					Id:           "command-2",
					DeploymentId: "deployment-1",
					StageId:      "stage-1",
					Type:         model.Command_APPROVE_STAGE,
					Commander:    "commander-2",
				},
				{
					Id:           "command-3",
					DeploymentId: "deployment-1",
					StageId:      "stage-1",
					Type:         model.Command_SKIP_STAGE,
				},
			},
		},
		{
			name:         "deploymentID not exist",
			deploymentID: "xxx",
			stageID:      "stage-1",
			want:         nil,
		},
		{
			name:         "stageID not exist",
			deploymentID: "deployment-1",
			stageID:      "stage-999",
			want:         nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := store.ListStageCommands(tc.deploymentID, tc.stageID)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestReportStageCommandsHandled(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := zap.NewNop()

	testCases := []struct {
		name          string
		deploymentID  string
		stageID       string
		stageCommands stageCommandsMap
		mockSetup     func(*mockAPIClient)
		expectedError bool
		expectedCalls int
	}{
		{
			name:         "successfully report multiple commands",
			deploymentID: "deployment-1",
			stageID:      "stage-1",
			stageCommands: stageCommandsMap{
				"deployment-1": {
					"stage-1": []*model.Command{
						{
							Id:           "command-1",
							DeploymentId: "deployment-1",
							StageId:      "stage-1",
							Type:         model.Command_APPROVE_STAGE,
							Commander:    "commander-1",
						},
						{
							Id:           "command-2",
							DeploymentId: "deployment-1",
							StageId:      "stage-1",
							Type:         model.Command_SKIP_STAGE,
							Commander:    "commander-2",
						},
					},
				},
			},
			mockSetup: func(m *mockAPIClient) {
				// No setup needed - mock will succeed by default
			},
			expectedError: false,
			expectedCalls: 2,
		},
		{
			name:         "no commands to report",
			deploymentID: "deployment-1",
			stageID:      "stage-1",
			stageCommands: stageCommandsMap{
				"deployment-1": {
					"stage-1": []*model.Command{}, // Empty slice
				},
			},
			mockSetup: func(m *mockAPIClient) {
				// No setup needed
			},
			expectedError: false,
			expectedCalls: 0,
		},
		{
			name:         "deployment not found",
			deploymentID: "deployment-1",
			stageID:      "stage-1",
			stageCommands: stageCommandsMap{
				"deployment-2": { // Different deployment ID
					"stage-1": []*model.Command{
						{
							Id:           "command-1",
							DeploymentId: "deployment-2",
							StageId:      "stage-1",
							Type:         model.Command_APPROVE_STAGE,
						},
					},
				},
			},
			mockSetup: func(m *mockAPIClient) {
				// No setup needed
			},
			expectedError: false,
			expectedCalls: 0,
		},
		{
			name:         "stage not found",
			deploymentID: "deployment-1",
			stageID:      "stage-1",
			stageCommands: stageCommandsMap{
				"deployment-1": {
					"stage-2": []*model.Command{ // Different stage ID
						{
							Id:           "command-1",
							DeploymentId: "deployment-1",
							StageId:      "stage-2",
							Type:         model.Command_APPROVE_STAGE,
						},
					},
				},
			},
			mockSetup: func(m *mockAPIClient) {
				// No setup needed
			},
			expectedError: false,
			expectedCalls: 0,
		},
		{
			name:         "API client error",
			deploymentID: "deployment-1",
			stageID:      "stage-1",
			stageCommands: stageCommandsMap{
				"deployment-1": {
					"stage-1": []*model.Command{
						{
							Id:           "command-1",
							DeploymentId: "deployment-1",
							StageId:      "stage-1",
							Type:         model.Command_APPROVE_STAGE,
						},
					},
				},
			},
			mockSetup: func(m *mockAPIClient) {
				// Mock API error
				m.setReportCommandHandledError(errors.New("API error"))
			},
			expectedError: true,
			expectedCalls: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create mock API client
			mockClient := &mockAPIClient{}
			tc.mockSetup(mockClient)

			// Create store with mock client
			store := &store{
				apiClient:       mockClient,
				stageCommands:   tc.stageCommands,
				handledCommands: make(map[string]time.Time),
				logger:          logger,
			}

			// Execute the function
			err := store.ReportStageCommandsHandled(ctx, tc.deploymentID, tc.stageID)

			// Assert results
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify mock calls
			calls := mockClient.getReportCommandHandledCalls()
			assert.Equal(t, tc.expectedCalls, len(calls), "Expected %d calls to ReportCommandHandled, got %d", tc.expectedCalls, len(calls))

			// Verify command details for successful cases
			if tc.expectedCalls > 0 {
				for i, call := range calls {
					assert.Equal(t, model.CommandStatus_COMMAND_SUCCEEDED, call.Status, "Call %d should have COMMAND_SUCCEEDED status", i)
					assert.NotEmpty(t, call.CommandId, "Call %d should have a command ID", i)
					assert.NotZero(t, call.HandledAt, "Call %d should have a handled timestamp", i)
				}
			}

			// Verify that commands are cleared from the map after successful reporting
			if !tc.expectedError && tc.expectedCalls > 0 {
				commands := store.stageCommands[tc.deploymentID][tc.stageID]
				assert.Empty(t, commands, "Commands should be cleared after successful reporting")
			}
		})
	}
}
