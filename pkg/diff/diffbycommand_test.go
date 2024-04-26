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
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectedDiff, string(actual))
		})
	}
}
