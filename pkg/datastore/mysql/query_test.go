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
			expectedQuery: "SELECT data FROM Project WHERE id = UUID_TO_BIN(?,true)",
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
			expectedQuery: "UPDATE Project SET data = ? WHERE id = UUID_TO_BIN(?,true)",
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
			expectedQuery: "INSERT INTO Project (id, data) VALUE (UUID_TO_BIN(?,true), ?)",
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
			expectedQuery: "INSERT INTO Project (id, data) VALUE (UUID_TO_BIN(?,true), ?) ON DUPLICATE KEY UPDATE data = ?",
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
	}{
		{
			name:          "query without filter and order",
			kind:          "Project",
			listOptions:   datastore.ListOptions{},
			expectedQuery: "SELECT data FROM Project",
		},
		{
			name: "query with one filter",
			kind: "Project",
			listOptions: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "extra",
						Operator: "LIKE",
						Value:    "app-1%",
					},
				},
			},
			expectedQuery: "SELECT data FROM Project WHERE extra LIKE ?",
		},
		{
			name: "query with multi filters",
			kind: "Project",
			listOptions: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "data->>\"$.name\"",
						Operator: "=",
						Value:    "app-123",
					},
					{
						Field:    "extra",
						Operator: "LIKE",
						Value:    "app-1%",
					},
				},
			},
			expectedQuery: "SELECT data FROM Project WHERE data->>\"$.name\" = ? AND extra LIKE ?",
		},
		{
			name: "query with one filter and one order by column",
			kind: "Project",
			listOptions: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "extra",
						Operator: "LIKE",
						Value:    "app-1%",
					},
				},
				Orders: []datastore.Order{
					{
						Field:     "updated_at",
						Direction: datastore.Desc,
					},
				},
			},
			expectedQuery: "SELECT data FROM Project WHERE extra LIKE ? ORDER BY updated_at DESC",
		},
		{
			name: "query with one filter and one order by on 2 columns",
			kind: "Project",
			listOptions: datastore.ListOptions{
				Filters: []datastore.ListFilter{
					{
						Field:    "extra",
						Operator: "LIKE",
						Value:    "app-1%",
					},
				},
				Orders: []datastore.Order{
					{
						Field:     "created_at",
						Direction: datastore.Asc,
					},
					{
						Field:     "updated_at",
						Direction: datastore.Desc,
					},
				},
			},
			expectedQuery: "SELECT data FROM Project WHERE extra LIKE ? ORDER BY created_at ASC, updated_at DESC",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			query := buildFindQuery(tc.kind, tc.listOptions)
			assert.Equal(t, tc.expectedQuery, query)
		})
	}
}
