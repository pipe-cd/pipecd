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
	"encoding/base64"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/cache/cachetest"
	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/datastore/datastoretest"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestAppendLabelMatchedEvents(t *testing.T) {
	labels := map[string]string{"env": "prod"}

	tests := []struct {
		name       string
		filtered   []*model.Event
		events     []*model.Event
		pageSize   int
		wantIDs    []string
		wantFull   bool
		wantCursor map[string]interface{}
	}{
		{
			name:     "dense matches stop at page size",
			pageSize: 2,
			events: []*model.Event{
				eventWithLabels("event-1", 300, labels),
				eventWithLabels("event-2", 200, labels),
				eventWithLabels("event-3", 100, labels),
			},
			wantIDs:  []string{"event-1", "event-2"},
			wantFull: true,
			wantCursor: map[string]interface{}{
				"Id":        "event-2",
				"UpdatedAt": float64(200),
			},
		},
		{
			name:     "sparse matches append until page size",
			pageSize: 2,
			filtered: []*model.Event{
				eventWithLabels("event-1", 400, labels),
			},
			events: []*model.Event{
				eventWithLabels("event-2", 300, map[string]string{"env": "dev"}),
				eventWithLabels("event-3", 200, labels),
				eventWithLabels("event-4", 100, labels),
			},
			wantIDs:  []string{"event-1", "event-3"},
			wantFull: true,
			wantCursor: map[string]interface{}{
				"Id":        "event-3",
				"UpdatedAt": float64(200),
			},
		},
		{
			name:     "zero matches does not fill page",
			pageSize: 2,
			events: []*model.Event{
				eventWithLabels("event-1", 300, map[string]string{"env": "dev"}),
				eventWithLabels("event-2", 200, map[string]string{"env": "staging"}),
			},
			wantIDs:  []string{},
			wantFull: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, cursor, full := appendLabelMatchedEvents(tt.filtered, tt.events, labels, tt.pageSize)
			assert.Equal(t, tt.wantIDs, eventIDs(got))
			assert.Equal(t, tt.wantFull, full)
			if tt.wantCursor == nil {
				assert.Empty(t, cursor)
				return
			}

			data, err := base64.StdEncoding.DecodeString(cursor)
			assert.NoError(t, err)
			gotCursor := make(map[string]interface{})
			assert.NoError(t, json.Unmarshal(data, &gotCursor))
			assert.Equal(t, tt.wantCursor, gotCursor)
		})
	}
}

func eventWithLabels(id string, updatedAt int64, labels map[string]string) *model.Event {
	return &model.Event{
		Id:        id,
		UpdatedAt: updatedAt,
		Labels:    labels,
	}
}

func eventIDs(events []*model.Event) []string {
	ids := make([]string, 0, len(events))
	for _, e := range events {
		ids = append(ids, e.Id)
	}
	return ids
}

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
