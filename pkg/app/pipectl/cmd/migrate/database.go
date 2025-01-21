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

package migrate

import (
	"context"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/cli"
)

type database struct {
	root *command

	applications []string
}

func newDatabaseCommand(root *command) *cobra.Command {
	c := &database{
		root: root,
	}
	cmd := &cobra.Command{
		Use:   "database",
		Short: "Do migration tasks for database.",
		RunE:  cli.WithContext(c.run),
	}

	cmd.Flags().StringSliceVar(&c.applications, "applications", c.applications, "The list of application to migrate database.")
	cmd.MarkFlagRequired("applications")
	return cmd
}

func (c *database) run(ctx context.Context, input cli.Input) error {
	client, err := c.root.clientOptions.NewClient(ctx)
	if err != nil {
		input.Logger.Error("failed to create client", zap.Error(err))
		return err
	}

	for _, appID := range c.applications {
		input.Logger.Info("migrating database", zap.String("application", appID))
		if err := c.migrateApplication(ctx, client, appID); err != nil {
			input.Logger.Error("failed to migrate database", zap.String("application", appID), zap.Error(err))
			return err
		}
		input.Logger.Info("successfully migrated database", zap.String("application", appID))
	}

	return nil
}

func (c *database) migrateApplication(ctx context.Context, client apiservice.Client, appID string) error {
	req := &apiservice.MigrateDatabaseRequest{
		Target: &apiservice.MigrateDatabaseRequest_Application_{
			Application: &apiservice.MigrateDatabaseRequest_Application{
				ApplicationId: appID,
			},
		},
	}
	if _, err := client.MigrateDatabase(ctx, req); err != nil {
		return err
	}
	return nil
}
