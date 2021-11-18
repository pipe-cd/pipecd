// Copyright 2021 The PipeCD Authors.
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

package appconfigreporter

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipe/pkg/app/api/service/pipedservice"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/git"
	"github.com/pipe-cd/pipe/pkg/model"
)

type apiClient interface {
	UpdateApplicationConfigurations(ctx context.Context, in *pipedservice.UpdateApplicationConfigurationsRequest, opts ...grpc.CallOption) (*pipedservice.UpdateApplicationConfigurationsResponse, error)
	ReportUnregisteredApplicationConfigurations(ctx context.Context, in *pipedservice.ReportUnregisteredApplicationConfigurationsRequest, opts ...grpc.CallOption) (*pipedservice.ReportUnregisteredApplicationConfigurationsResponse, error)
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

type gitRepo interface {
	GetPath() string
	ChangedFiles(ctx context.Context, from, to string) ([]string, error)
}

type applicationLister interface {
	List() []*model.Application
}

type environmentGetter interface {
	GetByName(ctx context.Context, name string) (*model.Environment, error)
}

type fileSystem struct{}

func (s *fileSystem) Open(name string) (fs.File, error) { return os.Open(name) }

type Reporter struct {
	apiClient         apiClient
	gitClient         gitClient
	applicationLister applicationLister
	envGetter         environmentGetter
	config            *config.PipedSpec
	gitRepos          map[string]git.Repo
	gracePeriod       time.Duration
	// Cache for the last scanned commit for each repository.
	// Not goroutine safe.
	lastScannedCommits map[string]string
	fileSystem         fs.FS
	logger             *zap.Logger
}

func NewReporter(
	apiClient apiClient,
	gitClient gitClient,
	appLister applicationLister,
	envGetter environmentGetter,
	cfg *config.PipedSpec,
	gracePeriod time.Duration,
	logger *zap.Logger,
) *Reporter {
	return &Reporter{
		apiClient:          apiClient,
		gitClient:          gitClient,
		applicationLister:  appLister,
		envGetter:          envGetter,
		config:             cfg,
		gracePeriod:        gracePeriod,
		lastScannedCommits: make(map[string]string),
		fileSystem:         &fileSystem{},
		logger:             logger.Named("app-config-reporter"),
	}
}

func (r *Reporter) Run(ctx context.Context) error {
	r.logger.Info("start running app-config-reporter")

	// Pre-clone to cache the registered git repositories.
	r.gitRepos = make(map[string]git.Repo, len(r.config.Repositories))
	for _, repoCfg := range r.config.Repositories {
		repo, err := r.gitClient.Clone(ctx, repoCfg.RepoID, repoCfg.Remote, repoCfg.Branch, "")
		if err != nil {
			r.logger.Error("failed to clone repository",
				zap.String("repo-id", repoCfg.RepoID),
				zap.Error(err),
			)
			return err
		}
		r.gitRepos[repoCfg.RepoID] = repo
	}

	// FIXME: Think about sync interval of app config reporter
	ticker := time.NewTicker(r.config.SyncInterval.Duration())
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := r.scanAppConfigs(ctx); err != nil {
				r.logger.Error("failed to check application configurations defined in Git", zap.Error(err))
			}
		case <-ctx.Done():
			r.logger.Info("app-config-reporter has been stopped")
			return nil
		}
	}
}

// scanAppConfigs checks and reports two types of applications.
// One is applications registered in Control-plane already, and another is ones that aren't registered yet.
func (r *Reporter) scanAppConfigs(ctx context.Context) error {
	if len(r.gitRepos) == 0 {
		r.logger.Info("no repositories were configured for this piped")
		return nil
	}

	// Make all repos up-to-date.
	for repoID, repo := range r.gitRepos {
		if err := repo.Pull(ctx, repo.GetClonedBranch()); err != nil {
			r.logger.Error("failed to update repo to latest",
				zap.String("repo-id", repoID),
				zap.Error(err),
			)
			return err
		}
	}

	// Create a map to determine from GitPath if the application is registered.
	apps := r.applicationLister.List()
	registeredAppPaths := make(map[string]struct{}, len(apps))
	for _, app := range apps {
		key := makeGitPathKey(app.GitPath.Repo.Id, filepath.Join(app.GitPath.Path, app.GitPath.ConfigFilename))
		registeredAppPaths[key] = struct{}{}
	}

	if err := r.updateRegisteredApps(ctx, registeredAppPaths); err != nil {
		return fmt.Errorf("failed to update registered applications: %w", err)
	}
	if err := r.updateUnregisteredApps(ctx, registeredAppPaths); err != nil {
		return fmt.Errorf("failed to update unregistered applications: %w", err)
	}

	return nil
}

// updateRegisteredApps sends application configurations that have changed since the last time to the control-plane.
func (r *Reporter) updateRegisteredApps(ctx context.Context, registeredAppPaths map[string]struct{}) (err error) {
	apps := make([]*model.ApplicationInfo, 0)
	headCommits := make(map[string]string, len(r.gitRepos))
	for repoID, repo := range r.gitRepos {
		var headCommit git.Commit
		headCommit, err = repo.GetLatestCommit(ctx)
		if err != nil {
			return fmt.Errorf("failed to get the latest commit of %s: %w", repoID, err)
		}
		// Skip if the head commit is already scanned.
		lastScannedCommit, ok := r.lastScannedCommits[repoID]
		if ok && headCommit.Hash == lastScannedCommit {
			continue
		}
		var as []*model.ApplicationInfo
		as, err = r.findRegisteredApps(ctx, repoID, repo, lastScannedCommit, headCommit.Hash, registeredAppPaths)
		if err != nil {
			return err
		}
		apps = append(apps, as...)
		headCommits[repoID] = headCommit.Hash
	}
	defer func() {
		if err == nil {
			for repoID, hash := range headCommits {
				r.lastScannedCommits[repoID] = hash
			}
		}
	}()
	if len(apps) == 0 {
		return
	}

	_, err = r.apiClient.UpdateApplicationConfigurations(
		ctx,
		&pipedservice.UpdateApplicationConfigurationsRequest{
			Applications: apps,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update application configurations: %w", err)
	}
	return
}

// findRegisteredApps finds out registered application info in the given git repository.
func (r *Reporter) findRegisteredApps(ctx context.Context, repoID string, repo gitRepo, lastScannedCommit, headCommitHash string, registeredAppPaths map[string]struct{}) ([]*model.ApplicationInfo, error) {
	if lastScannedCommit == "" {
		return r.scanAllFiles(ctx, repo.GetPath(), repoID, registeredAppPaths, true)
	}

	files, err := repo.ChangedFiles(ctx, lastScannedCommit, headCommitHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get files those were touched between two commits: %w", err)
	}
	if len(files) == 0 {
		// The case where all changes have been fully reverted.
		return []*model.ApplicationInfo{}, nil
	}
	apps := make([]*model.ApplicationInfo, 0)
	for _, filename := range files {
		if shouldSkip(repoID, filename, registeredAppPaths, true) {
			continue
		}
		appInfo, err := r.readApplicationInfo(ctx, filepath.Join(repo.GetPath(), filename))
		if err != nil {
			r.logger.Error("failed to read application info",
				zap.String("repo-id", repoID),
				zap.String("config-file-path", filename),
				zap.Error(err),
			)
			continue
		}
		apps = append(apps, appInfo)
	}
	return apps, nil
}

// updateUnregisteredApps sends all unregistered application configurations to the control-plane.
func (r *Reporter) updateUnregisteredApps(ctx context.Context, registeredAppPaths map[string]struct{}) error {
	apps := make([]*model.ApplicationInfo, 0)
	for repoID, repo := range r.gitRepos {
		as, err := r.findUnregisteredApps(ctx, repo.GetPath(), repoID, registeredAppPaths)
		if err != nil {
			return err
		}
		r.logger.Info(fmt.Sprintf("found out %d unregistered applications in repository %s", len(as), repoID))
		apps = append(apps, as...)
	}
	if len(apps) == 0 {
		return nil
	}

	_, err := r.apiClient.ReportUnregisteredApplicationConfigurations(
		ctx,
		&pipedservice.ReportUnregisteredApplicationConfigurationsRequest{
			Applications: apps,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to put the latest unregistered application configurations: %w", err)
	}
	return nil
}

// findUnregisteredApps finds out unregistered application info in the given git repository.
func (r *Reporter) findUnregisteredApps(ctx context.Context, repoPath, repoID string, registeredAppPaths map[string]struct{}) ([]*model.ApplicationInfo, error) {
	return r.scanAllFiles(ctx, repoPath, repoID, registeredAppPaths, false)
}

// scanAllFiles inspects all files under the root or the given repository.
// And gives back all application info as much as possible.
func (r *Reporter) scanAllFiles(ctx context.Context, repoRoot, repoID string, registeredAppPaths map[string]struct{}, wantRegistered bool) ([]*model.ApplicationInfo, error) {
	apps := make([]*model.ApplicationInfo, 0)
	err := fs.WalkDir(r.fileSystem, repoRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		cfgRelPath, err := filepath.Rel(repoRoot, path)
		if err != nil {
			return err
		}

		if shouldSkip(repoID, cfgRelPath, registeredAppPaths, wantRegistered) {
			return nil
		}

		appInfo, err := r.readApplicationInfo(ctx, path)
		if err != nil {
			r.logger.Error("failed to read application info",
				zap.String("repo-id", repoID),
				zap.String("config-file-path", cfgRelPath),
				zap.Error(err),
			)
			return nil
		}
		apps = append(apps, appInfo)
		// Continue reading so that it can return apps as much as possible.
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to inspect files under %s: %w", repoRoot, err)
	}
	return apps, nil
}

// cfgFilePath is a relative path from the repo root.
func shouldSkip(repoID, cfgFilePath string, registeredAppPaths map[string]struct{}, wantRegistered bool) bool {
	if !strings.HasSuffix(cfgFilePath, model.DefaultDeploymentConfigFileExtension) {
		return true
	}
	gitPathKey := makeGitPathKey(repoID, cfgFilePath)
	if _, registered := registeredAppPaths[gitPathKey]; registered != wantRegistered {
		return true
	}
	return false
}

// makeGitPathKey builds a unique path between repositories.
// cfgFilePath is a relative path from the repo root.
func makeGitPathKey(repoID, cfgFilePath string) string {
	return fmt.Sprintf("%s:%s", repoID, cfgFilePath)
}

func (r *Reporter) readApplicationInfo(ctx context.Context, cfgFilePath string) (*model.ApplicationInfo, error) {
	b, err := fs.ReadFile(r.fileSystem, cfgFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open the configuration file: %w", err)
	}
	cfg, err := config.DecodeYAML(b)
	if err != nil {
		return nil, fmt.Errorf("failed to decode configuration file: %w", err)
	}

	spec, ok := cfg.GetGenericDeployment()
	if !ok {
		return nil, fmt.Errorf("unsupported application kind %q", cfg.Kind)
	}
	env, err := r.envGetter.GetByName(ctx, spec.EnvName)
	if err != nil {
		return nil, fmt.Errorf("failed to get env by name: %w", err)
	}

	return &model.ApplicationInfo{
		Name: spec.Name,
		// TODO: Convert Kind string into dedicated type
		//Kind:           cfg.Kind,
		EnvId:          env.Id,
		Path:           filepath.Dir(cfgFilePath),
		ConfigFilename: filepath.Base(cfgFilePath),
		Labels:         spec.Labels,
	}, nil
}
