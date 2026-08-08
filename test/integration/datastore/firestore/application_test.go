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

package firestore

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestGetApplication(t *testing.T) {
	col := &collection{kind: "Application"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeApplication := &model.Application{
		Id:        "get-id",
		Name:      "name",
		PipedId:   "piped-id",
		ProjectId: "project-id",
		Kind:      model.ApplicationKind_KUBERNETES,
		GitPath: &model.ApplicationGitPath{
			Repo: &model.ApplicationGitRepository{Id: "id"},
			Path: "path",
		},
		CloudProvider: "cloud-provider",
		CreatedAt:     1,
		UpdatedAt:     1,
	}
	err := store.Create(ctx, col, "get-id", fakeApplication)
	require.NoError(t, err)

	testcases := []struct {
		name    string
		id      string
		want    *model.Application
		wantErr error
	}{
		{
			name:    "entity found",
			id:      "get-id",
			want:    fakeApplication,
			wantErr: nil,
		},
		{
			name:    "not found",
			id:      "id-wrong",
			want:    &model.Application{},
			wantErr: datastore.ErrNotFound,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := &model.Application{}
			err := store.Get(ctx, col, tc.id, got)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCreateApplication(t *testing.T) {
	col := &collection{kind: "Application"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeApplication := &model.Application{
		Id:        "create-id",
		Name:      "name",
		PipedId:   "piped-id",
		ProjectId: "project-id",
		Kind:      model.ApplicationKind_KUBERNETES,
		GitPath: &model.ApplicationGitPath{
			Repo: &model.ApplicationGitRepository{Id: "id"},
			Path: "path",
		},
		CloudProvider: "cloud-provider",
		CreatedAt:     1,
		UpdatedAt:     1,
	}
	err := store.Create(ctx, col, "create-id", fakeApplication)
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
			err := store.Create(ctx, col, tc.id, fakeApplication)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestUpdateApplication(t *testing.T) {
	col := &collection{
		kind: "Application",
		factory: func() interface{} {
			return &model.Application{}
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeApplication := &model.Application{
		Id:        "update-id",
		Name:      "name",
		PipedId:   "piped-id",
		ProjectId: "project-id",
		Kind:      model.ApplicationKind_KUBERNETES,
		GitPath: &model.ApplicationGitPath{
			Repo: &model.ApplicationGitRepository{Id: "id"},
			Path: "path",
		},
		CloudProvider: "cloud-provider",
		CreatedAt:     1,
		UpdatedAt:     1,
	}
	err := store.Create(ctx, col, "update-id", fakeApplication)
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
				v := e.(*model.Application)
				v.Name = "new-name"
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
