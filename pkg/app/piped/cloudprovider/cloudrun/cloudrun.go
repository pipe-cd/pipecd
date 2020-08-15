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

package cloudrun

import (
	"context"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/config"
)

const (
	DefaultServiceManifestFilename = "service.yaml"
)

type Provider interface {
	// LoadServiceManifest loads the service manifest
	// placing at the application configuration directory.
	LoadServiceManifest() (ServiceManifest, error)
	// Apply applies the service to the state specified in the give manifest.
	Apply(ctx context.Context, m ServiceManifest) error
}

type provider struct {
	appDir string
	input  config.CloudRunDeploymentInput
	logger *zap.Logger
}

func NewProvider(appDir string, input config.CloudRunDeploymentInput, logger *zap.Logger) *provider {
	return &provider{
		appDir: appDir,
		input:  input,
		logger: logger.Named("cloudrun-provider"),
	}
}

func (p *provider) LoadServiceManifest() (ServiceManifest, error) {
	filename := DefaultServiceManifestFilename
	if p.input.ServiceManifestFile != "" {
		filename = p.input.ServiceManifestFile
	}
	path := filepath.Join(p.appDir, filename)
	return LoadServiceManifest(path)
}

func (p *provider) Apply(ctx context.Context, sm ServiceManifest) error {
	cmd := NewGCloud("")
	return cmd.Apply(ctx, p.input.Platform, p.input.Region, sm)
}
