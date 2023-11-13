// Copyright 2023 The PipeCD Authors.
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
	"encoding/base64"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/app/server/commandstore"
	"github.com/pipe-cd/pipecd/pkg/app/server/stagelogstore"
	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/crypto"
	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/filestore"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type applicationGetter interface {
	Get(ctx context.Context, id string) (*model.Application, error)
}

type deploymentGetter interface {
	Get(ctx context.Context, id string) (*model.Deployment, error)
}

type pipedGetter interface {
	Get(ctx context.Context, id string) (*model.Piped, error)
}

func getPiped(ctx context.Context, store pipedGetter, id string, logger *zap.Logger) (*model.Piped, error) {
	piped, err := store.Get(ctx, id)
	if errors.Is(err, datastore.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "Piped is not found")
	}
	if err != nil {
		logger.Error("failed to get piped", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to get piped")
	}

	return piped, nil
}

func getApplication(ctx context.Context, store applicationGetter, id string, logger *zap.Logger) (*model.Application, error) {
	app, err := store.Get(ctx, id)
	if errors.Is(err, datastore.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "Application is not found")
	}
	if err != nil {
		logger.Error("failed to get application", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to get application")
	}

	return app, nil
}

func getDeployment(ctx context.Context, store deploymentGetter, id string, logger *zap.Logger) (*model.Deployment, error) {
	deployment, err := store.Get(ctx, id)
	if errors.Is(err, datastore.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "Deployment is not found")
	}
	if err != nil {
		logger.Error("failed to get deployment", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to get deployment")
	}

	return deployment, nil
}

func getCommand(ctx context.Context, store commandstore.Store, id string, logger *zap.Logger) (*model.Command, error) {
	cmd, err := store.GetCommand(ctx, id)
	if errors.Is(err, datastore.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "Command is not found")
	}
	if err != nil {
		logger.Error("failed to get command", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to get command")
	}

	return cmd, nil
}

func addCommand(ctx context.Context, store commandstore.Store, cmd *model.Command, logger *zap.Logger) error {
	if err := store.AddCommand(ctx, cmd); err != nil {
		logger.Error("failed to create command", zap.Error(err))
		return gRPCStoreError(err, "create command")
	}
	return nil
}

// makeGitPath returns an ApplicationGitPath by adding Repository info and GitPath URL to given args.
func makeGitPath(repoID, path, cfgFilename string, piped *model.Piped, logger *zap.Logger) (*model.ApplicationGitPath, error) {
	var repo *model.ApplicationGitRepository
	for _, r := range piped.Repositories {
		if r.Id == repoID {
			repo = r
			break
		}
	}
	if repo == nil {
		return nil, status.Error(codes.NotFound, "The requested repository is not found in the Piped configuration")
	}

	u, err := git.MakeDirURL(repo.Remote, path, repo.Branch)
	if err != nil {
		logger.Error("failed to make GitPath URL", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to make GitPath URL")
	}

	return &model.ApplicationGitPath{
		Repo:           repo,
		Path:           path,
		ConfigFilename: cfgFilename,
		Url:            u,
	}, nil
}

func encrypt(plaintext string, key []byte, base64Encoding bool, logger *zap.Logger) (string, error) {
	if base64Encoding {
		plaintext = base64.StdEncoding.EncodeToString([]byte(plaintext))
	}
	encrypter, err := crypto.NewHybridEncrypter(key)
	if err != nil {
		logger.Error("failed to initialize the crypter", zap.Error(err))
		return "", status.Error(codes.InvalidArgument, "Invalid public key")
	}
	ciphertext, err := encrypter.Encrypt(plaintext)
	if err != nil {
		logger.Error("failed to encrypt the secret", zap.Error(err))
		return "", status.Error(codes.FailedPrecondition, "Failed to encrypt the secret")
	}
	return ciphertext, nil
}

func getEncriptionKey(se *model.Piped_SecretEncryption) ([]byte, error) {
	if se == nil {
		return nil, status.Error(codes.FailedPrecondition, "The piped does not contain a public key")
	}
	switch model.SecretManagementType(se.Type) {
	case model.SecretManagementTypeKeyPair:
		if se.PublicKey == "" {
			return nil, status.Error(codes.FailedPrecondition, "The piped does not contain a public key")
		}
		return []byte(se.PublicKey), nil
	default:
		return nil, status.Error(codes.FailedPrecondition, "The piped does not contain a valid encryption type")
	}
}

func gRPCStoreError(err error, msg string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, datastore.ErrNotFound) || errors.Is(err, filestore.ErrNotFound) || errors.Is(err, stagelogstore.ErrNotFound) {
		return status.Error(codes.NotFound, fmt.Sprintf("Entity was not found to %s", msg))
	}
	if errors.Is(err, datastore.ErrInvalidArgument) {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("Invalid argument to %s", msg))
	}
	if errors.Is(err, datastore.ErrAlreadyExists) {
		return status.Error(codes.AlreadyExists, fmt.Sprintf("Entity already exists to %s", msg))
	}
	if errors.Is(err, datastore.ErrUserDefined) {
		return status.Error(codes.FailedPrecondition, err.Error())
	}
	return status.Error(codes.Internal, fmt.Sprintf("Failed to %s", msg))
}

func getPipedStatus(cs cache.Cache, id string) (model.Piped_ConnectionStatus, error) {
	pipedStatus, err := cs.Get(id)
	if errors.Is(err, cache.ErrNotFound) {
		return model.Piped_OFFLINE, nil
	}
	if err != nil {
		return model.Piped_UNKNOWN, err
	}

	ps := model.PipedStat{}
	if err = model.UnmarshalPipedStat(pipedStatus, &ps); err != nil {
		return model.Piped_UNKNOWN, err
	}
	if ps.IsStaled(model.PipedStatsRetention) {
		return model.Piped_OFFLINE, nil
	}
	return model.Piped_ONLINE, nil
}
