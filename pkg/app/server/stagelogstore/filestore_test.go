// Copyright 2023 The PipeCD Authors.
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
	"io"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/filestore"
	"github.com/pipe-cd/pipecd/pkg/filestore/filestoretest"
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
