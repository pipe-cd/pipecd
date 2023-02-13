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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"

	"github.com/pipe-cd/pipecd/pkg/config"
)

var ErrNoTargetGroup = errors.New("no target group")

func loadTargetGroups(targetGroups config.ECSTargetGroups) (*types.LoadBalancer, *types.LoadBalancer, error) {
	if len(targetGroups.Primary) == 0 {
		return nil, nil, ErrNoTargetGroup
	}

	// Decode Primary target group config.
	primary := &types.LoadBalancer{}
	primaryDecoder := json.NewDecoder(bytes.NewReader(targetGroups.Primary))
	primaryDecoder.DisallowUnknownFields()
	err := primaryDecoder.Decode(primary)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid primary target group definition given: %v", err)
	}

	canary := &types.LoadBalancer{}
	if len(targetGroups.Canary) > 0 {
		canaryDecoder := json.NewDecoder(bytes.NewReader(targetGroups.Canary))
		canaryDecoder.DisallowUnknownFields()
		err := canaryDecoder.Decode(canary)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid canary target group definition given: %v", err)
		}
	}

	return primary, canary, nil
}
