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
)

func TestNewStage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		stage string
		want  stage
	}{
		{
			name:  "one word",
			stage: "WAIT",
			want:  stage{Name: "WAIT", Const: "stageWait", Method: "executeWait", File: "wait_stage.go"},
		},
		{
			name:  "several words",
			stage: "SCRIPT_RUN",
			want:  stage{Name: "SCRIPT_RUN", Const: "stageScriptRun", Method: "executeScriptRun", File: "script_run_stage.go"},
		},
		{
			name:  "digits",
			stage: "STEP_2",
			want:  stage{Name: "STEP_2", Const: "stageStep2", Method: "executeStep2", File: "step_2_stage.go"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.want, newStage(tc.stage))
		})
	}
}

func TestStageFileName(t *testing.T) {
	t.Parallel()

	// None of these may end in a suffix that Go reads as a build constraint
	// (_test.go, _linux.go, _arm64.go) and none may be the file of the plugin.
	tests := []struct {
		name  string
		stage string
		want  string
	}{
		{name: "test suffix", stage: "RUN_TEST", want: "run_test_stage.go"},
		{name: "os suffix", stage: "RUN_LINUX", want: "run_linux_stage.go"},
		{name: "os and arch suffix", stage: "RUN_LINUX_ARM64", want: "run_linux_arm64_stage.go"},
		{name: "named like the plugin file", stage: "PLUGIN", want: "plugin_stage.go"},
		{name: "named like the main file", stage: "MAIN", want: "main_stage.go"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.want, stageFileName(tc.stage))
		})
	}
}

func TestStagePrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		plugin string
		want   string
	}{
		{name: "one word", plugin: "myplugin", want: "MYPLUGIN"},
		{name: "dash", plugin: "my-plugin", want: "MY_PLUGIN"},
		{name: "underscore", plugin: "my_plugin", want: "MY_PLUGIN"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.want, stagePrefix(tc.plugin))
		})
	}
}
