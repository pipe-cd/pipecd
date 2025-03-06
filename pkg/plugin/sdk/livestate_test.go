package sdk

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/livestate"
)

type mockLivestatePlugin struct {
	result *GetLivestateResponse
	err    error
}

func (m *mockLivestatePlugin) Name() string {
	return "mockLivestatePlugin"
}

func (m *mockLivestatePlugin) Version() string {
	return "v1.0.0"
}

func (m *mockLivestatePlugin) GetLivestate(ctx context.Context, config *struct{}, targets []*DeployTarget[struct{}], input *GetLivestateInput) (*GetLivestateResponse, error) {
	return m.result, m.err
}

func newTestLivestatePluginServer(t *testing.T, plugin *mockLivestatePlugin) *LivestatePluginServer[struct{}, struct{}] {
	return &LivestatePluginServer[struct{}, struct{}]{
		base: plugin,
		commonFields: commonFields{
			logger: zaptest.NewLogger(t),
		},
		deployTargets: map[string]*DeployTarget[struct{}]{
			"target1": {
				Name: "target1",
				Labels: map[string]string{
					"key1": "value1",
				},
			},
		},
	}
}

func TestLivestatePluginServer_GetLivestate(t *testing.T) {
	tests := []struct {
		name           string
		request        *livestate.GetLivestateRequest
		result         *GetLivestateResponse
		err            error
		expectedStatus codes.Code
		expectErr      bool
	}{
		{
			name: "success",
			request: &livestate.GetLivestateRequest{
				ApplicationId: "app1",
				DeployTargets: []string{"target1"},
			},
			result: &GetLivestateResponse{
				LiveState: ApplicationLiveState{
					Resources: []ResourceState{
						{
							ID:   "resource1",
							Name: "Resource 1",
						},
					},
					HealthStatus: ApplicationHealthStateHealthy,
				},
				SyncState: ApplicationSyncState{
					Status: ApplicationSyncStateSynced,
				},
			},
			expectedStatus: codes.OK,
		},
		{
			name: "failure when deploy target not found",
			request: &livestate.GetLivestateRequest{
				ApplicationId: "app1",
				DeployTargets: []string{"target2"},
			},
			result:         &GetLivestateResponse{},
			expectErr:      true,
			expectedStatus: codes.Internal,
		},
		{
			name:           "error",
			request:        &livestate.GetLivestateRequest{},
			result:         &GetLivestateResponse{},
			err:            errors.New("some error"),
			expectErr:      true,
			expectedStatus: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plugin := &mockLivestatePlugin{
				result: tt.result,
				err:    tt.err,
			}
			server := newTestLivestatePluginServer(t, plugin)

			response, err := server.GetLivestate(context.Background(), tt.request)
			if (err != nil) != tt.expectErr {
				t.Fatalf("unexpected error: %v", err)
			}

			if status.Code(err) != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, status.Code(err))
			}

			if response != nil && response.GetApplicationLiveState().GetResources()[0].GetId() != tt.result.LiveState.Resources[0].ID {
				t.Errorf("expected resource ID %v, got %v", tt.result.LiveState.Resources[0].ID, response.GetApplicationLiveState().GetResources()[0].GetId())
			}
		})
	}
}
