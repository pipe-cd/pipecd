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

package ecs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pipe-cd/pipecd/pkg/app/piped/livestatestore/ecs"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/ecs"
	"github.com/pipe-cd/pipecd/pkg/app/piped/sourcedecrypter"
	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
	"go.uber.org/zap"
)

type applicationLister interface {
	ListByPlatformProvider(name string) []*model.Application
}

type gitClient interface {
	Clone(ctx context.Context, repoId, remote, branch, destination string) (git.Repo, error)
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
	stateGetter       ecs.Getter
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
	stateGetter ecs.Getter,
	reporter reporter,
	appManifestsCache cache.Cache,
	cfg *config.PipedSpec,
	sd secretDecrypter,
	logger *zap.Logger,
) Detector {
	logger = logger.Named("ecs-detector").With(
		zap.String("ecs-provider", cp.Name),
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
	d.logger.Info("start running drift for ecs application")

	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			d.check(ctx)

		case <-ctx.Done():
			d.logger.Info("drift detector for ecs application has been stopped")
			return nil
		}
	}
}

func (d *detector) ProviderName() string {
	return d.provider.Name
}

func (d *detector) check(ctx context.Context) error {
	appsByRepo := d.listGroupedApplication()

	for repoID, apps := range appsByRepo {
		gitRepo, ok := d.gitRepos[repoID]
		if !ok {
			gr, err := d.cloneGitRepository(ctx, repoID)
			if err != nil {
				d.logger.Error("failed to clone git reposigory",
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

func (d *detector) cloneGitRepository(ctx context.Context, repoID string) (git.Repo, error) {
	repoCfg, ok := d.config.GetRepository(repoID)
	if !ok {
		return nil, fmt.Errorf("repository %s was not found in piped configuration", repoID)
	}
	return d.gitClient.Clone(ctx, repoID, repoCfg.Remote, repoCfg.Branch, "")
}

func (d *detector) checkApplication(ctx context.Context, app *model.Application, repo git.Repo, headCommit git.Commit) error {
	headManifests, err := d.loadHeadManifests(ctx, app, repo, headCommit)
	if err != nil {
		return err
	}
	d.logger.Info(fmt.Sprintf("application %s has %d manifests at commit %s", app.Id, len(headManifests), headCommit.Hash))

	liveManifests := d.stateGetter.GetAppLiveManifests(app.Id)
	return nil
}

func (d *detector) loadHeadManifests(ctx context.Context, app *model.Application, repo git.Repo, headCommit git.Commit) ([]provider.Manifest, error) {
	var (
		manifestCache = provider.AppManifestsCache{
			AppID:  app.Id,
			Cache:  d.appManifestsCache,
			Logger: d.logger,
		}
		repoDir = repo.GetPath()
		appDir  = filepath.Join(repoDir, app.GitPath.Path)
	)

	manifests, ok := manifestCache.Get(headCommit.Hash)
	if !ok {
		cfg, err := d.loadApplicationConfiguration(repoDir, app)
		if err != nil {
			return nil, fmt.Errorf("failed to load application configuration: %w", err)
		}
		gds, ok := cfg.GetGenericApplication()
		if !ok {
			return nil, fmt.Errorf("unsuppot application kind %s", cfg.Kind)
		}

		if d.secretDecrypter != nil && gds.Encryption != nil {
			// We have to copy repository into another directory because
			// decrypting the sealed secrets might change the git repository.
			dir, err := os.MkdirTemp("", "detector-git-decrypt")
			if err != nil {
				return nil, fmt.Errorf("failed to prepare a temporary directory for git repository (%w)", err)
			}
			defer os.RemoveAll(dir)

			repo, err = repo.Copy(filepath.Join(dir, "repo"))
			if err != nil {
				return nil, fmt.Errorf("failed to copy the cloned git repository (%w)", err)
			}
			repoDir = repo.GetPath()
			appDir = filepath.Join(repoDir, app.GitPath.Path)

			if err := sourcedecrypter.DecryptSecrets(appDir, *gds.Encryption, d.secretDecrypter); err != nil {
				return nil, fmt.Errorf("failed to decrypt secrets (%w)", err)
			}
		}

		loader := provider.NewLoader(
			app.Name,
			appDir,
			repoDir,
			app.GitPath.ConfigFilename,
			cfg.ECSApplicationSpec.Input,
			d.gitClient,
			d.logger,
		)
		manifests, err = loader.LoadManifests(ctx)
		if err != nil {
			err = fmt.Errorf("failed to load new manifests: %w", err)
			return nil, err
		}
		manifestCache.Put(headCommit.Hash, manifests)
	}

	return manifests, nil
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
