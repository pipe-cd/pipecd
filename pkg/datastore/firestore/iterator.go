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
	"encoding/json"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/pipe-cd/pipe/pkg/datastore"
)

type Iterator struct {
	it   *firestore.DocumentIterator
	last interface{}
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

func (it *Iterator) Cursor() (string, error) {
	if it.last == nil {
		return "", datastore.ErrInvalidCursor
	}

	cursorVal, err := json.Marshal(it.last)
	if err != nil {
		return "", err
	}

	// TODO: Encrypt cursor string value before pass it to the caller.
	return string(cursorVal), nil
}
