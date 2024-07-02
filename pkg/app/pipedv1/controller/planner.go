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
	"time"

	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/controller/controllermetrics"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/metadatastore"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/platform"
)

type planner struct {
	// Readonly deployment model.
	deployment                   *model.Deployment
	lastSuccessfulCommitHash     string
	lastSuccessfulConfigFilename string
	workingDir                   string
	pipedConfig                  []byte

	// The pluginClient is used to call pluggin that actually
	// performs planning deployment.
	pluginClient platform.PlatformPluginClient

	// The apiClient is used to report the deployment status.
	apiClient apiClient

	// The notifier and metadataStore are used for
	// notification features.
	notifier      notifier
	metadataStore metadatastore.MetadataStore

	// TODO: Find a way to show log from pluggin's planner
	logger *zap.Logger

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
	pluginClient platform.PlatformPluginClient,
	apiClient apiClient,
	notifier notifier,
	pipedConfig []byte,
	logger *zap.Logger,
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
		pluginClient:                 pluginClient,
		apiClient:                    apiClient,
		metadataStore:                metadatastore.NewMetadataStore(apiClient, d),
		notifier:                     notifier,
		pipedConfig:                  pipedConfig,
		doneDeploymentStatus:         d.Status,
		cancelledCh:                  make(chan *model.ReportableCommand, 1),
		nowFunc:                      time.Now,
		logger:                       logger,
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

// What planner does:
// - Wait until there is no PLANNED or RUNNING deployment
// - Pick the oldest PENDING deployment to plan its pipeline
// - <*> Perform planning a deployment by calling the pluggin's planner
// - Update the deployment status to PLANNED or not based on the result
func (p *planner) Run(ctx context.Context) error {
	p.logger.Info("start running planner")

	defer func() {
		p.doneTimestamp = p.nowFunc()
		p.done.Store(true)
	}()

	defer func() {
		controllermetrics.UpdateDeploymentStatus(p.deployment, p.doneDeploymentStatus)
	}()

	return nil

	// in := &platform.BuildPlanRequest{
	// 	Deployment:                   p.deployment,
	// 	WorkingDir:                   p.workingDir,
	// 	LastSuccessfulCommitHash:     p.lastSuccessfulCommitHash,
	// 	LastSuccessfulConfigFileName: p.lastSuccessfulConfigFilename,
	// 	PipedConfig:                  p.pipedConfig,
	// }

	// out, err := p.pluginClient.BuildPlan(ctx, in)

	// // If the deployment was already cancelled, we ignore the plan result.
	// select {
	// case cmd := <-p.cancelledCh:
	// 	if cmd != nil {
	// 		p.doneDeploymentStatus = model.DeploymentStatus_DEPLOYMENT_CANCELLED
	// 		desc := fmt.Sprintf("Deployment was cancelled by %s while planning", cmd.Commander)
	// 		p.reportDeploymentCancelled(ctx, cmd.Commander, desc)
	// 		return cmd.Report(ctx, model.CommandStatus_COMMAND_SUCCEEDED, nil, nil)
	// 	}
	// default:
	// }

	// if err != nil {
	// 	p.doneDeploymentStatus = model.DeploymentStatus_DEPLOYMENT_FAILURE
	// 	return p.reportDeploymentFailed(ctx, fmt.Sprintf("Unable to plan the deployment (%v)", err))
	// }

	// p.doneDeploymentStatus = model.DeploymentStatus_DEPLOYMENT_PLANNED
	// return p.reportDeploymentPlanned(ctx, out.Plan)
}

// func (p *planner) reportDeploymentPlanned(ctx context.Context, out *platform.DeploymentPlan) error {
// 	var (
// 		err   error
// 		retry = pipedservice.NewRetry(10)
// 		req   = &pipedservice.ReportDeploymentPlannedRequest{
// 			DeploymentId:              p.deployment.Id,
// 			Summary:                   out.Summary,
// 			StatusReason:              "The deployment has been planned",
// 			RunningCommitHash:         p.lastSuccessfulCommitHash,
// 			RunningConfigFilename:     p.lastSuccessfulConfigFilename,
// 			Versions:                  out.Versions,
// 			Stages:                    out.Stages,
// 			DeploymentChainId:         p.deployment.DeploymentChainId,
// 			DeploymentChainBlockIndex: p.deployment.DeploymentChainBlockIndex,
// 		}
// 	)

// 	accounts, err := p.getMentionedAccounts(model.NotificationEventType_EVENT_DEPLOYMENT_PLANNED)
// 	if err != nil {
// 		p.logger.Error("failed to get the list of accounts", zap.Error(err))
// 	}

// 	defer func() {
// 		p.notifier.Notify(model.NotificationEvent{
// 			Type: model.NotificationEventType_EVENT_DEPLOYMENT_PLANNED,
// 			Metadata: &model.NotificationEventDeploymentPlanned{
// 				Deployment:        p.deployment,
// 				Summary:           out.Summary,
// 				MentionedAccounts: accounts,
// 			},
// 		})
// 	}()

// 	for retry.WaitNext(ctx) {
// 		if _, err = p.apiClient.ReportDeploymentPlanned(ctx, req); err == nil {
// 			return nil
// 		}
// 		err = fmt.Errorf("failed to report deployment status to control-plane: %v", err)
// 	}

// 	if err != nil {
// 		p.logger.Error("failed to mark deployment to be planned", zap.Error(err))
// 	}
// 	return err
// }

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

	accounts, err := p.getMentionedAccounts(model.NotificationEventType_EVENT_DEPLOYMENT_FAILED)
	if err != nil {
		p.logger.Error("failed to get the list of accounts", zap.Error(err))
	}

	defer func() {
		p.notifier.Notify(model.NotificationEvent{
			Type: model.NotificationEventType_EVENT_DEPLOYMENT_FAILED,
			Metadata: &model.NotificationEventDeploymentFailed{
				Deployment:        p.deployment,
				Reason:            reason,
				MentionedAccounts: accounts,
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

	accounts, err := p.getMentionedAccounts(model.NotificationEventType_EVENT_DEPLOYMENT_CANCELLED)
	if err != nil {
		p.logger.Error("failed to get the list of accounts", zap.Error(err))
	}

	defer func() {
		p.notifier.Notify(model.NotificationEvent{
			Type: model.NotificationEventType_EVENT_DEPLOYMENT_CANCELLED,
			Metadata: &model.NotificationEventDeploymentCancelled{
				Deployment:        p.deployment,
				Commander:         commander,
				MentionedAccounts: accounts,
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

func (p *planner) getMentionedAccounts(event model.NotificationEventType) ([]string, error) {
	n, ok := p.metadataStore.Shared().Get(model.MetadataKeyDeploymentNotification)
	if !ok {
		return []string{}, nil
	}

	var notification config.DeploymentNotification
	if err := json.Unmarshal([]byte(n), &notification); err != nil {
		return nil, fmt.Errorf("could not extract mentions config: %w", err)
	}

	return notification.FindSlackAccountsAndGroups(event), nil
}
