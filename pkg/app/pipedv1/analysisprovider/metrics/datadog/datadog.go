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

package datadog

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/DataDog/datadog-api-client-go/api/v1/datadog"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/analysisprovider/metrics"
)

const (
	ProviderType   = "Datadog"
	defaultAddress = "datadoghq.com"
	defaultTimeout = 30 * time.Second
)

// Provider works as an HTTP client for datadog.
type Provider struct {
	client   *datadog.APIClient
	runQuery func(request datadog.ApiQueryMetricsRequest) (datadog.MetricsQueryResponse, *http.Response, error)

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
		client: datadog.NewAPIClient(datadog.NewConfiguration()),
		runQuery: func(request datadog.ApiQueryMetricsRequest) (datadog.MetricsQueryResponse, *http.Response, error) {
			return request.Execute()
		},
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

func (p *Provider) QueryPoints(ctx context.Context, query string, queryRange metrics.QueryRange) ([]metrics.DataPoint, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	if err := queryRange.Validate(); err != nil {
		return nil, err
	}
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

	req := p.client.MetricsApi.QueryMetrics(ctx).
		From(queryRange.From.Unix()).
		To(queryRange.To.Unix()).
		Query(query)
	resp, httpResp, err := p.runQuery(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call \"MetricsApi.QueryMetrics\": %w", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status code from %s: %d", httpResp.Request.URL, httpResp.StatusCode)
	}

	// Collect data points given by the provider.
	var size int
	for _, s := range *resp.Series {
		size += int(*s.Length)
	}
	out := make([]metrics.DataPoint, 0, size)
	for _, s := range *resp.Series {
		points := s.Pointlist
		if points == nil || len(*points) == 0 {
			return nil, fmt.Errorf("invalid response: no data points found within the queried range: %w", metrics.ErrNoDataFound)
		}
		for _, point := range *points {
			if len(point) < 2 {
				return nil, fmt.Errorf("invalid response: invalid data point found")
			}
			// NOTE: A data point is assumed to be kind of like [unix-time, value].
			out = append(out, metrics.DataPoint{
				Timestamp: int64(point[0]),
				Value:     point[1],
			})
		}
	}
	return out, nil
}
