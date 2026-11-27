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
	"bytes"
	"context"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// envBuildTest turns on TestGeneratedPluginBuilds. It needs the Go toolchain and
// a way to resolve github.com/pipe-cd/piped-plugin-sdk-go (the network or a warm
// module cache), so it is opt in:
//
//	PLUGIN_SCAFFOLD_BUILD_TEST=1 go test ./pkg/app/pipectl/pluginscaffold
const envBuildTest = "PLUGIN_SCAFFOLD_BUILD_TEST"

// TestGeneratedPluginBuilds is the test that matters most: whatever the scaffold
// writes has to compile, vet clean, be gofmt-ed and pass its own tests against
// the real SDK.
func TestGeneratedPluginBuilds(t *testing.T) {
	t.Parallel()

	if os.Getenv(envBuildTest) == "" {
		t.Skipf("set %s=1 to build the generated plugins", envBuildTest)
	}
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("the go command is not on PATH")
	}
	gofmtBin, err := exec.LookPath("gofmt")
	if err != nil {
		t.Skip("the gofmt command is not on PATH")
	}

	tests := []struct {
		name string
		opts Options
	}{
		{
			name: "default stage",
			opts: Options{Name: "myplugin", Module: "example.com/myplugin"},
		},
		{
			name: "several stages",
			opts: Options{Name: "my-plugin", Module: "example.com/my-plugin", Stages: []string{"MY_PREPARE", "MY_CHECK"}},
		},
		{
			name: "stages named like files that Go treats specially",
			opts: Options{Name: "myplugin", Module: "example.com/myplugin", Stages: []string{"RUN_TEST", "RUN_LINUX_ARM64", "PLUGIN", "MAIN"}},
		},
		{
			name: "a stage named like the unknown stage the generated test uses",
			opts: Options{Name: "myplugin", Module: "example.com/myplugin", Stages: []string{"NOT_A_STAGE_OF_THIS_PLUGIN"}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			files, err := Generate(tc.opts)
			require.NoError(t, err)
			dir := t.TempDir()
			require.NoError(t, Write(dir, files))

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			run := func(name string, args ...string) string {
				t.Helper()

				cmd := exec.CommandContext(ctx, name, args...)
				cmd.Dir = dir
				var out bytes.Buffer
				cmd.Stdout, cmd.Stderr = &out, &out
				require.NoError(t, cmd.Run(), "%s %v failed:\n%s", name, args, out.String())
				return out.String()
			}

			run(goBin, "mod", "tidy")
			run(goBin, "vet", "./...")
			run(goBin, "test", "./...")
			require.Empty(t, run(gofmtBin, "-l", "."), "these generated files are not gofmt-ed")
		})
	}
}
