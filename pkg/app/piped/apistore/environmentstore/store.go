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
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipe/pkg/app/api/service/pipedservice"
	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/model"
)

type apiClient interface {
	GetEnvironment(ctx context.Context, in *pipedservice.GetEnvironmentRequest, opts ...grpc.CallOption) (*pipedservice.GetEnvironmentResponse, error)
}

// Lister helps list and get Environment.
// All objects returned here must be treated as read-only.
type Lister interface {
	// Get retrieves a specifiec Environment for the given id.
	Get(id string) (*model.Environment, bool)
}

type Store struct {
	apiClient  apiClient
	cache      cache.Cache
	apiTimeout time.Duration
	logger     *zap.Logger
}

func NewStore(apiClient apiClient, cache cache.Cache, logger *zap.Logger) *Store {
	return &Store{
		apiClient:  apiClient,
		cache:      cache,
		apiTimeout: 10 * time.Second,
		logger:     logger.Named("environmentstore"),
	}
}

func (s *Store) Get(id string) (*model.Environment, bool) {
	env, err := s.cache.Get(id)
	if err == nil {
		return env.(*model.Environment), true
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.apiTimeout)
	defer cancel()

	resp, err := s.apiClient.GetEnvironment(ctx, &pipedservice.GetEnvironmentRequest{
		Id: id,
	})
	if err != nil {
		s.logger.Warn("unable to get environment from control plane",
			zap.String("env", id),
			zap.Error(err),
		)
		return nil, false
	}

	if err := s.cache.Put(id, resp.Environment); err != nil {
		s.logger.Warn("unable to put environment to cache", zap.Error(err))
	}
	return resp.Environment, true
}
