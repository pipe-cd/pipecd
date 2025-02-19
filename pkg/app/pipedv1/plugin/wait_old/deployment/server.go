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

package deployment

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister"
	"github.com/pipe-cd/pipecd/pkg/plugin/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/plugin/signalhandler"
)

type deploymentServiceServer struct {
	deployment.UnimplementedDeploymentServiceServer

	pluginConfig *config.PipedPlugin

	logger        *zap.Logger
	logPersister  logPersister
	metadataStore metadataStoreClient
}

type logPersister interface {
	StageLogPersister(deploymentID, stageID string) logpersister.StageLogPersister
}

type metadataStoreClient interface {
	GetStageMetadata(ctx context.Context, in *pipedservice.GetStageMetadataRequest, opts ...grpc.CallOption) (*pipedservice.GetStageMetadataResponse, error)
	PutStageMetadata(ctx context.Context, in *pipedservice.PutStageMetadataRequest, opts ...grpc.CallOption) (*pipedservice.PutStageMetadataResponse, error)
}

// NewDeploymentService creates a new deploymentServiceServer of Wait Stage plugin.
func NewDeploymentService(
	config *config.PipedPlugin,
	logger *zap.Logger,
	logPersister logPersister,
	metadataStore metadataStoreClient,
) *deploymentServiceServer {
	return &deploymentServiceServer{
		pluginConfig:  config,
		logger:        logger.Named("wait-stage-plugin"),
		logPersister:  logPersister,
		metadataStore: metadataStore,
	}
}

// Register registers all handling of this service into the specified gRPC server.
func (s *deploymentServiceServer) Register(server *grpc.Server) {
	deployment.RegisterDeploymentServiceServer(server, s)
}

// ExecuteStage implements deployment.ExecuteStage.
func (s *deploymentServiceServer) ExecuteStage(ctx context.Context, request *deployment.ExecuteStageRequest) (response *deployment.ExecuteStageResponse, err error) {
	slp := s.logPersister.StageLogPersister(request.Input.GetDeployment().GetId(), request.Input.GetStage().GetId())
	defer func() {
		// When termination signal received and the stage is not completed yet, we should not mark the log persister as completed.
		// This can occur when the piped is shutting down while the stage is still running.
		if !response.GetStatus().IsCompleted() && signalhandler.Terminated() {
			return
		}
		slp.Complete(time.Minute)
	}()

	status := s.execute(ctx, request.Input, slp)
	return &deployment.ExecuteStageResponse{
		Status: status,
	}, nil
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
