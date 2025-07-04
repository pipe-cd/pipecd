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

	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/common"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/planpreview"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockPlanPreviewPlugin struct {
	result *GetPlanPreviewResponse
	err    error
}

func (m *mockPlanPreviewPlugin) GetPlanPreview(ctx context.Context, config *struct{}, targets []*DeployTarget[struct{}], input *GetPlanPreviewInput[struct{}]) (*GetPlanPreviewResponse, error) {
	return m.result, m.err
}

func newTestPlanPreviewPluginServer(t *testing.T, plugin *mockPlanPreviewPlugin) *PlanPreviewPluginServer[struct{}, struct{}, struct{}] {
	return &PlanPreviewPluginServer[struct{}, struct{}, struct{}]{
		base: plugin,
		commonFields: commonFields{
			logger: zaptest.NewLogger(t),
			config: &config.PipedPlugin{
				Name: "mockPlanPreviewPlugin",
			},
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
		result         *GetPlanPreviewResponse
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
			result: &GetPlanPreviewResponse{
				Summary:  "summary",
				NoChange: true,
				Details:  []byte("details"),
			},
			expectedStatus: codes.OK,
		},
		{
			name: "failure when deploy target not found",
			request: &planpreview.GetPlanPreviewRequest{
				ApplicationId: "app1",
				DeployTargets: []string{"target2"},
			},
			result:         &GetPlanPreviewResponse{},
			expectErr:      true,
			expectedStatus: codes.Internal,
		},
		{
			name:           "error",
			request:        &planpreview.GetPlanPreviewRequest{},
			result:         &GetPlanPreviewResponse{},
			err:            errors.New("some error"),
			expectErr:      true,
			expectedStatus: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			plugin := &mockPlanPreviewPlugin{
				result: tt.result,
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

			if response != nil && response.GetSummary() != tt.result.Summary {
				t.Errorf("expected summary %v, got %v", tt.result.Summary, response.GetSummary())
			}
		})
	}
}
