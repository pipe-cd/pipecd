// Copyright 2026 The PipeCD Authors.
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
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	deploySourceTargetCommit  = "target-commit"
	deploySourceRunningCommit = "running-commit"
)

// deploySourceGitClient clones a repository with a valid application config,
// but fails to check out the commit given in failCommit.
type deploySourceGitClient struct {
	failCommit string
}

func (c *deploySourceGitClient) Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error) {
	if err := os.MkdirAll(destination, 0700); err != nil {
		return nil, err
	}
	path := filepath.Join(destination, "app.pipecd.yaml")
	if err := os.WriteFile(path, []byte(cancelRaceAppConfig), 0600); err != nil {
		return nil, err
	}
	return &deploySourceGitRepo{failCommit: c.failCommit}, nil
}

type deploySourceGitRepo struct {
	git.Repo
	failCommit string
}

func (r *deploySourceGitRepo) Checkout(ctx context.Context, commitish string) error {
	if commitish == r.failCommit {
		return errors.New("reference not found")
	}
	return nil
}

// deploySourceAPIClient records the deployments reported as completed.
type deploySourceAPIClient struct {
	apiClient
	mu        sync.Mutex
	completed []*pipedservice.ReportDeploymentCompletedRequest
}

func (c *deploySourceAPIClient) ReportDeploymentCompleted(ctx context.Context, req *pipedservice.ReportDeploymentCompletedRequest, opts ...grpc.CallOption) (*pipedservice.ReportDeploymentCompletedResponse, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.completed = append(c.completed, req)
	return &pipedservice.ReportDeploymentCompletedResponse{}, nil
}

func newDeploySourcePlanner(t *testing.T, gc gitClient, ac apiClient, n notifier, lastSuccessfulCommitHash string) *planner {
	d := &model.Deployment{
		Id:            "deployment-id",
		ApplicationId: "app-id",
		Status:        model.DeploymentStatus_DEPLOYMENT_PENDING,
		GitPath: &model.ApplicationGitPath{
			Repo: &model.ApplicationGitRepository{
				Id:     "repo-id",
				Remote: "remote",
				Branch: "branch",
			},
			ConfigFilename: "app.pipecd.yaml",
		},
		Trigger: &model.DeploymentTrigger{
			Commit: &model.Commit{Hash: deploySourceTargetCommit},
		},
	}
	return newPlanner(
		d,
		lastSuccessfulCommitHash,
		"",
		t.TempDir(),
		nil,
		ac,
		gc,
		n,
		nil,
		&cancelRaceMetadataStore{},
		zaptest.NewLogger(t),
		noop.NewTracerProvider(),
	)
}

func TestPlannerReportsFailureWhenDeploySourceCannotBePrepared(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name                     string
		failCommit               string
		lastSuccessfulCommitHash string
		wantReason               string
	}{
		{
			name:       "target deploy source",
			failCommit: deploySourceTargetCommit,
			wantReason: "Unable to prepare the target deploy source",
		},
		{
			name:                     "running deploy source",
			failCommit:               deploySourceRunningCommit,
			lastSuccessfulCommitHash: deploySourceRunningCommit,
			wantReason:               "Unable to prepare the running deploy source",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ac := &deploySourceAPIClient{}
			n := &fakeNotifier{}
			p := newDeploySourcePlanner(t, &deploySourceGitClient{failCommit: tc.failCommit}, ac, n, tc.lastSuccessfulCommitHash)

			require.NoError(t, p.Run(context.Background()))

			assert.True(t, p.IsDone())
			assert.Equal(t, model.DeploymentStatus_DEPLOYMENT_FAILURE, p.DoneDeploymentStatus())

			require.Len(t, ac.completed, 1)
			assert.Equal(t, "deployment-id", ac.completed[0].DeploymentId)
			assert.Equal(t, model.DeploymentStatus_DEPLOYMENT_FAILURE, ac.completed[0].Status)
			assert.Contains(t, ac.completed[0].StatusReason, tc.wantReason)

			require.Len(t, n.events, 1)
			assert.Equal(t, model.NotificationEventType_EVENT_DEPLOYMENT_FAILED, n.events[0].Type)
		})
	}
}

func TestPlannerKeepsDeploymentPendingWhenShuttingDown(t *testing.T) {
	t.Parallel()

	ac := &deploySourceAPIClient{}
	n := &fakeNotifier{}
	p := newDeploySourcePlanner(t, &deploySourceGitClient{failCommit: deploySourceTargetCommit}, ac, n, "")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	require.Error(t, p.Run(ctx))

	assert.True(t, p.IsDone())
	assert.Equal(t, model.DeploymentStatus_DEPLOYMENT_PENDING, p.DoneDeploymentStatus())
	assert.Empty(t, ac.completed)
	assert.Empty(t, n.events)
}
