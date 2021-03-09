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

	"github.com/pipe-cd/pipe/pkg/app/piped/analysisprovider/metrics/datadog"
	"github.com/pipe-cd/pipe/pkg/app/piped/analysisprovider/metrics/prometheus"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

type Factory struct {
	logger *zap.Logger
}

func NewFactory(logger *zap.Logger) *Factory {
	return &Factory{logger: logger}
}

// NewProvider generates an appropriate provider according to analysis provider config.
func (f *Factory) NewProvider(analysisTempCfg *config.TemplatableAnalysisMetrics, providerCfg *config.PipedAnalysisProvider) (Provider, error) {
	switch providerCfg.Type {
	case model.AnalysisProviderPrometheus:
		return prometheus.NewProvider(providerCfg.PrometheusConfig.Address, analysisTempCfg.Timeout.Duration(), f.logger)
	case model.AnalysisProviderDatadog:
		var apiKey, applicationKey string
		cfg := providerCfg.DatadogConfig
		if cfg.APIKeyFile != "" {
			a, err := ioutil.ReadFile(cfg.APIKeyFile)
			if err != nil {
				return nil, fmt.Errorf("failed to read the api-key file: %w", err)
			}
			apiKey = string(a)
		}
		if cfg.ApplicationKeyFile != "" {
			a, err := ioutil.ReadFile(cfg.ApplicationKeyFile)
			if err != nil {
				return nil, fmt.Errorf("failed to read the application-key file: %w", err)
			}
			applicationKey = string(a)
		}
		options := []datadog.Option{
			datadog.WithAddress(cfg.Address),
			datadog.WithTimeout(analysisTempCfg.Timeout.Duration()),
			datadog.WithLogger(f.logger),
		}
		return datadog.NewProvider(apiKey, applicationKey, analysisTempCfg.Interval.Duration(), options...)
	default:
		return nil, fmt.Errorf("any of providers config not found")
	}
}
