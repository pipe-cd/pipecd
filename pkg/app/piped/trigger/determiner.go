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
	logger       *zap.Logger
}

func NewDeterminer(repo git.Repo, targetCommit string, cg LastTriggeredCommitGetter, logger *zap.Logger) *Determiner {
	return &Determiner{
		repo:         repo,
		targetCommit: targetCommit,
		commitGetter: cg,
		logger:       logger.Named("determiner"),
	}
}

// ShouldTrigger decides whether a given application should be triggered or not.
// Flag `ignorable` set to `false` will force check changes and use it to determine
// the application deployment should be triggered or not, regardless of the user's configuration.
func (d *Determiner) ShouldTrigger(ctx context.Context, app *model.Application, ignorable bool) (bool, error) {
	logger := d.logger.With(
		zap.String("app", app.Name),
		zap.String("app-id", app.Id),
		zap.String("target-commit", d.targetCommit),
	)

	deployConfig, err := loadDeploymentConfiguration(d.repo.GetPath(), app)
	if err != nil {
		return false, err
	}

	// Not trigger in case users disable auto trigger deploy on change and the change is ignorable.
	if deployConfig.Trigger.DisableAutoDeployOnChange && ignorable {
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

	touched, err := isTouchedByChangedFiles(app.GitPath.Path, deployConfig.Trigger.Paths, changedFiles)
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
