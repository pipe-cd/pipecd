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

package datastore

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipe/pkg/model"
)

func TestAddImageReference(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	im := model.ImageReference{
		Id:        "id",
		RepoName:  "repo",
		Tag:       "tag",
		ProjectId: "projectId",
		CreatedAt: 12345,
		UpdatedAt: 12345,
	}

	testcases := []struct {
		name    string
		im      model.ImageReference
		ds      DataStore
		wantErr bool
	}{
		{
			name: "Invalid image metadata",
			im:   model.ImageReference{},
			ds: func() DataStore {
				return NewMockDataStore(ctrl)
			}(),
			wantErr: true,
		},
		{
			name: "OK to create",
			im:   im,
			ds: func() DataStore {
				ds := NewMockDataStore(ctrl)
				ds.EXPECT().
					Create(gomock.Any(), "ImageReference", im.Id, &im).
					Return(nil)
				return ds
			}(),
			wantErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewImageReferenceStore(tc.ds)
			err := s.AddImageReference(context.Background(), tc.im)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}
