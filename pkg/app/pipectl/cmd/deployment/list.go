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
	labels   []string
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

	cmd.Flags().StringSliceVar(&c.statuses, "status", c.statuses, fmt.Sprintf("The list of application statuses to filter. (%s)", strings.Join(model.DeploymentStatusStrings(), "|")))
	cmd.Flags().StringSliceVar(&c.appIds, "app-id", c.appIds, fmt.Sprintf("The list of application ids to filter. (%s)", strings.Join(model.ApplicationKindStrings(), "|")))
	cmd.Flags().StringSliceVar(&c.appKinds, "app-kind", c.appKinds, fmt.Sprintf("The list of application kinds to filter. (%s)", strings.Join(model.ApplicationKindStrings(), "|")))
	cmd.Flags().StringVar(&c.appName, "app-name", c.appName, "The application name to filter.")
	cmd.Flags().StringVar(&c.cursor, "cursor", c.cursor, "The cursor which returned by the previous request applications list.")
	cmd.Flags().Int32Var(&c.limit, "limit", 30, "Upper limit on the number of return values. Default value is 30.")
	cmd.Flags().StringSliceVar(&c.labels, "label", c.labels, "The list of labels to filter. Expect input in the form KEY:VALUE.")

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

	labels := map[string]string{}
	for _, label := range c.labels {
		sp := strings.SplitN(label, ":", 2)
		if len(sp) == 2 {
			labels[sp[0]] = sp[1]
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
		Labels:          labels,
	}

	resp, err := cli.ListDeployments(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to list deployment: %w", err)
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal deployments: %w", err)
	}

	fmt.Fprintln(c.stdout, string(bytes))
	return nil
}
