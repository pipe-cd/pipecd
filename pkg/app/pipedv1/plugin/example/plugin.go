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

type plugin struct{}

type config struct{}

type applicationConfigSpec struct{}

// BuildPipelineSyncStages implements sdk.StagePlugin.
func (p *plugin) BuildPipelineSyncStages(context.Context, *config, *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	return &sdk.BuildPipelineSyncStagesResponse{}, nil
}

// ExecuteStage implements sdk.StagePlugin.
func (p *plugin) ExecuteStage(context.Context, *config, sdk.DeployTargetsNone, *sdk.ExecuteStageInput[applicationConfigSpec]) (*sdk.ExecuteStageResponse, error) {
	return &sdk.ExecuteStageResponse{}, nil
}

// FetchDefinedStages implements sdk.StagePlugin.
func (p *plugin) FetchDefinedStages() []string {
	return []string{"EXAMPLE_PLAN", "EXAMPLE_APPLY"}
}
