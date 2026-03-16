// Copyright 2025 The PipeCD Authors.
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

package logpersister

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func newTestStageLogPersister() *stageLogPersister {
	apiClient := &fakeAPIClient{}
	p := NewPersister(apiClient, zap.NewNop())
	return &stageLogPersister{
		key:                     key{DeploymentID: "deploy-1", StageID: "stage-1"},
		doneCh:                  make(chan struct{}),
		checkpointFlushInterval: time.Minute,
		persister:               p,
		logger:                  zap.NewNop(),
	}
}

func TestStageLogPersister_Info(t *testing.T) {
	t.Parallel()

	sp := newTestStageLogPersister()
	sp.Info("hello info")

	require.Len(t, sp.blocks, 1)
	assert.Equal(t, "hello info", sp.blocks[0].Log)
	assert.Equal(t, model.LogSeverity_INFO, sp.blocks[0].Severity)
}

func TestStageLogPersister_Infof(t *testing.T) {
	t.Parallel()

	sp := newTestStageLogPersister()
	sp.Infof("hello %s %d", "world", 42)

	require.Len(t, sp.blocks, 1)
	assert.Equal(t, "hello world 42", sp.blocks[0].Log)
	assert.Equal(t, model.LogSeverity_INFO, sp.blocks[0].Severity)
}

func TestStageLogPersister_Success(t *testing.T) {
	t.Parallel()

	sp := newTestStageLogPersister()
	sp.Success("all good")

	require.Len(t, sp.blocks, 1)
	assert.Equal(t, "all good", sp.blocks[0].Log)
	assert.Equal(t, model.LogSeverity_SUCCESS, sp.blocks[0].Severity)
}

func TestStageLogPersister_Successf(t *testing.T) {
	t.Parallel()

	sp := newTestStageLogPersister()
	sp.Successf("deployed %s", "v1.2.3")

	require.Len(t, sp.blocks, 1)
	assert.Equal(t, "deployed v1.2.3", sp.blocks[0].Log)
	assert.Equal(t, model.LogSeverity_SUCCESS, sp.blocks[0].Severity)
}

func TestStageLogPersister_Error(t *testing.T) {
	t.Parallel()

	sp := newTestStageLogPersister()
	sp.Error("something failed")

	require.Len(t, sp.blocks, 1)
	assert.Equal(t, "something failed", sp.blocks[0].Log)
	assert.Equal(t, model.LogSeverity_ERROR, sp.blocks[0].Severity)
}

func TestStageLogPersister_Errorf(t *testing.T) {
	t.Parallel()

	sp := newTestStageLogPersister()
	sp.Errorf("exit code %d", 1)

	require.Len(t, sp.blocks, 1)
	assert.Equal(t, "exit code 1", sp.blocks[0].Log)
	assert.Equal(t, model.LogSeverity_ERROR, sp.blocks[0].Severity)
}

func TestStageLogPersister_Write(t *testing.T) {
	t.Parallel()

	sp := newTestStageLogPersister()
	n, err := sp.Write([]byte("written log"))

	require.NoError(t, err)
	assert.Equal(t, len("written log"), n)
	require.Len(t, sp.blocks, 1)
	assert.Equal(t, "written log", sp.blocks[0].Log)
	assert.Equal(t, model.LogSeverity_INFO, sp.blocks[0].Severity)
}

func TestStageLogPersister_AppendIncrementsIndex(t *testing.T) {
	t.Parallel()

	sp := newTestStageLogPersister()
	sp.Info("first")
	sp.Info("second")
	sp.Info("third")

	require.Len(t, sp.blocks, 3)
	assert.Less(t, sp.blocks[0].Index, sp.blocks[1].Index)
	assert.Less(t, sp.blocks[1].Index, sp.blocks[2].Index)
}

func TestStageLogPersister_IsStale_NotDoneNotCompleted(t *testing.T) {
	t.Parallel()

	sp := newTestStageLogPersister()
	assert.False(t, sp.isStale(time.Minute))
}

func TestStageLogPersister_IsStale_DoneIsTrue(t *testing.T) {
	t.Parallel()

	sp := newTestStageLogPersister()
	sp.done.Store(true)
	assert.True(t, sp.isStale(time.Minute))
}

func TestStageLogPersister_IsStale_CompletedAndExpired(t *testing.T) {
	t.Parallel()

	sp := newTestStageLogPersister()
	sp.mu.Lock()
	sp.completed = true
	sp.completedAt = time.Now().Add(-2 * time.Minute)
	sp.mu.Unlock()

	assert.True(t, sp.isStale(time.Minute))
}

func TestStageLogPersister_IsStale_CompletedButNotExpired(t *testing.T) {
	t.Parallel()

	sp := newTestStageLogPersister()
	sp.mu.Lock()
	sp.completed = true
	sp.completedAt = time.Now()
	sp.mu.Unlock()

	assert.False(t, sp.isStale(time.Minute))
}

func TestStageLogPersister_Complete_Timeout(t *testing.T) {
	t.Parallel()

	sp := newTestStageLogPersister()
	err := sp.Complete(10 * time.Millisecond)
	require.Error(t, err)
	assert.Equal(t, "timed out", err.Error())
}

func TestStageLogPersister_Complete_Success(t *testing.T) {
	t.Parallel()

	sp := newTestStageLogPersister()

	// Close doneCh to simulate successful flush completion.
	go func() {
		time.Sleep(10 * time.Millisecond)
		close(sp.doneCh)
	}()

	err := sp.Complete(time.Second)
	assert.NoError(t, err)
}
