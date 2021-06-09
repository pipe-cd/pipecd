// Copyright 2020 The PipeCD Authors.
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
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"

	"github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider"
	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/ecs"
	"github.com/pipe-cd/pipe/pkg/app/piped/deploysource"
	"github.com/pipe-cd/pipe/pkg/app/piped/executor"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

const canaryTaskSetARNKeyName = "canary-taskset-arn"

type registerer interface {
	Register(stage model.Stage, f executor.Factory) error
	RegisterRollback(kind model.ApplicationKind, f executor.Factory) error
}

func Register(r registerer) {
	f := func(in executor.Input) executor.Executor {
		return &deployExecutor{
			Input: in,
		}
	}
	r.Register(model.StageECSSync, f)
	r.Register(model.StageECSCanaryRollout, f)
	r.Register(model.StageECSPrimaryRollout, f)
	r.Register(model.StageECSCanaryClean, f)

	r.RegisterRollback(model.ApplicationKind_ECS, func(in executor.Input) executor.Executor {
		return &rollbackExecutor{
			Input: in,
		}
	})
}

func findCloudProvider(in *executor.Input) (name string, cfg *config.CloudProviderECSConfig, found bool) {
	name = in.Application.CloudProvider
	if name == "" {
		in.LogPersister.Errorf("Missing the CloudProvider name in the application configuration")
		return
	}

	cp, ok := in.PipedConfig.FindCloudProvider(name, model.CloudProviderECS)
	if !ok {
		in.LogPersister.Errorf("The specified cloud provider %q was not found in piped configuration", name)
		return
	}

	cfg = cp.ECSConfig
	found = true
	return
}

func loadServiceDefinition(in *executor.Input, serviceDefinitionFile string, ds *deploysource.DeploySource) (types.Service, bool) {
	in.LogPersister.Infof("Loading service manifest at the %s commit (%s)", ds.RevisionName, ds.RevisionName)

	serviceDefinition, err := provider.LoadServiceDefinition(ds.AppDir, serviceDefinitionFile)
	if err != nil {
		in.LogPersister.Errorf("Failed to load ECS service definition (%v)", err)
		return types.Service{}, false
	}

	in.LogPersister.Infof("Successfully loaded the ECS service definition at the %s commit", ds.RevisionName)
	return serviceDefinition, true
}

func loadTaskDefinition(in *executor.Input, taskDefinitionFile string, ds *deploysource.DeploySource) (types.TaskDefinition, bool) {
	in.LogPersister.Infof("Loading task definition manifest at the %s commit (%s)", ds.RevisionName, ds.RevisionName)

	taskDefinition, err := provider.LoadTaskDefinition(ds.AppDir, taskDefinitionFile)
	if err != nil {
		in.LogPersister.Errorf("Failed to load ECS task definition (%v)", err)
		return types.TaskDefinition{}, false
	}

	in.LogPersister.Infof("Successfully loaded the ECS task definition at the %s commit", ds.RevisionName)
	return taskDefinition, true
}

func loadTargetGroups(in *executor.Input, deployCfg *config.ECSDeploymentSpec, ds *deploysource.DeploySource) (*types.LoadBalancer, *types.LoadBalancer, bool) {
	in.LogPersister.Infof("Loading target groups config at the %s commit (%s)", ds.RevisionName, ds.RevisionName)

	primary, canary, err := provider.LoadTargetGroups(deployCfg.Input.TargetGroups)
	if err != nil {
		in.LogPersister.Errorf("Failed to load TargetGroups (%v)", err)
		return nil, nil, false
	}

	in.LogPersister.Infof("Successfully loaded the ECS target groups at the %s commit", ds.RevisionName)
	return primary, canary, true
}

func applyTaskDefinition(ctx context.Context, cli provider.Client, taskDefinition types.TaskDefinition) (*types.TaskDefinition, error) {
	td, err := cli.RegisterTaskDefinition(ctx, taskDefinition)
	if err != nil {
		return nil, err
	}
	return td, nil
}

func applyServiceDefinition(ctx context.Context, cli provider.Client, serviceDefinition types.Service) (*types.Service, error) {
	found, err := cli.ServiceExists(ctx, *serviceDefinition.ClusterArn, *serviceDefinition.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("unable to validate service name %s: %v", *serviceDefinition.ServiceName, err)
	}

	var service *types.Service
	if found {
		service, err = cli.UpdateService(ctx, serviceDefinition)
		if err != nil {
			return nil, fmt.Errorf("failed to update ECS service %s: %v", *serviceDefinition.ServiceName, err)
		}
	} else {
		service, err = cli.CreateService(ctx, serviceDefinition)
		if err != nil {
			return nil, fmt.Errorf("failed to create ECS service %s: %v", *serviceDefinition.ServiceName, err)
		}
	}

	return service, nil
}

func createPrimaryTaskSet(ctx context.Context, client provider.Client, service types.Service, taskDef types.TaskDefinition, targetGroup types.LoadBalancer) error {
	// Get current PRIMARY task set.
	prevPrimaryTaskSet, err := client.GetPrimaryTaskSet(ctx, service)
	// Ignore error in case it's not found error, the prevPrimaryTaskSet doesn't exist for newly created Service.
	if err != nil && !errors.Is(err, cloudprovider.ErrNotFound) {
		return err
	}

	// Create a task set in the specified cluster and service.
	taskSet, err := client.CreateTaskSet(ctx, service, taskDef, targetGroup)
	if err != nil {
		return err
	}

	// Make new taskSet as PRIMARY task set, so that it will handle production service.
	if _, err = client.UpdateServicePrimaryTaskSet(ctx, service, *taskSet); err != nil {
		return err
	}

	// Remove old taskSet if existed.
	if prevPrimaryTaskSet != nil {
		if err = client.DeleteTaskSet(ctx, service, *prevPrimaryTaskSet.TaskSetArn); err != nil {
			return err
		}
	}

	return nil
}

func sync(ctx context.Context, in *executor.Input, cloudProviderName string, cloudProviderCfg *config.CloudProviderECSConfig, taskDefinition types.TaskDefinition, serviceDefinition types.Service, targetGroup types.LoadBalancer) bool {
	client, err := provider.DefaultRegistry().Client(cloudProviderName, cloudProviderCfg, in.Logger)
	if err != nil {
		in.LogPersister.Errorf("Unable to create ECS client for the provider %s: %v", cloudProviderName, err)
		return false
	}

	in.LogPersister.Infof("Start applying the ECS task definition")
	td, err := applyTaskDefinition(ctx, client, taskDefinition)
	if err != nil {
		in.LogPersister.Errorf("Failed to register ECS task definition of family %s: %v", *taskDefinition.Family, err)
		return false
	}

	in.LogPersister.Infof("Start applying the ECS service definition")
	service, err := applyServiceDefinition(ctx, client, serviceDefinition)
	if err != nil {
		in.LogPersister.Errorf("Failed to apply service %s: %v", *serviceDefinition.ServiceName, err)
		return false
	}

	in.LogPersister.Infof("Start rolling out ECS task set")
	if err := createPrimaryTaskSet(ctx, client, *service, *td, targetGroup); err != nil {
		in.LogPersister.Errorf("Failed to rolling out ECS task set %s: %v", *serviceDefinition.ServiceName, err)
		return false
	}

	in.LogPersister.Infof("Successfully applied the service definition and the task definition for ECS service %s and task definition of family %s", *serviceDefinition.ServiceName, *taskDefinition.Family)
	return true
}

func rollout(ctx context.Context, in *executor.Input, cloudProviderName string, cloudProviderCfg *config.CloudProviderECSConfig, taskDefinition types.TaskDefinition, serviceDefinition types.Service, targetGroup types.LoadBalancer) bool {
	client, err := provider.DefaultRegistry().Client(cloudProviderName, cloudProviderCfg, in.Logger)
	if err != nil {
		in.LogPersister.Errorf("Unable to create ECS client for the provider %s: %v", cloudProviderName, err)
		return false
	}

	in.LogPersister.Infof("Start applying the ECS task definition")
	td, err := applyTaskDefinition(ctx, client, taskDefinition)
	if err != nil {
		in.LogPersister.Errorf("Failed to register ECS task definition of family %s: %v", *taskDefinition.Family, err)
		return false
	}

	in.LogPersister.Infof("Start applying the ECS service definition")
	service, err := applyServiceDefinition(ctx, client, serviceDefinition)
	if err != nil {
		in.LogPersister.Errorf("Failed to apply service %s: %v", *serviceDefinition.ServiceName, err)
		return false
	}

	// Create a task set in the specified cluster and service.
	in.LogPersister.Infof("Start rolling out ECS task set")
	if in.StageConfig.Name == model.StageECSPrimaryRollout {
		// Create PRIMARY task set in case of Primary rollout.
		if err := createPrimaryTaskSet(ctx, client, *service, *td, targetGroup); err != nil {
			in.LogPersister.Errorf("Failed to rolling out ECS task set %s: %v", *serviceDefinition.ServiceName, err)
			return false
		}
	} else {
		// Create ACTIVE task set in case of Canary rollout.
		taskSet, err := client.CreateTaskSet(ctx, *service, *td, targetGroup)
		if err != nil {
			in.LogPersister.Errorf("Failed to create ECS task set %s: %v", *serviceDefinition.ServiceName, err)
			return false
		}
		// Store created ACTIVE TaskSet (CANARY variant) to delete later.
		if err := in.MetadataStore.Set(ctx, canaryTaskSetARNKeyName, *taskSet.TaskSetArn); err != nil {
			in.LogPersister.Errorf("Unable to store created active taskSet to metadata store: %v", err)
			return false
		}
	}

	in.LogPersister.Infof("Successfully applied the service definition and the task definition for ECS service %s and task definition of family %s", *serviceDefinition.ServiceName, *taskDefinition.Family)
	return true
}

func clean(ctx context.Context, in *executor.Input, cloudProviderName string, cloudProviderCfg *config.CloudProviderECSConfig, service types.Service) bool {
	client, err := provider.DefaultRegistry().Client(cloudProviderName, cloudProviderCfg, in.Logger)
	if err != nil {
		in.LogPersister.Errorf("Unable to create ECS client for the provider %s: %v", cloudProviderName, err)
		return false
	}

	taskSetArn, ok := in.MetadataStore.Get(canaryTaskSetARNKeyName)
	if !ok {
		in.LogPersister.Errorf("Unable to restore CANARY task set to clean: Not found")
		return false
	}

	if err := client.DeleteTaskSet(ctx, service, taskSetArn); err != nil {
		in.LogPersister.Errorf("Failed to clean CANARY task set %s: %v", taskSetArn, err)
		return false
	}

	return true
}

func routing(ctx context.Context, in *executor.Input, cloudProviderName string, cloudProviderCfg *config.CloudProviderECSConfig, primaryTargetGroup types.LoadBalancer, canaryTargetGroup types.LoadBalancer) bool {
	client, err := provider.DefaultRegistry().Client(cloudProviderName, cloudProviderCfg, in.Logger)
	if err != nil {
		in.LogPersister.Errorf("Unable to create ECS client for the provider %s: %v", cloudProviderName, err)
		return false
	}

	primary, canary, err := determineTrafficAmount(in.StageConfig)
	if err != nil {
		in.LogPersister.Errorf("Malformed configuration for stage %s", in.Stage.Name)
		return false
	}
	routingTrafficCfg := provider.RoutingTrafficConfig{
		{
			TargetGroupArn: *primaryTargetGroup.TargetGroupArn,
			Weight:         primary,
		},
		{
			TargetGroupArn: *canaryTargetGroup.TargetGroupArn,
			Weight:         canary,
		},
	}

	currListenerArn, err := client.GetListener(ctx, primaryTargetGroup)
	if err != nil {
		in.LogPersister.Errorf("Failed to get current active listener: %v", err)
		return false
	}

	if err := client.ModifyListener(ctx, currListenerArn, routingTrafficCfg); err != nil {
		in.LogPersister.Errorf("Failed to routing traffic to canary variant: %v", err)
		return false
	}

	return true
}

func determineTrafficAmount(stageCfg config.PipelineStage) (primary int, canary int, err error) {
	switch stageCfg.Name {
	case model.StageECSCanaryRollout:
		options := stageCfg.ECSCanaryRolloutStageOptions
		if options == nil {
			err = fmt.Errorf("traffic configuration is missing")
			return
		}
		// TODO: Validate Traffic config is lower than 100.
		canary = options.Traffic
		primary = 100 - canary
		return
	case model.StageECSPrimaryRollout:
		primary = 100
		canary = 0
		return
	default:
		err = fmt.Errorf("unexpected stage given")
		return
	}
}
