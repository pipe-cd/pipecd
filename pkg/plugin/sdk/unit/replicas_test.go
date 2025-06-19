// Copyright 2025 The PipeCD Authors.
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

package unit

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReplicasMarshal(t *testing.T) {
	type wrapper struct {
		Replicas Replicas
	}

	testcases := []struct {
		name     string
		input    wrapper
		expected string
	}{
		{
			name: "normal number",
			input: wrapper{
				Replicas{
					Number:       1,
					IsPercentage: false,
				},
			},
			expected: "{\"Replicas\":\"1\"}",
		},
		{
			name: "percentage number",
			input: wrapper{
				Replicas{
					Number:       1,
					IsPercentage: true,
				},
			},
			expected: "{\"Replicas\":\"1%\"}",
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

func TestReplicasUnmarshal(t *testing.T) {
	type wrapper struct {
		Replicas Replicas
	}

	testcases := []struct {
		name        string
		input       string
		expected    *wrapper
		expectedErr error
	}{
		{
			name:  "normal number",
			input: "{\"Replicas\": 1}",
			expected: &wrapper{
				Replicas{
					Number:       1,
					IsPercentage: false,
				},
			},
			expectedErr: nil,
		},
		{
			name:  "normal number by string",
			input: "{\"Replicas\":\"1\"}",
			expected: &wrapper{
				Replicas{
					Number:       1,
					IsPercentage: false,
				},
			},
			expectedErr: nil,
		},
		{
			name:  "percentage number",
			input: "{\"Replicas\":\"1%\"}",
			expected: &wrapper{
				Replicas{
					Number:       1,
					IsPercentage: true,
				},
			},
			expectedErr: nil,
		},
		{
			name:        "wrong string format",
			input:       "{\"Replicas\":\"1a%\"}",
			expected:    nil,
			expectedErr: fmt.Errorf("invalid replicas: strconv.Atoi: parsing \"1a\": invalid syntax"),
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := &wrapper{}
			err := json.Unmarshal([]byte(tc.input), got)
			assert.Equal(t, tc.expectedErr, err)
			if tc.expected != nil {
				assert.Equal(t, tc.expected, got)
			}
		})
	}
}

func TestReplicasString(t *testing.T) {
	testcases := []struct {
		name     string
		input    Replicas
		expected string
	}{
		{
			name: "normal number",
			input: Replicas{
				Number:       5,
				IsPercentage: false,
			},
			expected: "5",
		},
		{
			name: "percentage number",
			input: Replicas{
				Number:       50,
				IsPercentage: true,
			},
			expected: "50%",
		},
		{
			name: "zero replicas",
			input: Replicas{
				Number:       0,
				IsPercentage: false,
			},
			expected: "0",
		},
		{
			name: "zero percentage",
			input: Replicas{
				Number:       0,
				IsPercentage: true,
			},
			expected: "0%",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.String()
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestReplicasCalculate(t *testing.T) {
	testcases := []struct {
		name         string
		input        Replicas
		total        int
		defaultValue int
		expected     int
	}{
		{
			name: "zero number returns default",
			input: Replicas{
				Number:       0,
				IsPercentage: false,
			},
			total:        10,
			defaultValue: 3,
			expected:     3,
		},
		{
			name: "zero percentage returns default",
			input: Replicas{
				Number:       0,
				IsPercentage: true,
			},
			total:        10,
			defaultValue: 5,
			expected:     5,
		},
		{
			name: "normal number",
			input: Replicas{
				Number:       7,
				IsPercentage: false,
			},
			total:        10,
			defaultValue: 3,
			expected:     7,
		},
		{
			name: "50 percent of 10",
			input: Replicas{
				Number:       50,
				IsPercentage: true,
			},
			total:        10,
			defaultValue: 3,
			expected:     5,
		},
		{
			name: "50 percent of 3 (rounds up)",
			input: Replicas{
				Number:       50,
				IsPercentage: true,
			},
			total:        3,
			defaultValue: 1,
			expected:     2,
		},
		{
			name: "33 percent of 10 (rounds up)",
			input: Replicas{
				Number:       33,
				IsPercentage: true,
			},
			total:        10,
			defaultValue: 1,
			expected:     4,
		},
		{
			name: "25 percent of 10",
			input: Replicas{
				Number:       25,
				IsPercentage: true,
			},
			total:        10,
			defaultValue: 1,
			expected:     3,
		},
		{
			name: "10 percent of 3 (rounds up to 1)",
			input: Replicas{
				Number:       10,
				IsPercentage: true,
			},
			total:        3,
			defaultValue: 0,
			expected:     1,
		},
		{
			name: "100 percent",
			input: Replicas{
				Number:       100,
				IsPercentage: true,
			},
			total:        7,
			defaultValue: 1,
			expected:     7,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.Calculate(tc.total, tc.defaultValue)
			assert.Equal(t, tc.expected, got)
		})
	}
}
