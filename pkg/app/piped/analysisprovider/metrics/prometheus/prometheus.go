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
	if timeout == 0 {
		timeout = defaultTimeout
	}
	return func(p *Provider) {
		p.timeout = timeout
	}
}

func WithLogger(logger *zap.Logger) Option {
	if logger == nil {
		return func(p *Provider) {}
	}
	return func(p *Provider) {
		p.logger = logger.Named("prometheus-provider")
	}
}

func (p *Provider) Type() string {
	return ProviderType
}

func (p *Provider) RunQuery(ctx context.Context, query string, evaluator metrics.Evaluator, queryRange metrics.QueryRange) (bool, error) {
	if err := queryRange.Validate(); err != nil {
		return false, err
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	p.logger.Info("run query", zap.String("query", query))
	// TODO: Use HTTP Basic Authentication with the username and password when needed.
	response, warnings, err := p.api.QueryRange(ctx, query, v1.Range{
		Start: queryRange.From,
		End:   queryRange.To,
		Step:  queryRange.Step,
	})
	if err != nil {
		return false, err
	}
	for _, w := range warnings {
		p.logger.Warn("non critical error occurred", zap.String("warning", w))
	}
	return p.evaluate(evaluator, response)
}

func (p *Provider) evaluate(evaluator metrics.Evaluator, response model.Value) (bool, error) {
	if err := evaluator.Validate(); err != nil {
		return false, err
	}

	// NOTE: Maybe it's enough to handle only matrix type as long as calling range queries endpoint.
	switch res := response.(type) {
	case *model.Scalar:
		value := float64(res.Value)
		if math.IsNaN(value) {
			return false, fmt.Errorf("the value is not a number")
		}
		return evaluator.InRange(value), nil
	case model.Vector:
		if len(res) == 0 {
			return false, fmt.Errorf("zero value in instant vector type returned")
		}
		// Check if all values are expected value.
		for _, s := range res {
			if s == nil {
				continue
			}
			value := float64(s.Value)
			if math.IsNaN(value) {
				return false, fmt.Errorf("the returned value is not a number")
			}
			if !evaluator.InRange(value) {
				return false, nil
			}
		}
		return true, nil
	case model.Matrix:
		if len(res) == 0 {
			return false, fmt.Errorf("no time series data points in range vector type")
		}
		// Check if all values are expected value.
		for _, r := range res {
			if len(r.Values) == 0 {
				return false, fmt.Errorf("zero value in range vector type returned")
			}
			for _, value := range r.Values {
				v := float64(value.Value)
				if math.IsNaN(v) {
					return false, fmt.Errorf("the returned value is not a number")
				}
				if !evaluator.InRange(v) {
					return false, nil
				}
			}
		}
		return true, nil
	default:
		return false, fmt.Errorf("unexpected data type returned")
	}
}
