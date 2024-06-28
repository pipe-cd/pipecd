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

package cache

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type getterFunc func(key string) (interface{}, error)

func (f getterFunc) Get(key string) (interface{}, error) {
	return f(key)
}

func (f getterFunc) GetAll() (map[string]interface{}, error) {
	return nil, ErrUnimplemented
}

func TestMultiGetter(t *testing.T) {
	value := "ok"
	err := errors.New("err")
	var calls int

	successGetter := getterFunc(func(key string) (interface{}, error) {
		calls++
		return value, nil
	})
	failureGetter := getterFunc(func(key string) (interface{}, error) {
		calls++
		return nil, err
	})
	testcases := []struct {
		name   string
		getter Getter
		err    error
		value  interface{}
		calls  int
	}{
		{
			getter: MultiGetter(successGetter),
			value:  value,
			calls:  1,
		},
		{
			getter: MultiGetter(successGetter, successGetter),
			value:  value,
			calls:  1,
		},
		{
			getter: MultiGetter(successGetter, failureGetter),
			value:  value,
			calls:  1,
		},
		{
			getter: MultiGetter(failureGetter, successGetter),
			value:  value,
			calls:  2,
		},
		{
			getter: MultiGetter(failureGetter, MultiGetter(failureGetter, successGetter)),
			value:  value,
			calls:  3,
		},
		{
			getter: MultiGetter(failureGetter, failureGetter),
			err:    err,
			calls:  2,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			calls = 0
			value, err := tc.getter.Get("")
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.value, value)
			assert.Equal(t, tc.calls, calls)
		})
	}
}
