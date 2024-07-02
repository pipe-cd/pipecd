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
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/crypto"
	"github.com/pipe-cd/pipecd/pkg/model"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type PluginAPI struct {
	service.PluginServiceServer

	cfg    *config.PipedSpec
	Logger *zap.Logger
}

// Register registers all handling of this service into the specified gRPC server.
func (a *PluginAPI) Register(server *grpc.Server) {
	service.RegisterPluginServiceServer(server, a)
}

func NewPluginAPI(cfg *config.PipedSpec, logger *zap.Logger) *PluginAPI {
	return &PluginAPI{
		cfg:    cfg,
		Logger: logger.Named("plugin-api"),
	}
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
