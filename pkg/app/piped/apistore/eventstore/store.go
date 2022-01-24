// Copyright 2021 The PipeCD Authors.
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

package eventstore

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/model"
)

// Getter helps get an event. All objects returned here must be treated as read-only.
type Getter interface {
	// GetLatest returns the latest event that meets the given conditions.
	GetLatest(ctx context.Context, name string, labels map[string]string, notUsingCache bool) (event *model.Event, cacheUsed bool, ok bool)
}

type Store interface {
	// Run starts syncing the event list with the control-plane.
	Run(ctx context.Context) error
	// Getter returns a getter for retrieving an event.
	Getter() Getter
}

type apiClient interface {
	GetLatestEvent(ctx context.Context, req *pipedservice.GetLatestEventRequest, opts ...grpc.CallOption) (*pipedservice.GetLatestEventResponse, error)
	ListEvents(ctx context.Context, req *pipedservice.ListEventsRequest, opts ...grpc.CallOption) (*pipedservice.ListEventsResponse, error)
}

type store struct {
	apiClient    apiClient
	syncInterval time.Duration
	gracePeriod  time.Duration
	logger       *zap.Logger

	// Mark that it has handled all events that was created before this UNIX time.
	milestone int64
	mu        sync.RWMutex
	// The key is supposed to be a string consists of name and labels.
	// And the value is the address to the latest Event.
	latestEvents map[string]*model.Event
}

const (
	defaultSyncInterval = time.Minute
)

// NewStore creates a new event store instance.
// This syncs with the control plane to keep the list of events for this runner up-to-date.
func NewStore(apiClient apiClient, gracePeriod time.Duration, logger *zap.Logger) Store {
	return &store{
		apiClient:    apiClient,
		syncInterval: defaultSyncInterval,
		gracePeriod:  gracePeriod,
		latestEvents: make(map[string]*model.Event),
		logger:       logger.Named("event-store"),
	}
}

// Run starts runner that periodically makes the Events in the cache up-to-date
// by fetching from the control-plane.
func (s *store) Run(ctx context.Context) error {
	s.logger.Info("start running event store")

	syncTicker := time.NewTicker(s.syncInterval)
	defer syncTicker.Stop()

	// Do first sync without waiting the first ticker.
	s.milestone = time.Now().Add(-time.Hour).Unix()
	s.sync(ctx)

	for {
		select {
		case <-syncTicker.C:
			if err := s.sync(ctx); err != nil {
				s.logger.Error("failed to sync events", zap.Error(err))
			}

		case <-ctx.Done():
			s.logger.Info("event store has been stopped")
			return nil
		}
	}
}

// sync fetches a list of events newly created after its own milestone,
// and updates the cache of latest events.
func (s *store) sync(ctx context.Context) error {
	resp, err := s.apiClient.ListEvents(ctx, &pipedservice.ListEventsRequest{
		From:  s.milestone,
		Order: pipedservice.ListOrder_ASC,
	})
	if err != nil {
		return fmt.Errorf("failed to list events: %w", err)
	}
	if len(resp.Events) == 0 {
		return nil
	}

	// Eliminate events that have duplicated key.
	filtered := make(map[string]*model.Event, len(resp.Events))
	for _, e := range resp.Events {
		filtered[e.EventKey] = e
	}
	// Make the cache up-to-date.
	s.mu.Lock()
	for key, event := range filtered {
		cached, ok := s.latestEvents[key]
		if ok && cached.CreatedAt > event.CreatedAt {
			continue
		}
		s.latestEvents[key] = event
	}
	s.mu.Unlock()

	// Set the latest one within the result as the next time's "from".
	s.milestone = resp.Events[len(resp.Events)-1].CreatedAt + 1
	return nil
}

func (s *store) Getter() Getter {
	return s
}

func (s *store) GetLatest(ctx context.Context, name string, labels map[string]string, notUsingCache bool) (event *model.Event, cacheUsed bool, ok bool) {
	key := model.MakeEventKey(name, labels)

	if !notUsingCache {
		s.mu.RLock()
		event, ok = s.latestEvents[key]
		s.mu.RUnlock()
		if ok {
			return event, true, true
		}
	}

	// If not found in the cache, fetch from the control-plane.
	resp, err := s.apiClient.GetLatestEvent(ctx, &pipedservice.GetLatestEventRequest{
		Name:   name,
		Labels: labels,
	})
	if status.Code(err) == codes.NotFound {
		s.logger.Info("event not found in control-plane",
			zap.String("event-name", name),
			zap.Any("labels", labels),
		)
		return nil, false, false
	}
	if err != nil {
		s.logger.Error("failed to get the latest event", zap.Error(err))
		return nil, false, false
	}

	s.mu.Lock()
	s.latestEvents[key] = resp.Event
	s.mu.Unlock()
	return resp.Event, false, true
}
