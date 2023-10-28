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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/pipe-cd/pipecd/pkg/datastore"
)

var operatorMap = map[datastore.Operator]string{
	datastore.OperatorEqual:              "=",
	datastore.OperatorNotEqual:           "!=",
	datastore.OperatorIn:                 "IN",
	datastore.OperatorNotIn:              "NOT IN",
	datastore.OperatorGreaterThan:        ">",
	datastore.OperatorGreaterThanOrEqual: ">=",
	datastore.OperatorLessThan:           "<",
	datastore.OperatorLessThanOrEqual:    "<=",
	datastore.OperatorContains:           "MEMBER OF",
}

func buildGetQuery(table string) string {
	return fmt.Sprintf("SELECT Data FROM %s WHERE Id = UUID_TO_BIN(?,true)", table)
}

func buildUpdateQuery(table string) string {
	return fmt.Sprintf("UPDATE %s SET Data = ? WHERE Id = UUID_TO_BIN(?,true)", table)
}

func buildCreateQuery(table string) string {
	return fmt.Sprintf("INSERT INTO %s (Id, Data) VALUE (UUID_TO_BIN(?,true), ?)", table)
}

func buildFindQuery(table string, ops datastore.ListOptions) (string, error) {
	filters := refineFiltersField(ops.Filters)

	whereClause, err := buildWhereClause(filters)
	if err != nil {
		return "", err
	}
	orderByClause, err := buildOrderByClause(refineOrdersField(ops.Orders))
	if err != nil {
		return "", err
	}
	rawQuery := fmt.Sprintf(
		"SELECT Data FROM %s %s %s %s %s",
		table,
		whereClause,
		buildPaginationCondition(ops),
		orderByClause,
		buildLimitClause(ops.Limit),
	)
	return strings.Join(strings.Fields(rawQuery), " "), nil
}

func buildWhereClause(filters []datastore.ListFilter) (string, error) {
	if len(filters) == 0 {
		return "", nil
	}

	conds := make([]string, len(filters))
	for i, filter := range filters {
		op, ok := operatorMap[filter.Operator]
		if !ok {
			return "", fmt.Errorf("unsupported operator given: %v", filter.Operator)
		}
		switch filter.Operator {
		case datastore.OperatorIn, datastore.OperatorNotIn:
			// Make string of (?,...) which contains the number of `?` equal to the element number of filter.Value
			valLength := reflect.ValueOf(filter.Value).Len()
			conds[i] = fmt.Sprintf("%s %s (?%s)", filter.Field, op, strings.Repeat(",?", valLength-1))
		case datastore.OperatorContains:
			conds[i] = fmt.Sprintf("? %s (%s)", op, filter.Field)
		default:
			conds[i] = fmt.Sprintf("%s %s ?", filter.Field, op)
		}
	}
	return fmt.Sprintf("WHERE %s", strings.Join(conds, " AND ")), nil
}

func buildPaginationCondition(opts datastore.ListOptions) string {
	// Skip on no cursor.
	if len(opts.Cursor) == 0 {
		return ""
	}

	// Build outer set condition. The outer set condition should be
	// in format:
	//   X < Vx AND Y < Vy AND ...
	// with x, y, etc is not Id field.
	outerSetConds := make([]string, len(opts.Orders)-1)
	for i, o := range opts.Orders {
		if o.Field == "Id" {
			continue
		}
		outerSetConds[i] = fmt.Sprintf("%s %s ?", o.Field, makeCompareOperatorForOuterSet(o.Direction))
	}

	// Build sub set condition. The sub set condition should be
	// in format:
	//   X = Vx AND Y = Vy AND ... AND Id <= last_iterated_id
	// with last_iterated_id from the given cursor.
	subSetConds := make([]string, len(opts.Orders))
	for i, o := range opts.Orders {
		if o.Field == "Id" {
			subSetConds[i] = fmt.Sprintf("%s %s UUID_TO_BIN(?,true)", o.Field, makeCompareOperatorForSubSet(o.Direction))
		} else {
			subSetConds[i] = fmt.Sprintf("%s = ?", o.Field)
		}
	}

	// If there is no filter, mean pagination condition should be treated as the only where condition.
	if len(opts.Filters) == 0 {
		return fmt.Sprintf("WHERE %s AND NOT (%s)", strings.Join(outerSetConds, " AND "), strings.Join(subSetConds, " AND "))
	}
	return fmt.Sprintf("AND %s AND NOT (%s)", strings.Join(outerSetConds, " AND "), strings.Join(subSetConds, " AND "))
}

func makeCompareOperatorForOuterSet(direction datastore.OrderDirection) string {
	if direction == datastore.Asc {
		return ">="
	}
	return "<="
}

func makeCompareOperatorForSubSet(direction datastore.OrderDirection) string {
	if direction == datastore.Asc {
		return "<="
	}
	return ">="
}

func buildOrderByClause(orders []datastore.Order) (string, error) {
	if len(orders) == 0 {
		return "", nil
	}

	conds := make([]string, len(orders))
	hasIDFieldInOrdering := false
	for i, ord := range orders {
		if ord.Field == "Id" {
			hasIDFieldInOrdering = true
		}
		conds[i] = fmt.Sprintf("%s %s", ord.Field, toMySQLDirection(ord.Direction))
	}

	if !hasIDFieldInOrdering {
		return "", fmt.Errorf("id field is required as ordering field")
	}

	return fmt.Sprintf("ORDER BY %s", strings.Join(conds, ", ")), nil
}

func buildLimitClause(limit int) string {
	var clause string
	if limit > 0 {
		clause = fmt.Sprintf("LIMIT %d ", limit)
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

func refineOrdersField(orders []datastore.Order) []datastore.Order {
	out := make([]datastore.Order, len(orders))
	for i, order := range orders {
		switch order.Field {
		case "SyncState.Status":
			order.Field = "SyncState_Status"
		default:
			break
		}
		out[i] = order
	}
	return out
}

func refineFiltersField(filters []datastore.ListFilter) []datastore.ListFilter {
	out := make([]datastore.ListFilter, len(filters))
	for i, filter := range filters {
		switch filter.Field {
		case "SyncState.Status":
			filter.Field = "SyncState_Status"
		default:
			break
		}
		out[i] = filter
	}
	return out
}

// refineFiltersValue destructs all slide/array type values and makes an array of all element values.
func refineFiltersValue(filters []datastore.ListFilter) []interface{} {
	var filtersVals []interface{}
	for _, filter := range filters {
		fv := reflect.ValueOf(filter.Value)
		switch fv.Kind() {
		case reflect.Slice, reflect.Array:
			for j := 0; j < fv.Len(); j++ {
				filtersVals = append(filtersVals, fv.Index(j).Interface())
			}
		default:
			filtersVals = append(filtersVals, filter.Value)
		}
	}
	return filtersVals
}

// makePaginationCursorValues builds array of element values used on pagination condition check.
func makePaginationCursorValues(opts datastore.ListOptions) ([]interface{}, error) {
	// Skip pagination on cursor is empty.
	if len(opts.Cursor) == 0 {
		return nil, nil
	}

	// Decode last object of previous page stored as opts.Cursor to string.
	data, err := base64.StdEncoding.DecodeString(opts.Cursor)
	if err != nil {
		return nil, err
	}
	// Encode cursor data string to map[string]interface{} format for further process.
	obj := make(map[string]interface{})
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	// The cursorVals contains values used for pagination condition.
	// For each field except Id, it should be duplicated as for using in outer set and subset.
	// The Id field value should be one, and it's the last value in this list.
	cursorVals := make([]interface{}, 0, 2*len(opts.Orders)-1)
	for _, o := range opts.Orders {
		// Skip the Id field value to add it at last.
		if o.Field == "Id" {
			continue
		}
		val, ok := obj[o.Field]
		if !ok {
			return nil, fmt.Errorf("cursor does not contain values that match to ordering field %s", o.Field)
		}
		cursorVals = append(cursorVals, val)
	}
	// Duplicate all values in added order.
	cursorVals = append(cursorVals, cursorVals...)

	// Add Id value at last.
	id, ok := obj["Id"]
	if !ok {
		return nil, fmt.Errorf("cursor does not contain required value Id")
	}
	cursorVals = append(cursorVals, id)

	return cursorVals, nil
}
