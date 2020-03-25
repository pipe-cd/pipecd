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

	"github.com/kapetaniosci/pipe/pkg/app/api/service"
)

// RunnerAPI implements the behaviors for the gRPC definitions of RunnerAPI.
type RunnerAPI struct {
	logger *zap.Logger
}

// NewRunnerAPIService creates a new service instance.
func NewRunnerAPIService(logger *zap.Logger) *RunnerAPI {
	a := &RunnerAPI{
		logger: logger.Named("runner-api"),
	}
	return a
}

// Register registers all handling of this service into the specified gRPC server.
func (a *RunnerAPI) Register(server *grpc.Server) {
	service.RegisterRunnerAPIServer(server, a)
}

// Ping is periodically sent by runner to report its status/stats to API.
// The received stats will be written to the cache immediately.
// The cache data may be lost anytime so we need a singleton Persister
// to persist those data into datastore every n minutes.
func (a *RunnerAPI) Ping(ctx context.Context, req *service.PingRequest) (*service.PingResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ListApplicationsByRunner returns a list of registered applications
// that should be managed by the requested runner.
// Disabled applications should not be included in the response.
// Runner uses this RPC to fetch and sync the application configuration into its local database.
func (a *RunnerAPI) ListApplicationsByRunner(ctx context.Context, req *service.ListApplicationsByRunnerRequest) (*service.ListApplicationsByRunnerResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// RegisterEvents is sent by runner to submit one or multiple events
// about executing pipelines and application resources.
// Control plane uses the received events to update the state of pipeline/application-resource-tree.
// We want to start by a simple solution at this initial stage of development,
// so the API server just handles as below:
// - loads the releated pipeline/application-resource-tree from datastore
// - checks and builds new state for the pipeline/application-resource-tree
// - updates new state into datastore and cache (cache data is for reading while handling web requests)
// In the future, we may want to redesign the behavior of this RPC by using pubsub/queue pattern.
// After receiving the events, all of them will be publish into a queue immediately,
// and then another Handler service will pick them inorder to apply to build new state.
// By that way we can control the traffic to the datastore in a better way.
func (a *RunnerAPI) RegisterEvents(ctx context.Context, req *service.RegisterEventsRequest) (*service.RegisterEventsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// SendStageCompleteLog is sent by runner to save the complete log
// of a pipeline stage after its completion.
// The complete log will be persisted into the filestore immediately.
func (a *RunnerAPI) SendStageCompleteLog(ctx context.Context, req *service.SendStageCompleteLogRequest) (*service.SendStageCompleteLogResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// GetCommands is periodically called by runner to obtain the commands
// that should be handled.
// Whenever an user makes an interaction from WebUI (cancel/approve/retry/sync)
// a new command with a unique identifier will be generated an saved into the datastore.
// Runner uses this RPC to list all still-not-handled commands to handle them,
// then report back the result to server.
// On other side, the web will periodically check the command status and feedback the result to user.
// In the future, we may need a solution to remove all old-handled commands from datastore for space.
func (a *RunnerAPI) GetCommands(ctx context.Context, req *service.GetCommandsRequest) (*service.GetCommandsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ReportCommandHandled is called by runner to mark a specific command as handled.
// The request payload will contain the handle status as well as any additional result data.
// The handle result should be updated to both datastore and cache (for reading from web).
func (a *RunnerAPI) ReportCommandHandled(ctx context.Context, req *service.ReportCommandHandledRequest) (*service.ReportCommandHandledResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ReportPipelineStarted used by runner to send initial state
// of the pipeline that has just been started.
func (a *RunnerAPI) ReportPipelineStarted(ctx context.Context, req *service.ReportPipelineStartedRequest) (*service.ReportPipelineStartedResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ReportPipelineCompleted used by runner to send the final state
// of the pipeline that has just been completed.
func (a *RunnerAPI) ReportPipelineCompleted(ctx context.Context, req *service.ReportPipelineCompletedRequest) (*service.ReportPipelineCompletedResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ReportApplicationResourceTree is periodically sent by runner to refresh the full tree of application resources.
// The tree data will be written into filestore and the cache inmmediately.
func (a *RunnerAPI) ReportApplicationResourceTree(ctx context.Context, req *service.ReportApplicationResourceTreeRequest) (*service.ReportApplicationResourceTreeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}
