// Copyright 2026 The PipeCD Authors.
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

package pluginscaffold

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewModel(t *testing.T) {
	t.Parallel()

	const module = "example.com/my-plugin"

	tests := []struct {
		name       string
		opts       Options
		wantStages []string
		wantErr    string
	}{
		{
			name:       "the default stage is named after the plugin",
			opts:       Options{Name: "my-plugin", Module: module},
			wantStages: []string{"MY_PLUGIN_RUN"},
		},
		{
			name:       "stages keep the given order",
			opts:       Options{Name: "myplugin", Module: module, Stages: []string{"B_STAGE", "A_STAGE"}},
			wantStages: []string{"B_STAGE", "A_STAGE"},
		},
		{
			name:       "spaces around a stage name are not part of it",
			opts:       Options{Name: "myplugin", Module: module, Stages: []string{"A_STAGE", " B_STAGE"}},
			wantStages: []string{"A_STAGE", "B_STAGE"},
		},
		{
			name:       "names that Go treats specially in a file name are fine",
			opts:       Options{Name: "myplugin", Module: module, Stages: []string{"RUN_TEST", "RUN_LINUX_ARM64", "PLUGIN", "MAIN"}},
			wantStages: []string{"RUN_TEST", "RUN_LINUX_ARM64", "PLUGIN", "MAIN"},
		},
		{
			name:    "plugin name is empty",
			opts:    Options{Module: module},
			wantErr: "plugin name must be set",
		},
		{
			name:    "plugin name has an uppercase letter",
			opts:    Options{Name: "MyPlugin", Module: module},
			wantErr: "plugin name \"MyPlugin\" is invalid",
		},
		{
			name:    "plugin name starts with a digit",
			opts:    Options{Name: "1plugin", Module: module},
			wantErr: "is invalid",
		},
		{
			name:    "plugin name has a path separator",
			opts:    Options{Name: "../plugin", Module: module},
			wantErr: "is invalid",
		},
		{
			name:    "module path is empty",
			opts:    Options{Name: "myplugin"},
			wantErr: "module path must be set",
		},
		{
			name:    "module path has a space",
			opts:    Options{Name: "myplugin", Module: "example.com/my plugin"},
			wantErr: "not allowed in a Go module path",
		},
		{
			name:    "module path has a line break",
			opts:    Options{Name: "myplugin", Module: "example.com/plugin\nreplace x => y"},
			wantErr: "not allowed in a Go module path",
		},
		{
			name:    "module path has an empty element",
			opts:    Options{Name: "myplugin", Module: "example.com//plugin"},
			wantErr: "must not be empty",
		},
		{
			name:    "module path has an element that starts with a dot",
			opts:    Options{Name: "myplugin", Module: "example.com/.hidden"},
			wantErr: "must not be empty or start or end with a dot",
		},
		{
			name:    "stage name is empty",
			opts:    Options{Name: "myplugin", Module: module, Stages: []string{"A_STAGE", " "}},
			wantErr: "stage name must not be empty",
		},
		{
			name:    "stage name is lowercase",
			opts:    Options{Name: "myplugin", Module: module, Stages: []string{"my_stage"}},
			wantErr: "use UPPER_SNAKE_CASE",
		},
		{
			name:    "stage name has a dash",
			opts:    Options{Name: "myplugin", Module: module, Stages: []string{"MY-STAGE"}},
			wantErr: "use UPPER_SNAKE_CASE",
		},
		{
			name:    "stage name has a double underscore",
			opts:    Options{Name: "myplugin", Module: module, Stages: []string{"MY__STAGE"}},
			wantErr: "use UPPER_SNAKE_CASE",
		},
		{
			name:    "stage is listed twice",
			opts:    Options{Name: "myplugin", Module: module, Stages: []string{"A_STAGE", "A_STAGE"}},
			wantErr: "listed more than once",
		},
		{
			name:    "two stages would generate the same identifier",
			opts:    Options{Name: "myplugin", Module: module, Stages: []string{"A1", "A_1"}},
			wantErr: "same Go identifier",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m, err := newModel(tc.opts)
			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
				return
			}
			require.NoError(t, err)

			got := make([]string, 0, len(m.Stages))
			for _, s := range m.Stages {
				got = append(got, s.Name)
			}
			assert.Equal(t, tc.wantStages, got)
			assert.Equal(t, defaultSDKVersion, m.SDKVersion)
			assert.Equal(t, defaultGoVersion, m.GoVersion)
		})
	}
}
