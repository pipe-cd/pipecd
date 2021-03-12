// Copyright 2021 The PipeCD Authors.
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

package insightcollector

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/filestore"
	"github.com/pipe-cd/pipe/pkg/insight"
	"github.com/pipe-cd/pipe/pkg/insight/insightstore"
	"github.com/pipe-cd/pipe/pkg/model"
)

// InsightCollector implements the behaviors for the gRPC definitions of InsightCollector.
type InsightCollector struct {
	projectStore     datastore.ProjectStore
	applicationStore datastore.ApplicationStore
	deploymentStore  datastore.DeploymentStore
	insightstore     insightstore.Store

	applicationsHandlers              []func(ctx context.Context, applications []*model.Application, target time.Time) error
	newlyCreatedDeploymentsHandlers   []func(ctx context.Context, developments []*model.Deployment, target time.Time) error
	newlyCompletedDeploymentsHandlers []func(ctx context.Context, developments []*model.Deployment, target time.Time) error

	logger *zap.Logger
}

// NewInsightCollector creates a new InsightCollector instance.
func NewInsightCollector(ds datastore.DataStore, fs filestore.Store, metrics CollectorMetrics, logger *zap.Logger) *InsightCollector {
	c := &InsightCollector{
		projectStore:     datastore.NewProjectStore(ds),
		applicationStore: datastore.NewApplicationStore(ds),
		deploymentStore:  datastore.NewDeploymentStore(ds),
		insightstore:     insightstore.NewStore(fs),
		logger:           logger.Named("insight-collector"),
	}
	c.setHandlers(metrics)

	return c
}

func (c *InsightCollector) setHandlers(metrics CollectorMetrics) {
	if metrics.IsEnabled(ApplicationCount) {
		c.applicationsHandlers = append(c.applicationsHandlers, c.collectApplicationCount)
	}
	if metrics.IsEnabled(DevelopmentFrequency) {
		c.newlyCreatedDeploymentsHandlers = append(c.newlyCreatedDeploymentsHandlers, c.collectDevelopmentFrequency)
	}
	if metrics.IsEnabled(ChangeFailureRate) {
		c.newlyCompletedDeploymentsHandlers = append(c.newlyCompletedDeploymentsHandlers, c.collectDeploymentChangeFailureRate)
	}
}

func (c *InsightCollector) ProcessNewlyCreatedDeployments(ctx context.Context) error {
	c.logger.Info("will retrieve newly created deployments to build insight data")
	if len(c.newlyCreatedDeploymentsHandlers) == 0 {
		c.logger.Info("skip building insight data for newly created deployments because there is no configured handlers")
		return nil
	}

	now := time.Now()
	targetDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	m, err := c.insightstore.LoadMilestone(ctx)
	if err != nil {
		if !errors.Is(err, filestore.ErrNotFound) {
			c.logger.Error("failed to load milestone", zap.Error(err))
			return err
		}
		m = &insight.Milestone{}
	}

	dc, err := c.findDeploymentsCreatedInRange(ctx, m.DeploymentCreatedAtMilestone, targetDate.Unix())
	if err != nil {
		c.logger.Error("failed to find newly created deployment", zap.Error(err))
		return err
	}

	var handleErr error
	for _, handler := range c.newlyCreatedDeploymentsHandlers {
		if err := handler(ctx, dc, targetDate); err != nil {
			c.logger.Error("failed to execute a handler for newly created deployments", zap.Error(err))
			// In order to give all handlers the chance to handle the received data, we do not return here.
			handleErr = err
		}
	}

	if handleErr != nil {
		return handleErr
	}

	m.DeploymentCreatedAtMilestone = targetDate.Unix()
	if err := c.insightstore.PutMilestone(ctx, m); err != nil {
		c.logger.Error("failed to store milestone", zap.Error(err))
		return err
	}

	return nil
}

func (c *InsightCollector) ProcessNewlyCompletedDeployments(ctx context.Context) error {
	c.logger.Info("will retrieve newly completed deployments to build insight data")
	if len(c.newlyCreatedDeploymentsHandlers) == 0 {
		c.logger.Info("skip building insight data for newly completed deployments because there is no configured handlers")
		return nil
	}

	now := time.Now()
	targetDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	m, err := c.insightstore.LoadMilestone(ctx)
	if err != nil {
		if !errors.Is(err, filestore.ErrNotFound) {
			c.logger.Error("failed to load milestone", zap.Error(err))
			return err
		}
		m = &insight.Milestone{}
	}

	dc, err := c.findDeploymentsCompletedInRange(ctx, m.DeploymentCompletedAtMilestone, targetDate.Unix())
	if err != nil {
		c.logger.Error("failed to find newly completed deployment", zap.Error(err))
		return err
	}

	var handleErr error
	for _, handler := range c.newlyCompletedDeploymentsHandlers {
		if err := handler(ctx, dc, targetDate); err != nil {
			c.logger.Error("failed to execute a handler for newly completed deployments", zap.Error(err))
			// In order to give all handlers the chance to handle the received data, we do not return here.
			handleErr = err
		}
	}

	if handleErr != nil {
		return handleErr
	}

	m.DeploymentCompletedAtMilestone = targetDate.Unix()
	if err := c.insightstore.PutMilestone(ctx, m); err != nil {
		c.logger.Error("failed to store milestone", zap.Error(err))
		return err
	}

	return nil
}

func (c *InsightCollector) ProcessApplications(ctx context.Context) error {
	c.logger.Info("will retrieve all applications to build insight data")
	if len(c.newlyCreatedDeploymentsHandlers) == 0 {
		c.logger.Info("skip building insight data for applications because there is no configured handlers")
		return nil
	}

	return errors.New("not implemented yet")
}
