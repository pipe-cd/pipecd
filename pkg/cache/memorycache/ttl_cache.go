// Copyright 2020 The PipeCD Authors.
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
	"sync"
	"time"

	"github.com/kapetaniosci/pipe/pkg/cache"
)

type entry struct {
	value      interface{}
	expiration time.Time
}

type TTLCache struct {
	entries sync.Map
	ttl     time.Duration
	ctx     context.Context
}

func NewTTLCache(ctx context.Context, ttl time.Duration, evictionInterval time.Duration) *TTLCache {
	c := &TTLCache{
		ttl: ttl,
		ctx: ctx,
	}
	if evictionInterval > 0 {
		go c.startEvicter(evictionInterval)
	}
	return c
}

func (c *TTLCache) startEvicter(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case now := <-ticker.C:
			c.evictExpired(now)
		case <-c.ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (c *TTLCache) evictExpired(t time.Time) {
	c.entries.Range(func(key interface{}, value interface{}) bool {
		e := value.(*entry)
		if e.expiration.Before(t) {
			c.entries.Delete(key)
		}
		return true
	})
}

func (c *TTLCache) Get(key interface{}) (interface{}, error) {
	item, ok := c.entries.Load(key)
	if !ok {
		return nil, cache.ErrNotFound
	}
	return item.(*entry).value, nil
}

func (c *TTLCache) Put(key interface{}, value interface{}) error {
	e := &entry{
		value:      value,
		expiration: time.Now().Add(c.ttl),
	}
	c.entries.Store(key, e)
	return nil
}

func (c *TTLCache) Delete(key interface{}) error {
	c.entries.Delete(key)
	return nil
}
