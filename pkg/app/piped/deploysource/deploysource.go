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

package deploysource

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/pipe-cd/pipecd/pkg/app/piped/sourceprocesser"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type EmbeddedType string

const (
	EmbeddedTypeSecretOnly     EmbeddedType = "EMBEDDED_TYPE_SECRET_ONLY"
	EmbeddedTypeAttachmentOnly EmbeddedType = "EMBEDDED_TYPE_ATTACHMENT_ONLY"
	EmbeddedTypeBoth           EmbeddedType = "EMBEDDED_TYPE_BOTH"
	EmbeddedTypeNone           EmbeddedType = "EMBEDDED_TYPE_NONE"
)

type embedded struct {
	secret     bool
	attachment bool
}

type DeploySource struct {
	RepoDir                  string
	AppDir                   string
	Revision                 string
	ApplicationConfig        *config.Config
	GenericApplicationConfig config.GenericApplicationSpec
}

type Provider interface {
	Revision() string
	Get(ctx context.Context, logWriter io.Writer) (*DeploySource, error)
	GetReadOnly(ctx context.Context, logWriter io.Writer) (*DeploySource, error)
}

type secretDecrypter interface {
	Decrypt(string) (string, error)
}

type provider struct {
	workingDir      string
	cloner          SourceCloner
	revisionName    string
	revision        string
	appGitPath      model.ApplicationGitPath
	secretDecrypter secretDecrypter

	done    bool
	source  *DeploySource
	err     error
	copyNum int
	mu      sync.Mutex
}

func NewProvider(
	workingDir string,
	cloner SourceCloner,
	appGitPath model.ApplicationGitPath,
	sd secretDecrypter,
) Provider {

	return &provider{
		workingDir:      workingDir,
		cloner:          cloner,
		revisionName:    cloner.RevisionName(),
		revision:        cloner.Revision(),
		appGitPath:      appGitPath,
		secretDecrypter: sd,
	}
}

func (p *provider) Revision() string {
	return p.cloner.Revision()
}

func (p *provider) Get(ctx context.Context, lw io.Writer) (*DeploySource, error) {
	fmt.Fprintf(lw, "Preparing deploy source at %s commit (%s)\n", p.revisionName, p.revision)

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

	fmt.Fprintf(lw, "Successfully prepared deploy source at %s commit (%s)\n", p.revisionName, p.revision)
	return ds, nil
}

func (p *provider) GetReadOnly(ctx context.Context, lw io.Writer) (*DeploySource, error) {
	fmt.Fprintf(lw, "Preparing deploy source at %s commit (%s)\n", p.revisionName, p.revision)

	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.done {
		p.source, p.err = p.prepare(ctx, lw)
		p.done = true
	}

	if p.err != nil {
		return nil, p.err
	}

	fmt.Fprintf(lw, "Successfully prepared deploy source at %s commit (%s)\n", p.revisionName, p.revision)
	return p.source, nil
}

func (p *provider) prepare(ctx context.Context, lw io.Writer) (*DeploySource, error) {
	// Ensure the existence of the working directory.
	if err := os.MkdirAll(p.workingDir, 0700); err != nil {
		fmt.Fprintf(lw, "Unable to create the working directory to store deploy source (%v)\n", err)
		return nil, err
	}

	// Create a temporary directory for storing the source.
	dir, err := os.MkdirTemp(p.workingDir, "deploysource")
	if err != nil {
		fmt.Fprintf(lw, "Unable to create a temp directory to store the deploy source (%v)\n", err)
		return nil, err
	}

	repoDir := filepath.Join(dir, "repo")
	appDir := filepath.Join(repoDir, p.appGitPath.Path)

	// Clone the specified revision of the repository.
	if err := p.cloner.Clone(ctx, repoDir); err != nil {
		fmt.Fprintf(lw, "Unable to clone the %s commit (%v)\n", p.revisionName, err)
		return nil, err
	}
	fmt.Fprintf(lw, "Successfully cloned the %s commit\n", p.revisionName)

	// Load the application configuration file.
	var (
		cfgFileRelPath = p.appGitPath.GetApplicationConfigFilePath()
		cfgFileAbsPath = filepath.Join(repoDir, cfgFileRelPath)
	)
	cfg, err := config.LoadFromYAML(cfgFileAbsPath)
	if err != nil {
		fmt.Fprintf(lw, "Unable to load the application configuration file at %s (%v)\n", cfgFileRelPath, err)

		if os.IsNotExist(err) {
			return nil, fmt.Errorf("application config file %s was not found", cfgFileRelPath)
		}
		return nil, err
	}

	gac, ok := cfg.GetGenericApplication()
	if !ok {
		fmt.Fprintf(lw, "Invalid application kind %s\n", cfg.Kind)
		return nil, fmt.Errorf("unsupport application kind %s", cfg.Kind)
	}
	fmt.Fprintln(lw, "Successfully loaded the application configuration file")

	var (
		embedded     embedded
		embeddedType = EmbeddedTypeNone
	)
	// Decrypt the sealed secrets if needed.
	if gac.Encryption != nil && p.secretDecrypter != nil && len(gac.Encryption.DecryptionTargets) > 0 {
		embedded.secret = true
	}

	if gac.Attachment != nil && len(gac.Attachment.Targets) > 0 {
		embedded.attachment = true
	}

	if embedded.secret && embedded.attachment {
		embeddedType = EmbeddedTypeBoth
	} else if embedded.secret {
		embeddedType = EmbeddedTypeSecretOnly
	} else if embedded.attachment {
		embeddedType = EmbeddedTypeAttachmentOnly
	}

	switch embeddedType {
	case EmbeddedTypeBoth:
		if err := sourceprocesser.EmbedCombination(appDir, *gac.Encryption, p.secretDecrypter, *gac.Attachment); err != nil {
			fmt.Fprintf(lw, "Unable to embed the secret and attachment data (%v)\n", err)
			return nil, err
		}
	case EmbeddedTypeSecretOnly:
		if err := sourceprocesser.DecryptSecrets(appDir, *gac.Encryption, p.secretDecrypter); err != nil {
			fmt.Fprintf(lw, "Unable to decrypt the secrets (%v)\n", err)
			return nil, err
		}
		fmt.Fprintf(lw, "Successfully decrypted secrets: %v\n", gac.Encryption.DecryptionTargets)
	case EmbeddedTypeAttachmentOnly:
		if err := sourceprocesser.AttachData(appDir, *gac.Attachment); err != nil {
			fmt.Fprintf(lw, "Unable to attach the data (%v)\n", err)
			return nil, err
		}
		fmt.Fprintf(lw, "Successfully attached data: %v\n", gac.Attachment.Targets)
	}

	return &DeploySource{
		RepoDir:                  repoDir,
		AppDir:                   appDir,
		Revision:                 p.revision,
		ApplicationConfig:        cfg,
		GenericApplicationConfig: gac,
	}, nil
}

func (p *provider) copy(lw io.Writer) (*DeploySource, error) {
	p.copyNum++

	dest := fmt.Sprintf("%s-%d", p.source.RepoDir, p.copyNum)
	cmd := exec.Command("cp", "-rf", p.source.RepoDir, dest)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(lw, "Unable to copy deploy source data (%v, %s)\n", err, string(out))
		return nil, err
	}

	return &DeploySource{
		RepoDir:                  dest,
		AppDir:                   filepath.Join(dest, p.appGitPath.Path),
		Revision:                 p.revision,
		ApplicationConfig:        p.source.ApplicationConfig,
		GenericApplicationConfig: p.source.GenericApplicationConfig,
	}, nil
}
