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
	"time"

	"github.com/pipe-cd/pipecd/pkg/app/piped/livestatestore/ecs"
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
	headManifests, err := d.loadHeadManifests(app, repo, headCommit)
}

func (d *detector) loadHeadManifests(app *model.Application, repo git.Repo, headCommit git.Commit) ([]provider., error) {

}
