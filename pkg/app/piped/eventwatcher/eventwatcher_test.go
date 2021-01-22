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

package eventwatcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetValue(t *testing.T) {
	testcases := []struct {
		name    string
		value   interface{}
		want    string
		wantErr bool
	}{
		{
			name:    "string",
			value:   "value",
			want:    "value",
			wantErr: false,
		},
		{
			name:    "int",
			value:   1,
			want:    "1",
			wantErr: false,
		},
		{
			name:    "int64",
			value:   int64(1),
			want:    "1",
			wantErr: false,
		},
		{
			name:    "uint64",
			value:   uint64(1),
			want:    "1",
			wantErr: false,
		},
		{
			name:    "float64",
			value:   1.1,
			want:    "1.1",
			wantErr: false,
		},
		{
			name:    "bool",
			value:   true,
			want:    "true",
			wantErr: false,
		},
		{
			name:    "map",
			value:   make(map[string]interface{}),
			want:    "",
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := convertStr(tc.value)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.want, got)
		})
	}
}
