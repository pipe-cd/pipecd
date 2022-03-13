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

	"github.com/pipe-cd/pipecd/pkg/datastore"
)

type filterable interface {
	Match(e interface{}, filters []datastore.ListFilter) (bool, error)
}

func filter(col datastore.Collection, e interface{}, filters []datastore.ListFilter) (bool, error) {
	fcol, ok := col.(filterable)
	if ok {
		return fcol.Match(e, filters)
	}

	// remarshal entity as map[string]interface{} struct.
	raw, _ := json.Marshal(e)
	var omap map[string]interface{}
	if err := json.Unmarshal(raw, &omap); err != nil {
		return false, err
	}

	for _, filter := range filters {
		if strings.Contains(filter.Field, ".") {
			// TODO: Handle nested field name such as SyncState.Status.
			return false, datastore.ErrUnsupported
		}

		val, ok := omap[filter.Field]
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
	switch op {
	case datastore.OperatorEqual:
		return val == operand, nil
	case datastore.OperatorNotEqual:
		return val != operand, nil
	case datastore.OperatorGreaterThan:
		return val.(int64) > operand.(int64), nil
	case datastore.OperatorGreaterThanOrEqual:
		return val.(int64) >= operand.(int64), nil
	case datastore.OperatorLessThan:
		return val.(int64) < operand.(int64), nil
	case datastore.OperatorLessThanOrEqual:
		return val.(int64) <= operand.(int64), nil
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
