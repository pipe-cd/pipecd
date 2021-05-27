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
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"go.uber.org/zap"
)

type client struct {
	client *ecs.Client
	logger *zap.Logger
}

func newClient(region, profile, credentialsFile, roleARN, tokenPath string, logger *zap.Logger) (Client, error) {
	if region == "" {
		return nil, fmt.Errorf("region is required field")
	}

	c := &client{
		logger: logger.Named("ecs"),
	}

	optFns := []func(*config.LoadOptions) error{config.WithRegion(region)}
	if credentialsFile != "" {
		optFns = append(optFns, config.WithSharedCredentialsFiles([]string{credentialsFile}))
	}
	if profile != "" {
		optFns = append(optFns, config.WithSharedConfigProfile(profile))
	}
	if tokenPath != "" && roleARN != "" {
		optFns = append(optFns, config.WithWebIdentityRoleCredentialOptions(func(v *stscreds.WebIdentityRoleOptions) {
			v.RoleARN = roleARN
			v.TokenRetriever = stscreds.IdentityTokenFile(tokenPath)
		}))
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), optFns...)
	if err != nil {
		return nil, fmt.Errorf("failed to load config to create ecs client: %w", err)
	}
	c.client = ecs.NewFromConfig(cfg)

	return c, nil
}

func (c *client) CreateService(ctx context.Context, service types.Service) (*types.Service, error) {
	if service.DeploymentController == nil || service.DeploymentController.Type != types.DeploymentControllerTypeExternal {
		return nil, fmt.Errorf("failed to create ECS service %s: deployment controller of type EXTERNAL is required", *service.ServiceName)
	}
	input := &ecs.CreateServiceInput{
		Cluster:                       service.ClusterArn,
		ServiceName:                   service.ServiceName,
		DesiredCount:                  aws.Int32(service.DesiredCount),
		DeploymentController:          service.DeploymentController,
		DeploymentConfiguration:       service.DeploymentConfiguration,
		EnableECSManagedTags:          service.EnableECSManagedTags,
		HealthCheckGracePeriodSeconds: service.HealthCheckGracePeriodSeconds,
		LoadBalancers:                 service.LoadBalancers,
		PlacementConstraints:          service.PlacementConstraints,
		PlacementStrategy:             service.PlacementStrategy,
		PlatformVersion:               service.PlatformVersion,
		PropagateTags:                 service.PropagateTags,
		Role:                          service.RoleArn,
		SchedulingStrategy:            service.SchedulingStrategy,
		ServiceRegistries:             service.ServiceRegistries,
		Tags:                          service.Tags,
	}

	output, err := c.client.CreateService(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create ECS service %s: %w", *service.ServiceName, err)
	}
	return output.Service, nil
}

func (c *client) UpdateService(ctx context.Context, service types.Service) (*types.Service, error) {
	input := &ecs.UpdateServiceInput{
		Cluster:           service.ClusterArn,
		Service:           service.ServiceName,
		DesiredCount:      aws.Int32(service.DesiredCount),
		PlacementStrategy: service.PlacementStrategy,
		// TODO: Support update other properties of service.
		// DeploymentConfiguration: service.DeploymentConfiguration,
		// NetworkConfiguration:    service.NetworkConfiguration,
		// PlacementConstraints:    service.PlacementConstraints,
	}
	output, err := c.client.UpdateService(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update ECS service %s: %w", *service.ServiceName, err)
	}
	return output.Service, nil
}

func (c *client) RegisterTaskDefinition(ctx context.Context, taskDefinition types.TaskDefinition) (*types.TaskDefinition, error) {
	input := &ecs.RegisterTaskDefinitionInput{
		Family:                  taskDefinition.Family,
		ContainerDefinitions:    taskDefinition.ContainerDefinitions,
		RequiresCompatibilities: taskDefinition.RequiresCompatibilities,
		ExecutionRoleArn:        taskDefinition.ExecutionRoleArn,
		NetworkMode:             taskDefinition.NetworkMode,
		Volumes:                 taskDefinition.Volumes,
		// Requires defined at task level in case Fargate is used.
		Cpu:    taskDefinition.Cpu,
		Memory: taskDefinition.Memory,
		// TODO: Support tags for registering task definition.
	}
	output, err := c.client.RegisterTaskDefinition(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to register ECS task definition of family %s: %w", *taskDefinition.Family, err)
	}
	return output.TaskDefinition, nil
}

func (c *client) CreateTaskSet(ctx context.Context, service types.Service, taskDefinition types.TaskDefinition) (*types.TaskSet, error) {
	if taskDefinition.TaskDefinitionArn == nil {
		return nil, fmt.Errorf("failed to create task set of task family %s: no task definition provided", *taskDefinition.Family)
	}
	input := &ecs.CreateTaskSetInput{
		Cluster:        service.ClusterArn,
		Service:        service.ServiceArn,
		TaskDefinition: taskDefinition.TaskDefinitionArn,
		// Always create a new taskSet which has as many tasks as desiredCount number set by service.
		Scale: &types.Scale{Unit: types.ScaleUnitPercent, Value: float64(100)},
		// If you specify the awsvpc network mode, the task is allocated an elastic network interface,
		// and you must specify a NetworkConfiguration when run a task with the task definition.
		// TODO: Find better way to get those 2 values instead of set it via service def.
		NetworkConfiguration: service.NetworkConfiguration,
		LaunchType:           service.LaunchType,
	}
	output, err := c.client.CreateTaskSet(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create ECS task set %s: %w", *taskDefinition.TaskDefinitionArn, err)
	}
	return output.TaskSet, nil
}

func (c *client) GetPrimaryTaskSet(ctx context.Context, service types.Service) (*types.TaskSet, error) {
	input := &ecs.DescribeServicesInput{
		Cluster: service.ClusterArn,
		Services: []string{
			*service.ServiceArn,
		},
	}
	output, err := c.client.DescribeServices(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get primary task set of service %s: %w", *service.ServiceName, err)
	}
	if len(output.Services) == 0 {
		return nil, fmt.Errorf("failed to get primary task set of service %s: services not found", *service.ServiceName)
	}
	taskSets := output.Services[0].TaskSets
	for _, taskSet := range taskSets {
		if aws.ToString(taskSet.Status) == "PRIMARY" {
			return &taskSet, nil
		}
	}
	return nil, nil
}

func (c *client) DeleteTaskSet(ctx context.Context, service types.Service, taskSet types.TaskSet) error {
	input := &ecs.DeleteTaskSetInput{
		Cluster: service.ClusterArn,
		Service: service.ServiceArn,
		TaskSet: taskSet.TaskSetArn,
	}
	if _, err := c.client.DeleteTaskSet(ctx, input); err != nil {
		return fmt.Errorf("failed to delete ECS task set %s: %w", *taskSet.TaskSetArn, err)
	}
	return nil
}

func (c *client) UpdateServicePrimaryTaskSet(ctx context.Context, service types.Service, taskSet types.TaskSet) (*types.TaskSet, error) {
	input := &ecs.UpdateServicePrimaryTaskSetInput{
		Cluster:        service.ClusterArn,
		Service:        service.ServiceArn,
		PrimaryTaskSet: taskSet.TaskSetArn,
	}
	output, err := c.client.UpdateServicePrimaryTaskSet(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update service primary ECS task set %s: %w", *taskSet.TaskSetArn, err)
	}
	return output.TaskSet, nil
}

func (c *client) ServiceExists(ctx context.Context, clusterName string, serviceName string) (bool, error) {
	input := &ecs.DescribeServicesInput{
		Cluster:  aws.String(clusterName),
		Services: []string{serviceName},
	}
	output, err := c.client.DescribeServices(ctx, input)
	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			// Only in case ResourceNotFound error occurred, the FunctionName is available for create so do not raise error.
			return false, nil
		}
		return false, err
	}
	// Note: In case of cluster's existing serviceName is set to inactive status, it's safe to recreate the service with the same serviceName.
	for _, service := range output.Services {
		if *service.ServiceName == serviceName && *service.Status == "ACTIVE" {
			return true, nil
		}
	}
	return false, nil
}
