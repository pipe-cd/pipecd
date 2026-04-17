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

// plugin-scaffold generates scaffolding for a new PipeCD deployment plugin.
//
// Usage:
//
//	go run ./hack/plugin-scaffold \
//	  --name myplatform \
//	  --module github.com/my-org/my-plugin \
//	  --stages "MY_DEPLOY:Deploy resources,MY_PROMOTE:Promote to production" \
//	  --rollback MY_ROLLBACK \
//	  --output ./output
package main

import (
	"bytes"
	"embed"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

//go:embed templates
var tmplFS embed.FS

// Stage holds parsed information about a single stage.
type Stage struct {
	// Name is the constant value, e.g. "MY_DEPLOY"
	Name string
	// Description is the human-readable description
	Description string
	// IsRollback marks this stage as the rollback stage
	IsRollback bool
}

// PluginData holds all data needed to render plugin templates.
type PluginData struct {
	// PluginName is the raw plugin name, e.g. "myplatform"
	PluginName string
	// PluginTitle is the title-cased plugin name, e.g. "Myplatform"
	PluginTitle string
	// Module is the Go module path, e.g. "github.com/my-org/my-plugin"
	Module string
	// Stages is the full list of stages including rollback
	Stages []Stage
	// DeployStages is stages excluding rollback
	DeployStages []Stage
	// RollbackStage is the rollback stage (may be zero value if none)
	RollbackStage *Stage
	// HasRollback is true if a rollback stage was specified
	HasRollback bool
	// HasLivestate is true if --livestate flag was set
	HasLivestate bool
	// HasPlanPreview is true if --planpreview flag was set
	HasPlanPreview bool
}

func main() {
	var (
		name        = flag.String("name", "", "Plugin name (lowercase, e.g. myplatform)")
		module      = flag.String("module", "", "Go module path (e.g. github.com/my-org/my-plugin)")
		stagesRaw   = flag.String("stages", "", "Comma-separated stages: NAME or NAME:Description")
		rollback    = flag.String("rollback", "", "Rollback stage name (optional, e.g. MY_ROLLBACK)")
		livestate   = flag.Bool("livestate", false, "Generate livestate/plugin.go stub")
		planpreview = flag.Bool("planpreview", false, "Generate planpreview/plugin.go stub")
		output      = flag.String("output", ".", "Output directory")
		force       = flag.Bool("force", false, "Overwrite existing output directory")
	)
	flag.Parse()

	if err := run(*name, *module, *stagesRaw, *rollback, *output, *livestate, *planpreview, *force); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(name, module, stagesRaw, rollbackName, output string, hasLivestate, hasPlanPreview, force bool) error {
	if name == "" {
		return errors.New("--name is required")
	}
	if module == "" {
		return errors.New("--module is required")
	}
	if stagesRaw == "" {
		return errors.New("--stages is required")
	}

	deployStages, err := parseStages(stagesRaw)
	if err != nil {
		return fmt.Errorf("parse stages: %w", err)
	}

	data := PluginData{
		PluginName:     name,
		PluginTitle:    titleCase(name),
		Module:         module,
		DeployStages:   deployStages,
		Stages:         deployStages,
		HasLivestate:   hasLivestate,
		HasPlanPreview: hasPlanPreview,
	}

	if rollbackName != "" {
		rb := Stage{
			Name:        rollbackName,
			Description: "Rollback to the previous version",
			IsRollback:  true,
		}
		data.RollbackStage = &rb
		data.HasRollback = true
		data.Stages = append(deployStages, rb)
	}

	root := filepath.Join(output, name)
	if _, err := os.Stat(root); err == nil {
		if !force {
			return fmt.Errorf("output directory %q already exists; use --force to overwrite", root)
		}
	}
	return scaffold(root, data)
}

func scaffold(root string, data PluginData) error {
	staticFiles := map[string]string{
		"main.go":                 "templates/main.go.tmpl",
		"go.mod":                  "templates/go.mod.tmpl",
		"config/plugin.go":        "templates/config/plugin.go.tmpl",
		"config/application.go":   "templates/config/application.go.tmpl",
		"config/deploy_target.go": "templates/config/deploy_target.go.tmpl",
		"deployment/plugin.go":    "templates/deployment/plugin.go.tmpl",
		"deployment/pipeline.go":  "templates/deployment/pipeline.go.tmpl",
	}

	for outPath, tmplPath := range staticFiles {
		if err := renderFile(filepath.Join(root, outPath), tmplPath, data); err != nil {
			return fmt.Errorf("render %s: %w", outPath, err)
		}
	}

	for _, stage := range data.DeployStages {
		stageData := struct {
			PluginData
			Stage Stage
		}{data, stage}
		outPath := filepath.Join(root, "deployment", stageNameToSnake(stage.Name)+".go")
		if err := renderFile(outPath, "templates/deployment/stage.go.tmpl", stageData); err != nil {
			return fmt.Errorf("render stage file %s: %w", outPath, err)
		}
	}

	if data.HasLivestate {
		outPath := filepath.Join(root, "livestate", "plugin.go")
		if err := renderFile(outPath, "templates/livestate/plugin.go.tmpl", data); err != nil {
			return fmt.Errorf("render livestate/plugin.go: %w", err)
		}
	}

	if data.HasPlanPreview {
		outPath := filepath.Join(root, "planpreview", "plugin.go")
		if err := renderFile(outPath, "templates/planpreview/plugin.go.tmpl", data); err != nil {
			return fmt.Errorf("render planpreview/plugin.go: %w", err)
		}
	}

	fmt.Printf("Plugin scaffolding generated at: %s\n", root)
	fmt.Println("\nFiles created:")
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			rel, _ := filepath.Rel(root, path)
			fmt.Printf("  %s\n", rel)
		}
		return nil
	})
}

func renderFile(outPath, tmplPath string, data any) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}
	tmplContent, err := tmplFS.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("read template %s: %w", tmplPath, err)
	}
	funcMap := template.FuncMap{
		"titleStage": stageNameToTitle,
		"snakeStage": stageNameToSnake,
		"lower":      strings.ToLower,
	}
	t, err := template.New(tmplPath).Funcs(funcMap).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("parse template %s: %w", tmplPath, err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Errorf("execute template %s: %w", tmplPath, err)
	}
	return os.WriteFile(outPath, buf.Bytes(), 0644)
}

// parseStages parses "NAME:Desc,NAME2:Desc2" into Stage slices.
func parseStages(raw string) ([]Stage, error) {
	parts := strings.Split(raw, ",")
	stages := make([]Stage, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		before, after, ok := strings.Cut(p, ":")
		var name, desc string
		if ok {
			name = strings.TrimSpace(before)
			desc = strings.TrimSpace(after)
		} else {
			name = p
			desc = fmt.Sprintf("Execute %s stage", strings.ToLower(strings.ReplaceAll(p, "_", " ")))
		}
		if name == "" {
			return nil, fmt.Errorf("empty stage name in %q", p)
		}
		stages = append(stages, Stage{Name: name, Description: desc})
	}
	if len(stages) == 0 {
		return nil, errors.New("at least one stage is required")
	}
	return stages, nil
}

// titleCase converts "myplatform" or "my-platform" to "Myplatform" / "MyPlatform".
func titleCase(s string) string {
	var b strings.Builder
	upper := true
	for _, r := range s {
		if r == '-' || r == '_' {
			upper = true
			continue
		}
		if upper {
			b.WriteRune(unicode.ToUpper(r))
			upper = false
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// stageNameToTitle converts "MY_SYNC" to "MySync".
func stageNameToTitle(s string) string {
	parts := strings.Split(s, "_")
	var b strings.Builder
	for _, p := range parts {
		if len(p) == 0 {
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]))
		b.WriteString(strings.ToLower(p[1:]))
	}
	return b.String()
}

// stageNameToSnake converts "MY_SYNC" to "my_sync".
func stageNameToSnake(s string) string {
	return strings.ToLower(s)
}
