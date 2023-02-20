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
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/filestore"
	"github.com/pipe-cd/pipecd/pkg/insight"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type completedDeploymentDataCollector struct {
	lister deploymentLister
	store  insight.CompletedDeploymentStore
	logger *zap.Logger
}

func newCompletedDeploymentDataCollector(lister deploymentLister, store insight.CompletedDeploymentStore, logger *zap.Logger) *completedDeploymentDataCollector {
	return &completedDeploymentDataCollector{
		lister: lister,
		store:  store,
		logger: logger,
	}
}

func (c *completedDeploymentDataCollector) Execute(ctx context.Context) {
	c.logger.Info("will retrieve newly completed deployments to build insight data")
	now := time.Now()

	m, err := c.store.GetMilestone(ctx)
	if err != nil {
		if !errors.Is(err, filestore.ErrNotFound) {
			c.logger.Error("failed to load milestone from storage", zap.Error(err))
			return
		}

		// There was no milestone in the storage.
		// It means this is the first time collector was run.
		// In that case, we collect data from 1 hour ago.
		m = &insight.Milestone{
			DeploymentCompletedAtMilestone: now.Add(-time.Hour).Unix(),
		}
	}

	var (
		dataRangeFrom = m.DeploymentCompletedAtMilestone
		dataRangeTo   = now.Unix()
	)

	const collectRangeLimit = 2 * 24 * 60 * 60 // 2 days
	if dataRangeTo-dataRangeFrom > collectRangeLimit {
		dataRangeFrom = dataRangeTo - collectRangeLimit
		c.logger.Warn(fmt.Sprintf("it seems collector had not been running for a long time, data between range [%d, %d] will be ignored",
			m.DeploymentCompletedAtMilestone,
			dataRangeFrom,
		))
	}

	ds, err := c.listCompletedDeployments(ctx, dataRangeFrom, dataRangeTo)
	if err != nil {
		c.logger.Error("failed to find newly completed deployment",
			zap.Duration("duration", time.Since(now)),
			zap.Error(err),
		)
		return
	}
	c.logger.Info(fmt.Sprintf("there are %d completed deployments to build insight data", len(ds)),
		zap.Duration("duration", time.Since(now)),
	)

	deploysByProject := make(map[string][]*insight.DeploymentData)
	for _, d := range ds {
		var (
			project = d.ProjectId
			data    = insight.BuildDeploymentData(d)
		)
		deploysByProject[project] = append(deploysByProject[project], &data)
	}

	var hasError bool
	for project, ds := range deploysByProject {
		if err := c.store.PutCompletedDeployments(ctx, project, ds); err != nil {
			c.logger.Error("failed to store deployment data",
				zap.String("project", project),
				zap.Error(err),
			)
			hasError = true
		}
	}

	if hasError {
		return
	}

	c.logger.Info(fmt.Sprintf("successfully built and stored deployment insight data for %d projects", len(deploysByProject)),
		zap.Duration("duration", time.Since(now)),
	)

	m.DeploymentCompletedAtMilestone = dataRangeTo
	if err := c.store.PutMilestone(ctx, m); err != nil {
		c.logger.Error("failed to store milestone", zap.Error(err))
		return
	}

	c.logger.Info("successfully stored a new milestone", zap.Int64("milestone", m.DeploymentCompletedAtMilestone))
}

func (c *completedDeploymentDataCollector) listCompletedDeployments(ctx context.Context, from, to int64) ([]*model.Deployment, error) {
	const callLimit = 50

	var (
		filters = []datastore.ListFilter{
			{
				Field:    "CompletedAt",
				Operator: datastore.OperatorLessThanOrEqual,
				Value:    to,
			},
			{
				Field:    "CompletedAt",
				Operator: datastore.OperatorGreaterThanOrEqual,
				Value:    from,
			},
		}
		orders = []datastore.Order{
			{
				Field:     "CompletedAt",
				Direction: datastore.Asc,
			},
			{
				Field:     "Id",
				Direction: datastore.Asc,
			},
		}
		deployments []*model.Deployment
		cursor      string
	)

	for {
		d, next, err := c.lister.List(ctx, datastore.ListOptions{
			Limit:   callLimit,
			Cursor:  cursor,
			Filters: filters,
			Orders:  orders,
		})
		if err != nil {
			return nil, err
		}

		deployments = append(deployments, d...)
		if next == "" {
			break
		}

		cursor = next
	}

	return deployments, nil
}
