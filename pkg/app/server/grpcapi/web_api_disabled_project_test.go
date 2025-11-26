package grpcapi

import (
	"context"
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
	"github.com/pipe-cd/pipecd/pkg/datastore/datastoretest"
	"github.com/pipe-cd/pipecd/pkg/jwt"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcauth"
)

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
		&model.Role{
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
		&model.Role{
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
		&model.Role{
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
		&model.Role{
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
		&model.Role{
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
