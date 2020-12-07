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

package grpcapi

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipe/pkg/app/api/commandstore"
	"github.com/pipe-cd/pipe/pkg/app/api/service/apiservice"
	"github.com/pipe-cd/pipe/pkg/datastore"
)

// API implements the behaviors for the gRPC definitions of API.
type API struct {
	applicationStore datastore.ApplicationStore
	deploymentStore  datastore.DeploymentStore
	pipedStore       datastore.PipedStore
	commandStore     commandstore.Store

	logger *zap.Logger
}

// NewAPI creates a new API instance.
func NewAPI(
	ds datastore.DataStore,
	cmds commandstore.Store,
	logger *zap.Logger,
) *API {
	a := &API{
		applicationStore: datastore.NewApplicationStore(ds),
		deploymentStore:  datastore.NewDeploymentStore(ds),
		pipedStore:       datastore.NewPipedStore(ds),
		commandStore:     cmds,
		logger:           logger.Named("api"),
	}
	return a
}

// Register registers all handling of this service into the specified gRPC server.
func (a *API) Register(server *grpc.Server) {
	apiservice.RegisterAPIServiceServer(server, a)
}

func (a *API) AddApplication(_ context.Context, _ *apiservice.AddApplicationRequest) (*apiservice.AddApplicationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *API) SyncApplication(_ context.Context, _ *apiservice.SyncApplicationRequest) (*apiservice.SyncApplicationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}
