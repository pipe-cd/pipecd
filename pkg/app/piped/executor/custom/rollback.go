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
	"context"

	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type customStageRollbackExecutor struct {
	executor.Input

	repoDir string
	appDir  string
}

func (e *customStageRollbackExecutor) Execute(sig executor.StopSignal) model.StageStatus {
	var (
		ctx            = sig.Context()
		originalStatus = e.Stage.Status
		status         model.StageStatus
	)

	switch model.Stage(e.Stage.Name) {
	case model.StageCustomStageRollback:
		status = e.ensureRollback(ctx)
	default:
		e.LogPersister.Errorf("Unsupported stage %s", e.Stage.Name)
		return model.StageStatus_STAGE_FAILURE
	}

	return executor.DetermineStageStatus(sig.Signal(), originalStatus, status)
}

func (e *customStageRollbackExecutor) ensureRollback(ctx context.Context) model.StageStatus {
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
	e.repoDir = runningDS.RepoDir
	e.appDir = runningDS.AppDir
	e.LogPersister.Infof("Start rollback for custom stages")

	customStageConfigs, ok := runningDS.GenericApplicationConfig.GetStagesFromName(model.StageCustomStage)
	if !ok || customStageConfigs == nil {
		e.LogPersister.Errorf("There are no custom stages in running commit")
	}
	if len(customStageConfigs) > 1 {
		e.LogPersister.Errorf("There are custom stages more than one stage.")
	}
	for _, customStageConfig := range customStageConfigs {
		e.LogPersister.Infof("Start rollback for custom stage (Name: %s Id: %s Desc: %s)", customStageConfig.Name, customStageConfig.Id, customStageConfig.Desc)

		if !ok {
			e.LogPersister.Errorf("Failed to get custom stage config")
			return model.StageStatus_STAGE_FAILURE
		}
		result := executeCommand(e.appDir, customStageConfig.CustomStageOptions, e.LogPersister)
		if !result {
			e.LogPersister.Errorf("Failed to execute command")
			return model.StageStatus_STAGE_FAILURE
		}
	}
	return model.StageStatus_STAGE_SUCCESS
}
