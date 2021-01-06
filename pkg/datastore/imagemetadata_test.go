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
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipe/pkg/model"
)

func TestPutImageMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now()
	im := model.ImageMetadata{
		Id:        "id",
		RepoName:  "repo",
		Tag:       "tag",
		ProjectId: "projectId",
	}
	testcases := []struct {
		name    string
		ds      DataStore
		wantErr bool
	}{
		{
			name: "First time to put",
			ds: func() DataStore {
				ds := NewMockDataStore(ctrl)
				ds.EXPECT().
					Get(gomock.Any(), "ImageMetadata", im.Id, &model.ImageMetadata{}).
					Return(ErrNotFound)

				expectedOne := im
				expectedOne.CreatedAt = now.Unix()
				expectedOne.UpdatedAt = now.Unix()

				ds.EXPECT().
					Put(gomock.Any(), "ImageMetadata", im.Id, &expectedOne).
					Return(nil)

				return ds
			}(),
			wantErr: false,
		},
		{
			name: "Put an existing one",
			ds: func() DataStore {
				ds := NewMockDataStore(ctrl)

				ds.EXPECT().
					Get(gomock.Any(), "ImageMetadata", im.Id, &model.ImageMetadata{}).
					DoAndReturn(func(_ context.Context, _, _ string, entity interface{}) error {
						e, ok := entity.(*model.ImageMetadata)
						if !ok {
							return fmt.Errorf("unexpected type, want ImageMetadata but got %t", entity)
						}
						e.CreatedAt = 12345
						e.UpdatedAt = 12345
						return nil
					})

				expectedOne := im
				expectedOne.CreatedAt = 12345
				expectedOne.UpdatedAt = now.Unix()

				ds.EXPECT().
					Put(gomock.Any(), "ImageMetadata", im.Id, &expectedOne).
					Return(nil)

				return ds
			}(),
			wantErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewImageMetadataStore(tc.ds)
			s.nowFunc = func() time.Time { return now }

			err := s.PutImageMetadata(context.Background(), im)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestGetImageMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		name    string
		id      string
		ds      DataStore
		wantErr bool
	}{
		{
			name: "successfully fetched from datastore",
			id:   "id",
			ds: func() DataStore {
				ds := NewMockDataStore(ctrl)
				ds.EXPECT().
					Get(gomock.Any(), "ImageMetadata", "id", &model.ImageMetadata{}).
					Return(nil)
				return ds
			}(),
			wantErr: false,
		},
		{
			name: "failed to fetch from datastore",
			id:   "id",
			ds: func() DataStore {
				ds := NewMockDataStore(ctrl)
				ds.EXPECT().
					Get(gomock.Any(), "ImageMetadata", "id", &model.ImageMetadata{}).
					Return(errors.New("test error"))
				return ds
			}(),
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewImageMetadataStore(tc.ds)
			_, err := s.GetImageMetadata(context.Background(), tc.id)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}
