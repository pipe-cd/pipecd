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

package stagelogstore

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/filestore"
	"github.com/pipe-cd/pipecd/pkg/model"
)

var (
	ErrNotFound         = errors.New("stage log was not found")
	ErrAlreadyCompleted = errors.New("stage log was already completed")
)

type logFragment struct {
	Blocks    []*model.LogBlock
	Completed bool
}

type Store interface {
	// FetchLogs get the specified stage logs which filtered by index.
	FetchLogs(ctx context.Context, deploymentID, stageID string, retriedCount int32, offsetIndex int64) ([]*model.LogBlock, bool, error)
	// AppendLogs appends the stage logs. The stage logs are deduplicated with index value.
	AppendLogs(ctx context.Context, deploymentID, stageID string, retriedCount int32, newBlocks []*model.LogBlock) error
	// AppendLogsFromLastCheckpoint appends the stage logs. The stage logs are deduplicated with index value.
	// If completed is true, flush all the logs to that point and cannot append it after this.
	AppendLogsFromLastCheckpoint(ctx context.Context, deploymentID, stageID string, retriedCount int32, newBlocks []*model.LogBlock, completed bool) error
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

func (s *store) FetchLogs(ctx context.Context, deploymentID, stageID string, retriedCount int32, offsetIndex int64) ([]*model.LogBlock, bool, error) {
	cf, err := s.cache.Get(deploymentID, stageID, retriedCount)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		s.logger.Error("failed to get stage log from cache", zap.Error(err))
	}

	if len(cf.Blocks) > 0 {
		blocks, completed := filterLogBlocks(&cf, offsetIndex)
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

	if err := s.cache.Put(deploymentID, stageID, retriedCount, &ff); err != nil {
		s.logger.Error("failed to put stage log to cache", zap.Error(err))
		return nil, false, err
	}

	blocks, completed := filterLogBlocks(&ff, offsetIndex)
	return blocks, completed, nil
}

func (s *store) AppendLogs(ctx context.Context, deploymentID, stageID string, retriedCount int32, newBlocks []*model.LogBlock) error {
	prev, err := s.cache.Get(deploymentID, stageID, retriedCount)
	if err != nil && err != cache.ErrNotFound {
		s.logger.Error("failed to get stage log from cache", zap.Error(err))
	}
	if prev.Completed {
		return ErrAlreadyCompleted
	}

	lf := logFragment{
		Blocks:    mergeBlocks(prev.Blocks, newBlocks),
		Completed: false,
	}
	if err := s.cache.Put(deploymentID, stageID, retriedCount, &lf); err != nil {
		s.logger.Error("failed to put stage log to cache", zap.Error(err))
	}
	return nil
}

func (s *store) AppendLogsFromLastCheckpoint(ctx context.Context, deploymentID, stageID string, retriedCount int32, newBlocks []*model.LogBlock, completed bool) error {
	prev, err := s.backend.Get(ctx, deploymentID, stageID, retriedCount)
	if err != nil && err != filestore.ErrNotFound {
		return err
	}
	if prev.Completed {
		return ErrAlreadyCompleted
	}

	lf := logFragment{
		Blocks:    mergeBlocks(prev.Blocks, newBlocks),
		Completed: completed,
	}
	if err := s.backend.Put(ctx, deploymentID, stageID, retriedCount, &lf); err != nil {
		s.logger.Error("failed to put stage log to filestore", zap.Error(err))
		return err
	}

	// AppendLogs should update to the cache after updating to the filestore. This order is safe.
	if err := s.cache.Put(deploymentID, stageID, retriedCount, &lf); err != nil {
		s.logger.Error("failed to put stage log to cache", zap.Error(err))
	}
	return nil
}

func mergeBlocks(prevs, news []*model.LogBlock) []*model.LogBlock {
	m := make(map[int64]*model.LogBlock, len(prevs))
	for _, lb := range prevs {
		m[lb.Index] = lb
	}

	merged := prevs
	for _, lb := range news {
		if _, ok := m[lb.Index]; !ok {
			merged = append(merged, lb)
		}
	}
	return merged
}

func filterLogBlocks(lf *logFragment, offsetIndex int64) ([]*model.LogBlock, bool) {
	blocks := make([]*model.LogBlock, 0)
	for i := range lf.Blocks {
		if lf.Blocks[i].Index >= offsetIndex {
			blocks = append(blocks, lf.Blocks[i])
		}
	}
	return blocks, lf.Completed
}
