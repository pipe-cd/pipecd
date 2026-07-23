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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

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

type mockJWTVerifier struct {
	claims *jwt.Claims
}

func (m *mockJWTVerifier) Verify(token string) (*jwt.Claims, error) {
	return m.claims, nil
}

type mockAuthorizer struct{}

func (m *mockAuthorizer) Authorize(ctx context.Context, method string, role model.Role) bool {
	return true
}

func TestSyncApplication_DisabledProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectID := "project-id"
	appID := "app-id"

	ctx := context.Background()
	// Mock JWT token in cookie
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"cookie": []string{"token=dummy-token"},
	})

	ps := datastoretest.NewMockProjectStore(ctrl)
	ps.EXPECT().Get(gomock.Any(), projectID).Return(&model.Project{
		Id:       projectID,
		Disabled: true,
	}, nil)

	as := datastoretest.NewMockApplicationStore(ctrl)
	as.EXPECT().Get(gomock.Any(), appID).Return(&model.Application{
		Id:        appID,
		ProjectId: projectID,
	}, nil)

	api := &WebAPI{
		projectStore:     &mockProjectStore{MockProjectStore: ps},
		applicationStore: as,
		logger:           zap.NewNop(),
	}

	claims := jwt.NewClaims(
		"sub",
		"avatar-url",
		10*time.Minute,
		model.Role{
			ProjectId:        projectID,
			ProjectRbacRoles: []string{"ADMIN"},
		},
	)
	verifier := &mockJWTVerifier{claims: claims}
	authorizer := &mockAuthorizer{}
	interceptor := rpcauth.JWTUnaryServerInterceptor(verifier, authorizer, zap.NewNop())

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return api.SyncApplication(ctx, req.(*webservice.SyncApplicationRequest))
	}

	resp, err := interceptor(ctx, &webservice.SyncApplicationRequest{
		ApplicationId: appID,
	}, &grpc.UnaryServerInfo{FullMethod: "/webservice.WebService/SyncApplication"}, handler)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.FailedPrecondition, st.Code())
	assert.Contains(t, st.Message(), "project is currently disabled")
}

func TestCancelDeployment_DisabledProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectID := "project-id"
	deploymentID := "deployment-id"

	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"cookie": []string{"token=dummy-token"},
	})

	ps := datastoretest.NewMockProjectStore(ctrl)
	ps.EXPECT().Get(gomock.Any(), projectID).Return(&model.Project{
		Id:       projectID,
		Disabled: true,
	}, nil)

	ds := datastoretest.NewMockDeploymentStore(ctrl)
	ds.EXPECT().Get(gomock.Any(), deploymentID).Return(&model.Deployment{
		Id:        deploymentID,
		ProjectId: projectID,
		Status:    model.DeploymentStatus_DEPLOYMENT_RUNNING,
	}, nil)

	api := &WebAPI{
		projectStore:    &mockProjectStore{MockProjectStore: ps},
		deploymentStore: ds,
		logger:          zap.NewNop(),
	}

	claims := jwt.NewClaims(
		"sub",
		"avatar-url",
		10*time.Minute,
		model.Role{
			ProjectId:        projectID,
			ProjectRbacRoles: []string{"ADMIN"},
		},
	)
	verifier := &mockJWTVerifier{claims: claims}
	authorizer := &mockAuthorizer{}
	interceptor := rpcauth.JWTUnaryServerInterceptor(verifier, authorizer, zap.NewNop())

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return api.CancelDeployment(ctx, req.(*webservice.CancelDeploymentRequest))
	}

	resp, err := interceptor(ctx, &webservice.CancelDeploymentRequest{
		DeploymentId: deploymentID,
	}, &grpc.UnaryServerInfo{FullMethod: "/webservice.WebService/CancelDeployment"}, handler)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.FailedPrecondition, st.Code())
	assert.Contains(t, st.Message(), "project is currently disabled")
}

func TestSkipStage_DisabledProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectID := "project-id"
	deploymentID := "deployment-id"
	stageID := "stage-id"

	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"cookie": []string{"token=dummy-token"},
	})

	ps := datastoretest.NewMockProjectStore(ctrl)
	ps.EXPECT().Get(gomock.Any(), projectID).Return(&model.Project{
		Id:       projectID,
		Disabled: true,
	}, nil)

	ds := datastoretest.NewMockDeploymentStore(ctrl)
	ds.EXPECT().Get(gomock.Any(), deploymentID).Return(&model.Deployment{
		Id:        deploymentID,
		ProjectId: projectID,
		Stages: []*model.PipelineStage{{
			Id:     stageID,
			Name:   "K8S_SYNC",
			Status: model.StageStatus_STAGE_RUNNING,
		}},
	}, nil)

	api := &WebAPI{
		projectStore:    &mockProjectStore{MockProjectStore: ps},
		deploymentStore: ds,
		logger:          zap.NewNop(),
	}

	claims := jwt.NewClaims(
		"sub",
		"avatar-url",
		10*time.Minute,
		model.Role{
			ProjectId:        projectID,
			ProjectRbacRoles: []string{"ADMIN"},
		},
	)
	verifier := &mockJWTVerifier{claims: claims}
	authorizer := &mockAuthorizer{}
	interceptor := rpcauth.JWTUnaryServerInterceptor(verifier, authorizer, zap.NewNop())

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return api.SkipStage(ctx, req.(*webservice.SkipStageRequest))
	}

	resp, err := interceptor(ctx, &webservice.SkipStageRequest{
		DeploymentId: deploymentID,
		StageId:      stageID,
	}, &grpc.UnaryServerInfo{FullMethod: "/webservice.WebService/SkipStage"}, handler)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.FailedPrecondition, st.Code())
	assert.Contains(t, st.Message(), "project is currently disabled")
}

func TestApproveStage_DisabledProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := "user-id"
	projectID := "project-id"
	deploymentID := "deployment-id"
	stageID := "stage-id"

	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"cookie": []string{"token=dummy-token"},
	})

	ps := datastoretest.NewMockProjectStore(ctrl)
	ps.EXPECT().Get(gomock.Any(), projectID).Return(&model.Project{
		Id:       projectID,
		Disabled: true,
	}, nil)

	ds := datastoretest.NewMockDeploymentStore(ctrl)
	ds.EXPECT().Get(gomock.Any(), deploymentID).Return(&model.Deployment{
		Id:        deploymentID,
		ProjectId: projectID,
		Stages: []*model.PipelineStage{{
			Id:     stageID,
			Name:   "WAIT_APPROVAL",
			Status: model.StageStatus_STAGE_RUNNING,
			Metadata: map[string]string{
				"Approvers": userID,
			},
		}},
	}, nil).Times(2)

	api := &WebAPI{
		projectStore:           &mockProjectStore{MockProjectStore: ps},
		deploymentStore:        ds,
		deploymentProjectCache: newMockCache(),
		logger:                 zap.NewNop(),
	}

	claims := jwt.NewClaims(
		userID,
		"avatar-url",
		10*time.Minute,
		model.Role{
			ProjectId:        projectID,
			ProjectRbacRoles: []string{"ADMIN"},
		},
	)
	verifier := &mockJWTVerifier{claims: claims}
	authorizer := &mockAuthorizer{}
	interceptor := rpcauth.JWTUnaryServerInterceptor(verifier, authorizer, zap.NewNop())

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return api.ApproveStage(ctx, req.(*webservice.ApproveStageRequest))
	}

	resp, err := interceptor(ctx, &webservice.ApproveStageRequest{
		DeploymentId: deploymentID,
		StageId:      stageID,
	}, &grpc.UnaryServerInfo{FullMethod: "/webservice.WebService/ApproveStage"}, handler)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.FailedPrecondition, st.Code())
	assert.Contains(t, st.Message(), "project is currently disabled")
}

func TestRestartPiped_DisabledProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectID := "project-id"
	pipedID := "piped-id"

	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"cookie": []string{"token=dummy-token"},
	})

	ps := datastoretest.NewMockProjectStore(ctrl)
	ps.EXPECT().Get(gomock.Any(), projectID).Return(&model.Project{
		Id:       projectID,
		Disabled: true,
	}, nil)

	pips := datastoretest.NewMockPipedStore(ctrl)
	pips.EXPECT().Get(gomock.Any(), pipedID).Return(&model.Piped{
		Id:        pipedID,
		ProjectId: projectID,
	}, nil)

	api := &WebAPI{
		projectStore: &mockProjectStore{MockProjectStore: ps},
		pipedStore:   pips,
		logger:       zap.NewNop(),
	}

	claims := jwt.NewClaims(
		"sub",
		"avatar-url",
		10*time.Minute,
		model.Role{
			ProjectId:        projectID,
			ProjectRbacRoles: []string{"ADMIN"},
		},
	)
	verifier := &mockJWTVerifier{claims: claims}
	authorizer := &mockAuthorizer{}
	interceptor := rpcauth.JWTUnaryServerInterceptor(verifier, authorizer, zap.NewNop())

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return api.RestartPiped(ctx, req.(*webservice.RestartPipedRequest))
	}

	resp, err := interceptor(ctx, &webservice.RestartPipedRequest{
		PipedId: pipedID,
	}, &grpc.UnaryServerInfo{FullMethod: "/webservice.WebService/RestartPiped"}, handler)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.FailedPrecondition, st.Code())
	assert.Contains(t, st.Message(), "project is currently disabled")
}

// mockCache is a simple mock implementation of cache.Cache for testing
type mockCache struct{}

func newMockCache() *mockCache {
	return &mockCache{}
}

func (c *mockCache) Get(key string) (interface{}, error) {
	return nil, assert.AnError
}

func (c *mockCache) Put(key string, value interface{}) error {
	return nil
}

func (c *mockCache) Delete(key string) error {
	return nil
}

func (c *mockCache) GetAll() (map[string]interface{}, error) {
	return nil, nil
}
