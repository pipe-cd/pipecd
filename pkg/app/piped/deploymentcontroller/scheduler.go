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

package deploymentcontroller

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/kapetaniosci/pipe/pkg/app/api/service/pipedservice"
	"github.com/kapetaniosci/pipe/pkg/app/piped/executor"
	"github.com/kapetaniosci/pipe/pkg/app/piped/logpersister"
	"github.com/kapetaniosci/pipe/pkg/config"
	"github.com/kapetaniosci/pipe/pkg/model"
)

var (
	workspaceGitRepoDirName = "repo"
	workspaceStagesDirName  = "stages"
)

type repoStore interface {
	CloneReadOnlyRepo(repo, branch, revision string) (string, error)
}

type executorRegistry interface {
	Executor(stage model.Stage, in executor.Input) (executor.Executor, error)
}

// scheduler is a dedicated object for a specific deployment of a single application.
type scheduler struct {
	deployment        *model.Deployment
	workingDir        string
	executorRegistry  executorRegistry
	apiClient         apiClient
	gitClient         gitClient
	commandStore      commandStore
	logPersister      logpersister.Persister
	metadataPersister metadataPersister
	pipedConfig       *config.PipedSpec
	logger            *zap.Logger

	// Deployment configuration for this application.
	appConfig *config.Config
	done      atomic.Bool
	nowFunc   func() time.Time
}

func newScheduler(
	d *model.Deployment,
	workingDir string,
	apiClient apiClient,
	gitClient gitClient,
	cmdStore commandStore,
	lp logpersister.Persister,
	mdp metadataPersister,
	pipedConfig *config.PipedSpec,
	logger *zap.Logger,
) *scheduler {

	logger = logger.Named("scheduler").With(
		zap.String("deployment-id", d.Id),
		zap.String("application-id", d.ApplicationId),
		zap.String("env-id", d.EnvId),
		zap.String("project-id", d.ProjectId),
		zap.String("application-kind", d.Kind.String()),
		zap.String("working-dir", workingDir),
	)
	return &scheduler{
		deployment:        d,
		pipedConfig:       pipedConfig,
		workingDir:        workingDir,
		executorRegistry:  executor.DefaultRegistry(),
		apiClient:         apiClient,
		gitClient:         gitClient,
		commandStore:      cmdStore,
		logPersister:      lp,
		metadataPersister: mdp,
		logger:            logger,
		nowFunc:           time.Now,
	}
}

// Id returns the id of scheduler.
// This is the same value with deployment ID.
func (s *scheduler) Id() string {
	return s.deployment.Id
}

// IsDone tells whether this scheduler is done it tasks or not.
// Returning true means this scheduler can be removable.
func (s *scheduler) IsDone() bool {
	return s.done.Load()
}

// Run starts running the scheduler.
func (s *scheduler) Run(ctx context.Context) (executeErr error) {
	s.logger.Info("start running a scheduler")
	defer func() {
		s.done.Store(true)
	}()

	defer func() {
		executeErr = s.end(ctx, executeErr)
		if executeErr != nil {
			s.logger.Error("a scheduler has been failed at end phase", zap.Error(executeErr))
		}
		s.logger.Info("a scheduler has been completed successfully")
	}()

	executeErr = s.start(ctx)
	if executeErr != nil {
		s.logger.Error("a scheduler has been failed at start phase", zap.Error(executeErr))
		return executeErr
	}

	executeErr = s.run(ctx)
	if executeErr != nil {
		s.logger.Error("a scheduler has been failed at run phase", zap.Error(executeErr))
		return executeErr
	}

	return nil
}

// start does all needed things before start executing the deployment.
func (s *scheduler) start(ctx context.Context) error {
	lp := s.logPersister.StageLogPersister(s.deployment.Id, model.StageStart.String())
	defer lp.Complete(ctx)

	// Update deployment status to RUNNING if needed.
	if s.deployment.CanUpdateStatus(model.DeploymentStatus_DEPLOYMENT_RUNNING) {
		err := s.reportDeploymentStatus(ctx, model.DeploymentStatus_DEPLOYMENT_RUNNING, "piped started handling this deployment")
		if err != nil {
			lp.AppendError(err.Error())
			return err
		}
	}

	// Clone repository and checkout to the target revision.
	var (
		appID       = s.deployment.ApplicationId
		repoID      = s.deployment.GitPath.RepoId
		repoDirPath = filepath.Join(s.workingDir, workspaceGitRepoDirName)
		revision    = s.deployment.Trigger.Commit.Revision
		repoCfg, ok = s.pipedConfig.GetRepository(repoID)
	)
	if !ok {
		err := fmt.Errorf("no registered repository id %s for application %s", repoID, appID)
		lp.AppendError(err.Error())
		return err
	}

	gitRepo, err := s.gitClient.Clone(ctx, repoCfg.RepoID, repoCfg.Remote, repoCfg.Branch, repoDirPath)
	if err != nil {
		err = fmt.Errorf("failed to clone repository %s for application %s", repoID, appID)
		lp.AppendError(err.Error())
		return err
	}

	if err = gitRepo.Checkout(ctx, revision); err != nil {
		err = fmt.Errorf("failed to clone repository %s for application %s", repoID, appID)
		lp.AppendError(err.Error())
		return err
	}
	lp.AppendSuccess(fmt.Sprintf("successfully cloned repository %s", repoID))

	// Restore previous executed state from deployment data.

	return nil
}

// run starts running from previous state.
func (s *scheduler) run(ctx context.Context) error {
	// Loop until one of the following conditions occurs:
	// - context has done
	// - no stage to execute
	// - executing stage has completed with an error
	// Determine the next stage that should be executed.
	var (
		stageName = model.Stage("")
		input     = executor.Input{
			Deployment:        s.deployment,
			AppConfig:         s.appConfig,
			WorkingDir:        s.workingDir,
			CommandStore:      s.commandStore,
			LogPersister:      s.logPersister.StageLogPersister("", ""),
			MetadataPersister: s.metadataPersister.StageMetadataPersister("", ""),
			Logger:            s.logger,
		}
	)
	ex, err := s.executorRegistry.Executor(stageName, input)
	if err != nil {
		return nil
	}
	_, err = ex.Execute(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *scheduler) end(ctx context.Context, executeErr error) error {
	// We check the runErr to decide adding a ROLLBACK stage or not.
	var (
		status     = model.DeploymentStatus_DEPLOYMENT_SUCCESS
		statusDesc string
	)
	if executeErr != nil {
		status = model.DeploymentStatus_DEPLOYMENT_FAILURE
	}
	return s.reportDeploymentStatus(ctx, status, statusDesc)
}

func (s *scheduler) determineNextStages() []string {
	return nil
}

func (s *scheduler) reportStageStatus(ctx context.Context, stageID string, status model.StageStatus) error {
	var (
		now = s.nowFunc()
		req = &pipedservice.ReportStageStatusChangedRequest{
			DeploymentId: s.deployment.Id,
			StageId:      stageID,
			Status:       status,
			CompletedAt:  now.Unix(),
		}
	)
	// TODO: Do this with exponential backoff.
	_, err := s.apiClient.ReportStageStatusChanged(ctx, req)
	if err != nil {
		err = fmt.Errorf("failed to report stage status to control-plane: %v", err)
	}

	// Update local deployment stage status?
	return err
}

func (s *scheduler) reportDeploymentStatus(ctx context.Context, status model.DeploymentStatus, desc string) error {
	var (
		now = s.nowFunc()
		req = &pipedservice.ReportDeploymentStatusChangedRequest{
			DeploymentId:      s.deployment.Id,
			Status:            status,
			StatusDescription: desc,
			CompletedAt:       now.Unix(),
		}
	)
	// TODO: Do this with exponential backoff.
	_, err := s.apiClient.ReportDeploymentStatusChanged(ctx, req)
	if err != nil {
		err = fmt.Errorf("failed to report deployment status to control-plane: %v", err)
	}

	// Update local deployment status?
	return err
}
