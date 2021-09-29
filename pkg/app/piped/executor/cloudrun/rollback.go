// Copyright 2020 The PipeCD Authors.
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

package cloudrun

import (
	"context"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/cloudrun"
	"github.com/pipe-cd/pipe/pkg/app/piped/executor"
	"github.com/pipe-cd/pipe/pkg/model"
)

type rollbackExecutor struct {
	executor.Input
	client provider.Client
}

func (e *rollbackExecutor) Execute(sig executor.StopSignal) model.StageStatus {
	var (
		ctx            = sig.Context()
		originalStatus = e.Stage.Status
		status         model.StageStatus
	)

	cpName, cpCfg, found := findCloudProvider(&e.Input)
	if !found {
		return model.StageStatus_STAGE_FAILURE
	}

	var err error
	e.client, err = provider.DefaultRegistry().Client(ctx, cpName, cpCfg, e.Logger)
	if err != nil {
		e.LogPersister.Errorf("Unable to create ClourRun client for the provider (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	switch model.Stage(e.Stage.Name) {
	case model.StageRollback:
		status = e.ensureRollback(ctx)

	default:
		e.LogPersister.Errorf("Unsupported stage %s for cloudrun application", e.Stage.Name)
		return model.StageStatus_STAGE_FAILURE
	}

	return executor.DetermineStageStatus(sig.Signal(), originalStatus, status)
}

func (e *rollbackExecutor) ensureRollback(ctx context.Context) model.StageStatus {
	// There is nothing to do if this is the first deployment.
	if e.Deployment.RunningCommitHash == "" {
		e.LogPersister.Errorf("Unable to determine the last deployed commit to rollback. It seems this is the first deployment.")
		return model.StageStatus_STAGE_FAILURE
	}

	runningDS, err := e.RunningDSP.GetReadOnly(ctx, e.LogPersister)
	if err != nil {
		e.LogPersister.Errorf("Failed to prepare running deploy source data (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	deployCfg := runningDS.DeploymentConfig.CloudRunDeploymentSpec
	if deployCfg == nil {
		e.LogPersister.Error("Malformed deployment configuration: missing CloudRunDeploymentSpec")
		return model.StageStatus_STAGE_FAILURE
	}

	sm, ok := loadServiceManifest(&e.Input, deployCfg.Input.ServiceManifestFile, runningDS)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	revision, ok := decideRevisionName(sm, e.Deployment.RunningCommitHash, e.LogPersister)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	traffics := []provider.RevisionTraffic{
		{
			RevisionName: revision,
			Percent:      100,
		},
	}
	if !configureServiceManifest(sm, revision, traffics, e.LogPersister) {
		return model.StageStatus_STAGE_FAILURE
	}

	if !apply(ctx, e.client, sm, e.LogPersister) {
		return model.StageStatus_STAGE_FAILURE
	}

	return model.StageStatus_STAGE_SUCCESS
}
