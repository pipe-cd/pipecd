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
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
	pluginapi "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
)

const cancelRaceAppConfig = `apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: cancel-race
  pipeline:
    stages:
      - name: stage-name
      - name: stage-name
`

type cancelRaceGitRepo struct {
	git.Repo
}

func (r *cancelRaceGitRepo) Checkout(ctx context.Context, commitish string) error {
	return nil
}

type cancelRaceGitClient struct{}

func (c *cancelRaceGitClient) Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error) {
	if err := os.MkdirAll(destination, 0700); err != nil {
		return nil, err
	}
	path := filepath.Join(destination, "app.pipecd.yaml")
	if err := os.WriteFile(path, []byte(cancelRaceAppConfig), 0600); err != nil {
		return nil, err
	}
	return &cancelRaceGitRepo{}, nil
}

type cancelRaceNotifier struct{}

func (n *cancelRaceNotifier) Notify(event model.NotificationEvent) {}

type cancelRaceCommandReporter struct{}

func (r *cancelRaceCommandReporter) ReportStageCommandsHandled(ctx context.Context, deploymentID, stageID string) error {
	return nil
}

type cancelRaceMetadataStore struct{}

func (s *cancelRaceMetadataStore) SharedGet(key string) (string, bool) { return "", false }

func (s *cancelRaceMetadataStore) StageGet(stageID, key string) (string, bool) { return "", false }

type cancelRaceAPIClient struct {
	apiClient
}

func (c *cancelRaceAPIClient) ReportDeploymentStatusChanged(ctx context.Context, req *pipedservice.ReportDeploymentStatusChangedRequest, opts ...grpc.CallOption) (*pipedservice.ReportDeploymentStatusChangedResponse, error) {
	return &pipedservice.ReportDeploymentStatusChangedResponse{}, nil
}

func (c *cancelRaceAPIClient) ReportDeploymentCompleted(ctx context.Context, req *pipedservice.ReportDeploymentCompletedRequest, opts ...grpc.CallOption) (*pipedservice.ReportDeploymentCompletedResponse, error) {
	return &pipedservice.ReportDeploymentCompletedResponse{}, nil
}

func (c *cancelRaceAPIClient) ReportStageStatusChanged(ctx context.Context, req *pipedservice.ReportStageStatusChangedRequest, opts ...grpc.CallOption) (*pipedservice.ReportStageStatusChangedResponse, error) {
	return &pipedservice.ReportStageStatusChangedResponse{}, nil
}

func (c *cancelRaceAPIClient) ReportApplicationMostRecentDeployment(ctx context.Context, req *pipedservice.ReportApplicationMostRecentDeploymentRequest, opts ...grpc.CallOption) (*pipedservice.ReportApplicationMostRecentDeploymentResponse, error) {
	return &pipedservice.ReportApplicationMostRecentDeploymentResponse{}, nil
}

// cancelRacePlugin records the id of every stage it is asked to execute so the
// test can tell whether the scheduler moved on to the next stage.
type cancelRacePlugin struct {
	pluginapi.PluginClient

	// onFirstStage runs just before the first stage reports success, so the
	// cancel lands while that stage is finishing.
	onFirstStage func()

	mu       sync.Mutex
	executed []string
}

func (p *cancelRacePlugin) Close() error { return nil }

func (p *cancelRacePlugin) FetchDefinedStages(ctx context.Context, req *deployment.FetchDefinedStagesRequest, opts ...grpc.CallOption) (*deployment.FetchDefinedStagesResponse, error) {
	return &deployment.FetchDefinedStagesResponse{Stages: []string{"stage-name"}}, nil
}

func (p *cancelRacePlugin) ExecuteStage(ctx context.Context, req *deployment.ExecuteStageRequest, opts ...grpc.CallOption) (*deployment.ExecuteStageResponse, error) {
	p.mu.Lock()
	p.executed = append(p.executed, req.Input.Stage.Id)
	first := len(p.executed) == 1
	p.mu.Unlock()

	if first && p.onFirstStage != nil {
		p.onFirstStage()
	}
	return &deployment.ExecuteStageResponse{Status: model.StageStatus_STAGE_SUCCESS}, nil
}

func (p *cancelRacePlugin) executedStages() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]string(nil), p.executed...)
}

func newCancelRaceDeployment() *model.Deployment {
	return &model.Deployment{
		Id:            "deployment-id",
		ApplicationId: "app-id",
		ProjectId:     "project-id",
		Status:        model.DeploymentStatus_DEPLOYMENT_PLANNED,
		GitPath: &model.ApplicationGitPath{
			Repo: &model.ApplicationGitRepository{
				Id:     "repo-id",
				Remote: "remote",
				Branch: "branch",
			},
			ConfigFilename: "app.pipecd.yaml",
		},
		Trigger: &model.DeploymentTrigger{
			Commit: &model.Commit{Hash: "commit-hash"},
		},
		Stages: []*model.PipelineStage{
			{
				Id:     "stage-1",
				Name:   "stage-name",
				Index:  0,
				Status: model.StageStatus_STAGE_NOT_STARTED_YET,
			},
			{
				Id:     "stage-2",
				Name:   "stage-name",
				Index:  1,
				Status: model.StageStatus_STAGE_NOT_STARTED_YET,
			},
		},
	}
}

// TestSchedulerCancelWhileStageFinishes covers the case where a cancel command
// arrives at the same moment a stage reports success. Both cancelledCh and the
// stage's doneCh are ready when the select runs, so the case that is chosen is
// not fixed. The deployment must end as cancelled either way and the next stage
// must not start. The body is repeated because a single run only exercises one
// of the two orders.
func TestSchedulerCancelWhileStageFinishes(t *testing.T) {
	logger := zaptest.NewLogger(t)

	for i := 0; i < 30; i++ {
		d := newCancelRaceDeployment()
		p := &cancelRacePlugin{}

		pr, err := plugin.NewPluginRegistry(context.TODO(), []plugin.Plugin{
			{Name: "stage-name", Cli: p},
		})
		require.NoError(t, err)

		s := newScheduler(
			d,
			t.TempDir(),
			&cancelRaceAPIClient{},
			&cancelRaceGitClient{},
			pr,
			&cancelRaceNotifier{},
			nil,
			&cancelRaceCommandReporter{},
			&cancelRaceMetadataStore{},
			logger,
			noop.NewTracerProvider(),
		)

		p.onFirstStage = func() {
			s.Cancel(model.ReportableCommand{
				Command: &model.Command{Id: "command-id", Commander: "commander"},
				Report: func(context.Context, model.CommandStatus, map[string]string, []byte) error {
					return nil
				},
			})
		}

		require.NoError(t, s.Run(context.Background()))

		// Run returns as soon as the loop breaks, so give an abandoned stage
		// goroutine time to record itself before checking that none started.
		time.Sleep(20 * time.Millisecond)

		assert.Equal(t, model.DeploymentStatus_DEPLOYMENT_CANCELLED, s.DoneDeploymentStatus())
		assert.Equal(t, []string{"stage-1"}, p.executedStages())
	}
}
