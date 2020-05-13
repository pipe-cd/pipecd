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
	"github.com/kapetaniosci/pipe/pkg/app/piped/executor/registry"
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
	done             atomic.Bool
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

// Run starts running the scheduler.
// It determines what stage should be executed next by which executor.
// The returning error does not mean that the pipeline was failed,
// but it means that the scheduler could not finish its job normally.
func (s *scheduler) Run(ctx context.Context) error {
	s.logger.Info("start running a scheduler")
	defer func() {
		s.done.Store(true)
	}()

	for _, ps := range s.deployment.Stages {
		if ps.Id == model.StageStart.String() {
			if err := s.executeStartStage(ctx); err != nil {
				return err
			}
			continue
		}
		if ps.Id == model.StageEnd.String() {
			if err := s.executeEndStage(ctx, nil); err != nil {
				return err
			}
			continue
		}

		// Handle user specified stages.
		input := executor.Input{
			Stage:             ps,
			Deployment:        s.deployment,
			DeploymentConfig:  s.deploymentConfig,
			PipedConfig:       s.pipedConfig,
			WorkingDir:        s.workingDir,
			CommandStore:      s.commandStore,
			LogPersister:      s.logPersister.StageLogPersister(s.deployment.Id, ps.Id),
			MetadataPersister: s.metadataPersister.StageMetadataPersister(s.deployment.Id, ps.Id),
			Logger:            s.logger,
		}
		ex, err := s.executorRegistry.Executor(model.Stage(ps.Name), input)
		if err != nil {
			s.logger.Error("no executor", zap.Error(err))
			return err
		}
		_, err = ex.Execute(ctx)
		if err != nil {
			s.logger.Error("failed to execute", zap.Error(err))
			return err
		}
	}

	return nil
}

// executeStartStage does all needed things before start executing the deployment.
func (s *scheduler) executeStartStage(ctx context.Context) error {
	lp := s.logPersister.StageLogPersister(s.deployment.Id, model.StageStart.String())
	defer lp.Complete(ctx)

	lp.AppendInfo("new scheduler has been created for this deployment")

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
		err = fmt.Errorf("failed to clone repository %s for application %s (%v)", repoID, appID, err)
		lp.AppendError(err.Error())
		return err
	}

	if err = gitRepo.Checkout(ctx, revision); err != nil {
		err = fmt.Errorf("failed to clone repository %s for application %s (%v)", repoID, appID, err)
		lp.AppendError(err.Error())
		return err
	}
	lp.AppendSuccess(fmt.Sprintf("successfully cloned repository %s", repoID))

	// Load deployment configuration for this application.
	cfg, err := s.loadDeploymentConfiguration(ctx, gitRepo.GetPath(), s.deployment)
	if err != nil {
		err = fmt.Errorf("failed to load deployment configuration (%v)", err)
	}
	s.deploymentConfig = cfg
	lp.AppendSuccess("successfully loaded deployment configuration")

	return nil
}

func (s *scheduler) executeEndStage(ctx context.Context, executeErr error) error {
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

func (s *scheduler) loadDeploymentConfiguration(ctx context.Context, repoPath string, d *model.Deployment) (*config.Config, error) {
	path := filepath.Join(repoPath, d.GetDeploymentConfigFilePath(config.DeploymentConfigurationFileName))
	cfg, err := config.LoadFromYAML(path)
	if err != nil {
		return nil, err
	}
	if appKind, ok := config.ToApplicationKind(cfg.Kind); !ok || appKind != d.Kind {
		return nil, fmt.Errorf("application in deployment configuration file is not match, got: %s, expected: %s", appKind, d.Kind)
	}
	return cfg, nil
}
