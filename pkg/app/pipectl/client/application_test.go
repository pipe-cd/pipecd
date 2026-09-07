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
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
)

func TestSyncApplicationRetriesUnavailableCommandUntilTimeout(t *testing.T) {
	const commandID = "command-id"
	calls := 0

	cli := &fakeAPIClient{
		syncApplication: func(context.Context, *apiservice.SyncApplicationRequest, ...grpc.CallOption) (*apiservice.SyncApplicationResponse, error) {
			return &apiservice.SyncApplicationResponse{
				CommandId: commandID,
			}, nil
		},
		getCommand: func(context.Context, *apiservice.GetCommandRequest, ...grpc.CallOption) (*apiservice.GetCommandResponse, error) {
			calls++
			return nil, status.Error(codes.Unavailable, "service unavailable")
		},
	}

	_, err := SyncApplication(
		context.Background(),
		cli,
		"app-id",
		time.Millisecond,
		50*time.Millisecond,
		zap.NewNop(),
	)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.GreaterOrEqual(t, calls, 2)
}

func TestSyncApplicationReturnsNotFoundForMissingCommand(t *testing.T) {
	const commandID = "missing-command"

	cli := &fakeAPIClient{
		syncApplication: func(context.Context, *apiservice.SyncApplicationRequest, ...grpc.CallOption) (*apiservice.SyncApplicationResponse, error) {
			return &apiservice.SyncApplicationResponse{
				CommandId: commandID,
			}, nil
		},
		getCommand: func(context.Context, *apiservice.GetCommandRequest, ...grpc.CallOption) (*apiservice.GetCommandResponse, error) {
			return nil, status.Error(codes.NotFound, "Command is not found")
		},
	}

	_, err := SyncApplication(
		context.Background(),
		cli,
		"app-id",
		time.Nanosecond,
		time.Hour,
		zap.NewNop(),
	)
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
	require.Contains(t, err.Error(), commandID)
}
