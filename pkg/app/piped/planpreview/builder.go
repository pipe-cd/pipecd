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

	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

type Builder interface {
	Build(ctx context.Context, id string, cmd model.Command_BuildPlanPreview) ([]*model.ApplicationPlanPreviewResult, error)
}

type builder struct {
	gitClient         gitClient
	applicationLister applicationLister
	config            *config.PipedSpec
	logger            *zap.Logger
}

func newBuilder(gc gitClient, al applicationLister, cfg *config.PipedSpec, logger *zap.Logger) *builder {
	return &builder{
		gitClient:         gc,
		applicationLister: al,
		config:            cfg,
		logger:            logger.Named("planpreview-builder"),
	}
}

func (b *builder) Build(ctx context.Context, id string, cmd model.Command_BuildPlanPreview) ([]*model.ApplicationPlanPreviewResult, error) {
	repoCfg, ok := b.config.GetRepository(cmd.RepositoryId)
	if !ok {
		return nil, fmt.Errorf("repository %s was not found in Piped config", cmd.RepositoryId)
	}
	if repoCfg.Branch != cmd.BaseBranch {
		return nil, fmt.Errorf("base branch repository %s was not correct, requested %s, expected %s", cmd.RepositoryId, cmd.BaseBranch, repoCfg.Branch)
	}

	apps := b.listApplications(repoCfg)
	if len(apps) == 0 {
		return nil, nil
	}

	repo, err := b.gitClient.Clone(ctx, repoCfg.RepoID, repoCfg.Remote, repoCfg.Branch, "")
	if err != nil {
		return nil, fmt.Errorf("failed to clone git repository %s", cmd.RepositoryId)
	}
	defer repo.Clean()

	// TODO: Implement planpreview builder.
	// 1. Fetch the source code at the head commit.
	// 2. Determine the list of applications that will be triggered.
	//    - Based on the changed files between 2 commits: head commit and mostRecentlyTriggeredCommit
	// 3. For each application:
	//    3.1. Start a builder to check what/why strategy will be used
	//    3.2. Check what resources should be added, deleted and modified
	//         - Terraform app: used terraform plan command
	//         - Kubernetes app: calculate the diff of resources at head commit and mostRecentlySuccessfulCommit
	return nil, fmt.Errorf("Not Implemented")
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
