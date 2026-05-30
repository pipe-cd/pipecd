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
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

// DefaultSDKVersion is the piped-plugin-sdk-go version written to generated go.mod files.
// Bump it when official plugins (for example wait) upgrade their SDK dependency.
const DefaultSDKVersion = "v0.3.0"

// Kind is the scaffold template set to use.
type Kind string

const (
	KindStage      Kind = "stage"
	KindDeployment Kind = "deployment"
)

// Options configures plugin scaffolding.
type Options struct {
	OutputDir  string
	PluginName string
	ModulePath string
	Kind       Kind
	Stages     []string
	DryRun     bool
	Force      bool
}

// File describes a generated file.
type File struct {
	Path    string
	Content []byte
}

// Generate renders plugin scaffold files without writing to disk.
func Generate(opts Options) ([]File, error) {
	if err := ValidatePluginName(opts.PluginName); err != nil {
		return nil, err
	}
	if err := ValidateStageNames(opts.Stages); err != nil {
		return nil, err
	}
	if opts.ModulePath == "" {
		opts.ModulePath = DefaultModulePath(opts.PluginName)
	}
	switch opts.Kind {
	case KindStage, KindDeployment:
	default:
		return nil, fmt.Errorf("kind must be %q or %q", KindStage, KindDeployment)
	}

	data := templateData{
		PluginName: opts.PluginName,
		TypePrefix: TypePrefix(opts.PluginName),
		ModulePath: opts.ModulePath,
		SDKVersion: DefaultSDKVersion,
		Stages:     append([]string(nil), opts.Stages...),
	}
	if rollback := FindRollbackStage(opts.Stages); rollback != "" {
		data.RollbackStage = stageEntry{
			Name:        rollback,
			FileBase:    StageFileBase(rollback),
			FuncName:    StageFuncName(rollback),
			ConstSuffix: StageConstSuffix(rollback),
		}
	}
	for _, stage := range opts.Stages {
		data.StageEntries = append(data.StageEntries, stageEntry{
			Name:        stage,
			FileBase:    StageFileBase(stage),
			FuncName:    StageFuncName(stage),
			ConstSuffix: StageConstSuffix(stage),
		})
	}

	templateSet := string(opts.Kind)
	templateRoot := filepath.Join("templates", templateSet)

	var files []File
	err := fs.WalkDir(templateFS, templateRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || !strings.HasSuffix(path, ".tmpl") {
			return nil
		}
		if filepath.Base(path) == "stage.go.tmpl" {
			return nil
		}
		rel, err := filepath.Rel(templateRoot, path)
		if err != nil {
			return err
		}
		outName := strings.TrimSuffix(rel, ".tmpl")
		content, err := renderTemplate(path, data)
		if err != nil {
			return err
		}
		files = append(files, File{Path: outName, Content: content})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk templates: %w", err)
	}

	if opts.Kind == KindDeployment {
		stageTmpl, err := templateFS.ReadFile(filepath.Join(templateRoot, "stage.go.tmpl"))
		if err != nil {
			return nil, fmt.Errorf("read stage template: %w", err)
		}
		tmpl, err := template.New("stage").Parse(string(stageTmpl))
		if err != nil {
			return nil, fmt.Errorf("parse stage template: %w", err)
		}
		for _, stage := range data.StageEntries {
			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, data.withStage(stage)); err != nil {
				return nil, fmt.Errorf("execute stage template for %s: %w", stage.Name, err)
			}
			files = append(files, File{
				Path:    filepath.Join("deployment", stage.FileBase+".go"),
				Content: buf.Bytes(),
			})
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})
	return files, nil
}

// Write writes generated files to opts.OutputDir.
func Write(opts Options, files []File) error {
	if _, err := os.Stat(opts.OutputDir); err == nil {
		if !opts.Force {
			return fmt.Errorf("output directory %q already exists (use --force to overwrite)", opts.OutputDir)
		}
		if err := os.RemoveAll(opts.OutputDir); err != nil {
			return fmt.Errorf("remove existing output directory: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat output directory: %w", err)
	}
	if err := os.MkdirAll(opts.OutputDir, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	for _, f := range files {
		path := filepath.Join(opts.OutputDir, filepath.FromSlash(f.Path))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, f.Content, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", f.Path, err)
		}
	}
	return nil
}

type templateData struct {
	PluginName    string
	TypePrefix    string
	ModulePath    string
	SDKVersion    string
	Stages        []string
	StageEntries  []stageEntry
	RollbackStage stageEntry
	CurrentStage  stageEntry
}

type stageEntry struct {
	Name        string
	FileBase    string
	FuncName    string
	ConstSuffix string
}

func (d templateData) withStage(s stageEntry) templateData {
	d.CurrentStage = s
	return d
}

func renderTemplate(path string, data templateData) ([]byte, error) {
	raw, err := templateFS.ReadFile(path)
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New(filepath.Base(path)).Parse(string(raw))
	if err != nil {
		return nil, fmt.Errorf("parse template %s: %w", path, err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("execute template %s: %w", path, err)
	}
	return buf.Bytes(), nil
}
