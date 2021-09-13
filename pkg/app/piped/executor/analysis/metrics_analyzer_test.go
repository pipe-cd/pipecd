package analysis

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/piped/analysisprovider/metrics"
	"github.com/pipe-cd/pipe/pkg/config"
)

type fakeMetricsProvider struct {
	points []metrics.DataPoint
	err    error
}

func (f *fakeMetricsProvider) Type() string { return "" }
func (f *fakeMetricsProvider) Evaluate(_ context.Context, _ string, _ metrics.QueryRange, _ metrics.Evaluator) (expected bool, reason string, err error) {
	return
}
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
