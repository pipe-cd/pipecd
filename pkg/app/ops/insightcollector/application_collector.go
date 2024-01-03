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

package insightcollector

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/insight"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type applicationDataCollector struct {
	lister applicationLister
	store  insight.ApplicationStore
	logger *zap.Logger
}

func newApplicationDataCollector(lister applicationLister, store insight.ApplicationStore, logger *zap.Logger) *applicationDataCollector {
	return &applicationDataCollector{
		lister: lister,
		store:  store,
		logger: logger,
	}
}

func (c *applicationDataCollector) Execute(ctx context.Context) {
	now := time.Now()
	c.logger.Info("will retrieve all applications to build insight data")

	apps, err := c.listApplications(ctx)
	if err != nil {
		c.logger.Error("failed to list applications", zap.Error(err))
		return
	}
	c.logger.Info(fmt.Sprintf("fetched %d applications to build insight data", len(apps)),
		zap.Duration("duration", time.Since(now)),
	)

	appsByProject := make(map[string][]*insight.ApplicationData)
	for _, a := range apps {
		var (
			project = a.ProjectId
			app     = insight.BuildApplicationData(a)
		)
		appsByProject[project] = append(appsByProject[project], &app)
	}

	var hasError bool
	for project, apps := range appsByProject {
		data := insight.BuildProjectApplicationData(apps, now.Unix())
		if err := c.store.PutApplications(ctx, project, &data); err != nil {
			c.logger.Error("failed to store application data",
				zap.String("project", project),
				zap.Error(err),
			)
			hasError = true
		}
	}

	if !hasError {
		c.logger.Info(fmt.Sprintf("successfully built and stored application insight data for %d projects", len(appsByProject)),
			zap.Duration("duration", time.Since(now)),
		)
	}
}

func (c *applicationDataCollector) listApplications(ctx context.Context) ([]*model.Application, error) {
	const limit = 100
	var cursor string
	var applications []*model.Application

	for {
		apps, next, err := c.lister.List(ctx, datastore.ListOptions{
			Filters: []datastore.ListFilter{
				{
					Field:    "Deleted",
					Operator: datastore.OperatorEqual,
					Value:    false,
				},
			},
			Orders: []datastore.Order{
				{
					Field:     "CreatedAt",
					Direction: datastore.Asc,
				},
				{
					Field:     "Id",
					Direction: datastore.Asc,
				},
			},
			Cursor: cursor,
			Limit:  limit,
		})
		if err != nil {
			return nil, err
		}

		applications = append(applications, apps...)
		if next == "" {
			break
		}
		cursor = next
	}

	return applications, nil
}
