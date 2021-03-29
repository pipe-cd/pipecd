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
	region string
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
	input := &ecs.CreateServiceInput{
		ServiceName:                   service.ServiceName,
		Cluster:                       service.ClusterArn,
		DeploymentConfiguration:       service.DeploymentConfiguration,
		DeploymentController:          service.DeploymentController,
		DesiredCount:                  aws.Int32(service.DesiredCount),
		EnableECSManagedTags:          service.EnableECSManagedTags,
		HealthCheckGracePeriodSeconds: service.HealthCheckGracePeriodSeconds,
		LoadBalancers:                 service.LoadBalancers,
		NetworkConfiguration:          service.NetworkConfiguration,
		PlacementConstraints:          service.PlacementConstraints,
		PlacementStrategy:             service.PlacementStrategy,
		PlatformVersion:               service.PlatformVersion,
		PropagateTags:                 service.PropagateTags,
		Role:                          service.RoleArn,
		SchedulingStrategy:            service.SchedulingStrategy,
		ServiceRegistries:             service.ServiceRegistries,
		Tags:                          service.Tags,
		TaskDefinition:                service.TaskDefinition,
	}
	if service.DeploymentController.Type != types.DeploymentControllerTypeExternal {
		input.LaunchType = service.LaunchType
	}
	output, err := c.client.CreateService(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create ECS service %s: %w", *service.ServiceName, err)
	}
	return output.Service, nil
}

func (c *client) UpdateService(ctx context.Context, service types.Service) (*types.Service, error) {
	input := &ecs.UpdateServiceInput{
		Service:                 service.ServiceName,
		Cluster:                 service.ClusterArn,
		DeploymentConfiguration: service.DeploymentConfiguration,
		DesiredCount:            aws.Int32(service.DesiredCount),
		NetworkConfiguration:    service.NetworkConfiguration,
		PlacementConstraints:    service.PlacementConstraints,
		PlacementStrategy:       service.PlacementStrategy,
		TaskDefinition:          service.TaskDefinition,
	}
	output, err := c.client.UpdateService(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update ECS service %s: %w", *service.ServiceName, err)
	}
	return output.Service, nil
}

func (c *client) RegisterTaskDefinition(ctx context.Context, taskDefinition types.TaskDefinition) (*types.TaskDefinition, error) {
	input := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: taskDefinition.ContainerDefinitions,
		Family:               taskDefinition.Family,
	}
	output, err := c.client.RegisterTaskDefinition(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to register ECS task definition %s: %w", *taskDefinition.TaskDefinitionArn, err)
	}
	return output.TaskDefinition, nil
}

func (c *client) DeregisterTaskDefinition(ctx context.Context, taskDefinition types.TaskDefinition) (*types.TaskDefinition, error) {
	input := &ecs.DeregisterTaskDefinitionInput{
		TaskDefinition: taskDefinition.TaskDefinitionArn,
	}
	output, err := c.client.DeregisterTaskDefinition(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to deregister ECS task definition %s: %w", *taskDefinition.TaskDefinitionArn, err)
	}
	return output.TaskDefinition, nil
}

func (c *client) CreateTaskSet(ctx context.Context, service types.Service, taskDefinition types.TaskDefinition) (*types.TaskSet, error) {
	input := &ecs.CreateTaskSetInput{
		Cluster:        service.ClusterArn,
		Service:        service.ServiceArn,
		TaskDefinition: taskDefinition.TaskDefinitionArn,
		Scale:          &types.Scale{Unit: types.ScaleUnitPercent, Value: 100},
	}
	output, err := c.client.CreateTaskSet(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create ECS task set %s: %w", *taskDefinition.TaskDefinitionArn, err)
	}
	return output.TaskSet, nil
}

func (c *client) DeleteTaskSet(ctx context.Context, service types.Service, taskSet types.TaskSet) (*types.TaskSet, error) {
	input := &ecs.DeleteTaskSetInput{
		Cluster: service.ClusterArn,
		Service: service.ServiceArn,
		TaskSet: taskSet.TaskSetArn,
	}
	output, err := c.client.DeleteTaskSet(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to delete ECS task set %s: %w", *taskSet.TaskSetArn, err)
	}
	return output.TaskSet, nil
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
