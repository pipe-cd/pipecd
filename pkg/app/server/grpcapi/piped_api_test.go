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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/cache/cachetest"
	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/datastore/datastoretest"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcauth"
)

type testPipedTokenVerifier struct {
	pipedKey string
}

func (v testPipedTokenVerifier) Verify(_ context.Context, _, _, pipedKey string) error {
	if pipedKey != v.pipedKey {
		return fmt.Errorf("invalid piped key, want: %s, got: %s", v.pipedKey, pipedKey)
	}
	return nil
}

func TestValidateAppBelongsToPiped(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name             string
		appID            string
		pipedID          string
		appPipedCache    cache.Cache
		applicationStore datastore.ApplicationStore
		wantErr          bool
	}{
		{
			name:    "valid with cached value",
			appID:   "appID",
			pipedID: "pipedID",
			appPipedCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("appID").Return("pipedID", nil)
				return c
			}(),
			wantErr: false,
		},
		{
			name:    "invalid with cached value",
			appID:   "appID",
			pipedID: "wrong",
			appPipedCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("appID").Return("pipedID", nil)
				return c
			}(),
			wantErr: true,
		},
		{
			name:    "valid with stored value",
			appID:   "appID",
			pipedID: "pipedID",
			appPipedCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("appID").Return("", errors.New("not found"))
				c.EXPECT().
					Put("appID", "pipedID").Return(nil)
				return c
			}(),
			applicationStore: func() datastore.ApplicationStore {
				s := datastoretest.NewMockApplicationStore(ctrl)
				s.EXPECT().
					Get(gomock.Any(), "appID").Return(&model.Application{PipedId: "pipedID"}, nil)
				return s
			}(),
			wantErr: false,
		},
		{
			name:    "invalid with stored value",
			appID:   "appID",
			pipedID: "wrong",
			appPipedCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("appID").Return("", errors.New("not found"))
				c.EXPECT().
					Put("appID", "pipedID").Return(nil)
				return c
			}(),
			applicationStore: func() datastore.ApplicationStore {
				s := datastoretest.NewMockApplicationStore(ctrl)
				s.EXPECT().
					Get(gomock.Any(), "appID").Return(&model.Application{PipedId: "pipedID"}, nil)
				return s
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &PipedAPI{
				appPipedCache:    tt.appPipedCache,
				applicationStore: tt.applicationStore,
			}
			err := api.validateAppBelongsToPiped(ctx, tt.appID, tt.pipedID)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestPipedAPIListApplicationsFetchesAllPages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := datastoretest.NewMockApplicationStore(ctrl)
	firstPage := make([]*model.Application, 0, 1)
	firstPage = append(firstPage, &model.Application{Id: "app-1"})
	secondPage := make([]*model.Application, 0, 1)
	secondPage = append(secondPage, &model.Application{Id: "app-2"})

	firstOpts := datastore.ListOptions{
		Filters: []datastore.ListFilter{
			{
				Field:    "ProjectId",
				Operator: datastore.OperatorEqual,
				Value:    "test-project-id",
			},
			{
				Field:    "PipedId",
				Operator: datastore.OperatorEqual,
				Value:    "test-piped-id",
			},
			{
				Field:    "Disabled",
				Operator: datastore.OperatorEqual,
				Value:    false,
			},
		},
	}
	secondOpts := firstOpts
	secondOpts.Cursor = "next-page"

	store.EXPECT().
		List(gomock.Any(), firstOpts).
		Return(firstPage, "next-page", nil)
	store.EXPECT().
		List(gomock.Any(), secondOpts).
		Return(secondPage, "", nil)

	api := &PipedAPI{
		applicationStore: store,
	}
	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{
		"authorization": []string{"PIPED-TOKEN " + rpcauth.MakePipedToken("test-project-id", "test-piped-id", "test-piped-key")},
	})

	got, err := rpcauth.PipedTokenUnaryServerInterceptor(testPipedTokenVerifier{pipedKey: "test-piped-key"}, zap.NewNop())(
		ctx,
		nil,
		nil,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return api.ListApplications(ctx, nil)
		},
	)

	assert.NoError(t, err)
	resp, ok := got.(*pipedservice.ListApplicationsResponse)
	assert.True(t, ok)
	assert.Equal(t, append(firstPage, secondPage...), resp.Applications)
}

func TestValidateDeploymentBelongsToPiped(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name                 string
		deploymentID         string
		pipedID              string
		deploymentPipedCache cache.Cache
		deploymentStore      datastore.DeploymentStore
		wantErr              bool
	}{
		{
			name:         "valid with cached value",
			deploymentID: "deploymentID",
			pipedID:      "pipedID",
			deploymentPipedCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("deploymentID").Return("pipedID", nil)
				return c
			}(),
			wantErr: false,
		},
		{
			name:         "invalid with cached value",
			deploymentID: "deploymentID",
			pipedID:      "wrong",
			deploymentPipedCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("deploymentID").Return("pipedID", nil)
				return c
			}(),
			wantErr: true,
		},
		{
			name:         "valid with stored value",
			deploymentID: "deploymentID",
			pipedID:      "pipedID",
			deploymentPipedCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("deploymentID").Return("", errors.New("not found"))
				c.EXPECT().
					Put("deploymentID", "pipedID").Return(nil)
				return c
			}(),
			deploymentStore: func() datastore.DeploymentStore {
				s := datastoretest.NewMockDeploymentStore(ctrl)
				s.EXPECT().
					Get(gomock.Any(), "deploymentID").Return(&model.Deployment{PipedId: "pipedID"}, nil)
				return s
			}(),
			wantErr: false,
		},
		{
			name:         "invalid with stored value",
			deploymentID: "deploymentID",
			pipedID:      "wrong",
			deploymentPipedCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("deploymentID").Return("", errors.New("not found"))
				c.EXPECT().
					Put("deploymentID", "pipedID").Return(nil)
				return c
			}(),
			deploymentStore: func() datastore.DeploymentStore {
				s := datastoretest.NewMockDeploymentStore(ctrl)
				s.EXPECT().
					Get(gomock.Any(), "deploymentID").Return(&model.Deployment{PipedId: "pipedID"}, nil)
				return s
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &PipedAPI{
				deploymentPipedCache: tt.deploymentPipedCache,
				deploymentStore:      tt.deploymentStore,
			}
			err := api.validateDeploymentBelongsToPiped(ctx, tt.deploymentID, tt.pipedID)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
