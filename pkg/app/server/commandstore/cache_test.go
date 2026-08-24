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
			name: "valid command is stored as before",
			command: &model.Command{
				Id:      "command-id",
				PipedId: "piped-id",
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
