package insightstore

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipe/pkg/model"
)

func Test_convertToInsightDataPoints(t *testing.T) {
	type args struct {
		report         Report
		from           time.Time
		dataPointCount int
		step           model.InsightStep
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
				report: func() Report {
					path := newYearlyFilePath("projectID", model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, "appID")
					expected := deployFrequencyReport{
						AccumulatedTo: 1609459200,
						Datapoints: deployFrequencyDataPoint{
							Yearly: map[string]deployFrequency{
								"2020": {DeployCount: 1000},
								"2021": {DeployCount: 3000},
							},
						},
						FilePath: path,
					}
					report, _ := toReport(&expected)
					return report
				}(),
				from:           time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				dataPointCount: 2,
				step:           model.InsightStep_YEARLY,
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
				report: func() Report {
					path := newYearlyFilePath("projectID", model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, "appID")
					expected := deployFrequencyReport{
						AccumulatedTo: 1609459200,
						Datapoints: deployFrequencyDataPoint{
							Monthly: map[string]deployFrequency{
								"2020-01": {DeployCount: 1000},
							},
						},
						FilePath: path,
					}
					report, _ := toReport(&expected)
					return report
				}(),
				from:           time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				dataPointCount: 1,
				step:           model.InsightStep_MONTHLY,
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
				report: func() Report {
					path := newYearlyFilePath("projectID", model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, "appID")
					expected := deployFrequencyReport{
						AccumulatedTo: 1609459200,
						Datapoints: deployFrequencyDataPoint{
							Weekly: map[string]deployFrequency{
								"2021-01-03": {DeployCount: 1000},
								"2021-01-10": {DeployCount: 3000},
							},
						},
						FilePath: path,
					}
					report, _ := toReport(&expected)
					return report
				}(),
				from:           time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
				dataPointCount: 2,
				step:           model.InsightStep_WEEKLY,
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
				report: func() Report {
					path := newYearlyFilePath("projectID", model.InsightMetricsKind_DEPLOYMENT_FREQUENCY, "appID")
					expected := deployFrequencyReport{
						AccumulatedTo: 1609459200,
						Datapoints: deployFrequencyDataPoint{
							Daily: map[string]deployFrequency{
								"2021-01-03": {DeployCount: 1000},
								"2021-01-04": {DeployCount: 3000},
							},
						},
						FilePath: path,
					}
					report, _ := toReport(&expected)
					return report
				}(),
				from:           time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
				dataPointCount: 2,
				step:           model.InsightStep_DAILY,
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
			got, err := convertToInsightDataPoints(tt.args.report, tt.args.from, tt.args.dataPointCount, tt.args.step)
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
