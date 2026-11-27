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
	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestWaitDeploymentStatusesSucceeds(t *testing.T) {
	const deploymentID = "deployment-id"

	cli := &fakeAPIClient{
		getDeployment: func(context.Context, *apiservice.GetDeploymentRequest, ...grpc.CallOption) (*apiservice.GetDeploymentResponse, error) {
			return &apiservice.GetDeploymentResponse{
				Deployment: &model.Deployment{
					Id:     deploymentID,
					Status: model.DeploymentStatus_DEPLOYMENT_SUCCESS,
				},
			}, nil
		},
	}

	err := WaitDeploymentStatuses(
		context.Background(),
		cli,
		deploymentID,
		[]model.DeploymentStatus{model.DeploymentStatus_DEPLOYMENT_SUCCESS},
		time.Hour,
		time.Hour,
		zap.NewNop(),
	)
	require.NoError(t, err)
}

func TestWaitDeploymentStatusesSucceedsAfterStatusTransition(t *testing.T) {
	const deploymentID = "deployment-id"
	calls := 0

	cli := &fakeAPIClient{
		getDeployment: func(context.Context, *apiservice.GetDeploymentRequest, ...grpc.CallOption) (*apiservice.GetDeploymentResponse, error) {
			calls++
			s := model.DeploymentStatus_DEPLOYMENT_RUNNING
			if calls >= 3 {
				s = model.DeploymentStatus_DEPLOYMENT_SUCCESS
			}
			return &apiservice.GetDeploymentResponse{
				Deployment: &model.Deployment{
					Id:     deploymentID,
					Status: s,
				},
			}, nil
		},
	}

	err := WaitDeploymentStatuses(
		context.Background(),
		cli,
		deploymentID,
		[]model.DeploymentStatus{model.DeploymentStatus_DEPLOYMENT_SUCCESS},
		time.Millisecond,
		time.Second,
		zap.NewNop(),
	)
	require.NoError(t, err)
	require.GreaterOrEqual(t, calls, 3)
}

func TestWaitDeploymentStatusesRejectsInvalidCheckInterval(t *testing.T) {
	tests := []struct {
		name          string
		checkInterval time.Duration
	}{
		{name: "zero", checkInterval: 0},
		{name: "negative", checkInterval: -time.Second},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			cli := &fakeAPIClient{
				getDeployment: func(context.Context, *apiservice.GetDeploymentRequest, ...grpc.CallOption) (*apiservice.GetDeploymentResponse, error) {
					calls++
					return &apiservice.GetDeploymentResponse{
						Deployment: &model.Deployment{Status: model.DeploymentStatus_DEPLOYMENT_SUCCESS},
					}, nil
				},
			}
			err := WaitDeploymentStatuses(
				context.Background(),
				cli,
				"deployment-id",
				[]model.DeploymentStatus{model.DeploymentStatus_DEPLOYMENT_SUCCESS},
				tc.checkInterval,
				time.Second,
				zap.NewNop(),
			)
			require.Error(t, err)
			require.Contains(t, err.Error(), "check-interval")
			require.Zero(t, calls)
		})
	}
}

func TestWaitDeploymentStatusesReturnsNotFound(t *testing.T) {
	const deploymentID = "missing-deployment"

	cli := &fakeAPIClient{
		getDeployment: func(context.Context, *apiservice.GetDeploymentRequest, ...grpc.CallOption) (*apiservice.GetDeploymentResponse, error) {
			return nil, status.Error(codes.NotFound, "Deployment is not found")
		},
	}

	err := WaitDeploymentStatuses(
		context.Background(),
		cli,
		deploymentID,
		[]model.DeploymentStatus{model.DeploymentStatus_DEPLOYMENT_SUCCESS},
		time.Hour,
		time.Hour,
		zap.NewNop(),
	)
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
	require.Contains(t, err.Error(), deploymentID)
}

func TestWaitDeploymentStatusesRetriesUnavailableUntilTimeout(t *testing.T) {
	calls := 0
	cli := &fakeAPIClient{
		getDeployment: func(context.Context, *apiservice.GetDeploymentRequest, ...grpc.CallOption) (*apiservice.GetDeploymentResponse, error) {
			calls++
			return nil, status.Error(codes.Unavailable, "service unavailable")
		},
	}

	err := WaitDeploymentStatuses(
		context.Background(),
		cli,
		"deployment-id",
		[]model.DeploymentStatus{model.DeploymentStatus_DEPLOYMENT_SUCCESS},
		time.Millisecond,
		50*time.Millisecond,
		zap.NewNop(),
	)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.GreaterOrEqual(t, calls, 2)
}
