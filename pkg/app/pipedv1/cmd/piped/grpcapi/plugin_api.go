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

package grpcapi

import (
	"context"
	"fmt"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/metadatastore"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/model"
	service "github.com/pipe-cd/pipecd/pkg/plugin/pipedservice"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PluginAPI struct {
	service.PluginServiceServer

	cfg       *config.PipedSpec
	apiClient apiClient

	toolRegistry          *toolRegistry
	Logger                *zap.Logger
	metadataStoreRegistry *metadatastore.MetadataStoreRegistry
	stageCommandLister    stageCommandLister
}

type apiClient interface {
	ReportStageLogs(ctx context.Context, req *pipedservice.ReportStageLogsRequest, opts ...grpc.CallOption) (*pipedservice.ReportStageLogsResponse, error)
	ReportStageLogsFromLastCheckpoint(ctx context.Context, in *pipedservice.ReportStageLogsFromLastCheckpointRequest, opts ...grpc.CallOption) (*pipedservice.ReportStageLogsFromLastCheckpointResponse, error)
	GetApplicationSharedObject(ctx context.Context, req *pipedservice.GetApplicationSharedObjectRequest, opts ...grpc.CallOption) (*pipedservice.GetApplicationSharedObjectResponse, error)
	PutApplicationSharedObject(ctx context.Context, req *pipedservice.PutApplicationSharedObjectRequest, opts ...grpc.CallOption) (*pipedservice.PutApplicationSharedObjectResponse, error)
}

type stageCommandLister interface {
	ListStageCommands(deploymentID, stageID string) []*model.Command
}

// Register registers all handling of this service into the specified gRPC server.
func (a *PluginAPI) Register(server *grpc.Server) {
	service.RegisterPluginServiceServer(server, a)
}

func NewPluginAPI(cfg *config.PipedSpec, apiClient apiClient, toolsDir string, logger *zap.Logger, metadataStoreRegistry *metadatastore.MetadataStoreRegistry, stageCommandLister stageCommandLister) (*PluginAPI, error) {
	toolRegistry, err := newToolRegistry(toolsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create tool registry: %w", err)
	}

	return &PluginAPI{
		cfg:                   cfg,
		apiClient:             apiClient,
		toolRegistry:          toolRegistry,
		Logger:                logger.Named("plugin-api"),
		metadataStoreRegistry: metadataStoreRegistry,
		stageCommandLister:    stageCommandLister,
	}, nil
}

// InstallTool installs the given tool.
// installed binary's filename becomes `name-version`.
func (a *PluginAPI) InstallTool(ctx context.Context, req *service.InstallToolRequest) (*service.InstallToolResponse, error) {
	p, err := a.toolRegistry.InstallTool(ctx, req.GetName(), req.GetVersion(), req.GetInstallScript())
	if err != nil {
		a.Logger.Error("failed to install tool",
			zap.String("name", req.GetName()),
			zap.String("version", req.GetVersion()),
			zap.Error(err))
		return nil, err
	}
	return &service.InstallToolResponse{
		InstalledPath: p,
	}, nil
}

func (a *PluginAPI) ReportStageLogs(ctx context.Context, req *service.ReportStageLogsRequest) (*service.ReportStageLogsResponse, error) {
	_, err := a.apiClient.ReportStageLogs(ctx, &pipedservice.ReportStageLogsRequest{
		DeploymentId: req.DeploymentId,
		StageId:      req.StageId,
		RetriedCount: req.RetriedCount,
		Blocks:       req.Blocks,
	})
	if err != nil {
		a.Logger.Error("failed to report stage logs",
			zap.String("deploymentID", req.DeploymentId),
			zap.String("stageID", req.StageId),
			zap.Error(err))
		return nil, err
	}

	return &service.ReportStageLogsResponse{}, nil
}

func (a *PluginAPI) ReportStageLogsFromLastCheckpoint(ctx context.Context, req *service.ReportStageLogsFromLastCheckpointRequest) (*service.ReportStageLogsFromLastCheckpointResponse, error) {
	_, err := a.apiClient.ReportStageLogsFromLastCheckpoint(ctx, &pipedservice.ReportStageLogsFromLastCheckpointRequest{
		DeploymentId: req.DeploymentId,
		StageId:      req.StageId,
		RetriedCount: req.RetriedCount,
		Blocks:       req.Blocks,
		Completed:    req.Completed,
	})
	if err != nil {
		a.Logger.Error("failed to report stage logs from last checkpoint",
			zap.String("deploymentID", req.DeploymentId),
			zap.String("stageID", req.StageId),
			zap.Error(err))
		return nil, err
	}

	return &service.ReportStageLogsFromLastCheckpointResponse{}, nil
}

func (a *PluginAPI) GetStageMetadata(ctx context.Context, req *service.GetStageMetadataRequest) (*service.GetStageMetadataResponse, error) {
	return a.metadataStoreRegistry.GetStageMetadata(ctx, req)
}

func (a *PluginAPI) PutStageMetadata(ctx context.Context, req *service.PutStageMetadataRequest) (*service.PutStageMetadataResponse, error) {
	return a.metadataStoreRegistry.PutStageMetadata(ctx, req)
}

func (a *PluginAPI) PutStageMetadataMulti(ctx context.Context, req *service.PutStageMetadataMultiRequest) (*service.PutStageMetadataMultiResponse, error) {
	return a.metadataStoreRegistry.PutStageMetadataMulti(ctx, req)
}

func (a *PluginAPI) GetDeploymentPluginMetadata(ctx context.Context, req *service.GetDeploymentPluginMetadataRequest) (*service.GetDeploymentPluginMetadataResponse, error) {
	return a.metadataStoreRegistry.GetDeploymentPluginMetadata(ctx, req)
}

func (a *PluginAPI) PutDeploymentPluginMetadata(ctx context.Context, req *service.PutDeploymentPluginMetadataRequest) (*service.PutDeploymentPluginMetadataResponse, error) {
	return a.metadataStoreRegistry.PutDeploymentPluginMetadata(ctx, req)
}

func (a *PluginAPI) PutDeploymentPluginMetadataMulti(ctx context.Context, req *service.PutDeploymentPluginMetadataMultiRequest) (*service.PutDeploymentPluginMetadataMultiResponse, error) {
	return a.metadataStoreRegistry.PutDeploymentPluginMetadataMulti(ctx, req)
}

func (a *PluginAPI) GetDeploymentSharedMetadata(ctx context.Context, req *service.GetDeploymentSharedMetadataRequest) (*service.GetDeploymentSharedMetadataResponse, error) {
	return a.metadataStoreRegistry.GetDeploymentSharedMetadata(ctx, req)
}

func (a *PluginAPI) ListStageCommands(ctx context.Context, req *service.ListStageCommandsRequest) (*service.ListStageCommandsResponse, error) {
	commands := a.stageCommandLister.ListStageCommands(req.DeploymentId, req.StageId)
	return &service.ListStageCommandsResponse{Commands: commands}, nil
}

func (a *PluginAPI) GetApplicationSharedObject(ctx context.Context, req *service.GetApplicationSharedObjectRequest) (*service.GetApplicationSharedObjectResponse, error) {
	resp, err := a.apiClient.GetApplicationSharedObject(ctx, &pipedservice.GetApplicationSharedObjectRequest{
		ApplicationId: req.ApplicationId,
		PluginName:    req.PluginName,
		Key:           req.Key,
	})
	if status.Code(err) == codes.NotFound {
		return nil, status.Error(codes.NotFound, "the requested application shared object was not found")
	}
	if err != nil {
		a.Logger.Error("failed to get application shared object",
			zap.String("applicationID", req.ApplicationId),
			zap.String("pluginName", req.PluginName),
			zap.String("key", req.Key),
			zap.Error(err))
		return nil, err
	}
	return &service.GetApplicationSharedObjectResponse{
		Object: resp.Object,
	}, nil
}

func (a *PluginAPI) PutApplicationSharedObject(ctx context.Context, req *service.PutApplicationSharedObjectRequest) (*service.PutApplicationSharedObjectResponse, error) {
	_, err := a.apiClient.PutApplicationSharedObject(ctx, &pipedservice.PutApplicationSharedObjectRequest{
		ApplicationId: req.ApplicationId,
		PluginName:    req.PluginName,
		Key:           req.Key,
		Object:        req.Object,
	})
	if err != nil {
		a.Logger.Error("failed to put application shared object",
			zap.String("applicationID", req.ApplicationId),
			zap.String("pluginName", req.PluginName),
			zap.String("key", req.Key),
			zap.Int("data-size", len(req.Object)),
			zap.Error(err))
		return nil, err
	}
	return &service.PutApplicationSharedObjectResponse{}, nil
}
