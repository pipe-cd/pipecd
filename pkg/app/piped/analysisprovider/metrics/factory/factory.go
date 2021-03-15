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

package factory

import (
	"fmt"
	"io/ioutil"
	"strings"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/piped/analysisprovider/metrics"
	"github.com/pipe-cd/pipe/pkg/app/piped/analysisprovider/metrics/datadog"
	"github.com/pipe-cd/pipe/pkg/app/piped/analysisprovider/metrics/prometheus"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

// NewProvider generates an appropriate provider according to analysis provider config.
func NewProvider(analysisTempCfg *config.TemplatableAnalysisMetrics, providerCfg *config.PipedAnalysisProvider, logger *zap.Logger) (metrics.Provider, error) {
	switch providerCfg.Type {
	case model.AnalysisProviderPrometheus:
		options := []prometheus.Option{
			prometheus.WithLogger(logger),
		}
		if t := analysisTempCfg.Timeout.Duration(); t > 0 {
			options = append(options, prometheus.WithTimeout(t))
		}
		return prometheus.NewProvider(providerCfg.PrometheusConfig.Address, options...)
	case model.AnalysisProviderDatadog:
		var apiKey, applicationKey string
		cfg := providerCfg.DatadogConfig
		if cfg.APIKeyFile != "" {
			a, err := ioutil.ReadFile(cfg.APIKeyFile)
			if err != nil {
				return nil, fmt.Errorf("failed to read the api-key file: %w", err)
			}
			apiKey = strings.TrimSpace(string(a))
		}
		if cfg.ApplicationKeyFile != "" {
			a, err := ioutil.ReadFile(cfg.ApplicationKeyFile)
			if err != nil {
				return nil, fmt.Errorf("failed to read the application-key file: %w", err)
			}
			applicationKey = strings.TrimSpace(string(a))
		}
		options := []datadog.Option{
			datadog.WithLogger(logger),
		}
		if cfg.Address != "" {
			options = append(options, datadog.WithAddress(cfg.Address))
		}
		if t := analysisTempCfg.Timeout.Duration(); t > 0 {
			options = append(options, datadog.WithTimeout(t))
		}
		return datadog.NewProvider(apiKey, applicationKey, options...)
	default:
		return nil, fmt.Errorf("any of providers config not found")
	}
}
