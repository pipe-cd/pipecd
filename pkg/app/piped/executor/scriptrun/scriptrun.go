// Copyright 2023 The PipeCD Authors.
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

	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type registerer interface {
	Register(stage model.Stage, f executor.Factory) error
	RegisterRollback(kind model.RollbackKind, f executor.Factory) error
}

type Executor struct {
	executor.Input
}

func (e *Executor) Execute(sig executor.StopSignal) model.StageStatus {
	e.LogPersister.Infof("Start executing the script run stage")

	opts := e.Input.StageConfig.ScriptRunStageOptions
	if opts == nil {
		e.LogPersister.Infof("option for script run stage not found")
		return model.StageStatus_STAGE_FAILURE
	}

	if opts.Run == "" {
		return model.StageStatus_STAGE_SUCCESS
	}

	envs := make([]string, 0, len(opts.Env))
	for key, value := range opts.Env {
		envs = append(envs, key+"="+value)
	}

	cmd := exec.Command("/bin/sh", "-l", "-c", opts.Run)
	cmd.Env = append(os.Environ(), envs...)
	cmd.Stdout = e.LogPersister
	cmd.Stderr = e.LogPersister

	e.LogPersister.Infof("executing script:")
	e.LogPersister.Infof(opts.Run)

	if err := cmd.Run(); err != nil {
		e.LogPersister.Errorf("failed to execute script: %w", err)
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
