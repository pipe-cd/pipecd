// Copyright 2024 The PipeCD Authors.
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
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbtypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider"
	"github.com/pipe-cd/pipecd/pkg/backoff"
	appconfig "github.com/pipe-cd/pipecd/pkg/config"
)

const (
	// ServiceStable's constants.
	retryServiceStable         = 40
	retryServiceStableInterval = 15 * time.Second

	// TaskSetStable's constants.
	retryTaskSetStable         = 40
	retryTaskSetStableInterval = 15 * time.Second
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
		EnableExecuteCommand:          service.EnableExecuteCommand,
		HealthCheckGracePeriodSeconds: service.HealthCheckGracePeriodSeconds,
		PlacementConstraints:          service.PlacementConstraints,
		PlacementStrategy:             service.PlacementStrategy,
		PlatformVersion:               service.PlatformVersion,
		PropagateTags:                 types.PropagateTagsService,
		Role:                          service.RoleArn,
		SchedulingStrategy:            service.SchedulingStrategy,
		Tags:                          service.Tags,
	}
	output, err := c.ecsClient.CreateService(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create ECS service %s: %w", *service.ServiceName, err)
	}

	// Hack: Since we use EXTERNAL deployment controller, the below configurations are not allowed to be passed
	// in CreateService step, but it required in further step (CreateTaskSet step). We reassign those values
	// as part of service definition for that purpose.
	// ref: https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_CreateService.html
	output.Service.LaunchType = service.LaunchType
	output.Service.NetworkConfiguration = service.NetworkConfiguration
	output.Service.ServiceRegistries = service.ServiceRegistries

	return output.Service, nil
}

func (c *client) PruneServiceTasks(ctx context.Context, service types.Service) error {
	input := &ecs.UpdateServiceInput{
		Cluster:      service.ClusterArn,
		Service:      service.ServiceName,
		DesiredCount: aws.Int32(0),
	}

	_, err := c.ecsClient.UpdateService(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update ECS service %s: %w", *service.ServiceName, err)
	}

	return nil
}

func (c *client) UpdateService(ctx context.Context, service types.Service) (*types.Service, error) {
	input := &ecs.UpdateServiceInput{
		Cluster:              service.ClusterArn,
		Service:              service.ServiceName,
		EnableExecuteCommand: aws.Bool(service.EnableExecuteCommand),
		PlacementStrategy:    service.PlacementStrategy,
		// TODO: Support update other properties of service.
		// PlacementConstraints:    service.PlacementConstraints,
	}

	// If desiredCount is 0 or not set, keep current desiredCount because a user might use AutoScaling.
	if service.DesiredCount != 0 {
		input.DesiredCount = aws.Int32(service.DesiredCount)
	}

	output, err := c.ecsClient.UpdateService(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update ECS service %s: %w", *service.ServiceName, err)
	}

	// Hack: Since we use EXTERNAL deployment controller, the below configurations are not allowed to be passed
	// in UpdateService step, but it required in further step (CreateTaskSet step). We reassign those values
	// as part of service definition for that purpose.
	// ref: https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_CreateService.html
	output.Service.LaunchType = service.LaunchType
	output.Service.NetworkConfiguration = service.NetworkConfiguration
	output.Service.ServiceRegistries = service.ServiceRegistries

	return output.Service, nil
}

func (c *client) GetTaskDefinition(ctx context.Context, taskDefinitionArn string) (*types.TaskDefinition, error) {
	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(taskDefinitionArn),
	}
	output, err := c.ecsClient.DescribeTaskDefinition(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get ECS task definition %s: %w", taskDefinitionArn, err)
	}
	return output.TaskDefinition, nil
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
		RuntimePlatform:         taskDefinition.RuntimePlatform,
		EphemeralStorage:        taskDefinition.EphemeralStorage,
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

func (c *client) RunTask(ctx context.Context, taskDefinition types.TaskDefinition, clusterArn string, launchType string, awsVpcConfiguration *appconfig.ECSVpcConfiguration, tags []types.Tag) error {
	if taskDefinition.TaskDefinitionArn == nil {
		return fmt.Errorf("failed to run task of task family %s: no task definition provided", *taskDefinition.Family)
	}

	input := &ecs.RunTaskInput{
		TaskDefinition: taskDefinition.Family,
		Cluster:        aws.String(clusterArn),
		LaunchType:     types.LaunchType(launchType),
		Tags:           tags,
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
		Tags:           service.Tags,
		// If you specify the awsvpc network mode, the task is allocated an elastic network interface,
		// and you must specify a NetworkConfiguration when run a task with the task definition.
		NetworkConfiguration: service.NetworkConfiguration,
		LaunchType:           service.LaunchType,
		ServiceRegistries:    service.ServiceRegistries,
	}
	if targetGroup != nil {
		input.LoadBalancers = []types.LoadBalancer{*targetGroup}
	}
	output, err := c.ecsClient.CreateTaskSet(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create ECS task set %s: %w", *taskDefinition.TaskDefinitionArn, err)
	}

	// Wait created TaskSet to be stable.
	waitInput := &ecs.DescribeTaskSetsInput{
		Cluster:  service.ClusterArn,
		Service:  service.ServiceArn,
		TaskSets: []string{*output.TaskSet.TaskSetArn},
	}

	retry := backoff.NewRetry(retryTaskSetStable, backoff.NewConstant(retryTaskSetStableInterval))
	_, err = retry.Do(ctx, func() (interface{}, error) {
		output, err := c.ecsClient.DescribeTaskSets(ctx, waitInput)
		if err != nil {
			return nil, fmt.Errorf("failed to get ECS task set %s: %w", *taskDefinition.TaskDefinitionArn, err)
		}
		if len(output.TaskSets) == 0 {
			return nil, fmt.Errorf("failed to get ECS task set %s: task sets empty", *taskDefinition.TaskDefinitionArn)
		}
		taskSet := output.TaskSets[0]
		if taskSet.StabilityStatus == types.StabilityStatusSteadyState {
			return nil, nil
		}
		return nil, fmt.Errorf("task set %s is not stable", *taskDefinition.TaskDefinitionArn)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to wait ECS task set %s stable: %w", *taskDefinition.TaskDefinitionArn, err)
	}

	return output.TaskSet, nil
}

func (c *client) GetServiceTaskSets(ctx context.Context, service types.Service) ([]*types.TaskSet, error) {
	input := &ecs.DescribeServicesInput{
		Cluster: service.ClusterArn,
		Services: []string{
			*service.ServiceArn,
		},
	}
	output, err := c.ecsClient.DescribeServices(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get task sets of service %s: %w", *service.ServiceName, err)
	}
	if len(output.Services) == 0 {
		return nil, fmt.Errorf("failed to get task sets of service %s: services empty", *service.ServiceName)
	}
	svc := output.Services[0]
	activeTaskSetArns := make([]string, 0, len(svc.TaskSets))
	for i := range svc.TaskSets {
		if aws.ToString(svc.TaskSets[i].Status) == "DRAINING" {
			continue
		}
		activeTaskSetArns = append(activeTaskSetArns, *svc.TaskSets[i].TaskSetArn)
	}

	// No primary or active task set found.
	if len(activeTaskSetArns) == 0 {
		return []*types.TaskSet{}, nil
	}

	tsInput := &ecs.DescribeTaskSetsInput{
		Cluster:  service.ClusterArn,
		Service:  service.ServiceArn,
		TaskSets: activeTaskSetArns,
		Include: []types.TaskSetField{
			types.TaskSetFieldTags,
		},
	}
	tsOutput, err := c.ecsClient.DescribeTaskSets(ctx, tsInput)
	if err != nil {
		return nil, fmt.Errorf("failed to get task sets of service %s: %w", *service.ServiceName, err)
	}
	taskSets := make([]*types.TaskSet, 0, len(tsOutput.TaskSets))
	for i := range tsOutput.TaskSets {
		if !IsPipeCDManagedTaskSet(&tsOutput.TaskSets[i]) {
			continue
		}
		taskSets = append(taskSets, &tsOutput.TaskSets[i])
	}

	return taskSets, nil
}

// WaitServiceStable blocks until the ECS service is stable.
// It returns nil if the service is stable, otherwise it returns an error.
// Note: This function follow the implementation of the AWS CLI.
// AWS does not public API for waiting service stable, thus we use describe-service and workaround instead.
// ref: https://docs.aws.amazon.com/cli/latest/reference/ecs/wait/services-stable.html
func (c *client) WaitServiceStable(ctx context.Context, service types.Service) error {
	input := &ecs.DescribeServicesInput{
		Cluster:  service.ClusterArn,
		Services: []string{*service.ServiceArn},
	}

	retry := backoff.NewRetry(retryServiceStable, backoff.NewConstant(retryServiceStableInterval))
	_, err := retry.Do(ctx, func() (interface{}, error) {
		output, err := c.ecsClient.DescribeServices(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to get service %s: %w", *service.ServiceName, err)
		}

		if len(output.Services) == 0 {
			return nil, platformprovider.ErrNotFound
		}

		svc := output.Services[0]
		if svc.PendingCount == 0 && svc.RunningCount >= svc.DesiredCount {
			return nil, nil
		}

		return nil, fmt.Errorf("service %s is not stable", *service.ServiceName)
	})

	return err
}

func (c *client) DeleteTaskSet(ctx context.Context, taskSet types.TaskSet) error {
	input := &ecs.DeleteTaskSetInput{
		Cluster: taskSet.ClusterArn,
		Service: taskSet.ServiceArn,
		TaskSet: taskSet.TaskSetArn,
	}
	if _, err := c.ecsClient.DeleteTaskSet(ctx, input); err != nil {
		return fmt.Errorf("failed to delete ECS task set %s: %w", *taskSet.TaskSetArn, err)
	}

	// Inactive deleted taskset's task definition.
	taskDefInput := &ecs.DeregisterTaskDefinitionInput{
		TaskDefinition: taskSet.TaskDefinition,
	}
	if _, err := c.ecsClient.DeregisterTaskDefinition(ctx, taskDefInput); err != nil {
		return fmt.Errorf("failed to inactive ECS task definition %s: %w", *taskSet.TaskDefinition, err)
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

func (c *client) GetListenerArns(ctx context.Context, targetGroup types.LoadBalancer) ([]string, error) {
	loadBalancerArn, err := c.getLoadBalancerArn(ctx, *targetGroup.TargetGroupArn)
	if err != nil {
		return nil, err
	}

	input := &elasticloadbalancingv2.DescribeListenersInput{
		LoadBalancerArn: aws.String(loadBalancerArn),
	}
	output, err := c.elbClient.DescribeListeners(ctx, input)
	if err != nil {
		return nil, err
	}
	if len(output.Listeners) == 0 {
		return nil, platformprovider.ErrNotFound
	}

	arns := make([]string, len(output.Listeners))
	for i := range output.Listeners {
		arns[i] = *output.Listeners[i].ListenerArn
	}

	return arns, nil
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

func (c *client) ModifyListeners(ctx context.Context, listenerArns []string, routingTrafficCfg RoutingTrafficConfig) error {
	if len(routingTrafficCfg) != 2 {
		return fmt.Errorf("invalid listener configuration: requires 2 target groups")
	}

	for _, listenerArn := range listenerArns {
		describeRulesOutput, err := c.elbClient.DescribeRules(ctx, &elasticloadbalancingv2.DescribeRulesInput{
			ListenerArn: aws.String(listenerArn),
		})
		if err != nil {
			return fmt.Errorf("failed to describe rules of listener %s: %w", listenerArn, err)
		}

		for _, rule := range describeRulesOutput.Rules {
			modifiedActions := make([]elbtypes.Action, 0, len(rule.Actions))
			for _, action := range rule.Actions {
				if action.Type == elbtypes.ActionTypeEnumForward && routingTrafficCfg.hasSameTargets(action.ForwardConfig.TargetGroups) {
					// Modify only the forward action which has the same target groups.
					modifiedAction := elbtypes.Action{
						Type:  elbtypes.ActionTypeEnumForward,
						Order: action.Order,
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
					}
					modifiedActions = append(modifiedActions, modifiedAction)
				} else {
					modifiedActions = append(modifiedActions, action)
				}
			}

			// The default rule needs to be modified by ModifyListener API.
			if rule.IsDefault {
				_, err := c.elbClient.ModifyListener(ctx, &elasticloadbalancingv2.ModifyListenerInput{
					ListenerArn:    &listenerArn,
					DefaultActions: modifiedActions,
				})
				if err != nil {
					return fmt.Errorf("failed to modify default rule %s: %w", *rule.RuleArn, err)
				}
			} else {
				_, err := c.elbClient.ModifyRule(ctx, &elasticloadbalancingv2.ModifyRuleInput{
					RuleArn: rule.RuleArn,
					Actions: modifiedActions,
				})
				if err != nil {
					return fmt.Errorf("failed to modify rule %s: %w", *rule.RuleArn, err)
				}
			}
		}
	}
	return nil
}

func (c *client) TagResource(ctx context.Context, resourceArn string, tags []types.Tag) error {
	input := &ecs.TagResourceInput{
		ResourceArn: aws.String(resourceArn),
		Tags:        tags,
	}
	_, err := c.ecsClient.TagResource(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update tag of resource %s: %w", resourceArn, err)
	}
	return nil
}

func (c *client) ListClusters(ctx context.Context) ([]string, error) {
	in := &ecs.ListClustersInput{
		MaxResults: aws.Int32(100),
	}
	clusters := []string{}
	for {
		out, err := c.ecsClient.ListClusters(ctx, in)
		if err != nil {
			return nil, fmt.Errorf("failed to list ECS clusters: %w", err)
		}
		clusters = append(clusters, out.ClusterArns...)
		if out.NextToken == nil {
			return clusters, nil
		}
		in.NextToken = out.NextToken
	}
}

func (c *client) GetServices(ctx context.Context, clusterName string) ([]*types.Service, error) {
	listIn := &ecs.ListServicesInput{
		Cluster:    aws.String(clusterName),
		MaxResults: aws.Int32(100),
	}
	var serviceArns []string
	for {
		listOut, err := c.ecsClient.ListServices(ctx, listIn)
		if err != nil {
			return nil, fmt.Errorf("failed to list services of cluster %s: %w", clusterName, err)
		}
		serviceArns = append(serviceArns, listOut.ServiceArns...)
		if listOut.NextToken == nil {
			break
		}
		listIn.NextToken = listOut.NextToken
	}

	if len(serviceArns) == 0 {
		return []*types.Service{}, nil
	}

	services := make([]*types.Service, 0, len(serviceArns))
	// Split serviceArns into chunks of 10 to avoid the limitation in a single request of DescribeServices.
	for i := 0; i < len(serviceArns); i += 10 {
		end := i + 10
		if end > len(serviceArns) {
			end = len(serviceArns)
		}
		describeIn := &ecs.DescribeServicesInput{
			Cluster:  aws.String(clusterName),
			Services: serviceArns[i:end],
			Include:  []types.ServiceField{types.ServiceFieldTags},
		}
		describeOut, err := c.ecsClient.DescribeServices(ctx, describeIn)
		if err != nil {
			return nil, fmt.Errorf("failed to describe services: %w", err)
		}

		for i := range describeOut.Services {
			services = append(services, &describeOut.Services[i])
		}
	}

	return services, nil
}

func (c *client) GetTaskSetTasks(ctx context.Context, taskSet types.TaskSet) ([]*types.Task, error) {
	listIn := &ecs.ListTasksInput{
		Cluster: taskSet.ClusterArn,
		// Service tasks have the deployment ID, which is the same as taskSet's ID, as `startedBy` field.
		StartedBy: taskSet.Id,
	}
	listOut, err := c.ecsClient.ListTasks(ctx, listIn)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks of task set %s: %w", *taskSet.TaskSetArn, err)
	}

	taskArns := listOut.TaskArns
	tasks := make([]*types.Task, 0, len(taskArns))
	// Split taskArns into chunks of 100 to avoid the limitation in a single request of DescribeTasks.
	for i := 0; i < len(taskArns); i += 100 {
		end := i + 100
		if end > len(taskArns) {
			end = len(taskArns)
		}

		describeIn := &ecs.DescribeTasksInput{
			Cluster: taskSet.ClusterArn,
			Tasks:   listOut.TaskArns[i:end],
		}
		out, err := c.ecsClient.DescribeTasks(ctx, describeIn)
		if err != nil {
			return nil, fmt.Errorf("failed to describe tasks: %w", err)
		}

		for i := range out.Tasks {
			tasks = append(tasks, &out.Tasks[i])
		}
	}

	return tasks, nil
}
