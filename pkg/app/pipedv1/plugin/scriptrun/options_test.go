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

package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecode(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name     string
		data     json.RawMessage
		expected scriptRunStageOptions
		wantErr  bool
	}{
		{
			name: "valid config",
			data: json.RawMessage(`{
				"run":"echo 1"
			}`),
			expected: scriptRunStageOptions{
				Run: "echo 1",
			},
			wantErr: false,
		},
		{
			name:     "invalid config",
			data:     json.RawMessage(`invalid`),
			expected: scriptRunStageOptions{},
			wantErr:  true,
		},
		{
			name:     "empty config",
			data:     json.RawMessage(`{}`),
			expected: scriptRunStageOptions{},
			wantErr:  true,
		},
		{
			name: "multiline onRollback",
			data: json.RawMessage(`{
			"run": "echo main",
			"onRollback": "echo rollback1\necho rollback2"
		}`),
			expected: scriptRunStageOptions{
				Run:        "echo main",
				OnRollback: "echo rollback1\necho rollback2",
			},
			wantErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := decode(tc.data)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.expected, got)
		})
	}
}
