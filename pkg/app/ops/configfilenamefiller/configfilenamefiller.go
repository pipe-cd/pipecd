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

package configfilenamefiller

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/datastore"
)

type Filler struct {
	applicationStore datastore.ApplicationStore
	logger           *zap.Logger
}

func NewFiller(ds datastore.DataStore, logger *zap.Logger) *Filler {
	return &Filler{
		applicationStore: datastore.NewApplicationStore(ds),
		logger:           logger.Named("fill-config-filename"),
	}
}

func (c *Filler) Run(ctx context.Context) error {
	c.logger.Info("start running Filler")
	defer c.logger.Info("Filler has been stopped")

	var limit, scans, fills = 50, 0, 0
	var sleepInternal = 100 * time.Millisecond
	var cursor string

	for {
		apps, next, err := c.applicationStore.ListApplications(ctx, datastore.ListOptions{
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
			c.logger.Error("failed to list apps", zap.Error(err))
			return err
		}

		c.logger.Info(fmt.Sprintf("found %d apps to fill", len(apps)))
		for _, app := range apps {
			scans++
			if app.MostRecentlySuccessfulDeployment == nil || app.MostRecentlySuccessfulDeployment.ConfigFilename != "" {
				c.logger.Info(fmt.Sprintf("there is no need to fill config filename for application %s", app.Id))
				continue
			}

			c.logger.Info(fmt.Sprintf("will fill config filename for application %s", app.Id))
			if err := c.applicationStore.FillConfigFilenameToDeploymentReference(ctx, app.Id); err != nil {
				c.logger.Error(fmt.Sprintf("failed to fill config filename for application %s", app.Id), zap.Error(err))
				return err
			}

			fills++
			c.logger.Info(fmt.Sprintf("filled config filename for application %s", app.Id))
			time.Sleep(sleepInternal)
		}

		if next == "" {
			c.logger.Info(fmt.Sprintf("successfully scanned %d apps and filled config filename for %d apps", scans, fills))
			return nil
		}
		cursor = next
	}
}
