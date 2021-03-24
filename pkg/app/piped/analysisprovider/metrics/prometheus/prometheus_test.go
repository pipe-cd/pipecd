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
				api: fakeAPI{
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
