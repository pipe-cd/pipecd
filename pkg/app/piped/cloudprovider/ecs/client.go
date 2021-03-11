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

func newClient(region, profile, credentialsFile, roleARN, tokenPath string, logger *zap.Logger) (*client, error) {
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

func (c *client) CreateService(ctx context.Context, service types.Service) (*ecs.CreateServiceOutput, error) {
	input := &ecs.CreateServiceInput{
		ServiceName:                   service.ServiceName,
		Cluster:                       service.ClusterArn,
		DeploymentConfiguration:       service.DeploymentConfiguration,
		DeploymentController:          service.DeploymentController,
		DesiredCount:                  aws.Int32(service.DesiredCount),
		EnableECSManagedTags:          service.EnableECSManagedTags,
		HealthCheckGracePeriodSeconds: service.HealthCheckGracePeriodSeconds,
		LaunchType:                    service.LaunchType,
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
	output, err := c.client.CreateService(ctx, input)
	if err != nil {
		return &ecs.CreateServiceOutput{}, fmt.Errorf("failed to create ECS service %s: %w", *service.ServiceName, err)
	}
	return output, nil
}

func (c *client) UpdateService(ctx context.Context, service types.Service) (*ecs.UpdateServiceOutput, error) {
	input := &ecs.UpdateServiceInput{
		Service:                       service.ServiceName,
		Cluster:                       service.ClusterArn,
		DeploymentConfiguration:       service.DeploymentConfiguration,
		DesiredCount:                  &service.DesiredCount,
		HealthCheckGracePeriodSeconds: service.HealthCheckGracePeriodSeconds,
		NetworkConfiguration:          service.NetworkConfiguration,
		PlacementConstraints:          service.PlacementConstraints,
		PlacementStrategy:             service.PlacementStrategy,
		PlatformVersion:               service.PlatformVersion,
		TaskDefinition:                service.TaskDefinition,
	}
	output, err := c.client.UpdateService(ctx, input)
	if err != nil {
		return &ecs.UpdateServiceOutput{}, fmt.Errorf("failed to update ECS service %s: %w", *service.ServiceName, err)
	}
	return output, nil
}

func (c *client) RegisterTaskDefinition(ctx context.Context, taskDefinition types.TaskDefinition) (*ecs.RegisterTaskDefinitionOutput, error) {
	input := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: taskDefinition.ContainerDefinitions,
		Family:               taskDefinition.Family,
	}
	output, err := c.client.RegisterTaskDefinition(ctx, input)
	if err != nil {
		return &ecs.RegisterTaskDefinitionOutput{}, fmt.Errorf("failed to register ECS task definition %s: %w", *taskDefinition.TaskDefinitionArn, err)
	}
	return output, nil
}

func (c *client) DeregisterTaskDefinition(ctx context.Context, taskDefinition types.TaskDefinition) (*ecs.DeregisterTaskDefinitionOutput, error) {
	input := &ecs.DeregisterTaskDefinitionInput{
		TaskDefinition: taskDefinition.TaskDefinitionArn,
	}
	output, err := c.client.DeregisterTaskDefinition(ctx, input)
	if err != nil {
		return &ecs.DeregisterTaskDefinitionOutput{}, fmt.Errorf("failed to deregister ECS task definition %s: %w", *taskDefinition.TaskDefinitionArn, err)
	}
	return output, nil
}

func (c *client) CreateTaskSet(ctx context.Context, service types.Service, taskDefinition types.TaskDefinition) (*ecs.CreateTaskSetOutput, error) {
	input := &ecs.CreateTaskSetInput{
		Cluster:        service.ClusterArn,
		Service:        service.ServiceArn,
		TaskDefinition: taskDefinition.TaskDefinitionArn,
	}
	output, err := c.client.CreateTaskSet(ctx, input)
	if err != nil {
		return &ecs.CreateTaskSetOutput{}, fmt.Errorf("failed to create ECS task set %s: %w", *taskDefinition.TaskDefinitionArn, err)
	}
	return output, nil
}

func (c *client) DeleteTaskSet(ctx context.Context, service types.Service, taskSet types.TaskSet) (*ecs.DeleteTaskSetOutput, error) {
	input := &ecs.DeleteTaskSetInput{
		Cluster: service.ClusterArn,
		Service: service.ServiceArn,
		TaskSet: taskSet.TaskSetArn,
	}
	output, err := c.client.DeleteTaskSet(ctx, input)
	if err != nil {
		return &ecs.DeleteTaskSetOutput{}, fmt.Errorf("failed to delete ECS task set %s: %w", *taskSet.TaskSetArn, err)
	}
	return output, nil
}

func (c *client) ServiceExists(ctx context.Context, clusterName string, services []string) (bool, error) {
	input := &ecs.DescribeServicesInput{
		Cluster:  aws.String(clusterName),
		Services: services,
	}
	_, err := c.client.DescribeServices(ctx, input)
	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			// Only in case ResourceNotFound error occurred, the FunctionName is available for create so do not raise error.
			return false, nil
		}
		return false, err
	}
	return true, nil
}
