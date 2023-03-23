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

package trigger

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/filematcher"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type Determiner interface {
	ShouldTrigger(ctx context.Context, app *model.Application, appCfg *config.GenericApplicationSpec) (bool, error)
}

type determiners struct {
	onCommand   Determiner
	onOutOfSync Determiner
	onCommit    Determiner
	onChain     Determiner
}

func (ds *determiners) Determiner(k model.TriggerKind) Determiner {
	switch k {
	case model.TriggerKind_ON_COMMAND:
		return ds.onCommand
	case model.TriggerKind_ON_OUT_OF_SYNC:
		return ds.onOutOfSync
	case model.TriggerKind_ON_CHAIN:
		return ds.onChain
	default:
		return ds.onCommit
	}
}

type OnCommandDeterminer struct {
}

func NewOnCommandDeterminer() *OnCommandDeterminer {
	return &OnCommandDeterminer{}
}

// ShouldTrigger decides whether a given application should be triggered or not.
func (d *OnCommandDeterminer) ShouldTrigger(_ context.Context, _ *model.Application, appCfg *config.GenericApplicationSpec) (bool, error) {
	if appCfg.Trigger.OnCommand.Disabled {
		return false, nil
	}

	return true, nil
}

type OnChainDeterminer struct {
}

func NewOnChainDeterminer() *OnChainDeterminer {
	return &OnChainDeterminer{}
}

func (d *OnChainDeterminer) ShouldTrigger(_ context.Context, _ *model.Application, appCfg *config.GenericApplicationSpec) (bool, error) {
	if *appCfg.Trigger.OnChain.Disabled {
		return false, nil
	}
	return true, nil
}

type OnOutOfSyncDeterminer struct {
	client apiClient
}

func NewOnOutOfSyncDeterminer(client apiClient) *OnOutOfSyncDeterminer {
	return &OnOutOfSyncDeterminer{
		client: client,
	}
}

// ShouldTrigger decides whether a given application should be triggered or not.
func (d *OnOutOfSyncDeterminer) ShouldTrigger(ctx context.Context, app *model.Application, appCfg *config.GenericApplicationSpec) (bool, error) {
	if *appCfg.Trigger.OnOutOfSync.Disabled {
		return false, nil
	}

	// Find the most recently triggered deployment.
	// Nil means it seems the application has been added recently
	// and no deployment was triggered yet.
	ref := app.MostRecentlyTriggeredDeployment
	if ref == nil {
		return true, nil
	}

	resp, err := d.client.GetDeployment(ctx, &pipedservice.GetDeploymentRequest{
		Id: ref.DeploymentId,
	})
	if err != nil {
		return false, err
	}
	deployment := resp.Deployment

	// Check if it was already completed or not.
	// Not yet completed means the application is deploying currently,
	// so no need to trigger a new deployment for it.
	if !deployment.Status.IsCompleted() {
		return false, nil
	}

	// Check the elapsed time since the last deployment.
	if time.Since(time.Unix(deployment.CompletedAt, 0)) < appCfg.Trigger.OnOutOfSync.MinWindow.Duration() {
		return false, nil
	}

	return true, nil
}

type LastTriggeredCommitGetter interface {
	Get(ctx context.Context, applicationID string) (string, error)
}

type OnCommitDeterminer struct {
	repo         git.Repo
	targetCommit string
	commitGetter LastTriggeredCommitGetter
	logger       *zap.Logger
}

func NewOnCommitDeterminer(repo git.Repo, targetCommit string, cg LastTriggeredCommitGetter, logger *zap.Logger) Determiner {
	return &OnCommitDeterminer{
		repo:         repo,
		targetCommit: targetCommit,
		commitGetter: cg,
		logger:       logger.Named("determiner"),
	}
}

// ShouldTrigger decides whether a given application should be triggered or not.
func (d *OnCommitDeterminer) ShouldTrigger(ctx context.Context, app *model.Application, appCfg *config.GenericApplicationSpec) (bool, error) {
	logger := d.logger.With(
		zap.String("app", app.Name),
		zap.String("app-id", app.Id),
		zap.String("target-commit", d.targetCommit),
	)

	// Not trigger in case users disable auto trigger deploy on change and the user config is unignorable.
	if appCfg.Trigger.OnCommit.Disabled {
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
		logger.Debug(fmt.Sprintf("no update to sync for application, hash: %s", d.targetCommit))
		return false, nil
	}

	// List the changed files between those two commits and
	// determine whether this application was touch by those changed files.
	changedFiles, err := d.repo.ChangedFiles(ctx, preCommit, d.targetCommit)
	if err != nil {
		return false, err
	}

	touched, err := isTouchedByChangedFiles(app.GitPath.Path, appCfg.Trigger.OnCommit.Paths, appCfg.Trigger.OnCommit.Ignores, changedFiles)
	if err != nil {
		return false, err
	}

	if !touched {
		logger.Info("application was not touched by any new commits", zap.String("last-triggered-commit", preCommit))
		return false, nil
	}

	return true, nil
}

// isTouchedByChangedFiles checks whether this application changed files can trigger a new deployment or not (considered as "touched")
// The logic of watching files pattern contains both "includes" and "excludes" filter and be implemented as flow:
//  1. If any of changed files are listed in excludes, app is NOT considered as touched
//  2. If pass (1) and any of changed files are listed in includes, app is considered as touched
//  3. If any changes are under the app dir, app is considered as touched
func isTouchedByChangedFiles(appDir string, includes, excludes []string, changedFiles []string) (bool, error) {
	if !strings.HasSuffix(appDir, "/") {
		appDir += "/"
	}

	// If any changed files matches the specified "excludes"
	// this application is consided as not touched.
	for _, change := range excludes {
		matcher, err := filematcher.NewPatternMatcher([]string{change})
		if err != nil {
			return false, err
		}
		if matcher.MatchesAny(changedFiles) {
			return false, nil
		}
	}

	// If all changed files do not match any specified "excludes",
	// then if any changed files match the specified "includes"
	// this application is consided as touched.
	for _, change := range includes {
		matcher, err := filematcher.NewPatternMatcher([]string{change})
		if err != nil {
			return false, err
		}
		if matcher.MatchesAny(changedFiles) {
			return true, nil
		}
	}

	// It's considered any files changed inside the application directory as touched.
	for _, cf := range changedFiles {
		if ok := strings.HasPrefix(cf, appDir); ok {
			return true, nil
		}
	}

	return false, nil
}
