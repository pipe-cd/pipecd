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

import (
	"encoding/json"
	"fmt"
)

const (
	maskString = "******"
)

type PluginConfig struct {
	// List of analysis providers can be used by this piped.
	AnalysisProviders []PipedAnalysisProvider `json:"analysisProviders,omitempty"`
}

type PipedAnalysisProvider struct {
	Name string               `json:"name"`
	Type AnalysisProviderType `json:"type"`

	PrometheusConfig  *AnalysisProviderPrometheusConfig
	DatadogConfig     *AnalysisProviderDatadogConfig
	StackdriverConfig *AnalysisProviderStackdriverConfig
}

func (p *PipedAnalysisProvider) Mask() {
	if p.PrometheusConfig != nil {
		p.PrometheusConfig.Mask()
	}
	if p.DatadogConfig != nil {
		p.DatadogConfig.Mask()
	}
	if p.StackdriverConfig != nil {
		p.StackdriverConfig.Mask()
	}
}

type genericPipedAnalysisProvider struct {
	Name   string               `json:"name"`
	Type   AnalysisProviderType `json:"type"`
	Config json.RawMessage      `json:"config"`
}

func (p *PipedAnalysisProvider) MarshalJSON() ([]byte, error) {
	var (
		err    error
		config json.RawMessage
	)

	switch p.Type {
	case AnalysisProviderDatadog:
		config, err = json.Marshal(p.DatadogConfig)
	case AnalysisProviderPrometheus:
		config, err = json.Marshal(p.PrometheusConfig)
	case AnalysisProviderStackdriver:
		config, err = json.Marshal(p.StackdriverConfig)
	default:
		err = fmt.Errorf("unsupported analysis provider type: %s", p.Name)
	}

	if err != nil {
		return nil, err
	}

	return json.Marshal(&genericPipedAnalysisProvider{
		Name:   p.Name,
		Type:   p.Type,
		Config: config,
	})
}

func (p *PipedAnalysisProvider) UnmarshalJSON(data []byte) error {
	var err error
	gp := genericPipedAnalysisProvider{}
	if err = json.Unmarshal(data, &gp); err != nil {
		return err
	}
	p.Name = gp.Name
	p.Type = gp.Type

	switch p.Type {
	case AnalysisProviderPrometheus:
		p.PrometheusConfig = &AnalysisProviderPrometheusConfig{}
		if len(gp.Config) > 0 {
			err = json.Unmarshal(gp.Config, p.PrometheusConfig)
		}
	case AnalysisProviderDatadog:
		p.DatadogConfig = &AnalysisProviderDatadogConfig{}
		if len(gp.Config) > 0 {
			err = json.Unmarshal(gp.Config, p.DatadogConfig)
		}
	case AnalysisProviderStackdriver:
		p.StackdriverConfig = &AnalysisProviderStackdriverConfig{}
		if len(gp.Config) > 0 {
			err = json.Unmarshal(gp.Config, p.StackdriverConfig)
		}
	default:
		err = fmt.Errorf("unsupported analysis provider type: %s", p.Name)
	}
	return err
}

func (p *PipedAnalysisProvider) Validate() error {
	switch p.Type {
	case AnalysisProviderPrometheus:
		return p.PrometheusConfig.Validate()
	case AnalysisProviderDatadog:
		return p.DatadogConfig.Validate()
	case AnalysisProviderStackdriver:
		return p.StackdriverConfig.Validate()
	default:
		return fmt.Errorf("unknow provider type: %s", p.Type)
	}
}

type AnalysisProviderPrometheusConfig struct {
	Address string `json:"address"`
	// The path to the username file.
	UsernameFile string `json:"usernameFile,omitempty"`
	// The path to the password file.
	PasswordFile string `json:"passwordFile,omitempty"`
}

func (a *AnalysisProviderPrometheusConfig) Validate() error {
	if a.Address == "" {
		return fmt.Errorf("prometheus analysis provider requires the address")
	}
	return nil
}

func (a *AnalysisProviderPrometheusConfig) Mask() {
	if len(a.PasswordFile) != 0 {
		a.PasswordFile = maskString
	}
}

type AnalysisProviderDatadogConfig struct {
	// The address of Datadog API server.
	// Only "datadoghq.com", "us3.datadoghq.com", "datadoghq.eu", "ddog-gov.com" are available.
	// Defaults to "datadoghq.com"
	Address string `json:"address,omitempty"`
	// Required: The path to the api key file.
	APIKeyFile string `json:"apiKeyFile"`
	// Required: The path to the application key file.
	ApplicationKeyFile string `json:"applicationKeyFile"`
	// Base64 API Key for Datadog API server.
	APIKeyData string `json:"apiKeyData,omitempty"`
	// Base64 Application Key for Datadog API server.
	ApplicationKeyData string `json:"applicationKeyData,omitempty"`
}

func (a *AnalysisProviderDatadogConfig) Validate() error {
	if a.APIKeyFile == "" && a.APIKeyData == "" {
		return fmt.Errorf("either datadog APIKeyFile or APIKeyData must be set")
	}
	if a.ApplicationKeyFile == "" && a.ApplicationKeyData == "" {
		return fmt.Errorf("either datadog ApplicationKeyFile or ApplicationKeyData must be set")
	}
	if a.APIKeyData != "" && a.APIKeyFile != "" {
		return fmt.Errorf("only datadog APIKeyFile or APIKeyData can be set")
	}
	if a.ApplicationKeyData != "" && a.ApplicationKeyFile != "" {
		return fmt.Errorf("only datadog ApplicationKeyFile or ApplicationKeyData can be set")
	}
	return nil
}

func (a *AnalysisProviderDatadogConfig) Mask() {
	if len(a.APIKeyFile) != 0 {
		a.APIKeyFile = maskString
	}
	if len(a.ApplicationKeyFile) != 0 {
		a.ApplicationKeyFile = maskString
	}
	if len(a.APIKeyData) != 0 {
		a.APIKeyData = maskString
	}
	if len(a.ApplicationKeyData) != 0 {
		a.ApplicationKeyData = maskString
	}
}

// func(a *AnalysisProviderDatadogConfig)

type AnalysisProviderStackdriverConfig struct {
	// The path to the service account file.
	ServiceAccountFile string `json:"serviceAccountFile"`
}

func (a *AnalysisProviderStackdriverConfig) Mask() {
	if len(a.ServiceAccountFile) != 0 {
		a.ServiceAccountFile = maskString
	}
}

func (a *AnalysisProviderStackdriverConfig) Validate() error {
	return nil
}
