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

package config

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPercentageMarshal(t *testing.T) {
	type wrapper struct {
		Percentage Percentage
	}

	testcases := []struct {
		name     string
		input    wrapper
		expected string
	}{
		{
			name: "normal number",
			input: wrapper{
				Percentage{
					Number:    10,
					HasSuffix: false,
				},
			},
			expected: `{"Percentage":"10"}`,
		},
		{
			name: "percentage number",
			input: wrapper{
				Percentage{
					Number:    15,
					HasSuffix: true,
				},
			},
			expected: `{"Percentage":"15%"}`,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, string(got))
		})
	}
}

func TestPercentageUnmarshal(t *testing.T) {
	type wrapper struct {
		Percentage Percentage
	}

	testcases := []struct {
		name        string
		input       string
		expected    *wrapper
		expectedErr bool
	}{
		{
			name:  "normal number",
			input: `{"Percentage": 10}`,
			expected: &wrapper{
				Percentage{
					Number: 10,
				},
			},
		},
		{
			name:  "normal number by string",
			input: `{"Percentage": "10"}`,
			expected: &wrapper{
				Percentage{
					Number: 10,
				},
			},
		},
		{
			name:  "percentage number",
			input: `{"Percentage": "10%"}`,
			expected: &wrapper{
				Percentage{
					Number:    10,
					HasSuffix: true,
				},
			},
		},
		{
			name:        "wrong string format",
			input:       `{"Percentage": "1a%"}`,
			expected:    nil,
			expectedErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := &wrapper{}
			err := json.Unmarshal([]byte(tc.input), got)
			assert.Equal(t, tc.expectedErr, err != nil)
			if tc.expected != nil {
				assert.Equal(t, tc.expected, got)
			}
		})
	}
}
