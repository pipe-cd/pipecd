// Copyright 2020 The PipeCD Authors.
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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipe/pkg/datastore"
	"github.com/pipe-cd/pipe/pkg/datastore/firestore"
)

func TestGet(t *testing.T) {
	type Entity struct {
		Name string
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	store, err := firestore.NewFireStore(ctx, "project", "namespace", "environment")
	require.NoError(t, err)

	err = store.Create(ctx, "kind", "id", &Entity{Name: "name"})
	require.NoError(t, err)

	testcases := []struct {
		name    string
		kind    string
		id      string
		want    *Entity
		wantErr error
	}{
		{
			name:    "entity found",
			kind:    "kind",
			id:      "id",
			want:    &Entity{Name: "name"},
			wantErr: nil,
		},
		{
			name:    "not found",
			kind:    "kind",
			id:      "wrong-id",
			want:    &Entity{},
			wantErr: datastore.ErrNotFound,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := &Entity{}
			err := store.Get(ctx, tc.kind, tc.id, got)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
