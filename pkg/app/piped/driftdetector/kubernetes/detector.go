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
	"time"

	"go.uber.org/zap"

	provider "github.com/kapetaniosci/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/kapetaniosci/pipe/pkg/app/piped/livestatestore/kubernetes"
	"github.com/kapetaniosci/pipe/pkg/cache"
	"github.com/kapetaniosci/pipe/pkg/config"
	"github.com/kapetaniosci/pipe/pkg/git"
	"github.com/kapetaniosci/pipe/pkg/model"
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
			state := makeDeployingState(app, headDeployment)
			d.reporter.ReportApplicationSyncState(ctx, app.Id, state)
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
	var (
		manifestCache = provider.AppManifestsCache{
			AppID:  app.Id,
			Cache:  d.appManifestsCache,
			Logger: d.logger,
		}
		repoDir = repo.GetPath()
		appDir  = filepath.Join(repoDir, app.GitPath.Path)
	)

	headManifests, ok := manifestCache.Get(headCommit.Hash)
	if !ok {
		// When the manifests were not in the cache we have to load them.
		cfg, err := d.loadDeploymentConfiguration(repoDir, app)
		if err != nil {
			err = fmt.Errorf("failed to load deployment configuration: %w", err)
			return err
		}
		loader := provider.NewManifestLoader(appDir, repoDir, cfg.KubernetesDeploymentSpec.Input, d.logger)
		headManifests, err = loader.LoadManifests(ctx)
		if err != nil {
			err = fmt.Errorf("failed to load new manifests: %w", err)
			return err
		}
		manifestCache.Put(headCommit.Hash, headManifests)
	}
	d.logger.Info(fmt.Sprintf("application %s has %d manifests at commit %s", app.Id, len(headManifests), headCommit.Hash))

	liveManifests := d.stateGetter.GetAppLiveManifests(app.Id)
	d.logger.Info(fmt.Sprintf("application %s has %d live manifests", app.Id, len(liveManifests)))

	//watchingResourceKinds := d.stateGetter.GetWatchingResourceKinds()
	return nil
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

func makeDeployingState(app *model.Application, deployment *model.Deployment) model.ApplicationSyncState {
	reason := deployment.Description
	if reason == "" {
		reason = "A deployment is running."
	}
	return model.ApplicationSyncState{
		Status:           model.ApplicationSyncStatus_DEPLOYING,
		ShortReason:      reason,
		Reason:           reason,
		HeadDeploymentId: deployment.Id,
		Timestamp:        time.Now().Unix(),
	}
}
