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

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/ecs"
	"github.com/pipe-cd/pipe/pkg/app/piped/deploysource"
	"github.com/pipe-cd/pipe/pkg/app/piped/executor"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

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
	in.LogPersister.Infof("Loading service manifest at the %s commit (%s)", ds.RevisionName, ds.RevisionName)

	taskDefinition, err := provider.LoadTaskDefinition(ds.AppDir, taskDefinitionFile)
	if err != nil {
		in.LogPersister.Errorf("Failed to load ECS service definition (%v)", err)
		return types.TaskDefinition{}, false
	}

	in.LogPersister.Infof("Successfully loaded the ECS service definition at the %s commit", ds.RevisionName)
	return taskDefinition, true
}

func sync(ctx context.Context, in *executor.Input, cloudProviderName string, cloudProviderCfg *config.CloudProviderECSConfig, taskDefinition types.TaskDefinition, serviceDefinition types.Service) bool {
	in.LogPersister.Infof("Start applying the ECS task definition")
	client, err := provider.DefaultRegistry().Client(cloudProviderName, cloudProviderCfg, in.Logger)
	if err != nil {
		in.LogPersister.Errorf("Unable to create ECS client for the provider %s: %v", cloudProviderName, err)
		return false
	}

	// Build and publish new version of ECS service and task definition.
	ok := build(ctx, in, client, taskDefinition, serviceDefinition)
	if !ok {
		in.LogPersister.Errorf("Failed to build new version for ECS %s", *serviceDefinition.ServiceName)
		return false
	}

	in.LogPersister.Infof("Successfully applied the service definition and the task definition for ECS service %s and task definition of family %s", *serviceDefinition.ServiceName, *taskDefinition.Family)
	return true
}

func build(ctx context.Context, in *executor.Input, client provider.Client, taskDefinition types.TaskDefinition, serviceDefinition types.Service) bool {
	td, err := client.RegisterTaskDefinition(ctx, taskDefinition)
	if err != nil {
		in.LogPersister.Errorf("Failed to register ECS task definition of family %s: %v", *taskDefinition.Family, err)
		return false
	}

	found, err := client.ServiceExists(ctx, *serviceDefinition.ClusterArn, *serviceDefinition.ServiceName)
	if err != nil {
		in.LogPersister.Errorf("Unable to validate service name %s: %v", *serviceDefinition.ServiceName, err)
		return false
	}
	var service *types.Service
	// if serviceDefinition.DeploymentController.Type != types.DeploymentControllerTypeExternal {
	// 	serviceDefinition.TaskDefinition = td.TaskDefinitionArn
	// }

	// If the task definition is specificed in service definition, should use that sepecificed version.
	// Consider check this before register new task definition revision.
	if serviceDefinition.TaskDefinition == nil {
		serviceDefinition.TaskDefinition = td.TaskDefinitionArn
	}
	if found {
		service, err = client.UpdateService(ctx, serviceDefinition)
		if err != nil {
			in.LogPersister.Errorf("Failed to update ECS service %s: %v", *serviceDefinition.ServiceName, err)
			return false
		}
	} else {
		service, err = client.CreateService(ctx, serviceDefinition)
		if err != nil {
			in.LogPersister.Errorf("Failed to create ECS service %s: %v", *serviceDefinition.ServiceName, err)
			return false
		}
	}
	// service.TaskDefinition = td.TaskDefinitionArn
	// service.LaunchType = serviceDefinition.LaunchType
	// service.LoadBalancers = serviceDefinition.LoadBalancers

	// Create a task set in the specified cluster and service and routing traffic to that task set.
	// This is used when a service uses the EXTERNAL deployment controller type.
	// For more information, see Amazon ECS Deployment Types
	// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/deployment-types.html
	if service.DeploymentController.Type == types.DeploymentControllerTypeExternal {
		taskSet, err := client.CreateTaskSet(ctx, *service, *td, 100)
		if err != nil {
			in.LogPersister.Errorf("Failed to create ECS task set %s: %v", *serviceDefinition.ServiceName, err)
			return false
		}

		if _, err = client.UpdateServicePrimaryTaskSet(ctx, *service, *taskSet); err != nil {
			in.LogPersister.Errorf("Failed to update service primary ECS task set %s: %v", *serviceDefinition.ServiceName, err)
			return false
		}
	}

	in.LogPersister.Info("Successfully applied the service definition and the task definition")
	return true
}
