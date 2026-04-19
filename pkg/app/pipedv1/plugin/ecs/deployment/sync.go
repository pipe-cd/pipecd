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

func (p *ECSPlugin) executeECSSyncStage(
	ctx context.Context,
	input *sdk.ExecuteStageInput[ecsconfig.ECSApplicationSpec],
	deployTarget *sdk.DeployTarget[ecsconfig.ECSDeployTargetConfig],
) sdk.StageStatus {
	lp := input.Client.LogPersister()
	lp.Info("Starting ECS sync stage")

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

	// If there is no service definition file and the standalone task flag is set,
	// run the task as a standalone task without creating a service
	if cfg.Spec.Input.ServiceDefinitionFile == "" && cfg.Spec.Input.RunStandaloneTask {
		lp.Info("Standalone task detected, no service definition file found")
		if err := runStandaloneTask(ctx, client, taskDef, input); err != nil {
			lp.Errorf("Failed to run standalone task: %v", err)
			return sdk.StageStatusFailure
		}
		lp.Success("Successfully run standalone task")
		return sdk.StageStatusSuccess
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
	if err := ctrl.Sync(ctx, lp, client, taskDef, serviceDef, primary, cfg.Spec.QuickSyncOptions.Recreate); err != nil {
		lp.Errorf("Failed to sync ECS service: %v", err)
		return sdk.StageStatusFailure
	}

	return sdk.StageStatusSuccess
}

func runStandaloneTask(
	ctx context.Context,
	client provider.Client,
	taskDef types.TaskDefinition,
	input *sdk.ExecuteStageInput[ecsconfig.ECSApplicationSpec],
) error {
	lp := input.Client.LogPersister()
	lp.Info("Running standalone task")

	lp.Info("Start applying the ECS task definition")
	td, err := applyTaskDefinition(ctx, client, taskDef)
	if err != nil {
		return fmt.Errorf("failed to apply task definition: %w", err)
	}

	deployInput := input.Request.TargetDeploymentSource.ApplicationConfig.Spec.Input
	tags := provider.MakeTags(map[string]string{
		provider.LabelManagedBy:   provider.ManagedByECSPlugin,
		provider.LabelPiped:       input.Request.Deployment.PipedID,
		provider.LabelApplication: input.Request.Deployment.ApplicationID,
		provider.LabelCommitHash:  input.Request.TargetDeploymentSource.CommitHash,
	})
	err = client.RunTask(ctx, *td, deployInput.ClusterARN, deployInput.LaunchType, &deployInput.AwsVpcConfiguration, tags)
	if err != nil {
		return fmt.Errorf("failed to run task: %w", err)
	}

	return nil
}

func applyTaskDefinition(
	ctx context.Context,
	client provider.Client,
	taskDef types.TaskDefinition,
) (*types.TaskDefinition, error) {
	td, err := client.RegisterTaskDefinition(ctx, taskDef)
	if err != nil {
		return nil, fmt.Errorf("failed to register task definition: %w", err)
	}
	return td, nil
}

// applyServiceDefinition creates or updates the ECS service based on its existence.
//
// newlyCreated is true when the service did not exist and was newly created.
func applyServiceDefinition(
	ctx context.Context,
	lp sdk.StageLogPersister,
	client provider.Client,
	serviceDef types.Service,
) (service *types.Service, newlyCreated bool, err error) {
	// Check whether the service already exists or not.
	// If it exists, update the service, otherwise create a new one.
	found, err := client.ServiceExists(ctx, *serviceDef.ClusterArn, *serviceDef.ServiceName)
	if err != nil {
		return nil, false, fmt.Errorf("failed to check service %s existence: %w", *serviceDef.ServiceName, err)
	}

	if found {
		svcStatus, err := client.GetServiceStatus(ctx, *serviceDef.ClusterArn, *serviceDef.ServiceName)
		if err != nil {
			return nil, false, fmt.Errorf("failed to get service %s status: %w", *serviceDef.ServiceName, err)
		}
		lp.Infof("Service %s already exists with status %s", *serviceDef.ServiceName, svcStatus)

		// Only update the service when it is in ACTIVE status
		// Nothing can be performed if the service is in DRAINING or INACTIVE status
		if svcStatus != "ACTIVE" {
			return nil, false, fmt.Errorf("service %s is in %s status, cannot be updated", *serviceDef.ServiceName, svcStatus)
		}

		lp.Infof("Updating service %s", *serviceDef.ServiceName)
		service, err = client.UpdateService(ctx, serviceDef)
		if err != nil {
			return nil, false, fmt.Errorf("failed to update service %s: %w", *serviceDef.ServiceName, err)
		}

		currentTags, err := client.ListTags(ctx, *service.ServiceArn)
		if err != nil {
			return nil, false, fmt.Errorf("failed to list tags for ECS service %s: %w", *serviceDef.ServiceName, err)
		}

		tagsToRemove := findTagsToRemove(currentTags, serviceDef.Tags)
		if len(tagsToRemove) > 0 {
			lp.Infof("Found tags to remove from service %s: %v", *serviceDef.ServiceName, tagsToRemove)
			if err := client.UntagResource(ctx, *service.ServiceArn, tagsToRemove); err != nil {
				return nil, false, fmt.Errorf("failed to remove tags from ECS service %s: %w", *serviceDef.ServiceName, err)
			}
		}
		if err := client.TagResource(ctx, *service.ServiceArn, serviceDef.Tags); err != nil {
			return nil, false, fmt.Errorf("failed to update tags of ECS service %s: %w", *serviceDef.ServiceName, err)
		}
		// Re-assign tags to service object because UpdateService API doesn't return tags.
		service.Tags = serviceDef.Tags
		return service, false, nil
	}

	lp.Infof("Service %s does not exist, creating a new service", *serviceDef.ServiceName)
	service, err = client.CreateService(ctx, serviceDef)
	if err != nil {
		return nil, false, fmt.Errorf("failed to create service %s: %w", *serviceDef.ServiceName, err)
	}
	return service, true, nil
}

func findTagsToRemove(currentTags, desiredTags []types.Tag) []string {
	var tagsToRemove []string

	// Mark all desired tags in a map for easier lookup
	desired := make(map[string]struct{})
	for _, t := range desiredTags {
		desired[*t.Key] = struct{}{}
	}

	for _, t := range currentTags {
		if _, exists := desired[*t.Key]; !exists {
			tagsToRemove = append(tagsToRemove, *t.Key)
		}
	}

	return tagsToRemove
}

func createPrimaryTaskSet(
	ctx context.Context,
	lp sdk.StageLogPersister,
	client provider.Client,
	service types.Service,
	taskDef types.TaskDefinition,
	primary *types.LoadBalancer,
) error {
	// Create a task set in the specified cluster and service.
	// In case of creating Primary taskset, the number of desired tasks scale is always set to 100
	// which means we create as many tasks as the current primary taskset has.
	lp.Infof("Creating primary task set for service %s", *service.ServiceName)
	taskSet, err := client.CreateTaskSet(ctx, service, taskDef, primary, 100)
	if err != nil {
		return fmt.Errorf("failed to create primary task set: %w", err)
	}

	// Mark the new task set as PRIMARY
	lp.Infof("Updating primary task set for service %s", *service.ServiceName)
	if _, err = client.UpdateServicePrimaryTaskSet(ctx, service, *taskSet); err != nil {
		return err
	}

	return nil
}

func deleteOldTaskSets(
	ctx context.Context,
	client provider.Client,
	service types.Service,
) error {
	// Get all TaskSets (with status PRIMARY, ACTIVE)
	taskSets, err := client.GetServiceTaskSets(ctx, service)
	if err != nil {
		return fmt.Errorf("failed to get task sets: %w", err)
	}

	// Delete old TaskSets (tasksets with status ACTIVE)
	for _, ts := range taskSets {
		if ts.Status != nil && *ts.Status != "PRIMARY" {
			if err = client.DeleteTaskSet(ctx, ts); err != nil {
				return fmt.Errorf("failed to delete old task set %s: %w", *ts.TaskSetArn, err)
			}
		}
	}

	return nil
}
