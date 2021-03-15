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

package mysqlensurer

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"

	datastore "github.com/pipe-cd/pipe/pkg/datastore/mysql"
)

var (
	mysqlDatabaseSchema  = mysqlProperties_1
	mysqlDatabaseIndexes = mysqlProperties_0
)

const mysqlErrorCodeDuplicateKeyName = 1061

type mysqlEnsurer struct {
	client       *sql.DB
	logger       *zap.Logger
	url          string
	database     string
	usernameFile string
	passwordFile string
}

func NewMySQLEnsurer(url, database, usernameFile, passwordFile string, logger *zap.Logger) SQLEnsurer {
	return &mysqlEnsurer{
		url:          url,
		database:     database,
		usernameFile: usernameFile,
		passwordFile: passwordFile,
		logger:       logger.Named("mysql-ensurer"),
	}
}

func (m *mysqlEnsurer) EnsureIndexes(ctx context.Context) error {
	if err := m.connect(); err != nil {
		return err
	}

	for _, stmt := range makeCreateIndexStatements(mysqlDatabaseIndexes) {
		_, err := m.client.ExecContext(ctx, stmt)
		// Ignore in case error duplicate key name occurred.
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == mysqlErrorCodeDuplicateKeyName {
			continue
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *mysqlEnsurer) EnsureSchema(ctx context.Context) error {
	if err := m.connect(); err != nil {
		return err
	}

	_, err := m.client.ExecContext(ctx, mysqlDatabaseSchema)
	if err != nil {
		return err
	}
	return nil
}

func (m *mysqlEnsurer) Close() error {
	if m.client == nil {
		return nil
	}
	return m.client.Close()
}

func (m *mysqlEnsurer) connect() error {
	if m.client != nil {
		return nil
	}

	dataSourceName, err := datastore.BuildDataSourceName(m.url, m.database, m.usernameFile, m.passwordFile)
	if err != nil {
		return fmt.Errorf("failed to connect to sql database: %w", err)
	}

	// Enable run multi statements at once.
	db, err := sql.Open("mysql", fmt.Sprintf("%s?multiStatements=true", dataSourceName))
	if err != nil {
		return fmt.Errorf("failed to connect to sql database: %w", err)
	}
	m.client = db
	return nil
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
