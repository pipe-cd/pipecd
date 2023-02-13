// Copyright 2023 The PipeCD Authors.
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

package kubernetes

import (
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/cache"
)

type AppManifestsCache struct {
	AppID  string
	Cache  cache.Cache
	Logger *zap.Logger
}

func (c AppManifestsCache) Get(commit string) ([]Manifest, bool) {
	key := appManifestsCacheKey(c.AppID, commit)
	item, err := c.Cache.Get(key)
	if err == nil {
		return item.([]Manifest), true
	}

	if errors.Is(err, cache.ErrNotFound) {
		c.Logger.Info("app manifests were not found in cache",
			zap.String("app-id", c.AppID),
			zap.String("commit-hash", commit),
		)
		return nil, false
	}

	c.Logger.Error("failed while retrieving app manifests from cache",
		zap.String("app-id", c.AppID),
		zap.String("commit-hash", commit),
		zap.Error(err),
	)
	return nil, false
}

func (c AppManifestsCache) Put(commit string, manifests []Manifest) {
	key := appManifestsCacheKey(c.AppID, commit)
	if err := c.Cache.Put(key, manifests); err != nil {
		c.Logger.Error("failed while putting app manifests from cache",
			zap.String("app-id", c.AppID),
			zap.String("commit-hash", commit),
			zap.Error(err),
		)
	}
}

func appManifestsCacheKey(appID, commit string) string {
	return fmt.Sprintf("%s/%s", appID, commit)
}
