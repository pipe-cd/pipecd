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

package filedb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestNormalizeFieldName(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "single camel",
			in:   "Id",
			out:  "id",
		},
		{
			// There will be no all rune upper cases like this
			// as generated code by protobuf, but we add this test
			// to mention derectly what is the expected output.
			name: "full of upper cases",
			in:   "API",
			out:  "aPI",
		},
		{
			name: "mix with full of upper cases word",
			in:   "ApiKey",
			out:  "apiKey",
		},
		{
			name: "formal camel",
			in:   "StaticAdminDisabled",
			out:  "staticAdminDisabled",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			out := normalizeFieldName(tc.in)
			assert.Equal(t, tc.out, out)
		})
	}
}

func TestCompare(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name      string
		val       interface{}
		operand   interface{}
		operator  datastore.Operator
		expect    bool
		expectErr bool
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
		{
			name:     "greater than int",
			val:      3,
			operand:  1,
			operator: datastore.OperatorGreaterThan,
			expect:   true,
		},
		{
			name:     "greater than or equal int",
			val:      3,
			operand:  3,
			operator: datastore.OperatorGreaterThanOrEqual,
			expect:   true,
		},
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
		{
			name:      "error on query for numeric only operator with wrong value",
			val:       "string_1",
			operand:   "string_0",
			operator:  datastore.OperatorGreaterThan,
			expectErr: true,
		},
		{
			name:      "error on query in operator with not operand of type slide/array",
			val:       1,
			operand:   1,
			operator:  datastore.OperatorIn,
			expectErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := compare(tc.val, tc.operand, tc.operator)
			require.Equal(t, tc.expectErr, err != nil)

			if err != nil {
				assert.Equal(t, tc.expect, res)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name    string
		entity  interface{}
		filters []datastore.ListFilter
		expect  bool
	}{
		{
			name:   "filter single condition - passed",
			entity: &model.Application{Id: "app_1"},
			filters: []datastore.ListFilter{
				{
					Field:    "Id",
					Operator: datastore.OperatorEqual,
					Value:    "app_1",
				},
			},
			expect: true,
		},
		{
			name:   "filter single condition - not passed",
			entity: &model.Application{Id: "app_1"},
			filters: []datastore.ListFilter{
				{
					Field:    "Id",
					Operator: datastore.OperatorEqual,
					Value:    "project_2",
				},
			},
			expect: false,
		},
		{
			name:   "filter multiple conditions - passed",
			entity: &model.Application{Id: "app_1", ProjectId: "project_1"},
			filters: []datastore.ListFilter{
				{
					Field:    "Id",
					Operator: datastore.OperatorEqual,
					Value:    "app_1",
				},
				{
					Field:    "ProjectId",
					Operator: datastore.OperatorEqual,
					Value:    "project_1",
				},
			},
			expect: true,
		},
		{
			name:   "filter multiple conditions with int zero value - passed",
			entity: &model.Application{Id: "app_1", Kind: model.ApplicationKind_KUBERNETES},
			filters: []datastore.ListFilter{
				{
					Field:    "Id",
					Operator: datastore.OperatorEqual,
					Value:    "app_1",
				},
				{
					Field:    "Kind",
					Operator: datastore.OperatorEqual,
					Value:    model.ApplicationKind_KUBERNETES,
				},
			},
			expect: true,
		},
		{
			name:   "filter multiple conditions with boolean zero value - passed",
			entity: &model.Application{Id: "app_1", Disabled: false},
			filters: []datastore.ListFilter{
				{
					Field:    "Id",
					Operator: datastore.OperatorEqual,
					Value:    "app_1",
				},
				{
					Field:    "Disabled",
					Operator: datastore.OperatorEqual,
					Value:    false,
				},
			},
			expect: true,
		},
		{
			name:   "filter multiple conditions - not passed",
			entity: &model.Application{Id: "app_1", ProjectId: "project_1"},
			filters: []datastore.ListFilter{
				{
					Field:    "Id",
					Operator: datastore.OperatorEqual,
					Value:    "app_1",
				},
				{
					Field:    "ProjectId",
					Operator: datastore.OperatorEqual,
					Value:    "project_2",
				},
			},
			expect: false,
		},
		{
			name:   "filter multiple conditions wrong type - not passed",
			entity: &model.Application{Id: "app_1", Disabled: false},
			filters: []datastore.ListFilter{
				{
					Field:    "Id",
					Operator: datastore.OperatorEqual,
					Value:    "app_1",
				},
				{
					Field:    "Disabled",
					Operator: datastore.OperatorEqual,
					Value:    0,
				},
			},
			expect: false,
		},
		{
			name:   "filter multiple conditions with numberic operator - passed",
			entity: &model.Application{Id: "app_1", UpdatedAt: 1649219699},
			filters: []datastore.ListFilter{
				{
					Field:    "Id",
					Operator: datastore.OperatorEqual,
					Value:    "app_1",
				},
				{
					Field:    "UpdatedAt",
					Operator: datastore.OperatorGreaterThanOrEqual,
					Value:    1646648937,
				},
			},
			expect: true,
		},
		{
			name:   "filter with IN operator - passed",
			entity: &model.Deployment{Status: model.DeploymentStatus_DEPLOYMENT_PENDING},
			filters: []datastore.ListFilter{
				{
					Field:    "Status",
					Operator: datastore.OperatorIn,
					Value:    model.GetNotCompletedDeploymentStatuses(),
				},
			},
			expect: true,
		},
		{
			name:   "filter with IN operator - not passed",
			entity: &model.Deployment{Status: model.DeploymentStatus_DEPLOYMENT_CANCELLED},
			filters: []datastore.ListFilter{
				{
					Field:    "Status",
					Operator: datastore.OperatorIn,
					Value:    model.GetNotCompletedDeploymentStatuses(),
				},
			},
			expect: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			passed, err := filter(nil, tc.entity, tc.filters)
			require.Nil(t, err)
			assert.Equal(t, tc.expect, passed)
		})
	}
}
