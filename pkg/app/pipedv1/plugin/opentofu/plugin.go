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
	"bytes"
	"context"
	"os/exec"

	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

// allow exec.CommandContext to be overridden in tests
var execCommand = exec.CommandContext

type plugin struct{}

type config struct{}

type applicationConfigSpec struct{}

const (
	stagePlan    = "OPENTOFU_PLAN"
	stageApply   = "OPENTOFU_APPLY"
	stageDestroy = "OPENTOFU_DESTROY"
)

// BuildPipelineSyncStages implements sdk.StagePlugin.
func (p *plugin) BuildPipelineSyncStages(context.Context, *config, *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	return &sdk.BuildPipelineSyncStagesResponse{}, nil
}

// ExecuteStage implements sdk.StagePlugin.
func (p *plugin) ExecuteStage(ctx context.Context, _ *config, _ sdk.DeployTargetsNone, input *sdk.ExecuteStageInput[applicationConfigSpec]) (*sdk.ExecuteStageResponse, error) {
	stage := input.Request.StageName
	var out bytes.Buffer
	var err error

	switch stage {
	case stagePlan:
		cmd := execCommand(ctx, "tofu", "plan")
		cmd.Stdout = &out
		cmd.Stderr = &out
		err = cmd.Run()
	case stageApply:
		cmd := execCommand(ctx, "tofu", "apply", "-auto-approve")
		cmd.Stdout = &out
		cmd.Stderr = &out
		err = cmd.Run()
	case stageDestroy:
		cmd := execCommand(ctx, "tofu", "destroy", "-auto-approve")
		cmd.Stdout = &out
		cmd.Stderr = &out
		err = cmd.Run()
	default:
		return &sdk.ExecuteStageResponse{Status: sdk.StageStatusFailure}, nil
	}

	if err != nil {
		return &sdk.ExecuteStageResponse{Status: sdk.StageStatusFailure}, nil
	}
	return &sdk.ExecuteStageResponse{Status: sdk.StageStatusSuccess}, nil
}

// FetchDefinedStages implements sdk.StagePlugin.
func (p *plugin) FetchDefinedStages() []string {
	return []string{stagePlan, stageApply, stageDestroy}
}
