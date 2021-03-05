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
