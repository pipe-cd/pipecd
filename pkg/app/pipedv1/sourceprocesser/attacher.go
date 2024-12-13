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

	config "github.com/pipe-cd/pipecd/pkg/configv1"
)

type attachmentProcessor struct {
	atc *config.Attachment
}

func NewAttachmentProcessor(atc *config.Attachment) *attachmentProcessor {
	return &attachmentProcessor{
		atc: atc,
	}
}

func (p *attachmentProcessor) BuildTemplateData(appDir string) (map[string]string, error) {
	if len(p.atc.Sources) == 0 {
		// Skip building no error.
		return nil, nil
	}

	src := make(map[string]string, len(p.atc.Sources))
	for k, v := range p.atc.Sources {
		srcPath := filepath.Join(appDir, v)
		buff, err := os.ReadFile(srcPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read data source to attach from file %s (%w)", v, err)
		}
		src[k] = string(buff)
	}
	return src, nil
}

func (p *attachmentProcessor) TemplateKey() string {
	return "attachment"
}

func (p *attachmentProcessor) TemplateSource(appDir string, data map[string]map[string]string) error {
	for _, t := range p.atc.Targets {
		targetPath := filepath.Join(appDir, t)
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
