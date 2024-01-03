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

package insightstore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/filestore"
	"github.com/pipe-cd/pipecd/pkg/insight"
)

type DeploymentBlockMetadata struct {
	BlockID       string                     `json:"block_id"`
	ChunkMetadata []*DeploymentChunkMetadata `json:"chunk_metadata"`
}

type DeploymentChunkMetadata struct {
	ChunkID      string `json:"chunk_id"`
	ChunkIndex   int    `json:"chunk_index"`
	MinTimestamp int64  `json:"min_timestamp"`
	MaxTimestamp int64  `json:"max_timestamp"`
	Count        int    `json:"count"`
	Completed    bool   `json:"completed"`
}

type DeploymentChunk struct {
	ChunkID     string                    `json:"chunk_id"`
	Deployments []*insight.DeploymentData `json:"deployments"`
}

func (c *DeploymentChunk) Size() (int64, error) {
	raw, err := json.Marshal(c)
	if err != nil {
		return 0, err
	}

	return int64(len(raw)), nil
}

func (c *DeploymentChunk) FindDeployments(from, to int64) []*insight.DeploymentData {
	out := make([]*insight.DeploymentData, 0, len(c.Deployments))
	for _, d := range c.Deployments {
		if from <= d.CompletedAt && d.CompletedAt <= to {
			out = append(out, d)
		}
	}
	return out
}

func (m *DeploymentBlockMetadata) FindChunks(from, to int64) []DeploymentChunkMetadata {
	var out []DeploymentChunkMetadata
	for _, m := range m.ChunkMetadata {
		if overlap(from, to, m.MinTimestamp, m.MaxTimestamp) {
			out = append(out, *m)
		}
	}
	return out
}

func overlap(firstFrom, firstTo, secondFrom, secondTo int64) bool {
	return !(firstTo < secondFrom || firstFrom > secondTo)
}

func (s *store) ListCompletedDeployments(ctx context.Context, projectID string, from, to int64) ([]*insight.DeploymentData, error) {
	const rangeLimit time.Duration = 2 * 365 * 24 * time.Hour // 2 year

	if from > to {
		return nil, errInvalidArg
	}
	if float64(to-from) > rangeLimit.Seconds() {
		return nil, errLargeDuration
	}

	var (
		fromYear = time.Unix(from, 0).Year()
		toYear   = time.Unix(to, 0).Year()
		out      = make([]*insight.DeploymentData, 0)
	)

	for year := fromYear; year <= toYear; year++ {
		blockID := makeDeploymentBlockID(year)
		blockMetadata, err := s.loadBlockMetadata(ctx, projectID, blockID)
		if errors.Is(err, filestore.ErrNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}

		chunkMDs := blockMetadata.FindChunks(from, to)
		if len(chunkMDs) == 0 {
			continue
		}

		for _, chunkMD := range chunkMDs {
			chunk, err := s.loadChunk(ctx, projectID, blockID, chunkMD.ChunkID, chunkMD.Completed)
			if errors.Is(err, filestore.ErrNotFound) {
				continue
			}
			if err != nil {
				return nil, err
			}

			deployments := chunk.FindDeployments(from, to)
			out = append(out, deployments...)
		}
	}

	return out, nil
}

func (s *store) PutCompletedDeployments(ctx context.Context, projectID string, ds []*insight.DeploymentData) error {
	var lastErr error
	blocks := groupCompletedDeploymentsByBlock(ds)

	for _, b := range blocks {
		if err := s.putCompletedDeploymentsBlock(ctx, projectID, b.BlockID, b.Deployments); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

type blockData struct {
	BlockID     string
	Deployments []*insight.DeploymentData
}

func groupCompletedDeploymentsByBlock(ds []*insight.DeploymentData) []*blockData {
	if len(ds) == 0 {
		return nil
	}

	var (
		out      []*blockData
		curBlock *blockData
	)

	// To be sure, sort the deployment list again in order of CompletedAt.
	sort.Slice(ds, func(i, j int) bool {
		return ds[i].CompletedAt < ds[j].CompletedAt
	})

	for _, d := range ds {
		var (
			timestamp = d.CompletedAt
			blockID   = makeDeploymentBlockID(time.Unix(timestamp, 0).Year())
		)
		if curBlock == nil || curBlock.BlockID != blockID {
			curBlock = &blockData{BlockID: blockID}
			out = append(out, curBlock)
		}
		curBlock.Deployments = append(curBlock.Deployments, d)
	}
	return out
}

func (s *store) putCompletedDeploymentsBlock(ctx context.Context, projectID, blockID string, ds []*insight.DeploymentData) error {
	block, err := s.loadBlockMetadata(ctx, projectID, blockID)
	if err != nil {
		if !errors.Is(err, filestore.ErrNotFound) {
			return err
		}
		block = &DeploymentBlockMetadata{
			BlockID: blockID,
		}
		if err := s.saveBlockMetadata(ctx, projectID, block); err != nil {
			return fmt.Errorf("could not save metadata for new block: %w", err)
		}
		s.logger.Info(fmt.Sprintf("saved a new block metadata file: %s", blockID), zap.String("project", projectID))
	}

	for index, chunkIndex := 0, 0; ; chunkIndex++ {
		// There is no more deployment to add.
		if index >= len(ds) {
			break
		}

		// Append new chunk if needed.
		if chunkIndex >= len(block.ChunkMetadata) {
			block.ChunkMetadata = append(block.ChunkMetadata, &DeploymentChunkMetadata{
				ChunkID:      makeDeploymentChunkID(chunkIndex),
				ChunkIndex:   chunkIndex,
				MinTimestamp: ds[index].CompletedAt,
			})
		}

		chunkMD := block.ChunkMetadata[chunkIndex]
		if chunkMD.Completed {
			continue
		}
		if chunkMD.Count >= s.chunkMaxCount {
			chunkMD.Completed = true
			continue
		}

		var (
			from     = index
			addCount = len(ds) - from
		)
		if chunkMD.Count+addCount > s.chunkMaxCount {
			addCount = s.chunkMaxCount - chunkMD.Count
		}

		index = from + addCount
		adds := ds[from:index]

		if err := s.putCompletedDeploymentsToChunk(ctx, projectID, blockID, chunkMD.ChunkID, adds); err != nil {
			s.logger.Error("could not put deployments to chunk",
				zap.String("project", projectID),
				zap.String("block", blockID),
				zap.String("chunk", chunkMD.ChunkID),
			)
			return fmt.Errorf("could not put deployments to chunk %s: %w", chunkMD.ChunkID, err)
		}

		// Update chunk metadata in block.
		chunkMD.Count += len(adds)
		if chunkMD.Count >= s.chunkMaxCount {
			chunkMD.Completed = true
		}
		if chunkMD.MinTimestamp > ds[from].CompletedAt {
			chunkMD.MinTimestamp = ds[from].CompletedAt
		}
		if chunkMD.MaxTimestamp < ds[index-1].CompletedAt {
			chunkMD.MaxTimestamp = ds[index-1].CompletedAt
		}
	}

	// Save block metadata.
	if err := s.saveBlockMetadata(ctx, projectID, block); err != nil {
		return fmt.Errorf("could not save block metadata: %w", err)
	}

	return nil
}

func (s *store) putCompletedDeploymentsToChunk(ctx context.Context, projectID, blockID, chunkID string, ds []*insight.DeploymentData) error {
	chunk, err := s.loadChunk(ctx, projectID, blockID, chunkID, false)
	if err != nil {
		if !errors.Is(err, filestore.ErrNotFound) {
			return err
		}
		chunk = &DeploymentChunk{
			ChunkID: chunkID,
		}
	}

	var (
		mergedList = make([]*insight.DeploymentData, 0, len(chunk.Deployments)+len(ds))
		dsMap      = make(map[string]struct{}, len(ds))
		duplicates = 0
	)
	for _, d := range chunk.Deployments {
		dsMap[d.ID] = struct{}{}
		mergedList = append(mergedList, d)
	}
	for _, d := range ds {
		if _, ok := dsMap[d.ID]; ok {
			duplicates++
			continue
		}
		mergedList = append(mergedList, d)
	}
	chunk.Deployments = mergedList

	if err := s.saveChunk(ctx, projectID, blockID, chunk); err != nil {
		return fmt.Errorf("could not save chunk: %w", err)
	}

	s.logger.Info(fmt.Sprintf("stored %d deployments (%d duplicates) into chunk file: %s", len(ds), duplicates, chunkID))
	return nil
}

func (s *store) loadBlockMetadata(ctx context.Context, projectID, blockID string) (*DeploymentBlockMetadata, error) {
	path := makeCompletedDeploymentsBlockMetaFilePath(projectID, blockID)
	data, err := s.fileStore.Get(ctx, path)
	if err != nil {
		return nil, err
	}

	var md DeploymentBlockMetadata
	if err := json.Unmarshal(data, &md); err != nil {
		return nil, err
	}

	return &md, nil
}

func (s *store) saveBlockMetadata(ctx context.Context, projectID string, block *DeploymentBlockMetadata) error {
	data, err := json.Marshal(block)
	if err != nil {
		return err
	}

	path := makeCompletedDeploymentsBlockMetaFilePath(projectID, block.BlockID)
	return s.fileStore.Put(ctx, path, data)
}

func (s *store) loadChunk(ctx context.Context, projectID, blockID, chunkID string, useCache bool) (*DeploymentChunk, error) {
	var (
		cacheKey = makeCompletedDeploymentChunkCacheKey(projectID, blockID, chunkID)
		path     = makeCompletedDeploymentsChunkFilePath(projectID, blockID, chunkID)
		data     []byte
	)

	if useCache {
		if cdata, err := s.deploymentChunkCache.Get(cacheKey); err == nil {
			data = cdata.([]byte)
			s.logger.Info("successfully loaded deployment chunk from cache", zap.String("key", cacheKey))
		} else {
			data, err = s.fileStore.Get(ctx, path)
			if err != nil {
				return nil, err
			}
			if err := s.deploymentChunkCache.Put(cacheKey, data); err != nil {
				s.logger.Error("failed to put deployment chunk to cache", zap.Error(err))
			}
		}
	} else {
		var err error
		data, err = s.fileStore.Get(ctx, path)
		if err != nil {
			return nil, err
		}
	}

	var c DeploymentChunk
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func (s *store) saveChunk(ctx context.Context, projectID, blockID string, chunk *DeploymentChunk) error {
	data, err := json.Marshal(chunk)
	if err != nil {
		return err
	}

	path := makeCompletedDeploymentsChunkFilePath(projectID, blockID, chunk.ChunkID)
	return s.fileStore.Put(ctx, path, data)
}
