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
