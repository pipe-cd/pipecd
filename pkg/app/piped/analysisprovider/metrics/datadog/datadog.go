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

	"github.com/pipe-cd/pipe/pkg/config"
)

const (
	ProviderType   = "Datadog"
	defaultAddress = "datadoghq.com"
	defaultTimeout = 30 * time.Second
)

// Provider works as an HTTP client for datadog.
type Provider struct {
	client *datadog.APIClient

	address           string
	apiKey            string
	applicationKey    string
	queriedTimePeriod int64
	timeout           time.Duration
	logger            *zap.Logger
}

func NewProvider(apiKey, applicationKey string, queriedTimePeriod time.Duration, opts ...Option) (*Provider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api-key is required")
	}
	if applicationKey == "" {
		return nil, fmt.Errorf("application-key is required")
	}
	if queriedTimePeriod == 0 {
		return nil, fmt.Errorf("aggregation period is required")
	}

	p := &Provider{
		client:            datadog.NewAPIClient(datadog.NewConfiguration()),
		address:           defaultAddress,
		apiKey:            apiKey,
		applicationKey:    applicationKey,
		queriedTimePeriod: int64(queriedTimePeriod.Seconds()),
		timeout:           defaultTimeout,
		logger:            zap.NewNop(),
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

// RunQuery issues an HTTP request to the API named "MetricsApi.QueryMetrics", then evaluate its response.
// It performs the given query for the period specified by its own queried time period using the current
// time as the end of the queried time period.
//
// See more: https://docs.datadoghq.com/api/latest/metrics/#query-timeseries-points
func (p *Provider) RunQuery(ctx context.Context, query string, expected config.AnalysisExpected) (bool, error) {
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

	from := time.Now().Unix() - p.queriedTimePeriod
	to := time.Now().Unix()

	resp, httpResp, err := p.client.MetricsApi.QueryMetrics(ctx).From(from).To(to).Query(query).Execute()
	if err != nil {
		return false, fmt.Errorf("failed to call \"MetricsApi.QueryMetrics\": %w", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected HTTP status code from %s: %d", httpResp.Request.URL, httpResp.StatusCode)
	}
	if resp.Series == nil || len(*resp.Series) == 0 {
		return false, fmt.Errorf("no timeseries queried found")
	}
	points := (*resp.Series)[0].Pointlist
	if points == nil || len(*points) == 0 {
		return false, fmt.Errorf("no data points of the time series found")
	}
	// TODO: Think about how to handle multiple data points
	point := (*points)[len(*points)-1]
	if len(point) < 2 {
		return false, fmt.Errorf("invalid data point found")
	}
	// A data point is assumed to be kind of like [unix-time, value].
	return p.evaluate(expected, point[1])
}

func (p *Provider) evaluate(expected config.AnalysisExpected, response float64) (bool, error) {
	// TODO: Implement evaluation of response from Datadog
	return false, nil
}
