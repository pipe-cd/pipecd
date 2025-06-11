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

package sdk

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/creasty/defaults"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/common"
)

// DeploymentSource represents the source of the deployment.
type DeploymentSource[Spec any] struct {
	// ApplicationDirectory is the directory where the source code is located.
	ApplicationDirectory string
	// CommitHash is the git commit hash of the source code.
	CommitHash string
	// ApplicationConfig is the configuration of the application.
	ApplicationConfig *ApplicationConfig[Spec]
	// ApplicationConfigFilename is the name of the file that contains the application configuration.
	// The plugins can use this to avoid mistakenly reading this file as a manifest to deploy.
	ApplicationConfigFilename string
}

// newDeploymentSource converts the common.DeploymentSource to the internal representation.
func newDeploymentSource[Spec any](pluginName string, source *common.DeploymentSource) (DeploymentSource[Spec], error) {
	cfg, err := config.DecodeYAML[*ApplicationConfig[Spec]](source.GetApplicationConfig())
	if err != nil {
		return DeploymentSource[Spec]{}, fmt.Errorf("failed to decode application config: %w", err)
	}

	if err := cfg.Spec.parsePluginConfig(pluginName); err != nil {
		return DeploymentSource[Spec]{}, fmt.Errorf("failed to parse plugin config: %w", err)
	}

	return DeploymentSource[Spec]{
		ApplicationDirectory:      source.GetApplicationDirectory(),
		CommitHash:                source.GetCommitHash(),
		ApplicationConfig:         cfg.Spec,
		ApplicationConfigFilename: source.GetApplicationConfigFilename(),
	}, nil
}

// AppConfig returns the application config.
func (d *DeploymentSource[Spec]) AppConfig() (*ApplicationConfig[Spec], error) {
	if d.ApplicationConfig == nil {
		return nil, fmt.Errorf("application config is not set")
	}
	return d.ApplicationConfig, nil
}

// ApplicationConfig is the configuration of the application.
type ApplicationConfig[Spec any] struct {
	// commonSpec is the common spec of the application.
	commonSpec *config.GenericApplicationSpec
	// pluginConfigs is the map of the plugin configs.
	// The key is the plugin name.
	// The value is the plugin config.
	pluginConfigs map[string]json.RawMessage
	// Spec is the plugin spec of the application.
	Spec *Spec
}

// LoadApplicationConfigForTest loads the application config from the given filename.
// When the error occurs, it will call t.Fatal/t.Fatalf and stop the test.
// This function is only used in the tests.
//
// NOTE: we want to put this function under package for testing like sdktest, but we can't do that
// because we don't want to make public parsePluginConfig in the ApplicationConfig struct.
func LoadApplicationConfigForTest[Spec any](t *testing.T, filename string, pluginName string) *ApplicationConfig[Spec] {
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read application config: %s", err)
	}
	cfg, err := config.DecodeYAML[*ApplicationConfig[Spec]](data)
	if err != nil {
		t.Fatalf("failed to decode application config: %s", err)
	}
	if cfg.Spec == nil {
		t.Fatal("application config is not set")
	}
	if err := cfg.Spec.parsePluginConfig(pluginName); err != nil {
		t.Fatalf("failed to parse plugin config: %s", err)
	}
	return cfg.Spec
}

// parsePluginConfig parses the plugin config for the given plugin name.
// It returns nil if no config is set for the plugin.
// After calling this method, the pluginConfigs is cleared to avoid leaking the internal data.
func (c *ApplicationConfig[Spec]) parsePluginConfig(pluginName string) error {
	defer func() {
		// Clear the plugin configs after using it to avoid leaking the internal data.
		c.pluginConfigs = nil
	}()

	// It is necessary to prepare config with default value when users doesn't set any config, or when developers implements custom unmarshalling logic.
	data := []byte("{}")

	if c.pluginConfigs != nil && c.pluginConfigs[pluginName] != nil {
		data = c.pluginConfigs[pluginName]
	}

	var spec Spec
	if err := json.Unmarshal(data, &spec); err != nil {
		return fmt.Errorf("failed to unmarshal application config: plugin spec: %w", err)
	}

	if err := defaults.Set(&spec); err != nil {
		return fmt.Errorf("failed to set default values for plugin spec: %w", err)
	}

	// Validate the spec if it implements the Validate method.
	if v, ok := any(spec).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("failed to validate plugin spec: %w", err)
		}
	}

	// Sometimes the receiver of Validate method is pointer to the spec.
	if v, ok := any(&spec).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("failed to validate plugin spec: %w", err)
		}
	}

	c.Spec = &spec

	return nil
}

func (c *ApplicationConfig[Spec]) UnmarshalJSON(data []byte) error {
	if c.commonSpec == nil {
		c.commonSpec = new(config.GenericApplicationSpec)
	}

	if err := json.Unmarshal(data, c.commonSpec); err != nil {
		return fmt.Errorf("failed to unmarshal application config: generic spec: %w", err)
	}

	type pluginSpecs struct {
		Plugins map[string]json.RawMessage `json:"plugins"`
	}

	var p pluginSpecs
	if err := json.Unmarshal(data, &p); err != nil {
		return fmt.Errorf("failed to unmarshal application config: plugin specs: %w", err)
	}

	c.pluginConfigs = p.Plugins

	return nil
}

func (c *ApplicationConfig[Spec]) Validate() error {
	if c.commonSpec == nil {
		return fmt.Errorf("application config is not initialized")
	}

	if err := c.commonSpec.Validate(); err != nil {
		return fmt.Errorf("validation failed on generic spec: %w", err)
	}

	return nil
}

// HasStage returns true if the application config has a stage with the given name.
func (c *ApplicationConfig[Spec]) HasStage(name string) bool {
	if c.commonSpec.Pipeline == nil {
		return false
	}
	// linear search is enough because the number of stages is limited.
	for _, stage := range c.commonSpec.Pipeline.Stages {
		if string(stage.Name) == name {
			return true
		}
	}
	return false
}
