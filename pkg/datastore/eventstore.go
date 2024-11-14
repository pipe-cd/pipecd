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

package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pipe-cd/pipecd/pkg/model"
)

type eventCollection struct {
	requestedBy Commander
}

func (e *eventCollection) Kind() string {
	return "Event"
}

func (e *eventCollection) Factory() Factory {
	return func() interface{} {
		return &model.Event{}
	}
}

func (e *eventCollection) ListInUsedShards() []Shard {
	return []Shard{
		AgentShard,
	}
}

func (e *eventCollection) GetUpdatableShard() (Shard, error) {
	switch e.requestedBy {
	case PipedCommander:
		return AgentShard, nil
	default:
		return "", ErrUnsupported
	}
}

func (e *eventCollection) Encode(entity interface{}) (map[Shard][]byte, error) {
	const errFmt = "failed while encode Event object: %s"

	me, ok := entity.(*model.Event)
	if !ok {
		return nil, fmt.Errorf(errFmt, "type not matched")
	}

	data, err := json.Marshal(me)
	if err != nil {
		return nil, fmt.Errorf(errFmt, "unable to marshal entity data")
	}
	return map[Shard][]byte{
		AgentShard: data,
	}, nil
}

type EventStore interface {
	Add(ctx context.Context, e model.Event) error
	List(ctx context.Context, opts ListOptions) ([]*model.Event, string, error)
	UpdateStatus(ctx context.Context, eventID string, status model.EventStatus, statusDescription string) error
}

type eventStore struct {
	backend
	commander Commander
	nowFunc   func() time.Time
}

func NewEventStore(ds DataStore, c Commander) EventStore {
	return &eventStore{
		backend: backend{
			ds:  ds,
			col: &eventCollection{requestedBy: c},
		},
		commander: c,
		nowFunc:   time.Now,
	}
}

func (s *eventStore) Add(ctx context.Context, e model.Event) error {
	now := s.nowFunc().Unix()
	if e.CreatedAt == 0 {
		e.CreatedAt = now
	}
	if e.UpdatedAt == 0 {
		e.UpdatedAt = now
	}
	if err := e.Validate(); err != nil {
		return fmt.Errorf("failed to validate event: %w: %w", ErrInvalidArgument, err)
	}
	return s.ds.Create(ctx, s.col, e.Id, &e)
}

func (s *eventStore) List(ctx context.Context, opts ListOptions) ([]*model.Event, string, error) {
	it, err := s.ds.Find(ctx, s.col, opts)
	if err != nil {
		return nil, "", err
	}
	es := make([]*model.Event, 0)
	for {
		var e model.Event
		err := it.Next(&e)
		if err == ErrIteratorDone {
			break
		}
		if err != nil {
			return nil, "", err
		}
		es = append(es, &e)
	}

	// In case there is no more elements found, cursor should be set to empty too.
	if len(es) == 0 {
		return es, "", nil
	}
	cursor, err := it.Cursor()
	if err != nil {
		return nil, "", err
	}
	return es, cursor, nil
}

func (s *eventStore) UpdateStatus(ctx context.Context, eventID string, status model.EventStatus, statusDescription string) error {
	return s.ds.Update(ctx, s.col, eventID, func(e interface{}) error {
		event := e.(*model.Event)
		event.Status = status
		event.StatusDescription = statusDescription
		if event.IsHandled() {
			now := s.nowFunc().Unix()
			event.HandledAt = now
			event.UpdatedAt = now
		}
		return nil
	})
}
