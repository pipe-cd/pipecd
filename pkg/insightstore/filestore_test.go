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
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipe/pkg/filestore"
	"github.com/pipe-cd/pipe/pkg/filestore/filestoretest"
	"github.com/pipe-cd/pipe/pkg/model"
)

func TestGetReports(t *testing.T) {
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
		expected       []Report
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
					"accumulated_to": 1609459200,
					"datapoints": {
						"daily": {
							"2021-01-31": {
								"deploy_count": 1000
							}
						}
					}
				}`,
				`{
					"accumulated_to": 1612123592,
					"datapoints": {
						"daily": {
							"2021-02-01": {
								"deploy_count": 3000
							}
						}
					}
				}`},
			expected: func() []Report {
				path := newMonthlyFilePath("projectID", model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, "appID", "2021-01")
				expected1 := DeployFrequencyReport{
					AccumulatedTo: 1609459200,
					Datapoints: DeployFrequencyDataPoint{
						Daily: map[string]DeployFrequency{
							"2021-01-31": {DeployCount: 1000},
						},
					},
					FilePath: path,
				}
				report1, _ := toReport(&expected1)
				path = newMonthlyFilePath("projectID", model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, "appID", "2021-02")
				expected2 := DeployFrequencyReport{
					AccumulatedTo: 1612123592,
					Datapoints: DeployFrequencyDataPoint{
						Daily: map[string]DeployFrequency{
							"2021-02-01": {DeployCount: 3000},
						},
					},
					FilePath: path,
				}
				report2, _ := toReport(&expected2)
				return []Report{report1, report2}
			}(),
		},
	}

	fs := insightFileStore{
		filestore: store,
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			paths := searchFilePaths(tc.projectID, tc.appID, tc.from, tc.dataPointCount, tc.kind, tc.step)
			if len(paths) != tc.fileCount {
				t.Fatalf("the count of path must be %d, but, %d", tc.fileCount, len(paths))
			}

			for i, c := range tc.contents {
				obj := filestore.Object{
					Content: []byte(c),
				}
				store.EXPECT().GetObject(context.TODO(), paths[i]).Return(obj, tc.readerErr)

			}

			rs, err := fs.GetReports(context.Background(), tc.projectID, tc.appID, tc.kind, tc.step, tc.from, tc.dataPointCount)
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

func TestGetReport(t *testing.T) {
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
		expected       Report
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
			expected: func() Report {
				path := newYearlyFilePath("projectID", model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, "appID")
				expected := DeployFrequencyReport{
					AccumulatedTo: 1609459200,
					Datapoints: DeployFrequencyDataPoint{
						Yearly: map[string]DeployFrequency{
							"2020": {DeployCount: 1000},
							"2021": {DeployCount: 3000},
						},
					},
					FilePath: path,
				}
				report, _ := toReport(&expected)
				return report
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
				"datapoints": {
					"monthly": {
						"2020-01": {
							"deploy_count": 1000
						}
					}
				}
			}`,
			expected: func() Report {
				path := newMonthlyFilePath("projectID", model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, "appID", "2020-01")
				expected := DeployFrequencyReport{
					AccumulatedTo: 1609459200,
					Datapoints: DeployFrequencyDataPoint{
						Monthly: map[string]DeployFrequency{
							"2020-01": {DeployCount: 1000},
						},
					},
					FilePath: path,
				}
				report, _ := toReport(&expected)
				return report
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
			expected: func() Report {
				path := newMonthlyFilePath("projectID", model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, "appID", "2021-01")
				expected := DeployFrequencyReport{
					AccumulatedTo: 1609459200,
					Datapoints: DeployFrequencyDataPoint{
						Weekly: map[string]DeployFrequency{
							"2021-01-03": {DeployCount: 1000},
							"2021-01-10": {DeployCount: 3000},
						},
					},
					FilePath: path,
				}
				report, _ := toReport(&expected)
				return report
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
			expected: func() Report {
				path := newMonthlyFilePath("projectID", model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, "appID", "2021-01")
				expected := DeployFrequencyReport{
					AccumulatedTo: 1609459200,
					Datapoints: DeployFrequencyDataPoint{
						Daily: map[string]DeployFrequency{
							"2021-01-03": {DeployCount: 1000},
							"2021-01-04": {DeployCount: 3000},
						},
					},
					FilePath: path,
				}
				report, _ := toReport(&expected)
				return report
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
			expected: func() Report {
				path := newYearlyFilePath("projectID", model.InsightMetricsKind_CHANGE_FAILURE_RATE, "appID")
				expected := ChangeFailureRateReport{
					AccumulatedTo: 1609459200,
					Datapoints: ChangeFailureRateDataPoint{
						Yearly: map[string]ChangeFailureRate{
							"2020": {Rate: 0.75, SuccessCount: 1000, FailureCount: 3000},
							"2021": {Rate: 0.50, SuccessCount: 1000, FailureCount: 1000},
						},
					},
					FilePath: path,
				}
				report, _ := toReport(&expected)
				return report
			}(),
		},
	}

	fs := insightFileStore{
		filestore: store,
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			path := searchFilePaths(tc.projectID, tc.appID, tc.from, tc.dataPointCount, tc.kind, tc.step)
			if len(path) != 1 {
				t.Fatalf("the count of path must be 1, but, %d", len(path))
			}
			obj := filestore.Object{
				Content: []byte(tc.content),
			}
			store.EXPECT().GetObject(context.TODO(), path[0]).Return(obj, tc.readerErr)
			idps, err := fs.getReport(context.Background(), path[0], tc.kind)
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
			got := formatFrom(tt.args.from, tt.args.step)
			assert.Equal(t, got, tt.want)
		})
	}
}
