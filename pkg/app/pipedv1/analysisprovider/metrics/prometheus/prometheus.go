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

package prometheus

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/analysisprovider/metrics"
)

const (
	ProviderType   = "Prometheus"
	defaultTimeout = 30 * time.Second
)

type client interface {
	QueryRange(ctx context.Context, query string, r v1.Range) (model.Value, v1.Warnings, error)
}

// Provider is a client for prometheus.
type Provider struct {
	api      client
	username string
	password string

	timeout time.Duration
	logger  *zap.Logger
}

func NewProvider(address string, opts ...Option) (*Provider, error) {
	if address == "" {
		return nil, fmt.Errorf("address is required")
	}

	p := &Provider{
		timeout: defaultTimeout,
		logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(p)
	}

	cfg := api.Config{
		Address: address,
	}
	if p.username != "" && p.password != "" {
		cfg.RoundTripper = config.NewBasicAuthRoundTripper(p.username, config.Secret(p.password), "", api.DefaultRoundTripper)
	}
	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	p.api = v1.NewAPI(client)
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

func WithBasicAuth(username, password string) Option {
	return func(p *Provider) {
		p.username = username
		p.password = password
	}
}

func (p *Provider) Type() string {
	return ProviderType
}

func (p *Provider) QueryPoints(ctx context.Context, query string, queryRange metrics.QueryRange) ([]metrics.DataPoint, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	if err := queryRange.Validate(); err != nil {
		return nil, err
	}
	// NOTE: Use 1m as a step but make sure the "step" is smaller than the query range.
	step := time.Minute
	if diff := queryRange.To.Sub(queryRange.From); diff < step {
		step = diff
	}

	p.logger.Info("run query", zap.String("query", query))
	response, warnings, err := p.api.QueryRange(ctx, query, v1.Range{
		Start: queryRange.From,
		End:   queryRange.To,
		Step:  step,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to run query for %s: %w", ProviderType, err)
	}
	for _, w := range warnings {
		p.logger.Warn("non critical error occurred", zap.String("warning", w))
	}

	// Collect data points given by the provider.
	// NOTE: Possibly, it's enough to handle only matrix type as long as calling range queries endpoint.
	switch res := response.(type) {
	case *model.Scalar:
		if math.IsNaN(float64(res.Value)) {
			return nil, fmt.Errorf("the value is not a number: %w", metrics.ErrNoDataFound)
		}
		return []metrics.DataPoint{
			{Timestamp: res.Timestamp.Unix(), Value: float64(res.Value)},
		}, nil
	case model.Vector:
		points := make([]metrics.DataPoint, 0, len(res))
		for _, s := range res {
			if s == nil {
				continue
			}
			if math.IsNaN(float64(s.Value)) {
				return nil, fmt.Errorf("the value is not a number: %w", metrics.ErrNoDataFound)
			}
			points = append(points, metrics.DataPoint{
				Timestamp: s.Timestamp.Unix(),
				Value:     float64(s.Value),
			})
		}
		return points, nil
	case model.Matrix:
		var size int
		for _, r := range res {
			size += len(r.Values)
		}
		points := make([]metrics.DataPoint, 0, size)
		for _, r := range res {
			if len(r.Values) == 0 {
				return nil, fmt.Errorf("zero value in range vector type returned: %w", metrics.ErrNoDataFound)
			}
			for _, point := range r.Values {
				if math.IsNaN(float64(point.Value)) {
					return nil, fmt.Errorf("the value is not a number: %w", metrics.ErrNoDataFound)
				}
				points = append(points, metrics.DataPoint{
					Timestamp: point.Timestamp.Unix(),
					Value:     float64(point.Value),
				})
			}
		}
		return points, nil
	default:
		return nil, fmt.Errorf("unexpected data type returned")
	}
}
