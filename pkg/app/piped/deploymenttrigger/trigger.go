// Copyright 2020 The PipeCD Authors.
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

// Package deploymenttrigger provides a piped component
// that detects a list of application should be synced
// and then trigger their deployments by calling to API to create a new Deployment model.
// Until V1, we detect based on the new merged commit and its changes.
// But in the next versions, we also want to enable the ability to detect
// based on the diff between the repo state (desired state) and cluster state (actual state).
package deploymenttrigger

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/kapetaniosci/pipe/pkg/app/api/service/pipedservice"
	"github.com/kapetaniosci/pipe/pkg/config"
	"github.com/kapetaniosci/pipe/pkg/model"
)

type apiClient interface {
	ListApplications(ctx context.Context, in *pipedservice.ListApplicationsRequest, opts ...grpc.CallOption) (*pipedservice.ListApplicationsResponse, error)
	CreateDeployment(ctx context.Context, in *pipedservice.CreateDeploymentRequest, opts ...grpc.CallOption) (*pipedservice.CreateDeploymentResponse, error)
}

type gitClient interface {
	GetLatestRemoteHashForBranch(ctx context.Context, remote, branch string) (string, error)
}

type applicationStore interface {
	ListApplications() []*model.Application
	GetApplication(id string) (*model.Application, bool)
}

type commandStore interface {
	ListApplicationCommands() []*model.Command
	ReportCommandHandled(ctx context.Context, c *model.Command, status model.CommandStatus, metadata map[string]string) error
}

type DeploymentTrigger struct {
	apiClient        apiClient
	gitClient        gitClient
	applicationStore applicationStore
	commandStore     commandStore
	config           *config.PipedSpec
	triggeredCommits map[string]string
	mu               sync.Mutex
	gracePeriod      time.Duration
	logger           *zap.Logger
}

// NewTrigger creates a new instance for DeploymentTrigger.
// What does this need to do its task?
// - A way to get commit/source-code of a specific repository
// - A way to get the current state of application
func NewTrigger(apiClient apiClient, gitClient gitClient, appStore applicationStore, cmdStore commandStore, cfg *config.PipedSpec, gracePeriod time.Duration, logger *zap.Logger) *DeploymentTrigger {
	return &DeploymentTrigger{
		apiClient:        apiClient,
		gitClient:        gitClient,
		applicationStore: appStore,
		commandStore:     cmdStore,
		config:           cfg,
		gracePeriod:      gracePeriod,
		logger:           logger.Named("deployment-trigger"),
	}
}

// Run starts running DeploymentTrigger until the specified context
// has done. This also waits for its cleaning up before returning.
func (t *DeploymentTrigger) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(t.config.SyncInterval))
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				t.check(ctx)
			}
		}
	}()

	return nil
}

func (t *DeploymentTrigger) check(ctx context.Context) error {
	// List all applications that should be handled by this piped
	// and then group them by repository.
	applications, err := t.listApplications(ctx)
	if err != nil {
		return err
	}

	repos := t.config.GetRepositoryMap()
	if len(repos) == 0 {
		t.logger.Info("no repositories were configured for this piped")
		return nil
	}

	for repoID, apps := range applications {
		repo, ok := repos[repoID]
		if !ok {
			t.logger.Warn("detected some applications are binding with an non existent repository",
				zap.String("repo-id", repoID),
				zap.String("application-id", apps[0].Id),
			)
			continue
		}

		// Get the head commit of the repository.
		headCommitSHA, err := t.gitClient.GetLatestRemoteHashForBranch(ctx, repo.Remote, repo.Branch)
		if err != nil {
			continue
		}

		for _, app := range apps {
			// Get the most recently applied commit of this application.
			// If it is not in the memory cache, we have to call the API to list the deployments
			// and use the commit sha of the most recent one.
			triggeredCommitSHA, ok := t.triggeredCommits[app.Id]
			if !ok {
				t.triggeredCommits[app.Id] = "retrieved-one"
			}

			// Check whether the most recently applied one is the head commit or not.
			// If not, nothing to do for this time.
			if triggeredCommitSHA == headCommitSHA {
				continue
			}

			// List the changed files between those two commits
			// Determine whether this application was touch by those changed files.

			// Send a request to API to create a new deployment.
			if err := t.triggerDeployment(ctx, app, headCommitSHA); err != nil {
				continue
			}
			t.triggeredCommits[app.Id] = headCommitSHA
		}
	}
	return nil
}

func (t *DeploymentTrigger) listApplications(ctx context.Context) (map[string][]*model.Application, error) {
	return nil, nil
}

func (t *DeploymentTrigger) triggerDeployment(ctx context.Context, app *model.Application, commit string) error {
	// Detect the update type (just scale or need rollout with pipeline) by checking the change
	deployment := &model.Deployment{
		Id:            uuid.New().String(),
		ApplicationId: "fake-application-id",
	}
	_, err := t.apiClient.CreateDeployment(ctx, &pipedservice.CreateDeploymentRequest{
		Deployment: deployment,
	})
	if err != nil {
		t.logger.Error("failed to create deployment", zap.Error(err))
	}
	return nil
}
