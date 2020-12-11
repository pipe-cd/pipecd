package insightstore

import (
	"reflect"
	"testing"
	"time"

	"github.com/pipe-cd/pipe/pkg/model"
)

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
			if got := newFilePaths(tt.args.projectID, tt.args.appID, tt.args.from, tt.args.dataPointCount, tt.args.metricsKind, tt.args.step); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newFilePaths() = %v, want %v", got, tt.want)
			}
		})
	}
}
