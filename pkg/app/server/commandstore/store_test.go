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

package commandstore

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/cache/cachetest"
	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/model"
)

// fakeCommandStore is a minimal hand-rolled datastore.CommandStore used to
// control exactly what the backend returns for Get/UpdateStatus in each test.
type fakeCommandStore struct {
	getFunc          func(ctx context.Context, id string) (*model.Command, error)
	updateStatusFunc func(ctx context.Context, id string, status model.CommandStatus, metadata map[string]string, handledAt int64) error
}

func (f *fakeCommandStore) Add(ctx context.Context, cmd *model.Command) error {
	return nil
}

func (f *fakeCommandStore) Get(ctx context.Context, id string) (*model.Command, error) {
	return f.getFunc(ctx, id)
}

func (f *fakeCommandStore) List(ctx context.Context, opts datastore.ListOptions) ([]*model.Command, error) {
	return nil, nil
}

func (f *fakeCommandStore) UpdateStatus(ctx context.Context, id string, status model.CommandStatus, metadata map[string]string, handledAt int64) error {
	return f.updateStatusFunc(ctx, id, status, metadata, handledAt)
}

// TestUpdateCommandHandled_BackendGetFailureAfterSuccessfulUpdateStatus reproduces
// the reported cache poisoning scenario: UpdateStatus succeeds, but the
// subsequent backend.Get fails transiently. It verifies that the failed Get
// never results in a command being written to the cache, that the now-stale
// cache entry is deleted instead, and that the original UpdateCommandHandled
// behavior (returning nil once the status update itself succeeded) is
// preserved.
func TestUpdateCommandHandled_BackendGetFailureAfterSuccessfulUpdateStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// No EXPECT() is set on Put: if the fix regresses and a nil/zero-value
	// command reaches the cache, this mock call will fail the test.
	c := cachetest.NewMockCache(ctrl)
	c.EXPECT().Delete(cacheKey("command-id")).Return(nil)

	backend := &fakeCommandStore{
		updateStatusFunc: func(ctx context.Context, id string, status model.CommandStatus, metadata map[string]string, handledAt int64) error {
			return nil
		},
		getFunc: func(ctx context.Context, id string) (*model.Command, error) {
			return nil, errors.New("transient datastore error")
		},
	}

	s := &store{
		backend: backend,
		cache:   &commandCache{backend: c},
		logger:  zap.NewNop(),
	}

	err := s.UpdateCommandHandled(context.Background(), "command-id", model.CommandStatus_COMMAND_SUCCEEDED, nil, 12345)
	assert.NoError(t, err)
}

// TestGetCommand_StaleCacheDeletedAfterUpdateHandledGetFailure exercises the
// full regression scenario reported against the earlier fix: an old command
// is already cached (as GetCommand would leave it, e.g. from
// ReportCommandHandled's lookup), UpdateStatus then succeeds against the
// datastore, but the subsequent backend.Get fails transiently. It proves the
// now-stale cache entry is deleted (rather than left in place until its TTL
// expires), so a later GetCommand call misses the cache, falls back to the
// datastore, and observes the updated command/status instead of the stale
// one.
func TestGetCommand_StaleCacheDeletedAfterUpdateHandledGetFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	oldCommand := &model.Command{
		Id:        "command-id",
		PipedId:   "legit-piped-id",
		Status:    model.CommandStatus_COMMAND_NOT_HANDLED_YET,
		CreatedAt: 1700000000,
		UpdatedAt: 1700000000,
	}
	oldCommandData, err := json.Marshal(oldCommand)
	assert.NoError(t, err)

	updatedCommand := &model.Command{
		Id:        "command-id",
		PipedId:   "legit-piped-id",
		Status:    model.CommandStatus_COMMAND_SUCCEEDED,
		CreatedAt: 1700000000,
		UpdatedAt: 1700000100,
	}

	getCalls := 0
	backend := &fakeCommandStore{
		updateStatusFunc: func(ctx context.Context, id string, status model.CommandStatus, metadata map[string]string, handledAt int64) error {
			return nil
		},
		getFunc: func(ctx context.Context, id string) (*model.Command, error) {
			getCalls++
			if getCalls == 1 {
				// The Get triggered by UpdateCommandHandled fails transiently,
				// even though UpdateStatus above it already succeeded.
				return nil, errors.New("transient datastore error")
			}
			// The later Get, triggered by GetCommand's cache-miss fallback,
			// succeeds and returns the now-updated command.
			return updatedCommand, nil
		},
	}

	c := cachetest.NewMockCache(ctrl)
	key := cacheKey("command-id")
	gomock.InOrder(
		// Old command is cached, e.g. by an earlier GetCommand call.
		c.EXPECT().Get(key).Return(oldCommandData, nil),
		// UpdateStatus succeeds but backend.Get fails: the stale entry must
		// be deleted instead of being left in place.
		c.EXPECT().Delete(key).Return(nil),
		// Subsequent GetCommand call misses the now-deleted cache entry.
		c.EXPECT().Get(key).Return(nil, cache.ErrNotFound),
		// ...and falls back to the datastore, re-caching the fresh result.
		c.EXPECT().Put(key, gomock.Any()).Return(nil),
	)

	s := &store{
		backend: backend,
		cache:   &commandCache{backend: c},
		logger:  zap.NewNop(),
	}

	got, err := s.GetCommand(context.Background(), "command-id")
	assert.NoError(t, err)
	assert.Equal(t, oldCommand, got)

	err = s.UpdateCommandHandled(context.Background(), "command-id", model.CommandStatus_COMMAND_SUCCEEDED, nil, 12345)
	assert.NoError(t, err)

	got, err = s.GetCommand(context.Background(), "command-id")
	assert.NoError(t, err)
	assert.Equal(t, updatedCommand, got)
}

// TestGetCommand_FallsBackToDatastoreWhenCacheEntryIsPoisoned covers a cache
// entry that was already poisoned before this fix (e.g. by the JSON "null"
// bug), independent of how it got there. It proves GetCommand never returns
// the poisoned zero-value Command and instead falls back to the datastore
// for the real one, which is exactly what ReportCommandHandled's
// pipedID != cmd.PipedId check depends on to avoid a spurious PermissionDenied.
func TestGetCommand_FallsBackToDatastoreWhenCacheEntryIsPoisoned(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	realCommand := &model.Command{
		Id:        "command-id",
		PipedId:   "legit-piped-id",
		Status:    model.CommandStatus_COMMAND_SUCCEEDED,
		CreatedAt: 1700000000,
		UpdatedAt: 1700000000,
	}

	backend := &fakeCommandStore{
		getFunc: func(ctx context.Context, id string) (*model.Command, error) {
			return realCommand, nil
		},
	}

	c := cachetest.NewMockCache(ctrl)
	key := cacheKey("command-id")
	// Simulate a previously poisoned entry: raw JSON null in the cache.
	c.EXPECT().Get(key).Return([]byte("null"), nil)
	c.EXPECT().Delete(key).Return(nil)
	c.EXPECT().Put(key, gomock.Any()).Return(nil)

	s := &store{
		backend: backend,
		cache:   &commandCache{backend: c},
		logger:  zap.NewNop(),
	}

	got, err := s.GetCommand(context.Background(), "command-id")
	assert.NoError(t, err)
	assert.Equal(t, realCommand, got)
	assert.NotEmpty(t, got.PipedId, "must not return the poisoned zero-value command with an empty PipedId")
}

// TestGetCommand_ReturnsValidCachedCommandWithoutQueryingDatastore ensures the
// fix doesn't overreach: a genuinely valid cache entry is still served
// directly from the cache, without hitting the datastore.
func TestGetCommand_ReturnsValidCachedCommandWithoutQueryingDatastore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	backend := &fakeCommandStore{
		getFunc: func(ctx context.Context, id string) (*model.Command, error) {
			t.Fatal("datastore must not be queried when the cache holds a valid command")
			return nil, nil
		},
	}

	c := cachetest.NewMockCache(ctrl)
	key := cacheKey("command-id")
	c.EXPECT().Get(key).Return([]byte(`{"id":"command-id","piped_id":"legit-piped-id","created_at":1700000000,"updated_at":1700000000}`), nil)

	s := &store{
		backend: backend,
		cache:   &commandCache{backend: c},
		logger:  zap.NewNop(),
	}

	got, err := s.GetCommand(context.Background(), "command-id")
	assert.NoError(t, err)
	assert.Equal(t, &model.Command{Id: "command-id", PipedId: "legit-piped-id", CreatedAt: 1700000000, UpdatedAt: 1700000000}, got)
}
