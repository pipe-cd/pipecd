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
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/executor"
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

	envs := make([]string, 0, len(opts.Env))
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
