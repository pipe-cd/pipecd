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
	"context"
	"errors"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	ecsconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/config"
)

var _ sdk.DeploymentPlugin[ecsconfig.ECSPluginConfig, ecsconfig.ECSDeployTargetConfig, ecsconfig.ECSApplicationSpec] = (*ECSPlugin)(nil)

var ErrUnsupportedStage = errors.New("unsupported stage")

type ECSPlugin struct {
}

// FetchDefinedStages returns the list of stages that the plugin can execute.
func (p *ECSPlugin) FetchDefinedStages() []string {
	return allStages
}

// BuildQuickSyncStages builds the stages that will be executed during the quick sync process.
func (p *ECSPlugin) BuildQuickSyncStages(
	ctx context.Context,
	cfg *ecsconfig.ECSPluginConfig,
	input *sdk.BuildQuickSyncStagesInput,
) (*sdk.BuildQuickSyncStagesResponse, error) {
	return &sdk.BuildQuickSyncStagesResponse{
		Stages: buildQuickSyncPipeline(input.Request.Rollback),
	}, nil
}

// BuildPipelineSyncStages builds the stages that will be executed by the plugin.
func (p *ECSPlugin) BuildPipelineSyncStages(
	_ context.Context,
	_ *ecsconfig.ECSPluginConfig,
	input *sdk.BuildPipelineSyncStagesInput,
) (*sdk.BuildPipelineSyncStagesResponse, error) {
	return &sdk.BuildPipelineSyncStagesResponse{
		Stages: buildPipelineStages(input),
	}, nil
}

// ExecuteStage executes the given stage.
func (p *ECSPlugin) ExecuteStage(
	ctx context.Context,
	cfg *ecsconfig.ECSPluginConfig,
	deployTargets []*sdk.DeployTarget[ecsconfig.ECSDeployTargetConfig],
	input *sdk.ExecuteStageInput[ecsconfig.ECSApplicationSpec],
) (*sdk.ExecuteStageResponse, error) {
	switch input.Request.StageName {
	case StageECSSync:
		return &sdk.ExecuteStageResponse{
			Status: p.executeECSSyncStage(ctx, input, deployTargets[0]),
		}, nil
	case StageECSRollback:
		return &sdk.ExecuteStageResponse{
			Status: p.executeECSRollbackStage(ctx, input, deployTargets[0]),
		}, nil
	default:
		return nil, ErrUnsupportedStage
	}
}

// DetermineVersions determines the versions of the resources that will be deployed.
func (p *ECSPlugin) DetermineVersions(
	ctx context.Context,
	cfg *ecsconfig.ECSPluginConfig,
	input *sdk.DetermineVersionsInput[ecsconfig.ECSApplicationSpec],
) (*sdk.DetermineVersionsResponse, error) {
	return &sdk.DetermineVersionsResponse{
		// TODO: Implement the logic to determine the versions of the resources that will be deployed.
		// This is just a placeholder
		Versions: []sdk.ArtifactVersion{
			{
				Version: "latest",
				Name:    "ecs-task",
				URL:     "",
			},
		},
	}, nil
}

// DetermineStrategy determines the strategy to deploy the resources.
func (p *ECSPlugin) DetermineStrategy(
	ctx context.Context,
	cfg *ecsconfig.ECSPluginConfig,
	input *sdk.DetermineStrategyInput[ecsconfig.ECSApplicationSpec],
) (*sdk.DetermineStrategyResponse, error) {
	// Use quick sync as the default strategy for ECS deployment.
	return &sdk.DetermineStrategyResponse{
		Strategy: sdk.SyncStrategyQuickSync,
		Summary:  "Use quick sync strategy for ECS deployment (work as ECS_SYNC stage)",
	}, nil
}
