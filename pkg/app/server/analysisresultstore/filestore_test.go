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

package analysisresultstore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		name          string
		applicationID string
		content       string
		readerErr     error

		expected    *model.AnalysisResult
		expectedErr error
	}{
		{
			name:          "file not found in filestore",
			applicationID: "application-id",
			content:       "",
			readerErr:     filestore.ErrNotFound,
			expectedErr:   filestore.ErrNotFound,
		},
		{
			name:          "found file in filestore",
			applicationID: "application-id",

			content: `{
				"start_time": 1590000000
			}`,

			expected: &model.AnalysisResult{
				StartTime: 1590000000,
			},
		},
	}

	fs := analysisFileStore{
		backend: store,
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			path := buildPath(tc.applicationID)
			content := []byte(tc.content)

			store.EXPECT().Get(context.TODO(), path).Return(content, tc.readerErr)
			metadata, err := fs.Get(context.TODO(), tc.applicationID)
			require.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expected, metadata)
		})
	}
}
