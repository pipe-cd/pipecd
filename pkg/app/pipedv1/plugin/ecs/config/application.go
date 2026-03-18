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

package config

import (
	"encoding/json"

	"github.com/creasty/defaults"
)

// ECSApplicationSpec defines the application specification for ECS plugin.
type ECSApplicationSpec struct {
	Input            ECSDeploymentInput  `json:"input"`
	QuickSyncOptions ECSSyncStageOptions `json:"quickSync"`
}

// ECSDeploymentInput defines the input for ECS deployment.
type ECSDeploymentInput struct {
	// TaskDefinitionFile is the name of task definition file placing in application directory
	// e.g., "taskdef.json" or "ecs/taskdef.yaml"
	// Default: taskdef.json
	TaskDefinitionFile string `json:"taskDefinitionFile,omitempty" default:"taskdef.json"`

	// ServiceDefinitionFile is the name of service definition file placing in application directory
	// e.g., "servicedef.json" or "ecs/servicedef.yaml"
	ServiceDefinitionFile string `json:"serviceDefinitionFile,omitempty"`

	// RunStandaloneTask indicates whether to run the task as a standalone task without creating/updating an ECS service.
	// If true, the plugin will run the task directly without managing it through an ECS service.
	// This is useful for running one-off tasks or jobs that do not require long-term management.
	// Default: false
	RunStandaloneTask bool `json:"runStandaloneTask,omitempty" default:"false"`

	// ClusterARN identifies the ECS cluster where the task and service will be deployed.
	ClusterARN string `json:"clusterArn,omitempty"`

	// LaunchType specifies the launch type on which to run your task.
	// Valid values: "EC2", "FARGATE"
	// Default: "FARGATE"
	LaunchType string `json:"launchType,omitempty" default:"FARGATE"`

	// AccessType specifies how the ECS service is accessed.
	// Valid values: "ELB", "SERVICE_DISCOVERY"
	// Default: "ELB"
	AccessType string `json:"accessType,omitempty" default:"ELB"`

	// AwsVpcConfiguration contains the VPC configuration for running ECS tasks.
	AwsVpcConfiguration ECSVpcConfiguration `json:"awsvpcConfiguration"`

	// TargetGroups contains the load balancer target groups for the ECS service.
	TargetGroups ECSTargetGroups `json:"targetGroups"`
}

func (di *ECSDeploymentInput) UnmarshalJSON(data []byte) error {
	type alias ECSDeploymentInput
	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	*di = ECSDeploymentInput(a)
	return defaults.Set(di)
}

// ECSVpcConfiguration contains the VPC configuration for running ECS tasks.
type ECSVpcConfiguration struct {
	// Subnets is a list of VPC subnet IDs where tasks will be launched.
	// Limit: 16 subnets per VPC configuration (https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/ecs@v1.46.2/types#AwsVpcConfiguration).
	// This field is required.
	Subnets []string `json:"subnets"`

	// AssignPublicIP indicates whether to assign a public IP address to the task's ENI
	// Valid values: "ENABLED","DISABLED"
	AssignPublicIP string `json:"assignPublicIp,omitempty"`

	// SecurityGroups is a list of security group IDs associated with the task's ENI
	// If not specified, the default security group for the VPC will be used.
	// Limit: 5 security groups per VPC configuration (https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/ecs@v1.46.2/types#AwsVpcConfiguration).
	// All security groups must be from the same VPC.
	SecurityGroups []string `json:"securityGroups,omitempty"`
}

// ECSTargetGroups represents the load balancer target groups.
type ECSTargetGroups struct {
	// Primary is the target group for the primary service.
	Primary *ECSTargetGroup `json:"primary,omitempty"`

	// Canary is the target group for the canary service (optional).
	Canary *ECSTargetGroup `json:"canary,omitempty"`
}

// ECSTargetGroup represents a single load balancer target group
type ECSTargetGroup struct {
	// TargetGroupARN is the ARN of the target group
	TargetGroupARN string `json:"targetGroupArn,omitempty"`

	// ContainerName is the name of the container to associate with the target group
	ContainerName string `json:"containerName,omitempty"`

	// ContainerPort is the port on the container to associate with the target group
	ContainerPort int32 `json:"containerPort,omitempty"`
}
