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

package cloudrun

import (
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/cache"
)

type ServiceManifestCache struct {
	AppID  string
	Cache  cache.Cache
	Logger *zap.Logger
}

func (c ServiceManifestCache) Get(commit string) (ServiceManifest, bool) {
	key := serviceManifestCacheKey(c.AppID, commit)
	item, err := c.Cache.Get(key)
	if err == nil {
		return item.(ServiceManifest), true
	}

	if errors.Is(err, cache.ErrNotFound) {
		c.Logger.Info("service manifest were not found in cache",
			zap.String("app-id", c.AppID),
			zap.String("commit-hash", commit),
		)
		return ServiceManifest{}, false
	}

	c.Logger.Error("failed while retrieving service manifest from cache",
		zap.String("app-id", c.AppID),
		zap.String("commit-hash", commit),
		zap.Error(err),
	)
	return ServiceManifest{}, false
}

func (c ServiceManifestCache) Put(commit string, sm ServiceManifest) {
	key := serviceManifestCacheKey(c.AppID, commit)
	if err := c.Cache.Put(key, sm); err != nil {
		c.Logger.Error("failed while putting service manifest into cache",
			zap.String("app-id", c.AppID),
			zap.String("commit-hash", commit),
			zap.Error(err),
		)
	}
}

func serviceManifestCacheKey(appID, commit string) string {
	return fmt.Sprintf("%s/%s", appID, commit)
}
