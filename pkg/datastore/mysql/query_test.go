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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipe/pkg/datastore"
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

func TestBuildPutQuery(t *testing.T) {
	testcases := []struct {
		name          string
		kind          string
		expectedQuery string
	}{
		{
			name:          "query for Project kind",
			kind:          "Project",
			expectedQuery: "INSERT INTO Project (Id, Data) VALUE (UUID_TO_BIN(?,true), ?) ON DUPLICATE KEY UPDATE Data = ?",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			query := buildPutQuery(tc.kind)
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
						Operator: "==",
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
						Operator: "==",
						Value:    1,
					},
				},
			},
			expectedQuery: "SELECT Data FROM Project WHERE SyncState = ?",
		},
		{
			name: "query with multi filters",
			kind: "Project",
			listOptions: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "Data->>\"$.name\"",
						Operator: "==",
						Value:    "app-123",
					},
					{
						Field:    "Extra",
						Operator: "==",
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
						Operator: "==",
						Value:    "app-1",
					},
				},
				Orders: []datastore.Order{
					{
						Field:     "UpdatedAt",
						Direction: datastore.Desc,
					},
				},
			},
			expectedQuery: "SELECT Data FROM Project WHERE Extra = ? ORDER BY UpdatedAt DESC",
		},
		{
			name: "query with wrapped filter field name as order by column",
			kind: "Project",
			listOptions: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "Extra",
						Operator: "==",
						Value:    "app-1",
					},
				},
				Orders: []datastore.Order{
					{
						Field:     "SyncState.Status",
						Direction: datastore.Desc,
					},
				},
			},
			expectedQuery: "SELECT Data FROM Project WHERE Extra = ? ORDER BY SyncState DESC",
		},
		{
			name: "query with one filter and one order by on 2 columns",
			kind: "Project",
			listOptions: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "Extra",
						Operator: "==",
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
				},
			},
			expectedQuery: "SELECT Data FROM Project WHERE Extra = ? ORDER BY CreatedAt ASC, UpdatedAt DESC",
		},
		{
			name: "query with limit",
			kind: "Project",
			listOptions: datastore.ListOptions{
				PageSize: 20,
			},
			expectedQuery: "SELECT Data FROM Project LIMIT 20",
		},
		{
			name: "query with limit offset",
			kind: "Project",
			listOptions: datastore.ListOptions{
				PageSize: 20,
				Page:     20,
			},
			expectedQuery: "SELECT Data FROM Project LIMIT 20 OFFSET 400",
		},
		{
			name: "query with unsupported operator",
			kind: "Project",
			listOptions: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "Extra",
						Operator: "like",
						Value:    "app-%",
					},
				},
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
