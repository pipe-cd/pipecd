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

package application

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/cli"
)

type get struct {
	root *command

	appID  string
	stdout io.Writer
}

func newGetCommand(root *command) *cobra.Command {
	c := &get{
		root:   root,
		stdout: os.Stdout,
	}
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Show the information about the specified application.",
		RunE:  cli.WithContext(c.run),
	}

	cmd.Flags().StringVar(&c.appID, "app-id", c.appID, "The application ID.")
	cmd.MarkFlagRequired("app-id")

	return cmd
}

func (c *get) run(ctx context.Context, _ cli.Input) error {
	cli, err := c.root.clientOptions.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}
	defer cli.Close()

	req := &apiservice.GetApplicationRequest{
		ApplicationId: c.appID,
	}

	resp, err := cli.GetApplication(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get application: %w", err)
	}

	bytes, err := json.Marshal(resp.Application)
	if err != nil {
		return fmt.Errorf("failed to marshal application: %w", err)
	}

	fmt.Fprintln(c.stdout, string(bytes))
	return nil
}
