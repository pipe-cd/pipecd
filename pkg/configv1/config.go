// Copyright 2024 The PipeCD Authors.
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
	"errors"
	"fmt"
	"os"

	"github.com/creasty/defaults"
	"sigs.k8s.io/yaml"
)

const (
	SharedConfigurationDirName = ".pipe"
	VersionV1Beta1             = "pipecd.dev/v1beta1"
)

// Kind represents the kind of configuration the data contains.
type Kind string

const (
	// KindKubernetesApp represents application configuration for a Kubernetes application.
	// This application can be a group of plain-YAML Kubernetes manifests,
	// or kustomization manifests or helm manifests.
	//
	// Deprecated: use KindApplication instead.
	KindKubernetesApp Kind = "KubernetesApp"
	// KindTerraformApp represents application configuration for a Terraform application.
	// This application contains a single workspace of a terraform root module.
	//
	// Deprecated: use KindApplication instead.
	KindTerraformApp Kind = "TerraformApp"
	// KindLambdaApp represents application configuration for an AWS Lambda application.
	//
	// Deprecated: use KindApplication instead.
	KindLambdaApp Kind = "LambdaApp"
	// KindCloudRunApp represents application configuration for a CloudRun application.
	//
	// Deprecated: use KindApplication instead.
	KindCloudRunApp Kind = "CloudRunApp"
	// KindECSApp represents application configuration for an AWS ECS.
	//
	// Deprecated: use KindApplication instead.
	KindECSApp Kind = "ECSApp"
	// KindApplication represents a generic application configuration.
	KindApplication Kind = "Application"
)

const (
	// KindPiped represents configuration for piped.
	// This configuration will be loaded while the piped is starting up.
	KindPiped Kind = "Piped"
	// KindControlPlane represents configuration for control plane's services.
	KindControlPlane Kind = "ControlPlane"
	// KindAnalysisTemplate represents shared analysis template for a repository.
	// This configuration file should be placed in .pipe directory
	// at the root of the repository.
	KindAnalysisTemplate Kind = "AnalysisTemplate"
	// KindEventWatcher represents configuration for Event Watcher.
	KindEventWatcher Kind = "EventWatcher"
)

var (
	ErrNotFound = errors.New("not found")
)

// Spec[T] represents both of follows
// - the type is pointer type of T
// - the type has Validate method
type Spec[T any] interface {
	*T
	Validate() error
}

// Config represents configuration data load from file.
// The spec is depend on the kind of configuration.
type Config[T Spec[RT], RT any] struct {
	Kind       Kind
	APIVersion string
	Spec       T
}

func (c *Config[T, RT]) UnmarshalJSON(data []byte) error {
	// Define a type alias Config[T, RT] to avoid infinite recursion.
	type alias Config[T, RT]
	a := alias{}
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	*c = Config[T, RT](a)

	// Set default values.
	if c.Spec == nil {
		c.Spec = new(RT)
	}

	return nil
}

// Validate validates the value of all fields.
func (c *Config[T, RT]) Validate() error {
	if c.APIVersion != VersionV1Beta1 {
		return fmt.Errorf("unsupported version: %s", c.APIVersion)
	}
	if c.Kind == "" {
		return fmt.Errorf("kind is required")
	}

	if err := c.Spec.Validate(); err != nil {
		return err
	}
	return nil
}

// LoadFromYAML reads and decodes a yaml file to construct the Config.
func LoadFromYAML[T Spec[RT], RT any](file string) (*Config[T, RT], error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return DecodeYAML[T, RT](data)
}

// DecodeYAML unmarshals config YAML data to config struct.
// It also validates the configuration after decoding.
func DecodeYAML[T Spec[RT], RT any](data []byte) (*Config[T, RT], error) {
	js, err := yaml.YAMLToJSON(data)
	if err != nil {
		return nil, err
	}
	c := &Config[T, RT]{}
	if err := json.Unmarshal(js, c); err != nil {
		return nil, err
	}
	if err := defaults.Set(c); err != nil {
		return nil, err
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return c, nil
}
