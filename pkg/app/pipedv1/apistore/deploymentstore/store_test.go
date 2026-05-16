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

// fakeAPIClient returns a pre-canned sequence of pages, one per call.
// If the test code calls beyond len(pages), an empty terminal response is
// returned so misconfigured tests fail with a clear assertion mismatch
// instead of an index-out-of-range panic.
type fakeAPIClient struct {
	pages []*pipedservice.ListNotCompletedDeploymentsResponse
	call  int
	err   error
}

func (f *fakeAPIClient) ListNotCompletedDeployments(_ context.Context, _ *pipedservice.ListNotCompletedDeploymentsRequest, _ ...grpc.CallOption) (*pipedservice.ListNotCompletedDeploymentsResponse, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.call >= len(f.pages) {
		return &pipedservice.ListNotCompletedDeploymentsResponse{}, nil
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
		name         string
		pages        []*pipedservice.ListNotCompletedDeploymentsResponse
		wantPendings []*model.Deployment
		wantPlanneds []*model.Deployment
		wantRunnings []*model.Deployment
		wantHeads    map[string]*model.Deployment
	}{
		{
			name: "empty_response",
			pages: []*pipedservice.ListNotCompletedDeploymentsResponse{
				{Deployments: nil, Cursor: ""},
			},
			wantPendings: nil,
			wantPlanneds: nil,
			wantRunnings: nil,
			wantHeads:    map[string]*model.Deployment{},
		},
		{
			name: "single_page_classifies_each_status",
			pages: []*pipedservice.ListNotCompletedDeploymentsResponse{
				{Deployments: []*model.Deployment{pending, planned, running, rollingBack}, Cursor: ""},
			},
			wantPendings: []*model.Deployment{pending},
			wantPlanneds: []*model.Deployment{planned},
			wantRunnings: []*model.Deployment{running, rollingBack},
			wantHeads: map[string]*model.Deployment{
				"app-1": pending,
				"app-2": planned,
				"app-3": running,
				"app-4": rollingBack,
			},
		},
		{
			name: "paginates_across_multiple_pages",
			pages: []*pipedservice.ListNotCompletedDeploymentsResponse{
				{Deployments: []*model.Deployment{pending}, Cursor: "page2"},
				{Deployments: []*model.Deployment{planned, running}, Cursor: ""},
			},
			wantPendings: []*model.Deployment{pending},
			wantPlanneds: []*model.Deployment{planned},
			wantRunnings: []*model.Deployment{running},
			wantHeads: map[string]*model.Deployment{
				"app-1": pending,
				"app-2": planned,
				"app-3": running,
			},
		},
		{
			name: "rolling_back_is_classified_as_running",
			pages: []*pipedservice.ListNotCompletedDeploymentsResponse{
				{Deployments: []*model.Deployment{rollingBack}, Cursor: ""},
			},
			wantPendings: nil,
			wantPlanneds: nil,
			wantRunnings: []*model.Deployment{rollingBack},
			wantHeads: map[string]*model.Deployment{
				"app-4": rollingBack,
			},
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
			assert.Equal(t, tc.wantHeads, s.ListAppHeadDeployments())
		})
	}
}
