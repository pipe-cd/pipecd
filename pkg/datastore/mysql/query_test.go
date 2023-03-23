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

package mysql

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/datastore"
)

func TestBuildGetQuery(t *testing.T) {
	testcases := []struct {
		name          string
		kind          string
		expectedQuery string
	}{
		{
			name:          "query for Project kind",
			kind:          "Project",
			expectedQuery: "SELECT Data FROM Project WHERE Id = UUID_TO_BIN(?,true)",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			query := buildGetQuery(tc.kind)
			assert.Equal(t, tc.expectedQuery, query)
		})
	}
}

func TestBuildUpdateQuery(t *testing.T) {
	testcases := []struct {
		name          string
		kind          string
		expectedQuery string
	}{
		{
			name:          "query for Project kind",
			kind:          "Project",
			expectedQuery: "UPDATE Project SET Data = ? WHERE Id = UUID_TO_BIN(?,true)",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			query := buildUpdateQuery(tc.kind)
			assert.Equal(t, tc.expectedQuery, query)
		})
	}
}

func TestBuildCreateQuery(t *testing.T) {
	testcases := []struct {
		name          string
		kind          string
		expectedQuery string
	}{
		{
			name:          "query for Project kind",
			kind:          "Project",
			expectedQuery: "INSERT INTO Project (Id, Data) VALUE (UUID_TO_BIN(?,true), ?)",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			query := buildCreateQuery(tc.kind)
			assert.Equal(t, tc.expectedQuery, query)
		})
	}
}

func TestBuildFindQuery(t *testing.T) {
	testcases := []struct {
		name          string
		kind          string
		listOptions   datastore.ListOptions
		expectedQuery string
		wantErr       bool
	}{
		{
			name:          "query without filter and order",
			kind:          "Project",
			listOptions:   datastore.ListOptions{},
			expectedQuery: "SELECT Data FROM Project",
		},
		{
			name: "query with one filter",
			kind: "Project",
			listOptions: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "Extra",
						Operator: datastore.OperatorEqual,
						Value:    "app-1",
					},
				},
			},
			expectedQuery: "SELECT Data FROM Project WHERE Extra = ?",
		},
		{
			name: "query with wrapped filter field name in where clause",
			kind: "Project",
			listOptions: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "SyncState.Status",
						Operator: datastore.OperatorEqual,
						Value:    1,
					},
				},
			},
			expectedQuery: "SELECT Data FROM Project WHERE SyncState_Status = ?",
		},
		{
			name: "query with multi filters",
			kind: "Project",
			listOptions: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "Data->>\"$.name\"",
						Operator: datastore.OperatorEqual,
						Value:    "app-123",
					},
					{
						Field:    "Extra",
						Operator: datastore.OperatorEqual,
						Value:    "app-1",
					},
				},
			},
			expectedQuery: "SELECT Data FROM Project WHERE Data->>\"$.name\" = ? AND Extra = ?",
		},
		{
			name: "query with one filter and one order by column",
			kind: "Project",
			listOptions: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "Extra",
						Operator: datastore.OperatorEqual,
						Value:    "app-1",
					},
				},
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
			},
			expectedQuery: "SELECT Data FROM Project WHERE Extra = ? ORDER BY UpdatedAt DESC, Id ASC",
		},
		{
			name: "query with wrapped filter field name as order by column",
			kind: "Project",
			listOptions: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "Extra",
						Operator: datastore.OperatorEqual,
						Value:    "app-1",
					},
				},
				Orders: []datastore.Order{
					{
						Field:     "SyncState.Status",
						Direction: datastore.Desc,
					},
					{
						Field:     "Id",
						Direction: datastore.Asc,
					},
				},
			},
			expectedQuery: "SELECT Data FROM Project WHERE Extra = ? ORDER BY SyncState_Status DESC, Id ASC",
		},
		{
			name: "query with one filter and one order by on 2 columns",
			kind: "Project",
			listOptions: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "Extra",
						Operator: datastore.OperatorEqual,
						Value:    "app-1",
					},
				},
				Orders: []datastore.Order{
					{
						Field:     "CreatedAt",
						Direction: datastore.Asc,
					},
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
			expectedQuery: "SELECT Data FROM Project WHERE Extra = ? ORDER BY CreatedAt ASC, UpdatedAt DESC, Id ASC",
		},
		{
			name: "query with unsupported operator",
			kind: "Project",
			listOptions: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "Extra",
						Operator: 0,
						Value:    "app-%",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "query with IN operator",
			kind: "Project",
			listOptions: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "Status",
						Operator: datastore.OperatorIn,
						Value:    []int32{1, 2, 3},
					},
				},
			},
			expectedQuery: "SELECT Data FROM Project WHERE Status IN (?,?,?)",
		},
		{
			name: "query with MEMBER OF operator",
			kind: "Piped",
			listOptions: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "EnvIds",
						Operator: datastore.OperatorContains,
						Value:    "xxx",
					},
				},
			},
			expectedQuery: "SELECT Data FROM Piped WHERE ? MEMBER OF (EnvIds)",
		},
		{
			name: "query with IN operator (one element)",
			kind: "Project",
			listOptions: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "Status",
						Operator: datastore.OperatorIn,
						Value:    []int32{1},
					},
				},
			},
			expectedQuery: "SELECT Data FROM Project WHERE Status IN (?)",
		},
		{
			name: "query with limit",
			kind: "Project",
			listOptions: datastore.ListOptions{
				Limit: 20,
			},
			expectedQuery: "SELECT Data FROM Project LIMIT 20",
		},
		{
			name: "query with pagination cursor",
			kind: "Application",
			listOptions: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "ProjectId",
						Operator: datastore.OperatorEqual,
					},
				},
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
				Limit: 20,
				Cursor: func() string {
					return base64.StdEncoding.EncodeToString([]byte(`{"Id":"object-id","UpdatedAt":100}`))
				}(),
			},
			expectedQuery: "SELECT Data FROM Application WHERE ProjectId = ? AND UpdatedAt <= ? AND NOT (UpdatedAt = ? AND Id <= UUID_TO_BIN(?,true)) ORDER BY UpdatedAt DESC, Id ASC LIMIT 20",
			wantErr:       false,
		},
		{
			name: "query with pagination cursor: no filter",
			kind: "Application",
			listOptions: datastore.ListOptions{
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
					return base64.StdEncoding.EncodeToString([]byte(`{"Id":"object-id","UpdatedAt":100}`))
				}(),
			},
			expectedQuery: "SELECT Data FROM Application WHERE UpdatedAt <= ? AND NOT (UpdatedAt = ? AND Id <= UUID_TO_BIN(?,true)) ORDER BY UpdatedAt DESC, Id ASC",
			wantErr:       false,
		},
		{
			name: "query with pagination cursor: more than 2 ordering fields",
			kind: "Application",
			listOptions: datastore.ListOptions{
				Orders: []datastore.Order{
					{
						Field:     "UpdatedAt",
						Direction: datastore.Desc,
					},
					{
						Field:     "CreatedAt",
						Direction: datastore.Desc,
					},
					{
						Field:     "Id",
						Direction: datastore.Asc,
					},
				},
				Cursor: func() string {
					return base64.StdEncoding.EncodeToString([]byte(`{"Id":"object-id","UpdatedAt":100,"CreatedAt":100}`))
				}(),
			},
			expectedQuery: "SELECT Data FROM Application WHERE UpdatedAt <= ? AND CreatedAt <= ? AND NOT (UpdatedAt = ? AND CreatedAt = ? AND Id <= UUID_TO_BIN(?,true)) ORDER BY UpdatedAt DESC, CreatedAt DESC, Id ASC",
			wantErr:       false,
		},
		{
			name: "query with cursor: missing Id from ordering fields",
			kind: "Project",
			listOptions: datastore.ListOptions{
				Orders: []datastore.Order{
					{
						Field:     "UpdatedAt",
						Direction: datastore.Desc,
					},
				},
				Cursor: func() string {
					return base64.StdEncoding.EncodeToString([]byte(`{"UpdatedAt":100}`))
				}(),
			},
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := buildFindQuery(tc.kind, tc.listOptions)
			assert.Equal(t, tc.expectedQuery, query)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestRefineFiltersValue(t *testing.T) {
	testcases := []struct {
		name               string
		filters            []datastore.ListFilter
		expectedFiltersVal []interface{}
	}{
		{
			name: "mixed types",
			filters: []datastore.ListFilter{
				{
					Value: 1,
				},
				{
					Value: "app-1",
				},
				{
					Value: []string{"app-1", "app-2", "app-3"},
				},
				{
					Value: []int32{1, 2, 3},
				},
				{
					Value: [3]int32{1, 2, 3},
				},
			},
			expectedFiltersVal: []interface{}{
				1,
				"app-1",
				"app-1", "app-2", "app-3",
				int32(1), int32(2), int32(3),
				int32(1), int32(2), int32(3),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			vals := refineFiltersValue(tc.filters)
			assert.Equal(t, tc.expectedFiltersVal, vals)
		})
	}
}

func TestMakeCompareOperatorForOuterSet(t *testing.T) {
	testcases := []struct {
		name      string
		direction datastore.OrderDirection
		expectOpe string
	}{
		{
			name:      "should return ope same direction with order direction: asc",
			direction: datastore.Asc,
			expectOpe: ">=",
		},
		{
			name:      "should return ope same direction with order direction: desc",
			direction: datastore.Desc,
			expectOpe: "<=",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ope := makeCompareOperatorForOuterSet(tc.direction)
			assert.Equal(t, tc.expectOpe, ope)
		})
	}
}

func TestMakeCompareOperatorForSubSet(t *testing.T) {
	testcases := []struct {
		name      string
		direction datastore.OrderDirection
		expectOpe string
	}{
		{
			name:      "should return ope in revert direction with order direction: asc",
			direction: datastore.Asc,
			expectOpe: "<=",
		},
		{
			name:      "should return ope in revert direction with order direction: desc",
			direction: datastore.Desc,
			expectOpe: ">=",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ope := makeCompareOperatorForSubSet(tc.direction)
			assert.Equal(t, tc.expectOpe, ope)
		})
	}
}

func TestMakePaginationCursorValues(t *testing.T) {
	testcases := []struct {
		name               string
		opts               datastore.ListOptions
		expectedCursorVals []interface{}
		wantErr            bool
	}{
		{
			name: "valid cursor with CamelCase key",
			opts: datastore.ListOptions{
				Orders: []datastore.Order{
					{
						Field:     "UpdatedAt",
						Direction: datastore.Desc,
					},
					{
						Field:     "CreatedAt",
						Direction: datastore.Desc,
					},
					{
						Field:     "Id",
						Direction: datastore.Asc,
					},
				},
				Cursor: func() string {
					return base64.StdEncoding.EncodeToString([]byte(`{"Id":"object-id","UpdatedAt":100,"CreatedAt":99}`))
				}(),
			},
			expectedCursorVals: []interface{}{
				float64(100),
				float64(99),
				float64(100),
				float64(99),
				"object-id",
			},
		},
		{
			name: "invalid cursor with snake_case key",
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
					return base64.StdEncoding.EncodeToString([]byte(`{"id":"object-id","updated_at":100}`))
				}(),
			},
			wantErr: true,
		},
		{
			name: "invalid cursor missing ordering field value",
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
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			vals, err := makePaginationCursorValues(tc.opts)
			assert.Equal(t, tc.expectedCursorVals, vals)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}
