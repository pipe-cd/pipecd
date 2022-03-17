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
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/filestore"
	"github.com/pipe-cd/pipecd/pkg/insight/insightstore"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type applicationStore interface {
	List(ctx context.Context, opts datastore.ListOptions) ([]*model.Application, string, error)
}

type deploymentStore interface {
	List(ctx context.Context, opts datastore.ListOptions) ([]*model.Deployment, string, error)
}

type Collector struct {
	applicationStore      applicationStore
	deploymentStore       deploymentStore
	applicationCountStore insightstore.ApplicationCountStore

	applicationsHandlers              []func(ctx context.Context, applications []*model.Application, target time.Time) error
	newlyCreatedDeploymentsHandlers   []func(ctx context.Context, developments []*model.Deployment, target time.Time) error
	newlyCompletedDeploymentsHandlers []func(ctx context.Context, developments []*model.Deployment, target time.Time) error

	config config.ControlPlaneInsightCollector
	logger *zap.Logger
}

func NewCollector(ds datastore.DataStore, fs filestore.Store, cfg config.ControlPlaneInsightCollector, logger *zap.Logger) *Collector {
	w := datastore.OpsCommander
	c := &Collector{
		applicationStore:      datastore.NewApplicationStore(ds, w),
		deploymentStore:       datastore.NewDeploymentStore(ds, w),
		applicationCountStore: insightstore.NewStore(fs),
		config:                cfg,
		logger:                logger.Named("insight-collector"),
	}

	if cfg.Application.Enabled {
		c.applicationsHandlers = append(c.applicationsHandlers, c.collectApplicationCount)
	}
	// if cfg.Deployment.Enabled {
	// 	c.newlyCreatedDeploymentsHandlers = append(c.newlyCreatedDeploymentsHandlers, c.collectDevelopmentFrequency)
	// 	c.newlyCompletedDeploymentsHandlers = append(c.newlyCompletedDeploymentsHandlers, c.collectDeploymentChangeFailureRate)
	// }

	return c
}

func (c *Collector) Run(ctx context.Context) error {
	c.logger.Info("start running insight collector", zap.Any("config", c.config))
	cr := cron.New(cron.WithLocation(time.UTC))

	if c.config.Application.Enabled {
		_, err := cr.AddFunc(c.config.Application.Schedule, func() {
			c.collectApplicationMetrics(ctx)
		})
		if err != nil {
			c.logger.Error("failed to configure cron job for collecting application metrics", zap.Error(err))
			return err
		}
		c.logger.Info("added a cron job for collecting application metrics")
	}

	// if c.config.Deployment.Enabled {
	// 	_, err := cr.AddFunc(c.config.Deployment.Schedule, func() {
	// 		c.collectDeploymentMetrics(ctx)
	// 	})
	// 	if err != nil {
	// 		c.logger.Error("failed to configure cron job for collecting deployment metrics", zap.Error(err))
	// 		return err
	// 	}
	// 	c.logger.Info("added a cron job for collecting deployment metrics")
	// }

	cr.Start()
	<-ctx.Done()
	cr.Stop()
	c.logger.Info("insight collector has been stopped")
	return nil
}

func (c *Collector) collectApplicationMetrics(ctx context.Context) {
	if !c.config.Application.Enabled {
		c.logger.Info("do not collecting application metrics because it was not enabled")
		return
	}

	start := time.Now()
	c.logger.Info("will retrieve all applications to build insight data")

	apps, err := c.listApplications(ctx, start)
	if err != nil {
		c.logger.Error("failed to list applications", zap.Error(err))
		return
	}
	c.logger.Info(fmt.Sprintf("there are %d applications to build insight data", len(apps)),
		zap.Duration("duration", time.Since(start)),
	)

	for _, handler := range c.applicationsHandlers {
		if err := handler(ctx, apps, start); err != nil {
			c.logger.Error("failed to execute a handler for applications", zap.Error(err))
			// In order to give all handlers the chance to handle the received data, we do not return here.
		}
	}
	return
}

// func (c *Collector) collectDeploymentMetrics(ctx context.Context) {
// 	if !c.config.Deployment.Enabled {
// 		c.logger.Info("do not collecting deployment metrics because it was not enabled")
// 		return
// 	}

// 	cfg := c.config.Deployment
// 	retry := backoff.NewRetry(cfg.Retries, backoff.NewConstant(cfg.RetryInterval.Duration()))

// 	var doneNewlyCompleted, doneNewlyCreated bool

// 	for retry.WaitNext(ctx) {
// 		if !doneNewlyCompleted {
// 			start := time.Now()
// 			if err := c.processNewlyCompletedDeployments(ctx); err != nil {
// 				c.logger.Error("failed to process the newly completed deployments", zap.Error(err))
// 			} else {
// 				c.logger.Info("successfully processed the newly completed deployments",
// 					zap.Duration("duration", time.Since(start)),
// 				)
// 				doneNewlyCompleted = true
// 			}
// 		}

// 		if !doneNewlyCreated {
// 			start := time.Now()
// 			if err := c.processNewlyCreatedDeployments(ctx); err != nil {
// 				c.logger.Error("failed to process the newly created deployments", zap.Error(err))
// 			} else {
// 				c.logger.Info("successfully processed the newly created deployments",
// 					zap.Duration("duration", time.Since(start)),
// 				)
// 				doneNewlyCreated = true
// 			}
// 		}

// 		if doneNewlyCompleted && doneNewlyCreated {
// 			return
// 		}
// 		c.logger.Info("will do another try to collect insight data")
// 	}
// }

// func (c *Collector) processNewlyCreatedDeployments(ctx context.Context) error {
// 	c.logger.Info("will retrieve newly created deployments to build insight data")
// 	now := time.Now()
// 	targetDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

// 	m, err := c.insightstore.LoadMilestone(ctx)
// 	if err != nil {
// 		if !errors.Is(err, filestore.ErrNotFound) {
// 			c.logger.Error("failed to load milestone", zap.Error(err))
// 			return err
// 		}
// 		m = &insight.Milestone{}
// 	}

// 	dc, err := c.findDeploymentsCreatedInRange(ctx, m.DeploymentCreatedAtMilestone, targetDate.Unix())
// 	if err != nil {
// 		c.logger.Error("failed to find newly created deployment", zap.Error(err))
// 		return err
// 	}

// 	var handleErr error
// 	for _, handler := range c.newlyCreatedDeploymentsHandlers {
// 		if err := handler(ctx, dc, targetDate); err != nil {
// 			c.logger.Error("failed to execute a handler for newly created deployments", zap.Error(err))
// 			// In order to give all handlers the chance to handle the received data, we do not return here.
// 			handleErr = err
// 		}
// 	}

// 	if handleErr != nil {
// 		return handleErr
// 	}

// 	m.DeploymentCreatedAtMilestone = targetDate.Unix()
// 	if err := c.insightstore.PutMilestone(ctx, m); err != nil {
// 		c.logger.Error("failed to store milestone", zap.Error(err))
// 		return err
// 	}

// 	return nil
// }

// func (c *Collector) processNewlyCompletedDeployments(ctx context.Context) error {
// 	c.logger.Info("will retrieve newly completed deployments to build insight data")
// 	now := time.Now()
// 	targetDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

// 	m, err := c.insightstore.LoadMilestone(ctx)
// 	if err != nil {
// 		if !errors.Is(err, filestore.ErrNotFound) {
// 			c.logger.Error("failed to load milestone", zap.Error(err))
// 			return err
// 		}
// 		m = &insight.Milestone{}
// 	}

// 	dc, err := c.findDeploymentsCompletedInRange(ctx, m.DeploymentCompletedAtMilestone, targetDate.Unix())
// 	if err != nil {
// 		c.logger.Error("failed to find newly completed deployment", zap.Error(err))
// 		return err
// 	}

// 	var handleErr error
// 	for _, handler := range c.newlyCompletedDeploymentsHandlers {
// 		if err := handler(ctx, dc, targetDate); err != nil {
// 			c.logger.Error("failed to execute a handler for newly completed deployments", zap.Error(err))
// 			// In order to give all handlers the chance to handle the received data, we do not return here.
// 			handleErr = err
// 		}
// 	}

// 	if handleErr != nil {
// 		return handleErr
// 	}

// 	m.DeploymentCompletedAtMilestone = targetDate.Unix()
// 	if err := c.insightstore.PutMilestone(ctx, m); err != nil {
// 		c.logger.Error("failed to store milestone", zap.Error(err))
// 		return err
// 	}

// 	return nil
// }
