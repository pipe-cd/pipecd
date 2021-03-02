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
	"database/sql"
	"fmt"
	"io/ioutil"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/pipe-cd/pipe/pkg/datastore"
	"go.uber.org/zap"
)

// MySQL client wrapper
type MySQL struct {
	client       *sql.DB
	logger       *zap.Logger
	usernameFile string
	passwordFile string
}

// Option for create MySQL typed instance
type Option func(*MySQL)

// WithLogger returns logger setup function
func WithLogger(logger *zap.Logger) Option {
	return func(m *MySQL) {
		m.logger = logger
	}
}

// WithAuthenticationFile returns auth info setup function
func WithAuthenticationFile(usernameFile, passwordFile string) Option {
	return func(m *MySQL) {
		m.usernameFile = usernameFile
		m.passwordFile = passwordFile
	}
}

// NewMySQL returns new MySQL instance
func NewMySQL(url, database string, opts ...Option) (*MySQL, error) {
	m := &MySQL{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(m)
	}

	dataSourceName, err := m.buildDataSourceName(url, database)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	m.client = db

	return m, nil
}

// Find implementation for MySQL
func (m *MySQL) Find(ctx context.Context, kind string, opts datastore.ListFilter) (datastore.Iterator, error) {
	return nil, datastore.ErrUnimplemented
}

// Get implementation for MySQL
func (m *MySQL) Get(ctx context.Context, kind, id string, v interface{}) error {
	return datastore.ErrUnimplemented
}

// Create implementation for MySQL
func (m *MySQL) Create(ctx context.Context, kind, id string, entity interface{}) error {
	return datastore.ErrUnimplemented
}

// Put implementation for MySQL
func (m *MySQL) Put(ctx context.Context, kind, id string, entity interface{}) error {
	return datastore.ErrUnimplemented
}

// Update implementation for MySQL
func (m *MySQL) Update(ctx context.Context, kind, id string, factory datastore.Factory, updater datastore.Updater) error {
	return datastore.ErrUnimplemented
}

// Close implementation for MySQL
func (m *MySQL) Close() error {
	return m.client.Close()
}

func (m *MySQL) buildDataSourceName(url, database string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("url is required field")
	}
	if database == "" {
		return "", fmt.Errorf("database is required field")
	}
	if m.usernameFile == "" || m.passwordFile == "" {
		return "", fmt.Errorf("credentials info are missing")
	}

	username, err := ioutil.ReadFile(m.usernameFile)
	if err != nil {
		return "", fmt.Errorf("failed to read username file: %w", err)
	}
	password, err := ioutil.ReadFile(m.passwordFile)
	if err != nil {
		return "", fmt.Errorf("failed to read password file: %w", err)
	}

	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, url, database), nil
}
