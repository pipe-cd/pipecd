package terraform

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"

	"github.com/pipe-cd/pipecd/pkg/git"
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
	Version string   `hcl:"version"`
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
	IsLocal bool
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
				IsLocal: isLocalModule(m.Source),
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

func isLocalModule(moduleSrc string) bool {
	return strings.HasPrefix(moduleSrc, "./") || strings.HasPrefix(moduleSrc, "../")
}

// LocalModuleSourceConverter is a converter from local module source to its URL.
type LocalModuleSourceConverter struct {
	// GitURL is a URL for git repository URL
	GitURL string
	// Branch is a git branch for the current terraform file repo.
	Branch string
	// RepoDir is a dir where the terraform file repo located.
	RepoDir string
	// AppDir is a dir where the terraform module located.
	AppDir string
}

func NewLocalModuleSourceConverter(gitURL, branch, repoDir, AppDir string) *LocalModuleSourceConverter {
	return &LocalModuleSourceConverter{
		GitURL:  gitURL,
		Branch:  branch,
		RepoDir: repoDir,
		AppDir:  AppDir,
	}
}

// MakeURL make a URL from a local module source.
func (l *LocalModuleSourceConverter) MakeURL(moduleSrc string) (string, error) {
	// resolve path for local module
	// moduleSrc is a relative path like "./" or "../" and so on.
	moduleDir := filepath.Join(l.AppDir, moduleSrc)
	dirFromRepo, err := filepath.Rel(l.RepoDir, moduleDir)
	if err != nil {
		return "", err
	}

	if strings.HasSuffix(dirFromRepo, "..") {
		return "", fmt.Errorf("can't resolve relative path on git repo")
	}

	return git.MakeDirURL(l.GitURL, dirFromRepo, l.Branch)
}
