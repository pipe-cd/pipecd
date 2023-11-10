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

package kubernetes

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/piped/livestatestore/kubernetes"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/kubernetes"
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
	stateGetter       kubernetes.Getter
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
	stateGetter kubernetes.Getter,
	reporter reporter,
	appManifestsCache cache.Cache,
	cfg *config.PipedSpec,
	sd secretDecrypter,
	logger *zap.Logger,
) Detector {

	logger = logger.Named("kubernetes-detector").With(
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
		syncStates:        make(map[string]model.ApplicationSyncState),
		logger:            logger,
	}
}

func (d *detector) Run(ctx context.Context) error {
	d.logger.Info("start running drift detector for kubernetes applications")

	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			d.check(ctx)

		case <-ctx.Done():
			d.logger.Info("drift detector for kubernetes applications has been stopped")
			return nil
		}
	}
}

func (d *detector) check(ctx context.Context) {
	appsByRepo := d.listGroupedApplication()

	for repoID, apps := range appsByRepo {
		gitRepo, ok := d.gitRepos[repoID]
		if !ok {
			// Clone repository for the first time.
			repoCfg, ok := d.config.GetRepository(repoID)
			if !ok {
				d.logger.Error(fmt.Sprintf("repository %s was not found in piped configuration", repoID))
				continue
			}
			gr, err := d.gitClient.Clone(ctx, repoID, repoCfg.Remote, repoCfg.Branch, "")
			if err != nil {
				d.logger.Error("failed to clone repository",
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

		// Start checking all applications in this repository.
		for _, app := range apps {
			if err := d.checkApplication(ctx, app, gitRepo, headCommit); err != nil {
				d.logger.Error(fmt.Sprintf("failed to check application: %s", app.Id), zap.Error(err))
			}
		}
	}
}

func (d *detector) checkApplication(ctx context.Context, app *model.Application, repo git.Repo, headCommit git.Commit) error {
	watchingResourceKinds := d.stateGetter.GetWatchingResourceKinds()
	headManifests, err := d.loadHeadManifests(ctx, app, repo, headCommit, watchingResourceKinds)
	if err != nil {
		return err
	}
	headManifests = filterIgnoringManifests(headManifests)
	d.logger.Debug(fmt.Sprintf("application %s has %d manifests at commit %s", app.Id, len(headManifests), headCommit.Hash))

	liveManifests := d.stateGetter.GetAppLiveManifests(app.Id)
	liveManifests = filterIgnoringManifests(liveManifests)
	d.logger.Debug(fmt.Sprintf("application %s has %d live manifests", app.Id, len(liveManifests)))

	ddCfg, err := d.getDriftDetectionConfig(repo.GetPath(), app)
	if err != nil {
		return err
	}

	ignoreConfig := make(map[string][]string, 0)
	if ddCfg != nil {
		for _, ignoreField := range ddCfg.IgnoreFields {
			// ignoreField is 'apiVersion:kind:namespace:name#fieldPath'
			splited := strings.Split(ignoreField, "#")
			key, ignoredPath := splited[0], splited[1]
			ignoreConfig[key] = append(ignoreConfig[key], ignoredPath)
		}
	}

	result, err := provider.DiffList(
		liveManifests,
		headManifests,
		d.logger,
		diff.WithEquateEmpty(),
		diff.WithIgnoreAddingMapKeys(),
		diff.WithCompareNumberAndNumericString(),
		diff.WithIgnoreConfig(ignoreConfig),
	)
	if err != nil {
		return err
	}

	state := makeSyncState(result, headCommit.Hash)

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
		cfg, err := d.loadApplicationConfiguration(repoDir, app)
		if err != nil {
			return nil, fmt.Errorf("failed to load application configuration: %w", err)
		}

		gds, ok := cfg.GetGenericApplication()
		if !ok {
			return nil, fmt.Errorf("unsupport application kind %s", cfg.Kind)
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
				return nil, fmt.Errorf("failed to prepare a temporary directory for git repository (%w)", err)
			}
			defer os.RemoveAll(dir)

			repo, err = repo.Copy(filepath.Join(dir, "repo"))
			if err != nil {
				return nil, fmt.Errorf("failed to copy the cloned git repository (%w)", err)
			}
			repoDir = repo.GetPath()
			appDir = filepath.Join(repoDir, app.GitPath.Path)
		}

		// Decrypting secrets to manifests.
		if encryptionUsed {
			if err := sourceprocesser.DecryptSecrets(appDir, *gds.Encryption, d.secretDecrypter); err != nil {
				return nil, fmt.Errorf("failed to decrypt secrets (%w)", err)
			}
		}
		// Then attaching configurated files to manifests.
		if attachmentUsed {
			if err := sourceprocesser.AttachData(appDir, *gds.Attachment); err != nil {
				return nil, fmt.Errorf("failed to attach files (%w)", err)
			}
		}

		loader := provider.NewLoader(app.Name, appDir, repoDir, app.GitPath.ConfigFilename, cfg.KubernetesApplicationSpec.Input, d.gitClient, d.logger)
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

// listGroupedApplication retrieves all applications those should be handled by this director
// and then groups them by repoID.
func (d *detector) listGroupedApplication() map[string][]*model.Application {
	var (
		apps = d.appLister.ListByPlatformProvider(d.provider.Name)
		m    = make(map[string][]*model.Application)
	)
	for _, app := range apps {
		repoID := app.GitPath.Repo.Id
		if _, ok := m[repoID]; !ok {
			m[repoID] = []*model.Application{app}
		} else {
			m[repoID] = append(m[repoID], app)
		}
	}
	return m
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

	if cfg.KubernetesApplicationSpec != nil && cfg.KubernetesApplicationSpec.Input.HelmChart != nil {
		chartRepoName := cfg.KubernetesApplicationSpec.Input.HelmChart.Repository
		if chartRepoName != "" {
			cfg.KubernetesApplicationSpec.Input.HelmChart.Insecure = d.config.IsInsecureChartRepository(chartRepoName)
		}
	}

	return cfg, nil
}

func (d *detector) ProviderName() string {
	return d.provider.Name
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

func makeSyncState(r *provider.DiffListResult, commit string) model.ApplicationSyncState {
	if r.NoChange() {
		return model.ApplicationSyncState{
			Status:      model.ApplicationSyncStatus_SYNCED,
			ShortReason: "",
			Reason:      "",
			Timestamp:   time.Now().Unix(),
		}
	}

	total := len(r.Adds) + len(r.Deletes) + len(r.Changes)
	shortReason := fmt.Sprintf("There are %d manifests not synced (%d adds, %d deletes, %d changes)", total, len(r.Adds), len(r.Deletes), len(r.Changes))
	if len(commit) >= 7 {
		commit = commit[:7]
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("Diff between the defined state in Git at commit %s and actual state in cluster:\n\n", commit))
	b.WriteString("--- Actual   (LiveState)\n+++ Expected (Git)\n\n")

	details := r.Render(provider.DiffRenderOptions{
		MaskSecret:          true,
		MaskConfigMap:       true,
		MaxChangedManifests: 3,
		// Currently, we do not use the diff command to render the result
		// because Kubernetes adds a large number of default values to the
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

func (d *detector) getDriftDetectionConfig(repoDir string, app *model.Application) (*config.DriftDetection, error) {
	cfg, err := d.loadApplicationConfiguration(repoDir, app)
	if err != nil {
		return nil, fmt.Errorf("failed to load application configuration: %w", err)
	}
	gds, ok := cfg.GetGenericApplication()
	if !ok {
		return nil, fmt.Errorf("unsupport application kind %s", cfg.Kind)
	}

	return gds.DriftDetection, nil
}
