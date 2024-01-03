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

package appconfigreporter

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
)

var (
	errMissingRequiredField = errors.New("missing required field")
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

type fileSystem struct{}

func (s *fileSystem) Open(name string) (fs.File, error) { return os.Open(name) }

type Reporter struct {
	apiClient         apiClient
	gitClient         gitClient
	applicationLister applicationLister
	config            *config.PipedSpec
	gitRepos          map[string]git.Repo
	gracePeriod       time.Duration
	// Cache for the last scanned commit for each registered application.
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

	// Scan them once first.
	if err := r.scanAppConfigs(ctx); err != nil {
		r.logger.Error("failed to check application configurations defined in Git", zap.Error(err))
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
	headCommits := make(map[string]string, len(r.gitRepos))
	for repoID, repo := range r.gitRepos {
		if err := repo.Pull(ctx, repo.GetClonedBranch()); err != nil {
			r.logger.Error("failed to update repo to latest",
				zap.String("repo-id", repoID),
				zap.Error(err),
			)
			return err
		}
		headCommit, err := repo.GetLatestCommit(ctx)
		if err != nil {
			return fmt.Errorf("failed to get the latest commit of %s: %w", repoID, err)
		}
		headCommits[repoID] = headCommit.Hash
	}

	if err := r.updateRegisteredApps(ctx, headCommits); err != nil {
		return err
	}
	if err := r.updateUnregisteredApps(ctx); err != nil {
		return err
	}

	return nil
}

// updateRegisteredApps sends application configurations that have changed since the last time to the control-plane.
func (r *Reporter) updateRegisteredApps(ctx context.Context, headCommits map[string]string) error {
	outOfSyncRegisteredApps := make([]*model.ApplicationInfo, 0)
	for repoID, repo := range r.gitRepos {
		headCommit := headCommits[repoID]
		rs := r.findOutOfSyncRegisteredApps(repo.GetPath(), repoID, headCommit)
		r.logger.Info(fmt.Sprintf("found out %d valid registered applications that config has been changed in repository %q", len(rs), repoID))
		outOfSyncRegisteredApps = append(outOfSyncRegisteredApps, rs...)
	}
	if len(outOfSyncRegisteredApps) == 0 {
		return nil
	}

	_, err := r.apiClient.UpdateApplicationConfigurations(
		ctx,
		&pipedservice.UpdateApplicationConfigurationsRequest{
			Applications: outOfSyncRegisteredApps,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update application configurations: %w", err)
	}

	// Memorize registered applications, which are updated above.
	for _, app := range outOfSyncRegisteredApps {
		r.lastScannedCommits[app.Id] = headCommits[app.RepoId]
	}

	return nil
}

// updateUnregisteredApps sends all unregistered application configurations to the control-plane.
func (r *Reporter) updateUnregisteredApps(ctx context.Context) error {
	unregisteredApps := make([]*model.ApplicationInfo, 0)
	for repoID, repo := range r.gitRepos {
		// The unregistered apps sent previously aren't persisted, that's why it has to send them again even if it's scanned one.
		us, err := r.findUnregisteredApps(repo.GetPath(), repoID)
		if err != nil {
			return err
		}
		r.logger.Info(fmt.Sprintf("found out %d valid unregistered applications in repository %q", len(us), repoID))
		unregisteredApps = append(unregisteredApps, us...)
	}

	// Even if the result is zero, we need to report at least once.
	// However, it should return after the second time which is unnecessary.
	if len(unregisteredApps) == 0 {
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
			Applications: unregisteredApps,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to put the latest unregistered application configurations: %w", err)
	}
	return nil
}

// findOutOfSyncRegisteredApps finds out registered application info that should be updated in the given git repository.
func (r *Reporter) findOutOfSyncRegisteredApps(repoPath, repoID, headCommit string) []*model.ApplicationInfo {
	// Compare the apps registered on Control-plane with the latest config file
	// and return only the ones that have been changed.
	apps := make([]*model.ApplicationInfo, 0)
	for _, app := range r.applicationLister.List() {
		if app.GitPath.Repo.Id != repoID {
			continue
		}

		// Skip if there is no new commit pushed from last scanned time for this application.
		if lc, ok := r.lastScannedCommits[app.Id]; ok && headCommit == lc {
			continue
		}

		appCfg, err := r.readApplicationInfo(repoPath, repoID, app.GitPath.GetApplicationConfigFilePath())
		if errors.Is(err, errMissingRequiredField) {
			// For historical reasons, we need to treat applications that don't define app config in a file as normal.
			r.logger.Warn("found a registered application config file that is missing a required field",
				zap.String("repo-id", repoID),
				zap.String("config-file-path", app.GitPath.GetApplicationConfigFilePath()),
				zap.Error(err),
			)
			continue
		}
		if err != nil {
			r.logger.Error("failed to read registered application config file",
				zap.String("repo-id", repoID),
				zap.String("config-file-path", app.GitPath.GetApplicationConfigFilePath()),
				zap.Error(err),
			)
			// Continue reading so that it can return apps as much as possible.
			continue
		}

		// Memorize the application last scanned commit in case the app is unchanged.
		if r.isSynced(appCfg, app) {
			r.lastScannedCommits[app.Id] = headCommit
			continue
		}
		appCfg.Id = app.Id
		apps = append(apps, appCfg)
	}
	return apps
}

func (r *Reporter) isSynced(appInfo *model.ApplicationInfo, app *model.Application) bool {
	if appInfo.Kind != app.Kind {
		r.logger.Warn("kind in application config has been changed which isn't allowed",
			zap.String("app-id", app.Id),
			zap.String("repo-id", app.GitPath.Repo.Id),
			zap.String("config-file-path", app.GitPath.GetApplicationConfigFilePath()),
		)
	}

	// TODO: Make it possible to follow the ApplicationInfo field changes
	if appInfo.Name != app.Name {
		return false
	}
	if appInfo.Description != app.Description {
		return false
	}
	if len(appInfo.Labels) != len(app.Labels) {
		return false
	}
	for key, value := range appInfo.Labels {
		if value != app.Labels[key] {
			return false
		}
	}
	return true
}

// findUnregisteredApps finds out unregistered application info in the given git repository.
// The file name must be default name in order to be recognized as an Application config.
func (r *Reporter) findUnregisteredApps(repoPath, repoID string) ([]*model.ApplicationInfo, error) {
	var (
		apps     = r.applicationLister.List()
		selector = r.config.AppSelector
	)
	// Create a map to determine the app is registered by GitPath.
	registeredAppPaths := make(map[string]struct{}, len(apps))
	for _, app := range apps {
		if app.GitPath.Repo.Id != repoID {
			continue
		}
		registeredAppPaths[app.GitPath.GetApplicationConfigFilePath()] = struct{}{}
	}

	out := make([]*model.ApplicationInfo, 0)
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
		if _, registered := registeredAppPaths[cfgRelPath]; registered {
			return nil
		}

		appInfo, err := r.readApplicationInfo(repoPath, repoID, cfgRelPath)
		if errors.Is(err, errMissingRequiredField) {
			r.logger.Warn("found an unregistered application config file that is missing a required field",
				zap.String("repo-id", repoID),
				zap.String("config-file-path", cfgRelPath),
				zap.Error(err),
			)
			return nil
		}
		if err != nil {
			r.logger.Error("failed to read unregistered application info",
				zap.String("repo-id", repoID),
				zap.String("config-file-path", cfgRelPath),
				zap.Error(err),
			)
			return nil
		}

		// Filter the apps by appSelector if appSelector set.
		if len(selector) != 0 && !appInfo.ContainLabels(selector) {
			return nil
		}

		out = append(out, appInfo)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to inspect files under %s: %w", repoPath, err)
	}
	return out, nil
}

func (r *Reporter) readApplicationInfo(repoDir, repoID, cfgRelPath string) (*model.ApplicationInfo, error) {
	b, err := fs.ReadFile(r.fileSystem, filepath.Join(repoDir, cfgRelPath))
	if err != nil {
		return nil, fmt.Errorf("failed to open the configuration file: %w", err)
	}
	cfg, err := config.DecodeYAML(b)
	if err != nil {
		return nil, fmt.Errorf("failed to decode configuration file: %w", err)
	}

	spec, ok := cfg.GetGenericApplication()
	if !ok {
		return nil, fmt.Errorf("unsupported application kind %q", cfg.Kind)
	}

	kind, ok := cfg.Kind.ToApplicationKind()
	if !ok {
		return nil, fmt.Errorf("%q is not application config kind", cfg.Kind)
	}
	if spec.Name == "" {
		return nil, fmt.Errorf("missing application name: %w", errMissingRequiredField)
	}
	return &model.ApplicationInfo{
		Name:           spec.Name,
		Kind:           kind,
		Labels:         spec.Labels,
		RepoId:         repoID,
		Path:           filepath.Dir(cfgRelPath),
		ConfigFilename: filepath.Base(cfgRelPath),
		PipedId:        r.config.PipedID,
		Description:    spec.Description,
	}, nil
}
