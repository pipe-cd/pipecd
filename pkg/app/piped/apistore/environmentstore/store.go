// Copyright 2020 The PipeCD Authors.
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

package environmentstore

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipe/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/model"
)

const (
	defaultAPITimeout = time.Minute
)

type apiClient interface {
	GetEnvironment(ctx context.Context, in *pipedservice.GetEnvironmentRequest, opts ...grpc.CallOption) (*pipedservice.GetEnvironmentResponse, error)
}

// Lister helps list and get Environment.
// All objects returned here must be treated as read-only.
type Lister interface {
	Get(ctx context.Context, id string) (*model.Environment, error)
	GetByName(ctx context.Context, name string) (*model.Environment, error)
}

type Store struct {
	apiClient apiClient
	// A goroutine-safe map from id to Environment.
	cache cache.Cache
	// A goroutine-safe map from name to Environment.
	cacheByName cache.Cache
	callGroup   *singleflight.Group
	logger      *zap.Logger
}

func NewStore(apiClient apiClient, cache, cacheByName cache.Cache, logger *zap.Logger) *Store {
	return &Store{
		apiClient:   apiClient,
		cache:       cache,
		cacheByName: cacheByName,
		callGroup:   &singleflight.Group{},
		logger:      logger.Named("environmentstore"),
	}
}

func (s *Store) Get(ctx context.Context, id string) (*model.Environment, error) {
	env, err := s.cache.Get(id)
	if err == nil {
		return env.(*model.Environment), nil
	}

	// Ensure that timeout is configured.
	ctx, cancel := context.WithTimeout(ctx, defaultAPITimeout)
	defer cancel()

	// Ensure that only one RPC call is executed for the given key at a time
	// and the newest data is stored in the cache.
	data, err, _ := s.callGroup.Do(id, func() (interface{}, error) {
		req := &pipedservice.GetEnvironmentRequest{
			Id: id,
		}
		resp, err := s.apiClient.GetEnvironment(ctx, req)
		if err != nil {
			s.logger.Warn("failed to get environment from control plane",
				zap.String("env", id),
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to get environment %s, %w", id, err)
		}

		if err := s.cache.Put(id, resp.Environment); err != nil {
			s.logger.Warn("failed to put environment to cache", zap.Error(err))
		}
		return resp.Environment, nil
	})

	if err != nil {
		return nil, err
	}
	return data.(*model.Environment), nil
}
