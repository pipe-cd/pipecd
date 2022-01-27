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
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	defaultSyncInterval = time.Minute
	cacheTTL            = 24 * time.Hour
)

type Store interface {
	// Run starts syncing the event list with the control-plane.
	Run(ctx context.Context) error
	// List gives back the not-handled event list which is sorted by newest.
	ListNotHandled(name string, labels map[string]string, limit int) []*model.Event
	// GetLatest returns the latest event that meets the given conditions.
	// All objects returned here must be treated as read-only.
	GetLatest(ctx context.Context, name string, labels map[string]string) (*model.Event, bool)
	// UpdateStatuses updates the status of the latest events.
	// The second arg supposed to be the latest event. If it's not the latest, it will be ignored.
	UpdateStatuses(ctx context.Context, latestEvents []model.Event) error
}

type apiClient interface {
	GetLatestEvent(ctx context.Context, req *pipedservice.GetLatestEventRequest, opts ...grpc.CallOption) (*pipedservice.GetLatestEventResponse, error)
	ListEvents(ctx context.Context, req *pipedservice.ListEventsRequest, opts ...grpc.CallOption) (*pipedservice.ListEventsResponse, error)
	ReportEventStatuses(ctx context.Context, req *pipedservice.ReportEventStatusesRequest, opts ...grpc.CallOption) (*pipedservice.ReportEventStatusesResponse, error)
}

type store struct {
	apiClient    apiClient
	syncInterval time.Duration
	gracePeriod  time.Duration
	logger       *zap.Logger

	mu sync.RWMutex
	// A list of events that are not being handled in spite of not the latest.
	// Only events created within 24h is cached.
	// They are grouped by key, sorted by newest.
	notHandledEvents atomic.Value
	latestEvents     map[string]*model.Event
}

// NewStore creates a new event store instance.
// This syncs with the control plane to keep the list of events for this runner up-to-date.
func NewStore(apiClient apiClient, gracePeriod time.Duration, logger *zap.Logger) Store {
	s := &store{
		apiClient:    apiClient,
		syncInterval: defaultSyncInterval,
		gracePeriod:  gracePeriod,
		latestEvents: make(map[string]*model.Event),
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
		From:  time.Now().Add(-cacheTTL).Unix(),
		Order: pipedservice.ListOrder_DESC,
	})
	if err != nil {
		return fmt.Errorf("failed to list events: %w", err)
	}
	if len(resp.Events) == 0 {
		return nil
	}

	eventsByKey := make(map[string][]*model.Event, 0)
	for _, e := range resp.Events {
		eventsByKey[e.EventKey] = append(eventsByKey[e.EventKey], e)
	}
	// Update two types of cache.
	notHandledEvents := make(map[string][]*model.Event)
	latestEvents := make(map[string]*model.Event, len(eventsByKey))
	for key, es := range eventsByKey {
		latestEvents[key] = es[0]
		if len(es) == 1 {
			continue
		}
		for i := 1; i < len(es); i++ {
			if es[i].Status != model.EventStatus_EVENT_NOT_HANDLED {
				continue
			}
			notHandledEvents[key] = append(notHandledEvents[key], es[i])
		}
	}
	s.mu.Lock()
	s.latestEvents = latestEvents
	s.mu.Unlock()
	s.notHandledEvents.Store(notHandledEvents)

	return nil
}

func (s *store) ListNotHandled(name string, labels map[string]string, limit int) []*model.Event {
	if limit == 0 {
		return []*model.Event{}
	}
	key := model.MakeEventKey(name, labels)
	notHandledEvents := s.notHandledEvents.Load()
	if notHandledEvents == nil {
		return []*model.Event{}
	}
	events := notHandledEvents.(map[string][]*model.Event)[key]
	if len(events) < limit {
		return events
	}
	return events[:limit]
}

func (s *store) GetLatest(ctx context.Context, name string, labels map[string]string) (*model.Event, bool) {
	key := model.MakeEventKey(name, labels)
	s.mu.RLock()
	event, ok := s.latestEvents[key]
	s.mu.RUnlock()
	if ok {
		return event, true
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
		return nil, false
	}
	if err != nil {
		s.logger.Error("failed to get the latest event", zap.Error(err))
		return nil, false
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	cached, ok := s.latestEvents[key]
	if ok && cached.CreatedAt > event.CreatedAt {
		// Prevent overwriting even though a more recent Event has been synced immediately before.
		return resp.Event, true
	}
	s.latestEvents[key] = resp.Event
	return resp.Event, true
}

func (s *store) UpdateStatuses(ctx context.Context, latestEvents []model.Event) error {
	es := make([]*pipedservice.ReportEventStatusesRequest_Event, 0, len(latestEvents))
	for _, e := range latestEvents {
		es = append(es, &pipedservice.ReportEventStatusesRequest_Event{
			Id:                e.Id,
			Status:            e.Status,
			StatusDescription: e.StatusDescription,
		})
	}
	if _, err := s.apiClient.ReportEventStatuses(ctx, &pipedservice.ReportEventStatusesRequest{Events: es}); err != nil {
		return fmt.Errorf("failed to report event statuses: %w", err)
	}

	// TODO: Stop updating latest event cache whenever reporting
	//   and change s.latestEvents to be atomic.Value
	// Update cached events.
	for _, e := range latestEvents {
		s.mu.Lock()
		cached, ok := s.latestEvents[e.EventKey]
		if !ok {
			s.latestEvents[e.EventKey] = &e
			s.mu.Unlock()
			continue
		}
		if cached.Id != e.Id {
			// There is already an event newer than the given one.
			s.mu.Unlock()
			continue
		}
		s.latestEvents[e.EventKey].Status = e.Status
		s.latestEvents[e.EventKey].StatusDescription = e.StatusDescription
		s.mu.Unlock()
	}
	return nil
}
