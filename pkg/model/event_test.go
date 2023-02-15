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

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeEventKey(t *testing.T) {
	testcases := []struct {
		testname string
		name1    string
		labels1  map[string]string
		name2    string
		labels2  map[string]string
		wantSame bool
	}{
		{
			testname: "no name and labels given",
			wantSame: true,
		},
		{
			testname: "no labels given",
			name1:    "name",
			name2:    "name",
			wantSame: true,
		},
		{
			testname: "no name given",
			labels1: map[string]string{
				"key1": "value1",
			},
			labels2: map[string]string{
				"key1": "value1",
			},
			wantSame: true,
		},
		{
			testname: "the exact same labels given",
			name1:    "name",
			labels1: map[string]string{
				"key1": "value",
				"key2": "value",
				"key3": "value",
			},
			name2: "name",
			labels2: map[string]string{
				"key2": "value",
				"key3": "value",
				"key1": "value",
			},
			wantSame: true,
		},
		{
			testname: "the sub match labels given",
			name1:    "name",
			labels1: map[string]string{
				"key1": "value",
				"key2": "value",
				"key3": "value",
			},
			name2: "name",
			labels2: map[string]string{
				"key1": "value",
			},
			wantSame: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.testname, func(t *testing.T) {
			// Check if we can get the exact same string.
			got1 := MakeEventKey(tc.name1, tc.labels1)
			got2 := MakeEventKey(tc.name2, tc.labels2)
			assert.NotEmpty(t, got1)
			assert.NotEmpty(t, got2)
			assert.Equal(t, tc.wantSame, got1 == got2)
		})
	}
}
