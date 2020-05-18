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
	"sync"
	"time"

	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/kapetaniosci/pipe/pkg/app/api/service/pipedservice"
	"github.com/kapetaniosci/pipe/pkg/app/piped/executor"
	"github.com/kapetaniosci/pipe/pkg/app/piped/executor/registry"
	"github.com/kapetaniosci/pipe/pkg/app/piped/logpersister"
	"github.com/kapetaniosci/pipe/pkg/backoff"
	"github.com/kapetaniosci/pipe/pkg/config"
	"github.com/kapetaniosci/pipe/pkg/git"
	"github.com/kapetaniosci/pipe/pkg/model"
)

var (
	workspaceGitRepoDirName = "repo"
	workspaceStagesDirName  = "stages"
)

type repoStore interface {
	CloneReadOnlyRepo(repo, branch, revision string) (string, error)
}

// scheduler is a dedicated object for a specific deployment of a single application.
type scheduler struct {
	deployment        *model.Deployment
	workingDir        string
	executorRegistry  registry.Registry
	apiClient         apiClient
	gitClient         gitClient
	commandStore      commandStore
	logPersister      logpersister.Persister
	metadataPersister metadataPersister
	pipedConfig       *config.PipedSpec
	logger            *zap.Logger

	// Deployment configuration for this application.
	deploymentConfig *config.Config
	prepareOnce      sync.Once
	done             atomic.Bool
	reported         atomic.Bool
	nowFunc          func() time.Time
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
		executorRegistry:  registry.DefaultRegistry(),
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

// IsReported tells whether this scheduler has already reported its state
// to the control-plane.
func (s *scheduler) IsReported() bool {
	return s.reported.Load()
}

// Run starts running the scheduler.
// It determines what stage should be executed next by which executor.
// The returning error does not mean that the pipeline was failed,
// but it means that the scheduler could not finish its job normally.
func (s *scheduler) Run(ctx context.Context) error {
	s.logger.Info("start running a scheduler")
	defer func() {
		s.done.Store(true)
	}()

	// If this deployment is already completed. Do nothing.
	if model.IsCompletedDeployment(s.deployment.Status) {
		s.logger.Info("this deployment is already completed")
		return nil
	}

	// Update deployment status to RUNNING if needed.
	if model.CanUpdateDeploymentStatus(s.deployment.Status, model.DeploymentStatus_DEPLOYMENT_RUNNING) {
		err := s.reportDeploymentStatus(ctx, model.DeploymentStatus_DEPLOYMENT_RUNNING, "The piped started handling this deployment")
		if err != nil {
			return err
		}
	}

	// Iterate all the stages and execute the uncompleted ones.
	for _, ps := range s.deployment.Stages {
		// This stage is already handed by a previous scheduler.
		if model.IsCompletedStage(ps.Status) {
			continue
		}

		status := s.executeStage(ctx, ps)
		if status == model.StageStatus_STAGE_SUCCESS {
			continue
		}
		if status == model.StageStatus_STAGE_FAILURE {
			s.reportDeploymentStatus(ctx, model.DeploymentStatus_DEPLOYMENT_FAILURE, fmt.Sprintf("Failed while executing stage %s", ps.Id))
			return nil
		}
		if status == model.StageStatus_STAGE_CANCELLED {
			s.reportDeploymentStatus(ctx, model.DeploymentStatus_DEPLOYMENT_CANCELLED, fmt.Sprintf("Deployment was cancelled while executing stage %s", ps.Id))
			return nil
		}
		return nil
	}

	s.reportDeploymentStatus(ctx, model.DeploymentStatus_DEPLOYMENT_SUCCESS, "")
	return nil
}

// executeStage finds the executor for the given stage and execute.
func (s *scheduler) executeStage(ctx context.Context, ps *model.PipelineStage) (status model.StageStatus) {
	lp := s.logPersister.StageLogPersister(s.deployment.Id, ps.Id)
	defer lp.Complete(ctx)

	// Update stage status to RUNNING if needed.
	if model.CanUpdateStageStatus(ps.Status, model.StageStatus_STAGE_RUNNING) {
		if err := s.reportStageStatus(ctx, ps.Id, model.StageStatus_STAGE_RUNNING); err != nil {
			status = model.StageStatus_STAGE_FAILURE
			return
		}
		status = model.StageStatus_STAGE_RUNNING
	}

	// Ensure that all needed things has been prepared before executing any stage.
	if err := s.ensurePreparing(ctx, lp); err != nil {
		status = model.StageStatus_STAGE_FAILURE
		return
	}

	input := executor.Input{
		Stage:             ps,
		Deployment:        s.deployment,
		DeploymentConfig:  s.deploymentConfig,
		PipedConfig:       s.pipedConfig,
		WorkingDir:        s.workingDir,
		RepoDir:           filepath.Join(s.workingDir, workspaceGitRepoDirName),
		StageWorkingDir:   filepath.Join(s.workingDir, workspaceStagesDirName, ps.Id),
		CommandStore:      s.commandStore,
		LogPersister:      lp,
		MetadataPersister: s.metadataPersister.StageMetadataPersister(s.deployment.Id, ps.Id),
		Logger:            s.logger,
	}

	// Find the executor for this stage.
	ex, ok := s.executorRegistry.Executor(model.Stage(ps.Name), input)
	if !ok {
		err := fmt.Errorf("no registered executor for stage %s", ps.Name)
		lp.AppendError(err.Error())
		s.reportStageStatus(ctx, ps.Id, model.StageStatus_STAGE_FAILURE)
		status = model.StageStatus_STAGE_FAILURE
		return
	}

	// Start running executor.
	status = ex.Execute(ctx)
	if err := s.reportStageStatus(ctx, ps.Id, status); err != nil {
		return
	}

	return status
}

// ensurePreparing ensures that all needed things should be prepared before executing any stages.
// The log of this preparing process will be written to the first executing stage
// when a new scheduler has been created.
func (s *scheduler) ensurePreparing(ctx context.Context, lp logpersister.StageLogPersister) error {
	var err error
	s.prepareOnce.Do(func() {
		lp.AppendInfo("START PREPARING")
		lp.AppendInfo("new scheduler has been created for this deployment so we need some preparation")

		// Clone repository and checkout to the target revision.
		var (
			appID       = s.deployment.ApplicationId
			repoID      = s.deployment.GitPath.RepoId
			repoDirPath = filepath.Join(s.workingDir, workspaceGitRepoDirName)
			revision    = s.deployment.Trigger.Commit.Hash
			repoCfg, ok = s.pipedConfig.GetRepository(repoID)
		)
		if !ok {
			err = fmt.Errorf("no registered repository id %s for application %s", repoID, appID)
			lp.AppendError(err.Error())
			return
		}

		var gitRepo git.Repo
		gitRepo, err = s.gitClient.Clone(ctx, repoCfg.RepoID, repoCfg.Remote, repoCfg.Branch, repoDirPath)
		if err != nil {
			err = fmt.Errorf("failed to clone repository %s for application %s (%v)", repoID, appID, err)
			lp.AppendError(err.Error())
			return
		}

		err = gitRepo.Checkout(ctx, revision)
		if err != nil {
			err = fmt.Errorf("failed to clone repository %s for application %s (%v)", repoID, appID, err)
			lp.AppendError(err.Error())
			return
		}
		lp.AppendSuccess(fmt.Sprintf("successfully cloned repository %s", repoID))

		// Load deployment configuration for this application.
		var cfg *config.Config
		cfg, err = s.loadDeploymentConfiguration(ctx, gitRepo.GetPath(), s.deployment)
		if err != nil {
			err = fmt.Errorf("failed to load deployment configuration (%v)", err)
		}
		s.deploymentConfig = cfg
		lp.AppendSuccess("successfully loaded deployment configuration")
		lp.AppendInfo("PREPARING COMPLETED")
	})
	return err
}

func (s *scheduler) reportStageStatus(ctx context.Context, stageID string, status model.StageStatus) error {
	var (
		err error
		now = s.nowFunc()
		req = &pipedservice.ReportStageStatusChangedRequest{
			DeploymentId: s.deployment.Id,
			StageId:      stageID,
			Status:       status,
			CompletedAt:  now.Unix(),
		}
		retry = newAPIRetry(10)
	)

	// Update stage status at local.
	for _, stage := range s.deployment.Stages {
		if stage.Id == stageID {
			stage.Status = status
			break
		}
	}

	// Update stage status on the remote.
	for retry.WaitNext(ctx) {
		_, err = s.apiClient.ReportStageStatusChanged(ctx, req)
		if err == nil {
			break
		}
		err = fmt.Errorf("failed to report stage status to control-plane: %v", err)
	}

	return err
}

func (s *scheduler) reportDeploymentStatus(ctx context.Context, status model.DeploymentStatus, desc string) error {
	var (
		err error
		now = s.nowFunc()
		req = &pipedservice.ReportDeploymentStatusChangedRequest{
			DeploymentId:      s.deployment.Id,
			Status:            status,
			StatusDescription: desc,
			StageStatuses:     s.deployment.StageStatusMap(),
			CompletedAt:       now.Unix(),
		}
		retry = newAPIRetry(10)
	)

	// Update deployment status at local.
	s.deployment.Status = status
	s.deployment.StatusDescription = desc

	// Update deployment status on remote.
	for retry.WaitNext(ctx) {
		_, err = s.apiClient.ReportDeploymentStatusChanged(ctx, req)
		if err == nil {
			break
		}
		err = fmt.Errorf("failed to report deployment status to control-plane: %v", err)
	}

	if err == nil && model.IsCompletedDeployment(status) {
		s.reported.Store(true)
	}
	return err
}

func (s *scheduler) loadDeploymentConfiguration(ctx context.Context, repoPath string, d *model.Deployment) (*config.Config, error) {
	path := filepath.Join(repoPath, d.GitPath.GetDeploymentConfigFilePath(config.DeploymentConfigurationFileName))
	cfg, err := config.LoadFromYAML(path)
	if err != nil {
		return nil, err
	}
	if appKind, ok := config.ToApplicationKind(cfg.Kind); !ok || appKind != d.Kind {
		return nil, fmt.Errorf("application in deployment configuration file is not match, got: %s, expected: %s", appKind, d.Kind)
	}
	return cfg, nil
}

// 0s 997.867435ms 2.015381172s 3.485134345s 4.389600179s 18.118099328s 48.73058264s
func newAPIRetry(maxRetries int) backoff.Retry {
	bo := backoff.NewExponential(2*time.Second, time.Minute)
	return backoff.NewRetry(maxRetries, bo)
}
