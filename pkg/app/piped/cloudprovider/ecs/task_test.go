// Copyright 2021 The PipeCD Authors.
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

func TestParseTaskDefinition(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		input       string
		expected    types.TaskDefinition
		expectedErr bool
	}{
		{
			name: "yaml format input",
			input: `
family: nginx-canary-fam-1
compatibilities:
  - FARGATE
networkMode: awsvpc
memory: 512
cpu: 256
`,
			expected: types.TaskDefinition{
				Family:          aws.String("nginx-canary-fam-1"),
				Compatibilities: []types.Compatibility{types.CompatibilityFargate},
				NetworkMode:     types.NetworkModeAwsvpc,
				Memory:          aws.String("512"),
				Cpu:             aws.String("256"),
			},
		},
		{
			name: "json format input",
			input: `
{
  "family": "nginx-canary-fam-1",
  "compatibilities": [
    "FARGATE"
  ],
  "networkMode": "awsvpc",
  "memory": 512,
  "cpu": 256
}
`,
			expected: types.TaskDefinition{
				Family:          aws.String("nginx-canary-fam-1"),
				Compatibilities: []types.Compatibility{types.CompatibilityFargate},
				NetworkMode:     types.NetworkModeAwsvpc,
				Memory:          aws.String("512"),
				Cpu:             aws.String("256"),
			},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseTaskDefinition([]byte(tc.input))
			assert.Equal(t, tc.expectedErr, err != nil)
			assert.Equal(t, tc.expected, got)
		})
	}
}
