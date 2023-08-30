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

	"github.com/pipe-cd/pipecd/pkg/config"
)

func TestLoadTargetGroup(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		cfg         config.ECSTargetGroups
		expected    []*types.LoadBalancer
		expectedErr bool
	}{
		{
			name:        "no target group",
			cfg:         config.ECSTargetGroups{},
			expected:    []*types.LoadBalancer{nil, nil},
			expectedErr: true,
		},
		{
			name: "primary target group only",
			cfg: config.ECSTargetGroups{
				Primary: []byte(`{"targetGroupArn": "primary-target-group-arn", "containerName": "primary-container-name", "containerPort": 80}`),
			},
			expected: []*types.LoadBalancer{
				{
					TargetGroupArn: aws.String("primary-target-group-arn"),
					ContainerName:  aws.String("primary-container-name"),
					ContainerPort:  aws.Int32(80),
				},
				nil,
			},
			expectedErr: false,
		},
		{
			name: "primary and canary target group",
			cfg: config.ECSTargetGroups{
				Primary: []byte(`{"targetGroupArn": "primary-target-group-arn", "containerName": "primary-container-name", "containerPort": 80}`),
				Canary:  []byte(`{"targetGroupArn": "canary-target-group-arn", "containerName": "canary-container-name", "containerPort": 80}`),
			},
			expected: []*types.LoadBalancer{
				{
					TargetGroupArn: aws.String("primary-target-group-arn"),
					ContainerName:  aws.String("primary-container-name"),
					ContainerPort:  aws.Int32(80),
				},
				{
					TargetGroupArn: aws.String("canary-target-group-arn"),
					ContainerName:  aws.String("canary-container-name"),
					ContainerPort:  aws.Int32(80),
				},
			},
			expectedErr: false,
		},
		{
			name: "invalid primary target group",
			cfg: config.ECSTargetGroups{
				Primary: []byte(`{"invalidField": "primary-target-group-arn"}`),
			},
			expected:    []*types.LoadBalancer{nil, nil},
			expectedErr: true,
		},
		{
			name: "invalid canary target group",
			cfg: config.ECSTargetGroups{
				Primary: []byte(`{"targetGroupArn": "primary-target-group-arn"`),
				Canary:  []byte(`{"invalidField": "canary-target-group-arn"}`),
			},
			expected:    []*types.LoadBalancer{nil, nil},
			expectedErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			primary, canary, err := loadTargetGroups(tc.cfg)
			assert.Equal(t, tc.expectedErr, err != nil)
			assert.Equal(t, tc.expected[0], primary)
			assert.Equal(t, tc.expected[1], canary)
		})
	}
}
