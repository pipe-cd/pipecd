// Copyright 2023 The PipeCD Authors.
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

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/cache/memorycache"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	ondemandCheckInterval               = 10 * time.Second
	defaultLastTriggeredCommitCacheSize = 500
)

type apiClient interface {
	GetApplicationMostRecentDeployment(ctx context.Context, req *pipedservice.GetApplicationMostRecentDeploymentRequest, opts ...grpc.CallOption) (*pipedservice.GetApplicationMostRecentDeploymentResponse, error)
	CreateDeployment(ctx context.Context, in *pipedservice.CreateDeploymentRequest, opts ...grpc.CallOption) (*pipedservice.CreateDeploymentResponse, error)
	GetDeployment(ctx context.Context, in *pipedservice.GetDeploymentRequest, opts ...grpc.CallOption) (*pipedservice.GetDeploymentResponse, error)
	ReportApplicationMostRecentDeployment(ctx context.Context, req *pipedservice.ReportApplicationMostRecentDeploymentRequest, opts ...grpc.CallOption) (*pipedservice.ReportApplicationMostRecentDeploymentResponse, error)
	CreateDeploymentChain(ctx context.Context, in *pipedservice.CreateDeploymentChainRequest, opts ...grpc.CallOption) (*pipedservice.CreateDeploymentChainResponse, error)
	ReportApplicationSyncState(ctx context.Context, in *pipedservice.ReportApplicationSyncStateRequest, opts ...grpc.CallOption) (*pipedservice.ReportApplicationSyncStateResponse, error)
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

type notifier interface {
	Notify(event model.NotificationEvent)
}

type candidate struct {
	application *model.Application
	kind        model.TriggerKind
	command     model.ReportableCommand
}

func (c *candidate) HasCommand() bool {
	return c.kind == model.TriggerKind_ON_COMMAND || c.kind == model.TriggerKind_ON_CHAIN
}

type Trigger struct {
	apiClient         apiClient
	gitClient         gitClient
	applicationLister applicationLister
	commandLister     commandLister
	notifier          notifier
	config            *config.PipedSpec
	commitStore       *lastTriggeredCommitStore
	gitRepos          map[string]git.Repo
	gracePeriod       time.Duration
	logger            *zap.Logger
}

func NewTrigger(
	apiClient apiClient,
	gitClient gitClient,
	appLister applicationLister,
	commandLister commandLister,
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
		notifier:          notifier,
		config:            cfg,
		commitStore:       commitStore,
		gitRepos:          make(map[string]git.Repo, len(cfg.Repositories)),
		gracePeriod:       gracePeriod,
		logger:            logger.Named("trigger"),
	}

	return t, nil
}

func (t *Trigger) Run(ctx context.Context) error {
	t.logger.Info("start running deployment trigger")

	// Pre cloning to cache the registered git repositories.
	t.gitRepos = make(map[string]git.Repo, len(t.config.Repositories))
	for _, r := range t.config.Repositories {
		repo, err := t.gitClient.Clone(ctx, r.RepoID, r.Remote, r.Branch, "")
		if err != nil {
			t.logger.Error(fmt.Sprintf("failed to clone git repository %s", r.RepoID), zap.Error(err))
			return err
		}
		t.gitRepos[r.RepoID] = repo
	}

	syncTicker := time.NewTicker(time.Duration(t.config.SyncInterval))
	defer syncTicker.Stop()

	ondemandTicker := time.NewTicker(ondemandCheckInterval)
	defer ondemandTicker.Stop()

	for {
		select {
		case <-syncTicker.C:
			var (
				commitCandidates    = t.listCommitCandidates()
				outOfSyncCandidates = t.listOutOfSyncCandidates()
				candidates          = append(commitCandidates, outOfSyncCandidates...)
			)
			t.logger.Info(fmt.Sprintf("found %d candidates: %d commit candidates and %d out_of_sync candidates",
				len(candidates),
				len(commitCandidates),
				len(outOfSyncCandidates),
			))
			t.checkCandidates(ctx, candidates)

		case <-ondemandTicker.C:
			candidates := t.listCommandCandidates()
			if len(candidates) > 0 {
				t.logger.Info(fmt.Sprintf("found %d command candidates", len(candidates)))
			} else {
				t.logger.Debug(fmt.Sprintf("found %d command candidates", len(candidates)))
			}
			t.checkCandidates(ctx, candidates)

		case <-ctx.Done():
			t.logger.Info("deployment trigger has been stopped")
			return nil
		}
	}
}

func (t *Trigger) checkCandidates(ctx context.Context, cs []candidate) (err error) {
	// Group candidates by repository to reduce the number of Git operations on each repo.
	csm := make(map[string][]candidate)
	for _, c := range cs {
		repoId := c.application.GitPath.Repo.Id
		if _, ok := csm[repoId]; !ok {
			csm[repoId] = []candidate{c}
			continue
		}
		csm[repoId] = append(csm[repoId], c)
	}

	// Iterate each repository and check its candidates.
	// Only the last error will be returned.
	for repoID, cs := range csm {
		if e := t.checkRepoCandidates(ctx, repoID, cs); e != nil {
			t.logger.Error(fmt.Sprintf("failed while checking applications in repo %s", repoID), zap.Error(e))
			err = e
		}
	}
	return
}

func (t *Trigger) checkRepoCandidates(ctx context.Context, repoID string, cs []candidate) error {
	gitRepo, branch, headCommit, err := t.updateRepoToLatest(ctx, repoID)
	if err != nil {
		// TODO: Find a better way to skip the CANCELLED error log while shutting down.
		if ctx.Err() != context.Canceled {
			t.logger.Error(fmt.Sprintf("failed to update git repository %s to latest", repoID), zap.Error(err))
		}
		return err
	}

	ds := &determiners{
		onCommand:   NewOnCommandDeterminer(),
		onOutOfSync: NewOnOutOfSyncDeterminer(t.apiClient),
		onCommit:    NewOnCommitDeterminer(gitRepo, headCommit.Hash, t.commitStore, t.logger),
		onChain:     NewOnChainDeterminer(),
	}
	triggered := make(map[string]struct{})

	for _, c := range cs {
		app := c.application

		// Avoid triggering multiple deployments for the same application in the same iteration.
		if _, ok := triggered[app.Id]; ok {
			continue
		}

		appCfg, err := config.LoadApplication(gitRepo.GetPath(), app.GitPath.GetApplicationConfigFilePath(), app.Kind)
		if err != nil {
			t.logger.Error("failed to load application config file",
				zap.String("app", app.Name),
				zap.String("app-id", app.Id),
				zap.String("commit", headCommit.Hash),
				zap.Error(err),
			)

			// Set ApplicationSyncState to INVALID_CONFIG when LoadApplication fails.
			req := &pipedservice.ReportApplicationSyncStateRequest{
				ApplicationId: app.Id,
				State: &model.ApplicationSyncState{
					Status:      model.ApplicationSyncStatus_INVALID_CONFIG,
					ShortReason: "failed to load application configuration",
					Reason:      err.Error(),
					Timestamp:   time.Now().Unix(),
				},
			}
			_, err := t.apiClient.ReportApplicationSyncState(ctx, req)
			if err != nil {
				msg := fmt.Sprintf("failed to report application sync state %s: %v", app.Id, err)
				t.logger.Error(msg, zap.Error(err))
			}

			continue
		}

		shouldTrigger, err := ds.Determiner(c.kind).ShouldTrigger(ctx, app, appCfg)
		if err != nil {
			msg := fmt.Sprintf("failed while determining whether application %s should be triggered or not: %s", app.Name, err)
			t.notifyDeploymentTriggerFailed(app, appCfg, msg, headCommit)
			t.logger.Error(msg, zap.Error(err))
			continue
		}

		if !shouldTrigger {
			t.commitStore.Put(app.Id, headCommit.Hash)
			continue
		}

		var (
			commander                 string
			strategy                  model.SyncStrategy
			strategySummary           string
			deploymentChainID         string
			deploymentChainBlockIndex uint32
		)

		switch c.kind {
		case model.TriggerKind_ON_COMMAND:
			strategy = c.command.GetSyncApplication().SyncStrategy
			commander = c.command.Commander
			if strategy == model.SyncStrategy_QUICK_SYNC {
				strategySummary = "Quick sync because piped received a command from user via web console or pipectl"
			} else {
				strategySummary = "Sync with the specified pipeline because piped received a command from user via web console or pipectl"
			}

		case model.TriggerKind_ON_CHAIN:
			strategy = c.command.GetChainSyncApplication().SyncStrategy
			commander = c.command.Commander
			strategySummary = "Sync application in chain"
			deploymentChainID = c.command.GetChainSyncApplication().DeploymentChainId
			deploymentChainBlockIndex = c.command.GetChainSyncApplication().BlockIndex

		case model.TriggerKind_ON_OUT_OF_SYNC:
			strategy = model.SyncStrategy_QUICK_SYNC
			strategySummary = "Quick sync to attempt to resolve the detected configuration drift"

		default:
			strategy = model.SyncStrategy_AUTO
		}

		// Build the deployment to trigger.
		deployment, err := buildDeployment(
			app,
			branch,
			headCommit,
			commander,
			strategy,
			strategySummary,
			time.Now(),
			appCfg.DeploymentNotification,
			deploymentChainID,
			deploymentChainBlockIndex,
		)
		if err != nil {
			msg := fmt.Sprintf("failed to build deployment for application %s: %v", app.Id, err)
			t.notifyDeploymentTriggerFailed(app, appCfg, msg, headCommit)
			t.logger.Error(msg, zap.Error(err))
			continue
		}

		// In case the triggered deployment is of application that can trigger a deployment chain
		// create a new deployment chain with its configuration besides with the first deployment
		// in that chain.
		if appCfg.PostSync != nil && appCfg.PostSync.DeploymentChain != nil {
			if err := t.triggerDeploymentChain(ctx, appCfg.PostSync.DeploymentChain, deployment); err != nil {
				msg := fmt.Sprintf("failed to trigger application %s and its deployment chain: %v", app.Id, err)
				t.notifyDeploymentTriggerFailed(app, appCfg, msg, headCommit)
				t.logger.Error(msg, zap.Error(err))
				continue
			}
		} else {
			// Send a request to API to create a new deployment.
			if err := t.triggerDeployment(ctx, deployment); err != nil {
				msg := fmt.Sprintf("failed to trigger application %s: %v", app.Id, err)
				t.notifyDeploymentTriggerFailed(app, appCfg, msg, headCommit)
				t.logger.Error(msg, zap.Error(err))
				continue
			}
		}

		// TODO: Find a better way to ensure that the application should be updated correctly
		// when the deployment was successfully triggered.
		// This error is ignored because the deployment was already registered successfully.
		if e := reportMostRecentlyTriggeredDeployment(ctx, t.apiClient, deployment); e != nil {
			t.logger.Error("failed to report most recently triggered deployment", zap.Error(e))
		}

		triggered[app.Id] = struct{}{}
		t.commitStore.Put(app.Id, headCommit.Hash)
		t.notifyDeploymentTriggered(ctx, appCfg, deployment)

		// Mask command as handled since the deployment has been triggered successfully.
		if c.HasCommand() {
			metadata := map[string]string{
				model.MetadataKeyTriggeredDeploymentID: deployment.Id,
			}
			if err := c.command.Report(ctx, model.CommandStatus_COMMAND_SUCCEEDED, metadata, nil); err != nil {
				t.logger.Error("failed to report command status", zap.Error(err))
			}
		}
	}

	return nil
}

// listCommandCandidates finds all applications that have been commanded to sync.
func (t *Trigger) listCommandCandidates() []candidate {
	var (
		cmds = t.commandLister.ListApplicationCommands()
		apps = make([]candidate, 0)
	)

	for _, cmd := range cmds {
		// Prepare to handle SYNC_APPLICATION command.
		if cmd.IsSyncApplicationCmd() {
			// Find the target application specified in command.
			app, ok := t.applicationLister.Get(cmd.ApplicationId)
			if !ok {
				t.logger.Warn("detected an AppSync command for an unregistered application",
					zap.String("command", cmd.Id),
					zap.String("app-id", cmd.ApplicationId),
					zap.String("commander", cmd.Commander),
				)
				continue
			}

			apps = append(apps, candidate{
				application: app,
				kind:        model.TriggerKind_ON_COMMAND,
				command:     cmd,
			})
		}

		// Prepare to handle CHAIN_SYNC_APPLICATION command.
		if cmd.IsChainSyncApplicationCmd() {
			// Find the target application specified in command.
			app, ok := t.applicationLister.Get(cmd.ApplicationId)
			if !ok {
				t.logger.Warn("detected an InChainAppSync command for an unregistered application",
					zap.String("command", cmd.Id),
					zap.String("app-id", cmd.ApplicationId),
					zap.String("commander", cmd.Commander),
				)
				continue
			}

			apps = append(apps, candidate{
				application: app,
				kind:        model.TriggerKind_ON_CHAIN,
				command:     cmd,
			})
		}
	}

	return apps
}

// listOutOfSyncCandidates finds all applications that are staying at OUT_OF_SYNC state.
func (t *Trigger) listOutOfSyncCandidates() []candidate {
	var (
		list = t.applicationLister.List()
		apps = make([]candidate, 0)
	)
	for _, app := range list {
		if !app.IsOutOfSync() {
			continue
		}
		apps = append(apps, candidate{
			application: app,
			kind:        model.TriggerKind_ON_OUT_OF_SYNC,
		})
	}
	return apps
}

// listCommitCandidates finds all applications that have potentiality
// to be candidates by the changes of new commits.
// They are all applications managed by this Piped.
func (t *Trigger) listCommitCandidates() []candidate {
	var (
		list = t.applicationLister.List()
		apps = make([]candidate, 0)
	)
	for _, app := range list {
		apps = append(apps, candidate{
			application: app,
			kind:        model.TriggerKind_ON_COMMIT,
		})
	}
	return apps
}

// updateRepoToLatest ensures that the local data of the given Git repository should be up-to-date.
func (t *Trigger) updateRepoToLatest(ctx context.Context, repoID string) (repo git.Repo, branch string, headCommit git.Commit, err error) {
	var ok bool

	// Find the repository from the previously loaded list.
	repo, ok = t.gitRepos[repoID]
	if !ok {
		err = fmt.Errorf("the repository was not registered in Piped configuration")
		return
	}
	branch = repo.GetClonedBranch()

	// Fetch to update the repository.
	err = repo.Pull(ctx, branch)
	if err != nil {
		return
	}

	// Get the head commit of the repository.
	headCommit, err = repo.GetLatestCommit(ctx)
	return
}

func (t *Trigger) GetLastTriggeredCommitGetter() LastTriggeredCommitGetter {
	return t.commitStore
}

func (t *Trigger) notifyDeploymentTriggered(ctx context.Context, appCfg *config.GenericApplicationSpec, d *model.Deployment) {
	var mentions []string
	if n := appCfg.DeploymentNotification; n != nil {
		mentions = n.FindSlackAccounts(model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGERED)
	}

	t.notifier.Notify(model.NotificationEvent{
		Type: model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGERED,
		Metadata: &model.NotificationEventDeploymentTriggered{
			Deployment:        d,
			MentionedAccounts: mentions,
		},
	})
}

func (t *Trigger) notifyDeploymentTriggerFailed(app *model.Application, appCfg *config.GenericApplicationSpec, reason string, commit git.Commit) {
	var mentions []string
	if n := appCfg.DeploymentNotification; n != nil {
		mentions = n.FindSlackAccounts(model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGER_FAILED)
	}

	t.notifier.Notify(model.NotificationEvent{
		Type: model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGER_FAILED,
		Metadata: &model.NotificationEventDeploymentTriggerFailed{
			Application:       app,
			CommitHash:        commit.Hash,
			MentionedAccounts: mentions,
			CommitMessage:     commit.Message,
			Reason:            reason,
		},
	})
}
