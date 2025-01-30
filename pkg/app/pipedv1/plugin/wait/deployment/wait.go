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
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/piped/logpersister"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/wait/config"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/pipedservice"
)

type Stage string

const (
	logInterval        = 10 * time.Second
	startTimeKey       = "startTime"
	stageWait    Stage = "WAIT"
)

// Execute starts waiting for the specified duration.
func (s *deploymentServiceServer) execute(ctx context.Context, in *deployment.ExecutePluginInput, slp logpersister.StageLogPersister) model.StageStatus {
	opts, err := config.Decode(in.StageConfig)
	if err != nil {
		slp.Errorf("failed to decode the stage config: %v", err)
		return model.StageStatus_STAGE_FAILURE
	}

	duration := opts.Duration.Duration()

	// Retrieve the saved initialStart from the previous run.
	initialStart := s.retrieveStartTime(ctx, in.Deployment.Id, in.Stage.Id)
	if initialStart.IsZero() {
		// When this is the first run.
		initialStart = time.Now()
	}
	s.saveStartTime(ctx, initialStart, in.Deployment.Id, in.Stage.Id)

	return wait(ctx, duration, initialStart, slp)
}

func wait(ctx context.Context, duration time.Duration, initialStart time.Time, slp logpersister.StageLogPersister) model.StageStatus {
	remaining := duration - time.Since(initialStart)
	if remaining <= 0 {
		// When this stage restarted and the duration has already passed.
		slp.Infof("Already waited for %v since %v", duration, initialStart.Local())
		return model.StageStatus_STAGE_SUCCESS
	}

	timer := time.NewTimer(remaining)
	defer timer.Stop()

	ticker := time.NewTicker(logInterval)
	defer ticker.Stop()

	slp.Infof("Waiting for %v since %v...", duration, initialStart.Local())
	for {
		select {
		case <-timer.C: // on completed
			slp.Infof("Waited for %v", duration)
			return model.StageStatus_STAGE_SUCCESS

		case <-ticker.C: // on interval elapsed
			slp.Infof("%v elapsed...", time.Since(initialStart))

		case <-ctx.Done(): // on cancelled
			slp.Info("Wait cancelled")
			return model.StageStatus_STAGE_CANCELLED
		}
	}
}

func (s *deploymentServiceServer) retrieveStartTime(ctx context.Context, deploymentID, stageID string) (t time.Time) {
	sec, err := s.metadataStore.GetStageMetadata(ctx, &pipedservice.GetStageMetadataRequest{
		DeploymentId: deploymentID,
		StageId:      stageID,
		Key:          startTimeKey,
	})
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to get stage metadata %s", startTimeKey), zap.Error(err))
		return
	}
	ut, err := strconv.ParseInt(sec.Value, 10, 64)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to parse stage metadata %s", startTimeKey), zap.Error(err))
		return
	}
	return time.Unix(ut, 0)
}

func (s *deploymentServiceServer) saveStartTime(ctx context.Context, t time.Time, deploymentID, stageID string) {
	req := &pipedservice.PutStageMetadataRequest{
		DeploymentId: deploymentID,
		StageId:      stageID,
		Key:          startTimeKey,
		Value:        strconv.FormatInt(t.Unix(), 10),
	}
	if _, err := s.metadataStore.PutStageMetadata(ctx, req); err != nil {
		s.logger.Error(fmt.Sprintf("failed to store %s as stage metadata %s", req.Value, startTimeKey), zap.Error(err))
	}
}
