// Copyright 2023 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package deployment

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type list struct {
	root *command

	statuses []string
	appKinds []string
	appIds   []string
	appName  string
	labels   map[string]string
	limit    int32

	cursor string
	stdout io.Writer
}

func newListCommand(root *command) *cobra.Command {
	c := &list{
		root:   root,
		stdout: os.Stdout,
	}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Show the list of deployments.",
		RunE:  cli.WithContext(c.run),
	}

	// TODO: Support pipectl to list application by Label.
	cmd.Flags().StringSliceVar(&c.statuses, "status", c.statuses, fmt.Sprintf("The list of waiting statuses. (%s)", strings.Join(model.DeploymentStatusStrings(), "|")))
	cmd.Flags().StringSliceVar(&c.appIds, "app-id", c.appIds, fmt.Sprintf("The application id. (%s)", strings.Join(model.ApplicationKindStrings(), "|")))
	cmd.Flags().StringSliceVar(&c.appKinds, "app-kind", c.appKinds, fmt.Sprintf("The kind of application. (%s)", strings.Join(model.ApplicationKindStrings(), "|")))
	cmd.Flags().StringVar(&c.appName, "app-name", c.appName, "The application name.")
	cmd.Flags().StringVar(&c.cursor, "cursor", c.cursor, "The cursor which returned by the previous request applications list.")
	cmd.Flags().Int32Var(&c.limit, "limit", c.limit, "")

	return cmd
}

func (c *list) run(ctx context.Context, _ cli.Input) error {
	for _, status := range c.statuses {
		if status != "" {
			if _, ok := model.DeploymentStatus_value[status]; !ok {
				return errors.Errorf("%s is invalid deployment status", status)
			}
		}
	}

	for _, kind := range c.appKinds {
		if kind != "" {
			if _, ok := model.ApplicationKind_value[kind]; !ok {
				return errors.Errorf("%s is invalid application kind", kind)
			}
		}
	}

	cli, err := c.root.clientOptions.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}
	defer cli.Close()

	req := &apiservice.ListDeploymentsRequest{
		Statuses:        c.statuses,
		Kinds:           c.appKinds,
		ApplicationIds:  c.appIds,
		ApplicationName: c.appName,
		Limit:           c.limit,
		Cursor:          c.cursor,
	}

	resp, err := cli.ListDeployments(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to list application: %w", err)
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal applications: %w", err)
	}

	fmt.Fprintln(c.stdout, string(bytes))
	return nil
}
