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

package application

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type list struct {
	root *command

	appName  string
	appKind  string
	disabled bool
	cursor   string
	stdout   io.Writer
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

	// TODO: Support pipectl to list application by Label.
	cmd.Flags().StringVar(&c.appName, "app-name", c.appName, "The application name.")
	cmd.Flags().StringVar(&c.appKind, "app-kind", c.appKind, fmt.Sprintf("The kind of application. (%s)", strings.Join(model.ApplicationKindStrings(), "|")))
	cmd.Flags().BoolVar(&c.disabled, "disabled", c.disabled, "True to show only disabled applications.")
	cmd.Flags().StringVar(&c.cursor, "cursor", c.cursor, "The cursor which returned by the previous request applications list.")

	return cmd
}

func (c *list) run(ctx context.Context, _ cli.Input) error {
	if c.appKind != "" {
		if _, ok := model.ApplicationKind_value[c.appKind]; !ok {
			return fmt.Errorf("invalid application kind")
		}
	}

	cli, err := c.root.clientOptions.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}
	defer cli.Close()

	req := &apiservice.ListApplicationsRequest{
		Name:     c.appName,
		Kind:     c.appKind,
		Disabled: c.disabled,
		Cursor:   c.cursor,
	}

	resp, err := cli.ListApplications(ctx, req)
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
