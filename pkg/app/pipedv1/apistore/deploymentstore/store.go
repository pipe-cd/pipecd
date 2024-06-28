// Copyright 2024 The PipeCD Authors.
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

package deploymentstore

import (
	"context"
	"time"

	"go.uber.org/atomic"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/model"
)

// Lister helps list and get deployment.
// All objects returned here must be treated as read-only.
type Lister interface {
	// ListPendings lists all pending deployments that should be handled by this piped.
	ListPendings() []*model.Deployment
	// ListPlanneds lists all planned deployments that should be handled by this piped.
	ListPlanneds() []*model.Deployment
	// ListRunnings lists all running deployments that should be handled by this piped.
	ListRunnings() []*model.Deployment
	// ListAppHeadDeployments returns the map from application ID to its head of deploying deployments.
	// Application that currently has not any uncompleted deployments will not be in this list.
	// Head deployment is same with the oldest uncompleted one.
	ListAppHeadDeployments() map[string]*model.Deployment
}

type apiClient interface {
	ListNotCompletedDeployments(ctx context.Context, in *pipedservice.ListNotCompletedDeploymentsRequest, opts ...grpc.CallOption) (*pipedservice.ListNotCompletedDeploymentsResponse, error)
}

type Store interface {
	// Run starts syncing the deployment list with the control-plane.
	Run(ctx context.Context) error
	// Lister returns a lister for retrieving deployments.
	Lister() Lister
}

type store struct {
	apiClient          apiClient
	pendingDeployments atomic.Value
	plannedDeployments atomic.Value
	runningDeployments atomic.Value
	headDeployments    atomic.Value
	syncInterval       time.Duration
	gracePeriod        time.Duration
	logger             *zap.Logger
}

var (
	defaultSyncInterval = 10 * time.Second
)

// NewStore creates a new deployment store instance.
// This syncs with the control plane to keep the list of deployments for this runner up-to-date.
func NewStore(apiClient apiClient, gracePeriod time.Duration, logger *zap.Logger) Store {
	return &store{
		apiClient:    apiClient,
		syncInterval: defaultSyncInterval,
		gracePeriod:  gracePeriod,
		logger:       logger.Named("deployment-store"),
	}
}

// Run starts syncing the deployment list with the control-plane.
func (s *store) Run(ctx context.Context) error {
	s.logger.Info("start running deployment store")

	syncTicker := time.NewTicker(s.syncInterval)
	defer syncTicker.Stop()

	for {
		select {
		case <-syncTicker.C:
			s.sync(ctx)

		case <-ctx.Done():
			s.logger.Info("deployment store has been stopped")
			return nil
		}
	}
}

// Lister returns a lister for retrieving deployments.
func (s *store) Lister() Lister {
	return s
}

func (s *store) sync(ctx context.Context) error {
	resp, err := s.apiClient.ListNotCompletedDeployments(ctx, &pipedservice.ListNotCompletedDeploymentsRequest{})
	if err != nil {
		s.logger.Error("failed to list unhandled deployment", zap.Error(err))
		return err
	}

	// TODO: Call ListNotCompletedDeployments itervally until all required deployments are fetched.
	var pendings, planneds, runnings []*model.Deployment
	for _, d := range resp.Deployments {
		switch d.Status {
		case model.DeploymentStatus_DEPLOYMENT_PENDING:
			pendings = append(pendings, d)
		case model.DeploymentStatus_DEPLOYMENT_PLANNED:
			planneds = append(planneds, d)
		case model.DeploymentStatus_DEPLOYMENT_RUNNING:
			runnings = append(runnings, d)
		}
	}

	headDeployments := make(map[string]*model.Deployment)
	for _, d := range pendings {
		headDeployments[d.ApplicationId] = d
	}
	for _, d := range planneds {
		headDeployments[d.ApplicationId] = d
	}
	for _, d := range runnings {
		headDeployments[d.ApplicationId] = d
	}

	s.plannedDeployments.Store(planneds)
	s.runningDeployments.Store(runnings)
	s.pendingDeployments.Store(pendings)
	s.headDeployments.Store(headDeployments)

	return nil
}

// ListPendings lists all pending deployments that should be handled by this piped.
func (s *store) ListPendings() []*model.Deployment {
	list := s.pendingDeployments.Load()
	if list == nil {
		return nil
	}
	return list.([]*model.Deployment)
}

// ListPlanneds lists all planned deployments that should be handled by this piped.
func (s *store) ListPlanneds() []*model.Deployment {
	list := s.plannedDeployments.Load()
	if list == nil {
		return nil
	}
	return list.([]*model.Deployment)
}

// ListRunnings lists all running deployments that should be handled by this piped.
func (s *store) ListRunnings() []*model.Deployment {
	list := s.runningDeployments.Load()
	if list == nil {
		return nil
	}
	return list.([]*model.Deployment)
}

func (s *store) ListAppHeadDeployments() map[string]*model.Deployment {
	m := s.headDeployments.Load()
	if m == nil {
		return nil
	}
	return m.(map[string]*model.Deployment)
}
