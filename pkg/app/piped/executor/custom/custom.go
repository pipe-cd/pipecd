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
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	"github.com/pipe-cd/pipecd/pkg/app/piped/toolregistry"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	defaultTimeout = 20 * time.Minute
)

type deployExecutor struct {
	executor.Input

	repoDir string
	appDir  string
}

type registerer interface {
	Register(stage model.Stage, f executor.Factory) error
	RegisterCustomStageRollback(f executor.Factory) error
}

// Register registers this executor factory into a given registerer.
func Register(r registerer) {
	f := func(in executor.Input) executor.Executor {
		return &deployExecutor{
			Input: in,
		}
	}
	r.Register(model.StageCustomSync, f)
	r.RegisterCustomStageRollback(func(in executor.Input) executor.Executor {
		return &customStageRollbackExecutor{
			Input: in,
		}
	})
}

// Execute starts waiting for the specified duration.
func (e *deployExecutor) Execute(sig executor.StopSignal) model.StageStatus {
	var (
		originalStatus = e.Stage.Status
		timeout        = defaultTimeout
	)
	ctx := sig.Context()
	ds, err := e.TargetDSP.Get(ctx, e.LogPersister)
	if err != nil {
		e.LogPersister.Errorf("Failed to prepare target deploy source data (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.repoDir = ds.RepoDir
	e.appDir = ds.AppDir
	if e.StageConfig.CustomSyncOptions.Timeout != 0 {
		timeout = e.StageConfig.CustomSyncOptions.Timeout.Duration()
	}
	fmt.Println(timeout)

	c := make(chan bool, 1)
	go func() {
		result := executeCommand(e.appDir, e.StageConfig.CustomSyncOptions, e.LogPersister)
		c <- result
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case result := <-c:
			if result {
				return model.StageStatus_STAGE_SUCCESS
			} else {
				return model.StageStatus_STAGE_FAILURE
			}
		case <-timer.C:
			e.LogPersister.Errorf("Canceled because of timeout")
			return model.StageStatus_STAGE_FAILURE

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

func executeCommand(appDir string, opts *config.CustomSyncOptions, lp executor.LogPersister) bool {

	binDir := toolregistry.DefaultRegistry().GetBinDir()
	pathFromOS := os.Getenv("PATH")

	path := binDir + ":" + pathFromOS
	envs := make([]string, len(opts.Env))
	for key, value := range opts.Env {
		envs = append(envs, key+"="+value)
	}
	for _, v := range opts.Runs {
		cmd := exec.Command("/bin/sh", "-c", v)
		lp.Infof("RUN %s (env: %v)", v, envs)
		cmd.Dir = appDir
		cmd.Env = append(os.Environ(), append(envs, "PATH="+path)...)
		out, err := cmd.CombinedOutput()
		if len(out) != 0 {
			lp.Infof("%s", out)
		}
		if err != nil {
			lp.Errorf("ERROR %v", err)
			return false
		}

	}
	return true
}
