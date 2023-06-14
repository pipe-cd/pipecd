// Copyright 2023 The PipeCD Authors.
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

package terraform

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"

	"github.com/pipe-cd/pipecd/pkg/model"
)

// FileMapping is a schema for Terraform file.
type FileMapping struct {
	ModuleMappings []*ModuleMapping `hcl:"module,block"`
	Remain         hcl.Body         `hcl:",remain"`
}

// ModuleMapping is a schema for "module" block in Terraform file.
type ModuleMapping struct {
	Name    string   `hcl:"name,label"`
	Source  string   `hcl:"source"`
	Version string   `hcl:"version,optional"`
	Remain  hcl.Body `hcl:",remain"`
}

// File represents a Terraform file.
type File struct {
	Modules []*Module
}

// Module represents a "module" block in Terraform file.
type Module struct {
	Name    string
	Source  string
	Version string
}

const tfFileExtension = ".tf"

// LoadTerraformFiles loads terraform files from a given dir.
func LoadTerraformFiles(dir string) ([]File, error) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filepaths := make([]string, 0)
	for _, f := range fileInfos {
		if f.IsDir() {
			continue
		}

		if ext := filepath.Ext(f.Name()); ext != tfFileExtension {
			continue
		}

		filepaths = append(filepaths, filepath.Join(dir, f.Name()))
	}

	if len(filepaths) == 0 {
		return nil, fmt.Errorf("couldn't find terraform module")
	}

	p := hclparse.NewParser()
	tfs := make([]File, 0, len(filepaths))
	for _, fp := range filepaths {
		f, diags := p.ParseHCLFile(fp)
		if diags.HasErrors() {
			return nil, diags
		}

		fm := &FileMapping{}
		diags = gohcl.DecodeBody(f.Body, nil, fm)
		if diags.HasErrors() {
			return nil, diags
		}

		tf := File{
			Modules: make([]*Module, 0, len(fm.ModuleMappings)),
		}
		for _, m := range fm.ModuleMappings {
			tf.Modules = append(tf.Modules, &Module{
				Name:    m.Name,
				Source:  m.Source,
				Version: m.Version,
			})
		}

		tfs = append(tfs, tf)
	}

	return tfs, nil
}

// FindArtifactVersions parses artifact versions from Terraform files.
// For Terraform, module version is an artifact version.
func FindArtifactVersions(tfs []File) ([]*model.ArtifactVersion, error) {
	versions := make([]*model.ArtifactVersion, 0)
	for _, tf := range tfs {
		for _, m := range tf.Modules {
			versions = append(versions, &model.ArtifactVersion{
				Kind:    model.ArtifactVersion_TERRAFORM_MODULE,
				Version: m.Version,
				Name:    m.Name,
				Url:     m.Source,
			})
		}
	}

	return versions, nil
}
