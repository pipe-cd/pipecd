// Copyright 2022 The PipeCD Authors.
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

	"google.golang.org/protobuf/proto"

	"github.com/pipe-cd/pipecd/pkg/filestore"
	"github.com/pipe-cd/pipecd/pkg/insight"
	"github.com/pipe-cd/pipecd/pkg/model"
)

var (
	errInvalidArg    = errors.New("invalid arg")
	errLargeDuration = errors.New("too large duration")
	errExceedMaxSize = errors.New("exceed max file size")
	errInconsistent  = errors.New("no consistency between meta and chunk")
)

const (
	maxChunkByteSize = 1 * 1024 * 1024 // 1MB
	metaFileName     = "meta.proto.bin"

	maxDuration time.Duration = 24 * 365 * 2 * time.Hour // 2 year
)

type DeploymentStore interface {
	// List returns slice of Deployment sorted by startedAt ASC.
	List(ctx context.Context, projectID string, from, to int64, minimumVersion model.InsightDeploymentVersion) ([]*model.InsightDeployment, error)

	Put(ctx context.Context, projectID string, deployments []*model.InsightDeployment, version model.InsightDeploymentVersion) error
}

// List returns slice of Deployment sorted by startedAt ASC.
func (s *store) List(ctx context.Context, projectID string, from, to int64, minimumVersion model.InsightDeploymentVersion) ([]*model.InsightDeployment, error) {
	fromTime := time.Unix(from, 0)
	fromYear := fromTime.Year()
	toTime := time.Unix(to, 0)
	toYear := toTime.Year()

	sub := toTime.Sub(fromTime)
	if sub < 0 {
		return nil, errInvalidArg
	}
	if sub > maxDuration {
		return nil, errLargeDuration
	}

	var result []*model.InsightDeployment
	for year := fromYear; year <= toYear; year++ {
		dirPath := determineDeploymentDirPath(year, projectID)
		meta := model.InsightDeploymentChunkMetadata{}
		err := s.loadDataFromFilestore(ctx, fmt.Sprintf("%s/%s", dirPath, metaFileName), &meta)
		if errors.Is(err, filestore.ErrNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}

		keys := findPathFromMeta(&meta, from, to)

		for _, key := range keys {
			data := model.InsightDeploymentChunk{}
			err := s.loadDataFromFilestore(ctx, fmt.Sprintf("%s/%s", dirPath, key), &data)
			if err != nil {
				return nil, err
			}

			// Maybe we should use metadata instead of loading chunkdata
			if data.Version >= minimumVersion {
				deployments := extractDeploymentsFromChunk(&data, from, to)
				result = append(result, deployments...)
			}
		}
	}

	return result, nil
}

// deployments must be sorted by startedAt,
func (s *store) Put(ctx context.Context, projectID string, deployments []*model.InsightDeployment, version model.InsightDeploymentVersion) error {
	dailyDeployments := insight.GroupDeploymentsByDaily(deployments)

	for _, daily := range dailyDeployments {
		err := s.putDeployments(ctx, projectID, daily, version, time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}

// deployments must be sorted by startedAt, and year must be same.
// if deployments is too large, this function cannot store deployments efficiently.
func (s *store) putDeployments(ctx context.Context, projectID string, deployments []*model.InsightDeployment, version model.InsightDeploymentVersion, updatedAt time.Time) error {
	if len(deployments) == 0 {
		return nil
	}

	var year = time.Unix(deployments[0].StartedAt, 0).Year()

	// Load chunk
	dirPath := determineDeploymentDirPath(year, projectID)
	meta := model.InsightDeploymentChunkMetadata{}

	metaPath := fmt.Sprintf("%s/%s", dirPath, metaFileName)
	err := s.loadDataFromFilestore(ctx, metaPath, &meta)
	if err != nil && !errors.Is(err, filestore.ErrNotFound) {
		return err
	}

	// create chunk and meta
	if len(meta.Chunks) == 0 {
		chunkKey := determineDeploymentChunkKey(1)
		newMeta, newChunk, err := createNewChunkAndMeta(deployments, version, chunkKey, updatedAt)
		if err != nil {
			return err
		}

		_, err = s.saveDataIntoFilestore(ctx, metaPath, newMeta)
		if err != nil {
			return err
		}

		_, err = s.saveDataIntoFilestore(ctx, fmt.Sprintf("%s/%s", dirPath, chunkKey), newChunk)
		if err != nil {
			return err
		}
		return nil
	}

	firstDeployment := deployments[0]

	latestChunkMeta := meta.Chunks[len(meta.Chunks)-1]
	if firstDeployment.StartedAt < latestChunkMeta.To {
		return errors.New("cannot overwrite past deployment")
	}

	// TODO refine this size check
	size, err := messageSize(&model.InsightDeploymentChunk{Deployments: deployments})
	if err != nil {
		return err
	}

	// append to current chunk
	if latestChunkMeta.Size+size < maxChunkByteSize {
		// Update chunk
		chunkData := model.InsightDeploymentChunk{}
		chunkPath := fmt.Sprintf("%s/%s", dirPath, latestChunkMeta.Name)
		err := s.loadDataFromFilestore(ctx, chunkPath, &chunkData)
		if err != nil {
			return err
		}

		if firstDeployment.StartedAt < chunkData.To {
			return errInconsistent
		}

		// if version is not equal, then skip this branch and create new chunk.
		if chunkData.Version == version {
			err := appendChunkAndUpdateMeta(&meta, &chunkData, deployments, updatedAt)
			if err != nil {
				return err
			}

			// Store metadata first
			_, err = s.saveDataIntoFilestore(ctx, metaPath, &meta)
			if err != nil {
				return err
			}

			_, err = s.saveDataIntoFilestore(ctx, chunkPath, &chunkData)
			if err != nil {
				return err
			}
			return err
		}
	}

	// Create new meta and chunk
	newChunkKey := determineDeploymentChunkKey(len(meta.Chunks) + 1)
	chunk, err := createNewChunkAndUpdateMeta(&meta, deployments, updatedAt)
	if err != nil {
		return err
	}

	// Store meta first
	_, err = s.saveDataIntoFilestore(ctx, metaPath, &meta)
	if err != nil {
		return err
	}

	_, err = s.saveDataIntoFilestore(ctx, fmt.Sprintf("%s/%s", dirPath, newChunkKey), chunk)
	return err
}

// deployments must not be empty
func createNewChunk(deployments []*model.InsightDeployment) (*model.InsightDeploymentChunk, int64, error) {
	from, to := deployments[0].StartedAt, deployments[len(deployments)-1].StartedAt
	// Create new meta and chunk
	chunkData := model.InsightDeploymentChunk{
		From:        from,
		To:          to,
		Deployments: deployments,
	}
	size, err := messageSize(&chunkData)
	if err != nil {
		return nil, 0, err
	}
	if size > maxChunkByteSize {
		return nil, 0, errExceedMaxSize
	}

	return &chunkData, size, nil
}

func createNewChunkAndMeta(deployments []*model.InsightDeployment, version model.InsightDeploymentVersion, key string, updatedAt time.Time) (*model.InsightDeploymentChunkMetadata, *model.InsightDeploymentChunk, error) {
	if len(deployments) == 0 {
		return nil, nil, errInvalidArg
	}

	from, to := deployments[0].StartedAt, deployments[len(deployments)-1].StartedAt
	// Create new meta and chunk
	chunkData := &model.InsightDeploymentChunk{
		From:        from,
		To:          to,
		Version:     version,
		Deployments: deployments,
	}
	size, err := messageSize(chunkData)
	if err != nil {
		return nil, nil, err
	}
	if size > maxChunkByteSize {
		return nil, nil, errExceedMaxSize
	}

	// Create meta
	newChunkMetaData := model.InsightDeploymentChunkMetadata_ChunkMeta{
		From:  from,
		To:    to,
		Name:  key,
		Size:  size,
		Count: int64(len(deployments)),
	}

	meta := &model.InsightDeploymentChunkMetadata{
		Chunks:    []*model.InsightDeploymentChunkMetadata_ChunkMeta{&newChunkMetaData},
		CreatedAt: updatedAt.Unix(),
		UpdatedAt: updatedAt.Unix(),
	}

	return meta, chunkData, nil

}

// deployments must not be empty
func createNewChunkAndUpdateMeta(curMeta *model.InsightDeploymentChunkMetadata, deployments []*model.InsightDeployment, updatedAt time.Time) (*model.InsightDeploymentChunk, error) {
	newChunk, size, err := createNewChunk(deployments)
	if err != nil {
		return nil, err
	}
	if size > maxChunkByteSize {
		return nil, errExceedMaxSize
	}

	// Create new meta and chunk
	newChunkKey := determineDeploymentChunkKey(len(curMeta.Chunks) + 1)

	// Create meta
	from, to := deployments[0].StartedAt, deployments[len(deployments)-1].StartedAt
	newChunkMetaData := model.InsightDeploymentChunkMetadata_ChunkMeta{
		From:  from,
		To:    to,
		Name:  newChunkKey,
		Size:  size,
		Count: int64(len(deployments)),
	}

	curMeta.Chunks = append(curMeta.Chunks, &newChunkMetaData)
	curMeta.UpdatedAt = updatedAt.Unix()

	return newChunk, nil
}

// deployments must not be empty
func appendChunkAndUpdateMeta(meta *model.InsightDeploymentChunkMetadata, curChunk *model.InsightDeploymentChunk, deployments []*model.InsightDeployment, updatedAt time.Time) error {
	firstDeployment := deployments[0]
	lastDeployment := deployments[len(deployments)-1]
	if firstDeployment.StartedAt < curChunk.To {
		return errInvalidArg
	}

	curChunk.Deployments = append(curChunk.Deployments, deployments...)
	curChunk.To = lastDeployment.StartedAt

	size, err := messageSize(curChunk)
	if err != nil {
		return err
	}

	// we cannot know exact size on caller side.
	// so maybe this condition become true.
	// this is safe if we think maxChunkByteSize as soft limit
	// if size > maxChunkByteSize {
	// 	return errExceedMaxSize
	// }

	latestChunkMeta := meta.Chunks[len(meta.Chunks)-1]

	// Update meta
	latestChunkMeta.Size = size
	latestChunkMeta.To = lastDeployment.StartedAt
	latestChunkMeta.Count += int64(len(curChunk.Deployments))
	meta.UpdatedAt = updatedAt.Unix()

	return nil
}

// Load proto message stored in path. If path does not exists, dest does not modified.
func (s *store) loadDataFromFilestore(ctx context.Context, path string, dest proto.Message) error {
	raw, err := s.filestore.Get(ctx, path)
	if err != nil {
		return err
	}

	err = proto.Unmarshal(raw, dest)
	if err != nil {
		return err
	}
	return nil
}

func (s *store) saveDataIntoFilestore(ctx context.Context, path string, data proto.Message) (dataSize int64, err error) {
	raw, err := proto.Marshal(data)
	if err != nil {
		return dataSize, err
	}

	err = s.filestore.Put(ctx, path, raw)
	if err != nil {
		return dataSize, err
	}
	return int64(len(raw)), nil
}

func findPathFromMeta(meta *model.InsightDeploymentChunkMetadata, from, to int64) []string {
	var paths []string
	for _, m := range meta.Chunks {
		if overlap(from, to, m.From, m.To) {
			paths = append(paths, m.Name)
		}
	}
	return paths
}

func messageSize(pbm proto.Message) (int64, error) {
	raw, err := proto.Marshal(pbm)
	if err != nil {
		return 0, err
	}

	return int64(len(raw)), nil
}

func overlap(lhsFrom, lhsTo, rhsFrom, rhsTo int64) bool {
	return (rhsFrom < lhsTo && lhsFrom < rhsTo)
}

func extractDeploymentsFromChunk(chunk *model.InsightDeploymentChunk, from, to int64) []*model.InsightDeployment {
	var result []*model.InsightDeployment
	for _, d := range chunk.Deployments {
		if from <= d.StartedAt && d.StartedAt <= to {
			result = append(result, d)
		}
	}
	return result
}

func determineDeploymentDirPath(year int, projectID string) string {
	return fmt.Sprintf("insights/deployments/%s/%d", projectID, year)
}

func determineDeploymentChunkKey(n int) string {
	return fmt.Sprintf("%04d.proto.bin", n)
}
