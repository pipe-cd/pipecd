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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func Test_determineFilePaths(t *testing.T) {
	t.Parallel()

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
			want: []string{
				"insights/projectID/deployment_frequency/appID/2020-01.json",
				"insights/projectID/deployment_frequency/appID/2020-02.json",
			},
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
			want: []string{
				"insights/projectID/deployment_frequency/appID/2020-01.json",
				"insights/projectID/deployment_frequency/appID/2020-02.json",
			},
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
			want: []string{
				"insights/projectID/deployment_frequency/appID/2020-01.json",
				"insights/projectID/deployment_frequency/appID/2020-02.json",
			},
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := DetermineFilePaths(tt.args.projectID, tt.args.appID, tt.args.metricsKind, tt.args.step, tt.args.from, tt.args.dataPointCount)
			assert.Equal(t, tt.want, got)
		})
	}
}
