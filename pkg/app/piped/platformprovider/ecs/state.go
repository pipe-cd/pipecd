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
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/pipe-cd/pipecd/pkg/model"
)

// MakeServiceResourceStates creates ECSResourceStates of Service, TaskSets, and Tasks.
// `taskSetTasks` is a map of TaskSetArn to Tasks.
func MakeServiceResourceStates(service *types.Service, taskSetTasks map[string][]*types.Task) []*model.ECSResourceState {
	states := []*model.ECSResourceState{}

	states = append(states, makeServiceResourceState(service))

	for _, taskSet := range service.TaskSets {
		states = append(states, makeTaskSetResourceState(&taskSet))

		for _, task := range taskSetTasks[*taskSet.TaskSetArn] {
			states = append(states, makeTaskResourceState(task, *taskSet.TaskSetArn))
		}
	}
	return states
}

func makeServiceResourceState(service *types.Service) *model.ECSResourceState {
	var healthStatus model.ECSResourceState_HealthStatus
	switch *service.Status {
	case "ACTIVE":
		healthStatus = model.ECSResourceState_HEALTHY
	case "DRAINING", "INACTIVE":
		healthStatus = model.ECSResourceState_OTHER
	default:
		healthStatus = model.ECSResourceState_UNKNOWN
	}

	createdAt := service.CreatedAt.Unix()
	return &model.ECSResourceState{
		Id:        *service.ServiceArn,
		OwnerIds:  []string{*service.ClusterArn},
		ParentIds: []string{*service.ClusterArn},
		Name:      *service.ServiceName,

		Kind: "Service",

		HealthStatus:      healthStatus,
		HealthDescription: fmt.Sprintf("Service's status is %s", *service.Status),

		CreatedAt: createdAt,
		// We use CreatedAt for UpdatedAt although Service is not immutable because Service does not have UpdatedAt field.
		UpdatedAt: createdAt,
	}
}

func makeTaskSetResourceState(taskSet *types.TaskSet) *model.ECSResourceState {
	var healthStatus model.ECSResourceState_HealthStatus
	switch *taskSet.Status {
	case "PRIMARY", "ACTIVE":
		healthStatus = model.ECSResourceState_HEALTHY
	case "DRAINING":
		healthStatus = model.ECSResourceState_OTHER
	default:
		healthStatus = model.ECSResourceState_UNKNOWN
	}

	return &model.ECSResourceState{
		Id:        *taskSet.TaskSetArn,
		OwnerIds:  []string{*taskSet.ServiceArn},
		ParentIds: []string{*taskSet.ServiceArn},
		Name:      *taskSet.Id,

		Kind: "TaskSet",

		HealthStatus:      healthStatus,
		HealthDescription: fmt.Sprintf("TaskSet's status is %s", *taskSet.Status),

		CreatedAt: taskSet.CreatedAt.Unix(),
		UpdatedAt: taskSet.UpdatedAt.Unix(),
	}
}

// `parentArn`: Specify taskSet's arn for service tasks, and specify cluster's arn for standalone tasks.
func makeTaskResourceState(task *types.Task, parentArn string) *model.ECSResourceState {
	var healthStatus model.ECSResourceState_HealthStatus
	switch task.HealthStatus {
	case types.HealthStatusHealthy:
		healthStatus = model.ECSResourceState_HEALTHY
	case types.HealthStatusUnhealthy:
		healthStatus = model.ECSResourceState_OTHER
	default:
		healthStatus = model.ECSResourceState_UNKNOWN
	}

	taskArnParts := strings.Split(*task.TaskArn, "/")
	taskId := taskArnParts[len(taskArnParts)-1]

	createdAt := task.CreatedAt.Unix()
	return &model.ECSResourceState{
		Id:        *task.TaskArn,
		OwnerIds:  []string{parentArn},
		ParentIds: []string{parentArn},
		Name:      taskId,

		Kind: "Task",

		HealthStatus:      healthStatus,
		HealthDescription: fmt.Sprintf("Task's last status is %s and the health status is %s", *task.LastStatus, task.HealthStatus),

		CreatedAt: createdAt,
		UpdatedAt: createdAt, // Task is immutable, so updatedAt is the same as createdAt.
	}
}
