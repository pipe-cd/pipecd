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

package deployment

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/pipe-cd/pipe/pkg/app/api/service/apiservice"
	"github.com/pipe-cd/pipe/pkg/cli"
	"github.com/pipe-cd/pipe/pkg/model"
)

type waitStatus struct {
	root *command

	deploymentID  string
	status        []string
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
		Short: "Wait until the specified status.",
		RunE:  cli.WithContext(c.run),
	}

	cmd.Flags().StringVar(&c.deploymentID, "deployment-id", c.deploymentID, "The deployment ID.")
	cmd.Flags().StringSliceVar(&c.status, "status", c.status, "The list of waiting statuses. (PENDING|PLANNED|RUNNING|ROLLING_BACK|SUCCESS|FAILURE|CANCELLED)")
	cmd.Flags().DurationVar(&c.checkInterval, "check-interval", c.checkInterval, "The interval of checking the deployment status.")
	cmd.Flags().DurationVar(&c.timeout, "timeout", c.timeout, "Maximum execution time.")

	cmd.MarkFlagRequired("deployment-id")
	cmd.MarkFlagRequired("status")

	return cmd
}

func (c *waitStatus) run(ctx context.Context, _ cli.Telemetry) error {
	logger := c.root.logOptions.NewLogger()

	statuses, err := makeDeploymentStatuses(c.status)
	if err != nil {
		logger.Fatal("Invalid deployment status: %v", err)
	}

	cli, err := c.root.clientOptions.NewClient(ctx)
	if err != nil {
		logger.Fatal("Failed to initialize client (%v)", err)
	}
	defer cli.Close()

	timer := time.NewTimer(c.timeout)
	defer timer.Stop()

	ticker := time.NewTicker(c.checkInterval)
	defer ticker.Stop()

	check := func() (status string, shouldRetry bool) {
		req := &apiservice.GetDeploymentRequest{
			DeploymentId: c.deploymentID,
		}
		resp, err := cli.GetDeployment(ctx, req)
		if err != nil {
			logger.Error("Failed while retrieving deployment information. Try again. (%v)", err)
			shouldRetry = true
			return
		}

		if _, ok := statuses[resp.Deployment.Status]; !ok {
			shouldRetry = true
			return
		}

		status = strings.TrimPrefix(resp.Deployment.Status.String(), "DEPLOYMENT_")
		return
	}

	// Do the first check immediately.
	status, shouldRetry := check()
	if !shouldRetry {
		logger.Info("Deployment is at %s status", status)
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-timer.C:
			logger.Fatal("Timed out %v", c.timeout)

		case <-ticker.C:
			status, shouldRetry := check()
			if shouldRetry {
				logger.Info("...")
				continue
			}

			logger.Info("Deployment is at %s status", status)
			return nil
		}
	}

	return nil
}

func makeDeploymentStatuses(statuses []string) (map[model.DeploymentStatus]struct{}, error) {
	out := make(map[model.DeploymentStatus]struct{}, len(statuses))
	for _, s := range statuses {
		status, ok := model.DeploymentStatus_value["DEPLOYMENT_"+s]
		if !ok {
			return nil, fmt.Errorf("bad status %s", s)
		}
		out[model.DeploymentStatus(status)] = struct{}{}
	}
	return out, nil
}
