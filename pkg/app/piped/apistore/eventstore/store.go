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
	"sync/atomic"
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

	// Mark that it has handled all events that was created before this value.
	milestone int64
	// The key is supposed to be event-definition ID, a string consists of name and labels.
	// And the value is the address to the latest Event.
	latestEventMap atomic.Value
}

const (
	defaultSyncInterval = time.Minute
)

// NewStore creates a new event store instance.
// This syncs with the control plane to keep the list of events for this runner up-to-date.
func NewStore(apiClient apiClient, gracePeriod time.Duration, logger *zap.Logger) Store {
	s := &store{
		apiClient:    apiClient,
		syncInterval: defaultSyncInterval,
		gracePeriod:  gracePeriod,
		logger:       logger.Named("event-store"),
	}
	s.latestEventMap.Store(make(map[string]*model.Event))
	return s
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
			s.sync(ctx)

		case <-ctx.Done():
			s.logger.Info("event store has been stopped")
			return nil
		}
	}
}

// sync fetches a list of events inside the range between own milestone and the current local time.
// Only this function takes responsibility for updating the cache.
func (s *store) sync(ctx context.Context) error {
	// TODO: Use UTC and let control-plane convert it into the control-plane's Local time
	// Unexpected behavior can be happened if Piped uses different timezone from the control-plane.
	to := time.Now().Unix()
	resp, err := s.apiClient.ListEvents(ctx, &pipedservice.ListEventsRequest{
		From:  s.milestone,
		To:    to,
		Order: pipedservice.ListEventsOrder_FROM_OLDEST,
	})
	if err != nil {
		s.logger.Error("failed to list events", zap.Error(err))
		return err
	}
	// Set this time's "to" as the next time's "from".
	s.milestone = to

	// Make the cache up-to-date by traversing events sorted by oldest first.
	eventMap := s.latestEventMap.Load().(map[string]*model.Event)
	for _, event := range resp.Events {
		id := model.EventDefinitionID(event.Name, event.Labels)
		eventMap[id] = event
	}
	s.latestEventMap.Store(eventMap)

	return nil
}

func (s *store) Getter() Getter {
	return s
}

func (s *store) GetLatest(ctx context.Context, name string, labels map[string]string) (*model.Event, bool) {
	eventMap := s.latestEventMap.Load().(map[string]*model.Event)
	id := model.EventDefinitionID(name, labels)
	event, ok := eventMap[id]
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

	// NOTE: Don't update the cache to prevent it from being overwritten by slightly older ones.
	return resp.Event, true
}
