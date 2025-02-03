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

package pipedsdk

import (
	"context"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister"
	"github.com/pipe-cd/pipecd/pkg/plugin/pipedapi"
	"google.golang.org/grpc"
)

var (
	deploymentServiceServer interface {
		Register(server *grpc.Server)
		deployment.DeploymentServiceServer
	}
)

// DeploymentPlugin is the interface that be implemented by a full-spec deployment plugin.
// This kind of plugin should implement all methods to manage resources and execute stages.
type DeploymentPlugin[Config any] interface {
	PipelineSyncPlugin[Config]

	// DetermineVersions determines the versions of the resources that will be deployed.
	DetermineVersions(context.Context, Config, Client, TODO) (TODO, error)
	// DetermineStrategy determines the strategy to deploy the resources.
	DetermineStrategy(context.Context, Config, Client, TODO) (TODO, error)
	// BuildQuickSyncStages builds the stages that will be executed during the quick sync process.
	BuildQuickSyncStages(context.Context, Config, Client, TODO) (TODO, error)
}

// PipelineSyncPlugin is the interface that be implemented by a pipeline sync plugin.
// This kind of plugin may not implement quick sync stages, and will not manage resources like deployment plugin.
// It only focuses on executing stages which is generic for all kinds of pipeline sync plugins.
type PipelineSyncPlugin[Config any] interface {
	// Name returns the name of the plugin.
	Name() string
	// FetchDefinedStages returns the list of stages that the plugin can execute.
	FetchDefinedStages() []string
	// BuildPipelineSyncStages builds the stages that will be executed by the plugin.
	BuildPipelineSyncStages(context.Context, Config, Client, TODO) (TODO, error)
	// ExecuteStage executes the given stage.
	ExecuteStage(context.Context, Config, Client, TODO) (TODO, error)
}

// RegisterDeploymentPlugin registers the given deployment plugin.
// It will be used when running the piped.
func RegisterDeploymentPlugin[Config any](plugin DeploymentPlugin[Config]) {
	deploymentServiceServer = &DeploymentPluginServiceServer[Config]{base: plugin}
}

// RegisterPipelineSyncPlugin registers the given pipeline sync plugin.
// It will be used when running the piped.
func RegisterPipelineSyncPlugin[Config any](plugin PipelineSyncPlugin[Config]) {
	deploymentServiceServer = &PipelineSyncPluginServiceServer[Config]{base: plugin}
}

// DeploymentPluginServiceServer is the gRPC server that handles requests from the piped.
type DeploymentPluginServiceServer[Config any] struct {
	deployment.UnimplementedDeploymentServiceServer

	base DeploymentPlugin[Config]
	config *config.PipedPlugin
	logPersister logpersister.Persister
	client pipedapi.PipedServiceClient
}

// Register registers the server to the given gRPC server.
func (s *DeploymentPluginServiceServer[Config]) Register(server *grpc.Server) {
	deployment.RegisterDeploymentServiceServer(server, s)
}

// PipelineSyncPluginServiceServer is the gRPC server that handles requests from the piped.
type PipelineSyncPluginServiceServer[Config any] struct {
	deployment.UnimplementedDeploymentServiceServer

	base PipelineSyncPlugin[Config]
	config *config.PipedPlugin
	logPersister logpersister.Persister
	client pipedapi.PipedServiceClient
}

// Register registers the server to the given gRPC server.
func (s *PipelineSyncPluginServiceServer[Config]) Register(server *grpc.Server) {
	deployment.RegisterDeploymentServiceServer(server, s)
}
