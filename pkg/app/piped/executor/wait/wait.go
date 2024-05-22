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

package wait

import (
	"context"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	defaultDuration = time.Minute
	logInterval     = 10 * time.Second
	startTimeKey    = "startTime"
)

type Executor struct {
	executor.Input
}

type registerer interface {
	Register(stage model.Stage, f executor.Factory) error
}

// Register registers this executor factory into a given registerer.
func Register(r registerer) {
	f := func(in executor.Input) executor.Executor {
		return &Executor{
			Input: in,
		}
	}
	r.Register(model.StageWait, f)
}

// Execute starts waiting for the specified duration.
func (e *Executor) Execute(sig executor.StopSignal) model.StageStatus {
	// Skip the stage if needed based on the skip config.
	skip, err := executor.CheckSkipStage(sig.Context(), e.Input, e.StageConfig.WaitStageOptions.SkipOptions)
	if err != nil {
		e.Logger.Error("failed to check whether skipping the stage", zap.Error(err))
		return model.StageStatus_STAGE_FAILURE
	}
	if skip {
		return model.StageStatus_STAGE_SKIPPED
	}

	var (
		originalStatus = e.Stage.Status
		duration       = defaultDuration
	)

	// Apply the stage configurations.
	if opts := e.StageConfig.WaitStageOptions; opts != nil {
		if opts.Duration > 0 {
			duration = opts.Duration.Duration()
		}
	}
	totalDuration := duration

	// Retrieve the saved startTime from the previous run.
	startTime := e.retrieveStartTime()
	if !startTime.IsZero() {
		duration -= time.Since(startTime)
		if duration < 0 {
			duration = 0
		}
	} else {
		startTime = time.Now()
	}
	defer e.saveStartTime(sig.Context(), startTime)

	timer := time.NewTimer(duration)
	defer timer.Stop()

	ticker := time.NewTicker(logInterval)
	defer ticker.Stop()

	e.LogPersister.Infof("Waiting for %v...", duration)
	for {
		select {
		case <-timer.C:
			e.LogPersister.Infof("Waited for %v", totalDuration)
			return model.StageStatus_STAGE_SUCCESS

		case <-ticker.C:
			e.LogPersister.Infof("%v elapsed...", time.Since(startTime))

		case s := <-sig.Ch():
			switch s {
			case executor.StopSignalCancel:
				return model.StageStatus_STAGE_CANCELLED
			case executor.StopSignalTerminate:
				return originalStatus
			default:
				return model.StageStatus_STAGE_FAILURE
			}
		}
	}
}

func (e *Executor) retrieveStartTime() (t time.Time) {
	s, ok := e.MetadataStore.Stage(e.Stage.Id).Get(startTimeKey)
	if !ok {
		return
	}
	ut, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return
	}
	return time.Unix(ut, 0)
}

func (e *Executor) saveStartTime(ctx context.Context, t time.Time) {
	metadata := map[string]string{
		startTimeKey: strconv.FormatInt(t.Unix(), 10),
	}
	if err := e.MetadataStore.Stage(e.Stage.Id).PutMulti(ctx, metadata); err != nil {
		e.Logger.Error("failed to store metadata", zap.Error(err))
	}
}
