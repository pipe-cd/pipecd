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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/datastore"
)

func TestProcessCursorArg(t *testing.T) {
	testcases := []struct {
		name       string
		opts       datastore.ListOptions
		expectVals []interface{}
		expectErr  bool
	}{
		{
			name: "contains required fields",
			opts: datastore.ListOptions{
				Orders: []datastore.Order{
					{
						Field:     "UpdatedAt",
						Direction: datastore.Desc,
					},
					{
						Field:     "Id",
						Direction: datastore.Asc,
					},
				},
				Cursor: func() string {
					return base64.StdEncoding.EncodeToString([]byte(`{"UpdatedAt":100,"Id":"object-id"}`))
				}(),
			},
			expectVals: []interface{}{float64(100), "object-id"},
			expectErr:  false,
		},
		{
			name: "missing id field in ordering options",
			opts: datastore.ListOptions{
				Orders: []datastore.Order{
					{
						Field:     "UpdatedAt",
						Direction: datastore.Desc,
					},
				},
				Cursor: func() string {
					return base64.StdEncoding.EncodeToString([]byte(`{"UpdatedAt":100,"Id":"object-id"}`))
				}(),
			},
			expectErr: true,
		},
		{
			name: "invalid cursor: does not contain fields from ordering list",
			opts: datastore.ListOptions{
				Orders: []datastore.Order{
					{
						Field:     "UpdatedAt",
						Direction: datastore.Desc,
					},
					{
						Field:     "Id",
						Direction: datastore.Asc,
					},
				},
				Cursor: func() string {
					return base64.StdEncoding.EncodeToString([]byte(`{"Id":"object-id"}`))
				}(),
			},
			expectErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			vals, err := makeCursorValues(tc.opts)
			assert.Equal(t, tc.expectErr, err != nil)
			assert.Equal(t, tc.expectVals, vals)
		})
	}
}
