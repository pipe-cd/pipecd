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
	"context"
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

func (r *Registry) Kubectl(ctx context.Context, version string) (string, error) {
	return r.client.InstallTool(ctx, "kubectl", version, kubectlInstallScript)
}

func (r *Registry) Kustomize(ctx context.Context, version string) (string, error) {
	return r.client.InstallTool(ctx, "kustomize", version, kustomizeInstallScript)
}

func (r *Registry) Helm(ctx context.Context, version string) (string, error) {
	return r.client.InstallTool(ctx, "helm", version, helmInstallScript)
}
