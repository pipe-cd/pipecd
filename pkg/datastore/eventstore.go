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

package datastore

import (
	"context"
	"time"

	"github.com/pipe-cd/pipecd/pkg/model"
)

const EventModelKind = "Event"

var (
	eventFactory = func() interface{} {
		return &model.Event{}
	}
)

type EventStore interface {
	AddEvent(ctx context.Context, e model.Event) error
	ListEvents(ctx context.Context, opts ListOptions) ([]*model.Event, error)
	UpdateEventStatus(ctx context.Context, eventID string, status model.EventStatus, statusDescription string) error
}

type eventStore struct {
	backend
	nowFunc func() time.Time
}

func NewEventStore(ds DataStore) EventStore {
	return &eventStore{
		backend: backend{
			ds: ds,
		},
		nowFunc: time.Now,
	}
}

func (s *eventStore) AddEvent(ctx context.Context, e model.Event) error {
	now := s.nowFunc().Unix()
	if e.CreatedAt == 0 {
		e.CreatedAt = now
	}
	if e.UpdatedAt == 0 {
		e.UpdatedAt = now
	}
	if err := e.Validate(); err != nil {
		return err
	}
	return s.ds.Create(ctx, EventModelKind, e.Id, &e)
}

func (s *eventStore) ListEvents(ctx context.Context, opts ListOptions) ([]*model.Event, error) {
	it, err := s.ds.Find(ctx, EventModelKind, opts)
	if err != nil {
		return nil, err
	}
	es := make([]*model.Event, 0)
	for {
		var e model.Event
		err := it.Next(&e)
		if err == ErrIteratorDone {
			break
		}
		if err != nil {
			return nil, err
		}
		es = append(es, &e)
	}
	return es, nil
}

func (s *eventStore) UpdateEventStatus(ctx context.Context, eventID string, status model.EventStatus, statusDescription string) error {
	return s.ds.Update(ctx, EventModelKind, eventID, eventFactory, func(e interface{}) error {
		event := e.(*model.Event)
		event.Status = status
		event.StatusDescription = statusDescription
		if status == model.EventStatus_EVENT_SUCCESS || status == model.EventStatus_EVENT_FAILURE {
			event.Handled = true
			event.HandledAt = s.nowFunc().Unix()
		}
		return nil
	})
}
