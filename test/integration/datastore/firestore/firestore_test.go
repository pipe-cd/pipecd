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

package firestore

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/datastore"
)

type Entity struct {
	Name string
}

type collection struct {
	kind    string
	factory datastore.Factory
}

func (c *collection) Kind() string {
	return c.kind
}

func (c *collection) Factory() datastore.Factory {
	return c.factory
}

func TestGet(t *testing.T) {
	col := &collection{kind: "GetEntity"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := store.Create(ctx, col, "id", &Entity{Name: "name"})
	require.NoError(t, err)

	testcases := []struct {
		name    string
		id      string
		want    *Entity
		wantErr error
	}{
		{
			name:    "entity found",
			id:      "id",
			want:    &Entity{Name: "name"},
			wantErr: nil,
		},
		{
			name:    "not found",
			id:      "id-wrong",
			want:    &Entity{},
			wantErr: datastore.ErrNotFound,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := &Entity{}
			err := store.Get(ctx, col, tc.id, got)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFind(t *testing.T) {
	col := &collection{kind: "FindEntity"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := store.Create(ctx, col, "id-1", &Entity{Name: "name-1"})
	require.NoError(t, err)
	err = store.Create(ctx, col, "id-2", &Entity{Name: "name-2"})
	require.NoError(t, err)

	testcases := []struct {
		name    string
		opts    datastore.ListOptions
		want    []*Entity
		wantErr bool
	}{
		{
			name: "fetch all",
			want: []*Entity{
				{
					Name: "name-1",
				},
				{
					Name: "name-2",
				},
			},
			wantErr: false,
		},
		{
			name: "fetch by name",
			opts: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "Name",
						Operator: datastore.OperatorEqual,
						Value:    "name-1",
					},
				},
			},
			want: []*Entity{
				{
					Name: "name-1",
				},
			},
			wantErr: false,
		},
		{
			name: "only cursor given",
			opts: datastore.ListOptions{
				Cursor: "cursor",
			},
			want:    []*Entity{},
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			it, err := store.Find(ctx, col, tc.opts)
			assert.Equal(t, tc.wantErr, err != nil)
			got, err := listEntities(it)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func listEntities(it datastore.Iterator) ([]*Entity, error) {
	entity := make([]*Entity, 0)
	if it == nil {
		return entity, nil
	}
	for {
		var e Entity
		err := it.Next(&e)
		if errors.Is(err, datastore.ErrIteratorDone) {
			break
		}
		if err != nil {
			return nil, err
		}
		entity = append(entity, &e)
	}
	return entity, nil
}

func TestCreate(t *testing.T) {
	col := &collection{kind: "CreateEntity"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := store.Create(ctx, col, "id", &Entity{Name: "name"})
	require.NoError(t, err)

	testcases := []struct {
		name    string
		id      string
		wantErr error
	}{
		{
			name:    "already exists",
			id:      "id",
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
			err := store.Create(ctx, col, tc.id, &Entity{Name: "name"})
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestUpdate(t *testing.T) {
	col := &collection{
		kind: "UpdateEntity",
		factory: func() interface{} {
			return &Entity{}
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := store.Create(ctx, col, "id", &Entity{Name: "name"})
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
			id:   "id",
			updater: func(interface{}) error {
				return fmt.Errorf("error")
			},
			wantErr: fmt.Errorf("error"),
		},
		{
			name: "successful update",
			id:   "id",
			updater: func(e interface{}) error {
				entity := e.(*Entity)
				entity.Name = "new-name"
				return nil
			},
			wantErr: nil,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := store.Update(ctx, col, tc.id, tc.updater)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
