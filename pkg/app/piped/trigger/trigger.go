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

// Package trigger provides a piped component
// that detects a list of application should be synced (by new commit, sync command or configuration drift)
// and then sends request to the control-plane to create a new Deployment.
package trigger

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipe/pkg/app/api/service/pipedservice"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/git"
	"github.com/pipe-cd/pipe/pkg/model"
)

type apiClient interface {
	GetMostRecentDeployment(ctx context.Context, req *pipedservice.GetMostRecentDeploymentRequest, opts ...grpc.CallOption) (*pipedservice.GetMostRecentDeploymentResponse, error)
	CreateDeployment(ctx context.Context, in *pipedservice.CreateDeploymentRequest, opts ...grpc.CallOption) (*pipedservice.CreateDeploymentResponse, error)
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

type applicationLister interface {
	List() []*model.Application
}

type commandLister interface {
	ListApplicationCommands() []model.ReportableCommand
}

type Trigger struct {
	apiClient                    apiClient
	gitClient                    gitClient
	applicationLister            applicationLister
	commandLister                commandLister
	config                       *config.PipedSpec
	mostRecentlyTriggeredCommits map[string]string
	gitRepos                     map[string]git.Repo
	gracePeriod                  time.Duration
	logger                       *zap.Logger
}

// NewTrigger creates a new instance for Trigger.
func NewTrigger(
	apiClient apiClient,
	gitClient gitClient,
	appLister applicationLister,
	commandLister commandLister,
	cfg *config.PipedSpec,
	gracePeriod time.Duration,
	logger *zap.Logger,
) *Trigger {

	return &Trigger{
		apiClient:                    apiClient,
		gitClient:                    gitClient,
		applicationLister:            appLister,
		commandLister:                commandLister,
		config:                       cfg,
		mostRecentlyTriggeredCommits: make(map[string]string),
		gitRepos:                     make(map[string]git.Repo, len(cfg.Repositories)),
		gracePeriod:                  gracePeriod,
		logger:                       logger.Named("trigger"),
	}
}

// Run starts running Trigger until the specified context has done.
// This also waits for its cleaning up before returning.
func (t *Trigger) Run(ctx context.Context) error {
	t.logger.Info("start running deployment trigger")

	// Pre-clone to cache the registered git repositories.
	t.gitRepos = make(map[string]git.Repo, len(t.config.Repositories))
	for _, r := range t.config.Repositories {
		repo, err := t.gitClient.Clone(ctx, r.RepoID, r.Remote, r.Branch, "")
		if err != nil {
			t.logger.Error("failed to clone repository",
				zap.String("repo-id", r.RepoID),
				zap.Error(err),
			)
			return err
		}
		t.gitRepos[r.RepoID] = repo
	}

	ticker := time.NewTicker(time.Duration(t.config.SyncInterval))
	defer ticker.Stop()

L:
	for {
		select {
		case <-ctx.Done():
			break L
		case <-ticker.C:
			t.check(ctx)
		}
	}

	t.logger.Info("deployment trigger has been stopped")
	return nil
}

func (t *Trigger) check(ctx context.Context) error {
	if len(t.gitRepos) == 0 {
		t.logger.Info("no repositories were configured for this piped")
		return nil
	}

	// List all applications that should be handled by this piped
	// and then group them by repository.
	var applications = t.listApplications()

	// ENHANCEMENT: We may want to apply worker model here to run them concurrently.
	for repoID, apps := range applications {
		gitRepo, ok := t.gitRepos[repoID]
		if !ok {
			t.logger.Warn("detected some applications are binding with an non existent repository",
				zap.String("repo-id", repoID),
				zap.String("first-application-id", apps[0].Id),
			)
			continue
		}
		branch := gitRepo.GetClonedBranch()

		// Fetch to update the repository and then
		if err := gitRepo.Pull(ctx, branch); err != nil {
			t.logger.Error("failed to update repository branch",
				zap.String("repo-id", repoID),
				zap.Error(err),
			)
			continue
		}

		// Get the head commit of the repository.
		headCommit, err := gitRepo.GetLatestCommit(ctx)
		if err != nil {
			t.logger.Error("failed to get head commit hash",
				zap.String("repo-id", repoID),
				zap.Error(err),
			)
			continue
		}

		for _, app := range apps {
			if err := t.checkApplication(ctx, app, gitRepo, branch, headCommit); err != nil {
				t.logger.Error(fmt.Sprintf("failed to check application: %s", app.Id), zap.Error(err))
			}
		}
	}
	return nil
}

func (t *Trigger) checkApplication(ctx context.Context, app *model.Application, repo git.Repo, branch string, headCommit git.Commit) error {
	logger := t.logger.With(
		zap.String("application", app.Id),
		zap.String("head-commit", headCommit.Hash),
	)

	// Get the most recently triggered commit of this application.
	// Most of the cases that data can be loaded from in-memory cache but
	// when the piped is restared that data will be cleared too.
	// So in that case, we have to make an API call.
	preCommitHash := t.mostRecentlyTriggeredCommits[app.Id]
	if preCommitHash == "" {
		mostRecent, err := t.getMostRecentDeployment(ctx, app.Id)
		if err == nil {
			preCommitHash = mostRecent.CommitHash()
			t.mostRecentlyTriggeredCommits[app.Id] = preCommitHash
		} else if status.Code(err) == codes.NotFound {
			logger.Info("there is no previously triggered commit for this application")
		} else {
			logger.Error("unabled to get the most recently triggered deployment", zap.Error(err))
		}
	}

	// Check whether the most recently applied one is the head commit or not.
	// If not, nothing to do for this time.
	if headCommit.Hash == preCommitHash {
		logger.Info(fmt.Sprintf("no update to sync for application: %s, hash: %s", app.Id, headCommit.Hash))
		return nil
	}

	trigger := func() error {
		// Build deployment model and send a request to API to create a new deployment.
		logger.Info(fmt.Sprintf("application %s should be synced because of the new commit", app.Id),
			zap.String("most-recently-triggered-commit", preCommitHash),
		)
		if err := t.triggerDeployment(ctx, app, repo, branch, headCommit); err != nil {
			return err
		}
		t.mostRecentlyTriggeredCommits[app.Id] = headCommit.Hash
		return nil
	}

	// There is no previous deployment so we don't need to check anymore.
	// Just do it.
	if preCommitHash == "" {
		return trigger()
	}

	// List the changed files between those two commits and
	// determine whether this application was touch by those changed files.
	changedFiles, err := repo.ChangedFiles(ctx, preCommitHash, headCommit.Hash)
	if err != nil {
		return err
	}
	if touched := isTouchedByChangedFiles(app.GitPath.Path, nil, changedFiles); !touched {
		logger.Info(fmt.Sprintf("application %s was not touched by the new commit", app.Id),
			zap.String("most-recently-triggered-commit", preCommitHash),
		)
		t.mostRecentlyTriggeredCommits[app.Id] = headCommit.Hash
		return nil
	}

	return trigger()
}

// listApplications retrieves all applications those should be handled by this piped
// and then groups them by repoID.
func (t *Trigger) listApplications() map[string][]*model.Application {
	var (
		apps = t.applicationLister.List()
		m    = make(map[string][]*model.Application)
	)
	for _, app := range apps {
		repoId := app.GitPath.RepoId
		if _, ok := m[repoId]; !ok {
			m[repoId] = []*model.Application{app}
		} else {
			m[repoId] = append(m[repoId], app)
		}
	}
	return m
}

func (t *Trigger) getMostRecentDeployment(ctx context.Context, applicationID string) (*model.Deployment, error) {
	var (
		err   error
		resp  *pipedservice.GetMostRecentDeploymentResponse
		retry = pipedservice.NewRetry(3)
		req   = &pipedservice.GetMostRecentDeploymentRequest{
			ApplicationId: applicationID,
		}
	)

	for retry.WaitNext(ctx) {
		if resp, err = t.apiClient.GetMostRecentDeployment(ctx, req); err == nil {
			return resp.Deployment, nil
		}
		if !pipedservice.Retriable(err) {
			return nil, err
		}
	}
	return nil, err
}

func isTouchedByChangedFiles(appDir string, dependencyDirs []string, changedFiles []string) bool {
	// If any files inside the application directory was changed
	// this application is considered as touched.
	for _, cf := range changedFiles {
		if ok := strings.HasPrefix(cf, appDir); ok {
			return true
		}
	}

	// If any files inside the app's dependencies was changed
	// this application is consided as touched too.
	for _, depDir := range dependencyDirs {
		for _, cf := range changedFiles {
			if ok := strings.HasPrefix(cf, depDir); ok {
				return true
			}
		}
	}

	return false
}
