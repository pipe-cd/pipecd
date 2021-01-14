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
	"sort"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipe/pkg/app/api/service/pipedservice"
	"github.com/pipe-cd/pipe/pkg/model"
)

// Getter helps get an event. All objects returned here must be treated as read-only.
type Getter interface {
	// GetLatest returns the latest event that meets the given conditions.
	GetLatest(ctx context.Context, name string, labels map[string]string) (*model.Event, bool)
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
	latestEventMap map[string]*model.Event
}

const (
	defaultSyncInterval = time.Minute
)

// NewStore creates a new event store instance.
// This syncs with the control plane to keep the list of events for this runner up-to-date.
func NewStore(apiClient apiClient, gracePeriod time.Duration, logger *zap.Logger) Store {
	return &store{
		apiClient:      apiClient,
		syncInterval:   defaultSyncInterval,
		gracePeriod:    gracePeriod,
		latestEventMap: make(map[string]*model.Event),
		logger:         logger.Named("event-store"),
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

// sync fetches a list of events newly created after its own milestone.
func (s *store) sync(ctx context.Context) error {
	resp, err := s.apiClient.ListEvents(ctx, &pipedservice.ListEventsRequest{
		From:  s.milestone,
		Order: pipedservice.ListOrder_ASC,
	})
	if err != nil {
		return fmt.Errorf("failed to list events: %w", err)
	}

	// Make the cache up-to-date by traversing events sorted by oldest first.
	var latestTime int64
	s.mu.Lock()
	for _, event := range resp.Events {
		id := eventDefinitionID(event.Name, event.Labels)
		s.latestEventMap[id] = event
		latestTime = event.CreatedAt
	}
	s.mu.Unlock()

	// Set the latest one within the result as the next time's "from".
	s.milestone = latestTime + 1
	return nil
}

func (s *store) Getter() Getter {
	return s
}

func (s *store) GetLatest(ctx context.Context, name string, labels map[string]string) (*model.Event, bool) {
	id := eventDefinitionID(name, labels)
	s.mu.RLock()
	event, ok := s.latestEventMap[id]
	s.mu.RUnlock()
	if ok {
		return event, true
	}

	// If not found in the cache, fetch from the control-plane.
	resp, err := s.apiClient.GetLatestEvent(ctx, &pipedservice.GetLatestEventRequest{
		Name:   name,
		Labels: labels,
	})
	if err != nil {
		s.logger.Error("failed to get the latest event", zap.Error(err))
		return nil, false
	}

	s.mu.Lock()
	s.latestEventMap[id] = resp.Event
	s.mu.Unlock()
	return resp.Event, true
}

// eventDefinitionID builds a unique identifier based on the given name and labels.
// It returns the exact same string as long as both are the same.
func eventDefinitionID(name string, labels map[string]string) string {
	if len(labels) == 0 {
		return name
	}

	var b strings.Builder
	b.WriteString(name)

	// Guarantee uniqueness by sorting by keys.
	keys := make([]string, 0, len(labels))
	for key := range labels {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		b.WriteString(fmt.Sprintf("/%s:%s", key, labels[key]))
	}
	return b.String()
}
