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
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipe/pkg/app/api/service/pipedservice"
	"github.com/pipe-cd/pipe/pkg/cache/memorycache"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/git"
	"github.com/pipe-cd/pipe/pkg/model"
)

const (
	commandCheckInterval                = 10 * time.Second
	defaultLastTriggeredCommitCacheSize = 500
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
	Get(ctx context.Context, id string) (*model.Environment, error)
}

type notifier interface {
	Notify(event model.NotificationEvent)
}

type Trigger struct {
	apiClient         apiClient
	gitClient         gitClient
	applicationLister applicationLister
	commandLister     commandLister
	environmentLister environmentLister
	notifier          notifier
	config            *config.PipedSpec
	commitStore       *lastTriggeredCommitStore
	gitRepos          map[string]git.Repo
	gracePeriod       time.Duration
	logger            *zap.Logger
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
) (*Trigger, error) {

	cache, err := memorycache.NewLRUCache(defaultLastTriggeredCommitCacheSize)
	if err != nil {
		return nil, err
	}
	commitStore := &lastTriggeredCommitStore{
		apiClient: apiClient,
		cache:     cache,
	}

	t := &Trigger{
		apiClient:         apiClient,
		gitClient:         gitClient,
		applicationLister: appLister,
		commandLister:     commandLister,
		environmentLister: environmentLister,
		notifier:          notifier,
		config:            cfg,
		commitStore:       commitStore,
		gitRepos:          make(map[string]git.Repo, len(cfg.Repositories)),
		gracePeriod:       gracePeriod,
		logger:            logger.Named("trigger"),
	}

	return t, nil
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
			t.checkNewCommands(ctx)

		case <-commitTicker.C:
			t.checkNewCommits(ctx)

		case <-ctx.Done():
			break L
		}
	}

	t.logger.Info("deployment trigger has been stopped")
	return nil
}

func (t *Trigger) GetLastTriggeredCommitGetter() LastTriggeredCommitGetter {
	return t.commitStore
}

func (t *Trigger) checkNewCommands(ctx context.Context) error {
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
			if err := cmd.Report(ctx, model.CommandStatus_COMMAND_FAILED, nil, nil); err != nil {
				t.logger.Error("failed to report command status", zap.Error(err))
			}
			continue
		}

		metadata := map[string]string{
			triggeredDeploymentIDKey: d.Id,
		}
		if err := cmd.Report(ctx, model.CommandStatus_COMMAND_SUCCEEDED, metadata, nil); err != nil {
			t.logger.Error("failed to report command status", zap.Error(err))
		}
	}

	return nil
}

func (t *Trigger) checkNewCommits(ctx context.Context) error {
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
		d := NewDeterminer(gitRepo, headCommit.Hash, t.commitStore, t.logger)

		for _, app := range apps {
			shouldTrigger, err := d.ShouldTrigger(ctx, app, false)
			if err != nil {
				t.logger.Error(fmt.Sprintf("failed to check application: %s", app.Id), zap.Error(err))
				continue
			}

			if !shouldTrigger {
				t.commitStore.Put(app.Id, headCommit.Hash)
				continue
			}

			// Build deployment model and send a request to API to create a new deployment.
			t.logger.Info("application should be synced because of the new commit")
			if _, err := t.triggerDeployment(ctx, app, branch, headCommit, "", model.SyncStrategy_AUTO); err != nil {
				t.logger.Error(fmt.Sprintf("failed to trigger application: %s", app.Id), zap.Error(err))
			}
			t.commitStore.Put(app.Id, headCommit.Hash)
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
	t.commitStore.Put(app.Id, headCommit.Hash)

	return d, nil
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
