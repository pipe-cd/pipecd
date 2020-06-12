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

package applicationlivestatestore

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/filestore"
	"github.com/pipe-cd/pipe/pkg/model"
)

type Store interface {
	// GetStateSnapshot get the specified application live state snapshot.
	GetStateSnapshot(ctx context.Context, applicationID string) (*model.ApplicationLiveStateSnapshot, error)
}

type store struct {
	backend *applicationLiveStateFileStore
	cache   *applicationLiveStateCache
	logger  *zap.Logger
}

func NewStore(fs filestore.Store, c cache.Cache, logger *zap.Logger) Store {
	return &store{
		backend: &applicationLiveStateFileStore{
			backend: fs,
		},
		cache: &applicationLiveStateCache{
			backend: c,
		},
		logger: logger.Named("application-live-state-store"),
	}
}

func (s *store) GetStateSnapshot(ctx context.Context, applicationID string) (*model.ApplicationLiveStateSnapshot, error) {
	cacheResp, err := s.cache.Get(applicationID)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		s.logger.Error("failed to get application live state from cache", zap.Error(err))
	}
	if cacheResp != nil {
		return cacheResp, nil
	}

	fileResp, err := s.backend.Get(ctx, applicationID)
	if err != nil {
		s.logger.Error("failed to get application live state to filestore", zap.Error(err))
		return nil, err
	}

	if err := s.cache.Put(applicationID, fileResp); err != nil {
		s.logger.Error("failed to put application live state to cache", zap.Error(err))
	}
	return fileResp, nil
}
