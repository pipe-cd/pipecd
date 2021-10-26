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

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeployment_ContainTags(t *testing.T) {
	testcases := []struct {
		name       string
		deployment *Deployment
		tags       []string
		want       bool
	}{
		{
			name:       "all given tags aren't contained",
			deployment: &Deployment{Tags: []string{"foo"}},
			tags:       []string{"foo", "bar"},
			want:       false,
		},
		{
			name:       "a tag is contained",
			deployment: &Deployment{Tags: []string{"foo", "bar"}},
			tags:       []string{"foo"},
			want:       true,
		},
		{
			name:       "all tags are contained",
			deployment: &Deployment{Tags: []string{"foo", "bar", "baz"}},
			tags:       []string{"baz", "foo", "bar"},
			want:       true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.deployment.ContainTags(tc.tags)
			assert.Equal(t, tc.want, got)
		})
	}
}
