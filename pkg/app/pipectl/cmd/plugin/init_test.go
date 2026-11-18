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

package plugin

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/cli"
)

func TestInit_run(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		cmd       initCommand
		wantFiles []string
		wantErr   bool
	}{
		{
			name:      "default stage",
			cmd:       initCommand{name: "my-plugin", module: "example.com/my-plugin"},
			wantFiles: []string{"go.mod", "main.go", "my_plugin_run_stage.go", "plugin.go", "plugin_test.go", "README.md"},
		},
		{
			name:      "several stages",
			cmd:       initCommand{name: "my-plugin", module: "example.com/my-plugin", stages: []string{"FIRST", "SECOND"}},
			wantFiles: []string{"go.mod", "main.go", "first_stage.go", "second_stage.go", "plugin.go", "plugin_test.go", "README.md"},
		},
		{
			name:    "invalid plugin name",
			cmd:     initCommand{name: "My Plugin", module: "example.com/my-plugin"},
			wantErr: true,
		},
		{
			name:    "invalid stage name",
			cmd:     initCommand{name: "my-plugin", module: "example.com/my-plugin", stages: []string{"not-a-stage"}},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// The plugin goes to --output-dir when it is set.
			dir := filepath.Join(t.TempDir(), "out")
			tc.cmd.outputDir = dir

			err := tc.cmd.run(context.Background(), cli.Input{Logger: zap.NewNop()})
			if tc.wantErr {
				assert.Error(t, err)
				assert.NoDirExists(t, dir)
				return
			}
			require.NoError(t, err)

			entries, err := os.ReadDir(dir)
			require.NoError(t, err)
			got := make([]string, 0, len(entries))
			for _, e := range entries {
				got = append(got, e.Name())
			}
			assert.ElementsMatch(t, tc.wantFiles, got)
		})
	}
}

func TestInit_runRefusesANonEmptyDirectory(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	existing := filepath.Join(dir, "notes.txt")
	require.NoError(t, os.WriteFile(existing, []byte("mine"), 0o644))

	c := initCommand{name: "my-plugin", module: "example.com/my-plugin", outputDir: dir}
	err := c.run(context.Background(), cli.Input{Logger: zap.NewNop()})

	assert.Error(t, err)
	assert.NoFileExists(t, filepath.Join(dir, "go.mod"))
}

func TestNewInitCommand_requiresNameAndModule(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
	}{
		{name: "no flags", args: []string{}},
		{name: "no module", args: []string{"--name", "my-plugin"}},
		{name: "no name", args: []string{"--module", "example.com/my-plugin"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cmd := newInitCommand(&command{})
			cmd.SetArgs(tc.args)
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			err := cmd.Execute()

			require.Error(t, err)
			assert.Contains(t, err.Error(), "required flag(s)")
		})
	}
}
