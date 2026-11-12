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
	"regexp"
	"strings"
)

var (
	// A plugin name is a lowercase word that may be joined by '-' or '_'.
	// Piped itself only requires the name to be non-empty, so this is the
	// scaffold's own convention: it has to be a safe directory name and it is
	// the prefix of the default stage name.
	pluginNameRe = regexp.MustCompile(`^[a-z][a-z0-9]*(?:[-_][a-z0-9]+)*$`)
	// A stage name is UPPER_SNAKE_CASE, like the stages of the official plugins
	// (WAIT, SCRIPT_RUN, TERRAFORM_APPLY).
	stageNameRe = regexp.MustCompile(`^[A-Z][A-Z0-9]*(?:_[A-Z0-9]+)*$`)
	// One element of a Go module path.
	modulePathElemRe = regexp.MustCompile(`^[A-Za-z0-9._~-]+$`)
)

func validatePluginName(name string) error {
	if name == "" {
		return errors.New("plugin name must be set")
	}
	if !pluginNameRe.MatchString(name) {
		return fmt.Errorf("plugin name %q is invalid: use lowercase letters and digits, optionally joined by '-' or '_', starting with a letter", name)
	}
	return nil
}

func validateStageName(name string) error {
	if name == "" {
		return errors.New("stage name must not be empty")
	}
	if !stageNameRe.MatchString(name) {
		return fmt.Errorf("stage name %q is invalid: use UPPER_SNAKE_CASE such as MY_STAGE", name)
	}
	return nil
}

func validateModulePath(module string) error {
	if module == "" {
		return errors.New("module path must be set, for example github.com/your-org/your-plugin")
	}
	for _, elem := range strings.Split(module, "/") {
		if elem == "" || strings.HasPrefix(elem, ".") || strings.HasSuffix(elem, ".") {
			return fmt.Errorf("module path %q is invalid: a path element must not be empty or start or end with a dot", module)
		}
		if !modulePathElemRe.MatchString(elem) {
			return fmt.Errorf("module path %q is invalid: %q contains characters that are not allowed in a Go module path", module, elem)
		}
	}
	return nil
}

// stagePrefix turns a plugin name into the prefix of its default stage name,
// e.g. "my-plugin" -> "MY_PLUGIN".
func stagePrefix(pluginName string) string {
	return strings.ToUpper(strings.ReplaceAll(pluginName, "-", "_"))
}

// pascal converts UPPER_SNAKE_CASE to PascalCase: "MY_STAGE" -> "MyStage". The
// stage name is validated before, so none of its words is empty.
func pascal(stage string) string {
	words := strings.Split(strings.ToLower(stage), "_")
	for i, w := range words {
		words[i] = strings.ToUpper(w[:1]) + w[1:]
	}
	return strings.Join(words, "")
}

// stageFileName returns the file name that holds the executor of a stage.
//
// The stage name is always followed by "_stage" so the file name can never end
// in a suffix that Go gives a special meaning: a stage called RUN_TEST or
// RUN_LINUX would otherwise become run_test.go (only built for tests) or
// run_linux.go (only built on Linux). It also keeps a stage named PLUGIN from
// colliding with plugin.go.
func stageFileName(stage string) string {
	return strings.ToLower(stage) + "_stage.go"
}
