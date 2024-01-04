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
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/pipe-cd/pipecd/pkg/datastore"
)

type dataConverter interface {
	Data() map[string]interface{}
}

// Iterator for MySQL result set
type Iterator struct {
	rows   *sql.Rows
	orders []datastore.Order
	last   dataConverter
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

	// Update last iterated item as last read row.
	it.last = &rowDataConverter{val: val}

	return decodeJSONValue(val, dst)
}

// Cursor builds a base64 string (encode from string in map[string]interface{} format).
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

type rowDataConverter struct {
	val string
}

// Data make JSON object with key in CamelCase format.
func (r *rowDataConverter) Data() map[string]interface{} {
	jsonRaw := convertKeys(json.RawMessage(r.val), convertSnakeToCamel)
	obj := make(map[string]interface{})
	json.Unmarshal(jsonRaw, &obj)
	return obj
}

// convertKeys convert all keys of json object with convert function.
func convertKeys(j json.RawMessage, convertFunc func(string) string) json.RawMessage {
	m := make(map[string]json.RawMessage)
	if err := json.Unmarshal([]byte(j), &m); err != nil {
		// Not a JSON object
		return j
	}

	for k, v := range m {
		fixed := convertFunc(k)
		delete(m, k)
		m[fixed] = convertKeys(v, convertFunc)
	}

	b, err := json.Marshal(m)
	if err != nil {
		return j
	}

	return json.RawMessage(b)
}

func convertSnakeToCamel(key string) string {
	var out string
	isToUpper := true
	for _, v := range key {
		if isToUpper {
			out += strings.ToUpper(string(v))
			isToUpper = false
			continue
		}
		if v == '_' {
			isToUpper = true
			continue
		}
		out += string(v)
	}
	return out
}
