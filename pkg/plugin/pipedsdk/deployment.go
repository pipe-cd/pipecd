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

package pipedsdk

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister"
	"github.com/pipe-cd/pipecd/pkg/plugin/pipedapi"
)

var (
	deploymentServiceServer interface {
		Plugin

		Register(server *grpc.Server)
		setCommonFields(commonFields)
		deployment.DeploymentServiceServer
	}
)

// DeploymentPlugin is the interface that be implemented by a full-spec deployment plugin.
// This kind of plugin should implement all methods to manage resources and execute stages.
// The Config parameter is the plugin's config defined in piped's config.
type DeploymentPlugin[Config any] interface {
	PipelineSyncPlugin[Config]

	// DetermineVersions determines the versions of the resources that will be deployed.
	DetermineVersions(context.Context, *Config, *Client, TODO) (TODO, error)
	// DetermineStrategy determines the strategy to deploy the resources.
	DetermineStrategy(context.Context, *Config, *Client, TODO) (TODO, error)
	// BuildQuickSyncStages builds the stages that will be executed during the quick sync process.
	BuildQuickSyncStages(context.Context, *Config, *Client, TODO) (TODO, error)
}

// PipelineSyncPlugin is the interface implemented by a pipeline sync plugin.
// This kind of plugin may not implement quick sync stages, and will not manage resources like deployment plugin.
// It only focuses on executing stages which is generic for all kinds of pipeline sync plugins.
type PipelineSyncPlugin[Config any] interface {
	Plugin

	// FetchDefinedStages returns the list of stages that the plugin can execute.
	FetchDefinedStages() []string
	// BuildPipelineSyncStages builds the stages that will be executed by the plugin.
	BuildPipelineSyncStages(context.Context, *Config, *Client, TODO) (TODO, error)
	// ExecuteStage executes the given stage.
	ExecuteStage(context.Context, *Config, *Client, logpersister.StageLogPersister, TODO) (TODO, error)
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

type logPersister interface {
	StageLogPersister(deploymentID, stageID string) logpersister.StageLogPersister
}

type commonFields struct {
	config       *config.PipedPlugin
	logger       *zap.Logger
	logPersister logPersister
	client       *pipedapi.PipedServiceClient
}

// DeploymentPluginServiceServer is the gRPC server that handles requests from the piped.
type DeploymentPluginServiceServer[Config any] struct {
	deployment.UnimplementedDeploymentServiceServer
	commonFields

	base DeploymentPlugin[Config]
}

// Name returns the name of the plugin.
func (s *DeploymentPluginServiceServer[Config]) Name() string {
	return s.base.Name()
}

func (s *DeploymentPluginServiceServer[Config]) Version() string {
	return s.base.Version()
}

// Register registers the server to the given gRPC server.
func (s *DeploymentPluginServiceServer[Config]) Register(server *grpc.Server) {
	deployment.RegisterDeploymentServiceServer(server, s)
}

func (s *DeploymentPluginServiceServer[Config]) setCommonFields(fields commonFields) {
	s.commonFields = fields
}

func (s *DeploymentPluginServiceServer[Config]) FetchDefinedStages(context.Context, *deployment.FetchDefinedStagesRequest) (*deployment.FetchDefinedStagesResponse, error) {
	return &deployment.FetchDefinedStagesResponse{Stages: s.base.FetchDefinedStages()}, nil
}
func (s *DeploymentPluginServiceServer[Config]) DetermineVersions(context.Context, *deployment.DetermineVersionsRequest) (*deployment.DetermineVersionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DetermineVersions not implemented")
}
func (s *DeploymentPluginServiceServer[Config]) DetermineStrategy(context.Context, *deployment.DetermineStrategyRequest) (*deployment.DetermineStrategyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DetermineStrategy not implemented")
}
func (s *DeploymentPluginServiceServer[Config]) BuildPipelineSyncStages(context.Context, *deployment.BuildPipelineSyncStagesRequest) (*deployment.BuildPipelineSyncStagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BuildPipelineSyncStages not implemented")
}
func (s *DeploymentPluginServiceServer[Config]) BuildQuickSyncStages(context.Context, *deployment.BuildQuickSyncStagesRequest) (*deployment.BuildQuickSyncStagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BuildQuickSyncStages not implemented")
}
func (s *DeploymentPluginServiceServer[Config]) ExecuteStage(context.Context, *deployment.ExecuteStageRequest) (*deployment.ExecuteStageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteStage not implemented")
}

// PipelineSyncPluginServiceServer is the gRPC server that handles requests from the piped.
type PipelineSyncPluginServiceServer[Config any] struct {
	deployment.UnimplementedDeploymentServiceServer
	commonFields

	base PipelineSyncPlugin[Config]
}

// Name returns the name of the plugin.
func (s *PipelineSyncPluginServiceServer[Config]) Name() string {
	return s.base.Name()
}

// Version returns the version of the plugin.
func (s *PipelineSyncPluginServiceServer[Config]) Version() string {
	return s.base.Version()
}

// Register registers the server to the given gRPC server.
func (s *PipelineSyncPluginServiceServer[Config]) Register(server *grpc.Server) {
	deployment.RegisterDeploymentServiceServer(server, s)
}

func (s *PipelineSyncPluginServiceServer[Config]) setCommonFields(fields commonFields) {
	s.commonFields = fields
}

func (s *PipelineSyncPluginServiceServer[Config]) FetchDefinedStages(context.Context, *deployment.FetchDefinedStagesRequest) (*deployment.FetchDefinedStagesResponse, error) {
	return &deployment.FetchDefinedStagesResponse{Stages: s.base.FetchDefinedStages()}, nil
}
func (s *PipelineSyncPluginServiceServer[Config]) DetermineVersions(context.Context, *deployment.DetermineVersionsRequest) (*deployment.DetermineVersionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DetermineVersions not implemented")
}
func (s *PipelineSyncPluginServiceServer[Config]) DetermineStrategy(context.Context, *deployment.DetermineStrategyRequest) (*deployment.DetermineStrategyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DetermineStrategy not implemented")
}
func (s *PipelineSyncPluginServiceServer[Config]) BuildPipelineSyncStages(context.Context, *deployment.BuildPipelineSyncStagesRequest) (*deployment.BuildPipelineSyncStagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BuildPipelineSyncStages not implemented")
}
func (s *PipelineSyncPluginServiceServer[Config]) BuildQuickSyncStages(context.Context, *deployment.BuildQuickSyncStagesRequest) (*deployment.BuildQuickSyncStagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BuildQuickSyncStages not implemented")
}
func (s *PipelineSyncPluginServiceServer[Config]) ExecuteStage(context.Context, *deployment.ExecuteStageRequest) (*deployment.ExecuteStageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteStage not implemented")
}
