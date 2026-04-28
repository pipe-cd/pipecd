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

package deploymentstore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/model"
)

// fakeAPIClient returns a sequence of pages, one per call.
type fakeAPIClient struct {
	pages []*pipedservice.ListNotCompletedDeploymentsResponse
	call  int
	err   error
}

func (f *fakeAPIClient) ListNotCompletedDeployments(_ context.Context, _ *pipedservice.ListNotCompletedDeploymentsRequest, _ ...grpc.CallOption) (*pipedservice.ListNotCompletedDeploymentsResponse, error) {
	if f.err != nil {
		return nil, f.err
	}
	resp := f.pages[f.call]
	f.call++
	return resp, nil
}

func makeDeployment(id, appID string, status model.DeploymentStatus) *model.Deployment {
	return &model.Deployment{Id: id, ApplicationId: appID, Status: status}
}

func TestSync(t *testing.T) {
	pending := makeDeployment("d-pending", "app-1", model.DeploymentStatus_DEPLOYMENT_PENDING)
	planned := makeDeployment("d-planned", "app-2", model.DeploymentStatus_DEPLOYMENT_PLANNED)
	running := makeDeployment("d-running", "app-3", model.DeploymentStatus_DEPLOYMENT_RUNNING)
	rollingBack := makeDeployment("d-rolling-back", "app-4", model.DeploymentStatus_DEPLOYMENT_ROLLING_BACK)

	tests := []struct {
		name           string
		pages          []*pipedservice.ListNotCompletedDeploymentsResponse
		wantPendings   []*model.Deployment
		wantPlanneds   []*model.Deployment
		wantRunnings   []*model.Deployment
		wantHeadAppIDs []string
	}{
		{
			name: "empty response",
			pages: []*pipedservice.ListNotCompletedDeploymentsResponse{
				{Deployments: nil, Cursor: ""},
			},
			wantPendings:   nil,
			wantPlanneds:   nil,
			wantRunnings:   nil,
			wantHeadAppIDs: nil,
		},
		{
			name: "single page with all statuses",
			pages: []*pipedservice.ListNotCompletedDeploymentsResponse{
				{Deployments: []*model.Deployment{pending, planned, running, rollingBack}, Cursor: ""},
			},
			wantPendings:   []*model.Deployment{pending},
			wantPlanneds:   []*model.Deployment{planned},
			wantRunnings:   []*model.Deployment{running, rollingBack},
			wantHeadAppIDs: []string{"app-1", "app-2", "app-3", "app-4"},
		},
		{
			name: "multiple pages are all fetched",
			pages: []*pipedservice.ListNotCompletedDeploymentsResponse{
				{Deployments: []*model.Deployment{pending}, Cursor: "page2"},
				{Deployments: []*model.Deployment{planned, running}, Cursor: ""},
			},
			wantPendings:   []*model.Deployment{pending},
			wantPlanneds:   []*model.Deployment{planned},
			wantRunnings:   []*model.Deployment{running},
			wantHeadAppIDs: []string{"app-1", "app-2", "app-3"},
		},
		{
			name: "rolling back is classified as running",
			pages: []*pipedservice.ListNotCompletedDeploymentsResponse{
				{Deployments: []*model.Deployment{rollingBack}, Cursor: ""},
			},
			wantPendings:   nil,
			wantPlanneds:   nil,
			wantRunnings:   []*model.Deployment{rollingBack},
			wantHeadAppIDs: []string{"app-4"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &store{
				apiClient: &fakeAPIClient{pages: tc.pages},
				logger:    zap.NewNop(),
			}

			err := s.sync(context.Background())
			require.NoError(t, err)

			assert.Equal(t, tc.wantPendings, s.ListPendings())
			assert.Equal(t, tc.wantPlanneds, s.ListPlanneds())
			assert.Equal(t, tc.wantRunnings, s.ListRunnings())

			heads := s.ListAppHeadDeployments()
			if len(tc.wantHeadAppIDs) == 0 {
				assert.Empty(t, heads)
			} else {
				assert.Len(t, heads, len(tc.wantHeadAppIDs))
				for _, appID := range tc.wantHeadAppIDs {
					assert.Contains(t, heads, appID)
				}
			}
		})
	}
}
