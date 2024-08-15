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
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestMakeFunctionResourceState(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		title              string
		state              types.State
		expectedStatus     model.LambdaResourceState_HealthStatus
		expectedHealthDesc string
	}{
		{
			title:              "active is healthy",
			state:              types.StateActive,
			expectedStatus:     model.LambdaResourceState_HEALTHY,
			expectedHealthDesc: "Function's state is Active.",
		},
		{
			title:              "pending is other",
			state:              types.StatePending,
			expectedStatus:     model.LambdaResourceState_OTHER,
			expectedHealthDesc: "Function's state is Pending.",
		},
		{
			title:              "inactive is other",
			state:              types.StateInactive,
			expectedStatus:     model.LambdaResourceState_OTHER,
			expectedHealthDesc: "Function's state is Inactive.",
		},
		{
			title:              "failed is other",
			state:              types.StateFailed,
			expectedStatus:     model.LambdaResourceState_OTHER,
			expectedHealthDesc: "Function's state is Failed.",
		},
		{
			title:              "else is unknown",
			state:              "dummy-status",
			expectedStatus:     model.LambdaResourceState_UNKNOWN,
			expectedHealthDesc: "Function's state is dummy-status.",
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()
			f := &types.FunctionConfiguration{
				State: tc.state,

				FunctionArn:  aws.String("test-function-arn"),
				FunctionName: aws.String("test-function-name"),
			}
			state := MakeFunctionResourceState(f)

			expected := &model.LambdaResourceState{
				Id:                "test-function-arn",
				Name:              "test-function-name",
				Kind:              "Function",
				HealthStatus:      tc.expectedStatus,
				HealthDescription: tc.expectedHealthDesc,
			}
			assert.Equal(t, expected, state)
		})
	}
}
