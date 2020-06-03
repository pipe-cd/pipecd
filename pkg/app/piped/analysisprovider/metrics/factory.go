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

package metrics

import (
	"fmt"
	"io/ioutil"

	"go.uber.org/zap"

	"github.com/kapetaniosci/pipe/pkg/app/piped/analysisprovider/metrics/datadog"
	"github.com/kapetaniosci/pipe/pkg/app/piped/analysisprovider/metrics/prometheus"
	"github.com/kapetaniosci/pipe/pkg/config"
	"github.com/kapetaniosci/pipe/pkg/model"
)

type Factory struct {
	logger *zap.Logger
}

func NewFactory(logger *zap.Logger) *Factory {
	return &Factory{logger: logger}
}

// NewProvider generates an appropriate provider according to analysis provider config.
func (f *Factory) NewProvider(providerCfg *config.PipedAnalysisProvider) (provider Provider, err error) {
	switch providerCfg.Type {
	case model.AnalysisProviderPrometheus:
		provider, err = prometheus.NewProvider(providerCfg.PrometheusConfig.Address, f.logger)
		if err != nil {
			return
		}
	case model.AnalysisProviderDatadog:
		var apiKey, applicationKey string
		cfg := providerCfg.DatadogConfig
		if cfg.APIKeyFile != "" {
			a, err := ioutil.ReadFile(cfg.APIKeyFile)
			if err != nil {
				return nil, err
			}
			apiKey = string(a)
		}
		if cfg.ApplicationKeyFile != "" {
			a, err := ioutil.ReadFile(cfg.ApplicationKeyFile)
			if err != nil {
				return nil, err
			}
			applicationKey = string(a)
		}
		provider, err = datadog.NewProvider(cfg.Address, apiKey, applicationKey)
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("any of providers config not found")
		return
	}
	return
}
