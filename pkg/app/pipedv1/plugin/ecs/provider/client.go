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

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"

	"github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider"
	"github.com/pipe-cd/pipecd/pkg/backoff"

	appconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/config"
)

const (
	// WaitServiceStable's constants.
	retryServiceStable         = 40
	retryServiceStableInterval = 15 * time.Second

	// WaitTaskSetStable's constants.
	maxTaskSetStableRetries    = 5
	retryTaskSetStableInterval = 30 * time.Second
)

type client struct {
	ecsClient *ecs.Client
}

func newClient(region, profile, credentialsFile, roleARN, tokenPath string) (Client, error) {
	if region == "" {
		return nil, fmt.Errorf("region is required field")
	}

	c := &client{}

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
		PropagateTags:                 service.PropagateTags,
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

func (c *client) UpdateService(ctx context.Context, service types.Service) (*types.Service, error) {
	// TODO: Support other properties (current only support the properties that v0 supports)
	// This should be delegated to user to decide which properties to update by defining in service definition file.
	input := &ecs.UpdateServiceInput{
		Cluster:              service.ClusterArn,
		Service:              service.ServiceName,
		EnableExecuteCommand: aws.Bool(service.EnableExecuteCommand),
		PlacementStrategy:    service.PlacementStrategy,
		PropagateTags:        service.PropagateTags,
		EnableECSManagedTags: aws.Bool(service.EnableECSManagedTags),
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

func (c *client) DescribeService(ctx context.Context, service types.Service) (*types.Service, error) {
	input := &ecs.DescribeServicesInput{
		Cluster: service.ClusterArn,
		Services: []string{
			*service.ServiceName,
		},
	}
	output, err := c.ecsClient.DescribeServices(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get service %s description: %w", *service.ServiceName, err)
	}

	if len(output.Services) == 0 {
		return nil, fmt.Errorf("services %s does not exist", *service.ServiceName)
	}

	return &output.Services[0], nil
}

func (c *client) GetServiceTaskSets(ctx context.Context, service types.Service) ([]types.TaskSet, error) {
	input := &ecs.DescribeServicesInput{
		Cluster: service.ClusterArn,
		Services: []string{
			*service.ServiceArn,
		},
	}
	output, err := c.ecsClient.DescribeServices(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get service %s description: %w", *service.ServiceName, err)
	}
	if len(output.Services) == 0 {
		return nil, fmt.Errorf("services %s does not exist", *service.ServiceName)
	}
	svc := output.Services[0]
	nonDrainTaskSetArns := make([]string, 0, len(svc.TaskSets))
	for i := range svc.TaskSets {
		if aws.ToString(svc.TaskSets[i].Status) == "DRAINING" {
			continue
		}
		nonDrainTaskSetArns = append(nonDrainTaskSetArns, *svc.TaskSets[i].TaskSetArn)
	}

	// There is no primary or active task set, return immediately
	if len(nonDrainTaskSetArns) == 0 {
		return []types.TaskSet{}, nil
	}

	// AWS does not return full TaskSet information in DescribeServices API
	// Need to call DescribeTaskSets API to get full information of TaskSet.

	// Need tags information to find out which TaskSet is managed by PipeCD ECS plugin.
	tsInput := &ecs.DescribeTaskSetsInput{
		Cluster:  service.ClusterArn,
		Service:  service.ServiceArn,
		TaskSets: nonDrainTaskSetArns,
		Include: []types.TaskSetField{
			types.TaskSetFieldTags,
		},
	}
	tsOutput, err := c.ecsClient.DescribeTaskSets(ctx, tsInput)
	if err != nil {
		return nil, fmt.Errorf("failed to get task sets of service %s: %w", *service.ServiceName, err)
	}
	taskSets := make([]types.TaskSet, 0, len(tsOutput.TaskSets))
	for i := range tsOutput.TaskSets {
		if !IsPipeCDManagedTaskSet(&tsOutput.TaskSets[i]) {
			continue
		}
		taskSets = append(taskSets, tsOutput.TaskSets[i])
	}

	return taskSets, nil
}

func (c *client) GetPrimaryTaskSet(ctx context.Context, service types.Service) (*types.TaskSet, error) {
	input := &ecs.DescribeServicesInput{
		Cluster:  service.ClusterArn,
		Services: []string{*service.ServiceArn},
	}

	output, err := c.ecsClient.DescribeServices(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get primary taskset of service %s: %w", *service.ServiceName, err)
	}
	if len(output.Services) == 0 {
		return nil, fmt.Errorf("failed to get primary task set of service %s: service not found", *service.ServiceName)
	}

	for _, ts := range output.Services[0].TaskSets {
		if aws.ToString(ts.Status) == "PRIMARY" {
			return &ts, nil
		}
	}

	// A newly created service may have no PRIMARY task set yet,
	return nil, nil
}

func (c *client) CreateTaskSet(ctx context.Context, service types.Service, taskDefinition types.TaskDefinition, targetGroup *types.LoadBalancer, scale float64) (*types.TaskSet, error) {
	if taskDefinition.TaskDefinitionArn == nil {
		return nil, fmt.Errorf("no task definition provided for family %s", *taskDefinition.Family)
	}

	input := &ecs.CreateTaskSetInput{
		Cluster:              service.ClusterArn,
		Service:              service.ServiceArn,
		TaskDefinition:       taskDefinition.TaskDefinitionArn,
		Scale:                &types.Scale{Unit: types.ScaleUnitPercent, Value: scale},
		Tags:                 service.Tags,
		NetworkConfiguration: service.NetworkConfiguration,
		LaunchType:           service.LaunchType,
		ServiceRegistries:    service.ServiceRegistries,
	}

	if targetGroup != nil {
		input.LoadBalancers = append(input.LoadBalancers, *targetGroup)
	}
	output, err := c.ecsClient.CreateTaskSet(ctx, input)
	if err != nil {
		return nil, err
	}

	// Wait created TaskSet to be stable.
	waitInput := &ecs.DescribeTaskSetsInput{
		Cluster:  service.ClusterArn,
		Service:  service.ServiceArn,
		TaskSets: []string{*output.TaskSet.TaskSetArn},
	}

	retry := backoff.NewRetry(maxTaskSetStableRetries, backoff.NewConstant(retryTaskSetStableInterval))
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

func (c *client) GetTasks(ctx context.Context, service types.Service) ([]types.Task, error) {
	// Get list of task ARN of the given service, using pagination here because max number of tasks return from ListTasks API is 100
	var taskArns []string
	paginator := ecs.NewListTasksPaginator(c.ecsClient, &ecs.ListTasksInput{
		Cluster:     service.ClusterArn,
		ServiceName: service.ServiceName,
	})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list tasks of service %s: %w", *service.ServiceName, err)
		}
		taskArns = append(taskArns, page.TaskArns...)
	}

	if len(taskArns) == 0 {
		return nil, nil
	}

	var tasks []types.Task
	// Max number of tasks in each run of DescribeTasks is 100
	const batchSize = 100
	for i := 0; i < len(taskArns); i += batchSize {
		end := i + batchSize
		if end > len(taskArns) {
			end = len(taskArns)
		}

		batch := taskArns[i:end]
		out, err := c.ecsClient.DescribeTasks(ctx, &ecs.DescribeTasksInput{
			Cluster: service.ClusterArn,
			Tasks:   batch,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to describe tasks: %w", err)
		}

		tasks = append(tasks, out.Tasks...)
	}
	return tasks, nil
}

func (c *client) ServiceExists(ctx context.Context, cluster, serviceName string) (bool, error) {
	input := &ecs.DescribeServicesInput{
		Cluster: aws.String(cluster), // cluster field can be ARN or name
		Services: []string{
			serviceName,
		},
	}
	output, err := c.ecsClient.DescribeServices(ctx, input)
	if err != nil {
		return false, err
	}
	if len(output.Services) == 0 {
		return false, nil
	}
	return true, nil
}

// WaitServiceStable blocks until the ECS service is stable.
// It returns nil if the service is stable, otherwise it returns an error.
// Note: This function follow the implementation of the AWS CLI.
// AWS does not public API for waiting service stable, thus we use describe-service and workaround instead.
// ref: https://docs.aws.amazon.com/cli/latest/reference/ecs/wait/services-stable.html
func (c *client) WaitServiceStable(ctx context.Context, clusterArn, service string) error {
	input := &ecs.DescribeServicesInput{
		Cluster:  aws.String(clusterArn),
		Services: []string{service},
	}

	retry := backoff.NewRetry(retryServiceStable, backoff.NewConstant(retryServiceStableInterval))
	_, err := retry.Do(ctx, func() (interface{}, error) {
		output, err := c.ecsClient.DescribeServices(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to get service %s: %w", service, err)
		}

		if len(output.Services) == 0 {
			return nil, platformprovider.ErrNotFound
		}

		svc := output.Services[0]
		if svc.PendingCount == 0 && svc.RunningCount >= svc.DesiredCount {
			return nil, nil
		}

		return nil, fmt.Errorf("service %s is not stable", service)
	})

	return err
}

func (c *client) GetServiceStatus(ctx context.Context, cluster, serviceName string) (string, error) {
	input := &ecs.DescribeServicesInput{
		Cluster:  aws.String(cluster),
		Services: []string{serviceName},
	}
	output, err := c.ecsClient.DescribeServices(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to get service %s description: %w", serviceName, err)
	}
	if len(output.Services) == 0 {
		return "", fmt.Errorf("services %s does not exist", serviceName)
	}
	return *output.Services[0].Status, nil
}

func (c *client) RegisterTaskDefinition(ctx context.Context, taskDef types.TaskDefinition) (*types.TaskDefinition, error) {
	input := &ecs.RegisterTaskDefinitionInput{
		Family:                  taskDef.Family,
		ContainerDefinitions:    taskDef.ContainerDefinitions,
		RequiresCompatibilities: taskDef.RequiresCompatibilities,
		ExecutionRoleArn:        taskDef.ExecutionRoleArn,
		TaskRoleArn:             taskDef.TaskRoleArn,
		NetworkMode:             taskDef.NetworkMode,
		Volumes:                 taskDef.Volumes,
		RuntimePlatform:         taskDef.RuntimePlatform,
		EphemeralStorage:        taskDef.EphemeralStorage,
		Cpu:                     taskDef.Cpu,
		Memory:                  taskDef.Memory,
		InferenceAccelerators:   taskDef.InferenceAccelerators,
		IpcMode:                 taskDef.IpcMode,
		PidMode:                 taskDef.PidMode,
		PlacementConstraints:    taskDef.PlacementConstraints,
		ProxyConfiguration:      taskDef.ProxyConfiguration,
	}
	output, err := c.ecsClient.RegisterTaskDefinition(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to register ECS task definition %s: %w", *taskDef.Family, err)
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

// PruneServiceTasks sets desired count of the service to 0 to stop all running tasks of the service.
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

// ListTags returns the list of tags for the given resource ARN.
func (c *client) ListTags(ctx context.Context, resourceArn string) ([]types.Tag, error) {
	input := &ecs.ListTagsForResourceInput{
		ResourceArn: aws.String(resourceArn),
	}
	output, err := c.ecsClient.ListTagsForResource(ctx, input)
	if err != nil {
		return nil, err
	}
	// If there is no tag, AWS returns nil
	// Return an empty slice instead
	if output.Tags == nil {
		return []types.Tag{}, nil
	}
	return output.Tags, nil
}

func (c *client) TagResource(ctx context.Context, resourceArn string, tags []types.Tag) error {
	input := &ecs.TagResourceInput{
		ResourceArn: aws.String(resourceArn),
		Tags:        tags,
	}
	_, err := c.ecsClient.TagResource(ctx, input)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) UntagResource(ctx context.Context, resourceArn string, tagKeys []string) error {
	input := &ecs.UntagResourceInput{
		ResourceArn: aws.String(resourceArn),
		TagKeys:     tagKeys,
	}
	_, err := c.ecsClient.UntagResource(ctx, input)
	if err != nil {
		return err
	}
	return nil
}
