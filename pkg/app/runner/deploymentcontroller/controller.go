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

// Package deploymentcontroller provides a runner component
// that managing all of the not completed deployments.
// This manages a pool of DeploymentExecutors.
// Whenever a new uncompleted Deployment is detected, this creates a new DeploymentExecutor
// for that Deployment to handle the deployment pipeline.
// The DeploymentExecutor will update the deployment status back to the API.
package deploymentcontroller

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/kapetaniosci/pipe/pkg/app/api/service/runnerservice"
	"github.com/kapetaniosci/pipe/pkg/config"
)

const (
	subsetLabel    = "pipecd.dev/subset"
	managedByLabel = "pipecd.dev/managed-by"
)

type apiClient interface {
	ListNotCompletedDeployments(ctx context.Context, in *runnerservice.ListNotCompletedDeploymentsRequest, opts ...grpc.CallOption) (*runnerservice.ListNotCompletedDeploymentsResponse, error)
	SendStageLog(ctx context.Context, in *runnerservice.SendStageLogRequest, opts ...grpc.CallOption) (*runnerservice.SendStageLogResponse, error)
	RegisterEvents(ctx context.Context, in *runnerservice.RegisterEventsRequest, opts ...grpc.CallOption) (*runnerservice.RegisterEventsResponse, error)
	GetCommands(ctx context.Context, in *runnerservice.GetCommandsRequest, opts ...grpc.CallOption) (*runnerservice.GetCommandsResponse, error)
	ReportCommandHandled(ctx context.Context, in *runnerservice.ReportCommandHandledRequest, opts ...grpc.CallOption) (*runnerservice.ReportCommandHandledResponse, error)
	ReportDeploymentCompleted(ctx context.Context, in *runnerservice.ReportDeploymentCompletedRequest, opts ...grpc.CallOption) (*runnerservice.ReportDeploymentCompletedResponse, error)
}

type DeploymentController struct {
	apiClient   apiClient
	config      *config.RunnerSpec
	executors   map[string]*executor
	mu          sync.Mutex
	gracePeriod time.Duration
	logger      *zap.Logger
}

// NewController creates a new instance for DeploymentController.
func NewController(apiClient apiClient, cfg *config.RunnerSpec, gracePeriod time.Duration, logger *zap.Logger) *DeploymentController {
	return &DeploymentController{
		apiClient:   apiClient,
		config:      cfg,
		gracePeriod: gracePeriod,
		logger:      logger.Named("deployment-controller"),
	}
}

// Run starts running DeploymentController until the specified context
// has done. This also waits for its cleaning up before returning.
func (c *DeploymentController) Run(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.syncExecutor(ctx)
			}
		}
	}()
	return nil
}

func (c *DeploymentController) syncExecutor(ctx context.Context) error {
	resp, err := c.apiClient.ListNotCompletedDeployments(ctx, &runnerservice.ListNotCompletedDeploymentsRequest{})
	if err != nil {
		return err
	}
	// Add missing executors.
	for _, d := range resp.Deployments {
		if _, ok := c.executors[d.Id]; ok {
			continue
		}
		e := newExecutor(d, c.logger)
		c.executors[e.Id()] = e
		go e.Run(ctx)
	}

	// Remove done executors.
	for id, e := range c.executors {
		if e.IsDone() {
			delete(c.executors, id)
		}
	}
	return nil
}
