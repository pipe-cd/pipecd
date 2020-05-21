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

// response represents a response from prometheus server.
type response struct {
	Data struct {
		Result []struct {
			Metric struct {
				Name string `json:"name"`
			}
			Value []interface{} `json:"value"`
		}
	}
}

func NewProvider(address, username, password string) (*Provider, error) {
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

func (p *Provider) RunQuery(ctx context.Context, query, expected string) (bool, error) {
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

func (p *Provider) evaluate(expected string, response model.Value) (bool, error) {
	switch value := response.(type) {
	case *model.Scalar:
		result := float64(value.Value)
		if math.IsNaN(result) {
			return false, fmt.Errorf("the result %v is not a number", result)
		}
		// FIXME: evaluate
		return false, nil
	case model.Vector:
		results := make([]float64, 0, len(value))
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
		// FIXME: evaluate
		return false, nil
	default:
		return false, fmt.Errorf("unsupported prometheus metric type")
	}
}
