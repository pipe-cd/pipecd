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
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/datastore"
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

	dataSourceName, err := buildDataSourceName(url, database, m.usernameFile, m.passwordFile)
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
func (m *MySQL) Find(ctx context.Context, kind string, opts datastore.ListOptions) (datastore.Iterator, error) {
	return nil, datastore.ErrUnimplemented
}

// Get implementation for MySQL
func (m *MySQL) Get(ctx context.Context, kind, id string, v interface{}) error {
	query := buildGetQuery(kind)
	row := m.client.QueryRowContext(ctx, query, id)
	var val string
	err := row.Scan(&val)
	if err == sql.ErrNoRows {
		return datastore.ErrNotFound
	}
	if err != nil {
		m.logger.Error("failed to get entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	return json.Unmarshal([]byte(val), v)
}

// Create implementation for MySQL
func (m *MySQL) Create(ctx context.Context, kind, id string, entity interface{}) error {
	stmt, err := m.client.PrepareContext(ctx, fmt.Sprintf("INSERT INTO %s (id, data) VALUE (UUID_TO_BIN(?,true), ?)", kind))
	if err != nil {
		m.logger.Error("failed to create entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}
	defer stmt.Close()

	val, err := json.Marshal(entity)
	if err != nil {
		return datastore.ErrInvalidArgument
	}

	// In case the given model is `project`, id is not in uuid type so just generate random uuid instead.
	if _, err := uuid.Parse(id); err != nil {
		id = uuid.New().String()
	}

	_, err = stmt.ExecContext(ctx, id, string(val))
	if err != nil && err.(*mysql.MySQLError).Number == 1062 {
		return datastore.ErrAlreadyExists
	}
	if err != nil {
		m.logger.Error("failed to create entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}
	return nil
}

// Put implementation for MySQL
func (m *MySQL) Put(ctx context.Context, kind, id string, entity interface{}) error {
	stmt, err := m.client.PrepareContext(ctx, fmt.Sprintf("INSERT INTO %s (id, data) VALUE (UUID_TO_BIN(?,true), ?) ON DUPLICATE KEY UPDATE data = ?", kind))
	if err != nil {
		m.logger.Error("failed to put entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}
	defer stmt.Close()

	val, err := json.Marshal(entity)
	if err != nil {
		return datastore.ErrInvalidArgument
	}

	// In case the given model is `project`, id is not in uuid type so just generate random uuid instead.
	if _, err := uuid.Parse(id); err != nil {
		id = uuid.New().String()
	}

	_, err = stmt.ExecContext(ctx, id, string(val), string(val))
	if err != nil {
		m.logger.Error("failed to put entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}
	return nil
}

// Update implementation for MySQL
func (m *MySQL) Update(ctx context.Context, kind, id string, factory datastore.Factory, updater datastore.Updater) error {
	// Start transaction with default isolation level.
	tx, err := m.client.BeginTx(ctx, nil)
	if err != nil {
		m.logger.Error("failed to update entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	query := buildGetQuery(kind)
	row := tx.QueryRowContext(ctx, query, id)
	var val string
	err = row.Scan(&val)
	if err == sql.ErrNoRows {
		tx.Rollback()
		return datastore.ErrNotFound
	}
	if err != nil {
		m.logger.Error("failed to update entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		tx.Rollback()
		return err
	}

	entity := factory()
	if err := json.Unmarshal([]byte(val), entity); err != nil {
		m.logger.Error("failed to update entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		tx.Rollback()
		return err
	}

	if err := updater(entity); err != nil {
		m.logger.Error("failed to update entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		tx.Rollback()
		return err
	}

	updateQuery := buildUpdateQuery(kind)
	encodedEntity, err := json.Marshal(entity)
	if err != nil {
		m.logger.Error("failed to update entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx, updateQuery, string(encodedEntity), id)
	if err != nil {
		m.logger.Error("failed to update entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Close implementation for MySQL
func (m *MySQL) Close() error {
	return m.client.Close()
}

func buildDataSourceName(url, database, usernameFile, passwordFile string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("url is required field")
	}
	if database == "" {
		return "", fmt.Errorf("database is required field")
	}
	// In case username and password files are not provided,
	// those values may be included in the URL already so just return the given URL attached with Database name.
	if usernameFile == "" || passwordFile == "" {
		return fmt.Sprintf("%s/%s", url, database), nil
	}

	username, err := ioutil.ReadFile(usernameFile)
	if err != nil {
		return "", fmt.Errorf("failed to read username file: %w", err)
	}
	password, err := ioutil.ReadFile(passwordFile)
	if err != nil {
		return "", fmt.Errorf("failed to read password file: %w", err)
	}

	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, url, database), nil
}

func buildGetQuery(table string) string {
	// TODO: make kinds from datastore package public
	if table == "Project" {
		return fmt.Sprintf("SELECT data FROM %s WHERE data->>\"$.id\" = ?", table)
	}
	return fmt.Sprintf("SELECT data FROM %s WHERE id = UUID_TO_BIN(?,true)", table)
}

func buildUpdateQuery(table string) string {
	// TODO: make kinds from datastore package public
	if table == "Project" {
		return fmt.Sprintf("UPDATE %s SET data = ? WHERE data->>\"$.id\" = ?", table)
	}
	return fmt.Sprintf("UPDATE %s SET data = ? WHERE id = UUID_TO_BIN(?,true)", table)
}
