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

package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/pipe-cd/pipe/pkg/datastore"
)

type Iterator struct {
	ctx context.Context
	cur *mongo.Cursor
}

func (it *Iterator) Next(dst interface{}) error {
	if !it.cur.Next(it.ctx) {
		return datastore.ErrIteratorDone
	}
	wrapper, err := wrapModel(dst)
	if err != nil {
		return err
	}
	if err := it.cur.Decode(wrapper); err != nil {
		return err
	}
	return extractModel(wrapper, dst)
}

func (it *Iterator) Cursor() (string, error) {
	// Note: Perhaps, not required.
	return "", datastore.ErrUnimplemented
}
