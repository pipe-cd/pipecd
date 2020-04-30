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

package api

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kapetaniosci/pipe/pkg/app/api/service/webservice"
)

// RunnerAPI implements the behaviors for the gRPC definitions of WebAPI.
type WebAPI struct {
	logger *zap.Logger
}

// NewWebAPI creates a new WebAPI instance.
func NewWebAPI(logger *zap.Logger) *WebAPI {
	a := &WebAPI{
		logger: logger.Named("web-api"),
	}
	return a
}

// Register registers all handling of this service into the specified gRPC server.
func (a *WebAPI) Register(server *grpc.Server) {
	webservice.RegisterWebServiceServer(server, a)
}

func (a *WebAPI) AddEnvironment(ctx context.Context, req *webservice.AddEnvironmentRequest) (*webservice.AddEnvironmentResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *WebAPI) UpdateEnvironmentDesc(ctx context.Context, req *webservice.UpdateEnvironmentDescRequest) (*webservice.UpdateEnvironmentDescResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *WebAPI) ListEnvironments(ctx context.Context, req *webservice.ListEnvironmentsRequest) (*webservice.ListEnvironmentsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *WebAPI) RegisterRunner(ctx context.Context, req *webservice.RegisterRunnerRequest) (*webservice.RegisterRunnerResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *WebAPI) DisableRunner(ctx context.Context, req *webservice.DisableRunnerRequest) (*webservice.DisableRunnerResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *WebAPI) ListRunners(ctx context.Context, req *webservice.ListRunnersRequest) (*webservice.ListRunnersResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *WebAPI) AddApplication(ctx context.Context, req *webservice.AddApplicationRequest) (*webservice.AddApplicationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *WebAPI) DisableApplication(ctx context.Context, req *webservice.DisableApplicationRequest) (*webservice.DisableApplicationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *WebAPI) ListApplications(ctx context.Context, req *webservice.ListApplicationsRequest) (*webservice.ListApplicationsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *WebAPI) SyncApplication(ctx context.Context, req *webservice.SyncApplicationRequest) (*webservice.SyncApplicationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *WebAPI) ListDeployments(ctx context.Context, req *webservice.ListDeploymentsRequest) (*webservice.ListDeploymentsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *WebAPI) GetDeployment(ctx context.Context, req *webservice.GetDeploymentRequest) (*webservice.GetDeploymentResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *WebAPI) CancelDeployment(ctx context.Context, req *webservice.CancelDeploymentRequest) (*webservice.CancelDeploymentResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *WebAPI) ApproveDeploymentStage(ctx context.Context, req *webservice.ApproveDeploymentStageRequest) (*webservice.ApproveDeploymentStageResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *WebAPI) RetryDeploymentStage(ctx context.Context, req *webservice.RetryDeploymentStageRequest) (*webservice.RetryDeploymentStageResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *WebAPI) GetApplicationState(ctx context.Context, req *webservice.GetApplicationStateRequest) (*webservice.GetApplicationStateResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *WebAPI) GetProject(ctx context.Context, req *webservice.GetProjectRequest) (*webservice.GetProjectResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *WebAPI) GetMe(ctx context.Context, req *webservice.GetMeRequest) (*webservice.GetMeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}
