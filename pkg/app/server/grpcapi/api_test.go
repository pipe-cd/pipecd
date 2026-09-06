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

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/datastore"
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

func TestListApplicationsCursor(t *testing.T) {
	const storeCursor = "next-page-cursor"

	testcases := []struct {
		name           string
		req            *apiservice.ListApplicationsRequest
		storedApps     []*model.Application
		expectedCursor string
	}{
		{
			name: "cursor is returned when no label filter is given",
			req:  &apiservice.ListApplicationsRequest{Limit: 1},
			storedApps: []*model.Application{
				{Id: "app-1", ProjectId: "project-id"},
			},
			expectedCursor: storeCursor,
		},
		{
			name: "cursor is returned when a label filter is given",
			req: &apiservice.ListApplicationsRequest{
				Limit:  1,
				Labels: map[string]string{"env": "prod"},
			},
			storedApps: []*model.Application{
				{Id: "app-1", ProjectId: "project-id", Labels: map[string]string{"env": "prod"}},
			},
			expectedCursor: storeCursor,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := datastoretest.NewMockApplicationStore(ctrl)
			store.EXPECT().
				List(gomock.Any(), gomock.Any()).
				DoAndReturn(func(context.Context, datastore.ListOptions) ([]*model.Application, string, error) {
					return tc.storedApps, storeCursor, nil
				}).
				AnyTimes()

			api := &API{
				applicationStore: store,
				logger:           zap.NewNop(),
			}
			ctx := rpcauth.ContextWithAPIKey(context.TODO(), &model.APIKey{
				ProjectId: "project-id",
				Role:      model.APIKey_READ_ONLY,
			})

			resp, err := api.ListApplications(ctx, tc.req)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCursor, resp.Cursor)
			assert.Len(t, resp.Applications, len(tc.storedApps))
		})
	}
}
