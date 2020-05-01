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

const versionV1Beta1 = "pipecd.dev/v1beta1"

// Kind represents the kind of configuration the data contains.
type Kind string

const (
	// KindKubernetesApp represents deployment configuration for a Kubernetes application.
	// This application can be a group of plain-YAML Kubernetes manifests,
	// or kustomization manifests or helm manifests.
	KindKubernetesApp Kind = "KubernetesApp"
	// KindTerraformApp represents deployment configuration for a Terraform application.
	// This application contains a single workspace of a terraform root module.
	KindTerraformApp Kind = "TerraformApp"
	// KindCrossplaneApp represents deployment configuration for a Crossplane applicaiton.
	KindCrossplaneApp Kind = "CrossplaneApp"
	// KindLambdaApp represents deployment configuration for an AWS Lambda application.
	KindLambdaApp Kind = "LambdaApp"
	// KindCloudRunApp represents deployment configuration for a CloudRun application.
	KindCloudRunApp Kind = "CloudRunApp"
)

const (
	// KindRunner represents configuration for runner.
	// This configuration will be loaded while the runner is starting up.
	KindRunner Kind = "Runner"
	// KindControlPlane represents configuration for control plane's services.
	KindControlPlane Kind = "ControlPlane"
	// KindNotification represents shared notification configuration for a repository.
	// This configuration file should be placed in .pipe directory
	// at the root of the repository.
	KindNotification Kind = "Notification"
	// KindAnalysisTemplate represents shared analysis template for a repository.
	// This configuration file should be placed in .pipe directory
	// at the root of the repository.
	KindAnalysisTemplate Kind = "AnalysisTemplate"
)

// Config represents configuration data load from file.
// The spec is depend on the kind of configuration.
type Config struct {
	Kind       Kind
	APIVersion string
	spec       interface{}

	// Application specs.
	KubernetesAppSpec *KubernetesAppSpec
	TerraformAppSpec  *TerraformAppSpec

	RunnerSpec           *RunnerSpec
	ControlPlaneSpec     *ControlPlaneSpec
	NotificationSpec     *NotificationSpec
	AnalysisTemplateSpec *AnalysisTemplateSpec
}

type genericConfig struct {
	Kind       Kind            `json:"kind"`
	APIVersion string          `json:"apiVersion,omitempty"`
	Spec       json.RawMessage `json:"spec"`
}

func (c *Config) init(kind Kind, apiVersion string) error {
	c.Kind = kind
	c.APIVersion = apiVersion

	switch kind {
	case KindKubernetesApp:
		c.KubernetesAppSpec = &KubernetesAppSpec{}
		c.spec = c.KubernetesAppSpec
	case KindTerraformApp:
		c.TerraformAppSpec = &TerraformAppSpec{}
		c.spec = c.TerraformAppSpec
	case KindRunner:
		c.RunnerSpec = &RunnerSpec{}
		c.spec = c.RunnerSpec
	case KindControlPlane:
		c.ControlPlaneSpec = &ControlPlaneSpec{}
		c.spec = c.ControlPlaneSpec
	case KindNotification:
		c.NotificationSpec = &NotificationSpec{}
		c.spec = c.NotificationSpec
	case KindAnalysisTemplate:
		c.AnalysisTemplateSpec = &AnalysisTemplateSpec{}
		c.spec = c.AnalysisTemplateSpec
	default:
		return fmt.Errorf("unsupported kind: %s", c.Kind)
	}
	return nil
}

// UnmarshalJSON customizes the way to unmarshal json data into Config struct.
// Firstly, this unmarshal to a generic config and then unmarshal the spec
// which depend on the kind of configuration.
func (c *Config) UnmarshalJSON(data []byte) error {
	var (
		err error
		gc  = genericConfig{}
	)
	if err = json.Unmarshal(data, &gc); err != nil {
		return err
	}
	if err = c.init(gc.Kind, gc.APIVersion); err != nil {
		return err
	}
	if len(gc.Spec) > 0 {
		err = json.Unmarshal(gc.Spec, c.spec)
	}
	return err
}

type validator interface {
	Validate() error
}

// Validate validates the value of all fields.
func (c *Config) Validate() error {
	if c.APIVersion != versionV1Beta1 {
		return fmt.Errorf("unsupported version: %s", c.APIVersion)
	}
	if spec, ok := c.spec.(validator); ok && spec != nil {
		if err := spec.Validate(); err != nil {
			return err
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
