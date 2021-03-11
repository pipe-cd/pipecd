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
	"fmt"
	"strings"

	"github.com/pipe-cd/pipe/pkg/datastore"
)

func buildGetQuery(table string) string {
	return fmt.Sprintf("SELECT data FROM %s WHERE id = UUID_TO_BIN(?,true)", table)
}

func buildUpdateQuery(table string) string {
	return fmt.Sprintf("UPDATE %s SET data = ? WHERE id = UUID_TO_BIN(?,true)", table)
}

func buildPutQuery(table string) string {
	return fmt.Sprintf("INSERT INTO %s (id, data) VALUE (UUID_TO_BIN(?,true), ?) ON DUPLICATE KEY UPDATE data = ?", table)
}

func buildCreateQuery(table string) string {
	return fmt.Sprintf("INSERT INTO %s (id, data) VALUE (UUID_TO_BIN(?,true), ?)", table)
}

func buildFindQuery(table string, ops datastore.ListOptions) string {
	rawQuery := fmt.Sprintf(
		"SELECT data FROM %s %s %s %s",
		table,
		buildWhereClause(ops.Filters),
		buildOrderByClause(ops.Orders),
		buildPaginationClause(ops.Page, ops.PageSize),
	)
	return strings.Join(strings.Fields(rawQuery), " ")
}

func buildWhereClause(filters []datastore.ListFilter) string {
	if len(filters) == 0 {
		return ""
	}

	conds := make([]string, len(filters))
	for i, filter := range filters {
		conds[i] = fmt.Sprintf("%s %s ?", filter.Field, filter.Operator)
	}
	return fmt.Sprintf("WHERE %s", strings.Join(conds[:], " AND "))
}

func buildOrderByClause(orders []datastore.Order) string {
	if len(orders) == 0 {
		return ""
	}

	conds := make([]string, len(orders))
	for i, ord := range orders {
		conds[i] = fmt.Sprintf("%s %s", ord.Field, toMySQLDirection(ord.Direction))
	}
	return fmt.Sprintf("ORDER BY %s", strings.Join(conds[:], ", "))
}

func buildPaginationClause(page, pageSize int) string {
	var clause string
	if pageSize > 0 {
		clause = fmt.Sprintf("LIMIT %d ", pageSize)
		if page > 0 {
			clause = fmt.Sprintf("%sOFFSET %d", clause, pageSize*page)
		}
	}
	return clause
}

func toMySQLDirection(d datastore.OrderDirection) string {
	switch d {
	case datastore.Asc:
		return "ASC"
	case datastore.Desc:
		return "DESC"
	default:
		return ""
	}
}
