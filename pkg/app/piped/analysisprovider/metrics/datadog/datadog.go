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

package datadog

import (
	"context"
	"time"

	"github.com/kapetaniosci/pipe/pkg/config"
)

const ProviderType = "Datadog"

// Provider is a client for datadog.
type Provider struct {
	metricsQueryEndpoint     string
	apiKeyValidationEndpoint string

	timeout        time.Duration
	apiKey         string
	applicationKey string
	fromDelta      int64
}

func NewProvider(address, apiKey, applicationKey string) (*Provider, error) {
	return &Provider{
		metricsQueryEndpoint: address,
		apiKey:               apiKey,
		applicationKey:       applicationKey,
	}, nil
}

// response represents a response from datadog server.
type response struct {
	Series []struct {
		Pointlist [][]float64 `json:"pointlist"`
	}
}

func (p *Provider) Type() string {
	return ProviderType
}

func (p *Provider) RunQuery(ctx context.Context, query string, expected config.AnalysisExpected) (bool, error) {
	return false, nil
}
