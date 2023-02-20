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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildDataSourceName(t *testing.T) {
	testcases := []struct {
		name           string
		url            string
		database       string
		usernameFile   string
		passwordFile   string
		dataSourceName string
		expectErr      bool
	}{
		{
			name:      "returns error on missing url",
			database:  "test",
			expectErr: true,
		},
		{
			name:      "returns error on missing database",
			url:       "test:test@tcp(localhost:3306)",
			expectErr: true,
		},
		{
			name:           "returns url with configured database",
			url:            "test:test@tcp(localhost:3306)/test",
			dataSourceName: "test:test@tcp(localhost:3306)/test",
			expectErr:      false,
		},
		{
			name:           "returns url/database in usernameFile or passwordFile are not provided",
			url:            "test:test@tcp(localhost:3306)",
			database:       "test",
			dataSourceName: "test:test@tcp(localhost:3306)/test",
			expectErr:      false,
		},
		{
			name:         "returns error on unable to read username or password from files",
			url:          "localhost:3306",
			database:     "test",
			usernameFile: "./testdata/not_existed_file_name",
			passwordFile: "./testdata/not_existed_file_name",
			expectErr:    true,
		},
		{
			name:           "returns data source name",
			url:            "localhost:3306",
			database:       "test",
			usernameFile:   "./testdata/username",
			passwordFile:   "./testdata/password",
			dataSourceName: "test:test@tcp(localhost:3306)/test",
			expectErr:      false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			dataSourceName, err := BuildDataSourceName(tc.url, tc.database, tc.usernameFile, tc.passwordFile)
			assert.Equal(t, tc.expectErr, err != nil)
			assert.Equal(t, tc.dataSourceName, dataSourceName)
		})
	}
}

func TestMakeRowID(t *testing.T) {
	testcases := []struct {
		name    string
		modelID string
		rowID   string
	}{
		{
			name:    "modelID is simple string, not UUID",
			modelID: "pipecd",
			rowID:   "1b247cf8-ee2c-56db-af91-be9e25ff3b6a",
		},
		{
			name:    "modelID is as same as previuos, ensure test",
			modelID: "pipecd",
			rowID:   "1b247cf8-ee2c-56db-af91-be9e25ff3b6a",
		},
		{
			name:    "modelID is UUID",
			modelID: "dfc55495-7dbd-11eb-8636-42010a920020",
			rowID:   "dfc55495-7dbd-11eb-8636-42010a920020",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rowID := makeRowID(tc.modelID)
			assert.Equal(t, tc.rowID, rowID)
		})
	}
}
