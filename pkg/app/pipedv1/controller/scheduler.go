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

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/controller/controllermetrics"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/deploysource"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/metadatastore"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
)

// scheduler is a dedicated object for a specific deployment of a single application.
type scheduler struct {
	deployment *model.Deployment
	workingDir string

	pluginRegistry plugin.PluginRegistry

	apiClient       apiClient
	gitClient       gitClient
	metadataStore   metadatastore.MetadataStore
	notifier        notifier
	secretDecrypter secretDecrypter

	targetDSP  deploysource.Provider
	runningDSP deploysource.Provider

	// Current status of each stages.
	// We stores their current statuses into this field
	// because the deployment model is readonly to avoid data race.
	// We may need a mutex for this field in the future
	// when the stages can be executed concurrently.
	stageStatuses            map[string]model.StageStatus
	genericApplicationConfig *config.GenericApplicationSpec

	logger *zap.Logger
	tracer trace.Tracer

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
	pluginRegistry plugin.PluginRegistry,
	notifier notifier,
	secretsDecrypter secretDecrypter,
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
		apiClient:            apiClient,
		gitClient:            gitClient,
		pluginRegistry:       pluginRegistry,
		metadataStore:        metadatastore.NewMetadataStore(apiClient, d),
		notifier:             notifier,
		secretDecrypter:      secretsDecrypter,
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
// It determines what stage should be executed next by which plugin.
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

		// Notify the deployment started event
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

	if s.deployment.RunningCommitHash != "" {
		s.runningDSP = deploysource.NewProvider(
			filepath.Join(s.workingDir, "running-deploysource"),
			deploysource.NewGitSourceCloner(s.gitClient, repoCfg, "running", s.deployment.RunningCommitHash),
			s.deployment.GetGitPath(),
			s.secretDecrypter,
		)
	}

	s.targetDSP = deploysource.NewProvider(
		filepath.Join(s.workingDir, "target-deploysource"),
		deploysource.NewGitSourceCloner(s.gitClient, repoCfg, "target", s.deployment.Trigger.Commit.Hash),
		s.deployment.GetGitPath(),
		s.secretDecrypter,
	)

	ds, err := s.targetDSP.Get(ctx, io.Discard)
	if err != nil {
		deploymentStatus = model.DeploymentStatus_DEPLOYMENT_FAILURE
		statusReason = fmt.Sprintf("Failed to get deploy source at target commit (%v)", err)
		s.reportDeploymentCompleted(ctx, deploymentStatus, statusReason, "")
		return err
	}
	cfg, err := config.DecodeYAML[*config.GenericApplicationSpec](ds.ApplicationConfig)
	if err != nil {
		deploymentStatus = model.DeploymentStatus_DEPLOYMENT_FAILURE
		statusReason = fmt.Sprintf("Failed to decode application configuration at target commit (%v)", err)
		s.reportDeploymentCompleted(ctx, deploymentStatus, statusReason, "")
		return err
	}
	s.genericApplicationConfig = cfg.Spec

	ctx, span := s.tracer.Start(
		newContextWithDeploymentSpan(ctx, s.deployment),
		"Deploy",
		trace.WithAttributes(
			attribute.String("application-id", s.deployment.ApplicationId),
			attribute.String("kind", s.deployment.Kind.String()),
			attribute.String("deployment-id", s.deployment.Id),
		))
	defer span.End()

	defer func() {
		switch deploymentStatus {
		case model.DeploymentStatus_DEPLOYMENT_SUCCESS:
			span.SetStatus(codes.Ok, statusReason)
		case model.DeploymentStatus_DEPLOYMENT_FAILURE, model.DeploymentStatus_DEPLOYMENT_CANCELLED:
			span.SetStatus(codes.Error, statusReason)
		}
	}()

	timer := time.NewTimer(s.genericApplicationConfig.Timeout.Duration())
	defer timer.Stop()

	// Iterate all the stages and execute the uncompleted ones.
	for i, ps := range s.deployment.Stages {
		lastStage = s.deployment.Stages[i]

		// Ignore the stage if it is already completed.
		if ps.Status == model.StageStatus_STAGE_SUCCESS {
			continue
		}
		// Ignore the rollback stage, we did it later by another loop.
		if ps.Rollback {
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
			sig, handler = NewStopSignal()
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

			result = s.executeStage(sig, ps)

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
			if sig.Signal() == StopSignalTimeout {
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
					sig, handler = NewStopSignal()
					doneCh       = make(chan struct{})
				)
				go func() {
					rbs := stage
					rbs.Requires = []string{lastStage.Id}

					_, span := s.tracer.Start(ctx, rbs.Name, trace.WithAttributes(
						attribute.String("application-id", s.deployment.ApplicationId),
						attribute.String("kind", s.deployment.Kind.String()),
						attribute.String("deployment-id", s.deployment.Id),
						attribute.String("stage-id", rbs.Id),
					))
					defer span.End()

					result := s.executeStage(sig, rbs)

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

// executeStage finds the plugin for the given stage and execute.
// At the time this executeStage is called, the stage status is before model.StageStatus_STAGE_RUNNING.
// As the first step, it updates the stage status to model.StageStatus_STAGE_RUNNING.
// And that will be treated as the original status of the given stage.
func (s *scheduler) executeStage(sig StopSignal, ps *model.PipelineStage) (finalStatus model.StageStatus) {
	var (
		ctx            = sig.Context()
		originalStatus = ps.Status
	)

	rds, err := s.runningDSP.Get(ctx, io.Discard)
	if err != nil {
		s.logger.Error("failed to get running deployment source", zap.String("stage-name", ps.Name), zap.Error(err))
		return model.StageStatus_STAGE_FAILURE
	}

	tds, err := s.targetDSP.Get(ctx, io.Discard)
	if err != nil {
		s.logger.Error("failed to get target deployment source", zap.String("stage-name", ps.Name), zap.Error(err))
		return model.StageStatus_STAGE_FAILURE
	}

	// Check whether to execute the script rollback stage or not.
	// If the base stage is executed, the script rollback stage will be executed.
	if ps.Rollback {
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

	// Find the executor plugin for this stage.
	plugin, err := s.pluginRegistry.GetPluginClientByStageName(ps.Name)
	if err != nil {
		s.logger.Error("failed to find the plugin for the stage", zap.String("stage-name", ps.Name), zap.Error(err))
		s.reportStageStatus(ctx, ps.Id, model.StageStatus_STAGE_FAILURE, ps.Requires)
		return model.StageStatus_STAGE_FAILURE
	}

	// Load the stage configuration.
	stageConfig, stageConfigFound := s.genericApplicationConfig.GetStageByte(ps.Index)
	if !stageConfigFound {
		s.logger.Error("Unable to find the stage configuration", zap.String("stage-name", ps.Name))
		if err := s.reportStageStatus(ctx, ps.Id, model.StageStatus_STAGE_FAILURE, ps.Requires); err != nil {
			s.logger.Error("failed to report stage status", zap.Error(err))
		}
		return model.StageStatus_STAGE_FAILURE
	}

	// Start running executor.
	res, err := plugin.ExecuteStage(ctx, &deployment.ExecuteStageRequest{
		Input: &deployment.ExecutePluginInput{
			Deployment:              s.deployment,
			Stage:                   ps,
			StageConfig:             stageConfig,
			RunningDeploymentSource: rds.ToPluginDeploySource(),
			TargetDeploymentSource:  tds.ToPluginDeploySource(),
		},
	})
	if err != nil {
		s.logger.Error("failed to execute stage", zap.String("stage-name", ps.Name), zap.Error(err))
		s.reportStageStatus(ctx, ps.Id, model.StageStatus_STAGE_FAILURE, ps.Requires)
		return model.StageStatus_STAGE_FAILURE
	}

	// Determine the final status of the stage.
	status := determineStageStatus(sig.Signal(), originalStatus, res.Status)

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

// determineStageStatus determines the final status of the stage based on the given stop signal.
// Normal is the case when the stop signal is StopSignalNone.
func determineStageStatus(sig StopSignalType, ori, got model.StageStatus) model.StageStatus {
	switch sig {
	case StopSignalNone:
		return got
	case StopSignalTerminate:
		return ori
	case StopSignalCancel:
		return model.StageStatus_STAGE_CANCELLED
	case StopSignalTimeout:
		return model.StageStatus_STAGE_FAILURE
	default:
		return model.StageStatus_STAGE_FAILURE
	}
}

func (s *scheduler) reportStageStatus(ctx context.Context, stageID string, status model.StageStatus, requires []string) error {
	var (
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

	_, err := retry.Do(ctx, func() (interface{}, error) {
		_, err := s.apiClient.ReportStageStatusChanged(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to report stage status to control-plane: %v", err)
		}
		return nil, nil
	})

	return err
}

func (s *scheduler) reportDeploymentStatusChanged(ctx context.Context, status model.DeploymentStatus, desc string) error {
	var (
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
	_, err := retry.Do(ctx, func() (interface{}, error) {
		_, err := s.apiClient.ReportDeploymentStatusChanged(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to report deployment status to control-plane: %v", err)
		}
		return nil, nil
	})

	return err
}

func (s *scheduler) reportDeploymentCompleted(ctx context.Context, status model.DeploymentStatus, desc, cancelCommander string) error {
	var (
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
			users, groups, err := s.getApplicationNotificationMentions(model.NotificationEventType_EVENT_DEPLOYMENT_CANCELLED)
			if err != nil {
				s.logger.Error("failed to get the list of users", zap.Error(err))
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
			users, groups, err := s.getApplicationNotificationMentions(model.NotificationEventType_EVENT_DEPLOYMENT_CANCELLED)
			if err != nil {
				s.logger.Error("failed to get the list of users", zap.Error(err))
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
	_, err := retry.Do(ctx, func() (interface{}, error) {
		_, err := s.apiClient.ReportDeploymentCompleted(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to report deployment status to control-plane: %v", err)
		}
		return nil, nil
	})

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

	_, err := retry.Do(ctx, func() (interface{}, error) {
		_, err := s.apiClient.ReportApplicationMostRecentDeployment(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to report most recent successful deployment: %v", err)
		}
		return nil, nil
	})

	return err
}
