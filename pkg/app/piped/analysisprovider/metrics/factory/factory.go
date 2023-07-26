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

package factory

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/piped/analysisprovider/metrics"
	"github.com/pipe-cd/pipecd/pkg/app/piped/analysisprovider/metrics/datadog"
	"github.com/pipe-cd/pipecd/pkg/app/piped/analysisprovider/metrics/prometheus"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

// NewProvider generates an appropriate provider according to analysis provider config.
func NewProvider(analysisTempCfg *config.TemplatableAnalysisMetrics, providerCfg *config.PipedAnalysisProvider, logger *zap.Logger) (metrics.Provider, error) {
	switch providerCfg.Type {
	case model.AnalysisProviderPrometheus:
		options := []prometheus.Option{
			prometheus.WithLogger(logger),
			prometheus.WithTimeout(analysisTempCfg.Timeout.Duration()),
		}
		cfg := providerCfg.PrometheusConfig
		if cfg.UsernameFile != "" && cfg.PasswordFile != "" {
			username, err := os.ReadFile(cfg.UsernameFile)
			if err != nil {
				return nil, fmt.Errorf("failed to read the username file: %w", err)
			}
			password, err := os.ReadFile(cfg.PasswordFile)
			if err != nil {
				return nil, fmt.Errorf("failed to read the password file: %w", err)
			}
			options = append(options, prometheus.WithBasicAuth(strings.TrimSpace(string(username)), strings.TrimSpace(string(password))))
		}
		return prometheus.NewProvider(providerCfg.PrometheusConfig.Address, options...)
	case model.AnalysisProviderDatadog:
		var apiKey, applicationKey string
		cfg := providerCfg.DatadogConfig
		if cfg.APIKeyFile != "" {
			a, err := os.ReadFile(cfg.APIKeyFile)
			if err != nil {
				return nil, fmt.Errorf("failed to read the api-key file: %w", err)
			}
			apiKey = strings.TrimSpace(string(a))
		}
		if cfg.ApplicationKeyFile != "" {
			a, err := os.ReadFile(cfg.ApplicationKeyFile)
			if err != nil {
				return nil, fmt.Errorf("failed to read the application-key file: %w", err)
			}
			applicationKey = strings.TrimSpace(string(a))
		}
		if cfg.APIKeyData != "" {
			a, err := base64.StdEncoding.DecodeString(cfg.APIKeyData)
			if err != nil {
				return nil, fmt.Errorf("failed to decode the api-key data: %w", err)
			}
			apiKey = string(a)
		}
		if cfg.ApplicationKeyData != "" {
			a, err := base64.StdEncoding.DecodeString(cfg.ApplicationKeyData)
			if err != nil {
				return nil, fmt.Errorf("failed to decode the application-key data: %w", err)
			}
			applicationKey = string(a)
		}
		options := []datadog.Option{
			datadog.WithLogger(logger),
			datadog.WithTimeout(analysisTempCfg.Timeout.Duration()),
		}
		if cfg.Address != "" {
			options = append(options, datadog.WithAddress(cfg.Address))
		}
		return datadog.NewProvider(apiKey, applicationKey, options...)
	default:
		return nil, fmt.Errorf("any of providers config not found")
	}
}
