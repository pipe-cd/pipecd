// Copyright 2020 The PipeCD Authors.
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

package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ImageWatcherSpec struct {
	Targets []ImageWatcherTarget `json:"targets"`
}

// ImageWatcherTarget provides information to compare the latest tags.
//
// Image watcher typically compares the "Image" in the "Provider" and
// an image defined at the "Field" in the "FilePath".
type ImageWatcherTarget struct {
	// The name of Image provider.
	Provider string `json:"provider"`
	// Fully qualified image name.
	Image string `json:"image"`
	// The path to the file to be updated.
	FilePath string `json:"filePath"`
	// The path to the field to be updated. It requires to start
	// with `$` which represents the root element. e.g. `$.foo.bar[0].baz`.
	Field string `json:"field"`
}

// LoadImageWatcher finds the config files for the image watcher in the .pipe
// directory first up. And returns parsed config after merging the targets.
// Only one of includes or excludes can be used. ErrNotFound is returned if not found.
func LoadImageWatcher(repoRoot string, includes, excludes []string) (*ImageWatcherSpec, error) {
	dir := filepath.Join(repoRoot, SharedConfigurationDirName)
	files, err := ioutil.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", dir, err)
	}

	spec := &ImageWatcherSpec{
		Targets: make([]ImageWatcherTarget, 0),
	}
	filtered, err := filterImageWatcherFiles(files, includes, excludes)
	if err != nil {
		return nil, fmt.Errorf("failed to filter image watcher files at %s: %w", dir, err)
	}
	for _, f := range filtered {
		if f.IsDir() {
			continue
		}
		path := filepath.Join(dir, f.Name())
		cfg, err := LoadFromYAML(path)
		if err != nil {
			return nil, fmt.Errorf("failed to load config file %s: %w", path, err)
		}
		if cfg.Kind == KindImageWatcher {
			spec.Targets = append(spec.Targets, cfg.ImageWatcherSpec.Targets...)
		}
	}
	if len(spec.Targets) == 0 {
		return nil, ErrNotFound
	}
	if err := spec.Validate(); err != nil {
		return nil, err
	}

	return spec, nil
}

// filterImageWatcherFiles filters the given files based on the given Includes and Excludes.
// Excludes are prioritized if both Excludes and Includes are given.
func filterImageWatcherFiles(files []os.FileInfo, includes, excludes []string) ([]os.FileInfo, error) {
	if len(includes) == 0 && len(excludes) == 0 {
		return files, nil
	}

	filtered := make([]os.FileInfo, 0, len(files))
	useWhitelist := len(includes) != 0 && len(excludes) == 0
	if useWhitelist {
		whiteList := make(map[string]struct{}, len(includes))
		for _, i := range includes {
			whiteList[i] = struct{}{}
		}
		for _, f := range files {
			if _, ok := whiteList[f.Name()]; ok {
				filtered = append(filtered, f)
			}
		}
		return filtered, nil
	}

	blackList := make(map[string]struct{}, len(excludes))
	for _, e := range excludes {
		blackList[e] = struct{}{}
	}
	for _, f := range files {
		if _, ok := blackList[f.Name()]; !ok {
			filtered = append(filtered, f)
		}
	}
	return filtered, nil
}

func (s *ImageWatcherSpec) Validate() error {
	for _, t := range s.Targets {
		if t.Provider == "" {
			return fmt.Errorf("provider must not be empty")
		}
		if t.Image == "" {
			return fmt.Errorf("image must not be empty")
		}
		if t.FilePath == "" {
			return fmt.Errorf("filePath must not be empty")
		}
		if t.Field == "" {
			return fmt.Errorf("field must not be empty")
		}
	}
	return nil
}
