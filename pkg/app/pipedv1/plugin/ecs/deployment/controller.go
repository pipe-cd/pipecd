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

// controllerType returns the deployment controller type of the given service definition.
// A service definition without deploymentController uses the ECS controller, which is the AWS default.
func controllerType(service types.Service) types.DeploymentControllerType {
	if service.DeploymentController == nil || service.DeploymentController.Type == "" {
		return types.DeploymentControllerTypeEcs
	}
	return service.DeploymentController.Type
}

// resolveController returns the deployment controller type declared by the service definition.
// It fails when the controller is not supported, or when the service already exists
// with a different controller, because the controller cannot be switched by PipeCD.
func resolveController(ctx context.Context, client provider.Client, serviceDef types.Service) (types.DeploymentControllerType, error) {
	desired := controllerType(serviceDef)
	switch desired {
	case types.DeploymentControllerTypeEcs, types.DeploymentControllerTypeExternal:
	case types.DeploymentControllerTypeCodeDeploy:
		return "", fmt.Errorf("deployment controller %s is not supported", desired)
	default:
		return "", fmt.Errorf("unknown deployment controller %q", desired)
	}

	found, err := client.ServiceExists(ctx, *serviceDef.ClusterArn, *serviceDef.ServiceName)
	if err != nil {
		return "", fmt.Errorf("failed to check service %s existence: %w", *serviceDef.ServiceName, err)
	}
	if !found {
		return desired, nil
	}

	live, err := client.DescribeService(ctx, serviceDef)
	if err != nil {
		return "", fmt.Errorf("failed to describe service %s: %w", *serviceDef.ServiceName, err)
	}
	if running := controllerType(*live); running != desired {
		return "", fmt.Errorf(
			"service %s uses deployment controller %s but the service definition declares %s: switching controller is not supported, recreate the service",
			*serviceDef.ServiceName, running, desired,
		)
	}
	return desired, nil
}

// rejectExternalOnlyStage returns an error when a stage that only works with the EXTERNAL controller
// (primary rollout, canary rollout/clean, traffic routing) is used with a service managed by another controller.
func rejectExternalOnlyStage(stageName string, serviceDef types.Service) error {
	if c := controllerType(serviceDef); c != types.DeploymentControllerTypeExternal {
		return fmt.Errorf(
			"stage %s requires the EXTERNAL deployment controller but the service definition uses %s: configure deploymentConfiguration in the service definition instead",
			stageName, c,
		)
	}
	return nil
}

// checkExternalOnlyStageAllowed loads the service definition of the target deployment source
// and rejects the stage when the service is not managed by the EXTERNAL controller.
func checkExternalOnlyStageAllowed(input *sdk.ExecuteStageInput[ecsconfig.ECSApplicationSpec]) error {
	ds := input.Request.TargetDeploymentSource
	cfg, err := ds.AppConfig()
	if err != nil {
		return fmt.Errorf("failed to load app config: %w", err)
	}
	if cfg.Spec.Input.ServiceDefinitionFile == "" {
		return nil
	}
	serviceDef, err := provider.ParseServiceDefinition(ds.ApplicationDirectory, cfg.Spec.Input.ServiceDefinitionFile)
	if err != nil {
		return fmt.Errorf("failed to load service definition: %w", err)
	}
	return rejectExternalOnlyStage(input.Request.StageName, serviceDef)
}
