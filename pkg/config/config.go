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

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"sigs.k8s.io/yaml"
)

const version = "v1"

// Kind represents the kind of configuration the data contains.
type Kind string

const (
	// KindK8sApp represents configuration for a Kubernetes application.
	// This application can be a group of plain-YAML Kubernetes manifests.
	KindK8sApp Kind = "K8sApp"
	// KindK8sKustomizationApp represents configuration for a Kubernetes application
	// that using kustomization for templating.
	KindK8sKustomizationApp Kind = "K8sKustomizationApp"
	// KindK8sHelmApp represents configuration for a Kubernetes application.
	// that using helm for templating.
	KindK8sHelmApp Kind = "K8sHelmApp"
	// KindTerraformApp represents configuration for a Terraform application.
	// This application contains a single workspace of a terraform root module.
	KindTerraformApp Kind = "TerraformApp"
	// KindLambdaApp represents configuration for an AWS Lambda application.
	KindLambdaApp Kind = "LambdaApp"
	// KindNotification represents shared notification configuration for a repository.
	// This configuration file should be placed in .pipe directory
	// at the root of the repository.
	KindNotification Kind = "Notification"
	// KindAnalysisTemplate represents shared analysis template for a repository.
	// This configuration file should be placed in .pipe directory
	// at the root of the repository.
	KindAnalysisTemplate Kind = "AnalysisTemplate"
	// KindRunner represents configuration for runner.
	// This configuration will be loaded while the runner is starting up.
	KindRunner Kind = "Runner"
	// KindControlPlane represents configuration for control plane's services.
	KindControlPlane Kind = "ControlPlane"
)

// Config represents configuration data load from file.
// The spec is depend on the kind of configuration.
type Config struct {
	Kind    Kind
	Version string

	// Application Specs.
	K8sAppSpec              *K8sAppSpec
	K8sKustomizationAppSpec *K8sKustomizationAppSpec
	K8sHelmAppSpec          *K8sHelmAppSpec
	TerraformAppSpec        *TerraformAppSpec

	// Repositori Shared Specs.
	NotificationSpec     *NotificationSpec
	AnalysisTemplateSpec *AnalysisTemplateSpec

	// Runner and Server Specs.
	RunnerSpec       *RunnerSpec
	ControlPlaneSpec *ControlPlaneSpec
}

type genericConfig struct {
	Kind    Kind            `json:"kind"`
	Version string          `json:"version,omitempty"`
	Spec    json.RawMessage `json:"spec"`
}

// Spec returns the actual spec inside for this kind of configuration.
func (c *Config) Spec() interface{} {
	switch c.Kind {
	case KindK8sApp:
		return c.K8sAppSpec
	case KindK8sKustomizationApp:
		return c.K8sKustomizationAppSpec
	case KindK8sHelmApp:
		return c.K8sHelmAppSpec
	case KindTerraformApp:
		return c.TerraformAppSpec
	case KindNotification:
		return c.NotificationSpec
	case KindAnalysisTemplate:
		return c.AnalysisTemplateSpec
	case KindRunner:
		return c.RunnerSpec
	case KindControlPlane:
		return c.ControlPlaneSpec
	}
	return nil
}

// UnmarshalJSON customizes the way to unmarshal json data into Config struct.
// Firstly, this unmarshal to a generic config and then unmarshal the spec
// which depend on the kind of configuration.
func (c *Config) UnmarshalJSON(data []byte) error {
	var err error
	gc := genericConfig{}
	if err = json.Unmarshal(data, &gc); err != nil {
		return err
	}
	c.Kind = gc.Kind
	c.Version = gc.Version

	switch gc.Kind {
	case KindK8sApp:
		c.K8sAppSpec = &K8sAppSpec{}
		if len(gc.Spec) > 0 {
			err = json.Unmarshal(gc.Spec, c.K8sAppSpec)
		}
	case KindK8sKustomizationApp:
		c.K8sKustomizationAppSpec = &K8sKustomizationAppSpec{}
		if len(gc.Spec) > 0 {
			err = json.Unmarshal(gc.Spec, c.K8sKustomizationAppSpec)
		}
	case KindK8sHelmApp:
		c.K8sHelmAppSpec = &K8sHelmAppSpec{}
		if len(gc.Spec) > 0 {
			err = json.Unmarshal(gc.Spec, c.K8sHelmAppSpec)
		}
	case KindTerraformApp:
		c.TerraformAppSpec = &TerraformAppSpec{}
		if len(gc.Spec) > 0 {
			err = json.Unmarshal(gc.Spec, c.TerraformAppSpec)
		}
	case KindNotification:
		c.NotificationSpec = &NotificationSpec{}
		if len(gc.Spec) > 0 {
			err = json.Unmarshal(gc.Spec, c.NotificationSpec)
		}
	case KindAnalysisTemplate:
		c.AnalysisTemplateSpec = &AnalysisTemplateSpec{}
		if len(gc.Spec) > 0 {
			err = json.Unmarshal(gc.Spec, c.AnalysisTemplateSpec)
		}
	case KindRunner:
		c.RunnerSpec = &RunnerSpec{}
		if len(gc.Spec) > 0 {
			err = json.Unmarshal(gc.Spec, c.RunnerSpec)
		}
	case KindControlPlane:
		c.ControlPlaneSpec = &ControlPlaneSpec{}
		if len(gc.Spec) > 0 {
			err = json.Unmarshal(gc.Spec, c.ControlPlaneSpec)
		}
	default:
		err = fmt.Errorf("unsupported kind: %s", c.Kind)
	}
	return err
}

type validator interface {
	Validate() error
}

// Validate validates the value of all fields.
func (c *Config) Validate() error {
	if c.Version != "v1" && c.Version != "" {
		return fmt.Errorf("unsupported version: %s", c.Version)
	}
	if spec := c.Spec(); spec != nil {
		if v, ok := spec.(validator); ok {
			if err := v.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
}

// LoadFromYAML reads and decodes a yaml file to construct the Config.
func LoadFromYAML(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return DecodeYAML(data)
}

// DecodeYAML unmarshals config YAML data to config struct.
// It also validates the configuration after decoding.
func DecodeYAML(data []byte) (*Config, error) {
	js, err := yaml.YAMLToJSON(data)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	if err := json.Unmarshal(js, c); err != nil {
		return nil, err
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return c, nil
}
