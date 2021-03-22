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

package prometheus

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/piped/analysisprovider/metrics"
)

const (
	ProviderType   = "Prometheus"
	defaultTimeout = 30 * time.Second
)

// Provider is a client for prometheus.
type Provider struct {
	api v1.API
	//username string
	//password string

	timeout time.Duration
	logger  *zap.Logger
}

func NewProvider(address string, opts ...Option) (*Provider, error) {
	if address == "" {
		return nil, fmt.Errorf("address is required")
	}
	client, err := api.NewClient(api.Config{
		Address: address,
	})
	if err != nil {
		return nil, err
	}

	p := &Provider{
		api:     v1.NewAPI(client),
		timeout: defaultTimeout,
		logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p, nil
}

type Option func(*Provider)

func WithTimeout(timeout time.Duration) Option {
	return func(p *Provider) {
		p.timeout = timeout
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(p *Provider) {
		p.logger = logger.Named("prometheus-provider")
	}
}

func (p *Provider) Type() string {
	return ProviderType
}

// Evaluate queries the range query endpoint and checks if values in all data points are within the expected range.
// For the range query endpoint, see: https://prometheus.io/docs/prometheus/latest/querying/api/#range-queries
func (p *Provider) Evaluate(ctx context.Context, query string, queryRange metrics.QueryRange, evaluator metrics.Evaluator) (bool, string, error) {
	if err := queryRange.Validate(); err != nil {
		return false, "", err
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	// NOTE: Use 1m as a step but make sure the "step" isn't smaller than the query range.
	step := time.Minute
	if diff := queryRange.To.Sub(queryRange.From); diff < step {
		step = diff
	}

	p.logger.Info("run query", zap.String("query", query))
	// TODO: Use HTTP Basic Authentication with the username and password when needed.
	response, warnings, err := p.api.QueryRange(ctx, query, v1.Range{
		Start: queryRange.From,
		End:   queryRange.To,
		Step:  step,
	})
	if err != nil {
		return false, "", err
	}
	for _, w := range warnings {
		p.logger.Warn("non critical error occurred", zap.String("warning", w))
	}
	return evaluate(evaluator, response)
}

func evaluate(evaluator metrics.Evaluator, response model.Value) (bool, string, error) {
	evaluateValue := func(value float64) (bool, error) {
		if math.IsNaN(value) {
			return false, fmt.Errorf("the value is not a number")
		}
		return evaluator.InRange(value), nil
	}

	// NOTE: Maybe it's enough to handle only matrix type as long as calling range queries endpoint.
	switch res := response.(type) {
	case *model.Scalar:
		expected, err := evaluateValue(float64(res.Value))
		if err != nil {
			return false, "", err
		}
		if !expected {
			reason := fmt.Sprintf("found a value (%g) that is out of the expected range (%s)", float64(res.Value), evaluator)
			return false, reason, nil
		}
	case model.Vector:
		if len(res) == 0 {
			return false, "", fmt.Errorf("zero value in instant vector type returned")
		}
		// Check if all values are expected value.
		for _, s := range res {
			if s == nil {
				continue
			}
			expected, err := evaluateValue(float64(s.Value))
			if err != nil {
				return false, "", err
			}
			if !expected {
				reason := fmt.Sprintf("found a value (%g) that is out of the expected range (%s)", float64(s.Value), evaluator)
				return false, reason, nil
			}
		}
	case model.Matrix:
		if len(res) == 0 {
			return false, "", fmt.Errorf("no time series data points in range vector type")
		}
		// Check if all values are expected value.
		for _, r := range res {
			if len(r.Values) == 0 {
				return false, "", fmt.Errorf("zero value in range vector type returned")
			}
			for _, value := range r.Values {
				expected, err := evaluateValue(float64(value.Value))
				if err != nil {
					return false, "", err
				}
				if !expected {
					reason := fmt.Sprintf("found a value (%g) that is out of the expected range (%s)", float64(value.Value), evaluator)
					return false, reason, nil
				}
			}
		}
	default:
		return false, "", fmt.Errorf("unexpected data type returned")
	}

	reason := fmt.Sprintf("all values are within the expected range (%s)", evaluator)
	return true, reason, nil
}
