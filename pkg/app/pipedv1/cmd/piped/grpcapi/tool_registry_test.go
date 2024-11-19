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

package grpcapi

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToolRegistry_InstallTool(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	toolsDir, err := os.MkdirTemp(t.TempDir(), "tools")
	require.NoError(t, err)
	tmpDir, err := os.MkdirTemp(t.TempDir(), "tmp")
	require.NoError(t, err)

	registry, err := newToolRegistry(toolsDir, tmpDir)
	require.NoError(t, err)

	tests := []struct {
		name        string
		toolName    string
		toolVersion string
		script      string
		wantErr     bool
	}{
		{
			name:        "valid script",
			toolName:    "tool-a",
			toolVersion: "1.0.0",
			script:      `touch {{ .OutPath }}`,
			wantErr:     false,
		},
		{
			name:        "output is not found",
			toolName:    "tool-b",
			toolVersion: "1.0.0",
			script:      "exit 0",
			wantErr:     true,
		},
		{
			name:        "script failed",
			toolName:    "tool-c",
			toolVersion: "1.0.0",
			script:      "exit 1",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			out, err := registry.InstallTool(ctx, tt.toolName, tt.toolVersion, tt.script)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.FileExists(t, out)
			assert.True(t, strings.HasSuffix(out, tt.toolName+"-"+tt.toolVersion), "output path should have the tool {name}-{version}, got %s", out)
		})
	}
}

func TestToolRegistry_InstallTool_CacheHit(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	toolsDir, err := os.MkdirTemp(t.TempDir(), "tools")
	require.NoError(t, err)
	tmpDir, err := os.MkdirTemp(t.TempDir(), "tmp")
	require.NoError(t, err)

	registry, err := newToolRegistry(toolsDir, tmpDir)
	require.NoError(t, err)

	toolName := "tool-a"
	toolVersion := "1.0.0"

	out, err := registry.InstallTool(ctx, toolName, toolVersion, "touch {{ .OutPath }}") // success
	require.NoError(t, err)
	assert.FileExists(t, out)

	// cache hit and should not run the script, so success again even if the script is invalid.
	// because the cache key is constructed from the tool name and version.
	// we don't expect the script to be changed between the first and second calls, it's just for testing.
	out2, err := registry.InstallTool(ctx, toolName, toolVersion, "exit 1")
	require.NoError(t, err)
	assert.FileExists(t, out2)

	assert.Equal(t, out, out2)
}
