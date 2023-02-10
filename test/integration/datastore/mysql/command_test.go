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

func TestGetCommand(t *testing.T) {
	col := &collection{kind: "Command"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeCommand := &model.Command{
		Id:            "get-id",
		ApplicationId: "app-id",
		PipedId:       "piped-id",
		ProjectId:     "project-id",
		CreatedAt:     1,
		UpdatedAt:     1,
	}
	err := client.Create(ctx, col, "get-id", fakeCommand)
	require.NoError(t, err)

	testcases := []struct {
		name    string
		id      string
		want    *model.Command
		wantErr error
	}{
		{
			name:    "entity found",
			id:      "get-id",
			want:    fakeCommand,
			wantErr: nil,
		},
		{
			name:    "not found",
			id:      "id-wrong",
			want:    &model.Command{},
			wantErr: datastore.ErrNotFound,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := &model.Command{}
			err := client.Get(ctx, col, tc.id, got)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFindCommand(t *testing.T) {
	col := &collection{kind: "Command"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeCommand := &model.Command{
		Id:            "find-id-1",
		ApplicationId: "app-id",
		PipedId:       "find-piped-id",
		ProjectId:     "project-id",
		CreatedAt:     1,
		UpdatedAt:     1,
	}
	fakeCommand2 := &model.Command{
		Id:            "find-id-2",
		ApplicationId: "app-id",
		PipedId:       "find-piped-id",
		ProjectId:     "project-id",
		CreatedAt:     1,
		UpdatedAt:     1,
	}
	err := client.Create(ctx, col, "find-id-1", fakeCommand)
	require.NoError(t, err)
	err = client.Create(ctx, col, "find-id-2", fakeCommand2)
	require.NoError(t, err)

	testcases := []struct {
		name    string
		opts    datastore.ListOptions
		want    []*model.Command
		wantErr bool
	}{
		{
			name: "fetch by piped_id",
			opts: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "PipedId",
						Operator: datastore.OperatorEqual,
						Value:    "find-piped-id",
					},
				},
			},
			want: []*model.Command{
				fakeCommand,
				fakeCommand2,
			},
			wantErr: false,
		},
		{
			name: "only cursor given",
			opts: datastore.ListOptions{
				Cursor: "cursor",
			},
			want:    []*model.Command{},
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			it, err := client.Find(ctx, col, tc.opts)
			assert.Equal(t, tc.wantErr, err != nil)
			got, err := listCommands(it)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func listCommands(it datastore.Iterator) ([]*model.Command, error) {
	ret := make([]*model.Command, 0)
	if it == nil {
		return ret, nil
	}
	for {
		var v model.Command
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

func TestCreateCommand(t *testing.T) {
	col := &collection{kind: "Command"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeCommand := &model.Command{
		Id:            "create-id",
		ApplicationId: "app-id",
		PipedId:       "piped-id",
		ProjectId:     "project-id",
		CreatedAt:     1,
		UpdatedAt:     1,
	}
	err := client.Create(ctx, col, "create-id", fakeCommand)
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
			err := client.Create(ctx, col, tc.id, fakeCommand)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestUpdateCommand(t *testing.T) {
	col := &collection{
		kind: "Command",
		factory: func() interface{} {
			return &model.Command{}
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeCommand := &model.Command{
		Id:            "update-id",
		ApplicationId: "app-id",
		PipedId:       "piped-id",
		ProjectId:     "project-id",
		CreatedAt:     1,
		UpdatedAt:     1,
	}
	err := client.Create(ctx, col, "update-id", fakeCommand)
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
				v := e.(*model.Command)
				v.Status = model.CommandStatus_COMMAND_SUCCEEDED
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
