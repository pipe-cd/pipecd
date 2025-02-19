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

package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister/logpersistertest"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

func TestWait_Complete(t *testing.T) {
	t.Parallel()

	duration := 50 * time.Millisecond

	resultCh := make(chan sdk.StageStatus)
	go func() {
		result := wait(context.Background(), duration, time.Now(), logpersistertest.NewTestLogPersister(t))
		resultCh <- result
	}()

	// Assert that wait() didn't end before the specified duration has passed.
	select {
	case <-resultCh:
		t.Error("wait() ended too early")
	case <-time.After(duration / 10):
	}

	// Assert that wait() ends after the specified duration has passed.
	select {
	case result := <-resultCh:
		assert.Equal(t, sdk.StageStatusSuccess, result)
	case <-time.After(duration):
		// Wait 1.1x duration in total to avoid flaky test.
		t.Error("wait() did not end even after the specified duration has passed")
	}
}

func TestWait_Cancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultCh := make(chan sdk.StageStatus)
	go func() {
		result := wait(ctx, 1*time.Second, time.Now(), logpersistertest.NewTestLogPersister(t))
		resultCh <- result
	}()

	cancel()

	select {
	case result := <-resultCh:
		assert.Equal(t, sdk.StageStatusCancelled, result)
	case <-time.After(1 * time.Second):
		t.Error("wait() did not ended even after the context was canceled")
	}
}

func TestWait_RestartAfterLongTime(t *testing.T) {
	t.Parallel()
	// Suppose this stage started 2 hours ago but it was interrupted.
	previousStart := time.Now().Add(-2 * time.Hour)

	result := wait(context.Background(), 1*time.Second, previousStart, logpersistertest.NewTestLogPersister(t))
	// Immediately return success because the duration has already passed.
	assert.Equal(t, sdk.StageStatusSuccess, result)
}

func TestWait_RestartAndContinue(t *testing.T) {
	t.Parallel()
	// Imagine this timeline:
	//   begin          interrupted  now         (1)not end  (2)end
	//   | <--------30ms--|--------> | <--10ms--> | <--15s--> |
	//   | <------------- 50ms ------------------------> |
	duration := 50 * time.Millisecond
	previousStart := time.Now().Add(-30 * time.Millisecond)

	resultCh := make(chan sdk.StageStatus)
	go func() {
		result := wait(context.Background(), duration, previousStart, logpersistertest.NewTestLogPersister(t))
		resultCh <- result
	}()

	// (1) Assert that wait() didn't end before the specified duration has passed.
	select {
	case <-resultCh:
		t.Error("wait() ended too early")
	case <-time.After(10 * time.Millisecond):
	}

	// (2) Assert that wait() ends after the specified duration has passed.
	select {
	case result := <-resultCh:
		assert.Equal(t, sdk.StageStatusSuccess, result)
	case <-time.After(15 * time.Millisecond): // Not 50ms
		// Wait 55ms in total to avoid flaky test.
		t.Error("wait() did not end even after the specified duration has passed")
	}
}
