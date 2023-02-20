// Copyright 2023 The PipeCD Authors.
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
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/insight"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type applicationLister interface {
	List(ctx context.Context, opts datastore.ListOptions) ([]*model.Application, string, error)
}

type deploymentLister interface {
	List(ctx context.Context, opts datastore.ListOptions) ([]*model.Deployment, string, error)
}

type Collector struct {
	applicationDataCol         *applicationDataCollector
	completedDeploymentDataCol *completedDeploymentDataCollector
	cfg                        config.ControlPlaneInsightCollector
	logger                     *zap.Logger
}

func NewCollector(ds datastore.DataStore, store insight.Store, cfg config.ControlPlaneInsightCollector, logger *zap.Logger) *Collector {
	logger = logger.Named("insight-collector")

	var (
		w            = datastore.OpsCommander
		appLister    = datastore.NewApplicationStore(ds, w)
		deployLister = datastore.NewDeploymentStore(ds, w)
		appCol       = newApplicationDataCollector(appLister, store, logger)
		comDepCol    = newCompletedDeploymentDataCollector(deployLister, store, logger)
	)

	return &Collector{
		applicationDataCol:         appCol,
		completedDeploymentDataCol: comDepCol,
		cfg:                        cfg,
		logger:                     logger,
	}
}

func (c *Collector) Run(ctx context.Context) error {
	c.logger.Info("start running insight collector", zap.Any("config", c.cfg))
	cr := cron.New(cron.WithLocation(time.UTC))

	if *c.cfg.Application.Enabled {
		_, err := cr.AddFunc(c.cfg.Application.Schedule, func() {
			c.applicationDataCol.Execute(ctx)
		})
		if err != nil {
			c.logger.Error("failed to configure cron job to collect app data", zap.Error(err))
			return err
		}
		c.logger.Info("added a cron job to collect app data")
	}

	if *c.cfg.Deployment.Enabled {
		_, err := cr.AddFunc(c.cfg.Deployment.Schedule, func() {
			c.completedDeploymentDataCol.Execute(ctx)
		})
		if err != nil {
			c.logger.Error("failed to configure cron job to collect deployment data", zap.Error(err))
			return err
		}
		c.logger.Info("added a cron job to collect deployment data")
	}

	cr.Start()
	<-ctx.Done()

	cr.Stop()
	c.logger.Info("insight collector has been stopped")
	return nil
}
