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

package unregisteredappstore

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/cache/rediscache"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/redis"
)

type Store interface {
	ListApplications(ctx context.Context, projectID string) ([]*model.ApplicationInfo, error)
	PutApplications(projectID, pipedID string, apps []*model.ApplicationInfo) error
}

type store struct {
	backend redis.Redis
	logger  *zap.Logger
}

func NewStore(r redis.Redis, logger *zap.Logger) Store {
	return &store{
		backend: r,
		logger:  logger,
	}
}

func (c *store) ListApplications(_ context.Context, projectID string) ([]*model.ApplicationInfo, error) {
	key := makeUnregisteredAppsCacheKey(projectID)
	hc := rediscache.NewHashCache(c.backend, key)

	// pipedToApps assumes to be a map["piped-id"][]byte(slice of *model.ApplicationInfo encoded by encoding/gob)
	pipedToApps, err := hc.GetAll()
	if errors.Is(err, cache.ErrNotFound) {
		return []*model.ApplicationInfo{}, nil
	}

	if err != nil {
		c.logger.Error("failed to get unregistered apps", zap.Error(err))
		return nil, err
	}

	// Integrate all apps cached for each Piped.
	allApps := make([]*model.ApplicationInfo, 0)
	for _, as := range pipedToApps {
		b, ok := as.([]byte)
		if !ok {
			return nil, errors.New("unexpected data cached")
		}

		dec := gob.NewDecoder(bytes.NewReader(b))
		var apps []*model.ApplicationInfo
		if err := dec.Decode(&apps); err != nil {
			c.logger.Error("failed to decode the unregistered apps", zap.Error(err))
			return nil, err
		}

		allApps = append(allApps, apps...)
	}

	return allApps, nil
}

func (c *store) PutApplications(projectID, pipedID string, apps []*model.ApplicationInfo) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(apps); err != nil {
		c.logger.Error("failed to encode the unregistered apps", zap.Error(err))
		return err
	}

	key := makeUnregisteredAppsCacheKey(projectID)
	hc := rediscache.NewHashCache(c.backend, key)
	if err := hc.Put(pipedID, buf.Bytes()); err != nil {
		return err
	}

	return nil
}

func makeUnregisteredAppsCacheKey(projectID string) string {
	return fmt.Sprintf("HASHKEY:UNREGISTERED_APPS:%s", projectID)
}
