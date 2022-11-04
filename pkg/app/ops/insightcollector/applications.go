// Copyright 2022 The PipeCD Authors.
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

// collectApplicationCount collects application count data.
func (c *Collector) collectApplicationCount(ctx context.Context, apps []*model.Application, target time.Time) error {
	var lastErr error
	appmap := groupApplicationsByProjectID(apps)

	for pid, apps := range appmap {
		if err := c.updateApplicationCounts(ctx, pid, apps, target); err != nil {
			c.logger.Error("failed to update ApplicationCounts data",
				zap.String("project", pid),
				zap.Error(err),
			)
			lastErr = err
		}
	}
	return lastErr
}

func (c *Collector) updateApplicationCounts(ctx context.Context, projectID string, apps []*model.Application, target time.Time) error {
	counts := insight.MakeApplicationCounts(apps, target)

	if err := c.insightstore.PutApplicationCounts(ctx, projectID, counts); err != nil {
		return fmt.Errorf("failed to put application counts: %w", err)
	}

	return nil
}

func (c *Collector) listApplications(ctx context.Context, to time.Time) ([]*model.Application, error) {
	const limit = 100
	var cursor string
	var applications []*model.Application

	for {
		apps, next, err := c.applicationStore.List(ctx, datastore.ListOptions{
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

// groupApplicationsByProjectID groups applications by projectID
func groupApplicationsByProjectID(applications []*model.Application) map[string][]*model.Application {
	apps := make(map[string][]*model.Application)
	for _, a := range applications {
		apps[a.ProjectId] = append(apps[a.ProjectId], a)
	}
	return apps
}

func findOldestApplication(apps []*model.Application) *model.Application {
	oldestApp := apps[0]
	for _, a := range apps {
		if a.CreatedAt < oldestApp.CreatedAt {
			oldestApp = a
		}
	}
	return oldestApp
}
