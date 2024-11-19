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

package deployment

import (
	"fmt"
	"slices"
	"time"

	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
)

type Stage string

const (
	// StageK8sSync represents the state where
	// all resources should be synced with the Git state.
	StageK8sSync Stage = "K8S_SYNC"
	// StageK8sPrimaryRollout represents the state where
	// the PRIMARY variant resources has been updated to the new version/configuration.
	StageK8sPrimaryRollout Stage = "K8S_PRIMARY_ROLLOUT"
	// StageK8sCanaryRollout represents the state where
	// the CANARY variant resources has been rolled out with the new version/configuration.
	StageK8sCanaryRollout Stage = "K8S_CANARY_ROLLOUT"
	// StageK8sCanaryClean represents the state where
	// the CANARY variant resources has been cleaned.
	StageK8sCanaryClean Stage = "K8S_CANARY_CLEAN"
	// StageK8sBaselineRollout represents the state where
	// the BASELINE variant resources has been rolled out.
	StageK8sBaselineRollout Stage = "K8S_BASELINE_ROLLOUT"
	// StageK8sBaselineClean represents the state where
	// the BASELINE variant resources has been cleaned.
	StageK8sBaselineClean Stage = "K8S_BASELINE_CLEAN"
	// StageK8sTrafficRouting represents the state where the traffic to application
	// should be splitted as the specified percentage to PRIMARY, CANARY, BASELINE variants.
	StageK8sTrafficRouting Stage = "K8S_TRAFFIC_ROUTING"
	// StageK8sRollback represents the state where all deployed resources should be rollbacked.
	StageK8sRollback Stage = "K8S_ROLLBACK"
)

var AllStages = []Stage{
	StageK8sSync,
	StageK8sPrimaryRollout,
	StageK8sCanaryRollout,
	StageK8sCanaryClean,
	StageK8sBaselineRollout,
	StageK8sBaselineClean,
	StageK8sTrafficRouting,
	StageK8sRollback,
}

func (s Stage) String() string {
	return string(s)
}

const (
	PredefinedStageK8sSync  = "K8sSync"
	PredefinedStageRollback = "K8sRollback"
)

var predefinedStages = map[string]*model.PipelineStage{
	PredefinedStageK8sSync: {
		Id:       PredefinedStageK8sSync,
		Name:     string(StageK8sSync),
		Desc:     "Sync by applying all manifests",
		Rollback: false,
	},
	PredefinedStageRollback: {
		Id:       PredefinedStageRollback,
		Name:     string(StageK8sRollback),
		Desc:     "Rollback the deployment",
		Rollback: true,
	},
}

// GetPredefinedStage finds and returns the predefined stage for the given id.
func GetPredefinedStage(id string) (*model.PipelineStage, bool) {
	stage, ok := predefinedStages[id]
	return stage, ok
}

func buildQuickSyncPipeline(autoRollback bool, now time.Time) []*model.PipelineStage {
	out := make([]*model.PipelineStage, 0, 2)

	stage, _ := GetPredefinedStage(PredefinedStageK8sSync)
	// we copy the predefined stage to avoid modifying the original one.
	out = append(out, &model.PipelineStage{
		Id:        stage.GetId(),
		Name:      stage.GetName(),
		Desc:      stage.GetDesc(),
		Rollback:  stage.GetRollback(),
		Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
		Metadata:  nil,
		CreatedAt: now.Unix(),
		UpdatedAt: now.Unix(),
	},
	)

	if autoRollback {
		s, _ := GetPredefinedStage(PredefinedStageRollback)
		// we copy the predefined stage to avoid modifying the original one.
		out = append(out, &model.PipelineStage{
			Id:        s.GetId(),
			Name:      s.GetName(),
			Desc:      s.GetDesc(),
			Rollback:  s.GetRollback(),
			Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
			CreatedAt: now.Unix(),
			UpdatedAt: now.Unix(),
		})
	}

	return out
}

func buildPipelineStages(stages []*deployment.BuildPipelineSyncStagesRequest_StageConfig, autoRollback bool, now time.Time) []*model.PipelineStage {
	out := make([]*model.PipelineStage, 0, len(stages)+1)

	for _, s := range stages {
		id := s.GetId()
		if id == "" {
			id = fmt.Sprintf("stage-%d", s.GetIndex())
		}
		stage := &model.PipelineStage{
			Id:        id,
			Name:      s.GetName(),
			Desc:      s.GetDesc(),
			Index:     s.GetIndex(),
			Rollback:  false,
			Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
			CreatedAt: now.Unix(),
			UpdatedAt: now.Unix(),
		}
		out = append(out, stage)
	}

	if autoRollback {
		// we set the index of the rollback stage to the minimum index of all stages.
		minIndex := slices.MinFunc(stages, func(a, b *deployment.BuildPipelineSyncStagesRequest_StageConfig) int {
			return int(a.GetIndex() - b.GetIndex())
		}).GetIndex()

		s, _ := GetPredefinedStage(PredefinedStageRollback)
		// we copy the predefined stage to avoid modifying the original one.
		out = append(out, &model.PipelineStage{
			Id:        s.GetId(),
			Name:      s.GetName(),
			Desc:      s.GetDesc(),
			Index:     minIndex,
			Rollback:  s.GetRollback(),
			Status:    model.StageStatus_STAGE_NOT_STARTED_YET,
			CreatedAt: now.Unix(),
			UpdatedAt: now.Unix(),
		})
	}

	return out
}
