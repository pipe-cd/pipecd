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

	"github.com/example/piped-plugin-demo/config"
)

var ErrUnsupportedStage = errors.New("unsupported stage")

var _ sdk.DeploymentPlugin[config.DemoPluginConfig, config.DemoDeployTargetConfig, config.DemoApplicationSpec] = (*DemoPlugin)(nil)

type DemoPlugin struct{}

func (p *DemoPlugin) FetchDefinedStages() []string {
	return allStages
}

func (p *DemoPlugin) BuildQuickSyncStages(
	_ context.Context,
	_ *config.DemoPluginConfig,
	input *sdk.BuildQuickSyncStagesInput,
) (*sdk.BuildQuickSyncStagesResponse, error) {
	return &sdk.BuildQuickSyncStagesResponse{
		Stages: buildQuickSyncPipeline(input.Request.Rollback),
	}, nil
}

func (p *DemoPlugin) BuildPipelineSyncStages(
	_ context.Context,
	_ *config.DemoPluginConfig,
	input *sdk.BuildPipelineSyncStagesInput,
) (*sdk.BuildPipelineSyncStagesResponse, error) {
	return &sdk.BuildPipelineSyncStagesResponse{
		Stages: buildPipelineStages(input),
	}, nil
}

func (p *DemoPlugin) ExecuteStage(
	ctx context.Context,
	cfg *config.DemoPluginConfig,
	deployTargets []*sdk.DeployTarget[config.DemoDeployTargetConfig],
	input *sdk.ExecuteStageInput[config.DemoApplicationSpec],
) (*sdk.ExecuteStageResponse, error) {
	if len(deployTargets) == 0 {
		input.Logger.Error("no deploy targets")
		return &sdk.ExecuteStageResponse{Status: sdk.StageStatusFailure}, nil
	}
	switch input.Request.StageName {
	case StageDemoSync:
		return &sdk.ExecuteStageResponse{
			Status: p.executeDemoSync(ctx, cfg, deployTargets[0], input),
		}, nil
	case StageDemoRollback:
		return &sdk.ExecuteStageResponse{
			Status: p.executeDemoRollback(ctx, cfg, deployTargets[0], input),
		}, nil
	default:
		return nil, ErrUnsupportedStage
	}
}

func (p *DemoPlugin) DetermineVersions(
	_ context.Context,
	_ *config.DemoPluginConfig,
	_ *sdk.DetermineVersionsInput[config.DemoApplicationSpec],
) (*sdk.DetermineVersionsResponse, error) {
	return &sdk.DetermineVersionsResponse{Versions: []sdk.ArtifactVersion{}}, nil
}

func (p *DemoPlugin) DetermineStrategy(
	_ context.Context,
	_ *config.DemoPluginConfig,
	_ *sdk.DetermineStrategyInput[config.DemoApplicationSpec],
) (*sdk.DetermineStrategyResponse, error) {
	return &sdk.DetermineStrategyResponse{
		Strategy: sdk.SyncStrategyPipelineSync,
		Summary:  "pipeline sync",
	}, nil
}
