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
	"testing"

	"github.com/stretchr/testify/assert"

	model "github.com/pipe-cd/pipecd/pkg/model"
)

func newValidDeployment() *model.Deployment {
	return &model.Deployment{
		Id:              "dep-id",
		ApplicationId:   "app-id",
		ApplicationName: "app-name",
		PipedId:         "piped-id",
		ProjectId:       "project-id",
		GitPath: &model.ApplicationGitPath{
			Repo: &model.ApplicationGitRepository{
				Id: "repo-id",
			},
			Path: "path",
		},
		Trigger: &model.DeploymentTrigger{
			Commit: &model.Commit{
				Hash:      "hash",
				Message:   "msg",
				Author:    "auth",
				Branch:    "branch",
				CreatedAt: 123456,
			},
			Commander: "commander",
			Timestamp: 123456,
		},
	}
}

func newValidPipelineStage() *model.PipelineStage {
	return &model.PipelineStage{
		Id:        "stage-id",
		Name:      "stage-name",
		CreatedAt: 123456,
		UpdatedAt: 123456,
	}
}

func TestDetermineVersionsRequestValidate(t *testing.T) {
	tests := []struct {
		name        string
		req         *DetermineVersionsRequest
		expectedErr string
	}{
		{
			name:        "nil struct is valid",
			req:         nil,
			expectedErr: "",
		},
		{
			name:        "invalid: missing input",
			req:         &DetermineVersionsRequest{},
			expectedErr: "DetermineVersionsRequest.Input",
		},
		{
			name: "valid: input provided",
			req: &DetermineVersionsRequest{
				Input: &PlanPluginInput{
					Deployment: newValidDeployment(),
				},
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			}
		})
	}
}

func TestDetermineStrategyRequestValidate(t *testing.T) {
	tests := []struct {
		name        string
		req         *DetermineStrategyRequest
		expectedErr string
	}{
		{
			name:        "nil struct is valid",
			req:         nil,
			expectedErr: "",
		},
		{
			name:        "invalid: missing input",
			req:         &DetermineStrategyRequest{},
			expectedErr: "DetermineStrategyRequest.Input",
		},
		{
			name: "valid: input provided",
			req: &DetermineStrategyRequest{
				Input: &PlanPluginInput{
					Deployment: newValidDeployment(),
				},
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			}
		})
	}
}

func TestExecuteStageRequestValidate(t *testing.T) {
	tests := []struct {
		name        string
		req         *ExecuteStageRequest
		expectedErr string
	}{
		{
			name:        "nil struct is valid",
			req:         nil,
			expectedErr: "",
		},
		{
			name:        "invalid: missing input",
			req:         &ExecuteStageRequest{},
			expectedErr: "ExecuteStageRequest.Input",
		},
		{
			name: "valid: input provided",
			req: &ExecuteStageRequest{
				Input: &ExecutePluginInput{
					Deployment: newValidDeployment(),
					Stage:      newValidPipelineStage(),
				},
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			}
		})
	}
}

func TestBuildPipelineSyncStagesRequestStageConfigValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      *BuildPipelineSyncStagesRequest_StageConfig
		expectedErr string
	}{
		{
			name:        "nil struct is valid",
			config:      nil,
			expectedErr: "",
		},
		{
			name: "invalid: empty name",
			config: &BuildPipelineSyncStagesRequest_StageConfig{
				Name:  "",
				Index: 0,
			},
			expectedErr: "BuildPipelineSyncStagesRequest_StageConfig.Name",
		},
		{
			name: "invalid: negative index",
			config: &BuildPipelineSyncStagesRequest_StageConfig{
				Name:  "stage-1",
				Index: -1,
			},
			expectedErr: "BuildPipelineSyncStagesRequest_StageConfig.Index",
		},
		{
			name: "valid",
			config: &BuildPipelineSyncStagesRequest_StageConfig{
				Name:  "stage-1",
				Index: 0,
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			}
		})
	}
}

func TestPlanPluginInputValidate(t *testing.T) {
	tests := []struct {
		name        string
		input       *PlanPluginInput
		expectedErr string
	}{
		{
			name:        "nil struct is valid",
			input:       nil,
			expectedErr: "",
		},
		{
			name:        "invalid: missing deployment",
			input:       &PlanPluginInput{},
			expectedErr: "PlanPluginInput.Deployment",
		},
		{
			name: "valid",
			input: &PlanPluginInput{
				Deployment: newValidDeployment(),
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			}
		})
	}
}

func TestExecutePluginInputValidate(t *testing.T) {
	tests := []struct {
		name        string
		input       *ExecutePluginInput
		expectedErr string
	}{
		{
			name:        "nil struct is valid",
			input:       nil,
			expectedErr: "",
		},
		{
			name:        "invalid: missing deployment and stage",
			input:       &ExecutePluginInput{},
			expectedErr: "ExecutePluginInput.Deployment",
		},
		{
			name: "invalid: missing stage only",
			input: &ExecutePluginInput{
				Deployment: newValidDeployment(),
			},
			expectedErr: "ExecutePluginInput.Stage",
		},
		{
			name: "valid",
			input: &ExecutePluginInput{
				Deployment: newValidDeployment(),
				Stage:      newValidPipelineStage(),
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			}
		})
	}
}
