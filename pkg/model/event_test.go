// Copyright 2021 The PipeCD Authors.
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
		name     string
		labels   map[string]string
		want     string
	}{
		{
			testname: "no name and labels given",
			want:     "",
		},
		{
			testname: "no labels given",
			name:     "name1",
			want:     "name1",
		},
		{
			testname: "no name given",
			labels: map[string]string{
				"key1": "value1",
			},
			want: "/key1:value1",
		},
		{
			testname: "labels given",
			name:     "name1",
			labels: map[string]string{
				"key1": "value",
				"key2": "value",
				"key3": "value",
			},
			want: "name1/key1:value/key2:value/key3:value",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.testname, func(t *testing.T) {
			got := MakeEventKey(tc.name, tc.labels)
			assert.Equal(t, tc.want, got)
		})
	}
}
