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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipe/pkg/model"
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
			dataSourceName, err := buildDataSourceName(tc.url, tc.database, tc.usernameFile, tc.passwordFile)
			assert.Equal(t, tc.expectErr, err != nil)
			assert.Equal(t, tc.dataSourceName, dataSourceName)
		})
	}
}

func TestBuildGetQuery(t *testing.T) {
	testcases := []struct {
		name          string
		kind          string
		expectedQuery string
	}{
		{
			name:          "query for Project kind",
			kind:          "Project",
			expectedQuery: "SELECT data FROM Project WHERE extra = ?",
		},
		{
			name:          "query for other kinds than Project",
			kind:          "Application",
			expectedQuery: "SELECT data FROM Application WHERE id = UUID_TO_BIN(?,true)",
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
			expectedQuery: "UPDATE Project SET data = ? WHERE extra = ?",
		},
		{
			name:          "query for other kinds than Project",
			kind:          "Application",
			expectedQuery: "UPDATE Application SET data = ? WHERE id = UUID_TO_BIN(?,true)",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			query := buildUpdateQuery(tc.kind)
			assert.Equal(t, tc.expectedQuery, query)
		})
	}
}

func makeTestableMySQLClient() *MySQL {
	options := []Option{
		WithAuthenticationFile("/Users/s12228/workspace/pipe-cd/pipe/.dev/mysql_username", "/Users/s12228/workspace/pipe-cd/pipe/.dev/mysql_password"),
	}
	m, _ := NewMySQL("127.0.0.1:3307", "pipecd", options...)
	return m
}

func _TestPut(t *testing.T) {
	m := makeTestableMySQLClient()
	err := m.Put(context.TODO(), "Project", "proj-1", &model.Project{
		Id:        "proj-1",
		CreatedAt: 100,
		UpdatedAt: 102,
	})
	assert.Equal(t, nil, err)
}

func _TestCreate(t *testing.T) {
	m := makeTestableMySQLClient()
	err := m.Create(context.TODO(), "Project", "proj-1", &model.Project{
		Id:        "proj-1",
		CreatedAt: 100,
		UpdatedAt: 100,
	})
	assert.Equal(t, nil, err)
}

func _TestGet(t *testing.T) {
	m := makeTestableMySQLClient()
	p := &model.Project{}
	err := m.Get(context.TODO(), "Project", "proj-1", p)
	assert.Equal(t, nil, err)
	assert.Equal(t, "proj-1", p.Id)
}

var projectFactory = func() interface{} {
	return &model.Project{}
}

var projectUpdater = func(p *model.Project) error {
	p.UpdatedAt++
	return nil
}

func _TestUpdate(t *testing.T) {
	m := makeTestableMySQLClient()
	err := m.Update(context.TODO(), "Project", "proj-1", projectFactory, func(e interface{}) error {
		d := e.(*model.Project)
		return projectUpdater(d)
	})
	assert.Nil(t, err)
}
