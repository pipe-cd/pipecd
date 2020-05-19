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
	"net/url"
	"time"
)

const ProviderType = "Prometheus"

// Provider is a client for prometheus.
type Provider struct {
	timeout  time.Duration
	address  *url.URL
	username string
	password string
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
	u, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	return &Provider{
		address:  u,
		username: username,
		password: password,
	}, nil
}

func (p *Provider) Type() string {
	return ProviderType
}

func (p *Provider) RunQuery(query, expected string) (bool, error) {
	return false, nil
}
