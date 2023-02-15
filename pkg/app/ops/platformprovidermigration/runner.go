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

package platformprovidermigration

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	migrationRunInterval = 30 * time.Minute
)

type applicationStore interface {
	List(ctx context.Context, opts datastore.ListOptions) ([]*model.Application, string, error)
	UpdatePlatformProvider(ctx context.Context, id string, provider string) error
}

type Runner struct {
	applicationStore applicationStore
	logger           *zap.Logger
}

func NewRunner(ds datastore.DataStore, logger *zap.Logger) *Runner {
	w := datastore.OpsCommander
	return &Runner{
		applicationStore: datastore.NewApplicationStore(ds, w),
		logger:           logger.Named("platform-provider-migrate-runner"),
	}
}

func (r *Runner) Migrate(ctx context.Context) error {
	r.logger.Info("start running application migration task")

	// Run migration task once on start this ops migration.
	cursor, err := r.migrate(ctx, "")
	if err == nil {
		r.logger.Info("application migration task finished successfully")
		return nil
	}

	r.logger.Error("unable to finish application platform provider migration task in first run", zap.Error(err))

	taskRunTicker := time.NewTicker(migrationRunInterval)
	defer taskRunTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-taskRunTicker.C:
			cursor, err = r.migrate(ctx, cursor)
			if err == nil {
				r.logger.Info("application migration task finished successfully")
				return nil
			}
		}
	}
}

// migrate runs the migration task to update all applications in the database.
// In case of error occurred, it returns error and a cursor string which contains
// the information so that next time we can pass that value to keep migrating from
// failed application, not from start.
func (r *Runner) migrate(ctx context.Context, cursor string) (string, error) {
	const limit = 100

	for {
		apps, nextCur, err := r.applicationStore.List(ctx, datastore.ListOptions{
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
			Limit:  limit,
			Cursor: cursor,
		})
		if err != nil {
			r.logger.Error("failed to fetch applications to run migrate task", zap.Error(err))
			return cursor, err
		}

		if len(apps) == 0 {
			return "", nil
		}

		r.logger.Info(fmt.Sprintf("migrate platform provider value for %d application(s)", len(apps)))

		for _, app := range apps {
			if app.PlatformProvider != "" {
				continue
			}

			//lint:ignore SA1019 app.CloudProvider is deprecated.
			if err = r.applicationStore.UpdatePlatformProvider(ctx, app.Id, app.CloudProvider); err != nil {
				r.logger.Error("failed to update application platform provider value",
					zap.String("id", app.Id),
					//lint:ignore SA1019 app.CloudProvider is deprecated.
					zap.String("provider", app.CloudProvider),
					zap.Error(err),
				)
				return cursor, err
			}
		}

		cursor = nextCur
	}
}
