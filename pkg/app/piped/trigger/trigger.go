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
// that detects a list of application should be synced (by new commit, sync command or configuration driff)
// and then sends request to the control-plane to create a new Deployment.
package trigger

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/kapetaniosci/pipe/pkg/app/api/service/pipedservice"
	"github.com/kapetaniosci/pipe/pkg/config"
	"github.com/kapetaniosci/pipe/pkg/git"
	"github.com/kapetaniosci/pipe/pkg/model"
)

type apiClient interface {
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

type DeploymentTrigger struct {
	apiClient         apiClient
	gitClient         gitClient
	applicationLister applicationLister
	commandLister     commandLister
	config            *config.PipedSpec
	triggeredCommits  map[string]string
	gitRepos          map[string]git.Repo
	gracePeriod       time.Duration
	logger            *zap.Logger
}

// NewTrigger creates a new instance for DeploymentTrigger.
func NewTrigger(
	apiClient apiClient,
	gitClient gitClient,
	appLister applicationLister,
	commandLister commandLister,
	cfg *config.PipedSpec,
	gracePeriod time.Duration,
	logger *zap.Logger,
) *DeploymentTrigger {

	return &DeploymentTrigger{
		apiClient:         apiClient,
		gitClient:         gitClient,
		applicationLister: appLister,
		commandLister:     commandLister,
		config:            cfg,
		triggeredCommits:  make(map[string]string),
		gitRepos:          make(map[string]git.Repo, len(cfg.Repositories)),
		gracePeriod:       gracePeriod,
		logger:            logger.Named("deployment-trigger"),
	}
}

// Run starts running DeploymentTrigger until the specified context has done.
// This also waits for its cleaning up before returning.
func (t *DeploymentTrigger) Run(ctx context.Context) error {
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

func (t *DeploymentTrigger) check(ctx context.Context) error {
	if len(t.gitRepos) == 0 {
		t.logger.Info("no repositories were configured for this piped")
		return nil
	}

	// List all applications that should be handled by this piped
	// and then group them by repository.
	var applications = t.listApplications(ctx)

	// TODO: We may want to apply worker model here to run them concurrently.
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

func (t *DeploymentTrigger) checkApplication(ctx context.Context, app *model.Application, repo git.Repo, branch string, headCommit git.Commit) error {
	// Get the most recently applied commit of this application.
	// TODO: If it is not in the memory cache, we have to call the API to list the deployments
	// and use the commit sha of the most recent one.
	triggeredCommitHash, ok := t.triggeredCommits[app.Id]
	if !ok {
		t.triggeredCommits[app.Id] = "retrieved-one"
	}

	// Check whether the most recently applied one is the head commit or not.
	// If not, nothing to do for this time.
	if triggeredCommitHash == headCommit.Hash {
		t.logger.Info(fmt.Sprintf("no update to sync for application: %s, hash: %s", app.Id, headCommit.Hash))
		return nil
	}

	trigger := func() error {
		// Build deployment model and send a request to API to create a new deployment.
		t.logger.Info(fmt.Sprintf("application %s should be synced because of the new commit", app.Id),
			zap.String("previous-commit-hash", triggeredCommitHash),
			zap.String("head-commit-hash", headCommit.Hash),
		)
		if err := t.triggerDeployment(ctx, app, repo, branch, headCommit); err != nil {
			return err
		}
		t.triggeredCommits[app.Id] = headCommit.Hash
		return nil
	}

	// There is no previous deployment so we don't need to check anymore.
	// Just do it.
	if triggeredCommitHash == "" {
		return trigger()
	}

	// List the changed files between those two commits and
	// determine whether this application was touch by those changed files.
	changedFiles, err := repo.ChangedFiles(ctx, triggeredCommitHash, headCommit.Hash)
	if err != nil {
		return err
	}
	if touched := isTouchedByChangedFiles(app.GitPath.Path, nil, changedFiles); !touched {
		t.logger.Info(fmt.Sprintf("application %s was not touched by the new commit", app.Id),
			zap.String("previous-commit-hash", triggeredCommitHash),
			zap.String("head-commit-hash", headCommit.Hash),
		)
		t.triggeredCommits[app.Id] = headCommit.Hash
		return nil
	}

	return trigger()
}

// listApplications retrieves all applications those should be handled by this piped
// and then groups them by repoID.
func (t *DeploymentTrigger) listApplications(ctx context.Context) map[string][]*model.Application {
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
