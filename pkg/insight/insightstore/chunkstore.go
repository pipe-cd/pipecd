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

package insightstore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/filestore"
	"github.com/pipe-cd/pipe/pkg/insight"
	"github.com/pipe-cd/pipe/pkg/model"
)

type store struct {
	filestore filestore.Store
}

func NewStore(fs filestore.Store) Store {
	return &store{
		filestore: fs,
	}
}

// LoadChunks returns all needed chunks for the specified kind and time range.
func (s *store) LoadChunks(
	ctx context.Context,
	projectID, appID string,
	kind model.InsightMetricsKind,
	step model.InsightStep,
	from time.Time,
	count int,
) (insight.Chunks, error) {
	from = insight.NormalizeTime(from, step)
	paths := insight.DetermineFilePaths(projectID, appID, kind, step, from, count)
	var chunks []insight.Chunk
	for _, p := range paths {
		c, err := s.getChunk(ctx, p, kind)
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, c)
	}

	return chunks, nil
}

// PutChunk creates or updates chunk.
func (s *store) PutChunk(ctx context.Context, chunk insight.Chunk) error {
	data, err := json.Marshal(chunk)
	if err != nil {
		return err
	}
	path := chunk.GetFilePath()
	if path == "" {
		return fmt.Errorf("filepath not found on chunk struct")
	}
	return s.filestore.Put(ctx, path, data)
}

func LoadChunksFromCache(cache cache.Cache, projectID, appID string, kind model.InsightMetricsKind, step model.InsightStep, from time.Time, count int) (insight.Chunks, error) {
	paths := insight.DetermineFilePaths(projectID, appID, kind, step, from, count)
	chunks := make([]insight.Chunk, 0, len(paths))
	for _, p := range paths {
		c, err := cache.Get(p)
		if err != nil {
			return nil, err
		}

		chunk, ok := c.(insight.Chunk)
		if !ok {
			return nil, errors.New("malformed chunk data in cache")
		}
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

func PutChunksToCache(cache cache.Cache, chunks insight.Chunks) error {
	var err error
	for _, c := range chunks {
		// continue process even if an error occurs.
		if e := cache.Put(c.GetFilePath(), c); e != nil {
			err = e
		}
	}
	return err
}

func (s *store) getChunk(ctx context.Context, path string, kind model.InsightMetricsKind) (insight.Chunk, error) {
	var c interface{}
	switch kind {
	case model.InsightMetricsKind_DEPLOYMENT_FREQUENCY:
		c = &insight.DeployFrequencyChunk{}
	case model.InsightMetricsKind_CHANGE_FAILURE_RATE:
		c = &insight.ChangeFailureRateChunk{}
	default:
		return nil, fmt.Errorf("unimpremented insight kind: %s", kind)
	}

	obj, err := s.filestore.Get(ctx, path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(obj.Content, c)
	if err != nil {
		return nil, err
	}
	chunk, err := insight.ToChunk(c)
	if err != nil {
		return nil, err
	}

	chunk.SetFilePath(path)
	return chunk, nil
}
