// Copyright 2024 The PipeCD Authors.
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
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/piped/livestatestore/ecs"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/ecs"
	"github.com/pipe-cd/pipecd/pkg/app/piped/sourceprocesser"
	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/diff"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type applicationLister interface {
	ListByPlatformProvider(name string) []*model.Application
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
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

	gitRepos map[string]git.Repo
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
		zap.String("platform-provider", cp.Name),
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
		logger:            logger,
	}
}

func (d *detector) Run(ctx context.Context) error {
	d.logger.Info("start running drift detector for ecs applications")

	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			d.logger.Info("drift detector for ecs applications has been stopped")
			return nil

		case <-ticker.C:
			d.check(ctx)
		}
	}
}

func (d *detector) ProviderName() string {
	return d.provider.Name
}

func (d *detector) check(ctx context.Context) {
	appsByRepo := d.listGroupedApplication()

	for repoID, apps := range appsByRepo {
		gitRepo, ok := d.gitRepos[repoID]
		if !ok {
			// Clone repository for the first time.
			gr, err := d.cloneGitRepository(ctx, repoID)
			if err != nil {
				d.logger.Error("failed to clone git repository",
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
}

func (d *detector) cloneGitRepository(ctx context.Context, repoID string) (git.Repo, error) {
	repoCfg, ok := d.config.GetRepository(repoID)
	if !ok {
		return nil, fmt.Errorf("repository %s was not found in piped configuration", repoID)
	}
	return d.gitClient.Clone(ctx, repoID, repoCfg.Remote, repoCfg.Branch, "")
}

// listGroupedApplication retrieves all applications those should be handled by this director
// and then groups them by repoID.
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

func (d *detector) checkApplication(ctx context.Context, app *model.Application, repo git.Repo, headCommit git.Commit) error {
	headManifests, err := d.loadConfigs(app, repo, headCommit)
	if err != nil {
		return err
	}
	d.logger.Info(fmt.Sprintf("application %s has ecs manifests at commit %s", app.Id, headCommit.Hash))

	liveManifests, ok := d.stateGetter.GetECSManifests(app.Id)
	if !ok {
		return fmt.Errorf("failed to get live ecs manifests")
	}
	d.logger.Info(fmt.Sprintf("application %s has live ecs manifests", app.Id))

	result, err := provider.Diff(
		liveManifests,
		headManifests,
		diff.WithEquateEmpty(),
		diff.WithIgnoreAddingMapKeys(),
		diff.WithCompareNumberAndNumericString(),
	)
	if err != nil {
		return err
	}

	state := makeSyncState(result, headCommit.Hash)

	return d.reporter.ReportApplicationSyncState(ctx, app.Id, state)
}

func (d *detector) loadConfigs(app *model.Application, repo git.Repo, headCommit git.Commit) (provider.ECSManifests, error) {
	var (
		manifestCache = provider.ECSManifestsCache{
			AppID:  app.Id,
			Cache:  d.appManifestsCache,
			Logger: d.logger,
		}
		repoDir = repo.GetPath()
		appDir  = filepath.Join(repoDir, app.GitPath.Path)
	)

	manifest, ok := manifestCache.Get(headCommit.Hash)
	if ok {
		return manifest, nil
	}
	// When the manifests were not in the cache we have to load them.
	cfg, err := d.loadApplicationConfiguration(repoDir, app)
	if err != nil {
		return provider.ECSManifests{}, fmt.Errorf("failed to load application configuration: %w", err)
	}

	gds, ok := cfg.GetGenericApplication()
	if !ok {
		return provider.ECSManifests{}, fmt.Errorf("unsupported application kind %s", cfg.Kind)
	}

	var (
		encryptionUsed = d.secretDecrypter != nil && gds.Encryption != nil
		attachmentUsed = gds.Attachment != nil
	)

	// We have to copy repository into another directory because
	// decrypting the sealed secrets or attaching files might change the git repository.
	if attachmentUsed || encryptionUsed {
		dir, err := os.MkdirTemp("", "detector-git-processing")
		if err != nil {
			return provider.ECSManifests{}, fmt.Errorf("failed to prepare a temporary directory for git repository (%w)", err)
		}
		defer os.RemoveAll(dir)

		repo, err = repo.Copy(filepath.Join(dir, "repo"))
		if err != nil {
			return provider.ECSManifests{}, fmt.Errorf("failed to copy the cloned git repository (%w)", err)
		}
		repoDir := repo.GetPath()
		appDir = filepath.Join(repoDir, app.GitPath.Path)
	}

	var templProcessors []sourceprocesser.SourceTemplateProcessor
	// Decrypting secrets to manifests.
	if encryptionUsed {
		templProcessors = append(templProcessors, sourceprocesser.NewSecretDecrypterProcessor(gds.Encryption, d.secretDecrypter))
	}
	// Then attaching configurated files to manifests.
	if attachmentUsed {
		templProcessors = append(templProcessors, sourceprocesser.NewAttachmentProcessor(gds.Attachment))
	}
	if len(templProcessors) > 0 {
		sp := sourceprocesser.NewSourceProcessor(appDir, templProcessors...)
		if err := sp.Process(); err != nil {
			return provider.ECSManifests{}, fmt.Errorf("failed to process source files: %w", err)
		}
	}

	var serviceDefFile string
	var taskDefFile string
	if cfg.ECSApplicationSpec != nil {
		serviceDefFile = cfg.ECSApplicationSpec.Input.ServiceDefinitionFile
		taskDefFile = cfg.ECSApplicationSpec.Input.TaskDefinitionFile
	}
	serviceDef, err := provider.LoadServiceDefinition(appDir, serviceDefFile)
	if err != nil {
		return provider.ECSManifests{}, fmt.Errorf("failed to load new service definition: %w", err)
	}
	taskDef, err := provider.LoadTaskDefinition(appDir, taskDefFile)
	if err != nil {
		return provider.ECSManifests{}, fmt.Errorf("failed to load new task definition: %w", err)
	}

	manifests := provider.ECSManifests{
		ServiceDefinition: &serviceDef,
		TaskDefinition:    &taskDef,
	}
	manifestCache.Put(headCommit.Hash, manifests)

	return manifest, nil
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

func makeSyncState(r *provider.DiffResult, commit string) model.ApplicationSyncState {
	if r.NoChange() {
		return model.ApplicationSyncState{
			Status:    model.ApplicationSyncStatus_SYNCED,
			Timestamp: time.Now().Unix(),
		}
	}

	if onlyDesiredCountChanged(r) {
		return model.ApplicationSyncState{
			Status:      model.ApplicationSyncStatus_SYNCED,
			ShortReason: "Ignore diff of `desiredCount`.",
			Reason:      "`desiredCount` is 0 or not defined in your config (which means ignoring updating desiredCount) and only `desiredCount` is changed.",
			Timestamp:   time.Now().Unix(),
		}
	}

	shortReason := "The ecs manifests is not synced"
	if len(commit) >= 7 {
		commit = commit[:7]
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("Diff between the defined state in Git at commit %s and actual live state:\n\n", commit))
	b.WriteString("--- Actual   (LiveState)\n+++ Expected (Git)\n\n")

	details := r.Render(provider.DiffRenderOptions{
		// Currently, we do not use the diff command to render the result
		// because ECS adds a large number of default values to the
		// running manifest that causes a wrong diff text.
		UseDiffCommand: false,
	})
	b.WriteString(details)

	return model.ApplicationSyncState{
		Status:      model.ApplicationSyncStatus_OUT_OF_SYNC,
		ShortReason: shortReason,
		Reason:      b.String(),
		Timestamp:   time.Now().Unix(),
	}
}

func onlyDesiredCountChanged(r *provider.DiffResult) bool {
	return r.Diff.NumNodes() == 1 && r.Old.ServiceDefinition.DesiredCount != r.New.ServiceDefinition.DesiredCount
}
