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
	return fmt.Sprintf("SELECT Data FROM %s WHERE Id = UUID_TO_BIN(?,true)", table)
}

func buildUpdateQuery(table string) string {
	return fmt.Sprintf("UPDATE %s SET Data = ? WHERE Id = UUID_TO_BIN(?,true)", table)
}

func buildPutQuery(table string) string {
	return fmt.Sprintf("INSERT INTO %s (Id, Data) VALUE (UUID_TO_BIN(?,true), ?) ON DUPLICATE KEY UPDATE Data = ?", table)
}

func buildCreateQuery(table string) string {
	return fmt.Sprintf("INSERT INTO %s (Id, Data) VALUE (UUID_TO_BIN(?,true), ?)", table)
}

func buildFindQuery(table string, ops datastore.ListOptions) string {
	rawQuery := fmt.Sprintf(
		"SELECT Data FROM %s %s %s %s",
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

	conds := make([]string, 0, len(filters))
	for _, filter := range filters {
		switch filter.Operator {
		case "==":
			conds = append(conds, fmt.Sprintf("%s = ?", filter.Field))
		case "in":
			conds = append(conds, fmt.Sprintf("%s IN ?", filter.Field))
		case "not-in":
			conds = append(conds, fmt.Sprintf("%s NOT IN ?", filter.Field))
		case "!=", ">", ">=", "<", "<=":
			conds = append(conds, fmt.Sprintf("%s %s ?", filter.Field, filter.Operator))
		default:
			// Skip if unsupported operator is passed.
			continue
		}
	}

	if len(conds) == 0 {
		return ""
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
