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

package controller

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/api/service/pipedservice"
	pln "github.com/pipe-cd/pipe/pkg/app/piped/planner"
	"github.com/pipe-cd/pipe/pkg/app/piped/planner/registry"
	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

// What planner does:
// - Wait until there is no PLANNED or RUNNING deployment
// - Pick the oldest PENDING deployment to plan its pipeline
// - Compare with the last successful commit
// - Decide the pipeline should be executed (scale, progressive, rollback)
// - Update the pipeline stages and change the deployment status to PLANNED
type planner struct {
	// Readonly deployment model.
	deployment               *model.Deployment
	lastSuccessfulCommitHash string
	workingDir               string
	apiClient                apiClient
	gitClient                gitClient
	plannerRegistry          registry.Registry
	pipedConfig              *config.PipedSpec
	appManifestsCache        cache.Cache
	logger                   *zap.Logger

	deploymentConfig *config.Config
	done             atomic.Bool
	doneTimestamp    time.Time
	cancelled        bool
	cancelledCh      chan *model.ReportableCommand

	nowFunc func() time.Time
}

func newPlanner(
	d *model.Deployment,
	lastSuccessfulCommitHash string,
	workingDir string,
	apiClient apiClient,
	gitClient gitClient,
	pipedConfig *config.PipedSpec,
	appManifestsCache cache.Cache,
	logger *zap.Logger,
) *planner {

	logger = logger.Named("planner").With(
		zap.String("deployment-id", d.Id),
		zap.String("application-id", d.ApplicationId),
		zap.String("env-id", d.EnvId),
		zap.String("project-id", d.ProjectId),
		zap.String("application-kind", d.Kind.String()),
		zap.String("working-dir", workingDir),
	)

	p := &planner{
		deployment:               d,
		lastSuccessfulCommitHash: lastSuccessfulCommitHash,
		workingDir:               workingDir,
		apiClient:                apiClient,
		gitClient:                gitClient,
		pipedConfig:              pipedConfig,
		plannerRegistry:          registry.DefaultRegistry(),
		appManifestsCache:        appManifestsCache,
		nowFunc:                  time.Now,
		logger:                   logger,
	}
	return p
}

// ID returns the id of planner.
// This is the same value with deployment ID.
func (p *planner) ID() string {
	return p.deployment.Id
}

// IsDone tells whether this planner is done it tasks or not.
// Returning true means this planner can be removable.
func (p *planner) IsDone() bool {
	return p.done.Load()
}

// DoneTimestamp returns the time when planner has done.
func (p *planner) DoneTimestamp() time.Time {
	return p.doneTimestamp
}

func (p *planner) Cancel(cmd model.ReportableCommand) {
	if p.cancelled {
		return
	}
	p.cancelled = true
	p.cancelledCh <- &cmd
	close(p.cancelledCh)
}

func (p *planner) Run(ctx context.Context) error {
	p.logger.Info("start running a planner")
	defer func() {
		p.doneTimestamp = p.nowFunc()
		p.done.Store(true)
	}()

	planner, ok := p.plannerRegistry.Planner(p.deployment.Kind)
	if !ok {
		return fmt.Errorf("no registered planner for application %v", p.deployment.Kind)
	}

	// Clone repository and checkout to the target revision.
	repoDirPath := filepath.Join(p.workingDir, workspaceGitRepoDirName)
	gitRepo, err := prepareDeployRepository(ctx, p.deployment, p.gitClient, repoDirPath, p.pipedConfig)
	if err != nil {
		return err
	}

	// Load deployment configuration for this application.
	cfg, err := loadDeploymentConfiguration(gitRepo.GetPath(), p.deployment)
	if err != nil {
		return fmt.Errorf("failed to load deployment configuration (%v)", err)
	}
	p.deploymentConfig = cfg

	in := pln.Input{
		Deployment:                     p.deployment,
		MostRecentSuccessfulCommitHash: p.lastSuccessfulCommitHash,
		DeploymentConfig:               cfg,
		Repo:                           gitRepo,
		RepoDir:                        gitRepo.GetPath(),
		AppDir:                         filepath.Join(gitRepo.GetPath(), p.deployment.GitPath.Path),
		AppManifestsCache:              p.appManifestsCache,
		Logger:                         p.logger,
	}
	out, err := planner.Plan(ctx, in)

	// If the deployment was already cancelled, we ignore the plan result.
	select {
	case cmd := <-p.cancelledCh:
		if cmd != nil {
			return p.reportDeploymentCancelled(ctx, cmd)
		}
	default:
	}

	if err == nil {
		return p.reportDeploymentPlanned(ctx, out)
	}
	return p.reportDeploymentCompleted(ctx, model.DeploymentStatus_DEPLOYMENT_FAILURE, err.Error())
}

func (p *planner) reportDeploymentCancelled(ctx context.Context, cmd *model.ReportableCommand) error {
	var (
		status = model.DeploymentStatus_DEPLOYMENT_CANCELLED
		desc   = fmt.Sprintf("Deployment was cancelled by %s while planning", cmd.Commander)
	)
	p.reportDeploymentCompleted(ctx, status, desc)
	return cmd.Report(ctx, model.CommandStatus_COMMAND_SUCCEEDED, nil)
}

func (p *planner) reportDeploymentPlanned(ctx context.Context, out pln.Output) error {
	var (
		err   error
		retry = pipedservice.NewRetry(10)
		req   = &pipedservice.ReportDeploymentPlannedRequest{
			DeploymentId:      p.deployment.Id,
			Description:       out.Description,
			StatusDescription: "Deployment pipeline has been planned",
			Stages:            out.Stages,
		}
	)

	for retry.WaitNext(ctx) {
		if _, err = p.apiClient.ReportDeploymentPlanned(ctx, req); err == nil {
			return nil
		}
		err = fmt.Errorf("failed to report deployment status to control-plane: %v", err)
	}

	if err != nil {
		p.logger.Error("failed to mark deployment to be planned", zap.Error(err))
	}
	return err
}

func (p *planner) reportDeploymentCompleted(ctx context.Context, status model.DeploymentStatus, desc string) error {
	var (
		err error
		now = p.nowFunc()
		req = &pipedservice.ReportDeploymentCompletedRequest{
			DeploymentId:      p.deployment.Id,
			Status:            status,
			StatusDescription: desc,
			StageStatuses:     nil,
			CompletedAt:       now.Unix(),
		}
		retry = pipedservice.NewRetry(10)
	)

	for retry.WaitNext(ctx) {
		if _, err = p.apiClient.ReportDeploymentCompleted(ctx, req); err == nil {
			return nil
		}
		err = fmt.Errorf("failed to report deployment status to control-plane: %v", err)
	}

	if err != nil {
		p.logger.Error("failed to mark deployment to be failed", zap.Error(err))
	}
	return err
}
