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
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/datastore"
)

const mysqlErrorCodeDuplicateEntry = 1062
const mysqlErrorCodeUserDefined = 1644

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

	dataSourceName, err := BuildDataSourceName(url, database, m.usernameFile, m.passwordFile)
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
func (m *MySQL) Find(ctx context.Context, col datastore.Collection, opts datastore.ListOptions) (datastore.Iterator, error) {
	kind := col.Kind()
	if opts.Cursor != "" && len(opts.Orders) == 0 {
		return nil, errors.New("opts.Cursor also requires Orders to be set")
	}

	query, err := buildFindQuery(kind, opts)
	if err != nil {
		m.logger.Error("failed to build find entities query",
			zap.String("kind", kind),
			zap.Error(err),
		)
		return nil, err
	}

	whereConditionVals := refineFiltersValue(opts.Filters)
	cursorVals, err := makePaginationCursorValues(opts)
	if err != nil {
		return nil, err
	}
	whereConditionVals = append(whereConditionVals, cursorVals...)

	rows, err := m.client.QueryContext(ctx, query, whereConditionVals...)
	if err != nil {
		m.logger.Error("failed to find entities",
			zap.String("kind", kind),
			zap.String("query", query),
			zap.Any("whereConditionValues", whereConditionVals),
			zap.Error(err),
		)
		return nil, err
	}
	return &Iterator{
		rows:   rows,
		orders: opts.Orders,
	}, nil
}

// Get implementation for MySQL
func (m *MySQL) Get(ctx context.Context, col datastore.Collection, id string, v interface{}) error {
	kind := col.Kind()
	row := m.client.QueryRowContext(ctx, buildGetQuery(kind), makeRowID(id))
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

	return decodeJSONValue(val, v)
}

// Create implementation for MySQL
func (m *MySQL) Create(ctx context.Context, col datastore.Collection, id string, entity interface{}) error {
	kind := col.Kind()
	stmt, err := m.client.PrepareContext(ctx, buildCreateQuery(kind))
	if err != nil {
		m.logger.Error("failed to create entity: failed to prepare query",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}
	defer stmt.Close()

	data, err := encodeJSONValue(entity)
	if err != nil {
		m.logger.Error("failed to create entity: failed to encode json data",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	_, err = stmt.ExecContext(ctx, makeRowID(id), data)
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if mysqlErr.Number == mysqlErrorCodeDuplicateEntry {
			return datastore.ErrAlreadyExists
		}
		if mysqlErr.Number == mysqlErrorCodeUserDefined {
			return fmt.Errorf("%w: %s", datastore.ErrUserDefined, mysqlErr.Message)
		}
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

// Update implementation for MySQL
func (m *MySQL) Update(ctx context.Context, col datastore.Collection, id string, updater datastore.Updater) error {
	kind := col.Kind()
	// Start transaction with default isolation level.
	tx, err := m.client.BeginTx(ctx, nil)
	if err != nil {
		m.logger.Error("failed to update entity: failed to start transaction",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	row := tx.QueryRowContext(ctx, buildGetQuery(kind), makeRowID(id))
	var val string
	err = row.Scan(&val)
	if err == sql.ErrNoRows {
		tx.Rollback()
		return datastore.ErrNotFound
	}
	if err != nil {
		m.logger.Error("failed to update entity: failed to get entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		tx.Rollback()
		return err
	}

	entity := col.Factory()()
	if err := decodeJSONValue(val, entity); err != nil {
		m.logger.Error("failed to update entity: failed to decode data",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		tx.Rollback()
		return err
	}

	if err := updater(entity); err != nil {
		m.logger.Error("failed to update entity: failed to apply updater",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		tx.Rollback()
		return err
	}

	data, err := encodeJSONValue(entity)
	if err != nil {
		m.logger.Error("failed to update entity: failed to encode json data",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx, buildUpdateQuery(kind), data, makeRowID(id))
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

// Ping implementation for MySQL
func (m *MySQL) Ping() error {
	return m.client.Ping()
}

// BuildDataSourceName returns source name to make connection to database.
func BuildDataSourceName(url, database, usernameFile, passwordFile string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("url is required field")
	}
	if database == "" {
		if !strings.Contains(url, "/") {
			return "", fmt.Errorf("database is not set")
		}
		return url, nil
	}
	// In case username and password files are not provided,
	// those values may be included in the URL already so just return the given URL attached with Database name.
	if usernameFile == "" || passwordFile == "" {
		return fmt.Sprintf("%s/%s", url, database), nil
	}

	username, err := os.ReadFile(usernameFile)
	if err != nil {
		return "", fmt.Errorf("failed to read username file: %w", err)
	}
	password, err := os.ReadFile(passwordFile)
	if err != nil {
		return "", fmt.Errorf("failed to read password file: %w", err)
	}

	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, url, database), nil
}

// makeRowID converts a given string which not in UUID format to UUID string.
// Otherwise, return the id itself.
func makeRowID(id string) string {
	_, err := uuid.Parse(id)
	if err == nil {
		return id
	}
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(id)).String()
}
