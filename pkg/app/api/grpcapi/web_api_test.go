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
	"fmt"
	"reflect"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipe/pkg/app/api/service/webservice"
	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/cache/cachetest"
	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/datastore/datastoretest"
	"github.com/pipe-cd/pipe/pkg/model"
)

func Test_filterDeploymentConfigTemplates(t *testing.T) {
	type args struct {
		labels    []webservice.DeploymentConfigTemplateLabel
		templates []*webservice.DeploymentConfigTemplate
	}
	tests := []struct {
		name string
		args args
		want []*webservice.DeploymentConfigTemplate
	}{
		{
			name: "Specify just one label",
			args: args{
				labels: []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_CANARY},
				templates: []*webservice.DeploymentConfigTemplate{
					{
						Name:   "Canary",
						Labels: []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_CANARY},
					},
					{
						Name:   "Blue/Green",
						Labels: []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_BLUE_GREEN},
					},
				},
			},
			want: []*webservice.DeploymentConfigTemplate{
				{
					Name:   "Canary",
					Labels: []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_CANARY},
				},
			},
		},
		{
			name: "Two labels specified, non-existent",
			args: args{
				labels: []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_CANARY, webservice.DeploymentConfigTemplateLabel_BLUE_GREEN},
				templates: []*webservice.DeploymentConfigTemplate{
					{
						Name:   "Canary",
						Labels: []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_CANARY},
					},
					{
						Name:   "Blue/Green",
						Labels: []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_BLUE_GREEN},
					},
				},
			},
			want: []*webservice.DeploymentConfigTemplate{},
		},
		{
			name: "Two labels specified, existent",
			args: args{
				labels: []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_CANARY, webservice.DeploymentConfigTemplateLabel_BLUE_GREEN},
				templates: []*webservice.DeploymentConfigTemplate{
					{
						Name:   "Canary Blue/Green",
						Labels: []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_CANARY, webservice.DeploymentConfigTemplateLabel_BLUE_GREEN},
					},
					{
						Name:   "Blue/Green",
						Labels: []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_BLUE_GREEN},
					},
				},
			},
			want: []*webservice.DeploymentConfigTemplate{
				{
					Name:   "Canary Blue/Green",
					Labels: []webservice.DeploymentConfigTemplateLabel{webservice.DeploymentConfigTemplateLabel_CANARY, webservice.DeploymentConfigTemplateLabel_BLUE_GREEN},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterDeploymentConfigTemplates(tt.args.templates, tt.args.labels); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterDeploymentConfigTemplates() = %v, want %v", got, tt.want)
			}
		})
	}
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
					GetApplication(gomock.Any(), "appID").Return(&model.Application{ProjectId: "projectID"}, nil)
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
					GetApplication(gomock.Any(), "appID").Return(&model.Application{ProjectId: "projectID"}, nil)
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
					GetDeployment(gomock.Any(), "deploymentID").Return(&model.Deployment{ProjectId: "projectID"}, nil)
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
					GetDeployment(gomock.Any(), "deploymentID").Return(&model.Deployment{ProjectId: "projectID"}, nil)
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
					GetPiped(gomock.Any(), "pipedID").Return(&model.Piped{ProjectId: "projectID"}, nil)
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
					GetPiped(gomock.Any(), "pipedID").Return(&model.Piped{ProjectId: "projectID"}, nil)
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

func TestGetInsightDataForDeployFrequency(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name              string
		pipedID           string
		projectID         string
		pipedProjectCache cache.Cache
		deploymentStore   datastore.DeploymentStore
		req               *webservice.GetInsightDataRequest
		res               *webservice.GetInsightDataResponse
		wantErr           bool
	}{
		{
			name:      "valid with InsightStep_DAILY",
			pipedID:   "pipedID",
			projectID: "projectID",
			deploymentStore: func() datastore.DeploymentStore {
				target := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
				targetNextDate := target.AddDate(0, 0, 1)
				s := datastoretest.NewMockDeploymentStore(ctrl)
				s.EXPECT().
					ListDeployments(gomock.Any(), datastore.ListOptions{
						Filters: []datastore.ListFilter{
							{
								Field:    "ProjectId",
								Operator: "==",
								Value:    "projectID",
							},
							{
								Field:    "CreatedAt",
								Operator: ">=",
								Value:    target.Unix(),
							},
							{
								Field:    "CreatedAt",
								Operator: "<",
								Value:    targetNextDate.Unix(),
							},
							{
								Field:    "ApplicationId",
								Operator: "==",
								Value:    "ApplicationId",
							},
						},
					}).Return([]*model.Deployment{
					{
						Id: "id1",
					},
					{
						Id: "id2",
					},
				}, nil)

				target = time.Date(2020, 1, 2, 0, 0, 0, 0, time.Local)
				targetNextDate = target.AddDate(0, 0, 1)
				s.EXPECT().
					ListDeployments(gomock.Any(), datastore.ListOptions{
						Filters: []datastore.ListFilter{
							{
								Field:    "ProjectId",
								Operator: "==",
								Value:    "projectID",
							},
							{
								Field:    "CreatedAt",
								Operator: ">=",
								Value:    target.Unix(),
							},
							{
								Field:    "CreatedAt",
								Operator: "<",
								Value:    targetNextDate.Unix(),
							},
							{
								Field:    "ApplicationId",
								Operator: "==",
								Value:    "ApplicationId",
							},
						},
					}).Return([]*model.Deployment{
					{
						Id: "id1",
					},
					{
						Id: "id2",
					},
					{
						Id: "id3",
					},
				}, nil)

				return s
			}(),
			req: &webservice.GetInsightDataRequest{
				MetricsKind:    model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
				Step:           model.InsightStep_DAILY,
				RangeFrom:      time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local).Unix(),
				DataPointCount: 2,
				ApplicationId:  "ApplicationId",
			},
			res: &webservice.GetInsightDataResponse{
				UpdatedAt: time.Now().Unix(),
				DataPoints: []*model.InsightDataPoint{
					{
						Value:     2,
						Timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local).Unix(),
					},
					{
						Value:     3,
						Timestamp: time.Date(2020, 1, 2, 0, 0, 0, 0, time.Local).Unix(),
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "valid with InsightStep_WEEKLY and rangeFrom is Sunday",
			pipedID:   "pipedID",
			projectID: "projectID",
			deploymentStore: func() datastore.DeploymentStore {
				target := time.Date(2020, 1, 5, 0, 0, 0, 0, time.Local)
				targetNextWeek := target.AddDate(0, 0, 7)
				s := datastoretest.NewMockDeploymentStore(ctrl)
				s.EXPECT().
					ListDeployments(gomock.Any(), datastore.ListOptions{
						Filters: []datastore.ListFilter{
							{
								Field:    "ProjectId",
								Operator: "==",
								Value:    "projectID",
							},
							{
								Field:    "CreatedAt",
								Operator: ">=",
								Value:    target.Unix(),
							},
							{
								Field:    "CreatedAt",
								Operator: "<",
								Value:    targetNextWeek.Unix(),
							},
							{
								Field:    "ApplicationId",
								Operator: "==",
								Value:    "ApplicationId",
							},
						},
					}).Return([]*model.Deployment{
					{
						Id: "id1",
					},
					{
						Id: "id2",
					},
				}, nil)

				target = time.Date(2020, 1, 12, 0, 0, 0, 0, time.Local)
				targetNextWeek = target.AddDate(0, 0, 7)
				s.EXPECT().
					ListDeployments(gomock.Any(), datastore.ListOptions{
						Filters: []datastore.ListFilter{
							{
								Field:    "ProjectId",
								Operator: "==",
								Value:    "projectID",
							},
							{
								Field:    "CreatedAt",
								Operator: ">=",
								Value:    target.Unix(),
							},
							{
								Field:    "CreatedAt",
								Operator: "<",
								Value:    targetNextWeek.Unix(),
							},
							{
								Field:    "ApplicationId",
								Operator: "==",
								Value:    "ApplicationId",
							},
						},
					}).Return([]*model.Deployment{
					{
						Id: "id1",
					},
					{
						Id: "id2",
					},
					{
						Id: "id3",
					},
				}, nil)

				return s
			}(),
			req: &webservice.GetInsightDataRequest{
				MetricsKind:    model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
				Step:           model.InsightStep_WEEKLY,
				RangeFrom:      time.Date(2020, 1, 5, 0, 0, 0, 0, time.Local).Unix(), // 2020/01/05 is Sunday
				DataPointCount: 2,
				ApplicationId:  "ApplicationId",
			},
			res: &webservice.GetInsightDataResponse{
				UpdatedAt: time.Now().Unix(),
				DataPoints: []*model.InsightDataPoint{
					{
						Value:     2,
						Timestamp: time.Date(2020, 1, 5, 0, 0, 0, 0, time.Local).Unix(),
					},
					{
						Value:     3,
						Timestamp: time.Date(2020, 1, 12, 0, 0, 0, 0, time.Local).Unix(),
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "valid with InsightStep_WEEKLY and rangeFrom is not Sunday",
			pipedID:   "pipedID",
			projectID: "projectID",
			deploymentStore: func() datastore.DeploymentStore {
				target := time.Date(2020, 1, 5, 0, 0, 0, 0, time.Local)
				targetNextWeek := target.AddDate(0, 0, 7)
				s := datastoretest.NewMockDeploymentStore(ctrl)
				s.EXPECT().
					ListDeployments(gomock.Any(), datastore.ListOptions{
						Filters: []datastore.ListFilter{
							{
								Field:    "ProjectId",
								Operator: "==",
								Value:    "projectID",
							},
							{
								Field:    "CreatedAt",
								Operator: ">=",
								Value:    target.Unix(),
							},
							{
								Field:    "CreatedAt",
								Operator: "<",
								Value:    targetNextWeek.Unix(),
							},
							{
								Field:    "ApplicationId",
								Operator: "==",
								Value:    "ApplicationId",
							},
						},
					}).Return([]*model.Deployment{
					{
						Id: "id1",
					},
					{
						Id: "id2",
					},
				}, nil)

				target = time.Date(2020, 1, 12, 0, 0, 0, 0, time.Local)
				targetNextWeek = target.AddDate(0, 0, 7)
				s.EXPECT().
					ListDeployments(gomock.Any(), datastore.ListOptions{
						Filters: []datastore.ListFilter{
							{
								Field:    "ProjectId",
								Operator: "==",
								Value:    "projectID",
							},
							{
								Field:    "CreatedAt",
								Operator: ">=",
								Value:    target.Unix(),
							},
							{
								Field:    "CreatedAt",
								Operator: "<",
								Value:    targetNextWeek.Unix(),
							},
							{
								Field:    "ApplicationId",
								Operator: "==",
								Value:    "ApplicationId",
							},
						},
					}).Return([]*model.Deployment{
					{
						Id: "id1",
					},
					{
						Id: "id2",
					},
					{
						Id: "id3",
					},
				}, nil)

				return s
			}(),
			req: &webservice.GetInsightDataRequest{
				MetricsKind:    model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
				Step:           model.InsightStep_WEEKLY,
				RangeFrom:      time.Date(2020, 1, 6, 0, 0, 0, 0, time.Local).Unix(), // 2020/01/06 is Monday
				DataPointCount: 2,
				ApplicationId:  "ApplicationId",
			},
			res: &webservice.GetInsightDataResponse{
				UpdatedAt: time.Now().Unix(),
				DataPoints: []*model.InsightDataPoint{
					{
						Value:     2,
						Timestamp: time.Date(2020, 1, 5, 0, 0, 0, 0, time.Local).Unix(),
					},
					{
						Value:     3,
						Timestamp: time.Date(2020, 1, 12, 0, 0, 0, 0, time.Local).Unix(),
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "valid with InsightStep_MONTHLY",
			pipedID:   "pipedID",
			projectID: "projectID",
			deploymentStore: func() datastore.DeploymentStore {
				target := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
				targetNextMonth := target.AddDate(0, 1, 0)
				s := datastoretest.NewMockDeploymentStore(ctrl)
				s.EXPECT().
					ListDeployments(gomock.Any(), datastore.ListOptions{
						Filters: []datastore.ListFilter{
							{
								Field:    "ProjectId",
								Operator: "==",
								Value:    "projectID",
							},
							{
								Field:    "CreatedAt",
								Operator: ">=",
								Value:    target.Unix(),
							},
							{
								Field:    "CreatedAt",
								Operator: "<",
								Value:    targetNextMonth.Unix(),
							},
							{
								Field:    "ApplicationId",
								Operator: "==",
								Value:    "ApplicationId",
							},
						},
					}).Return([]*model.Deployment{
					{
						Id: "id1",
					},
					{
						Id: "id2",
					},
				}, nil)

				target = time.Date(2020, 2, 1, 0, 0, 0, 0, time.Local)
				targetNextMonth = target.AddDate(0, 1, 0)
				s.EXPECT().
					ListDeployments(gomock.Any(), datastore.ListOptions{
						Filters: []datastore.ListFilter{
							{
								Field:    "ProjectId",
								Operator: "==",
								Value:    "projectID",
							},
							{
								Field:    "CreatedAt",
								Operator: ">=",
								Value:    target.Unix(),
							},
							{
								Field:    "CreatedAt",
								Operator: "<",
								Value:    targetNextMonth.Unix(),
							},
							{
								Field:    "ApplicationId",
								Operator: "==",
								Value:    "ApplicationId",
							},
						},
					}).Return([]*model.Deployment{
					{
						Id: "id1",
					},
					{
						Id: "id2",
					},
					{
						Id: "id3",
					},
				}, nil)

				return s
			}(),
			req: &webservice.GetInsightDataRequest{
				MetricsKind:    model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
				Step:           model.InsightStep_MONTHLY,
				RangeFrom:      time.Date(2020, 1, 4, 0, 0, 0, 0, time.Local).Unix(),
				DataPointCount: 2,
				ApplicationId:  "ApplicationId",
			},
			res: &webservice.GetInsightDataResponse{
				UpdatedAt: time.Now().Unix(),
				DataPoints: []*model.InsightDataPoint{
					{
						Value:     2,
						Timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local).Unix(),
					},
					{
						Value:     3,
						Timestamp: time.Date(2020, 2, 1, 0, 0, 0, 0, time.Local).Unix(),
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "valid with InsightStep_YEARLY",
			pipedID:   "pipedID",
			projectID: "projectID",
			deploymentStore: func() datastore.DeploymentStore {
				target := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
				targetNextYear := target.AddDate(1, 0, 0)
				s := datastoretest.NewMockDeploymentStore(ctrl)
				s.EXPECT().
					ListDeployments(gomock.Any(), datastore.ListOptions{
						Filters: []datastore.ListFilter{
							{
								Field:    "ProjectId",
								Operator: "==",
								Value:    "projectID",
							},
							{
								Field:    "CreatedAt",
								Operator: ">=",
								Value:    target.Unix(),
							},
							{
								Field:    "CreatedAt",
								Operator: "<",
								Value:    targetNextYear.Unix(),
							},
							{
								Field:    "ApplicationId",
								Operator: "==",
								Value:    "ApplicationId",
							},
						},
					}).Return([]*model.Deployment{
					{
						Id: "id1",
					},
					{
						Id: "id2",
					},
				}, nil)

				target = time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local)
				targetNextYear = target.AddDate(1, 0, 0)
				s.EXPECT().
					ListDeployments(gomock.Any(), datastore.ListOptions{
						Filters: []datastore.ListFilter{
							{
								Field:    "ProjectId",
								Operator: "==",
								Value:    "projectID",
							},
							{
								Field:    "CreatedAt",
								Operator: ">=",
								Value:    target.Unix(),
							},
							{
								Field:    "CreatedAt",
								Operator: "<",
								Value:    targetNextYear.Unix(),
							},
							{
								Field:    "ApplicationId",
								Operator: "==",
								Value:    "ApplicationId",
							},
						},
					}).Return([]*model.Deployment{
					{
						Id: "id1",
					},
					{
						Id: "id2",
					},
					{
						Id: "id3",
					},
				}, nil)

				return s
			}(),
			req: &webservice.GetInsightDataRequest{
				MetricsKind:    model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
				Step:           model.InsightStep_YEARLY,
				RangeFrom:      time.Date(2020, 1, 4, 0, 0, 0, 0, time.Local).Unix(),
				DataPointCount: 2,
				ApplicationId:  "ApplicationId",
			},
			res: &webservice.GetInsightDataResponse{
				UpdatedAt: time.Now().Unix(),
				DataPoints: []*model.InsightDataPoint{
					{
						Value:     2,
						Timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local).Unix(),
					},
					{
						Value:     3,
						Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local).Unix(),
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "return error when something wrong happen on ListDeployments",
			pipedID:   "pipedID",
			projectID: "projectID",
			deploymentStore: func() datastore.DeploymentStore {
				target := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
				targetNextYear := target.AddDate(1, 0, 0)
				s := datastoretest.NewMockDeploymentStore(ctrl)
				s.EXPECT().
					ListDeployments(gomock.Any(), datastore.ListOptions{
						Filters: []datastore.ListFilter{
							{
								Field:    "ProjectId",
								Operator: "==",
								Value:    "projectID",
							},
							{
								Field:    "CreatedAt",
								Operator: ">=",
								Value:    target.Unix(),
							},
							{
								Field:    "CreatedAt",
								Operator: "<",
								Value:    targetNextYear.Unix(),
							},
							{
								Field:    "ApplicationId",
								Operator: "==",
								Value:    "ApplicationId",
							},
						},
					}).Return([]*model.Deployment{}, fmt.Errorf("something wrong happens in ListDeployments"))
				return s
			}(),
			req: &webservice.GetInsightDataRequest{
				MetricsKind:    model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
				Step:           model.InsightStep_YEARLY,
				RangeFrom:      time.Date(2020, 1, 4, 0, 0, 0, 0, time.Local).Unix(),
				DataPointCount: 2,
				ApplicationId:  "ApplicationId",
			},
			res:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &WebAPI{
				pipedProjectCache: tt.pipedProjectCache,
				deploymentStore:   tt.deploymentStore,
				logger:            zap.NewNop(),
			}
			res, err := api.getInsightDataForDeployFrequency(ctx, tt.projectID, tt.req)
			assert.Equal(t, tt.wantErr, err != nil)
			if err == nil {
				assert.Equal(t, tt.res.DataPoints, res.DataPoints)
			}
		})
	}
}
