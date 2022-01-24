// Copyright 2022 The PipeCD Authors.
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

package datastore

import "context"

type collectionStore struct {
	col Collection
	ds  DataStore
}

func (c *collectionStore) Find(ctx context.Context, opts ListOptions) (Iterator, error) {
	return c.ds.Find(ctx, c.col.Kind(), opts)
}

func (c *collectionStore) Get(ctx context.Context, id string, entity interface{}) error {
	return c.ds.Get(ctx, c.col.Kind(), id, entity)
}

func (c *collectionStore) Create(ctx context.Context, id string, entity interface{}) error {
	return c.ds.Create(ctx, c.col.Kind(), id, entity)
}

func (c *collectionStore) Put(ctx context.Context, id string, entity interface{}) error {
	return c.ds.Put(ctx, c.col.Kind(), id, entity)
}

func (c *collectionStore) Update(ctx context.Context, id string, updater Updater) error {
	return c.ds.Update(ctx, c.col.Kind(), id, c.col.Factory(), updater)
}
