// Copyright 2022 The PipeCD Authors.
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
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/piped/livestatestore/terraform"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/terraform"
	"github.com/pipe-cd/pipecd/pkg/app/piped/sourcedecrypter"
	"github.com/pipe-cd/pipecd/pkg/app/piped/toolregistry"
	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type applicationLister interface {
	ListByPlatformProvider(name string) []*model.Application
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

type secretDecrypter interface {
	Decrypt(string) (string, error)
}

type reporter interface {
	ReportApplicationSyncState(ctx context.Context, appID string, state model.ApplicationSyncState) error
}

type Detector interface {
	Run(ctx context.Context) error
	ProviderName() string
}

type detector struct {
	provider          config.PipedPlatformProvider
	appLister         applicationLister
	gitClient         gitClient
	stateGetter       terraform.Getter
	reporter          reporter
	appManifestsCache cache.Cache
	interval          time.Duration
	config            *config.PipedSpec
	secretDecrypter   secretDecrypter
	logger            *zap.Logger

	gitRepos   map[string]git.Repo
	syncStates map[string]model.ApplicationSyncState
}

func NewDetector(
	cp config.PipedPlatformProvider,
	appLister applicationLister,
	gitClient gitClient,
	stateGetter terraform.Getter,
	reporter reporter,
	appManifestsCache cache.Cache,
	cfg *config.PipedSpec,
	sd secretDecrypter,
	logger *zap.Logger,
) Detector {

	logger = logger.Named("terraform-detector").With(
		zap.String("cloud-provider", cp.Name),
	)
	return &detector{
		provider:          cp,
		appLister:         appLister,
		gitClient:         gitClient,
		stateGetter:       stateGetter,
		reporter:          reporter,
		appManifestsCache: appManifestsCache,
		interval:          time.Minute,
		config:            cfg,
		secretDecrypter:   sd,
		gitRepos:          make(map[string]git.Repo),
		syncStates:        make(map[string]model.ApplicationSyncState),
		logger:            logger,
	}
}

func (d *detector) Run(ctx context.Context) error {
	d.logger.Info("start running drift detector for terraform applications")

	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			d.check(ctx)

		case <-ctx.Done():
			d.logger.Info("drift detector for terraform applications has been stopped")
			return nil
		}
	}
}

func (d *detector) check(ctx context.Context) error {
	appsByRepo := d.listGroupedApplication()

	for repoID, apps := range appsByRepo {
		gitRepo, ok := d.gitRepos[repoID]
		if !ok {
			// Clone repository for the first time.
			gr, err := d.cloneGitRepository(ctx, repoID)
			if err != nil {
				d.logger.Error("failed to clone git repository",
					zap.String("repo-id", repoID),
					zap.Error(err),
				)
				continue
			}
			gitRepo = gr
			d.gitRepos[repoID] = gitRepo
		}

		// Fetch the latest commit to compare the states.
		branch := gitRepo.GetClonedBranch()
		if err := gitRepo.Pull(ctx, branch); err != nil {
			d.logger.Error("failed to pull repository branch",
				zap.String("repo-id", repoID),
				zap.Error(err),
			)
			continue
		}

		// Get the head commit of the repository.
		headCommit, err := gitRepo.GetLatestCommit(ctx)
		if err != nil {
			d.logger.Error("failed to get head commit hash",
				zap.String("repo-id", repoID),
				zap.Error(err),
			)
			continue
		}

		// Start checking all applications in this repository.
		for _, app := range apps {
			if err := d.checkApplication(ctx, app, gitRepo, headCommit); err != nil {
				d.logger.Error(fmt.Sprintf("failed to check application: %s", app.Id), zap.Error(err))
			}
		}
	}

	return nil
}

func (d *detector) checkApplication(ctx context.Context, app *model.Application, repo git.Repo, headCommit git.Commit) error {
	var (
		repoDir = repo.GetPath()
		appDir  = filepath.Join(repoDir, app.GitPath.Path)
	)

	// Load config
	cpCfg := d.provider.TerraformConfig
	cfg, err := d.loadApplicationConfiguration(repoDir, app)
	appCfg := cfg.TerraformApplicationSpec
	if err != nil {
		return fmt.Errorf("failed to load application configuration: %w", err)
	}

	gds, ok := cfg.GetGenericApplication()
	if !ok {
		return fmt.Errorf("unsupport application kind %s", cfg.Kind)
	}

	if d.secretDecrypter != nil && gds.Encryption != nil {
		// We have to copy repository into another directory because
		// decrypting the sealed secrets might change the git repository.
		dir, err := os.MkdirTemp("", "detector-git-decrypt")
		if err != nil {
			return fmt.Errorf("failed to prepare a temporary directory for git repository (%w)", err)
		}
		defer os.RemoveAll(dir)

		repo, err = repo.Copy(filepath.Join(dir, "repo"))
		if err != nil {
			return fmt.Errorf("failed to copy the cloned git repository (%w)", err)
		}
		repoDir = repo.GetPath()
		appDir = filepath.Join(repoDir, app.GitPath.Path)

		if err := sourcedecrypter.DecryptSecrets(appDir, *gds.Encryption, d.secretDecrypter); err != nil {
			return fmt.Errorf("failed to decrypt secrets (%w)", err)
		}
	}

	// Set up terraform
	version := appCfg.Input.TerraformVersion
	terraformPath, _, err := toolregistry.DefaultRegistry().Terraform(ctx, version)
	if err != nil {
		return err
	}

	vars := make([]string, 0, len(cpCfg.Vars)+len(appCfg.Input.Vars))
	vars = append(vars, cpCfg.Vars...)
	vars = append(vars, appCfg.Input.Vars...)
	flags := appCfg.Input.CommandFlags
	envs := appCfg.Input.CommandEnvs

	executor := provider.NewTerraform(
		terraformPath,
		appDir,
		provider.WithoutColor(),
		provider.WithVars(vars),
		provider.WithVarFiles(appCfg.Input.VarFiles),
		provider.WithAdditionalFlags(flags.Shared, flags.Init, flags.Plan, flags.Apply),
		provider.WithAdditionalEnvs(envs.Shared, envs.Init, envs.Plan, envs.Apply),
	)

	buf := new(bytes.Buffer)
	if err := executor.Init(ctx, buf); err != nil {
		fmt.Fprintf(buf, "failed while executing terraform init (%v)\n", err)
		return err
	}

	if ws := appCfg.Input.Workspace; ws != "" {
		if err := executor.SelectWorkspace(ctx, ws); err != nil {
			fmt.Fprintf(buf, "failed to select workspace %q (%v). You might need to create the workspace before using by command %q\n",
				ws,
				err,
				"terraform workspace new "+ws,
			)
			return err
		}
		fmt.Fprintf(buf, "selected workspace %q\n", ws)
	}

	result, err := executor.Plan(ctx, buf)
	if err != nil {
		fmt.Fprintf(buf, "failed while executing terraform plan (%v)\n", err)
		return err
	}

	state := makeSyncState(result, headCommit.Hash)

	return d.reporter.ReportApplicationSyncState(ctx, app.Id, state)
}

func makeSyncState(r provider.PlanResult, commit string) model.ApplicationSyncState {
	if r.NoChanges() {
		return model.ApplicationSyncState{
			Status:      model.ApplicationSyncStatus_SYNCED,
			ShortReason: "",
			Reason:      "",
			Timestamp:   time.Now().Unix(),
		}
	}

	total := r.Adds + r.Destroys + r.Changes
	shortReason := fmt.Sprintf("There are %d manifests not synced (%d adds, %d deletes, %d changes)", total, r.Adds, r.Destroys, r.Changes)
	if len(commit) >= 7 {
		commit = commit[:7]
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("Diff between the defined state in Git at commit %s and actual state in cluster:\n\n", commit))
	b.WriteString("--- Expected\n+++ Actual\n\n")

	return model.ApplicationSyncState{
		Status:      model.ApplicationSyncStatus_OUT_OF_SYNC,
		ShortReason: shortReason,
		Reason:      b.String(),
		Timestamp:   time.Now().Unix(),
	}
}

func (d *detector) cloneGitRepository(ctx context.Context, repoID string) (git.Repo, error) {
	repoCfg, ok := d.config.GetRepository(repoID)
	if !ok {
		return nil, fmt.Errorf("repository %s was not found in piped configuration", repoID)
	}
	return d.gitClient.Clone(ctx, repoID, repoCfg.Remote, repoCfg.Branch, "")
}

// listGroupedApplication retrieves all applications those should be handled by this director
// and then groups them by repoID.
func (d *detector) listGroupedApplication() map[string][]*model.Application {
	var (
		apps = d.appLister.ListByPlatformProvider(d.provider.Name)
		m    = make(map[string][]*model.Application)
	)
	for _, app := range apps {
		repoID := app.GitPath.Repo.Id
		m[repoID] = append(m[repoID], app)
	}
	return m
}

func (d *detector) loadApplicationConfiguration(repoPath string, app *model.Application) (*config.Config, error) {
	path := filepath.Join(repoPath, app.GitPath.GetApplicationConfigFilePath())
	cfg, err := config.LoadFromYAML(path)
	if err != nil {
		return nil, err
	}
	if appKind, ok := cfg.Kind.ToApplicationKind(); !ok || appKind != app.Kind {
		return nil, fmt.Errorf("application in application configuration file is not match, got: %s, expected: %s", appKind, app.Kind)
	}
	return cfg, nil
}

func (d *detector) ProviderName() string {
	return d.provider.Name
}
