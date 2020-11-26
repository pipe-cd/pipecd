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

type ImageWatcherSpec struct {
	Targets []ImageWatcherTarget `json:"targets"`
}

type ImageWatcherTarget struct {
	Provider string `json:"provider"`
	Image    string `json:"image"`
	FilePath string `json:"filePath"`
	Field    string `json:"field"`
}

// LoadImageWatcher finds the config file for the image watcher in the .pipe directory first up.
// And returns parsed config, False is returned as the second returned value if not found.
func LoadImageWatcher(repoRoot string) (*ImageWatcherSpec, bool, error) {
	// TODO: Load image watcher config, referring to AnalysisTemplateSpec
	return nil, false, nil
}

func (s *ImageWatcherSpec) Validate() error {
	return nil
}
