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
// Whenever a new uncompleted Deployment is detected, this creates a new DeploymentScheduler
// for that Deployment to handle the deployment pipeline.
package deploymentcontroller

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/kapetaniosci/pipe/pkg/app/api/service/pipedservice"
	"github.com/kapetaniosci/pipe/pkg/app/piped/logpersister"
	"github.com/kapetaniosci/pipe/pkg/config"
)

type apiClient interface {
	ListNotCompletedDeployments(ctx context.Context, in *pipedservice.ListNotCompletedDeploymentsRequest, opts ...grpc.CallOption) (*pipedservice.ListNotCompletedDeploymentsResponse, error)
	SaveStageMetadata(ctx context.Context, in *pipedservice.SaveStageMetadataRequest, opts ...grpc.CallOption) (*pipedservice.SaveStageMetadataResponse, error)
	ReportStageStatusChanged(ctx context.Context, in *pipedservice.ReportStageStatusChangedRequest, opts ...grpc.CallOption) (*pipedservice.ReportStageStatusChangedResponse, error)
	ReportStageLog(ctx context.Context, in *pipedservice.ReportStageLogRequest, opts ...grpc.CallOption) (*pipedservice.ReportStageLogResponse, error)
	ReportDeploymentCompleted(ctx context.Context, in *pipedservice.ReportDeploymentCompletedRequest, opts ...grpc.CallOption) (*pipedservice.ReportDeploymentCompletedResponse, error)
	GetCommands(ctx context.Context, in *pipedservice.GetCommandsRequest, opts ...grpc.CallOption) (*pipedservice.GetCommandsResponse, error)
	ReportCommandHandled(ctx context.Context, in *pipedservice.ReportCommandHandledRequest, opts ...grpc.CallOption) (*pipedservice.ReportCommandHandledResponse, error)
}

type DeploymentController struct {
	apiClient         apiClient
	config            *config.PipedSpec
	schedulers        map[string]*scheduler
	logPersister      logpersister.Persister
	metadataPersister metadataPersister
	mu                sync.Mutex
	gracePeriod       time.Duration
	logger            *zap.Logger
}

// NewController creates a new instance for DeploymentController.
func NewController(apiClient apiClient, cfg *config.PipedSpec, gracePeriod time.Duration, logger *zap.Logger) *DeploymentController {
	return &DeploymentController{
		apiClient:         apiClient,
		config:            cfg,
		logPersister:      logpersister.NewPersister(apiClient, logger),
		metadataPersister: metadataPersister{apiClient: apiClient},
		gracePeriod:       gracePeriod,
		logger:            logger.Named("deployment-controller"),
	}
}

// Run starts running DeploymentController until the specified context
// has done. This also waits for its cleaning up before returning.
func (c *DeploymentController) Run(ctx context.Context) error {
	c.logger.Info("start running deployment controller")
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	go c.logPersister.Run(ctx)

L:
	for {
		select {
		case <-ctx.Done():
			break L
		case <-ticker.C:
			c.syncScheduler(ctx)
		}
	}

	// TODO: Wait for graceful shutdowns of all components.
	c.logger.Info("deployment controller has been stopped")
	return nil
}

// syncScheduler adds new scheduler for newly aaded deployments
// as well as removes the removeable deployments.
func (c *DeploymentController) syncScheduler(ctx context.Context) error {
	resp, err := c.apiClient.ListNotCompletedDeployments(ctx, &pipedservice.ListNotCompletedDeploymentsRequest{})
	if err != nil {
		return err
	}

	// Add missing schedulers.
	for _, d := range resp.Deployments {
		if _, ok := c.schedulers[d.Id]; ok {
			continue
		}
		e := newScheduler(d, c.logPersister, c.metadataPersister, c.logger)
		c.schedulers[e.Id()] = e
		go e.Run(ctx)
	}

	// Remove done schedulers.
	for id, e := range c.schedulers {
		if e.IsDone() {
			delete(c.schedulers, id)
		}
	}
	return nil
}
