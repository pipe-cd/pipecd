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
	"fmt"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"iter"

	"slices"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/pipe-cd/piped-plugin-sdk-go/logpersister/logpersistertest"
	"github.com/stretchr/testify/assert"
)

// mockStageClient satisfies StageClient for tests.
type mockStageClient struct {
	lp            sdk.StageLogPersister
	stageMetadata map[string]string
	commands      []*sdk.StageCommand
	// keepPolling makes ListStageCommands wait for new commands until the
	// context is done after returning the existing ones, as the SDK client does.
	keepPolling bool
	// waiting is closed when ListStageCommands starts waiting for new commands.
	waiting chan struct{}
}

func newMockStageClient(t *testing.T) *mockStageClient {
	t.Helper()
	return &mockStageClient{
		lp:            logpersistertest.NewTestLogPersister(t),
		stageMetadata: map[string]string{},
	}
}

func (m *mockStageClient) GetStageMetadata(_ context.Context, key string) (string, bool, error) {
	v, ok := m.stageMetadata[key]
	return v, ok, nil
}

func (m *mockStageClient) PutStageMetadata(_ context.Context, key, value string) error {
	m.stageMetadata[key] = value
	return nil
}

func (m *mockStageClient) ListStageCommands(ctx context.Context, _ ...sdk.CommandType) iter.Seq2[*sdk.StageCommand, error] {
	return func(yield func(*sdk.StageCommand, error) bool) {
		for _, c := range m.commands {
			if !yield(c, nil) {
				return
			}
		}
		if !m.keepPolling {
			return
		}
		if m.waiting != nil {
			close(m.waiting)
		}
		<-ctx.Done()
		yield(nil, ctx.Err())
	}
}

// recordingLogPersister records the info logs written to the stage log.
type recordingLogPersister struct {
	sdk.StageLogPersister
	mu    sync.Mutex
	infos []string
}

func (lp *recordingLogPersister) Infof(format string, a ...interface{}) {
	lp.mu.Lock()
	lp.infos = append(lp.infos, fmt.Sprintf(format, a...))
	lp.mu.Unlock()
	lp.StageLogPersister.Infof(format, a...)
}

func (lp *recordingLogPersister) infoLogs() []string {
	lp.mu.Lock()
	defer lp.mu.Unlock()
	return slices.Clone(lp.infos)
}

func Test_checkApproval(t *testing.T) {
	t.Parallel()
	p := &plugin{}

	testcases := []struct {
		name         string
		minApprovers int
		existing     string
		commands     []*sdk.StageCommand
		keepPolling  bool
		wantApproved bool
		wantUsers    []string
	}{
		{
			name:         "single approval meets threshold",
			minApprovers: 1,
			existing:     "",
			commands:     []*sdk.StageCommand{{Commander: "alice", Type: sdk.CommandTypeApproveStage}},
			wantApproved: true,
			wantUsers:    []string{"alice"},
		},
		{
			name:         "merge existing and new unique approver",
			minApprovers: 2,
			existing:     "alice",
			commands: []*sdk.StageCommand{
				{Commander: "alice", Type: sdk.CommandTypeApproveStage},
				{Commander: "bob", Type: sdk.CommandTypeApproveStage},
			},
			wantApproved: true,
			wantUsers:    []string{"alice", "bob"},
		},
		{
			name:         "not enough approvals persists progress",
			minApprovers: 2,
			existing:     "",
			commands:     []*sdk.StageCommand{{Commander: "alice", Type: sdk.CommandTypeApproveStage}},
			wantApproved: false,
			wantUsers:    []string{"alice"},
		},
		{
			name:         "enough approvals while waiting for more commands",
			minApprovers: 2,
			existing:     "",
			commands: []*sdk.StageCommand{
				{Commander: "alice", Type: sdk.CommandTypeApproveStage},
				{Commander: "bob", Type: sdk.CommandTypeApproveStage},
			},
			keepPolling:  true,
			wantApproved: true,
			wantUsers:    []string{"alice", "bob"},
		},
		{
			name:         "existing approvals already meet threshold",
			minApprovers: 2,
			existing:     "alice, bob",
			commands: []*sdk.StageCommand{
				{Commander: "alice", Type: sdk.CommandTypeApproveStage},
				{Commander: "bob", Type: sdk.CommandTypeApproveStage},
			},
			keepPolling:  true,
			wantApproved: true,
			wantUsers:    []string{"alice", "bob"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			m := newMockStageClient(t)
			if tc.existing != "" {
				m.stageMetadata[sdk.MetadataKeyStageApprovedUsers] = tc.existing
			}
			m.commands = tc.commands
			m.keepPolling = tc.keepPolling

			ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
			defer cancel()

			approved := p.checkApproval(ctx, tc.minApprovers, m.lp, m)
			assert.Equal(t, tc.wantApproved, approved)

			gotUsers := strings.Split(m.stageMetadata[sdk.MetadataKeyStageApprovedUsers], ", ")
			sort.Strings(gotUsers)
			wantUsers := append([]string(nil), tc.wantUsers...)
			sort.Strings(wantUsers)
			assert.Equal(t, wantUsers, gotUsers)

			disp := m.stageMetadata[sdk.MetadataKeyStageDisplay]
			for _, u := range tc.wantUsers {
				assert.Contains(t, disp, u)
			}
		})
	}
}

func Test_checkApproval_reportsProgressWhileWaiting(t *testing.T) {
	t.Parallel()
	p := &plugin{}

	m := newMockStageClient(t)
	m.commands = []*sdk.StageCommand{{Commander: "alice", Type: sdk.CommandTypeApproveStage}}
	m.keepPolling = true
	m.waiting = make(chan struct{})
	lp := &recordingLogPersister{StageLogPersister: m.lp}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan bool, 1)
	go func() {
		done <- p.checkApproval(ctx, 2, lp, m)
	}()

	select {
	case <-m.waiting:
	case <-time.After(5 * time.Second):
		t.Fatal("ListStageCommands did not start waiting for new commands")
	}

	// The progress must be reported while ListStageCommands is still waiting.
	assert.Equal(t, "alice", m.stageMetadata[sdk.MetadataKeyStageApprovedUsers])
	assert.Equal(t, "Approved by: alice", m.stageMetadata[sdk.MetadataKeyStageDisplay])
	assert.Contains(t, lp.infoLogs(), "Waiting for 1 more approver(s)...")

	cancel()
	assert.False(t, <-done)
}

func Test_getApprovedUsers(t *testing.T) {
	t.Parallel()
	p := &plugin{}
	ctx := context.Background()

	testcases := []struct {
		name     string
		existing string
		want     []string
	}{
		{name: "empty", existing: "", want: []string{}},
		{name: "existing two", existing: "alice, bob", want: []string{"alice", "bob"}},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			m := newMockStageClient(t)
			if tc.existing != "" {
				m.stageMetadata[sdk.MetadataKeyStageApprovedUsers] = tc.existing
			}
			users, err := p.getApprovedUsers(ctx, m)
			assert.NoError(t, err)
			if len(tc.want) == 0 {
				assert.Empty(t, users)
				return
			}
			sort.Strings(users)
			want := slices.Clone(tc.want)
			sort.Strings(want)
			assert.Equal(t, want, users)
		})
	}
}
