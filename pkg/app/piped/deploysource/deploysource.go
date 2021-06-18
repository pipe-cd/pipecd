// Copyright 2020 The PipeCD Authors.
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

package deploysource

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/pipe-cd/pipe/pkg/app/piped/sourcedecrypter"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/git"
	"github.com/pipe-cd/pipe/pkg/model"
)

type DeploySource struct {
	RepoDir                 string
	AppDir                  string
	RevisionName            string
	Revision                string
	DeploymentConfig        *config.Config
	GenericDeploymentConfig config.GenericDeploymentSpec
}

type Provider interface {
	Get(ctx context.Context, logWriter io.Writer) (*DeploySource, error)
	GetReadOnly(ctx context.Context, logWriter io.Writer) (*DeploySource, error)
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

type secretDecrypter interface {
	Decrypt(string) (string, error)
}

type provider struct {
	workingDir      string
	repoConfig      config.PipedRepository
	revisionName    string
	revision        string
	gitClient       gitClient
	appGitPath      *model.ApplicationGitPath
	secretDecrypter secretDecrypter

	done    bool
	source  *DeploySource
	err     error
	copyNum int
	mu      sync.Mutex
}

func NewProvider(
	workingDir string,
	repoConfig config.PipedRepository,
	revisionName string,
	revision string,
	gitClient gitClient,
	appGitPath *model.ApplicationGitPath,
	sd secretDecrypter,
) Provider {
	return &provider{
		workingDir:      workingDir,
		repoConfig:      repoConfig,
		revisionName:    revisionName,
		revision:        revision,
		gitClient:       gitClient,
		appGitPath:      appGitPath,
		secretDecrypter: sd,
	}
}

func (p *provider) Get(ctx context.Context, lw io.Writer) (*DeploySource, error) {
	writeLog(lw, "Preparing deploy source at %s commit (%s)", p.revisionName, p.revision)

	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.done {
		p.source, p.err = p.prepare(ctx, lw)
		p.done = true
	}

	if p.err != nil {
		return nil, p.err
	}

	ds, err := p.copy(lw)
	if err != nil {
		return nil, err
	}

	writeLog(lw, "Successfully prepared deploy source at %s commit (%s)", p.revisionName, p.revision)
	return ds, nil
}

func (p *provider) GetReadOnly(ctx context.Context, lw io.Writer) (*DeploySource, error) {
	writeLog(lw, "Preparing deploy source at %s commit (%s)", p.revisionName, p.revision)

	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.done {
		p.source, p.err = p.prepare(ctx, lw)
		p.done = true
	}

	if p.err != nil {
		return nil, p.err
	}

	writeLog(lw, "Successfully prepared deploy source at %s commit (%s)", p.revisionName, p.revision)
	return p.source, nil
}

func (p *provider) prepare(ctx context.Context, lw io.Writer) (*DeploySource, error) {
	// Ensure the existence of the working directory.
	if err := os.MkdirAll(p.workingDir, 0700); err != nil {
		writeLog(lw, "Unable to create the working directory to store deploy source (%v)", err)
		return nil, err
	}

	// Create a temporary directory for storing the source.
	dir, err := ioutil.TempDir(p.workingDir, "deploysource")
	if err != nil {
		writeLog(lw, "Unable to create a temp directory to store the deploy source (%v)", err)
		return nil, err
	}

	repoDir := filepath.Join(dir, "repo")
	appDir := filepath.Join(repoDir, p.appGitPath.Path)

	// Clone the specified revision of the repository.
	gitRepo, err := p.gitClient.Clone(ctx, p.repoConfig.RepoID, p.repoConfig.Remote, p.repoConfig.Branch, repoDir)
	if err != nil {
		writeLog(lw, "Unable to clone the branch %s of the repository %s (%v)", p.repoConfig.Branch, p.repoConfig.RepoID, err)
		return nil, err
	}
	if err := gitRepo.Checkout(ctx, p.revision); err != nil {
		writeLog(lw, "Unable to checkout the %s commit %s (%v)", p.revisionName, p.revision, err)
		return nil, err
	}
	writeLog(lw, "Successfully cloned the %s commit %s of the repository %s", p.revisionName, p.revision, p.repoConfig.RepoID)

	// Load the deployment configuration file.
	configFileRelativePath := p.appGitPath.GetDeploymentConfigFilePath()
	configFilePath := filepath.Join(repoDir, configFileRelativePath)
	cfg, err := config.LoadFromYAML(configFilePath)
	if err != nil {
		writeLog(lw, "Unable to load the deployment configuration file at %s (%v)", configFileRelativePath, err)
		return nil, err
	}

	gdc, ok := cfg.GetGenericDeployment()
	if !ok {
		writeLog(lw, "Invalid application kind %s", cfg.Kind)
		return nil, fmt.Errorf("unsupport application kind %s", cfg.Kind)
	}
	writeLog(lw, "Successfully loaded the deployment configuration file")

	// Decrypt the sealed secrets if needed.
	if len(gdc.SealedSecrets) > 0 && p.secretDecrypter != nil {
		if err := sourcedecrypter.DecryptSealedSecrets(appDir, gdc.SealedSecrets, p.secretDecrypter); err != nil {
			writeLog(lw, "Unable to decrypt the sealed secrets (%v)", err)
			return nil, err
		}
		writeLog(lw, "Successfully decrypted %d sealed secrets", len(gdc.SealedSecrets))
	}
	if gdc.Encryption != nil && p.secretDecrypter != nil {
		if err := sourcedecrypter.DecryptSecrets(appDir, *gdc.Encryption, p.secretDecrypter); err != nil {
			writeLog(lw, "Unable to decrypt the secrets (%v)", err)
			return nil, err
		}
		writeLog(lw, "Successfully decrypted secrets: %v", gdc.Encryption.DecryptionTargets)
	}

	return &DeploySource{
		RepoDir:                 repoDir,
		AppDir:                  appDir,
		RevisionName:            p.revisionName,
		Revision:                p.revision,
		DeploymentConfig:        cfg,
		GenericDeploymentConfig: gdc,
	}, nil
}

func (p *provider) copy(lw io.Writer) (*DeploySource, error) {
	p.copyNum++

	dest := fmt.Sprintf("%s-%d", p.source.RepoDir, p.copyNum)
	cmd := exec.Command("cp", "-rf", p.source.RepoDir, dest)
	out, err := cmd.CombinedOutput()
	if err != nil {
		writeLog(lw, "Unable to copy deploy source data (%v, %s)", err, string(out))
		return nil, err
	}

	return &DeploySource{
		RepoDir:                 dest,
		AppDir:                  filepath.Join(dest, p.appGitPath.Path),
		RevisionName:            p.revisionName,
		Revision:                p.revision,
		DeploymentConfig:        p.source.DeploymentConfig,
		GenericDeploymentConfig: p.source.GenericDeploymentConfig,
	}, nil
}

func writeLog(w io.Writer, format string, a ...interface{}) {
	io.WriteString(w, fmt.Sprintf(format, a...))
}
