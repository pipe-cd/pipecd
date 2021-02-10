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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipe/pkg/filestore"
	"github.com/pipe-cd/pipe/pkg/filestore/filestoretest"
	"github.com/pipe-cd/pipe/pkg/insight"
)

func TestStore_LoadApplicationCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := filestoretest.NewMockStore(ctrl)

	fs := Store{
		filestore: store,
	}

	tests := []struct {
		name        string
		projectID   string
		content     string
		readerErr   error
		want        *insight.ApplicationCount
		expectedErr error
	}{
		{
			name:        "file not found in filestore",
			projectID:   "pid1",
			content:     "",
			readerErr:   filestore.ErrNotFound,
			expectedErr: filestore.ErrNotFound,
		},
		{
			name:      "success",
			projectID: "pid1",
			content: `{
				"accumulated_to": 1609459200,
				"accumulated_from": 1609459100,
				"counts": [
					{
						"label_set": {
							"kind": "CLOUDRUN",
							"status": "deploying"
						},
						"count": 2
					},
					{
						"label_set": {
							"kind": "CLOUDRUN",
							"status": "deleted"
						},
						"count": 1
					}
				]
			}`,
			want: &insight.ApplicationCount{
				AccumulatedTo:   1609459200,
				AccumulatedFrom: 1609459100,
				Counts: []insight.ApplicationCountByLabelSet{
					{
						LabelSet: insight.ApplicationCountLabelSet{
							Kind:   "CLOUDRUN",
							Status: "deploying",
						},
						Count: 2,
					},
					{
						LabelSet: insight.ApplicationCountLabelSet{
							Kind:   "CLOUDRUN",
							Status: "deleted",
						},
						Count: 1,
					},
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := determineFilePath(tc.projectID)
			obj := filestore.Object{
				Content: []byte(tc.content),
			}
			store.EXPECT().GetObject(context.TODO(), path).Return(obj, tc.readerErr)
			ac, err := fs.LoadApplicationCount(context.TODO(), tc.projectID)
			if err != nil {
				if tc.expectedErr == nil {
					assert.NoError(t, err)
					return
				}
				assert.Error(t, err, tc.expectedErr)
				return
			}
			assert.Equal(t, tc.want, ac)
		})
	}
}
