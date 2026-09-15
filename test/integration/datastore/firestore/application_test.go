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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestFindApplication(t *testing.T) {
	col := &collection{kind: "Application"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fakeApplication := &model.Application{
		Id:        "find-id-1",
		Name:      "name-1",
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
	fakeApplication2 := &model.Application{
		Id:        "find-id-2",
		Name:      "name-2",
		PipedId:   "piped-id",
		ProjectId: "project-id",
		Kind:      model.ApplicationKind_KUBERNETES,
		GitPath: &model.ApplicationGitPath{
			Repo: &model.ApplicationGitRepository{Id: "id"},
			Path: "path",
		},
		CloudProvider: "cloud-provider",
		CreatedAt:     2,
		UpdatedAt:     2,
	}
	err := store.Create(ctx, col, "find-id-1", fakeApplication)
	require.NoError(t, err)
	err = store.Create(ctx, col, "find-id-2", fakeApplication2)
	require.NoError(t, err)

	testcases := []struct {
		name    string
		opts    datastore.ListOptions
		want    []*model.Application
		wantErr bool
	}{
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
			want: []*model.Application{
				fakeApplication,
			},
			wantErr: false,
		},
		{
			name: "only cursor given",
			opts: datastore.ListOptions{
				Cursor: "cursor",
			},
			want:    []*model.Application{},
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			it, err := store.Find(ctx, col, tc.opts)
			assert.Equal(t, tc.wantErr, err != nil)
			got, err := listApplications(it)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func listApplications(it datastore.Iterator) ([]*model.Application, error) {
	ret := make([]*model.Application, 0)
	if it == nil {
		return ret, nil
	}
	for {
		var v model.Application
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
