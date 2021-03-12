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
	"fmt"
	"net/http"
	"time"

	"github.com/DataDog/datadog-api-client-go/api/v1/datadog"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/piped/analysisprovider/metrics"
)

const (
	ProviderType   = "Datadog"
	defaultAddress = "datadoghq.com"
	defaultTimeout = 30 * time.Second
)

// Provider works as an HTTP client for datadog.
type Provider struct {
	client *datadog.APIClient

	address        string
	apiKey         string
	applicationKey string
	timeout        time.Duration
	logger         *zap.Logger
}

func NewProvider(apiKey, applicationKey string, opts ...Option) (*Provider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api-key is required")
	}
	if applicationKey == "" {
		return nil, fmt.Errorf("application-key is required")
	}

	p := &Provider{
		client:         datadog.NewAPIClient(datadog.NewConfiguration()),
		address:        defaultAddress,
		apiKey:         apiKey,
		applicationKey: applicationKey,
		timeout:        defaultTimeout,
		logger:         zap.NewNop(),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p, nil
}

type Option func(*Provider)

func WithAddress(address string) Option {
	if address == "" {
		address = defaultAddress
	}
	return func(p *Provider) {
		p.address = address
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(p *Provider) {
		p.logger = logger.Named("datadog-provider")
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(p *Provider) {
		p.timeout = timeout
	}
}

func (p *Provider) Type() string {
	return ProviderType
}

// RunQuery issues an HTTP request to the API named "MetricsApi.QueryMetrics", then evaluate its response.
// See more: https://docs.datadoghq.com/api/latest/metrics/#query-timeseries-points
func (p *Provider) RunQuery(ctx context.Context, query string, queryRange metrics.QueryRange, evaluator metrics.Evaluator) (bool, error) {
	if err := queryRange.Validate(); err != nil {
		return false, err
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()
	ctx = context.WithValue(
		ctx,
		datadog.ContextServerVariables,
		map[string]string{"site": p.address},
	)
	ctx = context.WithValue(
		ctx,
		datadog.ContextAPIKeys,
		map[string]datadog.APIKey{
			"apiKeyAuth": {
				Key: p.apiKey,
			},
			"appKeyAuth": {
				Key: p.applicationKey,
			},
		},
	)

	resp, httpResp, err := p.client.MetricsApi.QueryMetrics(ctx).
		From(queryRange.From.Unix()).
		To(queryRange.To.Unix()).
		Query(query).
		Execute()
	if err != nil {
		return false, fmt.Errorf("failed to call \"MetricsApi.QueryMetrics\": %w", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected HTTP status code from %s: %d", httpResp.Request.URL, httpResp.StatusCode)
	}
	if resp.Series == nil || len(*resp.Series) == 0 {
		return false, metrics.ErrNoValuesFound
	}
	return evaluate(evaluator, *resp.Series)
}

// evaluate checks if all data points for all time series are within the expected range.
func evaluate(evaluator metrics.Evaluator, series []datadog.MetricsQueryMetadata) (bool, error) {
	if err := evaluator.Validate(); err != nil {
		return false, err
	}

	for _, s := range series {
		points := s.Pointlist
		if points == nil || len(*points) == 0 {
			return false, fmt.Errorf("invalid response: no data points of the time series found")
		}
		for _, point := range *points {
			if len(point) < 2 {
				return false, fmt.Errorf("invalid response: invalid data point found")
			}
			// NOTE: A data point is assumed to be kind of like [unix-time, value].
			value := point[1]
			if !evaluator.InRange(value) {
				return false, nil
			}
		}
	}
	return true, nil
}
