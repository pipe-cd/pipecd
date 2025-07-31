// Copyright 2025 The PipeCD Authors.
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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeVars(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		deployTargetVars []string
		appVars          []string
		want             []string
	}{
		{
			name:             "empty vars",
			deployTargetVars: []string{},
			appVars:          []string{},
			want:             []string{},
		},
		{
			name:             "only deploy target vars",
			deployTargetVars: []string{"key1=value1", "key2=value2"},
			appVars:          []string{},
			want:             []string{"key1=value1", "key2=value2"},
		},
		{
			name:             "only app vars",
			deployTargetVars: []string{},
			appVars:          []string{"key3=value3", "key4=value4"},
			want:             []string{"key3=value3", "key4=value4"},
		},
		{
			name:             "both deploy target and app vars",
			deployTargetVars: []string{"key1=value1", "key2=value2"},
			appVars:          []string{"key3=value3", "key4=value4"},
			want:             []string{"key1=value1", "key2=value2", "key3=value3", "key4=value4"},
		},
		{
			// TODO: Validate duplication ans the want should not contain "key2=valueX"
			name:             "duplicate vars",
			deployTargetVars: []string{"key1=value1", "key2=value2"},
			appVars:          []string{"key2=valueX", "key3=value3"},
			want:             []string{"key1=value1", "key2=value2", "key2=valueX", "key3=value3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := mergeVars(tt.deployTargetVars, tt.appVars)
			assert.Equal(t, tt.want, got)
		})
	}
}
