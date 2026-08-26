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

package client

import (
	"context"

	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
)

type fakeAPIClient struct {
	apiservice.APIServiceClient

	syncApplication func(context.Context, *apiservice.SyncApplicationRequest, ...grpc.CallOption) (*apiservice.SyncApplicationResponse, error)
	getDeployment   func(context.Context, *apiservice.GetDeploymentRequest, ...grpc.CallOption) (*apiservice.GetDeploymentResponse, error)
	getCommand      func(context.Context, *apiservice.GetCommandRequest, ...grpc.CallOption) (*apiservice.GetCommandResponse, error)
}

func (f *fakeAPIClient) SyncApplication(ctx context.Context, req *apiservice.SyncApplicationRequest, opts ...grpc.CallOption) (*apiservice.SyncApplicationResponse, error) {
	return f.syncApplication(ctx, req, opts...)
}

func (f *fakeAPIClient) GetDeployment(ctx context.Context, req *apiservice.GetDeploymentRequest, opts ...grpc.CallOption) (*apiservice.GetDeploymentResponse, error) {
	return f.getDeployment(ctx, req, opts...)
}

func (f *fakeAPIClient) GetCommand(ctx context.Context, req *apiservice.GetCommandRequest, opts ...grpc.CallOption) (*apiservice.GetCommandResponse, error) {
	return f.getCommand(ctx, req, opts...)
}

func (f *fakeAPIClient) Close() error {
	return nil
}
