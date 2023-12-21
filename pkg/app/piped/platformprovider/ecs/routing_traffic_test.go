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
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/stretchr/testify/assert"
)

func TestHasSameTargets(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		cfg           RoutingTrafficConfig
		actionTargets []types.TargetGroupTuple
		expected      bool
	}{
		{
			name: "has the same target groups in the same order",
			cfg: RoutingTrafficConfig{
				{
					TargetGroupArn: "arn:aws:elasticloadbalancing:<region>:<account-id>:targetgroup/xxx/yyy1",
					Weight:         100,
				},
				{
					TargetGroupArn: "arn:aws:elasticloadbalancing:<region>:<account-id>:targetgroup/xxx/yyy2",
					Weight:         0,
				},
			},
			actionTargets: []types.TargetGroupTuple{
				{
					TargetGroupArn: aws.String("arn:aws:elasticloadbalancing:<region>:<account-id>:targetgroup/xxx/yyy1"),
					Weight:         aws.Int32(100),
				},
				{
					TargetGroupArn: aws.String("arn:aws:elasticloadbalancing:<region>:<account-id>:targetgroup/xxx/yyy2"),
					Weight:         aws.Int32(0),
				},
			},
			expected: true,
		},
		{
			name: "has the same target groups in the different order",
			cfg: RoutingTrafficConfig{
				{
					TargetGroupArn: "arn:aws:elasticloadbalancing:<region>:<account-id>:targetgroup/xxx/yyy1",
					Weight:         100,
				},
				{
					TargetGroupArn: "arn:aws:elasticloadbalancing:<region>:<account-id>:targetgroup/xxx/yyy2",
					Weight:         0,
				},
			},
			actionTargets: []types.TargetGroupTuple{
				{
					TargetGroupArn: aws.String("arn:aws:elasticloadbalancing:<region>:<account-id>:targetgroup/xxx/yyy2"),
					Weight:         aws.Int32(0),
				},
				{
					TargetGroupArn: aws.String("arn:aws:elasticloadbalancing:<region>:<account-id>:targetgroup/xxx/yyy1"),
					Weight:         aws.Int32(100),
				},
			},
			expected: true,
		},
		{
			name: "the number of target groups are different",
			cfg: RoutingTrafficConfig{
				{
					TargetGroupArn: "arn:aws:elasticloadbalancing:<region>:<account-id>:targetgroup/xxx/yyy1",
					Weight:         100,
				},
			},
			actionTargets: []types.TargetGroupTuple{
				{
					TargetGroupArn: aws.String("arn:aws:elasticloadbalancing:<region>:<account-id>:targetgroup/xxx/yyy1"),
					Weight:         aws.Int32(0),
				},
				{
					TargetGroupArn: aws.String("arn:aws:elasticloadbalancing:<region>:<account-id>:targetgroup/xxx/yyy2"),
					Weight:         aws.Int32(100),
				},
			},
			expected: false,
		},
		{
			name: "has a different target group",
			cfg: RoutingTrafficConfig{
				{
					TargetGroupArn: "arn:aws:elasticloadbalancing:<region>:<account-id>:targetgroup/xxx/yyy1",
					Weight:         100,
				},
				{
					TargetGroupArn: "arn:aws:elasticloadbalancing:<region>:<account-id>:targetgroup/xxx/yyy2",
					Weight:         0,
				},
			},
			actionTargets: []types.TargetGroupTuple{
				{
					TargetGroupArn: aws.String("arn:aws:elasticloadbalancing:<region>:<account-id>:targetgroup/xxx/yyy1"),
					Weight:         aws.Int32(0),
				},
				{
					TargetGroupArn: aws.String("arn:aws:elasticloadbalancing:<region>:<account-id>:targetgroup/xxx/yyy3"),
					Weight:         aws.Int32(100),
				},
			},
			expected: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			hasSame := tc.cfg.hasSameTargets(tc.actionTargets)
			assert.Equal(t, tc.expected, hasSame)
		})
	}
}
