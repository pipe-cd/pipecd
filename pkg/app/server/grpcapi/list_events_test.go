// Copyright 2026 The PipeCD Authors.
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

package grpcapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	service "github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type fakeEventStore struct {
	pipedAPIEventStore
	pages   [][]*model.Event
	gotOpts []datastore.ListOptions
}

func (f *fakeEventStore) List(ctx context.Context, opts datastore.ListOptions) ([]*model.Event, string, error) {
	f.gotOpts = append(f.gotOpts, opts)
	if len(f.pages) == 0 {
		return nil, "", nil
	}
	page := f.pages[0]
	f.pages = f.pages[1:]
	cursor := ""
	if len(f.pages) > 0 {
		cursor = "next"
	}
	return page, cursor, nil
}

func TestListEventsPagination(t *testing.T) {
	ctx := pipedAuthContext(t)

	t.Run("defaults to newest first when order is unset", func(t *testing.T) {
		store := &fakeEventStore{
			pages: [][]*model.Event{
				{{Id: "event-1"}},
			},
		}

		api := &PipedAPI{eventStore: store}
		resp, err := api.ListEvents(ctx, &service.ListEventsRequest{})
		assert.NoError(t, err)
		assert.Len(t, resp.Events, 1)

		opts := store.gotOpts[0]
		assert.NotEmpty(t, opts.Orders, "paging requires a stable order")
		assert.Equal(t, datastore.Desc, opts.Orders[0].Direction)
		assert.Equal(t, listEventsPageSize, opts.Limit)
	})

	t.Run("unknown order value is rejected", func(t *testing.T) {
		store := &fakeEventStore{}

		api := &PipedAPI{eventStore: store}
		resp, err := api.ListEvents(ctx, &service.ListEventsRequest{Order: service.ListOrder(99)})
		assert.Nil(t, resp)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("multiple pages are aggregated", func(t *testing.T) {
		store := &fakeEventStore{
			pages: [][]*model.Event{
				{{Id: "event-1"}},
				{{Id: "event-2"}},
			},
		}

		api := &PipedAPI{eventStore: store}
		resp, err := api.ListEvents(ctx, &service.ListEventsRequest{Order: service.ListOrder_ASC})
		assert.NoError(t, err)
		assert.Len(t, resp.Events, 2)
		assert.Equal(t, "event-1", resp.Events[0].Id)
		assert.Equal(t, "event-2", resp.Events[1].Id)

		// First call must not carry a cursor; the second must continue from
		// the cursor returned by the first.
		assert.Empty(t, store.gotOpts[0].Cursor)
		assert.NotEmpty(t, store.gotOpts[1].Cursor)
	})
}
