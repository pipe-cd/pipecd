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
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestAddApplication(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		name        string
		application *model.Application
		dsFactory   func(*model.Application) DataStore
		wantErr     bool
	}{
		{
			name:        "Invalid application",
			application: &model.Application{},
			dsFactory:   func(d *model.Application) DataStore { return nil },
			wantErr:     true,
		},
		{
			name: "Valid application",
			application: &model.Application{
				Id:        "id",
				Name:      "name",
				PipedId:   "piped-id",
				ProjectId: "project-id",
				Kind:      model.ApplicationKind_KUBERNETES,
				GitPath: &model.ApplicationGitPath{
					Repo: &model.ApplicationGitRepository{Id: "id"},
					Path: "path",
				},
				CloudProvider: "cloud-provider",

				CreatedAt: 1,
				UpdatedAt: 1,
			},
			dsFactory: func(d *model.Application) DataStore {
				ds := NewMockDataStore(ctrl)
				ds.EXPECT().Create(gomock.Any(), gomock.Any(), d.Id, d)
				return ds
			},
			wantErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewApplicationStore(tc.dsFactory(tc.application), TestCommander)
			err := s.Add(context.Background(), tc.application)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestGetApplication(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		name    string
		id      string
		ds      DataStore
		wantErr bool
	}{
		{
			name: "successful fetch from datastore",
			id:   "id",
			ds: func() DataStore {
				ds := NewMockDataStore(ctrl)
				ds.EXPECT().
					Get(gomock.Any(), gomock.Any(), "id", &model.Application{}).
					Return(nil)
				return ds
			}(),
			wantErr: false,
		},
		{
			name: "failed fetch from datastore",
			id:   "id",
			ds: func() DataStore {
				ds := NewMockDataStore(ctrl)
				ds.EXPECT().
					Get(gomock.Any(), gomock.Any(), "id", &model.Application{}).
					Return(fmt.Errorf("err"))
				return ds
			}(),
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewApplicationStore(tc.ds, TestCommander)
			_, err := s.Get(context.Background(), tc.id)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestListApplications(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		name    string
		opts    ListOptions
		ds      DataStore
		wantErr bool
	}{
		{
			name: "iterator done",
			opts: ListOptions{},
			ds: func() DataStore {
				it := NewMockIterator(ctrl)
				it.EXPECT().
					Next(&model.Application{}).
					Return(ErrIteratorDone)

				ds := NewMockDataStore(ctrl)
				ds.EXPECT().
					Find(gomock.Any(), gomock.Any(), ListOptions{}).
					Return(it, nil)
				return ds
			}(),
			wantErr: false,
		},
		{
			name: "unexpected error occurred",
			opts: ListOptions{},
			ds: func() DataStore {
				it := NewMockIterator(ctrl)
				it.EXPECT().
					Next(&model.Application{}).
					Return(fmt.Errorf("err"))

				ds := NewMockDataStore(ctrl)
				ds.EXPECT().
					Find(gomock.Any(), gomock.Any(), ListOptions{}).
					Return(it, nil)
				return ds
			}(),
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewApplicationStore(tc.ds, TestCommander)
			_, _, err := s.List(context.Background(), tc.opts)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestApplicationDecode(t *testing.T) {
	col := &applicationCollection{requestedBy: TestCommander}

	testcases := []struct {
		name      string
		parts     map[Shard][]byte
		expectApp *model.Application
		expectErr bool
	}{
		{
			name:      "shard count miss matched",
			parts:     make(map[Shard][]byte),
			expectErr: true,
		},
		{
			name: "decode correctly",
			parts: map[Shard][]byte{
				ClientShard: []byte(`{"kind":1,"name":"name","piped_id":"new_piped","platform_provider":"new_provider","git_path":{"config_filename":"new_file"},"updated_at":123}`),
				AgentShard:  []byte(`{"kind":0,"name":"new_name","piped_id":"piped","platform_provider":"provider","git_path":{"config_filename":"file"},"updated_at":1}`),
			},
			expectApp: &model.Application{
				Kind:             model.ApplicationKind_KUBERNETES,
				Name:             "new_name",
				PipedId:          "new_piped",
				PlatformProvider: "new_provider",
				GitPath: &model.ApplicationGitPath{
					ConfigFilename: "new_file",
				},
				UpdatedAt: 123,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			app := &model.Application{}
			err := col.Decode(app, tc.parts)
			require.Equal(t, tc.expectErr, err != nil)

			if err == nil {
				assert.Equal(t, tc.expectApp, app)
			}
		})
	}
}
