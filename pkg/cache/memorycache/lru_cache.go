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
	lru "github.com/hashicorp/golang-lru"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/cache/cachemetrics"
)

type LRUCache struct {
	cache *lru.Cache
}

func NewLRUCache(size int) (*LRUCache, error) {
	cache, err := lru.New(size)
	if err != nil {
		return nil, err
	}
	return &LRUCache{
		cache: cache,
	}, nil
}

func (c *LRUCache) Get(key string) (interface{}, error) {
	item, ok := c.cache.Get(key)
	if !ok {
		cachemetrics.IncGetOperationCounter(
			cachemetrics.LabelSourceInmemory,
			cachemetrics.LabelStatusMiss,
		)
		return nil, cache.ErrNotFound
	}
	cachemetrics.IncGetOperationCounter(
		cachemetrics.LabelSourceInmemory,
		cachemetrics.LabelStatusHit,
	)
	return item, nil
}

func (c *LRUCache) Put(key string, value interface{}) error {
	c.cache.Add(key, value)
	return nil
}

func (c *LRUCache) Delete(key string) error {
	c.cache.Remove(key)
	return nil
}

func (c *LRUCache) GetAll() (map[string]interface{}, error) {
	return nil, cache.ErrUnimplemented
}
