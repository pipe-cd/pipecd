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

package ecs

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestMakeServiceResourceState(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		title          string
		status         string
		expectedStatus model.ECSResourceState_HealthStatus
	}{
		{
			title:          "active is healthy",
			status:         "ACTIVE",
			expectedStatus: model.ECSResourceState_HEALTHY,
		},
		{
			title:          "draining is other",
			status:         "DRAINING",
			expectedStatus: model.ECSResourceState_OTHER,
		},
		{
			title:          "inactive is other",
			status:         "INACTIVE",
			expectedStatus: model.ECSResourceState_OTHER,
		},
		{
			title:          "else is unknown",
			status:         "dummy-status",
			expectedStatus: model.ECSResourceState_UNKNOWN,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()
			service := &types.Service{
				Status: &tc.status,
				// Folowing fields are required to avoid nil pointer panic.
				ServiceArn:  aws.String("test-service-arn"),
				ClusterArn:  aws.String("test-cluster-arn"),
				ServiceName: aws.String("test-service-name"),
				CreatedAt:   aws.Time(time.Now()),
			}
			state := makeServiceResourceState(service)

			expected := &model.ECSResourceState{
				Id:                "test-service-arn",
				OwnerIds:          []string{"test-cluster-arn"},
				ParentIds:         []string{"test-cluster-arn"},
				Name:              "test-service-name",
				ApiVersion:        "pipecd.dev/v1beta1",
				Kind:              "Service",
				HealthStatus:      tc.expectedStatus,
				HealthDescription: fmt.Sprintf("Service's status is %s", tc.status),
				CreatedAt:         service.CreatedAt.Unix(),
				UpdatedAt:         service.CreatedAt.Unix(),
			}

			assert.Equal(t, expected, state)
		})
	}
}

func TestMakeTaskSetResourceState(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		title          string
		status         string
		expectedStatus model.ECSResourceState_HealthStatus
	}{
		{
			title:          "primary is healthy",
			status:         "PRIMARY",
			expectedStatus: model.ECSResourceState_HEALTHY,
		},
		{
			title:          "active is healthy",
			status:         "ACTIVE",
			expectedStatus: model.ECSResourceState_HEALTHY,
		},
		{
			title:          "draining is other",
			status:         "DRAINING",
			expectedStatus: model.ECSResourceState_OTHER,
		},
		{
			title:          "else is unknown",
			status:         "dummy-status",
			expectedStatus: model.ECSResourceState_UNKNOWN,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()
			now := time.Now()
			taskSet := &types.TaskSet{
				Status: &tc.status,
				// Folowing fields are required to avoid nil pointer panic.
				TaskSetArn: aws.String("test-task-set-arn"),
				Id:         aws.String("test-task-set-id"),
				ServiceArn: aws.String("test-service-arn"),
				CreatedAt:  aws.Time(now),
				UpdatedAt:  aws.Time(now),
			}
			state := makeTaskSetResourceState(taskSet)

			expected := &model.ECSResourceState{
				Id:                "test-task-set-arn",
				OwnerIds:          []string{"test-service-arn"},
				ParentIds:         []string{"test-service-arn"},
				Name:              "test-task-set-id",
				ApiVersion:        "pipecd.dev/v1beta1",
				Kind:              "TaskSet",
				HealthStatus:      tc.expectedStatus,
				HealthDescription: fmt.Sprintf("TaskSet's status is %s", tc.status),
				CreatedAt:         now.Unix(),
				UpdatedAt:         now.Unix(),
			}

			assert.Equal(t, expected, state)
		})
	}
}

func TestMakeTaskResourceState(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		title          string
		healthStatus   types.HealthStatus
		expectedStatus model.ECSResourceState_HealthStatus
	}{
		{
			title:          "healthy",
			healthStatus:   types.HealthStatusHealthy,
			expectedStatus: model.ECSResourceState_HEALTHY,
		},
		{
			title:          "unhealthy",
			healthStatus:   types.HealthStatusUnhealthy,
			expectedStatus: model.ECSResourceState_OTHER,
		},
		{
			title:          "unknown",
			healthStatus:   types.HealthStatusUnknown,
			expectedStatus: model.ECSResourceState_UNKNOWN,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()
			now := time.Now()
			task := &types.Task{
				HealthStatus: tc.healthStatus,
				// Folowing fields are required to avoid nil pointer panic.
				LastStatus: aws.String("test-last-status"),
				TaskArn:    aws.String("arn:aws:ecs:region:account-id:task/test-cluster/test-task-id"),
				CreatedAt:  aws.Time(now),
			}
			state := makeTaskResourceState(task, "test-cluster-arn")

			expected := &model.ECSResourceState{
				Id:                "arn:aws:ecs:region:account-id:task/test-cluster/test-task-id",
				OwnerIds:          []string{"test-cluster-arn"},
				ParentIds:         []string{"test-cluster-arn"},
				Name:              "test-task-id",
				ApiVersion:        "pipecd.dev/v1beta1",
				Kind:              "Task",
				HealthStatus:      tc.expectedStatus,
				HealthDescription: fmt.Sprintf("Task's last status is test-last-status and the health status is %s", string(tc.healthStatus)),
				CreatedAt:         now.Unix(),
				UpdatedAt:         now.Unix(),
			}

			assert.Equal(t, expected, state)
		})
	}
}
