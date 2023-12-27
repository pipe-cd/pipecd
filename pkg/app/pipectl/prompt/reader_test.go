// Copyright 2023 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prompt

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadString(t *testing.T) {
	testcases := []struct {
		name        string
		input       string
		expected    string
		expectedErr bool
	}{
		{
			name:        "valid input",
			input:       "foo\n",
			expected:    "foo",
			expectedErr: false,
		},
		{
			name:        "empty input",
			input:       "\n",
			expected:    "",
			expectedErr: false,
		},
		{
			name:        "two words",
			input:       "foo bar\n",
			expected:    "",
			expectedErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := bytes.NewBufferString(tc.input)
			r := stdinReader{in: in}
			actual, e := r.ReadString("anyPrompt")
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedErr, e != nil)
		})
	}
}

func TestReadStrings(t *testing.T) {
	testcases := []struct {
		name        string
		input       string
		expected    []string
		expectedErr bool
	}{
		{
			name:        "valid input",
			input:       "foo\n",
			expected:    []string{"foo"},
			expectedErr: false,
		},
		{
			name:        "empty input",
			input:       "\n",
			expected:    []string{},
			expectedErr: false,
		},
		{
			name:        "two words",
			input:       "foo bar\n",
			expected:    []string{"foo", "bar"},
			expectedErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := bytes.NewBufferString(tc.input)
			r := stdinReader{in: in}
			actual, e := r.ReadStrings("anyPrompt")
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedErr, e != nil)
		})
	}
}

func TestReadInt(t *testing.T) {
	testcases := []struct {
		name        string
		input       string
		expected    int
		expectedErr bool
	}{
		{
			name:        "valid input",
			input:       "123\n",
			expected:    123,
			expectedErr: false,
		},
		{
			name:        "empty input",
			input:       "\n",
			expected:    0,
			expectedErr: false,
		},
		{
			name:        "string input",
			input:       "abc\n",
			expected:    0,
			expectedErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := bytes.NewBufferString(tc.input)
			r := stdinReader{in: in}
			actual, e := r.ReadInt("anyPrompt")
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedErr, e != nil)
		})
	}
}

func TestReadStringRequired(t *testing.T) {
	testcases := []struct {
		name        string
		input       string
		expected    string
		expectedErr bool
	}{
		{
			name:        "valid input",
			input:       "foo\n",
			expected:    "foo",
			expectedErr: false,
		},
		{
			name:        "empty input",
			input:       "\n",
			expected:    "",
			expectedErr: true,
		},
		{
			name:        "two words",
			input:       "foo bar\n",
			expected:    "",
			expectedErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := bytes.NewBufferString(tc.input)
			r := stdinReader{in: in}
			actual, e := r.ReadStringRequired("anyPrompt")
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedErr, e != nil)
		})
	}
}
