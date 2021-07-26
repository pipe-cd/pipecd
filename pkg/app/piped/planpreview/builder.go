// Copyright 2021 The PipeCD Authors.
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
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipe/pkg/app/api/service/pipedservice"
	"github.com/pipe-cd/pipe/pkg/app/piped/deploysource"
	"github.com/pipe-cd/pipe/pkg/app/piped/planner"
	"github.com/pipe-cd/pipe/pkg/app/piped/planner/registry"
	"github.com/pipe-cd/pipe/pkg/app/piped/trigger"
	"github.com/pipe-cd/pipe/pkg/backoff"
	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/git"
	"github.com/pipe-cd/pipe/pkg/model"
	"github.com/pipe-cd/pipe/pkg/regexpool"
)

const (
	workspacePattern = "plan-preview-builder-*"
)

var (
	defaultPlannerRegistry = registry.DefaultRegistry()
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
	environmentGetter environmentGetter
	commitGetter      lastTriggeredCommitGetter
	secretDecrypter   secretDecrypter
	appManifestsCache cache.Cache
	regexPool         *regexpool.Pool
	pipedCfg          *config.PipedSpec
	logger            *zap.Logger

	workingDir string
	repoCfg    config.PipedRepository
}

func newBuilder(
	gc gitClient,
	ac apiClient,
	al applicationLister,
	eg environmentGetter,
	cg lastTriggeredCommitGetter,
	sd secretDecrypter,
	amc cache.Cache,
	rp *regexpool.Pool,
	cfg *config.PipedSpec,
	logger *zap.Logger,
) *builder {

	return &builder{
		gitClient:         gc,
		apiClient:         ac,
		applicationLister: al,
		environmentGetter: eg,
		commitGetter:      cg,
		secretDecrypter:   sd,
		appManifestsCache: amc,
		regexPool:         rp,
		pipedCfg:          cfg,
		logger:            logger.Named("plan-preview-builder"),
	}
}

func (b *builder) Build(ctx context.Context, id string, cmd model.Command_BuildPlanPreview) (results []*model.ApplicationPlanPreviewResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("an unexpected panic occurred (%w)", r)
			b.logger.Error("unexpected panic", zap.Error(err))
		}
	}()

	return b.build(ctx, id, cmd)
}

func (b *builder) build(ctx context.Context, id string, cmd model.Command_BuildPlanPreview) ([]*model.ApplicationPlanPreviewResult, error) {
	b.logger.Info(fmt.Sprintf("start building planpreview result for command %s", id))

	// Ensure the existence of the working directory.
	workingDir, err := ioutil.TempDir("", workspacePattern)
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
		b.logger.Info(fmt.Sprintf("there is no target application for command %s", id))
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
	triggerApps, failedResults, err := b.findTriggerApps(ctx, repo, apps, mergedCommit.Hash)
	if err != nil {
		return nil, err
	}
	results := failedResults

	// Plan the trigger applications for more detailed feedback.
	for _, app := range triggerApps {
		b.logger.Info("will decide sync strategy for an application",
			zap.String("id", app.Id),
			zap.String("name", app.Name),
			zap.String("kind", app.Kind.String()),
		)

		// We only need the environment name
		// so the returned error can be ignorable.
		var envName string
		if env, err := b.environmentGetter.Get(ctx, app.EnvId); err == nil {
			envName = env.Name
		}

		r := model.MakeApplicationPlanPreviewResult(*app, envName)
		results = append(results, r)

		var preCommit string
		// Find the commit of the last successful deployment.
		if deploy, err := b.getMostRecentlySuccessfulDeployment(ctx, app.Id); err == nil {
			preCommit = deploy.Trigger.Commit.Hash
		} else if status.Code(err) != codes.NotFound {
			r.Error = fmt.Sprintf("failed while finding the last successful deployment (%w)", err)
			continue
		}

		targetDSP := deploysource.NewProvider(
			b.workingDir,
			deploysource.NewLocalSourceCloner(repo, "target", mergedCommit.Hash),
			*app.GitPath,
			b.secretDecrypter,
		)

		strategy, err := b.plan(ctx, app, targetDSP, preCommit)
		if err != nil {
			r.Error = fmt.Sprintf("failed while planning, %v", err)
			continue
		}
		r.SyncStrategy = strategy

		b.logger.Info("successfully decided sync strategy for a application",
			zap.String("id", app.Id),
			zap.String("name", app.Name),
			zap.String("strategy", strategy.String()),
			zap.String("kind", app.Kind.String()),
		)

		var buf bytes.Buffer
		var summary string

		switch app.Kind {
		case model.ApplicationKind_KUBERNETES:
			summary, err = b.kubernetesDiff(ctx, app, targetDSP, preCommit, &buf)
		case model.ApplicationKind_TERRAFORM:
			summary, err = b.terraformDiff(ctx, app, targetDSP, &buf)
		default:
			// TODO: Calculating planpreview's diff for other application kinds.
			summary = fmt.Sprintf("%s application is not implemented yet (coming soon)", app.Kind.String())
		}

		r.PlanSummary = []byte(summary)
		r.PlanDetails = buf.Bytes()
		if err != nil {
			r.Error = fmt.Sprintf("failed while calculating diff, %v", err)
			continue
		}
	}

	return results, nil
}

func (b *builder) cloneHeadCommit(ctx context.Context, headBranch, headCommit string) (git.Repo, error) {
	dir, err := ioutil.TempDir(b.workingDir, "")
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

func (b *builder) findTriggerApps(ctx context.Context, repo git.Repo, apps []*model.Application, headCommit string) (triggerApps []*model.Application, failedResults []*model.ApplicationPlanPreviewResult, err error) {
	d := trigger.NewDeterminer(repo, headCommit, b.commitGetter, b.logger)
	for _, app := range apps {
		shouldTrigger, err := d.ShouldTrigger(ctx, app)
		if err != nil {
			// We only need the environment name
			// so the returned error can be ignorable.
			var envName string
			if env, err := b.environmentGetter.Get(ctx, app.EnvId); err == nil {
				envName = env.Name
			}

			r := model.MakeApplicationPlanPreviewResult(*app, envName)
			r.Error = fmt.Sprintf("failed while determining the application should be triggered or not, %v", err)
			failedResults = append(failedResults, r)
			continue
		}

		if shouldTrigger {
			triggerApps = append(triggerApps, app)
		}
	}
	return
}

func (b *builder) plan(ctx context.Context, app *model.Application, targetDSP deploysource.Provider, lastSuccessfulCommit string) (strategy model.SyncStrategy, err error) {
	p, ok := defaultPlannerRegistry.Planner(app.Kind)
	if !ok {
		err = fmt.Errorf("application kind %s is not supported yet", app.Kind.String())
		return
	}

	in := planner.Input{
		ApplicationID:   app.Id,
		ApplicationName: app.Name,
		GitPath:         *app.GitPath,
		Trigger: model.DeploymentTrigger{
			Commit: &model.Commit{
				Branch: b.repoCfg.Branch,
				Hash:   targetDSP.Revision(),
			},
			Commander: "pipectl",
		},
		TargetDSP:                      targetDSP,
		MostRecentSuccessfulCommitHash: lastSuccessfulCommit,
		PipedConfig:                    b.pipedCfg,
		AppManifestsCache:              b.appManifestsCache,
		RegexPool:                      b.regexPool,
		Logger:                         b.logger,
	}

	if lastSuccessfulCommit != "" {
		in.RunningDSP = deploysource.NewProvider(
			b.workingDir,
			deploysource.NewGitSourceCloner(b.gitClient, b.repoCfg, "running", lastSuccessfulCommit),
			*app.GitPath,
			b.secretDecrypter,
		)
	}

	out, err := p.Plan(ctx, in)
	if err != nil {
		return
	}

	strategy = out.SyncStrategy
	return
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
