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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/datastore"
)

func TestCompare(t *testing.T) {
	testcases := []struct {
		name     string
		val      interface{}
		operand  interface{}
		operator datastore.Operator
		expect   bool
	}{
		{
			name:     "equal number int",
			val:      5,
			operand:  5,
			operator: datastore.OperatorEqual,
			expect:   true,
		},
		{
			name:     "equal string",
			val:      "text",
			operand:  "text",
			operator: datastore.OperatorEqual,
			expect:   true,
		},
		{
			name:     "not equal int",
			val:      3,
			operand:  2,
			operator: datastore.OperatorNotEqual,
			expect:   true,
		},
		{
			name:     "not equal string",
			val:      "text_val",
			operand:  "text_operand",
			operator: datastore.OperatorNotEqual,
			expect:   true,
		},
		// {
		// 	name:     "greater than int",
		// 	val:      3,
		// 	operand:  1,
		// 	operator: datastore.OperatorGreaterThan,
		// 	expect:   true,
		// },
		// {
		// 	name:     "greater than string",
		// 	val:      "text_2",
		// 	operand:  "text_1",
		// 	operator: datastore.OperatorGreaterThan,
		// 	expect:   true,
		// },
		{
			name:     "in int",
			val:      1,
			operand:  []int{1, 2, 3},
			operator: datastore.OperatorIn,
			expect:   true,
		},
		{
			name:     "in int false",
			val:      4,
			operand:  []int{1, 2, 3},
			operator: datastore.OperatorIn,
			expect:   false,
		},
		{
			name:     "not in int",
			val:      4,
			operand:  []int{1, 2, 3},
			operator: datastore.OperatorNotIn,
			expect:   true,
		},
		{
			name:     "not in int false",
			val:      1,
			operand:  []int{1, 2, 3},
			operator: datastore.OperatorNotIn,
			expect:   false,
		},
		{
			name:     "contains int",
			val:      []int{1, 2, 3},
			operand:  1,
			operator: datastore.OperatorContains,
			expect:   true,
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res, err := compare(tc.val, tc.operand, tc.operator)
			require.Nil(t, err)
			assert.Equal(t, tc.expect, res)
		})
	}
}
