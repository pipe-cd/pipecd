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

func TestGetAPIKey(t *testing.T) {
	col := &collection{kind: "APIKey"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeAPIKey := &model.APIKey{
		Id:        "get-id",
		Name:      "name",
		KeyHash:   "keyHash",
		ProjectId: "project-id",
		Role:      model.APIKey_READ_ONLY,
		Creator:   "user",
		CreatedAt: 1,
		UpdatedAt: 1,
	}
	err := client.Create(ctx, col, "get-id", fakeAPIKey)
	require.NoError(t, err)

	testcases := []struct {
		name    string
		id      string
		want    *model.APIKey
		wantErr error
	}{
		{
			name:    "entity found",
			id:      "get-id",
			want:    fakeAPIKey,
			wantErr: nil,
		},
		{
			name:    "not found",
			id:      "id-wrong",
			want:    &model.APIKey{},
			wantErr: datastore.ErrNotFound,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := &model.APIKey{}
			err := client.Get(ctx, col, tc.id, got)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFindAPIKey(t *testing.T) {
	col := &collection{kind: "APIKey"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeAPIKey := &model.APIKey{
		Id:        "find-id-1",
		Name:      "name-1",
		KeyHash:   "keyHash",
		ProjectId: "project-id",
		Role:      model.APIKey_READ_ONLY,
		Creator:   "user",
		CreatedAt: 1,
		UpdatedAt: 1,
	}
	fakeAPIKey2 := &model.APIKey{
		Id:        "find-id-2",
		Name:      "name-2",
		KeyHash:   "keyHash",
		ProjectId: "project-id",
		Role:      model.APIKey_READ_ONLY,
		Creator:   "user",
		CreatedAt: 1,
		UpdatedAt: 1,
	}
	err := client.Create(ctx, col, "find-id-1", fakeAPIKey)
	require.NoError(t, err)
	err = client.Create(ctx, col, "find-id-2", fakeAPIKey2)
	require.NoError(t, err)

	testcases := []struct {
		name    string
		opts    datastore.ListOptions
		want    []*model.APIKey
		wantErr bool
	}{
		{
			name: "no index",
			opts: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "Name",
						Operator: datastore.OperatorEqual,
						Value:    "name-1",
					},
				},
			},
			want:    []*model.APIKey{},
			wantErr: true,
		},
		{
			name: "only cursor given",
			opts: datastore.ListOptions{
				Cursor: "cursor",
			},
			want:    []*model.APIKey{},
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			it, err := client.Find(ctx, col, tc.opts)
			assert.Equal(t, tc.wantErr, err != nil)
			got, err := listAPIKeys(it)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func listAPIKeys(it datastore.Iterator) ([]*model.APIKey, error) {
	ret := make([]*model.APIKey, 0)
	if it == nil {
		return ret, nil
	}
	for {
		var v model.APIKey
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

func TestCreateAPIKey(t *testing.T) {
	col := &collection{kind: "APIKey"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeAPIKey := &model.APIKey{
		Id:        "create-id",
		Name:      "name",
		KeyHash:   "keyHash",
		ProjectId: "project-id",
		Role:      model.APIKey_READ_ONLY,
		Creator:   "user",
		CreatedAt: 1,
		UpdatedAt: 1,
	}
	err := client.Create(ctx, col, "create-id", fakeAPIKey)
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
			err := client.Create(ctx, col, tc.id, fakeAPIKey)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestUpdateAPIKey(t *testing.T) {
	col := &collection{
		kind: "APIKey",
		factory: func() interface{} {
			return &model.APIKey{}
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeAPIKey := &model.APIKey{
		Id:        "update-id",
		Name:      "name",
		KeyHash:   "keyHash",
		ProjectId: "project-id",
		Role:      model.APIKey_READ_ONLY,
		Creator:   "user",
		CreatedAt: 1,
		UpdatedAt: 1,
	}
	err := client.Create(ctx, col, "update-id", fakeAPIKey)
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
				v := e.(*model.APIKey)
				v.Name = "new-name"
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
