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
	"fmt"
	"time"

	"github.com/pipe-cd/pipe/pkg/filestore"
	"github.com/pipe-cd/pipe/pkg/model"
)

type Store struct {
	filestore filestore.Store
}

func NewStore(fs filestore.Store) Store {
	return Store{filestore: fs}
}

// LoadChunks returns all needed chunks for the specified kind and time range.
func (s *Store) LoadChunks(
	ctx context.Context,
	projectID, appID string,
	kind model.InsightMetricsKind,
	step model.InsightStep,
	from time.Time,
	count int,
) ([]Chunk, error) {
	from = normalizeTime(from, step)
	paths := determineFilePaths(projectID, appID, kind, step, from, count)
	var chunks []Chunk
	for _, p := range paths {
		c, err := s.getChunk(ctx, p, kind)
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, c)
	}

	return chunks, nil
}

// PutChunk create or update chunk.
func (s *Store) PutChunk(ctx context.Context, chunk Chunk) error {
	data, err := json.Marshal(chunk)
	if err != nil {
		return err
	}
	path := chunk.GetFilePath()
	if path == "" {
		return fmt.Errorf("filepath not found on chunk struct")
	}
	return s.filestore.PutObject(ctx, path, data)
}

func (s *Store) getChunk(ctx context.Context, path string, kind model.InsightMetricsKind) (Chunk, error) {
	obj, err := s.filestore.GetObject(ctx, path)
	if err != nil {
		return nil, err
	}

	var c interface{}
	switch kind {
	case model.InsightMetricsKind_DEPLOYMENT_FREQUENCY:
		c = &DeployFrequencyChunk{}
	case model.InsightMetricsKind_CHANGE_FAILURE_RATE:
		c = &ChangeFailureRateChunk{}
	default:
		return nil, fmt.Errorf("unimpremented insight kind: %s", kind)
	}

	err = json.Unmarshal(obj.Content, c)
	if err != nil {
		return nil, err
	}
	chunk, err := toChunk(c)
	if err != nil {
		return nil, err
	}

	chunk.SetFilePath(path)
	return chunk, nil
}

func normalizeTime(from time.Time, step model.InsightStep) time.Time {
	var formattedTime time.Time
	switch step {
	case model.InsightStep_DAILY:
		formattedTime = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	case model.InsightStep_WEEKLY:
		// Sunday in the week of rangeFrom
		sunday := from.AddDate(0, 0, -int(from.Weekday()))
		formattedTime = time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 0, 0, 0, 0, time.UTC)
	case model.InsightStep_MONTHLY:
		formattedTime = time.Date(from.Year(), from.Month(), 1, 0, 0, 0, 0, time.UTC)
	case model.InsightStep_YEARLY:
		formattedTime = time.Date(from.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	}
	return formattedTime
}
