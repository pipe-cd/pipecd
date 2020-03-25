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
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kapetaniosci/pipe/pkg/app/api/service"
)

// RunnerAPI implements the behaviors for the gRPC definitions of WebAPI.
type WebAPI struct {
	logger *zap.Logger
}

// NewWebAPIService creates a new service instance.
func NewWebAPIService(logger *zap.Logger) *WebAPI {
	a := &WebAPI{
		logger: logger.Named("web-api"),
	}
	return a
}

// Register registers all handling of this service into the specified gRPC server.
func (a *WebAPI) Register(server *grpc.Server) {
	//service.RegisterWebAPIServer(server, a)
}

func (a *WebAPI) RegisterRunner(ctx context.Context, req *service.RegisterRunnerRequest) (*service.RegisterRunnerResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}
