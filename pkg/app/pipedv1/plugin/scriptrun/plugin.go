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

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

const (
	stageScriptRun         = "SCRIPT_RUN"
	stageScriptRunRollback = "SCRIPT_RUN_ROLLBACK"
)

type plugin struct{}

func (p *plugin) BuildPipelineSyncStages(_ context.Context, _ sdk.ConfigNone, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	stages := make([]sdk.PipelineStage, 0, len(input.Request.Stages)*2)
	for _, rs := range input.Request.Stages {
		stages = append(stages, sdk.PipelineStage{
			Index:              rs.Index,
			Name:               rs.Name,
			Rollback:           false,
			Metadata:           map[string]string{},
			AvailableOperation: sdk.ManualOperationNone,
		})
		if rs.Name != stageScriptRun {
			continue
		}
		opts, err := decode(rs.Config)
		if err != nil {
			return nil, err
		}
		if opts.OnRollback != "" {
			stages = append(stages, sdk.PipelineStage{
				Index:              rs.Index,
				Name:               stageScriptRunRollback,
				Rollback:           true,
				Metadata:           map[string]string{},
				AvailableOperation: sdk.ManualOperationNone,
			})
		}
	}

	return &sdk.BuildPipelineSyncStagesResponse{
		Stages: stages,
	}, nil
}
func (p *plugin) ExecuteStage(ctx context.Context, _ sdk.ConfigNone, _ sdk.DeployTargetsNone, input *sdk.ExecuteStageInput[struct{}]) (*sdk.ExecuteStageResponse, error) {
	//TODO: later
	return &sdk.ExecuteStageResponse{
		Status: sdk.StageStatusSuccess,
	}, nil
}

func (p *plugin) FetchDefinedStages() []string {
	return []string{stageScriptRun, stageScriptRunRollback}
}
