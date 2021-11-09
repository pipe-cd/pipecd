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

package trigger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/filematcher"
	"github.com/pipe-cd/pipe/pkg/git"
	"github.com/pipe-cd/pipe/pkg/model"
)

type LastTriggeredCommitGetter interface {
	Get(ctx context.Context, applicationID string) (string, error)
}

type Determiner struct {
	repo         git.Repo
	targetCommit string
	commitGetter LastTriggeredCommitGetter
	// Flag `ignoreUserConfig` set to `true` will force check changes and use it to determine
	// the application deployment should be triggered or not, regardless of the user's configuration.
	ignoreUserConfig bool
	logger           *zap.Logger
}

func NewDeterminer(repo git.Repo, targetCommit string, cg LastTriggeredCommitGetter, ignoreUserConfig bool, logger *zap.Logger) *Determiner {
	return &Determiner{
		repo:             repo,
		targetCommit:     targetCommit,
		commitGetter:     cg,
		ignoreUserConfig: ignoreUserConfig,
		logger:           logger.Named("determiner"),
	}
}

// ShouldTrigger decides whether a given application should be triggered or not.
func (d *Determiner) ShouldTrigger(ctx context.Context, app *model.Application) (bool, error) {
	logger := d.logger.With(
		zap.String("app", app.Name),
		zap.String("app-id", app.Id),
		zap.String("target-commit", d.targetCommit),
	)

	// TODO: Add logic to determine trigger or not based on other configuration than onCommit.
	return d.shouldTriggerOnCommit(ctx, app, logger)
}

func (d *Determiner) shouldTriggerOnCommit(ctx context.Context, app *model.Application, logger *zap.Logger) (bool, error) {
	deployConfig, err := loadDeploymentConfiguration(d.repo.GetPath(), app)
	if err != nil {
		return false, err
	}

	// Not trigger in case users disable auto trigger deploy on change and the user config is unignorable.
	if deployConfig.Trigger.OnCommit.Disabled && !d.ignoreUserConfig {
		logger.Info(fmt.Sprintf("auto trigger deployment disabled for application, hash: %s", d.targetCommit))
		return false, nil
	}

	preCommit, err := d.commitGetter.Get(ctx, app.Id)
	if err != nil {
		logger.Error("failed to get last triggered commit", zap.Error(err))
		return false, err
	}

	// There is no previous deployment so we don't need to check anymore.
	// Just do it.
	if preCommit == "" {
		logger.Info("no previously triggered deployment was found")
		return true, nil
	}

	// Check whether the most recently applied one is the target commit or not.
	// If so, nothing to do for this time.
	if preCommit == d.targetCommit {
		logger.Info(fmt.Sprintf("no update to sync for application, hash: %s", d.targetCommit))
		return false, nil
	}

	// List the changed files between those two commits and
	// determine whether this application was touch by those changed files.
	changedFiles, err := d.repo.ChangedFiles(ctx, preCommit, d.targetCommit)
	if err != nil {
		return false, err
	}

	// TODO: Remove deprecated `deployConfig.TriggerPaths` configuration.
	checkingPaths := make([]string, 0, len(deployConfig.Trigger.OnCommit.Paths)+len(deployConfig.TriggerPaths))
	// Note: deployConfig.TriggerPaths or deployConfig.Trigger.OnCommit.Paths may contain "" (empty string)
	// in case users use one of them without the other, that cause unexpected "" path in the checkingPaths list
	// leads to always trigger deployment since "" path matched all other paths.
	// The below logic is to remove that "" path from checking path list, will remove after remove the
	// deprecated deployConfig.TriggerPaths.
	for _, p := range deployConfig.Trigger.OnCommit.Paths {
		if p != "" {
			checkingPaths = append(checkingPaths, p)
		}
	}
	for _, p := range deployConfig.TriggerPaths {
		if p != "" {
			checkingPaths = append(checkingPaths, p)
		}
	}

	touched, err := isTouchedByChangedFiles(app.GitPath.Path, checkingPaths, changedFiles)
	if err != nil {
		return false, err
	}

	if !touched {
		logger.Info("application was not touched by any new commits", zap.String("last-triggered-commit", preCommit))
		return false, nil
	}

	return true, nil
}

func loadDeploymentConfiguration(repoPath string, app *model.Application) (*config.GenericDeploymentSpec, error) {
	var (
		relPath = app.GitPath.GetDeploymentConfigFilePath()
		absPath = filepath.Join(repoPath, relPath)
	)

	cfg, err := config.LoadFromYAML(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("deployment config file %s was not found", relPath)
		}
		return nil, err
	}
	if appKind, ok := config.ToApplicationKind(cfg.Kind); !ok || appKind != app.Kind {
		return nil, fmt.Errorf("invalid application kind in the deployment config file, got: %s, expected: %s", appKind, app.Kind)
	}

	spec, ok := cfg.GetGenericDeployment()
	if !ok {
		return nil, fmt.Errorf("unsupported application kind: %s", app.Kind)
	}

	return &spec, nil
}

func isTouchedByChangedFiles(appDir string, changes []string, changedFiles []string) (bool, error) {
	if !strings.HasSuffix(appDir, "/") {
		appDir += "/"
	}

	// If any files inside the application directory was changed
	// this application is considered as touched.
	for _, cf := range changedFiles {
		if ok := strings.HasPrefix(cf, appDir); ok {
			return true, nil
		}
	}

	// If any changed files matches the specified "changes"
	// this application is consided as touched too.
	for _, change := range changes {
		matcher, err := filematcher.NewPatternMatcher([]string{change})
		if err != nil {
			return false, err
		}
		if matcher.MatchesAny(changedFiles) {
			return true, nil
		}
	}

	return false, nil
}
