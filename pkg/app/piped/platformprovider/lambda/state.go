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

package lambda

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/lambda/types"

	"github.com/pipe-cd/pipecd/pkg/model"
)

// MakeFunctionResourceState creates LambdaResourceState of a Function.
func MakeFunctionResourceState(fc *types.FunctionConfiguration) *model.LambdaResourceState {
	var healthStatus model.LambdaResourceState_HealthStatus

	switch fc.State {
	case types.StateActive:
		healthStatus = model.LambdaResourceState_HEALTHY
	case types.StatePending, types.StateInactive, types.StateFailed:
		healthStatus = model.LambdaResourceState_OTHER
	default:
		healthStatus = model.LambdaResourceState_UNKNOWN
	}

	healthDesc := fmt.Sprintf("Function's state is %s.", fc.State)
	if fc.StateReason != nil {
		healthDesc = fmt.Sprintf("%s StateReason: %s, StateReasonCode: %s", healthDesc, *fc.StateReason, fc.StateReasonCode)
	}

	return &model.LambdaResourceState{
		Id:   *fc.FunctionArn,
		Name: *fc.FunctionName,

		Kind: "Function",

		HealthStatus:      healthStatus,
		HealthDescription: healthDesc,
	}
}
