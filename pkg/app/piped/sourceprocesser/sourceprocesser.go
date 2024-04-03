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

package sourceprocesser

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

type SourceTemplateProcessor interface {
	// BuildTemplateData returns the data that will be used to template target files
	BuildTemplateData(appDir string) (map[string]string, error)
	// TargetFilePaths returns the paths of target files that will be templated
	TargetFilePaths() []string
	// TemplateKey returns the key that will be used to store the data in the template
	TemplateKey() string
}

type SourceProcessor interface {
	Process() error
}

type processor struct {
	appDir             string
	templateProcessors []SourceTemplateProcessor
}

func NewSourceProcessor(appDir string, templateProcessors ...SourceTemplateProcessor) SourceProcessor {
	return &processor{
		appDir:             appDir,
		templateProcessors: templateProcessors,
	}
}

func (p *processor) Process() error {
	if len(p.templateProcessors) == 0 {
		return nil
	}

	var targets []string
	for _, tp := range p.templateProcessors {
		targets = append(targets, tp.TargetFilePaths()...)
	}
	if len(targets) == 0 {
		return fmt.Errorf("no target file path was specified")
	}

	data := make(map[string](map[string]string))
	for _, tp := range p.templateProcessors {
		pdata, err := tp.BuildTemplateData(p.appDir)
		if err != nil {
			return err
		}
		if len(pdata) == 0 {
			continue
		}
		data[tp.TemplateKey()] = pdata
	}

	for _, t := range targets {
		targetPath := filepath.Join(p.appDir, t)
		fileName := filepath.Base(targetPath)
		tmpl := template.New(fileName).Funcs(sprig.TxtFuncMap()).Option("missingkey=error")
		tmpl, err := tmpl.ParseFiles(targetPath)
		if err != nil {
			return fmt.Errorf("failed to parse target file %s (%w)", t, err)
		}

		f, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("failed to open target file %s (%w)", t, err)
		}

		if err := tmpl.Execute(f, data); err != nil {
			f.Close()
			return fmt.Errorf("failed to render target file %s (%w)", t, err)
		}

		if err := f.Close(); err != nil {
			return fmt.Errorf("failed to close target file %s (%w)", t, err)
		}
	}

	return nil
}
