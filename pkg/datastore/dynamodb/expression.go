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

package dynamodb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/pipe-cd/pipe/pkg/datastore"
)

func buildDynamoDBCondition(f datastore.ListFilter) (expression.ConditionBuilder, error) {
	switch f.Operator {
	case "==":
		return expression.Name(f.Field).Equal(expression.Value(f.Value)), nil
	case "!=":
		return expression.Name(f.Field).NotEqual(expression.Value(f.Value)), nil
	case ">":
		return expression.Name(f.Field).GreaterThan(expression.Value(f.Value)), nil
	case ">=":
		return expression.Name(f.Field).GreaterThanEqual(expression.Value(f.Value)), nil
	case "in":
		return expression.Name(f.Field).In(expression.Value(f.Value)), nil
	case "<":
		return expression.Name(f.Field).LessThan(expression.Value(f.Value)), nil
	case "<=":
		return expression.Name(f.Field).LessThanEqual(expression.Value(f.Value)), nil
	default:
		return expression.ConditionBuilder{}, fmt.Errorf("unacceptable expression for dynamodb: %s %s %v", f.Field, f.Operator, f.Value)
	}
}

func buildDynamoDBExpression(opts datastore.ListOptions) (expression.Expression, error) {
	var expr expression.Expression
	if ok := hasHashKeyField(opts); !ok {
		return expr, fmt.Errorf("missing required field: %s", hashKeyFieldName)
	}

	// TODO: Build KeyConditionExpression.

	// Build FilterExpression.
	ops := make([]expression.ConditionBuilder, len(opts.Filters))
	for i, f := range opts.Filters {
		op, err := buildDynamoDBCondition(f)
		if err != nil {
			return expr, err
		}
		ops[i] = op
	}
	// if len(ops) == 0 {
	// 	return expr, fmt.Errorf("missing expression for dynamodb")
	// }
	var filter expression.ConditionBuilder
	switch len(ops) {
	case 0:
		// if there is no filter, just pass nil filter is okay.
		break
	case 1:
		filter = ops[0]
	case 2:
		filter = expression.And(ops[0], ops[1])
	default:
		filter = expression.And(ops[0], ops[1], ops[2:]...)
	}
	return expression.NewBuilder().WithFilter(filter).Build()
}

func buildDynamoDBKeyExistedExpression(key string, value interface{}) (expression.Expression, error) {
	return expression.NewBuilder().WithCondition(expression.Name(key).Equal(expression.Value(value))).Build()
}

func buildDynamoDBKeyNotExistedExpression(key string, value interface{}) (expression.Expression, error) {
	return expression.NewBuilder().WithCondition(expression.Name(key).NotEqual(expression.Value(value))).Build()
}

func hasHashKeyField(opts datastore.ListOptions) bool {
	for _, f := range opts.Filters {
		if f.Field == hashKeyFieldName {
			return true
		}
	}
	return false
}
