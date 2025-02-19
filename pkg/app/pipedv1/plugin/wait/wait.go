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

package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/piped/logpersister"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

const (
	logInterval  = 10 * time.Second
	startTimeKey = "startTime"
)

// executeWait starts waiting for the specified duration.
func (p *plugin) executeWait(ctx context.Context, in *sdk.ExecuteStageInput) sdk.StageStatus {
	opts, err := decode(in.Request.StageConfig)
	if err != nil {
		in.Client.LogPersister.Errorf("failed to decode the stage config: %v", err)
		return sdk.StageStatusFailure
	}

	duration := opts.Duration.Duration()

	// Retrieve the saved initialStart from the previous run.
	initialStart := p.retrieveStartTime(ctx, in.Client, in.Logger)
	if initialStart.IsZero() {
		// When this is the first run.
		initialStart = time.Now()
	}
	p.saveStartTime(ctx, in.Client, initialStart, in.Logger)

	return wait(ctx, duration, initialStart, in.Client.LogPersister)
}

func wait(ctx context.Context, duration time.Duration, initialStart time.Time, slp logpersister.StageLogPersister) sdk.StageStatus {
	remaining := duration - time.Since(initialStart)
	if remaining <= 0 {
		// When this stage restarted and the duration has already passed.
		slp.Infof("Already waited for %v since %v", duration, initialStart.Local())
		return sdk.StageStatusSuccess
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
			return sdk.StageStatusSuccess

		case <-ticker.C: // on interval elapsed
			slp.Infof("%v elapsed...", time.Since(initialStart))

		case <-ctx.Done(): // on cancelled
			slp.Info("Wait cancelled")
			return sdk.StageStatusCancelled
		}
	}
}

func (p *plugin) retrieveStartTime(ctx context.Context, client *sdk.Client, logger *zap.Logger) time.Time {
	sec, err := client.GetStageMetadata(ctx, startTimeKey)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get stage metadata %s", startTimeKey), zap.Error(err))
		return time.Time{}
	}

	ut, err := strconv.ParseInt(sec, 10, 64)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to parse stage metadata %s", startTimeKey), zap.Error(err))
		return time.Time{}
	}
	return time.Unix(ut, 0)
}

func (p *plugin) saveStartTime(ctx context.Context, client *sdk.Client, t time.Time, logger *zap.Logger) {
	value := strconv.FormatInt(t.Unix(), 10)
	if err := client.PutStageMetadata(ctx, startTimeKey, value); err != nil {
		logger.Error(fmt.Sprintf("failed to store %s as stage metadata %s", value, startTimeKey), zap.Error(err))
	}
}
