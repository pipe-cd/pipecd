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

package diff

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderNodeValue(t *testing.T) {
	var (
		mapOfPrimative = map[string]string{
			"one": "1",
			"two": "2",
		}
		mapOfMap = map[string]interface{}{
			"one": map[string]string{
				"one": "1-1",
				"two": "1-2",
			},
			"two": map[string]string{
				"one": "2-1",
				"two": "2-2",
			},
		}
		mapOfSlice = map[string]interface{}{
			"one": []string{"one-1", "one-2"},
			"two": []string{"two-1", "two-2"},
		}
	)

	testcases := []struct {
		name     string
		value    reflect.Value
		expected string
	}{
		{
			name:     "int value",
			value:    reflect.ValueOf(1),
			expected: "1",
		},
		{
			name:     "float value",
			value:    reflect.ValueOf(1.25),
			expected: "1.25",
		},
		{
			name:     "string value",
			value:    reflect.ValueOf("hello"),
			expected: "hello",
		},
		{
			name: "slice of primitive elements",
			value: func() reflect.Value {
				v := []int{1, 2, 3}
				return reflect.ValueOf(v)
			}(),
			expected: `- 1
- 2
- 3`,
		},
		{
			name: "slice of interface",
			value: func() reflect.Value {
				v := []interface{}{
					map[string]int{
						"1-one": 1,
						"2-two": 2,
					},
					map[string]int{
						"3-three": 3,
						"4-four":  4,
					},
				}
				return reflect.ValueOf(v)
			}(),
			expected: `- 1-one: 1
  2-two: 2
- 3-three: 3
  4-four: 4`,
		},
		{
			name: "simple map",
			value: reflect.ValueOf(map[string]string{
				"one": "one-value",
				"two": "two-value",
			}),
			expected: `one: one-value
two: two-value`,
		},
		{
			name: "nested map",
			value: func() reflect.Value {
				v := map[string]interface{}{
					"1-number":           1,
					"2-string":           "hello",
					"3-map-of-primitive": mapOfPrimative,
					"4-map-of-map":       mapOfMap,
					"5-map-of-slice":     mapOfSlice,
					"6-slice":            []string{"a", "b"},
					"7-string":           "hi",
				}
				return reflect.ValueOf(v)
			}(),
			expected: `1-number: 1
2-string: hello
3-map-of-primitive:
  one: 1
  two: 2
4-map-of-map:
  one:
    one: 1-1
    two: 1-2
  two:
    one: 2-1
    two: 2-2
5-map-of-slice:
  one:
    - one-1
    - one-2
  two:
    - two-1
    - two-2
6-slice:
  - a
  - b
7-string: hi`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := renderNodeValue(tc.value, "")
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestRenderNodeValueComplex(t *testing.T) {
	// Complex node. Note that the keys in the yaml file must be in order.
	objs, err := loadUnstructureds("testdata/complex-node.yaml")
	require.NoError(t, err)
	require.Equal(t, 1, len(objs))

	root := reflect.ValueOf(objs[0].Object)
	got, _ := renderNodeValue(root, "")

	data, err := os.ReadFile("testdata/complex-node.yaml")
	require.NoError(t, err)
	assert.Equal(t, string(data), got)
}

func TestRenderPrimitiveValue(t *testing.T) {
	testcases := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{
			name:     "string",
			value:    "hello",
			expected: "hello",
		},
		{
			name:     "int",
			value:    1,
			expected: "1",
		},
		{
			name:     "float",
			value:    1.25,
			expected: "1.25",
		},
		{
			name: "map",
			value: map[string]int{
				"one": 1,
			},
			expected: "<map[string]int Value>",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(tc.value)
			got := RenderPrimitiveValue(v)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestRenderByCommand(t *testing.T) {
	t.Parallel()

	diffCommand := "diff"

	type SubStruct struct {
		FieldSub string
	}
	type TestStruct struct {
		StringV     string
		StringP     *string
		IntV        int
		BoolV       bool
		TaggedV     int `json:"taggedv_tag"`
		StringSlice []string
		Map         map[string]string
		SubStruct   SubStruct
	}

	testcases := []struct {
		name         string
		command      string
		old          TestStruct
		new          TestStruct
		expectedDiff string
		expectErr    bool
	}{
		{
			name:         "invalid command",
			command:      "invalid-command",
			old:          TestStruct{},
			new:          TestStruct{},
			expectedDiff: "",
			expectErr:    true,
		},
		{
			name:    "no diff",
			command: diffCommand,
			old: TestStruct{
				StringV:     "value1",
				StringP:     nil,
				IntV:        1,
				BoolV:       true,
				TaggedV:     1,
				StringSlice: []string{"a", "b"},
				Map: map[string]string{
					"key1": "value1",
				},
				SubStruct: SubStruct{
					FieldSub: "valueSub",
				},
			},
			new: TestStruct{
				StringV:     "value1",
				StringP:     nil,
				IntV:        1,
				BoolV:       true,
				TaggedV:     1,
				StringSlice: []string{"a", "b"},
				Map: map[string]string{
					"key1": "value1",
				},
				SubStruct: SubStruct{
					FieldSub: "valueSub",
				},
			},
			expectedDiff: "",
			expectErr:    false,
		},
		{
			name:    "has diff",
			command: diffCommand,
			old: TestStruct{
				StringV:     "value1",
				StringP:     &[]string{"a"}[0],
				IntV:        1,
				BoolV:       true,
				TaggedV:     1,
				StringSlice: []string{"a", "b"},
				Map: map[string]string{
					"key1": "value1",
				},
				SubStruct: SubStruct{
					FieldSub: "valueSub",
				},
			},
			new: TestStruct{
				StringV:     "value2",
				StringP:     &[]string{"b"}[0],
				IntV:        2,
				BoolV:       false,
				TaggedV:     10,
				StringSlice: []string{"a", "b", "c"},
				Map: map[string]string{
					"key1": "valueXXX",
				},
				SubStruct: SubStruct{
					FieldSub: "xxx",
				},
			},
			expectedDiff: `@@ -1,12 +1,13 @@
-BoolV: true
-IntV: 1
+BoolV: false
+IntV: 2
 Map:
-  key1: value1
-StringP: a
+  key1: valueXXX
+StringP: b
 StringSlice:
 - a
 - b
-StringV: value1
+- c
+StringV: value2
 SubStruct:
-  FieldSub: valueSub
-taggedv_tag: 1
+  FieldSub: xxx
+taggedv_tag: 10`,
			expectErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := RenderByCommand(tc.command, tc.old, tc.new)
			assert.Equal(t, tc.expectErr, err != nil)
			assert.Equal(t, tc.expectedDiff, string(actual))
		})
	}
}
