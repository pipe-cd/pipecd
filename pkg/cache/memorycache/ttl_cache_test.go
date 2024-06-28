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

package memorycache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/cache"
)

func TestTTLCache(t *testing.T) {
	t.Run("remain a value after context canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		c := NewTTLCache(ctx, 0, 5*time.Second)
		err := c.Put("key-1", "value-1")
		require.NoError(t, err)
		cancel()
		<-time.After(6 * time.Second)
		value, err := c.Get("key-1")
		require.NoError(t, err)
		assert.Equal(t, "value-1", value)
	})
	t.Run("test eviction", func(t *testing.T) {
		c := NewTTLCache(context.TODO(), 0, 5*time.Second)
		err := c.Put("key-1", "value-1")
		require.NoError(t, err)
		value, err := c.Get("key-1")
		require.NoError(t, err)
		assert.Equal(t, "value-1", value)

		c.evictExpired(time.Now())
		value, err = c.Get("key-1")
		assert.Equal(t, cache.ErrNotFound, err)
		assert.Equal(t, nil, value)
	})
}
