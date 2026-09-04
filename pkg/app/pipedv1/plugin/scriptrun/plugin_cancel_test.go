// Copyright 2026 The PipeCD Authors.
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
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/pipe-cd/piped-plugin-sdk-go/logpersister/logpersistertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// alive reports whether pid is still running. Signal 0 performs the permission
// and existence checks without delivering anything.
func alive(pid int) bool {
	return syscall.Kill(pid, 0) == nil
}

func waitGone(pid int, d time.Duration) bool {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if !alive(pid) {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return !alive(pid)
}

// A cancelled stage must take the whole process tree with it. The script here
// backgrounds a child that traps SIGTERM and keeps running, which is what a
// naive kill of /bin/sh alone would leave behind holding cluster state.
func TestExecuteCommandKillsProcessGroupOnCancel(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	pidFile := filepath.Join(dir, "child.pid")

	script := fmt.Sprintf(`
trap '' TERM
( trap '' TERM; echo $$ > %q; while true; do sleep 0.05; done ) &
while true; do sleep 0.05; done
`, pidFile)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan sdk.StageStatus, 1)
	go func() {
		done <- executeCommand(ctx, script, nil, sdk.ExecuteStageRequest[struct{}]{
			StageName:  stageScriptRun,
			Deployment: sdk.Deployment{ID: "deployment-1", ApplicationID: "app-1"},
		}, logpersistertest.NewTestLogPersister(t))
	}()

	// Wait for the grandchild to record its pid.
	var childPID int
	require.Eventually(t, func() bool {
		b, err := os.ReadFile(pidFile)
		if err != nil {
			return false
		}
		_, err = fmt.Sscanf(string(b), "%d", &childPID)
		return err == nil && childPID > 0
	}, 10*time.Second, 20*time.Millisecond, "grandchild never started")

	require.True(t, alive(childPID), "grandchild should be running before cancel")

	cancel()

	select {
	case <-done:
	case <-time.After(30 * time.Second):
		t.Fatal("executeCommand did not return after cancel")
	}

	assert.True(t, waitGone(childPID, 15*time.Second),
		"grandchild %d survived cancellation; the process group was not terminated", childPID)
}

// The normal path must be unaffected: a command that finishes on its own still
// reports success and does not wait out the grace period.
func TestExecuteCommandSucceedsWithoutCancel(t *testing.T) {
	t.Parallel()

	start := time.Now()
	status := executeCommand(context.Background(), "echo hello", nil, sdk.ExecuteStageRequest[struct{}]{
		StageName:  stageScriptRun,
		Deployment: sdk.Deployment{ID: "deployment-1", ApplicationID: "app-1"},
	}, logpersistertest.NewTestLogPersister(t))

	assert.Equal(t, sdk.StageStatusSuccess, status)
	assert.Less(t, time.Since(start), commandTerminationGracePeriod,
		"a command that exits on its own must not wait for the termination grace period")
}
