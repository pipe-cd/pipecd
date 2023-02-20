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

package mysqlensurer

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/datastore/mysql/ensurer"
)

type mysqlEnsurer struct {
	exec ensurer.SQLEnsurer
}

func NewMySQLEnsurer(url, database, usernameFile, passwordFile string, logger *zap.Logger) (SQLEnsurer, error) {
	executor, err := ensurer.NewMySQLEnsurer(url, database, usernameFile, passwordFile, logger.Named("mysql-ensurer"))
	if err != nil {
		return nil, fmt.Errorf("failed to create mysql ensurer executor: %w", err)
	}

	return &mysqlEnsurer{
		exec: executor,
	}, nil
}

func (m *mysqlEnsurer) Run(ctx context.Context) error {
	err := m.exec.EnsureSchema(ctx)
	if err != nil {
		return fmt.Errorf("failed to prepare sql database: %w", err)
	}

	// No need to run this create indexes operation in routine because it runs asynchronously.
	// ref: https://dev.mysql.com/doc/refman/8.0/en/innodb-online-ddl-operations.html#online-ddl-index-operations
	err = m.exec.EnsureIndexes(ctx)
	if err != nil {
		return fmt.Errorf("failed to create required indexes on sql database: %w", err)
	}

	return nil
}

func (m *mysqlEnsurer) Close() error {
	return m.exec.Close()
}
