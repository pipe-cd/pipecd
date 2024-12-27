// Copyright 2024 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package scriptrun

import (
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type registerer interface {
	Register(stage model.Stage, f executor.Factory) error
	RegisterRollback(kind model.RollbackKind, f executor.Factory) error
}

type Executor struct {
	executor.Input

	appDir string
}

func (e *Executor) Execute(sig executor.StopSignal) model.StageStatus {
	e.LogPersister.Infof("Start executing the script run stage")

	opts := e.Input.StageConfig.ScriptRunStageOptions
	if opts == nil {
		e.LogPersister.Error("option for script run stage not found")
		return model.StageStatus_STAGE_FAILURE
	}

	if opts.Run == "" {
		return model.StageStatus_STAGE_SUCCESS
	}

	var originalStatus = e.Stage.Status
	ds, err := e.TargetDSP.Get(sig.Context(), e.LogPersister)
	if err != nil {
		e.LogPersister.Errorf("Failed to prepare target deploy source data (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.appDir = ds.AppDir

	timeout := e.StageConfig.ScriptRunStageOptions.Timeout.Duration()

	c := make(chan model.StageStatus, 1)
	go func() {
		c <- e.executeCommand()
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case result := <-c:
			return result
		case <-timer.C:
			e.LogPersister.Errorf("Canceled because of timeout")
			return model.StageStatus_STAGE_FAILURE

		case s := <-sig.Ch():
			switch s {
			case executor.StopSignalCancel:
				e.LogPersister.Info("Canceled by user")
				return model.StageStatus_STAGE_CANCELLED
			case executor.StopSignalTerminate:
				e.LogPersister.Info("Terminated by system")
				return originalStatus
			default:
				e.LogPersister.Error("Unexpected")
				return model.StageStatus_STAGE_FAILURE
			}
		}
	}
}

func (e *Executor) executeCommand() model.StageStatus {
	opts := e.StageConfig.ScriptRunStageOptions

	e.LogPersister.Infof("Runnnig commands...")
	for _, v := range strings.Split(opts.Run, "\n") {
		if v != "" {
			e.LogPersister.Infof("   %s", v)
		}
	}

	ci := NewContextInfo(e.Deployment, false)
	ciEnv, err := ci.BuildEnv()
	if err != nil {
		e.LogPersister.Errorf("failed to build srcipt run context info: %w", err)
		return model.StageStatus_STAGE_FAILURE
	}

	envs := make([]string, 0, len(ciEnv)+len(opts.Env))
	for key, value := range ciEnv {
		envs = append(envs, key+"="+value)
	}

	for key, value := range opts.Env {
		envs = append(envs, key+"="+value)
	}

	cmd := exec.Command("/bin/sh", "-l", "-c", opts.Run)
	cmd.Dir = e.appDir
	cmd.Env = append(os.Environ(), envs...)
	cmd.Stdout = e.LogPersister
	cmd.Stderr = e.LogPersister
	if err := cmd.Run(); err != nil {
		e.LogPersister.Errorf("failed to exec command: %w", err)
		return model.StageStatus_STAGE_FAILURE
	}
	return model.StageStatus_STAGE_SUCCESS
}

// ContextInfo is the information that will be passed to the script run stage.
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

// NewContextInfo creates a new ContextInfo from the given deployment.
func NewContextInfo(d *model.Deployment, isRollback bool) *ContextInfo {
	return &ContextInfo{
		DeploymentID:        d.Id,
		ApplicationID:       d.ApplicationId,
		ApplicationName:     d.ApplicationName,
		TriggeredAt:         d.Trigger.Timestamp,
		TriggeredCommitHash: d.Trigger.Commit.Hash,
		TriggeredCommander:  d.Trigger.Commander,
		RepositoryURL:       d.GitPath.Repo.Remote,
		Summary:             d.Summary,
		Labels:              d.Labels,
		IsRollback:          isRollback,
	}
}

// BuildEnv builds the environment variables from the context info.
func (src *ContextInfo) BuildEnv() (map[string]string, error) {
	b, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}

	envs := map[string]string{
		"SR_DEPLOYMENT_ID":         src.DeploymentID,
		"SR_APPLICATION_ID":        src.ApplicationID,
		"SR_APPLICATION_NAME":      src.ApplicationName,
		"SR_TRIGGERED_AT":          strconv.FormatInt(src.TriggeredAt, 10),
		"SR_TRIGGERED_COMMIT_HASH": src.TriggeredCommitHash,
		"SR_TRIGGERED_COMMANDER":   src.TriggeredCommander,
		"SR_REPOSITORY_URL":        src.RepositoryURL,
		"SR_SUMMARY":               src.Summary,
		"SR_IS_ROLLBACK":           strconv.FormatBool(src.IsRollback),
		"SR_CONTEXT_RAW":           string(b), // Add the raw json string as an environment variable.
	}

	for k, v := range src.Labels {
		eName := "SR_LABELS_" + strings.ToUpper(k)
		envs[eName] = v
	}

	return envs, nil
}

type RollbackExecutor struct {
	executor.Input
}

func (e *RollbackExecutor) Execute(sig executor.StopSignal) model.StageStatus {
	e.LogPersister.Infof("Unimplement: rollbacking the script run stage")
	return model.StageStatus_STAGE_FAILURE
}

// Register registers this executor factory into a given registerer.
func Register(r registerer) {
	r.Register(model.StageScriptRun, func(in executor.Input) executor.Executor {
		return &Executor{
			Input: in,
		}
	})
}
