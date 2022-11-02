// Copyright 2022 The PipeCD Authors.
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

// import (
// 	"context"
// 	"fmt"
// 	"testing"
// 	"time"

// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// 	"go.uber.org/zap"

// 	"github.com/pipe-cd/pipecd/pkg/datastore"
// 	"github.com/pipe-cd/pipecd/pkg/datastore/datastoretest"
// 	"github.com/pipe-cd/pipecd/pkg/filestore/filestoretest"
// 	"github.com/pipe-cd/pipecd/pkg/insight"
// 	"github.com/pipe-cd/pipecd/pkg/insight/insightstore"
// 	"github.com/pipe-cd/pipecd/pkg/model"
// )

// func TestUpdateDataPoints(t *testing.T) {
// 	type args struct {
// 		chunk         insight.Chunk
// 		step          model.InsightStep
// 		updatedps     []insight.DataPoint
// 		accumulatedTo int64
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    insight.Chunk
// 		wantErr bool
// 	}{
// 		{
// 			name: "success with daily and deploy frequency",
// 			args: args{
// 				chunk: func() insight.Chunk {
// 					df := &insight.DeployFrequencyChunk{
// 						AccumulatedTo: time.Date(2020, 10, 11, 1, 0, 0, 0, time.UTC).Unix(),
// 						DataPoints: insight.DeployFrequencyDataPoint{
// 							Daily: []*insight.DeployFrequency{
// 								{
// 									Timestamp:   time.Date(2020, 10, 10, 0, 0, 0, 0, time.UTC).Unix(),
// 									DeployCount: 10,
// 								},
// 							},
// 							Weekly:  nil,
// 							Monthly: nil,
// 							Yearly:  nil,
// 						},
// 						FilePath: "",
// 					}
// 					c, e := insight.ToChunk(df)
// 					if e != nil {
// 						t.Fatalf("error when convert to chunk: %v", e)
// 					}
// 					return c
// 				}(),
// 				step: model.InsightStep_DAILY,
// 				updatedps: func() []insight.DataPoint {
// 					daily := []*insight.DeployFrequency{
// 						{
// 							Timestamp:   time.Date(2020, 10, 11, 0, 0, 0, 0, time.UTC).Unix(),
// 							DeployCount: 3,
// 						},
// 						{
// 							Timestamp:   time.Date(2020, 10, 12, 0, 0, 0, 0, time.UTC).Unix(),
// 							DeployCount: 2,
// 						},
// 						{
// 							Timestamp:   time.Date(2020, 10, 13, 0, 0, 0, 0, time.UTC).Unix(),
// 							DeployCount: 1,
// 						},
// 					}
// 					dps, e := insight.ToDataPoints(daily)
// 					if e != nil {
// 						t.Fatalf("error when convert to data points: %v", e)
// 					}
// 					return dps
// 				}(),
// 				accumulatedTo: time.Date(2020, 10, 13, 1, 0, 0, 0, time.UTC).Unix(),
// 			},
// 			want: func() insight.Chunk {
// 				df := &insight.DeployFrequencyChunk{
// 					AccumulatedTo: time.Date(2020, 10, 13, 1, 0, 0, 0, time.UTC).Unix(),
// 					DataPoints: insight.DeployFrequencyDataPoint{
// 						Daily: []*insight.DeployFrequency{
// 							{
// 								Timestamp:   time.Date(2020, 10, 10, 0, 0, 0, 0, time.UTC).Unix(),
// 								DeployCount: 10,
// 							},
// 							{
// 								Timestamp:   time.Date(2020, 10, 11, 0, 0, 0, 0, time.UTC).Unix(),
// 								DeployCount: 3,
// 							},
// 							{
// 								Timestamp:   time.Date(2020, 10, 12, 0, 0, 0, 0, time.UTC).Unix(),
// 								DeployCount: 2,
// 							},
// 							{
// 								Timestamp:   time.Date(2020, 10, 13, 0, 0, 0, 0, time.UTC).Unix(),
// 								DeployCount: 1,
// 							},
// 						},
// 						Weekly:  nil,
// 						Monthly: nil,
// 						Yearly:  nil,
// 					},
// 					FilePath: "",
// 				}
// 				c, e := insight.ToChunk(df)
// 				if e != nil {
// 					t.Fatalf("error when convert to chunk: %v", e)
// 				}
// 				return c
// 			}(),
// 		},
// 		{
// 			name: "success with weekly and deploy frequency",
// 			args: args{
// 				chunk: func() insight.Chunk {
// 					df := &insight.DeployFrequencyChunk{
// 						AccumulatedTo: time.Date(2020, 10, 11, 1, 0, 0, 0, time.UTC).Unix(),
// 						DataPoints: insight.DeployFrequencyDataPoint{
// 							Weekly: []*insight.DeployFrequency{
// 								{
// 									Timestamp:   time.Date(2020, 10, 4, 0, 0, 0, 0, time.UTC).Unix(),
// 									DeployCount: 10,
// 								},
// 							},
// 							Daily:   nil,
// 							Monthly: nil,
// 							Yearly:  nil,
// 						},
// 						FilePath: "",
// 					}
// 					c, e := insight.ToChunk(df)
// 					if e != nil {
// 						t.Fatalf("error when convert to chunk: %v", e)
// 					}
// 					return c
// 				}(),
// 				step: model.InsightStep_WEEKLY,
// 				updatedps: func() []insight.DataPoint {
// 					df := []*insight.DeployFrequency{
// 						{
// 							Timestamp:   time.Date(2020, 10, 11, 0, 0, 0, 0, time.UTC).Unix(),
// 							DeployCount: 7,
// 						},
// 					}
// 					dps, e := insight.ToDataPoints(df)
// 					if e != nil {
// 						t.Fatalf("error when convert to data points: %v", e)
// 					}
// 					return dps
// 				}(),
// 				accumulatedTo: time.Date(2020, 10, 13, 3, 0, 0, 0, time.UTC).Unix(),
// 			},
// 			want: func() insight.Chunk {
// 				df := &insight.DeployFrequencyChunk{
// 					AccumulatedTo: time.Date(2020, 10, 13, 3, 0, 0, 0, time.UTC).Unix(),
// 					DataPoints: insight.DeployFrequencyDataPoint{
// 						Weekly: []*insight.DeployFrequency{
// 							{
// 								Timestamp:   time.Date(2020, 10, 4, 0, 0, 0, 0, time.UTC).Unix(),
// 								DeployCount: 10,
// 							},
// 							{
// 								Timestamp:   time.Date(2020, 10, 11, 0, 0, 0, 0, time.UTC).Unix(),
// 								DeployCount: 7,
// 							},
// 						},
// 						Daily:   nil,
// 						Monthly: nil,
// 						Yearly:  nil,
// 					},
// 					FilePath: "",
// 				}
// 				c, e := insight.ToChunk(df)
// 				if e != nil {
// 					t.Fatalf("error when convert to chunk: %v", e)
// 				}
// 				return c
// 			}(),
// 		},
// 		{
// 			name: "success with monthly and deploy frequency",
// 			args: args{
// 				chunk: func() insight.Chunk {
// 					df := &insight.DeployFrequencyChunk{
// 						AccumulatedTo: time.Date(2020, 10, 11, 1, 0, 0, 0, time.UTC).Unix(),
// 						DataPoints: insight.DeployFrequencyDataPoint{
// 							Monthly: []*insight.DeployFrequency{
// 								{
// 									Timestamp:   time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 									DeployCount: 10,
// 								},
// 							},
// 							Daily:  nil,
// 							Weekly: nil,
// 							Yearly: nil,
// 						},
// 						FilePath: "",
// 					}
// 					c, e := insight.ToChunk(df)
// 					if e != nil {
// 						t.Fatalf("error when convert to chunk: %v", e)
// 					}
// 					return c
// 				}(),
// 				step: model.InsightStep_MONTHLY,
// 				updatedps: func() []insight.DataPoint {
// 					df := []*insight.DeployFrequency{
// 						{
// 							Timestamp:   time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 							DeployCount: 3,
// 						},
// 						{
// 							Timestamp:   time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 							DeployCount: 7,
// 						},
// 					}
// 					dps, e := insight.ToDataPoints(df)
// 					if e != nil {
// 						t.Fatalf("error when convert to data points: %v", e)
// 					}
// 					return dps
// 				}(),
// 				accumulatedTo: time.Date(2020, 11, 13, 3, 0, 0, 0, time.UTC).Unix(),
// 			},
// 			want: func() insight.Chunk {
// 				df := &insight.DeployFrequencyChunk{
// 					AccumulatedTo: time.Date(2020, 11, 13, 3, 0, 0, 0, time.UTC).Unix(),
// 					DataPoints: insight.DeployFrequencyDataPoint{
// 						Monthly: []*insight.DeployFrequency{
// 							{
// 								Timestamp:   time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 								DeployCount: 13,
// 							},
// 							{
// 								Timestamp:   time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 								DeployCount: 7,
// 							},
// 						},
// 						Daily:  nil,
// 						Weekly: nil,
// 						Yearly: nil,
// 					},
// 					FilePath: "",
// 				}
// 				c, e := insight.ToChunk(df)
// 				if e != nil {
// 					t.Fatalf("error when convert to chunk: %v", e)
// 				}
// 				return c
// 			}(),
// 		},
// 		{
// 			name: "success with yearly and deploy frequency",
// 			args: args{
// 				chunk: func() insight.Chunk {
// 					df := &insight.DeployFrequencyChunk{
// 						AccumulatedTo: time.Date(2020, 10, 11, 1, 0, 0, 0, time.UTC).Unix(),
// 						DataPoints: insight.DeployFrequencyDataPoint{
// 							Yearly: []*insight.DeployFrequency{
// 								{
// 									Timestamp:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 									DeployCount: 10,
// 								},
// 							},
// 							Daily:   nil,
// 							Weekly:  nil,
// 							Monthly: nil,
// 						},
// 						FilePath: "",
// 					}
// 					c, e := insight.ToChunk(df)
// 					if e != nil {
// 						t.Fatalf("error when convert to chunk: %v", e)
// 					}
// 					return c
// 				}(),
// 				step: model.InsightStep_YEARLY,
// 				updatedps: func() []insight.DataPoint {
// 					df := []*insight.DeployFrequency{
// 						{
// 							Timestamp:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 							DeployCount: 3,
// 						},
// 						{
// 							Timestamp:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 							DeployCount: 7,
// 						},
// 					}
// 					dps, e := insight.ToDataPoints(df)
// 					if e != nil {
// 						t.Fatalf("error when convert to data points: %v", e)
// 					}
// 					return dps
// 				}(),
// 				accumulatedTo: time.Date(2021, 1, 13, 3, 0, 0, 0, time.UTC).Unix(),
// 			},
// 			want: func() insight.Chunk {
// 				df := &insight.DeployFrequencyChunk{
// 					AccumulatedTo: time.Date(2021, 1, 13, 3, 0, 0, 0, time.UTC).Unix(),
// 					DataPoints: insight.DeployFrequencyDataPoint{
// 						Yearly: []*insight.DeployFrequency{
// 							{
// 								Timestamp:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 								DeployCount: 13,
// 							},
// 							{
// 								Timestamp:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 								DeployCount: 7,
// 							},
// 						},
// 						Daily:   nil,
// 						Weekly:  nil,
// 						Monthly: nil,
// 					},
// 					FilePath: "",
// 				}
// 				c, e := insight.ToChunk(df)
// 				if e != nil {
// 					t.Fatalf("error when convert to chunk: %v", e)
// 				}
// 				return c
// 			}(),
// 		},
// 		{
// 			name: "success with daily and change failure rate",
// 			args: args{
// 				chunk: func() insight.Chunk {
// 					df := &insight.ChangeFailureRateChunk{
// 						AccumulatedTo: time.Date(2020, 10, 11, 1, 0, 0, 0, time.UTC).Unix(),
// 						DataPoints: insight.ChangeFailureRateDataPoint{
// 							Daily: []*insight.ChangeFailureRate{
// 								{
// 									Timestamp:    time.Date(2020, 10, 10, 0, 0, 0, 0, time.UTC).Unix(),
// 									Rate:         0,
// 									SuccessCount: 10,
// 									FailureCount: 0,
// 								},
// 							},
// 							Weekly:  nil,
// 							Monthly: nil,
// 							Yearly:  nil,
// 						},
// 						FilePath: "",
// 					}
// 					c, e := insight.ToChunk(df)
// 					if e != nil {
// 						t.Fatalf("error when convert to chunk: %v", e)
// 					}
// 					return c
// 				}(),
// 				step: model.InsightStep_DAILY,
// 				updatedps: func() []insight.DataPoint {
// 					daily := []*insight.ChangeFailureRate{
// 						{
// 							Timestamp:    time.Date(2020, 10, 11, 0, 0, 0, 0, time.UTC).Unix(),
// 							Rate:         0.5,
// 							SuccessCount: 2,
// 							FailureCount: 2,
// 						},
// 						{
// 							Timestamp:    time.Date(2020, 10, 12, 0, 0, 0, 0, time.UTC).Unix(),
// 							Rate:         0.25,
// 							SuccessCount: 3,
// 							FailureCount: 1,
// 						},
// 						{
// 							Timestamp:    time.Date(2020, 10, 13, 0, 0, 0, 0, time.UTC).Unix(),
// 							Rate:         0,
// 							SuccessCount: 1,
// 							FailureCount: 0,
// 						},
// 					}
// 					dps, e := insight.ToDataPoints(daily)
// 					if e != nil {
// 						t.Fatalf("error when convert to data points: %v", e)
// 					}
// 					return dps
// 				}(),
// 				accumulatedTo: time.Date(2020, 10, 13, 8, 0, 0, 0, time.UTC).Unix(),
// 			},
// 			want: func() insight.Chunk {
// 				df := &insight.ChangeFailureRateChunk{
// 					AccumulatedTo: time.Date(2020, 10, 13, 8, 0, 0, 0, time.UTC).Unix(),
// 					DataPoints: insight.ChangeFailureRateDataPoint{
// 						Daily: []*insight.ChangeFailureRate{
// 							{
// 								Timestamp:    time.Date(2020, 10, 10, 0, 0, 0, 0, time.UTC).Unix(),
// 								Rate:         0,
// 								SuccessCount: 10,
// 								FailureCount: 0,
// 							},
// 							{
// 								Timestamp:    time.Date(2020, 10, 11, 0, 0, 0, 0, time.UTC).Unix(),
// 								Rate:         0.5,
// 								SuccessCount: 2,
// 								FailureCount: 2,
// 							},
// 							{
// 								Timestamp:    time.Date(2020, 10, 12, 0, 0, 0, 0, time.UTC).Unix(),
// 								Rate:         0.25,
// 								SuccessCount: 3,
// 								FailureCount: 1,
// 							},
// 							{
// 								Timestamp:    time.Date(2020, 10, 13, 0, 0, 0, 0, time.UTC).Unix(),
// 								Rate:         0,
// 								SuccessCount: 1,
// 								FailureCount: 0,
// 							},
// 						},
// 						Weekly:  nil,
// 						Monthly: nil,
// 						Yearly:  nil,
// 					},
// 					FilePath: "",
// 				}
// 				c, e := insight.ToChunk(df)
// 				if e != nil {
// 					t.Fatalf("error when convert to chunk: %v", e)
// 				}
// 				return c
// 			}(),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := updateDataPoints(tt.args.chunk, tt.args.step, tt.args.updatedps, tt.args.accumulatedTo)
// 			if (err != nil) != tt.wantErr {
// 				if !tt.wantErr {
// 					assert.NoError(t, err)
// 					return
// 				}
// 				assert.Error(t, err, tt.wantErr)
// 				return
// 			}

// 			assert.Equal(t, tt.want, got)
// 		})
// 	}
// }
// func TestFindDeploymentsCreatedInRange(t *testing.T) {
// 	type args struct {
// 		from int64
// 		to   int64
// 	}
// 	tests := []struct {
// 		name                   string
// 		args                   args
// 		prepareMockDataStoreFn func(m *datastoretest.MockDeploymentStore)
// 		want                   []*model.Deployment
// 		wantErr                bool
// 	}{
// 		{
// 			name: "success with multiple pages",
// 			args: args{
// 				from: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 				to:   time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC).Unix(),
// 			},
// 			prepareMockDataStoreFn: func(m *datastoretest.MockDeploymentStore) {
// 				m.EXPECT().List(gomock.Any(), datastore.ListOptions{
// 					Limit: 50,
// 					Filters: []datastore.ListFilter{
// 						{
// 							Field:    "CreatedAt",
// 							Operator: datastore.OperatorGreaterThanOrEqual,
// 							Value:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 						},
// 						{
// 							Field:    "CreatedAt",
// 							Operator: datastore.OperatorLessThan,
// 							Value:    time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC).Unix(),
// 						},
// 					},
// 					Orders: []datastore.Order{
// 						{
// 							Field:     "CreatedAt",
// 							Direction: datastore.Desc,
// 						},
// 						{
// 							Field:     "Id",
// 							Direction: datastore.Asc,
// 						},
// 					},
// 				}).Return([]*model.Deployment{
// 					{
// 						Id:        "4",
// 						CreatedAt: time.Date(2020, 1, 3, 5, 0, 0, 0, time.UTC).Unix(),
// 					},
// 					{
// 						Id:        "5",
// 						CreatedAt: time.Date(2020, 1, 3, 5, 0, 0, 0, time.UTC).Unix(),
// 					},
// 					{
// 						Id:        "6",
// 						CreatedAt: time.Date(2020, 1, 3, 10, 0, 0, 0, time.UTC).Unix(),
// 					},
// 				}, "", nil)
// 				m.EXPECT().List(gomock.Any(), datastore.ListOptions{
// 					Limit: 50,
// 					Filters: []datastore.ListFilter{
// 						{
// 							Field:    "CreatedAt",
// 							Operator: datastore.OperatorGreaterThanOrEqual,
// 							Value:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 						},
// 						{
// 							Field:    "CreatedAt",
// 							Operator: datastore.OperatorLessThan,
// 							Value:    time.Date(2020, 1, 3, 10, 0, 0, 0, time.UTC).Unix(),
// 						},
// 					},
// 					Orders: []datastore.Order{
// 						{
// 							Field:     "CreatedAt",
// 							Direction: datastore.Desc,
// 						},
// 						{
// 							Field:     "Id",
// 							Direction: datastore.Asc,
// 						},
// 					},
// 				}).Return([]*model.Deployment{
// 					{
// 						Id:        "1",
// 						CreatedAt: time.Date(2020, 1, 1, 5, 0, 0, 0, time.UTC).Unix(),
// 					},
// 					{
// 						Id:        "2",
// 						CreatedAt: time.Date(2020, 1, 1, 5, 0, 0, 0, time.UTC).Unix(),
// 					},
// 					{
// 						Id:        "3",
// 						CreatedAt: time.Date(2020, 1, 3, 5, 0, 0, 0, time.UTC).Unix(),
// 					},
// 				}, "", nil)
// 				m.EXPECT().List(gomock.Any(), datastore.ListOptions{
// 					Limit: 50,
// 					Filters: []datastore.ListFilter{
// 						{
// 							Field:    "CreatedAt",
// 							Operator: datastore.OperatorGreaterThanOrEqual,
// 							Value:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 						},
// 						{
// 							Field:    "CreatedAt",
// 							Operator: datastore.OperatorLessThan,
// 							Value:    time.Date(2020, 1, 3, 5, 0, 0, 0, time.UTC).Unix(),
// 						},
// 					},
// 					Orders: []datastore.Order{
// 						{
// 							Field:     "CreatedAt",
// 							Direction: datastore.Desc,
// 						},
// 						{
// 							Field:     "Id",
// 							Direction: datastore.Asc,
// 						},
// 					},
// 					Cursor: "",
// 				}).Return([]*model.Deployment{}, "", nil)
// 			},
// 			want: []*model.Deployment{
// 				{
// 					Id:        "4",
// 					CreatedAt: time.Date(2020, 1, 3, 5, 0, 0, 0, time.UTC).Unix(),
// 				},
// 				{
// 					Id:        "5",
// 					CreatedAt: time.Date(2020, 1, 3, 5, 0, 0, 0, time.UTC).Unix(),
// 				},
// 				{
// 					Id:        "6",
// 					CreatedAt: time.Date(2020, 1, 3, 10, 0, 0, 0, time.UTC).Unix(),
// 				},
// 				{
// 					Id:        "1",
// 					CreatedAt: time.Date(2020, 1, 1, 5, 0, 0, 0, time.UTC).Unix(),
// 				},
// 				{
// 					Id:        "2",
// 					CreatedAt: time.Date(2020, 1, 1, 5, 0, 0, 0, time.UTC).Unix(),
// 				},
// 				{
// 					Id:        "3",
// 					CreatedAt: time.Date(2020, 1, 3, 5, 0, 0, 0, time.UTC).Unix(),
// 				},
// 			},
// 		},
// 		{
// 			name: "success with multiple pages",
// 			args: args{
// 				from: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 				to:   time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC).Unix(),
// 			},
// 			prepareMockDataStoreFn: func(m *datastoretest.MockDeploymentStore) {
// 				m.EXPECT().List(gomock.Any(), datastore.ListOptions{
// 					Limit: 50,
// 					Filters: []datastore.ListFilter{
// 						{
// 							Field:    "CreatedAt",
// 							Operator: datastore.OperatorGreaterThanOrEqual,
// 							Value:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 						},
// 						{
// 							Field:    "CreatedAt",
// 							Operator: datastore.OperatorLessThan,
// 							Value:    time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC).Unix(),
// 						},
// 					},
// 					Orders: []datastore.Order{
// 						{
// 							Field:     "CreatedAt",
// 							Direction: datastore.Desc,
// 						},
// 						{
// 							Field:     "Id",
// 							Direction: datastore.Asc,
// 						},
// 					},
// 				}).Return([]*model.Deployment{}, "", fmt.Errorf("something wrong happens in ListDeployments"))
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mock := datastoretest.NewMockDeploymentStore(ctrl)
// 			tt.prepareMockDataStoreFn(mock)

// 			a := &Collector{
// 				applicationStore: nil,
// 				deploymentStore:  mock,
// 				insightstore:     insightstore.NewStore(filestoretest.NewMockStore(ctrl)),
// 				logger:           zap.NewNop(),
// 			}
// 			got, err := a.findDeploymentsCreatedInRange(context.Background(), tt.args.from, tt.args.to)
// 			if (err != nil) != tt.wantErr {
// 				if !tt.wantErr {
// 					assert.NoError(t, err)
// 					return
// 				}
// 				assert.Error(t, err, tt.wantErr)
// 				return
// 			}
// 			assert.Equal(t, tt.want, got)
// 		})
// 	}
// }
// func TestFindDeploymentsCompletedInRange(t *testing.T) {
// 	type args struct {
// 		from int64
// 		to   int64
// 	}
// 	tests := []struct {
// 		name                   string
// 		args                   args
// 		prepareMockDataStoreFn func(m *datastoretest.MockDeploymentStore)
// 		want                   []*model.Deployment
// 		wantErr                bool
// 	}{
// 		{
// 			name: "success with multiple pages",
// 			args: args{
// 				from: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 				to:   time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC).Unix(),
// 			},
// 			prepareMockDataStoreFn: func(m *datastoretest.MockDeploymentStore) {
// 				m.EXPECT().List(gomock.Any(), datastore.ListOptions{
// 					Limit: 50,
// 					Filters: []datastore.ListFilter{
// 						{
// 							Field:    "CompletedAt",
// 							Operator: datastore.OperatorGreaterThanOrEqual,
// 							Value:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 						},
// 						{
// 							Field:    "CompletedAt",
// 							Operator: datastore.OperatorLessThan,
// 							Value:    time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC).Unix(),
// 						},
// 					},
// 					Orders: []datastore.Order{
// 						{
// 							Field:     "CompletedAt",
// 							Direction: datastore.Desc,
// 						},
// 						{
// 							Field:     "Id",
// 							Direction: datastore.Asc,
// 						},
// 					},
// 				}).Return([]*model.Deployment{
// 					{
// 						Id:          "4",
// 						CompletedAt: time.Date(2020, 1, 3, 5, 0, 0, 0, time.UTC).Unix(),
// 					},
// 					{
// 						Id:          "5",
// 						CompletedAt: time.Date(2020, 1, 3, 5, 0, 0, 0, time.UTC).Unix(),
// 					},
// 					{
// 						Id:          "6",
// 						CompletedAt: time.Date(2020, 1, 3, 10, 0, 0, 0, time.UTC).Unix(),
// 					},
// 				}, "", nil)
// 				m.EXPECT().List(gomock.Any(), datastore.ListOptions{
// 					Limit: 50,
// 					Filters: []datastore.ListFilter{
// 						{
// 							Field:    "CompletedAt",
// 							Operator: datastore.OperatorGreaterThanOrEqual,
// 							Value:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 						},
// 						{
// 							Field:    "CompletedAt",
// 							Operator: datastore.OperatorLessThan,
// 							Value:    time.Date(2020, 1, 3, 10, 0, 0, 0, time.UTC).Unix(),
// 						},
// 					},
// 					Orders: []datastore.Order{
// 						{
// 							Field:     "CompletedAt",
// 							Direction: datastore.Desc,
// 						},
// 						{
// 							Field:     "Id",
// 							Direction: datastore.Asc,
// 						},
// 					},
// 				}).Return([]*model.Deployment{
// 					{
// 						Id:          "1",
// 						CompletedAt: time.Date(2020, 1, 1, 5, 0, 0, 0, time.UTC).Unix(),
// 					},
// 					{
// 						Id:          "2",
// 						CompletedAt: time.Date(2020, 1, 1, 5, 0, 0, 0, time.UTC).Unix(),
// 					},
// 					{
// 						Id:          "3",
// 						CompletedAt: time.Date(2020, 1, 3, 5, 0, 0, 0, time.UTC).Unix(),
// 					},
// 				}, "", nil)
// 				m.EXPECT().List(gomock.Any(), datastore.ListOptions{
// 					Limit: 50,
// 					Filters: []datastore.ListFilter{
// 						{
// 							Field:    "CompletedAt",
// 							Operator: datastore.OperatorGreaterThanOrEqual,
// 							Value:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 						},
// 						{
// 							Field:    "CompletedAt",
// 							Operator: datastore.OperatorLessThan,
// 							Value:    time.Date(2020, 1, 3, 5, 0, 0, 0, time.UTC).Unix(),
// 						},
// 					},
// 					Orders: []datastore.Order{
// 						{
// 							Field:     "CompletedAt",
// 							Direction: datastore.Desc,
// 						},
// 						{
// 							Field:     "Id",
// 							Direction: datastore.Asc,
// 						},
// 					},
// 					Cursor: "",
// 				}).Return([]*model.Deployment{}, "", nil)
// 			},
// 			want: []*model.Deployment{
// 				{
// 					Id:          "4",
// 					CompletedAt: time.Date(2020, 1, 3, 5, 0, 0, 0, time.UTC).Unix(),
// 				},
// 				{
// 					Id:          "5",
// 					CompletedAt: time.Date(2020, 1, 3, 5, 0, 0, 0, time.UTC).Unix(),
// 				},
// 				{
// 					Id:          "6",
// 					CompletedAt: time.Date(2020, 1, 3, 10, 0, 0, 0, time.UTC).Unix(),
// 				},
// 				{
// 					Id:          "1",
// 					CompletedAt: time.Date(2020, 1, 1, 5, 0, 0, 0, time.UTC).Unix(),
// 				},
// 				{
// 					Id:          "2",
// 					CompletedAt: time.Date(2020, 1, 1, 5, 0, 0, 0, time.UTC).Unix(),
// 				},
// 				{
// 					Id:          "3",
// 					CompletedAt: time.Date(2020, 1, 3, 5, 0, 0, 0, time.UTC).Unix(),
// 				},
// 			},
// 		},
// 		{
// 			name: "success with multiple pages",
// 			args: args{
// 				from: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 				to:   time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC).Unix(),
// 			},
// 			prepareMockDataStoreFn: func(m *datastoretest.MockDeploymentStore) {
// 				m.EXPECT().List(gomock.Any(), datastore.ListOptions{
// 					Limit: 50,
// 					Filters: []datastore.ListFilter{
// 						{
// 							Field:    "CompletedAt",
// 							Operator: datastore.OperatorGreaterThanOrEqual,
// 							Value:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
// 						},
// 						{
// 							Field:    "CompletedAt",
// 							Operator: datastore.OperatorLessThan,
// 							Value:    time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC).Unix(),
// 						},
// 					},
// 					Orders: []datastore.Order{
// 						{
// 							Field:     "CompletedAt",
// 							Direction: datastore.Desc,
// 						},
// 						{
// 							Field:     "Id",
// 							Direction: datastore.Asc,
// 						},
// 					},
// 				}).Return([]*model.Deployment{}, "", fmt.Errorf("something wrong happens in ListDeployments"))
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mock := datastoretest.NewMockDeploymentStore(ctrl)
// 			tt.prepareMockDataStoreFn(mock)

// 			a := &Collector{
// 				applicationStore: nil,
// 				deploymentStore:  mock,
// 				insightstore:     insightstore.NewStore(filestoretest.NewMockStore(ctrl)),
// 				logger:           zap.NewNop(),
// 			}
// 			got, err := a.findDeploymentsCompletedInRange(context.Background(), tt.args.from, tt.args.to)
// 			if (err != nil) != tt.wantErr {
// 				if !tt.wantErr {
// 					assert.NoError(t, err)
// 					return
// 				}
// 				assert.Error(t, err, tt.wantErr)
// 				return
// 			}
// 			assert.Equal(t, tt.want, got)
// 		})
// 	}
// }

// func TestExtractDailyInsightDataPoints(t *testing.T) {
// 	type args struct {
// 		kind        model.InsightMetricsKind
// 		deployments []*model.Deployment
// 		rangeFrom   time.Time
// 		rangeTo     time.Time
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    []insight.DataPoint
// 		wantErr bool
// 	}{
// 		{
// 			name: "Deploy Frequency / DAILY",
// 			args: args{
// 				deployments: []*model.Deployment{
// 					{
// 						Id:        "1",
// 						CreatedAt: time.Date(2020, 10, 11, 5, 0, 0, 0, time.UTC).Unix(),
// 					},
// 					{
// 						Id:        "2",
// 						CreatedAt: time.Date(2020, 10, 11, 5, 0, 0, 0, time.UTC).Unix(),
// 					},
// 					{
// 						Id:        "3",
// 						CreatedAt: time.Date(2020, 10, 11, 5, 0, 0, 0, time.UTC).Unix(),
// 					},
// 					{
// 						Id:        "4",
// 						CreatedAt: time.Date(2020, 10, 12, 1, 0, 0, 0, time.UTC).Unix(),
// 					},
// 					{
// 						Id:        "5",
// 						CreatedAt: time.Date(2020, 10, 12, 1, 0, 0, 0, time.UTC).Unix(),
// 					},
// 					{
// 						Id:        "6",
// 						CreatedAt: time.Date(2020, 10, 13, 1, 0, 0, 0, time.UTC).Unix(),
// 					},
// 				},
// 				kind:      model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
// 				rangeFrom: time.Date(2020, 10, 11, 4, 0, 0, 0, time.UTC),
// 				rangeTo:   time.Date(2020, 10, 14, 0, 0, 0, 0, time.UTC),
// 			},
// 			want: func() []insight.DataPoint {
// 				daily := []*insight.DeployFrequency{
// 					{
// 						Timestamp:   time.Date(2020, 10, 11, 0, 0, 0, 0, time.UTC).Unix(),
// 						DeployCount: 3,
// 					},
// 					{
// 						Timestamp:   time.Date(2020, 10, 12, 0, 0, 0, 0, time.UTC).Unix(),
// 						DeployCount: 2,
// 					},
// 					{
// 						Timestamp:   time.Date(2020, 10, 13, 0, 0, 0, 0, time.UTC).Unix(),
// 						DeployCount: 1,
// 					},
// 				}
// 				dps, e := insight.ToDataPoints(daily)
// 				if e != nil {
// 					t.Fatalf("error when convert to data points: %v", e)
// 				}
// 				return dps
// 			}(),
// 			wantErr: false,
// 		},
// 		{
// 			name: "Change Failure Rate/ DAILY",
// 			args: args{
// 				deployments: []*model.Deployment{
// 					{
// 						Status:      model.DeploymentStatus_DEPLOYMENT_FAILURE,
// 						CompletedAt: time.Date(2020, 10, 11, 5, 0, 0, 0, time.UTC).Unix(),
// 					},
// 					{
// 						Status:      model.DeploymentStatus_DEPLOYMENT_FAILURE,
// 						CompletedAt: time.Date(2020, 10, 11, 5, 0, 0, 0, time.UTC).Unix(),
// 					},
// 					{
// 						Status:      model.DeploymentStatus_DEPLOYMENT_SUCCESS,
// 						CompletedAt: time.Date(2020, 10, 11, 5, 0, 0, 0, time.UTC).Unix(),
// 					},
// 					{
// 						Status:      model.DeploymentStatus_DEPLOYMENT_SUCCESS,
// 						CompletedAt: time.Date(2020, 10, 11, 5, 0, 0, 0, time.UTC).Unix(),
// 					},
// 					{
// 						Status:      model.DeploymentStatus_DEPLOYMENT_FAILURE,
// 						CompletedAt: time.Date(2020, 10, 12, 5, 0, 0, 0, time.UTC).Unix(),
// 					},
// 					{
// 						Status:      model.DeploymentStatus_DEPLOYMENT_SUCCESS,
// 						CompletedAt: time.Date(2020, 10, 12, 5, 0, 0, 0, time.UTC).Unix(),
// 					},
// 					{
// 						Status:      model.DeploymentStatus_DEPLOYMENT_SUCCESS,
// 						CompletedAt: time.Date(2020, 10, 12, 5, 0, 0, 0, time.UTC).Unix(),
// 					},
// 					{
// 						Status:      model.DeploymentStatus_DEPLOYMENT_SUCCESS,
// 						CompletedAt: time.Date(2020, 10, 12, 5, 0, 0, 0, time.UTC).Unix(),
// 					},
// 					{
// 						Status:      model.DeploymentStatus_DEPLOYMENT_SUCCESS,
// 						CompletedAt: time.Date(2020, 10, 13, 8, 0, 0, 0, time.UTC).Unix(),
// 					},
// 				},
// 				kind:      model.InsightMetricsKind_CHANGE_FAILURE_RATE,
// 				rangeFrom: time.Date(2020, 10, 11, 4, 0, 0, 0, time.UTC),
// 				rangeTo:   time.Date(2020, 10, 14, 0, 0, 0, 0, time.UTC),
// 			},
// 			want: func() []insight.DataPoint {
// 				daily := []*insight.ChangeFailureRate{
// 					{
// 						Timestamp:    time.Date(2020, 10, 11, 0, 0, 0, 0, time.UTC).Unix(),
// 						Rate:         0.5,
// 						SuccessCount: 2,
// 						FailureCount: 2,
// 					},
// 					{
// 						Timestamp:    time.Date(2020, 10, 12, 0, 0, 0, 0, time.UTC).Unix(),
// 						Rate:         0.25,
// 						SuccessCount: 3,
// 						FailureCount: 1,
// 					},
// 					{
// 						Timestamp:    time.Date(2020, 10, 13, 0, 0, 0, 0, time.UTC).Unix(),
// 						Rate:         0,
// 						SuccessCount: 1,
// 						FailureCount: 0,
// 					},
// 				}
// 				dps, e := insight.ToDataPoints(daily)
// 				if e != nil {
// 					t.Fatalf("error when convert to data points: %v", e)
// 				}
// 				return dps
// 			}(),
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := extractDailyInsightDataPoints(tt.args.deployments, tt.args.kind, tt.args.rangeFrom, tt.args.rangeTo)
// 			if (err != nil) != tt.wantErr {
// 				if !tt.wantErr {
// 					assert.NoError(t, err)
// 					return
// 				}
// 				assert.Error(t, err, tt.wantErr)
// 				return
// 			}
// 			assert.Equal(t, tt.want, got)
// 		})
// 	}
// }

// func TestGroupDeployments(t *testing.T) {

// 	var (
// 		d111 = &model.Deployment{
// 			Id:            "deployment-1-1-1",
// 			ProjectId:     "project-1",
// 			ApplicationId: "application-1-1",
// 		}
// 		d112 = &model.Deployment{
// 			Id:            "deployment-1-1-2",
// 			ProjectId:     "project-1",
// 			ApplicationId: "application-1-1",
// 		}
// 		d121 = &model.Deployment{
// 			Id:            "deployment-1-2-1",
// 			ProjectId:     "project-1",
// 			ApplicationId: "application-1-2",
// 		}
// 		d211 = &model.Deployment{
// 			Id:            "deployment-2-1-1",
// 			ProjectId:     "project-2",
// 			ApplicationId: "application-2-1",
// 		}
// 	)

// 	testcases := []struct {
// 		name        string
// 		deployments []*model.Deployment
// 		apps        map[string][]*model.Deployment
// 		projects    map[string][]*model.Deployment
// 	}{
// 		{
// 			name:     "no deployment",
// 			apps:     map[string][]*model.Deployment{},
// 			projects: map[string][]*model.Deployment{},
// 		},
// 		{
// 			name: "multiple deployments",
// 			deployments: []*model.Deployment{
// 				d111,
// 				d112,
// 				d121,
// 				d211,
// 			},
// 			apps: map[string][]*model.Deployment{
// 				"application-1-1": {
// 					d111,
// 					d112,
// 				},
// 				"application-1-2": {
// 					d121,
// 				},
// 				"application-2-1": {
// 					d211,
// 				},
// 			},
// 			projects: map[string][]*model.Deployment{
// 				"project-1": {
// 					d111,
// 					d112,
// 					d121,
// 				},
// 				"project-2": {
// 					d211,
// 				},
// 			},
// 		},
// 	}

// 	for _, tc := range testcases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			apps, projects := groupDeployments(tc.deployments)
// 			assert.Equal(t, tc.apps, apps)
// 			assert.Equal(t, tc.projects, projects)
// 		})
// 	}
// }
