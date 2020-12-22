// Copyright 2020 The PipeCD Authors.
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

package application

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pipe-cd/pipe/pkg/app/api/service/apiservice"
	"github.com/pipe-cd/pipe/pkg/cli"
	"github.com/pipe-cd/pipe/pkg/model"
)

type list struct {
	root *command

	appName string
	envId   string
	appKind string
	stdout  io.Writer
}

func newListCommand(root *command) *cobra.Command {
	c := &list{
		root:   root,
		stdout: os.Stdout,
	}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Show the list of applications. Currently, the maximum number of returned applications is 10.",
		RunE:  cli.WithContext(c.run),
	}

	cmd.Flags().StringVar(&c.appName, "app-name", c.appName, "The application name.")
	cmd.Flags().StringVar(&c.envId, "env-id", c.envId, "The environment ID.")
	cmd.Flags().StringVar(&c.appKind, "app-kind", c.appKind, fmt.Sprintf("The kind of application. (%s)", strings.Join(model.ApplicationKindStrings(), "|")))

	return cmd
}

func (c *list) run(ctx context.Context, _ cli.Telemetry) error {
	if c.appKind != "" {
		if _, ok := model.ApplicationKind_value[c.appKind]; !ok {
			return fmt.Errorf("invalid applicaiton kind")
		}
	}

	cli, err := c.root.clientOptions.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}
	defer cli.Close()

	req := &apiservice.ListApplicationsRequest{
		Name:  c.appName,
		EnvId: c.envId,
		Kind:  c.appKind,
	}

	resp, err := cli.ListApplications(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to list application: %w", err)
	}

	bytes, err := json.Marshal(resp.Applications)
	if err != nil {
		return fmt.Errorf("failed to marshal applications: %w", err)
	}

	fmt.Fprintln(c.stdout, string(bytes))
	return nil
}
