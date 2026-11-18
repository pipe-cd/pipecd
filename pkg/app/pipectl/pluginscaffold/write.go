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
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	dirMode  = 0o755
	fileMode = 0o644
)

// Write writes files below dir.
//
// The directory must not exist or must be empty, so a scaffold never lands on
// top of existing work. That check runs before the first byte is written. An
// I/O error while writing (a full disk, say) can still stop it half way.
func Write(dir string, files []File) error {
	if dir == "" {
		return errors.New("output directory must be set")
	}

	info, err := os.Stat(dir)
	switch {
	case err == nil:
		if !info.IsDir() {
			return fmt.Errorf("%s exists and is not a directory", dir)
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", dir, err)
		}
		if len(entries) > 0 {
			return fmt.Errorf("directory %s is not empty", dir)
		}
	case errors.Is(err, fs.ErrNotExist):
		// Created below.
	default:
		return fmt.Errorf("failed to inspect %s: %w", dir, err)
	}

	for _, f := range files {
		target := filepath.Join(dir, filepath.FromSlash(f.Path))
		if err := os.MkdirAll(filepath.Dir(target), dirMode); err != nil {
			return fmt.Errorf("failed to create the directory of %s: %w", f.Path, err)
		}
		if err := os.WriteFile(target, f.Content, fileMode); err != nil {
			return fmt.Errorf("failed to write %s: %w", f.Path, err)
		}
	}
	return nil
}
