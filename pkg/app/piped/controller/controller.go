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

// Package controller provides a piped component
// that managing all of the PLANNED and not completed deployments.
// This manages a pool of DeploymentSchedulers.
// Whenever a new PLANNED Deployment is detected,
// this creates a new DeploymentScheduler
// for that Deployment to handle the deployment pipeline.
package controller

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/kapetaniosci/pipe/pkg/app/api/service/pipedservice"
	"github.com/kapetaniosci/pipe/pkg/app/piped/logpersister"
	"github.com/kapetaniosci/pipe/pkg/config"
	"github.com/kapetaniosci/pipe/pkg/git"
	"github.com/kapetaniosci/pipe/pkg/model"
)

type apiClient interface {
	ReportDeploymentPlanned(ctx context.Context, req *pipedservice.ReportDeploymentPlannedRequest, opts ...grpc.CallOption) (*pipedservice.ReportDeploymentPlannedResponse, error)
	ReportDeploymentRunning(ctx context.Context, req *pipedservice.ReportDeploymentRunningRequest, opts ...grpc.CallOption) (*pipedservice.ReportDeploymentRunningResponse, error)
	ReportDeploymentCompleted(ctx context.Context, req *pipedservice.ReportDeploymentCompletedRequest, opts ...grpc.CallOption) (*pipedservice.ReportDeploymentCompletedResponse, error)
	SaveDeploymentMetadata(ctx context.Context, req *pipedservice.SaveDeploymentMetadataRequest, opts ...grpc.CallOption) (*pipedservice.SaveDeploymentMetadataResponse, error)

	SaveStageMetadata(ctx context.Context, req *pipedservice.SaveStageMetadataRequest, opts ...grpc.CallOption) (*pipedservice.SaveStageMetadataResponse, error)
	ReportStageLog(ctx context.Context, req *pipedservice.ReportStageLogRequest, opts ...grpc.CallOption) (*pipedservice.ReportStageLogResponse, error)
	ReportStageStatusChanged(ctx context.Context, req *pipedservice.ReportStageStatusChangedRequest, opts ...grpc.CallOption) (*pipedservice.ReportStageStatusChangedResponse, error)
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

type deploymentLister interface {
	ListPlanneds() []*model.Deployment
}

type commandLister interface {
	ListDeploymentCommands(deploymentID string) []model.ReportableCommand
	ListStageCommands(deploymentID, stageID string) []model.ReportableCommand
}

type DeploymentController interface {
	Run(ctx context.Context) error
}

type controller struct {
	apiClient        apiClient
	gitClient        gitClient
	deploymentLister deploymentLister
	commandLister    commandLister
	pipedConfig      *config.PipedSpec
	logPersister     logpersister.Persister

	schedulers            map[string]*scheduler
	unreportedDeployments map[string]*model.Deployment
	wg                    sync.WaitGroup

	workspaceDir string
	syncInternal time.Duration
	gracePeriod  time.Duration
	logger       *zap.Logger
}

// NewController creates a new instance for DeploymentController.
func NewController(
	apiClient apiClient,
	gitClient gitClient,
	deploymentLister deploymentLister,
	commandLister commandLister,
	pipedConfig *config.PipedSpec,
	gracePeriod time.Duration,
	logger *zap.Logger,
) DeploymentController {

	var (
		lp = logpersister.NewPersister(apiClient, logger)
		lg = logger.Named("deployment-controller")
	)
	return &controller{
		apiClient:             apiClient,
		gitClient:             gitClient,
		deploymentLister:      deploymentLister,
		commandLister:         commandLister,
		pipedConfig:           pipedConfig,
		logPersister:          lp,
		schedulers:            make(map[string]*scheduler),
		unreportedDeployments: make(map[string]*model.Deployment),
		syncInternal:          10 * time.Second,
		gracePeriod:           gracePeriod,
		logger:                lg,
	}
}

// Run starts running DeploymentController until the specified context has done.
// This also waits for its cleaning up before returning.
func (c *controller) Run(ctx context.Context) error {
	c.logger.Info("start running deployment controller")

	// Make sure the existence of the workspace directory.
	// Each scheduler/deployment will have an working directory inside this workspace.
	dir, err := ioutil.TempDir("", "workspace")
	if err != nil {
		c.logger.Error("failed to create workspace directory", zap.Error(err))
		return err
	}
	c.workspaceDir = dir
	c.logger.Info(fmt.Sprintf("workspace directory was configured to %s", c.workspaceDir))

	// Start running log persister to buffer and flush the log blocks.
	// We do not use the passed ctx directly because we want log persister
	// component to be stopped at the last order to avoid lossing log from other components.
	var (
		lpStoppedCh     = make(chan error, 1)
		lpCtx, lpCancel = context.WithCancel(context.Background())
	)
	go func() {
		lpStoppedCh <- c.logPersister.Run(lpCtx)
		close(lpStoppedCh)
	}()

	ticker := time.NewTicker(c.syncInternal)
	defer ticker.Stop()

	c.logger.Info("start syncing schedulers")
L:
	for {
		select {
		case <-ctx.Done():
			break L
		case <-ticker.C:
			c.syncScheduler(ctx)
		}
	}
	c.logger.Info("stop syncing schedulers")

	c.logger.Info("waiting for stopping all executors")
	c.wg.Wait()

	// Stop log persiter and wait for its stopping.
	lpCancel()
	err = <-lpStoppedCh

	c.logger.Info("deployment controller has been stopped")
	return err
}

// syncScheduler adds new scheduler for newly added deployments
// as well as removes the removable deployments.
func (c *controller) syncScheduler(ctx context.Context) error {
	planneds := c.deploymentLister.ListPlanneds()

	c.logger.Info(fmt.Sprintf("there are %d planned deployments for this piped", len(planneds)),
		zap.Int("scheduler-count", len(c.schedulers)),
	)

	// Add missing schedulers.
	for _, d := range planneds {
		if _, ok := c.schedulers[d.Id]; ok {
			continue
		}
		if err := c.startNewScheduler(ctx, d); err != nil {
			continue
		}
	}

	// Remove done schedulers.
	for id, e := range c.schedulers {
		if !e.IsDone() {
			continue
		}
		c.unreportedDeployments[id] = e.deployment

		c.logger.Info("deleted done scheduler",
			zap.String("deployment-id", id),
			zap.String("application-id", e.deployment.ApplicationId),
		)
		delete(c.schedulers, id)
	}

	return nil
}

// startNewScheduler creates and starts running a new scheduler
// for a specific PLANNED deployment.
// This adds the newly created one to the scheduler list
// for tracking its lifetime periodically later.
func (c *controller) startNewScheduler(ctx context.Context, d *model.Deployment) error {
	logger := c.logger.With(
		zap.String("deployment-id", d.Id),
		zap.String("application-id", d.ApplicationId),
	)
	logger.Info("will add a new scheduler")

	// Ensure the existence of the working directory for the deployment.
	workingDir, err := ioutil.TempDir(c.workspaceDir, d.Id+"-*")
	if err != nil {
		logger.Error("failed to create working directory for scheduler", zap.Error(err))
		return err
	}
	logger.Info("created working directory for scheduler", zap.String("working-dir", workingDir))

	// Create a new scheduler and append to the list for tracking.
	e := newScheduler(
		d,
		workingDir,
		c.apiClient,
		c.gitClient,
		c.commandLister,
		c.logPersister,
		c.pipedConfig,
		c.logger,
	)
	c.schedulers[e.Id()] = e
	logger.Info("added a new scheduler", zap.Int("scheduler-count", len(c.schedulers)))

	cleanup := func() {
		logger.Info("cleaning up working directory for scheduler", zap.String("working-dir", workingDir))
		err := os.RemoveAll(workingDir)
		if err == nil {
			return
		}
		logger.Warn("failed to clean working directory",
			zap.String("working-dir", workingDir),
			zap.Error(err),
		)
	}

	// Start running executor.
	go func() {
		c.wg.Add(1)
		defer c.wg.Done()
		defer cleanup()
		e.Run(ctx)
	}()

	return nil
}
