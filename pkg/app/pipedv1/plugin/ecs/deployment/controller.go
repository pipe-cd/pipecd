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

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/provider"
)

// deploymentController is the strategy interface for deploying ECS services.
//
// Each ECS deployment controller type (EXTERNAL, ECS) has its own implementation.
type deploymentController interface {
	// Sync performs a full sync of the ECS service (used by ECS_SYNC stage).
	Sync(ctx context.Context, lp sdk.StageLogPersister, client provider.Client,
		taskDef types.TaskDefinition, serviceDef types.Service,
		primary *types.LoadBalancer, recreate bool) error

	// PrimaryRollout rolls out the new task definition as the primary (used by ECS_PRIMARY_ROLLOUT stage).
	PrimaryRollout(ctx context.Context, lp sdk.StageLogPersister, client provider.Client,
		taskDef types.TaskDefinition, serviceDef types.Service,
		primary *types.LoadBalancer) error

	// Rollback restores the service to the state of the running deployment source (used by ECS_ROLLBACK stage).
	Rollback(ctx context.Context, lp sdk.StageLogPersister, client provider.Client,
		taskDef types.TaskDefinition, serviceDef types.Service,
		primary *types.LoadBalancer) error
}

// newDeploymentController returns the appropriate deploymentController based on the deployment controller type declared in the service definition.
//
// Defaults to externalController when the type is EXTERNAL or unset.
func newDeploymentController(serviceDef types.Service) deploymentController {
	if serviceDef.DeploymentController != nil &&
		serviceDef.DeploymentController.Type == types.DeploymentControllerTypeEcs {
		return &ecsController{}
	}
	return &externalController{}
}

// isECSControllerType returns true when the service definition uses the native ECS deployment controller.
func isECSControllerType(serviceDef types.Service) bool {
	return serviceDef.DeploymentController != nil &&
		serviceDef.DeploymentController.Type == types.DeploymentControllerTypeEcs
}

// externalController implements deploymentController for EXTERNAL deployment controller type.
type externalController struct{}

func (e *externalController) Sync(
	ctx context.Context,
	lp sdk.StageLogPersister,
	client provider.Client,
	taskDef types.TaskDefinition,
	serviceDef types.Service,
	primary *types.LoadBalancer,
	recreate bool,
) error {
	lp.Info("Start applying the ECS task definition")
	td, err := applyTaskDefinition(ctx, client, taskDef)
	if err != nil {
		lp.Errorf("Failed to apply task definition: %v", err)
		return fmt.Errorf("failed to apply task definition: %w", err)
	}

	lp.Info("Start applying the ECS service definition")
	service, _, err := applyServiceDefinition(ctx, lp, client, serviceDef)
	if err != nil {
		lp.Errorf("Failed to apply service definition: %v", err)
		return fmt.Errorf("failed to apply service definition: %w", err)
	}

	if recreate {
		cnt := service.DesiredCount
		lp.Info("Recreate option is enabled, stop all running tasks before creating new task set")
		if err := client.PruneServiceTasks(ctx, *service); err != nil {
			lp.Errorf("Failed to prune service tasks: %v", err)
			return fmt.Errorf("failed to prune service tasks: %w", err)
		}

		lp.Info("Start rolling out ECS TaskSet for the new task definition")
		if err = createPrimaryTaskSet(ctx, lp, client, *service, *td, primary); err != nil {
			lp.Errorf("Failed to rollout ECS TaskSet for service %s: %v", *service.ServiceName, err)
			return fmt.Errorf("failed to create primary task set: %w", err)
		}

		lp.Info("Deleting old ECS TaskSets")
		if err = deleteOldTaskSets(ctx, client, *service); err != nil {
			lp.Errorf("Failed to delete old Tasksets of service %s: %v", *service.ServiceName, err)
			return fmt.Errorf("failed to delete old tasksets: %w", err)
		}

		lp.Infof("Scale up ECS desired tasks count back to %d", cnt)
		service.DesiredCount = cnt
		if _, err = client.UpdateService(ctx, *service); err != nil {
			lp.Errorf("Failed to revive service tasks: %v", err)
			return fmt.Errorf("failed to revive service tasks: %w", err)
		}
	} else {
		lp.Info("Start rolling out ECS TaskSet for the new task definition")
		if err = createPrimaryTaskSet(ctx, lp, client, *service, *td, primary); err != nil {
			lp.Errorf("Failed to rollout ECS TaskSet for service %s: %v", *service.ServiceName, err)
			return fmt.Errorf("failed to create primary task set: %w", err)
		}

		lp.Info("Deleting old ECS TaskSets")
		if err = deleteOldTaskSets(ctx, client, *service); err != nil {
			lp.Errorf("Failed to delete old Tasksets of service %s: %v", *service.ServiceName, err)
			return fmt.Errorf("failed to delete old tasksets: %w", err)
		}
	}

	lp.Infof("Wait service %s to reach stable state", *service.ServiceName)
	if err := client.WaitServiceStable(ctx, *service.ClusterArn, *service.ServiceName); err != nil {
		lp.Errorf("Failed to wait for service to be stable: %v", err)
		return err
	}

	return nil
}

func (e *externalController) PrimaryRollout(
	ctx context.Context,
	lp sdk.StageLogPersister,
	client provider.Client,
	taskDef types.TaskDefinition,
	serviceDef types.Service,
	primary *types.LoadBalancer,
) error {
	lp.Info("Start applying the ECS task definition")
	td, err := applyTaskDefinition(ctx, client, taskDef)
	if err != nil {
		return fmt.Errorf("failed to apply task definition: %w", err)
	}

	lp.Info("Start applying the ECS service definition")
	service, _, err := applyServiceDefinition(ctx, lp, client, serviceDef)
	if err != nil {
		return fmt.Errorf("failed to apply service definition: %w", err)
	}

	lp.Infof("Get current PRIMARY taskset")
	currPrimaryTs, err := client.GetPrimaryTaskSet(ctx, *service)
	if err != nil {
		return fmt.Errorf("failed to get current primary taskset: %w", err)
	}

	lp.Infof("Rolling out new PRIMARY taskset for service %s", *service.ServiceName)
	if err = createPrimaryTaskSet(ctx, lp, client, *service, *td, primary); err != nil {
		return fmt.Errorf("failed to create primary taskset for service %s: %w", *service.ServiceName, err)
	}

	lp.Infof("Deleting old PRIMARY taskset")
	if currPrimaryTs != nil {
		if err = client.DeleteTaskSet(ctx, *currPrimaryTs); err != nil {
			return fmt.Errorf("failed to delete old primary taskset: %w", err)
		}
	}

	lp.Infof("Waiting for service %s to reach stable state", *service.ServiceName)
	if err := client.WaitServiceStable(ctx, *service.ClusterArn, *service.ServiceName); err != nil {
		return fmt.Errorf("service %s did not reach stable state: %w", *service.ServiceName, err)
	}

	lp.Successf("Successfully rolled out PRIMARY task set for service %s", *service.ServiceName)
	return nil
}

func (e *externalController) Rollback(
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
	service, _, err := applyServiceDefinition(ctx, lp, client, serviceDef)
	if err != nil {
		return fmt.Errorf("failed to apply service definition for service %s: %w", *serviceDef.ServiceName, err)
	}

	lp.Infof("Getting current task sets for service %s", *service.ServiceName)
	prevTaskSets, err := client.GetServiceTaskSets(ctx, *service)
	if err != nil {
		return fmt.Errorf("failed to get task sets for service %s: %w", *service.ServiceName, err)
	}

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

// ecsController implements deploymentController for the native ECS deployment controller type.
//
// Deployments are triggered by calling UpdateService with a new task definition and ForceNewDeployment=true
type ecsController struct{}

func (e *ecsController) Sync(
	ctx context.Context,
	lp sdk.StageLogPersister,
	client provider.Client,
	taskDef types.TaskDefinition,
	serviceDef types.Service,
	_ *types.LoadBalancer,
	_ bool,
) error {
	return e.deploy(ctx, lp, client, taskDef, serviceDef)
}

func (e *ecsController) PrimaryRollout(
	ctx context.Context,
	lp sdk.StageLogPersister,
	client provider.Client,
	taskDef types.TaskDefinition,
	serviceDef types.Service,
	_ *types.LoadBalancer,
) error {
	return e.deploy(ctx, lp, client, taskDef, serviceDef)
}

func (e *ecsController) Rollback(
	ctx context.Context,
	lp sdk.StageLogPersister,
	client provider.Client,
	taskDef types.TaskDefinition,
	serviceDef types.Service,
	_ *types.LoadBalancer,
) error {
	return e.deploy(ctx, lp, client, taskDef, serviceDef)
}

// deploy is the shared deployment flow for all ECS controller stages:
//
// register task definition -> apply service -> force new deployment -> wait stable.
func (e *ecsController) deploy(
	ctx context.Context,
	lp sdk.StageLogPersister,
	client provider.Client,
	taskDef types.TaskDefinition,
	serviceDef types.Service,
) error {
	lp.Info("Start applying the ECS task definition")
	td, err := applyTaskDefinition(ctx, client, taskDef)
	if err != nil {
		return fmt.Errorf("failed to apply task definition: %w", err)
	}

	// Inject the registered task definition ARN so CreateService can include it.
	// ECS deployment controller requires a task definition at service creation time,
	// unlike EXTERNAL controller which sets it per-task-set via CreateTaskSet.
	serviceDef.TaskDefinition = td.TaskDefinitionArn

	lp.Info("Start applying the ECS service definition")
	service, newlyCreated, err := applyServiceDefinition(ctx, lp, client, serviceDef)
	if err != nil {
		return fmt.Errorf("failed to apply service definition: %w", err)
	}

	if !newlyCreated {
		// For existing services, trigger a new deployment with the updated task definition.
		// When the service was just created, CreateService already starts the first deployment automatically
		// (calling ForceNewDeployment would trigger a second redundant deployment).
		lp.Infof("Triggering new deployment for service %s with task definition %s", *service.ServiceName, *td.TaskDefinitionArn)
		if _, err := client.ForceNewDeployment(ctx, *service, *td); err != nil {
			return fmt.Errorf("failed to force new deployment for service %s: %w", *service.ServiceName, err)
		}
	}

	lp.Infof("Waiting for service %s to reach stable state", *service.ServiceName)
	if err := client.WaitServiceStable(ctx, *service.ClusterArn, *service.ServiceName); err != nil {
		return fmt.Errorf("service %s did not reach stable state: %w", *service.ServiceName, err)
	}

	return nil
}
