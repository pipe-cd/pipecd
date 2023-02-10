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

package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveMapFields(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		config   map[string]interface{}
		live     map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name:     "Empty map",
			config:   make(map[string]interface{}, 0),
			live:     make(map[string]interface{}, 0),
			expected: make(map[string]interface{}, 0),
		},
		{
			name: "Not nested 1",
			config: map[string]interface{}{
				"key a": "value a",
			},
			live: map[string]interface{}{
				"key a": "value a",
				"key b": "value b",
			},
			expected: map[string]interface{}{
				"key a": "value a",
			},
		},
		{
			name: "Not nested 2",
			config: map[string]interface{}{
				"key a": "value a",
				"key b": "value b",
			},
			live: map[string]interface{}{
				"key a": "value a",
			},
			expected: map[string]interface{}{
				"key a": "value a",
			},
		},
		{
			name: "Nested live deleted",
			config: map[string]interface{}{
				"key a": "value a",
			},
			live: map[string]interface{}{
				"key a": "value a",
				"key b": map[string]interface{}{
					"nested key a": "nested value a",
				},
			},
			expected: map[string]interface{}{
				"key a": "value a",
			},
		},
		{
			name: "Nested same",
			config: map[string]interface{}{
				"key a": "value a",
				"key b": map[string]interface{}{
					"nested key a": "nested value a",
				},
			},
			live: map[string]interface{}{
				"key a": "value a",
				"key b": map[string]interface{}{
					"nested key a": "nested value a",
				},
			},
			expected: map[string]interface{}{
				"key a": "value a",
				"key b": map[string]interface{}{
					"nested key a": "nested value a",
				},
			},
		},
		{
			name: "Nested nested live deleted",
			config: map[string]interface{}{
				"key a": "value a",
				"key b": map[string]interface{}{
					"nested key a": "nested value a",
				},
			},
			live: map[string]interface{}{
				"key a": "value a",
				"key b": map[string]interface{}{
					"nested key a": "nested value a",
					"nested key b": "nested value b",
				},
			},
			expected: map[string]interface{}{
				"key a": "value a",
				"key b": map[string]interface{}{
					"nested key a": "nested value a",
				},
			},
		},
		{
			name: "Nested array",
			config: map[string]interface{}{
				"key a": "value a",
				"key b": []interface{}{
					"a", "b", 3,
				},
			},
			live: map[string]interface{}{
				"key a": "value a",
				"key b": []interface{}{
					"a", "b", 3,
				},
			},
			expected: map[string]interface{}{
				"key a": "value a",
				"key b": []interface{}{
					"a", "b", 3,
				},
			},
		},
		{
			name: "Nested array 2",
			config: map[string]interface{}{
				"key a": "value a",
				"key b": []interface{}{
					"a", "b", 3, 4,
				},
			},
			live: map[string]interface{}{
				"key a": "value a",
				"key b": []interface{}{
					"a", "b", 3,
				},
			},
			expected: map[string]interface{}{
				"key a": "value a",
				"key b": []interface{}{
					"a", "b", 3,
				},
			},
		},
		{
			name: "Nested array remain",
			config: map[string]interface{}{
				"key a": "value a",
				"key b": []interface{}{
					"a", "b",
				},
			},
			live: map[string]interface{}{
				"key a": "value a",
				"key b": []interface{}{
					"a", "b", map[string]interface{}{
						"aa": "aa",
					},
				},
			},
			expected: map[string]interface{}{
				"key a": "value a",
				"key b": []interface{}{
					"a", "b", map[string]interface{}{
						"aa": "aa",
					},
				},
			},
		},
		{
			name: "Nested array same",
			config: map[string]interface{}{
				"key a": "value a",
				"key b": []interface{}{
					"a", "b", 3,
				},
			},
			live: map[string]interface{}{
				"key a": "value a",
				"key b": []interface{}{
					"b", "a", 3,
				},
			},
			expected: map[string]interface{}{
				"key a": "value a",
				"key b": []interface{}{
					"b", "a", 3,
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			removed := removeMapFields(tc.config, tc.live)
			assert.Equal(t, tc.expected, removed)
		})
	}
}
