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
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	ecsconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/provider"
)

const canaryTaskSetMetadataKey = "canary-task-set"

func (p *ECSPlugin) executeECSCanaryRolloutStage(
	ctx context.Context,
	input *sdk.ExecuteStageInput[ecsconfig.ECSApplicationSpec],
	deployTarget *sdk.DeployTarget[ecsconfig.ECSDeployTargetConfig],
) sdk.StageStatus {
	lp := input.Client.LogPersister()

	cfg, err := input.Request.TargetDeploymentSource.AppConfig()
	if err != nil {
		lp.Errorf("Failed to load app config: %v", err)
		return sdk.StageStatusFailure
	}

	var options ecsconfig.ECSCanaryRolloutStageOptions
	if len(input.Request.StageConfig) > 0 {
		if err := json.Unmarshal(input.Request.StageConfig, &options); err != nil {
			lp.Errorf("Failed to parse canary rollout stage options: %v", err)
			return sdk.StageStatusFailure
		}
	}
	client, err := provider.DefaultRegistry().Client(deployTarget.Name, deployTarget.Config)
	if err != nil {
		lp.Errorf("Failed to get ECS client for deploy target %s: %v", deployTarget.Name, err)
		return sdk.StageStatusFailure
	}

	taskDef, err := provider.LoadTaskDefinition(input.Request.TargetDeploymentSource.ApplicationDirectory, cfg.Spec.Input.TaskDefinitionFile)
	if err != nil {
		lp.Errorf("Failed to load task definition: %v", err)
		return sdk.StageStatusFailure
	}

	serviceDef, err := provider.LoadServiceDefinition(
		input.Request.TargetDeploymentSource.ApplicationDirectory,
		cfg.Spec.Input.ServiceDefinitionFile,
		input,
	)
	if err != nil {
		lp.Errorf("Failed to load service definition: %v", err)
		return sdk.StageStatusFailure
	}

	var canary *types.LoadBalancer
	if cfg.Spec.Input.AccessType == "ELB" {
		_, canary, err = provider.LoadTargetGroups(cfg.Spec.Input.TargetGroups)
		if err != nil {
			lp.Errorf("Failed to load target groups: %v", err)
			return sdk.StageStatusFailure
		}
		if canary == nil {
			lp.Error("Canary target group is required for ELB access type in ECS_CANARY_ROLLOUT stage")
			return sdk.StageStatusFailure
		}
	}

	taskSet, err := canaryRollout(ctx, lp, client, taskDef, serviceDef, canary, options.Scale)
	if err != nil {
		lp.Errorf("Failed to roll out ECS canary task set: %v", err)
		return sdk.StageStatusFailure
	}

	// Persist the canary task set so ECS_CANARY_CLEAN can delete it later
	taskSetData, err := json.Marshal(taskSet)
	if err != nil {
		lp.Errorf("Failed to marshal canary task set for metadata store: %v", err)
		return sdk.StageStatusFailure
	}
	if err := input.Client.PutDeploymentPluginMetadata(ctx, canaryTaskSetMetadataKey, string(taskSetData)); err != nil {
		lp.Errorf("Failed to store canary task set to metadata store: %v", err)
		return sdk.StageStatusFailure
	}

	return sdk.StageStatusSuccess
}

// canaryRollout performs the canary rollout workflow:
//
// 1. Registers the task definition
//
// 2. Applies the service definition (creates or updates the service)
//
// 3. Creates a new CANARY task set at specified scale
//
// 4. Waits for the service to reach stable state
func canaryRollout(
	ctx context.Context,
	lp sdk.StageLogPersister,
	client provider.Client,
	taskDef types.TaskDefinition,
	serviceDef types.Service,
	canary *types.LoadBalancer,
	scale float64,
) (*types.TaskSet, error) {
	lp.Info("Start applying the ECS task definition")
	td, err := applyTaskDefinition(ctx, client, taskDef)
	if err != nil {
		return nil, fmt.Errorf("failed to apply task definition: %w", err)
	}

	lp.Info("Start applying the ECS service definition")
	service, err := applyServiceDefinition(ctx, lp, client, serviceDef)
	if err != nil {
		return nil, fmt.Errorf("failed to apply service definition: %w", err)
	}

	lp.Infof("Creating CANARY task set for service %s at scale %.0f%%", *service.ServiceName, scale)
	taskSet, err := client.CreateTaskSet(ctx, *service, *td, canary, scale)
	if err != nil {
		return nil, fmt.Errorf("failed to create canary task set for service %s: %w", *service.ServiceName, err)
	}

	lp.Infof("Waiting for service %s to reach stable state", *service.ServiceName)
	if err := client.WaitServiceStable(ctx, *service.ClusterArn, *service.ServiceName); err != nil {
		return nil, fmt.Errorf("service %s did not reach stable state: %w", *service.ServiceName, err)
	}

	lp.Successf("Successfully rolled out CANARY task set %s for service %s", *taskSet.TaskSetArn, *service.ServiceName)
	return taskSet, nil
}
