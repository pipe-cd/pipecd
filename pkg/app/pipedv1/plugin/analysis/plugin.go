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

package main

import (
	"context"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/analysis/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/analysis/executestage"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

const (
	stageAnalysis = "ANALYSIS"
)

type plugin struct{}

var _ sdk.StagePlugin[config.PluginConfig, struct{}, struct{}] = (*plugin)(nil)

func (p *plugin) BuildPipelineSyncStages(_ context.Context, _ *config.PluginConfig, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	stages := make([]sdk.PipelineStage, 0, len(input.Request.Stages))
	for _, rs := range input.Request.Stages {
		stages = append(stages, sdk.PipelineStage{
			Index:              rs.Index,
			Name:               rs.Name,
			Rollback:           false,
			Metadata:           map[string]string{},
			AvailableOperation: sdk.ManualOperationSkip,
		})
	}

	return &sdk.BuildPipelineSyncStagesResponse{
		Stages: stages,
	}, nil
}

func (p *plugin) ExecuteStage(ctx context.Context, pluginCfg *config.PluginConfig, _ sdk.DeployTargetsNone, input *sdk.ExecuteStageInput[struct{}]) (*sdk.ExecuteStageResponse, error) {
	return &sdk.ExecuteStageResponse{
		Status: executestage.ExecuteAnalysisStage(ctx, input, pluginCfg),
	}, nil
}

func (p *plugin) FetchDefinedStages() []string {
	return []string{stageAnalysis}
}
