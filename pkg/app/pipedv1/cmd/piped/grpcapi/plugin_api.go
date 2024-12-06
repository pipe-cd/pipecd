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

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/cmd/piped/service"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/crypto"
	"github.com/pipe-cd/pipecd/pkg/model"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type PluginAPI struct {
	service.PluginServiceServer

	cfg       *config.PipedSpec
	apiClient apiClient

	toolRegistry *toolRegistry
	Logger       *zap.Logger
}

type apiClient interface {
	ReportStageLogs(ctx context.Context, req *pipedservice.ReportStageLogsRequest, opts ...grpc.CallOption) (*pipedservice.ReportStageLogsResponse, error)
	ReportStageLogsFromLastCheckpoint(ctx context.Context, in *pipedservice.ReportStageLogsFromLastCheckpointRequest, opts ...grpc.CallOption) (*pipedservice.ReportStageLogsFromLastCheckpointResponse, error)
}

// Register registers all handling of this service into the specified gRPC server.
func (a *PluginAPI) Register(server *grpc.Server) {
	service.RegisterPluginServiceServer(server, a)
}

func NewPluginAPI(cfg *config.PipedSpec, apiClient apiClient, toolsDir string, logger *zap.Logger) (*PluginAPI, error) {
	toolRegistry, err := newToolRegistry(toolsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create tool registry: %w", err)
	}

	return &PluginAPI{
		cfg:          cfg,
		apiClient:    apiClient,
		toolRegistry: toolRegistry,
		Logger:       logger.Named("plugin-api"),
	}, nil
}

func (a *PluginAPI) DecryptSecret(ctx context.Context, req *service.DecryptSecretRequest) (*service.DecryptSecretResponse, error) {
	decrypter, err := initializeSecretDecrypter(a.cfg.SecretManagement)
	if err != nil {
		a.Logger.Error("failed to initialize secret decrypter", zap.Error(err))
		return nil, err
	}

	// Return the secret as is in case of no decrypter configured.
	if decrypter == nil {
		return &service.DecryptSecretResponse{
			DecryptedSecret: req.Secret,
		}, nil
	}

	decrypted, err := decrypter.Decrypt(req.Secret)
	if err != nil {
		a.Logger.Error("failed to decrypt the secret", zap.Error(err))
		return nil, err
	}

	return &service.DecryptSecretResponse{
		DecryptedSecret: decrypted,
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

func initializeSecretDecrypter(sm *config.SecretManagement) (crypto.Decrypter, error) {
	if sm == nil {
		return nil, nil
	}

	switch sm.Type {
	case model.SecretManagementTypeNone:
		return nil, nil

	case model.SecretManagementTypeKeyPair:
		key, err := sm.KeyPair.LoadPrivateKey()
		if err != nil {
			return nil, err
		}
		decrypter, err := crypto.NewHybridDecrypter(key)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize decrypter (%w)", err)
		}
		return decrypter, nil

	case model.SecretManagementTypeGCPKMS:
		return nil, fmt.Errorf("type %q is not implemented yet", sm.Type.String())

	case model.SecretManagementTypeAWSKMS:
		return nil, fmt.Errorf("type %q is not implemented yet", sm.Type.String())

	default:
		return nil, fmt.Errorf("unsupported secret management type: %s", sm.Type.String())
	}
}
