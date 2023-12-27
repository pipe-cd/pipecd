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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockReadString(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		inputs        []string
		expected      string
		expectedArray []string
		expectedErr   bool
	}{
		{
			name:          "valid input",
			inputs:        []string{"abc"},
			expected:      "abc",
			expectedArray: []string{},
			expectedErr:   false,
		},
		{
			name:          "empty input",
			inputs:        []string{""},
			expected:      "",
			expectedArray: []string{},
			expectedErr:   false,
		},
		{
			name:          "too many arguments",
			inputs:        []string{"abc def"},
			expected:      "",
			expectedArray: []string{},
			expectedErr:   true,
		},
		{
			name:          "do not read next element",
			inputs:        []string{"abc", "def"},
			expected:      "abc",
			expectedArray: []string{"def"},
			expectedErr:   false,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r := mockReader{inputs: tc.inputs}
			actual, e := r.ReadString("anyPrompt")
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedArray, r.inputs)
			assert.Equal(t, tc.expectedErr, e != nil)
		})
	}

}

func TestMockReadStrings(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		inputs        []string
		expected      []string
		expectedArray []string
		expectedErr   bool
	}{
		{
			name:          "valid input",
			inputs:        []string{"abc def"},
			expected:      []string{"abc", "def"},
			expectedArray: []string{},
			expectedErr:   false,
		},
		{
			name:          "empty input",
			inputs:        []string{""},
			expected:      []string{},
			expectedArray: []string{},
			expectedErr:   false,
		},
		{
			name:          "do not read next element",
			inputs:        []string{"abc def", "ghi jkl"},
			expected:      []string{"abc", "def"},
			expectedArray: []string{"ghi jkl"},
			expectedErr:   false,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r := mockReader{inputs: tc.inputs}
			actual, e := r.ReadStrings("anyPrompt")
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedArray, r.inputs)
			assert.Equal(t, tc.expectedErr, e != nil)
		})
	}
}

func TestMockReadInt(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		inputs        []string
		expected      int
		expectedArray []string
		expectedErr   bool
	}{
		{
			name:          "valid input",
			inputs:        []string{"123"},
			expected:      123,
			expectedArray: []string{},
			expectedErr:   false,
		},
		{
			name:          "empty input",
			inputs:        []string{""},
			expected:      0,
			expectedArray: []string{},
			expectedErr:   false,
		},
		{
			name:          "string input",
			inputs:        []string{"abc"},
			expected:      0,
			expectedArray: []string{},
			expectedErr:   true,
		},
		{
			name:          "do not read next element",
			inputs:        []string{"1", "23"},
			expected:      1,
			expectedArray: []string{"23"},
			expectedErr:   false,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r := mockReader{inputs: tc.inputs}
			actual, e := r.ReadInt("anyPrompt")
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedArray, r.inputs)
			assert.Equal(t, tc.expectedErr, e != nil)
		})
	}
}

func TestMockReadStringRequired(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		inputs        []string
		expected      string
		expectedArray []string
		expectedErr   bool
	}{
		{
			name:          "valid input",
			inputs:        []string{"abc"},
			expected:      "abc",
			expectedArray: []string{},
			expectedErr:   false,
		},
		{
			name:          "empty input",
			inputs:        []string{""},
			expected:      "",
			expectedArray: []string{},
			expectedErr:   true,
		},
		{
			name:          "too many arguments",
			inputs:        []string{"abc def"},
			expected:      "",
			expectedArray: []string{},
			expectedErr:   true,
		},
		{
			name:          "do not read next element",
			inputs:        []string{"abc", "def"},
			expected:      "abc",
			expectedArray: []string{"def"},
			expectedErr:   false,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r := mockReader{inputs: tc.inputs}
			actual, e := r.ReadStringRequired("anyPrompt")
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedArray, r.inputs)
			assert.Equal(t, tc.expectedErr, e != nil)
		})
	}
}
