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

package objectcache

import (
	"encoding/json"
	"fmt"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/datastore"
)

type Cache interface {
	Get(shard datastore.Shard, id, etag string) ([]byte, error)
	Put(shard datastore.Shard, id, etag string, val []byte) error
}

type objectCache struct {
	backend cache.Cache
}

func NewCache(c cache.Cache) Cache {
	return &objectCache{backend: c}
}

type objectValue struct {
	Etag string `json:"etag"`
	Data []byte `json:"data"`
}

func (o *objectCache) Get(shard datastore.Shard, id, etag string) ([]byte, error) {
	raw, err := o.backend.Get(makeObjectKey(shard, id))
	if err != nil {
		return nil, err
	}

	var obj objectValue
	if err = json.Unmarshal(raw.([]byte), &obj); err != nil {
		return nil, err
	}

	if etag == obj.Etag {
		return obj.Data, nil
	}
	return nil, cache.ErrNotFound
}

func (o *objectCache) Put(shard datastore.Shard, id, etag string, val []byte) error {
	obj := &objectValue{
		Etag: etag,
		Data: val,
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	return o.backend.Put(makeObjectKey(shard, id), data)
}

func makeObjectKey(shard datastore.Shard, id string) string {
	return fmt.Sprintf("FILEDB:OBJECT:%s:%s", shard, id)
}
