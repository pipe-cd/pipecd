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

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/provider"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
	"go.uber.org/zap"
)

const (
	// stageSync executes `terraform plan` and `terraform apply`.
	stageSync = "TERRAFORM_SYNC"
	// stagePlan shows `terraform plan` result.
	stagePlan = "TERRAFORM_PLAN"
	// stageApply executes `terraform apply`.
	stageApply = "TERRAFORM_APPLY"
	// stageRollback rollbacks the deployment by applying terraform files of the previous commit.
	stageRollback = "TERRAFORM_ROLLBACK"
)

func determineVersions(ds sdk.DeploymentSource[config.TerraformApplicationSpec], logger *zap.Logger) ([]sdk.ArtifactVersion, error) {
	files, err := provider.LoadTerraformFiles(ds.ApplicationDirectory)
	if err != nil {
		return nil, err
	}

	if versions, e := provider.FindArtifactVersions(files); e != nil || len(versions) == 0 {
		logger.Warn("failed to determine artifact versions", zap.Error(e))
		return []sdk.ArtifactVersion{
			{
				Kind:    sdk.ArtifactKindUnknown,
				Version: "unknown",
			},
		}, nil
	} else {
		return versions, nil
	}
}

func determineStrategy() (strategy sdk.SyncStrategy, summary string) {
	return sdk.SyncStrategyPipelineSync, "PipelineSync with the specified pipeline"
}

func buildQuickSyncStages(autoRollback bool) []sdk.QuickSyncStage {
	out := make([]sdk.QuickSyncStage, 0, 2)

	out = append(out, sdk.QuickSyncStage{
		Name:               stageSync,
		Description:        "Sync by executing 'terraform plan' and 'terraform apply'",
		Rollback:           false,
		Metadata:           make(map[string]string, 0),
		AvailableOperation: sdk.ManualOperationNone,
	},
	)

	if autoRollback {
		out = append(out, sdk.QuickSyncStage{
			Name:               stageRollback,
			Description:        "Rollback the deployment by applying terraform files of the previous commit",
			Rollback:           true,
			Metadata:           make(map[string]string, 0),
			AvailableOperation: sdk.ManualOperationNone,
		})
	}

	return out
}

func buildPipelineSyncStages(stages []sdk.StageConfig, autoRollback bool) []sdk.PipelineStage {
	out := make([]sdk.PipelineStage, 0, len(stages))
	for _, s := range stages {
		stage := sdk.PipelineStage{
			Index:              s.Index,
			Name:               s.Name,
			Rollback:           false,
			Metadata:           map[string]string{},
			AvailableOperation: sdk.ManualOperationNone,
		}
		out = append(out, stage)
	}

	if autoRollback {
		// we set the index of the rollback stage to the minimum index of all stages.
		minIndex := slices.MinFunc(stages, func(a, b sdk.StageConfig) int {
			return a.Index - b.Index
		}).Index

		out = append(out, sdk.PipelineStage{
			Index:    minIndex,
			Name:     stageRollback,
			Rollback: true,
		})
	}

	return out
}

func fetchDefinedStages() []string {
	return []string{
		stageSync,
		stagePlan,
		stageApply,
		stageRollback,
	}
}
