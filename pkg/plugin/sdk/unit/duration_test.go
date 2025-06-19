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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDuration(t *testing.T) {
	testcases := []struct {
		name     string
		input    Duration
		expected time.Duration
	}{
		{
			name:     "zero duration",
			input:    Duration(0),
			expected: 0,
		},
		{
			name:     "one second",
			input:    Duration(time.Second),
			expected: time.Second,
		},
		{
			name:     "one minute",
			input:    Duration(time.Minute),
			expected: time.Minute,
		},
		{
			name:     "complex duration",
			input:    Duration(2*time.Hour + 30*time.Minute + 15*time.Second),
			expected: 2*time.Hour + 30*time.Minute + 15*time.Second,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.Duration()
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestDurationMarshal(t *testing.T) {
	type wrapper struct {
		Duration Duration
	}

	testcases := []struct {
		name     string
		input    wrapper
		expected string
	}{
		{
			name: "zero duration",
			input: wrapper{
				Duration: Duration(0),
			},
			expected: `{"Duration":"0s"}`,
		},
		{
			name: "one second",
			input: wrapper{
				Duration: Duration(time.Second),
			},
			expected: `{"Duration":"1s"}`,
		},
		{
			name: "one minute",
			input: wrapper{
				Duration: Duration(time.Minute),
			},
			expected: `{"Duration":"1m0s"}`,
		},
		{
			name: "one hour",
			input: wrapper{
				Duration: Duration(time.Hour),
			},
			expected: `{"Duration":"1h0m0s"}`,
		},
		{
			name: "complex duration",
			input: wrapper{
				Duration: Duration(2*time.Hour + 30*time.Minute + 15*time.Second),
			},
			expected: `{"Duration":"2h30m15s"}`,
		},
		{
			name: "milliseconds",
			input: wrapper{
				Duration: Duration(500 * time.Millisecond),
			},
			expected: `{"Duration":"500ms"}`,
		},
		{
			name: "microseconds",
			input: wrapper{
				Duration: Duration(100 * time.Microsecond),
			},
			expected: `{"Duration":"100µs"}`,
		},
		{
			name: "nanoseconds",
			input: wrapper{
				Duration: Duration(50),
			},
			expected: `{"Duration":"50ns"}`,
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

func TestDurationUnmarshal(t *testing.T) {
	type wrapper struct {
		Duration Duration
	}

	testcases := []struct {
		name        string
		input       string
		expected    *wrapper
		expectedErr bool
	}{
		{
			name:  "duration as float64 nanoseconds",
			input: `{"Duration": 1000000000}`,
			expected: &wrapper{
				Duration: Duration(time.Second),
			},
			expectedErr: false,
		},
		{
			name:  "duration as string seconds",
			input: `{"Duration": "1s"}`,
			expected: &wrapper{
				Duration: Duration(time.Second),
			},
			expectedErr: false,
		},
		{
			name:  "duration as string minutes",
			input: `{"Duration": "5m"}`,
			expected: &wrapper{
				Duration: Duration(5 * time.Minute),
			},
			expectedErr: false,
		},
		{
			name:  "duration as string hours",
			input: `{"Duration": "2h"}`,
			expected: &wrapper{
				Duration: Duration(2 * time.Hour),
			},
			expectedErr: false,
		},
		{
			name:  "complex duration string",
			input: `{"Duration": "2h30m45s"}`,
			expected: &wrapper{
				Duration: Duration(2*time.Hour + 30*time.Minute + 45*time.Second),
			},
			expectedErr: false,
		},
		{
			name:  "duration with milliseconds",
			input: `{"Duration": "1.5s"}`,
			expected: &wrapper{
				Duration: Duration(1500 * time.Millisecond),
			},
			expectedErr: false,
		},
		{
			name:  "duration with microseconds",
			input: `{"Duration": "100us"}`,
			expected: &wrapper{
				Duration: Duration(100 * time.Microsecond),
			},
			expectedErr: false,
		},
		{
			name:  "zero duration as float64",
			input: `{"Duration": 0}`,
			expected: &wrapper{
				Duration: Duration(0),
			},
			expectedErr: false,
		},
		{
			name:  "zero duration as string",
			input: `{"Duration": "0s"}`,
			expected: &wrapper{
				Duration: Duration(0),
			},
			expectedErr: false,
		},
		{
			name:        "invalid duration string",
			input:       `{"Duration": "invalid"}`,
			expected:    nil,
			expectedErr: true,
		},
		{
			name:        "invalid type bool",
			input:       `{"Duration": true}`,
			expected:    nil,
			expectedErr: true,
		},
		{
			name:        "invalid type array",
			input:       `{"Duration": [1, 2, 3]}`,
			expected:    nil,
			expectedErr: true,
		},
		{
			name:        "invalid type object",
			input:       `{"Duration": {"value": 1}}`,
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
