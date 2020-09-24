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
	assert.Nil(t, err)

	err = store.Create(ctx, "kind", "id", &Entity{Name: "name"})
	assert.Nil(t, err)

	tests := []struct {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &Entity{}
			err := store.Get(ctx, tt.kind, tt.id, got)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
