// Copyright 2022 The PipeCD Authors.
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
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/piped/deploysource"
	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	"github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/ecs"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	canaryTaskSetARNKeyName = "canary-taskset-arn"
	canaryServiceKeyName    = "canary-service-object"
	// Stage metadata keys.
	trafficRoutePrimaryMetadataKey = "primary-percentage"
	trafficRouteCanaryMetadataKey  = "canary-percentage"
	canaryScaleMetadataKey         = "canary-scale"
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
	r.Register(model.StageECSCanaryRollout, f)
	r.Register(model.StageECSPrimaryRollout, f)
	r.Register(model.StageECSCanaryClean, f)
	r.Register(model.StageECSTrafficRouting, f)

	r.RegisterRollback(model.ApplicationKind_ECS, func(in executor.Input) executor.Executor {
		return &rollbackExecutor{
			Input: in,
		}
	})
}

func findPlatformProvider(in *executor.Input) (name string, cfg *config.PlatformProviderECSConfig, found bool) {
	name = in.Application.PlatformProvider
	if name == "" {
		in.LogPersister.Errorf("Missing the PlatformProvider name in the application configuration")
		return
	}

	cp, ok := in.PipedConfig.FindPlatformProvider(name, model.ApplicationKind_ECS)
	if !ok {
		in.LogPersister.Errorf("The specified platform provider %q was not found in piped configuration", name)
		return
	}

	cfg = cp.ECSConfig
	found = true
	return
}

func loadServiceDefinition(in *executor.Input, serviceDefinitionFile string, ds *deploysource.DeploySource) (types.Service, bool) {
	in.LogPersister.Infof("Loading service manifest at commit %s", ds.Revision)

	serviceDefinition, err := provider.LoadServiceDefinition(ds.AppDir, serviceDefinitionFile)
	if err != nil {
		in.LogPersister.Errorf("Failed to load ECS service definition (%v)", err)
		return types.Service{}, false
	}

	serviceDefinition.Tags = append(
		serviceDefinition.Tags,
		provider.MakeTags(map[string]string{
			provider.LabelManagedBy:   provider.ManagedByPiped,
			provider.LabelPiped:       in.PipedConfig.PipedID,
			provider.LabelApplication: in.Deployment.ApplicationId,
			provider.LabelCommitHash:  in.Deployment.CommitHash(),
		})...,
	)

	in.LogPersister.Infof("Successfully loaded the ECS service definition at commit %s", ds.Revision)
	return serviceDefinition, true
}

func loadTaskDefinition(in *executor.Input, taskDefinitionFile string, ds *deploysource.DeploySource) (types.TaskDefinition, bool) {
	in.LogPersister.Infof("Loading task definition manifest at commit %s", ds.Revision)

	taskDefinition, err := provider.LoadTaskDefinition(ds.AppDir, taskDefinitionFile)
	if err != nil {
		in.LogPersister.Errorf("Failed to load ECS task definition (%v)", err)
		return types.TaskDefinition{}, false
	}

	in.LogPersister.Infof("Successfully loaded the ECS task definition at commit %s", ds.Revision)
	return taskDefinition, true
}

func loadTargetGroups(in *executor.Input, appCfg *config.ECSApplicationSpec, ds *deploysource.DeploySource) (*types.LoadBalancer, *types.LoadBalancer, bool) {
	in.LogPersister.Infof("Loading target groups config at the commit %s", ds.Revision)

	primary, canary, err := provider.LoadTargetGroups(appCfg.Input.TargetGroups)
	if err != nil && !errors.Is(err, provider.ErrNoTargetGroup) {
		in.LogPersister.Errorf("Failed to load TargetGroups (%v)", err)
		return nil, nil, false
	}

	if errors.Is(err, provider.ErrNoTargetGroup) {
		in.LogPersister.Infof("No target groups were set at commit %s", ds.Revision)
		return nil, nil, true
	}

	in.LogPersister.Infof("Successfully loaded the ECS target groups at commit %s", ds.Revision)
	return primary, canary, true
}

func applyTaskDefinition(ctx context.Context, cli provider.Client, taskDefinition types.TaskDefinition) (*types.TaskDefinition, error) {
	td, err := cli.RegisterTaskDefinition(ctx, taskDefinition)
	if err != nil {
		return nil, fmt.Errorf("unable to register ECS task definition of family %s: %v", *taskDefinition.Family, err)
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
		if err := cli.TagResource(ctx, *service.ServiceArn, serviceDefinition.Tags); err != nil {
			return nil, fmt.Errorf("failed to update tags of ECS service %s: %v", *serviceDefinition.ServiceName, err)
		}

	} else {
		service, err = cli.CreateService(ctx, serviceDefinition)
		if err != nil {
			return nil, fmt.Errorf("failed to create ECS service %s: %v", *serviceDefinition.ServiceName, err)
		}
	}

	return service, nil
}

func runStandaloneTask(
	ctx context.Context,
	in *executor.Input,
	cloudProviderName string,
	cloudProviderCfg *config.PlatformProviderECSConfig,
	taskDefinition types.TaskDefinition,
	ecsInput *config.ECSDeploymentInput,
) bool {
	client, err := provider.DefaultRegistry().Client(cloudProviderName, cloudProviderCfg, in.Logger)
	if err != nil {
		in.LogPersister.Errorf("Unable to create ECS client for the provider %s: %v", cloudProviderName, err)
		return false
	}

	in.LogPersister.Infof("Start applying the ECS task definition")
	tags := provider.MakeTags(map[string]string{
		provider.LabelManagedBy:   provider.ManagedByPiped,
		provider.LabelPiped:       in.PipedConfig.PipedID,
		provider.LabelApplication: in.Deployment.ApplicationId,
		provider.LabelCommitHash:  in.Deployment.CommitHash(),
	})
	td, err := applyTaskDefinition(ctx, client, taskDefinition)
	if err != nil {
		in.LogPersister.Errorf("Failed to apply ECS task definition: %v", err)
		return false
	}

	err = client.RunTask(
		ctx,
		*td,
		ecsInput.ClusterArn,
		ecsInput.LaunchType,
		&ecsInput.AwsVpcConfiguration,
		tags,
	)
	if err != nil {
		in.LogPersister.Errorf("Failed to run ECS task: %v", err)
		return false
	}
	return true
}

func createPrimaryTaskSet(ctx context.Context, client provider.Client, service types.Service, taskDef types.TaskDefinition, targetGroup *types.LoadBalancer) error {
	// Get current PRIMARY task set.
	prevPrimaryTaskSet, err := client.GetPrimaryTaskSet(ctx, service)
	// Ignore error in case it's not found error, the prevPrimaryTaskSet doesn't exist for newly created Service.
	if err != nil && !errors.Is(err, platformprovider.ErrNotFound) {
		return err
	}

	// Create a task set in the specified cluster and service.
	// In case of creating Primary taskset, the number of desired tasks scale is always set to 100
	// which means we create as many tasks as the current primary taskset has.
	taskSet, err := client.CreateTaskSet(ctx, service, taskDef, targetGroup, 100)
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

func sync(ctx context.Context, in *executor.Input, platformProviderName string, platformProviderCfg *config.PlatformProviderECSConfig, taskDefinition types.TaskDefinition, serviceDefinition types.Service, targetGroup *types.LoadBalancer) bool {
	client, err := provider.DefaultRegistry().Client(platformProviderName, platformProviderCfg, in.Logger)
	if err != nil {
		in.LogPersister.Errorf("Unable to create ECS client for the provider %s: %v", platformProviderName, err)
		return false
	}

	in.LogPersister.Infof("Start applying the ECS task definition")
	td, err := applyTaskDefinition(ctx, client, taskDefinition)
	if err != nil {
		in.LogPersister.Errorf("Failed to apply ECS task definition: %v", err)
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
		in.LogPersister.Errorf("Failed to rolling out ECS task set for service %s: %v", *serviceDefinition.ServiceName, err)
		return false
	}

	in.LogPersister.Infof("Successfully applied the service definition and the task definition for ECS service %s and task definition of family %s", *serviceDefinition.ServiceName, *taskDefinition.Family)
	return true
}

func rollout(ctx context.Context, in *executor.Input, platformProviderName string, platformProviderCfg *config.PlatformProviderECSConfig, taskDefinition types.TaskDefinition, serviceDefinition types.Service, targetGroup *types.LoadBalancer) bool {
	client, err := provider.DefaultRegistry().Client(platformProviderName, platformProviderCfg, in.Logger)
	if err != nil {
		in.LogPersister.Errorf("Unable to create ECS client for the provider %s: %v", platformProviderName, err)
		return false
	}

	in.LogPersister.Infof("Start applying the ECS task definition")
	td, err := applyTaskDefinition(ctx, client, taskDefinition)
	if err != nil {
		in.LogPersister.Errorf("Failed to apply ECS task definition: %v", err)
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
			in.LogPersister.Errorf("Failed to rolling out ECS task set for service %s: %v", *serviceDefinition.ServiceName, err)
			return false
		}
	} else {
		// Load Canary rollout stage options to get scale configuration.
		options := in.StageConfig.ECSCanaryRolloutStageOptions
		if options == nil {
			in.LogPersister.Errorf("Malformed configuration for stage %s", in.Stage.Name)
			return false
		}

		metadata := map[string]string{
			canaryScaleMetadataKey: strconv.FormatInt(int64(options.Scale.Int()), 10),
		}
		if err := in.MetadataStore.Stage(in.Stage.Id).PutMulti(ctx, metadata); err != nil {
			in.Logger.Error("Failed to store canary scale infor to metadata store", zap.Error(err))
		}

		// Create ACTIVE task set in case of Canary rollout.
		taskSet, err := client.CreateTaskSet(ctx, *service, *td, targetGroup, options.Scale.Int())
		if err != nil {
			in.LogPersister.Errorf("Failed to create ECS task set for service %s: %v", *serviceDefinition.ServiceName, err)
			return false
		}
		// Store created ACTIVE TaskSet (CANARY variant) to delete later.
		if err := in.MetadataStore.Shared().Put(ctx, canaryTaskSetARNKeyName, *taskSet.TaskSetArn); err != nil {
			in.LogPersister.Errorf("Unable to store created active taskSet to metadata store: %v", err)
			return false
		}
		// Store applied Service (CANARY variant) to delete its TaskSet later.
		serviceObjData, err := json.Marshal(service)
		if err != nil {
			in.LogPersister.Errorf("Unable to store applied service to metadata store: %v", err)
			return false
		}
		if err := in.MetadataStore.Shared().Put(ctx, canaryServiceKeyName, string(serviceObjData)); err != nil {
			in.LogPersister.Errorf("Unable to store applied service to metadata store: %v", err)
			return false
		}
	}

	in.LogPersister.Infof("Successfully applied the service definition and the task definition for ECS service %s and task definition of family %s", *serviceDefinition.ServiceName, *taskDefinition.Family)
	return true
}

func clean(ctx context.Context, in *executor.Input, platformProviderName string, platformProviderCfg *config.PlatformProviderECSConfig) bool {
	client, err := provider.DefaultRegistry().Client(platformProviderName, platformProviderCfg, in.Logger)
	if err != nil {
		in.LogPersister.Errorf("Unable to create ECS client for the provider %s: %v", platformProviderName, err)
		return false
	}

	taskSetArn, ok := in.MetadataStore.Shared().Get(canaryTaskSetARNKeyName)
	if !ok {
		in.LogPersister.Errorf("Unable to restore CANARY task set to clean: Not found")
		return false
	}
	serviceObjData, ok := in.MetadataStore.Shared().Get(canaryServiceKeyName)
	if !ok {
		in.LogPersister.Errorf("Unable to restore CANARY service to clean: Not found")
		return false
	}
	service := &types.Service{}
	if err := json.Unmarshal([]byte(serviceObjData), service); err != nil {
		in.LogPersister.Errorf("Unable to restore CANARY service to clean: %v", err)
		return false
	}

	if err := client.DeleteTaskSet(ctx, *service, taskSetArn); err != nil {
		in.LogPersister.Errorf("Failed to clean CANARY task set %s: %v", taskSetArn, err)
		return false
	}

	return true
}

func routing(ctx context.Context, in *executor.Input, platformProviderName string, platformProviderCfg *config.PlatformProviderECSConfig, primaryTargetGroup types.LoadBalancer, canaryTargetGroup types.LoadBalancer) bool {
	client, err := provider.DefaultRegistry().Client(platformProviderName, platformProviderCfg, in.Logger)
	if err != nil {
		in.LogPersister.Errorf("Unable to create ECS client for the provider %s: %v", platformProviderName, err)
		return false
	}

	options := in.StageConfig.ECSTrafficRoutingStageOptions
	if options == nil {
		in.LogPersister.Errorf("Malformed configuration for stage %s", in.Stage.Name)
		return false
	}
	primary, canary := options.Percentage()
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

	metadata := map[string]string{
		trafficRoutePrimaryMetadataKey: strconv.FormatInt(int64(primary), 10),
		trafficRouteCanaryMetadataKey:  strconv.FormatInt(int64(canary), 10),
	}
	if err := in.MetadataStore.Stage(in.Stage.Id).PutMulti(ctx, metadata); err != nil {
		in.Logger.Error("Failed to store traffic routing config to metadata store", zap.Error(err))
	}

	currListenerArn, err := client.GetListener(ctx, primaryTargetGroup)
	if err != nil {
		in.LogPersister.Errorf("Failed to get current active listener: %v", err)
		return false
	}

	if err := client.ModifyListener(ctx, currListenerArn, routingTrafficCfg); err != nil {
		in.LogPersister.Errorf("Failed to routing traffic to CANARY variant: %v", err)
		return false
	}

	return true
}
