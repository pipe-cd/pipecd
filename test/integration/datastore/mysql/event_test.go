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

package mysql

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestGetEvent(t *testing.T) {
	col := &collection{kind: "Event"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeEvent := &model.Event{
		Id:        "get-id",
		Name:      "name",
		Data:      "data",
		ProjectId: "project-id",
		EventKey:  "event-key",
		CreatedAt: 1,
		UpdatedAt: 1,
	}
	err := client.Create(ctx, col, "get-id", fakeEvent)
	require.NoError(t, err)

	testcases := []struct {
		name    string
		id      string
		want    *model.Event
		wantErr error
	}{
		{
			name:    "entity found",
			id:      "get-id",
			want:    fakeEvent,
			wantErr: nil,
		},
		{
			name:    "not found",
			id:      "id-wrong",
			want:    &model.Event{},
			wantErr: datastore.ErrNotFound,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := &model.Event{}
			err := client.Get(ctx, col, tc.id, got)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFindEvent(t *testing.T) {
	col := &collection{kind: "Event"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeEvent := &model.Event{
		Id:        "find-id-1",
		Name:      "name",
		Data:      "data",
		ProjectId: "find-project-id",
		EventKey:  "event-key",
		CreatedAt: 1,
		UpdatedAt: 1,
	}
	fakeEvent2 := &model.Event{
		Id:        "find-id-2",
		Name:      "name",
		Data:      "data",
		ProjectId: "find-project-id",
		EventKey:  "event-key",
		CreatedAt: 1,
		UpdatedAt: 1,
	}
	err := client.Create(ctx, col, "find-id-1", fakeEvent)
	require.NoError(t, err)
	err = client.Create(ctx, col, "find-id-2", fakeEvent2)
	require.NoError(t, err)

	testcases := []struct {
		name    string
		opts    datastore.ListOptions
		want    []*model.Event
		wantErr bool
	}{
		{
			name: "fetch by project_id",
			opts: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "ProjectId",
						Operator: datastore.OperatorEqual,
						Value:    "find-project-id",
					},
				},
			},
			want: []*model.Event{
				fakeEvent,
				fakeEvent2,
			},
			wantErr: false,
		},
		{
			name: "only cursor given",
			opts: datastore.ListOptions{
				Cursor: "cursor",
			},
			want:    []*model.Event{},
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			it, err := client.Find(ctx, col, tc.opts)
			assert.Equal(t, tc.wantErr, err != nil)
			got, err := listEvents(it)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func listEvents(it datastore.Iterator) ([]*model.Event, error) {
	ret := make([]*model.Event, 0)
	if it == nil {
		return ret, nil
	}
	for {
		var v model.Event
		err := it.Next(&v)
		if errors.Is(err, datastore.ErrIteratorDone) {
			break
		}
		if err != nil {
			return nil, err
		}
		ret = append(ret, &v)
	}
	return ret, nil
}

func TestCreateEvent(t *testing.T) {
	col := &collection{kind: "Event"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeEvent := &model.Event{
		Id:        "create-id",
		Name:      "name",
		Data:      "data",
		ProjectId: "project-id",
		EventKey:  "event-key",
		CreatedAt: 1,
		UpdatedAt: 1,
	}
	err := client.Create(ctx, col, "create-id", fakeEvent)
	require.NoError(t, err)

	testcases := []struct {
		name    string
		id      string
		wantErr error
	}{
		{
			name:    "already exists",
			id:      "create-id",
			wantErr: datastore.ErrAlreadyExists,
		},
		{
			name:    "successful create",
			id:      "id-new",
			wantErr: nil,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := client.Create(ctx, col, tc.id, fakeEvent)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestUpdateEvent(t *testing.T) {
	col := &collection{
		kind: "Event",
		factory: func() interface{} {
			return &model.Event{}
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeEvent := &model.Event{
		Id:        "update-id",
		Name:      "name",
		Data:      "data",
		ProjectId: "project-id",
		EventKey:  "event-key",
		CreatedAt: 1,
		UpdatedAt: 1,
	}
	err := client.Create(ctx, col, "update-id", fakeEvent)
	require.NoError(t, err)

	testcases := []struct {
		name    string
		id      string
		updater func(interface{}) error
		wantErr error
	}{
		{
			name:    "not found",
			id:      "id-wrong",
			wantErr: datastore.ErrNotFound,
		},
		{
			name: "unable to update",
			id:   "update-id",
			updater: func(interface{}) error {
				return fmt.Errorf("error")
			},
			wantErr: fmt.Errorf("error"),
		},
		{
			name: "successful update",
			id:   "update-id",
			updater: func(e interface{}) error {
				v := e.(*model.Event)
				v.Status = model.EventStatus_EVENT_SUCCESS
				return nil
			},
			wantErr: nil,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := client.Update(ctx, col, tc.id, tc.updater)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
