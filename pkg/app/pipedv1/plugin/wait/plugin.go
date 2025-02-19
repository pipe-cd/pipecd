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

	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

const (
	stageWait string = "WAIT"
)

type plugin struct{}

type config struct{}

// Name implements sdk.Plugin.
func (p *plugin) Name() string {
	return "wait"
}

// Version implements sdk.Plugin.
func (p *plugin) Version() string {
	return "0.0.1" // TODO
}

// BuildPipelineSyncStages implements sdk.PipelineSyncPlugin.
func (p *plugin) BuildPipelineSyncStages(ctx context.Context, cfg *config, client *sdk.Client, request *sdk.BuildPipelineSyncStagesRequest) (*sdk.BuildPipelineSyncStagesResponse, error) {
	stages := make([]sdk.PipelineStage, 0, len(request.Stages))
	for _, rs := range request.Stages {
		stage := sdk.PipelineStage{
			Index:              rs.Index,
			Name:               rs.Name,
			Rollback:           false,
			Metadata:           map[string]string{},
			AvailableOperation: sdk.ManualOperationNone,
		}
		stages = append(stages, stage)
	}

	return &sdk.BuildPipelineSyncStagesResponse{
		Stages: stages,
	}, nil
}

// ExecuteStage implements sdk.PipelineSyncPlugin.
func (p *plugin) ExecuteStage(ctx context.Context, cfg *config, targets sdk.DeployTargetsNone, client *sdk.Client, request *sdk.ExecuteStageRequest) (*sdk.ExecuteStageResponse, error) {
	return &sdk.ExecuteStageResponse{}, nil
}

// FetchDefinedStages implements sdk.PipelineSyncPlugin.
func (p *plugin) FetchDefinedStages() []string {
	return []string{stageWait}
}
