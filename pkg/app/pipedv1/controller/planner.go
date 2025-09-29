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

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
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
	pluginapi "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/common"
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

	// The apiClient is used to report the deployment status.
	apiClient apiClient

	// The gitClient is used to perform git commands.
	gitClient gitClient

	// The notifier is used for notification features.
	notifier notifier

	// The secretDecrypter is used to decrypt secrets
	// which encrypted using PipeCD built-in secret management.
	secretDecrypter secretDecrypter

	// The metadataStore is used to get the application notification config.
	medatadaStore metadatastore.MetadataStore

	// The pluginRegistry is used to determine which plugins to be used
	pluginRegistry plugin.PluginRegistry

	logger *zap.Logger
	tracer trace.Tracer

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
	pluginRegistry plugin.PluginRegistry,
	apiClient apiClient,
	gitClient gitClient,
	notifier notifier,
	secretDecrypter secretDecrypter,
	metadataStore metadatastore.MetadataStore,
	logger *zap.Logger,
	tracerProvider trace.TracerProvider,
) *planner {

	logger = logger.Named("planner").With(
		zap.String("deployment-id", d.Id),
		zap.String("app-id", d.ApplicationId),
		zap.String("project-id", d.ProjectId),
		zap.String("labels", d.GetLabelsString()),
		zap.String("working-dir", workingDir),
	)

	p := &planner{
		deployment:                   d,
		lastSuccessfulCommitHash:     lastSuccessfulCommitHash,
		lastSuccessfulConfigFilename: lastSuccessfulConfigFilename,
		workingDir:                   workingDir,
		pluginRegistry:               pluginRegistry,
		apiClient:                    apiClient,
		gitClient:                    gitClient,
		notifier:                     notifier,
		secretDecrypter:              secretDecrypter,
		medatadaStore:                metadataStore,
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
			attribute.String("deployment-id", p.deployment.Id),
			attribute.String("labels", p.deployment.GetLabelsString()),
		))
	defer span.End()

	defer func() {
		controllermetrics.UpdateDeploymentStatus(p.deployment, p.doneDeploymentStatus)
	}()

	// Prepare running deploy source and target deploy source.
	var runningDS, targetDS *common.DeploymentSource

	repoCfg := config.PipedRepository{
		RepoID: p.deployment.GitPath.Repo.Id,
		Remote: p.deployment.GitPath.Repo.Remote,
		Branch: p.deployment.GitPath.Repo.Branch,
	}

	targetDSP := deploysource.NewProvider(
		filepath.Join(p.workingDir, "target-deploysource"),
		deploysource.NewGitSourceCloner(p.gitClient, repoCfg, "target", p.deployment.Trigger.Commit.Hash),
		p.deployment.GetGitPath(),
		p.secretDecrypter,
	)
	tds, err := targetDSP.Get(ctx, io.Discard)
	if err != nil {
		p.logger.Error("error while preparing target deploy source data", zap.Error(err))
		return err
	}
	targetDS = tds.ToPluginDeploySource()

	if p.lastSuccessfulCommitHash != "" {
		runningDSP := deploysource.NewProvider(
			filepath.Join(p.workingDir, "running-deploysource"),
			deploysource.NewGitSourceCloner(p.gitClient, repoCfg, "running", p.lastSuccessfulCommitHash),
			p.deployment.GetGitPath(),
			p.secretDecrypter,
		)
		rds, err := runningDSP.Get(ctx, io.Discard)
		if err != nil {
			p.logger.Error("error while preparing running deploy source data", zap.Error(err))
			return err
		}
		runningDS = rds.ToPluginDeploySource()
	}

	out, err := p.buildPlan(ctx, runningDS, targetDS)

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
func (p *planner) buildPlan(ctx context.Context, runningDS, targetDS *common.DeploymentSource) (*plannerOutput, error) {
	out := &plannerOutput{}

	input := &deployment.PlanPluginInput{
		Deployment:              p.deployment,
		RunningDeploymentSource: runningDS,
		TargetDeploymentSource:  targetDS,
	}

	cfg, err := config.DecodeYAML[*config.GenericApplicationSpec](targetDS.GetApplicationConfig())
	if err != nil {
		p.logger.Error("unable to parse application config", zap.Error(err))
		return nil, err
	}
	spec := cfg.Spec

	plugins, err := p.pluginRegistry.GetPluginClientsByAppConfig(spec)
	if err != nil {
		p.logger.Error("unable to get plugin clients by app config", zap.Error(err))
		return nil, err
	}

	// Build deployment target versions.
	for _, plg := range plugins {
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

	strategy, summary, err := DetermineStrategy(ctx, spec, plugins, input, p.deployment.GetTrigger(), p.lastSuccessfulCommitHash, p.logger)
	if err != nil {
		p.logger.Error("unable to determine strategy", zap.Error(err))
		return nil, err
	}

	switch strategy {
	case model.SyncStrategy_QUICK_SYNC:
		if stages, err := p.buildQuickSyncStages(ctx, spec); err == nil {
			out.SyncStrategy = model.SyncStrategy_QUICK_SYNC
			out.Summary = summary
			out.Stages = stages
			return out, nil
		}
	case model.SyncStrategy_PIPELINE:
		if stages, err := p.buildPipelineSyncStages(ctx, spec); err == nil {
			out.SyncStrategy = model.SyncStrategy_PIPELINE
			out.Summary = summary
			out.Stages = stages
			return out, nil
		}
	}

	return nil, fmt.Errorf("unable to plan the deployment")
}

func DetermineStrategy(
	ctx context.Context,
	spec *config.GenericApplicationSpec,
	plugins []pluginapi.PluginClient,
	input *deployment.PlanPluginInput,
	trigger *model.DeploymentTrigger,
	lastSuccessfulCommitHash string,
	logger *zap.Logger,
) (strategy model.SyncStrategy, summary string, err error) {
	// In case the strategy has been decided by trigger.
	// For example: user triggered the deployment via web console.
	switch trigger.SyncStrategy {
	case model.SyncStrategy_QUICK_SYNC:
		return model.SyncStrategy_QUICK_SYNC, trigger.StrategySummary, nil
	case model.SyncStrategy_PIPELINE:
		return model.SyncStrategy_PIPELINE, trigger.StrategySummary, nil
	}

	// When no pipeline was configured, do the quick sync.
	if spec.Pipeline == nil || len(spec.Pipeline.Stages) == 0 {
		return model.SyncStrategy_QUICK_SYNC, "Quick sync due to the pipeline was not configured", nil
	}

	// Force to use pipeline when the `spec.planner.alwaysUsePipeline` was configured.
	if spec.Planner.AlwaysUsePipeline {
		return model.SyncStrategy_PIPELINE, "Sync with the specified pipeline (alwaysUsePipeline was set)", nil
	}

	regexPool := regexpool.DefaultPool()

	// This deployment is triggered by a commit with the intent to perform pipeline.
	// Commit Matcher will be ignored when triggered by a command.
	if pattern := spec.CommitMatcher.Pipeline; pattern != "" && trigger.Commander == "" {
		if pipelineRegex, err := regexPool.Get(pattern); err == nil &&
			pipelineRegex.MatchString(trigger.Commit.Message) {
			return model.SyncStrategy_PIPELINE, fmt.Sprintf("Sync progressively because the commit message was matching %q", pattern), nil
		}
	}

	// This deployment is triggered by a commit with the intent to synchronize.
	// Commit Matcher will be ignored when triggered by a command.
	if pattern := spec.CommitMatcher.QuickSync; pattern != "" && trigger.Commander == "" {
		if syncRegex, err := regexPool.Get(pattern); err == nil &&
			syncRegex.MatchString(trigger.Commit.Message) {
			return model.SyncStrategy_QUICK_SYNC, fmt.Sprintf("Quick sync because the commit message was matching %q", pattern), nil
		}
	}

	// Quick sync if this is the first time to deploy this application or it was unable to retrieve running commit hash.
	if lastSuccessfulCommitHash == "" {
		return model.SyncStrategy_QUICK_SYNC, "Quick sync, it seems this is the first deployment of the application", nil
	}

	// Build plan based on plugins determined strategy
	for _, plg := range plugins {
		res, err := plg.DetermineStrategy(ctx, &deployment.DetermineStrategyRequest{Input: input})
		if err != nil {
			logger.Warn("Unable to determine strategy using current plugin", zap.Error(err))
			continue
		}
		// If the plugin does not support DetermineStrategy(), then ignore.
		if res.Unsupported {
			continue
		}
		strategy = res.SyncStrategy
		summary = res.Summary
		// If one of plugins returns PIPELINE_SYNC, use that as strategy intermediately
		if strategy == model.SyncStrategy_PIPELINE {
			break
		}
	}

	return strategy, summary, nil
}

// buildQuickSyncStages requests all plugins and returns quick sync stage
// from each plugins to build the deployment pipeline.
// NOTE:
//   - For quick sync, we expect all stages given by plugins can be performed in random order.
//   - Rollback stage will always be added as the trail.
func (p *planner) buildQuickSyncStages(ctx context.Context, cfg *config.GenericApplicationSpec) ([]*model.PipelineStage, error) {
	var (
		stages         = []*model.PipelineStage{}
		rollbackStages = []*model.PipelineStage{}
		rollback       = *cfg.Planner.AutoRollback
	)

	plugins, err := p.pluginRegistry.GetPluginClientsByAppConfig(cfg)
	if err != nil {
		p.logger.Error("failed to get plugin clients by app config", zap.Error(err))
		return nil, err
	}
	for _, plg := range plugins {
		res, err := plg.BuildQuickSyncStages(ctx, &deployment.BuildQuickSyncStagesRequest{Rollback: rollback})
		if err != nil {
			p.logger.Error("failed to build quick sync stages for deployment", zap.Error(err))
			return nil, err
		}
		for i := range res.Stages {
			res.Stages[i].Id = uuid.New().String()
			if res.Stages[i].Rollback {
				rollbackStages = append(rollbackStages, res.Stages[i])
			} else {
				stages = append(stages, res.Stages[i])
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
		plg, err := p.pluginRegistry.GetPluginClientByStageName(stageCfg.Name.String())
		if err != nil {
			p.logger.Error("failed to get plugin client by stage name", zap.Error(err))
			return nil, err
		}

		stagesCfgPerPlugin[plg] = append(stagesCfgPerPlugin[plg], &deployment.BuildPipelineSyncStagesRequest_StageConfig{
			Name:   stageCfg.Name.String(),
			Desc:   stageCfg.Desc,
			Index:  int32(i),
			Config: stageCfg.With,
		})
	}

	// Request each plugin to build stages.
	for plg, stageCfgs := range stagesCfgPerPlugin {
		res, err := plg.BuildPipelineSyncStages(ctx, &deployment.BuildPipelineSyncStagesRequest{
			Stages:   stageCfgs,
			Rollback: rollback,
		})
		if err != nil {
			p.logger.Error("failed to build pipeline sync stages for deployment", zap.Error(err))
			return nil, err
		}
		if err := validateStageIndexes(stageCfgs, res.Stages); err != nil {
			p.logger.Error("invalid stage index was returned from a plugin", zap.Error(err))
			return nil, err
		}

		for i := range res.Stages {
			res.Stages[i].Id = uuid.New().String()
			if res.Stages[i].Rollback {
				rollbackStages = append(rollbackStages, res.Stages[i])
			} else {
				stages = append(stages, res.Stages[i])
			}
		}
	}

	// Sort stages by index.
	sort.Sort(model.PipelineStages(stages))
	sort.Sort(model.PipelineStages(rollbackStages))

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

// validateStageIndexes validates the response stage indexes, including rollback stages, for two criteria:
//   - duplication: Indexes of the response stages must not be duplicated within non-rollback stages and rollback stages.
//     A non-rollback stage and a rollback stage can have the same index.
//   - range: Each response stage must have a index defined in the request.
func validateStageIndexes(req []*deployment.BuildPipelineSyncStagesRequest_StageConfig, res []*model.PipelineStage) error {
	// check duplication
	resIndexes := make(map[int32]struct{})
	resRollbackIndexes := make(map[int32]struct{})
	for _, resStage := range res {
		if resStage.Rollback {
			if _, ok := resRollbackIndexes[resStage.Index]; ok {
				return fmt.Errorf("rollback stage index %d from plugin is duplicated", resStage.Index)
			}
			resRollbackIndexes[resStage.Index] = struct{}{}
		} else {
			if _, ok := resIndexes[resStage.Index]; ok {
				return fmt.Errorf("stage index %d from plugin is duplicated", resStage.Index)
			}
			resIndexes[resStage.Index] = struct{}{}
		}
	}

	// check range
	reqIndexes := make(map[int32]struct{})
	for _, reqStage := range req {
		reqIndexes[reqStage.Index] = struct{}{}
	}
	for _, resStage := range res {
		if _, ok := reqIndexes[resStage.Index]; !ok {
			return fmt.Errorf("stage index %d from plugin is not defined in the request", resStage.Index)
		}
	}
	return nil
}

func (p *planner) reportDeploymentPlanned(ctx context.Context, out *plannerOutput) error {
	users, groups, err := p.getApplicationNotificationMentions(model.NotificationEventType_EVENT_DEPLOYMENT_PLANNED)
	if err != nil {
		p.logger.Error("failed to get the list of users or groups", zap.Error(err))
	}

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

	req := &pipedservice.ReportDeploymentPlannedRequest{
		DeploymentId:              p.deployment.Id,
		SyncStrategy:              out.SyncStrategy,
		Summary:                   out.Summary,
		StatusReason:              "The deployment has been planned",
		RunningCommitHash:         p.lastSuccessfulCommitHash,
		RunningConfigFilename:     p.lastSuccessfulConfigFilename,
		Versions:                  out.Versions,
		Stages:                    out.Stages,
		DeploymentChainId:         p.deployment.DeploymentChainId,
		DeploymentChainBlockIndex: p.deployment.DeploymentChainBlockIndex,
	}

	_, err = pipedservice.NewRetry(10).Do(ctx, func() (interface{}, error) {
		_, err := p.apiClient.ReportDeploymentPlanned(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to report deployment status to control-plane: %w", err)
		}
		return nil, nil
	})
	if err != nil {
		p.logger.Error("failed to mark deployment to be planned", zap.Error(err))
	}
	return err
}

func (p *planner) reportDeploymentFailed(ctx context.Context, reason string) error {
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

	req := &pipedservice.ReportDeploymentCompletedRequest{
		DeploymentId:              p.deployment.Id,
		Status:                    model.DeploymentStatus_DEPLOYMENT_FAILURE,
		StatusReason:              reason,
		StageStatuses:             nil,
		DeploymentChainId:         p.deployment.DeploymentChainId,
		DeploymentChainBlockIndex: p.deployment.DeploymentChainBlockIndex,
		CompletedAt:               p.nowFunc().Unix(),
	}

	_, err = pipedservice.NewRetry(10).Do(ctx, func() (interface{}, error) {
		_, err := p.apiClient.ReportDeploymentCompleted(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to report deployment status to control-plane: %w", err)
		}
		return nil, nil
	})

	if err != nil {
		p.logger.Error("failed to mark deployment to be failed", zap.Error(err))
	}
	return err
}

func (p *planner) reportDeploymentCancelled(ctx context.Context, commander, reason string) error {
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

	req := &pipedservice.ReportDeploymentCompletedRequest{
		DeploymentId:              p.deployment.Id,
		Status:                    model.DeploymentStatus_DEPLOYMENT_CANCELLED,
		StatusReason:              reason,
		StageStatuses:             nil,
		DeploymentChainId:         p.deployment.DeploymentChainId,
		DeploymentChainBlockIndex: p.deployment.DeploymentChainBlockIndex,
		CompletedAt:               p.nowFunc().Unix(),
	}

	_, err = pipedservice.NewRetry(10).Do(ctx, func() (interface{}, error) {
		_, err := p.apiClient.ReportDeploymentCompleted(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to report deployment status to control-plane: %w", err)
		}
		return nil, nil
	})

	if err != nil {
		p.logger.Error("failed to mark deployment to be cancelled", zap.Error(err))
	}
	return err
}

// getApplicationNotificationMentions returns the list of users groups who should be mentioned in the notification.
func (p *planner) getApplicationNotificationMentions(event model.NotificationEventType) ([]string, []string, error) {
	n, ok := p.medatadaStore.SharedGet(model.MetadataKeyDeploymentNotification)
	if !ok {
		return []string{}, []string{}, nil
	}

	var notification config.DeploymentNotification
	if err := json.Unmarshal([]byte(n), &notification); err != nil {
		return nil, nil, fmt.Errorf("could not extract mentions config: %w", err)
	}

	return notification.FindSlackUsers(event), notification.FindSlackGroups(event), nil
}
