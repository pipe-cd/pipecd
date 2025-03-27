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

// Package toolregistry installs and manages the needed tools
// such as kubectl, helm... for executing tasks in pipeline.
package toolregistry

import (
	"cmp"
	"context"
)

const (
	defaultKubectlVersion   = "1.18.2"
	defaultKustomizeVersion = "3.8.1"
	defaultHelmVersion      = "3.8.2"
)

type client interface {
	InstallTool(ctx context.Context, name, version, script string) (string, error)
}

func NewRegistry(client client) *Registry {
	return &Registry{
		client: client,
	}
}

// Registry provides functions to get path to the needed tools.
type Registry struct {
	client client
}

// Kubectl installs the kubectl tool with the given version and return the path to the installed binary.
// If the version is empty, the default version will be used.
func (r *Registry) Kubectl(ctx context.Context, version string) (string, error) {
	return r.client.InstallTool(ctx, "kubectl", cmp.Or(version, defaultKubectlVersion), kubectlInstallScript)
}

// Kustringize installs the kustomize tool with the given version and return the path to the installed binary.
// If the version is empty, the default version will be used.
func (r *Registry) Kustomize(ctx context.Context, version string) (string, error) {
	return r.client.InstallTool(ctx, "kustomize", cmp.Or(version, defaultKustomizeVersion), kustomizeInstallScript)
}

// Helm installs the helm tool with the given version and return the path to the installed binary.
// If the version is empty, the default version will be used.
func (r *Registry) Helm(ctx context.Context, version string) (string, error) {
	return r.client.InstallTool(ctx, "helm", cmp.Or(version, defaultHelmVersion), helmInstallScript)
}
