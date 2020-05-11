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

// Package deploymentcontroller provides a piped component
// that managing all of the not completed deployments.
// This manages a pool of DeploymentSchedulers.
// Whenever a new uncompleted Deployment is detected,
// this creates a new DeploymentScheduler
// for that Deployment to handle the deployment pipeline.
package deploymentcontroller

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	ListNotCompletedDeployments(ctx context.Context, in *pipedservice.ListNotCompletedDeploymentsRequest, opts ...grpc.CallOption) (*pipedservice.ListNotCompletedDeploymentsResponse, error)
	SaveStageMetadata(ctx context.Context, in *pipedservice.SaveStageMetadataRequest, opts ...grpc.CallOption) (*pipedservice.SaveStageMetadataResponse, error)
	ReportStageLog(ctx context.Context, in *pipedservice.ReportStageLogRequest, opts ...grpc.CallOption) (*pipedservice.ReportStageLogResponse, error)
	ReportStageStatusChanged(ctx context.Context, in *pipedservice.ReportStageStatusChangedRequest, opts ...grpc.CallOption) (*pipedservice.ReportStageStatusChangedResponse, error)
	ReportDeploymentCompleted(ctx context.Context, in *pipedservice.ReportDeploymentCompletedRequest, opts ...grpc.CallOption) (*pipedservice.ReportDeploymentCompletedResponse, error)
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

type commandStore interface {
	ListDeploymentCommands(deploymentID string) []*model.Command
	ReportCommandHandled(ctx context.Context, c *model.Command, status model.CommandStatus, metadata map[string]string) error
}

type DeploymentController struct {
	apiClient         apiClient
	gitClient         gitClient
	commandStore      commandStore
	pipedConfig       *config.PipedSpec
	logPersister      logpersister.Persister
	metadataPersister metadataPersister

	schedulers map[string]*scheduler
	wg         sync.WaitGroup
	mu         sync.Mutex

	workingDir   string
	syncInternal time.Duration
	gracePeriod  time.Duration
	logger       *zap.Logger
}

// NewController creates a new instance for DeploymentController.
func NewController(apiClient apiClient, gitClient gitClient, cmdStore commandStore, cfg *config.PipedSpec, gracePeriod time.Duration, logger *zap.Logger) *DeploymentController {
	var (
		lp  = logpersister.NewPersister(apiClient, logger)
		mdp = metadataPersister{apiClient: apiClient}
		lg  = logger.Named("deployment-controller")
	)
	return &DeploymentController{
		apiClient:         apiClient,
		gitClient:         gitClient,
		commandStore:      cmdStore,
		pipedConfig:       cfg,
		logPersister:      lp,
		metadataPersister: mdp,
		schedulers:        make(map[string]*scheduler),
		syncInternal:      30 * time.Second,
		gracePeriod:       gracePeriod,
		logger:            lg,
	}
}

// Run starts running DeploymentController until the specified context has done.
// This also waits for its cleaning up before returning.
func (c *DeploymentController) Run(ctx context.Context) error {
	c.logger.Info("start running deployment controller")

	dir, err := ioutil.TempDir("", "working")
	if err != nil {
		c.logger.Error("failed to create working directory", zap.Error(err))
		return err
	}
	c.workingDir = dir
	c.logger.Info(fmt.Sprintf("working directory was configured to %s", c.workingDir))

	ticker := time.NewTicker(c.syncInternal)
	defer ticker.Stop()

	doneCh := make(chan error, 1)
	go func() {
		doneCh <- c.logPersister.Run(ctx)
		close(doneCh)
	}()

L:
	for {
		select {
		case <-ctx.Done():
			break L
		case <-ticker.C:
			c.syncScheduler(ctx)
		}
	}

	// Waiting for stopping of log persister.
	err = <-doneCh

	c.logger.Info("waiting for stopping all executors")
	c.wg.Wait()

	c.logger.Info("deployment controller has been stopped")
	return err
}

// syncScheduler adds new scheduler for newly added deployments
// as well as removes the removable deployments.
func (c *DeploymentController) syncScheduler(ctx context.Context) error {
	resp, err := c.apiClient.ListNotCompletedDeployments(ctx, &pipedservice.ListNotCompletedDeploymentsRequest{})
	if err != nil {
		c.logger.Warn("failed to list uncompleted deployments", zap.Error(err))
		return err
	}
	c.logger.Info(fmt.Sprintf("there are %d uncompleted deployments for this piped", len(resp.Deployments)),
		zap.Int("scheduler-count", len(c.schedulers)),
	)

	// Add missing schedulers.
	for _, d := range resp.Deployments {
		if _, ok := c.schedulers[d.Id]; ok {
			continue
		}
		if err := c.startNewScheduler(ctx, d); err != nil {
			continue
		}
	}

	// Remove done schedulers.
	for id, e := range c.schedulers {
		if e.IsDone() {
			c.logger.Info("deleted finished scheduler",
				zap.String("deployment-id", id),
				zap.String("application-id", e.deployment.ApplicationId),
			)
			delete(c.schedulers, id)
		}
	}

	return nil
}

func (c *DeploymentController) startNewScheduler(ctx context.Context, d *model.Deployment) error {
	logger := c.logger.With(
		zap.String("deployment-id", d.Id),
		zap.String("application-id", d.ApplicationId),
	)
	logger.Info("will add a new scheduler")

	// Ensure the existence of the working directory for the deployment.
	workingDir := filepath.Join(c.workingDir, d.Id)
	if err := os.MkdirAll(workingDir, os.ModePerm); err != nil {
		logger.Error("failed to create working directory for deployment",
			zap.String("working-dir", workingDir),
			zap.Error(err),
		)
		return err
	}
	logger.Info("created working directory for deployment", zap.String("working-dir", workingDir))

	// Create a new scheduler and append to the list for tracking.
	e := newScheduler(
		d,
		c.pipedConfig,
		workingDir,
		c.gitClient,
		c.commandStore,
		c.logPersister,
		c.metadataPersister,
		c.logger,
	)
	c.schedulers[e.Id()] = e
	logger.Info("added a new scheduler", zap.Int("scheduler-count", len(c.schedulers)))

	// Start running executor.
	cleanup := func() {
		logger.Info("cleaning up working directory for deployment", zap.String("working-dir", workingDir))
		err := os.RemoveAll(workingDir)
		if err == nil {
			return
		}
		logger.Warn("failed to clean working directory",
			zap.String("working-dir", workingDir),
			zap.Error(err),
		)
	}
	go func() {
		c.wg.Add(1)
		defer c.wg.Done()
		defer cleanup()
		e.Run(ctx)
	}()

	return nil
}
