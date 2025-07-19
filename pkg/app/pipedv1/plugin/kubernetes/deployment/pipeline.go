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
	"encoding/json"
	"slices"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"go.uber.org/zap"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
)

const (
	// StageK8sSync represents the state where
	// all resources should be synced with the Git state.
	StageK8sSync = "K8S_SYNC"
	// StageK8sPrimaryRollout represents the state where
	// the PRIMARY variant resources has been updated to the new version/configuration.
	StageK8sPrimaryRollout = "K8S_PRIMARY_ROLLOUT"
	// StageK8sCanaryRollout represents the state where
	// the CANARY variant resources has been rolled out with the new version/configuration.
	StageK8sCanaryRollout = "K8S_CANARY_ROLLOUT"
	// StageK8sCanaryClean represents the state where
	// the CANARY variant resources has been cleaned.
	StageK8sCanaryClean = "K8S_CANARY_CLEAN"
	// StageK8sBaselineRollout represents the state where
	// the BASELINE variant resources has been rolled out.
	StageK8sBaselineRollout = "K8S_BASELINE_ROLLOUT"
	// StageK8sBaselineClean represents the state where
	// the BASELINE variant resources has been cleaned.
	StageK8sBaselineClean = "K8S_BASELINE_CLEAN"
	// StageK8sTrafficRouting represents the state where the traffic to application
	// should be splitted as the specified percentage to PRIMARY, CANARY, BASELINE variants.
	StageK8sTrafficRouting = "K8S_TRAFFIC_ROUTING"
	// StageK8sRollback represents the state where all deployed resources should be rollbacked.
	StageK8sRollback = "K8S_ROLLBACK"
)

var allStages = []string{
	StageK8sSync,
	StageK8sPrimaryRollout,
	StageK8sCanaryRollout,
	StageK8sCanaryClean,
	StageK8sBaselineRollout,
	StageK8sBaselineClean,
	StageK8sTrafficRouting,
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
func buildPipelineStages(input *sdk.BuildPipelineSyncStagesInput) ([]sdk.PipelineStage, error) {
	stages := input.Request.Stages
	autoRollback := input.Request.Rollback
	logger := input.Logger

	out := make([]sdk.PipelineStage, 0, len(stages)+1)

	for _, s := range stages {
		metadata, err := initialMetadata(s, logger)
		if err != nil {
			return nil, err
		}
		out = append(out, sdk.PipelineStage{
			Name:               s.Name,
			Index:              s.Index,
			Rollback:           false,
			Metadata:           metadata,
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

	return out, nil
}

func initialMetadata(s sdk.StageConfig, logger *zap.Logger) (map[string]string, error) {
	switch s.Name {
	case StageK8sTrafficRouting:
		stageCfg := kubeconfig.K8sTrafficRoutingStageOptions{}
		if err := json.Unmarshal(s.Config, &stageCfg); err != nil {
			logger.Error("failed to unmarshal stage config", zap.Error(err))
			return nil, err
		}
		return map[string]string{
			sdk.MetadataKeyStageDisplay: stageCfg.DisplayString(),
		}, nil
	default:
		return make(map[string]string), nil
	}
}
