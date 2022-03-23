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

// TerraformFileMapping is a schema for Terraform file.
type TerraformFileMapping struct {
	ModuleMappings []*ModuleMapping `hcl:"module,block"`
	Remain         hcl.Body         `hcl:",remain"`
}

// ModuleMapping is a schema for "module" block in Terraform file.
type ModuleMapping struct {
	Name    string   `hcl:"name,label"`
	Source  string   `hcl:"source"`
	Version string   `hcl:"version"`
	Remain  hcl.Body `hcl:",remain"`
}

// TerraformFile represents a Terraform file.
type TerraformFile struct {
	Modules []*Module
}

// Module represents a "module" block in Terraform file.
type Module struct {
	Name    string
	Source  string
	Version string
}

// LoadTerraformFiles loads terraform files from a given dir.
func LoadTerraformFiles(dir string) ([]*TerraformFile, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filenames := []string{}
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if ext := filepath.Ext(f.Name()); ext != ".tf" {
			continue
		}

		filenames = append(filenames, f.Name())
	}

	if len(filenames) == 0 {
		return nil, fmt.Errorf("couldn't find terraform module")
	}

	p := hclparse.NewParser()
	tfs := make([]*TerraformFile, 0, len(filenames))
	for _, fn := range filenames {
		f, diags := p.ParseHCLFile(filepath.Join(dir, fn))
		if diags.HasErrors() {
			return nil, diags
		}

		tfm := &TerraformFileMapping{}
		diags = gohcl.DecodeBody(f.Body, nil, tfm)
		if diags.HasErrors() {
			return nil, diags
		}

		tf := TerraformFile{
			Modules: make([]*Module, 0, len(tfm.ModuleMappings)),
		}
		for _, m := range tfm.ModuleMappings {
			tf.Modules = append(tf.Modules, &Module{
				Name:    m.Name,
				Source:  m.Source,
				Version: m.Version,
			})
		}

		tfs = append(tfs, &tf)
	}

	return tfs, nil
}

// FindArtifactVersions parses artifact versions from Terraform files.
// For Terraform, module version is an artifact version.
func FindArtifactVersions(tfs []*TerraformFile) ([]*model.ArtifactVersion, error) {
	var modules []*Module
	for _, tf := range tfs {
		modules = append(modules, tf.Modules...)
	}

	versions := make([]*model.ArtifactVersion, 0, len(modules))
	for _, m := range modules {
		versions = append(versions, &model.ArtifactVersion{
			Kind:    model.ArtifactVersion_TERRAFORM_MODULE,
			Version: m.Version,
			Name:    m.Name,
			Url:     m.Source,
		})
	}

	return versions, nil
}
