// Copyright 2023 The PipeCD Authors.
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
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"
)

func TestParseServiceDefinition(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		input       string
		expected    types.Service
		expectedErr bool
	}{
		{
			name: "yaml format input",
			input: `
cluster: arn:aws:ecs:ap-northeast-1:XXXX:cluster/YYYY
serviceName: nginx-external-canary
desiredCount: 2
role: arn:aws:iam::XXXXX:role/ecsTaskExecutionRole
deploymentConfiguration:
  maximumPercent: 200
  minimumHealthyPercent: 0
schedulingStrategy: REPLICA
deploymentController:
  type: EXTERNAL
`,
			expected: types.Service{
				ClusterArn:   aws.String("arn:aws:ecs:ap-northeast-1:XXXX:cluster/YYYY"),
				ServiceName:  aws.String("nginx-external-canary"),
				DesiredCount: 2,
				RoleArn:      aws.String("arn:aws:iam::XXXXX:role/ecsTaskExecutionRole"),
				DeploymentConfiguration: &types.DeploymentConfiguration{
					MaximumPercent:        aws.Int32(200),
					MinimumHealthyPercent: aws.Int32(0),
				},
				SchedulingStrategy: types.SchedulingStrategyReplica,
				DeploymentController: &types.DeploymentController{
					Type: types.DeploymentControllerTypeExternal,
				},
			},
		},
		{
			name: "yaml format input with roleArn field name",
			input: `
cluster: arn:aws:ecs:ap-northeast-1:XXXX:cluster/YYYY
serviceName: nginx-external-canary
desiredCount: 2
roleArn: arn:aws:iam::XXXXX:role/ecsTaskExecutionRole
deploymentConfiguration:
  maximumPercent: 200
  minimumHealthyPercent: 0
schedulingStrategy: REPLICA
deploymentController:
  type: EXTERNAL
`,
			expected: types.Service{
				ClusterArn:   aws.String("arn:aws:ecs:ap-northeast-1:XXXX:cluster/YYYY"),
				ServiceName:  aws.String("nginx-external-canary"),
				DesiredCount: 2,
				RoleArn:      aws.String("arn:aws:iam::XXXXX:role/ecsTaskExecutionRole"),
				DeploymentConfiguration: &types.DeploymentConfiguration{
					MaximumPercent:        aws.Int32(200),
					MinimumHealthyPercent: aws.Int32(0),
				},
				SchedulingStrategy: types.SchedulingStrategyReplica,
				DeploymentController: &types.DeploymentController{
					Type: types.DeploymentControllerTypeExternal,
				},
			},
		},
		{
			name: "json format input",
			input: `
{
  "cluster": "arn:aws:ecs:ap-northeast-1:XXXX:cluster/YYYY",
  "serviceName": "nginx-external-canary",
  "desiredCount": 2,
  "role": "arn:aws:iam::XXXXX:role/ecsTaskExecutionRole",
  "deploymentConfiguration": {
    "maximumPercent": 200,
    "minimumHealthyPercent": 0
  },
  "schedulingStrategy": "REPLICA",
  "deploymentController": {
    "type": "EXTERNAL"
  }
}
`,
			expected: types.Service{
				ClusterArn:   aws.String("arn:aws:ecs:ap-northeast-1:XXXX:cluster/YYYY"),
				ServiceName:  aws.String("nginx-external-canary"),
				DesiredCount: 2,
				RoleArn:      aws.String("arn:aws:iam::XXXXX:role/ecsTaskExecutionRole"),
				DeploymentConfiguration: &types.DeploymentConfiguration{
					MaximumPercent:        aws.Int32(200),
					MinimumHealthyPercent: aws.Int32(0),
				},
				SchedulingStrategy: types.SchedulingStrategyReplica,
				DeploymentController: &types.DeploymentController{
					Type: types.DeploymentControllerTypeExternal,
				},
			},
		},
		{
			name: "json format input with clusterArn field name",
			input: `
{
  "clusterArn": "arn:aws:ecs:ap-northeast-1:XXXX:cluster/YYYY",
  "serviceName": "nginx-external-canary",
  "desiredCount": 2,
  "role": "arn:aws:iam::XXXXX:role/ecsTaskExecutionRole",
  "deploymentConfiguration": {
    "maximumPercent": 200,
    "minimumHealthyPercent": 0
  },
  "schedulingStrategy": "REPLICA",
  "deploymentController": {
    "type": "EXTERNAL"
  }
}
`,
			expected: types.Service{
				ClusterArn:   aws.String("arn:aws:ecs:ap-northeast-1:XXXX:cluster/YYYY"),
				ServiceName:  aws.String("nginx-external-canary"),
				DesiredCount: 2,
				RoleArn:      aws.String("arn:aws:iam::XXXXX:role/ecsTaskExecutionRole"),
				DeploymentConfiguration: &types.DeploymentConfiguration{
					MaximumPercent:        aws.Int32(200),
					MinimumHealthyPercent: aws.Int32(0),
				},
				SchedulingStrategy: types.SchedulingStrategyReplica,
				DeploymentController: &types.DeploymentController{
					Type: types.DeploymentControllerTypeExternal,
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseServiceDefinition([]byte(tc.input))
			assert.Equal(t, tc.expectedErr, err != nil)
			assert.Equal(t, tc.expected, got)
		})
	}
}
