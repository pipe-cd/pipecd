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

	"github.com/kapetaniosci/pipe/pkg/app/piped/executor"
	"github.com/kapetaniosci/pipe/pkg/config"
)

const (
	ProviderType   = "Prometheus"
	defaultTimeout = 30 * time.Second
)

// Provider is a client for prometheus.
type Provider struct {
	api      v1.API
	username string
	password string

	timeout      time.Duration
	logPersister executor.LogPersister
}

func NewProvider(address, username, password string) (*Provider, error) {
	// TODO: Decide the way to authenticate.
	client, err := api.NewClient(api.Config{
		Address: address,
	})
	if err != nil {
		return nil, err
	}

	return &Provider{
		api:      v1.NewAPI(client),
		username: username,
		password: password,
		timeout:  defaultTimeout,
	}, nil
}

func (p *Provider) Type() string {
	return ProviderType
}

func (p *Provider) RunQuery(ctx context.Context, query string, expected config.AnalysisExpected) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	response, warnings, err := p.api.Query(ctx, query, time.Now())
	if err != nil {
		return false, err
	}
	for _, w := range warnings {
		p.logPersister.AppendInfo(w)
	}
	// TODO: Address the case of comparing with baseline
	return p.evaluate(expected, response)
}

func (p *Provider) evaluate(expected config.AnalysisExpected, response model.Value) (bool, error) {
	switch value := response.(type) {
	case *model.Scalar:
		result := float64(value.Value)
		if math.IsNaN(result) {
			return false, fmt.Errorf("the result %v is not a number", result)
		}
		return inRange(expected, result)
	case model.Vector:
		lv := len(value)
		if lv == 0 {
			return false, fmt.Errorf("zero value returned")
		}
		results := make([]float64, 0, lv)
		for _, s := range value {
			if s == nil {
				continue
			}
			result := float64(s.Value)
			if math.IsNaN(result) {
				return false, fmt.Errorf("the result %v is not a number", result)
			}
			results = append(results, result)
		}
		// TODO: Consider the case of multiple results.
		return inRange(expected, results[0])
	default:
		return false, fmt.Errorf("unsupported prometheus metric type")
	}
}

// TODO: Move to common package.
func inRange(expected config.AnalysisExpected, value float64) (bool, error) {
	if expected.Min == nil && expected.Max == nil {
		return false, fmt.Errorf("expected range is undefined")
	}
	if min := expected.Min; min != nil && *min > value {
		return false, nil
	}
	if max := expected.Min; max != nil && *max > value {
		return false, nil
	}
	return true, nil
}
