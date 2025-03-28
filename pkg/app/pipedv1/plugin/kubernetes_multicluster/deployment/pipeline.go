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
	// StageK8sSync represents the state where
	// all resources should be synced with the Git state.
	StageK8sSync = "K8S_SYNC"
	// StageK8sRollback represents the state where all deployed resources should be rollbacked.
	StageK8sRollback = "K8S_ROLLBACK"
)

var allStages = []string{
	StageK8sSync,
	StageK8sRollback,
}

const (
	// StageDescriptionK8sSync represents the description of the K8sSync stage.
	StageDescriptionK8sSync = "Sync by applying all manifests"
	// StageDescriptionK8sRollback represents the description of the K8sRollback stage.
	StageDescriptionK8sRollback = "Rollback the deployment"
)

func buildQuickSyncPipeline(autoRollback bool) []sdk.QuickSyncStage {
	out := make([]sdk.QuickSyncStage, 0, 2)

	out = append(out, sdk.QuickSyncStage{
		Name:               StageK8sSync,
		Description:        StageDescriptionK8sSync,
		Rollback:           false,
		Metadata:           make(map[string]string, 0),
		AvailableOperation: sdk.ManualOperationNone,
	},
	)

	if autoRollback {
		out = append(out, sdk.QuickSyncStage{
			Name:               StageK8sRollback,
			Description:        StageDescriptionK8sRollback,
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
			Name:               StageK8sRollback,
			Index:              minIndex,
			Rollback:           true,
			Metadata:           make(map[string]string, 0),
			AvailableOperation: sdk.ManualOperationNone,
		})
	}

	return out
}
