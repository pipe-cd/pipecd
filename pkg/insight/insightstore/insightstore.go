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
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/cache/rediscache"
	"github.com/pipe-cd/pipecd/pkg/insight"
	"github.com/pipe-cd/pipecd/pkg/redis"
)

var (
	errInvalidArg    = errors.New("invalid arg")
	errLargeDuration = errors.New("too large duration")
)

type fileStore interface {
	Get(ctx context.Context, path string) ([]byte, error)
	Put(ctx context.Context, path string, content []byte) error
}

type store struct {
	fileStore            fileStore
	chunkMaxCount        int
	deploymentChunkCache cache.Cache
	logger               *zap.Logger
}

func NewStore(fs fileStore, chunkMaxCount int, rd redis.Redis, logger *zap.Logger) insight.Store {
	return &store{
		fileStore:            fs,
		chunkMaxCount:        chunkMaxCount,
		deploymentChunkCache: rediscache.NewTTLCache(rd, 7*24*time.Hour),
		logger:               logger.Named("insight-store"),
	}
}

// File hierarchy structure inside storage:
//
//	insights
//	├─ {projectId}
//	  ├─ applications
//	     ├─ applications.json
//	  ├─ completed-deployments
//	     ├─ block-2021
//	     ├─ block-2022
//		     ├─ block_meta.json
//		     ├─ chunk_0.json
//		     ├─ chunk_1.json

func makeApplicationsFilePath(projectID string) string {
	return fmt.Sprintf("insights/%s/applications/applications.json", projectID)
}

func makeCompletedDeploymentsBlockPath(projectID, blockID string) string {
	return fmt.Sprintf("insights/%s/completed-deployments/%s", projectID, blockID)
}

func makeCompletedDeploymentsBlockMetaFilePath(projectID, blockID string) string {
	dir := makeCompletedDeploymentsBlockPath(projectID, blockID)
	return fmt.Sprintf("%s/block_meta.json", dir)
}

func makeCompletedDeploymentsChunkFilePath(projectID, blockID, chunkID string) string {
	dir := makeCompletedDeploymentsBlockPath(projectID, blockID)
	return fmt.Sprintf("%s/%s.json", dir, chunkID)
}

func makeDeploymentBlockID(year int) string {
	return fmt.Sprintf("block_%d", year)
}

func makeDeploymentChunkID(index int) string {
	return fmt.Sprintf("chunk_%d", index)
}

func makeCompletedDeploymentChunkCacheKey(projectID, blockID, chunkID string) string {
	return fmt.Sprintf("insights-chunk/%s/%s/%s", projectID, blockID, chunkID)
}
