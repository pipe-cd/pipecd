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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var goldenStageOpts = Options{
	PluginName: "demo",
	ModulePath: "github.com/example/piped-plugin-demo",
	Kind:       KindStage,
	Stages:     []string{"DEMO_WAIT"},
}

var goldenDeploymentOpts = Options{
	PluginName: "demo",
	ModulePath: "github.com/example/piped-plugin-demo",
	Kind:       KindDeployment,
	Stages:     []string{"DEMO_SYNC", "DEMO_ROLLBACK"},
}

func TestGenerate_golden(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		opts      Options
		goldenDir string
	}{
		{
			name:      "stage",
			opts:      goldenStageOpts,
			goldenDir: "testdata/golden/stage",
		},
		{
			name:      "deployment",
			opts:      goldenDeploymentOpts,
			goldenDir: "testdata/golden/deployment",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			files, err := Generate(tc.opts)
			require.NoError(t, err)

			if os.Getenv("UPDATE_GOLDEN") != "" {
				updateGolden(t, tc.goldenDir, files)
			}

			for _, f := range files {
				goldenPath := filepath.Join(tc.goldenDir, filepath.FromSlash(f.Path))
				want, err := os.ReadFile(goldenPath)
				require.NoError(t, err, "missing golden file %s (run UPDATE_GOLDEN=1 go test ./pkg/app/pipectl/pluginscaffold -run Golden -count=1)", goldenPath)
				require.Equal(t, string(want), string(f.Content), "file %s", f.Path)
			}

			var goldenCount int
			err = filepath.WalkDir(tc.goldenDir, func(path string, d os.DirEntry, walkErr error) error {
				if walkErr != nil {
					return walkErr
				}
				if !d.IsDir() {
					goldenCount++
				}
				return nil
			})
			require.NoError(t, err)
			require.Equal(t, len(files), goldenCount, "golden file count mismatch for %s", tc.goldenDir)
		})
	}
}

func updateGolden(t *testing.T, goldenDir string, files []File) {
	t.Helper()

	for _, f := range files {
		path := filepath.Join(goldenDir, filepath.FromSlash(f.Path))
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, f.Content, 0o644))
	}
}
