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

package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/ecs"
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
		e.LogPersister.Errorf("Unsupported stage %s for ECS application", e.Stage.Name)
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

	deployCfg := runningDS.DeploymentConfig.ECSDeploymentSpec
	if deployCfg == nil {
		e.LogPersister.Errorf("Malformed deployment configuration: missing ECSDeploymentSpec")
		return model.StageStatus_STAGE_FAILURE
	}

	cloudProviderName, cloudProviderCfg, found := findCloudProvider(&e.Input)
	if !found {
		return model.StageStatus_STAGE_FAILURE
	}

	taskDefinition, ok := loadTaskDefinition(&e.Input, deployCfg.Input.TaskDefinitionFile, runningDS)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}
	serviceDefinition, ok := loadServiceDefinition(&e.Input, deployCfg.Input.ServiceDefinitionFile, runningDS)
	if !ok {
		return model.StageStatus_STAGE_FAILURE
	}

	if !rollback(ctx, &e.Input, cloudProviderName, cloudProviderCfg, taskDefinition, serviceDefinition) {
		return model.StageStatus_STAGE_FAILURE
	}

	return model.StageStatus_STAGE_SUCCESS
}

func rollback(ctx context.Context, in *executor.Input, cloudProviderName string, cloudProviderCfg *config.CloudProviderECSConfig, taskDefinition types.TaskDefinition, serviceDefinition types.Service) bool {
	in.LogPersister.Infof("Start rollback the ECS service and task definition: %s and %s to original stage", serviceDefinition.ServiceName, *taskDefinition.TaskDefinitionArn)
	client, err := provider.DefaultRegistry().Client(cloudProviderName, cloudProviderCfg, in.Logger)
	if err != nil {
		in.LogPersister.Errorf("Unable to create ECS client for the provider %s: %v", cloudProviderName, err)
		return false
	}

	td, err := client.RegisterTaskDefinition(ctx, taskDefinition)
	if err != nil {
		in.LogPersister.Errorf("Failed to register ECS task definition %s: %v", taskDefinition.TaskDefinitionArn, err)
		return false
	}
	serviceDefinition.TaskDefinition = td.TaskDefinitionArn
	// Rollback ECS service configuration to previous state.
	if _, err := client.UpdateService(ctx, serviceDefinition); err != nil {
		in.LogPersister.Errorf("Unable to rollback ECS service %s configuration to previous stage: %w", serviceDefinition.ServiceName, err)
		return false
	}

	if _, err := client.CreateTaskSet(ctx, serviceDefinition, taskDefinition); err != nil {
		in.LogPersister.Errorf("Failed to create ECS task set %s: %v", serviceDefinition.ServiceName, err)
		return false
	}

	in.LogPersister.Infof("Rolled back the ECS service %s and task definition %s configuration to original stage", serviceDefinition.ServiceName, *taskDefinition.TaskDefinitionArn)
	return true
}
