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
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipe/pkg/filestore"
	"github.com/pipe-cd/pipe/pkg/filestore/filestoretest"
	"github.com/pipe-cd/pipe/pkg/model"
)

func TestGetInsightDataPoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := filestoretest.NewMockStore(ctrl)

	testcases := []struct {
		name           string
		content        string
		from           time.Time
		dataPointCount int
		step           model.InsightStep
		kind           model.InsightMetricsKind
		readerErr      error
		expected       []*model.InsightDataPoint
		expectedErr    error
	}{
		{
			name:           "file not found in filestore",
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
			step:           model.InsightStep_YEARLY,
			from:           time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			dataPointCount: 2,
			kind:           model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
			content: `{
				"accumulated_to": 1609459200,
				"datapoints": {
					"yearly": {
						"2020": {
							"deploy_count": 1000
						},
						"2021": {
							"deploy_count": 3000
						}
					}
				}
			}`,
			expected: []*model.InsightDataPoint{
				{
					Value:     1000,
					Timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},
				{
					Value:     3000,
					Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},
			},
		},
		{
			name:           "[deploy frequency] success in monthly",
			step:           model.InsightStep_MONTHLY,
			from:           time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			dataPointCount: 2,
			kind:           model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
			content: `{
				"accumulated_to": 1609459200,
				"datapoints": {
					"monthly": {
						"2020-01": {
							"deploy_count": 1000
						},
						"2020-02": {
							"deploy_count": 3000
						}
					}
				}
			}`,
			expected: []*model.InsightDataPoint{
				{
					Value:     1000,
					Timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},
				{
					Value:     3000,
					Timestamp: time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},
			},
		},
		{
			name:           "[deploy frequency] success in weekly",
			step:           model.InsightStep_WEEKLY,
			from:           time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
			dataPointCount: 2,
			kind:           model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
			content: `{
				"accumulated_to": 1609459200,
				"datapoints": {
					"weekly": {
						"2021-01-03": {
							"deploy_count": 1000
						},
						"2021-01-10": {
							"deploy_count": 3000
						}
					}
				}
			}`,
			expected: []*model.InsightDataPoint{
				{
					Value:     1000,
					Timestamp: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC).Unix(),
				},
				{
					Value:     3000,
					Timestamp: time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC).Unix(),
				},
			},
		},
		{
			name:           "[deploy frequency] success in daily",
			step:           model.InsightStep_DAILY,
			from:           time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
			dataPointCount: 2,
			kind:           model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
			content: `{
				"accumulated_to": 1609459200,
				"datapoints": {
					"daily": {
						"2021-01-03": {
							"deploy_count": 1000
						},
						"2021-01-04": {
							"deploy_count": 3000
						}
					}
				}
			}`,
			expected: []*model.InsightDataPoint{
				{
					Value:     1000,
					Timestamp: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC).Unix(),
				},
				{
					Value:     3000,
					Timestamp: time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC).Unix(),
				},
			},
		},
		{
			name:           "[change failure rate] success",
			step:           model.InsightStep_YEARLY,
			from:           time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			dataPointCount: 2,
			kind:           model.InsightMetricsKind_CHANGE_FAILURE_RATE,
			content: `{
				"accumulated_to": 1609459200,
				"datapoints": {
					"yearly": {
						"2020": {
							"rate": 0.75,
							"success_count": 1000,
							"failure_count": 3000
						},
						"2021": {
							"rate": 0.50,
							"success_count": 1000,
							"failure_count": 1000
						}
					}
				}
			}`,
			expected: []*model.InsightDataPoint{
				{
					Value:     0.75,
					Timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},
				{
					Value:     0.50,
					Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},
			},
		},
	}

	fs := insightFileStore{
		filestore: store,
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			obj := filestore.Object{
				Content: []byte(tc.content),
			}
			idps, err := fs.getInsightDataPoints(obj, tc.from, tc.dataPointCount, tc.step, tc.kind)
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
			if got := formatFrom(tt.args.from, tt.args.step); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("formatFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_insightFilePaths(t *testing.T) {
	type args struct {
		projectID      string
		appID          string
		from           time.Time
		dataPointCount int
		metricsKind    model.InsightMetricsKind
		step           model.InsightStep
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "return correct path with daily",
			args: args{
				projectID:      "projectID",
				appID:          "appID",
				from:           time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				dataPointCount: 2,
				metricsKind:    model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
				step:           model.InsightStep_DAILY,
			},
			want: []string{"insights/projectID/deployment_frequency/appID/2020-01.json"},
		},
		{
			name: "return correct path with daily and dates that straddles months",
			args: args{
				projectID:      "projectID",
				appID:          "appID",
				from:           time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				dataPointCount: 50,
				metricsKind:    model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
				step:           model.InsightStep_DAILY,
			},
			want: []string{"insights/projectID/deployment_frequency/appID/2020-01.json", "insights/projectID/deployment_frequency/appID/2020-02.json"},
		},
		{
			name: "return correct path with weekly",
			args: args{
				projectID:      "projectID",
				appID:          "appID",
				from:           time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC),
				dataPointCount: 2,
				metricsKind:    model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
				step:           model.InsightStep_WEEKLY,
			},
			want: []string{"insights/projectID/deployment_frequency/appID/2020-01.json"},
		},
		{
			name: "return correct path with weekly and weeks that straddles months",
			args: args{
				projectID:      "projectID",
				appID:          "appID",
				from:           time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC),
				dataPointCount: 6,
				metricsKind:    model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
				step:           model.InsightStep_WEEKLY,
			},
			want: []string{"insights/projectID/deployment_frequency/appID/2020-01.json", "insights/projectID/deployment_frequency/appID/2020-02.json"},
		},
		{
			name: "return correct path with monthly",
			args: args{
				projectID:      "projectID",
				appID:          "appID",
				from:           time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				dataPointCount: 2,
				metricsKind:    model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
				step:           model.InsightStep_MONTHLY,
			},
			want: []string{"insights/projectID/deployment_frequency/appID/2020-01.json", "insights/projectID/deployment_frequency/appID/2020-02.json"},
		},
		{
			name: "return correct path with yearly",
			args: args{
				projectID:      "projectID",
				appID:          "appID",
				from:           time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				dataPointCount: 2,
				metricsKind:    model.InsightMetricsKind_DEPLOYMENT_FREQUENCY,
				step:           model.InsightStep_YEARLY,
			},
			want: []string{"insights/projectID/deployment_frequency/appID/years.json"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := insightFilePaths(tt.args.projectID, tt.args.appID, tt.args.from, tt.args.dataPointCount, tt.args.metricsKind, tt.args.step); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("insightFilePaths() = %v, want %v", got, tt.want)
			}
		})
	}
}
