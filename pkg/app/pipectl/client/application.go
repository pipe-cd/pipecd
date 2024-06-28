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

package client

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/model"
)

// SyncApplication sents a command to sync a given application and waits until it has been triggered.
// The deployment ID will be returned or an error.
func SyncApplication(
	ctx context.Context,
	cli apiservice.Client,
	appID string,
	checkInterval, timeout time.Duration,
	logger *zap.Logger,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req := &apiservice.SyncApplicationRequest{
		ApplicationId: appID,
	}
	resp, err := cli.SyncApplication(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to sync application %w", err)
	}

	logger.Info("Sent a request to sync application and waiting to be accepted...")

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	check := func() (deploymentID string, shouldRetry bool) {
		cmd, err := getCommand(ctx, cli, resp.CommandId)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed while retrieving command information. Try again. (%v)", err))
			shouldRetry = true
			return
		}

		if cmd.Type != model.Command_SYNC_APPLICATION {
			logger.Error(fmt.Sprintf("Unexpected command type, want: %s, got: %s", model.Command_SYNC_APPLICATION.String(), cmd.Type.String()))
			return
		}

		switch cmd.Status {
		case model.CommandStatus_COMMAND_SUCCEEDED:
			deploymentID = cmd.Metadata[model.MetadataKeyTriggeredDeploymentID]
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
			return "", ctx.Err()

		case <-ticker.C:
			deploymentID, shouldRetry := check()
			if shouldRetry {
				logger.Info("...")
				continue
			}
			if deploymentID == "" {
				return "", fmt.Errorf("failed to detect the triggered deployment ID")
			}
			return deploymentID, nil
		}
	}
}
