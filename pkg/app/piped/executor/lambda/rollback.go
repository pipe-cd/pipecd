// Copyright 2021 The PipeCD Authors.
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

package lambda

import (
	"context"
	"encoding/json"
	"fmt"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/lambda"
	"github.com/pipe-cd/pipe/pkg/app/piped/executor"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

type rollbackExecutor struct {
	executor.Input
}

func (e *rollbackExecutor) Execute(sig executor.StopSignal) model.StageStatus {
	var (
		ctx            = sig.Context()
		originalStatus = e.Stage.Status
		status         model.StageStatus
	)

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
	// Not rollback in case this is the first deployment.
	if e.Deployment.RunningCommitHash == "" {
		e.LogPersister.Errorf("Unable to determine the last deployed commit to rollback. It seems this is the first deployment.")
		return model.StageStatus_STAGE_FAILURE
	}

	runningDS, err := e.RunningDSP.GetReadOnly(ctx, e.LogPersister)
	if err != nil {
		e.LogPersister.Errorf("Failed to prepare running deploy source data (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	deployCfg := runningDS.DeploymentConfig.LambdaDeploymentSpec
	if deployCfg == nil {
		e.LogPersister.Errorf("Malformed deployment configuration: missing LambdaDeploymentSpec")
		return model.StageStatus_STAGE_FAILURE
	}

	cloudProviderName, cloudProviderCfg, found := findCloudProvider(&e.Input)
	if !found {
		return model.StageStatus_STAGE_FAILURE
	}

	fm, ok := loadFunctionManifest(&e.Input, deployCfg.Input.FunctionManifestFile, runningDS)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	if !rollback(ctx, &e.Input, cloudProviderName, cloudProviderCfg, fm) {
		return model.StageStatus_STAGE_FAILURE
	}

	return model.StageStatus_STAGE_SUCCESS
}

func rollback(ctx context.Context, in *executor.Input, cloudProviderName string, cloudProviderCfg *config.CloudProviderLambdaConfig, fm provider.FunctionManifest) bool {
	in.LogPersister.Infof("Start rollback the lambda function: %s to original stage", fm.Spec.Name)
	client, err := provider.DefaultRegistry().Client(cloudProviderName, cloudProviderCfg, in.Logger)
	if err != nil {
		in.LogPersister.Errorf("Unable to create Lambda client for the provider %s: %v", cloudProviderName, err)
		return false
	}

	originalTrafficKeyName := fmt.Sprintf("%s-%s-original", fm.Spec.Name, in.Deployment.RunningCommitHash)
	originalTrafficCfgData, ok := in.MetadataStore.Get(originalTrafficKeyName)
	if !ok {
		in.LogPersister.Errorf("Unable to prepare original traffic config to rollback Lambda function %s: not found", fm.Spec.Name)
		return false
	}

	originalTrafficCfg := provider.RoutingTrafficConfig{}
	if err := json.Unmarshal([]byte(originalTrafficCfgData), &originalTrafficCfg); err != nil {
		in.LogPersister.Errorf("Unable to prepare original traffic config to rollback Lambda function %s: %v", fm.Spec.Name, err)
		return false
	}

	// TODO: fix case == 1
	if len(originalTrafficCfg) != 2 {
		in.LogPersister.Errorf("Unable to prepare original traffic config to rollback Lambda function %s: invalid traffic config stored", fm.Spec.Name)
		return false
	}

	if err = client.UpdateTrafficConfig(ctx, fm, originalTrafficCfg); err != nil {
		in.LogPersister.Errorf("Failed to rollback original traffic config for Lambda function %s: %v", fm.Spec.Name, err)
		return false
	}

	return true
}
