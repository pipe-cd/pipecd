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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
)

const (
	runningCommit = "running-commit"
	targetCommit  = "target-commit"

	renamedAppConfig = `apiVersion: pipecd.dev/v1beta1
kind: Application
spec:
  name: renamed
  plugins:
    plugin-1: {}
  pipeline:
    stages:
      - name: plugin-1-stage-1
`
)

// renamedConfigGitRepo writes the application config file under the name it has at the checked out commit.
type renamedConfigGitRepo struct {
	git.Repo
	dir             string
	configFilenames map[string]string
}

func (r *renamedConfigGitRepo) Checkout(ctx context.Context, commitish string) error {
	return os.WriteFile(filepath.Join(r.dir, r.configFilenames[commitish]), []byte(renamedAppConfig), 0600)
}

type renamedConfigGitClient struct {
	configFilenames map[string]string
}

func (c *renamedConfigGitClient) Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error) {
	if err := os.MkdirAll(destination, 0700); err != nil {
		return nil, err
	}
	return &renamedConfigGitRepo{dir: destination, configFilenames: c.configFilenames}, nil
}

type renamedConfigAPIClient struct {
	cancelRaceAPIClient
	planned bool
}

func (c *renamedConfigAPIClient) ReportDeploymentPlanned(ctx context.Context, req *pipedservice.ReportDeploymentPlannedRequest, opts ...grpc.CallOption) (*pipedservice.ReportDeploymentPlannedResponse, error) {
	c.planned = true
	return &pipedservice.ReportDeploymentPlannedResponse{}, nil
}

// renamedConfigPlugin records the config filename of the running deploy source it receives.
type renamedConfigPlugin struct {
	fakePlugin
	runningConfigFilename string
}

func (p *renamedConfigPlugin) ExecuteStage(ctx context.Context, req *deployment.ExecuteStageRequest, opts ...grpc.CallOption) (*deployment.ExecuteStageResponse, error) {
	p.runningConfigFilename = req.Input.RunningDeploymentSource.GetApplicationConfigFilename()
	return &deployment.ExecuteStageResponse{Status: model.StageStatus_STAGE_SUCCESS}, nil
}

func newRenamedConfigDeployment() *model.Deployment {
	d := newCancelRaceDeployment()
	d.GitPath.ConfigFilename = "new.pipecd.yaml"
	d.Trigger.Commit.Hash = targetCommit
	return d
}

func TestPlannerLoadsRunningSourceWithItsConfigFilename(t *testing.T) {
	pr, err := plugin.NewPluginRegistry(context.TODO(), []plugin.Plugin{{
		Name: "plugin-1",
		Cli: &fakePlugin{
			quickStages: []*model.PipelineStage{{Name: "plugin-1-stage-1"}},
		},
	}})
	require.NoError(t, err)

	d := newRenamedConfigDeployment()
	d.Status = model.DeploymentStatus_DEPLOYMENT_PENDING
	d.Stages = nil
	d.Trigger.SyncStrategy = model.SyncStrategy_QUICK_SYNC

	gc := &renamedConfigGitClient{configFilenames: map[string]string{
		runningCommit: "app.pipecd.yaml",
		targetCommit:  "new.pipecd.yaml",
	}}
	api := &renamedConfigAPIClient{}

	p := newPlanner(
		d,
		runningCommit,
		"app.pipecd.yaml",
		t.TempDir(),
		pr,
		api,
		gc,
		&cancelRaceNotifier{},
		nil,
		&cancelRaceMetadataStore{},
		zaptest.NewLogger(t),
		noop.NewTracerProvider(),
	)

	require.NoError(t, p.Run(context.Background()))
	assert.True(t, api.planned)
	assert.Equal(t, model.DeploymentStatus_DEPLOYMENT_PLANNED, p.DoneDeploymentStatus())
}

func TestSchedulerLoadsRunningSourceWithItsConfigFilename(t *testing.T) {
	pl := &renamedConfigPlugin{
		fakePlugin: fakePlugin{
			pipelineStages: []*model.PipelineStage{{Name: "plugin-1-stage-1"}},
		},
	}
	pr, err := plugin.NewPluginRegistry(context.TODO(), []plugin.Plugin{{Name: "plugin-1", Cli: pl}})
	require.NoError(t, err)

	d := newRenamedConfigDeployment()
	d.RunningCommitHash = runningCommit
	d.RunningConfigFilename = "app.pipecd.yaml"
	d.Stages = []*model.PipelineStage{{Id: "stage-1", Name: "plugin-1-stage-1", Index: 0}}

	gc := &renamedConfigGitClient{configFilenames: map[string]string{
		runningCommit: "app.pipecd.yaml",
		targetCommit:  "new.pipecd.yaml",
	}}

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
		zaptest.NewLogger(t),
		noop.NewTracerProvider(),
	)

	require.NoError(t, s.Run(context.Background()))
	assert.Equal(t, model.DeploymentStatus_DEPLOYMENT_SUCCESS, s.DoneDeploymentStatus())
	assert.Equal(t, "app.pipecd.yaml", pl.runningConfigFilename)
}

func TestRunningGitPath(t *testing.T) {
	gp := &model.ApplicationGitPath{Path: "app", ConfigFilename: "new.pipecd.yaml"}

	assert.Same(t, gp, runningGitPath(gp, ""))

	got := runningGitPath(gp, "app.pipecd.yaml")
	assert.Equal(t, "app.pipecd.yaml", got.ConfigFilename)
	assert.Equal(t, "app", got.Path)
	assert.Equal(t, "new.pipecd.yaml", gp.ConfigFilename)
}
