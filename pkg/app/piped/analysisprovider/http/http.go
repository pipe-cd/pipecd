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

// Package http provides a way to analyze with http requests.
// This allows you to do smoke tests, load tests and so on, at your leisure.
package http

import (
	"context"
	"fmt"
	"net/http"
	"time"
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

// Run sends the given HTTP request and then evaluate whether the response is expected one.
func (p *Provider) Run(ctx context.Context, req *http.Request, expectedCode int, expectedResponse string) (bool, error) {
	req = req.WithContext(ctx)
	res, err := p.client.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode != expectedCode {
		return false, fmt.Errorf("unexpected status code %d", expectedCode)
	}
	// TODO: Decide how to check if the body is expected one.
	return true, nil
}
