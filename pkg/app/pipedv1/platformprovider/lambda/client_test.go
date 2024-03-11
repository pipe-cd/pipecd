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

	"github.com/stretchr/testify/assert"
)

func TestMakeFlowControlTagsMap(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name                 string
		remoteTags           map[string]string
		definedTags          map[string]string
		wantedNewDefinedTags map[string]string
		wantedUpdatedTags    map[string]string
		wantedRemovedTags    map[string]string
	}{
		{
			name: "has only updated tags",
			remoteTags: map[string]string{
				"app":      "simple",
				"function": "code",
			},
			definedTags: map[string]string{
				"app":      "simple-app",
				"function": "code",
			},
			wantedNewDefinedTags: map[string]string{},
			wantedUpdatedTags: map[string]string{
				"app": "simple-app",
			},
			wantedRemovedTags: map[string]string{},
		},
		{
			name: "has only remove tags",
			remoteTags: map[string]string{
				"app":      "simple",
				"function": "code",
			},
			definedTags: map[string]string{
				"app": "simple",
			},
			wantedNewDefinedTags: map[string]string{},
			wantedUpdatedTags:    map[string]string{},
			wantedRemovedTags: map[string]string{
				"function": "code",
			},
		},
		{
			name: "has only newly defined tags",
			remoteTags: map[string]string{
				"app": "simple",
			},
			definedTags: map[string]string{
				"app":      "simple",
				"function": "code",
			},
			wantedNewDefinedTags: map[string]string{
				"function": "code",
			},
			wantedUpdatedTags: map[string]string{},
			wantedRemovedTags: map[string]string{},
		},
		{
			name: "complex defined tags",
			remoteTags: map[string]string{
				"app":      "simple",
				"function": "code",
			},
			definedTags: map[string]string{
				"foo": "bar",
				"app": "simple-app",
				"bar": "foo",
			},
			wantedNewDefinedTags: map[string]string{
				"foo": "bar",
				"bar": "foo",
			},
			wantedUpdatedTags: map[string]string{
				"app": "simple-app",
			},
			wantedRemovedTags: map[string]string{
				"function": "code",
			},
		},
		{
			name: "defined tags is nil",
			remoteTags: map[string]string{
				"app":      "simple",
				"function": "code",
			},
			definedTags:          nil,
			wantedNewDefinedTags: map[string]string{},
			wantedUpdatedTags:    map[string]string{},
			wantedRemovedTags: map[string]string{
				"app":      "simple",
				"function": "code",
			},
		},
		{
			name:                 "remote tags is empty and defined tags is nil",
			remoteTags:           map[string]string{},
			definedTags:          nil,
			wantedNewDefinedTags: map[string]string{},
			wantedUpdatedTags:    map[string]string{},
			wantedRemovedTags:    map[string]string{},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			newDefinedTags, updatedTags, removedTags := makeFlowControlTagsMaps(tc.remoteTags, tc.definedTags)
			assert.Equal(t, tc.wantedNewDefinedTags, newDefinedTags)
			assert.Equal(t, tc.wantedUpdatedTags, updatedTags)
			assert.Equal(t, tc.wantedRemovedTags, removedTags)
		})
	}
}
