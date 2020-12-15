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

type ImageWatcherTarget struct {
	Provider string `json:"provider"`
	Image    string `json:"image"`
	FilePath string `json:"filePath"`
	Field    string `json:"field"`
}

// LoadImageWatcher finds the config files for the image watcher in the .pipe
// directory first up. And returns parsed config after merging the targets.
// Only one of includes or excludes can be used.
// False is returned as the second returned value if not found.
func LoadImageWatcher(repoRoot string, includes, excludes []string) (*ImageWatcherSpec, bool, error) {
	dir := filepath.Join(repoRoot, SharedConfigurationDirName)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, false, nil
	}

	spec := &ImageWatcherSpec{
		Targets: make([]ImageWatcherTarget, 0),
	}
	filtered, err := filterImageWatcherFiles(files, includes, excludes)
	if err != nil {
		return nil, false, fmt.Errorf("failed to filter image watcher files at %s: %w", dir, err)
	}
	for _, f := range filtered {
		if f.IsDir() {
			continue
		}
		path := filepath.Join(dir, f.Name())
		cfg, err := LoadFromYAML(path)
		if err != nil {
			return nil, false, fmt.Errorf("failed to load config file %s: %w", path, err)
		}
		if cfg.Kind == KindImageWatcher {
			spec.Targets = append(spec.Targets, cfg.ImageWatcherSpec.Targets...)
		}
	}
	if len(spec.Targets) == 0 {
		return nil, false, nil
	}

	return spec, true, nil
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
	return nil
}
