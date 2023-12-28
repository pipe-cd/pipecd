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
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"

	"github.com/pipe-cd/pipecd/pkg/config"
)

var ErrNoTargetGroup = errors.New("no target group")

func loadTargetGroups(targetGroups config.ECSTargetGroups) (*types.LoadBalancer, *types.LoadBalancer, error) {
	if targetGroups.Primary == nil {
		return nil, nil, ErrNoTargetGroup
	}

	primary := &types.LoadBalancer{
		TargetGroupArn:   aws.String(targetGroups.Primary.TargetGroupArn),
		ContainerName:    aws.String(targetGroups.Primary.ContainerName),
		ContainerPort:    aws.Int32(int32(targetGroups.Primary.ContainerPort)),
		LoadBalancerName: aws.String(targetGroups.Primary.LoadBalancerName),
	}

	var canary *types.LoadBalancer
	if targetGroups.Canary != nil {
		canary = &types.LoadBalancer{
			TargetGroupArn:   aws.String(targetGroups.Canary.TargetGroupArn),
			ContainerName:    aws.String(targetGroups.Canary.ContainerName),
			ContainerPort:    aws.Int32(int32(targetGroups.Canary.ContainerPort)),
			LoadBalancerName: aws.String(targetGroups.Canary.LoadBalancerName),
		}
	}

	return primary, canary, nil
}
