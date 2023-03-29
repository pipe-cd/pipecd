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

package customsync

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	"github.com/pipe-cd/pipecd/pkg/app/piped/toolregistry"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type deployExecutor struct {
	executor.Input

	repoDir string
	appDir  string
}

type registerer interface {
	Register(stage model.Stage, f executor.Factory) error
	RegisterRollback(kind model.RollbackKind, f executor.Factory) error
}

// Register registers this executor factory into a given registerer.
func Register(r registerer) {
	f := func(in executor.Input) executor.Executor {
		return &deployExecutor{
			Input: in,
		}
	}
	r.Register(model.StageCustomSync, f)
	r.RegisterRollback(model.RollbackKind_Rollback_CUSTOM_SYNC, func(in executor.Input) executor.Executor {
		return &rollbackExecutor{
			Input: in,
		}
	})
}

// Execute exec the user-defined scripts in timeout duration.
func (e *deployExecutor) Execute(sig executor.StopSignal) model.StageStatus {
	var originalStatus = e.Stage.Status
	ctx := sig.Context()
	ds, err := e.TargetDSP.Get(ctx, e.LogPersister)
	if err != nil {
		e.LogPersister.Errorf("Failed to prepare target deploy source data (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.repoDir = ds.RepoDir
	e.appDir = ds.AppDir

	e.LogPersister.Infof("Prepare external tools...")
	for _, config := range e.StageConfig.CustomSyncOptions.ExternalTools {
		e.LogPersister.Infof(fmt.Sprintf("Check %s %s", config.Package, config.Version))
		addedPlugin, installed, err := toolregistry.DefaultRegistry().ExternalTool(ctx, e.appDir, config)
		if addedPlugin {
			e.LogPersister.Infof(fmt.Sprintf(" plugin %s has just been added", config.Package))
		}
		if installed {
			e.LogPersister.Infof(fmt.Sprintf(" %s %s has just been installed", config.Package, config.Version))
		}
		if err != nil {
			e.LogPersister.Errorf(fmt.Sprintf(" unable to prepare %s %s (%v)", config.Package, config.Version, err))
			continue
		}
		e.LogPersister.Infof(fmt.Sprintf(" %s %s has just been locally set to application directory", config.Package, config.Version))
	}

	timeout := e.StageConfig.CustomSyncOptions.Timeout.Duration()

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
				return model.StageStatus_STAGE_CANCELLED
			case executor.StopSignalTerminate:
				return originalStatus
			default:
				return model.StageStatus_STAGE_FAILURE
			}
		}
	}
}

func (e *deployExecutor) executeCommand() model.StageStatus {
	opts := e.StageConfig.CustomSyncOptions

	e.LogPersister.Infof("Runnnig commands...")
	for _, v := range strings.Split(opts.Run, "\n") {
		if v != "" {
			e.LogPersister.Infof("   %s", v)
		}
	}

	envs := make([]string, 0, len(opts.Envs))
	for key, value := range opts.Envs {
		envs = append(envs, key+"="+value)
	}

	cmd := exec.Command("/bin/sh", "-c", opts.Run)
	cmd.Dir = e.appDir
	cmd.Env = append(os.Environ(), envs...)
	cmd.Stdout = e.LogPersister
	cmd.Stderr = e.LogPersister
	if err := cmd.Run(); err != nil {
		return model.StageStatus_STAGE_FAILURE
	}
	return model.StageStatus_STAGE_SUCCESS
}
