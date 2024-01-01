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
	"sync"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/cache/cachemetrics"
)

type Cache struct {
	values sync.Map
}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) Get(key string) (interface{}, error) {
	item, ok := c.values.Load(key)
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

func (c *Cache) Put(key string, value interface{}) error {
	c.values.Store(key, value)
	return nil
}

func (c *Cache) Delete(key string) error {
	c.values.Delete(key)
	return nil
}

func (c *Cache) GetAll() (map[string]interface{}, error) {
	return nil, cache.ErrUnimplemented
}
