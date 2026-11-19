// Copyright 2025 The PipeCD Authors.
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

// Package http provides a way to analyze with http requests.
// This allows you to do smoke tests, load tests and so on, at your leisure.
package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/analysis/config"
)

const (
	ProviderType        = "HTTP"
	defaultTimeout      = 30 * time.Second
	maxResponseBodySize = 1 << 20 // 1 MiB
)

type Provider struct {
	client *http.Client
}

func (p *Provider) Type() string {
	return ProviderType
}

func NewProvider(timeout time.Duration) *Provider {
	if timeout == 0 {
		timeout = defaultTimeout
	}
	return &Provider{
		client: &http.Client{Timeout: timeout},
	}
}

// Run sends an HTTP request and then evaluate whether the response is expected one.
func (p *Provider) Run(ctx context.Context, cfg *config.AnalysisHTTP) (bool, string, error) {
	if len(cfg.ExpectedResponse) > maxResponseBodySize {
		return false, "", fmt.Errorf("expected response exceeds maximum size of %d bytes", maxResponseBodySize)
	}

	req, err := p.makeRequest(ctx, cfg)
	if err != nil {
		return false, "", err
	}

	res, err := p.client.Do(req)
	if err != nil {
		return false, "", err
	}
	defer res.Body.Close()

	if res.StatusCode != cfg.ExpectedCode {
		return false, "", fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

	if cfg.ExpectedResponse == "" {
		return true, "", nil
	}

	body, err := io.ReadAll(io.LimitReader(res.Body, maxResponseBodySize+1))
	if err != nil {
		return false, "", fmt.Errorf("failed to read response body: %w", err)
	}
	if len(body) > maxResponseBodySize {
		return false, "", fmt.Errorf("response body exceeds maximum size of %d bytes", maxResponseBodySize)
	}
	if string(body) != cfg.ExpectedResponse {
		return false, "", fmt.Errorf("unexpected response body")
	}

	return true, "", nil
}

func (p *Provider) makeRequest(ctx context.Context, cfg *config.AnalysisHTTP) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, cfg.Method, cfg.URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header = make(http.Header, len(cfg.Headers))
	for _, h := range cfg.Headers {
		req.Header.Set(h.Key, h.Value)
	}
	return req, nil
}
