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

package config

import (
	"fmt"

	"github.com/pipe-cd/pipecd/pkg/filematcher"
)

type EventWatcherConfig struct {
	// Matcher represents which event will be handled.
	Matcher EventWatcherMatcher `json:"matcher"`
	// Handler represents how the matched event will be handled.
	Handler EventWatcherHandler `json:"handler"`
}

type EventWatcherMatcher struct {
	// The handled event name.
	Name string `json:"name"`
	// Additional attributes of event. This can make an event definition
	// unique even if the one with the same name exists.
	Labels map[string]string `json:"labels"`
}

type EventWatcherHandler struct {
	// The handler type of event watcher.
	Type EventWatcherHandlerType `json:"type,omitempty"`
	// The config for event watcher handler.
	Config EventWatcherHandlerConfig `json:"config"`
}

type EventWatcherHandlerConfig struct {
	// The commit message used to push after replacing values.
	// Default message is used if not given.
	CommitMessage string `json:"commitMessage,omitempty"`
	// Whether to create a new branch or not when event watcher commits changes.
	MakePullRequest bool `json:"makePullRequest,omitempty"`
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
	// The regex string specifying what should be replaced.
	// Only the first capturing group enclosed by `()` will be replaced with the new value.
	// e.g. "host.xz/foo/bar:(v[0-9].[0-9].[0-9])"
	Regex string `json:"regex"`
}

// EventWatcherHandlerType represents the type of an event watcher handler.
type EventWatcherHandlerType string

const (
	// EventWatcherHandlerTypeGitUpdate represents the handler type for git updating.
	EventWatcherHandlerTypeGitUpdate = "GIT_UPDATE"
)

// filterEventWatcherFiles filters the given files based on the given Includes and Excludes.
// Excludes are prioritized if both Excludes and Includes are given.
func filterEventWatcherFiles(files, includePatterns, excludePatterns []string) ([]string, error) {
	if len(includePatterns) == 0 && len(excludePatterns) == 0 {
		return files, nil
	}

	filtered := make([]string, 0, len(files))

	// Use include patterns
	if len(includePatterns) != 0 && len(excludePatterns) == 0 {
		matcher, err := filematcher.NewPatternMatcher(includePatterns)
		if err != nil {
			return nil, fmt.Errorf("failed to create a matcher object: %w", err)
		}
		for _, f := range files {
			if matcher.Matches(f) {
				filtered = append(filtered, f)
			}
		}
		return filtered, nil
	}

	// Use exclude patterns
	matcher, err := filematcher.NewPatternMatcher(excludePatterns)
	if err != nil {
		return nil, fmt.Errorf("failed to create a matcher object: %w", err)
	}
	for _, f := range files {
		if matcher.Matches(f) {
			continue
		}
		filtered = append(filtered, f)
	}
	return filtered, nil
}
