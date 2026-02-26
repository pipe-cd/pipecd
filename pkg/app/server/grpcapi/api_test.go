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

func TestRequireAPIKey(t *testing.T) {
	testcases := []struct {
		name        string
		key         *model.APIKey
		requireRole model.APIKey_Role
		expectedKey *model.APIKey
		expectedErr string
	}{
		{
			name: "ok: using READ_ONLY to read",
			key: &model.APIKey{
				Role: model.APIKey_READ_ONLY,
			},
			requireRole: model.APIKey_READ_ONLY,
			expectedKey: &model.APIKey{
				Role: model.APIKey_READ_ONLY,
			},
		},
		{
			name: "ok: using READ_WRITE to read",
			key: &model.APIKey{
				Role: model.APIKey_READ_WRITE,
			},
			requireRole: model.APIKey_READ_ONLY,
			expectedKey: &model.APIKey{
				Role: model.APIKey_READ_WRITE,
			},
		},
		{
			name: "ok: using READ_WRITE to write",
			key: &model.APIKey{
				Role: model.APIKey_READ_WRITE,
			},
			requireRole: model.APIKey_READ_WRITE,
			expectedKey: &model.APIKey{
				Role: model.APIKey_READ_WRITE,
			},
		},
		{
			name: "invalid: using READ_ONLY to write",
			key: &model.APIKey{
				Role: model.APIKey_READ_ONLY,
			},
			requireRole: model.APIKey_READ_WRITE,
			expectedErr: "rpc error: code = PermissionDenied desc = Permission denied",
		},
		{
			name: "invalid: invalid role",
			key: &model.APIKey{
				Role: -1,
			},
			requireRole: model.APIKey_READ_ONLY,
			expectedErr: "rpc error: code = PermissionDenied desc = Invalid role",
		},
		{
			name:        "invalid: api key was not included",
			requireRole: model.APIKey_READ_ONLY,
			expectedErr: "rpc error: code = Unauthenticated desc = Unauthenticated",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.TODO()
			if tc.key != nil {
				ctx = rpcauth.ContextWithAPIKey(ctx, tc.key)
			}
			key, err := requireAPIKey(ctx, tc.requireRole, zap.NewNop())
			assert.Equal(t, tc.expectedKey, key)
			if err != nil {
				assert.Equal(t, tc.expectedErr, err.Error())
			} else {
				assert.Equal(t, tc.expectedErr, "")
			}
		})
	}
}

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
