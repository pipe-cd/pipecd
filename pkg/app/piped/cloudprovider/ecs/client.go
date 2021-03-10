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
		region: region,
		logger: logger.Named("ecs"),
	}

	optFns := []func(*config.LoadOptions) error{config.WithRegion(region)}
	if credentialsFile != "" {
		optFns = append(optFns, config.WithSharedCredentialsFiles([]string{credentialsFile}))
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

func (c *client) CreateService(ctx context.Context, service types.Service) error {
	input := &ecs.CreateServiceInput{
		ServiceName:                   service.ServiceName,
		Cluster:                       service.ClusterArn,
		DeploymentConfiguration:       service.DeploymentConfiguration,
		DeploymentController:          service.DeploymentController,
		DesiredCount:                  &service.DesiredCount,
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
	_, err := c.client.CreateService(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create ECS service %s: %w", *service.ServiceName, err)
	}
	return nil
}

func (c *client) UpdateService(ctx context.Context, service types.Service) error {
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
	_, err := c.client.UpdateService(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update ECS service %s: %w", *service.ServiceName, err)
	}
	return nil
}

func (c *client) ServiceExist(ctx context.Context, clusterName string, services []string) (bool, error) {
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
