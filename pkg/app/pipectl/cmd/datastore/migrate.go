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

package datastore

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pipe-cd/pipe/pkg/cli"
	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/datastore/migration"
	"github.com/pipe-cd/pipe/pkg/datastore/mongodb"
	"github.com/pipe-cd/pipe/pkg/datastore/mysql"
	"github.com/pipe-cd/pipe/pkg/datastore/mysql/ensurer"
)

type migrate struct {
	root *command

	upstreamDataSrc   string
	downstreamDataSrc string
	database          string
	models            []string
	stdout            io.Writer
}

func newMigrateCommand(root *command) *cobra.Command {
	m := &migrate{
		root:   root,
		stdout: os.Stdout,
	}
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate data to MySQL datastore.",
		Long:  "Migrate data from MongoDB to MySQL datastore.\nBoth upstream (MongoDB) and downstream (MySQL) datastore have to available to use this command.",
		RunE:  cli.WithContext(m.run),
	}

	cmd.Flags().StringVar(&m.upstreamDataSrc, "upstream-data-src", m.upstreamDataSrc, "The URL to connect to upstream datastore (MongoDB).\n Format: mongodb://username:password@hostname:27017/database")
	cmd.Flags().StringVar(&m.downstreamDataSrc, "downstream-data-src", m.downstreamDataSrc, "The URL to connect to downstream datastore (MySQL).\n Format: username:password@tcp(hostname:3306)")
	cmd.Flags().StringVar(&m.database, "database", m.database, "The SQL database name.")
	cmd.Flags().StringSliceVar(&m.models, "models", m.models, fmt.Sprintf("The list of migrating models. If nothing is passed, all models will be migrated.\n (%s)", strings.Join(datastore.MigratableModelKinds, " | ")))

	cmd.MarkFlagRequired("upstream-data-src")
	cmd.MarkFlagRequired("downstream-data-src")
	cmd.MarkFlagRequired("database")

	return cmd
}

func (m *migrate) run(ctx context.Context, t cli.Telemetry) error {
	ensurerExec, err := ensurer.NewMySQLEnsurer(m.downstreamDataSrc, m.database, "", "", t.Logger)
	if err != nil {
		return fmt.Errorf("failed to create SQL ensurer instance: %w", err)
	}
	defer func() {
		if err := ensurerExec.Close(); err != nil {
			fmt.Fprintf(m.stdout, "error occurred while close the ensurer database connection: %v", err)
		}
	}()
	// Ensure SQL schema on the new datastore.
	if err = ensurerExec.EnsureSchema(ctx); err != nil {
		return fmt.Errorf("failed to prepare schema on the new datastore: %w", err)
	}

	mongodbDatastore, err := mongodb.NewMongoDB(ctx, m.upstreamDataSrc, m.database, mongodb.WithLogger(t.Logger))
	if err != nil {
		return fmt.Errorf("failed to connect to upstream datastore: %w", err)
	}

	mysqlDatastore, err := mysql.NewMySQL(m.downstreamDataSrc, m.database, mysql.WithLogger(t.Logger))
	if err != nil {
		return fmt.Errorf("failed to connect to downstream datastore: %w", err)
	}

	modelsNameList, err := makeMigrateModelList(m.models)
	if err != nil {
		return fmt.Errorf("failed to migrate datastore: %w", err)
	}

	dataTransfer := migration.NewDataTransfer(mongodbDatastore, mysqlDatastore)
	if err = dataTransfer.TransferMulti(ctx, modelsNameList); err != nil {
		return fmt.Errorf("failed to migrate datastore: %w", err)
	}

	fmt.Fprintln(m.stdout, "the migration process was completed succesfully.")
	return nil
}

func makeMigrateModelList(modelsName []string) ([]string, error) {
	if len(modelsName) == 0 {
		return datastore.MigratableModelKinds, nil
	}

	nameMap := make(map[string]struct{}, len(datastore.MigratableModelKinds))
	for _, name := range datastore.MigratableModelKinds {
		nameMap[name] = struct{}{}
	}
	// validate
	for _, name := range modelsName {
		if _, ok := nameMap[name]; !ok {
			return nil, fmt.Errorf("invalid model name passed: %s", name)
		}
	}

	return modelsName, nil
}
