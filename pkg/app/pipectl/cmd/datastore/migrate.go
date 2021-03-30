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
	"strings"

	"github.com/spf13/cobra"

	"github.com/pipe-cd/pipe/pkg/app/api/service/apiservice"
	"github.com/pipe-cd/pipe/pkg/cli"
)

var availableModelsStringName = []string{
	"Project",
	"Application",
	"Command",
	"Deployment",
	"Environment",
	"Piped",
	"APIKey",
	"Event",
}

type migrate struct {
	root *command

	downstreamDataSrc string
	models            []string
}

func newMigrateCommand(root *command) *cobra.Command {
	m := &migrate{
		root: root,
	}
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate data to MySQL datastore.",
		Long:  "Migrate data to MySQL datastore.\nBoth upstream (MongoDB/Firestore) and downstream (MySQL) datastore have to available to connect from control-plane to use this command.",
		RunE:  cli.WithContext(m.run),
	}

	cmd.Flags().StringVar(&m.downstreamDataSrc, "downstream-data-src", m.downstreamDataSrc, "The URL to connect to downstream datastore (MySQL).\n Format: username:password@tcp(hostname:3306)/database")
	cmd.Flags().StringSliceVar(&m.models, "models", m.models, fmt.Sprintf("The list of migrating models. If nothing is passed, all models will be migrated.\n (%s)", strings.Join(availableModelsStringName, " | ")))

	cmd.MarkFlagRequired("downstream-data-src")

	return cmd
}

func (m *migrate) run(ctx context.Context, _ cli.Telemetry) error {
	cli, err := m.root.clientOptions.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}
	defer cli.Close()

	modelsNameList, err := makeMigrateModelsList(m.models)
	if err != nil {
		return fmt.Errorf("failed to migrate datastore: %w", err)
	}

	req := &apiservice.MigrateDatastoreRequest{
		DownstreamDataSrc: m.downstreamDataSrc,
		Models:            modelsNameList,
	}

	if _, err := cli.MigrateDatastore(ctx, req); err != nil {
		return fmt.Errorf("failed to migrate datastore: %w", err)
	}
	return nil
}

func makeMigrateModelsList(modelsName []string) ([]string, error) {
	if len(modelsName) == 0 {
		return availableModelsStringName, nil
	}

	nameMap := make(map[string]interface{})
	for _, name := range availableModelsStringName {
		nameMap[name] = nil
	}
	// validate
	for _, name := range modelsName {
		if _, ok := nameMap[name]; !ok {
			return nil, fmt.Errorf("invalid model name passed: %s", name)
		}
	}

	return modelsName, nil
}
