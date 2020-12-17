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

package insightstore

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipe/pkg/model"
)

func Test_ExtractDataPoints(t *testing.T) {
	type args struct {
		chunk Chunk
		from  time.Time
		to    time.Time
		step  model.InsightStep
	}
	tests := []struct {
		name    string
		args    args
		want    []*model.InsightDataPoint
		wantErr bool
	}{
		{
			name: "success with yearly",
			args: args{
				chunk: func() Chunk {
					path := makeYearsFilePath("projectID", model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, "appID")
					expected := DeployFrequencyChunk{
						AccumulatedTo: 1609459200,
						DataPoints: DeployFrequencyDataPoint{
							Yearly: []DeployFrequency{
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
					chunk, _ := toChunk(&expected)
					return chunk
				}(),
				from: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				step: model.InsightStep_YEARLY,
			},
			want: []*model.InsightDataPoint{
				{
					Timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
					Value:     1000,
				},
				{
					Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
					Value:     3000,
				},
			},
		},
		{
			name: "success with monthly",
			args: args{
				chunk: func() Chunk {
					path := makeYearsFilePath("projectID", model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, "appID")
					expected := DeployFrequencyChunk{
						AccumulatedTo: 1609459200,
						DataPoints: DeployFrequencyDataPoint{
							Monthly: []DeployFrequency{
								{
									DeployCount: 1000,
									Timestamp:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
								},
							},
						},
						FilePath: path,
					}
					chunk, _ := toChunk(&expected)
					return chunk
				}(),
				from: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				step: model.InsightStep_MONTHLY,
			},
			want: []*model.InsightDataPoint{
				{
					Timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
					Value:     1000,
				},
			},
		},
		{
			name: "success with weekly",
			args: args{
				chunk: func() Chunk {
					path := makeYearsFilePath("projectID", model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, "appID")
					expected := DeployFrequencyChunk{
						AccumulatedTo: 1609459200,
						DataPoints: DeployFrequencyDataPoint{
							Weekly: []DeployFrequency{
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
					chunk, _ := toChunk(&expected)
					return chunk
				}(),
				from: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC),
				step: model.InsightStep_WEEKLY,
			},
			want: []*model.InsightDataPoint{
				{
					Timestamp: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC).Unix(),
					Value:     1000,
				},
				{
					Timestamp: time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC).Unix(),
					Value:     3000,
				},
			},
		},
		{
			name: "success with daily",
			args: args{
				chunk: func() Chunk {
					path := makeYearsFilePath("projectID", model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, "appID")
					expected := DeployFrequencyChunk{
						AccumulatedTo: 1609459200,
						DataPoints: DeployFrequencyDataPoint{
							Daily: []DeployFrequency{
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
					chunk, _ := toChunk(&expected)
					return chunk
				}(),
				from: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC),
				step: model.InsightStep_DAILY,
			},
			want: []*model.InsightDataPoint{
				{
					Timestamp: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC).Unix(),
					Value:     1000,
				},
				{
					Timestamp: time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC).Unix(),
					Value:     3000,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractDataPoints(tt.args.chunk, tt.args.step, tt.args.from, tt.args.to)
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

func TestChunksToDataPoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		name           string
		projectID      string
		appID          string
		chunks         Chunks
		from           time.Time
		dataPointCount int
		fileCount      int
		step           model.InsightStep
		kind           model.InsightMetricsKind
		readerErr      error
		expected       []*model.InsightDataPoint
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
			chunks: func() []Chunk {
				path := makeChunkFilePath("projectID", model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, "appID", "2021-01")
				expected1 := DeployFrequencyChunk{
					AccumulatedTo: 1609459200,
					DataPoints: DeployFrequencyDataPoint{
						Daily: []DeployFrequency{
							{
								DeployCount: 1000,
								Timestamp:   1612051200,
							},
						},
					},
					FilePath: path,
				}
				chunk1, _ := toChunk(&expected1)
				path = makeChunkFilePath("projectID", model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, "appID", "2021-02")
				expected2 := DeployFrequencyChunk{
					AccumulatedTo: 1612123592,
					DataPoints: DeployFrequencyDataPoint{
						Daily: []DeployFrequency{
							{
								DeployCount: 3000,
								Timestamp:   1612137600,
							},
						},
					},
					FilePath: path,
				}
				chunk2, _ := toChunk(&expected2)
				return []Chunk{chunk1, chunk2}
			}(),
			expected: []*model.InsightDataPoint{
				{
					Timestamp: 1612051200,
					Value:     1000,
				},
				{
					Timestamp: 1612137600,
					Value:     3000,
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rs, err := tc.chunks.ExtractDataPoints(tc.step, tc.from, tc.dataPointCount)
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
