// Copyright 2025 The PipeCD Authors.
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
	"slices"

	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

type Stage string

const (
	// StageK8sSync represents the state where
	// all resources should be synced with the Git state.
	StageK8sMultiSync Stage = "K8S_MULTI_SYNC"
	// StageK8sRollback represents the state where all deployed resources should be rollbacked.
	StageK8sMultiRollback Stage = "K8S_MULTI_ROLLBACK"
)

var AllStages = []Stage{
	StageK8sMultiSync,
	StageK8sMultiRollback,
}

func (s Stage) String() string {
	return string(s)
}

const (
	PredefinedStageK8sMultiSync     = "K8sMultiSync"
	PredefinedStageK8sMultiRollback = "K8sMultiRollback"
)

var predefinedStages = map[string]*model.PipelineStage{
	PredefinedStageK8sMultiSync: {
		Id:       PredefinedStageK8sMultiSync,
		Name:     string(StageK8sMultiSync),
		Desc:     "Sync by applying all manifests",
		Rollback: false,
	},
	PredefinedStageK8sMultiRollback: {
		Id:       PredefinedStageK8sMultiRollback,
		Name:     string(StageK8sMultiRollback),
		Desc:     "Rollback the deployment",
		Rollback: true,
	},
}

// GetPredefinedStage finds and returns the predefined stage for the given id.
func GetPredefinedStage(id string) (*model.PipelineStage, bool) {
	stage, ok := predefinedStages[id]
	return stage, ok
}

func BuildPipelineStages(input *sdk.BuildPipelineSyncStagesInput) []sdk.PipelineStage {
	out := make([]sdk.PipelineStage, 0, len(input.Request.Stages)+1)

	for _, s := range input.Request.Stages {
		stage := sdk.PipelineStage{
			Index:              s.Index,
			Name:               s.Name,
			Rollback:           false,
			Metadata:           make(map[string]string, 0),
			AvailableOperation: sdk.ManualOperationNone,
		}
		out = append(out, stage)
	}

	if input.Request.Rollback {
		// we set the index of the rollback stage to the minimum index of all stages.
		minIndex := slices.MinFunc(out, func(a, b sdk.PipelineStage) int {
			return a.Index - b.Index
		}).Index

		s, _ := GetPredefinedStage(PredefinedStageK8sMultiRollback)
		// we copy the predefined stage to avoid modifying the original one.
		out = append(out, sdk.PipelineStage{
			Index:              minIndex,
			Name:               s.Name,
			Rollback:           true,
			Metadata:           make(map[string]string, 0),
			AvailableOperation: sdk.ManualOperationNone,
		})
	}

	return out
}
