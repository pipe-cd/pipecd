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

package kubernetes

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"go.uber.org/zap"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/pipe-cd/pipe/pkg/app/piped/livestatestore/kubernetes"
	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/git"
	"github.com/pipe-cd/pipe/pkg/model"
)

type applicationLister interface {
	ListByCloudProvider(name string) []*model.Application
}

type deploymentLister interface {
	ListAppHeadDeployments() map[string]*model.Deployment
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

type reporter interface {
	ReportApplicationSyncState(ctx context.Context, appID string, state model.ApplicationSyncState) error
}

type detector struct {
	provider          config.PipedCloudProvider
	appLister         applicationLister
	deploymentLister  deploymentLister
	gitClient         gitClient
	stateGetter       kubernetes.Getter
	reporter          reporter
	appManifestsCache cache.Cache
	interval          time.Duration
	config            *config.PipedSpec
	logger            *zap.Logger

	gitRepos   map[string]git.Repo
	syncStates map[string]model.ApplicationSyncState
}

func NewDetector(
	cp config.PipedCloudProvider,
	appLister applicationLister,
	deploymentLister deploymentLister,
	gitClient gitClient,
	stateGetter kubernetes.Getter,
	reporter reporter,
	appManifestsCache cache.Cache,
	cfg *config.PipedSpec,
	logger *zap.Logger,
) *detector {

	logger = logger.Named("kubernetes-detector").With(
		zap.String("cloud-provider", cp.Name),
	)
	return &detector{
		provider:          cp,
		appLister:         appLister,
		deploymentLister:  deploymentLister,
		gitClient:         gitClient,
		stateGetter:       stateGetter,
		reporter:          reporter,
		appManifestsCache: appManifestsCache,
		interval:          time.Minute,
		config:            cfg,
		gitRepos:          make(map[string]git.Repo),
		syncStates:        make(map[string]model.ApplicationSyncState),
		logger:            logger,
	}
}

func (d *detector) Run(ctx context.Context) error {
	d.logger.Info("start running drift detector for kubernetes applications")

	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()

L:
	for {
		select {
		case <-ticker.C:
			d.check(ctx)

		case <-ctx.Done():
			break L
		}
	}

	d.logger.Info("drift detector for kubernetes applications has been stopped")
	return nil
}

func (d *detector) check(ctx context.Context) error {
	var (
		err             error
		applications    = d.listApplications()
		headDeployments = d.deploymentLister.ListAppHeadDeployments()
	)

	for repoID, apps := range applications {
		var notDeployingApps []*model.Application

		// Firstly, handle all deploying applications
		// and remove them from the list.
		for _, app := range apps {
			headDeployment, ok := headDeployments[app.Id]
			if !ok {
				notDeployingApps = append(notDeployingApps, app)
				continue
			}
			state := makeDeployingState(headDeployment)
			if err := d.reporter.ReportApplicationSyncState(ctx, app.Id, state); err != nil {
				d.logger.Error("failed to report application sync state", zap.Error(err))
			}
		}

		if len(notDeployingApps) == 0 {
			continue
		}

		// Next, we have to clone the lastest commit of repository
		// to compare the states.
		gitRepo, ok := d.gitRepos[repoID]
		if !ok {
			// Clone repository for the first time.
			repoCfg, ok := d.config.GetRepository(repoID)
			if !ok {
				d.logger.Error(fmt.Sprintf("repository %s was not found in piped configuration", repoID))
				continue
			}
			gitRepo, err = d.gitClient.Clone(ctx, repoID, repoCfg.Remote, repoCfg.Branch, "")
			if err != nil {
				d.logger.Error("failed to clone repository",
					zap.String("repo-id", repoID),
					zap.Error(err),
				)
				continue
			}
			d.gitRepos[repoID] = gitRepo
		}

		// Fetch to update the repository.
		branch := gitRepo.GetClonedBranch()
		if err := gitRepo.Pull(ctx, branch); err != nil {
			d.logger.Error("failed to update repository branch",
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

		for _, app := range notDeployingApps {
			if err := d.checkApplication(ctx, app, gitRepo, headCommit); err != nil {
				d.logger.Error(fmt.Sprintf("failed to check application: %s", app.Id), zap.Error(err))
			}
		}
	}

	return nil
}

func (d *detector) checkApplication(ctx context.Context, app *model.Application, repo git.Repo, headCommit git.Commit) error {
	watchingResourceKinds := d.stateGetter.GetWatchingResourceKinds()
	headManifests, err := d.loadHeadManifests(ctx, app, repo, headCommit, watchingResourceKinds)
	if err != nil {
		return err
	}
	headManifests = filterIgnoringManifests(headManifests)
	sort.Slice(headManifests, func(i, j int) bool {
		return headManifests[i].Key.IsLessWithIgnoringNamespace(headManifests[j].Key)
	})
	d.logger.Info(fmt.Sprintf("application %s has %d manifests at commit %s", app.Id, len(headManifests), headCommit.Hash))

	liveManifests := d.stateGetter.GetAppLiveManifests(app.Id)
	sort.Slice(liveManifests, func(i, j int) bool {
		return liveManifests[i].Key.IsLessWithIgnoringNamespace(liveManifests[j].Key)
	})
	liveManifests = filterIgnoringManifests(liveManifests)
	d.logger.Info(fmt.Sprintf("application %s has %d live manifests", app.Id, len(liveManifests)))

	var missings, redundancies []provider.Manifest
	var h, l int
	for {
		if h >= len(headManifests) || l >= len(liveManifests) {
			break
		}
		if headManifests[h].Key.IsEqualWithIgnoringNamespace(liveManifests[l].Key) {
			h++
			l++
			continue
		}
		// Has in head but not in live so this should be a missing one.
		if headManifests[h].Key.IsLessWithIgnoringNamespace(liveManifests[l].Key) {
			missings = append(missings, headManifests[h])
			h++
			continue
		}
		// Has in live but not in head so this should be a redundant one.
		redundancies = append(redundancies, liveManifests[l])
		l++
	}
	if len(headManifests) > h {
		missings = append(missings, headManifests[h:]...)
	}
	if len(liveManifests) > h {
		redundancies = append(redundancies, liveManifests[l:]...)
	}
	if len(missings)+len(redundancies) > 0 {
		state := makeOutOfSyncStateBecauseMissingOrRedundant(missings, redundancies)
		return d.reporter.ReportApplicationSyncState(ctx, app.Id, state)
	}

	// All manifest keys are matched. Now we will go to check the diff of each manifest pair.
	var (
		diffs       = make(map[int]provider.DiffResultList)
		diffOptions = []provider.DiffOption{
			provider.WithDiffIgnoreMissingFields(),
		}
		secretDiffOptions = []provider.DiffOption{
			provider.WithDiffIgnoreMissingFields(),
			provider.WithDiffRedactPathPrefix("data", "secret-data-in-configuration", "secret-data-in-cluster"),
		}
		configMapDiffOptions = []provider.DiffOption{
			provider.WithDiffIgnoreMissingFields(),
			provider.WithDiffRedactPathPrefix("data", "configmap-data-in-configuration", "configmap-data-in-cluster"),
		}
	)

	for i := 0; i < len(headManifests); i++ {
		var options []provider.DiffOption
		if headManifests[i].Key.IsSecret() {
			options = secretDiffOptions
		} else if headManifests[i].Key.IsConfigMap() {
			options = configMapDiffOptions
		} else {
			options = diffOptions
		}

		result := provider.Diff(headManifests[i], liveManifests[i], options...)
		if len(result) > 0 {
			diffs[i] = result
		}
	}

	// No diffs means this application is in SYNCED state.
	if len(diffs) == 0 {
		state := makeSyncedState()
		return d.reporter.ReportApplicationSyncState(ctx, app.Id, state)
	}

	state := makeOutOfSyncState(headManifests, liveManifests, diffs)
	return d.reporter.ReportApplicationSyncState(ctx, app.Id, state)
}

func (d *detector) loadHeadManifests(ctx context.Context, app *model.Application, repo git.Repo, headCommit git.Commit, watchingResourceKinds []provider.APIVersionKind) ([]provider.Manifest, error) {
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
		// When the manifests were not in the cache we have to load them.
		cfg, err := d.loadDeploymentConfiguration(repoDir, app)
		if err != nil {
			err = fmt.Errorf("failed to load deployment configuration: %w", err)
			return nil, err
		}
		loader := provider.NewManifestLoader(app.Name, appDir, repoDir, cfg.KubernetesDeploymentSpec.Input, d.logger)
		manifests, err = loader.LoadManifests(ctx)
		if err != nil {
			err = fmt.Errorf("failed to load new manifests: %w", err)
			return nil, err
		}
		manifestCache.Put(headCommit.Hash, manifests)
	}

	watchingMap := make(map[provider.APIVersionKind]struct{}, len(watchingResourceKinds))
	for _, k := range watchingResourceKinds {
		watchingMap[k] = struct{}{}
	}

	filtered := make([]provider.Manifest, 0, len(manifests))
	for _, m := range manifests {
		_, ok := watchingMap[provider.APIVersionKind{
			APIVersion: m.Key.APIVersion,
			Kind:       m.Key.Kind,
		}]
		if ok {
			filtered = append(filtered, m)
		}
	}

	return filtered, nil
}

// listApplications retrieves all applications those should be handled by this director
// and then groups them by repoID.
func (d *detector) listApplications() map[string][]*model.Application {
	var (
		apps = d.appLister.ListByCloudProvider(d.provider.Name)
		m    = make(map[string][]*model.Application)
	)
	for _, app := range apps {
		repoID := app.GitPath.RepoId
		if _, ok := m[repoID]; !ok {
			m[repoID] = []*model.Application{app}
		} else {
			m[repoID] = append(m[repoID], app)
		}
	}
	return m
}

func (d *detector) loadDeploymentConfiguration(repoPath string, app *model.Application) (*config.Config, error) {
	path := filepath.Join(repoPath, app.GitPath.GetDeploymentConfigFilePath(config.DeploymentConfigurationFileName))
	cfg, err := config.LoadFromYAML(path)
	if err != nil {
		return nil, err
	}
	if appKind, ok := config.ToApplicationKind(cfg.Kind); !ok || appKind != app.Kind {
		return nil, fmt.Errorf("application in deployment configuration file is not match, got: %s, expected: %s", appKind, app.Kind)
	}
	return cfg, nil
}

func (d *detector) ProviderName() string {
	return d.provider.Name
}

func makeDeployingState(deployment *model.Deployment) model.ApplicationSyncState {
	var (
		shortReason = "A deployment of this application is running"
		reason      = deployment.Summary
	)
	if reason == "" {
		reason = shortReason
	}
	return model.ApplicationSyncState{
		Status:           model.ApplicationSyncStatus_DEPLOYING,
		ShortReason:      shortReason,
		Reason:           reason,
		HeadDeploymentId: deployment.Id,
		Timestamp:        time.Now().Unix(),
	}
}

func makeSyncedState() model.ApplicationSyncState {
	return model.ApplicationSyncState{
		Status:      model.ApplicationSyncStatus_SYNCED,
		ShortReason: "",
		Reason:      "",
		Timestamp:   time.Now().Unix(),
	}
}

func makeOutOfSyncStateBecauseMissingOrRedundant(missings, redundancies []provider.Manifest) model.ApplicationSyncState {
	shortReason := fmt.Sprintf("There are %d missing manifests and %d redundant manifests.", len(missings), len(redundancies))
	var b strings.Builder
	if len(missings) > 0 {
		b.WriteString(fmt.Sprintf("The following %d manifests are defined in Git, but NOT appearing in the cluster:\n", len(missings)))
		for _, m := range missings {
			b.WriteString(fmt.Sprintf("- %s\n", m.Key.ReadableString()))
		}
	}
	if len(redundancies) > 0 {
		b.WriteString(fmt.Sprintf("The following %d manifests are NOT defined in Git, but appearing in the cluster:\n", len(redundancies)))
		for _, m := range redundancies {
			b.WriteString(fmt.Sprintf("- %s\n", m.Key.ReadableString()))
		}
	}
	return model.ApplicationSyncState{
		Status:      model.ApplicationSyncStatus_OUT_OF_SYNC,
		ShortReason: shortReason,
		Reason:      b.String(),
		Timestamp:   time.Now().Unix(),
	}
}

func makeOutOfSyncState(headManifests, liveManifests []provider.Manifest, diffs map[int]provider.DiffResultList) model.ApplicationSyncState {
	const maxPrintDiffs = 3

	shortReason := fmt.Sprintf("There are %d manifests are not synced.", len(diffs))
	var b strings.Builder
	b.WriteString(fmt.Sprintf("There are %d manifests are not synced:\n", len(diffs)))

	var prints = 0
	for i, list := range diffs {
		b.WriteString(fmt.Sprintf("%d. %s\n", i+1, headManifests[i].Key.ReadableString()))
		b.WriteString(list.String())
		b.WriteString("\n")
		prints++
		if prints >= maxPrintDiffs {
			break
		}
	}

	if prints < len(diffs) {
		b.WriteString(fmt.Sprintf("... (diffs from %d other manifests are omitted\n", len(diffs)-prints))
	}

	return model.ApplicationSyncState{
		Status:      model.ApplicationSyncStatus_OUT_OF_SYNC,
		ShortReason: shortReason,
		Reason:      b.String(),
		Timestamp:   time.Now().Unix(),
	}
}

func filterIgnoringManifests(manifests []provider.Manifest) []provider.Manifest {
	out := make([]provider.Manifest, 0, len(manifests))
	for _, m := range manifests {
		annotations := m.GetAnnotations()
		if annotations[provider.LabelIgnoreDriftDirection] == provider.IgnoreDriftDetectionTrue {
			continue
		}
		out = append(out, m)
	}
	return out
}
