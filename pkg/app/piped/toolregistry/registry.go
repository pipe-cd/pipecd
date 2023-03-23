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

// Package toolregistry installs and manages the needed tools
// such as kubectl, helm... for executing tasks in pipeline.
package toolregistry

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

// Registry provides functions to get path to the needed tools.
type Registry interface {
	Kubectl(ctx context.Context, version string) (string, bool, error)
	Kustomize(ctx context.Context, version string) (string, bool, error)
	Helm(ctx context.Context, version string) (string, bool, error)
	Terraform(ctx context.Context, version string) (string, bool, error)
	GetBinDir() string
}

var defaultRegistry *registry

// DefaultRegistry returns the shared registry.
func DefaultRegistry() Registry {
	return defaultRegistry
}

// InitDefaultRegistry initializes the default registry.
// This also preloads the pre-installed tools in the binDir.
func InitDefaultRegistry(binDir string, logger *zap.Logger) error {
	logger = logger.Named("tool-registry")
	if err := os.MkdirAll(binDir, os.ModePerm); err != nil {
		return err
	}

	tools, err := loadPreinstalledTool(binDir)
	if err != nil {
		return err
	}
	logger.Info("successfully loaded the pre-installed tools", zap.Any("tools", tools))

	defaultRegistry = &registry{
		binDir:       binDir,
		versions:     tools,
		installGroup: &singleflight.Group{},
		logger:       logger,
	}

	return nil
}

func loadPreinstalledTool(binDir string) (map[string]struct{}, error) {
	tools := make(map[string]struct{})
	err := filepath.Walk(binDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == binDir {
			return nil
		}
		if info.IsDir() {
			return filepath.SkipDir
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		name := filepath.Base(path)
		tools[name] = struct{}{}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tools, nil
}

const (
	kubectlPrefix   = "kubectl"
	kustomizePrefix = "kustomize"
	helmPrefix      = "helm"
	terraformPrefix = "terraform"
)

type registry struct {
	binDir       string
	versions     map[string]struct{}
	mu           sync.RWMutex
	installGroup *singleflight.Group
	logger       *zap.Logger
}

func (r *registry) Kubectl(ctx context.Context, version string) (string, bool, error) {
	name := kubectlPrefix
	if version != "" {
		name = fmt.Sprintf("%s-%s", kubectlPrefix, version)
	}
	path := filepath.Join(r.binDir, name)

	r.mu.RLock()
	_, ok := r.versions[name]
	r.mu.RUnlock()
	if ok {
		return path, false, nil
	}

	_, err, _ := r.installGroup.Do(name, func() (interface{}, error) {
		return nil, r.installKubectl(ctx, version)
	})
	if err != nil {
		return "", true, err
	}

	r.mu.Lock()
	r.versions[name] = struct{}{}
	r.mu.Unlock()

	return path, true, nil
}

func (r *registry) Kustomize(ctx context.Context, version string) (string, bool, error) {
	name := kustomizePrefix
	if version != "" {
		name = fmt.Sprintf("%s-%s", kustomizePrefix, version)
	}
	path := filepath.Join(r.binDir, name)

	r.mu.RLock()
	_, ok := r.versions[name]
	r.mu.RUnlock()
	if ok {
		return path, false, nil
	}

	_, err, _ := r.installGroup.Do(name, func() (interface{}, error) {
		return nil, r.installKustomize(ctx, version)
	})
	if err != nil {
		return "", true, err
	}

	r.mu.Lock()
	r.versions[name] = struct{}{}
	r.mu.Unlock()

	return path, true, nil
}

func (r *registry) Helm(ctx context.Context, version string) (string, bool, error) {
	name := helmPrefix
	if version != "" {
		name = fmt.Sprintf("%s-%s", helmPrefix, version)
	}
	path := filepath.Join(r.binDir, name)

	r.mu.RLock()
	_, ok := r.versions[name]
	r.mu.RUnlock()
	if ok {
		return path, false, nil
	}

	_, err, _ := r.installGroup.Do(name, func() (interface{}, error) {
		return nil, r.installHelm(ctx, version)
	})
	if err != nil {
		return "", true, err
	}

	r.mu.Lock()
	r.versions[name] = struct{}{}
	r.mu.Unlock()

	return path, true, nil
}

func (r *registry) Terraform(ctx context.Context, version string) (string, bool, error) {
	name := terraformPrefix
	if version != "" {
		name = fmt.Sprintf("%s-%s", terraformPrefix, version)
	}
	path := filepath.Join(r.binDir, name)

	r.mu.RLock()
	_, ok := r.versions[name]
	r.mu.RUnlock()
	if ok {
		return path, false, nil
	}

	_, err, _ := r.installGroup.Do(name, func() (interface{}, error) {
		return nil, r.installTerraform(ctx, version)
	})
	if err != nil {
		return "", true, err
	}

	r.mu.Lock()
	r.versions[name] = struct{}{}
	r.mu.Unlock()

	return path, true, nil
}

func (r *registry) GetBinDir() string {
	return r.binDir
}
