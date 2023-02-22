// Copyright 2023 The PipeCD Authors.
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

package custom

import (
	"os"
	"os/exec"
	"time"

	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	"github.com/pipe-cd/pipecd/pkg/app/piped/toolregistry"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	defaultDuration = 1 * time.Minute
)

type Executor struct {
	executor.Input
}

type registerer interface {
	Register(stage model.Stage, f executor.Factory) error
}

// Register registers this executor factory into a given registerer.
func Register(r registerer) {
	f := func(in executor.Input) executor.Executor {
		return &Executor{
			Input: in,
		}
	}
	r.Register(model.StageCustomStage, f)
}

// Execute starts waiting for the specified duration.
func (e *Executor) Execute(sig executor.StopSignal) model.StageStatus {
	var (
		originalStatus = e.Stage.Status
		duration       = defaultDuration
	)
	c := make(chan model.StageStatus, 1)
	go func() {
		result := e.executeCommand(e.StageConfig.CustomStageOptions)
		c <- result
	}()

	timer := time.NewTimer(duration)
	defer timer.Stop()

	for {
		select {
		case result := <-c:
			return result
		case <-timer.C:
			return model.StageStatus_STAGE_SUCCESS

		case s := <-sig.Ch():
			switch s {
			case executor.StopSignalCancel:
				return model.StageStatus_STAGE_CANCELLED
			case executor.StopSignalTerminate:
				return originalStatus
			default:
				return model.StageStatus_STAGE_FAILURE
			}
		}
	}
}

func (e *Executor) executeCommand(opts *config.CustomStageOptions) model.StageStatus {
	workingDir, err := os.MkdirTemp("", "custom-stage")
	if err != nil {
		e.LogPersister.Errorf("failed to make working directory, %v", err)
		return model.StageStatus_STAGE_FAILURE
	}
	defer os.RemoveAll(workingDir)

	binDir := toolregistry.DefaultRegistry().GetBinDir()
	pathFromOS := os.Getenv("PATH")

	path := binDir + ":" + pathFromOS
	var envs []string
	for key, value := range opts.Env {
		envs = append(envs, key+"="+value)
	}
	for _, v := range opts.Runs {
		cmd := exec.Command("/bin/sh", "-c", v)
		e.LogPersister.Infof("RUN %s (env: %v)", v, envs)
		cmd.Dir = workingDir
		cmd.Env = append(os.Environ(), append(envs, "PATH="+path)...)
		out, err := cmd.CombinedOutput()
		e.LogPersister.Infof("%s", out)
		if err != nil {
			e.LogPersister.Errorf("ERROR %v", err)
			return model.StageStatus_STAGE_FAILURE
		}

	}
	return model.StageStatus_STAGE_SUCCESS
}
