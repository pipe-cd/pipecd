// Copyright 2026 The PipeCD Authors.
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

package executestage

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"

	"github.com/pipe-cd/piped-plugin-sdk-go/unit"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/analysis/analysisprovider/metrics"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/analysis/config"
)

type metricsProviderFunc func(context.Context, string, metrics.QueryRange) ([]metrics.DataPoint, error)

func (f metricsProviderFunc) Type() string { return "test" }
func (f metricsProviderFunc) QueryPoints(ctx context.Context, query string, queryRange metrics.QueryRange) ([]metrics.DataPoint, error) {
	return f(ctx, query, queryRange)
}

func TestAnalyzers_run(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		duration     time.Duration
		interval     time.Duration
		failureLimit int
		skipOnNoData bool
		unexpected   bool
		queryErr     error
		blockOnCall  int
		canceled     bool
		wantCalls    int
		wantErr      string
	}{
		{
			name: "duration shorter than interval", duration: time.Second, interval: time.Minute,
			wantErr: "no query completed",
		},
		{
			name: "duration already elapsed", interval: time.Minute,
			wantErr: "no query completed",
		},
		{
			name: "duration equal to interval", duration: time.Second, interval: time.Second,
			unexpected: true, wantCalls: -1, wantErr: "analysis '",
		},
		{
			name: "first query exceeds stage deadline", duration: 2500 * time.Millisecond, interval: time.Second,
			blockOnCall: 1, wantCalls: 1, wantErr: "no query completed",
		},
		{
			name: "provider reports cancellation at stage deadline", duration: 2500 * time.Millisecond, interval: time.Second,
			blockOnCall: 1, queryErr: context.Canceled, failureLimit: 1, wantCalls: 1, wantErr: "no query completed",
		},
		{
			name: "successful checks", duration: 2500 * time.Millisecond, interval: time.Second,
			wantCalls: 2,
		},
		{
			name: "deadline after a completed check", duration: 2500 * time.Millisecond, interval: time.Second,
			blockOnCall: 2, wantCalls: 2,
		},
		{
			name: "unexpected result", duration: 2500 * time.Millisecond, interval: time.Second,
			unexpected: true, wantCalls: 1, wantErr: "failure number exceeded the failure limit (0)",
		},
		{
			name: "tolerated failures", duration: 2500 * time.Millisecond, interval: time.Second,
			unexpected: true, failureLimit: 2, wantCalls: 2,
		},
		{
			name: "failure limit exceeded", duration: 3500 * time.Millisecond, interval: time.Second,
			unexpected: true, failureLimit: 1, wantCalls: 2, wantErr: "failure number exceeded the failure limit (1)",
		},
		{
			name: "skip no data", duration: 2500 * time.Millisecond, interval: time.Second,
			queryErr: metrics.ErrNoDataFound, skipOnNoData: true, wantCalls: 2,
		},
		{
			name: "reject no data", duration: 2500 * time.Millisecond, interval: time.Second,
			queryErr: metrics.ErrNoDataFound, wantCalls: 1, wantErr: "failure number exceeded the failure limit (0)",
		},
		{
			name: "query timeout before stage deadline", duration: 2500 * time.Millisecond, interval: time.Second,
			queryErr: context.DeadlineExceeded, wantCalls: 1, wantErr: "failure number exceeded the failure limit (0)",
		},
		{
			name: "provider error", duration: 2500 * time.Millisecond, interval: time.Second,
			queryErr: errors.New("provider unavailable"), wantCalls: 1, wantErr: "failure number exceeded the failure limit (0)",
		},
		{
			name: "provider error at stage deadline", duration: 2500 * time.Millisecond, interval: time.Second,
			queryErr: errors.New("provider unavailable"), blockOnCall: 2, failureLimit: 1,
			wantCalls: 2, wantErr: "failure number exceeded the failure limit (1)",
		},
		{
			name: "canceled by skip command", duration: time.Minute, interval: time.Second,
			canceled: true,
		},
	}

	for _, tc := range tests {
		for _, kind := range []string{"metrics", "http/log"} {
			t.Run(tc.name+"/"+kind, func(t *testing.T) {
				t.Parallel()
				synctest.Test(t, func(t *testing.T) {
					ctx, cancel := context.WithTimeout(t.Context(), tc.duration)
					defer cancel()
					if tc.canceled {
						cancel()
					}

					calls := 0
					evaluate := func(ctx context.Context) (bool, error) {
						calls++
						if calls == tc.blockOnCall {
							<-ctx.Done()
							if tc.queryErr != nil {
								return false, tc.queryErr
							}
							return false, ctx.Err()
						}
						return !tc.unexpected && tc.queryErr == nil, tc.queryErr
					}

					var err error
					if kind == "metrics" {
						provider := metricsProviderFunc(func(ctx context.Context, _ string, _ metrics.QueryRange) ([]metrics.DataPoint, error) {
							expected, err := evaluate(ctx)
							value := 0.5
							if !expected {
								value = 2
							}
							return []metrics.DataPoint{{Value: value}}, err
						})
						a := newMetricsAnalyzer("metrics-0", config.AnalysisMetrics{
							Strategy: config.AnalysisStrategyThreshold, Query: "query",
							Interval: unit.Duration(tc.interval), Expected: config.AnalysisExpected{Max: floatToPointer(1)},
							FailureLimit: tc.failureLimit, SkipOnNoData: tc.skipOnNoData,
						}, time.Now(), provider, nil, argsTemplate{}, zap.NewNop(), &fakeLogPersister{})
						err = a.run(ctx)
					} else {
						runner := func(ctx context.Context, query string) (bool, string, error) {
							expected, err := evaluate(ctx)
							return expected, query, err
						}
						a := newAnalyzer("http-0", "http", "query", runner, tc.interval, tc.failureLimit, tc.skipOnNoData, zap.NewNop(), &fakeLogPersister{})
						err = a.run(ctx)
					}

					if tc.wantErr == "" {
						assert.NoError(t, err)
					} else {
						assert.ErrorContains(t, err, tc.wantErr)
					}
					if tc.wantCalls < 0 {
						// At equal duration, either the ticker or the deadline may be selected first.
						assert.LessOrEqual(t, calls, 1)
					} else {
						assert.Equal(t, tc.wantCalls, calls)
					}
				})
			})
		}
	}
}
