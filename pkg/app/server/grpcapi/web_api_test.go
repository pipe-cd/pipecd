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

package grpcapi

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/webservice"
	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/cache/cachetest"
	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/datastore/datastoretest"
	"github.com/pipe-cd/pipecd/pkg/jwt"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcauth"
)

func TestValidateAppBelongsToProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name             string
		appID            string
		projectID        string
		appProjectCache  cache.Cache
		applicationStore datastore.ApplicationStore
		wantErr          bool
	}{
		{
			name:      "valid with cached value",
			appID:     "appID",
			projectID: "projectID",
			appProjectCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("appID").Return("projectID", nil)
				return c
			}(),
			wantErr: false,
		},
		{
			name:      "invalid with cached value",
			appID:     "appID",
			projectID: "wrong",
			appProjectCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("appID").Return("projectID", nil)
				return c
			}(),
			wantErr: true,
		},
		{
			name:      "valid with stored value",
			appID:     "appID",
			projectID: "projectID",
			appProjectCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("appID").Return("", errors.New("not found"))
				c.EXPECT().
					Put("appID", "projectID").Return(nil)
				return c
			}(),
			applicationStore: func() datastore.ApplicationStore {
				s := datastoretest.NewMockApplicationStore(ctrl)
				s.EXPECT().
					Get(gomock.Any(), "appID").Return(&model.Application{ProjectId: "projectID"}, nil)
				return s
			}(),
			wantErr: false,
		},
		{
			name:      "invalid with stored value",
			appID:     "appID",
			projectID: "wrong",
			appProjectCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("appID").Return("", errors.New("not found"))
				c.EXPECT().
					Put("appID", "projectID").Return(nil)
				return c
			}(),
			applicationStore: func() datastore.ApplicationStore {
				s := datastoretest.NewMockApplicationStore(ctrl)
				s.EXPECT().
					Get(gomock.Any(), "appID").Return(&model.Application{ProjectId: "projectID"}, nil)
				return s
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &WebAPI{
				appProjectCache:  tt.appProjectCache,
				applicationStore: tt.applicationStore,
			}
			err := api.validateAppBelongsToProject(ctx, tt.appID, tt.projectID)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestValidateDeploymentBelongsToProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name                   string
		deploymentID           string
		projectID              string
		deploymentProjectCache cache.Cache
		deploymentStore        datastore.DeploymentStore
		wantErr                bool
	}{
		{
			name:         "valid with cached value",
			deploymentID: "deploymentID",
			projectID:    "projectID",
			deploymentProjectCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("deploymentID").Return("projectID", nil)
				return c
			}(),
			wantErr: false,
		},
		{
			name:         "invalid with cached value",
			deploymentID: "deploymentID",
			projectID:    "wrong",
			deploymentProjectCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("deploymentID").Return("projectID", nil)
				return c
			}(),
			wantErr: true,
		},
		{
			name:         "valid with stored value",
			deploymentID: "deploymentID",
			projectID:    "projectID",
			deploymentProjectCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("deploymentID").Return("", errors.New("not found"))
				c.EXPECT().
					Put("deploymentID", "projectID").Return(nil)
				return c
			}(),
			deploymentStore: func() datastore.DeploymentStore {
				s := datastoretest.NewMockDeploymentStore(ctrl)
				s.EXPECT().
					Get(gomock.Any(), "deploymentID").Return(&model.Deployment{ProjectId: "projectID"}, nil)
				return s
			}(),
			wantErr: false,
		},
		{
			name:         "invalid with stored value",
			deploymentID: "deploymentID",
			projectID:    "wrong",
			deploymentProjectCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("deploymentID").Return("", errors.New("not found"))
				c.EXPECT().
					Put("deploymentID", "projectID").Return(nil)
				return c
			}(),
			deploymentStore: func() datastore.DeploymentStore {
				s := datastoretest.NewMockDeploymentStore(ctrl)
				s.EXPECT().
					Get(gomock.Any(), "deploymentID").Return(&model.Deployment{ProjectId: "projectID"}, nil)
				return s
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &WebAPI{
				deploymentProjectCache: tt.deploymentProjectCache,
				deploymentStore:        tt.deploymentStore,
			}
			err := api.validateDeploymentBelongsToProject(ctx, tt.deploymentID, tt.projectID)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestValidatePipedBelongsToProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name              string
		pipedID           string
		projectID         string
		pipedProjectCache cache.Cache
		pipedStore        datastore.PipedStore
		wantErr           bool
	}{
		{
			name:      "valid with cached value",
			pipedID:   "pipedID",
			projectID: "projectID",
			pipedProjectCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("pipedID").Return("projectID", nil)
				return c
			}(),
			wantErr: false,
		},
		{
			name:      "invalid with cached value",
			pipedID:   "pipedID",
			projectID: "wrong",
			pipedProjectCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("pipedID").Return("projectID", nil)
				return c
			}(),
			wantErr: true,
		},
		{
			name:      "valid with stored value",
			pipedID:   "pipedID",
			projectID: "projectID",
			pipedProjectCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("pipedID").Return("", errors.New("not found"))
				c.EXPECT().
					Put("pipedID", "projectID").Return(nil)
				return c
			}(),
			pipedStore: func() datastore.PipedStore {
				s := datastoretest.NewMockPipedStore(ctrl)
				s.EXPECT().
					Get(gomock.Any(), "pipedID").Return(&model.Piped{ProjectId: "projectID"}, nil)
				return s
			}(),
			wantErr: false,
		},
		{
			name:      "invalid with stored value",
			pipedID:   "pipedID",
			projectID: "wrong",
			pipedProjectCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("pipedID").Return("", errors.New("not found"))
				c.EXPECT().
					Put("pipedID", "projectID").Return(nil)
				return c
			}(),
			pipedStore: func() datastore.PipedStore {
				s := datastoretest.NewMockPipedStore(ctrl)
				s.EXPECT().
					Get(gomock.Any(), "pipedID").Return(&model.Piped{ProjectId: "projectID"}, nil)
				return s
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &WebAPI{
				pipedProjectCache: tt.pipedProjectCache,
				pipedStore:        tt.pipedStore,
			}
			err := api.validatePipedBelongsToProject(ctx, tt.pipedID, tt.projectID)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestValidateApprover(t *testing.T) {
	tests := []struct {
		name      string
		stages    []*model.PipelineStage
		commander string
		stageID   string
		wantErr   bool
	}{
		{
			name: "valid if a commander is included in approvers",
			stages: []*model.PipelineStage{
				{
					Id: "stage-id",
					AuthorizedOperators: []string{
						"user1",
						"user2",
					},
				},
			},
			commander: "user1",
			stageID:   "stage-id",
			wantErr:   false,
		},
		{
			name: "valid if a commander match approvers",
			stages: []*model.PipelineStage{
				{
					Id: "stage-id",
					AuthorizedOperators: []string{
						"user1",
					},
				},
			},
			commander: "user1",
			stageID:   "stage-id",
			wantErr:   false,
		},
		{
			name: "invalid if a commander isn't included in approvers",
			stages: []*model.PipelineStage{
				{
					Id: "stage-id",
					AuthorizedOperators: []string{
						"user2",
						"user3",
					},
				},
			},
			commander: "user1",
			stageID:   "stage-id",
			wantErr:   true,
		},
		{
			name: "valid if the AuthorizedOperators is empty",
			stages: []*model.PipelineStage{
				{
					Id: "stage-id",
				},
			},
			commander: "user1",
			stageID:   "stage-id",
			wantErr:   false,
		},
		{
			name: "invalid if a commander isn't included in approvers metadata for pipedv0 compatibility",
			stages: []*model.PipelineStage{
				{
					Id: "stage-id",
					Metadata: map[string]string{
						"Approvers": "user2,user3",
					},
				},
			},
			commander: "user1",
			stageID:   "stage-id",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateApprover(tt.stages, tt.commander, tt.stageID)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestListDeploymentsWithLabelsReturnsStableFilteredPages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	firstOpts := datastore.ListOptions{
		Filters: []datastore.ListFilter{
			{
				Field:    "ProjectId",
				Operator: datastore.OperatorEqual,
				Value:    "project-id",
			},
			{
				Field:    "UpdatedAt",
				Operator: datastore.OperatorGreaterThanOrEqual,
				Value:    int64(0),
			},
		},
		Orders: []datastore.Order{
			{
				Field:     "UpdatedAt",
				Direction: datastore.Desc,
			},
			{
				Field:     "Id",
				Direction: datastore.Asc,
			},
		},
		Limit:  2,
		Cursor: "",
	}
	secondOpts := firstOpts
	secondOpts.Cursor = "page-2"
	thirdOpts := firstOpts
	thirdOpts.Cursor = "page-3"

	store := datastoretest.NewMockDeploymentStore(ctrl)
	store.EXPECT().
		List(gomock.Any(), firstOpts).
		Return([]*model.Deployment{
			testDeployment("dep-1", 300, map[string]string{"team": "backend"}),
			testDeployment("dep-2", 299, map[string]string{"team": "frontend"}),
		}, "page-2", nil)
	store.EXPECT().
		List(gomock.Any(), secondOpts).
		Return([]*model.Deployment{
			testDeployment("dep-3", 298, map[string]string{"team": "backend"}),
			testDeployment("dep-4", 297, map[string]string{"team": "backend"}),
		}, "page-3", nil)
	store.EXPECT().
		List(gomock.Any(), secondOpts).
		Return([]*model.Deployment{
			testDeployment("dep-3", 298, map[string]string{"team": "backend"}),
			testDeployment("dep-4", 297, map[string]string{"team": "backend"}),
		}, "page-3", nil)
	store.EXPECT().
		List(gomock.Any(), thirdOpts).
		Return([]*model.Deployment{}, "", nil)

	api := &WebAPI{
		deploymentStore: store,
		logger:          zap.NewNop(),
	}

	firstResp, err := listDeploymentsWithClaims(t, api, &webservice.ListDeploymentsRequest{
		PageSize: 2,
		Options: &webservice.ListDeploymentsRequest_Options{
			Labels: map[string]string{"team": "backend"},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, []string{"dep-1", "dep-3"}, deploymentIDs(firstResp.Deployments))

	cursor, ok, err := decodeListDeploymentsCursor(firstResp.Cursor)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, listDeploymentsCursor{
		DatastoreCursor: "page-2",
		Skip:            1,
	}, cursor)

	secondResp, err := listDeploymentsWithClaims(t, api, &webservice.ListDeploymentsRequest{
		PageSize: 2,
		Cursor:   firstResp.Cursor,
		Options: &webservice.ListDeploymentsRequest_Options{
			Labels: map[string]string{"team": "backend"},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, []string{"dep-4"}, deploymentIDs(secondResp.Deployments))
	assert.Equal(t, "", secondResp.Cursor)
}

func TestListDeploymentsWithLabelsKeepsDatastoreCursorAtPageBoundary(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	opts := datastore.ListOptions{
		Filters: []datastore.ListFilter{
			{
				Field:    "ProjectId",
				Operator: datastore.OperatorEqual,
				Value:    "project-id",
			},
			{
				Field:    "UpdatedAt",
				Operator: datastore.OperatorGreaterThanOrEqual,
				Value:    int64(0),
			},
		},
		Orders: []datastore.Order{
			{
				Field:     "UpdatedAt",
				Direction: datastore.Desc,
			},
			{
				Field:     "Id",
				Direction: datastore.Asc,
			},
		},
		Limit:  2,
		Cursor: "",
	}

	store := datastoretest.NewMockDeploymentStore(ctrl)
	store.EXPECT().
		List(gomock.Any(), opts).
		Return([]*model.Deployment{
			testDeployment("dep-1", 300, map[string]string{"team": "backend"}),
			testDeployment("dep-2", 299, map[string]string{"team": "backend"}),
		}, "page-2", nil)

	api := &WebAPI{
		deploymentStore: store,
		logger:          zap.NewNop(),
	}

	resp, err := listDeploymentsWithClaims(t, api, &webservice.ListDeploymentsRequest{
		PageSize: 2,
		Options: &webservice.ListDeploymentsRequest_Options{
			Labels: map[string]string{"team": "backend"},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, []string{"dep-1", "dep-2"}, deploymentIDs(resp.Deployments))
	assert.Equal(t, "page-2", resp.Cursor)
}

func listDeploymentsWithClaims(t *testing.T, api *WebAPI, req *webservice.ListDeploymentsRequest) (*webservice.ListDeploymentsResponse, error) {
	t.Helper()

	interceptor := rpcauth.JWTUnaryServerInterceptor(
		staticJWTVerifier{
			claims: jwt.NewClaims("user-id", "", time.Minute, model.Role{ProjectId: "project-id"}),
		},
		allowAllAuthorizer{},
		zap.NewNop(),
	)

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("cookie", "token=test-token"))
	resp, err := interceptor(
		ctx,
		req,
		&grpc.UnaryServerInfo{FullMethod: "/webservice.WebService/ListDeployments"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return api.ListDeployments(ctx, req.(*webservice.ListDeploymentsRequest))
		},
	)
	if err != nil {
		return nil, err
	}
	return resp.(*webservice.ListDeploymentsResponse), nil
}

func deploymentIDs(deployments []*model.Deployment) []string {
	ids := make([]string, 0, len(deployments))
	for _, d := range deployments {
		ids = append(ids, d.Id)
	}
	return ids
}

func testDeployment(id string, updatedAt int64, labels map[string]string) *model.Deployment {
	return &model.Deployment{
		Id:        id,
		UpdatedAt: updatedAt,
		Labels:    labels,
	}
}

type staticJWTVerifier struct {
	claims *jwt.Claims
}

func (v staticJWTVerifier) Verify(string) (*jwt.Claims, error) {
	return v.claims, nil
}

type allowAllAuthorizer struct{}

func (allowAllAuthorizer) Authorize(context.Context, string, model.Role) bool {
	return true
}
