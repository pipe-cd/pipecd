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

package grpcapi

import (
	"fmt"
	"sync"
	"time"
)

const planPreviewRepositoryCachePrefix = "plan-preview-repository"

type planPreviewRepositoryCache struct {
	mu         sync.Mutex
	maxEntries int
	ttl        time.Duration
	entries    map[string]planPreviewRepositoryCacheEntry
	nowFunc    func() time.Time
}

type planPreviewRepositoryCacheEntry struct {
	repositories map[string]string
	expiresAt    time.Time
}

func newPlanPreviewRepositoryCache(maxEntries int, ttl time.Duration) *planPreviewRepositoryCache {
	return &planPreviewRepositoryCache{
		maxEntries: maxEntries,
		ttl:        ttl,
		entries:    make(map[string]planPreviewRepositoryCacheEntry, maxEntries),
		nowFunc:    time.Now,
	}
}

func planPreviewRepositoryCacheKey(projectID, remoteURL, baseBranch string) string {
	return fmt.Sprintf("%s:%d:%s:%d:%s:%d:%s", planPreviewRepositoryCachePrefix, len(projectID), projectID, len(remoteURL), remoteURL, len(baseBranch), baseBranch)
}

func (c *planPreviewRepositoryCache) Get(key string) (map[string]string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	if !entry.expiresAt.After(c.nowFunc()) {
		delete(c.entries, key)
		return nil, false
	}
	return cloneStringMap(entry.repositories), true
}

func (c *planPreviewRepositoryCache) Put(key string, repositories map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.nowFunc()
	if len(c.entries) >= c.maxEntries {
		c.evictExpired(now)
	}
	for len(c.entries) >= c.maxEntries {
		for k := range c.entries {
			delete(c.entries, k)
			break
		}
	}
	c.entries[key] = planPreviewRepositoryCacheEntry{
		repositories: cloneStringMap(repositories),
		expiresAt:    now.Add(c.ttl),
	}
}

func (c *planPreviewRepositoryCache) evictExpired(now time.Time) {
	for key, entry := range c.entries {
		if !entry.expiresAt.After(now) {
			delete(c.entries, key)
		}
	}
}

func cloneStringMap(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
