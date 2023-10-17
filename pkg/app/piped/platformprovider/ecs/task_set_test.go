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

func TestIsPipeCDManagedTaskSet(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		ts       *types.TaskSet
		expected bool
	}{
		{
			name: "managed by piped",
			ts: &types.TaskSet{Tags: []types.Tag{
				{Key: aws.String(LabelManagedBy), Value: aws.String(ManagedByPiped)},
			}},
			expected: true,
		},
		{
			name:     "nil tags",
			ts:       &types.TaskSet{},
			expected: false,
		},
		{
			name: "not managed by piped",
			ts: &types.TaskSet{Tags: []types.Tag{
				{Key: aws.String(LabelManagedBy), Value: aws.String("other")},
				{Key: aws.String("hoge"), Value: aws.String("fuga")},
			}},
			expected: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := IsPipeCDManagedTaskSet(tc.ts)
			assert.Equal(t, tc.expected, got)
		})
	}
}
