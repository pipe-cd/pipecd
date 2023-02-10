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

package deployment

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

type logs struct {
	root *command

	deploymentID string
	stageID      string
	stdout       io.Writer
}

func newLogsCommand(root *command) *cobra.Command {
	c := &logs{
		root:   root,
		stdout: os.Stdout,
	}
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Get deployment stage logs.",
		RunE:  cli.WithContext(c.run),
	}

	cmd.Flags().StringVar(&c.deploymentID, "deployment-id", c.deploymentID, "The deployment ID.")

	cmd.MarkFlagRequired("deployment-id")

	return cmd
}

func (c *logs) run(ctx context.Context, input cli.Input) error {
	cli, err := c.root.clientOptions.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}
	defer cli.Close()

	req := &apiservice.ListStageLogRequest{
		DeploymentId: c.deploymentID,
	}

	resp, err := cli.ListStageLog(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get stage log: %w", err)
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal stage log: %w", err)
	}

	fmt.Fprintln(c.stdout, string(bytes))
	return nil
}
