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

	"go.uber.org/zap"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	ecsconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/provider"
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
	case StageECSPrimaryRollout:
		return &sdk.ExecuteStageResponse{
			Status: p.executeECSPrimaryRolloutStage(ctx, input, deployTargets[0]),
		}, nil
	case StageECSCanaryRollout:
		return &sdk.ExecuteStageResponse{
			Status: p.executeECSCanaryRolloutStage(ctx, input, deployTargets[0]),
		}, nil
	case StageECSCanaryClean:
		return &sdk.ExecuteStageResponse{
			Status: p.executeECSCanaryCleanStage(ctx, input, deployTargets[0]),
		}, nil
	case StageECSTrafficRouting:
		return &sdk.ExecuteStageResponse{
			Status: p.executeECSTrafficRouting(ctx, input, deployTargets[0]),
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
	appCfg, err := input.Request.DeploymentSource.AppConfig()
	if err != nil {
		input.Logger.Error("failed to load application config", zap.Error(err))
		return nil, err
	}

	taskDef, err := provider.LoadTaskDefinition(
		input.Request.DeploymentSource.ApplicationDirectory,
		appCfg.Spec.Input.TaskDefinitionFile,
	)
	if err != nil {
		input.Logger.Error("failed to load task definition", zap.Error(err))
		return nil, err
	}

	return &sdk.DetermineVersionsResponse{
		Versions: determineVersions(taskDef),
	}, nil
}

// DetermineStrategy determines the strategy to deploy the resources.
//
// Use PipelineSync if any container image added, removed, or changed.
//
// Use QuickSync if no image difference.
func (p *ECSPlugin) DetermineStrategy(
	ctx context.Context,
	cfg *ecsconfig.ECSPluginConfig,
	input *sdk.DetermineStrategyInput[ecsconfig.ECSApplicationSpec],
) (*sdk.DetermineStrategyResponse, error) {
	targetAppCfg, err := input.Request.TargetDeploymentSource.AppConfig()
	if err != nil {
		input.Logger.Error("failed to load target application config", zap.Error(err))
		return nil, err
	}

	taskDefFile := targetAppCfg.Spec.Input.TaskDefinitionFile

	targetTaskDef, err := provider.LoadTaskDefinition(
		input.Request.TargetDeploymentSource.ApplicationDirectory,
		taskDefFile,
	)
	if err != nil {
		input.Logger.Error("failed to load target task definition", zap.Error(err))
		return nil, err
	}

	if input.Request.RunningDeploymentSource.ApplicationDirectory == "" {
		return &sdk.DetermineStrategyResponse{
			Strategy: sdk.SyncStrategyPipelineSync,
			Summary:  "Sync with the specified pipeline (no running deployment source)",
		}, nil
	}

	runningTaskDef, err := provider.LoadTaskDefinition(
		input.Request.RunningDeploymentSource.ApplicationDirectory,
		taskDefFile,
	)
	if err != nil {
		input.Logger.Warn("failed to load running task definition, falling back to pipeline sync", zap.Error(err))
		return &sdk.DetermineStrategyResponse{
			Strategy: sdk.SyncStrategyPipelineSync,
			Summary:  "Sync with the specified pipeline (unable to load running task definition)",
		}, nil
	}

	return determineStrategy(runningTaskDef, targetTaskDef), nil
}
