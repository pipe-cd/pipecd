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

package lambda

import (
	"context"
	"fmt"

	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/lambda"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
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
		e.LogPersister.Errorf("Unsupported stage %s for lambda application", e.Stage.Name)
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

	appCfg := runningDS.ApplicationConfig.LambdaApplicationSpec
	if appCfg == nil {
		e.LogPersister.Errorf("Malformed application configuration: missing LambdaApplicationSpec")
		return model.StageStatus_STAGE_FAILURE
	}

	platformProviderName, platformProviderCfg, found := findPlatformProvider(&e.Input)
	if !found {
		return model.StageStatus_STAGE_FAILURE
	}

	fm, ok := loadFunctionManifest(&e.Input, appCfg.Input.FunctionManifestFile, runningDS)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	if !rollback(ctx, &e.Input, platformProviderName, platformProviderCfg, fm) {
		return model.StageStatus_STAGE_FAILURE
	}

	return model.StageStatus_STAGE_SUCCESS
}

func rollback(ctx context.Context, in *executor.Input, platformProviderName string, platformProviderCfg *config.PlatformProviderLambdaConfig, fm provider.FunctionManifest) bool {
	in.LogPersister.Infof("Start rollback the lambda function: %s to original stage", fm.Spec.Name)
	client, err := provider.DefaultRegistry().Client(platformProviderName, platformProviderCfg, in.Logger)
	if err != nil {
		in.LogPersister.Errorf("Unable to create Lambda client for the provider %s: %v", platformProviderName, err)
		return false
	}

	// Rollback Lambda application configuration to previous state.
	if err := client.UpdateFunction(ctx, fm); err != nil {
		in.LogPersister.Errorf("Unable to rollback Lambda function %s configuration to previous stage: %v", fm.Spec.Name, err)
		return false
	}
	in.LogPersister.Infof("Rolled back the lambda function %s configuration to original stage", fm.Spec.Name)

	// Rollback traffic routing to previous state.
	// Restore original traffic config from metadata store.
	originalTrafficKeyName := fmt.Sprintf("original-traffic-%s", in.Deployment.RunningCommitHash)
	originalTrafficCfgData, ok := in.MetadataStore.Shared().Get(originalTrafficKeyName)
	if !ok {
		in.LogPersister.Errorf("Unable to prepare original traffic config to rollback Lambda function %s. No traffic changes have been committed yet.", fm.Spec.Name)
		return false
	}

	originalTrafficCfg := provider.RoutingTrafficConfig{}
	if err := originalTrafficCfg.Decode([]byte(originalTrafficCfgData)); err != nil {
		in.LogPersister.Errorf("Unable to prepare original traffic config to rollback Lambda function %s: %v", fm.Spec.Name, err)
		return false
	}

	// Restore promoted traffic config from metadata store.
	promotedTrafficKeyName := fmt.Sprintf("latest-promote-traffic-%s", in.Deployment.RunningCommitHash)
	promotedTrafficCfgData, ok := in.MetadataStore.Shared().Get(promotedTrafficKeyName)
	// If there is no previous promoted traffic config, which mean no promote run previously so no need to do anything to rollback.
	if !ok {
		in.LogPersister.Info("It seems the traffic has not been changed during the deployment process. No need to rollback the traffic config.")
		return true
	}

	promotedTrafficCfg := provider.RoutingTrafficConfig{}
	if err := promotedTrafficCfg.Decode([]byte(promotedTrafficCfgData)); err != nil {
		in.LogPersister.Errorf("Unable to prepare promoted traffic config to rollback Lambda function %s: %v", fm.Spec.Name, err)
		return false
	}

	switch len(originalTrafficCfg) {
	// Original traffic config has both PRIMARY and SECONDARY version config.
	case 2:
		if err = client.UpdateTrafficConfig(ctx, fm, originalTrafficCfg); err != nil {
			in.LogPersister.Errorf("Failed to rollback original traffic config for Lambda function %s: %v", fm.Spec.Name, err)
			return false
		}
		return true
	// Original traffic config is PRIMARY ONLY config,
	// we need to reset any others SECONDARY created by previous (until failed) PROMOTE stages.
	case 1:
		// Validate stored original traffic config, since it PRIMARY ONLY, the percent must be float64(100)
		primary, ok := originalTrafficCfg[provider.TrafficPrimaryVersionKeyName]
		if !ok || primary.Percent != float64(100) {
			in.LogPersister.Errorf("Unable to prepare original traffic config: invalid original traffic config stored")
			return false
		}

		// Update promoted traffic config by add 0% SECONDARY for reset remote promoted version config.
		if !configureTrafficRouting(promotedTrafficCfg, primary.Version, 100) {
			in.LogPersister.Errorf("Unable to prepare traffic config to rollback Lambda function %s: can not reset promoted version", fm.Spec.Name)
			return false
		}

		if err = client.UpdateTrafficConfig(ctx, fm, promotedTrafficCfg); err != nil {
			in.LogPersister.Errorf("Failed to rollback original traffic config for Lambda function %s: %v", fm.Spec.Name, err)
			return false
		}
		return true
	default:
		in.LogPersister.Errorf("Unable to prepare original traffic config: invalid original traffic config stored")
		return false
	}
}
