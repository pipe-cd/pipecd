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

	"github.com/kapetaniosci/pipe/pkg/app/piped/executor"
	"github.com/kapetaniosci/pipe/pkg/model"
)

type Executor struct {
	executor.Input
}

func init() {
	var (
		f = func(in executor.Input) executor.Executor {
			return &Executor{
				Input: in,
			}
		}
		r = executor.DefaultRegistry()
	)
	r.Register(model.StageWaitApproval, f)
}

func (e *Executor) Execute(ctx context.Context) (model.StageStatus, error) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	e.LogPersister.Append("Waiting for an approval...")
	for {
		select {
		case <-ticker.C:
			if ok := e.checkApproval(ctx); !ok {
				continue
			}
			e.LogPersister.Append("Got an approval from abc")
			return model.StageStatus_STAGE_SUCCESS, nil

		case <-ctx.Done():
			return model.StageStatus_STAGE_CANCELLED, fmt.Errorf("context cancelled")
		}
	}
}

func (e *Executor) checkApproval(ctx context.Context) bool {
	var (
		command  *model.Command
		commands = e.CommandStore.ListDeploymentCommands(e.Deployment.Id)
	)

	for _, cmd := range commands {
		c := cmd.GetApproveStage()
		if c == nil {
			continue
		}
		if c.StageId != "e.Stage.Id" {
			continue
		}
		command = cmd
		break
	}
	if command == nil {
		return false
	}

	e.CommandStore.ReportCommandHandled(ctx, command, model.CommandStatus_COMMAND_SUCCEEDED, nil)
	return true
}
