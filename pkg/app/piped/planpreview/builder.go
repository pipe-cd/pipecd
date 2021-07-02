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
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/piped/trigger"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/git"
	"github.com/pipe-cd/pipe/pkg/model"
)

type lastTriggeredCommitGetter interface {
	Get(ctx context.Context, applicationID string) (string, error)
}

type Builder interface {
	Build(ctx context.Context, id string, cmd model.Command_BuildPlanPreview) ([]*model.ApplicationPlanPreviewResult, error)
}

type builder struct {
	gitClient         gitClient
	applicationLister applicationLister
	environmentGetter environmentGetter
	commitGetter      lastTriggeredCommitGetter
	config            *config.PipedSpec
	logger            *zap.Logger
}

func newBuilder(gc gitClient, al applicationLister, eg environmentGetter, cg lastTriggeredCommitGetter, cfg *config.PipedSpec, logger *zap.Logger) *builder {
	return &builder{
		gitClient:         gc,
		applicationLister: al,
		environmentGetter: eg,
		commitGetter:      cg,
		config:            cfg,
		logger:            logger.Named("planpreview-builder"),
	}
}

func (b *builder) Build(ctx context.Context, id string, cmd model.Command_BuildPlanPreview) ([]*model.ApplicationPlanPreviewResult, error) {
	b.logger.Info(fmt.Sprintf("start building planpreview result for command %s", id))

	// Find the registered repository in Piped config and validate the command's payload against it.
	repoCfg, ok := b.config.GetRepository(cmd.RepositoryId)
	if !ok {
		return nil, fmt.Errorf("repository %s was not found in Piped config", cmd.RepositoryId)
	}
	if repoCfg.Branch != cmd.BaseBranch {
		return nil, fmt.Errorf("base branch repository %s was not matched, requested %s, expected %s", cmd.RepositoryId, cmd.BaseBranch, repoCfg.Branch)
	}

	// List all applications that belong to this Piped
	// and are placed in the given repository.
	apps := b.listApplications(repoCfg)
	if len(apps) == 0 {
		return nil, nil
	}

	// Clone the source code and checkout to the given branch, commit.
	repo, err := b.gitClient.Clone(ctx, repoCfg.RepoID, repoCfg.Remote, repoCfg.Branch, "")
	if err != nil {
		return nil, fmt.Errorf("failed to clone git repository %s", cmd.RepositoryId)
	}
	defer repo.Clean()

	if err := repo.Checkout(ctx, cmd.HeadCommit); err != nil {
		return nil, fmt.Errorf("failed to checkout the head commit %s: %w", cmd.HeadCommit, err)
	}

	// Compared to the total number of applications,
	// the number of applications that should be triggered will be very smaller
	// therefore we do not explicitly specify the capacity for these slices.
	triggerApps := make([]*model.Application, 0)
	results := make([]*model.ApplicationPlanPreviewResult, 0)

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
			r.Error = fmt.Sprintf("Failed while determining the application should be triggered or not, %v", err)
			results = append(results, r)
			continue
		}

		if shouldTrigger {
			triggerApps = append(triggerApps, app)
		}
	}

	// All triggered applications will be passed to plan.
	for _, app := range triggerApps {
		// We only need the environment name
		// so the returned error can be ignorable.
		var envName string
		if env, err := b.environmentGetter.Get(ctx, app.EnvId); err == nil {
			envName = env.Name
		}

		r := model.MakeApplicationPlanPreviewResult(*app, envName)
		results = append(results, r)

		strategy, changes, err := b.plan(repo, app, cmd)
		if err != nil {
			r.Error = fmt.Sprintf("Failed while planning, %v", err)
			continue
		}

		r.SyncStrategy = strategy
		r.Changes = changes
	}

	return results, nil
}

func (b *builder) plan(repo git.Repo, app *model.Application, cmd model.Command_BuildPlanPreview) (model.SyncStrategy, []byte, error) {
	// TODO: Implement planpreview plan.
	// 1. Start a planner to check what/why strategy will be used
	// 2. Check what resources should be added, deleted and modified
	//    - Terraform app: used terraform plan command
	//    - Kubernetes app: calculate the diff of resources at head commit and mostRecentlySuccessfulCommit

	return model.SyncStrategy_QUICK_SYNC, []byte("NOT IMPLEMENTED"), nil
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
