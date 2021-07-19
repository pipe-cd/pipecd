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
	"github.com/pipe-cd/pipe/pkg/cache"
	"github.com/pipe-cd/pipe/pkg/config"
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

	// Find all applications that should be triggered.
	triggerApps, failedResults, err := b.findTriggerApps(ctx, apps, cmd)
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

		strategy, err := b.plan(ctx, app, cmd, preCommit)
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
			summary, err = b.kubernetesDiff(ctx, app, cmd, preCommit, &buf)
		case model.ApplicationKind_TERRAFORM:
			summary, err = b.terraformDiff(ctx, app, cmd, &buf)
		default:
			// TODO: Calculating planpreview's diff for other application kinds.
			err = fmt.Errorf("%s application is not implemented yet (coming soon)", app.Kind.String())
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

func (b *builder) findTriggerApps(ctx context.Context, apps []*model.Application, cmd model.Command_BuildPlanPreview) (triggerApps []*model.Application, failedResults []*model.ApplicationPlanPreviewResult, err error) {
	// Clone the source code and checkout to the given branch, commit.
	dir, err := ioutil.TempDir(b.workingDir, "")
	if err != nil {
		err = fmt.Errorf("failed to create temporary directory %w", err)
		return
	}
	repo, err := b.gitClient.Clone(ctx, b.repoCfg.RepoID, b.repoCfg.Remote, cmd.HeadBranch, dir)
	if err != nil {
		err = fmt.Errorf("failed to clone git repository %s", cmd.RepositoryId)
		return
	}
	defer repo.Clean()

	err = repo.Checkout(ctx, cmd.HeadCommit)
	if err != nil {
		err = fmt.Errorf("failed to checkout the head commit %s: %w", cmd.HeadCommit, err)
		return
	}

	d := trigger.NewDeterminer(repo, cmd.HeadCommit, b.commitGetter, b.logger)
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

func (b *builder) plan(ctx context.Context, app *model.Application, cmd model.Command_BuildPlanPreview, lastSuccessfulCommit string) (strategy model.SyncStrategy, err error) {
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
				Branch: cmd.HeadBranch,
				Hash:   cmd.HeadCommit,
			},
			Commander: "pipectl",
		},
		MostRecentSuccessfulCommitHash: lastSuccessfulCommit,
		PipedConfig:                    b.pipedCfg,
		AppManifestsCache:              b.appManifestsCache,
		RegexPool:                      b.regexPool,
		Logger:                         b.logger,
	}

	repoCfg := b.repoCfg
	repoCfg.Branch = cmd.HeadBranch

	in.TargetDSP = deploysource.NewProvider(
		b.workingDir,
		repoCfg,
		"target",
		cmd.HeadCommit,
		b.gitClient,
		app.GitPath,
		b.secretDecrypter,
	)

	if lastSuccessfulCommit != "" {
		in.RunningDSP = deploysource.NewProvider(
			b.workingDir,
			repoCfg,
			"running",
			lastSuccessfulCommit,
			b.gitClient,
			app.GitPath,
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
	var (
		err   error
		resp  *pipedservice.GetApplicationMostRecentDeploymentResponse
		retry = pipedservice.NewRetry(3)
		req   = &pipedservice.GetApplicationMostRecentDeploymentRequest{
			ApplicationId: applicationID,
			Status:        model.DeploymentStatus_DEPLOYMENT_SUCCESS,
		}
	)

	for retry.WaitNext(ctx) {
		if resp, err = b.apiClient.GetApplicationMostRecentDeployment(ctx, req); err == nil {
			return resp.Deployment, nil
		}
		if !pipedservice.Retriable(err) {
			return nil, err
		}
	}
	return nil, err
}
