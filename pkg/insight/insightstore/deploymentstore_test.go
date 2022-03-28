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
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pipe-cd/pipecd/pkg/filestore/filestoretest"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestGetDailyDeployments(t *testing.T) {
	t.Parallel()

	const (
		testDateUnix    = 1647820800 // 2022/03/21:00:00:00 UTC
		dayUnix         = int64(time.Hour * 24 / time.Second)
		endYearDateUnix = 1640908800 // 2021:12:31:00:00:00 UTC
	)

	type args struct {
		from, to  int64
		projectID string // currently this is not considered in test
	}
	testcases := []struct {
		name         string
		args         args
		storedMeta   []*model.InsightDeploymentChunkMetadata
		storedChunks [][]*model.InsightDeploymentChunk
		expected     []*model.InsightDeployment
		expectedErr  error
	}{
		{
			name: "Single Chunk",
			args: args{
				from: testDateUnix,
				to:   testDateUnix + dayUnix,
			},
			storedMeta: []*model.InsightDeploymentChunkMetadata{{
				Chunks: []*model.InsightDeploymentChunkMetadata_ChunkMeta{
					{
						From: testDateUnix,
						To:   testDateUnix + dayUnix,
						Name: "key1",
						Size: 0,
					},
				},
			}},
			storedChunks: [][]*model.InsightDeploymentChunk{{
				{
					From: testDateUnix,
					To:   testDateUnix,
					Deployments: []*model.InsightDeployment{
						{
							StartedAt: testDateUnix,
						},
					},
				},
			}},
			expected: []*model.InsightDeployment{
				{
					StartedAt: testDateUnix,
				},
			},
			expectedErr: nil,
		},
		{
			name: "Multi Chunk",
			args: args{
				from: testDateUnix - dayUnix,
				to:   testDateUnix + 2*dayUnix,
			},
			storedMeta: []*model.InsightDeploymentChunkMetadata{{
				Chunks: []*model.InsightDeploymentChunkMetadata_ChunkMeta{
					{
						From: testDateUnix - 2*dayUnix,
						To:   testDateUnix,
						Name: "key1",
					},
					{
						From: testDateUnix,
						To:   testDateUnix + dayUnix,
						Name: "key2",
					},
				},
			}},
			storedChunks: [][]*model.InsightDeploymentChunk{{
				{
					From: testDateUnix - 2*dayUnix,
					To:   testDateUnix,
					Deployments: []*model.InsightDeployment{
						{
							StartedAt: testDateUnix - 2*dayUnix + 1,
						},
						{
							StartedAt: testDateUnix - dayUnix + 2,
						},
					},
				},
				{
					From: testDateUnix,
					To:   testDateUnix + dayUnix,
					Deployments: []*model.InsightDeployment{
						{
							StartedAt: testDateUnix + 1,
						},
						{
							StartedAt: testDateUnix + dayUnix,
						},
					},
				},
			}},
			expected: []*model.InsightDeployment{
				{
					StartedAt: testDateUnix - dayUnix + 2,
				},
				{
					StartedAt: testDateUnix + 1,
				},
				{
					StartedAt: testDateUnix + dayUnix,
				},
			},
			expectedErr: nil,
		},
		{
			name: "Multi Chunk Multi Year",
			args: args{
				from: endYearDateUnix - dayUnix*200,
				to:   endYearDateUnix + dayUnix*300,
			},
			storedMeta: []*model.InsightDeploymentChunkMetadata{
				{
					Chunks: []*model.InsightDeploymentChunkMetadata_ChunkMeta{
						{
							From: endYearDateUnix - 300*dayUnix,
							To:   endYearDateUnix - 150*dayUnix,
							Name: "key1",
						},
						{
							From: endYearDateUnix - 149*dayUnix,
							To:   endYearDateUnix,
							Name: "key2",
						},
					},
				},
				{
					Chunks: []*model.InsightDeploymentChunkMetadata_ChunkMeta{
						{
							From: endYearDateUnix + dayUnix,
							To:   endYearDateUnix + 100*dayUnix,
							Name: "key3",
						},
						{
							From: endYearDateUnix + 101*dayUnix,
							To:   endYearDateUnix + 366*dayUnix,
							Name: "key4",
						},
					},
				},
			},
			storedChunks: [][]*model.InsightDeploymentChunk{
				{
					{
						From: endYearDateUnix - 300*dayUnix,
						To:   endYearDateUnix - 150*dayUnix,
						Deployments: []*model.InsightDeployment{
							{
								StartedAt: endYearDateUnix - 200*dayUnix,
							},
							{
								StartedAt: endYearDateUnix - 150*dayUnix,
							},
						},
					},
					{
						From: endYearDateUnix - 149*dayUnix,
						To:   endYearDateUnix,
						Deployments: []*model.InsightDeployment{
							{
								StartedAt: endYearDateUnix - 140*dayUnix,
							},
							{
								StartedAt: endYearDateUnix,
							},
						},
					},
				},
				{
					{
						From: endYearDateUnix + dayUnix,
						To:   endYearDateUnix + 100*dayUnix,
						Deployments: []*model.InsightDeployment{
							{
								StartedAt: endYearDateUnix + dayUnix,
							},
							{
								StartedAt: endYearDateUnix + 100*dayUnix,
							},
						},
					},
					{
						From: endYearDateUnix + 101*dayUnix,
						To:   endYearDateUnix + 366*dayUnix,
						Deployments: []*model.InsightDeployment{
							{
								StartedAt: endYearDateUnix + 100*dayUnix,
							},
							{
								StartedAt: endYearDateUnix + 366*dayUnix,
							},
						},
					},
				},
			},
			expected: []*model.InsightDeployment{
				{
					StartedAt: endYearDateUnix - 200*dayUnix,
				},
				{
					StartedAt: endYearDateUnix - 150*dayUnix,
				},
				{
					StartedAt: endYearDateUnix - 140*dayUnix,
				},
				{
					StartedAt: endYearDateUnix,
				},
				{
					StartedAt: endYearDateUnix + dayUnix,
				},
				{
					StartedAt: endYearDateUnix + 100*dayUnix,
				},
				{
					StartedAt: endYearDateUnix + 100*dayUnix,
				},
			},
			expectedErr: nil,
		},
		{
			name: "Err Long Duration",
			args: args{
				from: testDateUnix - 365*2*dayUnix - 1,
				to:   testDateUnix,
			},
			expectedErr: errLargeDuration,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			fs := filestoretest.NewMockStore(ctrl)
			for i, meta := range tc.storedMeta {
				fs.EXPECT().Get(gomock.Any(), gomock.Any()).Return(proto.Marshal(meta))

				keys := findPathFromMeta(meta, tc.args.from, tc.args.to)
				assert.Equal(t, len(tc.storedChunks[i]), len(keys))
				for _, v := range tc.storedChunks[i] {
					fs.EXPECT().Get(gomock.Any(), gomock.Any()).Return(proto.Marshal(v))
				}
			}

			store := NewStore(fs)
			got, err := store.List(context.Background(), tc.args.projectID, tc.args.from, tc.args.to, *model.InsightDeploymentVersion_V0.Enum())
			assert.Equal(t, tc.expectedErr, err)
			if tc.expectedErr == nil {
				assert.Equal(t, tc.expected, got)
			}
		})
	}
}

func TestPutDeployment(t *testing.T) {
	t.Parallel()
	now := time.Now()

	const (
		testDateUnix = 1647820800 // 2022/03/21:00:00:00 UTC
		dayUnix      = int64(time.Hour * 24 / time.Second)
	)

	type args struct {
		dailyDeployment []*model.InsightDeployment
		projectID       string // currently, this is not considered in test
	}
	type chunks struct {
		storedMeta     *model.InsightDeploymentChunkMetadata
		storedChunk    *model.InsightDeploymentChunk
		willStoreMeta  *model.InsightDeploymentChunkMetadata
		willStoreChunk *model.InsightDeploymentChunk
	}
	testcases := []struct {
		name        string
		args        args
		chunks      chunks
		setup       func(*filestoretest.MockStore, *chunks)
		expectedErr error
	}{
		{
			name: "No chunk exists",
			args: args{
				dailyDeployment: []*model.InsightDeployment{
					{
						StartedAt: testDateUnix,
					},
					{
						StartedAt: testDateUnix + 2,
					},
				},
			},
			chunks: chunks{
				willStoreMeta: &model.InsightDeploymentChunkMetadata{
					CreatedAt: now.Unix(),
					UpdatedAt: now.Unix(),
					Chunks: []*model.InsightDeploymentChunkMetadata_ChunkMeta{
						{
							From:  testDateUnix,
							To:    testDateUnix + 2,
							Name:  determineDeploymentChunkKey(1),
							Size:  28, // binary size of willStoreChunk
							Count: 2,
						},
					},
				},
				willStoreChunk: &model.InsightDeploymentChunk{
					From:    testDateUnix,
					To:      testDateUnix + 2,
					Version: *model.InsightDeploymentVersion_V0.Enum(),

					Deployments: []*model.InsightDeployment{
						{
							StartedAt: testDateUnix,
						},
						{
							StartedAt: testDateUnix + 2,
						},
					},
				},
			},

			setup: func(fs *filestoretest.MockStore, c *chunks) {
				// get meta
				fs.EXPECT().Get(gomock.Any(), gomock.Any()).Return(proto.Marshal(c.storedMeta))

				// store meta
				raw, err := proto.Marshal(c.willStoreMeta)
				assert.NoError(t, err)
				fs.EXPECT().Put(gomock.Any(), gomock.Any(), raw).Return(nil)

				// store chunk
				raw, err = proto.Marshal(c.willStoreChunk)
				assert.NoError(t, err)
				fs.EXPECT().Put(gomock.Any(), gomock.Any(), raw).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Append to current chunk",
			args: args{
				dailyDeployment: []*model.InsightDeployment{
					{
						StartedAt: testDateUnix + 1,
					},
					{
						StartedAt: testDateUnix + 2,
					},
				},
			},

			chunks: chunks{
				storedMeta: &model.InsightDeploymentChunkMetadata{
					CreatedAt: now.Unix() - 20,
					UpdatedAt: now.Unix() - 20,
					Chunks: []*model.InsightDeploymentChunkMetadata_ChunkMeta{
						{
							From: testDateUnix - dayUnix,
							To:   testDateUnix,

							Name: determineDeploymentChunkKey(1),
						},
					},
				},
				storedChunk: &model.InsightDeploymentChunk{
					From: testDateUnix - dayUnix,
					To:   testDateUnix,

					Deployments: []*model.InsightDeployment{
						{
							StartedAt: testDateUnix - dayUnix,
						},
						{
							StartedAt: testDateUnix,
						},
					},
				},
				willStoreMeta: &model.InsightDeploymentChunkMetadata{
					CreatedAt: now.Unix() - 20,
					UpdatedAt: now.Unix(),

					Chunks: []*model.InsightDeploymentChunkMetadata_ChunkMeta{
						{
							From: testDateUnix - dayUnix,
							To:   testDateUnix + 2,

							Name:  determineDeploymentChunkKey(1),
							Size:  44, // binary size of willStoreChunk
							Count: 4,
						},
					},
				},
				willStoreChunk: &model.InsightDeploymentChunk{
					From: testDateUnix - dayUnix,
					To:   testDateUnix + 2,

					Deployments: []*model.InsightDeployment{
						{
							StartedAt: testDateUnix - dayUnix,
						},
						{
							StartedAt: testDateUnix,
						},
						{
							StartedAt: testDateUnix + 1,
						},
						{
							StartedAt: testDateUnix + 2,
						},
					},
				},
			},

			setup: func(fs *filestoretest.MockStore, c *chunks) {
				// get meta
				fs.EXPECT().Get(gomock.Any(), gomock.Any()).Return(proto.Marshal(c.storedMeta))

				// get chunk
				fs.EXPECT().Get(gomock.Any(), gomock.Any()).Return(proto.Marshal(c.storedChunk))

				// store meta
				raw, err := proto.Marshal(c.willStoreMeta)
				assert.NoError(t, err)
				fs.EXPECT().Put(gomock.Any(), gomock.Any(), raw).Return(nil)

				// store chunk
				raw, err = proto.Marshal(c.willStoreChunk)
				assert.NoError(t, err)
				fs.EXPECT().Put(gomock.Any(), gomock.Any(), raw).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Size over and create new chunk",
			args: args{
				dailyDeployment: []*model.InsightDeployment{
					{
						StartedAt: testDateUnix + 1,
					},
					{
						StartedAt: testDateUnix + 2,
					},
				},
			},

			chunks: chunks{
				storedMeta: &model.InsightDeploymentChunkMetadata{
					CreatedAt: now.Unix() - 20,
					UpdatedAt: now.Unix() - 20,
					Chunks: []*model.InsightDeploymentChunkMetadata_ChunkMeta{
						{
							From: testDateUnix - dayUnix,
							To:   testDateUnix,

							Name: determineDeploymentChunkKey(1),
							Size: maxChunkByteSize - 1,
						},
					},
				},
				storedChunk: &model.InsightDeploymentChunk{
					From: testDateUnix - dayUnix,
					To:   testDateUnix,

					Deployments: []*model.InsightDeployment{
						{
							StartedAt: testDateUnix - dayUnix - 2,
						},
						{
							StartedAt: testDateUnix - dayUnix - 1,
						},
					},
				},
				willStoreMeta: &model.InsightDeploymentChunkMetadata{
					CreatedAt: now.Unix() - 20,
					UpdatedAt: now.Unix(),
					Chunks: []*model.InsightDeploymentChunkMetadata_ChunkMeta{
						{
							From: testDateUnix - dayUnix,
							To:   testDateUnix,

							Name: determineDeploymentChunkKey(1),
							Size: maxChunkByteSize - 1,
						},
						{
							From: testDateUnix + 1,
							To:   testDateUnix + 2,

							Name:  determineDeploymentChunkKey(2),
							Size:  28, // binary size of willStoreChunk
							Count: 2,
						},
					},
				},
				willStoreChunk: &model.InsightDeploymentChunk{
					From: testDateUnix + 1,
					To:   testDateUnix + 2,

					Deployments: []*model.InsightDeployment{
						{
							StartedAt: testDateUnix + 1,
						},
						{
							StartedAt: testDateUnix + 2,
						},
					},
				},
			},

			setup: func(fs *filestoretest.MockStore, c *chunks) {
				// get meta
				fs.EXPECT().Get(gomock.Any(), gomock.Any()).Return(proto.Marshal(c.storedMeta))

				// get chunk
				fs.EXPECT().Get(gomock.Any(), gomock.Any()).Return(proto.Marshal(c.storedChunk))

				// store meta
				raw, err := proto.Marshal(c.willStoreMeta)
				assert.NoError(t, err)
				fs.EXPECT().Put(gomock.Any(), gomock.Any(), raw).Return(nil)

				// store chunk
				raw, err = proto.Marshal(c.willStoreChunk)
				assert.NoError(t, err)
				fs.EXPECT().Put(gomock.Any(), gomock.Any(), raw).Return(nil)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			fs := filestoretest.NewMockStore(ctrl)

			tc.setup(fs, &tc.chunks)

			store := NewStore(fs).(*store)
			err := store.putDeployments(context.Background(), tc.args.projectID, tc.args.dailyDeployment, *model.InsightDeploymentVersion_V0.Enum(), now)

			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestExtractDailyDeploymentFromChunk(t *testing.T) {
	t.Parallel()

	const (
		testDateUnix = 1647820800 // 2022/03/21:00:00:00 UTC
		dayUnix      = int64(time.Hour * 24 / time.Second)
	)

	testcases := []struct {
		name                     string
		chunk                    *model.InsightDeploymentChunk
		from, to                 int64
		expectedDailyDeployments []*model.InsightDeployment
	}{
		{
			name: "SingleDeployments-SingleDeployments",
			chunk: &model.InsightDeploymentChunk{
				From: testDateUnix, To: testDateUnix + 2,
				Deployments: []*model.InsightDeployment{
					{StartedAt: testDateUnix + 1},
					{StartedAt: testDateUnix + 2},
				},
			},
			from: testDateUnix, to: testDateUnix + 2,
			expectedDailyDeployments: []*model.InsightDeployment{
				{StartedAt: testDateUnix + 1},
				{StartedAt: testDateUnix + 2},
			},
		},
		{
			name: "MultipleDeployments-SingleDeployments",
			chunk: &model.InsightDeploymentChunk{
				From: testDateUnix - dayUnix, To: testDateUnix + dayUnix + 2,
				Deployments: []*model.InsightDeployment{
					{StartedAt: testDateUnix - dayUnix + 1},
					{StartedAt: testDateUnix - dayUnix + 2},
					{StartedAt: testDateUnix + 1},
					{StartedAt: testDateUnix + 2},
					{StartedAt: testDateUnix + dayUnix + 1},
					{StartedAt: testDateUnix + dayUnix + 2},
				},
			},
			from: testDateUnix, to: testDateUnix + dayUnix,
			expectedDailyDeployments: []*model.InsightDeployment{
				{StartedAt: testDateUnix + 1},
				{StartedAt: testDateUnix + 2},
			},
		},
		{
			name: "MultipleDeployments-MultipleDeployments",
			chunk: &model.InsightDeploymentChunk{
				From: testDateUnix - dayUnix, To: testDateUnix + dayUnix + 2,
				Deployments: []*model.InsightDeployment{
					{StartedAt: testDateUnix - dayUnix + 1},
					{StartedAt: testDateUnix - dayUnix + 2},
					{StartedAt: testDateUnix + 1},
					{StartedAt: testDateUnix + 2},
					{StartedAt: testDateUnix + dayUnix + 1},
					{StartedAt: testDateUnix + dayUnix + 2},
				},
			},
			from: testDateUnix - dayUnix, to: testDateUnix + 2*dayUnix,
			expectedDailyDeployments: []*model.InsightDeployment{
				{StartedAt: testDateUnix - dayUnix + 1},
				{StartedAt: testDateUnix - dayUnix + 2},
				{StartedAt: testDateUnix + 1},
				{StartedAt: testDateUnix + 2},
				{StartedAt: testDateUnix + dayUnix + 1},
				{StartedAt: testDateUnix + dayUnix + 2},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := extractDeploymentsFromChunk(tc.chunk, tc.from, tc.to)
			assert.Equal(t, tc.expectedDailyDeployments, got)
		})
	}
}

func TestOverlap(t *testing.T) {
	t.Parallel()

	const (
		testDateUnix = 1647820800 // 2022/03/21:00:00:00 UTC
		dayUnix      = int64(time.Hour * 24 / time.Second)
	)

	testcases := []struct {
		name             string
		lhs_from, lhs_to int64
		rhs_from, rhs_to int64
		expected         bool
	}{
		{
			name:     "No overlap 1",
			lhs_from: testDateUnix - dayUnix,
			lhs_to:   testDateUnix,
			rhs_from: testDateUnix,
			rhs_to:   testDateUnix + dayUnix,
			expected: false,
		},
		{
			name:     "No overlap 2",
			lhs_from: testDateUnix,
			lhs_to:   testDateUnix + dayUnix,
			rhs_from: testDateUnix - dayUnix,
			rhs_to:   testDateUnix,
			expected: false,
		},
		{
			name:     "Overlap same day",
			lhs_from: testDateUnix,
			lhs_to:   testDateUnix + dayUnix,
			rhs_from: testDateUnix,
			rhs_to:   testDateUnix + dayUnix,
			expected: true,
		},
		{
			name:     "Overlap",
			lhs_from: testDateUnix - dayUnix,
			lhs_to:   testDateUnix + dayUnix,
			rhs_from: testDateUnix,
			rhs_to:   testDateUnix + dayUnix,
			expected: true,
		},
		{
			name:     "Overlap contain",
			lhs_from: testDateUnix - dayUnix,
			lhs_to:   testDateUnix + dayUnix,
			rhs_from: testDateUnix,
			rhs_to:   testDateUnix,
			expected: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := overlap(tc.lhs_from, tc.lhs_to, tc.rhs_from, tc.rhs_to)
			assert.Equal(t, tc.expected, got)
		})
	}
}
