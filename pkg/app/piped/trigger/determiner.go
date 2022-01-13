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
	if !model.IsCompletedDeployment(deployment.Status) {
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
		logger.Info(fmt.Sprintf("no update to sync for application, hash: %s", d.targetCommit))
		return false, nil
	}

	// List the changed files between those two commits and
	// determine whether this application was touch by those changed files.
	changedFiles, err := d.repo.ChangedFiles(ctx, preCommit, d.targetCommit)
	if err != nil {
		return false, err
	}

	// TODO: Remove deprecated `appCfg.TriggerPaths` configuration.
	checkingPaths := make([]string, 0, len(appCfg.Trigger.OnCommit.Paths)+len(appCfg.TriggerPaths))
	// Note: appCfg.TriggerPaths or appCfg.Trigger.OnCommit.Paths may contain "" (empty string)
	// in case users use one of them without the other, that cause unexpected "" path in the checkingPaths list
	// leads to always trigger deployment since "" path matched all other paths.
	// The below logic is to remove that "" path from checking path list, will remove after remove the
	// deprecated appCfg.TriggerPaths.
	for _, p := range appCfg.Trigger.OnCommit.Paths {
		if p != "" {
			checkingPaths = append(checkingPaths, p)
		}
	}
	for _, p := range appCfg.TriggerPaths {
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
