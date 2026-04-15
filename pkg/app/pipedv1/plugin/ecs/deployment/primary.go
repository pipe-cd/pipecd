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

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	ecsconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/provider"
)

func (p *ECSPlugin) executeECSPrimaryRolloutStage(
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

	var primary *types.LoadBalancer
	// When the services is not accessed via ELB, the target group is not used
	if cfg.Spec.Input.AccessType == "ELB" {
		primary, _, err = provider.LoadTargetGroups(cfg.Spec.Input.TargetGroups)
		if err != nil {
			lp.Errorf("Failed to load target groups: %v", err)
			return sdk.StageStatusFailure
		}
	}

	ctrl := newDeploymentController(serviceDef)
	if err := ctrl.PrimaryRollout(ctx, lp, client, taskDef, serviceDef, primary); err != nil {
		lp.Errorf("Failed to roll out ECS primary task set: %v", err)
		return sdk.StageStatusFailure
	}

	return sdk.StageStatusSuccess
}
