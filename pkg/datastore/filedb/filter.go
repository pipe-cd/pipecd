// Copyright 2022 The PipeCD Authors.
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

package filedb

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/pipe-cd/pipecd/pkg/datastore"
)

type filterable interface {
	Match(e interface{}, filters []datastore.ListFilter) (bool, error)
}

func filter(col datastore.Collection, e interface{}, filters []datastore.ListFilter) (bool, error) {
	// If the collection implement filterable interface, use it.
	fcol, ok := col.(filterable)
	if ok {
		return fcol.Match(e, filters)
	}

	// remarshal entity as map[string]interface{} struct.
	raw, err := json.Marshal(e)
	if err != nil {
		return false, err
	}
	var omap map[string]interface{}
	if err := json.Unmarshal(raw, &omap); err != nil {
		return false, err
	}

	for _, filter := range filters {
		field := convertCamelToSnake(filter.Field)
		if strings.Contains(field, ".") {
			// TODO: Handle nested field name such as SyncState.Status.
			return false, datastore.ErrUnsupported
		}

		val, ok := omap[field]
		// If the object does not contain given field name in filter, return false immidiately.
		if !ok {
			return false, nil
		}

		cmp, err := compare(val, filter.Value, filter.Operator)
		if err != nil {
			return false, err
		}

		if !cmp {
			return false, nil
		}
	}

	return true, nil
}

func compare(val, operand interface{}, op datastore.Operator) (bool, error) {
	var valNum, operandNum int64
	switch v := val.(type) {
	case int, int8, int16, int32, int64:
		valNum = reflect.ValueOf(v).Int()
	case uint, uint8, uint16, uint32:
		valNum = int64(reflect.ValueOf(v).Uint())
	default:
		if op.IsNumericOperator() {
			return false, fmt.Errorf("value of type unsupported")
		}
	}
	switch o := operand.(type) {
	case int, int8, int16, int32, int64:
		operandNum = reflect.ValueOf(o).Int()
	case uint, uint8, uint16, uint32:
		operandNum = int64(reflect.ValueOf(o).Uint())
	default:
		if op.IsNumericOperator() {
			return false, fmt.Errorf("operand of type unsupported")
		}
	}

	switch op {
	case datastore.OperatorEqual:
		return val == operand, nil
	case datastore.OperatorNotEqual:
		return val != operand, nil
	case datastore.OperatorGreaterThan:
		return valNum > operandNum, nil
	case datastore.OperatorGreaterThanOrEqual:
		return valNum >= operandNum, nil
	case datastore.OperatorLessThan:
		return valNum < operandNum, nil
	case datastore.OperatorLessThanOrEqual:
		return valNum <= operandNum, nil
	case datastore.OperatorIn:
		os, err := makeSliceOfInterfaces(operand)
		if err != nil {
			return false, fmt.Errorf("operand error: %w", err)
		}

		for _, o := range os {
			if o == val {
				return true, nil
			}
		}
		return false, nil
	case datastore.OperatorNotIn:
		os, err := makeSliceOfInterfaces(operand)
		if err != nil {
			return false, fmt.Errorf("operand error: %w", err)
		}

		for _, o := range os {
			if o == val {
				return false, nil
			}
		}
		return true, nil
	case datastore.OperatorContains:
		vs, err := makeSliceOfInterfaces(val)
		if err != nil {
			return false, fmt.Errorf("value error: %w", err)
		}

		for _, v := range vs {
			if v == operand {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, datastore.ErrUnsupported
	}
}

func makeSliceOfInterfaces(v interface{}) ([]interface{}, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil, fmt.Errorf("value is not a slide or array")
	}

	vs := make([]interface{}, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		vs[i] = rv.Index(i).Interface()
	}

	return vs, nil
}

func convertCamelToSnake(key string) string {
	runeToLower := func(r rune) string {
		return strings.ToLower(string(r))
	}

	var out string
	for i, v := range key {
		if i == 0 {
			out += runeToLower(v)
			continue
		}

		if i == len(key)-1 {
			out += runeToLower(v)
			break
		}

		if unicode.IsUpper(v) && unicode.IsLower(rune(key[i+1])) {
			out += fmt.Sprintf("_%s", runeToLower(v))
			continue
		}

		if unicode.IsUpper(v) {
			out += runeToLower(v)
			continue
		}

		out += string(v)
	}
	return out
}
