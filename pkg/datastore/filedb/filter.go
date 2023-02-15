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

package filedb

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/pipe-cd/pipecd/pkg/datastore"
)

type filterable interface {
	Match(e interface{}, filters []datastore.ListFilter) (bool, error)
}

func filter(col datastore.Collection, e interface{}, filters []datastore.ListFilter) (bool, error) {
	// Always pass, if there is no filter.
	if len(filters) == 0 {
		return true, nil
	}

	// If the collection implement filterable interface, use it.
	fcol, ok := col.(filterable)
	if ok {
		return fcol.Match(e, filters)
	}

	pe, ok := e.(proto.Message)
	if !ok {
		return false, datastore.ErrUnsupported
	}

	// remarshal entity as map[string]interface{} struct.
	m := protojson.MarshalOptions{
		EmitUnpopulated: true,
		UseEnumNumbers:  true,
	}
	raw, err := m.Marshal(pe)
	if err != nil {
		return false, err
	}
	var omap map[string]interface{}
	if err := json.Unmarshal(raw, &omap); err != nil {
		return false, err
	}

	for _, filter := range filters {
		field := normalizeFieldName(filter.Field)
		if strings.Contains(field, ".") {
			// TODO: Handle nested field name such as SyncState.Status.
			return false, datastore.ErrUnsupported
		}

		val, ok := omap[field]
		// If the object does not contain given field name in filter, return false immidiately.
		if !ok {
			return false, nil
		}

		operand, err := normalizeFieldValue(filter.Value)
		if err != nil {
			return false, err
		}

		cmp, err := compare(val, operand, filter.Operator)
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
	var (
		err                error
		valNum, operandNum float64
		valCasted          = true
		operandCasted      = true
	)
	switch v := val.(type) {
	case float32, float64:
		valNum = reflect.ValueOf(v).Float()
	case int, int8, int16, int32, int64:
		valNum = float64(reflect.ValueOf(v).Int())
	case uint, uint8, uint16, uint32, uint64:
		valNum = float64(reflect.ValueOf(v).Uint())
	case string:
		if !op.IsNumericOperator() {
			valCasted = false
			break
		}
		valNum, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return false, err
		}
	default:
		if op.IsNumericOperator() {
			return false, fmt.Errorf("value of type unsupported: %v", reflect.TypeOf(v))
		}
		valCasted = false
	}
	switch o := operand.(type) {
	case float32, float64:
		operandNum = reflect.ValueOf(o).Float()
	case int, int8, int16, int32, int64:
		operandNum = float64(reflect.ValueOf(o).Int())
	case uint, uint8, uint16, uint32, uint64:
		operandNum = float64(reflect.ValueOf(o).Uint())
	case string:
		if !op.IsNumericOperator() {
			operandCasted = false
			break
		}
		operandNum, err = strconv.ParseFloat(o, 64)
		if err != nil {
			return false, err
		}
	default:
		if op.IsNumericOperator() {
			return false, fmt.Errorf("operand of type unsupported: %v", reflect.TypeOf(o))
		}
		operandCasted = false
	}

	switch op {
	case datastore.OperatorEqual:
		if valCasted && operandCasted {
			return valNum == operandNum, nil
		}
		return val == operand, nil
	case datastore.OperatorNotEqual:
		if valCasted && operandCasted {
			return valNum != operandNum, nil
		}
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
		return nil, fmt.Errorf("value type %v is not a slide or array", rv.Kind())
	}

	vs := make([]interface{}, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		vs[i] = rv.Index(i).Interface()
	}

	return vs, nil
}

func normalizeFieldName(key string) string {
	if len(key) == 1 {
		return strings.ToLower(key)
	}
	return strings.ToLower(string(key[0])) + key[1:]
}

// normalizeFieldValue converts value of any type to the primitive type
// Note: Find a better way to handle this instead of marshal/unmarshal.
func normalizeFieldValue(val interface{}) (interface{}, error) {
	var needConvert = false
	switch val.(type) {
	case int, int8, int16, int32, int64:
	case uint, uint8, uint16, uint32, uint64:
	case float32, float64:
	case string:
	case bool:
	default:
		needConvert = true
	}

	if !needConvert {
		return val, nil
	}

	raw, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}

	var out interface{}
	if err = json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return out, nil
}
