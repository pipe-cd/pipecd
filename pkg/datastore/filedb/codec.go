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

package filedb

import (
	"encoding/json"

	"github.com/pipe-cd/pipecd/pkg/datastore"
)

// decode checks for the given collection object. If the given collection
// implements the `datastore.ShardDecoder` interface, its implementation will
// be used. If not, time order regardless merge logic will be used.
func decode(col datastore.Collection, e interface{}, parts map[datastore.Shard][]byte) error {
	dcol, ok := col.(datastore.ShardDecoder)
	if ok {
		return dcol.Decode(e, parts)
	}

	ps := make([][]byte, 0, len(parts))
	for _, part := range parts {
		ps = append(ps, part)
	}

	// In case it's single part contained object, unmarshal it directly.
	if len(ps) == 1 {
		return json.Unmarshal(ps[0], e)
	}

	return merge(e, ps...)
}

type mergeable interface {
	GetUpdatedAt() int64
	SetUpdatedAt(t int64)
}

// merge function unmarshal all parts of the given data to entity e.
// The data will be merged regardless of its time order, after be merged,
// the latest UpdatedAt time will be used as the entity UpdatedAt value.
func merge(e interface{}, parts ...[]byte) error {
	me, ok := e.(mergeable)
	if !ok {
		return datastore.ErrUnsupported
	}

	var latest int64
	for _, p := range parts {
		if err := json.Unmarshal(p, e); err != nil {
			return err
		}
		if latest < me.GetUpdatedAt() {
			latest = me.GetUpdatedAt()
		}
	}
	// Fixme: Find a way to set updated_at value without force models having this setter.
	me.SetUpdatedAt(latest)
	return nil
}

// encode checks for the given collection object. If the given collection
// implements the `datastore.ShardEncoder` interface, its implementation will
// be used. If not, `datastore.ErrUnsupported` error will be raised.
func encode(col datastore.Collection, e interface{}) (map[datastore.Shard][]byte, error) {
	ecol, ok := col.(datastore.ShardEncoder)
	if !ok {
		return nil, datastore.ErrUnsupported
	}

	return ecol.Encode(e)
}
