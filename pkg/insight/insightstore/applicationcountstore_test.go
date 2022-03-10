// Copyright 2021 The PipeCD Authors.
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

	"github.com/pipe-cd/pipecd/pkg/model"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/filestore"
	"github.com/pipe-cd/pipecd/pkg/filestore/filestoretest"
	"github.com/pipe-cd/pipecd/pkg/insight"
)

func TestLoadApplicationCounts(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fs := filestoretest.NewMockStore(ctrl)
	s := &store{filestore: fs}

	tests := []struct {
		name      string
		projectID string
		content   string
		readerErr error

		expectedCounts *insight.ApplicationCounts
		expectedErr    error
	}{
		{
			name:        "not found in filestore",
			projectID:   "pid1",
			content:     "",
			readerErr:   filestore.ErrNotFound,
			expectedErr: filestore.ErrNotFound,
		},
		{
			name:      "successfully loaded from filestore",
			projectID: "pid1",
			content: `{
				"updated_at": 1609459200,
				"counts": [
					{
						"labels": {
							"key1": "value1",
							"key2": "value2"
						},
						"count": 2
					},
					{
						"labels": {
							"key3": "value3"
						},
						"count": 1
					}
				]
			}`,
			expectedCounts: &insight.ApplicationCounts{
				Counts: []model.InsightApplicationCount{
					{
						Labels: map[string]string{
							"key1": "value1",
							"key2": "value2",
						},
						Count: 2,
					},
					{
						Labels: map[string]string{
							"key3": "value3",
						},
						Count: 1,
					},
				},
				UpdatedAt: 1609459200,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := determineFilePath(tc.projectID)
			obj := []byte(tc.content)

			fs.EXPECT().Get(context.TODO(), path).Return(obj, tc.readerErr)

			counts, err := s.LoadApplicationCounts(context.TODO(), tc.projectID)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedCounts, counts)
		})
	}
}
