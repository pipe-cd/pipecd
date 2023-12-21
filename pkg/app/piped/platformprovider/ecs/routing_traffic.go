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
	"sort"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
)

type RoutingTrafficConfig []targetGroupWeight

type targetGroupWeight struct {
	TargetGroupArn string
	Weight         int
}

func (c RoutingTrafficConfig) hasSameTargets(forwardActionTargets []types.TargetGroupTuple) bool {
	if len(c) != len(forwardActionTargets) {
		return false
	}

	// Sort to the both slices have the same target groups.
	sort.Slice(c, func(i, j int) bool {
		return c[i].TargetGroupArn < c[j].TargetGroupArn
	})
	sort.Slice(forwardActionTargets, func(i, j int) bool {
		return *forwardActionTargets[i].TargetGroupArn < *forwardActionTargets[j].TargetGroupArn
	})

	for i := range c {
		if c[i].TargetGroupArn != *forwardActionTargets[i].TargetGroupArn {
			return false
		}
	}

	return true
}
