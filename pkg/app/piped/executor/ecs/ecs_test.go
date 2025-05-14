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

package ecs

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"

	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/ecs"
)

func TestFindRemovedTags(t *testing.T) {
	currentTags := []types.Tag{
		{Key: strPtr(provider.LabelManagedBy), Value: strPtr("piped")},
		{Key: strPtr("region"), Value: strPtr("us-west-1")},
		{Key: strPtr("project"), Value: strPtr("abc")},
	}

	desiredTags := []types.Tag{
		{Key: strPtr("project"), Value: strPtr("abc")},
	}

	got := findRemovedTags(currentTags, desiredTags)
	assert.ElementsMatch(t, []string{"region"}, got)
}

func strPtr(s string) *string {
	return &s
}
