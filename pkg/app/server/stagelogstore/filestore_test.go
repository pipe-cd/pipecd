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
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/pipe-cd/pipecd/pkg/filestore"
	"github.com/pipe-cd/pipecd/pkg/filestore/filestoretest"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestFileStoreGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := filestoretest.NewMockStore(ctrl)

	testcases := []struct {
		name         string
		deploymentID string
		stageID      string
		retriedCount int32

		content   string
		readerErr error

		expectedCompleted bool
		expectedRowLength int
		expectedErr       error
	}{
		{
			name:         "file not found in filestore",
			deploymentID: "deployment-id",
			stageID:      "stage-id",
			retriedCount: 0,

			content:   "",
			readerErr: filestore.ErrNotFound,

			expectedErr: filestore.ErrNotFound,
		},
		{
			name:         "incomplete logs",
			deploymentID: "deployment-id",
			stageID:      "stage-id",
			retriedCount: 0,

			content: `
				{"index":1,"log":"Hello 1","severity":0,"created_at":1590499431}
				{"index":2,"log":"Hello 2","severity":0,"created_at":1590499432}`,
			expectedRowLength: 2,
			expectedCompleted: false,
			expectedErr:       nil,
		},
		{
			name:         "incomplete multiple line logs",
			deploymentID: "deployment-id",
			stageID:      "stage-id",
			retriedCount: 0,

			content: `
				{"index":1,"log":"Hello 1\nWorld","severity":0,"created_at":1590499431}
				{"index":2,"log":"Hello 2\nPiped,\nThank you.","severity":0,"created_at":1590499432}`,
			expectedRowLength: 2,
			expectedCompleted: false,
			expectedErr:       nil,
		},
		{
			name:         "complete logs",
			deploymentID: "deployment-id",
			stageID:      "stage-id",
			retriedCount: 0,

			content: `
				{"index":1,"log":"Hello 1","severity":1,"created_at":1590499431}
				{"index":2,"log":"Hello 2","severity":1,"created_at":1590499432}
EOL`,
			expectedRowLength: 2,
			expectedCompleted: true,
			expectedErr:       nil,
		},
	}

	fs := stageLogFileStore{
		filestore: store,
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			path := stageLogPath(tc.deploymentID, tc.stageID, tc.retriedCount)
			reader := io.NopCloser(strings.NewReader(tc.content))
			store.EXPECT().GetReader(context.TODO(), path).Return(reader, tc.readerErr)
			lf, err := fs.Get(context.TODO(), tc.deploymentID, tc.stageID, tc.retriedCount)
			if err != nil {
				if tc.expectedErr == nil {
					assert.NoError(t, err)
					return
				}
				assert.Error(t, err, tc.expectedErr)
				return
			}
			assert.Equal(t, tc.expectedRowLength, len(lf.Blocks))
			assert.Equal(t, tc.expectedCompleted, lf.Completed)
		})
	}
}

func TestFileStorePut(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := filestoretest.NewMockStore(ctrl)

	fs := stageLogFileStore{
		filestore: store,
	}

	lf := &logFragment{
		Blocks: []*model.LogBlock{
			{
				Index:     1,
				Log:       "Hello 1",
				Severity:  model.LogSeverity_INFO,
				CreatedAt: 1590499431,
			},
			{
				Index:     2,
				Log:       "Hello 2\nWorld",
				Severity:  model.LogSeverity_ERROR,
				CreatedAt: 1590499432,
			},
		},
		Completed: true,
	}

	store.EXPECT().
		Put(context.TODO(), "log/deployment-id/stage-id/0.txt", []byte("{\"index\":1,\"log\":\"Hello 1\",\"created_at\":1590499431}\n{\"index\":2,\"log\":\"Hello 2\\nWorld\",\"severity\":2,\"created_at\":1590499432}\nEOL")).
		Return(nil)

	assert.NoError(t, fs.Put(context.TODO(), "deployment-id", "stage-id", 0, lf))
}

func BenchmarkStageLogFileStorePut(b *testing.B) {
	for _, blockCount := range []int{10, 100, 1000} {
		lf := benchmarkLogFragment(blockCount, true)
		b.Run("legacy/"+strconv.Itoa(blockCount), func(b *testing.B) {
			store := stageLogFileStore{filestore: benchmarkFilestore{}}
			ctx := context.Background()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := benchmarkLegacyPut(&store, ctx, "deployment-id", "stage-id", 0, lf); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run("current/"+strconv.Itoa(blockCount), func(b *testing.B) {
			store := stageLogFileStore{filestore: benchmarkFilestore{}}
			ctx := context.Background()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := store.Put(ctx, "deployment-id", "stage-id", 0, lf); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func benchmarkLegacyPut(f *stageLogFileStore, ctx context.Context, deploymentID, stageID string, retriedCount int32, lf *logFragment) error {
	path := stageLogPath(deploymentID, stageID, retriedCount)
	var buf bytes.Buffer
	for _, lb := range lf.Blocks {
		raw, err := json.Marshal(lb)
		if err != nil {
			return err
		}
		buf.Write(raw)
		buf.WriteString("\n")
	}

	if lf.Completed {
		buf.Write(eol)
	}
	return f.filestore.Put(ctx, path, buf.Bytes())
}

func benchmarkLogFragment(blockCount int, completed bool) *logFragment {
	blocks := make([]*model.LogBlock, 0, blockCount)
	for i := 0; i < blockCount; i++ {
		blocks = append(blocks, &model.LogBlock{
			Index:     int64(i),
			Log:       strings.Repeat("benchmark-stage-log-line-", 4) + strconv.Itoa(i),
			Severity:  model.LogSeverity_INFO,
			CreatedAt: 1590499431 + int64(i),
		})
	}
	return &logFragment{
		Blocks:    blocks,
		Completed: completed,
	}
}

type benchmarkFilestore struct{}

func (benchmarkFilestore) Get(context.Context, string) ([]byte, error) {
	panic("not implemented")
}

func (benchmarkFilestore) GetReader(context.Context, string) (io.ReadCloser, error) {
	panic("not implemented")
}

func (benchmarkFilestore) Put(context.Context, string, []byte) error {
	return nil
}

func (benchmarkFilestore) List(context.Context, string) ([]filestore.ObjectAttrs, error) {
	panic("not implemented")
}

func (benchmarkFilestore) Delete(context.Context, string) error {
	panic("not implemented")
}

func (benchmarkFilestore) Close() error {
	return nil
}
