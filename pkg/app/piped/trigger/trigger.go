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
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipe/pkg/app/api/service/pipedservice"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/filematcher"
	"github.com/pipe-cd/pipe/pkg/git"
	"github.com/pipe-cd/pipe/pkg/model"
)

var (
	commandCheckInterval = 10 * time.Second
)

const (
	triggeredDeploymentIDKey = "TriggeredDeploymentID"
)

type apiClient interface {
	GetApplicationMostRecentDeployment(ctx context.Context, req *pipedservice.GetApplicationMostRecentDeploymentRequest, opts ...grpc.CallOption) (*pipedservice.GetApplicationMostRecentDeploymentResponse, error)
	CreateDeployment(ctx context.Context, in *pipedservice.CreateDeploymentRequest, opts ...grpc.CallOption) (*pipedservice.CreateDeploymentResponse, error)
	ReportApplicationMostRecentDeployment(ctx context.Context, req *pipedservice.ReportApplicationMostRecentDeploymentRequest, opts ...grpc.CallOption) (*pipedservice.ReportApplicationMostRecentDeploymentResponse, error)
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

type applicationLister interface {
	Get(id string) (*model.Application, bool)
	List() []*model.Application
}

type commandLister interface {
	ListApplicationCommands() []model.ReportableCommand
}

type environmentLister interface {
	Get(id string) (*model.Environment, bool)
}

type notifier interface {
	Notify(event model.NotificationEvent)
}

type Trigger struct {
	apiClient                    apiClient
	gitClient                    gitClient
	applicationLister            applicationLister
	commandLister                commandLister
	environmentLister            environmentLister
	notifier                     notifier
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
	environmentLister environmentLister,
	notifier notifier,
	cfg *config.PipedSpec,
	gracePeriod time.Duration,
	logger *zap.Logger,
) *Trigger {

	return &Trigger{
		apiClient:                    apiClient,
		gitClient:                    gitClient,
		applicationLister:            appLister,
		commandLister:                commandLister,
		environmentLister:            environmentLister,
		notifier:                     notifier,
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

	commitTicker := time.NewTicker(time.Duration(t.config.SyncInterval))
	defer commitTicker.Stop()

	commandTicker := time.NewTicker(commandCheckInterval)
	defer commandTicker.Stop()

L:
	for {
		select {

		case <-commandTicker.C:
			t.checkCommand(ctx)

		case <-commitTicker.C:
			t.checkCommit(ctx)

		case <-ctx.Done():
			break L
		}
	}

	t.logger.Info("deployment trigger has been stopped")
	return nil
}

func (t *Trigger) checkCommand(ctx context.Context) error {
	commands := t.commandLister.ListApplicationCommands()
	for _, cmd := range commands {
		syncCmd := cmd.GetSyncApplication()
		if syncCmd == nil {
			continue
		}
		app, ok := t.applicationLister.Get(syncCmd.ApplicationId)
		if !ok {
			t.logger.Warn("detected an AppSync command for an unregistered application",
				zap.String("command", cmd.Id),
				zap.String("app-id", syncCmd.ApplicationId),
				zap.String("commander", cmd.Commander),
			)
			continue
		}
		d, err := t.syncApplication(ctx, app, cmd.Commander, syncCmd.SyncStrategy)
		if err != nil {
			t.logger.Error("failed to sync application",
				zap.String("app-id", app.Id),
				zap.Error(err),
			)
			if err := cmd.Report(ctx, model.CommandStatus_COMMAND_FAILED, nil); err != nil {
				t.logger.Error("failed to report command status", zap.Error(err))
			}
			continue
		}

		metadata := map[string]string{
			triggeredDeploymentIDKey: d.Id,
		}
		if err := cmd.Report(ctx, model.CommandStatus_COMMAND_SUCCEEDED, metadata); err != nil {
			t.logger.Error("failed to report command status", zap.Error(err))
		}
	}
	return nil
}

func (t *Trigger) syncApplication(ctx context.Context, app *model.Application, commander string, syncStrategy model.SyncStrategy) (*model.Deployment, error) {
	_, branch, headCommit, err := t.updateRepoToLatest(ctx, app.GitPath.Repo.Id)
	if err != nil {
		return nil, err
	}

	// Build deployment model and send a request to API to create a new deployment.
	t.logger.Info(fmt.Sprintf("application %s will be synced because of a sync command", app.Id),
		zap.String("head-commit", headCommit.Hash),
	)
	d, err := t.triggerDeployment(ctx, app, branch, headCommit, commander, syncStrategy)
	if err != nil {
		return nil, err
	}
	t.mostRecentlyTriggeredCommits[app.Id] = headCommit.Hash

	return d, nil
}

func (t *Trigger) checkCommit(ctx context.Context) error {
	if len(t.gitRepos) == 0 {
		t.logger.Info("no repositories were configured for this piped")
		return nil
	}

	// List all applications that should be handled by this piped
	// and then group them by repository.
	var applications = t.listApplications()

	// ENHANCEMENT: We may want to apply worker model here to run them concurrently.
	for repoID, apps := range applications {
		gitRepo, branch, headCommit, err := t.updateRepoToLatest(ctx, repoID)
		if err != nil {
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
		zap.String("app", app.Name),
		zap.String("app-id", app.Id),
		zap.String("head-commit", headCommit.Hash),
	)

	// Get the most recently triggered commit of this application.
	// Most of the cases that data can be loaded from in-memory cache but
	// when the piped is restared that data will be cleared too.
	// So in that case, we have to make an API call.
	preCommitHash := t.mostRecentlyTriggeredCommits[app.Id]
	if preCommitHash == "" {
		mostRecent, err := t.getMostRecentlyTriggeredDeployment(ctx, app.Id)
		switch {
		case err == nil:
			preCommitHash = mostRecent.Trigger.Commit.Hash
			t.mostRecentlyTriggeredCommits[app.Id] = preCommitHash

		case status.Code(err) == codes.NotFound:
			logger.Info("there is no previously triggered commit for this application")

		default:
			logger.Error("unable to get the most recently triggered deployment", zap.Error(err))
			return err
		}
	}

	// Check whether the most recently applied one is the head commit or not.
	// If so, nothing to do for this time.
	if headCommit.Hash == preCommitHash {
		logger.Info(fmt.Sprintf("no update to sync for application, hash: %s", headCommit.Hash))
		return nil
	}

	trigger := func() error {
		// Build deployment model and send a request to API to create a new deployment.
		logger.Info("application should be synced because of the new commit",
			zap.String("most-recently-triggered-commit", preCommitHash),
		)
		if _, err := t.triggerDeployment(ctx, app, branch, headCommit, "", model.SyncStrategy_AUTO); err != nil {
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

	deployConfig, err := loadDeploymentConfiguration(repo.GetPath(), app)
	if err != nil {
		return err
	}

	touched, err := isTouchedByChangedFiles(app.GitPath.Path, deployConfig.TriggerPaths, changedFiles)
	if err != nil {
		return err
	}
	if !touched {
		logger.Info("application was not touched by the new commit",
			zap.String("most-recently-triggered-commit", preCommitHash),
		)
		t.mostRecentlyTriggeredCommits[app.Id] = headCommit.Hash
		return nil
	}

	return trigger()
}

func (t *Trigger) updateRepoToLatest(ctx context.Context, repoID string) (repo git.Repo, branch string, headCommit git.Commit, err error) {
	var ok bool

	// Find the application repo from pre-loaded ones.
	repo, ok = t.gitRepos[repoID]
	if !ok {
		t.logger.Warn("detected some applications binding with a non existent repository", zap.String("repo-id", repoID))
		err = fmt.Errorf("missing repository")
		return
	}
	branch = repo.GetClonedBranch()

	// Fetch to update the repository and then
	if err = repo.Pull(ctx, branch); err != nil {
		if ctx.Err() != context.Canceled {
			t.logger.Error("failed to update repository branch",
				zap.String("repo-id", repoID),
				zap.Error(err),
			)
		}
		return
	}

	// Get the head commit of the repository.
	headCommit, err = repo.GetLatestCommit(ctx)
	if err != nil {
		// TODO: Find a better way to skip the CANCELLED error log while shutting down.
		if ctx.Err() != context.Canceled {
			t.logger.Error("failed to get head commit hash",
				zap.String("repo-id", repoID),
				zap.Error(err),
			)
		}
		return
	}

	return
}

// listApplications retrieves all applications those should be handled by this piped
// and then groups them by repoID.
func (t *Trigger) listApplications() map[string][]*model.Application {
	var (
		apps = t.applicationLister.List()
		m    = make(map[string][]*model.Application)
	)
	for _, app := range apps {
		repoId := app.GitPath.Repo.Id
		if _, ok := m[repoId]; !ok {
			m[repoId] = []*model.Application{app}
		} else {
			m[repoId] = append(m[repoId], app)
		}
	}
	return m
}

func (t *Trigger) getMostRecentlyTriggeredDeployment(ctx context.Context, applicationID string) (*model.ApplicationDeploymentReference, error) {
	var (
		err   error
		resp  *pipedservice.GetApplicationMostRecentDeploymentResponse
		retry = pipedservice.NewRetry(3)
		req   = &pipedservice.GetApplicationMostRecentDeploymentRequest{
			ApplicationId: applicationID,
			Status:        model.DeploymentStatus_DEPLOYMENT_PENDING,
		}
	)

	for retry.WaitNext(ctx) {
		if resp, err = t.apiClient.GetApplicationMostRecentDeployment(ctx, req); err == nil {
			return resp.Deployment, nil
		}
		if !pipedservice.Retriable(err) {
			return nil, err
		}
	}
	return nil, err
}

func loadDeploymentConfiguration(repoPath string, app *model.Application) (*config.GenericDeploymentSpec, error) {
	path := filepath.Join(repoPath, app.GitPath.GetDeploymentConfigFilePath())
	cfg, err := config.LoadFromYAML(path)
	if err != nil {
		return nil, err
	}
	if appKind, ok := config.ToApplicationKind(cfg.Kind); !ok || appKind != app.Kind {
		return nil, fmt.Errorf("invalid application kind in the deployment config file, got: %s, expected: %s", appKind, app.Kind)
	}

	spec, ok := cfg.GetGenericDeployment()
	if !ok {
		return nil, fmt.Errorf("unsupported application kind: %s", app.Kind)
	}

	return &spec, nil
}

func isTouchedByChangedFiles(appDir string, changes []string, changedFiles []string) (bool, error) {
	if !strings.HasSuffix(appDir, "/") {
		appDir += "/"
	}

	// If any files inside the application directory was changed
	// this application is considered as touched.
	for _, cf := range changedFiles {
		if ok := strings.HasPrefix(cf, appDir); ok {
			return true, nil
		}
	}

	// If any changed files matches the specified "changes"
	// this application is consided as touched too.
	for _, change := range changes {
		matcher, err := filematcher.NewPatternMatcher([]string{change})
		if err != nil {
			return false, err
		}
		if matcher.MatchesAny(changedFiles) {
			return true, nil
		}
	}

	return false, nil
}
