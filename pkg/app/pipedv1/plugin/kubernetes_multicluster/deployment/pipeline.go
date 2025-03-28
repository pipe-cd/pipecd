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
	"slices"

	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

const (
	// StageK8sMultiSync represents the state where
	// all resources should be synced with the Git state.
	StageK8sMultiSync = "K8S_MULTI_SYNC"
	// StageK8sMultiRollback represents the state where all deployed resources should be rollbacked.
	StageK8sMultiRollback = "K8S_MULTI_ROLLBACK"
)

var allStages = []string{
	StageK8sMultiSync,
	StageK8sMultiRollback,
}

const (
	// StageDescriptionK8sMultiSync represents the description of the K8sSync stage.
	StageDescriptionK8sMultiSync = "Sync by applying all manifests"
	// StageDescriptionK8sMultiRollback represents the description of the K8sRollback stage.
	StageDescriptionK8sMultiRollback = "Rollback the deployment"
)

func buildQuickSyncPipeline(autoRollback bool) []sdk.QuickSyncStage {
	out := make([]sdk.QuickSyncStage, 0, 2)

	out = append(out, sdk.QuickSyncStage{
		Name:               StageK8sMultiSync,
		Description:        StageDescriptionK8sMultiSync,
		Rollback:           false,
		Metadata:           make(map[string]string, 0),
		AvailableOperation: sdk.ManualOperationNone,
	},
	)

	if autoRollback {
		out = append(out, sdk.QuickSyncStage{
			Name:               StageK8sMultiRollback,
			Description:        StageDescriptionK8sMultiRollback,
			Rollback:           true,
			Metadata:           make(map[string]string, 0),
			AvailableOperation: sdk.ManualOperationNone,
		})
	}

	return out
}

// buildPipelineStages builds the pipeline stages with the given SDK stages.
func buildPipelineStages(stages []sdk.StageConfig, autoRollback bool) []sdk.PipelineStage {
	out := make([]sdk.PipelineStage, 0, len(stages)+1)

	for _, s := range stages {
		out = append(out, sdk.PipelineStage{
			Name:               s.Name,
			Index:              s.Index,
			Rollback:           false,
			Metadata:           make(map[string]string, 0),
			AvailableOperation: sdk.ManualOperationNone,
		})
	}

	if autoRollback {
		// we set the index of the rollback stage to the minimum index of all stages.
		minIndex := slices.MinFunc(stages, func(a, b sdk.StageConfig) int {
			return a.Index - b.Index
		}).Index

		out = append(out, sdk.PipelineStage{
			Name:               StageK8sMultiRollback,
			Index:              minIndex,
			Rollback:           true,
			Metadata:           make(map[string]string, 0),
			AvailableOperation: sdk.ManualOperationNone,
		})
	}

	return out
}
