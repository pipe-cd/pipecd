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

package firestore

import (
	"encoding/base64"
	"encoding/json"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/pipe-cd/pipecd/pkg/datastore"
)

type dataConverter interface {
	Data() map[string]interface{}
}

type Iterator struct {
	it     *firestore.DocumentIterator
	orders []datastore.Order
	last   dataConverter
}

func (it *Iterator) Next(dst interface{}) error {
	doc, err := it.it.Next()
	if err != nil {
		if err == iterator.Done {
			return datastore.ErrIteratorDone
		}
		return err
	}

	// Update last iterated item as last read doc.
	it.last = doc

	return doc.DataTo(dst)
}

// Cursor builds a base 64 string (encode from string in map[string]interface{} format).
// The cursor contains only values attached with the fields used
// as ordering fields.
func (it *Iterator) Cursor() (string, error) {
	if it.last == nil {
		return "", datastore.ErrInvalidCursor
	}

	lastObjData := it.last.Data()

	cursor := make(map[string]interface{}, len(it.orders))
	for _, o := range it.orders {
		val, ok := lastObjData[o.Field]
		if !ok {
			return "", datastore.ErrInvalidCursor
		}
		// TODO: Support build cursor from nested Ordering field.
		cursor[o.Field] = val
	}

	b, _ := json.Marshal(cursor)
	return base64.StdEncoding.EncodeToString(b), nil
}
