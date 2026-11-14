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

// Package pluginscaffold generates the source of a new piped plugin.
//
// Generate renders every file in memory and Write is the only function that
// touches the file system, so the output can be tested without a disk.
package pluginscaffold

import (
	"fmt"
	"strings"
)

const (
	// defaultSDKVersion is the version of github.com/pipe-cd/piped-plugin-sdk-go
	// that generated plugins depend on. It has to match what the official
	// plugins (for example pkg/app/pipedv1/plugin/wait) pin.
	defaultSDKVersion = "v0.4.0"
	// defaultGoVersion is the "go" directive written to the generated go.mod.
	// It has to match the official plugins as well.
	defaultGoVersion = "1.26.2"
	// tutorialURL is the tutorial that builds a plugin step by step. The
	// generated README links to it.
	tutorialURL = "https://pipecd.dev/docs-v1.0.x/plugins/creating-a-plugin/"

	// pluginVersion is the version string the generated main.go passes to
	// sdk.NewPlugin. It is the same placeholder the official plugins use.
	pluginVersion = "0.0.1"
)

// Options describes the plugin to generate.
type Options struct {
	// Name is the name of the plugin, as it is registered in the piped config.
	// Required.
	Name string
	// Module is the Go module path of the generated plugin. Required.
	Module string
	// Stages are the stages the plugin defines, in UPPER_SNAKE_CASE. Defaults to
	// a single stage called <NAME>_RUN.
	Stages []string
}

// stage is one stage of the generated plugin with every derived name the
// templates need.
type stage struct {
	// Name is the stage name as piped sees it, e.g. MY_STAGE.
	Name string
	// Const is the name of the Go constant holding Name, e.g. stageMyStage.
	Const string
	// Method is the name of the executor method, e.g. executeMyStage.
	Method string
	// File is the name of the file holding the executor, e.g. my_stage_stage.go.
	File string
}

func newStage(name string) stage {
	return stage{
		Name:   name,
		Const:  "stage" + pascal(name),
		Method: "execute" + pascal(name),
		File:   stageFileName(name),
	}
}

// model is everything the templates are rendered with.
type model struct {
	Name          string
	Module        string
	SDKVersion    string
	GoVersion     string
	PluginVersion string
	TutorialURL   string
	Stages        []stage
}

// newModel validates opts, fills in the defaults and derives every name.
func newModel(opts Options) (*model, error) {
	if err := validatePluginName(opts.Name); err != nil {
		return nil, err
	}
	if err := validateModulePath(opts.Module); err != nil {
		return nil, err
	}

	names := opts.Stages
	if len(names) == 0 {
		names = []string{stagePrefix(opts.Name) + "_RUN"}
	}

	m := &model{
		Name:          opts.Name,
		Module:        opts.Module,
		SDKVersion:    defaultSDKVersion,
		GoVersion:     defaultGoVersion,
		PluginVersion: pluginVersion,
		TutorialURL:   tutorialURL,
		Stages:        make([]stage, 0, len(names)),
	}
	// "--stages A, B" arrives as "A" and " B", the spaces are not part of a name.
	for _, n := range names {
		n = strings.TrimSpace(n)
		if err := validateStageName(n); err != nil {
			return nil, err
		}
		m.Stages = append(m.Stages, newStage(n))
	}

	if err := checkStageCollisions(m.Stages); err != nil {
		return nil, err
	}
	return m, nil
}

// checkStageCollisions rejects a stage that is listed twice, and stages whose
// Go identifiers would clash even though the names differ: A1 and A_1 both
// become stageA1. The file name follows the stage name, so it cannot clash alone.
func checkStageCollisions(stages []stage) error {
	seen := make(map[string]string, len(stages))
	for _, s := range stages {
		other, ok := seen[s.Const]
		if !ok {
			seen[s.Const] = s.Name
			continue
		}
		if other == s.Name {
			return fmt.Errorf("stage %q is listed more than once", s.Name)
		}
		return fmt.Errorf("stages %q and %q would generate the same Go identifier, rename one of them", other, s.Name)
	}
	return nil
}
