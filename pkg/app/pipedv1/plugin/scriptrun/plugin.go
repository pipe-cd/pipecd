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
	"os"
	"os/exec"
	"strconv"
	"strings"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

const (
	stageScriptRun         = "SCRIPT_RUN"
	stageScriptRunRollback = "SCRIPT_RUN_ROLLBACK"
)

type ContextInfo struct {
	DeploymentID        string            `json:"deploymentID,omitempty"`
	ApplicationID       string            `json:"applicationID,omitempty"`
	ApplicationName     string            `json:"applicationName,omitempty"`
	TriggeredAt         int64             `json:"triggeredAt,omitempty"`
	TriggeredCommitHash string            `json:"triggeredCommitHash,omitempty"`
	TriggeredCommander  string            `json:"triggeredCommander,omitempty"`
	RepositoryURL       string            `json:"repositoryURL,omitempty"`
	Summary             string            `json:"summary,omitempty"`
	Labels              map[string]string `json:"labels,omitempty"`
	IsRollback          bool              `json:"isRollback,omitempty"`
}
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
	switch input.Request.StageName {
	case stageScriptRun:
		return &sdk.ExecuteStageResponse{
			Status: executeScriptRun(ctx, input.Request, input.Client.LogPersister()),
		}, nil
	case stageScriptRunRollback:
		panic("unimplemented")
	}
	return nil, fmt.Errorf("unsupported stage %s", input.Request.StageName)
}

func executeScriptRun(ctx context.Context, request sdk.ExecuteStageRequest[struct{}], lp sdk.StageLogPersister) sdk.StageStatus {
	lp.Infof("Start executing the script run stage")
	opts, err := decode(request.StageConfig)
	if err != nil {
		lp.Errorf("failed to decode the stage config: %v", err)
		return sdk.StageStatusFailure
	}
	if opts.Run == "" {
		return sdk.StageStatusSuccess
	}
	c := make(chan sdk.StageStatus, 1)
	go func() {
		c <- executeCommand(opts, request, lp)
	}()
	select {
	case result := <-c:
		return result
	case <-ctx.Done():
		lp.Info("ScriptRun cancelled")
		// We can return any status here because the piped handles this case as cancelled by a user,
		// ignoring the result from a plugin.
		return sdk.StageStatusFailure
	}
}
func (p *plugin) FetchDefinedStages() []string {
	return []string{stageScriptRun, stageScriptRunRollback}
}

func executeCommand(opts scriptRunStageOptions, request sdk.ExecuteStageRequest[struct{}], lp sdk.StageLogPersister) sdk.StageStatus {
	lp.Infof("Running commands...")
	for _, v := range strings.Split(opts.Run, "\n") {
		if v != "" {
			lp.Infof("   %s", v)
		}
	}
	ci := ContextInfo{
		DeploymentID:        request.Deployment.ID,
		ApplicationID:       request.Deployment.ApplicationID,
		ApplicationName:     request.Deployment.ApplicationName,
		TriggeredAt:         request.Deployment.CreatedAt,
		TriggeredCommitHash: request.TargetDeploymentSource.CommitHash,
		TriggeredCommander:  request.Deployment.TriggeredBy,
		RepositoryURL:       request.Deployment.RepositoryURL,
		Summary:             request.Deployment.Summary,
		Labels:              request.Deployment.Labels,
		IsRollback:          request.StageName == stageScriptRunRollback,
	}
	ciEnv, err := ci.buildEnv()
	if err != nil {
		lp.Errorf("failed to encode the stage config: %v", err)
		return sdk.StageStatusFailure
	}
	envs := make([]string, 0, len(ciEnv)+len(opts.Env))
	for key, value := range ciEnv {
		envs = append(envs, key+"="+value)
	}

	for key, value := range opts.Env {
		envs = append(envs, key+"="+value)
	}

	cmd := exec.Command("/bin/sh", "-l", "-c", opts.Run)
	cmd.Env = append(os.Environ(), envs...)
	cmd.Dir = request.TargetDeploymentSource.ApplicationDirectory
	cmd.Stdout = lp
	cmd.Stderr = lp
	if err := cmd.Run(); err != nil {
		lp.Errorf("failed to exec command: %w", err)
		return sdk.StageStatusFailure
	} else {
		return sdk.StageStatusSuccess
	}
}

func (ci *ContextInfo) buildEnv() (map[string]string, error) {
	b, err := json.Marshal(ci)
	if err != nil {
		return nil, err
	}
	envs := map[string]string{
		"SR_DEPLOYMENT_ID":         ci.DeploymentID,
		"SR_APPLICATION_ID":        ci.ApplicationID,
		"SR_APPLICATION_NAME":      ci.ApplicationName,
		"SR_TRIGGERED_AT":          strconv.FormatInt(ci.TriggeredAt, 10),
		"SR_TRIGGERED_COMMIT_HASH": ci.TriggeredCommitHash,
		"SR_TRIGGERED_COMMANDER":   ci.TriggeredCommander,
		"SR_REPOSITORY_URL":        ci.RepositoryURL,
		"SR_SUMMARY":               ci.Summary,
		"SR_IS_ROLLBACK":           strconv.FormatBool(ci.IsRollback),
		"SR_CONTEXT_RAW":           string(b), // Add the raw json string as an environment variable.
	}
	for k, v := range ci.Labels {
		eName := "SR_LABELS_" + strings.ToUpper(k)
		envs[eName] = v
	}
	return envs, nil
}
