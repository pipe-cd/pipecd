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
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbtypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider"
	appconfig "github.com/pipe-cd/pipecd/pkg/config"
)

type client struct {
	ecsClient *ecs.Client
	elbClient *elasticloadbalancingv2.Client
	logger    *zap.Logger
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
	c.ecsClient = ecs.NewFromConfig(cfg)
	c.elbClient = elasticloadbalancingv2.NewFromConfig(cfg)

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
		PlacementConstraints:          service.PlacementConstraints,
		PlacementStrategy:             service.PlacementStrategy,
		PlatformVersion:               service.PlatformVersion,
		PropagateTags:                 service.PropagateTags,
		Role:                          service.RoleArn,
		SchedulingStrategy:            service.SchedulingStrategy,
		ServiceRegistries:             service.ServiceRegistries,
		Tags:                          service.Tags,
	}

	output, err := c.ecsClient.CreateService(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create ECS service %s: %w", *service.ServiceName, err)
	}

	// Hack: Since we use EXTERNAL deployment controller, the below configurations are not allowed to be passed
	// in CreateService step, but it required in further step (CreateTaskSet step). We reassign those values
	// as part of service definition for that purpose.
	output.Service.LaunchType = service.LaunchType
	output.Service.NetworkConfiguration = service.NetworkConfiguration

	return output.Service, nil
}

func (c *client) UpdateService(ctx context.Context, service types.Service) (*types.Service, error) {
	input := &ecs.UpdateServiceInput{
		Cluster:           service.ClusterArn,
		Service:           service.ServiceName,
		DesiredCount:      aws.Int32(service.DesiredCount),
		PlacementStrategy: service.PlacementStrategy,
		// TODO: Support update other properties of service.
		// PlacementConstraints:    service.PlacementConstraints,
	}
	output, err := c.ecsClient.UpdateService(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update ECS service %s: %w", *service.ServiceName, err)
	}

	// Hack: Since we use EXTERNAL deployment controller, the below configurations are not allowed to be passed
	// in UpdateService step, but it required in further step (CreateTaskSet step). We reassign those values
	// as part of service definition for that purpose.
	output.Service.LaunchType = service.LaunchType
	output.Service.NetworkConfiguration = service.NetworkConfiguration

	return output.Service, nil
}

func (c *client) RegisterTaskDefinition(ctx context.Context, taskDefinition types.TaskDefinition) (*types.TaskDefinition, error) {
	input := &ecs.RegisterTaskDefinitionInput{
		Family:                  taskDefinition.Family,
		ContainerDefinitions:    taskDefinition.ContainerDefinitions,
		RequiresCompatibilities: taskDefinition.RequiresCompatibilities,
		ExecutionRoleArn:        taskDefinition.ExecutionRoleArn,
		TaskRoleArn:             taskDefinition.TaskRoleArn,
		NetworkMode:             taskDefinition.NetworkMode,
		Volumes:                 taskDefinition.Volumes,
		// Requires defined at task level in case Fargate is used.
		Cpu:    taskDefinition.Cpu,
		Memory: taskDefinition.Memory,
		// TODO: Support tags for registering task definition.
	}
	output, err := c.ecsClient.RegisterTaskDefinition(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to register ECS task definition of family %s: %w", *taskDefinition.Family, err)
	}
	return output.TaskDefinition, nil
}

func (c *client) RunTask(ctx context.Context, taskDefinition types.TaskDefinition, clusterArn string, launchType string, awsVpcConfiguration *appconfig.ECSVpcConfiguration) error {
	if taskDefinition.TaskDefinitionArn == nil {
		return fmt.Errorf("failed to run task of task family %s: no task definition provided", *taskDefinition.Family)
	}

	input := &ecs.RunTaskInput{
		TaskDefinition: taskDefinition.Family,
		Cluster:        aws.String(clusterArn),
		LaunchType:     types.LaunchType(launchType),
	}

	if len(awsVpcConfiguration.Subnets) > 0 {
		input.NetworkConfiguration = &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				Subnets:        awsVpcConfiguration.Subnets,
				AssignPublicIp: types.AssignPublicIp(awsVpcConfiguration.AssignPublicIP),
				SecurityGroups: awsVpcConfiguration.SecurityGroups,
			},
		}
	}

	_, err := c.ecsClient.RunTask(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to run ECS task %s: %w", *taskDefinition.TaskDefinitionArn, err)
	}
	return nil
}

func (c *client) CreateTaskSet(ctx context.Context, service types.Service, taskDefinition types.TaskDefinition, targetGroup *types.LoadBalancer, scale int) (*types.TaskSet, error) {
	if taskDefinition.TaskDefinitionArn == nil {
		return nil, fmt.Errorf("failed to create task set of task family %s: no task definition provided", *taskDefinition.Family)
	}
	input := &ecs.CreateTaskSetInput{
		Cluster:        service.ClusterArn,
		Service:        service.ServiceArn,
		TaskDefinition: taskDefinition.TaskDefinitionArn,
		Scale:          &types.Scale{Unit: types.ScaleUnitPercent, Value: float64(scale)},
		// If you specify the awsvpc network mode, the task is allocated an elastic network interface,
		// and you must specify a NetworkConfiguration when run a task with the task definition.
		NetworkConfiguration: service.NetworkConfiguration,
		LaunchType:           service.LaunchType,
	}
	if targetGroup != nil {
		input.LoadBalancers = []types.LoadBalancer{*targetGroup}
	}
	output, err := c.ecsClient.CreateTaskSet(ctx, input)
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
	output, err := c.ecsClient.DescribeServices(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get primary task set of service %s: %w", *service.ServiceName, err)
	}
	if len(output.Services) == 0 {
		return nil, fmt.Errorf("failed to get primary task set of service %s: services empty", *service.ServiceName)
	}
	taskSets := output.Services[0].TaskSets
	for _, taskSet := range taskSets {
		if aws.ToString(taskSet.Status) == "PRIMARY" {
			return &taskSet, nil
		}
	}
	return nil, platformprovider.ErrNotFound
}

func (c *client) DeleteTaskSet(ctx context.Context, service types.Service, taskSetArn string) error {
	input := &ecs.DeleteTaskSetInput{
		Cluster: service.ClusterArn,
		Service: service.ServiceArn,
		TaskSet: aws.String(taskSetArn),
	}
	if _, err := c.ecsClient.DeleteTaskSet(ctx, input); err != nil {
		return fmt.Errorf("failed to delete ECS task set %s: %w", taskSetArn, err)
	}
	return nil
}

func (c *client) UpdateServicePrimaryTaskSet(ctx context.Context, service types.Service, taskSet types.TaskSet) (*types.TaskSet, error) {
	input := &ecs.UpdateServicePrimaryTaskSetInput{
		Cluster:        service.ClusterArn,
		Service:        service.ServiceArn,
		PrimaryTaskSet: taskSet.TaskSetArn,
	}
	output, err := c.ecsClient.UpdateServicePrimaryTaskSet(ctx, input)
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
	output, err := c.ecsClient.DescribeServices(ctx, input)
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

func (c *client) GetListener(ctx context.Context, targetGroup types.LoadBalancer) (string, error) {
	loadBalancerArn, err := c.getLoadBalancerArn(ctx, *targetGroup.TargetGroupArn)
	if err != nil {
		return "", err
	}

	input := &elasticloadbalancingv2.DescribeListenersInput{
		LoadBalancerArn: aws.String(loadBalancerArn),
	}
	output, err := c.elbClient.DescribeListeners(ctx, input)
	if err != nil {
		return "", err
	}
	if len(output.Listeners) == 0 {
		return "", platformprovider.ErrNotFound
	}
	// Note: Suppose the load balancer only have one listener.
	// TODO: Support multi listeners pattern.
	if len(output.Listeners) > 1 {
		return "", fmt.Errorf("invalid listener configuration pointed to %s target group", *targetGroup.TargetGroupArn)
	}

	return *output.Listeners[0].ListenerArn, nil
}

func (c *client) getLoadBalancerArn(ctx context.Context, targetGroupArn string) (string, error) {
	input := &elasticloadbalancingv2.DescribeTargetGroupsInput{
		TargetGroupArns: []string{targetGroupArn},
	}
	output, err := c.elbClient.DescribeTargetGroups(ctx, input)
	if err != nil {
		return "", err
	}
	if len(output.TargetGroups) == 0 {
		return "", platformprovider.ErrNotFound
	}
	// Note: Currently, only support TargetGroup which serves traffic from one Load Balancer.
	return output.TargetGroups[0].LoadBalancerArns[0], nil
}

func (c *client) ModifyListener(ctx context.Context, listenerArn string, routingTrafficCfg RoutingTrafficConfig) error {
	if len(routingTrafficCfg) != 2 {
		return fmt.Errorf("invalid listener configuration: requires 2 target groups")
	}
	input := &elasticloadbalancingv2.ModifyListenerInput{
		ListenerArn: aws.String(listenerArn),
		DefaultActions: []elbtypes.Action{
			{
				Type: elbtypes.ActionTypeEnumForward,
				ForwardConfig: &elbtypes.ForwardActionConfig{
					TargetGroups: []elbtypes.TargetGroupTuple{
						{
							TargetGroupArn: aws.String(routingTrafficCfg[0].TargetGroupArn),
							Weight:         aws.Int32(int32(routingTrafficCfg[0].Weight)),
						},
						{
							TargetGroupArn: aws.String(routingTrafficCfg[1].TargetGroupArn),
							Weight:         aws.Int32(int32(routingTrafficCfg[1].Weight)),
						},
					},
				},
			},
		},
	}
	_, err := c.elbClient.ModifyListener(ctx, input)
	return err
}
