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

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/toolregistry"
)

// Registry provides functions to get path to the needed tools.
type Registry interface {
	Kubectl(ctx context.Context, version string) (string, error)
	Kustomize(ctx context.Context, version string) (string, error)
	Helm(ctx context.Context, version string) (string, error)
}

func NewRegistry(client toolregistry.ToolRegistry) Registry {
	return &registry{
		client: client,
	}
}

type registry struct {
	client toolregistry.ToolRegistry
}

func (r *registry) Kubectl(ctx context.Context, version string) (string, error) {
	return r.client.InstallTool(ctx, "kubectl", version, kubectlInstallScript)
}

func (r *registry) Kustomize(ctx context.Context, version string) (string, error) {
	return r.client.InstallTool(ctx, "kustomize", version, kustomizeInstallScript)
}

func (r *registry) Helm(ctx context.Context, version string) (string, error) {
	return r.client.InstallTool(ctx, "helm", version, helmInstallScript)
}
