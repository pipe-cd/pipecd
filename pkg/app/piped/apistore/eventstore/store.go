// Copyright 2023 The PipeCD Authors.
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
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	defaultSyncInterval = time.Minute
	cacheTTL            = 24 * time.Hour
)

type Lister interface {
	// List gives back the not-handled event list which is sorted by newest.
	ListNotHandled(name string, labels map[string]string, minCreatedAt int64, limit int) []*model.Event
}

type Store interface {
	// Run starts syncing the event list with the control-plane.
	Run(ctx context.Context) error
	Lister() Lister
}

type apiClient interface {
	ListEvents(ctx context.Context, req *pipedservice.ListEventsRequest, opts ...grpc.CallOption) (*pipedservice.ListEventsResponse, error)
}

type store struct {
	apiClient    apiClient
	syncInterval time.Duration
	gracePeriod  time.Duration
	logger       *zap.Logger

	// A list of events that are not being handled in spite of not the latest.
	// Only events created within 24h is cached.
	// They are grouped by key, sorted by newest.
	notHandledEvents atomic.Value
}

// NewStore creates a new event store instance.
// This syncs with the control plane to keep the list of events for this runner up-to-date.
func NewStore(apiClient apiClient, gracePeriod time.Duration, logger *zap.Logger) Store {
	s := &store{
		apiClient:    apiClient,
		syncInterval: defaultSyncInterval,
		gracePeriod:  gracePeriod,
		logger:       logger.Named("event-store"),
	}
	s.notHandledEvents.Store(make(map[string][]*model.Event))
	return s
}

// Run starts runner that periodically makes the Events in the cache up-to-date
// by fetching from the control-plane.
func (s *store) Run(ctx context.Context) error {
	s.logger.Info("start running event store")

	syncTicker := time.NewTicker(s.syncInterval)
	defer syncTicker.Stop()

	// Do first sync without waiting the first ticker.
	if err := s.sync(ctx); err != nil {
		return fmt.Errorf("failed to sync events first time: %w", err)
	}

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
	// Fetch events in descending order.
	resp, err := s.apiClient.ListEvents(ctx, &pipedservice.ListEventsRequest{
		From:   time.Now().Add(-cacheTTL).Unix(),
		Order:  pipedservice.ListOrder_DESC,
		Status: pipedservice.ListEventsRequest_NOT_HANDLED,
	})
	if err != nil {
		return fmt.Errorf("failed to list events: %w", err)
	}
	if len(resp.Events) == 0 {
		s.notHandledEvents.Store(make(map[string][]*model.Event))
		return nil
	}

	notHandledEvents := make(map[string][]*model.Event)
	for _, e := range resp.Events {
		notHandledEvents[e.EventKey] = append(notHandledEvents[e.EventKey], e)
	}
	s.notHandledEvents.Store(notHandledEvents)

	return nil
}

func (s *store) ListNotHandled(name string, labels map[string]string, minCreatedAt int64, limit int) []*model.Event {
	if limit == 0 {
		return []*model.Event{}
	}
	key := model.MakeEventKey(name, labels)
	notHandledEvents := s.notHandledEvents.Load()
	if notHandledEvents == nil {
		return []*model.Event{}
	}

	events := notHandledEvents.(map[string][]*model.Event)[key]
	out := make([]*model.Event, 0, len(events))
	for i, e := range events {
		if e.CreatedAt < minCreatedAt {
			break
		}
		if i >= limit {
			break
		}
		out = append(out, e)
	}
	return out
}

func (s *store) Lister() Lister {
	return s
}
