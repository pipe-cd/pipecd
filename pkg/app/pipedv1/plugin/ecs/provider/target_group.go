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
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/config"
)

var ErrNoTargetGroup = errors.New("no target group")

// loadTargetGroups returns the primary and canary target groups from the given config.
func loadTargetGroups(targetGroups config.ECSTargetGroups) (*types.LoadBalancer, *types.LoadBalancer, error) {
	if targetGroups.Primary == nil {
		return nil, nil, ErrNoTargetGroup
	}

	primary := &types.LoadBalancer{
		TargetGroupArn: aws.String(targetGroups.Primary.TargetGroupARN),
		ContainerName:  aws.String(targetGroups.Primary.ContainerName),
		ContainerPort:  aws.Int32(targetGroups.Primary.ContainerPort),
	}

	var canary *types.LoadBalancer
	if targetGroups.Canary != nil {
		canary = &types.LoadBalancer{
			TargetGroupArn: aws.String(targetGroups.Canary.TargetGroupARN),
			ContainerName:  aws.String(targetGroups.Canary.ContainerName),
			ContainerPort:  aws.Int32(targetGroups.Canary.ContainerPort),
		}
	}

	return primary, canary, nil
}
