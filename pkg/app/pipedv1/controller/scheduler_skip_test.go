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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
	pluginapi "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
)

const (
	skipOnPrefixAppConfig = `apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: skip-on
  pipeline:
    stages:
      - name: STAGE_A
        skipOn:
          commitMessagePrefixes: ["[skip-a]"]
      - name: STAGE_B
`
	skipOnPathsAppConfig = `apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: skip-on
  pipeline:
    stages:
      - name: STAGE_A
        skipOn:
          paths: ["docs/**"]
      - name: STAGE_B
`
)

type skipOnGitRepo struct {
	git.Repo
	commitMessage string
	changedFiles  []string
}

func (r *skipOnGitRepo) Checkout(ctx context.Context, commitish string) error { return nil }

func (r *skipOnGitRepo) Clean() error { return nil }

func (r *skipOnGitRepo) GetCommitForRev(ctx context.Context, rev string) (git.Commit, error) {
	return git.Commit{Message: r.commitMessage}, nil
}

func (r *skipOnGitRepo) ChangedFiles(ctx context.Context, from, to string) ([]string, error) {
	return r.changedFiles, nil
}

// skipOnGitClient writes the application config into the deploy source
// directories and returns the same repo for the clone of the skip check,
// which has no destination.
type skipOnGitClient struct {
	appConfig string
	repo      *skipOnGitRepo
}

func (c *skipOnGitClient) Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error) {
	if destination != "" {
		if err := os.MkdirAll(destination, 0700); err != nil {
			return nil, err
		}
		if err := os.WriteFile(filepath.Join(destination, "app.pipecd.yaml"), []byte(c.appConfig), 0600); err != nil {
			return nil, err
		}
	}
	return c.repo, nil
}

// skipOnPlugin records the id of every stage it executes and fails the stages in failStages.
type skipOnPlugin struct {
	pluginapi.PluginClient
	failStages map[string]bool

	mu       sync.Mutex
	executed []string
}

func (p *skipOnPlugin) Close() error { return nil }

func (p *skipOnPlugin) FetchDefinedStages(ctx context.Context, req *deployment.FetchDefinedStagesRequest, opts ...grpc.CallOption) (*deployment.FetchDefinedStagesResponse, error) {
	return &deployment.FetchDefinedStagesResponse{Stages: []string{"STAGE_A", "STAGE_B", "SYNC", "ROLLBACK"}}, nil
}

func (p *skipOnPlugin) ExecuteStage(ctx context.Context, req *deployment.ExecuteStageRequest, opts ...grpc.CallOption) (*deployment.ExecuteStageResponse, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.executed = append(p.executed, req.Input.Stage.Id)
	if p.failStages[req.Input.Stage.Id] {
		return &deployment.ExecuteStageResponse{Status: model.StageStatus_STAGE_FAILURE}, nil
	}
	return &deployment.ExecuteStageResponse{Status: model.StageStatus_STAGE_SUCCESS}, nil
}

func TestSchedulerSkipOn(t *testing.T) {
	logger := zaptest.NewLogger(t)

	// The rollback stage of a plugin reuses the minimum index of its stages,
	// and quick sync stages always have index 0.
	pipelineStages := func() []*model.PipelineStage {
		return []*model.PipelineStage{
			{Id: "stage-a", Name: "STAGE_A", Index: 0},
			{Id: "stage-b", Name: "STAGE_B", Index: 1},
			{Id: "rollback", Name: "ROLLBACK", Index: 0, Rollback: true},
		}
	}
	quickSyncStages := func() []*model.PipelineStage {
		return []*model.PipelineStage{
			{Id: "sync", Name: "SYNC", Index: 0},
			{Id: "rollback", Name: "ROLLBACK", Index: 0, Rollback: true},
		}
	}

	testcases := []struct {
		name              string
		appConfig         string
		commitMessage     string
		syncStrategy      model.SyncStrategy
		runningCommitHash string
		stages            []*model.PipelineStage
		failStages        map[string]bool
		expectedStatus    model.DeploymentStatus
		expectedExecuted  []string
		expectedStages    map[string]model.StageStatus
	}{
		{
			name:              "pipeline stage is skipped by its skip config",
			appConfig:         skipOnPrefixAppConfig,
			commitMessage:     "[skip-a] change",
			syncStrategy:      model.SyncStrategy_PIPELINE,
			runningCommitHash: "running-commit-hash",
			stages:            pipelineStages(),
			expectedStatus:    model.DeploymentStatus_DEPLOYMENT_SUCCESS,
			expectedExecuted:  []string{"stage-b"},
			expectedStages: map[string]model.StageStatus{
				"stage-a":  model.StageStatus_STAGE_SKIPPED,
				"stage-b":  model.StageStatus_STAGE_SUCCESS,
				"rollback": model.StageStatus_STAGE_NOT_STARTED_YET,
			},
		},
		{
			name:              "rollback stage is not skipped by the skip config of the pipeline stage with the same index",
			appConfig:         skipOnPrefixAppConfig,
			commitMessage:     "[skip-a] change",
			syncStrategy:      model.SyncStrategy_PIPELINE,
			runningCommitHash: "running-commit-hash",
			stages:            pipelineStages(),
			failStages:        map[string]bool{"stage-b": true},
			expectedStatus:    model.DeploymentStatus_DEPLOYMENT_FAILURE,
			expectedExecuted:  []string{"stage-b", "rollback"},
			expectedStages: map[string]model.StageStatus{
				"stage-a":  model.StageStatus_STAGE_SKIPPED,
				"stage-b":  model.StageStatus_STAGE_FAILURE,
				"rollback": model.StageStatus_STAGE_SUCCESS,
			},
		},
		{
			name:              "quick sync stage is not skipped by the skip config of the first pipeline stage",
			appConfig:         skipOnPrefixAppConfig,
			commitMessage:     "[skip-a] change",
			syncStrategy:      model.SyncStrategy_QUICK_SYNC,
			runningCommitHash: "running-commit-hash",
			stages:            quickSyncStages(),
			expectedStatus:    model.DeploymentStatus_DEPLOYMENT_SUCCESS,
			expectedExecuted:  []string{"sync"},
			expectedStages: map[string]model.StageStatus{
				"sync":     model.StageStatus_STAGE_SUCCESS,
				"rollback": model.StageStatus_STAGE_NOT_STARTED_YET,
			},
		},
		{
			name:             "first deployment with quick sync",
			appConfig:        skipOnPrefixAppConfig,
			commitMessage:    "change",
			syncStrategy:     model.SyncStrategy_QUICK_SYNC,
			stages:           quickSyncStages(),
			expectedStatus:   model.DeploymentStatus_DEPLOYMENT_SUCCESS,
			expectedExecuted: []string{"sync"},
			expectedStages: map[string]model.StageStatus{
				"sync":     model.StageStatus_STAGE_SUCCESS,
				"rollback": model.StageStatus_STAGE_NOT_STARTED_YET,
			},
		},
		{
			name:             "first deployment with pipeline is not skipped by path pattern",
			appConfig:        skipOnPathsAppConfig,
			commitMessage:    "change",
			syncStrategy:     model.SyncStrategy_PIPELINE,
			stages:           pipelineStages(),
			expectedStatus:   model.DeploymentStatus_DEPLOYMENT_SUCCESS,
			expectedExecuted: []string{"stage-a", "stage-b"},
			expectedStages: map[string]model.StageStatus{
				"stage-a":  model.StageStatus_STAGE_SUCCESS,
				"stage-b":  model.StageStatus_STAGE_SUCCESS,
				"rollback": model.StageStatus_STAGE_NOT_STARTED_YET,
			},
		},
		{
			name:             "first deployment with pipeline is skipped by commit message prefix",
			appConfig:        skipOnPrefixAppConfig,
			commitMessage:    "[skip-a] change",
			syncStrategy:     model.SyncStrategy_PIPELINE,
			stages:           pipelineStages(),
			expectedStatus:   model.DeploymentStatus_DEPLOYMENT_SUCCESS,
			expectedExecuted: []string{"stage-b"},
			expectedStages: map[string]model.StageStatus{
				"stage-a":  model.StageStatus_STAGE_SKIPPED,
				"stage-b":  model.StageStatus_STAGE_SUCCESS,
				"rollback": model.StageStatus_STAGE_NOT_STARTED_YET,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			d := newCancelRaceDeployment()
			d.RunningCommitHash = tc.runningCommitHash
			d.Trigger.SyncStrategy = tc.syncStrategy
			d.Stages = tc.stages

			p := &skipOnPlugin{failStages: tc.failStages}
			pr, err := plugin.NewPluginRegistry(context.TODO(), []plugin.Plugin{
				{Name: "plugin", Cli: p},
			})
			require.NoError(t, err)

			gc := &skipOnGitClient{
				appConfig: tc.appConfig,
				repo: &skipOnGitRepo{
					commitMessage: tc.commitMessage,
					changedFiles:  []string{"docs/README.md"},
				},
			}

			s := newScheduler(
				d,
				t.TempDir(),
				&cancelRaceAPIClient{},
				gc,
				pr,
				&cancelRaceNotifier{},
				nil,
				&cancelRaceCommandReporter{},
				&cancelRaceMetadataStore{},
				logger,
				noop.NewTracerProvider(),
			)

			require.NoError(t, s.Run(context.Background()))

			assert.Equal(t, tc.expectedStatus, s.DoneDeploymentStatus())
			assert.Equal(t, tc.expectedExecuted, p.executed)
			assert.Equal(t, tc.expectedStages, s.stageStatuses)
		})
	}
}
