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

func paths(files []File) []string {
	out := make([]string, 0, len(files))
	for _, f := range files {
		out = append(out, f.Path)
	}
	return out
}

func TestGenerate_files(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts Options
		want []string
	}{
		{
			name: "default stage",
			opts: Options{Name: "my-plugin", Module: "example.com/my-plugin"},
			want: []string{"README.md", "go.mod", "main.go", "my_plugin_run_stage.go", "plugin.go", "plugin_test.go"},
		},
		{
			name: "one executor file per stage",
			opts: Options{Name: "myplugin", Module: "example.com/myplugin", Stages: []string{"FIRST", "SECOND"}},
			want: []string{"README.md", "first_stage.go", "go.mod", "main.go", "plugin.go", "plugin_test.go", "second_stage.go"},
		},
		{
			name: "stages named like files that Go treats specially",
			opts: Options{Name: "myplugin", Module: "example.com/myplugin", Stages: []string{"RUN_TEST", "RUN_LINUX_ARM64", "PLUGIN", "MAIN"}},
			want: []string{
				"README.md", "go.mod", "main.go", "main_stage.go", "plugin.go", "plugin_stage.go",
				"plugin_test.go", "run_linux_arm64_stage.go", "run_test_stage.go",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			files, err := Generate(tc.opts)
			require.NoError(t, err)

			assert.Equal(t, tc.want, paths(files))
			for _, f := range files {
				assert.NotEmpty(t, f.Content, f.Path)
				// A leftover action means a template used a field that did not render.
				assert.NotContains(t, string(f.Content), "{{", f.Path)
				assert.NotContains(t, string(f.Content), "<no value>", f.Path)
			}
		})
	}
}

func TestGenerate_content(t *testing.T) {
	t.Parallel()

	files, err := Generate(Options{
		Name:   "my-plugin",
		Module: "example.com/my-plugin",
		Stages: []string{"MY_PREPARE", "MY_CHECK"},
	})
	require.NoError(t, err)

	content := make(map[string]string, len(files))
	for _, f := range files {
		content[f.Path] = string(f.Content)
	}

	assert.Contains(t, content["go.mod"], "module example.com/my-plugin\n")
	assert.Contains(t, content["go.mod"], "go "+defaultGoVersion+"\n")
	assert.Contains(t, content["go.mod"], "github.com/pipe-cd/piped-plugin-sdk-go "+defaultSDKVersion)
	assert.NotContains(t, content["go.mod"], "replace")

	// Every stage is defined, executed and has its own file.
	for _, stageName := range []string{"MY_PREPARE", "MY_CHECK"} {
		assert.Contains(t, content["plugin.go"], `"`+stageName+`"`)
	}
	assert.Contains(t, content["plugin.go"], "case stageMyPrepare:")
	assert.Contains(t, content["plugin.go"], "case stageMyCheck:")
	assert.Contains(t, content["my_prepare_stage.go"], "func (p *plugin) executeMyPrepare(")
	assert.Contains(t, content["my_check_stage.go"], "func (p *plugin) executeMyCheck(")

	assert.Contains(t, content["README.md"], tutorialURL)
	assert.Contains(t, content["README.md"], `- name: "MY_PREPARE"`)
}

func TestGenerate_isDeterministic(t *testing.T) {
	t.Parallel()

	opts := Options{Name: "myplugin", Module: "example.com/myplugin", Stages: []string{"B_STAGE", "A_STAGE"}}
	first, err := Generate(opts)
	require.NoError(t, err)
	second, err := Generate(opts)
	require.NoError(t, err)

	assert.Equal(t, first, second)
}

func TestGenerate_invalidOptions(t *testing.T) {
	t.Parallel()

	files, err := Generate(Options{Name: "My Plugin", Module: "example.com/myplugin"})

	assert.Error(t, err)
	assert.Nil(t, files)
}

func TestGenerate_readmeQuotesNames(t *testing.T) {
	t.Parallel()

	// Unquoted, YAML reads these as a boolean or null, not as the name.
	files, err := Generate(Options{Name: "yes", Module: "example.com/yes", Stages: []string{"NULL"}})
	require.NoError(t, err)

	var readme string
	for _, f := range files {
		if f.Path == "README.md" {
			readme = string(f.Content)
		}
	}
	assert.Contains(t, readme, `- name: "yes"`)
	assert.Contains(t, readme, `- name: "NULL"`)
}
