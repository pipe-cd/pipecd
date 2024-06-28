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
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/filestore"
	"github.com/pipe-cd/pipecd/pkg/insight"
)

type fakeFileStore struct {
	storage *sync.Map
}

func newFakeFileStore() fileStore {
	return &fakeFileStore{
		storage: &sync.Map{},
	}
}

func (s *fakeFileStore) Get(ctx context.Context, path string) ([]byte, error) {
	value, ok := s.storage.Load(path)
	if !ok {
		return nil, filestore.ErrNotFound
	}
	return value.([]byte), nil
}

func (s *fakeFileStore) Put(ctx context.Context, path string, content []byte) error {
	s.storage.Store(path, content)
	return nil
}

func TestBlockMetadataStoring(t *testing.T) {
	var (
		fs = newFakeFileStore()
		s  = &store{
			fileStore: fs,
			logger:    zap.NewNop(),
		}
		projectID = "test-project"
		blockID   = "block-2022"
		block     = DeploymentBlockMetadata{
			BlockID: blockID,
			ChunkMetadata: []*DeploymentChunkMetadata{
				&DeploymentChunkMetadata{
					ChunkID: "chunk-0",
					Count:   10,
				},
			},
		}
		ctx = context.TODO()
	)

	_, err := s.loadBlockMetadata(ctx, projectID, blockID)
	require.Equal(t, filestore.ErrNotFound, err)

	err = s.saveBlockMetadata(ctx, projectID, &block)
	require.NoError(t, err)

	got, err := s.loadBlockMetadata(ctx, projectID, blockID)
	require.NoError(t, err)
	assert.Equal(t, block, *got)
}

func TestChunkStoring(t *testing.T) {
	var (
		fs = newFakeFileStore()
		s  = &store{
			fileStore: fs,
			logger:    zap.NewNop(),
		}
		projectID = "test-project"
		blockID   = "block-2022"
		chunkID   = "chunk-10"
		chunk     = DeploymentChunk{
			ChunkID: chunkID,
			Deployments: []*insight.DeploymentData{
				&insight.DeploymentData{
					ID: "deployment-1",
				},
				&insight.DeploymentData{
					ID: "deployment-2",
				},
			},
		}
		ctx = context.TODO()
	)

	_, err := s.loadChunk(ctx, projectID, blockID, chunkID, false)
	require.Equal(t, filestore.ErrNotFound, err)

	err = s.saveChunk(ctx, projectID, blockID, &chunk)
	require.NoError(t, err)

	got, err := s.loadChunk(ctx, projectID, blockID, chunkID, false)
	require.NoError(t, err)
	assert.Equal(t, chunk, *got)
}

func TestPutCompletedDeploymentsToChunk(t *testing.T) {
	var (
		fs = newFakeFileStore()
		s  = &store{
			fileStore: fs,
			logger:    zap.NewNop(),
		}
		projectID = "test-project"
		blockID   = "block-2022"
		chunkID   = "chunk-10"
		ctx       = context.TODO()
	)

	// Add to a non-existing chunk.
	_, err := s.loadChunk(ctx, projectID, blockID, chunkID, false)
	require.Equal(t, filestore.ErrNotFound, err)

	ds := []*insight.DeploymentData{
		&insight.DeploymentData{
			ID: "deployment-1",
		},
		&insight.DeploymentData{
			ID: "deployment-2",
		},
	}
	err = s.putCompletedDeploymentsToChunk(ctx, projectID, blockID, chunkID, ds)
	require.NoError(t, err)

	expected := &DeploymentChunk{
		ChunkID: chunkID,
		Deployments: []*insight.DeploymentData{
			&insight.DeploymentData{
				ID: "deployment-1",
			},
			&insight.DeploymentData{
				ID: "deployment-2",
			},
		},
	}
	got, err := s.loadChunk(ctx, projectID, blockID, chunkID, false)
	require.NoError(t, err)
	assert.Equal(t, expected, got)

	// Add to an existing chunk.
	ds = []*insight.DeploymentData{
		&insight.DeploymentData{
			ID: "deployment-2",
		},
		&insight.DeploymentData{
			ID: "deployment-4",
		},
		&insight.DeploymentData{
			ID: "deployment-5",
		},
	}
	err = s.putCompletedDeploymentsToChunk(ctx, projectID, blockID, chunkID, ds)
	require.NoError(t, err)

	expected = &DeploymentChunk{
		ChunkID: chunkID,
		Deployments: []*insight.DeploymentData{
			&insight.DeploymentData{
				ID: "deployment-1",
			},
			&insight.DeploymentData{
				ID: "deployment-2",
			},
			&insight.DeploymentData{
				ID: "deployment-4",
			},
			&insight.DeploymentData{
				ID: "deployment-5",
			},
		},
	}
	got, err = s.loadChunk(ctx, projectID, blockID, chunkID, false)
	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestPutCompletedDeploymentsBlock(t *testing.T) {
	var (
		fs = newFakeFileStore()
		s  = &store{
			fileStore:     fs,
			chunkMaxCount: 3,
			logger:        zap.NewNop(),
		}
		projectID = "test-project"
		blockID   = "block_2022"
		ctx       = context.TODO()
	)

	// Add to a non-existing block.
	_, err := s.loadBlockMetadata(ctx, projectID, blockID)
	require.Equal(t, filestore.ErrNotFound, err)

	ds := []*insight.DeploymentData{
		&insight.DeploymentData{
			ID:          "deployment-1",
			CompletedAt: 100,
		},
	}
	err = s.putCompletedDeploymentsBlock(ctx, projectID, blockID, ds)
	require.NoError(t, err)

	expectedBlock := &DeploymentBlockMetadata{
		BlockID: blockID,
		ChunkMetadata: []*DeploymentChunkMetadata{
			&DeploymentChunkMetadata{
				ChunkID:      "chunk_0",
				Count:        1,
				MinTimestamp: 100,
				MaxTimestamp: 100,
				Completed:    false,
			},
		},
	}
	gotBlock, err := s.loadBlockMetadata(ctx, projectID, blockID)
	require.NoError(t, err)
	assert.Equal(t, expectedBlock, gotBlock)

	expectedChunks := []*DeploymentChunk{
		&DeploymentChunk{
			ChunkID: "chunk_0",
			Deployments: []*insight.DeploymentData{
				&insight.DeploymentData{
					ID:          "deployment-1",
					CompletedAt: 100,
				},
			},
		},
	}
	for _, expChunk := range expectedChunks {
		got, err := s.loadChunk(ctx, projectID, blockID, expChunk.ChunkID, false)
		require.NoError(t, err)
		assert.Equal(t, expChunk, got)
	}

	// Add to require a new chunk because of size limit.
	ds = []*insight.DeploymentData{
		&insight.DeploymentData{
			ID:          "deployment-2",
			CompletedAt: 200,
		},
		&insight.DeploymentData{
			ID:          "deployment-3",
			CompletedAt: 300,
		},
		&insight.DeploymentData{
			ID:          "deployment-4",
			CompletedAt: 400,
		},
	}
	err = s.putCompletedDeploymentsBlock(ctx, projectID, blockID, ds)
	require.NoError(t, err)

	expectedBlock = &DeploymentBlockMetadata{
		BlockID: blockID,
		ChunkMetadata: []*DeploymentChunkMetadata{
			&DeploymentChunkMetadata{
				ChunkID:      "chunk_0",
				ChunkIndex:   0,
				Count:        3,
				MinTimestamp: 100,
				MaxTimestamp: 300,
				Completed:    true,
			},
			&DeploymentChunkMetadata{
				ChunkID:      "chunk_1",
				ChunkIndex:   1,
				Count:        1,
				MinTimestamp: 400,
				MaxTimestamp: 400,
				Completed:    false,
			},
		},
	}
	gotBlock, err = s.loadBlockMetadata(ctx, projectID, blockID)
	require.NoError(t, err)
	assert.Equal(t, expectedBlock, gotBlock)

	expectedChunks = []*DeploymentChunk{
		&DeploymentChunk{
			ChunkID: "chunk_0",
			Deployments: []*insight.DeploymentData{
				&insight.DeploymentData{
					ID:          "deployment-1",
					CompletedAt: 100,
				},
				&insight.DeploymentData{
					ID:          "deployment-2",
					CompletedAt: 200,
				},
				&insight.DeploymentData{
					ID:          "deployment-3",
					CompletedAt: 300,
				},
			},
		},
		&DeploymentChunk{
			ChunkID: "chunk_1",
			Deployments: []*insight.DeploymentData{
				&insight.DeploymentData{
					ID:          "deployment-4",
					CompletedAt: 400,
				},
			},
		},
	}
	for _, expChunk := range expectedChunks {
		got, err := s.loadChunk(ctx, projectID, blockID, expChunk.ChunkID, false)
		require.NoError(t, err)
		assert.Equal(t, expChunk, got)
	}

	// Add to the limit of chunk to test the edge case.
	ds = []*insight.DeploymentData{
		&insight.DeploymentData{
			ID:          "deployment-5",
			CompletedAt: 500,
		},
		&insight.DeploymentData{
			ID:          "deployment-6",
			CompletedAt: 600,
		},
	}
	err = s.putCompletedDeploymentsBlock(ctx, projectID, blockID, ds)
	require.NoError(t, err)

	expectedBlock = &DeploymentBlockMetadata{
		BlockID: blockID,
		ChunkMetadata: []*DeploymentChunkMetadata{
			&DeploymentChunkMetadata{
				ChunkID:      "chunk_0",
				ChunkIndex:   0,
				Count:        3,
				MinTimestamp: 100,
				MaxTimestamp: 300,
				Completed:    true,
			},
			&DeploymentChunkMetadata{
				ChunkID:      "chunk_1",
				ChunkIndex:   1,
				Count:        3,
				MinTimestamp: 400,
				MaxTimestamp: 600,
				Completed:    true,
			},
		},
	}
	gotBlock, err = s.loadBlockMetadata(ctx, projectID, blockID)
	require.NoError(t, err)
	assert.Equal(t, expectedBlock, gotBlock)

	expectedChunks = []*DeploymentChunk{
		&DeploymentChunk{
			ChunkID: "chunk_0",
			Deployments: []*insight.DeploymentData{
				&insight.DeploymentData{
					ID:          "deployment-1",
					CompletedAt: 100,
				},
				&insight.DeploymentData{
					ID:          "deployment-2",
					CompletedAt: 200,
				},
				&insight.DeploymentData{
					ID:          "deployment-3",
					CompletedAt: 300,
				},
			},
		},
		&DeploymentChunk{
			ChunkID: "chunk_1",
			Deployments: []*insight.DeploymentData{
				&insight.DeploymentData{
					ID:          "deployment-4",
					CompletedAt: 400,
				},
				&insight.DeploymentData{
					ID:          "deployment-5",
					CompletedAt: 500,
				},
				&insight.DeploymentData{
					ID:          "deployment-6",
					CompletedAt: 600,
				},
			},
		},
	}
	for _, expChunk := range expectedChunks {
		got, err := s.loadChunk(ctx, projectID, blockID, expChunk.ChunkID, false)
		require.NoError(t, err)
		assert.Equal(t, expChunk, got)
	}
}

func TestGroupCompletedDeploymentsByBlock(t *testing.T) {
	date2022 := time.Date(2022, 11, 5, 6, 0, 0, 0, time.UTC)
	date2023 := time.Date(2023, 1, 5, 6, 0, 0, 0, time.UTC)

	testcases := []struct {
		name     string
		ds       []*insight.DeploymentData
		expected []*blockData
	}{
		{
			name:     "nil",
			ds:       nil,
			expected: nil,
		},
		{
			name:     "empty",
			ds:       []*insight.DeploymentData{},
			expected: nil,
		},
		{
			name: "one block",
			ds: []*insight.DeploymentData{
				&insight.DeploymentData{
					ID:          "deployment-1",
					CompletedAt: date2022.Unix(),
				},
				&insight.DeploymentData{
					ID:          "deployment-3",
					CompletedAt: date2022.Unix() + 200,
				},
				&insight.DeploymentData{
					ID:          "deployment-2",
					CompletedAt: date2022.Unix() + 100,
				},
			},
			expected: []*blockData{
				&blockData{
					BlockID: "block_2022",
					Deployments: []*insight.DeploymentData{
						&insight.DeploymentData{
							ID:          "deployment-1",
							CompletedAt: date2022.Unix(),
						},
						&insight.DeploymentData{
							ID:          "deployment-2",
							CompletedAt: date2022.Unix() + 100,
						},
						&insight.DeploymentData{
							ID:          "deployment-3",
							CompletedAt: date2022.Unix() + 200,
						},
					},
				},
			},
		},
		{
			name: "multiple blocks",
			ds: []*insight.DeploymentData{
				&insight.DeploymentData{
					ID:          "deployment-1",
					CompletedAt: date2022.Unix(),
				},
				&insight.DeploymentData{
					ID:          "deployment-4",
					CompletedAt: date2023.Unix(),
				},
				&insight.DeploymentData{
					ID:          "deployment-3",
					CompletedAt: date2022.Unix() + 200,
				},
				&insight.DeploymentData{
					ID:          "deployment-2",
					CompletedAt: date2022.Unix() + 100,
				},
			},
			expected: []*blockData{
				&blockData{
					BlockID: "block_2022",
					Deployments: []*insight.DeploymentData{
						&insight.DeploymentData{
							ID:          "deployment-1",
							CompletedAt: date2022.Unix(),
						},
						&insight.DeploymentData{
							ID:          "deployment-2",
							CompletedAt: date2022.Unix() + 100,
						},
						&insight.DeploymentData{
							ID:          "deployment-3",
							CompletedAt: date2022.Unix() + 200,
						},
					},
				},
				&blockData{
					BlockID: "block_2023",
					Deployments: []*insight.DeploymentData{
						&insight.DeploymentData{
							ID:          "deployment-4",
							CompletedAt: date2023.Unix(),
						},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := groupCompletedDeploymentsByBlock(tc.ds)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestOverlap(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name                 string
		firstFrom, firstTo   int64
		secondFrom, secondTo int64
		expected             bool
	}{
		{
			name:       "no overlap",
			firstFrom:  2,
			firstTo:    5,
			secondFrom: 15,
			secondTo:   20,
			expected:   false,
		},
		{
			name:       "overlap the boundary",
			firstFrom:  2,
			firstTo:    5,
			secondFrom: 5,
			secondTo:   20,
			expected:   true,
		},
		{
			name:       "overlap a part",
			firstFrom:  2,
			firstTo:    10,
			secondFrom: 5,
			secondTo:   20,
			expected:   true,
		},
		{
			name:       "overlap fully containing",
			firstFrom:  2,
			firstTo:    10,
			secondFrom: 5,
			secondTo:   8,
			expected:   true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := overlap(tc.firstFrom, tc.firstTo, tc.secondFrom, tc.secondTo)
			assert.Equal(t, tc.expected, got)

			got = overlap(tc.secondFrom, tc.secondTo, tc.firstFrom, tc.firstTo)
			assert.Equal(t, tc.expected, got)
		})
	}
}
