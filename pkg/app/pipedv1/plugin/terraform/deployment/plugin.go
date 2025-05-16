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

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/terraform/config"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

// Plugin implements the sdk.DeploymentPlugin interface.
type Plugin struct {
}

var _ sdk.DeploymentPlugin[sdk.ConfigNone, config.TerraformDeployTargetConfig, config.TerraformApplicationSpec] = (*Plugin)(nil)

// FetchDefinedStages returns the defined stages for this plugin.
func (p *Plugin) FetchDefinedStages() []string {
	return fetchDefinedStages()
}

// BuildPipelineSyncStages returns the stages for the pipeline sync strategy.
func (p *Plugin) BuildPipelineSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	return &sdk.BuildPipelineSyncStagesResponse{
		Stages: buildPipelineSyncStages(input.Request.Stages, input.Request.Rollback),
	}, nil
}

// ExecuteStage executes the stage.
func (p *Plugin) ExecuteStage(ctx context.Context, _ *sdk.ConfigNone, dts []*sdk.DeployTarget[config.TerraformDeployTargetConfig], input *sdk.ExecuteStageInput[config.TerraformApplicationSpec]) (*sdk.ExecuteStageResponse, error) {
	switch input.Request.StageName {
	case stageSync:
		return &sdk.ExecuteStageResponse{
			Status: executeSyncStage(ctx, input, dts),
		}, nil
	case stagePlan:
		return &sdk.ExecuteStageResponse{
			Status: executePlanStage(ctx, input, dts),
		}, nil
	case stageApply:
		return &sdk.ExecuteStageResponse{
			Status: executeApplyStage(ctx, input, dts),
		}, nil
	case stageRollback:
		return &sdk.ExecuteStageResponse{
			Status: executeRollbackStage(ctx, input, dts),
		}, nil
	default:
		return nil, errors.New("unsupported stage")
	}
}

// DetermineVersions determines the versions of the application.
func (p *Plugin) DetermineVersions(ctx context.Context, _ *sdk.ConfigNone, input *sdk.DetermineVersionsInput[config.TerraformApplicationSpec]) (*sdk.DetermineVersionsResponse, error) {
	versions, err := determineVersions(input.Request.DeploymentSource, input.Logger)
	if err != nil {
		input.Logger.Error("Failed while determining versions", zap.Error(err))
		return nil, err
	}

	return &sdk.DetermineVersionsResponse{
		Versions: versions,
	}, nil
}

// DetermineStrategy determines the strategy for the deployment.
func (p *Plugin) DetermineStrategy(ctx context.Context, _ *sdk.ConfigNone, _ *sdk.DetermineStrategyInput[config.TerraformApplicationSpec]) (*sdk.DetermineStrategyResponse, error) {
	strategy, summary := determineStrategy()
	return &sdk.DetermineStrategyResponse{
		Strategy: strategy,
		Summary:  summary,
	}, nil
}

// BuildQuickSyncStages returns the stages for the quick sync strategy.
func (p *Plugin) BuildQuickSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildQuickSyncStagesInput) (*sdk.BuildQuickSyncStagesResponse, error) {
	return &sdk.BuildQuickSyncStagesResponse{
		Stages: buildQuickSyncStages(input.Request.Rollback),
	}, nil
}
