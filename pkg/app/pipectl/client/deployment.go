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

package client

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/model"
)

// WaitDeploymentStatuses waits a given deployment until it reaches one of the specified statuses.
func WaitDeploymentStatuses(
	ctx context.Context,
	cli apiservice.Client,
	deploymentID string,
	statuses []model.DeploymentStatus,
	checkInterval, timeout time.Duration,
	logger *zap.Logger,
) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	statusMap := makeDeploymentStatusesMap(statuses)
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	check := func() (status string, shouldRetry bool) {
		req := &apiservice.GetDeploymentRequest{
			DeploymentId: deploymentID,
		}
		resp, err := cli.GetDeployment(ctx, req)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed while retrieving deployment information. Try again. (%v)", err))
			shouldRetry = true
			return
		}

		if _, ok := statusMap[resp.Deployment.Status]; !ok {
			shouldRetry = true
			return
		}

		status = resp.Deployment.Status.String()
		return
	}

	// Do the first check immediately.
	status, shouldRetry := check()
	if !shouldRetry {
		logger.Info(fmt.Sprintf("Deployment is at %s status", status))
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			status, shouldRetry := check()
			if shouldRetry {
				logger.Info("...")
				continue
			}

			logger.Info(fmt.Sprintf("Deployment is at %s status", status))
			return nil
		}
	}
}

func makeDeploymentStatusesMap(statuses []model.DeploymentStatus) map[model.DeploymentStatus]struct{} {
	out := make(map[model.DeploymentStatus]struct{}, len(statuses))
	for _, s := range statuses {
		out[model.DeploymentStatus(s)] = struct{}{}
	}
	return out
}
