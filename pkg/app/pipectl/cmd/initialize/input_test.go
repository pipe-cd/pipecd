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

package initialize

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPromptString(t *testing.T) {
	testcases := []struct {
		name        string
		inputs      string
		expected    string
		expectedErr bool
	}{
		{
			name:        "valid input",
			inputs:      "foo\n",
			expected:    "foo",
			expectedErr: false,
		},
		{
			name:        "empty input",
			inputs:      "\n",
			expected:    "",
			expectedErr: false,
		},
		{
			name:        "two words",
			inputs:      "foo bar \n",
			expected:    "",
			expectedErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// in := bytes.NewBufferString(tc.inputs)
			in := strings.NewReader(tc.inputs)
			actual, e := promptString("any-prompt", in)
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedErr, e != nil)
		})
	}
}

func TestPromptStringSequence1(t *testing.T) {
	txt := "xxx \n aaa \n ccc"
	in := strings.NewReader(txt)
	s1, e1 := promptString("", in)
	s2, e2 := promptString("", in)
	s3, e3 := promptString("", in)
	assert.Equal(t, "xxx", s1)
	assert.Equal(t, "aaa", s2)
	assert.Equal(t, "ccc", s3)

	fmt.Printf("e1: %v\n", e1)
	fmt.Printf("e2: %v\n", e2)
	fmt.Printf("e3: %v\n", e3)

	assert.Nil(t, e1)
	assert.Nil(t, e2)
	assert.Nil(t, e3)
}

func TestPromptStringSequence(t *testing.T) {
	txt := "xxx\naaa  bbb \nccc"
	in := bytes.NewBufferString(txt)
	s1, e1 := promptString("", in)
	s2, e2 := promptStrings("", in)
	s3, e3 := promptString("", in)
	assert.Equal(t, "xxx", s1)
	assert.Equal(t, []string{"aaa", "bbb"}, s2)
	assert.Equal(t, "ccc", s3)

	assert.Nil(t, e1)
	assert.Nil(t, e2)
	assert.Nil(t, e3)
}

func TestPromptStrings(t *testing.T) {
	testcases := []struct {
		name     string
		inputs   string
		expected []string
	}{
		{
			name:     "two words",
			inputs:   "foo bar\n",
			expected: []string{"foo", "bar"},
		},
		{
			name:     "empty input",
			inputs:   "\n",
			expected: []string{},
		},
		{
			name:     "longer blank space",
			inputs:   "   foo     bar   \n",
			expected: []string{"foo", "bar"},
		},
		{
			name:     "don't read next line",
			inputs:   "aaa bbb\nccc",
			expected: []string{"aaa", "bbb"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := bytes.NewBufferString(tc.inputs)
			actual, e := promptStrings("any-prompt", in)
			assert.Equal(t, tc.expected, actual)
			assert.Nil(t, e)
		})
	}
}

func TestPromptInt(t *testing.T) {
	testcases := []struct {
		name        string
		inputs      string
		expected    int
		expectedErr bool
	}{
		{
			name:        "valid input",
			inputs:      "123\n",
			expected:    123,
			expectedErr: false,
		},
		{
			name:        "empty input",
			inputs:      "\n",
			expected:    0,
			expectedErr: false,
		},
		{
			name:        "string input",
			inputs:      "abc\n",
			expected:    0,
			expectedErr: true,
		},
		{
			name:        "number with blank",
			inputs:      " 123 \n",
			expected:    123,
			expectedErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := bytes.NewBufferString(tc.inputs)
			actual, e := promptInt("any-prompt", in)
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedErr, e != nil)
		})
	}
}
