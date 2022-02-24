// Copyright 2022 The PipeCD Authors.
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

package unregisteredappcache

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/cache/rediscache"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/redis"
)

type Cache interface {
	ListUnregisteredApplications(ctx context.Context, projectID string) ([]*model.ApplicationInfo, error)
	PutUnregisteredApplications(projectID, pipedID string, apps []*model.ApplicationInfo) error
}

type unregisteredApplicationCache struct {
	backend redis.Redis
	logger  *zap.Logger
}

func NewCache(r redis.Redis, logger *zap.Logger) Cache {
	return &unregisteredApplicationCache{
		backend: r,
		logger:  logger,
	}
}

func (c *unregisteredApplicationCache) ListUnregisteredApplications(ctx context.Context, projectID string) ([]*model.ApplicationInfo, error) {
	key := makeUnregisteredAppsCacheKey(projectID)
	hc := rediscache.NewHashCache(c.backend, key)

	// pipedToApps assumes to be a map["piped-id"][]byte(slice of *model.ApplicationInfo encoded by encoding/gob)
	pipedToApps, err := hc.GetAll()
	if errors.Is(err, cache.ErrNotFound) {
		return []*model.ApplicationInfo{}, nil
	}

	if err != nil {
		c.logger.Error("failed to get unregistered apps", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to get unregistered apps")
	}

	// Integrate all apps cached for each Piped.
	allApps := make([]*model.ApplicationInfo, 0)
	for _, as := range pipedToApps {
		b, ok := as.([]byte)
		if !ok {
			return nil, status.Error(codes.Internal, "Unexpected data cached")
		}

		dec := gob.NewDecoder(bytes.NewReader(b))
		var apps []*model.ApplicationInfo
		if err := dec.Decode(&apps); err != nil {
			c.logger.Error("failed to decode the unregistered apps", zap.Error(err))
			return nil, status.Error(codes.Internal, "failed to decode the unregistered apps")
		}

		allApps = append(allApps, apps...)
	}

	return allApps, nil
}

func (c *unregisteredApplicationCache) PutUnregisteredApplications(projectID, pipedID string, apps []*model.ApplicationInfo) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(apps); err != nil {
		c.logger.Error("failed to encode the unregistered apps", zap.Error(err))
		return status.Error(codes.Internal, "failed to encode the unregistered apps")
	}

	key := makeUnregisteredAppsCacheKey(projectID)
	hc := rediscache.NewHashCache(c.backend, key)
	if err := hc.Put(pipedID, buf.Bytes()); err != nil {
		return status.Error(codes.Internal, "failed to put the unregistered apps to the cache")
	}

	return nil
}

func makeUnregisteredAppsCacheKey(projectID string) string {
	return fmt.Sprintf("HASHKEY:UNREGISTERED_APPS:%s", projectID)
}
