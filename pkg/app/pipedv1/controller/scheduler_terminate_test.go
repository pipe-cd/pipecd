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
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/model"
	pluginapi "github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
)

// controlPlaneState keeps the statuses the scheduler reported, the way the
// control-plane would. Like a real gRPC call, a report fails once the caller's
// context is done.
type controlPlaneState struct {
	apiClient

	mu               sync.Mutex
	stageReports     []model.StageStatus
	stages           map[string]model.StageStatus
	deployment       model.DeploymentStatus
	completedStatus  model.DeploymentStatus
	completedReason  string
	deploymentClosed bool
}

func newControlPlaneState() *controlPlaneState {
	return &controlPlaneState{stages: map[string]model.StageStatus{}}
}

func (c *controlPlaneState) ReportDeploymentStatusChanged(ctx context.Context, req *pipedservice.ReportDeploymentStatusChangedRequest, opts ...grpc.CallOption) (*pipedservice.ReportDeploymentStatusChangedResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.FromContextError(err).Err()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.deployment = req.Status
	return &pipedservice.ReportDeploymentStatusChangedResponse{}, nil
}

func (c *controlPlaneState) ReportDeploymentCompleted(ctx context.Context, req *pipedservice.ReportDeploymentCompletedRequest, opts ...grpc.CallOption) (*pipedservice.ReportDeploymentCompletedResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.FromContextError(err).Err()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.deployment = req.Status
	c.completedStatus = req.Status
	c.completedReason = req.StatusReason
	c.deploymentClosed = true
	return &pipedservice.ReportDeploymentCompletedResponse{}, nil
}

func (c *controlPlaneState) ReportStageStatusChanged(ctx context.Context, req *pipedservice.ReportStageStatusChangedRequest, opts ...grpc.CallOption) (*pipedservice.ReportStageStatusChangedResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.FromContextError(err).Err()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.stageReports = append(c.stageReports, req.Status)
	c.stages[req.StageId] = req.Status
	return &pipedservice.ReportStageStatusChangedResponse{}, nil
}

func (c *controlPlaneState) ReportApplicationMostRecentDeployment(ctx context.Context, req *pipedservice.ReportApplicationMostRecentDeploymentRequest, opts ...grpc.CallOption) (*pipedservice.ReportApplicationMostRecentDeploymentResponse, error) {
	return &pipedservice.ReportApplicationMostRecentDeploymentResponse{}, nil
}

// terminatePlugin records the stages it is asked to execute. When block is set,
// ExecuteStage runs until its context is cancelled and then fails the way a
// gRPC call does.
type terminatePlugin struct {
	pluginapi.PluginClient

	block   bool
	started chan struct{}
	once    sync.Once

	mu       sync.Mutex
	executed []string
}

func newTerminatePlugin(block bool) *terminatePlugin {
	return &terminatePlugin{block: block, started: make(chan struct{})}
}

func (p *terminatePlugin) Close() error { return nil }

func (p *terminatePlugin) FetchDefinedStages(ctx context.Context, req *deployment.FetchDefinedStagesRequest, opts ...grpc.CallOption) (*deployment.FetchDefinedStagesResponse, error) {
	return &deployment.FetchDefinedStagesResponse{Stages: []string{"stage-name", "rollback-name"}}, nil
}

func (p *terminatePlugin) ExecuteStage(ctx context.Context, req *deployment.ExecuteStageRequest, opts ...grpc.CallOption) (*deployment.ExecuteStageResponse, error) {
	p.mu.Lock()
	p.executed = append(p.executed, req.Input.Stage.Id)
	p.mu.Unlock()

	if !p.block {
		return &deployment.ExecuteStageResponse{Status: model.StageStatus_STAGE_SUCCESS}, nil
	}
	p.once.Do(func() { close(p.started) })
	<-ctx.Done()
	return nil, status.Error(grpccodes.Canceled, "context canceled")
}

func (p *terminatePlugin) executedStages() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]string(nil), p.executed...)
}

func newTerminateDeployment(deploymentStatus model.DeploymentStatus, stageStatus model.StageStatus) *model.Deployment {
	d := newCancelRaceDeployment()
	d.Status = deploymentStatus
	d.Stages = []*model.PipelineStage{
		{
			Id:     "stage-1",
			Name:   "stage-name",
			Index:  0,
			Status: stageStatus,
		},
		{
			Id:       "rollback",
			Name:     "rollback-name",
			Index:    0,
			Rollback: true,
			Status:   model.StageStatus_STAGE_NOT_STARTED_YET,
		},
	}
	return d
}

func newTerminateScheduler(t *testing.T, d *model.Deployment, cp *controlPlaneState, p *terminatePlugin) *scheduler {
	t.Helper()
	pr, err := plugin.NewPluginRegistry(context.TODO(), []plugin.Plugin{
		{Name: "plugin", Cli: p},
	})
	require.NoError(t, err)

	return newScheduler(
		d,
		t.TempDir(),
		cp,
		&cancelRaceGitClient{},
		pr,
		&cancelRaceNotifier{},
		nil,
		&cancelRaceCommandReporter{},
		&cancelRaceMetadataStore{},
		zaptest.NewLogger(t),
		noop.NewTracerProvider(),
	)
}

func waitStarted(t *testing.T, p *terminatePlugin) {
	t.Helper()
	select {
	case <-p.started:
	case <-time.After(10 * time.Second):
		t.Fatal("the stage was not started")
	}
}

// TestExecuteStage_StoppedWhileRunning stops a stage while its plugin is still
// executing it, once for each stop signal.
func TestExecuteStage_StoppedWhileRunning(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name         string
		stop         func(StopSignalHandler)
		wantStatus   model.StageStatus
		wantReported []model.StageStatus
	}{
		{
			name:       "terminated: the stage is left running so that it is resumed",
			stop:       func(h StopSignalHandler) { h.Terminate() },
			wantStatus: model.StageStatus_STAGE_RUNNING,
			wantReported: []model.StageStatus{
				model.StageStatus_STAGE_RUNNING,
			},
		},
		{
			name:       "timed out: the stage fails",
			stop:       func(h StopSignalHandler) { h.Timeout() },
			wantStatus: model.StageStatus_STAGE_FAILURE,
			wantReported: []model.StageStatus{
				model.StageStatus_STAGE_RUNNING,
				model.StageStatus_STAGE_FAILURE,
			},
		},
		{
			name:       "cancelled: the stage is cancelled",
			stop:       func(h StopSignalHandler) { h.Cancel() },
			wantStatus: model.StageStatus_STAGE_CANCELLED,
			wantReported: []model.StageStatus{
				model.StageStatus_STAGE_RUNNING,
				model.StageStatus_STAGE_CANCELLED,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cp := newControlPlaneState()
			p := newTerminatePlugin(true)
			s := newTerminateScheduler(t, newTerminateDeployment(model.DeploymentStatus_DEPLOYMENT_RUNNING, model.StageStatus_STAGE_NOT_STARTED_YET), cp, p)
			// newScheduler leaves the deploy sources to Run.
			s.targetDSP = &fakeDeploySourceProvider{}
			s.genericApplicationConfig = &config.GenericApplicationSpec{
				Pipeline: &config.DeploymentPipeline{
					Stages: []config.PipelineStage{{Name: "stage-name"}},
				},
			}

			sig, handler := NewStopSignal()
			done := make(chan model.StageStatus, 1)
			go func() {
				done <- s.executeStage(sig, s.deployment.Stages[0])
			}()
			waitStarted(t, p)
			tc.stop(handler)

			assert.Equal(t, tc.wantStatus, <-done)
			cp.mu.Lock()
			defer cp.mu.Unlock()
			assert.Equal(t, tc.wantReported, cp.stageReports)
		})
	}
}

// TestSchedulerResumesStageAfterPipedRestart stops piped while a stage is
// running and then starts a new scheduler from what the control-plane holds,
// as piped does after it restarts. The stage must be executed again, and the
// deployment must not be cancelled or rolled back.
func TestSchedulerResumesStageAfterPipedRestart(t *testing.T) {
	t.Parallel()

	cp := newControlPlaneState()

	// Piped is stopped while stage-1 is running.
	p1 := newTerminatePlugin(true)
	s1 := newTerminateScheduler(t, newTerminateDeployment(model.DeploymentStatus_DEPLOYMENT_PLANNED, model.StageStatus_STAGE_NOT_STARTED_YET), cp, p1)

	ctx, stopPiped := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- s1.Run(ctx)
	}()
	waitStarted(t, p1)
	stopPiped()
	<-done

	cp.mu.Lock()
	restarted := newTerminateDeployment(cp.deployment, cp.stages["stage-1"])
	cp.mu.Unlock()
	require.Equal(t, model.DeploymentStatus_DEPLOYMENT_RUNNING, restarted.Status)

	// Piped starts again and resumes the deployment.
	p2 := newTerminatePlugin(false)
	s2 := newTerminateScheduler(t, restarted, cp, p2)
	require.NoError(t, s2.Run(context.Background()))

	assert.Equal(t, []string{"stage-1"}, p2.executedStages())
	cp.mu.Lock()
	defer cp.mu.Unlock()
	assert.True(t, cp.deploymentClosed)
	assert.Equal(t, model.DeploymentStatus_DEPLOYMENT_SUCCESS, cp.completedStatus)
}
