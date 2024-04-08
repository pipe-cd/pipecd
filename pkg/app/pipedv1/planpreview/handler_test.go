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

package planpreview

import (
	"context"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/model"
)

type testCommandLister struct {
	commands []model.Command
}

func (l *testCommandLister) ListBuildPlanPreviewCommands() []model.ReportableCommand {
	out := make([]model.ReportableCommand, 0, len(l.commands))
	for i := range l.commands {
		out = append(out, model.ReportableCommand{
			Command: &l.commands[i],
			Report: func(ctx context.Context, status model.CommandStatus, metadata map[string]string, output []byte) error {
				return nil
			},
		})
	}
	return out
}

type testBuilder struct {
	recorder func(id string)
}

func (b *testBuilder) Build(ctx context.Context, id string, cmd model.Command_BuildPlanPreview) ([]*model.ApplicationPlanPreviewResult, error) {
	b.recorder(id)
	return nil, nil
}

func TestHandler(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cl := &testCommandLister{}
	handledCommands := make([]string, 0)
	var mu sync.Mutex
	var wg sync.WaitGroup

	handler := NewHandler(nil, nil, cl, nil, nil, nil, nil, nil,
		WithWorkerNum(2),
		// Use a long interval because we will directly call enqueueNewCommands function in this test.
		WithCommandCheckInterval(time.Hour),
	)
	handler.builderFactory = func() Builder {
		return &testBuilder{
			recorder: func(id string) {
				defer wg.Done()
				mu.Lock()
				defer mu.Unlock()
				handledCommands = append(handledCommands, id)
				sort.Strings(handledCommands)
			},
		}
	}
	go handler.Run(ctx)

	// CommandLister returns no command,
	// then there is no new command.
	handler.enqueueNewCommands(ctx)

	require.Equal(t, []string{}, handledCommands)

	// CommandLister returns 2 commands: 1, 2.
	// both of them will be considered as new commands.
	wg.Add(2)
	cl.commands = []model.Command{
		{
			Id:               "1",
			Type:             model.Command_BUILD_PLAN_PREVIEW,
			BuildPlanPreview: &model.Command_BuildPlanPreview{},
		},
		{
			Id:               "2",
			Type:             model.Command_BUILD_PLAN_PREVIEW,
			BuildPlanPreview: &model.Command_BuildPlanPreview{},
		},
	}
	handler.enqueueNewCommands(ctx)
	wg.Wait()
	require.Equal(t, []string{"1", "2"}, handledCommands)

	// CommandLister returns the same command list
	// so no new command will be added.
	handler.enqueueNewCommands(ctx)
	require.Equal(t, []string{"1", "2"}, handledCommands)

	// CommandLister returns commands: 2, 3.
	// then 3 will be considered as a new command.
	wg.Add(1)
	cl.commands = []model.Command{
		{
			Id:               "2",
			Type:             model.Command_BUILD_PLAN_PREVIEW,
			BuildPlanPreview: &model.Command_BuildPlanPreview{},
		},
		{
			Id:               "3",
			Type:             model.Command_BUILD_PLAN_PREVIEW,
			BuildPlanPreview: &model.Command_BuildPlanPreview{},
		},
	}
	handler.enqueueNewCommands(ctx)
	wg.Wait()
	require.Equal(t, []string{"1", "2", "3"}, handledCommands)
}
