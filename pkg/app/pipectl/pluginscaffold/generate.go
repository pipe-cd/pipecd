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
	"embed"
	"fmt"
	"go/format"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

//go:embed templates
var templatesFS embed.FS

// templates holds the templates of a stage plugin.
var templates = template.Must(template.New("").
	Funcs(template.FuncMap{"quote": strconv.Quote}).
	Option("missingkey=error").
	ParseFS(templatesFS, "templates/stage/*.tmpl"))

// File is one generated file.
type File struct {
	// Path is relative to the plugin directory and always uses forward slashes.
	Path string
	// Content is the rendered content. Go files are already gofmt-ed.
	Content []byte
}

// pluginFiles lists the templates that are rendered once, with the file each
// one produces. The stage executors are rendered once per stage, see Generate.
var pluginFiles = []struct {
	tmpl string
	out  string
}{
	{tmpl: "go.mod.tmpl", out: "go.mod"},
	{tmpl: "main.go.tmpl", out: "main.go"},
	{tmpl: "plugin.go.tmpl", out: "plugin.go"},
	{tmpl: "plugin_test.go.tmpl", out: "plugin_test.go"},
	{tmpl: "README.md.tmpl", out: "README.md"},
}

// stageData is what the stage template is rendered with.
type stageData struct {
	*model
	Stage stage
}

// Generate renders the plugin described by opts. Nothing is written to disk.
// The returned files are sorted by path.
func Generate(opts Options) ([]File, error) {
	m, err := newModel(opts)
	if err != nil {
		return nil, err
	}

	files := make([]File, 0, len(pluginFiles)+len(m.Stages))
	for _, spec := range pluginFiles {
		f, err := render(spec.tmpl, spec.out, m)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	for _, s := range m.Stages {
		f, err := render("stage.go.tmpl", s.File, stageData{model: m, Stage: s})
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files, nil
}

func render(tmplName, out string, data any) (File, error) {
	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, tmplName, data); err != nil {
		return File{}, fmt.Errorf("failed to render %s: %w", out, err)
	}

	content := buf.Bytes()
	if strings.HasSuffix(out, ".go") {
		formatted, err := format.Source(content)
		if err != nil {
			// A broken template is a bug in the scaffolder, not in the input.
			return File{}, fmt.Errorf("generated %s is not valid Go (this is a bug in the templates): %w", out, err)
		}
		content = formatted
	}
	return File{Path: out, Content: content}, nil
}
