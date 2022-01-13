// Copyright 2020 The PipeCD Authors.
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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/cache/cachetest"
	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/datastore/datastoretest"
	"github.com/pipe-cd/pipe/pkg/model"
)

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
					GetApplication(gomock.Any(), "appID").Return(&model.Application{PipedId: "pipedID"}, nil)
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
					GetApplication(gomock.Any(), "appID").Return(&model.Application{PipedId: "pipedID"}, nil)
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
					GetDeployment(gomock.Any(), "deploymentID").Return(&model.Deployment{PipedId: "pipedID"}, nil)
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
					GetDeployment(gomock.Any(), "deploymentID").Return(&model.Deployment{PipedId: "pipedID"}, nil)
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

func TestValidateEnvBelongsToProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name             string
		envID            string
		projectID        string
		envProjectCache  cache.Cache
		environmentStore datastore.EnvironmentStore
		wantErr          bool
	}{
		{
			name:      "valid with cached value",
			envID:     "envID",
			projectID: "projectID",
			envProjectCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("envID").Return("projectID", nil)
				return c
			}(),
			wantErr: false,
		},
		{
			name:      "invalid with cached value",
			envID:     "envID",
			projectID: "wrong",
			envProjectCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("envID").Return("projectID", nil)
				return c
			}(),
			wantErr: true,
		},
		{
			name:      "valid with stored value",
			envID:     "envID",
			projectID: "projectID",
			envProjectCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("envID").Return("", errors.New("not found"))
				c.EXPECT().
					Put("envID", "projectID").Return(nil)
				return c
			}(),
			environmentStore: func() datastore.EnvironmentStore {
				s := datastoretest.NewMockEnvironmentStore(ctrl)
				s.EXPECT().
					GetEnvironment(gomock.Any(), "envID").Return(&model.Environment{ProjectId: "projectID"}, nil)
				return s
			}(),
			wantErr: false,
		},
		{
			name:      "invalid with stored value",
			envID:     "envID",
			projectID: "wrong",
			envProjectCache: func() cache.Cache {
				c := cachetest.NewMockCache(ctrl)
				c.EXPECT().
					Get("envID").Return("", errors.New("not found"))
				c.EXPECT().
					Put("envID", "projectID").Return(nil)
				return c
			}(),
			environmentStore: func() datastore.EnvironmentStore {
				s := datastoretest.NewMockEnvironmentStore(ctrl)
				s.EXPECT().
					GetEnvironment(gomock.Any(), "envID").Return(&model.Environment{ProjectId: "projectID"}, nil)
				return s
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &PipedAPI{
				envProjectCache:  tt.envProjectCache,
				environmentStore: tt.environmentStore,
			}
			err := api.validateEnvBelongsToProject(ctx, tt.envID, tt.projectID)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
