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

package waitapproval

import (
	"context"
	"fmt"
	"time"

	"github.com/pipe-cd/pipe/pkg/app/piped/executor"
	"github.com/pipe-cd/pipe/pkg/model"
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
	r.Register(model.StageWaitApproval, f)
}

// Execute starts waiting until an approval from one of the specified users.
func (e *Executor) Execute(sig executor.StopSignal) model.StageStatus {
	var (
		originalStatus = e.Stage.Status
		ctx            = sig.Context()
		ticker         = time.NewTicker(5 * time.Second)
	)
	defer ticker.Stop()

	e.LogPersister.AppendInfo("Waiting for an approval...")
	for {
		select {
		case <-ticker.C:
			if commander, ok := e.checkApproval(ctx); ok {
				e.LogPersister.AppendInfo(fmt.Sprintf("Got an approval from %s", commander))
				return model.StageStatus_STAGE_SUCCESS
			}

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

func (e *Executor) checkApproval(ctx context.Context) (string, bool) {
	commands := e.CommandLister.ListCommands()

	for _, cmd := range commands {
		c := cmd.GetApproveStage()
		if c == nil {
			continue
		}

		if err := cmd.Report(ctx, model.CommandStatus_COMMAND_SUCCEEDED, nil); err == nil {
			return cmd.Commander, true
		}
		return "", false
	}

	return "", false
}
