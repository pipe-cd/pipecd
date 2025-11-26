package grpcapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/datastore/datastoretest"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcauth"
)

type mockProjectStore struct {
	*datastoretest.MockProjectStore
}

func (m *mockProjectStore) DisableProject(ctx context.Context, id string) error {
	return nil
}

func (m *mockProjectStore) EnableProject(ctx context.Context, id string) error {
	return nil
}

type mockVerifier struct{}

func (m *mockVerifier) Verify(ctx context.Context, projectID, pipedID, pipedKey string) error {
	return nil
}

func TestListApplications_DisabledProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectID := "project-id"
	pipedID := "piped-id"
	pipedKey := "piped-key"

	ctx := context.Background()
	token := rpcauth.MakePipedToken(projectID, pipedID, pipedKey)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"authorization": []string{"PIPED-TOKEN " + token},
	})

	ps := datastoretest.NewMockProjectStore(ctrl)
	ps.EXPECT().Get(gomock.Any(), projectID).Return(&model.Project{
		Id:       projectID,
		Disabled: true,
	}, nil)

	api := &PipedAPI{
		projectStore: &mockProjectStore{MockProjectStore: ps},
		logger:       zap.NewNop(),
	}

	verifier := &mockVerifier{}
	interceptor := rpcauth.PipedTokenUnaryServerInterceptor(verifier, zap.NewNop())

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return api.ListApplications(ctx, req.(*pipedservice.ListApplicationsRequest))
	}

	resp, err := interceptor(ctx, &pipedservice.ListApplicationsRequest{}, &grpc.UnaryServerInfo{}, handler)
	assert.NoError(t, err)
	assert.Empty(t, resp.(*pipedservice.ListApplicationsResponse).Applications)
}

func TestCreateDeployment_DisabledProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectID := "project-id"
	pipedID := "piped-id"
	pipedKey := "piped-key"

	ctx := context.Background()
	token := rpcauth.MakePipedToken(projectID, pipedID, pipedKey)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"authorization": []string{"PIPED-TOKEN " + token},
	})

	ps := datastoretest.NewMockProjectStore(ctrl)
	ps.EXPECT().Get(gomock.Any(), projectID).Return(&model.Project{
		Id:       projectID,
		Disabled: true,
	}, nil)

	api := &PipedAPI{
		projectStore: &mockProjectStore{MockProjectStore: ps},
		logger:       zap.NewNop(),
	}

	verifier := &mockVerifier{}
	interceptor := rpcauth.PipedTokenUnaryServerInterceptor(verifier, zap.NewNop())

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return api.CreateDeployment(ctx, req.(*pipedservice.CreateDeploymentRequest))
	}

	resp, err := interceptor(ctx, &pipedservice.CreateDeploymentRequest{
		Deployment: &model.Deployment{
			ApplicationId: "app-id",
		},
	}, &grpc.UnaryServerInfo{}, handler)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestListUnhandledCommands_DisabledProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectID := "project-id"
	pipedID := "piped-id"
	pipedKey := "piped-key"

	ctx := context.Background()
	token := rpcauth.MakePipedToken(projectID, pipedID, pipedKey)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"authorization": []string{"PIPED-TOKEN " + token},
	})

	ps := datastoretest.NewMockProjectStore(ctrl)
	ps.EXPECT().Get(gomock.Any(), projectID).Return(&model.Project{
		Id:       projectID,
		Disabled: true,
	}, nil)

	api := &PipedAPI{
		projectStore: &mockProjectStore{MockProjectStore: ps},
		logger:       zap.NewNop(),
	}

	verifier := &mockVerifier{}
	interceptor := rpcauth.PipedTokenUnaryServerInterceptor(verifier, zap.NewNop())

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return api.ListUnhandledCommands(ctx, req.(*pipedservice.ListUnhandledCommandsRequest))
	}

	resp, err := interceptor(ctx, &pipedservice.ListUnhandledCommandsRequest{}, &grpc.UnaryServerInfo{}, handler)

	assert.NoError(t, err)
	assert.Empty(t, resp.(*pipedservice.ListUnhandledCommandsResponse).Commands)
}
