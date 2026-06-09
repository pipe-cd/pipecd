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
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var stageNamePattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// ValidatePluginName checks the plugin name is a valid identifier segment.
func ValidatePluginName(name string) error {
	if name == "" {
		return fmt.Errorf("plugin name must not be empty")
	}
	for i, r := range name {
		if i == 0 {
			if !unicode.IsLetter(r) {
				return fmt.Errorf("plugin name must start with a letter")
			}
			continue
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' {
			continue
		}
		return fmt.Errorf("plugin name contains invalid character %q", r)
	}
	return nil
}

// ValidateStageNames checks stage names are UPPER_SNAKE_CASE.
func ValidateStageNames(stages []string) error {
	if len(stages) == 0 {
		return fmt.Errorf("at least one stage is required")
	}
	seen := make(map[string]struct{}, len(stages))
	for _, stage := range stages {
		if stage == "" {
			return fmt.Errorf("stage name must not be empty")
		}
		if !stageNamePattern.MatchString(stage) {
			return fmt.Errorf("stage name %q must be UPPER_SNAKE_CASE (e.g. MY_SYNC)", stage)
		}
		if _, ok := seen[stage]; ok {
			return fmt.Errorf("duplicate stage name %q", stage)
		}
		seen[stage] = struct{}{}
	}
	return nil
}

// TypePrefix converts a plugin name to an exported Go type prefix (e.g. my-plugin -> MyPlugin).
func TypePrefix(pluginName string) string {
	parts := strings.FieldsFunc(pluginName, func(r rune) bool {
		return r == '-' || r == '_'
	})
	var b strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		b.WriteString(strings.ToUpper(part[:1]))
		if len(part) > 1 {
			b.WriteString(strings.ToLower(part[1:]))
		}
	}
	if b.Len() == 0 {
		return "Plugin"
	}
	return b.String()
}

// StageFileBase returns a lowercase file-safe name for a stage.
func StageFileBase(stage string) string {
	return strings.ToLower(stage)
}

// StageConstSuffix returns a suffix for stage const identifiers (e.g. DEMO_SYNC -> DemoSync).
func StageConstSuffix(stage string) string {
	return strings.TrimPrefix(StageFuncName(stage), "execute")
}

// StageFuncName returns an exported handler name (e.g. DEMO_SYNC -> executeDemoSync).
func StageFuncName(stage string) string {
	parts := strings.Split(stage, "_")
	var b strings.Builder
	b.WriteString("execute")
	for _, part := range parts {
		if part == "" {
			continue
		}
		b.WriteString(strings.ToUpper(part[:1]))
		if len(part) > 1 {
			b.WriteString(strings.ToLower(part[1:]))
		}
	}
	return b.String()
}

// DefaultModulePath returns the default Go module path for a community plugin.
func DefaultModulePath(pluginName string) string {
	return fmt.Sprintf("github.com/example/piped-plugin-%s", pluginName)
}

// FindRollbackStage returns the first stage name containing "ROLLBACK", or empty if none.
func FindRollbackStage(stages []string) string {
	for _, stage := range stages {
		if strings.Contains(stage, "ROLLBACK") {
			return stage
		}
	}
	return ""
}
