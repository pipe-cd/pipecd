// Copyright 2025 The PipeCD Authors.
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

package planpreview

import (
	"context"
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/controller"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/deploysource"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/trigger"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/backoff"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
	pluginapi "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/common"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	planpreviewapi "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/planpreview"
	"github.com/pipe-cd/pipecd/pkg/regexpool"
)

const (
	workspacePattern    = "plan-preview-builder-*"
	defaultWorkerAppNum = 3
	maxWorkerNum        = 100
)

type lastTriggeredCommitGetter interface {
	Get(ctx context.Context, applicationID string) (string, error)
}

type Builder interface {
	Build(ctx context.Context, id string, cmd model.Command_BuildPlanPreview) ([]*model.ApplicationPlanPreviewResult, error)
}

type builder struct {
	gitClient         gitClient
	apiClient         apiClient
	applicationLister applicationLister
	commitGetter      lastTriggeredCommitGetter
	secretDecrypter   secretDecrypter
	regexPool         *regexpool.Pool
	pipedCfg          *config.PipedSpec
	pluginRegistry    plugin.PluginRegistry
	logger            *zap.Logger

	workingDir string
	repoCfg    config.PipedRepository
}

func newBuilder(
	gc gitClient,
	ac apiClient,
	al applicationLister,
	cg lastTriggeredCommitGetter,
	sd secretDecrypter,
	rp *regexpool.Pool,
	cfg *config.PipedSpec,
	pr plugin.PluginRegistry,
	logger *zap.Logger,
) *builder {

	return &builder{
		gitClient:         gc,
		apiClient:         ac,
		applicationLister: al,
		commitGetter:      cg,
		secretDecrypter:   sd,
		regexPool:         rp,
		pipedCfg:          cfg,
		pluginRegistry:    pr,
		logger:            logger.Named("plan-preview-builder"),
	}
}

func (b *builder) Build(ctx context.Context, id string, cmd model.Command_BuildPlanPreview) (results []*model.ApplicationPlanPreviewResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("an unexpected panic occurred (%v)", r)
			b.logger.Error("unexpected panic", zap.Error(err))
		}
	}()

	return b.build(ctx, id, cmd)
}

func (b *builder) build(ctx context.Context, id string, cmd model.Command_BuildPlanPreview) ([]*model.ApplicationPlanPreviewResult, error) {
	logger := b.logger.With(zap.String("command", id))
	logger.Info(fmt.Sprintf("start building planpreview result for command %s", id))

	// Ensure the existence of the working directory.
	workingDir, err := os.MkdirTemp("", workspacePattern)
	if err != nil {
		return nil, fmt.Errorf("failed to create working directory (%w)", err)
	}
	defer os.RemoveAll(workingDir)
	b.workingDir = workingDir

	// Find the registered repository in Piped config and validate the command's payload against it.
	repoCfg, ok := b.pipedCfg.GetRepository(cmd.RepositoryId)
	if !ok {
		return nil, fmt.Errorf("repository %s was not found in Piped config", cmd.RepositoryId)
	}
	if repoCfg.Branch != cmd.BaseBranch {
		return nil, fmt.Errorf("base branch of repository %s was not matched, requested %s, expected %s", cmd.RepositoryId, cmd.BaseBranch, repoCfg.Branch)
	}
	b.repoCfg = repoCfg

	// List all applications that belong to this Piped
	// and are placed in the given repository.
	apps := b.listApplications(repoCfg)
	if len(apps) == 0 {
		logger.Info(fmt.Sprintf("there is no target application for command %s", id))
		return nil, nil
	}

	// Prepare source code at the head commit.
	// This clones the base branch and merges the head branch into it for correct data.
	// Because new changes might be added into the base branch after the head branch had checked out.
	repo, err := b.cloneHeadCommit(ctx, cmd.HeadBranch, cmd.HeadCommit)
	if err != nil {
		return nil, err
	}

	// We added a merge commit so the commit ID was changed.
	mergedCommit, err := repo.GetLatestCommit(ctx)
	if err != nil {
		return nil, err
	}

	// Find all applications that should be triggered.
	triggerApps, failedResults := b.findTriggerApps(ctx, repo, apps, mergedCommit.Hash)
	results := failedResults

	if len(triggerApps) == 0 {
		return results, nil
	}

	// Plan the trigger applications for more detailed feedback.
	var (
		numApps  = len(triggerApps)
		appCh    = make(chan *model.Application, numApps)
		resultCh = make(chan *model.ApplicationPlanPreviewResult, numApps)
	)
	// Optimize the number of workers.
	numWorkers := numApps / defaultWorkerAppNum
	if numWorkers < 1 {
		numWorkers = numApps
	}
	if numWorkers > maxWorkerNum {
		numWorkers = maxWorkerNum
	}

	// Start some workers to speed up building time.
	logger.Info(fmt.Sprintf("start %d workers for building plan-preview results for %d applications", numWorkers, numApps))
	for w := 0; w < numWorkers; w++ {
		go func(wid int) {
			logger.Info("app worker for plan-preview started", zap.Int("worker", wid))
			for app := range appCh {
				resultCh <- b.buildApp(ctx, wid, id, app, repo, mergedCommit.Hash)
			}
			logger.Info("app worker for plan-preview stopped", zap.Int("worker", wid))
		}(w)
	}

	// Add all applications into the channel for start handling.
	for i := 0; i < numApps; i++ {
		appCh <- triggerApps[i]
	}
	close(appCh)

	// Wait and collect all results.
	for i := 0; i < numApps; i++ {
		r := <-resultCh
		results = append(results, r)
	}

	logger.Info("successfully collected plan-preview results of all applications")
	return results, nil
}

// TODO: add tests
func (b *builder) buildApp(ctx context.Context, worker int, command string, app *model.Application, repo git.Repo, mergedCommit string) (result *model.ApplicationPlanPreviewResult) {
	defer func() {
		// to distinguish that the result is generated by pipedv1
		if len(result.GetPluginNames()) == 0 {
			result.PluginNames = []string{"<unknown>"}
		}
	}()

	logger := b.logger.With(
		zap.Int("worker", worker),
		zap.String("command", command),
		zap.String("app-id", app.Id),
		zap.String("app-name", app.Name),
		zap.String("labels", app.GetLabelsString()),
	)

	logger.Info("will decide sync strategy for an application")

	result = model.MakeApplicationPlanPreviewResult(*app)

	var preCommit string
	// Find the commit of the last successful deployment.
	if deploy, err := b.getMostRecentlySuccessfulDeployment(ctx, app.Id); err == nil {
		preCommit = deploy.Trigger.Commit.Hash
	} else if status.Code(err) != codes.NotFound {
		result.Error = fmt.Sprintf("failed while finding the last successful deployment (%v)", err)
		return
	}

	targetDSP := deploysource.NewProvider(
		b.workingDir,
		deploysource.NewLocalSourceCloner(repo, "target", mergedCommit),
		app.GitPath,
		b.secretDecrypter,
	)
	targetDS, err := targetDSP.Get(ctx, io.Discard)
	if err != nil {
		result.Error = fmt.Sprintf("failed to get the target deploy source, %v", err)
		return
	}
	pluginTargetDS := targetDS.ToPluginDeploySource()
	targetAppCfg, err := config.DecodeYAML[*config.GenericApplicationSpec](pluginTargetDS.GetApplicationConfig())
	if err != nil {
		result.Error = fmt.Sprintf("failed to parse application config, %v", err)
		return
	}

	plugins, err := b.pluginRegistry.GetPluginClientsByAppConfig(targetAppCfg.Spec)
	if err != nil {
		result.Error = fmt.Sprintf("failed to get plugin clients, %v", err)
		return
	}

	strategy, errMsg := b.determineStrategy(ctx, app, pluginTargetDS, targetAppCfg.Spec, plugins, repo, preCommit, targetDSP.Revision())
	if errMsg != "" {
		result.Error = errMsg
		return
	}
	result.SyncStrategy = strategy

	var pluginRunningDS *common.DeploymentSource

	if preCommit != "" {
		runningDSP := deploysource.NewProvider(
			b.workingDir,
			deploysource.NewLocalSourceCloner(repo, "running", preCommit),
			app.GitPath,
			b.secretDecrypter,
		)
		runningDS, err := runningDSP.Get(ctx, io.Discard)
		if err != nil {
			result.Error = fmt.Sprintf("failed to get the running deploy source, %v", err)
			return
		}
		pluginRunningDS = runningDS.ToPluginDeploySource()
	}

	logger.Info("successfully decided sync strategy for the application", zap.String("strategy", result.SyncStrategy.String()))

	errors := ""
	for _, plugin := range plugins {
		result.PluginNames = append(result.PluginNames, plugin.Name())
		res, err := plugin.GetPlanPreview(ctx, &planpreviewapi.GetPlanPreviewRequest{
			ApplicationId:           app.Id,
			ApplicationName:         app.Name,
			PipedId:                 b.pipedCfg.PipedID,
			DeployTargets:           app.GetDeployTargets(),
			TargetDeploymentSource:  pluginTargetDS,
			RunningDeploymentSource: pluginRunningDS,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.Unimplemented {
				logger.Info(fmt.Sprintf("plugin '%s' does not support plan-preview feature", plugin.Name()))
				continue
			}
			errors = fmt.Sprintf("%s\n[%s] %v", errors, plugin.Name(), err)
		}
		for _, r := range res.GetResults() {
			result.PluginPlanResults = append(result.PluginPlanResults, &model.PluginPlanPreviewResult{
				PluginName:   plugin.Name(),
				DeployTarget: r.GetDeployTarget(),
				PlanSummary:  []byte(r.GetSummary()),
				PlanDetails:  r.GetDetails(),
				DiffLanguage: r.GetDiffLanguage(),
			})
			r.NoChange = r.NoChange && result.GetNoChange()
		}
	}

	if len(errors) > 0 {
		result.Error = fmt.Sprintf("failed to get plan preview, %+v", errors)
		return
	}

	return
}

func (b *builder) cloneHeadCommit(ctx context.Context, headBranch, headCommit string) (git.Repo, error) {
	dir, err := os.MkdirTemp(b.workingDir, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory %w", err)
	}

	var (
		remote     = b.repoCfg.Remote
		baseBranch = b.repoCfg.Branch
	)
	repo, err := b.gitClient.Clone(ctx, b.repoCfg.RepoID, remote, baseBranch, dir)
	if err != nil {
		return nil, fmt.Errorf("failed to clone git repository %s at branch %s", b.repoCfg.RepoID, baseBranch)
	}

	mergeCommitMessage := fmt.Sprintf("Plan-preview: merged %s commit from %s branch into %s base branch", headCommit, headBranch, baseBranch)
	if err := repo.MergeRemoteBranch(ctx, headBranch, headCommit, mergeCommitMessage); err != nil {
		return nil, fmt.Errorf("detected conflicts between commit %s at %s branch and the base branch %s (%w)", headCommit, headBranch, baseBranch, err)
	}

	return repo, nil
}

func (b *builder) findTriggerApps(ctx context.Context, repo git.Repo, apps []*model.Application, headCommit string) (triggerApps []*model.Application, failedResults []*model.ApplicationPlanPreviewResult) {
	d := trigger.NewOnCommitDeterminer(repo, headCommit, b.commitGetter, b.logger)
	determine := func(app *model.Application) (bool, error) {
		appCfg, err := config.LoadApplication(repo.GetPath(), app.GitPath.GetApplicationConfigFilePath())
		if err != nil {
			return false, err
		}
		return d.ShouldTrigger(ctx, app, appCfg)
	}

	for _, app := range apps {
		shouldTrigger, err := determine(app)
		if shouldTrigger {
			triggerApps = append(triggerApps, app)
			continue
		}
		if err == nil {
			continue
		}

		r := model.MakeApplicationPlanPreviewResult(*app)
		r.Error = fmt.Sprintf("failed while determining the application should be triggered or not, %v", err)
		failedResults = append(failedResults, r)
	}
	return
}

func (b *builder) determineStrategy(
	ctx context.Context,
	app *model.Application,
	pluginTargetDS *common.DeploymentSource,
	targetAppSpec *config.GenericApplicationSpec,
	plugins []pluginapi.PluginClient,
	repo git.Repo,
	preCommit string,
	revision string,
) (strategy model.SyncStrategy, errorMessage string) {

	var (
		runningDS *deploysource.DeploySource
		err       error
	)

	if preCommit != "" {
		runningDSP := deploysource.NewProvider(
			b.workingDir,
			deploysource.NewLocalSourceCloner(repo, "running", preCommit),
			app.GitPath,
			b.secretDecrypter,
		)
		runningDS, err = runningDSP.Get(ctx, io.Discard)
		if err != nil {
			return model.SyncStrategy_PIPELINE, fmt.Sprintf("failed to get the running deploy source, %v", err)
		}
	}

	// ensure pass nil as running deployment source in case of the first deployment
	// because the running deployment source is not available at this time.
	var pRds *common.DeploymentSource
	if runningDS != nil {
		pRds = runningDS.ToPluginDeploySource()
	}

	planPluginIn := &deployment.PlanPluginInput{
		Deployment: &model.Deployment{
			PipedId:         b.pipedCfg.PipedID,
			ApplicationId:   app.Id,
			ApplicationName: app.Name,
			// Add other fields when needed.
		},
		RunningDeploymentSource: pRds,
		TargetDeploymentSource:  pluginTargetDS,
	}
	trigger := &model.DeploymentTrigger{
		Commit: &model.Commit{
			Branch: b.repoCfg.Branch,
			Hash:   revision,
		},
		Commander: "pipectl",
	}
	strategy, _, err = controller.DetermineStrategy(ctx, targetAppSpec, plugins, planPluginIn, trigger, preCommit, b.logger)
	if err != nil {
		return model.SyncStrategy_PIPELINE, fmt.Sprintf("failed while determining sync strategy, %v", err)
	}
	return strategy, ""
}

func (b *builder) listApplications(repo config.PipedRepository) []*model.Application {
	apps := b.applicationLister.List()
	out := make([]*model.Application, 0, len(apps))

	for _, app := range apps {
		if app.GitPath.Repo.Id != repo.RepoID {
			continue
		}
		if app.GitPath.Repo.Remote != repo.Remote {
			continue
		}
		if app.GitPath.Repo.Branch != repo.Branch {
			continue
		}
		out = append(out, app)
	}

	return out
}

func (b *builder) getMostRecentlySuccessfulDeployment(ctx context.Context, applicationID string) (*model.ApplicationDeploymentReference, error) {
	retry := pipedservice.NewRetry(3)

	deploy, err := retry.Do(ctx, func() (interface{}, error) {
		resp, err := b.apiClient.GetApplicationMostRecentDeployment(ctx, &pipedservice.GetApplicationMostRecentDeploymentRequest{
			ApplicationId: applicationID,
			Status:        model.DeploymentStatus_DEPLOYMENT_SUCCESS,
		})
		if err != nil {
			return nil, backoff.NewError(err, pipedservice.Retriable(err))
		}
		return resp.Deployment, nil
	})
	if err != nil {
		return nil, err
	}

	return deploy.(*model.ApplicationDeploymentReference), nil
}
