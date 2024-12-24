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

package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/piped/controller/controllermetrics"
	"github.com/pipe-cd/pipecd/pkg/app/piped/deploysource"
	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	"github.com/pipe-cd/pipecd/pkg/app/piped/executor/registry"
	"github.com/pipe-cd/pipecd/pkg/app/piped/logpersister"
	"github.com/pipe-cd/pipecd/pkg/app/piped/metadatastore"
	pln "github.com/pipe-cd/pipecd/pkg/app/piped/planner"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

// scheduler is a dedicated object for a specific deployment of a single application.
type scheduler struct {
	// Readonly deployment model.
	deployment          *model.Deployment
	workingDir          string
	executorRegistry    registry.Registry
	apiClient           apiClient
	gitClient           gitClient
	commandLister       commandLister
	applicationLister   applicationLister
	liveResourceLister  liveResourceLister
	analysisResultStore analysisResultStore
	logPersister        logpersister.Persister
	metadataStore       metadatastore.MetadataStore
	notifier            notifier
	secretDecrypter     secretDecrypter
	pipedConfig         *config.PipedSpec
	appManifestsCache   cache.Cache
	logger              *zap.Logger
	tracer              trace.Tracer

	targetDSP  deploysource.Provider
	runningDSP deploysource.Provider

	// Current status of each stages.
	// We stores their current statuses into this field
	// because the deployment model is readonly to avoid data race.
	// We may need a mutex for this field in the future
	// when the stages can be executed concurrently.
	stageStatuses            map[string]model.StageStatus
	genericApplicationConfig config.GenericApplicationSpec

	done                 atomic.Bool
	doneTimestamp        time.Time
	doneDeploymentStatus model.DeploymentStatus
	cancelled            bool
	cancelledCh          chan *model.ReportableCommand

	nowFunc func() time.Time
}

func newScheduler(
	d *model.Deployment,
	workingDir string,
	apiClient apiClient,
	gitClient gitClient,
	commandLister commandLister,
	applicationLister applicationLister,
	liveResourceLister liveResourceLister,
	analysisResultStore analysisResultStore,
	lp logpersister.Persister,
	notifier notifier,
	sd secretDecrypter,
	pipedConfig *config.PipedSpec,
	appManifestsCache cache.Cache,
	logger *zap.Logger,
	tracerProvider trace.TracerProvider,
) *scheduler {
	logger = logger.Named("scheduler").With(
		zap.String("deployment-id", d.Id),
		zap.String("app-id", d.ApplicationId),
		zap.String("project-id", d.ProjectId),
		zap.String("app-kind", d.Kind.String()),
		zap.String("working-dir", workingDir),
	)

	s := &scheduler{
		deployment:           d,
		workingDir:           workingDir,
		executorRegistry:     registry.DefaultRegistry(),
		apiClient:            apiClient,
		gitClient:            gitClient,
		commandLister:        commandLister,
		applicationLister:    applicationLister,
		liveResourceLister:   liveResourceLister,
		analysisResultStore:  analysisResultStore,
		logPersister:         lp,
		metadataStore:        metadatastore.NewMetadataStore(apiClient, d),
		notifier:             notifier,
		secretDecrypter:      sd,
		pipedConfig:          pipedConfig,
		appManifestsCache:    appManifestsCache,
		doneDeploymentStatus: d.Status,
		cancelledCh:          make(chan *model.ReportableCommand, 1),
		logger:               logger,
		tracer:               tracerProvider.Tracer("controller/scheduler"),
		nowFunc:              time.Now,
	}

	// Initialize the map of current status of all stages.
	s.stageStatuses = make(map[string]model.StageStatus, len(d.Stages))
	for _, stage := range d.Stages {
		s.stageStatuses[stage.Id] = stage.Status
	}

	return s
}

// ID returns the id of scheduler.
// This is the same value with deployment ID.
func (s *scheduler) ID() string {
	return s.deployment.Id
}

// CommitHash returns the hash value of deploying commit.
func (s *scheduler) CommitHash() string {
	return s.deployment.CommitHash()
}

// ConfigFilename returns the config filename of the deployment.
func (s *scheduler) ConfigFilename() string {
	return s.deployment.GitPath.GetApplicationConfigFilename()
}

// IsDone tells whether this scheduler is done it tasks or not.
// Returning true means this scheduler can be removable.
func (s *scheduler) IsDone() bool {
	return s.done.Load()
}

// DoneTimestamp returns the time when scheduler has done.
// This can be used only after IsDone() returns true.
func (s *scheduler) DoneTimestamp() time.Time {
	if !s.IsDone() {
		return time.Now().AddDate(1, 0, 0)
	}
	return s.doneTimestamp
}

// DoneDeploymentStatus returns the deployment status when scheduler has done.
// This can be used only after IsDone() returns true.
func (s *scheduler) DoneDeploymentStatus() model.DeploymentStatus {
	if !s.IsDone() {
		return s.deployment.Status
	}
	return s.doneDeploymentStatus
}

func (s *scheduler) Cancel(cmd model.ReportableCommand) {
	if s.cancelled {
		return
	}
	s.cancelled = true
	s.cancelledCh <- &cmd
	close(s.cancelledCh)
}

// Run starts running the scheduler.
// It determines what stage should be executed next by which executor.
// The returning error does not mean that the pipeline was failed,
// but it means that the scheduler could not finish its job normally.
func (s *scheduler) Run(ctx context.Context) error {
	s.logger.Info("start running scheduler")
	deploymentStatus := s.deployment.Status

	defer func() {
		s.doneTimestamp = s.nowFunc()
		s.doneDeploymentStatus = deploymentStatus
		controllermetrics.UpdateDeploymentStatus(s.deployment, deploymentStatus)
		s.done.Store(true)
	}()

	// If this deployment is already completed. Do nothing.
	if s.deployment.Status.IsCompleted() {
		s.logger.Info("this deployment is already completed")
		return nil
	}

	// Update deployment status to RUNNING if needed.
	if model.CanUpdateDeploymentStatus(s.deployment.Status, model.DeploymentStatus_DEPLOYMENT_RUNNING) {
		err := s.reportDeploymentStatusChanged(ctx, model.DeploymentStatus_DEPLOYMENT_RUNNING, "The piped started handling this deployment")
		if err != nil {
			return err
		}
		controllermetrics.UpdateDeploymentStatus(s.deployment, model.DeploymentStatus_DEPLOYMENT_RUNNING)

		// notify the deployment started event
		users, groups, err := s.getApplicationNotificationMentions(model.NotificationEventType_EVENT_DEPLOYMENT_STARTED)
		if err != nil {
			s.logger.Error("failed to get the list of users or groups", zap.Error(err))
		}

		s.notifier.Notify(model.NotificationEvent{
			Type: model.NotificationEventType_EVENT_DEPLOYMENT_STARTED,
			Metadata: &model.NotificationEventDeploymentStarted{
				Deployment:        s.deployment,
				MentionedAccounts: users,
				MentionedGroups:   groups,
			},
		})
	}

	var (
		cancelCommand   *model.ReportableCommand
		cancelCommander string
		lastStage       *model.PipelineStage
		statusReason    = "The deployment was completed successfully"
	)
	deploymentStatus = model.DeploymentStatus_DEPLOYMENT_SUCCESS

	repoCfg := config.PipedRepository{
		RepoID: s.deployment.GitPath.Repo.Id,
		Remote: s.deployment.GitPath.Repo.Remote,
		Branch: s.deployment.GitPath.Repo.Branch,
	}

	s.targetDSP = deploysource.NewProvider(
		filepath.Join(s.workingDir, "target-deploysource"),
		deploysource.NewGitSourceCloner(s.gitClient, repoCfg, "target", s.deployment.Trigger.Commit.Hash),
		*s.deployment.GitPath,
		s.secretDecrypter,
	)

	if s.deployment.RunningCommitHash != "" {
		gp := *s.deployment.GitPath
		gp.ConfigFilename = s.deployment.RunningConfigFilename

		s.runningDSP = deploysource.NewProvider(
			filepath.Join(s.workingDir, "running-deploysource"),
			deploysource.NewGitSourceCloner(s.gitClient, repoCfg, "running", s.deployment.RunningCommitHash),
			gp,
			s.secretDecrypter,
		)
	}

	// We use another deploy source provider to load the application configuration at the target commit.
	// This provider is configured with a nil secretDecrypter
	// because decrypting the sealed secrets is not required.
	// We need only the application configuration spec.
	configDSP := deploysource.NewProvider(
		filepath.Join(s.workingDir, "target-config"),
		deploysource.NewGitSourceCloner(s.gitClient, repoCfg, "target", s.deployment.Trigger.Commit.Hash),
		*s.deployment.GitPath,
		nil,
	)
	ds, err := configDSP.GetReadOnly(ctx, io.Discard)
	if err != nil {
		deploymentStatus = model.DeploymentStatus_DEPLOYMENT_FAILURE
		statusReason = fmt.Sprintf("Unable to prepare application configuration source data at target commit (%v)", err)
		s.reportDeploymentCompleted(ctx, deploymentStatus, statusReason, "")
		return err
	}
	s.genericApplicationConfig = ds.GenericApplicationConfig

	ctx, span := s.tracer.Start(
		newContextWithDeploymentSpan(ctx, s.deployment),
		"Deploy",
		trace.WithAttributes(
			attribute.String("application-id", s.deployment.ApplicationId),
			attribute.String("kind", s.deployment.Kind.String()),
			attribute.String("deployment-id", s.deployment.Id),
		))
	defer func() {
		switch deploymentStatus {
		case model.DeploymentStatus_DEPLOYMENT_SUCCESS:
			span.SetStatus(codes.Ok, statusReason)
		case model.DeploymentStatus_DEPLOYMENT_FAILURE, model.DeploymentStatus_DEPLOYMENT_CANCELLED:
			span.SetStatus(codes.Error, statusReason)
		}

		span.End()
	}()

	timer := time.NewTimer(s.genericApplicationConfig.Timeout.Duration())
	defer timer.Stop()

	// Iterate all the stages and execute the uncompleted ones.
	for i, ps := range s.deployment.Stages {
		lastStage = s.deployment.Stages[i]

		if ps.Status == model.StageStatus_STAGE_SUCCESS {
			continue
		}
		if !ps.Visible || ps.Name == model.StageRollback.String() {
			continue
		}

		// This stage is already completed by a previous scheduler.
		if ps.Status == model.StageStatus_STAGE_CANCELLED {
			deploymentStatus = model.DeploymentStatus_DEPLOYMENT_CANCELLED
			statusReason = fmt.Sprintf("Deployment was cancelled while executing stage %s", ps.Id)
			break
		}
		if ps.Status == model.StageStatus_STAGE_FAILURE {
			deploymentStatus = model.DeploymentStatus_DEPLOYMENT_FAILURE
			statusReason = fmt.Sprintf("Failed while executing stage %s", ps.Id)
			break
		}

		var (
			result       model.StageStatus
			sig, handler = executor.NewStopSignal()
			doneCh       = make(chan struct{})
		)

		go func() {
			_, span := s.tracer.Start(ctx, ps.Name, trace.WithAttributes(
				attribute.String("application-id", s.deployment.ApplicationId),
				attribute.String("kind", s.deployment.Kind.String()),
				attribute.String("deployment-id", s.deployment.Id),
				attribute.String("stage-id", ps.Id),
			))
			defer span.End()

			s.notifier.Notify(model.NotificationEvent{
				Type: model.NotificationEventType_EVENT_STAGE_STARTED,
				Metadata: &model.NotificationEventStageStarted{
					Deployment: s.deployment,
					Stage:      ps,
				},
			})

			result = s.executeStage(sig, *ps, func(in executor.Input) (executor.Executor, bool) {
				return s.executorRegistry.Executor(model.Stage(ps.Name), in)
			})

			switch result {
			case model.StageStatus_STAGE_SUCCESS, model.StageStatus_STAGE_EXITED: // Exit stage is treated as success.
				s.notifier.Notify(model.NotificationEvent{
					Type: model.NotificationEventType_EVENT_STAGE_SUCCEEDED,
					Metadata: &model.NotificationEventStageSucceeded{
						Deployment: s.deployment,
						Stage:      ps,
					},
				})
			case model.StageStatus_STAGE_FAILURE:
				s.notifier.Notify(model.NotificationEvent{
					Type: model.NotificationEventType_EVENT_STAGE_FAILED,
					Metadata: &model.NotificationEventStageFailed{
						Deployment: s.deployment,
						Stage:      ps,
					},
				})
			case model.StageStatus_STAGE_CANCELLED:
				s.notifier.Notify(model.NotificationEvent{
					Type: model.NotificationEventType_EVENT_STAGE_CANCELLED,
					Metadata: &model.NotificationEventStageCancelled{
						Deployment: s.deployment,
						Stage:      ps,
					},
				})
			case model.StageStatus_STAGE_SKIPPED:
				s.notifier.Notify(model.NotificationEvent{
					Type: model.NotificationEventType_EVENT_STAGE_SKIPPED,
					Metadata: &model.NotificationEventStageSkipped{
						Deployment: s.deployment,
						Stage:      ps,
					},
				})
			}

			switch result {
			case model.StageStatus_STAGE_SUCCESS:
				span.SetStatus(codes.Ok, statusReason)
			case model.StageStatus_STAGE_FAILURE, model.StageStatus_STAGE_CANCELLED:
				span.SetStatus(codes.Error, statusReason)
			}

			close(doneCh)
		}()

		select {
		case <-ctx.Done():
			handler.Terminate()
			<-doneCh

		case <-timer.C:
			handler.Timeout()
			<-doneCh

		case cmd := <-s.cancelledCh:
			if cmd != nil {
				cancelCommand = cmd
				cancelCommander = cmd.Commander
				handler.Cancel()
				<-doneCh
			}

		case <-doneCh:
			break
		}

		// If all operations of the stage were completed successfully or skipped by a web user
		// handle the next stage.
		if result == model.StageStatus_STAGE_SUCCESS || result == model.StageStatus_STAGE_SKIPPED {
			continue
		}

		// If the stage was completed with exited stage, exit this deployment with success.
		if result == model.StageStatus_STAGE_EXITED {
			deploymentStatus = model.DeploymentStatus_DEPLOYMENT_SUCCESS
			break
		}

		// The deployment was cancelled by a web user.
		if result == model.StageStatus_STAGE_CANCELLED {
			deploymentStatus = model.DeploymentStatus_DEPLOYMENT_CANCELLED
			statusReason = fmt.Sprintf("Cancelled by %s while executing stage %s", cancelCommander, ps.Id)
			break
		}

		if result == model.StageStatus_STAGE_FAILURE {
			deploymentStatus = model.DeploymentStatus_DEPLOYMENT_FAILURE
			// The stage was failed because of timing out.
			if sig.Signal() == executor.StopSignalTimeout {
				statusReason = fmt.Sprintf("Timed out while executing stage %s", ps.Id)
			} else {
				statusReason = fmt.Sprintf("Failed while executing stage %s", ps.Id)
			}
			break
		}

		// The deployment was cancelled at the previous stage and this stage was terminated before run.
		if result == model.StageStatus_STAGE_NOT_STARTED_YET && cancelCommand != nil {
			deploymentStatus = model.DeploymentStatus_DEPLOYMENT_CANCELLED
			statusReason = fmt.Sprintf("Cancelled by %s while executing the previous stage of %s", cancelCommander, ps.Id)
			break
		}

		s.logger.Info("stop scheduler because of temination signal", zap.String("stage-id", ps.Id))
		return nil
	}

	// When the deployment has completed but not successful,
	// we start rollback stage if the auto-rollback option is true.
	if deploymentStatus == model.DeploymentStatus_DEPLOYMENT_CANCELLED ||
		deploymentStatus == model.DeploymentStatus_DEPLOYMENT_FAILURE {

		if rollbackStages, ok := s.deployment.FindRollbackStages(); ok {
			// Update to change deployment status to ROLLING_BACK.
			if err := s.reportDeploymentStatusChanged(ctx, model.DeploymentStatus_DEPLOYMENT_ROLLING_BACK, statusReason); err != nil {
				return err
			}

			for _, stage := range rollbackStages {
				// Start running rollback stage.
				var (
					sig, handler = executor.NewStopSignal()
					doneCh       = make(chan struct{})
				)
				go func() {
					rbs := *stage
					rbs.Requires = []string{lastStage.Id}

					_, span := s.tracer.Start(ctx, rbs.Name, trace.WithAttributes(
						attribute.String("application-id", s.deployment.ApplicationId),
						attribute.String("kind", s.deployment.Kind.String()),
						attribute.String("deployment-id", s.deployment.Id),
						attribute.String("stage-id", rbs.Id),
					))
					defer span.End()

					s.notifier.Notify(model.NotificationEvent{
						Type: model.NotificationEventType_EVENT_STAGE_STARTED,
						Metadata: &model.NotificationEventStageStarted{
							Deployment: s.deployment,
							Stage:      &rbs,
						},
					})

					result := s.executeStage(sig, rbs, func(in executor.Input) (executor.Executor, bool) {
						return s.executorRegistry.RollbackExecutor(s.deployment.Kind, in)
					})

					switch result {
					case model.StageStatus_STAGE_SUCCESS, model.StageStatus_STAGE_EXITED: // Exit stage is treated as success.
						s.notifier.Notify(model.NotificationEvent{
							Type: model.NotificationEventType_EVENT_STAGE_SUCCEEDED,
							Metadata: &model.NotificationEventStageSucceeded{
								Deployment: s.deployment,
								Stage:      &rbs,
							},
						})
					case model.StageStatus_STAGE_FAILURE:
						s.notifier.Notify(model.NotificationEvent{
							Type: model.NotificationEventType_EVENT_STAGE_FAILED,
							Metadata: &model.NotificationEventStageFailed{
								Deployment: s.deployment,
								Stage:      &rbs,
							},
						})
					case model.StageStatus_STAGE_CANCELLED:
						s.notifier.Notify(model.NotificationEvent{
							Type: model.NotificationEventType_EVENT_STAGE_CANCELLED,
							Metadata: &model.NotificationEventStageCancelled{
								Deployment: s.deployment,
								Stage:      &rbs,
							},
						})
					case model.StageStatus_STAGE_SKIPPED:
						s.notifier.Notify(model.NotificationEvent{
							Type: model.NotificationEventType_EVENT_STAGE_SKIPPED,
							Metadata: &model.NotificationEventStageSkipped{
								Deployment: s.deployment,
								Stage:      &rbs,
							},
						})
					}

					switch result {
					case model.StageStatus_STAGE_SUCCESS:
						span.SetStatus(codes.Ok, statusReason)
					case model.StageStatus_STAGE_FAILURE, model.StageStatus_STAGE_CANCELLED:
						span.SetStatus(codes.Error, statusReason)
					}

					close(doneCh)
				}()

				select {
				case <-ctx.Done():
					handler.Terminate()
					<-doneCh
					return nil

				case <-doneCh:
					break
				}
			}
		}
	}

	if deploymentStatus.IsCompleted() {
		err := s.reportDeploymentCompleted(ctx, deploymentStatus, statusReason, cancelCommander)
		if err == nil && deploymentStatus == model.DeploymentStatus_DEPLOYMENT_SUCCESS {
			s.reportMostRecentlySuccessfulDeployment(ctx)
		}
	}

	if cancelCommand != nil {
		if err := cancelCommand.Report(ctx, model.CommandStatus_COMMAND_SUCCEEDED, nil, nil); err != nil {
			s.logger.Error("failed to report command status", zap.Error(err))
		}
	}

	return nil
}

// executeStage finds the executor for the given stage and execute.
func (s *scheduler) executeStage(sig executor.StopSignal, ps model.PipelineStage, executorFactory func(executor.Input) (executor.Executor, bool)) (finalStatus model.StageStatus) {
	var (
		ctx            = sig.Context()
		originalStatus = ps.Status
		lp             = s.logPersister.StageLogPersister(s.deployment.Id, ps.Id)
	)
	defer func() {
		// When the piped has been terminated (PS kill) while the stage is still running
		// we should not mark the log persister as completed.
		if !finalStatus.IsCompleted() && sig.Terminated() {
			return
		}
		lp.Complete(time.Minute)
	}()

	// Check whether to execute the script rollback stage or not.
	// If the base stage is executed, the script rollback stage will be executed.
	if ps.Name == model.StageScriptRunRollback.String() {
		baseStageID := ps.Metadata["baseStageID"]
		if baseStageID == "" {
			return
		}

		baseStageStatus, ok := s.stageStatuses[baseStageID]
		if !ok {
			return
		}

		if baseStageStatus == model.StageStatus_STAGE_NOT_STARTED_YET || baseStageStatus == model.StageStatus_STAGE_SKIPPED {
			return
		}
	}

	// Update stage status to RUNNING if needed.
	if model.CanUpdateStageStatus(ps.Status, model.StageStatus_STAGE_RUNNING) {
		if err := s.reportStageStatus(ctx, ps.Id, model.StageStatus_STAGE_RUNNING, ps.Requires); err != nil {
			return model.StageStatus_STAGE_FAILURE
		}
		originalStatus = model.StageStatus_STAGE_RUNNING
	}

	// Check the existence of the specified cloud provider.
	if !s.pipedConfig.HasPlatformProvider(s.deployment.PlatformProvider, s.deployment.Kind) {
		lp.Errorf("This piped is not having the specified platform provider in this deployment: %v", s.deployment.PlatformProvider)
		if err := s.reportStageStatus(ctx, ps.Id, model.StageStatus_STAGE_FAILURE, ps.Requires); err != nil {
			s.logger.Error("failed to report stage status", zap.Error(err))
		}
		return model.StageStatus_STAGE_FAILURE
	}

	// Load the stage configuration.
	var stageConfig config.PipelineStage
	var stageConfigFound bool
	if ps.Predefined {
		stageConfig, stageConfigFound = pln.GetPredefinedStage(ps.Id)
	} else {
		stageConfig, stageConfigFound = s.genericApplicationConfig.GetStage(ps.Index)
	}

	if !stageConfigFound {
		lp.Error("Unable to find the stage configuration")
		if err := s.reportStageStatus(ctx, ps.Id, model.StageStatus_STAGE_FAILURE, ps.Requires); err != nil {
			s.logger.Error("failed to report stage status", zap.Error(err))
		}
		return model.StageStatus_STAGE_FAILURE
	}

	app, ok := s.applicationLister.Get(s.deployment.ApplicationId)
	if !ok {
		lp.Errorf("Application %s for this deployment was not found (Maybe it was disabled).", s.deployment.ApplicationId)
		s.reportStageStatus(ctx, ps.Id, model.StageStatus_STAGE_FAILURE, ps.Requires)
		return model.StageStatus_STAGE_FAILURE
	}

	cmdLister := stageCommandLister{
		lister:       s.commandLister,
		deploymentID: s.deployment.Id,
		stageID:      ps.Id,
	}
	alrLister := appLiveResourceLister{
		lister:           s.liveResourceLister,
		platformProvider: app.PlatformProvider,
		appID:            app.Id,
	}
	aStore := appAnalysisResultStore{
		store:         s.analysisResultStore,
		applicationID: app.Id,
	}
	input := executor.Input{
		Stage:                 &ps,
		StageConfig:           stageConfig,
		Deployment:            s.deployment,
		Application:           app,
		PipedConfig:           s.pipedConfig,
		TargetDSP:             s.targetDSP,
		RunningDSP:            s.runningDSP,
		GitClient:             s.gitClient,
		CommandLister:         cmdLister,
		LogPersister:          lp,
		MetadataStore:         s.metadataStore,
		AppManifestsCache:     s.appManifestsCache,
		AppLiveResourceLister: alrLister,
		AnalysisResultStore:   aStore,
		Logger:                s.logger,
		Notifier:              s.notifier,
	}

	// Skip the stage if needed based on the skip config.
	skip, err := s.shouldSkipStage(sig.Context(), input)
	if err != nil {
		lp.Errorf("failed to check whether skipping the stage: %w", err.Error())
		if err := s.reportStageStatus(ctx, ps.Id, model.StageStatus_STAGE_FAILURE, ps.Requires); err != nil {
			s.logger.Error("failed to report stage status", zap.Error(err))
		}
		return model.StageStatus_STAGE_FAILURE
	}
	if skip {
		if err := s.reportStageStatus(ctx, ps.Id, model.StageStatus_STAGE_SKIPPED, ps.Requires); err != nil {
			s.logger.Error("failed to report stage status", zap.Error(err))
			return model.StageStatus_STAGE_FAILURE
		}
		lp.Info("The stage was successfully skipped due to the skip configuration of the stage.")
		return model.StageStatus_STAGE_SKIPPED
	}

	// Find the executor for this stage.
	ex, ok := executorFactory(input)
	if !ok {
		err := fmt.Errorf("no registered executor for stage %s", ps.Name)
		lp.Error(err.Error())
		s.reportStageStatus(ctx, ps.Id, model.StageStatus_STAGE_FAILURE, ps.Requires)
		return model.StageStatus_STAGE_FAILURE
	}

	// Start running executor.
	status := ex.Execute(sig)

	// Commit deployment state status in the following cases:
	// - Apply state successfully.
	// - State was canceled while running (cancel via Controlpane).
	// - Apply state failed but not because of terminating piped process.
	// - State was skipped via Controlpane (currently supports only ANALYSIS stage).
	// - Apply state was exited.
	if status == model.StageStatus_STAGE_SUCCESS ||
		status == model.StageStatus_STAGE_CANCELLED ||
		status == model.StageStatus_STAGE_SKIPPED ||
		status == model.StageStatus_STAGE_EXITED ||
		(status == model.StageStatus_STAGE_FAILURE && !sig.Terminated()) {

		s.reportStageStatus(ctx, ps.Id, status, ps.Requires)
		return status
	}

	// In case piped process got killed (Terminated signal occurred)
	// the original state status will be returned.
	return originalStatus
}

func (s *scheduler) reportStageStatus(ctx context.Context, stageID string, status model.StageStatus, requires []string) error {
	var (
		err error
		now = s.nowFunc()
		req = &pipedservice.ReportStageStatusChangedRequest{
			DeploymentId: s.deployment.Id,
			StageId:      stageID,
			Status:       status,
			Requires:     requires,
			Visible:      true,
			CompletedAt:  now.Unix(),
		}
		retry = pipedservice.NewRetry(10)
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

func (s *scheduler) reportDeploymentStatusChanged(ctx context.Context, status model.DeploymentStatus, desc string) error {
	var (
		err   error
		retry = pipedservice.NewRetry(10)
		req   = &pipedservice.ReportDeploymentStatusChangedRequest{
			DeploymentId:              s.deployment.Id,
			Status:                    status,
			StatusReason:              desc,
			DeploymentChainId:         s.deployment.DeploymentChainId,
			DeploymentChainBlockIndex: s.deployment.DeploymentChainBlockIndex,
		}
	)

	// Update deployment status on remote.
	for retry.WaitNext(ctx) {
		if _, err = s.apiClient.ReportDeploymentStatusChanged(ctx, req); err == nil {
			return nil
		}
		err = fmt.Errorf("failed to report deployment status to control-plane: %v", err)
	}

	return err
}

func (s *scheduler) reportDeploymentCompleted(ctx context.Context, status model.DeploymentStatus, desc, cancelCommander string) error {
	var (
		err error
		now = s.nowFunc()
		req = &pipedservice.ReportDeploymentCompletedRequest{
			DeploymentId:              s.deployment.Id,
			Status:                    status,
			StatusReason:              desc,
			StageStatuses:             s.stageStatuses,
			DeploymentChainId:         s.deployment.DeploymentChainId,
			DeploymentChainBlockIndex: s.deployment.DeploymentChainBlockIndex,
			CompletedAt:               now.Unix(),
		}
		retry = pipedservice.NewRetry(10)
	)

	defer func() {
		switch status {
		case model.DeploymentStatus_DEPLOYMENT_SUCCESS:
			users, groups, err := s.getApplicationNotificationMentions(model.NotificationEventType_EVENT_DEPLOYMENT_SUCCEEDED)
			if err != nil {
				s.logger.Error("failed to get the list of users or groups", zap.Error(err))
			}

			s.notifier.Notify(model.NotificationEvent{
				Type: model.NotificationEventType_EVENT_DEPLOYMENT_SUCCEEDED,
				Metadata: &model.NotificationEventDeploymentSucceeded{
					Deployment:        s.deployment,
					MentionedAccounts: users,
					MentionedGroups:   groups,
				},
			})

		case model.DeploymentStatus_DEPLOYMENT_FAILURE:
			users, groups, err := s.getApplicationNotificationMentions(model.NotificationEventType_EVENT_DEPLOYMENT_FAILED)
			if err != nil {
				s.logger.Error("failed to get the list of users or groups", zap.Error(err))
			}

			s.notifier.Notify(model.NotificationEvent{
				Type: model.NotificationEventType_EVENT_DEPLOYMENT_FAILED,
				Metadata: &model.NotificationEventDeploymentFailed{
					Deployment:        s.deployment,
					Reason:            desc,
					MentionedAccounts: users,
					MentionedGroups:   groups,
				},
			})

		case model.DeploymentStatus_DEPLOYMENT_CANCELLED:
			users, groups, err := s.getApplicationNotificationMentions(model.NotificationEventType_EVENT_DEPLOYMENT_CANCELLED)
			if err != nil {
				s.logger.Error("failed to get the list of users", zap.Error(err))
			}

			s.notifier.Notify(model.NotificationEvent{
				Type: model.NotificationEventType_EVENT_DEPLOYMENT_CANCELLED,
				Metadata: &model.NotificationEventDeploymentCancelled{
					Deployment:        s.deployment,
					Commander:         cancelCommander,
					MentionedAccounts: users,
					MentionedGroups:   groups,
				},
			})
		}
	}()

	// Update deployment status on remote.
	for retry.WaitNext(ctx) {
		if _, err = s.apiClient.ReportDeploymentCompleted(ctx, req); err == nil {
			return nil
		}
		err = fmt.Errorf("failed to report deployment status to control-plane: %w", err)
	}

	return err
}

// getApplicationNotificationMentions returns the list of users groups who should be mentioned in the notification.
func (s *scheduler) getApplicationNotificationMentions(event model.NotificationEventType) ([]string, []string, error) {
	n, ok := s.metadataStore.Shared().Get(model.MetadataKeyDeploymentNotification)
	if !ok {
		return []string{}, []string{}, nil
	}
	var notification config.DeploymentNotification
	if err := json.Unmarshal([]byte(n), &notification); err != nil {
		return nil, nil, fmt.Errorf("could not extract mentions config: %w", err)
	}

	return notification.FindSlackUsers(event), notification.FindSlackGroups(event), nil
}

func (s *scheduler) reportMostRecentlySuccessfulDeployment(ctx context.Context) error {
	var (
		err error
		req = &pipedservice.ReportApplicationMostRecentDeploymentRequest{
			ApplicationId: s.deployment.ApplicationId,
			Status:        model.DeploymentStatus_DEPLOYMENT_SUCCESS,
			Deployment: &model.ApplicationDeploymentReference{
				DeploymentId:   s.deployment.Id,
				Trigger:        s.deployment.Trigger,
				Summary:        s.deployment.Summary,
				Version:        s.deployment.Version,
				Versions:       s.deployment.Versions,
				ConfigFilename: s.deployment.GitPath.GetApplicationConfigFilename(),
				StartedAt:      s.deployment.CreatedAt,
				CompletedAt:    s.deployment.CompletedAt,
			},
		}
		retry = pipedservice.NewRetry(10)
	)

	for retry.WaitNext(ctx) {
		if _, err = s.apiClient.ReportApplicationMostRecentDeployment(ctx, req); err == nil {
			return nil
		}
		err = fmt.Errorf("failed to report most recent successful deployment: %w", err)
	}

	return err
}

type stageCommandLister struct {
	lister       commandLister
	deploymentID string
	stageID      string
}

func (s stageCommandLister) ListCommands() []model.ReportableCommand {
	return s.lister.ListStageCommands(s.deploymentID, s.stageID)
}

type appAnalysisResultStore struct {
	store         analysisResultStore
	applicationID string
}

func (a appAnalysisResultStore) GetLatestAnalysisResult(ctx context.Context) (*model.AnalysisResult, error) {
	return a.store.GetLatestAnalysisResult(ctx, a.applicationID)
}

func (a appAnalysisResultStore) PutLatestAnalysisResult(ctx context.Context, analysisResult *model.AnalysisResult) error {
	return a.store.PutLatestAnalysisResult(ctx, a.applicationID, analysisResult)
}
