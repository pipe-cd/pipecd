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
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"go.uber.org/zap"
)

const (
	stageAnalysis = "ANALYSIS"
)

type analysisPlugin struct{}

// AnalysisStageOptions represents the configuration for ANALYSIS stage
type AnalysisStageOptions struct {
	Timeout  int                      `json:"timeout"`
	Duration time.Duration            `json:"duration"`
	Metrics  []map[string]interface{} `json:"metrics,omitempty"`
	HTTP     []map[string]interface{} `json:"http,omitempty"`
	Logs     []map[string]interface{} `json:"logs,omitempty"`
}

// AnalysisState represents the current state of the analysis
type AnalysisState struct {
	StartTime     time.Time `json:"startTime"`
	LastRunTime   time.Time `json:"lastRunTime"`
	Completed     bool      `json:"completed"`
	Success       bool      `json:"success"`
	FailureReason string    `json:"failureReason,omitempty"`
}

// BuildPipelineSyncStages implements sdk.StagePlugin.
func (p *analysisPlugin) BuildPipelineSyncStages(ctx context.Context, _ sdk.ConfigNone, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	stages := make([]sdk.PipelineStage, 0, len(input.Request.Stages))
	for _, rs := range input.Request.Stages {
		stage := sdk.PipelineStage{
			Index:    rs.Index,
			Name:     rs.Name,
			Rollback: false,
			Metadata: map[string]string{
				sdk.MetadataKeyStageDisplay: "Analysis",
			},
			AvailableOperation: sdk.ManualOperationNone,
		}
		stages = append(stages, stage)
	}

	return &sdk.BuildPipelineSyncStagesResponse{
		Stages: stages,
	}, nil
}

// ExecuteStage implements sdk.StagePlugin.
func (p *analysisPlugin) ExecuteStage(ctx context.Context, _ sdk.ConfigNone, _ sdk.DeployTargetsNone, input *sdk.ExecuteStageInput[struct{}]) (*sdk.ExecuteStageResponse, error) {
	// Parse stage configuration
	var opts AnalysisStageOptions
	if err := json.Unmarshal(input.Request.StageConfig, &opts); err != nil {
		input.Logger.Error("failed to parse analysis stage config", zap.Error(err))
		return &sdk.ExecuteStageResponse{Status: sdk.StageStatusFailure}, nil
	}

	// Create a unique key for this deployment stage
	stateKey := fmt.Sprintf("analysis-%s-%d", input.Request.Deployment.ID, input.Request.StageIndex)

	// Try to get analysis state from shared object store
	var state AnalysisState
	data, err := input.Client.GetApplicationSharedObject(ctx, stateKey)
	if err != nil {
		// Create new state if not found
		state = AnalysisState{
			StartTime: time.Now(),
			Completed: false,
			Success:   false,
		}
	} else {
		if err := json.Unmarshal(data, &state); err != nil {
			input.Logger.Error("failed to unmarshal analysis state", zap.Error(err))
			return &sdk.ExecuteStageResponse{Status: sdk.StageStatusFailure}, nil
		}
	}

	// Check if already completed
	if state.Completed {
		if state.Success {
			return &sdk.ExecuteStageResponse{Status: sdk.StageStatusSuccess}, nil
		}
		return &sdk.ExecuteStageResponse{Status: sdk.StageStatusFailure}, nil
	}

	// Check for stage commands (approvals, skips, etc.) using ListStageCommands
	for stageCommand, err := range input.Client.ListStageCommands(ctx) {
		if err != nil {
			input.Logger.Error("failed to list stage commands", zap.Error(err))
			continue
		}

		input.Logger.Info("processing stage command",
			zap.String("commander", stageCommand.Commander),
			zap.Any("type", stageCommand.Type))

		switch stageCommand.Type {
		case sdk.CommandTypeApproveStage:
			input.Logger.Info("analysis stage approved", zap.String("commander", stageCommand.Commander))
			state.Completed = true
			state.Success = true
			state.LastRunTime = time.Now()

			// Save state and return success
			stateData, err := json.Marshal(state)
			if err != nil {
				input.Logger.Error("failed to marshal analysis state", zap.Error(err))
				return &sdk.ExecuteStageResponse{Status: sdk.StageStatusFailure}, nil
			}
			if err := input.Client.PutApplicationSharedObject(ctx, stateKey, stateData); err != nil {
				input.Logger.Error("failed to save analysis state", zap.Error(err))
				return &sdk.ExecuteStageResponse{Status: sdk.StageStatusFailure}, nil
			}
			return &sdk.ExecuteStageResponse{Status: sdk.StageStatusSuccess}, nil

		case sdk.CommandTypeSkipStage:
			input.Logger.Info("analysis stage skipped", zap.String("commander", stageCommand.Commander))
			state.Completed = true
			state.Success = true // Skip counts as success
			state.LastRunTime = time.Now()

			// Save state and return success
			stateData, err := json.Marshal(state)
			if err != nil {
				input.Logger.Error("failed to marshal analysis state", zap.Error(err))
				return &sdk.ExecuteStageResponse{Status: sdk.StageStatusFailure}, nil
			}
			if err := input.Client.PutApplicationSharedObject(ctx, stateKey, stateData); err != nil {
				input.Logger.Error("failed to save analysis state", zap.Error(err))
				return &sdk.ExecuteStageResponse{Status: sdk.StageStatusFailure}, nil
			}
			return &sdk.ExecuteStageResponse{Status: sdk.StageStatusSuccess}, nil
		}
	}

	// Simulate analysis execution
	input.Logger.Info("executing analysis stage")

	// For now, just mark as successful after some time
	if time.Since(state.StartTime) > 30*time.Second {
		state.Completed = true
		state.Success = true
		state.LastRunTime = time.Now()
	} else {
		state.LastRunTime = time.Now()
	}

	// Save state back to shared object store
	stateData, err := json.Marshal(state)
	if err != nil {
		input.Logger.Error("failed to marshal analysis state", zap.Error(err))
		return &sdk.ExecuteStageResponse{Status: sdk.StageStatusFailure}, nil
	}

	if err := input.Client.PutApplicationSharedObject(ctx, stateKey, stateData); err != nil {
		input.Logger.Error("failed to save analysis state", zap.Error(err))
		return &sdk.ExecuteStageResponse{Status: sdk.StageStatusFailure}, nil
	}

	if state.Completed {
		if state.Success {
			input.Logger.Info("analysis completed successfully")
			return &sdk.ExecuteStageResponse{Status: sdk.StageStatusSuccess}, nil
		}
		input.Logger.Error("analysis failed", zap.String("reason", state.FailureReason))
		return &sdk.ExecuteStageResponse{Status: sdk.StageStatusFailure}, nil
	}

	// Continue running
	input.Logger.Info("analysis still running...")
	return &sdk.ExecuteStageResponse{Status: sdk.StageStatusFailure}, nil // Will retry
}

// FetchDefinedStages implements sdk.StagePlugin.
func (p *analysisPlugin) FetchDefinedStages() []string {
	return []string{stageAnalysis}
}
