package grpcapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/datastore/datastoretest"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcauth"
)

func TestAPISyncApplication_DisabledProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectID := "project-id"
	appID := "app-id"
	apiKeyID := "apikey-id"

	ctx := context.Background()
	// Mock API key authentication
	ctx = rpcauth.ContextWithAPIKey(ctx, &model.APIKey{
		Id:        apiKeyID,
		ProjectId: projectID,
		Role:      model.APIKey_READ_WRITE,
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

	api := &API{
		projectStore:     &mockProjectStore{MockProjectStore: ps},
		applicationStore: as,
		logger:           zap.NewNop(),
	}

	resp, err := api.SyncApplication(ctx, &apiservice.SyncApplicationRequest{
		ApplicationId: appID,
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.FailedPrecondition, st.Code())
	assert.Contains(t, st.Message(), "project is currently disabled")
}
