// Copyright 2021 The PipeCD Authors.
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
	"database/sql"

	"github.com/pipe-cd/pipe/pkg/datastore"
)

// Iterator for MySQL result set
type Iterator struct {
	rows *sql.Rows
}

// Next implementation for MySQL Iterator
func (it *Iterator) Next(dst interface{}) error {
	if !it.rows.Next() {
		return datastore.ErrIteratorDone
	}
	var val string
	err := it.rows.Scan(&val)
	if err != nil {
		return err
	}
	return decodeJSONValue(val, dst)
}

// Cursor implementation for MySQL Iterator
func (it *Iterator) Cursor() (string, error) {
	return "", datastore.ErrUnimplemented
}
