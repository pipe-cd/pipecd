// Copyright 2021 The PipeCD Authors.
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

package latestanalysisstore

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/filestore"
	"github.com/pipe-cd/pipe/pkg/model"
)

type Store interface {
	// GetMostRecentSuccessfulAnalysisMetadata gives back the most recent successful analysis metadata of the specified application.
	GetMostRecentSuccessfulAnalysisMetadata(ctx context.Context, applicationID string) (*model.AnalysisMetadata, error)
	// PutStateSnapshot updates the most recent successful analysis metadata of the specified application.
	PutMostRecentSuccessfulAnalysisMetadata(ctx context.Context, applicationID string, snapshot *model.AnalysisMetadata) error
}

type store struct {
	backend *analysisFileStore
	cache   *analysisCache
	logger  *zap.Logger
}

func NewStore(fs filestore.Store, c cache.Cache, logger *zap.Logger) Store {
	return &store{
		backend: &analysisFileStore{
			backend: fs,
		},
		cache: &analysisCache{
			backend: c,
		},
		logger: logger.Named("latest-analysis-store"),
	}
}

func (s *store) GetMostRecentSuccessfulAnalysisMetadata(ctx context.Context, applicationID string) (*model.AnalysisMetadata, error) {
	cacheResp, err := s.cache.Get(applicationID)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		s.logger.Error("failed to get the most recent successful analysis metadata from cache", zap.Error(err))
	}
	if cacheResp != nil {
		return cacheResp, nil
	}

	fileResp, err := s.backend.Get(ctx, applicationID)
	if err != nil {
		s.logger.Error("failed to get the most recent successful analysis metadata from filestore", zap.Error(err))
		return nil, err
	}

	if err := s.cache.Put(applicationID, fileResp); err != nil {
		s.logger.Error("failed to put the most recent successful analysis metadata to cache", zap.Error(err))
	}
	return fileResp, nil
}

func (s *store) PutMostRecentSuccessfulAnalysisMetadata(ctx context.Context, applicationID string, snapshot *model.AnalysisMetadata) error {
	if err := s.backend.Put(ctx, applicationID, snapshot); err != nil {
		s.logger.Error("failed to put the most recent successful analysis metadata to filestore", zap.Error(err))
		return err
	}

	if err := s.cache.Put(applicationID, snapshot); err != nil {
		s.logger.Error("failed to put the most recent successful analysis metadata to cache", zap.Error(err))
	}
	return nil
}
