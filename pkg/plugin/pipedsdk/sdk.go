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

// package pipedsdk provides software development kits for building PipeCD piped plugins.
package pipedsdk

import "context"

// TODO is a placeholder for the real type.
// This type will be replaced by the real type when implementing the sdk.
type TODO struct{}

// Client is a placeholder for the real client.
// This type will be replaced by the real client to gRPC piped-plugin-service when implementing the sdk.
type Client struct{}

// DeploymentPlugin is the interface that be implemented by a full-spec deployment plugin.
// This kind of plugin should implement all methods to manage resources and execute stages.
type DeploymentPlugin[Config any] interface {
	PipelineSyncPlugin[Config]

	DetermineVersions(context.Context, Config, Client, TODO) (TODO, error)
	DetermineStrategy(context.Context, Config, Client, TODO) (TODO, error)
	BuildQuickSyncStages(context.Context, Config, Client, TODO) (TODO, error)
}

// PipelineSyncPlugin is the interface that be implemented by a pipeline sync plugin.
// This kind of plugin may not implement quick sync stages, and will not manage resources like deployment plugin.
// It only focuses on executing stages which is generic for all kinds of pipeline sync plugins.
type PipelineSyncPlugin[Config any] interface {
	FetchDefinedStages() []string
	BuildPipelineSyncStages(context.Context, Config, Client, TODO) (TODO, error)
	ExecuteStage(context.Context, Config, Client, TODO) (TODO, error)
}

// RegisterDeploymentPlugin registers the given deployment plugin.
// It will be used when running the piped.
func RegisterDeploymentPlugin[Config any](plugin DeploymentPlugin[Config]) {
	panic("implement me")
}

// RegisterPipelineSyncPlugin registers the given pipeline sync plugin.
// It will be used when running the piped.
func RegisterPipelineSyncPlugin[Config any](plugin PipelineSyncPlugin[Config]) {
	panic("implement me")
}

// Run runs the registered plugins.
// It will listen the gRPC server and handle all requests from piped.
func Run() error {
	panic("implement me")
}
