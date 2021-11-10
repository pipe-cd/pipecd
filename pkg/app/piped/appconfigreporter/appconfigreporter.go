// Copyright 2021 The PipeCD Authors.
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

package appconfigreporter

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipe/pkg/app/api/service/pipedservice"
	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/cache/memorycache"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/git"
	"github.com/pipe-cd/pipe/pkg/model"
)

const (
	defaultLastFetchedCommitCacheSize = 500
)

type apiClient interface {
	UpdateApplicationConfigurations(ctx context.Context, in *pipedservice.UpdateApplicationConfigurationsRequest, opts ...grpc.CallOption) (*pipedservice.UpdateApplicationConfigurationsResponse, error)
	PutUnregisteredApplicationConfigurations(ctx context.Context, in *pipedservice.PutUnregisteredApplicationConfigurationsRequest, opts ...grpc.CallOption) (*pipedservice.PutUnregisteredApplicationConfigurationsResponse, error)
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

type applicationLister interface {
	List() []*model.Application
}

type Reporter struct {
	apiClient         apiClient
	gitClient         gitClient
	applicationLister applicationLister
	config            *config.PipedSpec
	gitRepos          map[string]git.Repo
	gracePeriod       time.Duration
	// Cache for the last fetched commit for each repository.
	lastFetchedCommitCache cache.Cache
	logger                 *zap.Logger
}

func NewReporter(
	apiClient apiClient,
	gitClient gitClient,
	appLister applicationLister,
	cfg *config.PipedSpec,
	gracePeriod time.Duration,
	logger *zap.Logger,
) (*Reporter, error) {
	cache, err := memorycache.NewLRUCache(defaultLastFetchedCommitCacheSize)
	if err != nil {
		return nil, err
	}
	return &Reporter{
		apiClient:              apiClient,
		gitClient:              gitClient,
		applicationLister:      appLister,
		config:                 cfg,
		gracePeriod:            gracePeriod,
		lastFetchedCommitCache: cache,
		logger:                 logger.Named("app-config-reporter"),
	}, nil
}

func (r *Reporter) Run(ctx context.Context) error {
	r.logger.Info("start running app-config-reporter")

	// Pre-clone to cache the registered git repositories.
	r.gitRepos = make(map[string]git.Repo, len(r.config.Repositories))
	for _, repoCfg := range r.config.Repositories {
		repo, err := r.gitClient.Clone(ctx, repoCfg.RepoID, repoCfg.Remote, repoCfg.Branch, "")
		if err != nil {
			r.logger.Error("failed to clone repository",
				zap.String("repo-id", repoCfg.RepoID),
				zap.Error(err),
			)
			return err
		}
		r.gitRepos[repoCfg.RepoID] = repo
	}

	// FIXME: Think about sync interval of app config reporter
	ticker := time.NewTicker(r.config.SyncInterval.Duration())
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := r.checkApps(ctx); err != nil {
				r.logger.Error("failed to check application configurations defined in Git", zap.Error(err))
			}
		case <-ctx.Done():
			r.logger.Info("app-config-reporter has been stopped")
			return nil
		}
	}
}

// checkApps checks and reports two types of applications.
// One is applications registered in Control-plane already, and another is ones that aren't registered yet.
func (r *Reporter) checkApps(ctx context.Context) (err error) {
	if len(r.gitRepos) == 0 {
		r.logger.Info("no repositories were configured for this piped")
		return
	}

	var (
		unusedApps      = make([]*pipedservice.ApplicationConfiguration, 0)
		appsToBeUpdated = make([]*pipedservice.ApplicationConfiguration, 0)
		appsMap         = r.listApplications()
	)
	for repoID, repo := range r.gitRepos {
		if err = repo.Pull(ctx, repo.GetClonedBranch()); err != nil {
			r.logger.Error("failed to update repo to latest",
				zap.String("repo-id", repoID),
				zap.Error(err),
			)
			return
		}

		// TODO: Collect unused application configurations that aren't used yet
		//   Currently, it could be thought the best to open files that suffixed by .pipe.yaml

		var headCommit git.Commit
		// Get the head commit of the repository.
		headCommit, err = repo.GetLatestCommit(ctx)
		if err != nil {
			return
		}
		var lastFetchedCommit interface{}
		lastFetchedCommit, err = r.lastFetchedCommitCache.Get(repoID)
		if err != nil && !errors.Is(err, cache.ErrNotFound) {
			r.logger.Error("failed to get the last fetched commit from cache", zap.Error(err))
		}
		if headCommit.Hash == lastFetchedCommit.(string) {
			continue
		}
		apps, ok := appsMap[repoID]
		if !ok {
			continue
		}
		for _, app := range apps {
			gitPath := app.GetGitPath()
			_ = filepath.Join(repo.GetPath(), gitPath.Path, gitPath.ConfigFilename)
			// TODO: Collect applications that need to be updated
		}

		defer func() {
			if err == nil {
				r.lastFetchedCommitCache.Put(repoID, headCommit)
			}
		}()
	}
	if len(unusedApps) > 0 {
		_, err = r.apiClient.PutUnregisteredApplicationConfigurations(
			ctx,
			&pipedservice.PutUnregisteredApplicationConfigurationsRequest{
				Applications: unusedApps,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to put the latest unregistered application configurations: %w", err)
		}
	}

	if len(appsToBeUpdated) == 0 {
		return nil
	}
	_, err = r.apiClient.UpdateApplicationConfigurations(
		ctx,
		&pipedservice.UpdateApplicationConfigurationsRequest{
			Applications: appsToBeUpdated,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update application configurations: %w", err)
	}

	return
}

// listApplications retrieves all applications that should be handled by this piped
// and then groups them by repoID.
func (r *Reporter) listApplications() map[string][]*model.Application {
	var (
		apps       = r.applicationLister.List()
		repoToApps = make(map[string][]*model.Application)
	)
	for _, app := range apps {
		repoId := app.GitPath.Repo.Id
		if _, ok := repoToApps[repoId]; !ok {
			repoToApps[repoId] = []*model.Application{app}
		} else {
			repoToApps[repoId] = append(repoToApps[repoId], app)
		}
	}
	return repoToApps
}
