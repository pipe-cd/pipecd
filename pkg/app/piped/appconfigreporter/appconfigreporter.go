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

type fileSystem struct{}

func (s *fileSystem) Open(name string) (fs.File, error) { return os.Open(name) }

type Reporter struct {
	apiClient         apiClient
	gitClient         gitClient
	applicationLister applicationLister
	config            *config.PipedSpec
	gitRepos          map[string]git.Repo
	gracePeriod       time.Duration
	// Cache for the last scanned commit for each repository.
	lastScannedCommits map[string]string
	fileSystem         fs.FS
	logger             *zap.Logger

	// Whether it already swept all unregistered apps from control-plane.
	sweptUnregisteredApps bool
}

func NewReporter(
	apiClient apiClient,
	gitClient gitClient,
	appLister applicationLister,
	cfg *config.PipedSpec,
	gracePeriod time.Duration,
	logger *zap.Logger,
) *Reporter {
	return &Reporter{
		apiClient:          apiClient,
		gitClient:          gitClient,
		applicationLister:  appLister,
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

	ticker := time.NewTicker(r.config.AppConfigSyncInterval.Duration())
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

	apps := r.applicationLister.List()
	// Create a map from Git path key to app id, to determine if the application is registered.
	registeredAppPaths := make(map[string]string, len(apps))
	for _, app := range apps {
		key := makeGitPathKey(app.GitPath.Repo.Id, filepath.Join(app.GitPath.Path, app.GitPath.ConfigFilename))
		registeredAppPaths[key] = app.Id
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
func (r *Reporter) updateRegisteredApps(ctx context.Context, registeredAppPaths map[string]string) (err error) {
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
func (r *Reporter) findRegisteredApps(ctx context.Context, repoID string, repo gitRepo, lastScannedCommit, headCommitHash string, registeredAppPaths map[string]string) ([]*model.ApplicationInfo, error) {
	if lastScannedCommit == "" {
		// Scan all files because it seems to be the first commit since Piped starts.
		apps := make([]*model.ApplicationInfo, 0)
		err := fs.WalkDir(r.fileSystem, repo.GetPath(), func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			cfgRelPath, err := filepath.Rel(repo.GetPath(), path)
			if err != nil {
				return err
			}
			gitPathKey := makeGitPathKey(repoID, cfgRelPath)
			appID, registered := registeredAppPaths[gitPathKey]
			if !registered {
				return nil
			}

			appInfo, err := r.readApplicationInfo(repo.GetPath(), repoID, filepath.Dir(cfgRelPath), filepath.Base(cfgRelPath))
			if err != nil {
				r.logger.Error("failed to read application info",
					zap.String("repo-id", repoID),
					zap.String("config-file-path", cfgRelPath),
					zap.Error(err),
				)
				return nil
			}
			appInfo.Id = appID
			apps = append(apps, appInfo)
			// Continue reading so that it can return apps as much as possible.
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to inspect files under %s: %w", repo.GetPath(), err)
		}
		return apps, nil
	}

	filePaths, err := repo.ChangedFiles(ctx, lastScannedCommit, headCommitHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get files those were touched between two commits: %w", err)
	}
	if len(filePaths) == 0 {
		// The case where all changes have been fully reverted.
		return []*model.ApplicationInfo{}, nil
	}
	apps := make([]*model.ApplicationInfo, 0)
	for _, path := range filePaths {
		gitPathKey := makeGitPathKey(repoID, path)
		appID, registered := registeredAppPaths[gitPathKey]
		if !registered {
			continue
		}
		appInfo, err := r.readApplicationInfo(repo.GetPath(), repoID, filepath.Dir(path), filepath.Base(path))
		if err != nil {
			r.logger.Error("failed to read application info",
				zap.String("repo-id", repoID),
				zap.String("config-file-path", path),
				zap.Error(err),
			)
			continue
		}
		appInfo.Id = appID
		apps = append(apps, appInfo)
	}
	return apps, nil
}

// updateUnregisteredApps sends all unregistered application configurations to the control-plane.
func (r *Reporter) updateUnregisteredApps(ctx context.Context, registeredAppPaths map[string]string) error {
	apps := make([]*model.ApplicationInfo, 0)
	for repoID, repo := range r.gitRepos {
		as, err := r.findUnregisteredApps(repo.GetPath(), repoID, registeredAppPaths)
		if err != nil {
			return err
		}
		r.logger.Info(fmt.Sprintf("found out %d unregistered applications in repository %s", len(as), repoID))
		apps = append(apps, as...)
	}
	// Prevent Piped from continuing to send meaningless requests
	// even though Control-plane already knows that there are zero unregistered apps.
	if len(apps) == 0 {
		if r.sweptUnregisteredApps {
			return nil
		}
		r.sweptUnregisteredApps = true
	} else {
		r.sweptUnregisteredApps = false
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
// The file name must be default name in order to be recognized as an Application config.
func (r *Reporter) findUnregisteredApps(repoPath, repoID string, registeredAppPaths map[string]string) ([]*model.ApplicationInfo, error) {
	apps := make([]*model.ApplicationInfo, 0)
	// Scan all files under the repository.
	err := fs.WalkDir(r.fileSystem, repoPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		cfgRelPath, err := filepath.Rel(repoPath, path)
		if err != nil {
			return err
		}
		if !model.IsApplicationConfigFile(filepath.Base(cfgRelPath)) {
			return nil
		}
		gitPathKey := makeGitPathKey(repoID, cfgRelPath)
		if _, registered := registeredAppPaths[gitPathKey]; registered {
			return nil
		}

		appInfo, err := r.readApplicationInfo(repoPath, repoID, filepath.Dir(cfgRelPath), filepath.Base(cfgRelPath))
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
		return nil, fmt.Errorf("failed to inspect files under %s: %w", repoPath, err)
	}
	return apps, nil
}

func (r *Reporter) readApplicationInfo(repoDir, repoID, appDirRelPath, cfgFilename string) (*model.ApplicationInfo, error) {
	b, err := fs.ReadFile(r.fileSystem, filepath.Join(repoDir, appDirRelPath, cfgFilename))
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

	if spec.Name == "" {
		return nil, fmt.Errorf("missing application name")
	}
	kind, ok := cfg.Kind.ToApplicationKind()
	if !ok {
		return nil, fmt.Errorf("%q is not application config kind", cfg.Kind)
	}
	return &model.ApplicationInfo{
		Name:           spec.Name,
		Kind:           kind,
		Labels:         spec.Labels,
		RepoId:         repoID,
		Path:           appDirRelPath,
		ConfigFilename: cfgFilename,
	}, nil
}

// makeGitPathKey builds a unique path between repositories.
// cfgFilePath is a relative path from the repo root.
func makeGitPathKey(repoID, cfgFilePath string) string {
	return fmt.Sprintf("%s:%s", repoID, cfgFilePath)
}
