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
	"path/filepath"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/piped/controller/controllermetrics"
	"github.com/pipe-cd/pipecd/pkg/app/piped/deploysource"
	"github.com/pipe-cd/pipecd/pkg/app/piped/metadatastore"
	pln "github.com/pipe-cd/pipecd/pkg/app/piped/planner"
	"github.com/pipe-cd/pipecd/pkg/app/piped/planner/registry"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/regexpool"
)

// What planner does:
// - Wait until there is no PLANNED or RUNNING deployment
// - Pick the oldest PENDING deployment to plan its pipeline
// - Compare with the last successful commit
// - Decide the pipeline should be executed (scale, progressive, rollback)
// - Update the pipeline stages and change the deployment status to PLANNED
type planner struct {
	// Readonly deployment model.
	deployment                   *model.Deployment
	lastSuccessfulCommitHash     string
	lastSuccessfulConfigFilename string
	workingDir                   string
	apiClient                    apiClient
	gitClient                    gitClient
	metadataStore                metadatastore.MetadataStore
	notifier                     notifier
	secretDecrypter              secretDecrypter
	plannerRegistry              registry.Registry
	pipedConfig                  *config.PipedSpec
	appManifestsCache            cache.Cache
	logger                       *zap.Logger
	tracer                       trace.Tracer

	done                 atomic.Bool
	doneTimestamp        time.Time
	doneDeploymentStatus model.DeploymentStatus
	cancelled            bool
	cancelledCh          chan *model.ReportableCommand

	nowFunc func() time.Time
}

func newPlanner(
	d *model.Deployment,
	lastSuccessfulCommitHash string,
	lastSuccessfulConfigFilename string,
	workingDir string,
	apiClient apiClient,
	gitClient gitClient,
	notifier notifier,
	sd secretDecrypter,
	pipedConfig *config.PipedSpec,
	appManifestsCache cache.Cache,
	logger *zap.Logger,
	tracerProvider trace.TracerProvider,
) *planner {

	logger = logger.Named("planner").With(
		zap.String("deployment-id", d.Id),
		zap.String("app-id", d.ApplicationId),
		zap.String("project-id", d.ProjectId),
		zap.String("app-kind", d.Kind.String()),
		zap.String("working-dir", workingDir),
	)

	p := &planner{
		deployment:                   d,
		lastSuccessfulCommitHash:     lastSuccessfulCommitHash,
		lastSuccessfulConfigFilename: lastSuccessfulConfigFilename,
		workingDir:                   workingDir,
		apiClient:                    apiClient,
		gitClient:                    gitClient,
		metadataStore:                metadatastore.NewMetadataStore(apiClient, d),
		notifier:                     notifier,
		secretDecrypter:              sd,
		pipedConfig:                  pipedConfig,
		plannerRegistry:              registry.DefaultRegistry(),
		appManifestsCache:            appManifestsCache,
		doneDeploymentStatus:         d.Status,
		cancelledCh:                  make(chan *model.ReportableCommand, 1),
		nowFunc:                      time.Now,
		logger:                       logger,
		tracer:                       tracerProvider.Tracer("controller/planner"),
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

// DoneDeploymentStatus returns the deployment status when planner has done.
// This can be used only after IsDone() returns true.
func (p *planner) DoneDeploymentStatus() model.DeploymentStatus {
	return p.doneDeploymentStatus
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
	p.logger.Info("start running planner")

	defer func() {
		p.doneTimestamp = p.nowFunc()
		p.done.Store(true)
	}()

	ctx, span := p.tracer.Start(
		newContextWithDeploymentSpan(ctx, p.deployment),
		"Plan",
		trace.WithAttributes(
			attribute.String("application-id", p.deployment.ApplicationId),
			attribute.String("kind", p.deployment.Kind.String()),
			attribute.String("deployment-id", p.deployment.Id),
		))
	defer span.End()

	repoCfg := config.PipedRepository{
		RepoID: p.deployment.GitPath.Repo.Id,
		Remote: p.deployment.GitPath.Repo.Remote,
		Branch: p.deployment.GitPath.Repo.Branch,
	}

	in := pln.Input{
		ApplicationID:                  p.deployment.ApplicationId,
		ApplicationName:                p.deployment.ApplicationName,
		GitPath:                        *p.deployment.GitPath,
		Trigger:                        *p.deployment.Trigger,
		MostRecentSuccessfulCommitHash: p.lastSuccessfulCommitHash,
		PipedConfig:                    p.pipedConfig,
		AppManifestsCache:              p.appManifestsCache,
		RegexPool:                      regexpool.DefaultPool(),
		GitClient:                      p.gitClient,
		Logger:                         p.logger,
	}

	in.TargetDSP = deploysource.NewProvider(
		filepath.Join(p.workingDir, "target-deploysource"),
		deploysource.NewGitSourceCloner(p.gitClient, repoCfg, "target", p.deployment.Trigger.Commit.Hash),
		*p.deployment.GitPath,
		p.secretDecrypter,
	)

	if p.lastSuccessfulCommitHash != "" {
		gp := *p.deployment.GitPath
		gp.ConfigFilename = p.lastSuccessfulConfigFilename

		in.RunningDSP = deploysource.NewProvider(
			filepath.Join(p.workingDir, "running-deploysource"),
			deploysource.NewGitSourceCloner(p.gitClient, repoCfg, "running", p.lastSuccessfulCommitHash),
			gp,
			p.secretDecrypter,
		)
	}

	defer func() {
		controllermetrics.UpdateDeploymentStatus(p.deployment, p.doneDeploymentStatus)
	}()

	planner, ok := p.plannerRegistry.Planner(p.deployment.Kind)
	if !ok {
		p.doneDeploymentStatus = model.DeploymentStatus_DEPLOYMENT_FAILURE
		p.reportDeploymentFailed(ctx, "Unable to find the planner for this application kind")
		return fmt.Errorf("unable to find the planner for application %v", p.deployment.Kind)
	}

	out, err := planner.Plan(ctx, in)

	// If the deployment was already cancelled, we ignore the plan result.
	select {
	case cmd := <-p.cancelledCh:
		if cmd != nil {
			p.doneDeploymentStatus = model.DeploymentStatus_DEPLOYMENT_CANCELLED
			desc := fmt.Sprintf("Deployment was cancelled by %s while planning", cmd.Commander)
			span.SetStatus(codes.Error, desc)
			p.reportDeploymentCancelled(ctx, cmd.Commander, desc)
			return cmd.Report(ctx, model.CommandStatus_COMMAND_SUCCEEDED, nil, nil)
		}
	default:
	}

	if err != nil {
		p.doneDeploymentStatus = model.DeploymentStatus_DEPLOYMENT_FAILURE
		span.SetStatus(codes.Error, err.Error())
		return p.reportDeploymentFailed(ctx, fmt.Sprintf("Unable to plan the deployment (%v)", err))
	}

	span.SetStatus(codes.Ok, "The deployment has been planned")
	p.doneDeploymentStatus = model.DeploymentStatus_DEPLOYMENT_PLANNED
	return p.reportDeploymentPlanned(ctx, out)
}

func (p *planner) reportDeploymentPlanned(ctx context.Context, out pln.Output) error {
	var (
		err   error
		retry = pipedservice.NewRetry(10)
		req   = &pipedservice.ReportDeploymentPlannedRequest{
			DeploymentId:              p.deployment.Id,
			Summary:                   out.Summary,
			StatusReason:              "The deployment has been planned",
			RunningCommitHash:         p.lastSuccessfulCommitHash,
			RunningConfigFilename:     p.lastSuccessfulConfigFilename,
			Version:                   out.Version,
			Versions:                  out.Versions,
			Stages:                    out.Stages,
			DeploymentChainId:         p.deployment.DeploymentChainId,
			DeploymentChainBlockIndex: p.deployment.DeploymentChainBlockIndex,
		}
	)

	users, groups, err := p.getApplicationNotificationMentions(model.NotificationEventType_EVENT_DEPLOYMENT_PLANNED)

	defer func() {
		p.notifier.Notify(model.NotificationEvent{
			Type: model.NotificationEventType_EVENT_DEPLOYMENT_PLANNED,
			Metadata: &model.NotificationEventDeploymentPlanned{
				Deployment:        p.deployment,
				Summary:           out.Summary,
				MentionedAccounts: users,
				MentionedGroups:   groups,
			},
		})
	}()

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

func (p *planner) reportDeploymentFailed(ctx context.Context, reason string) error {
	var (
		err error
		now = p.nowFunc()
		req = &pipedservice.ReportDeploymentCompletedRequest{
			DeploymentId:              p.deployment.Id,
			Status:                    model.DeploymentStatus_DEPLOYMENT_FAILURE,
			StatusReason:              reason,
			StageStatuses:             nil,
			DeploymentChainId:         p.deployment.DeploymentChainId,
			DeploymentChainBlockIndex: p.deployment.DeploymentChainBlockIndex,
			CompletedAt:               now.Unix(),
		}
		retry = pipedservice.NewRetry(10)
	)

	users, groups, err := p.getApplicationNotificationMentions(model.NotificationEventType_EVENT_DEPLOYMENT_FAILED)
	if err != nil {
		p.logger.Error("failed to get the list of users or groups", zap.Error(err))
	}

	defer func() {
		p.notifier.Notify(model.NotificationEvent{
			Type: model.NotificationEventType_EVENT_DEPLOYMENT_FAILED,
			Metadata: &model.NotificationEventDeploymentFailed{
				Deployment:        p.deployment,
				Reason:            reason,
				MentionedAccounts: users,
				MentionedGroups:   groups,
			},
		})
	}()

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

func (p *planner) reportDeploymentCancelled(ctx context.Context, commander, reason string) error {
	var (
		err error
		now = p.nowFunc()
		req = &pipedservice.ReportDeploymentCompletedRequest{
			DeploymentId:              p.deployment.Id,
			Status:                    model.DeploymentStatus_DEPLOYMENT_CANCELLED,
			StatusReason:              reason,
			StageStatuses:             nil,
			DeploymentChainId:         p.deployment.DeploymentChainId,
			DeploymentChainBlockIndex: p.deployment.DeploymentChainBlockIndex,
			CompletedAt:               now.Unix(),
		}
		retry = pipedservice.NewRetry(10)
	)

	users, groups, err := p.getApplicationNotificationMentions(model.NotificationEventType_EVENT_DEPLOYMENT_CANCELLED)
	if err != nil {
		p.logger.Error("failed to get the list of users or groups", zap.Error(err))
	}

	defer func() {
		p.notifier.Notify(model.NotificationEvent{
			Type: model.NotificationEventType_EVENT_DEPLOYMENT_CANCELLED,
			Metadata: &model.NotificationEventDeploymentCancelled{
				Deployment:        p.deployment,
				Commander:         commander,
				MentionedAccounts: users,
				MentionedGroups:   groups,
			},
		})
	}()

	for retry.WaitNext(ctx) {
		if _, err = p.apiClient.ReportDeploymentCompleted(ctx, req); err == nil {
			return nil
		}
		err = fmt.Errorf("failed to report deployment status to control-plane: %v", err)
	}

	if err != nil {
		p.logger.Error("failed to mark deployment to be cancelled", zap.Error(err))
	}
	return err
}

// getApplicationNotificationMentions returns the list of users groups who should be mentioned in the notification.
func (p *planner) getApplicationNotificationMentions(event model.NotificationEventType) ([]string, []string, error) {
	n, ok := p.metadataStore.Shared().Get(model.MetadataKeyDeploymentNotification)
	if !ok {
		return []string{}, []string{}, nil
	}

	var notification config.DeploymentNotification
	if err := json.Unmarshal([]byte(n), &notification); err != nil {
		return nil, nil, fmt.Errorf("could not extract mentions config: %w", err)
	}

	return notification.FindSlackUsers(event), notification.FindSlackGroups(event), nil
}
