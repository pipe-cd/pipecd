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
	"github.com/pipe-cd/pipecd/pkg/config"
)

func PrepareAttachmentData(appDir string, atc config.Attachment) (map[string](map[string]string), error) {
	if len(atc.Targets) == 0 {
		return make(map[string]map[string]string), nil
	}
	if len(atc.Sources) == 0 {
		return make(map[string]map[string]string), nil
	}

	src := make(map[string]string, len(atc.Sources))
	for k, v := range atc.Sources {
		srcPath := filepath.Join(appDir, v)
		buff, err := os.ReadFile(srcPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read data source to attach from file %s (%w)", v, err)
		}
		src[k] = string(buff)
	}
	data := map[string](map[string]string){
		"attachment": src,
	}

	return data, nil
}

func EmbedAttach(appDir string, targets []string, data map[string](map[string]string)) error {
	for _, t := range targets {
		targetPath := filepath.Join(appDir, t)
		fileName := filepath.Base(targetPath)
		tmpl := template.
			New(fileName).
			Funcs(sprig.TxtFuncMap()).
			Option("missingkey=error")
		tmpl, err := tmpl.ParseFiles(targetPath)
		if err != nil {
			return fmt.Errorf("failed to parse attaching target %s (%w)", t, err)
		}

		f, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("failed to open attaching target %s (%w)", t, err)
		}

		if err := tmpl.Execute(f, data); err != nil {
			f.Close()
			return fmt.Errorf("failed to render attaching target %s (%w)", t, err)
		}

		if err := f.Close(); err != nil {
			return fmt.Errorf("failed to close attached target %s (%w)", t, err)
		}
	}

	return nil
}

func AttachData(appDir string, atc config.Attachment) error {
	data, err := PrepareAttachmentData(appDir, atc)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	if err := EmbedAttach(appDir, atc.Targets, data); err != nil {
		return err
	}

	return nil
}
