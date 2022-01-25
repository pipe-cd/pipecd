// Copyright 2022 The PipeCD Authors.
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

// Package eventstatusupdater provides an ability to periodically update the status of outdated events.
package eventstatusupdater

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	interval = time.Hour
	// A month.
	durationToFetchAtStartUp = 720 * time.Hour
)

type Updater struct {
	// Mark that it has handled all events that was created before this UNIX time.
	milestone int64

	eventStore datastore.EventStore
	logger     *zap.Logger
}

func NewUpdater(ds datastore.DataStore, logger *zap.Logger) *Updater {
	return &Updater{
		milestone:  time.Now().Add(-durationToFetchAtStartUp).Unix(),
		eventStore: datastore.NewEventStore(ds),
		logger:     logger,
	}
}

func (u *Updater) Run(ctx context.Context) error {
	u.logger.Info("start running EventStatusUpdater")

	t := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			u.logger.Info("EventStatusUpdater has been stopped")
			return nil

		case <-t.C:
			num, err := u.makeEventsOutdated(ctx)
			if err != nil {
				u.logger.Error("failed to make events outdated", zap.Error(err))
				continue
			}
			if num > 0 {
				u.logger.Info(fmt.Sprintf("successfully made %d events outdated", num))
			}
		}
	}
}

func (u *Updater) makeEventsOutdated(ctx context.Context) (int, error) {
	opts := datastore.ListOptions{
		Filters: []datastore.ListFilter{
			{
				Field:    "CreatedAt",
				Operator: datastore.OperatorGreaterThanOrEqual,
				Value:    u.milestone,
			},
			{
				Field:    "Status",
				Operator: datastore.OperatorNotEqual,
				Value:    model.EventStatus_EVENT_OUTDATED,
			},
		},
		Orders: []datastore.Order{
			{
				Field:     "CreatedAt",
				Direction: datastore.Asc,
			},
			{
				Field:     "Id",
				Direction: datastore.Asc,
			},
		},
	}
	events, _, err := u.eventStore.ListEvents(ctx, opts)
	if err != nil {
		return 0, fmt.Errorf("failed to list events: %w", err)
	}
	if len(events) == 0 {
		return 0, nil
	}

	updatedNum := 0
	// Eliminate events that have duplicated key.
	groupedEvents := make(map[string][]*model.Event, len(events))
	for _, e := range events {
		groupedEvents[e.EventKey] = append(groupedEvents[e.EventKey], e)
	}
	for _, es := range groupedEvents {
		// Start updating other than the latest one.
		for i := 0; i < len(es)-1; i++ {
			if es[i].Status != model.EventStatus_EVENT_NOT_HANDLED {
				continue
			}
			err := u.eventStore.UpdateEventStatus(ctx, es[i].Id, model.EventStatus_EVENT_OUTDATED, "The new event has been created")
			if err != nil {
				return 0, fmt.Errorf("failed to update event %q", es[i].Id)
			}
			updatedNum++
		}
	}

	// Set the latest one within the result as the next time's "from".
	u.milestone = events[len(events)-1].CreatedAt + 1
	return updatedNum, nil
}
