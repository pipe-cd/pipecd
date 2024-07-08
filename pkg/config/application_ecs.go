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

package config

import (
	"fmt"
)

const (
	AccessTypeELB              string = "ELB"
	AccessTypeServiceDiscovery string = "SERVICE_DISCOVERY"
)

// ECSApplicationSpec represents an application configuration for ECS application.
type ECSApplicationSpec struct {
	GenericApplicationSpec
	// Input for ECS deployment such as where to fetch source code...
	Input ECSDeploymentInput `json:"input"`
	// Configuration for quick sync.
	QuickSync ECSSyncStageOptions `json:"quickSync"`
}

// Validate returns an error if any wrong configuration value was found.
func (s *ECSApplicationSpec) Validate() error {
	if err := s.GenericApplicationSpec.Validate(); err != nil {
		return err
	}

	if err := s.Input.validate(); err != nil {
		return err
	}

	return nil
}

type ECSDeploymentInput struct {
	// The Amazon Resource Name (ARN) that identifies the cluster.
	ClusterArn string `json:"clusterArn,omitempty"`
	// The launch type on which to run your task.
	// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/launch_types.html
	// Default is FARGATE
	LaunchType string `json:"launchType,omitempty" default:"FARGATE"`
	// VpcConfiguration ECSVpcConfiguration `json:"awsvpcConfiguration"`
	AwsVpcConfiguration ECSVpcConfiguration `json:"awsvpcConfiguration,omitempty" default:""`
	// The name of service definition file placing in application directory.
	ServiceDefinitionFile string `json:"serviceDefinitionFile"`
	// The name of task definition file placing in application directory.
	// Default is taskdef.json
	TaskDefinitionFile string `json:"taskDefinitionFile" default:"taskdef.json"`
	// ECSTargetGroups
	TargetGroups ECSTargetGroups `json:"targetGroups,omitempty"`
	// Automatically reverts all changes from all stages when one of them failed.
	// Default is true.
	AutoRollback *bool `json:"autoRollback,omitempty" default:"true"`
	// Run standalone task during deployment.
	// Default is true.
	RunStandaloneTask *bool `json:"runStandaloneTask,omitempty" default:"true"`
	// How the ECS service is accessed.
	// Possible values are:
	//  - ELB -  The service is accessed via ELB and target groups.
	//  - SERVICE_DISCOVERY -  The service is accessed via ECS Service Discovery.
	// Default is ELB.
	AccessType string `json:"accessType,omitempty" default:"ELB"`
	// Whether to ignore the desiredCount of the service definition when updating the service.
	// If this is set as true, the desiredCount will be updated when creating the service, but will not be updated when updating the service.
	// Default is false.
	IgnoreDesiredCountOnUpdate *bool `json:"ignoreDesiredCountOnUpdate,omitempty" default:"false"`
}

func (in *ECSDeploymentInput) IsStandaloneTask() bool {
	return in.ServiceDefinitionFile == ""
}

func (in *ECSDeploymentInput) IsAccessedViaELB() bool {
	return in.AccessType == AccessTypeELB
}

type ECSVpcConfiguration struct {
	Subnets        []string `json:"subnets,omitempty"`
	AssignPublicIP string   `json:"assignPublicIp,omitempty"`
	SecurityGroups []string `json:"securityGroups,omitempty"`
}

type ECSTargetGroups struct {
	Primary *ECSTargetGroup `json:"primary,omitempty"`
	Canary  *ECSTargetGroup `json:"canary,omitempty"`
}

type ECSTargetGroup struct {
	TargetGroupArn   string `json:"targetGroupArn,omitempty"`
	ContainerName    string `json:"containerName,omitempty"`
	ContainerPort    int    `json:"containerPort,omitempty"`
	LoadBalancerName string `json:"loadBalancerName,omitempty"`
}

// ECSSyncStageOptions contains all configurable values for a ECS_SYNC stage.
type ECSSyncStageOptions struct {
	// Whether to delete old tasksets before creating new ones or not.
	// If this is set, the application may be unavailable for a short of time during the deployment.
	// Default is false.
	Recreate bool `json:"recreate"`
}

// ECSCanaryRolloutStageOptions contains all configurable values for a ECS_CANARY_ROLLOUT stage.
type ECSCanaryRolloutStageOptions struct {
	// Scale represents the amount of desired task that should be rolled out as CANARY variant workload.
	Scale Percentage `json:"scale"`
}

// ECSPrimaryRolloutStageOptions contains all configurable values for a ECS_PRIMARY_ROLLOUT stage.
type ECSPrimaryRolloutStageOptions struct {
}

// ECSCanaryCleanStageOptions contains all configurable values for a ECS_CANARY_CLEAN stage.
type ECSCanaryCleanStageOptions struct {
}

// ECSTrafficRoutingStageOptions contains all configurable values for ECS_TRAFFIC_ROUTING stage.
type ECSTrafficRoutingStageOptions struct {
	// Canary represents the amount of traffic that the rolled out CANARY variant will serve.
	Canary Percentage `json:"canary,omitempty"`
	// Primary represents the amount of traffic that the rolled out CANARY variant will serve.
	Primary Percentage `json:"primary,omitempty"`
}

func (opts ECSTrafficRoutingStageOptions) Percentage() (primary, canary int) {
	primary = opts.Primary.Int()
	if primary > 0 && primary <= 100 {
		canary = 100 - primary
		return
	}

	canary = opts.Canary.Int()
	if canary > 0 && canary <= 100 {
		primary = 100 - canary
		return
	}
	// As default, Primary variant will receive 100% of traffic.
	primary = 100
	canary = 0
	return
}

func (in *ECSDeploymentInput) validate() error {
	switch in.AccessType {
	case AccessTypeELB, AccessTypeServiceDiscovery:
		break
	default:
		return fmt.Errorf("invalid accessType: %s", in.AccessType)
	}
	return nil
}
