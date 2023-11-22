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

package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/pipe-cd/pipecd/pkg/app/piped/deploysource"
	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type deployExecutor struct {
	executor.Input

	deploySource         *deploysource.DeploySource
	appCfg               *config.ECSApplicationSpec
	platformProviderName string
	platformProviderCfg  *config.PlatformProviderECSConfig
}

func (e *deployExecutor) Execute(sig executor.StopSignal) model.StageStatus {
	ctx := sig.Context()
	ds, err := e.TargetDSP.GetReadOnly(ctx, e.LogPersister)
	if err != nil {
		e.LogPersister.Errorf("Failed to prepare target deploy source data (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	e.deploySource = ds
	e.appCfg = ds.ApplicationConfig.ECSApplicationSpec
	if e.appCfg == nil {
		e.LogPersister.Errorf("Malformed application configuration: missing ECSApplicationSpec")
		return model.StageStatus_STAGE_FAILURE
	}

	var found bool
	e.platformProviderName, e.platformProviderCfg, found = findPlatformProvider(&e.Input)
	if !found {
		return model.StageStatus_STAGE_FAILURE
	}

	var (
		originalStatus = e.Stage.Status
		status         model.StageStatus
	)

	switch model.Stage(e.Stage.Name) {
	case model.StageECSSync:
		status = e.ensureSync(ctx)
	case model.StageECSCanaryRollout:
		status = e.ensureCanaryRollout(ctx)
	case model.StageECSPrimaryRollout:
		status = e.ensurePrimaryRollout(ctx)
	case model.StageECSCanaryClean:
		status = e.ensureCanaryClean(ctx)
	case model.StageECSTrafficRouting:
		status = e.ensureTrafficRouting(ctx)
	default:
		e.LogPersister.Errorf("Unsupported stage %s for ECS application", e.Stage.Name)
		return model.StageStatus_STAGE_FAILURE
	}

	return executor.DetermineStageStatus(sig.Signal(), originalStatus, status)
}

func (e *deployExecutor) ensureSync(ctx context.Context) model.StageStatus {
	ecsInput := e.appCfg.Input

	taskDefinition, ok := loadTaskDefinition(&e.Input, ecsInput.TaskDefinitionFile, e.deploySource)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	if ecsInput.IsStandaloneTask() {
		if !runStandaloneTask(ctx, &e.Input, e.platformProviderName, e.platformProviderCfg, taskDefinition, &ecsInput) {
			return model.StageStatus_STAGE_FAILURE
		}
		return model.StageStatus_STAGE_SUCCESS
	}

	servicedefinition, ok := loadServiceDefinition(&e.Input, ecsInput.ServiceDefinitionFile, e.deploySource)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	var primary *types.LoadBalancer
	// When the service is accessed via ELB, the target group is not used.
	if ecsInput.IsAccessedViaELB() {
		primary, _, ok = loadTargetGroups(&e.Input, e.appCfg, e.deploySource)
		if !ok {
			return model.StageStatus_STAGE_FAILURE
		}
	}

	recreate := e.appCfg.QuickSync.Recreate
	if !sync(ctx, &e.Input, e.platformProviderName, e.platformProviderCfg, recreate, taskDefinition, servicedefinition, primary) {
		return model.StageStatus_STAGE_FAILURE
	}

	return model.StageStatus_STAGE_SUCCESS
}

func (e *deployExecutor) ensurePrimaryRollout(ctx context.Context) model.StageStatus {
	taskDefinition, ok := loadTaskDefinition(&e.Input, e.appCfg.Input.TaskDefinitionFile, e.deploySource)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}
	servicedefinition, ok := loadServiceDefinition(&e.Input, e.appCfg.Input.ServiceDefinitionFile, e.deploySource)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	primary, _, ok := loadTargetGroups(&e.Input, e.appCfg, e.deploySource)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}
	if primary == nil {
		e.LogPersister.Error("Primary target group is required to enable rolling out PRIMARY variant")
		return model.StageStatus_STAGE_FAILURE
	}

	if !rollout(ctx, &e.Input, e.platformProviderName, e.platformProviderCfg, taskDefinition, servicedefinition, primary) {
		return model.StageStatus_STAGE_FAILURE
	}

	return model.StageStatus_STAGE_SUCCESS
}

func (e *deployExecutor) ensureCanaryRollout(ctx context.Context) model.StageStatus {
	taskDefinition, ok := loadTaskDefinition(&e.Input, e.appCfg.Input.TaskDefinitionFile, e.deploySource)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}
	servicedefinition, ok := loadServiceDefinition(&e.Input, e.appCfg.Input.ServiceDefinitionFile, e.deploySource)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	_, canary, ok := loadTargetGroups(&e.Input, e.appCfg, e.deploySource)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}
	if canary == nil {
		e.LogPersister.Error("Canary target group is required to enable rolling out CANARY variant")
		return model.StageStatus_STAGE_FAILURE
	}

	if !rollout(ctx, &e.Input, e.platformProviderName, e.platformProviderCfg, taskDefinition, servicedefinition, canary) {
		return model.StageStatus_STAGE_FAILURE
	}

	return model.StageStatus_STAGE_SUCCESS
}

func (e *deployExecutor) ensureTrafficRouting(ctx context.Context) model.StageStatus {
	primary, canary, ok := loadTargetGroups(&e.Input, e.appCfg, e.deploySource)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}
	if primary == nil || canary == nil {
		e.LogPersister.Error("Primary/Canary target group are required to enable traffic routing")
		return model.StageStatus_STAGE_FAILURE
	}

	if !routing(ctx, &e.Input, e.platformProviderName, e.platformProviderCfg, *primary, *canary) {
		return model.StageStatus_STAGE_FAILURE
	}
	return model.StageStatus_STAGE_SUCCESS
}

func (e *deployExecutor) ensureCanaryClean(ctx context.Context) model.StageStatus {
	if !clean(ctx, &e.Input, e.platformProviderName, e.platformProviderCfg) {
		return model.StageStatus_STAGE_FAILURE
	}

	return model.StageStatus_STAGE_SUCCESS
}
