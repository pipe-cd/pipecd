// Copyright 2025 The PipeCD Authors.
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

package config

import "time"

// Duration wraps time.Duration and provides a Duration() helper
// compatible with the upstream config types.
type Duration time.Duration

func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}

// AnalysisProviderType represents supported provider kinds.
type AnalysisProviderType string

const (
	AnalysisProviderPrometheus  AnalysisProviderType = "PROMETHEUS"
	AnalysisProviderDatadog     AnalysisProviderType = "DATADOG"
	AnalysisProviderStackdriver AnalysisProviderType = "STACKDRIVER"
)

// TemplatableAnalysisMetrics is a minimal subset used by the plugin factory.
type TemplatableAnalysisMetrics struct {
	Timeout Duration `json:"timeout"`
}

// PipedAnalysisProvider is a minimal subset of Piped config used by factories.
type PipedAnalysisProvider struct {
	Type              AnalysisProviderType               `json:"type"`
	PrometheusConfig  *AnalysisProviderPrometheusConfig  `json:"prometheus,omitempty"`
	DatadogConfig     *AnalysisProviderDatadogConfig     `json:"datadog,omitempty"`
	StackdriverConfig *AnalysisProviderStackdriverConfig `json:"stackdriver,omitempty"`
}

type AnalysisProviderPrometheusConfig struct {
	Address      string `json:"address"`
	UsernameFile string `json:"usernameFile,omitempty"`
	PasswordFile string `json:"passwordFile,omitempty"`
}

type AnalysisProviderDatadogConfig struct {
	Address            string `json:"address,omitempty"`
	APIKeyFile         string `json:"apiKeyFile,omitempty"`
	ApplicationKeyFile string `json:"applicationKeyFile,omitempty"`
	APIKeyData         string `json:"apiKeyData,omitempty"`
	ApplicationKeyData string `json:"applicationKeyData,omitempty"`
}

type AnalysisProviderStackdriverConfig struct {
	ServiceAccountFile string `json:"serviceAccountFile"`
}

// AnalysisHTTP holds settings for HTTP analysis.
type AnalysisHTTP struct {
	URL              string               `json:"url"`
	Method           string               `json:"method"`
	Headers          []AnalysisHTTPHeader `json:"headers"`
	ExpectedCode     int                  `json:"expectedCode"`
	ExpectedResponse string               `json:"expectedResponse,omitempty"`
}

type AnalysisHTTPHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
