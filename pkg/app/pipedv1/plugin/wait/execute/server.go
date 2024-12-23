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

package execute

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/metadatastore"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister"
)

type deploymentServiceServer struct {
	deployment.UnimplementedDeploymentServiceServer

	pluginConfig *config.PipedPlugin

	metadataStore metadatastore.MetadataStore
	logger        *zap.Logger
	logPersister  logPersister
}

type logPersister interface {
	StageLogPersister(deploymentID, stageID string) logpersister.StageLogPersister
}

// NewDeploymentService creates a new planService.
func NewDeploymentService(
	config *config.PipedPlugin,
	logger *zap.Logger,
	logPersister logPersister,
) *deploymentServiceServer {
	return &deploymentServiceServer{
		pluginConfig: config,
		// TODO: Add metadataStore? or not?
		logger:       logger.Named("planner"), // TODO: Is this 'planner'?
		logPersister: logPersister,
	}
}

// Register registers all handling of this service into the specified gRPC server.
func (a *deploymentServiceServer) Register(server *grpc.Server) {
	deployment.RegisterDeploymentServiceServer(server, a)
}

// ExecuteStage implements deployment.ExecuteStage.
func (s *deploymentServiceServer) ExecuteStage(ctx context.Context, request *deployment.ExecuteStageRequest) (*deployment.ExecuteStageResponse, error) {
	slp := s.logPersister.StageLogPersister(request.Input.GetDeployment().GetId(), request.Input.GetStage().GetId())
	return s.execute(ctx, request.Input, slp)
}

// FetchDefinedStages implements deployment.FetchDefinedStages.
func (s *deploymentServiceServer) FetchDefinedStages(ctx context.Context, request *deployment.FetchDefinedStagesRequest) (*deployment.FetchDefinedStagesResponse, error) {
	return &deployment.FetchDefinedStagesResponse{
		Stages: []string{string(stageWait)},
	}, nil
}

// DetermineVersions implements deployment.DeploymentServiceServer.
func (s *deploymentServiceServer) DetermineVersions(ctx context.Context, request *deployment.DetermineVersionsRequest) (*deployment.DetermineVersionsResponse, error) {
	// TODO: Implement this func
	return &deployment.DetermineVersionsResponse{}, nil
}

// DetermineStrategy implements deployment.DeploymentServiceServer.
func (s *deploymentServiceServer) DetermineStrategy(ctx context.Context, request *deployment.DetermineStrategyRequest) (*deployment.DetermineStrategyResponse, error) {
	// TODO: Implement this func
	return &deployment.DetermineStrategyResponse{}, nil
}

// BuildPipelineSyncStages implements deployment.BuildPipelineSyncStages.
func (s *deploymentServiceServer) BuildPipelineSyncStages(ctx context.Context, request *deployment.BuildPipelineSyncStagesRequest) (*deployment.BuildPipelineSyncStagesResponse, error) {
	// TODO: Implement this func
	return &deployment.BuildPipelineSyncStagesResponse{}, nil
}

// BuildQuickSyncStages implements deployment.BuildQuickSyncStages.
func (s *deploymentServiceServer) BuildQuickSyncStages(ctx context.Context, request *deployment.BuildQuickSyncStagesRequest) (*deployment.BuildQuickSyncStagesResponse, error) {
	// TODO: Implement this func
	return &deployment.BuildQuickSyncStagesResponse{}, nil
}
