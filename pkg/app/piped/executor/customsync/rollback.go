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
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	"github.com/pipe-cd/pipecd/pkg/app/piped/toolregistry"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type rollbackExecutor struct {
	executor.Input

	repoDir string
	appDir  string
}

func (e *rollbackExecutor) Execute(sig executor.StopSignal) model.StageStatus {
	var (
		ctx            = sig.Context()
		originalStatus = e.Stage.Status
		status         model.StageStatus
	)

	switch model.Stage(e.Stage.Name) {
	case model.StageCustomSyncRollback:
		status = e.ensureRollback(ctx)
	default:
		e.LogPersister.Errorf("Unsupported stage %s", e.Stage.Name)
		return model.StageStatus_STAGE_FAILURE
	}

	return executor.DetermineStageStatus(sig.Signal(), originalStatus, status)
}

func (e *rollbackExecutor) ensureRollback(ctx context.Context) model.StageStatus {
	// Not rollback in case this is the first deployment.
	if e.Deployment.RunningCommitHash == "" {
		e.LogPersister.Errorf("Unable to determine the last deployed commit to rollback. It seems this is the first deployment.")
		return model.StageStatus_STAGE_FAILURE
	}

	runningDS, err := e.RunningDSP.Get(ctx, e.LogPersister)
	if err != nil {
		e.LogPersister.Errorf("Failed to prepare running deploy source data (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}
	e.appDir = runningDS.AppDir

	if len(runningDS.GenericApplicationConfig.Pipeline.Stages) > 1 {
		e.LogPersister.Errorf("There are more than one custom sync stages in the running commit.")
		return model.StageStatus_STAGE_FAILURE
	}

	if runningDS.GenericApplicationConfig.Pipeline.Stages[0].Name.String() != string(model.StageCustomSync) {
		e.LogPersister.Errorf("There are no custom sync in the running commit")
		return model.StageStatus_STAGE_FAILURE
	}
	e.LogPersister.Infof("Start rollback for custom sync")

	e.LogPersister.Infof("Prepare external tools...")
	for _, config := range runningDS.GenericApplicationConfig.Pipeline.Stages[0].CustomSyncOptions.ExternalTools {
		e.LogPersister.Infof(fmt.Sprintf("Check %s %s", config.Package, config.Version))
		addedPlugin, installed, err := toolregistry.DefaultRegistry().ExternalTool(ctx, e.appDir, config)
		if addedPlugin {
			e.LogPersister.Infof(fmt.Sprintf("plugin %s has just added", config.Package))
		}
		if installed {
			e.LogPersister.Infof(fmt.Sprintf(" %s %s has just been installed", config.Package, config.Version))
		}
		if err != nil {
			e.LogPersister.Errorf(fmt.Sprintf("unable to set %s %s (%v)", config.Package, config.Version, err))
			continue
		}
		e.LogPersister.Infof(fmt.Sprintf(" %s %s has just been locally set to application directory", config.Package, config.Version))
	}

	return e.executeCommand(runningDS.GenericApplicationConfig.Pipeline.Stages[0])
}

func (e *rollbackExecutor) executeCommand(config config.PipelineStage) model.StageStatus {
	opts := config.CustomSyncOptions
	binDir := toolregistry.DefaultRegistry().GetBinDir()
	pathFromOS := os.Getenv("PATH")

	e.LogPersister.Infof("Runnnig commands...")
	for _, v := range strings.Split(opts.Run, "\n") {
		if v != "" {
			e.LogPersister.Infof("   %s", v)
		}
	}

	path := binDir + ":" + pathFromOS
	envs := make([]string, 0, len(opts.Envs))
	for key, value := range opts.Envs {
		envs = append(envs, key+"="+value)
	}

	cmd := exec.Command("/bin/sh", "-c", opts.Run)
	cmd.Dir = e.appDir
	cmd.Env = append(os.Environ(), append(envs, "PATH="+path)...)
	cmd.Stdout = e.LogPersister
	cmd.Stderr = e.LogPersister
	if err := cmd.Run(); err != nil {
		return model.StageStatus_STAGE_FAILURE
	}
	return model.StageStatus_STAGE_SUCCESS
}
