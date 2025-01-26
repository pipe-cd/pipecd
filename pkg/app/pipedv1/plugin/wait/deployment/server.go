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
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister"
	"github.com/pipe-cd/pipecd/pkg/plugin/signalhandler"
)

type Stage string

const (
	stageWait Stage = "WAIT"
)

type deploymentServiceServer struct {
	deployment.UnimplementedDeploymentServiceServer

	pluginConfig *config.PipedPlugin

	logger       *zap.Logger
	logPersister logPersister
}

type logPersister interface {
	StageLogPersister(deploymentID, stageID string) logpersister.StageLogPersister
}

// NewDeploymentService creates a new deploymentServiceServer of Wait Stage plugin.
func NewDeploymentService(
	config *config.PipedPlugin,
	logger *zap.Logger,
	logPersister logPersister,
) *deploymentServiceServer {
	return &deploymentServiceServer{
		pluginConfig: config,
		logger:       logger.Named("wait-stage-plugin"),
		logPersister: logPersister,
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
	// Wait Stage does not have any versioned resources.
	return &deployment.DetermineVersionsResponse{
		Versions: []*model.ArtifactVersion{},
	}, nil
}

// DetermineStrategy implements deployment.DeploymentServiceServer.
func (s *deploymentServiceServer) DetermineStrategy(ctx context.Context, request *deployment.DetermineStrategyRequest) (*deployment.DetermineStrategyResponse, error) {
	return &deployment.DetermineStrategyResponse{Unsupported: true}, nil
}

// BuildPipelineSyncStages implements deployment.BuildPipelineSyncStages.
func (s *deploymentServiceServer) BuildPipelineSyncStages(ctx context.Context, request *deployment.BuildPipelineSyncStagesRequest) (*deployment.BuildPipelineSyncStagesResponse, error) {
	stages := make([]*model.PipelineStage, 0, len(request.GetStages()))
	for _, stage := range request.GetStages() {
		waitStage := newWaitStage()

		id := stage.GetId()
		if id == "" {
			id = fmt.Sprintf("stage-%d", stage.GetIndex())
		}
		waitStage.Id = id
		waitStage.Index = stage.GetIndex()
		stages = append(stages, waitStage)
	}

	return &deployment.BuildPipelineSyncStagesResponse{Stages: stages}, nil
}

// BuildQuickSyncStages implements deployment.BuildQuickSyncStages.
func (s *deploymentServiceServer) BuildQuickSyncStages(ctx context.Context, request *deployment.BuildQuickSyncStagesRequest) (*deployment.BuildQuickSyncStagesResponse, error) {
	return &deployment.BuildQuickSyncStagesResponse{
		Stages: []*model.PipelineStage{},
	}, nil
}

// newWaitStage returns a new WAIT stage with the current time.
// WAIT Stages is not used in the Rollback.
func newWaitStage() *model.PipelineStage {
	now := time.Now()
	return &model.PipelineStage{
		Id:        string(stageWait),
		Name:      string(stageWait),
		Desc:      "Wait for the specified duration",
		Rollback:  false,
		Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
		CreatedAt: now.Unix(),
		UpdatedAt: now.Unix(),
	}
}
