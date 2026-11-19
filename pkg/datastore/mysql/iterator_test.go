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
	"encoding/base64"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/datastore"
)

type dummyDoc struct {
	val map[string]any
	err error
}

func (d *dummyDoc) Data() (map[string]any, error) {
	return d.val, d.err
}

func TestCursor(t *testing.T) {
	testcases := []struct {
		name         string
		iter         Iterator
		expectCursor string
		expectErr    bool
	}{
		{
			name:      "invalid cursor error returns on last is nil",
			iter:      Iterator{},
			expectErr: true,
		},
		{
			name: "error on converting the last row data",
			iter: Iterator{
				last: &dummyDoc{
					err: errors.New("failed to unmarshal row data"),
				},
			},
			expectErr: true,
		},
		{
			name: "valid last cursor",
			iter: Iterator{
				last: &dummyDoc{
					val: map[string]any{
						"Id":        "object-id",
						"CreatedAt": 100,
						"UpdatedAt": 100,
					},
				},
				orders: []datastore.Order{
					{
						Field:     "UpdatedAt",
						Direction: datastore.Desc,
					},
					{
						Field:     "Id",
						Direction: datastore.Asc,
					},
				},
			},
			expectCursor: func() string {
				return base64.StdEncoding.EncodeToString([]byte(`{"Id":"object-id","UpdatedAt":100}`))
			}(),
			expectErr: false,
		},
		{
			name: "invalid last cursor: field name of cursor data in snake_case",
			iter: Iterator{
				last: &dummyDoc{
					val: map[string]any{
						"id":         "object-id",
						"created_at": 100,
						"updated_at": 100,
					},
				},
				orders: []datastore.Order{
					{
						Field:     "UpdatedAt",
						Direction: datastore.Desc,
					},
					{
						Field:     "Id",
						Direction: datastore.Asc,
					},
				},
			},
			expectErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			cursor, err := tc.iter.Cursor()
			assert.Equal(t, tc.expectCursor, cursor)
			assert.Equal(t, tc.expectErr, err != nil)
		})
	}
}

func TestData(t *testing.T) {
	testcases := []struct {
		name         string
		rowData      string
		expectedData map[string]any
		expectErr    bool
	}{
		{
			name:    "valid data",
			rowData: `{"id": "object-id", "name": "app-1", "updated_at": 100, "created_at": 100}`,
			expectedData: map[string]any{
				"Id":        "object-id",
				"Name":      "app-1",
				"UpdatedAt": float64(100),
				"CreatedAt": float64(100),
			},
		},
		{
			name:    "valid nested data",
			rowData: `{"id": "object-id", "sync_state": { "status": 1 }, "updated_at": 100, "created_at": 100}`,
			expectedData: map[string]any{
				"Id": "object-id",
				"SyncState": map[string]any{
					"Status": float64(1),
				},
				"UpdatedAt": float64(100),
				"CreatedAt": float64(100),
			},
		},
		{
			name:      "invalid json data",
			rowData:   `{"id": "object-id", "name": "app-1", "updated_at": 100, "created_at": 100`,
			expectErr: true,
		},
		{
			name:      "data is not a json object",
			rowData:   `["object-id"]`,
			expectErr: true,
		},
		{
			name:      "data is a json null",
			rowData:   `null`,
			expectErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			converter := &rowDataConverter{val: tc.rowData}
			data, err := converter.Data()
			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, data)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedData, data)
		})
	}
}
