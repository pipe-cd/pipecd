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
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate_stage(t *testing.T) {
	t.Parallel()

	files, err := Generate(Options{
		PluginName: "demo",
		ModulePath: "github.com/example/piped-plugin-demo",
		Kind:       KindStage,
		Stages:     []string{"DEMO_WAIT"},
	})
	require.NoError(t, err)
	assert.Contains(t, filePaths(files), "main.go")
	assert.Contains(t, filePaths(files), "plugin.go")
}

func TestGenerate_deployment(t *testing.T) {
	t.Parallel()

	files, err := Generate(Options{
		PluginName: "demo",
		ModulePath: "github.com/example/piped-plugin-demo",
		Kind:       KindDeployment,
		Stages:     []string{"DEMO_SYNC", "DEMO_ROLLBACK"},
	})
	require.NoError(t, err)
	paths := filePaths(files)
	assert.Contains(t, paths, filepath.Join("deployment", "plugin.go"))
	assert.Contains(t, paths, filepath.Join("deployment", "demo_sync.go"))
	assert.Contains(t, paths, filepath.Join("deployment", "demo_rollback.go"))
}

func TestWrite_and_build_stage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration build in short mode")
	}
	t.Parallel()

	testWriteAndBuild(t, goldenStageOpts)
}

func TestWrite_and_build_deployment(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration build in short mode")
	}
	t.Parallel()

	testWriteAndBuild(t, goldenDeploymentOpts)
}

func testWriteAndBuild(t *testing.T, opts Options) {
	t.Helper()

	dir := filepath.Join(t.TempDir(), "plugin")
	opts.OutputDir = dir

	files, err := Generate(opts)
	require.NoError(t, err)
	require.NoError(t, Write(opts, files))

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))

	cmd = exec.Command("go", "build", "-o", filepath.Join(dir, "plugin"), ".")
	cmd.Dir = dir
	out, err = cmd.CombinedOutput()
	require.NoError(t, err, string(out))

	_, err = os.Stat(filepath.Join(dir, "plugin"))
	require.NoError(t, err)
}

func filePaths(files []File) []string {
	out := make([]string, len(files))
	for i, f := range files {
		out[i] = f.Path
	}
	return out
}
