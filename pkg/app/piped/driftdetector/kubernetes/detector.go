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
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"go.uber.org/zap"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/pipe-cd/pipe/pkg/app/piped/livestatestore/kubernetes"
	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/diff"
	"github.com/pipe-cd/pipe/pkg/git"
	"github.com/pipe-cd/pipe/pkg/model"
)

type applicationLister interface {
	ListByCloudProvider(name string) []*model.Application
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

type sealedSecretDecrypter interface {
	Decrypt(string) (string, error)
}

type reporter interface {
	ReportApplicationSyncState(ctx context.Context, appID string, state model.ApplicationSyncState) error
}

type detector struct {
	provider              config.PipedCloudProvider
	appLister             applicationLister
	gitClient             gitClient
	stateGetter           kubernetes.Getter
	reporter              reporter
	appManifestsCache     cache.Cache
	interval              time.Duration
	config                *config.PipedSpec
	sealedSecretDecrypter sealedSecretDecrypter
	logger                *zap.Logger

	gitRepos   map[string]git.Repo
	syncStates map[string]model.ApplicationSyncState
}

func NewDetector(
	cp config.PipedCloudProvider,
	appLister applicationLister,
	gitClient gitClient,
	stateGetter kubernetes.Getter,
	reporter reporter,
	appManifestsCache cache.Cache,
	cfg *config.PipedSpec,
	ssd sealedSecretDecrypter,
	logger *zap.Logger,
) *detector {

	logger = logger.Named("kubernetes-detector").With(
		zap.String("cloud-provider", cp.Name),
	)
	return &detector{
		provider:              cp,
		appLister:             appLister,
		gitClient:             gitClient,
		stateGetter:           stateGetter,
		reporter:              reporter,
		appManifestsCache:     appManifestsCache,
		interval:              time.Minute,
		config:                cfg,
		sealedSecretDecrypter: ssd,
		gitRepos:              make(map[string]git.Repo),
		syncStates:            make(map[string]model.ApplicationSyncState),
		logger:                logger,
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

	return nil
}

func (d *detector) checkApplication(ctx context.Context, app *model.Application, repo git.Repo, headCommit git.Commit) error {
	watchingResourceKinds := d.stateGetter.GetWatchingResourceKinds()
	headManifests, err := d.loadHeadManifests(ctx, app, repo, headCommit, watchingResourceKinds)
	if err != nil {
		return err
	}
	headManifests = filterIgnoringManifests(headManifests)
	d.logger.Info(fmt.Sprintf("application %s has %d manifests at commit %s", app.Id, len(headManifests), headCommit.Hash))

	liveManifests := d.stateGetter.GetAppLiveManifests(app.Id)
	liveManifests = filterIgnoringManifests(liveManifests)
	d.logger.Info(fmt.Sprintf("application %s has %d live manifests", app.Id, len(liveManifests)))

	// Divide manifests into separate groups.
	adds, deletes, headInters, liveInters := groupManifests(headManifests, liveManifests)

	// Now we will go to check the diff intersection group.
	changes := make(map[provider.Manifest]*diff.Result)
	for i := 0; i < len(headInters); i++ {
		result, err := provider.Diff(headInters[i], liveInters[i],
			diff.WithEquateEmpty(),
			diff.WithIgnoreAddingMapKeys(),
			diff.WithCompareNumberAndNumericString(),
		)
		if err != nil {
			d.logger.Error("failed to calculate the diff of manifests", zap.Error(err))
			return err
		}
		if !result.HasDiff() {
			continue
		}
		changes[headInters[i]] = result
	}

	// No diffs means this application is in SYNCED state.
	if len(adds) == 0 && len(deletes) == 0 && len(changes) == 0 {
		state := makeSyncedState()
		return d.reporter.ReportApplicationSyncState(ctx, app.Id, state)
	}

	state := makeOutOfSyncState(adds, deletes, changes, headCommit.Hash)
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
			return nil, fmt.Errorf("failed to load deployment configuration: %w", err)
		}

		gds, ok := cfg.GetGenericDeployment()
		if !ok {
			return nil, fmt.Errorf("unsupport application kind %s", cfg.Kind)
		}

		if d.sealedSecretDecrypter != nil && len(gds.SealedSecrets) > 0 {
			// We have to copy repository into another directory because
			// decrypting the sealed secrets might change the git repository.
			dir, err := ioutil.TempDir("", "detector-git-decrypt")
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

			if err := decryptSealedSecrets(appDir, gds.SealedSecrets, d.sealedSecretDecrypter); err != nil {
				return nil, fmt.Errorf("failed to decrypt sealed secrets (%w)", err)
			}
		}

		loader := provider.NewManifestLoader(app.Name, appDir, repoDir, app.GitPath.ConfigFilename, cfg.KubernetesDeploymentSpec.Input, d.logger)
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

func decryptSealedSecrets(appDir string, secrets []config.SealedSecretMapping, dcr sealedSecretDecrypter) error {
	for _, s := range secrets {
		secretPath := filepath.Join(appDir, s.Path)
		cfg, err := config.LoadFromYAML(secretPath)
		if err != nil {
			return fmt.Errorf("unable to read sealed secret file %s (%w)", s.Path, err)
		}
		if cfg.Kind != config.KindSealedSecret {
			return fmt.Errorf("unexpected kind in sealed secret file %s, want %q but got %q", s.Path, config.KindSealedSecret, cfg.Kind)
		}

		content, err := cfg.SealedSecretSpec.RenderOriginalContent(dcr)
		if err != nil {
			return fmt.Errorf("unable to render the original content of the sealed secret file %s (%w)", s.Path, err)
		}

		outDir, outFile := filepath.Split(s.Path)
		if s.OutFilename != "" {
			outFile = s.OutFilename
		}
		if s.OutDir != "" {
			outDir = s.OutDir
		}
		// TODO: Ensure that the output directory must be inside the application directory.
		if outDir != "" {
			if err := os.MkdirAll(filepath.Join(appDir, outDir), 0700); err != nil {
				return fmt.Errorf("unable to write decrypted content of sealed secret file %s to directory %s (%w)", s.Path, outDir, err)
			}
		}
		outPath := filepath.Join(appDir, outDir, outFile)

		if err := ioutil.WriteFile(outPath, content, 0644); err != nil {
			return fmt.Errorf("unable to write decrypted content of sealed secret file %s (%w)", s.Path, err)
		}
	}
	return nil
}

// listGroupedApplication retrieves all applications those should be handled by this director
// and then groups them by repoID.
func (d *detector) listGroupedApplication() map[string][]*model.Application {
	var (
		apps = d.appLister.ListByCloudProvider(d.provider.Name)
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

func (d *detector) loadDeploymentConfiguration(repoPath string, app *model.Application) (*config.Config, error) {
	path := filepath.Join(repoPath, app.GitPath.GetDeploymentConfigFilePath())
	cfg, err := config.LoadFromYAML(path)
	if err != nil {
		return nil, err
	}
	if appKind, ok := config.ToApplicationKind(cfg.Kind); !ok || appKind != app.Kind {
		return nil, fmt.Errorf("application in deployment configuration file is not match, got: %s, expected: %s", appKind, app.Kind)
	}

	if cfg.KubernetesDeploymentSpec != nil && cfg.KubernetesDeploymentSpec.Input.HelmChart != nil {
		chartRepoName := cfg.KubernetesDeploymentSpec.Input.HelmChart.Repository
		if chartRepoName != "" {
			cfg.KubernetesDeploymentSpec.Input.HelmChart.Insecure = d.config.IsInsecureChartRepository(chartRepoName)
		}
	}

	return cfg, nil
}

func (d *detector) ProviderName() string {
	return d.provider.Name
}

// groupManifests compares the given head and live manifests to divide them into three groups:
// - adds: contains all manifests that appear in lives but not in heads
// - deletes: contains all manifests that appear in heads but not in lives
// - inters: pairs of manifests that appear in both heads and lives.
func groupManifests(heads, lives []provider.Manifest) (adds, deletes, headInters, liveInters []provider.Manifest) {
	// Sort the manifests before comparing.
	sort.Slice(heads, func(i, j int) bool {
		return heads[i].Key.IsLessWithIgnoringNamespace(heads[j].Key)
	})
	sort.Slice(lives, func(i, j int) bool {
		return lives[i].Key.IsLessWithIgnoringNamespace(lives[j].Key)
	})

	var h, l int
	for {
		if h >= len(heads) || l >= len(lives) {
			break
		}
		if heads[h].Key.IsEqualWithIgnoringNamespace(lives[l].Key) {
			headInters = append(headInters, heads[h])
			liveInters = append(liveInters, lives[l])
			h++
			l++
			continue
		}
		// Has in head but not in live so this should be a deleted one.
		if heads[h].Key.IsLessWithIgnoringNamespace(lives[l].Key) {
			deletes = append(deletes, heads[h])
			h++
			continue
		}
		// Has in live but not in head so this should be an added one.
		adds = append(adds, lives[l])
		l++
	}

	if len(heads) > h {
		deletes = append(deletes, heads[h:]...)
	}
	if len(lives) > l {
		adds = append(adds, lives[l:]...)
	}
	return
}

func makeSyncedState() model.ApplicationSyncState {
	return model.ApplicationSyncState{
		Status:      model.ApplicationSyncStatus_SYNCED,
		ShortReason: "",
		Reason:      "",
		Timestamp:   time.Now().Unix(),
	}
}

func makeOutOfSyncState(adds, deletes []provider.Manifest, changes map[provider.Manifest]*diff.Result, commit string) model.ApplicationSyncState {
	total := len(adds) + len(deletes) + len(changes)
	shortReason := fmt.Sprintf("There are %d manifests not synced (%d adds, %d deletes, %d changes)", total, len(adds), len(deletes), len(changes))

	var b strings.Builder
	if len(commit) >= 7 {
		commit = commit[:7]
	}
	b.WriteString(fmt.Sprintf("Diff between the running resources and the definitions in Git at commit %q:\n", commit))
	b.WriteString("--- Git\n+++ Cluster\n\n")

	index := 0
	for _, delete := range deletes {
		index++
		b.WriteString(fmt.Sprintf("- %d. %s\n\n", index, delete.Key.ReadableString()))
	}
	for _, add := range adds {
		index++
		b.WriteString(fmt.Sprintf("+ %d. %s\n\n", index, add.Key.ReadableString()))
	}

	const maxPrintDiffs = 3
	var prints = 0
	for m, d := range changes {
		opts := []diff.RenderOption{
			diff.WithLeftPadding(1),
		}
		switch {
		case m.Key.IsSecret():
			opts = append(opts, diff.WithMaskPath("data"))
		case m.Key.IsConfigMap():
			opts = append(opts, diff.WithMaskPath("data"))
		}
		renderer := diff.NewRenderer(opts...)

		index++
		b.WriteString(fmt.Sprintf("* %d. %s\n\n", index, m.Key.ReadableString()))
		b.WriteString(renderer.Render(d.Nodes()))
		b.WriteString("\n")

		prints++
		if prints >= maxPrintDiffs {
			break
		}
	}

	if prints < len(changes) {
		b.WriteString(fmt.Sprintf("... (diffs from %d other manifests are omitted\n", len(changes)-prints))
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
