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
	workspaceGitRepoDirName  = "repo"
	workspaceStagesDirName   = "stages"
	defaultDeploymentTimeout = time.Hour
)

type repoStore interface {
	CloneReadOnlyRepo(repo, branch, revision string) (string, error)
}

// scheduler is a dedicated object for a specific deployment of a single application.
type scheduler struct {
	// Readonly deployment model.
	deployment       *model.Deployment
	workingDir       string
	executorRegistry registry.Registry
	apiClient        apiClient
	gitClient        gitClient
	commandLister    commandLister
	logPersister     logpersister.Persister
	metadataStore    *metadataStore
	pipedConfig      *config.PipedSpec
	logger           *zap.Logger

	deploymentConfig *config.Config
	prepareOnce      sync.Once
	// Current status of each stages.
	// We stores their current statuses into this field
	// because the deployment model is readonly to avoid data race.
	// We may need a mutex for this field in the future
	// when the stages can be executed concurrently.
	stageStatuses map[string]model.StageStatus
	done          atomic.Bool
	reported      atomic.Bool
	nowFunc       func() time.Time
}

func newScheduler(
	d *model.Deployment,
	workingDir string,
	apiClient apiClient,
	gitClient gitClient,
	commandLister commandLister,
	lp logpersister.Persister,
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

	s := &scheduler{
		deployment:       d,
		pipedConfig:      pipedConfig,
		workingDir:       workingDir,
		executorRegistry: registry.DefaultRegistry(),
		apiClient:        apiClient,
		gitClient:        gitClient,
		commandLister:    commandLister,
		logPersister:     lp,
		metadataStore:    NewMetadataStore(apiClient, d),
		logger:           logger,
		nowFunc:          time.Now,
	}

	// Initialize the map of current status of all stages.
	s.stageStatuses = make(map[string]model.StageStatus, len(d.Stages))
	for _, stage := range d.Stages {
		s.stageStatuses[stage.Id] = stage.Status
	}

	return s
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
		err := s.reportDeploymentRunning(ctx, "The piped started handling this deployment")
		if err != nil {
			return err
		}
	}

	var (
		deploymentStatus  = model.DeploymentStatus_DEPLOYMENT_SUCCESS
		statusDescription = "Completed Successfully"
		cancelCommand     *model.ReportableCommand
		cancelledCh       = make(chan struct{})
		timer             = time.NewTimer(defaultDeploymentTimeout)
	)
	defer timer.Stop()

	// Watch the cancel command from command lister.
	// TODO: In the future we may want to change the design of command lister
	// to support subscribing a specific command type.
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				commands := s.commandLister.ListDeploymentCommands(s.deployment.Id)
				for _, cmd := range commands {
					c := cmd.GetCancelDeployment()
					if c == nil {
						continue
					}
					cancelCommand = &cmd
					close(cancelledCh)
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Iterate all the stages and execute the uncompleted ones.
	for _, ps := range s.deployment.Stages {
		if ps.Status == model.StageStatus_STAGE_SUCCESS {
			continue
		}

		// This stage is already completed by a previous scheduler.
		if ps.Status == model.StageStatus_STAGE_CANCELLED {
			deploymentStatus = model.DeploymentStatus_DEPLOYMENT_CANCELLED
			statusDescription = fmt.Sprintf("Deployment was cancelled while executing stage %s", ps.Id)
			break
		}
		if ps.Status == model.StageStatus_STAGE_FAILURE {
			deploymentStatus = model.DeploymentStatus_DEPLOYMENT_FAILURE
			statusDescription = fmt.Sprintf("Failed while executing stage %s", ps.Id)
			break
		}

		var (
			result       model.StageStatus
			sig, handler = executor.NewStopSignal()
			doneCh       = make(chan struct{})
		)
		go func() {
			result = s.executeStage(sig, ps)
			close(doneCh)
		}()

		select {
		case <-ctx.Done():
			handler.Terminate()
			<-doneCh

		case <-timer.C:
			handler.Timeout()
			<-doneCh

		case <-cancelledCh:
			handler.Cancel()
			<-doneCh

		case <-doneCh:
			break
		}

		// If all operations of the stage were completed successfully
		// go the next stage to handle.
		if result == model.StageStatus_STAGE_SUCCESS {
			continue
		}

		sigType := sig.Signal()

		// The deployment was cancelled by a web user.
		if sigType == executor.StopSignalCancel {
			deploymentStatus = model.DeploymentStatus_DEPLOYMENT_CANCELLED
			statusDescription = fmt.Sprintf("Deployment was cancelled while executing stage %s", ps.Id)
			break
		}

		// The stage was failed but not caused by the stop signal.
		if result == model.StageStatus_STAGE_FAILURE && sigType == executor.StopSignalNone {
			deploymentStatus = model.DeploymentStatus_DEPLOYMENT_FAILURE
			statusDescription = fmt.Sprintf("Failed while executing stage %s", ps.Id)
			break
		}

		return nil
	}

	if model.IsCompletedDeployment(deploymentStatus) {
		s.reportDeploymentCompleted(ctx, deploymentStatus, statusDescription)
	}

	if cancelCommand != nil {
		cancelCommand.Report(ctx, model.CommandStatus_COMMAND_SUCCEEDED, nil)
	}

	return nil
}

// executeStage finds the executor for the given stage and execute.
func (s *scheduler) executeStage(sig executor.StopSignal, ps *model.PipelineStage) model.StageStatus {
	var (
		ctx            = sig.Context()
		originalStatus = ps.Status
		lp             = s.logPersister.StageLogPersister(s.deployment.Id, ps.Id)
	)
	defer lp.Complete(time.Minute)

	// Update stage status to RUNNING if needed.
	if model.CanUpdateStageStatus(ps.Status, model.StageStatus_STAGE_RUNNING) {
		if err := s.reportStageStatus(ctx, ps.Id, model.StageStatus_STAGE_RUNNING); err != nil {
			return model.StageStatus_STAGE_FAILURE
		}
		originalStatus = model.StageStatus_STAGE_RUNNING
	}

	// Ensure that all needed things has been prepared before executing any stage.
	if err := s.ensurePreparing(ctx, lp); err != nil {
		if !sig.Stopped() {
			s.reportStageStatus(ctx, ps.Id, model.StageStatus_STAGE_FAILURE)
			return model.StageStatus_STAGE_FAILURE
		}
		return originalStatus
	}

	input := executor.Input{
		Stage:            ps,
		Deployment:       s.deployment,
		DeploymentConfig: s.deploymentConfig,
		PipedConfig:      s.pipedConfig,
		WorkingDir:       s.workingDir,
		RepoDir:          filepath.Join(s.workingDir, workspaceGitRepoDirName),
		StageWorkingDir:  filepath.Join(s.workingDir, workspaceStagesDirName, ps.Id),
		CommandLister: stageCommandLister{
			lister:       s.commandLister,
			deploymentID: s.deployment.Id,
			stageID:      ps.Id,
		},
		LogPersister:  lp,
		MetadataStore: s.metadataStore,
		Logger:        s.logger,
	}

	// Find the executor for this stage.
	ex, ok := s.executorRegistry.Executor(model.Stage(ps.Name), input)
	if !ok {
		err := fmt.Errorf("no registered executor for stage %s", ps.Name)
		lp.AppendError(err.Error())
		s.reportStageStatus(ctx, ps.Id, model.StageStatus_STAGE_FAILURE)
		return model.StageStatus_STAGE_FAILURE
	}

	// Start running executor.
	status := ex.Execute(sig)

	if status == model.StageStatus_STAGE_SUCCESS ||
		status == model.StageStatus_STAGE_CANCELLED ||
		(status == model.StageStatus_STAGE_FAILURE && !sig.Stopped()) {

		s.reportStageStatus(ctx, ps.Id, status)
		return status
	}

	return originalStatus
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
	s.stageStatuses[stageID] = status

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

func (s *scheduler) reportDeploymentRunning(ctx context.Context, desc string) error {
	var (
		err   error
		retry = newAPIRetry(10)
		req   = &pipedservice.ReportDeploymentRunningRequest{
			DeploymentId:      s.deployment.Id,
			StatusDescription: desc,
		}
	)

	// Update deployment status on remote.
	for retry.WaitNext(ctx) {
		if _, err = s.apiClient.ReportDeploymentRunning(ctx, req); err == nil {
			break
		}
		err = fmt.Errorf("failed to report deployment status to control-plane: %v", err)
	}

	return err
}

func (s *scheduler) reportDeploymentCompleted(ctx context.Context, status model.DeploymentStatus, desc string) error {
	var (
		err error
		now = s.nowFunc()
		req = &pipedservice.ReportDeploymentCompletedRequest{
			DeploymentId:      s.deployment.Id,
			Status:            status,
			StatusDescription: desc,
			StageStatuses:     s.stageStatuses,
			CompletedAt:       now.Unix(),
		}
		retry = newAPIRetry(10)
	)

	// Update deployment status on remote.
	for retry.WaitNext(ctx) {
		if _, err = s.apiClient.ReportDeploymentCompleted(ctx, req); err == nil {
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

type stageCommandLister struct {
	lister       commandLister
	deploymentID string
	stageID      string
}

func (s stageCommandLister) ListCommands() []model.ReportableCommand {
	return s.lister.ListStageCommands(s.deploymentID, s.stageID)
}
