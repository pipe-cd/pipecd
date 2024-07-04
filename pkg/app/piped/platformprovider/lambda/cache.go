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

package lambda

import (
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/cache"
)

type FunctionManifestCache struct {
	AppID  string
	Cache  cache.Cache
	Logger *zap.Logger
}

func (c FunctionManifestCache) Get(commit string) (FunctionManifest, bool) {
	key := manifestCacheKey(c.AppID, commit)
	item, err := c.Cache.Get(key)
	if err == nil {
		return item.(FunctionManifest), true
	}

	if errors.Is(err, cache.ErrNotFound) {
		c.Logger.Info("function manifest wes not found in cache",
			zap.String("app-id", c.AppID),
			zap.String("commit-hash", commit),
		)
		return FunctionManifest{}, false
	}

	c.Logger.Error("failed while retrieving function manifest from cache",
		zap.String("app-id", c.AppID),
		zap.String("commit-hash", commit),
		zap.Error(err),
	)
	return FunctionManifest{}, false
}

func (c FunctionManifestCache) Put(commit string, sm FunctionManifest) {
	key := manifestCacheKey(c.AppID, commit)
	if err := c.Cache.Put(key, sm); err != nil {
		c.Logger.Error("failed while putting function manifest into cache",
			zap.String("app-id", c.AppID),
			zap.String("commit-hash", commit),
			zap.Error(err),
		)
	}
}

func manifestCacheKey(appID, commit string) string {
	return fmt.Sprintf("%s/%s", appID, commit)
}
