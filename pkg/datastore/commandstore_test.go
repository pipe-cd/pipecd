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

func TestAddCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		name      string
		command   *model.Command
		dsFactory func(*model.Command) DataStore
		wantErr   bool
	}{
		{
			name:      "Invalid command",
			command:   &model.Command{},
			dsFactory: func(d *model.Command) DataStore { return nil },
			wantErr:   true,
		},
		{
			name: "Valid command",
			command: &model.Command{
				Id:            "id",
				ApplicationId: "app-id",
				PipedId:       "piped-id",
				CreatedAt:     1,
				UpdatedAt:     1,
			},
			dsFactory: func(d *model.Command) DataStore {
				ds := NewMockDataStore(ctrl)
				ds.EXPECT().Create(gomock.Any(), gomock.Any(), d.Id, d)
				return ds
			},
			wantErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewCommandStore(tc.dsFactory(tc.command), TestCommander)
			err := s.Add(context.Background(), tc.command)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestGetCommand(t *testing.T) {
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
					Get(gomock.Any(), gomock.Any(), "id", &model.Command{}).
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
					Get(gomock.Any(), gomock.Any(), "id", &model.Command{}).
					Return(fmt.Errorf("err"))
				return ds
			}(),
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewCommandStore(tc.ds, TestCommander)
			_, err := s.Get(context.Background(), tc.id)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestListCommands(t *testing.T) {
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
					Next(&model.Command{}).
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
					Next(&model.Command{}).
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
			s := NewCommandStore(tc.ds, TestCommander)
			_, err := s.List(context.Background(), tc.opts)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestCommandDecode(t *testing.T) {
	col := &commandCollection{requestedBy: TestCommander}

	testcases := []struct {
		name      string
		parts     map[Shard][]byte
		expectCmd *model.Command
		expectErr bool
	}{
		{
			name:      "parts count miss matched",
			parts:     make(map[Shard][]byte),
			expectErr: true,
		},
		{
			name: "should merge correctly",
			parts: map[Shard][]byte{
				AgentShard: []byte(`{"id":"1","status":3,"updated_at":4}`),
				OpsShard:   []byte(`{"id":"1","status":0,"updated_at":1}`),
			},
			expectCmd: &model.Command{
				Id:        "1",
				Status:    model.CommandStatus_COMMAND_TIMEOUT,
				UpdatedAt: 4,
			},
			expectErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := &model.Command{}
			err := col.Decode(cmd, tc.parts)
			require.Equal(t, tc.expectErr, err != nil)

			if err == nil {
				assert.Equal(t, tc.expectCmd, cmd)
			}
		})
	}
}
