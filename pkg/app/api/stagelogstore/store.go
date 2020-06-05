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

package stagelogstore

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/kapetaniosci/pipe/pkg/cache"
	"github.com/kapetaniosci/pipe/pkg/filestore"
	"github.com/kapetaniosci/pipe/pkg/model"
)

var (
	ErrNotFound = errors.New("stage log not found")
)

type logFragment struct {
	Blocks    []*model.LogBlock
	Completed bool
}

type Store interface {
	// FetchLogs get the specified stage logs which filtered by timestamp.
	FetchLogs(ctx context.Context, deploymentID, stageID string, retriedCount int32, offsetTimestamp int64) ([]*model.LogBlock, bool, error)
}

type store struct {
	backend *stageLogFileStore
	cache   *stageLogCache
	logger  *zap.Logger
}

func NewStore(fs filestore.Store, c cache.Cache, logger *zap.Logger) Store {
	return &store{
		backend: &stageLogFileStore{
			filestore: fs,
		},
		cache: &stageLogCache{
			cache: c,
		},
		logger: logger.Named("stage-log-store"),
	}
}

func (s *store) FetchLogs(ctx context.Context, deploymentID, stageID string, retriedCount int32, offsetTimestamp int64) ([]*model.LogBlock, bool, error) {
	cf, err := s.cache.Get(deploymentID, stageID, retriedCount)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		s.logger.Error("failed to get stage log from cache", zap.Error(err))
	}

	if cf != nil && len(cf.Blocks) > 0 {
		blocks, completed := filterLogBlocks(cf, offsetTimestamp)
		return blocks, completed, nil
	}

	ff, err := s.backend.Get(ctx, deploymentID, stageID, retriedCount)
	if errors.Is(err, filestore.ErrNotFound) {
		return nil, false, ErrNotFound
	}
	if err != nil {
		s.logger.Error("failed to get stage log from filestore", zap.Error(err))
		return nil, false, err
	}

	if err := s.cache.Put(deploymentID, stageID, retriedCount, ff); err != nil {
		s.logger.Error("failed to put stage log to filestore", zap.Error(err))
		return nil, false, err
	}
	blocks, completed := filterLogBlocks(ff, offsetTimestamp)
	return blocks, completed, nil
}

func filterLogBlocks(lf *logFragment, offsetTimestamp int64) ([]*model.LogBlock, bool) {
	blocks := make([]*model.LogBlock, 0)
	for i := range lf.Blocks {
		if lf.Blocks[i].CreatedAt >= offsetTimestamp {
			blocks = append(blocks, lf.Blocks[i])
		}
	}
	return blocks, lf.Completed
}
