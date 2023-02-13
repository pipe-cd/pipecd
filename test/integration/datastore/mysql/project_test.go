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

func TestGetProject(t *testing.T) {
	col := &collection{kind: "Project"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeProject := &model.Project{
		Id: "get-id",
		StaticAdmin: &model.ProjectStaticUser{
			Username:     "username",
			PasswordHash: "password-hash",
		},
		CreatedAt: 1,
		UpdatedAt: 1,
	}
	err := client.Create(ctx, col, "get-id", fakeProject)
	require.NoError(t, err)

	testcases := []struct {
		name    string
		id      string
		want    *model.Project
		wantErr error
	}{
		{
			name:    "entity found",
			id:      "get-id",
			want:    fakeProject,
			wantErr: nil,
		},
		{
			name:    "not found",
			id:      "id-wrong",
			want:    &model.Project{},
			wantErr: datastore.ErrNotFound,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := &model.Project{}
			err := client.Get(ctx, col, tc.id, got)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFindProject(t *testing.T) {
	col := &collection{kind: "Project"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeProject := &model.Project{
		Id: "find-id-1",
		StaticAdmin: &model.ProjectStaticUser{
			Username:     "username",
			PasswordHash: "password-hash",
		},
		CreatedAt: 1,
		UpdatedAt: 1,
	}
	fakeProject2 := &model.Project{
		Id: "find-id-2",
		StaticAdmin: &model.ProjectStaticUser{
			Username:     "username",
			PasswordHash: "password-hash",
		},
		CreatedAt: 2,
		UpdatedAt: 2,
	}
	err := client.Create(ctx, col, "find-id-1", fakeProject)
	require.NoError(t, err)
	err = client.Create(ctx, col, "find-id-2", fakeProject2)
	require.NoError(t, err)

	testcases := []struct {
		name    string
		opts    datastore.ListOptions
		want    []*model.Project
		wantErr bool
	}{
		{
			name: "no index",
			opts: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "UpdateAt",
						Operator: datastore.OperatorEqual,
						Value:    1,
					},
				},
			},
			want:    []*model.Project{},
			wantErr: true,
		},
		{
			name: "only cursor given",
			opts: datastore.ListOptions{
				Cursor: "cursor",
			},
			want:    []*model.Project{},
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			it, err := client.Find(ctx, col, tc.opts)
			assert.Equal(t, tc.wantErr, err != nil)
			got, err := listProjects(it)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func listProjects(it datastore.Iterator) ([]*model.Project, error) {
	ret := make([]*model.Project, 0)
	if it == nil {
		return ret, nil
	}
	for {
		var v model.Project
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

func TestCreateProject(t *testing.T) {
	col := &collection{kind: "Project"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeProject := &model.Project{
		Id: "create-id",
		StaticAdmin: &model.ProjectStaticUser{
			Username:     "username",
			PasswordHash: "password-hash",
		},
		CreatedAt: 1,
		UpdatedAt: 1,
	}
	err := client.Create(ctx, col, "create-id", fakeProject)
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
			err := client.Create(ctx, col, tc.id, fakeProject)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestUpdateProject(t *testing.T) {
	col := &collection{
		kind: "Project",
		factory: func() interface{} {
			return &model.Project{}
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeProject := &model.Project{
		Id: "update-id",
		StaticAdmin: &model.ProjectStaticUser{
			Username:     "username",
			PasswordHash: "password-hash",
		},
		CreatedAt: 1,
		UpdatedAt: 1,
	}
	err := client.Create(ctx, col, "update-id", fakeProject)
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
				v := e.(*model.Project)
				v.Desc = "new-desc"
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
