// Copyright 2020 The PipeCD Authors.
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
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipe/pkg/app/api/commandstore"
	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/git"
	"github.com/pipe-cd/pipe/pkg/model"
)

func getPiped(ctx context.Context, store datastore.PipedStore, id string, logger *zap.Logger) (*model.Piped, error) {
	piped, err := store.GetPiped(ctx, id)
	if errors.Is(err, datastore.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "Piped is not found")
	}
	if err != nil {
		logger.Error("failed to get piped", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to get piped")
	}

	return piped, nil
}

func getApplication(ctx context.Context, store datastore.ApplicationStore, id string, logger *zap.Logger) (*model.Application, error) {
	app, err := store.GetApplication(ctx, id)
	if errors.Is(err, datastore.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "Application is not found")
	}
	if err != nil {
		logger.Error("failed to get application", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to get application")
	}

	return app, nil
}

func addCommand(ctx context.Context, store commandstore.Store, cmd *model.Command, logger *zap.Logger) error {
	if err := store.AddCommand(ctx, cmd); err != nil {
		logger.Error("failed to create command", zap.Error(err))
		return status.Error(codes.Internal, "Failed to create command")
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
		return nil, status.Error(codes.Internal, "The requested repository is not found")
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
