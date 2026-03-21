// Copyright 2026 The PipeCD Authors.
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

package deployment

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	ecsconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/provider"
)

func (p *ECSPlugin) executeECSRollbackStage(
	ctx context.Context,
	input *sdk.ExecuteStageInput[ecsconfig.ECSApplicationSpec],
	deployTarget *sdk.DeployTarget[ecsconfig.ECSDeployTargetConfig],
) sdk.StageStatus {
	lp := input.Client.LogPersister()

	runningSource := input.Request.RunningDeploymentSource
	if runningSource.CommitHash == "" {
		lp.Error("Unable to determine the last deployed commit to rollback. It seems this is the first deployment.")
		return sdk.StageStatusFailure
	}

	cfg, err := runningSource.AppConfig()
	if err != nil {
		lp.Errorf("Failed to load app config from running deployment source: %v", err)
		return sdk.StageStatusFailure
	}

	client, err := provider.DefaultRegistry().Client(deployTarget.Name, deployTarget.Config)
	if err != nil {
		lp.Errorf("Failed to get ECS client for deploy target %s: %v", deployTarget.Name, err)
		return sdk.StageStatusFailure
	}

	taskDef, err := provider.LoadTaskDefinition(runningSource.ApplicationDirectory, cfg.Spec.Input.TaskDefinitionFile)
	if err != nil {
		lp.Errorf("Failed to load task definition from running deployment source: %v", err)
		return sdk.StageStatusFailure
	}

	// Standalone task mode does not manage an ECS service so rollback is a no-op
	if cfg.Spec.Input.ServiceDefinitionFile == "" && cfg.Spec.Input.RunStandaloneTask {
		lp.Info("Standalone task mode: no service to rollback")
		return sdk.StageStatusSuccess
	}

	serviceDef, err := provider.LoadServiceDefinition(
		runningSource.ApplicationDirectory,
		cfg.Spec.Input.ServiceDefinitionFile,
		input,
	)
	if err != nil {
		lp.Errorf("Failed to load service definition from running deployment source: %v", err)
		return sdk.StageStatusFailure
	}

	var primary *types.LoadBalancer
	if cfg.Spec.Input.AccessType == "ELB" {
		primary, _, err = provider.LoadTargetGroups(cfg.Spec.Input.TargetGroups)
		if err != nil {
			lp.Errorf("Failed to load target groups from running deployment source: %v", err)
			return sdk.StageStatusFailure
		}
	}

	lp.Infof("Rolling back ECS service %s and task definition family %s", *serviceDef.ServiceName, *taskDef.Family)
	if err := rollback(ctx, lp, client, taskDef, serviceDef, primary); err != nil {
		lp.Errorf("Failed to rollback ECS service: %v", err)
		return sdk.StageStatusFailure
	}

	lp.Successf("Successfully rolled back ECS service %s to commit %s", *serviceDef.ServiceName, runningSource.CommitHash)
	return sdk.StageStatusSuccess
}

// rollback restores the ECS service and task set to the state defined in the running deployment source.
func rollback(
	ctx context.Context,
	lp sdk.StageLogPersister,
	client provider.Client,
	taskDef types.TaskDefinition,
	serviceDef types.Service,
	primary *types.LoadBalancer,
) error {
	lp.Infof("Registering task definition family %s", *taskDef.Family)
	td, err := client.RegisterTaskDefinition(ctx, taskDef)
	if err != nil {
		return fmt.Errorf("failed to register task definition %s: %w", *taskDef.Family, err)
	}

	lp.Infof("Applying service definition for service %s", *serviceDef.ServiceName)
	service, err := applyServiceDefinition(ctx, lp, client, serviceDef)
	if err != nil {
		return fmt.Errorf("failed to apply service definition for service %s: %w", *serviceDef.ServiceName, err)
	}

	// Capture existing task sets before creating the rollback task set
	lp.Infof("Getting current task sets for service %s", *service.ServiceName)
	prevTaskSets, err := client.GetServiceTaskSets(ctx, *service)
	if err != nil {
		return fmt.Errorf("failed to get task sets for service %s: %w", *service.ServiceName, err)
	}

	// Create a new task set at 100% scale to restore the original state
	lp.Infof("Creating rollback task set for service %s", *service.ServiceName)
	taskSet, err := client.CreateTaskSet(ctx, *service, *td, primary, 100)
	if err != nil {
		return fmt.Errorf("failed to create task set for service %s: %w", *service.ServiceName, err)
	}

	// Promote the new task set to PRIMARY
	lp.Infof("Promoting rollback task set to PRIMARY for service %s", *service.ServiceName)
	if _, err = client.UpdateServicePrimaryTaskSet(ctx, *service, *taskSet); err != nil {
		return fmt.Errorf("failed to update primary task set for service %s: %w", *service.ServiceName, err)
	}

	// TODO: Rollback ELB listener weights (100% primary, 0% canary)
	// once the GetListenerArns and ModifyListeners being implemented

	// Delete all previous task sets including any remaining canary tasksets
	lp.Info("Deleting previous task sets")
	for _, ts := range prevTaskSets {
		lp.Infof("Deleting task set %s", *ts.TaskSetArn)
		if err := client.DeleteTaskSet(ctx, ts); err != nil {
			return fmt.Errorf("failed to delete task set %s: %w", *ts.TaskSetArn, err)
		}
	}

	return nil
}
