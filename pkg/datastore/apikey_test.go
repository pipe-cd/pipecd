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

package datastore

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestAddAPIKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		name      string
		apiKey    *model.APIKey
		dsFactory func(*model.APIKey) DataStore
		wantErr   bool
	}{
		{
			name:      "Invalid apiKey",
			apiKey:    &model.APIKey{},
			dsFactory: func(d *model.APIKey) DataStore { return nil },
			wantErr:   true,
		},
		{
			name: "Valid apiKey",
			apiKey: &model.APIKey{
				Id:        "id",
				Name:      "name",
				KeyHash:   "keyHash",
				ProjectId: "project-id",
				Role:      model.APIKey_READ_ONLY,
				Creator:   "user",
				CreatedAt: 1,
				UpdatedAt: 1,
			},
			dsFactory: func(d *model.APIKey) DataStore {
				ds := NewMockDataStore(ctrl)
				ds.EXPECT().Create(gomock.Any(), gomock.Any(), d.Id, d)
				return ds
			},
			wantErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewAPIKeyStore(tc.dsFactory(tc.apiKey), TestCommander)
			err := s.Add(context.Background(), tc.apiKey)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestListAPIKeys(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		name    string
		opts    ListOptions
		ds      DataStore
		wantErr error
	}{
		{
			name: "iterator done",
			opts: ListOptions{},
			ds: func() DataStore {
				it := NewMockIterator(ctrl)
				it.EXPECT().
					Next(&model.APIKey{}).
					Return(ErrIteratorDone)

				ds := NewMockDataStore(ctrl)
				ds.EXPECT().
					Find(gomock.Any(), gomock.Any(), ListOptions{}).
					Return(it, nil)
				return ds
			}(),
			wantErr: nil,
		},
		{
			name: "unexpected error occurred",
			opts: ListOptions{},
			ds: func() DataStore {
				it := NewMockIterator(ctrl)
				it.EXPECT().
					Next(&model.APIKey{}).
					Return(errors.New("test-error"))

				ds := NewMockDataStore(ctrl)
				ds.EXPECT().
					Find(gomock.Any(), gomock.Any(), ListOptions{}).
					Return(it, nil)
				return ds
			}(),
			wantErr: errors.New("test-error"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewAPIKeyStore(tc.ds, TestCommander)
			_, err := s.List(context.Background(), tc.opts)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
