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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testFiles = []File{
	{Path: "go.mod", Content: []byte("module example.com/x\n")},
	{Path: "sub/dir/file.go", Content: []byte("package dir\n")},
}

func TestWrite(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// dir is the directory to write to, relative to a fresh temp dir.
		dir string
		// exists creates dir, empty, before the write.
		exists bool
	}{
		{name: "a directory that does not exist", dir: filepath.Join("new", "plugin")},
		{name: "an empty directory", dir: "empty", exists: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dir := filepath.Join(t.TempDir(), tc.dir)
			if tc.exists {
				require.NoError(t, os.Mkdir(dir, 0o755))
			}

			require.NoError(t, Write(dir, testFiles))

			for _, f := range testFiles {
				got, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(f.Path)))
				require.NoError(t, err)
				assert.Equal(t, f.Content, got, f.Path)
			}
		})
	}
}

func TestWrite_refusesANonEmptyDirectory(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	existing := filepath.Join(dir, "notes.txt")
	require.NoError(t, os.WriteFile(existing, []byte("mine"), 0o644))

	err := Write(dir, testFiles)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "is not empty")
	// Nothing was written and the existing file is untouched.
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	assert.Len(t, entries, 1)
	got, err := os.ReadFile(existing)
	require.NoError(t, err)
	assert.Equal(t, "mine", string(got))
}

func TestWrite_refusesAFile(t *testing.T) {
	t.Parallel()

	target := filepath.Join(t.TempDir(), "plugin")
	require.NoError(t, os.WriteFile(target, []byte("a file"), 0o644))

	err := Write(target, testFiles)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "is not a directory")
}

func TestWrite_requiresADirectory(t *testing.T) {
	t.Parallel()

	assert.Error(t, Write("", testFiles))
}
