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

package terraform

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlanHasChangeRegex(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "older than v1.5.0",
			input:    "Plan: 1 to add, 2 to change, 3 to destroy.",
			expected: []string{"Plan: 1 to add, 2 to change, 3 to destroy.", "", "1", "2", "3"},
		},
		{
			name:     "later than v1.5.0",
			input:    "Plan: 0 to import, 1 to add, 2 to change, 3 to destroy.",
			expected: []string{"Plan: 0 to import, 1 to add, 2 to change, 3 to destroy.", "0", "1", "2", "3"},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, planHasChangeRegex.FindStringSubmatch(tc.input))
		})
	}
}

func TestParsePlanResult(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		input       string
		expected    PlanResult
		expectedErr bool
	}{
		{
			name:        "older than v1.5.0",
			input:       `Plan: 1 to add, 2 to change, 3 to destroy.`,
			expected:    PlanResult{Adds: 1, Changes: 2, Destroys: 3, PlanOutput: "Plan: 1 to add, 2 to change, 3 to destroy."},
			expectedErr: false,
		},
		{
			name:        "later than v1.5.0",
			input:       `Plan: 1 to import, 1 to add, 2 to change, 3 to destroy.`,
			expected:    PlanResult{Imports: 1, Adds: 1, Changes: 2, Destroys: 3, PlanOutput: "Plan: 1 to import, 1 to add, 2 to change, 3 to destroy."},
			expectedErr: false,
		},
		{
			name:        "Invalid number of changes",
			input:       `Plan: a to add, 2 to change, 3 to destroy.`,
			expectedErr: true,
		},
		{
			name:        "Invalid plan result output",
			input:       `Plan: 1 to add, 2 to change.`,
			expectedErr: true,
		},
		{
			name:        "No changes",
			input:       `No changes. Infrastructure is up-to-date.`,
			expectedErr: false,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, err := parsePlanResult(tc.input, false)
			assert.Equal(t, tc.expectedErr, err != nil)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRender(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		expected    string
		planResult  *PlanResult
		expectedErr bool
	}{
		{
			name: "success",
			planResult: &PlanResult{
				Imports:  1,
				Adds:     2,
				Changes:  3,
				Destroys: 4,
				PlanOutput: `
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:
  + resource "test-add" "test" {
      + id    = (known after apply)
    }
  - resource "test-del" "test" {
      + id    = "foo"
    }

Plan: 1 to import, 2 to add, 3 to change, 4 to destroy.
`,
			},
			expected: `    resource "test-add" "test" {
+       id    = (known after apply)
    }
    resource "test-del" "test" {
+       id    = "foo"
    }
Plan: 1 to import, 2 to add, 3 to change, 4 to destroy.
`,
			expectedErr: false,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual, err := tc.planResult.Render()
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedErr, err != nil)
		})
	}
}
