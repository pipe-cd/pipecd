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
	"testing"
	"time"

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

func TestEntity(t *testing.T) {
	col := &collection{kind: "Entity"}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := client.Create(ctx, col, "id", &Entity{Name: "name"})
	require.Error(t, err)
}
