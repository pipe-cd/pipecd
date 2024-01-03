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

package ensurer

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"strings"

	driver "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/datastore/mysql"
)

var (
	//go:embed schema.sql
	mysqlDatabaseSchema string
	//go:embed indexes.sql
	mysqlDatabaseIndexes string
)

const (
	mysqlErrorCodeDuplicateColumnName  = 1060
	mysqlErrorCodeDuplicateKeyName     = 1061
	mysqlErrorCodeColumnOrKeyNotExists = 1091
)

type mysqlEnsurer struct {
	client *sql.DB
	logger *zap.Logger
}

func NewMySQLEnsurer(url, database, usernameFile, passwordFile string, logger *zap.Logger) (SQLEnsurer, error) {
	dataSourceName, err := mysql.BuildDataSourceName(url, database, usernameFile, passwordFile)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to sql database: %w", err)
	}

	// Enable run multi statements at once.
	db, err := sql.Open("mysql", fmt.Sprintf("%s?multiStatements=true", dataSourceName))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to sql database: %w", err)
	}

	return &mysqlEnsurer{
		client: db,
		logger: logger.Named("mysql-ensurer"),
	}, nil
}

func (m *mysqlEnsurer) EnsureIndexes(ctx context.Context) error {
	for _, stmt := range makeCreateIndexStatements(mysqlDatabaseIndexes) {
		_, err := m.client.ExecContext(ctx, stmt)
		// Ignore in case error inc case:
		// - Duplicate key name or column name occurred.
		// - Try to remove not exists column or key.
		if mysqlErr, ok := err.(*driver.MySQLError); ok && (mysqlErr.Number == mysqlErrorCodeDuplicateKeyName ||
			mysqlErr.Number == mysqlErrorCodeDuplicateColumnName ||
			mysqlErr.Number == mysqlErrorCodeColumnOrKeyNotExists) {
			continue
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *mysqlEnsurer) EnsureSchema(ctx context.Context) error {
	_, err := m.client.ExecContext(ctx, mysqlDatabaseSchema)
	if err != nil {
		return err
	}
	return nil
}

func (m *mysqlEnsurer) Close() error {
	return m.client.Close()
}

func (m *mysqlEnsurer) Ping() error {
	return m.client.Ping()
}

func makeCreateIndexStatements(indexesStatements string) []string {
	items := strings.Split(strings.TrimSpace(indexesStatements), ";")
	statements := make([]string, 0, len(items))
	for _, item := range items {
		// Ignore dummy statement.
		if item == "" {
			continue
		}
		statements = append(statements, strings.TrimSpace(item))
	}
	return statements
}
