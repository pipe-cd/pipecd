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

package deployment

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/pipe-cd/pipecd/pkg/app/pipectl/client"
	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type waitStatus struct {
	root *command

	deploymentID  string
	statuses      []string
	checkInterval time.Duration
	timeout       time.Duration
}

func newWaitStatusCommand(root *command) *cobra.Command {
	c := &waitStatus{
		root:          root,
		checkInterval: 15 * time.Second,
		timeout:       15 * time.Minute,
	}
	cmd := &cobra.Command{
		Use:   "wait-status",
		Short: "Wait for one of the specified statuses.",
		RunE:  cli.WithContext(c.run),
	}

	cmd.Flags().StringVar(&c.deploymentID, "deployment-id", c.deploymentID, "The deployment ID.")
	cmd.Flags().StringSliceVar(&c.statuses, "status", c.statuses, fmt.Sprintf("The list of waiting statuses. (%s)", strings.Join(model.DeploymentStatusStrings(), "|")))
	cmd.Flags().DurationVar(&c.checkInterval, "check-interval", c.checkInterval, "The interval of checking the deployment status.")
	cmd.Flags().DurationVar(&c.timeout, "timeout", c.timeout, "Maximum execution time.")

	cmd.MarkFlagRequired("deployment-id")
	cmd.MarkFlagRequired("status")

	return cmd
}

func (c *waitStatus) run(ctx context.Context, input cli.Input) error {
	statuses, err := model.DeploymentStatusesFromStrings(c.statuses)
	if err != nil {
		return fmt.Errorf("invalid deployment status: %w", err)
	}

	cli, err := c.root.clientOptions.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}
	defer cli.Close()

	return client.WaitDeploymentStatuses(
		ctx,
		cli,
		c.deploymentID,
		statuses,
		c.checkInterval,
		c.timeout,
		input.Logger,
	)
}
