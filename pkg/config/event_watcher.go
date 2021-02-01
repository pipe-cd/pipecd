// Copyright 2021 The PipeCD Authors.
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

type EventWatcherSpec struct {
	Events []EventWatcherEvent `json:"events"`
}

// EventWatcherEvent defines which file will be replaced when the given event happened.
type EventWatcherEvent struct {
	// The event name.
	Name string `json:"name"`
	// Additional attributes of event. This can make an event definition
	// unique even if the one with the same name exists.
	Labels map[string]string `json:"labels"`
	// List of places where will be replaced when the new event matches.
	Replacements []EventWatcherReplacement `json:"replacements"`
}

type EventWatcherReplacement struct {
	// The path to the file to be updated.
	File string `json:"file"`
	// The field to be updated. Only one of these can be used.
	//
	// The YAML path to the field to be updated. It requires to start
	// with `$` which represents the root element. e.g. `$.foo.bar[0].baz`.
	YAMLField string `json:"yamlField"`
	// The JSON path to the field to be updated.
	JSONField string `json:"jsonField"`
	// The HCL path to the field to be updated.
	HCLField string `json:"HCLField"`
}

// LoadEventWatcher gives back parsed EventWatcher config after merging config files placed under
// the .pipe directory. With "includes" and "excludes", you can filter the files included the result.
// "excludes" are prioritized if both "excludes" and "includes" are given. ErrNotFound is returned if not found.
func LoadEventWatcher(repoRoot string, includes, excludes []string) (*EventWatcherSpec, error) {
	dir := filepath.Join(repoRoot, SharedConfigurationDirName)
	files, err := ioutil.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", dir, err)
	}

	// Start merging events defined across multiple files.
	spec := &EventWatcherSpec{
		Events: make([]EventWatcherEvent, 0),
	}
	filtered, err := filterEventWatcherFiles(files, includes, excludes)
	if err != nil {
		return nil, fmt.Errorf("failed to filter event watcher files at %s: %w", dir, err)
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
		if cfg.Kind == KindEventWatcher {
			spec.Events = append(spec.Events, cfg.EventWatcherSpec.Events...)
		}
	}

	if err := spec.Validate(); err != nil {
		return nil, err
	}

	return spec, nil
}

// filterEventWatcherFiles filters the given files based on the given Includes and Excludes.
// Excludes are prioritized if both Excludes and Includes are given.
func filterEventWatcherFiles(files []os.FileInfo, includes, excludes []string) ([]os.FileInfo, error) {
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

func (s *EventWatcherSpec) Validate() error {
	for _, e := range s.Events {
		if err := e.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (e *EventWatcherEvent) Validate() error {
	if e.Name == "" {
		return fmt.Errorf("event name must not be empty")
	}
	if len(e.Replacements) == 0 {
		return fmt.Errorf("there must be at least one replacement to an event")
	}
	for _, r := range e.Replacements {
		if r.File == "" {
			return fmt.Errorf("event %q has a replacement with no file name", e.Name)
		}
		if r.YAMLField == "" && r.JSONField == "" && r.HCLField == "" {
			return fmt.Errorf("event %q has a replacement with no field", e.Name)
		}
		// Check if multiple fields aren't given.
		given := r.YAMLField != ""
		if given && r.JSONField != "" {
			return fmt.Errorf("event %q has multiple fields", e.Name)
		}
		if r.JSONField != "" {
			given = true
		}
		if given && r.HCLField != "" {
			return fmt.Errorf("event %q has multiple fields", e.Name)
		}
	}
	return nil
}
