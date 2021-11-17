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

type applicationLister interface {
	List() []*model.Application
}

type environmentGetter interface {
	Get(ctx context.Context, id string) (*model.Environment, error)
	GetByName(ctx context.Context, name string) (*model.Environment, error)
}

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
	fsys               fs.FS
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
		// FIXME: Find another way to specify the root dir for a file system
		fsys:   os.DirFS("/"),
		logger: logger.Named("app-config-reporter"),
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
		id := model.BuildGitPathID(app.GitPath.Repo.Id, app.GitPath.Path, app.GitPath.ConfigFilename)
		registeredAppPaths[id] = struct{}{}
	}

	if err := r.updateUnregisteredApps(ctx, registeredAppPaths); err != nil {
		return fmt.Errorf("failed to update unregistered applications: %w", err)
	}
	if err := r.updateRegisteredApps(ctx, registeredAppPaths); err != nil {
		return fmt.Errorf("failed to update registered applications: %w", err)
	}

	return nil
}

// updateUnregisteredApps sends all unregistered application configurations to the control-plane.
func (r *Reporter) updateUnregisteredApps(ctx context.Context, registeredAppPaths map[string]struct{}) error {
	apps := make([]*model.ApplicationInfo, 0)
	for repoID, repo := range r.gitRepos {
		as, err := r.findUnregisteredApps(ctx, repo.GetPath(), repoID, registeredAppPaths)
		if err != nil {
			return err
		}
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
	apps := make([]*model.ApplicationInfo, 0)
	err := fs.WalkDir(r.fsys, strings.TrimPrefix(repoPath, "/"), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		if shouldSkip(repoID, repoPath, filepath.Base(path), registeredAppPaths, false) {
			return nil
		}

		appInfo, err := r.readApplicationInfo(ctx, repoPath, path)
		if err != nil {
			r.logger.Error("failed to read application info",
				zap.String("repo-id", repoID),
				zap.String("config-file-path", path),
				zap.Error(err),
			)
			return nil
		}
		apps = append(apps, appInfo)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to inspect files under %s: %w", repoPath, err)
	}
	return apps, nil
}

// updateRegisteredApps sends application configurations that have changed since the last time to the control-plane.
func (r *Reporter) updateRegisteredApps(ctx context.Context, registeredAppPaths map[string]struct{}) (err error) {
	apps := make([]*model.ApplicationInfo, 0)
	for repoID, repo := range r.gitRepos {
		var headCommit git.Commit
		headCommit, err = repo.GetLatestCommit(ctx)
		if err != nil {
			return fmt.Errorf("failed to get the latest commit of %s: %w", repoID, err)
		}
		as, err := r.findRegisteredApps(ctx, repoID, repo, headCommit.Hash, registeredAppPaths)
		if err != nil {
			return err
		}
		apps = append(apps, as...)
		id := repoID
		defer func() {
			if err == nil {
				r.lastScannedCommits[id] = headCommit.Hash
			}
		}()
	}
	if len(apps) == 0 {
		return nil
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
	return nil
}

type gitRepo interface {
	GetPath() string
	ChangedFiles(ctx context.Context, from, to string) ([]string, error)
}

// findRegisteredApps finds out registered application info in the given git repository.
func (r *Reporter) findRegisteredApps(ctx context.Context, repoID string, repo gitRepo, headCommitHash string, registeredAppPaths map[string]struct{}) ([]*model.ApplicationInfo, error) {
	// Skip if the head commit is already scanned.
	lastScannedCommit, ok := r.lastScannedCommits[repoID]
	if ok && headCommitHash == lastScannedCommit {
		return []*model.ApplicationInfo{}, nil
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
		if shouldSkip(repoID, repo.GetPath(), filename, registeredAppPaths, true) {
			continue
		}
		appInfo, err := r.readApplicationInfo(ctx, repo.GetPath(), filepath.Join(repo.GetPath(), filename))
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

func shouldSkip(repoID, path, cfgFilename string, registeredAppPaths map[string]struct{}, wantRegistered bool) bool {
	if !strings.HasSuffix(cfgFilename, model.DefaultDeploymentConfigFileExtension) {
		return true
	}
	gitPathID := model.BuildGitPathID(repoID, path, cfgFilename)
	if _, registered := registeredAppPaths[gitPathID]; registered != wantRegistered {
		return true
	}
	return false
}

func (r *Reporter) readApplicationInfo(ctx context.Context, path, cfgFilePath string) (appInfo *model.ApplicationInfo, err error) {
	b, err := fs.ReadFile(r.fsys, cfgFilePath)
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
		Path:           path,
		ConfigFilename: filepath.Base(cfgFilePath),
		Labels:         spec.Labels,
	}, nil
}
