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

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/filestore"
	"github.com/pipe-cd/pipe/pkg/insight"
	"github.com/pipe-cd/pipe/pkg/model"
)

// collectApplicationCount collects application count data.
func (i *InsightCollector) collectApplicationCount(ctx context.Context, apps []*model.Application, target time.Time) error {
	appmap := groupApplicationsByProjectID(apps)
	var updateErr error
	for pid, apps := range appmap {
		if err := i.updateApplicationCount(ctx, apps, pid, target); err != nil {
			updateErr = err
		}
	}
	return updateErr
}

func (i *InsightCollector) updateApplicationCount(ctx context.Context, apps []*model.Application, pid string, target time.Time) error {
	a, err := i.insightstore.LoadApplicationCount(ctx, pid)
	if err != nil {
		if err == filestore.ErrNotFound {
			a = insight.NewApplicationCount()
			oldestApp := findOldestApplication(apps)
			a.AccumulatedFrom = oldestApp.CreatedAt
		} else {
			return fmt.Errorf("load application count: %w", err)
		}
	}

	a.MigrateApplicationCount()

	// update application count
	a.UpdateCount(apps)

	a.AccumulatedTo = target.Unix()

	if err := i.insightstore.PutApplicationCount(ctx, a, pid); err != nil {
		return fmt.Errorf("put application count: %w", err)
	}

	return nil
}

func (i *InsightCollector) getApplications(ctx context.Context, to time.Time) ([]*model.Application, error) {
	var applications []*model.Application
	maxCreatedAt := to.Unix()
	for {
		apps, _, err := i.applicationStore.ListApplications(ctx, datastore.ListOptions{
			Limit: limit,
			Filters: []datastore.ListFilter{
				{
					Field:    "CreatedAt",
					Operator: "<",
					Value:    maxCreatedAt,
				},
			},
			Orders: []datastore.Order{
				{
					Field:     "CreatedAt",
					Direction: datastore.Desc,
				},
			},
		})
		if err != nil {
			i.logger.Error("failed to fetch applications", zap.Error(err))
			return nil, err
		}

		applications = append(applications, apps...)
		if len(apps) < limit {
			break
		}
		maxCreatedAt = apps[len(apps)-1].CreatedAt
	}
	return applications, nil
}

// groupApplicationsByProjectID groups applications by projectID
func groupApplicationsByProjectID(applications []*model.Application) map[string][]*model.Application {
	apps := map[string][]*model.Application{}
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
