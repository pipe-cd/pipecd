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

package insightcollector

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipe/pkg/filestore/filestoretest"
	"github.com/pipe-cd/pipe/pkg/insight"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/datastore/datastoretest"

	"github.com/golang/mock/gomock"

	"github.com/pipe-cd/pipe/pkg/model"

	"github.com/pipe-cd/pipe/pkg/datastore"
)

func TestInsightCollector_getDailyInsightData(t *testing.T) {
	type args struct {
		projectID string
		appID     string
		kind      model.InsightMetricsKind
		rangeFrom time.Time
		rangeTo   time.Time
	}
	tests := []struct {
		name                   string
		prepareMockDataStoreFn func(m *datastoretest.MockDeploymentStore)
		args                   args
		want                   []insight.DataPoint
		wantAccumulatedTo      int64
		wantErr                bool
	}{
		{
			name: "Deploy Frequency / DAILY",
			prepareMockDataStoreFn: func(m *datastoretest.MockDeploymentStore) {
				m.EXPECT().ListDeployments(gomock.Any(), datastore.ListOptions{
					PageSize: 50,
					Filters: []datastore.ListFilter{
						{
							Field:    "CreatedAt",
							Operator: ">=",
							Value:    time.Date(2020, 10, 11, 4, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "CreatedAt",
							Operator: "<",
							Value:    time.Date(2020, 10, 12, 0, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "ApplicationId",
							Operator: "==",
							Value:    "appID",
						},
					},
					Orders: nil,
					Cursor: "",
				}).Return([]*model.Deployment{
					{
						Id:        "1",
						CreatedAt: time.Date(2020, 10, 11, 5, 0, 0, 0, time.UTC).Unix(),
					},
					{
						Id:        "2",
						CreatedAt: time.Date(2020, 10, 11, 5, 0, 0, 0, time.UTC).Unix(),
					},
					{
						Id:        "3",
						CreatedAt: time.Date(2020, 10, 11, 5, 0, 0, 0, time.UTC).Unix(),
					},
				}, nil)
				m.EXPECT().ListDeployments(gomock.Any(), datastore.ListOptions{
					PageSize: 50,
					Filters: []datastore.ListFilter{
						{
							Field:    "CreatedAt",
							Operator: ">=",
							Value:    time.Date(2020, 10, 11, 5, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "CreatedAt",
							Operator: "<",
							Value:    time.Date(2020, 10, 12, 0, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "ApplicationId",
							Operator: "==",
							Value:    "appID",
						},
					},
					Orders: nil,
					Cursor: "",
				}).Return([]*model.Deployment{}, nil)
				m.EXPECT().ListDeployments(gomock.Any(), datastore.ListOptions{
					PageSize: 50,
					Filters: []datastore.ListFilter{
						{
							Field:    "CreatedAt",
							Operator: ">=",
							Value:    time.Date(2020, 10, 12, 0, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "CreatedAt",
							Operator: "<",
							Value:    time.Date(2020, 10, 13, 0, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "ApplicationId",
							Operator: "==",
							Value:    "appID",
						},
					},
					Orders: nil,
					Cursor: "",
				}).Return([]*model.Deployment{
					{
						Id:        "1",
						CreatedAt: time.Date(2020, 10, 12, 1, 0, 0, 0, time.UTC).Unix(),
					},
					{
						Id:        "2",
						CreatedAt: time.Date(2020, 10, 12, 1, 0, 0, 0, time.UTC).Unix(),
					},
				}, nil)
				m.EXPECT().ListDeployments(gomock.Any(), datastore.ListOptions{
					PageSize: 50,
					Filters: []datastore.ListFilter{
						{
							Field:    "CreatedAt",
							Operator: ">=",
							Value:    time.Date(2020, 10, 12, 1, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "CreatedAt",
							Operator: "<",
							Value:    time.Date(2020, 10, 13, 0, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "ApplicationId",
							Operator: "==",
							Value:    "appID",
						},
					},
					Orders: nil,
					Cursor: "",
				}).Return([]*model.Deployment{}, nil)
				m.EXPECT().ListDeployments(gomock.Any(), datastore.ListOptions{
					PageSize: 50,
					Filters: []datastore.ListFilter{
						{
							Field:    "CreatedAt",
							Operator: ">=",
							Value:    time.Date(2020, 10, 13, 0, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "CreatedAt",
							Operator: "<",
							Value:    time.Date(2020, 10, 14, 0, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "ApplicationId",
							Operator: "==",
							Value:    "appID",
						},
					},
					Orders: nil,
					Cursor: "",
				}).Return([]*model.Deployment{
					{
						Id:        "1",
						CreatedAt: time.Date(2020, 10, 13, 1, 0, 0, 0, time.UTC).Unix(),
					},
				}, nil)
				m.EXPECT().ListDeployments(gomock.Any(), datastore.ListOptions{
					PageSize: 50,
					Filters: []datastore.ListFilter{
						{
							Field:    "CreatedAt",
							Operator: ">=",
							Value:    time.Date(2020, 10, 13, 1, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "CreatedAt",
							Operator: "<",
							Value:    time.Date(2020, 10, 14, 0, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "ApplicationId",
							Operator: "==",
							Value:    "appID",
						},
					},
					Orders: nil,
					Cursor: "",
				}).Return([]*model.Deployment{}, nil)
				m.EXPECT().ListDeployments(gomock.Any(), datastore.ListOptions{
					PageSize: 50,
					Filters: []datastore.ListFilter{
						{
							Field:    "CreatedAt",
							Operator: ">=",
							Value:    time.Date(2020, 10, 14, 0, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "CreatedAt",
							Operator: "<",
							Value:    time.Date(2020, 10, 15, 0, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "ApplicationId",
							Operator: "==",
							Value:    "appID",
						},
					},
					Orders: nil,
					Cursor: "",
				}).Return([]*model.Deployment{}, nil)
			},
			args: args{
				appID:     "appID",
				kind:      model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
				rangeFrom: time.Date(2020, 10, 11, 4, 0, 0, 0, time.UTC),
				rangeTo:   time.Date(2020, 10, 14, 0, 0, 0, 0, time.UTC),
			},
			want: func() []insight.DataPoint {
				daily := []*insight.DeployFrequency{
					{
						Timestamp:   time.Date(2020, 10, 11, 0, 0, 0, 0, time.UTC).Unix(),
						DeployCount: 3,
					},
					{
						Timestamp:   time.Date(2020, 10, 12, 0, 0, 0, 0, time.UTC).Unix(),
						DeployCount: 2,
					},
					{
						Timestamp:   time.Date(2020, 10, 13, 0, 0, 0, 0, time.UTC).Unix(),
						DeployCount: 1,
					},
				}
				dps, e := insight.ToDataPoints(daily)
				if e != nil {
					t.Fatalf("error when convert to data points: %v", e)
				}
				return dps
			}(),
			wantAccumulatedTo: time.Date(2020, 10, 13, 1, 0, 0, 0, time.UTC).Unix(),
			wantErr:           false,
		},
		{
			name: "Change Failure Rate/ DAILY",
			prepareMockDataStoreFn: func(m *datastoretest.MockDeploymentStore) {
				m.EXPECT().ListDeployments(gomock.Any(), datastore.ListOptions{
					PageSize: 50,
					Filters: []datastore.ListFilter{
						{
							Field:    "CompletedAt",
							Operator: ">=",
							Value:    time.Date(2020, 10, 11, 4, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "CompletedAt",
							Operator: "<",
							Value:    time.Date(2020, 10, 12, 0, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "ApplicationId",
							Operator: "==",
							Value:    "appID",
						},
					},
					Orders: nil,
					Cursor: "",
				}).Return([]*model.Deployment{
					{
						Status:      model.DeploymentStatus_DEPLOYMENT_FAILURE,
						CompletedAt: time.Date(2020, 10, 11, 5, 0, 0, 0, time.UTC).Unix(),
					},
					{
						Status:      model.DeploymentStatus_DEPLOYMENT_FAILURE,
						CompletedAt: time.Date(2020, 10, 11, 5, 0, 0, 0, time.UTC).Unix(),
					},
					{
						Status:      model.DeploymentStatus_DEPLOYMENT_SUCCESS,
						CompletedAt: time.Date(2020, 10, 11, 5, 0, 0, 0, time.UTC).Unix(),
					},
					{
						Status:      model.DeploymentStatus_DEPLOYMENT_SUCCESS,
						CompletedAt: time.Date(2020, 10, 11, 5, 0, 0, 0, time.UTC).Unix(),
					},
				}, nil)
				m.EXPECT().ListDeployments(gomock.Any(), datastore.ListOptions{
					PageSize: 50,
					Filters: []datastore.ListFilter{
						{
							Field:    "CompletedAt",
							Operator: ">=",
							Value:    time.Date(2020, 10, 11, 5, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "CompletedAt",
							Operator: "<",
							Value:    time.Date(2020, 10, 12, 0, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "ApplicationId",
							Operator: "==",
							Value:    "appID",
						},
					},
					Orders: nil,
					Cursor: "",
				}).Return([]*model.Deployment{}, nil)
				m.EXPECT().ListDeployments(gomock.Any(), datastore.ListOptions{
					PageSize: 50,
					Filters: []datastore.ListFilter{
						{
							Field:    "CompletedAt",
							Operator: ">=",
							Value:    time.Date(2020, 10, 12, 0, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "CompletedAt",
							Operator: "<",
							Value:    time.Date(2020, 10, 13, 0, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "ApplicationId",
							Operator: "==",
							Value:    "appID",
						},
					},
					Orders: nil,
					Cursor: "",
				}).Return([]*model.Deployment{
					{
						Status:      model.DeploymentStatus_DEPLOYMENT_FAILURE,
						CompletedAt: time.Date(2020, 10, 12, 5, 0, 0, 0, time.UTC).Unix(),
					},
					{
						Status:      model.DeploymentStatus_DEPLOYMENT_SUCCESS,
						CompletedAt: time.Date(2020, 10, 12, 5, 0, 0, 0, time.UTC).Unix(),
					},
					{
						Status:      model.DeploymentStatus_DEPLOYMENT_SUCCESS,
						CompletedAt: time.Date(2020, 10, 12, 5, 0, 0, 0, time.UTC).Unix(),
					},
					{
						Status:      model.DeploymentStatus_DEPLOYMENT_SUCCESS,
						CompletedAt: time.Date(2020, 10, 12, 5, 0, 0, 0, time.UTC).Unix(),
					},
				}, nil)
				m.EXPECT().ListDeployments(gomock.Any(), datastore.ListOptions{
					PageSize: 50,
					Filters: []datastore.ListFilter{
						{
							Field:    "CompletedAt",
							Operator: ">=",
							Value:    time.Date(2020, 10, 12, 5, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "CompletedAt",
							Operator: "<",
							Value:    time.Date(2020, 10, 13, 0, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "ApplicationId",
							Operator: "==",
							Value:    "appID",
						},
					},
					Orders: nil,
					Cursor: "",
				}).Return([]*model.Deployment{}, nil)
				m.EXPECT().ListDeployments(gomock.Any(), datastore.ListOptions{
					PageSize: 50,
					Filters: []datastore.ListFilter{
						{
							Field:    "CompletedAt",
							Operator: ">=",
							Value:    time.Date(2020, 10, 13, 0, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "CompletedAt",
							Operator: "<",
							Value:    time.Date(2020, 10, 14, 0, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "ApplicationId",
							Operator: "==",
							Value:    "appID",
						},
					},
					Orders: nil,
					Cursor: "",
				}).Return([]*model.Deployment{
					{
						Status:      model.DeploymentStatus_DEPLOYMENT_SUCCESS,
						CompletedAt: time.Date(2020, 10, 13, 8, 0, 0, 0, time.UTC).Unix(),
					},
				}, nil)
				m.EXPECT().ListDeployments(gomock.Any(), datastore.ListOptions{
					PageSize: 50,
					Filters: []datastore.ListFilter{
						{
							Field:    "CompletedAt",
							Operator: ">=",
							Value:    time.Date(2020, 10, 13, 8, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "CompletedAt",
							Operator: "<",
							Value:    time.Date(2020, 10, 14, 0, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "ApplicationId",
							Operator: "==",
							Value:    "appID",
						},
					},
					Orders: nil,
					Cursor: "",
				}).Return([]*model.Deployment{}, nil)
				m.EXPECT().ListDeployments(gomock.Any(), datastore.ListOptions{
					PageSize: 50,
					Filters: []datastore.ListFilter{
						{
							Field:    "CompletedAt",
							Operator: ">=",
							Value:    time.Date(2020, 10, 14, 0, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "CompletedAt",
							Operator: "<",
							Value:    time.Date(2020, 10, 15, 0, 0, 0, 0, time.UTC).Unix(),
						},
						{
							Field:    "ApplicationId",
							Operator: "==",
							Value:    "appID",
						},
					},
					Orders: nil,
					Cursor: "",
				}).Return([]*model.Deployment{}, nil)
			},
			args: args{
				appID:     "appID",
				kind:      model.InsightMetricsKind_CHANGE_FAILURE_RATE,
				rangeFrom: time.Date(2020, 10, 11, 4, 0, 0, 0, time.UTC),
				rangeTo:   time.Date(2020, 10, 14, 0, 0, 0, 0, time.UTC),
			},
			want: func() []insight.DataPoint {
				daily := []*insight.ChangeFailureRate{
					{
						Timestamp:    time.Date(2020, 10, 11, 0, 0, 0, 0, time.UTC).Unix(),
						Rate:         0.5,
						SuccessCount: 2,
						FailureCount: 2,
					},
					{
						Timestamp:    time.Date(2020, 10, 12, 0, 0, 0, 0, time.UTC).Unix(),
						Rate:         0.25,
						SuccessCount: 3,
						FailureCount: 1,
					},
					{
						Timestamp:    time.Date(2020, 10, 13, 0, 0, 0, 0, time.UTC).Unix(),
						Rate:         0,
						SuccessCount: 1,
						FailureCount: 0,
					},
				}
				dps, e := insight.ToDataPoints(daily)
				if e != nil {
					t.Fatalf("error when convert to data points: %v", e)
				}
				return dps
			}(),
			wantAccumulatedTo: time.Date(2020, 10, 13, 8, 0, 0, 0, time.UTC).Unix(),
			wantErr:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := datastoretest.NewMockDeploymentStore(ctrl)
			tt.prepareMockDataStoreFn(mock)

			a := &InsightCollector{
				applicationStore: nil,
				deploymentStore:  mock,
				insightstore:     insight.NewStore(filestoretest.NewMockStore(ctrl), nil),
				logger:           zap.NewNop(),
			}
			got, accumulatedTo, err := a.getDailyInsightData(context.Background(), tt.args.projectID, tt.args.appID, tt.args.kind, tt.args.rangeFrom, tt.args.rangeTo)
			if (err != nil) != tt.wantErr {
				if !tt.wantErr {
					assert.NoError(t, err)
					return
				}
				assert.Error(t, err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantAccumulatedTo, accumulatedTo)
		})
	}
}

func TestGetInsightDataForDeployFrequency(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	PageSizeForListDeployments := 50
	tests := []struct {
		name            string
		projectID       string
		applicationID   string
		targetRangeFrom time.Time
		targetRangeTo   time.Time
		targetTimestamp int64
		deploymentStore datastore.DeploymentStore
		dataPoint       *insight.DeployFrequency
		accumulateTo    time.Time
		wantErr         bool
	}{
		{
			name:            "valid with InsightStep_DAILY",
			applicationID:   "ApplicationId",
			targetRangeFrom: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			targetRangeTo:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			targetTimestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			deploymentStore: func() datastore.DeploymentStore {
				s := datastoretest.NewMockDeploymentStore(ctrl)
				s.EXPECT().
					ListDeployments(gomock.Any(), datastore.ListOptions{
						PageSize: PageSizeForListDeployments,
						Filters: []datastore.ListFilter{
							{
								Field:    "CreatedAt",
								Operator: ">=",
								Value:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
							},
							{
								Field:    "CreatedAt",
								Operator: "<",
								Value:    time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC).Unix(),
							},
							{
								Field:    "ApplicationId",
								Operator: "==",
								Value:    "ApplicationId",
							},
						},
					}).Return([]*model.Deployment{
					{
						Id:        "id1",
						CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC).Unix(),
					},
					{
						Id:        "id2",
						CreatedAt: time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC).Unix(),
					},
					{
						Id:        "id3",
						CreatedAt: time.Date(2020, 1, 1, 6, 0, 0, 0, time.UTC).Unix(),
					},
				}, nil)
				return s
			}(),
			dataPoint: &insight.DeployFrequency{
				Timestamp:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				DeployCount: 3,
			},
			accumulateTo: time.Date(2020, 1, 1, 6, 0, 0, 0, time.UTC),
			wantErr:      false,
		},
		{
			name:            "return error when something wrong happen on ListDeployments",
			applicationID:   "ApplicationId",
			targetRangeFrom: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			targetRangeTo:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			deploymentStore: func() datastore.DeploymentStore {
				s := datastoretest.NewMockDeploymentStore(ctrl)
				s.EXPECT().
					ListDeployments(gomock.Any(), datastore.ListOptions{
						PageSize: PageSizeForListDeployments,
						Filters: []datastore.ListFilter{
							{
								Field:    "CreatedAt",
								Operator: ">=",
								Value:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
							},
							{
								Field:    "CreatedAt",
								Operator: "<",
								Value:    time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC).Unix(),
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
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InsightCollector{
				deploymentStore: tt.deploymentStore,
				logger:          zap.NewNop(),
			}
			value, accumulatedTo, err := i.getInsightDataForDeployFrequency(ctx, tt.projectID, tt.applicationID, tt.targetTimestamp, tt.targetRangeFrom, tt.targetRangeTo)
			assert.Equal(t, tt.wantErr, err != nil)
			if err == nil {
				assert.Equal(t, tt.dataPoint, value)
				assert.Equal(t, tt.accumulateTo, accumulatedTo)
			}
		})
	}
}
func TestGetInsightDataForChangeFailureRate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	PageSizeForListDeployments := 50
	tests := []struct {
		name            string
		projectID       string
		applicationID   string
		targetRangeFrom time.Time
		targetRangeTo   time.Time
		targetTimestamp int64
		deploymentStore datastore.DeploymentStore
		dataPoint       *insight.ChangeFailureRate
		accumulatedTo   time.Time
		wantErr         bool
	}{
		{
			name:            "valid with InsightStep_DAILY and app id",
			applicationID:   "ApplicationId",
			targetRangeFrom: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			targetRangeTo:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			targetTimestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			deploymentStore: func() datastore.DeploymentStore {
				s := datastoretest.NewMockDeploymentStore(ctrl)
				s.EXPECT().
					ListDeployments(gomock.Any(), datastore.ListOptions{
						PageSize: PageSizeForListDeployments,
						Filters: []datastore.ListFilter{
							{
								Field:    "CompletedAt",
								Operator: ">=",
								Value:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
							},
							{
								Field:    "CompletedAt",
								Operator: "<",
								Value:    time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC).Unix(),
							},
							{
								Field:    "ApplicationId",
								Operator: "==",
								Value:    "ApplicationId",
							},
						},
					}).Return([]*model.Deployment{
					{
						Id:          "id1",
						Status:      model.DeploymentStatus_DEPLOYMENT_SUCCESS,
						CompletedAt: time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC).Unix(),
					},
					{
						Id:          "id2",
						Status:      model.DeploymentStatus_DEPLOYMENT_SUCCESS,
						CompletedAt: time.Date(2020, 1, 1, 5, 0, 0, 0, time.UTC).Unix(),
					},
					{
						Id:          "id3",
						Status:      model.DeploymentStatus_DEPLOYMENT_FAILURE,
						CompletedAt: time.Date(2020, 1, 1, 5, 0, 0, 0, time.UTC).Unix(),
					},
					{
						Id:          "id4",
						Status:      model.DeploymentStatus_DEPLOYMENT_SUCCESS,
						CompletedAt: time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC).Unix(),
					},
				}, nil)

				return s
			}(),
			dataPoint: &insight.ChangeFailureRate{
				Timestamp:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				Rate:         0.25,
				SuccessCount: 3,
				FailureCount: 1,
			},
			accumulatedTo: time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC),
			wantErr:       false,
		},
		{
			name:            "valid with InsightStep_DAILY and projectId",
			projectID:       "ProjectId",
			targetRangeFrom: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			targetRangeTo:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			targetTimestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			deploymentStore: func() datastore.DeploymentStore {
				s := datastoretest.NewMockDeploymentStore(ctrl)
				s.EXPECT().
					ListDeployments(gomock.Any(), datastore.ListOptions{
						PageSize: PageSizeForListDeployments,
						Filters: []datastore.ListFilter{
							{
								Field:    "CompletedAt",
								Operator: ">=",
								Value:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
							},
							{
								Field:    "CompletedAt",
								Operator: "<",
								Value:    time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC).Unix(),
							},
							{
								Field:    "ProjectId",
								Operator: "==",
								Value:    "ProjectId",
							},
						},
					}).Return([]*model.Deployment{
					{
						Id:          "id1",
						Status:      model.DeploymentStatus_DEPLOYMENT_SUCCESS,
						CompletedAt: time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC).Unix(),
					},
					{
						Id:          "id2",
						Status:      model.DeploymentStatus_DEPLOYMENT_SUCCESS,
						CompletedAt: time.Date(2020, 1, 1, 5, 0, 0, 0, time.UTC).Unix(),
					},
					{
						Id:          "id3",
						Status:      model.DeploymentStatus_DEPLOYMENT_FAILURE,
						CompletedAt: time.Date(2020, 1, 1, 5, 0, 0, 0, time.UTC).Unix(),
					},
					{
						Id:          "id4",
						Status:      model.DeploymentStatus_DEPLOYMENT_SUCCESS,
						CompletedAt: time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC).Unix(),
					},
				}, nil)

				return s
			}(),
			dataPoint: &insight.ChangeFailureRate{
				Timestamp:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				Rate:         0.25,
				SuccessCount: 3,
				FailureCount: 1,
			},
			accumulatedTo: time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC),
			wantErr:       false,
		},
		{
			name:            "return error when something wrong happen on ListDeployments",
			applicationID:   "ApplicationId",
			targetRangeFrom: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			targetRangeTo:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			deploymentStore: func() datastore.DeploymentStore {
				s := datastoretest.NewMockDeploymentStore(ctrl)
				s.EXPECT().
					ListDeployments(gomock.Any(), datastore.ListOptions{
						PageSize: PageSizeForListDeployments,
						Filters: []datastore.ListFilter{
							{
								Field:    "CompletedAt",
								Operator: ">=",
								Value:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
							},
							{
								Field:    "CompletedAt",
								Operator: "<",
								Value:    time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC).Unix(),
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
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InsightCollector{
				deploymentStore: tt.deploymentStore,
				logger:          zap.NewNop(),
			}
			value, accumulatedTo, err := i.getInsightDataForChangeFailureRate(ctx, tt.projectID, tt.applicationID, tt.targetTimestamp, tt.targetRangeFrom, tt.targetRangeTo)
			assert.Equal(t, tt.wantErr, err != nil)
			if err == nil {
				assert.Equal(t, tt.dataPoint, value)
				assert.Equal(t, tt.accumulatedTo, accumulatedTo)
			}
		})
	}
}

func TestInsightCollector_updateDataPoints(t *testing.T) {
	type args struct {
		chunk         insight.Chunk
		step          model.InsightStep
		updatedps     []insight.DataPoint
		accumulatedTo int64
	}
	tests := []struct {
		name    string
		args    args
		want    insight.Chunk
		wantErr bool
	}{
		{
			name: "success with daily and deploy frequency",
			args: args{
				chunk: func() insight.Chunk {
					df := &insight.DeployFrequencyChunk{
						AccumulatedTo: time.Date(2020, 10, 11, 1, 0, 0, 0, time.UTC).Unix(),
						DataPoints: insight.DeployFrequencyDataPoint{
							Daily: []*insight.DeployFrequency{
								{
									Timestamp:   time.Date(2020, 10, 10, 0, 0, 0, 0, time.UTC).Unix(),
									DeployCount: 10,
								},
							},
							Weekly:  nil,
							Monthly: nil,
							Yearly:  nil,
						},
						FilePath: "",
					}
					c, e := insight.ToChunk(df)
					if e != nil {
						t.Fatalf("error when convert to chunk: %v", e)
					}
					return c
				}(),
				step: model.InsightStep_DAILY,
				updatedps: func() []insight.DataPoint {
					daily := []*insight.DeployFrequency{
						{
							Timestamp:   time.Date(2020, 10, 11, 0, 0, 0, 0, time.UTC).Unix(),
							DeployCount: 3,
						},
						{
							Timestamp:   time.Date(2020, 10, 12, 0, 0, 0, 0, time.UTC).Unix(),
							DeployCount: 2,
						},
						{
							Timestamp:   time.Date(2020, 10, 13, 0, 0, 0, 0, time.UTC).Unix(),
							DeployCount: 1,
						},
					}
					dps, e := insight.ToDataPoints(daily)
					if e != nil {
						t.Fatalf("error when convert to data points: %v", e)
					}
					return dps
				}(),
				accumulatedTo: time.Date(2020, 10, 13, 1, 0, 0, 0, time.UTC).Unix(),
			},
			want: func() insight.Chunk {
				df := &insight.DeployFrequencyChunk{
					AccumulatedTo: time.Date(2020, 10, 13, 1, 0, 0, 0, time.UTC).Unix(),
					DataPoints: insight.DeployFrequencyDataPoint{
						Daily: []*insight.DeployFrequency{
							{
								Timestamp:   time.Date(2020, 10, 10, 0, 0, 0, 0, time.UTC).Unix(),
								DeployCount: 10,
							},
							{
								Timestamp:   time.Date(2020, 10, 11, 0, 0, 0, 0, time.UTC).Unix(),
								DeployCount: 3,
							},
							{
								Timestamp:   time.Date(2020, 10, 12, 0, 0, 0, 0, time.UTC).Unix(),
								DeployCount: 2,
							},
							{
								Timestamp:   time.Date(2020, 10, 13, 0, 0, 0, 0, time.UTC).Unix(),
								DeployCount: 1,
							},
						},
						Weekly:  nil,
						Monthly: nil,
						Yearly:  nil,
					},
					FilePath: "",
				}
				c, e := insight.ToChunk(df)
				if e != nil {
					t.Fatalf("error when convert to chunk: %v", e)
				}
				return c
			}(),
		},
		{
			name: "success with weekly and deploy frequency",
			args: args{
				chunk: func() insight.Chunk {
					df := &insight.DeployFrequencyChunk{
						AccumulatedTo: time.Date(2020, 10, 11, 1, 0, 0, 0, time.UTC).Unix(),
						DataPoints: insight.DeployFrequencyDataPoint{
							Weekly: []*insight.DeployFrequency{
								{
									Timestamp:   time.Date(2020, 10, 4, 0, 0, 0, 0, time.UTC).Unix(),
									DeployCount: 10,
								},
							},
							Daily:   nil,
							Monthly: nil,
							Yearly:  nil,
						},
						FilePath: "",
					}
					c, e := insight.ToChunk(df)
					if e != nil {
						t.Fatalf("error when convert to chunk: %v", e)
					}
					return c
				}(),
				step: model.InsightStep_WEEKLY,
				updatedps: func() []insight.DataPoint {
					df := []*insight.DeployFrequency{
						{
							Timestamp:   time.Date(2020, 10, 11, 0, 0, 0, 0, time.UTC).Unix(),
							DeployCount: 7,
						},
					}
					dps, e := insight.ToDataPoints(df)
					if e != nil {
						t.Fatalf("error when convert to data points: %v", e)
					}
					return dps
				}(),
				accumulatedTo: time.Date(2020, 10, 13, 3, 0, 0, 0, time.UTC).Unix(),
			},
			want: func() insight.Chunk {
				df := &insight.DeployFrequencyChunk{
					AccumulatedTo: time.Date(2020, 10, 13, 3, 0, 0, 0, time.UTC).Unix(),
					DataPoints: insight.DeployFrequencyDataPoint{
						Weekly: []*insight.DeployFrequency{
							{
								Timestamp:   time.Date(2020, 10, 4, 0, 0, 0, 0, time.UTC).Unix(),
								DeployCount: 10,
							},
							{
								Timestamp:   time.Date(2020, 10, 11, 0, 0, 0, 0, time.UTC).Unix(),
								DeployCount: 7,
							},
						},
						Daily:   nil,
						Monthly: nil,
						Yearly:  nil,
					},
					FilePath: "",
				}
				c, e := insight.ToChunk(df)
				if e != nil {
					t.Fatalf("error when convert to chunk: %v", e)
				}
				return c
			}(),
		},
		{
			name: "success with monthly and deploy frequency",
			args: args{
				chunk: func() insight.Chunk {
					df := &insight.DeployFrequencyChunk{
						AccumulatedTo: time.Date(2020, 10, 11, 1, 0, 0, 0, time.UTC).Unix(),
						DataPoints: insight.DeployFrequencyDataPoint{
							Monthly: []*insight.DeployFrequency{
								{
									Timestamp:   time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC).Unix(),
									DeployCount: 10,
								},
							},
							Daily:  nil,
							Weekly: nil,
							Yearly: nil,
						},
						FilePath: "",
					}
					c, e := insight.ToChunk(df)
					if e != nil {
						t.Fatalf("error when convert to chunk: %v", e)
					}
					return c
				}(),
				step: model.InsightStep_MONTHLY,
				updatedps: func() []insight.DataPoint {
					df := []*insight.DeployFrequency{
						{
							Timestamp:   time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC).Unix(),
							DeployCount: 3,
						},
						{
							Timestamp:   time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC).Unix(),
							DeployCount: 7,
						},
					}
					dps, e := insight.ToDataPoints(df)
					if e != nil {
						t.Fatalf("error when convert to data points: %v", e)
					}
					return dps
				}(),
				accumulatedTo: time.Date(2020, 11, 13, 3, 0, 0, 0, time.UTC).Unix(),
			},
			want: func() insight.Chunk {
				df := &insight.DeployFrequencyChunk{
					AccumulatedTo: time.Date(2020, 11, 13, 3, 0, 0, 0, time.UTC).Unix(),
					DataPoints: insight.DeployFrequencyDataPoint{
						Monthly: []*insight.DeployFrequency{
							{
								Timestamp:   time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC).Unix(),
								DeployCount: 13,
							},
							{
								Timestamp:   time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC).Unix(),
								DeployCount: 7,
							},
						},
						Daily:  nil,
						Weekly: nil,
						Yearly: nil,
					},
					FilePath: "",
				}
				c, e := insight.ToChunk(df)
				if e != nil {
					t.Fatalf("error when convert to chunk: %v", e)
				}
				return c
			}(),
		},
		{
			name: "success with yearly and deploy frequency",
			args: args{
				chunk: func() insight.Chunk {
					df := &insight.DeployFrequencyChunk{
						AccumulatedTo: time.Date(2020, 10, 11, 1, 0, 0, 0, time.UTC).Unix(),
						DataPoints: insight.DeployFrequencyDataPoint{
							Yearly: []*insight.DeployFrequency{
								{
									Timestamp:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
									DeployCount: 10,
								},
							},
							Daily:   nil,
							Weekly:  nil,
							Monthly: nil,
						},
						FilePath: "",
					}
					c, e := insight.ToChunk(df)
					if e != nil {
						t.Fatalf("error when convert to chunk: %v", e)
					}
					return c
				}(),
				step: model.InsightStep_YEARLY,
				updatedps: func() []insight.DataPoint {
					df := []*insight.DeployFrequency{
						{
							Timestamp:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
							DeployCount: 3,
						},
						{
							Timestamp:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
							DeployCount: 7,
						},
					}
					dps, e := insight.ToDataPoints(df)
					if e != nil {
						t.Fatalf("error when convert to data points: %v", e)
					}
					return dps
				}(),
				accumulatedTo: time.Date(2021, 1, 13, 3, 0, 0, 0, time.UTC).Unix(),
			},
			want: func() insight.Chunk {
				df := &insight.DeployFrequencyChunk{
					AccumulatedTo: time.Date(2021, 1, 13, 3, 0, 0, 0, time.UTC).Unix(),
					DataPoints: insight.DeployFrequencyDataPoint{
						Yearly: []*insight.DeployFrequency{
							{
								Timestamp:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
								DeployCount: 13,
							},
							{
								Timestamp:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
								DeployCount: 7,
							},
						},
						Daily:   nil,
						Weekly:  nil,
						Monthly: nil,
					},
					FilePath: "",
				}
				c, e := insight.ToChunk(df)
				if e != nil {
					t.Fatalf("error when convert to chunk: %v", e)
				}
				return c
			}(),
		},
		{
			name: "success with daily and change failure rate",
			args: args{
				chunk: func() insight.Chunk {
					df := &insight.ChangeFailureRateChunk{
						AccumulatedTo: time.Date(2020, 10, 11, 1, 0, 0, 0, time.UTC).Unix(),
						DataPoints: insight.ChangeFailureRateDataPoint{
							Daily: []*insight.ChangeFailureRate{
								{
									Timestamp:    time.Date(2020, 10, 10, 0, 0, 0, 0, time.UTC).Unix(),
									Rate:         0,
									SuccessCount: 10,
									FailureCount: 0,
								},
							},
							Weekly:  nil,
							Monthly: nil,
							Yearly:  nil,
						},
						FilePath: "",
					}
					c, e := insight.ToChunk(df)
					if e != nil {
						t.Fatalf("error when convert to chunk: %v", e)
					}
					return c
				}(),
				step: model.InsightStep_DAILY,
				updatedps: func() []insight.DataPoint {
					daily := []*insight.ChangeFailureRate{
						{
							Timestamp:    time.Date(2020, 10, 11, 0, 0, 0, 0, time.UTC).Unix(),
							Rate:         0.5,
							SuccessCount: 2,
							FailureCount: 2,
						},
						{
							Timestamp:    time.Date(2020, 10, 12, 0, 0, 0, 0, time.UTC).Unix(),
							Rate:         0.25,
							SuccessCount: 3,
							FailureCount: 1,
						},
						{
							Timestamp:    time.Date(2020, 10, 13, 0, 0, 0, 0, time.UTC).Unix(),
							Rate:         0,
							SuccessCount: 1,
							FailureCount: 0,
						},
					}
					dps, e := insight.ToDataPoints(daily)
					if e != nil {
						t.Fatalf("error when convert to data points: %v", e)
					}
					return dps
				}(),
				accumulatedTo: time.Date(2020, 10, 13, 8, 0, 0, 0, time.UTC).Unix(),
			},
			want: func() insight.Chunk {
				df := &insight.ChangeFailureRateChunk{
					AccumulatedTo: time.Date(2020, 10, 13, 8, 0, 0, 0, time.UTC).Unix(),
					DataPoints: insight.ChangeFailureRateDataPoint{
						Daily: []*insight.ChangeFailureRate{
							{
								Timestamp:    time.Date(2020, 10, 10, 0, 0, 0, 0, time.UTC).Unix(),
								Rate:         0,
								SuccessCount: 10,
								FailureCount: 0,
							},
							{
								Timestamp:    time.Date(2020, 10, 11, 0, 0, 0, 0, time.UTC).Unix(),
								Rate:         0.5,
								SuccessCount: 2,
								FailureCount: 2,
							},
							{
								Timestamp:    time.Date(2020, 10, 12, 0, 0, 0, 0, time.UTC).Unix(),
								Rate:         0.25,
								SuccessCount: 3,
								FailureCount: 1,
							},
							{
								Timestamp:    time.Date(2020, 10, 13, 0, 0, 0, 0, time.UTC).Unix(),
								Rate:         0,
								SuccessCount: 1,
								FailureCount: 0,
							},
						},
						Weekly:  nil,
						Monthly: nil,
						Yearly:  nil,
					},
					FilePath: "",
				}
				c, e := insight.ToChunk(df)
				if e != nil {
					t.Fatalf("error when convert to chunk: %v", e)
				}
				return c
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &InsightCollector{}
			got, err := a.updateDataPoints(tt.args.chunk, tt.args.step, tt.args.updatedps, tt.args.accumulatedTo)
			if (err != nil) != tt.wantErr {
				if !tt.wantErr {
					assert.NoError(t, err)
					return
				}
				assert.Error(t, err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
