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

package insight

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipe/pkg/filestore"
	"github.com/pipe-cd/pipe/pkg/filestore/filestoretest"
	"github.com/pipe-cd/pipe/pkg/model"
)

func TestGetChunks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := filestoretest.NewMockStore(ctrl)

	testcases := []struct {
		name           string
		projectID      string
		appID          string
		contents       []string
		from           time.Time
		dataPointCount int
		fileCount      int
		step           model.InsightStep
		kind           model.InsightMetricsKind
		readerErr      error
		expected       Chunks
		expectedErr    error
	}{
		{
			name:           "[deploy frequency] success in daily with dates that straddles months",
			projectID:      "projectID",
			appID:          "appID",
			step:           model.InsightStep_DAILY,
			from:           time.Date(2021, 1, 31, 0, 0, 0, 0, time.UTC),
			dataPointCount: 2,
			fileCount:      2,
			kind:           model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
			contents: []string{
				`{
					"accumulated_to": 1612051200,
					"data_points": {
						"daily": [
							{
								"deploy_count": 1000,
								"timestamp": 1612051200
							}
						]
					}
				}`,
				`{
					"accumulated_to": 1612137600,
					"data_points": {
						"daily": [
							{
								"deploy_count": 1000,
								"timestamp": 1612137600
							}
						]
					}
				}`},
			expected: func() []Chunk {
				path := makeChunkFilePath("projectID", model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, "appID", "2021-01")
				expected1 := DeployFrequencyChunk{
					AccumulatedTo: 1612051200,
					DataPoints: DeployFrequencyDataPoint{
						Daily: []*DeployFrequency{
							{
								DeployCount: 1000,
								Timestamp:   time.Date(2021, 1, 31, 0, 0, 0, 0, time.UTC).Unix(),
							},
						},
					},
					FilePath: path,
				}
				chunk1, _ := ToChunk(&expected1)
				path = makeChunkFilePath("projectID", model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, "appID", "2021-02")
				expected2 := DeployFrequencyChunk{
					AccumulatedTo: 1612137600,
					DataPoints: DeployFrequencyDataPoint{
						Daily: []*DeployFrequency{
							{
								DeployCount: 1000,
								Timestamp:   time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC).Unix(),
							},
						},
					},
					FilePath: path,
				}
				chunk2, _ := ToChunk(&expected2)
				return []Chunk{chunk1, chunk2}
			}(),
		},
	}

	fs := Store{
		filestore: store,
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			paths := determineFilePaths(tc.projectID, tc.appID, tc.kind, tc.step, tc.from, tc.dataPointCount)
			if len(paths) != tc.fileCount {
				t.Fatalf("the count of path must be %d, but, %d : &v", tc.fileCount, len(paths), paths)
			}

			for i, c := range tc.contents {
				obj := filestore.Object{
					Content: []byte(c),
				}
				store.EXPECT().GetObject(context.TODO(), paths[i]).Return(obj, tc.readerErr)

			}

			rs, err := fs.LoadChunks(context.Background(), tc.projectID, tc.appID, tc.kind, tc.step, tc.from, tc.dataPointCount)
			if err != nil {
				if tc.expectedErr == nil {
					assert.NoError(t, err)
					return
				}
				assert.Error(t, err, tc.expectedErr)
				return
			}
			assert.Equal(t, tc.expected, rs)
		})
	}
}

func TestGetChunk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := filestoretest.NewMockStore(ctrl)

	testcases := []struct {
		name           string
		projectID      string
		appID          string
		content        string
		from           time.Time
		dataPointCount int
		step           model.InsightStep
		kind           model.InsightMetricsKind
		readerErr      error
		expected       Chunk
		expectedErr    error
	}{
		{
			name:           "file not found in filestore",
			projectID:      "projectID",
			appID:          "appID",
			from:           time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			dataPointCount: 7,
			step:           model.InsightStep_DAILY,
			kind:           model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
			content:        "",
			readerErr:      filestore.ErrNotFound,
			expectedErr:    filestore.ErrNotFound,
		},
		{
			name:           "[deploy frequency] success in yearly",
			projectID:      "projectID",
			appID:          "appID",
			step:           model.InsightStep_YEARLY,
			from:           time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			dataPointCount: 2,
			kind:           model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
			content: `{
				"accumulated_to": 1609459200,
				"data_points": {
					"yearly": [
						{
							"deploy_count": 1000,
							"timestamp": 1577836800
						},
						{
							"deploy_count": 3000,
							"timestamp": 1609459200
						}
					]
				}
			}`,
			expected: func() Chunk {
				path := makeYearsFilePath("projectID", model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, "appID")
				expected := DeployFrequencyChunk{
					AccumulatedTo: 1609459200,
					DataPoints: DeployFrequencyDataPoint{
						Yearly: []*DeployFrequency{
							{
								DeployCount: 1000,
								Timestamp:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
							},
							{
								DeployCount: 3000,
								Timestamp:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
							},
						},
					},
					FilePath: path,
				}
				chunk, _ := ToChunk(&expected)
				return chunk
			}(),
		},
		{
			name:           "[deploy frequency] success in monthly",
			projectID:      "projectID",
			appID:          "appID",
			step:           model.InsightStep_MONTHLY,
			from:           time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			dataPointCount: 1,
			kind:           model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
			content: `{
				"accumulated_to": 1609459200,
				"data_points": {
					"monthly": [
						{
							"deploy_count": 1000,
							"timestamp": 1577836800
						}
					]
				}
			}`,
			expected: func() Chunk {
				path := makeChunkFilePath("projectID", model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, "appID", "2020-01")
				expected := DeployFrequencyChunk{
					AccumulatedTo: 1609459200,
					DataPoints: DeployFrequencyDataPoint{
						Monthly: []*DeployFrequency{
							{
								DeployCount: 1000,
								Timestamp:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
							},
						},
					},
					FilePath: path,
				}
				chunk, _ := ToChunk(&expected)
				return chunk
			}(),
		},
		{
			name:           "[deploy frequency] success in weekly",
			projectID:      "projectID",
			appID:          "appID",
			step:           model.InsightStep_WEEKLY,
			from:           time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
			dataPointCount: 2,
			kind:           model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
			content: `{
				"accumulated_to": 1609459200,
				"data_points": {
					"weekly": [
						{
							"deploy_count": 1000,
							"timestamp": 1609632000
						},
						{
							"deploy_count": 3000,
							"timestamp": 1610236800
						}
					]
				}
			}`,
			expected: func() Chunk {
				path := makeChunkFilePath("projectID", model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, "appID", "2021-01")
				expected := DeployFrequencyChunk{
					AccumulatedTo: 1609459200,
					DataPoints: DeployFrequencyDataPoint{
						Weekly: []*DeployFrequency{
							{
								DeployCount: 1000,
								Timestamp:   time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC).Unix(),
							},
							{
								DeployCount: 3000,
								Timestamp:   time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC).Unix(),
							},
						},
					},
					FilePath: path,
				}
				chunk, _ := ToChunk(&expected)
				return chunk
			}(),
		},
		{
			name:           "[deploy frequency] success in daily",
			projectID:      "projectID",
			appID:          "appID",
			step:           model.InsightStep_DAILY,
			from:           time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
			dataPointCount: 2,
			kind:           model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
			content: `{
				"accumulated_to": 1609459200,
				"data_points": {
					"daily": [
						{
							"deploy_count": 1000,
							"timestamp": 1609632000
						},
						{
							"deploy_count": 3000,
							"timestamp": 1609718400
						}
					]
				}
			}`,
			expected: func() Chunk {
				path := makeChunkFilePath("projectID", model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, "appID", "2021-01")
				expected := DeployFrequencyChunk{
					AccumulatedTo: 1609459200,
					DataPoints: DeployFrequencyDataPoint{
						Daily: []*DeployFrequency{
							{
								DeployCount: 1000,
								Timestamp:   time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC).Unix(),
							},
							{
								DeployCount: 3000,
								Timestamp:   time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC).Unix(),
							},
						},
					},
					FilePath: path,
				}
				chunk, _ := ToChunk(&expected)
				return chunk
			}(),
		},

		{
			name:           "[change failure rate] success",
			projectID:      "projectID",
			appID:          "appID",
			step:           model.InsightStep_YEARLY,
			from:           time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			dataPointCount: 2,
			kind:           model.InsightMetricsKind_CHANGE_FAILURE_RATE,
			content: `{
				"accumulated_to": 1609459200,
				"data_points": {
					"yearly": [
						{
							"rate": 0.75,
							"success_count": 1000,
							"failure_count": 3000,
							"timestamp": 1609632000
						},
						{
							"rate": 0.50,
							"success_count": 1000,
							"failure_count": 1000,
							"timestamp": 1609718400
						}
					]
				}
			}`,
			expected: func() Chunk {
				path := makeYearsFilePath("projectID", model.InsightMetricsKind_CHANGE_FAILURE_RATE, "appID")
				expected := ChangeFailureRateChunk{
					AccumulatedTo: 1609459200,
					DataPoints: ChangeFailureRateDataPoint{
						Yearly: []*ChangeFailureRate{
							{
								Rate:         0.75,
								SuccessCount: 1000,
								FailureCount: 3000,
								Timestamp:    time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC).Unix(),
							},
							{
								Rate:         0.50,
								SuccessCount: 1000,
								FailureCount: 1000,
								Timestamp:    time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC).Unix(),
							},
						},
					},
					FilePath: path,
				}
				chunk, _ := ToChunk(&expected)
				return chunk
			}(),
		},
	}

	fs := Store{
		filestore: store,
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			path := determineFilePaths(tc.projectID, tc.appID, tc.kind, tc.step, tc.from, tc.dataPointCount)
			if len(path) != 1 {
				t.Fatalf("the count of path must be 1, but, %d", len(path))
			}
			obj := filestore.Object{
				Content: []byte(tc.content),
			}
			store.EXPECT().GetObject(context.TODO(), path[0]).Return(obj, tc.readerErr)
			idps, err := fs.getChunk(context.Background(), path[0], tc.kind)
			if err != nil {
				if tc.expectedErr == nil {
					assert.NoError(t, err)
					return
				}
				assert.Error(t, err, tc.expectedErr)
				return
			}
			assert.Equal(t, tc.expected, idps)
		})
	}
}

func TestFormatFrom(t *testing.T) {
	type args struct {
		from time.Time
		step model.InsightStep
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "formatted correctly with daily",
			args: args{
				from: time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC),
				step: model.InsightStep_DAILY,
			},
			want: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "formatted correctly with weekly",
			args: args{
				from: time.Date(2020, 1, 10, 1, 1, 1, 1, time.UTC),
				step: model.InsightStep_WEEKLY,
			},
			want: time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "formatted correctly with monthly",
			args: args{
				from: time.Date(2020, 1, 7, 1, 1, 1, 1, time.UTC),
				step: model.InsightStep_MONTHLY,
			},
			want: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "formatted correctly with yearly",
			args: args{
				from: time.Date(2020, 7, 7, 1, 1, 1, 1, time.UTC),
				step: model.InsightStep_YEARLY,
			},
			want: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeTime(tt.args.from, tt.args.step)
			assert.Equal(t, got, tt.want)
		})
	}
}
