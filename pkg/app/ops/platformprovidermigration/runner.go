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
	appCnt, err := r.migrate(ctx)
	if err != nil {
		r.logger.Error("unable to finish application platform provider migration task in first run", zap.Error(err))
	}

	if appCnt == 0 {
		r.logger.Info("application platform provider migration task done successfully")
		return nil
	}

	taskRunTicker := time.NewTicker(migrationRunInterval)
	defer taskRunTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-taskRunTicker.C:
			appCnt, _ := r.migrate(ctx)
			r.logger.Warn(fmt.Sprintf("there are %d application(s) remained from platform provider migration task running", appCnt))
		}
	}
}

// migrate runs the migration task to update all applications in the database.
// It returns the number represents how many application remained/failed to update/migrate.
// We expect the number is 0 (zero) means all applications are migrated. If it returns -1
// means error occurred before the number of actually remained apps is not yet be calculated.
func (r *Runner) migrate(ctx context.Context) (int, error) {
	opts := datastore.ListOptions{
		Orders: []datastore.Order{
			{
				Field:     "UpdatedAt",
				Direction: datastore.Desc,
			},
			{
				Field:     "Id",
				Direction: datastore.Asc,
			},
		},
	}

	apps, _, err := r.applicationStore.List(ctx, opts)
	if err != nil {
		r.logger.Error("failed to fetch all applications to run migrate task", zap.Error(err))
		return -1, err
	}

	appCnt := len(apps)
	for _, app := range apps {
		if app.PlatformProvider == "" {
			//lint:ignore SA1019 app.CloudProvider is deprecated.
			if err = r.applicationStore.UpdatePlatformProvider(ctx, app.Id, app.CloudProvider); err != nil {
				r.logger.Error("failed to update application platform provider value",
					zap.String("id", app.Id),
					//lint:ignore SA1019 app.CloudProvider is deprecated.
					zap.String("provider", app.CloudProvider),
					zap.Error(err),
				)
				continue
			}
		}
		appCnt--
	}

	return appCnt, nil
}
