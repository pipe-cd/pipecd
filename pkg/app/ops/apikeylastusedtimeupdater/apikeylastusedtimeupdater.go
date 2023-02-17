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

package apikeylastusedtimeupdater

import (
	"context"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/cache/rediscache"
	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/redis"
)

var (
	commandTimeOut = 24 * time.Hour
	interval       = 1 * time.Minute
)

type apiKeyStore interface {
	UpdateLastUsedAt(ctx context.Context, id string, time int64) error
}

type apiKeyLastUsedTimeCache interface {
	GetAll() (map[string]interface{}, error)
}

type APIKeyLastUsedTimeUpdater struct {
	apiKeyStore             apiKeyStore
	apiKeyLastUsedTimeCache apiKeyLastUsedTimeCache
	logger                  *zap.Logger
}

const apiKeyLastUsedCacheHashKey = "HASHKEY:PIPED:API_KEYS" //nolint:gosec

func NewAPIKeyLastUsedTimeUpdater(
	ds datastore.DataStore,
	rd redis.Redis,
	logger *zap.Logger,
) *APIKeyLastUsedTimeUpdater {
	return &APIKeyLastUsedTimeUpdater{
		apiKeyStore:             datastore.NewAPIKeyStore(ds, datastore.OpsCommander),
		apiKeyLastUsedTimeCache: rediscache.NewHashCache(rd, apiKeyLastUsedCacheHashKey),
		logger:                  logger.Named("api-key-last-used-time-updater"),
	}
}

func (c *APIKeyLastUsedTimeUpdater) Run(ctx context.Context) error {
	c.logger.Info("start running APIKeyLastUsedTimeUpdater")

	t := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("APIKeyLastUsedTimeUpdater has been stopped")
			return nil

		case <-t.C:
			start := time.Now()
			if err := c.updateAPIKeyLastUsedTime(ctx); err == nil {
				c.logger.Info("successfully update api key last used time", zap.Duration("duration", time.Since(start)))
			}
		}
	}
}

func (c *APIKeyLastUsedTimeUpdater) updateAPIKeyLastUsedTime(ctx context.Context) error {
	keys, err := c.apiKeyLastUsedTimeCache.GetAll()
	if err != nil {
		c.logger.Info("there are no cache of api key last used time on redis")
	}

	for id, time := range keys {
		lastUsedTime := bytes2int64(time.([]byte))
		if err := c.apiKeyStore.UpdateLastUsedAt(ctx, id, lastUsedTime); err != nil {
			c.logger.Error("failed to update last used time",
				zap.String("id", id),
				zap.Error(err),
			)
			return err
		}
	}

	return nil
}

func bytes2int64(bytes []byte) int64 {
	var numString string
	for i := range bytes {
		numString += string(bytes[i])
	}
	num, _ := strconv.ParseInt(numString, 10, 64)
	return num
}
