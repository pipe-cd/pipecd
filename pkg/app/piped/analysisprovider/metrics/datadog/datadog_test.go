// Copyright 2021 The PipeCD Authors.
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

package datadog

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/DataDog/datadog-api-client-go/api/v1/datadog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/piped/analysisprovider/metrics"
)

type fakeEvaluator struct {
	expected bool
}

func (f *fakeEvaluator) InRange(_ float64) bool {
	return f.expected
}

func (f *fakeEvaluator) String() string {
	return ""
}

func TestEvaluate(t *testing.T) {
	testcases := []struct {
		name      string
		evaluator metrics.Evaluator
		series    []datadog.MetricsQueryMetadata
		want      bool
		wantErr   bool
		errNoData bool
	}{
		{
			name:      "no data points found",
			evaluator: &fakeEvaluator{},
			series: []datadog.MetricsQueryMetadata{
				{
					Pointlist: nil,
				},
			},
			want:      false,
			wantErr:   true,
			errNoData: true,
		},
		{
			name:      "invalid data point format",
			evaluator: &fakeEvaluator{},
			series: []datadog.MetricsQueryMetadata{
				{
					Pointlist: &[][]float64{
						{0},
					},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name:      "out of range",
			evaluator: &fakeEvaluator{expected: false},
			series: []datadog.MetricsQueryMetadata{
				{
					Pointlist: &[][]float64{
						{0, 1},
					},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name:      "within the range",
			evaluator: &fakeEvaluator{expected: true},
			series: []datadog.MetricsQueryMetadata{
				{
					Pointlist: &[][]float64{
						{0, 1},
					},
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, _, err := evaluate(tc.evaluator, tc.series)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.errNoData, errors.Is(err, metrics.ErrNoDataFound))
		})
	}
}

func TestProviderQueryPoints(t *testing.T) {
	toInt64Pointer := func(i int64) *int64 { return &i }
	type queryResponse struct {
		res        datadog.MetricsQueryResponse
		httpStatus int
		err        error
	}
	testcases := []struct {
		name          string
		queryResponse queryResponse
		query         string
		queryRange    metrics.QueryRange
		want          []metrics.DataPoint
		wantErr       bool
	}{
		{
			name: "query failed",
			queryResponse: queryResponse{
				res:        datadog.MetricsQueryResponse{},
				httpStatus: http.StatusOK,
				err:        fmt.Errorf("query failed"),
			},
			query: "foo",
			queryRange: metrics.QueryRange{
				From: time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC),
				To:   time.Date(2009, time.January, 1, 0, 5, 0, 0, time.UTC),
			},
			wantErr: true,
		},
		{
			name: "unexpected HTTP status given",
			queryResponse: queryResponse{
				res:        datadog.MetricsQueryResponse{},
				httpStatus: http.StatusInternalServerError,
			},
			query: "foo",
			queryRange: metrics.QueryRange{
				From: time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC),
				To:   time.Date(2009, time.January, 1, 0, 5, 0, 0, time.UTC),
			},
			wantErr: true,
		},
		{
			name: "multiple data points given",
			queryResponse: queryResponse{
				res: datadog.MetricsQueryResponse{
					Series: &[]datadog.MetricsQueryMetadata{
						{
							Length: toInt64Pointer(2),
							Pointlist: &[][]float64{
								{1600000000, 0.1},
								{1600000001, 0.2},
							},
						},
					},
				},
				httpStatus: http.StatusOK,
			},
			query: "foo",
			queryRange: metrics.QueryRange{
				From: time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC),
				To:   time.Date(2009, time.January, 1, 0, 5, 0, 0, time.UTC),
			},
			want: []metrics.DataPoint{
				{Timestamp: 1600000000, Value: 0.1},
				{Timestamp: 1600000001, Value: 0.2},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			provider := Provider{
				client: datadog.NewAPIClient(datadog.NewConfiguration()),
				runQuery: func(_ datadog.ApiQueryMetricsRequest) (datadog.MetricsQueryResponse, *http.Response, error) {
					return tc.queryResponse.res, &http.Response{StatusCode: tc.queryResponse.httpStatus, Request: &http.Request{}}, tc.queryResponse.err
				},
				timeout: defaultTimeout,
				logger:  zap.NewNop(),
			}
			got, err := provider.QueryPoints(context.Background(), tc.query, tc.queryRange)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.want, got)
		})
	}
}
