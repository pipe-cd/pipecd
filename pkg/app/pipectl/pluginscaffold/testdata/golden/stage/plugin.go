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

package main

import (
	"context"

	"go.uber.org/zap"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

type plugin struct{}

// FetchDefinedStages implements sdk.StagePlugin.
func (p *plugin) FetchDefinedStages() []string {
	return []string{
		"DEMO_WAIT",
	}
}

// BuildPipelineSyncStages implements sdk.StagePlugin.
func (p *plugin) BuildPipelineSyncStages(_ context.Context, _ sdk.ConfigNone, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	stages := make([]sdk.PipelineStage, 0, len(input.Request.Stages))
	for _, rs := range input.Request.Stages {
		stages = append(stages, sdk.PipelineStage{
			Index:              rs.Index,
			Name:               rs.Name,
			Rollback:           false,
			Metadata:           map[string]string{},
			AvailableOperation: sdk.ManualOperationNone,
		})
	}
	return &sdk.BuildPipelineSyncStagesResponse{Stages: stages}, nil
}

// ExecuteStage implements sdk.StagePlugin.
func (p *plugin) ExecuteStage(ctx context.Context, _ sdk.ConfigNone, _ sdk.DeployTargetsNone, input *sdk.ExecuteStageInput[struct{}]) (*sdk.ExecuteStageResponse, error) {
	switch input.Request.StageName {
	case "DEMO_WAIT":
		return &sdk.ExecuteStageResponse{Status: p.executeDemoWait(ctx, input)}, nil
	default:
		input.Logger.Error("unsupported stage", zap.String("stage", input.Request.StageName))
		return &sdk.ExecuteStageResponse{Status: sdk.StageStatusFailure}, nil
	}
}
func (p *plugin) executeDemoWait(_ context.Context, input *sdk.ExecuteStageInput[struct{}]) sdk.StageStatus {
	lp, err := input.Client.StageLogPersister()
	if err != nil {
		input.Logger.Error("no stage log persister", zap.Error(err))
		return sdk.StageStatusFailure
	}
	lp.Info("TODO: implement DEMO_WAIT stage logic")
	return sdk.StageStatusFailure
}
