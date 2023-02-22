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
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/redis"
)

var (
	interval = 10 * time.Minute
)

type apiKeyStore interface {
	List(ctx context.Context, opts datastore.ListOptions) ([]*model.APIKey, error)
	UpdateLastUsedAt(ctx context.Context, id string, time int64) error
}

type apiKeyLastUsedTimeCache interface {
	Get(k string) (interface{}, error)
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

	opts := datastore.ListOptions{
		Filters: []datastore.ListFilter{
			{
				Field:    "Disabled",
				Operator: datastore.OperatorEqual,
				Value:    false,
			},
		},
	}

	apiKeys, err := c.apiKeyStore.List(ctx, opts)
	if err != nil {
		c.logger.Error("failed to list API key", zap.Error(err))
		return err
	}

	for _, apiKey := range apiKeys {
		cachedLastUse, err := c.apiKeyLastUsedTimeCache.Get(apiKey.Id)
		if err != nil {
			c.logger.Error("failed to fetch last used time from cache",
				zap.String("id", apiKey.Id),
				zap.Error(err),
			)
			continue
		}

		lastUsedTime, err := strconv.ParseInt(string(cachedLastUse.([]byte)), 10, 64)
		if err != nil {
			c.logger.Error("failed to fetch last used time from cache",
				zap.String("id", apiKey.Id),
				zap.Error(err),
			)
			continue
		}

		// Skip update last_used_at in database if no changed.
		if lastUsedTime == apiKey.LastUsedAt {
			continue
		}
		if err := c.apiKeyStore.UpdateLastUsedAt(ctx, apiKey.Id, lastUsedTime); err != nil {
			c.logger.Error("failed to update last used time",
				zap.String("id", apiKey.Id),
				zap.Error(err),
			)
			continue
		}
	}

	return nil
}
