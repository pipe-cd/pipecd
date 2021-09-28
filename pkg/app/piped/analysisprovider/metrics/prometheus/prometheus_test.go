package prometheus

import (
	"context"
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/prometheus/common/model"
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

func TestType(t *testing.T) {
	p := Provider{}
	assert.Equal(t, ProviderType, p.Type())
}

func TestProviderEvaluate(t *testing.T) {
	cases := []struct {
		name       string
		queryError error
		wantErr    bool
	}{
		{
			name:       "query error occurred",
			queryError: fmt.Errorf("error"),
			wantErr:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := Provider{
				api: fakeClient{
					err: tc.queryError,
				},
				timeout: defaultTimeout,
				logger:  zap.NewNop(),
			}
			_, _, err := p.Evaluate(context.Background(), "query", metrics.QueryRange{From: time.Now()}, &fakeEvaluator{expected: true})
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}

}

func TestEvaluate(t *testing.T) {
	testcases := []struct {
		name      string
		evaluator metrics.Evaluator
		response  model.Value
		want      bool
		wantErr   bool
		errNoData bool
	}{
		{
			name:      "no data points found in the range vector response",
			evaluator: &fakeEvaluator{},
			response:  model.Matrix{},
			want:      false,
			wantErr:   true,
			errNoData: true,
		},
		{
			name:      "one of the instant vector within the range vector has no value",
			evaluator: &fakeEvaluator{},
			response: model.Matrix([]*model.SampleStream{
				{
					Values: nil,
				},
			}),
			want:      false,
			wantErr:   true,
			errNoData: true,
		},
		{
			name:      "NaN found in the range vector",
			evaluator: &fakeEvaluator{expected: false},
			response: model.Matrix([]*model.SampleStream{
				{
					Values: []model.SamplePair{
						{
							Value: model.SampleValue(math.NaN()),
						},
					},
				},
			}),
			want:      false,
			wantErr:   true,
			errNoData: true,
		},
		{
			name:      "value of the type range vector is out of range",
			evaluator: &fakeEvaluator{expected: false},
			response: model.Matrix([]*model.SampleStream{
				{
					Values: []model.SamplePair{
						{
							Value: 1,
						},
					},
				},
			}),
			want:    false,
			wantErr: false,
		},
		{
			name:      "value of the type range vector is within the expected range",
			evaluator: &fakeEvaluator{expected: true},
			response: model.Matrix([]*model.SampleStream{
				{
					Values: []model.SamplePair{
						{
							Value: 1,
						},
					},
				},
			}),
			want:    true,
			wantErr: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, _, err := evaluate(tc.evaluator, tc.response)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.errNoData, errors.Is(err, metrics.ErrNoDataFound))
		})
	}
}

func TestProviderQueryPoints(t *testing.T) {
	testcases := []struct {
		name       string
		client     client
		query      string
		queryRange metrics.QueryRange
		want       []metrics.DataPoint
		wantErr    bool
	}{
		{
			name: "query failed",
			client: &fakeClient{
				err: fmt.Errorf("query error"),
			},
			query: "foo",
			queryRange: metrics.QueryRange{
				From: time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC),
				To:   time.Date(2009, time.January, 1, 0, 5, 0, 0, time.UTC),
			},
			wantErr: true,
		},
		{
			name: "scalar data point returned",
			client: &fakeClient{
				value: &model.Scalar{Timestamp: model.Time(1600000000), Value: model.SampleValue(0.1)},
			},
			query: "foo",
			queryRange: metrics.QueryRange{
				From: time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC),
				To:   time.Date(2009, time.January, 1, 0, 5, 0, 0, time.UTC),
			},
			want: []metrics.DataPoint{
				{Timestamp: 1600000, Value: 0.1},
			},
		},
		{
			name: "vector data points returned",
			client: &fakeClient{
				value: model.Vector([]*model.Sample{
					{
						Timestamp: model.Time(1600000000),
						Value:     model.SampleValue(0.1),
					},
					{
						Timestamp: model.Time(1600001000),
						Value:     model.SampleValue(0.2),
					},
				}),
			},
			query: "foo",
			queryRange: metrics.QueryRange{
				From: time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC),
				To:   time.Date(2009, time.January, 1, 0, 5, 0, 0, time.UTC),
			},
			want: []metrics.DataPoint{
				{Timestamp: 1600000, Value: 0.1},
				{Timestamp: 1600001, Value: 0.2},
			},
		},
		{
			name: "matrix data points returned",
			client: &fakeClient{
				value: model.Matrix([]*model.SampleStream{
					{
						Values: []model.SamplePair{
							{
								Timestamp: model.Time(1600000000),
								Value:     model.SampleValue(0.1),
							},
							{
								Timestamp: model.Time(1600001000),
								Value:     model.SampleValue(0.2),
							},
						},
					},
				}),
			},
			query: "foo",
			queryRange: metrics.QueryRange{
				From: time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC),
				To:   time.Date(2009, time.January, 1, 0, 5, 0, 0, time.UTC),
			},
			want: []metrics.DataPoint{
				{Timestamp: 1600000, Value: 0.1},
				{Timestamp: 1600001, Value: 0.2},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			provider := &Provider{
				api:     tc.client,
				timeout: defaultTimeout,
				logger:  zap.NewNop(),
			}
			got, err := provider.QueryPoints(context.Background(), tc.query, tc.queryRange)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.want, got)
		})
	}
}
