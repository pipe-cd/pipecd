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
	"os/exec"
	"testing"

	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

// mockExecCommand is used to mock exec.CommandContext
var mockExecCommand func(ctx context.Context, name string, args ...string) *exec.Cmd

func TestMain(m *testing.M) {
	execCommand = func(ctx context.Context, name string, args ...string) *exec.Cmd {
		if mockExecCommand != nil {
			return mockExecCommand(ctx, name, args...)
		}
		return exec.CommandContext(ctx, name, args...)
	}
	m.Run()
}

// override exec.CommandContext for testing
// var execCommand = exec.CommandContext

func (p *plugin) testableExecuteStage(ctx context.Context, stage string) (*sdk.ExecuteStageResponse, error) {
	input := &sdk.ExecuteStageInput[applicationConfigSpec]{
		Request: sdk.ExecuteStageRequest[applicationConfigSpec]{
			StageName: stage,
		},
	}
	return p.ExecuteStage(ctx, nil, nil, input)
}

func TestExecuteStage_Success(t *testing.T) {
	p := &plugin{}
	mockExecCommand = func(ctx context.Context, name string, args ...string) *exec.Cmd {
		return exec.CommandContext(ctx, "echo", "success")
	}
	defer func() { mockExecCommand = nil }()

	stages := []string{stagePlan, stageApply, stageDestroy}
	for _, stage := range stages {
		resp, err := p.testableExecuteStage(context.Background(), stage)
		if err != nil {
			t.Errorf("unexpected error for stage %s: %v", stage, err)
		}
		if resp.Status != sdk.StageStatusSuccess {
			t.Errorf("expected success for stage %s, got %v", stage, resp.Status)
		}
	}
}

func TestExecuteStage_UnknownStage(t *testing.T) {
	p := &plugin{}
	resp, err := p.testableExecuteStage(context.Background(), "UNKNOWN_STAGE")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp.Status != sdk.StageStatusFailure {
		t.Errorf("expected failure for unknown stage, got %v", resp.Status)
	}
}
