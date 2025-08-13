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

	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/common"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/planpreview"
)

type mockPlanPreviewPlugin struct {
	result *GetPlanPreviewResponse
	err    error
}

func (m *mockPlanPreviewPlugin) GetPlanPreview(ctx context.Context, config ConfigNone, targets DeployTargetsNone, input *GetPlanPreviewInput[struct{}]) (*GetPlanPreviewResponse, error) {
	return m.result, m.err
}

func newTestPlanPreviewPluginServer(t *testing.T, plugin *mockPlanPreviewPlugin) *PlanPreviewPluginServer[struct{}, struct{}, struct{}] {
	return &PlanPreviewPluginServer[struct{}, struct{}, struct{}]{
		base: plugin,
		commonFields: commonFields[struct{}, struct{}]{
			logger: zaptest.NewLogger(t),
			config: &config.PipedPlugin{
				Name: "mockPlanPreviewPlugin",
			},
			deployTargets: map[string]*DeployTarget[struct{}]{
				"target1": {
					Name: "target1",
					Labels: map[string]string{
						"key1": "value1",
					},
				},
			},
		},
	}
}

func TestPlanPreviewPluginServer_GetPlanPreview(t *testing.T) {
	t.Parallel()

	validConfig := strings.TrimSpace(`
apiVersion: pipecd.dev/v1beta1
kind: Appilcation
spec: {}
`)

	tests := []struct {
		name           string
		request        *planpreview.GetPlanPreviewRequest
		mockResp       *GetPlanPreviewResponse
		err            error
		expectedStatus codes.Code
		expectErr      bool
	}{
		{
			name: "success",
			request: &planpreview.GetPlanPreviewRequest{
				ApplicationId: "app1",
				DeployTargets: []string{"target1"},
				TargetDeploymentSource: &common.DeploymentSource{
					ApplicationDirectory:      "app-dir",
					CommitHash:                "commit-hash",
					ApplicationConfig:         []byte(validConfig),
					ApplicationConfigFilename: "app-config-filename",
				},
			},
			mockResp: &GetPlanPreviewResponse{
				Results: []PlanPreviewResult{
					{
						DeployTarget: "target1",
						Summary:      "summary",
						NoChange:     true,
						Details:      []byte("details"),
						DiffLanguage: "diff",
					},
				},
			},
			expectedStatus: codes.OK,
		},
		{
			name: "failure when deploy target not found",
			request: &planpreview.GetPlanPreviewRequest{
				ApplicationId: "app1",
				DeployTargets: []string{"target2"},
			},
			mockResp:       &GetPlanPreviewResponse{},
			expectErr:      true,
			expectedStatus: codes.Internal,
		},
		{
			name:           "error",
			request:        &planpreview.GetPlanPreviewRequest{},
			mockResp:       &GetPlanPreviewResponse{},
			err:            errors.New("some error"),
			expectErr:      true,
			expectedStatus: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			plugin := &mockPlanPreviewPlugin{
				result: tt.mockResp,
				err:    tt.err,
			}
			server := newTestPlanPreviewPluginServer(t, plugin)

			response, err := server.GetPlanPreview(context.Background(), tt.request)
			if (err != nil) != tt.expectErr {
				t.Fatalf("unexpected error: %v", err)
			}

			if status.Code(err) != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, status.Code(err))
			}

			if !tt.expectErr {
				require.NotNil(t, response)
				assert.Equal(t, tt.mockResp.toProto(), response)
			}
		})
	}
}

func TestGetPlanPreviewResponse_toProto(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		response *GetPlanPreviewResponse
		expected *planpreview.GetPlanPreviewResponse
	}{
		{
			name: "success",
			response: &GetPlanPreviewResponse{
				Results: []PlanPreviewResult{
					{
						DeployTarget: "target-1",
						Summary:      "summary-1",
						NoChange:     true,
						Details:      []byte("details-1"),
						DiffLanguage: "diff",
					},
					{
						DeployTarget: "target-2",
						Summary:      "summary-2",
						NoChange:     false,
						Details:      []byte("details-2"),
						// not specify DiffLanguage
					},
				},
			},
			expected: &planpreview.GetPlanPreviewResponse{
				Results: []*planpreview.PlanPreviewResult{
					{
						DeployTarget: "target-1",
						Summary:      "summary-1",
						NoChange:     true,
						Details:      []byte("details-1"),
						DiffLanguage: "diff",
					},
					{
						DeployTarget: "target-2",
						Summary:      "summary-2",
						NoChange:     false,
						Details:      []byte("details-2"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			response := tt.response.toProto()
			assert.Equal(t, tt.expected, response)
		})
	}
}
