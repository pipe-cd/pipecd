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
	"sort"
	"time"

	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/controller/controllermetrics"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/deploysource"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/metadatastore"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
	pluginapi "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/regexpool"
)

const (
	versionUnknown = "unknown"
)

type plannerOutput struct {
	Versions     []*model.ArtifactVersion
	SyncStrategy model.SyncStrategy
	Summary      string
	Stages       []*model.PipelineStage
}

type planner struct {
	// Readonly deployment model.
	deployment                   *model.Deployment
	lastSuccessfulCommitHash     string
	lastSuccessfulConfigFilename string
	workingDir                   string
	pipedConfig                  []byte

	// The plugin clients are used to call plugin that actually
	// performs planning deployment.
	plugins []pluginapi.PluginClient
	// The map used to know which plugin is incharged for a given stage
	// of the current deployment.
	stageBasedPluginsMap map[string]pluginapi.PluginClient

	// The apiClient is used to report the deployment status.
	apiClient apiClient

	// The gitClient is used to perform git commands.
	gitClient gitClient

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
	pluginClient pluginapi.PluginClient,
	apiClient apiClient,
	gitClient gitClient,
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

	// TODO: Fix this. Passed by args
	tmp := make(map[string]pluginapi.PluginClient)
	tmp["K8S_SYNC"] = pluginClient

	plugins := make([]pluginapi.PluginClient, 0, len(tmp))
	for _, v := range tmp {
		plugins = append(plugins, v)
	}

	p := &planner{
		deployment:                   d,
		lastSuccessfulCommitHash:     lastSuccessfulCommitHash,
		lastSuccessfulConfigFilename: lastSuccessfulConfigFilename,
		workingDir:                   workingDir,
		stageBasedPluginsMap:         tmp,
		plugins:                      plugins,
		apiClient:                    apiClient,
		gitClient:                    gitClient,
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

func (p *planner) Run(ctx context.Context) error {
	p.logger.Info("start running planner")

	defer func() {
		p.doneTimestamp = p.nowFunc()
		p.done.Store(true)
	}()

	defer func() {
		controllermetrics.UpdateDeploymentStatus(p.deployment, p.doneDeploymentStatus)
	}()

	repoCfg := config.PipedRepository{
		RepoID: p.deployment.GitPath.Repo.Id,
		Remote: p.deployment.GitPath.Repo.Remote,
		Branch: p.deployment.GitPath.Repo.Branch,
	}

	// Prepare target deploy source.
	targetDSP := deploysource.NewProvider(
		filepath.Join(p.workingDir, "deploysource"),
		deploysource.NewGitSourceCloner(p.gitClient, repoCfg, "target", p.deployment.Trigger.Commit.Hash),
		*p.deployment.GitPath,
		nil, // TODO: Revise this secret decryter, is this need?
	)

	targetDS, err := targetDSP.Get(ctx, io.Discard)
	if err != nil {
		return fmt.Errorf("error while preparing deploy source data (%v)", err)
	}

	// TODO: Pass running DS as well if need?
	out, err := p.buildPlan(ctx, targetDS)

	// If the deployment was already cancelled, we ignore the plan result.
	select {
	case cmd := <-p.cancelledCh:
		if cmd != nil {
			p.doneDeploymentStatus = model.DeploymentStatus_DEPLOYMENT_CANCELLED
			desc := fmt.Sprintf("Deployment was cancelled by %s while planning", cmd.Commander)
			p.reportDeploymentCancelled(ctx, cmd.Commander, desc)
			return cmd.Report(ctx, model.CommandStatus_COMMAND_SUCCEEDED, nil, nil)
		}
	default:
	}

	if err != nil {
		p.doneDeploymentStatus = model.DeploymentStatus_DEPLOYMENT_FAILURE
		return p.reportDeploymentFailed(ctx, fmt.Sprintf("Unable to plan the deployment (%v)", err))
	}

	p.doneDeploymentStatus = model.DeploymentStatus_DEPLOYMENT_PLANNED
	return p.reportDeploymentPlanned(ctx, out)
}

// buildPlan builds the deployment plan.
// The strategy determination logic is based on the following order:
//   - Direct trigger via web console
//   - Force quick sync if there is no pipeline specified
//   - Force pipeline if the `spec.planner.alwaysUsePipeline` was configured
//   - CommitMatcher ensure pipeline/quick sync based on the commit message
//   - Force quick sync if there is no previous deployment (aka. this is the first deploy)
//   - Based on PlannerService.DetermineStrategy returned by plugins
func (p *planner) buildPlan(ctx context.Context, targetDS *deploysource.DeploySource) (*plannerOutput, error) {
	out := &plannerOutput{}

	input := &deployment.PlanPluginInput{
		Deployment: p.deployment,
		// TODO: Add more planner input fields.
		// NOTE: As discussed we pass targetDS & runningDS here.
	}

	// Build deployment target versions.
	for _, plg := range p.plugins {
		vRes, err := plg.DetermineVersions(ctx, &deployment.DetermineVersionsRequest{Input: input})
		if err != nil {
			p.logger.Warn("unable to determine versions", zap.Error(err))
			continue
		}
		out.Versions = append(out.Versions, vRes.Versions...)
	}
	if len(out.Versions) == 0 {
		out.Versions = []*model.ArtifactVersion{
			{
				Kind:    model.ArtifactVersion_UNKNOWN,
				Version: versionUnknown,
			},
		}
	}

	cfg := targetDS.GenericApplicationConfig

	// In case the strategy has been decided by trigger.
	// For example: user triggered the deployment via web console.
	switch p.deployment.Trigger.SyncStrategy {
	case model.SyncStrategy_QUICK_SYNC:
		if stages, err := p.buildQuickSyncStages(ctx, cfg); err == nil {
			out.SyncStrategy = model.SyncStrategy_QUICK_SYNC
			out.Summary = p.deployment.Trigger.StrategySummary
			out.Stages = stages
			return out, nil
		}
	case model.SyncStrategy_PIPELINE:
		if stages, err := p.buildPipelineSyncStages(ctx, cfg); err == nil {
			out.SyncStrategy = model.SyncStrategy_PIPELINE
			out.Summary = p.deployment.Trigger.StrategySummary
			out.Stages = stages
			return out, nil
		}
	}

	// When no pipeline was configured, do the quick sync.
	if cfg.Pipeline == nil || len(cfg.Pipeline.Stages) == 0 {
		if stages, err := p.buildQuickSyncStages(ctx, cfg); err == nil {
			out.SyncStrategy = model.SyncStrategy_QUICK_SYNC
			out.Summary = "Quick sync due to the pipeline was not configured"
			out.Stages = stages
			return out, nil
		}
	}

	// Force to use pipeline when the `spec.planner.alwaysUsePipeline` was configured.
	if cfg.Planner.AlwaysUsePipeline {
		if stages, err := p.buildPipelineSyncStages(ctx, cfg); err == nil {
			out.SyncStrategy = model.SyncStrategy_PIPELINE
			out.Summary = "Sync with the specified pipeline (alwaysUsePipeline was set)"
			out.Stages = stages
			return out, nil
		}
	}

	regexPool := regexpool.DefaultPool()

	// This deployment is triggered by a commit with the intent to perform pipeline.
	// Commit Matcher will be ignored when triggered by a command.
	if pattern := cfg.CommitMatcher.Pipeline; pattern != "" && p.deployment.Trigger.Commander == "" {
		if pipelineRegex, err := regexPool.Get(pattern); err == nil &&
			pipelineRegex.MatchString(p.deployment.Trigger.Commit.Message) {
			if stages, err := p.buildPipelineSyncStages(ctx, cfg); err == nil {
				out.SyncStrategy = model.SyncStrategy_PIPELINE
				out.Summary = fmt.Sprintf("Sync progressively because the commit message was matching %q", pattern)
				out.Stages = stages
				return out, nil
			}
		}
	}

	// This deployment is triggered by a commit with the intent to synchronize.
	// Commit Matcher will be ignored when triggered by a command.
	if pattern := cfg.CommitMatcher.QuickSync; pattern != "" && p.deployment.Trigger.Commander == "" {
		if syncRegex, err := regexPool.Get(pattern); err == nil &&
			syncRegex.MatchString(p.deployment.Trigger.Commit.Message) {
			if stages, err := p.buildQuickSyncStages(ctx, cfg); err == nil {
				out.SyncStrategy = model.SyncStrategy_QUICK_SYNC
				out.Summary = fmt.Sprintf("Quick sync because the commit message was matching %q", pattern)
				out.Stages = stages
				return out, nil
			}
		}
	}

	// Quick sync if this is the first time to deploy this application or it was unable to retrieve running commit hash.
	if p.lastSuccessfulCommitHash == "" {
		if stages, err := p.buildQuickSyncStages(ctx, cfg); err == nil {
			out.SyncStrategy = model.SyncStrategy_QUICK_SYNC
			out.Summary = "Quick sync, it seems this is the first deployment of the application"
			out.Stages = stages
			return out, nil
		}
	}

	var (
		strategy model.SyncStrategy
		summary  string
	)
	// Build plan based on plugins determined strategy
	for _, plg := range p.plugins {
		res, err := plg.DetermineStrategy(ctx, &deployment.DetermineStrategyRequest{Input: input})
		if err != nil {
			p.logger.Warn("Unable to determine strategy using current plugin", zap.Error(err))
			continue
		}
		strategy = res.SyncStrategy
		summary = res.Summary
		// If one of plugins returns PIPELINE_SYNC, use that as strategy intermediately
		if strategy == model.SyncStrategy_PIPELINE {
			break
		}
	}

	switch strategy {
	case model.SyncStrategy_QUICK_SYNC:
		if stages, err := p.buildQuickSyncStages(ctx, cfg); err == nil {
			out.SyncStrategy = model.SyncStrategy_QUICK_SYNC
			out.Summary = summary
			out.Stages = stages
			return out, nil
		}
	case model.SyncStrategy_PIPELINE:
		if stages, err := p.buildPipelineSyncStages(ctx, cfg); err == nil {
			out.SyncStrategy = model.SyncStrategy_PIPELINE
			out.Summary = summary
			out.Stages = stages
			return out, nil
		}
	}

	return nil, fmt.Errorf("unable to plan the deployment")
}

// buildQuickSyncStages requests all plugins and returns quick sync stage
// from each plugins to build the deployment pipeline.
// NOTE:
//   - For quick sync, we expect all stages given by plugins can be performed
//     at once regradless its order (aka. no `Stage.Requires` specified)
//   - Rollback stage will always be added as the trail.
func (p *planner) buildQuickSyncStages(ctx context.Context, cfg *config.GenericApplicationSpec) ([]*model.PipelineStage, error) {
	var (
		stages         = []*model.PipelineStage{}
		rollbackStages = []*model.PipelineStage{}
		rollback       = *cfg.Planner.AutoRollback
	)
	for _, plg := range p.plugins {
		res, err := plg.BuildQuickSyncStages(ctx, &deployment.BuildQuickSyncStagesRequest{Rollback: rollback})
		if err != nil {
			return nil, fmt.Errorf("failed to build quick sync stage deployment (%w)", err)
		}
		for i := range res.Stages {
			// TODO: Consider add Stage.Rollback to specify a stage is a rollback stage or forward stage instead.
			if res.Stages[i].Visible {
				stages = append(stages, res.Stages[i])
			} else {
				rollbackStages = append(rollbackStages, res.Stages[i])
			}
		}
	}
	stages = append(stages, rollbackStages...)
	if len(stages) == 0 {
		return nil, fmt.Errorf("unable to build quick sync stages for deployment")
	}
	return stages, nil
}

// buildPipelineSyncStages requests all plugins and returns built stages which be used
// to combine the deployment pipeline based on the application configuration `spec.pipeline.stages`.
// NOTE:
//   - The order of stages is determined by the application configuration `spec.pipeline.stages`.
//   - The `Stage.Requires` field is used to specify by the order of stages.
//   - Rollback stage will always be added as the trail.
func (p *planner) buildPipelineSyncStages(ctx context.Context, cfg *config.GenericApplicationSpec) ([]*model.PipelineStage, error) {
	var (
		rollback  = *cfg.Planner.AutoRollback
		stagesCfg = cfg.Pipeline.Stages

		stages         = make([]*model.PipelineStage, 0, len(stagesCfg))
		rollbackStages = []*model.PipelineStage{}

		stagesCfgPerPlugin = make(map[pluginapi.PluginClient][]*deployment.BuildPipelineSyncStagesRequest_StageConfig)
	)

	// Build stages config for each plugin.
	for i := range stagesCfg {
		stageCfg := stagesCfg[i]
		plg, ok := p.stageBasedPluginsMap[stageCfg.Name.String()]
		if !ok {
			return nil, fmt.Errorf("unable to find plugin for stage %q", stageCfg.Name.String())
		}

		stagesCfgPerPlugin[plg] = append(stagesCfgPerPlugin[plg], &deployment.BuildPipelineSyncStagesRequest_StageConfig{
			Id:      stageCfg.ID,
			Name:    stageCfg.Name.String(),
			Desc:    stageCfg.Desc,
			Timeout: stageCfg.Timeout.Duration().String(),
			Index:   int32(i),
			Config:  stageCfg.With,
		})
	}

	// Request each plugin to build stages.
	for plg, stageCfgs := range stagesCfgPerPlugin {
		res, err := plg.BuildPipelineSyncStages(ctx, &deployment.BuildPipelineSyncStagesRequest{
			Stages:   stageCfgs,
			Rollback: rollback,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to build pipeline sync stages for deployment (%w)", err)
		}
		for i := range res.Stages {
			if res.Stages[i].Visible {
				stages = append(stages, res.Stages[i])
			} else {
				rollbackStages = append(rollbackStages, res.Stages[i])
			}
		}
	}

	// Sort stages by index.
	sort.Sort(model.PipelineStages(stages))

	// Build requires for each stage.
	preStageID := ""
	for _, s := range stages {
		if preStageID != "" {
			s.Requires = []string{preStageID}
		}
		preStageID = s.Id
	}

	// Append all rollback stages as the trail.
	stages = append(stages, rollbackStages...)

	if len(stages) == 0 {
		return nil, fmt.Errorf("unable to build pipeline sync stages for deployment")
	}
	return stages, nil
}

func (p *planner) reportDeploymentPlanned(ctx context.Context, out *plannerOutput) error {
	var (
		err   error
		retry = pipedservice.NewRetry(10)
		req   = &pipedservice.ReportDeploymentPlannedRequest{
			DeploymentId:              p.deployment.Id,
			Summary:                   out.Summary,
			StatusReason:              "The deployment has been planned",
			RunningCommitHash:         p.lastSuccessfulCommitHash,
			RunningConfigFilename:     p.lastSuccessfulConfigFilename,
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
