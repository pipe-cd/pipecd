// Copyright 2024 The PipeCD Authors.
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

package analysis

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/analysisprovider/metrics"
	"github.com/pipe-cd/pipecd/pkg/config"
)

type fakeMetricsProvider struct {
	points []metrics.DataPoint
	err    error
}

func (f *fakeMetricsProvider) Type() string { return "" }
func (f *fakeMetricsProvider) QueryPoints(_ context.Context, _ string, _ metrics.QueryRange) ([]metrics.DataPoint, error) {
	return f.points, f.err
}

type fakeLogPersister struct{}

func (l *fakeLogPersister) Write(_ []byte) (int, error)         { return 0, nil }
func (l *fakeLogPersister) Info(_ string)                       {}
func (l *fakeLogPersister) Infof(_ string, _ ...interface{})    {}
func (l *fakeLogPersister) Success(_ string)                    {}
func (l *fakeLogPersister) Successf(_ string, _ ...interface{}) {}
func (l *fakeLogPersister) Error(_ string)                      {}
func (l *fakeLogPersister) Errorf(_ string, _ ...interface{})   {}

func floatToPointer(n float64) *float64 { return &n }

func Test_metricsAnalyzer_analyzeWithThreshold(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name            string
		metricsAnalyzer *metricsAnalyzer
		want            bool
		wantErr         bool
	}{
		{
			name: "no expected field given",
			metricsAnalyzer: &metricsAnalyzer{
				id: "id",
				cfg: config.AnalysisMetrics{
					Provider: "provider",
					Query:    "query",
				},
				provider: &fakeMetricsProvider{},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "query failed",
			metricsAnalyzer: &metricsAnalyzer{
				id: "id",
				cfg: config.AnalysisMetrics{
					Provider: "provider",
					Query:    "query",
					Expected: config.AnalysisExpected{Max: floatToPointer(1)},
				},
				provider: &fakeMetricsProvider{
					err: fmt.Errorf("query failed"),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "there is a point outside the expected range",
			metricsAnalyzer: &metricsAnalyzer{
				id: "id",
				cfg: config.AnalysisMetrics{
					Provider: "provider",
					Query:    "query",
					Expected: config.AnalysisExpected{Max: floatToPointer(1)},
				},
				provider: &fakeMetricsProvider{
					points: []metrics.DataPoint{
						{Value: 0.9},
						{Value: 1.1},
						{Value: 0.8},
					},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "all points are expected ones",
			metricsAnalyzer: &metricsAnalyzer{
				id: "id",
				cfg: config.AnalysisMetrics{
					Provider: "provider",
					Query:    "query",
					Expected: config.AnalysisExpected{Max: floatToPointer(1)},
				},
				provider: &fakeMetricsProvider{
					points: []metrics.DataPoint{
						{Value: 0.9},
						{Value: 0.9},
						{Value: 0.8},
					},
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.metricsAnalyzer.logger = zap.NewNop()
			tc.metricsAnalyzer.logPersister = &fakeLogPersister{}
			got, err := tc.metricsAnalyzer.analyzeWithThreshold(context.Background())
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_metricsAnalyzer_compare(t *testing.T) {
	t.Parallel()

	type args struct {
		experiment []float64
		control    []float64
		deviation  string
	}
	testcases := []struct {
		name            string
		metricsAnalyzer *metricsAnalyzer
		args            args
		wantExpected    bool
		wantErr         bool
	}{
		{
			name: "empty data points given",
			metricsAnalyzer: &metricsAnalyzer{
				id: "id",
				cfg: config.AnalysisMetrics{
					Provider: "provider",
					Query:    "query",
				},
				provider:     &fakeMetricsProvider{},
				logger:       zap.NewNop(),
				logPersister: &fakeLogPersister{},
			},
			args: args{
				experiment: []float64{},
				control:    []float64{0.1, 0.2, 0.3, 0.4, 0.5},
				deviation:  "EITHER",
			},
			wantExpected: false,
			wantErr:      true,
		},
		{
			name: "no significance",
			metricsAnalyzer: &metricsAnalyzer{
				id: "id",
				cfg: config.AnalysisMetrics{
					Provider: "provider",
					Query:    "query",
				},
				provider:     &fakeMetricsProvider{},
				logger:       zap.NewNop(),
				logPersister: &fakeLogPersister{},
			},
			args: args{
				experiment: []float64{0.1, 0.2, 0.3, 0.4, 0.5},
				control:    []float64{0.1, 0.2, 0.3, 0.4, 0.5},
				deviation:  "EITHER",
			},
			wantExpected: true,
			wantErr:      false,
		},
		{
			name: "deviation on high direction as expected",
			metricsAnalyzer: &metricsAnalyzer{
				id: "id",
				cfg: config.AnalysisMetrics{
					Provider: "provider",
					Query:    "query",
				},
				provider:     &fakeMetricsProvider{},
				logger:       zap.NewNop(),
				logPersister: &fakeLogPersister{},
			},
			args: args{
				experiment: []float64{10.1, 10.2, 10.3, 10.4, 10.5},
				control:    []float64{0.1, 0.2, 0.3, 0.4, 0.5},
				deviation:  "LOW",
			},
			wantExpected: true,
			wantErr:      false,
		},
		{
			name: "deviation on low direction as expected",
			metricsAnalyzer: &metricsAnalyzer{
				id: "id",
				cfg: config.AnalysisMetrics{
					Provider: "provider",
					Query:    "query",
				},
				provider:     &fakeMetricsProvider{},
				logger:       zap.NewNop(),
				logPersister: &fakeLogPersister{},
			},
			args: args{
				experiment: []float64{0.1, 0.2, 0.3, 0.4, 0.5},
				control:    []float64{10.1, 10.2, 10.3, 10.4, 10.5},
				deviation:  "HIGH",
			},
			wantExpected: true,
			wantErr:      false,
		},
		{
			name: "deviation on high direction as unexpected",
			metricsAnalyzer: &metricsAnalyzer{
				id: "id",
				cfg: config.AnalysisMetrics{
					Provider: "provider",
					Query:    "query",
				},
				provider:     &fakeMetricsProvider{},
				logger:       zap.NewNop(),
				logPersister: &fakeLogPersister{},
			},
			args: args{
				experiment: []float64{10.1, 10.2, 10.3, 10.4, 10.5},
				control:    []float64{0.1, 0.2, 0.3, 0.4, 0.5},
				deviation:  "HIGH",
			},
			wantExpected: false,
			wantErr:      false,
		},
		{
			name: "deviation on low direction as unexpected",
			metricsAnalyzer: &metricsAnalyzer{
				id: "id",
				cfg: config.AnalysisMetrics{
					Provider: "provider",
					Query:    "query",
				},
				provider:     &fakeMetricsProvider{},
				logger:       zap.NewNop(),
				logPersister: &fakeLogPersister{},
			},
			args: args{
				experiment: []float64{0.1, 0.2, 0.3, 0.4, 0.5},
				control:    []float64{10.1, 10.2, 10.3, 10.4, 10.5},
				deviation:  "LOW",
			},
			wantExpected: false,
			wantErr:      false,
		},
		{
			name: "deviation as unexpected",
			metricsAnalyzer: &metricsAnalyzer{
				id: "id",
				cfg: config.AnalysisMetrics{
					Provider: "provider",
					Query:    "query",
				},
				provider:     &fakeMetricsProvider{},
				logger:       zap.NewNop(),
				logPersister: &fakeLogPersister{},
			},
			args: args{
				experiment: []float64{0.1, 0.2, 5.3, 0.2, 0.5},
				control:    []float64{0.1, 0.1, 0.1, 0.1, 0.1},
				deviation:  "EITHER",
			},
			wantExpected: false,
			wantErr:      false,
		},
		{
			name: "the data points is empty",
			metricsAnalyzer: &metricsAnalyzer{
				id: "id",
				cfg: config.AnalysisMetrics{
					Provider: "provider",
					Query:    "query",
				},
				provider:     &fakeMetricsProvider{},
				logger:       zap.NewNop(),
				logPersister: &fakeLogPersister{},
			},
			args: args{
				experiment: nil,
				control:    nil,
				deviation:  "EITHER",
			},
			wantExpected: true,
			wantErr:      false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.metricsAnalyzer.compare(tc.args.experiment, tc.args.control, tc.args.deviation)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.wantExpected, got)
		})
	}
}

func Test_metricsAnalyzer_renderQuery(t *testing.T) {
	t.Parallel()

	type args struct {
		queryTemplate     string
		variantCustomArgs map[string]string
		variant           string
	}
	testcases := []struct {
		name            string
		metricsAnalyzer *metricsAnalyzer
		args            args
		want            string
		wantErr         bool
	}{
		{
			name: "using only variant built in args",
			args: args{
				queryTemplate: `variant="{{ .Variant.Name }}"`,
				variant:       "canary",
			},
			metricsAnalyzer: &metricsAnalyzer{},
			want:            `variant="canary"`,
			wantErr:         false,
		},
		{
			name: "using variant and app built in args",
			args: args{
				queryTemplate: `variant="{{ .Variant.Name }}", app="{{ .App.Name }}"`,
				variant:       "canary",
			},
			metricsAnalyzer: &metricsAnalyzer{
				argsTemplate: argsTemplate{
					App: appArgs{
						Name: "app-1",
					},
				},
			},
			want:    `variant="canary", app="app-1"`,
			wantErr: false,
		},
		{
			name: "using variant and app built in and custom args",
			args: args{
				queryTemplate:     `variant="{{ .Variant.Name }}", app="{{ .App.Name }}", pod="{{ .VariantCustomArgs.pod }}", id="{{ .AppCustomArgs.id }}"`,
				variantCustomArgs: map[string]string{"pod": "1234"},
				variant:           "canary",
			},
			metricsAnalyzer: &metricsAnalyzer{
				argsTemplate: argsTemplate{
					App: appArgs{
						Name: "app-1",
					},
					AppCustomArgs: map[string]string{"id": "xxxx"},
				},
			},
			want:    `variant="canary", app="app-1", pod="1234", id="xxxx"`,
			wantErr: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.metricsAnalyzer.renderQuery(tc.args.queryTemplate, tc.args.variantCustomArgs, tc.args.variant)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.want, got)
		})
	}
}
