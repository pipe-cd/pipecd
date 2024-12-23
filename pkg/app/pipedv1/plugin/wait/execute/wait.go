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

package execute

import (
	"context"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/piped/logpersister"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/wait/config"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
)

type Stage string

const (
	defaultDuration       = time.Minute
	logInterval           = 10 * time.Second
	startTimeKey          = "startTime"
	stageWait       Stage = "WAIT"
)

// Execute starts waiting for the specified duration.
func (s *deploymentServiceServer) execute(ctx context.Context, in *deployment.ExecutePluginInput, slp logpersister.StageLogPersister) model.StageStatus {
	var (
		// originalStatus = in.Stage.Status
		duration = defaultDuration
	)

	opts, err := config.DecodeStageOptionsYAML[waitStageOptions](in.StageConfig)
	if err != nil {
		slp.Errorf("failed to decode the stage configuration: %v", err)
		return model.StageStatus_STAGE_FAILURE
	}

	// Apply the stage configurations.
	if opts != nil {
		if opts.Duration > 0 {
			duration = opts.Duration.Duration()
		}
	}
	totalDuration := duration

	// Retrieve the saved startTime from the previous run.
	startTime := s.retrieveStartTime(in.Stage.Id)
	if !startTime.IsZero() {
		duration -= time.Since(startTime)
		if duration < 0 {
			duration = 0
		}
	} else {
		startTime = time.Now()
	}
	defer s.saveStartTime(ctx, startTime, in.Stage.Id)

	timer := time.NewTimer(duration)
	defer timer.Stop()

	ticker := time.NewTicker(logInterval)
	defer ticker.Stop()

	slp.Infof("Waiting for %v...", duration)
	for {
		select {
		case <-timer.C:
			slp.Infof("Waited for %v", totalDuration)
			return model.StageStatus_STAGE_SUCCESS

		case <-ticker.C:
			slp.Infof("%v elapsed...", time.Since(startTime))

			/** TODO: handle StopSignal
			case s := <-sig.Ch():
				switch s {
				case executor.StopSignalCancel:
					return &deployment.ExecuteStageResponse{
						Status: model.StageStatus_STAGE_CANCELLED,
					}, nil
				case executor.StopSignalTerminate:
					return &deployment.ExecuteStageResponse{
						Status: originalStatus,
					}, nil
				default:
					return &deployment.ExecuteStageResponse{
						Status: model.StageStatus_STAGE_FAILURE,
					}, nil // TODO: Return an error message like "received an unknown signal".
				}
			}
			*/
		}
	}
}

func (s *deploymentServiceServer) retrieveStartTime(stageId string) (t time.Time) {
	sec, ok := s.metadataStore.Stage(stageId).Get(startTimeKey)
	if !ok {
		return
	}
	ut, err := strconv.ParseInt(sec, 10, 64)
	if err != nil {
		return
	}
	return time.Unix(ut, 0)
}

func (s *deploymentServiceServer) saveStartTime(ctx context.Context, t time.Time, stageId string) {
	metadata := map[string]string{
		startTimeKey: strconv.FormatInt(t.Unix(), 10),
	}
	if err := s.metadataStore.Stage(stageId).PutMulti(ctx, metadata); err != nil {
		s.logger.Error("failed to store metadata", zap.Error(err))
	}
}
