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

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/regexpool"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type toolRegistry interface {
	InstallTool(ctx context.Context, name, version string) (path string, err error)
}

type loader interface {
	// LoadManifests renders and loads all manifests for application.
	LoadManifests(ctx context.Context, input provider.LoaderInput) ([]provider.Manifest, error)
}

type DeploymentService struct {
	deployment.UnimplementedDeploymentServiceServer

	RegexPool    *regexpool.Pool
	Logger       *zap.Logger
	ToolRegistry toolRegistry
	Loader       loader
}

// NewDeploymentService creates a new planService.
func NewDeploymentService(
	logger *zap.Logger,
) *DeploymentService {
	return &DeploymentService{
		RegexPool:    regexpool.DefaultPool(),
		Logger:       logger.Named("planner"),
		ToolRegistry: nil, // TODO: set the tool registry
	}
}

// Register registers all handling of this service into the specified gRPC server.
func (a *DeploymentService) Register(server *grpc.Server) {
	deployment.RegisterDeploymentServiceServer(server, a)
}

// DetermineStrategy implements deployment.DeploymentServiceServer.
func (a *DeploymentService) DetermineStrategy(context.Context, *deployment.DetermineStrategyRequest) (*deployment.DetermineStrategyResponse, error) {
	panic("unimplemented")
}

// DetermineVersions implements deployment.DeploymentServiceServer.
func (a *DeploymentService) DetermineVersions(ctx context.Context, request *deployment.DetermineVersionsRequest) (*deployment.DetermineVersionsResponse, error) {
	manifests, err := a.Loader.LoadManifests(ctx, provider.LoaderInput{
		// TODO: fill the input
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	versions, err := determineVersions(manifests)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &deployment.DetermineVersionsResponse{
		Versions: versions,
	}, nil
}

// BuildPipelineSyncStages implements deployment.DeploymentServiceServer.
func (a *DeploymentService) BuildPipelineSyncStages(context.Context, *deployment.BuildPipelineSyncStagesRequest) (*deployment.BuildPipelineSyncStagesResponse, error) {
	panic("unimplemented")
}

// BuildQuickSyncStages implements deployment.DeploymentServiceServer.
func (a *DeploymentService) BuildQuickSyncStages(ctx context.Context, request *deployment.BuildQuickSyncStagesRequest) (*deployment.BuildQuickSyncStagesResponse, error) {
	now := time.Now()
	stages := buildQuickSyncPipeline(request.GetStageIndex(), request.GetRollback(), now)
	return &deployment.BuildQuickSyncStagesResponse{
		Stages: stages,
	}, nil
}

// FetchDefinedStages implements deployment.DeploymentServiceServer.
func (a *DeploymentService) FetchDefinedStages(context.Context, *deployment.FetchDefinedStagesRequest) (*deployment.FetchDefinedStagesResponse, error) {
	stages := make([]string, 0, len(AllStages))
	for _, s := range AllStages {
		stages = append(stages, string(s))
	}

	return &deployment.FetchDefinedStagesResponse{
		Stages: stages,
	}, nil
}
