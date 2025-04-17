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
func newDeploymentSource[Spec any](source *common.DeploymentSource) (DeploymentSource[Spec], error) {
	cfg, err := config.DecodeYAML[*ApplicationConfig[Spec]](source.GetApplicationConfig())
	if err != nil {
		return DeploymentSource[Spec]{}, fmt.Errorf("failed to decode application config: %w", err)
	}
	return DeploymentSource[Spec]{
		ApplicationDirectory:      source.GetApplicationDirectory(),
		CommitHash:                source.GetCommitHash(),
		ApplicationConfig:         cfg.Spec,
		ApplicationConfigFilename: source.GetApplicationConfigFilename(),
	}, nil
}

// LoadApplicationConfig loads the application config from the given filename.
// This is intended to use in the tests.
func LoadApplicationConfig[Spec any](filename string) (*ApplicationConfig[Spec], error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read application config: %w", err)
	}
	cfg, err := config.DecodeYAML[*ApplicationConfig[Spec]](data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode application config: %w", err)
	}
	if cfg.Spec == nil {
		return nil, fmt.Errorf("application config is not set")
	}
	return cfg.Spec, nil
}

// ApplicationConfig is the configuration of the application.
type ApplicationConfig[Spec any] struct {
	// commonSpec is the common spec of the application.
	commonSpec *config.GenericApplicationSpec
	// Spec is the plugin spec of the application.
	Spec *Spec
}

func (c *ApplicationConfig[Spec]) UnmarshalJSON(data []byte) error {
	if c.commonSpec == nil {
		c.commonSpec = new(config.GenericApplicationSpec)
	}

	if err := json.Unmarshal(data, c.commonSpec); err != nil {
		return fmt.Errorf("failed to unmarshal application config: generic spec: %w", err)
	}

	if c.Spec == nil {
		c.Spec = new(Spec)
	}

	if err := json.Unmarshal(data, c.Spec); err != nil {
		return fmt.Errorf("failed to unmarshal application config: plugin spec: %w", err)
	}

	return nil
}

func (c *ApplicationConfig[Spec]) Validate() error {
	if c.commonSpec == nil {
		return fmt.Errorf("application config is not initialized")
	}

	if c.Spec == nil {
		return fmt.Errorf("plugin spec is not initialized")
	}

	if err := c.commonSpec.Validate(); err != nil {
		return fmt.Errorf("validation failed on generic spec: %w", err)
	}

	if v, ok := any(c.Spec).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("validation failed on plugin spec: %w", err)
		}
	}

	return nil
}
