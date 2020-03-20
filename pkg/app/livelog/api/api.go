// Copyright 2020 The Pipe Authors.
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

package api

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/kapetaniosci/pipe/pkg/app/livelog/service"
)

type api struct {
	dataDir string
	logger  *zap.Logger
}

type Option func(*api)

func WithLogger(logger *zap.Logger) Option {
	return func(a *api) {
		a.logger = logger.Named("api")
	}
}

// NewLiveLogService creates a new service that is implementing LiveLog gRPC interface.
func NewLiveLogService(dataDir string, opts ...Option) *api {
	a := &api{
		dataDir: dataDir,
		logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *api) Register(server *grpc.Server) {
	service.RegisterLiveLogServer(server, a)
}
