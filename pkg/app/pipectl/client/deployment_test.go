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

func TestWaitDeploymentStatusesReturnsErrorWhenCompletedWithOtherStatus(t *testing.T) {
	const deploymentID = "deployment-id"

	testcases := []struct {
		name          string
		statuses      []model.DeploymentStatus
		statusReason  string
		expectedError string
	}{
		{
			name:          "already failed at the first check",
			statuses:      []model.DeploymentStatus{model.DeploymentStatus_DEPLOYMENT_FAILURE},
			statusReason:  "Failed while executing stage stage-1",
			expectedError: `deployment "deployment-id" has already completed with status DEPLOYMENT_FAILURE, so it will never reach any of the waiting statuses (reason: Failed while executing stage stage-1)`,
		},
		{
			name: "cancelled while waiting",
			statuses: []model.DeploymentStatus{
				model.DeploymentStatus_DEPLOYMENT_RUNNING,
				model.DeploymentStatus_DEPLOYMENT_ROLLING_BACK,
				model.DeploymentStatus_DEPLOYMENT_CANCELLED,
			},
			expectedError: `deployment "deployment-id" has already completed with status DEPLOYMENT_CANCELLED, so it will never reach any of the waiting statuses`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			cli := &fakeAPIClient{
				getDeployment: func(context.Context, *apiservice.GetDeploymentRequest, ...grpc.CallOption) (*apiservice.GetDeploymentResponse, error) {
					s := tc.statuses[min(calls, len(tc.statuses)-1)]
					calls++
					return &apiservice.GetDeploymentResponse{
						Deployment: &model.Deployment{
							Id:           deploymentID,
							Status:       s,
							StatusReason: tc.statusReason,
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
				time.Minute,
				zap.NewNop(),
			)
			require.EqualError(t, err, tc.expectedError)
			require.Equal(t, len(tc.statuses), calls)
		})
	}
}

func TestWaitDeploymentStatusesWaitsUntilReachingStatus(t *testing.T) {
	statuses := []model.DeploymentStatus{
		model.DeploymentStatus_DEPLOYMENT_PENDING,
		model.DeploymentStatus_DEPLOYMENT_RUNNING,
		model.DeploymentStatus_DEPLOYMENT_ROLLING_BACK,
		model.DeploymentStatus_DEPLOYMENT_FAILURE,
	}
	calls := 0
	cli := &fakeAPIClient{
		getDeployment: func(context.Context, *apiservice.GetDeploymentRequest, ...grpc.CallOption) (*apiservice.GetDeploymentResponse, error) {
			s := statuses[min(calls, len(statuses)-1)]
			calls++
			return &apiservice.GetDeploymentResponse{
				Deployment: &model.Deployment{
					Id:     "deployment-id",
					Status: s,
				},
			}, nil
		},
	}

	err := WaitDeploymentStatuses(
		context.Background(),
		cli,
		"deployment-id",
		[]model.DeploymentStatus{
			model.DeploymentStatus_DEPLOYMENT_SUCCESS,
			model.DeploymentStatus_DEPLOYMENT_FAILURE,
		},
		time.Millisecond,
		time.Minute,
		zap.NewNop(),
	)
	require.NoError(t, err)
	require.Equal(t, len(statuses), calls)
}
