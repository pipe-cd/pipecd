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
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/pipe-cd/pipe/pkg/app/api/service/apiservice"
	"github.com/pipe-cd/pipe/pkg/cli"
	"github.com/pipe-cd/pipe/pkg/model"
)

type sync struct {
	root *command

	appID         string
	checkInterval time.Duration
	timeout       time.Duration
}

func newSyncCommand(root *command) *cobra.Command {
	c := &sync{
		root:          root,
		checkInterval: 15 * time.Second,
		timeout:       5 * time.Minute,
	}
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync an application.",
		RunE:  cli.WithContext(c.run),
	}

	cmd.Flags().StringVar(&c.appID, "app-id", c.appID, "The application ID.")
	cmd.Flags().DurationVar(&c.checkInterval, "check-interval", c.checkInterval, "The interval of checking the requested command.")
	cmd.Flags().DurationVar(&c.timeout, "timeout", c.timeout, "Maximum execution time.")

	cmd.MarkFlagRequired("app-id")

	return cmd
}

func (c *sync) run(ctx context.Context, _ cli.Telemetry) error {
	logger := c.root.logOptions.NewLogger()

	cli, err := c.root.clientOptions.NewClient(ctx)
	if err != nil {
		logger.Fatal("Failed to initialize client: %v", err)
	}
	defer cli.Close()

	req := &apiservice.SyncApplicationRequest{
		ApplicationId: c.appID,
	}
	resp, err := cli.SyncApplication(ctx, req)
	if err != nil {
		logger.Fatal("Failed to sync application: %v", err)
	}

	logger.Info("Sent a request to sync application and waiting to be accepted...")

	timer := time.NewTimer(c.timeout)
	defer timer.Stop()

	ticker := time.NewTicker(c.checkInterval)
	defer ticker.Stop()

	check := func() (deploymentID string, shouldRetry bool) {
		const triggeredDeploymentIDKey = "TriggeredDeploymentID"
		cmd, err := retrieveSyncCommand(ctx, cli, resp.CommandId)
		if err != nil {
			logger.Error("Failed while retrieving command information. Try again. (%v)", err)
			shouldRetry = true
			return
		}

		if cmd.Type != model.Command_SYNC_APPLICATION {
			logger.Error("Unexpected command type, want: %s, got: %s", model.Command_SYNC_APPLICATION.String(), cmd.Type.String())
			return
		}

		switch cmd.Status {
		case model.CommandStatus_COMMAND_SUCCEEDED:
			deploymentID = cmd.Metadata[triggeredDeploymentIDKey]
			return

		case model.CommandStatus_COMMAND_FAILED:
			logger.Error("The request was unable to handle")
			return

		case model.CommandStatus_COMMAND_TIMEOUT:
			logger.Error("The request was timed out")
			return

		default:
			shouldRetry = true
			return
		}
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-timer.C:
			logger.Fatal("Timed out %v", c.timeout)

		case <-ticker.C:
			deploymentID, shouldRetry := check()
			if shouldRetry {
				logger.Info("...")
				continue
			}
			if deploymentID == "" {
				logger.Fatal("Unable to detect the triggered deployment ID")
			}

			logger.Info("Successfully triggered deployment %s", deploymentID)
			fmt.Println(deploymentID)
			return nil
		}
	}
}

func retrieveSyncCommand(ctx context.Context, cli apiservice.Client, cmdID string) (*model.Command, error) {
	req := &apiservice.GetCommandRequest{
		CommandId: cmdID,
	}
	resp, err := cli.GetCommand(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Command, nil
}
