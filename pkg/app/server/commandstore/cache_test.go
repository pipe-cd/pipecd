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

package commandstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/cache/cachetest"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestCommandCachePut(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		name    string
		command *model.Command
		prepare func(c *cachetest.MockCache)
		wantErr bool
	}{
		{
			name:    "nil command is rejected before reaching the cache backend",
			command: nil,
			prepare: func(c *cachetest.MockCache) {
				// No EXPECT() set: the mock fails the test if the
				// backend is ever called with a nil-derived payload.
			},
			wantErr: true,
		},
		{
			name: "command failing Validate() is rejected before reaching the cache backend",
			command: &model.Command{
				Id: "command-id",
				// PipedId, CreatedAt and UpdatedAt are left unset, so
				// Validate() fails.
			},
			prepare: func(c *cachetest.MockCache) {
				// No EXPECT() set: the mock fails the test if an invalid
				// command ever reaches the backend.
			},
			wantErr: true,
		},
		{
			name: "valid command is stored as before",
			command: &model.Command{
				Id:        "command-id",
				PipedId:   "piped-id",
				CreatedAt: 1700000000,
				UpdatedAt: 1700000000,
			},
			prepare: func(c *cachetest.MockCache) {
				c.EXPECT().Put(cacheKey("command-id"), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := cachetest.NewMockCache(ctrl)
			tc.prepare(c)

			cc := &commandCache{backend: c}
			err := cc.Put("command-id", tc.command)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

// TestCommandCacheGet_InvalidEntry covers cache entries that were already
// poisoned before commandCache.Put started rejecting nil commands: a raw
// JSON "null" (unmarshals to a zero-value Command with an empty Id and
// PipedId) and any other entry that fails model.Command.Validate() (the same
// validation datastore.commandStore.Add already enforces before a command is
// ever persisted, so no legitimately-stored command can fail it). Both must
// be treated as a cache miss, not returned as if they were real commands,
// and the poisoned key should be evicted on a best-effort basis.
func TestCommandCacheGet_InvalidEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		name string
		data string
	}{
		{
			name: "previously poisoned null entry",
			data: "null",
		},
		{
			name: "cached command missing required PipedId",
			data: `{"id":"command-id"}`,
		},
		{
			name: "cached command missing required Id",
			data: `{"piped_id":"piped-id"}`,
		},
		{
			name: "cached command missing required CreatedAt/UpdatedAt",
			data: `{"id":"command-id","piped_id":"piped-id"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := cachetest.NewMockCache(ctrl)
			key := cacheKey("command-id")
			c.EXPECT().Get(key).Return([]byte(tc.data), nil)
			// The invalid entry must be evicted on a best-effort basis.
			c.EXPECT().Delete(key).Return(nil)

			cc := &commandCache{backend: c}
			got, err := cc.Get("command-id")
			assert.Nil(t, got)
			assert.ErrorIs(t, err, cache.ErrNotFound)
		})
	}
}

// TestCommandCacheGet_ValidEntry ensures a genuinely valid cached command is
// still returned directly, without being mistaken for a poisoned entry.
func TestCommandCacheGet_ValidEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := cachetest.NewMockCache(ctrl)
	key := cacheKey("command-id")
	c.EXPECT().Get(key).Return([]byte(`{"id":"command-id","piped_id":"piped-id","created_at":1700000000,"updated_at":1700000000}`), nil)

	cc := &commandCache{backend: c}
	got, err := cc.Get("command-id")
	assert.NoError(t, err)
	assert.Equal(t, &model.Command{Id: "command-id", PipedId: "piped-id", CreatedAt: 1700000000, UpdatedAt: 1700000000}, got)
}
