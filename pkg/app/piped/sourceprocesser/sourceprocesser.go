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
)

type SourceTemplateProcessor interface {
	// BuildTemplateData returns the data that will be used to template target files
	BuildTemplateData(appDir string) (map[string]string, error)
	// TemplateSource performs the templating prepared data to the source files
	TemplateSource(appDir string, data map[string](map[string]string)) error
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
		return fmt.Errorf("no template processor was specified")
	}

	// Build the initial data for the template.
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

	for _, tp := range p.templateProcessors {
		// Rebuid the data for the template to get up to date
		// data after previous template processing.
		pdata, err := tp.BuildTemplateData(p.appDir)
		if err != nil {
			return err
		}
		if len(pdata) == 0 {
			continue
		}
		data[tp.TemplateKey()] = pdata

		// Perform the templating.
		if err := tp.TemplateSource(p.appDir, data); err != nil {
			return err
		}
	}

	return nil
}
