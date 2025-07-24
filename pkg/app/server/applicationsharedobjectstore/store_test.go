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

package applicationsharedobjectstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildPath(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name       string
		appID      string
		pluginName string
		key        string
		want       string
	}{
		{
			name:       "simple path",
			appID:      "app-1",
			pluginName: "plugin-1",
			key:        "key-1",
			want:       "application-shared-objects/app-1/plugin-1/key-1.json",
		},
		{
			name:       "path with special characters",
			appID:      "app/1",
			pluginName: "plugin.1",
			key:        "key_1",
			want:       "application-shared-objects/app/1/plugin.1/key_1.json",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildPath(tc.appID, tc.pluginName, tc.key)
			assert.Equal(t, tc.want, got)
		})
	}
}
