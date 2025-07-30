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
	"context"
	"errors"
	"slices"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/provider"
)

// Stage names for Terraform plugin.
const (
	// TERRAFORM_PLAN stage executes `terraform plan`.
	stagePlan = "TERRAFORM_PLAN"
	// TERRAFORM_APPLY stage executes `terraform apply`.
	stageApply = "TERRAFORM_APPLY"
	// TERRAFORM_ROLLBACK stage rollbacks by executing `terraform apply` for the previous commit.`
	stageRollback = "TERRAFORM_ROLLBACK"
)

// Plugin implements sdk.DeploymentPlugin for Terraform.
type Plugin struct{}

var _ sdk.DeploymentPlugin[config.Config, config.DeployTargetConfig, config.ApplicationConfigSpec] = (*Plugin)(nil)

// BuildPipelineSyncStages implements sdk.DeploymentPlugin.
func (p *Plugin) BuildPipelineSyncStages(ctx context.Context, _ *config.Config, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	reqStages := input.Request.Stages
	out := make([]sdk.PipelineStage, 0, len(reqStages)+1)

	for _, s := range reqStages {
		out = append(out, sdk.PipelineStage{
			Name:               s.Name,
			Index:              s.Index,
			Rollback:           false,
			Metadata:           make(map[string]string),
			AvailableOperation: sdk.ManualOperationNone,
		})
	}
	if input.Request.Rollback {
		minIndex := slices.MinFunc(reqStages, func(a, b sdk.StageConfig) int { return a.Index - b.Index }).Index
		out = append(out, sdk.PipelineStage{
			Name:               stageRollback,
			Index:              minIndex,
			Rollback:           true,
			Metadata:           make(map[string]string),
			AvailableOperation: sdk.ManualOperationNone,
		})
	}
	return &sdk.BuildPipelineSyncStagesResponse{
		Stages: out,
	}, nil
}

// BuildQuickSyncStages implements sdk.DeploymentPlugin.
func (p *Plugin) BuildQuickSyncStages(ctx context.Context, _ *config.Config, input *sdk.BuildQuickSyncStagesInput) (*sdk.BuildQuickSyncStagesResponse, error) {
	stages := make([]sdk.QuickSyncStage, 0, 2)
	stages = append(stages, sdk.QuickSyncStage{
		Name:               stageApply,
		Description:        "Sync by applying any detected changes",
		Rollback:           false,
		Metadata:           map[string]string{},
		AvailableOperation: sdk.ManualOperationNone,
	})

	if input.Request.Rollback {
		stages = append(stages, sdk.QuickSyncStage{
			Name:               stageRollback,
			Description:        "Rollback by applying the previous Terraform files",
			Rollback:           true,
			Metadata:           map[string]string{},
			AvailableOperation: sdk.ManualOperationNone,
		})
	}
	return &sdk.BuildQuickSyncStagesResponse{
		Stages: stages,
	}, nil
}

// DetermineStrategy implements sdk.DeploymentPlugin.
// It returns (nil, nil) because this plugin does not have specific logic for DetermineStrategy.
func (p *Plugin) DetermineStrategy(ctx context.Context, _ *config.Config, input *sdk.DetermineStrategyInput[config.ApplicationConfigSpec]) (*sdk.DetermineStrategyResponse, error) {
	return nil, nil
}

// DetermineVersions implements sdk.DeploymentPlugin.
func (p *Plugin) DetermineVersions(ctx context.Context, _ *config.Config, input *sdk.DetermineVersionsInput[config.ApplicationConfigSpec]) (*sdk.DetermineVersionsResponse, error) {
	files, err := provider.LoadTerraformFiles(input.Request.DeploymentSource.ApplicationDirectory)
	if err != nil {
		input.Logger.Error("failed to load Terraform files", zap.Error(err))
		return nil, err
	}

	versions, err := provider.FindArtifactVersions(files)
	if err != nil || len(versions) == 0 {
		input.Logger.Warn("unable to determine target versions", zap.Error(err))
		versions = []sdk.ArtifactVersion{{Version: "unknown"}}
	}

	return &sdk.DetermineVersionsResponse{
		Versions: versions,
	}, nil
}

// ExecuteStage implements sdk.DeploymentPlugin.
func (p *Plugin) ExecuteStage(ctx context.Context, _ *config.Config, dts []*sdk.DeployTarget[config.DeployTargetConfig], input *sdk.ExecuteStageInput[config.ApplicationConfigSpec]) (*sdk.ExecuteStageResponse, error) {
	switch input.Request.StageName {
	case stagePlan:
		return &sdk.ExecuteStageResponse{
			Status: p.executePlanStage(ctx, input, dts),
		}, nil
	case stageApply:
		panic("unimplemented")
	case stageRollback:
		panic("unimplemented")
	}
	return nil, errors.New("unimplemented or unsupported stage")
}

// FetchDefinedStages implements sdk.DeploymentPlugin.
func (p *Plugin) FetchDefinedStages() []string {
	return []string{
		stagePlan,
		stageApply,
		stageRollback,
	}
}
