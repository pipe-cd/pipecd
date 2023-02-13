// Copyright 2023 The PipeCD Authors.
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
	"net/http"
	"time"

	"github.com/pipe-cd/pipecd/pkg/config"
)

const (
	ProviderType   = "HTTP"
	defaultTimeout = 30 * time.Second
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
	// TODO: Decide how to check if the body is expected one.
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
