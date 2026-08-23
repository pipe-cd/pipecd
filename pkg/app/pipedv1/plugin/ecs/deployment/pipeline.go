// Copyright 2026 The PipeCD Authors.
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

package deployment

import (
	"errors"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

// ErrRollbackRequiresStages is returned when a rollback is requested but no
// stage was provided in the request.
var ErrRollbackRequiresStages = errors.New("rollback requires at least one stage")

const (
	// StageECSSync represents the ECS sync stage.
	StageECSSync = "ECS_SYNC"

	// StageECSPrimaryRollout represents the ECS primary rollout stage.
	StageECSPrimaryRollout = "ECS_PRIMARY_ROLLOUT"

	// StageECSCanaryRollout represents the ECS canary rollout stage.
	StageECSCanaryRollout = "ECS_CANARY_ROLLOUT"

	// StageECSCanaryClean represents the ECS canary clean stage.
	StageECSCanaryClean = "ECS_CANARY_CLEAN"

	// StageECSTrafficRouting represents the ECS traffic routing stage.
	StageECSTrafficRouting = "ECS_TRAFFIC_ROUTING"

	// StageECSRollback represents the ECS rollback stage.
	StageECSRollback = "ECS_ROLLBACK"
)

var allStages = []string{
	StageECSSync,
	StageECSPrimaryRollout,
	StageECSCanaryRollout,
	StageECSCanaryClean,
	StageECSTrafficRouting,
	StageECSRollback,
}

const (
	StageECSSyncDescription           = "Sync ECS service with given task definition"
	StageECSPrimaryRolloutDescription = "Roll out new task set as primary"
	StageECSCanaryRolloutDescription  = "Roll out new task set as canary"
	StageECSCanaryCleanDescription    = "Clean up canary task set"
	StageECSTrafficRoutingDescription = "Route traffic between primary and canary task sets"
	StageECSRollbackDescription       = "Rollback to previous task set"
)

func buildQuickSyncPipeline(autoRollback bool) []sdk.QuickSyncStage {
	out := []sdk.QuickSyncStage{
		{
			Name:        StageECSSync,
			Description: StageECSSyncDescription,
			Rollback:    false,
		},
	}

	if autoRollback {
		out = append(out, sdk.QuickSyncStage{
			Name:        StageECSRollback,
			Description: StageECSRollbackDescription,
			Rollback:    true,
		})
	}

	return out
}

func buildPipelineStages(input *sdk.BuildPipelineSyncStagesInput) ([]sdk.PipelineStage, error) {
	stages := input.Request.Stages

	out := make([]sdk.PipelineStage, 0, len(stages)+1)

	for _, s := range stages {
		out = append(out, sdk.PipelineStage{
			Name:     s.Name,
			Index:    s.Index,
			Rollback: false,
		})
	}

	if input.Request.Rollback {
		// Guard against a rollback request without any stage, which would
		// panic on stages[0] below.
		if len(stages) == 0 {
			return nil, ErrRollbackRequiresStages
		}
		// The rollback stage must reuse one of the requested indexes and
		// piped runs rollback stages as a trail sorted by index.
		// Use the smallest requested index so our rollback runs first among all plugins' rollbacks.
		minIndex := stages[0].Index
		for _, s := range stages[1:] {
			if s.Index < minIndex {
				minIndex = s.Index
			}
		}
		out = append(out, sdk.PipelineStage{
			Name:     StageECSRollback,
			Index:    minIndex,
			Rollback: true,
		})
	}

	return out, nil
}
